package authconfig

import (
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authruntime"
)

type AuthConfig struct {
	OIDC   *OIDCConfig      `yaml:"oidc"`
	Gitlab *GitlabConfig    `yaml:"gitlab"`
	Basic  *BasicAuthConfig `yaml:"basic"`
}

func (c *AuthConfig) IsEnabled() bool {
	return c.OIDC != nil || c.Gitlab != nil || c.Basic != nil
}

func (c *AuthConfig) DesiresSigninPage() bool {
	// Basic auth is the only that truly needs the signin button
	return c.Basic != nil
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
