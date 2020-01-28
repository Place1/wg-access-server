package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Provider struct {
	RegisterRoutes func(*mux.Router, *ProviderRuntime) error
}

type ProviderRuntime struct {
	store sessions.Store
}

func NewProviderRuntime(store sessions.Store) *ProviderRuntime {
	return &ProviderRuntime{store}
}

func (p *ProviderRuntime) SetSession(w http.ResponseWriter, r *http.Request, s *AuthSession) error {
	return setSession(p.store, r, w, s)
}

func (p *ProviderRuntime) GetSession(r *http.Request) (*AuthSession, error) {
	return getSession(p.store, r)
}

func (p *ProviderRuntime) Done(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
