package auth

import (
	"fmt"

	"github.com/dexidp/dex/storage"
)

type Config struct {
	StaticUsers []StaticUser
	Connectors  []AuthConnector
}

type StaticUser struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type AuthConnector interface {
	toDexConnector(externalAddress string) storage.Connector
}

// implements toDexConnector
type OIDCConfig struct {
	Name         string `yaml:"name"`
	Issuer       string `yaml:"issuer"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

func (c *OIDCConfig) toDexConnector(externalAddr string) storage.Connector {
	return storage.Connector{
		ID:              storage.NewID(),
		Type:            "oidc",
		Name:            c.Name,
		ResourceVersion: "1",
		Config: []byte(fmt.Sprintf(`{
			"issuer": "%s",
			"redirectURI": "%s/auth/callback",
			"clientID": "%s",
			"clientSecret": "%s"
		}`, c.Issuer, externalAddr, c.ClientID, c.ClientSecret)),
	}
}

// implements toDexConnector
type GitlabConfig struct {
	Name         string `yaml:"name"`
	BaseURL      string `yaml:"baseURL"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

func (c *GitlabConfig) toDexConnector(externalAddr string) storage.Connector {
	return storage.Connector{
		ID:              "gitlab",
		Type:            "gitlab",
		Name:            c.Name,
		ResourceVersion: "1",
		Config: []byte(fmt.Sprintf(`{
			"baseURL": "%s",
			"redirectURI": "%s/auth/callback",
			"clientID": "%s",
			"clientSecret": "%s"
		}`, c.BaseURL, externalAddr, c.ClientID, c.ClientSecret)),
	}
}

// TODO: others
