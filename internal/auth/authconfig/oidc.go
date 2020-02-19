package authconfig

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/place1/wireguard-access-server/internal/auth/authruntime"
	"github.com/place1/wireguard-access-server/internal/auth/authsession"
	"github.com/place1/wireguard-access-server/internal/auth/authutil"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	Name         string   `yaml:"name"`
	Issuer       string   `yaml:"issuer"`
	ClientID     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	Scopes       []string `yaml:"scopes"`
	RedirectURL  string   `yaml:"redirectURL"`
}

func (c *OIDCConfig) Provider() *authruntime.Provider {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	provider, err := oidc.NewProvider(ctx, c.Issuer)
	if err != nil {
		logrus.Fatal(err)
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
			loginHandler(runtime, oauthConfig)(w, r)
		},
		RegisterRoutes: func(router *mux.Router, runtime *authruntime.ProviderRuntime) error {
			router.HandleFunc(redirectURL.Path, callbackHandler(runtime, oauthConfig, provider))
			return nil
		},
	}
}

func loginHandler(runtime *authruntime.ProviderRuntime, oauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthStateString := authutil.RandomString(32)
		runtime.SetSession(w, r, &authsession.AuthSession{
			Nonce: &oauthStateString,
		})
		url := oauthConfig.AuthCodeURL(oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func callbackHandler(runtime *authruntime.ProviderRuntime, oauthConfig *oauth2.Config, provider *oidc.Provider) http.HandlerFunc {
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
		token, _ := oauthConfig.Exchange(r.Context(), code)
		info, err := provider.UserInfo(r.Context(), oauthConfig.TokenSource(r.Context(), token))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		runtime.SetSession(w, r, &authsession.AuthSession{
			Identity: &authsession.Identity{
				Subject: info.Subject,
			},
		})

		runtime.Done(w, r)
	}
}
