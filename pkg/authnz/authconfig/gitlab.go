package authconfig

import (
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authruntime"
)

const GitlabAuthProvider = "gitlab"

type GitlabConfig struct {
	Name         string   `yaml:"name"`
	BaseURL      string   `yaml:"baseURL"`
	ClientID     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	RedirectURL  string   `yaml:"redirectURL"`
	EmailDomains []string `yaml:"emailDomains"`
}

func (c *GitlabConfig) Provider() *authruntime.Provider {
	o := OIDCConfig {
		Name:         c.Name,
		Issuer:       c.BaseURL,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       []string{"openid"},
		EmailDomains: c.EmailDomains,
	}
	p := o.Provider()
	p.Type = GitlabAuthProvider
	p.Name = c.Name
	return p
}
