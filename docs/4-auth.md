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

| Backend        | Use Case                                                                                      | Notes                                                               |
|----------------|-----------------------------------------------------------------------------------------------|---------------------------------------------------------------------|
| Simple Auth    | Deployments with a static list of users. Simple and great for self-hosters and home use-cases | Recommended, default for the admin account                          |
| Basic Auth     | Like Simple Auth, but using HTTP Basic Auth for login                                         | Logout does not work because browsers caches Basic Auth credentials |
| OpenID Connect | For delegating authentication to an existing identity solution                                |                                                                     |
| Gitlab         | For delegating authentication to gitlab. Supports self-hosted Gitlab.                         |                                                                     |

If `adminPassword` is set, an administrator account will be added with the username of `adminUsername` (default `admin`)
to the Simple Auth or Basic Auth backend; whichever is enabled, automatically enabling Simple if both are unset,
preferring Simple to Basic if both are enabled.

## Configuration

Currently authentication providers are only configurable via the wg-access-server
config file (config.yaml).

Below is an annotated example config section that can be used as a starting point.

```yaml
# You can disable the builtin admin account by leaving out 'adminPassword'. Requires another backend to be configured.
adminPassword: "<admin password>"
# adminUsername sets the user for the Basic/Simple Auth admin account if adminPassword is set.
# Every user of the basic and simple backend with a username matching adminUsername will have admin privileges.
adminUsername: "admin"
# Configure zero or more authentication backends
auth:
  sessionStore:
    # 32 random bytes in hexadecimal encoding (64 chars) used to sign session cookies. It's generated randomly
    # if not present. Need to be set when running in HA setup (more than one replica)
    secret: "<session store secret>"
  simple:
    # Users is a list of htpasswd encoded username:password pairs
    # supports BCrypt, Sha, Ssha, Md5
    # You can create a user using "htpasswd -nB <username>"
    users: []
  # HTTP Basic Authentication
  basic:
    # Users is a list of htpasswd encoded username:password pairs
    # supports BCrypt, Sha, Ssha, Md5
    # You can create a user using "htpasswd -nB <username>"
    users: []
  oidc:
    # A name for the backend (is shown on the login page and possibly in the devices list of the 'all devices' admin page)
    name: "My OIDC Backend"
    # Should point to the OIDC Issuer (excluding /.well-known/openid-configuration)
    issuer: "https://identity.example.com"
    # Your OIDC client credentials which would be provided by your OIDC provider
    clientID: "<client-id>"
    clientSecret: "<client-secret>"
    # The full redirect URL
    # The path can be almost anything as long as it doesn't
    # conflict with a path that the web UI uses.
    # /callback is recommended.
    redirectURL: "https://wg-access-server.example.com/callback"
    # List of scopes to request claims for. Must include 'openid'.
    # Must include 'email' if 'emailDomains' is used. Can include 'profile' to show the user's name in the UI.
    # Add custom ones if required for 'claimMapping'.
    # Defaults to ["openid"]
    scopes:
      - openid
      - profile
      - email
    # You can optionally restrict access to users with an email address
    # that matches an allowed domain.
    # If empty or omitted then all email domains will be allowed.
    emailDomains:
      - example.com
    # This is an advanced feature that allows you to define OIDC claim mapping expressions.
    # This feature is used to define wg-access-server admins based off a claim in your OIDC token.
    # A JSON-like object of claimKey: claimValue pairs as returned by the issuer is passed to the evaluation function. 
    # See https://github.com/Knetic/govaluate/blob/9aa49832a739dcd78a5542ff189fb82c3e423116/MANUAL.md for the syntax.
    claimMapping:
      # This example works if you have a custom group_membership claim which is a list of strings 
      admin: "'WireguardAdmins' in group_membership"
      access: "'WireguardAccess' in group_membership"
    # Let wg-access-server retrieve the claims from the ID Token instead of querying the UserInfo endpoint.
    # Some OIDC authorization provider implementations (e.g. ADFS) only publish claims in the ID Token.
    claimsFromIDToken: false
    # require this claim to be "true" to allow access for the user
    accessClaim: "access"
  gitlab:
    name: "My Gitlab Backend"
    baseURL: "https://mygitlab.example.com"
    clientID: "<client-id>"
    clientSecret: "<client-secret>"
    redirectURL: "https:///wg-access-server.example.com/callback"
    emailDomains:
      - example.com
```

## OIDC Provider specifics

### Active Directory Federation Services (ADFS)

Please see [this helpful issue comment](https://github.com/freifunkMUC/wg-access-server/issues/213#issuecomment-1172656633) for instructions for ADFS 2016 and above.
