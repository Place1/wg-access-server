package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type AuthMiddleware struct {
	config *AuthConfig
}

func New(config *AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{config}
}

func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {

	runtime := NewProviderRuntime(sessions.NewCookieStore([]byte(randomString(32))))
	router := mux.NewRouter()

	for _, p := range m.config.Providers() {
		p.RegisterRoutes(router, runtime)
	}

	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s, err := runtime.GetSession(r); err == nil {
			next.ServeHTTP(w, r.WithContext(setIdentityCtx(r.Context(), s)))
		} else {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
	}))

	return router
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Index")
}
