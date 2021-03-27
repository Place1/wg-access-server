package authconfig

import (
	"fmt"
	"net/http"

	"github.com/place1/wg-access-server/pkg/authnz/authruntime"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
)

type ProxyAuthConfig struct {
	Headers struct {
		User  string `yaml:"user"`
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	}
}

func (c *ProxyAuthConfig) Provider() *authruntime.Provider {
	return &authruntime.Provider{
		Type: "Proxy",
		Invoke: func(w http.ResponseWriter, r *http.Request, runtime *authruntime.ProviderRuntime) {
			proxyAuthLogin(c, runtime)(w, r)
		},
	}
}

func proxyAuthLogin(c *ProxyAuthConfig, runtime *authruntime.ProviderRuntime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var username string
		var name string
		var email string
		if c.Headers.User != "" {
			username = r.Header.Get(c.Headers.User)
		} else {
			username = r.Header.Get("Remote-User")
		}

		if username != "" {

			if c.Headers.Name != "" {
				name = r.Header.Get(c.Headers.Name)
			} else {
				name = r.Header.Get("Remote-Name")
			}

			if name == "" {
				name = username
			}

			if c.Headers.Email != "" {
				email = r.Header.Get(c.Headers.Email)
			} else {
				email = r.Header.Get("Remote-Email")
			}
			runtime.SetSession(w, r, &authsession.AuthSession{
				Identity: &authsession.Identity{
					Provider: "proxy",
					Subject:  username,
					Name:     name,
					Email:    email,
				},
			})
			runtime.Done(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "unauthorized")
		return
	}
}
