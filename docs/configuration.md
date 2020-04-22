# Configuration

## Environment Variables

| Variable              | Description |
|-----------------------|-------------|
| CONFIG                | Set the config file path |
| WIREGUARD_PRIVATE_KEY | Set the wireguard private key |
| STORAGE_DIRECTORY     | Set the directory where device config will be persisted |
| ADMIN_USERNAME        | Set the username (subject) for the admin account |
| ADMIN_PASSWORD        | Set the admin account's password. The admin account will be a basic-auth user. Leave blank if your admin username authenticates via a configured authentication backend. |
| UPSTREAM_DNS          | Set the upstream DNS server to proxy client DNS requests to. If empty, resolv.conf will be respected. |
| LOG_LEVEL             | Set the server's log level (debug, **info**, error, critical) |
| DISABLE_METADATA      | If true, the server will not record device level metadata such as the last handshake time, tx/rx data size |

## CLI Flags

All environment variables can be configured via a
CLI flag as well.

For example you can configure `STORAGE_DIRECTORY` by passing `--storage-directory="<value>"`.

## Config File (config.yaml)

Here's an annotated config file example:

```yaml
# The application's log level.
# Can be debug, info, error
# Optional, defaults to info
loglevel: info
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
storage:
  # Directory that VPN devices (WireGuard peers)
  # should be saved under.
  # If this value is empty then an InMemory storage
  # backend will be used (not recommended).
  # Optional
  # Defaults to in-memory
  # The docker container sets this value to /data automatically
  directory: /data
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
  # Rules allows you to configure what level
  # of network isolation should be enfoced.
  rules:
    # AllowVPNLAN enables routing between VPN clients
    # i.e. allows the VPN to work like a LAN.
    # true by default
    # Optional
    allowVPNLAN: true
    # AllowServerLAN enables routing to private IPv4
    # address ranges. Enabling this will allow VPN clients
    # to access networks on the server's LAN.
    # true by default
    # Optional
    allowServerLAN: true
    # AllowInternet enables routing of all traffic
    # to the public internet.
    # true by default
    # Optional
    allowInternet: true
    # AllowedNetworks allows you to whitelist a partcular
    # network CIDR. This is useful if you want to block
    # access to the Server's LAN but allow access to a few
    # specific IPs or a small range.
    # e.g. "192.0.2.0/24" or "192.0.2.10/32".
    # no networks are whitelisted by default (empty array)
    # Optional
    allowedNetworks: []
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
