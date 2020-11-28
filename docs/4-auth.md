# Authentication

Authentication is pluggable in wg-access-server. Community contributions are welcome
for supporting new authentication backends.

If you're just getting started you can skip over this section and rely on the default
admin account instead.

If your authentication system is not yet supported and you aren't quite ready to
contribute you could try using a project like [dex](https://github.com/dexidp/dex)
or SaaS provider like [Auth0](https://auth0.com/) which supports a wider variety of
authentication protocols. wg-access-server can happily be an OpenID Connect client
to a larger solution like this.

The following authentication backends are currently supported:

| Backend        | Use Case                                                                                      | Notes                                                         |
| -------------- | --------------------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| Basic Auth     | Deployments with a static list of users. Simple and great for self-hosters and home use-cases | The wg-access-server admin account is powered by this backend |
| OpenID Connect | For delegating authentication to an existing identity solution                                |                                                               |
| Gitlab         | For delegating authentication to gitlab. Supports self-hosted Gitlab.                         |                                                               |

## Configuration

Currently authentication providers are only configurable via the wg-access-server
config file (config.yaml).

Below is an annotated example config section that can be used as a starting point.

```yaml
# Configure zero or more authentication backends
auth:
  # HTTP Basic Authentication
  basic:
    # Users is a list of htpasswd encoded username:password pairs
    # supports BCrypt, Sha, Ssha, Md5
    # You can create a user using "htpasswd -nB <username>"
    users: []
  oidc:
    # A name for the backend (can be anything you want)
    name: "My OIDC Backend"
    # Should point to the OIDC Issuer (excluding /.well-known/openid-configuration)
    issuer: "https://identity.example.com"
    # Your OIDC client credentials which would be provided by your OIDC provider
    clientID: "<client-id>"
    clientSecret: "<client-secret>"
    # List of scopes to request defaults to ["openid"]
    scopes:
      - openid
    # The full redirect URL
    # The path can be almost anything as long as it doesn't
    # conflict with a path that the web UI uses.
    # /callback is recommended.
    redirectURL: "https://wg-access-server.example.com/callback"
    # You can optionally restrict access to users with an email address
    # that matches an allowed domain.
    # If empty or omitted then all email domains will be allowed.
    emailDomains:
      - example.com
    # This is an advanced feature that allows you to define
    # OIDC claim mapping expressions.
    # This feature is used to define wg-access-server admins
    # based off a claim in your OIDC token
    # See https://github.com/Knetic/govaluate/blob/9aa49832a739dcd78a5542ff189fb82c3e423116/MANUAL.md for how to write rules
    claimMapping:
      admin: "'WireguardAdmins' in group_membership"
  gitlab:
    name: "My Gitlab Backend"
    baseURL: "https://mygitlab.example.com"
    clientID: "<client-id>"
    clientSecret: "<client-secret>"
    redirectURL: "https:///wg-access-server.example.com/callback"
    emailDomains:
      - example.com
```
