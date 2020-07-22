package authnz

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"

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
	router           *mux.Router
	runtime          *authruntime.ProviderRuntime
}

func New(config authconfig.AuthConfig, claimsMiddleware authsession.ClaimsMiddleware) *AuthMiddleware {
	router := mux.NewRouter()
	store := sessions.NewCookieStore([]byte(authutil.RandomString(32)))
	runtime := authruntime.NewProviderRuntime(store)
	providers := config.Providers()

	for _, p := range providers {
		if p.RegisterRoutes != nil {
			p.RegisterRoutes(router, runtime)
		}
	}

	router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		banner, _ := runtime.GetBanner(r)
		fmt.Fprint(w, authtemplates.RenderLoginPage(w, authtemplates.LoginPage{
			Providers: providers,
			// TODO: make configurable (branding)
			Title:  "Welcome to WireGuard Access Portal",
			Banner: banner,
		}))
	})

	router.HandleFunc("/signin/{index}", func(w http.ResponseWriter, r *http.Request) {
		index, err := strconv.Atoi(mux.Vars(r)["index"])
		if err != nil || index < 0 || len(providers) <= index {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "unknown provider")
			return
		}
		provider := providers[index]
		provider.Invoke(w, r, runtime)
	})

	router.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {
		runtime.ClearSession(w, r)
		runtime.Restart(w, r)
	})

	return &AuthMiddleware{
		config,
		claimsMiddleware,
		router,
		runtime,
	}
}

func NewMiddleware(config authconfig.AuthConfig, claimsMiddleware authsession.ClaimsMiddleware) mux.MiddlewareFunc {
	return New(config, claimsMiddleware).Middleware
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if the request is for an auth
		// related page i.e. /signin
		// to be handled by our own router
		if ok := m.router.Match(r, &mux.RouteMatch{}); ok {
			m.router.ServeHTTP(w, r)
			return
		}

		// otherwise we apply the standard middleware
		// functionality i.e. annotate the request context
		// with the request user (identity)
		if s, err := m.runtime.GetSession(r); err == nil {
			if m.claimsMiddleware != nil {
				if err := m.claimsMiddleware(s.Identity); err != nil {
					ctxlogrus.Extract(r.Context()).Error(errors.Wrap(err, "authz middleware failure"))
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(authsession.SetIdentityCtx(r.Context(), s)))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authsession.Authenticated(r.Context()) {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		}
	})
}
