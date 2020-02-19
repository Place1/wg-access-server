package authconfig

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/place1/wireguard-access-server/internal/auth/authruntime"
	"github.com/place1/wireguard-access-server/internal/auth/authsession"
	"github.com/tg123/go-htpasswd"
)

type BasicAuthConfig struct {
	// Users is a list of htpasswd encoded username:password pairs
	// supports BCrypt, Sha, Ssha, Md5
	// example: "htpasswd -nB <username>"
	// copy the result into your user's array
	Users []string `yaml:"users"`
}

func (c *BasicAuthConfig) Provider() *authruntime.Provider {
	return &authruntime.Provider{
		Type: "Basic",
		Invoke: func(w http.ResponseWriter, r *http.Request, runtime *authruntime.ProviderRuntime) {
			basicAuthLogin(c, runtime)(w, r)
		},
	}
}

func basicAuthLogin(c *BasicAuthConfig, runtime *authruntime.ProviderRuntime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="site"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "unauthorized")
			return
		}

		if ok := checkCreds(c.Users, u, p); ok {
			runtime.SetSession(w, r, &authsession.AuthSession{
				Identity: &authsession.Identity{
					Subject: u,
				},
			})
		}

		runtime.Done(w, r)
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
