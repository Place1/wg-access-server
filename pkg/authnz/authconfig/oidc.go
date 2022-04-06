package authconfig

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authruntime"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authutil"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/Knetic/govaluate.v2"
	"gopkg.in/yaml.v2"
)

type OIDCConfig struct {
	Name         string                    `yaml:"name"`
	Issuer       string                    `yaml:"issuer"`
	ClientID     string                    `yaml:"clientID"`
	ClientSecret string                    `yaml:"clientSecret"`
	Scopes       []string                  `yaml:"scopes"`
	RedirectURL  string                    `yaml:"redirectURL"`
	EmailDomains []string                  `yaml:"emailDomains"`
	ClaimMapping map[string]ruleExpression `yaml:"claimMapping"`
}

func (c *OIDCConfig) Provider() *authruntime.Provider {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	provider, err := oidc.NewProvider(ctx, c.Issuer)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create oidc provider"))
	}

	if c.Scopes == nil {
		c.Scopes = []string{"openid"}
	}

	oauthConfig := &oauth2.Config{
		RedirectURL:  c.RedirectURL,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       c.Scopes,
		Endpoint:     provider.Endpoint(),
	}

	redirectURL, err := url.Parse(c.RedirectURL)
	if err != nil {
		panic(errors.Wrapf(err, "redirect url is not valid: %s", c.RedirectURL))
	}

	return &authruntime.Provider{
		Type: "OIDC",
		Invoke: func(w http.ResponseWriter, r *http.Request, runtime *authruntime.ProviderRuntime) {
			c.loginHandler(runtime, oauthConfig)(w, r)
		},
		RegisterRoutes: func(router *mux.Router, runtime *authruntime.ProviderRuntime) error {
			router.HandleFunc(redirectURL.Path, c.callbackHandler(runtime, oauthConfig, provider))
			return nil
		},
	}
}

func (c *OIDCConfig) loginHandler(runtime *authruntime.ProviderRuntime, oauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthStateString := authutil.RandomString(32)
		err := runtime.SetSession(w, r, &authsession.AuthSession{
			Nonce: &oauthStateString,
		})
		if err != nil {
			http.Error(w, "no session", http.StatusUnauthorized)
			return
		}
		url := oauthConfig.AuthCodeURL(oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (c *OIDCConfig) callbackHandler(runtime *authruntime.ProviderRuntime, oauthConfig *oauth2.Config, provider *oidc.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := runtime.GetSession(r)
		if err != nil {
			http.Error(w, "no session", http.StatusBadRequest)
			return
		}

		state := r.FormValue("state")
		if s.Nonce == nil || *s.Nonce != state {
			http.Error(w, "bad nonce", http.StatusBadRequest)
			return
		}

		code := r.FormValue("code")

		token, err := oauthConfig.Exchange(r.Context(), code)
		if err != nil {
			panic(errors.Wrap(err, "Unable to exchange tokens"))
		}

		info, err := provider.UserInfo(r.Context(), oauthConfig.TokenSource(r.Context(), token))
		if err != nil {
			panic(errors.Wrap(err, "Unable to get user info"))
		}

		if msg, valid := verifyEmailDomain(c.EmailDomains, info.Email); !valid {
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		oidcProfileData := make(map[string]interface{})
		err = info.Claims(&oidcProfileData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims := &authsession.Claims{}
		for claimName, rule := range c.ClaimMapping {
			result, err := rule.Evaluate(oidcProfileData)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// If result is 'false' or an empty string then don't include the Claim
			if val, ok := result.(bool); ok && val {
				claims.Add(claimName, strconv.FormatBool(val))
			} else if val, ok := result.(string); ok && len(val) > 0 {
				claims.Add(claimName, val)
			}
		}

		identity := &authsession.Identity{
			Provider: c.Name,
			Subject:  info.Subject,
			Email:    info.Email,
			Claims:   *claims,
		}
		if name, ok := oidcProfileData["name"].(string); ok {
			identity.Name = name
		}

		err = runtime.SetSession(w, r, &authsession.AuthSession{
			Identity: identity,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		runtime.Done(w, r)
	}
}

func verifyEmailDomain(allowedDomains []string, email string) (string, bool) {
	if len(allowedDomains) == 0 {
		return "", true
	}

	parsed := strings.Split(email, "@")

	// check we have 2 parts i.e. <user>@<domain>
	if len(parsed) != 2 {
		return "missing or invalid email address", false
	}

	// match the domain against the list of allowed domains
	for _, domain := range allowedDomains {
		if domain == parsed[1] {
			return "", true
		}
	}

	return "email domain not authorized", false
}

type ruleExpression struct {
	*govaluate.EvaluableExpression
}

// MarshalYAML will encode a RuleExpression/govalidate into yaml string
func (r ruleExpression) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(r.String())
}

// UnmarshalYAML will decode a RuleExpression/govalidate into yaml string
func (r *ruleExpression) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ruleStr string
	if err := unmarshal(&ruleStr); err != nil {
		return err
	}
	parsedRule, err := govaluate.NewEvaluableExpression(ruleStr)
	if err != nil {
		return errors.Wrap(err, "Unable to process oidc rule")
	}
	ruleExpression := &ruleExpression{parsedRule}
	*r = *ruleExpression
	return nil
}
