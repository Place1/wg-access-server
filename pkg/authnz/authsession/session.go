package authsession

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AuthSession struct {
	Nonce    *string
	Identity *Identity
}

type authSessionKey string

var sessionKey authSessionKey = "auth-session"

func GetSession(store sessions.Store, r *http.Request) (*AuthSession, error) {
	session, _ := store.Get(r, string(sessionKey))
	if data, ok := session.Values[string(sessionKey)].([]byte); ok {
		s := &AuthSession{}
		if err := json.Unmarshal(data, s); err != nil {
			return nil, errors.Wrap(err, "failed to parse session")
		}
		return s, nil
	}
	return nil, errors.New("Session not authenticated")
}

func SetSession(store sessions.Store, r *http.Request, w http.ResponseWriter, s *AuthSession) error {
	data, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "failed to marshal session")
	}
	session, _ := store.Get(r, string(sessionKey))
	session.Values[string(sessionKey)] = data
	if err := session.Save(r, w); err != nil {
		return err
	}

	if s.Identity != nil {
		logrus.Infof("Creating web session with provider '%s' for user '%s' (remote address: %s)", s.Identity.Provider, s.Identity.Name, r.RemoteAddr)
	}

	return nil
}

func ClearSession(store sessions.Store, r *http.Request, w http.ResponseWriter) error {
	session, _ := store.Get(r, string(sessionKey))
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func SetIdentityCtx(parent context.Context, session *AuthSession) context.Context {
	return context.WithValue(parent, sessionKey, session)
}

func CurrentUser(ctx context.Context) (*Identity, error) {
	if session, ok := ctx.Value(sessionKey).(*AuthSession); ok {
		if session.Identity != nil {
			return session.Identity, nil
		}
	}
	return nil, errors.New("Unauthenticated")
}

func Authenticated(ctx context.Context) bool {
	_, err := CurrentUser(ctx)
	return err == nil
}
