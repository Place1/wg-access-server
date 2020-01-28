package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
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

func (c *OIDCConfig) Provider() *Provider {
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

	return &Provider{
		RegisterRoutes: func(router *mux.Router, runtime *ProviderRuntime) error {
			router.HandleFunc("/login", loginHandler(runtime, oauthConfig))
			router.HandleFunc("/callback", callbackHandler(runtime, oauthConfig, provider))
			return nil
		},
	}
}

func loginHandler(runtime *ProviderRuntime, oauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthStateString := randomString(32)
		runtime.SetSession(w, r, &AuthSession{
			Nonce: &oauthStateString,
		})
		url := oauthConfig.AuthCodeURL(oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func callbackHandler(runtime *ProviderRuntime, oauthConfig *oauth2.Config, provider *oidc.Provider) http.HandlerFunc {
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

		runtime.SetSession(w, r, &AuthSession{
			Identity: &Identity{
				Subject: info.Subject,
			},
		})

		runtime.Done(w, r)
	}
}

func randomString(size int) string {
	blk := make([]byte, size)
	_, err := rand.Read(blk)
	if err != nil {
		logrus.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(blk)
}
