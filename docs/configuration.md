# Configuration

## Environment Variables

| Variable                       | Default        | Description |
|--------------------------------|----------------|-------------|
| WGAS_LOG_LEVEL                 | `info`         | Set the server's log level (debug, **info**, error, critical) |
| WGAS_DISABLE_METADATA          | `false`        | If true, the server will not record device level metadata such as the last handshake time, tx/rx data size |
| WGAS_ADMIN_SUBJECT             | ``             | Set the username (subject) for the admin account |
| WGAS_ADMIN_PASSWORD            | ``             | Set the admin account's password. The admin account will be a basic-auth user. Leave blank if your admin username authenticates via a configured authentication backend. |
| WGAS_WEB_PORT                  | `8000`         | Set the port that the web UI will listen on |
| WGAS_STORAGE                   | `memory://`    | Set the directory where device config will be persisted |
| WGAS_WIREGUARD_INTERFACE_NAME  | `wg0`          | Set the network interface name of the WireGuard netword device |
| WGAS_WIREGUARD_PRIVATE_KEY     | ``             | Set the wireguard private key |
| WGAS_WIREGUARD_EXTERNAL_HOST   | ``             | ExternalAddress is the address that clients use to connect to the wireguard interface |
| WGAS_WIREGUARD_PORT            | `51820`        | Set the WireGuard ListenPort |
| WGAS_VPN_CIDR                  | `10.44.0.0/24` | CIDR configures a network address space that client (WireGuard peers) will be allocated an IP address from |
| WGAS_VPN_GATEWAY_INTERFACE     | ``             | GatewayInterface will be used in iptable forwarding rules that send VPN traffic from clients to this interface Most use-cases will want this interface to have access to the outside internet |
| WGAS_VPN_ALLOWED_IPS           | `0.0.0.0/0`    | The "AllowedIPs" for VPN clients. This value will be included in client config files and in server-side iptable rules to enforce network access. | 
| WGAS_DNS_ENABLED               | `true`         | Enabled allows you to turn on/off the VPN DNS proxy feature. |
| WGAS_DNS_UPSTREAM              | ``             | Set the upstream DNS server to proxy client DNS requests to. If empty, resolv.conf will be respected. |
| WGAS_AUTH_OIDC_NAME            | ``             | |
| WGAS_AUTH_OIDC_ISSUER          | ``             | |
| WGAS_AUTH_OIDC_CLIENT_ID       | ``             | |
| WGAS_AUTH_OIDC_CLIENT_SECRET   | ``             | |
| WGAS_AUTH_OIDC_SCOPES          | ``             | |
| WGAS_AUTH_OIDC_REDIRECT_URL    | ``             | |
| WGAS_AUTH_OIDC_EMAIL_DOMAINS   | ``             | |
| WGAS_AUTH_OIDC_CLAIM_MAPPING   | ``             | |
| WGAS_AUTH_GITLAB_NAME          | ``             | |
| WGAS_AUTH_GITLAB_BASE_URL      | ``             | |
| WGAS_AUTH_GITLAB_CLIENT_ID     | ``             | |
| WGAS_AUTH_GITLAB_CLIENT_SECRET | ``             | |
| WGAS_AUTH_GITLAB_REDIRECT_URL  | ``             | |
| WGAS_AUTH_GITLAB_EMAIL_DOMAINS | ``             | |
| WGAS_AUTH_BASIC_USERS          | ``             | |

## Config File (config.yaml)

Here's an annotated config file example:

```yaml
# The application's log level.
# Can be debug, info, error
# Optional, defaults to info
logLevel: info
# Disable device metadata storage.
# Device metadata includes the last handshake time,
# total sent/received bytes count, their endpoint IP.
# This metadata is captured from wireguard itself.
# Disabling this flag will not stop wireguard from capturing
# this data.
# See stored data here: https://github.com/Place1/wg-access-server/blob/master/internal/storage/contracts.go#L14
# Optional, defaults to false.
disableMetadata: false
# The port that the web ui server (http) will listen on.
# Optional, defaults to 8000
port: 8000
# Directory that VPN devices (WireGuard peers)
# What type of storage do you want? inmemory (default), file:///some/directory, or postgresql, mysql, sqlite3
storage: "memory://"
wireguard:
  # The network interface name for wireguard
  # Optional, defaults to wg0
  interfaceName: wg0
  # The WireGuard PrivateKey
  # You can generate this value using "$ wg genkey"
  # If this value is empty then the server will use an in-memory
  # generated key
  privateKey: ""
  # ExternalAddress is the address (without port) that clients use to connect to the wireguard interface
  # By default, this will be empty and the web ui
  # will use the current page's origin i.e. window.location.origin
  # Optional
  externalHost: ""
  # The WireGuard ListenPort
  # Optional, defaults to 51820
  port: 51820
vpn:
  # CIDR configures a network address space
  # that client (WireGuard peers) will be allocated
  # an IP address from.
  # Optional
  cidr: "10.44.0.0/24"
  # GatewayInterface will be used in iptable forwarding
  # rules that send VPN traffic from clients to this interface
  # Most use-cases will want this interface to have access
  # to the outside internet
  # If not configured then the server will select the default
  # network interface e.g. eth0
  # Optional
  gatewayInterface: ""
  # The "AllowedIPs" for VPN clients.
  # This value will be included in client config
  # files and in server-side iptable rules
  # to enforce network access.
  # Optional
  allowedIPs:
    - "0.0.0.0/0"
dns:
  # Enable a DNS proxy for VPN clients.
  # Optional, Defaults to true
  enabled: true
  # upstream DNS servers.
  # that the server-side DNS proxy will forward requests to.
  # By default /etc/resolv.conf will be used to find upstream
  # DNS servers.
  # Optional
  upstream:
    - "1.1.1.1"
# Auth configures optional authentication backends
# to controll access to the web ui.
# Devices will be managed on a per-user basis if any
# auth backends are configured.
# If no authentication backends are configured then
# the server will not require any authentication.
# It's recommended to make use of basic authentication
# or use an upstream HTTP proxy that enforces authentication
# Optional
auth:
  # HTTP Basic Authentication
  basic:
    # Users is a list of htpasswd encoded username:password pairs
    # supports BCrypt, Sha, Ssha, Md5
    # You can create a user using "htpasswd -nB <username>"
    users: []
  oidc:
    name: "" # anything you want
    issuer: "" # Should point to the oidc url without .well-known
    clientID: ""
    clientSecret: ""
    scopes: null  # list of scopes, defaults to ["openid"]
    redirectURL: "" # full url you want the oidc to redirect to, example: https://vpn-admin.example.com/finish-signin
    # See https://github.com/Knetic/govaluate/blob/9aa49832a739dcd78a5542ff189fb82c3e423116/MANUAL.md for how to write rules
    userClaimsRules:
      admin: "'WireguardAdmins' in group_membership"
    # Optionally restrict login to users with an allowed email domain
    # if empty or omitted, any email domain will be allowed.
    emailDomains:
      - example.com
  gitlab:
    name: ""
    baseURL: ""
    clientID: ""
    clientSecret: ""
    redirectURL: ""
    # Optionally restrict login to users with an allowed email domain
    # if empty or omitted, any email domain will be allowed.
    emailDomains:
      - example.com
```
