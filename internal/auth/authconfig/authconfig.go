package authconfig

import "github.com/place1/wireguard-access-server/internal/auth/authruntime"

type AuthConfig struct {
	OIDC   *OIDCConfig      `yaml:"oidc"`
	Gitlab *GitlabConfig    `yaml:"gitlab"`
	Basic  *BasicAuthConfig `yaml:"basic"`
}

func (c *AuthConfig) Providers() []*authruntime.Provider {
	providers := []*authruntime.Provider{}

	if c.OIDC != nil {
		providers = append(providers, c.OIDC.Provider())
	}

	if c.Gitlab != nil {
		providers = append(providers, c.Gitlab.Provider())
	}

	if c.Basic != nil {
		providers = append(providers, c.Basic.Provider())
	}

	return providers
}
