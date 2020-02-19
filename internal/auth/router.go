package auth

import (
	"fmt"
	"net/http"

	"github.com/place1/wireguard-access-server/internal/auth/authconfig"
	"github.com/place1/wireguard-access-server/internal/auth/authruntime"
	"github.com/place1/wireguard-access-server/internal/auth/authsession"
	"github.com/place1/wireguard-access-server/internal/auth/authtemplates"
	"github.com/place1/wireguard-access-server/internal/auth/authutil"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type AuthMiddleware struct {
	config *authconfig.AuthConfig
}

func New(config *authconfig.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{config}
}

func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {

	runtime := authruntime.NewProviderRuntime(sessions.NewCookieStore([]byte(authutil.RandomString(32))))
	router := mux.NewRouter()

	for _, p := range m.config.Providers() {
		p.RegisterRoutes(router, runtime)
	}

	router.PathPrefix("/signin").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, authtemplates.RenderLoginPage(w, authtemplates.LoginPage{
			Config: m.config,
		}))
	}))

	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s, err := runtime.GetSession(r); err == nil {
			next.ServeHTTP(w, r.WithContext(authsession.SetIdentityCtx(r.Context(), s)))
		} else {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
	}))

	return router
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Index")
}
