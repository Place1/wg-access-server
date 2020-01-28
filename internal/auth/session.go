package auth

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

type Identity struct {
	Subject string
}

type authSessionKey string

var sessionKey authSessionKey = "auth-session"

func getSession(store sessions.Store, r *http.Request) (*AuthSession, error) {
	session, _ := store.Get(r, string(sessionKey))
	if data, ok := session.Values[string(sessionKey)].([]byte); ok {
		s := &AuthSession{}
		err := json.Unmarshal(data, s)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse session")
		}
		return s, nil
	}
	return nil, errors.New("session not authenticated")
}

func setSession(store sessions.Store, r *http.Request, w http.ResponseWriter, s *AuthSession) error {
	data, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "failed to marshal session")
	}
	session, _ := store.Get(r, string(sessionKey))
	session.Values[string(sessionKey)] = data
	err = session.Save(r, w)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to save session"))
		return err
	}
	return nil
}

func setIdentityCtx(parent context.Context, session *AuthSession) context.Context {
	return context.WithValue(parent, sessionKey, session)
}

func CurrentUser(ctx context.Context) (*Identity, error) {
	if session, ok := ctx.Value(sessionKey).(*AuthSession); ok {
		if session.Identity != nil {
			return session.Identity, nil
		}
	}
	return nil, errors.New("unauthenticated")
}
