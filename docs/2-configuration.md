# Configuration

You can configure wg-access-server using environment variables, cli flags or a config file
taking precedence over one another in that order.

The default configuration should work out of the box if you're just looking to try it out.

The only required configuration is a wireguard private key.
You can generate a wireguard private key by [following the official docs](https://www.wireguard.com/quickstart/#key-generation).

TLDR:

```bash
wg genkey
```

The config file format is `yaml` and an example is provided [below](#the-config-file-configyaml).

The format for specifying multiple values for options that allow it is:
* as commandline flags:
  * repeat the flag (e.g. `--dns-upstream 2001:db8::1 --dns-upstream 192.0.2.1`)
  * separate the values with a comma (e.g. `--dns-upstream 2001:db8::1,192.0.2.1`)
* as environment variables:
  * separate with a comma (e.g. `WG_DNS_UPSTREAM="2001:db8::1,192.0.2.1"`)
  * separate with a new line char (e.g. `WG_DNS_UPSTREAM=$'2001:db8::1\n192.0.2.1'`)
* in the config file as YAML list.

Here's what you can configure:

| Environment Variable       | CLI Flag                   | Config File Path       | Required | Default (docker)                             | Description                                                                                                                                                                        |
|----------------------------|----------------------------|------------------------|----------|----------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `WG_CONFIG`                | `--config`                 |                        |          |                                              | The path to a wg-access-server config.yaml file                                                                                                                                    |
| `WG_LOG_LEVEL`             | `--log-level`              | `logLevel`             |          | `info`                                       | The global log level                                                                                                                                                               |
| `WG_ADMIN_USERNAME`        | `--admin-username`         | `adminUsername`        |          | `admin`                                      | The admin account username                                                                                                                                                         |
| `WG_ADMIN_PASSWORD`        | `--admin-password`         | `adminPassword`        | Yes      |                                              | The admin account password                                                                                                                                                         |
| `WG_PORT`                  | `--port`                   | `port`                 |          | `8000`                                       | The port the web ui will listen on (http)                                                                                                                                          |
| `WG_EXTERNAL_HOST`         | `--external-host`          | `externalHost`         |          |                                              | The external domain for the server (e.g. www.mydomain.com)                                                                                                                         |
| `WG_STORAGE`               | `--storage`                | `storage`              |          | `sqlite3:///data/db.sqlite3`                 | A storage backend connection string. See [storage docs](./3-storage.md)                                                                                                            |
| `WG_DISABLE_METADATA`      | `--disable-metadata`       | `disableMetadata`      |          | `false`                                      | Turn off collection of device metadata logging. Includes last handshake time and RX/TX bytes only.                                                                                 |
| `WG_WIREGUARD_ENABLED`     | `--[no-]wireguard-enabled` | `wireguard.enabled`    |          | `true`                                       | Enable/disable the wireguard server. Useful for development on non-linux machines.                                                                                                 |
| `WG_WIREGUARD_INTERFACE`   | `--wireguard-interface`    | `wireguard.interface`  |          | `wg0`                                        | The wireguard network interface name                                                                                                                                               |
| `WG_WIREGUARD_PRIVATE_KEY` | `--wireguard-private-key`  | `wireguard.privateKey` | Yes      |                                              | The wireguard private key. This value is required and must be stable. If this value changes all devices must re-register.                                                          |
| `WG_WIREGUARD_PORT`        | `--wireguard-port`         | `wireguard.port`       |          | `51820`                                      | The wireguard server port (udp)                                                                                                                                                    |
| `WG_VPN_CIDR`              | `--vpn-cidr`               | `vpn.cidr`             |          | `10.44.0.0/24`                               | The VPN IPv4 network range. VPN clients will be assigned IP addresses in this range. Set to `0` to disable IPv4.                                                                   |
| `WG_IPV4_NAT_ENABLED`      | `--vpn-nat44-enabled`      | `vpn.nat44`            |          | `true`                                       | Disables NAT for IPv4                                                                                                                                                              |
| `WG_IPV6_NAT_ENABLED`      | `--vpn-nat66-enabled`      | `vpn.nat66`            |          | `true`                                       | Disables NAT for IPv6                                                                                                                                                              |
| `WG_VPN_CLIENT_ISOLATION`  | `--vpn-client-isolation`   | `vpn.clientIsolation`  |          | `false`                                      | BLock or allow traffic between client devices (client isolation)                                                                                                                   |
| `WG_VPN_CIDRV6`            | `--vpn-cidrv6`             | `vpn.cidrv6`           |          | `fd48:4c4:7aa9::/64`                         | The VPN IPv6 network range. VPN clients will be assigned IP addresses in this range. Set to `0` to disable IPv6.                                                                   |
| `WG_VPN_GATEWAY_INTERFACE` | `--vpn-gateway-interface`  | `vpn.gatewayInterface` |          | _default gateway interface (e.g. eth0)_      | The VPN gateway interface. VPN client traffic will be forwarded to this interface.                                                                                                 |
| `WG_VPN_ALLOWED_IPS`       | `--vpn-allowed-ips`        | `vpn.allowedIPs`       |          | `0.0.0.0/0, ::/0`                            | Allowed IPs that clients may route through this VPN. This will be set in the client's WireGuard connection file and routing is also enforced by the server using iptables.         |
| `WG_DNS_ENABLED`           | `--[no-]dns-enabled`       | `dns.enabled`          |          | `true`                                       | Enable/disable the embedded DNS proxy server. This is enabled by default and allows VPN clients to avoid DNS leaks by sending all DNS requests to wg-access-server itself.         |
| `WG_DNS_UPSTREAM`          | `--dns-upstream`           | `dns.upstream`         |          | _resolvconf autodetection or Cloudflare DNS_ | The upstream DNS servers to proxy DNS requests to. By default the host machine's resolveconf configuration is used to find its upstream DNS server, with a fallback to Cloudflare. |
| `WG_DNS_DOMAIN`            | `--dns-domain`             | `dns.domain`           |          |                                              | A domain to serve configured devices authoritatively. Queries for names in the format <device>.<user>.<domain> will be answered with the device's IP addresses.                    |

## The Config File (config.yaml)

Here's an example config file to get started with.

```yaml
loglevel: info
storage: sqlite3:///data/db.sqlite3
wireguard:
  privateKey: "<some-key>"
dns:
  upstream:
    - "2001:4860:4860::8888"
    - "8.8.8.8"
```
