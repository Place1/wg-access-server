package authconfig

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tg123/go-htpasswd"

	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authruntime"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"
)

const BasicAuthProvider = "basic"

type BasicAuthConfig struct {
	// Users is a list of htpasswd encoded username:password pairs
	// supports BCrypt, Sha, Ssha, Md5
	// example: "htpasswd -nB <username>"
	// copy the result into your user's array
	Users []string `yaml:"users"`
}

func (c *BasicAuthConfig) Provider() *authruntime.Provider {
	return &authruntime.Provider{
		Type: BasicAuthProvider,
		Invoke: func(w http.ResponseWriter, r *http.Request, runtime *authruntime.ProviderRuntime) {
			basicAuthLogin(c, runtime)(w, r)
		},
	}
}

func basicAuthLogin(c *BasicAuthConfig, runtime *authruntime.ProviderRuntime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok {
			if ok := checkCreds(c.Users, u, p); ok {
				err := runtime.SetSession(w, r, &authsession.AuthSession{
					Identity: &authsession.Identity{
						Provider: BasicAuthProvider,
						Subject:  u,
						Name:     u,
						Email:    "", // basic auth has no email
					},
				})
				if err == nil {
					runtime.Done(w, r)
					return
				}
			}
		}

		// If we're here something went wrong, return StatusUnauthorized
		w.Header().Set("WWW-Authenticate", `Basic realm="site"`)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "unauthorized")
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
