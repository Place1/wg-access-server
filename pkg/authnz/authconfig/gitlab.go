package authconfig

import "github.com/place1/wg-access-server/pkg/authnz/authruntime"

type GitlabConfig struct {
	Enabled      bool     `yaml:"enabled"`
	Name         string   `yaml:"name"`
	BaseURL      string   `yaml:"baseURL" split_words:"true"`
	ClientID     string   `yaml:"clientID" split_words:"true"`
	ClientSecret string   `yaml:"clientSecret" split_words:"true"`
	RedirectURL  string   `yaml:"redirectURL" split_words:"true"`
	EmailDomains []string `yaml:"emailDomains" split_words:"true"`
}

func (c *GitlabConfig) Provider() *authruntime.Provider {
	o := OIDCConfig{
		Enabled:      c.Enabled,
		Name:         c.Name,
		Issuer:       c.BaseURL,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       []string{"openid"},
		EmailDomains: c.EmailDomains,
	}
	p := o.Provider()
	p.Type = "Gitlab"
	return p
}
