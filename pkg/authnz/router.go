package authnz

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/freifunkMUC/wg-access-server/internal/traces"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authconfig"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authruntime"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authtemplates"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authutil"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

type AuthMiddleware struct {
	config           authconfig.AuthConfig
	claimsMiddleware authsession.ClaimsMiddleware
	router           *mux.Router
	runtime          *authruntime.ProviderRuntime
}

func New(config authconfig.AuthConfig, claimsMiddleware authsession.ClaimsMiddleware) (*AuthMiddleware, error) {
	router := mux.NewRouter()
	var storeSecret []byte
	if config.SessionStore == nil || config.SessionStore.Secret == "" {
		storeSecret = []byte(authutil.RandomString(32))
	} else {
		var err error
		storeSecret, err = hex.DecodeString(config.SessionStore.Secret)
		if err != nil {
			return nil, err
		}
		if len(storeSecret) != 32 {
			return nil, errors.New("session store secret must be 32 bytes long")
		}
	}
	store := sessions.NewCookieStore(storeSecret)
	runtime := authruntime.NewProviderRuntime(store)
	providers := config.Providers()

	for _, p := range providers {
		if p.RegisterRoutes != nil {
			err := p.RegisterRoutes(router, runtime)
			if err != nil {
				return nil, err
			}
		}
	}

	router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("signout") != "1" && !config.DesiresSigninPage() && len(providers) == 1 {
			// we only have one provider, so jump directly to that
			providers[0].Invoke(w, r, runtime)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, authtemplates.RenderLoginPage(w, authtemplates.LoginPage{
			Providers: providers,
		}))
	})

	router.HandleFunc("/signin/{index}", func(w http.ResponseWriter, r *http.Request) {
		index, err := strconv.Atoi(mux.Vars(r)["index"])
		if err != nil || index < 0 || len(providers) <= index {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "unknown provider")
			return
		}
		provider := providers[index]
		provider.Invoke(w, r, runtime)
	})

	router.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {
		_ = runtime.ClearSession(w, r)
		runtime.Restart(w, r)
	})

	return &AuthMiddleware{
		config,
		claimsMiddleware,
		router,
		runtime,
	}, nil
}

func NewMiddleware(config authconfig.AuthConfig, claimsMiddleware authsession.ClaimsMiddleware) (mux.MiddlewareFunc, error) {
	authMiddleware, err := New(config, claimsMiddleware)
	if err != nil {
		return nil, err
	}
	return authMiddleware.Middleware, nil
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
			if s.Identity == nil {
				// Can happen due to an aborted or failed login at the OIDC provider
				// Redirect the user to the signin page, so they can redo the login
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}
			if m.claimsMiddleware != nil {
				if err := m.claimsMiddleware(s.Identity); err != nil {
					traces.Logger(r.Context()).Error(errors.Wrap(err, "authnz middleware failure"))
					http.Redirect(w, r, "/signin", http.StatusSeeOther)
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(authsession.SetIdentityCtx(r.Context(), s)))
		} else {
			// GetSession() errors e.g. after the server restarted, because old session cookies are no longer trusted
			// The RequireAuthentication() middleware will be next in line and prompt the user to log in
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
