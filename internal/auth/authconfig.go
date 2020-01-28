package auth

type AuthConfig struct {
	OIDC   *OIDCConfig      `yaml:"oidc"`
	Gitlab *GitlabConfig    `yaml:"gitlab"`
	Basic  *BasicAuthConfig `yaml:"basic"`
}

func (c *AuthConfig) Providers() []*Provider {
	providers := []*Provider{}

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
