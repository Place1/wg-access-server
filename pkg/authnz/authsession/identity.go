package authsession

type Identity struct {
	// Provider is the name of the authentication provider
	// that authenticated (created) this Identity struct.
	Provider string
	// Subject is the canonical identifier for this Identity.
	Subject string
	// Name is the name of the person this Identity refers to.
	// It may be empty.
	Name string
	// Email is the email address of the person this Identity refers to.
	// It may be empty.
	Email string
	// Claims are any additional claims that middleware have
	// added to this Identity.
	Claims Claims
}
