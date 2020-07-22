package authruntime

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/place1/wg-access-server/internal/traces"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
)

type Provider struct {
	// Name is a name for this provider
	// for display to users.
	Name string
	// Type is the name for the specific type
	// of provider and must be unique.
	Type           string
	Invoke         func(http.ResponseWriter, *http.Request, *ProviderRuntime)
	RegisterRoutes func(*mux.Router, *ProviderRuntime) error
	Branding       ProviderBranding
}

type ProviderBranding struct {
	Background string
	Color      string
	Icon       string
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

func (p *ProviderRuntime) ShowBanner(w http.ResponseWriter, r *http.Request, banner authsession.Banner) {
	data, err := json.Marshal(banner)
	if err != nil {
		traces.Logger(r.Context()).Error(errors.Wrap(err, "failed to serialize banner message"))
	}
	authsession.AddFlash(p.store, r, w, "banner", string(data))
	http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
}

func (p *ProviderRuntime) GetBanner(r *http.Request) (*authsession.Banner, bool) {
	if v, found := authsession.GetFlash(p.store, r, "banner"); found {
		banner := &authsession.Banner{}
		if err := json.Unmarshal([]byte(v), banner); err == nil {
			return banner, true
		}
	}
	return nil, false
}
