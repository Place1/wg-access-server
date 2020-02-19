package auth

import (
	"fmt"
	"net/http"
	"strconv"

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
