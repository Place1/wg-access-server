package authsession

type ClaimsMiddleware func(user *Identity) error
