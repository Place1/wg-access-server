package authnz

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/place1/wg-access-server/pkg/authnz/authconfig"
	"github.com/place1/wg-access-server/pkg/authnz/authruntime"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
	"github.com/place1/wg-access-server/pkg/authnz/authtemplates"
	"github.com/place1/wg-access-server/pkg/authnz/authutil"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type AuthMiddleware struct {
	config           authconfig.AuthConfig
	claimsMiddleware authsession.ClaimsMiddleware
}

func New(config authconfig.AuthConfig, claimsMiddleware authsession.ClaimsMiddleware) *AuthMiddleware {
	return &AuthMiddleware{config, claimsMiddleware}
}

func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {

	runtime := authruntime.NewProviderRuntime(sessions.NewCookieStore([]byte(authutil.RandomString(32))))
	router := mux.NewRouter()

	providers := m.config.Providers()

	for _, p := range providers {
		if p.RegisterRoutes != nil {
			p.RegisterRoutes(router, runtime)
		}
	}

	router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, authtemplates.RenderLoginPage(w, authtemplates.LoginPage{
			Providers: providers,
		}))
	})

	router.HandleFunc("/signin/{index}", func(w http.ResponseWriter, r *http.Request) {
		index, err := strconv.Atoi(mux.Vars(r)["index"])
		if err != nil || (index < 0 || index >= len(providers)) {
			fmt.Fprintf(w, "unknown provider")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		provider := providers[index]
		provider.Invoke(w, r, runtime)
	})

	router.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {
		runtime.ClearSession(w, r)
		runtime.Restart(w, r)
	})

	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s, err := runtime.GetSession(r); err == nil {
			if m.claimsMiddleware != nil {
				if err := m.claimsMiddleware(s.Identity); err != nil {
					logrus.Error(errors.Wrap(err, "authz middleware failure"))
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(authsession.SetIdentityCtx(r.Context(), s)))
		} else {
			next.ServeHTTP(w, r)
		}
	}))

	return router
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Index")
}
