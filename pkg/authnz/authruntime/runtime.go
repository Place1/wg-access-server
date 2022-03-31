package authruntime

import (
	"net/http"

	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Provider struct {
	Type           string
	Invoke         func(http.ResponseWriter, *http.Request, *ProviderRuntime)
	RegisterRoutes func(*mux.Router, *ProviderRuntime) error
}

type ProviderRuntime struct {
	store sessions.Store
}

func NewProviderRuntime(store sessions.Store) *ProviderRuntime {
	return &ProviderRuntime{store}
}

func (p *ProviderRuntime) SetSession(w http.ResponseWriter, r *http.Request, s *authsession.AuthSession) error {
	return authsession.SetSession(p.store, r, w, s)
}

func (p *ProviderRuntime) GetSession(r *http.Request) (*authsession.AuthSession, error) {
	return authsession.GetSession(p.store, r)
}

func (p *ProviderRuntime) ClearSession(w http.ResponseWriter, r *http.Request) error {
	return authsession.ClearSession(p.store, r, w)
}

func (p *ProviderRuntime) Restart(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
}

func (p *ProviderRuntime) Done(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
