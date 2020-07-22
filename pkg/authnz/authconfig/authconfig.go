package authconfig

import (
	"github.com/place1/wg-access-server/pkg/authnz/authproviders"
	"github.com/place1/wg-access-server/pkg/authnz/authruntime"
)

type AuthConfig struct {
	OIDC   *authproviders.OIDCConfig      `yaml:"oidc"`
	Gitlab *authproviders.GitlabConfig    `yaml:"gitlab"`
	Basic  *authproviders.BasicAuthConfig `yaml:"basic"`
}

func (c *AuthConfig) IsEnabled() bool {
	return c.OIDC != nil || c.Gitlab != nil || c.Basic != nil
}

func (c *AuthConfig) Providers() []*authruntime.Provider {
	providers := []*authruntime.Provider{}

	// must be first because the login page looks nicer
	// with basic auth appearing first.
	if c.Basic != nil {
		providers = append(providers, c.Basic.Provider())
	}

	if c.OIDC != nil {
		providers = append(providers, c.OIDC.Provider())
	}

	if c.Gitlab != nil {
		providers = append(providers, c.Gitlab.Provider())
	}

	return providers
}
