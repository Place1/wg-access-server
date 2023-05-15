package authconfig

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authruntime"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authtemplates"
)

const SimpleAuthProvider = "simple"

// SimpleAuthConfig is an alternative to BasicAuthConfig where the login happens through a login page and a POST request.
type SimpleAuthConfig struct {
	// Users is a list of htpasswd encoded username:password pairs
	// supports BCrypt, Sha, Ssha, Md5
	// example: "htpasswd -nB <username>"
	// copy the result into your user's array
	Users []string `yaml:"users"`
}

const postURL = "/signin/simpleauth"

func (c *SimpleAuthConfig) Provider() *authruntime.Provider {
	return &authruntime.Provider{
		Type: SimpleAuthProvider,
		Name: SimpleAuthProvider,
		// The flow is as follows: /signin page -> navigation to /signin/{index}
		// -> Invoke / simpleAuthLogin() renders login form -> POST to postURL / simpleAuthPostEndpoint()
		// -> redirect to /
		Invoke: func(w http.ResponseWriter, r *http.Request, runtime *authruntime.ProviderRuntime) {
			simpleAuthLogin()(w, r)
		},
		RegisterRoutes: func(router *mux.Router, runtime *authruntime.ProviderRuntime) error {
			router.HandleFunc(postURL, simpleAuthPostEndpoint(c, runtime))
			return nil
		},
	}
}

func simpleAuthLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// A login page with a username and a password field
		w.WriteHeader(http.StatusOK)
		err := authtemplates.RenderSimpleAuthPage(w, authtemplates.SimpleAuthPage{PostURL: postURL})
		if err != nil {
			logrus.Error(errors.Wrap(err, "failed to render simple auth login page"))
			return
		}
	}
}

func simpleAuthPostEndpoint(c *SimpleAuthConfig, runtime *authruntime.ProviderRuntime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}
		u := r.PostForm.Get("username")
		p := r.PostForm.Get("password")
		if u != "" && p != "" && checkCreds(c.Users, u, p) {
			err = runtime.SetSession(w, r, &authsession.AuthSession{
				Identity: &authsession.Identity{
					Provider: SimpleAuthProvider,
					Subject:  u,
					Name:     u,
					Email:    "", // simple auth has no email
				},
			})
			if err == nil {
				runtime.Done(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusForbidden)
		err = authtemplates.RenderSimpleAuthPage(w, authtemplates.SimpleAuthPage{
			PostURL:      postURL,
			ErrorMessage: "Invalid username or password",
		})
		if err != nil {
			logrus.Error(errors.Wrap(err, "failed to render simple auth login page"))
			return
		}
	}
}
