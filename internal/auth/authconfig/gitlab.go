package authconfig

import "github.com/place1/wireguard-access-server/internal/auth/authruntime"

type GitlabConfig struct {
	Name         string `yaml:"name"`
	BaseURL      string `yaml:"baseURL"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectURL  string `yaml:"redirectURL"`
}

func (c *GitlabConfig) Provider() *authruntime.Provider {
	o := OIDCConfig{
		Name:         c.Name,
		Issuer:       c.BaseURL,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       []string{"openid"},
	}
	return o.Provider()
}
