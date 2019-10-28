package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/alexedwards/scs/v2"
	"github.com/coreos/go-oidc"
	"github.com/dexidp/dex/storage"
	"github.com/dexidp/dex/storage/memory"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/dexidp/dex/server"
)

type DexIntegration struct {
	router       *mux.Router
	externalAddr string
	port         int
}

func NewDexServer(session *scs.SessionManager, externalAddr string, port int, config Config) (*DexIntegration, error) {

	users := []storage.Password{}
	if config.StaticUsers != nil {
		for _, u := range config.StaticUsers {
			h, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
			if err != nil {
				return nil, errors.Wrap(err, "failed to hash password for static user")
			}
			users = append(users, storage.Password{
				Email: u.Email,
				Hash:  h,
			})
		}
	}

	connectors := []storage.Connector{}
	if config.Connectors != nil {
		for _, c := range config.Connectors {
			connectors = append(connectors, c.toDexConnector(externalAddr))
		}
	}

	s := storage.WithStaticClients(memory.New(logrus.New()), []storage.Client{
		storage.Client{
			ID:           "internal",
			Name:         "internal",
			RedirectURIs: []string{fmt.Sprintf("%s/auth/-/callback", externalAddr)},
			Secret:       "dummy-secret",
		},
	})
	if len(connectors) > 0 {
		s = storage.WithStaticConnectors(s, connectors)
	}
	if len(users) > 0 {
		s = storage.WithStaticPasswords(s, users, logrus.New())
	}

	serv, err := server.NewServer(context.TODO(), server.Config{
		Logger:             logrus.New(),
		Issuer:             fmt.Sprintf("%s/auth", externalAddr),
		PrometheusRegistry: prometheus.NewRegistry(),
		Storage:            s,
		Web: server.WebConfig{
			Dir:     "dex-web",
			LogoURL: "todo",
		},
		SkipApprovalScreen: true,
	})
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter().PathPrefix("/auth").Subrouter()

	dex := &DexIntegration{
		router:       router,
		externalAddr: externalAddr,
		port:         port,
	}

	router.HandleFunc("/login", dex.handleLogin).Methods("GET")
	router.HandleFunc("/-/callback", dex.handleCallback(session)).Methods("GET")
	router.PathPrefix("/").Handler(serv)

	return dex, nil
}

func (d *DexIntegration) Router() *mux.Router {
	return d.router
}

func (d *DexIntegration) oauthConfig() (*oauth2.Config, *oidc.IDTokenVerifier) {
	provider, err := oidc.NewProvider(context.TODO(), fmt.Sprintf("%s/auth", d.externalAddr))
	if err != nil {
		panic(err)
	}
	return &oauth2.Config{
		ClientID:     "internal",
		ClientSecret: "dummy-secret",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{"openid", "email", "profile"},
		RedirectURL:  fmt.Sprintf("%s/auth/-/callback", d.externalAddr),
	}, provider.Verifier(&oidc.Config{ClientID: "internal"})
}

func (d *DexIntegration) handleLogin(w http.ResponseWriter, r *http.Request) {
	logrus.Info("handling login")
	c, _ := d.oauthConfig()
	authCodeURL := c.AuthCodeURL("dummy-state")
	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}

func (d *DexIntegration) handleCallback(session *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("handling callback")
		state := r.URL.Query().Get("state")

		// Verify state.
		if state != "dummy-state" {
			panic("boom")
		}

		c, idTokenVerifier := d.oauthConfig()
		oauth2Token, err := c.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			panic(err)
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			panic("missing token")
		}

		// Parse and verify ID Token payload.
		idToken, err := idTokenVerifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			// handle error
			panic(err)
		}

		session.Put(r.Context(), "auth/subject", idToken.Subject)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
