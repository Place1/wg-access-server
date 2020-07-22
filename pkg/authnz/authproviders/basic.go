package authproviders

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/place1/wg-access-server/pkg/authnz/authruntime"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
	"github.com/tg123/go-htpasswd"
)

type BasicAuthConfig struct {
	Name string `yaml:"name"`
	// Users is a list of htpasswd encoded username:password pairs
	// supports BCrypt, Sha, Ssha, Md5
	// example: "htpasswd -nB <username>"
	// copy the result into your user's array
	Users []string `yaml:"users"`
}

func (c *BasicAuthConfig) Provider() *authruntime.Provider {
	return &authruntime.Provider{
		Name: c.Name,
		Type: "Basic",
		Invoke: func(w http.ResponseWriter, r *http.Request, runtime *authruntime.ProviderRuntime) {
			basicAuthLogin(c, runtime)(w, r)
		},
	}
}

func basicAuthLogin(c *BasicAuthConfig, runtime *authruntime.ProviderRuntime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// accept standard basic auth challenges
		u, p, isBasic := r.BasicAuth()

		if !isBasic {
			// we'll handle form submissions and direct
			// browser challenges
			u = r.FormValue("username")
			p = r.FormValue("password")
		}

		if ok := checkCreds(c.Users, u, p); ok {
			runtime.SetSession(w, r, &authsession.AuthSession{
				Identity: &authsession.Identity{
					Provider: "basic",
					Subject:  u,
					Name:     u,
					Email:    "", // basic auth has no email
				},
			})
			runtime.Done(w, r)
			return
		}

		if !isBasic {
			runtime.ShowBanner(w, r, authsession.Banner{
				Text:   "Invalid username or password",
				Intent: "danger",
			})
		} else {
			// challenge browser
			w.Header().Set("WWW-Authenticate", `Basic realm="site"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "unauthorized")
		}
	}
}

func checkCreds(users []string, username string, password string) bool {
	for _, user := range users {
		if u, p, ok := parsehtpassword(user); ok {
			if u == username {
				return checkhtpasswd(p, password)
			}
		}
	}
	return false
}

func parsehtpassword(user string) (string, string, bool) {
	segments := strings.SplitN(user, ":", 2)
	if len(segments) >= 1 {
		return segments[0], segments[1], true
	}
	return "", "", false
}

func checkhtpasswd(required string, given string) bool {
	if encoded, err := htpasswd.AcceptBcrypt(required); encoded != nil && err == nil {
		return encoded.MatchesPassword(given)
	}
	if encoded, err := htpasswd.AcceptSha(required); encoded != nil && err == nil {
		return encoded.MatchesPassword(given)
	}
	if encoded, err := htpasswd.AcceptSsha(required); encoded != nil && err == nil {
		return encoded.MatchesPassword(given)
	}
	if encoded, err := htpasswd.AcceptMd5(required); encoded != nil && err == nil {
		return encoded.MatchesPassword(given)
	}
	return false
}
