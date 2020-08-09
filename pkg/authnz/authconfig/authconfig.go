package authconfig

import (
	"github.com/place1/wg-access-server/pkg/authnz/authruntime"
)

type AuthConfig struct {
	OIDC   *OIDCConfig      `yaml:"oidc"`
	Gitlab *GitlabConfig    `yaml:"gitlab"`
	Basic  *BasicAuthConfig `yaml:"basic"`
}

func (c *AuthConfig) IsEnabled() bool {
	return c.OIDC.Enabled || c.Gitlab.Enabled || c.Basic.Enabled
}

func (c *AuthConfig) Providers() []*authruntime.Provider {
	providers := []*authruntime.Provider{}

	if c.OIDC.Enabled {
		providers = append(providers, c.OIDC.Provider())
	}

	if c.Gitlab.Enabled {
		providers = append(providers, c.Gitlab.Provider())
	}

	if c.Basic.Enabled {
		providers = append(providers, c.Basic.Provider())
	}

	return providers
}
