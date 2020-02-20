# wg-access-server

## What is this

This project aims to create a simple VPN solution for developers,
homelab enthusiasts and anyone else feeling adventurous.

This project offers a single docker container that provides a WireGuard
VPN server and device management web ui.

You can use wg-access-server's web ui to connect your Linux/Mac/Windows/iOS/Android
devices. The server automatically configure iptables rules to ensure that client VPN traffic
can access the internet via the server's default gateway or configured gateway NIC.
Currently, all VPN clients can route traffic to each other. VPN client isolation via
iptables can be added if there's demand for it.

wg-access-server embeds a user-space wireguard implementation to simplify
deployment - you just run the container, no kernel setup required.

Support for the kernal's wireguard implementation could be added if
there's demand for it.

Currently wg-access-server requires `NET_ADMIN` and access to `/dev/net/tun` to create
a user-space virtual network interface ([wikipedia](https://en.wikipedia.org/wiki/TUN/TAP)).

wg-access-server also configures iptables and network routes within it's own network
namespace to route client VPN traffic. The container doesn't require host networking
but it can be enabled if you want client VPN traffic to be able to access the host's
network as well.

## Running with Docker

Here's a quick command to run the server to try it out.

If you open your browser using your LAN ip address you can even connect your
phone to try it out: for example, i'll open my browser at http://192.168.0.15:8000
using my laptop's LAN IP address.

```
docker run \
  -it \
  --rm \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  -v wg-access-server-data:/data \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  place1/wg-access-server
```

## Configuration

You can configure the server using a yaml configuration file. Just mount the file into the container like this:

```
docker run \
  ... \
  -v $(pwd)/config.yaml:/config.yaml \
  place1/wg-access-server
```

Here's and example showing the recommended config:

```yaml
wireguard:
  // The WireGuard PrivateKey
  // You can generate this value using "$ wg genkey"
  // If this value is empty then the server will use an in-memory
  // generated key
  privateKey: ""
// Auth configures optional authentication backends
// to controll access to the web ui.
// Devices will be managed on a per-user basis if any
// auth backends are configured.
// If no authentication backends are configured then
// the server will not require any authentication.
// It's recommended to make use of basic authentication
// or use an upstream HTTP proxy that enforces authentication
// Optional
auth:
  // HTTP Basic Authentication
  basic:
    // Users is a list of htpasswd encoded username:password pairs
    // supports BCrypt, Sha, Ssha, Md5
    // You can create a user using "htpasswd -nB <username>"
    users: []
```

Here's an example showing the all config values:

```yaml
loglevel: debug
storage:
  // Directory that VPN devices (WireGuard peers)
  // should be saved under.
  // If this value is empty then an InMemory storage
  // backend will be used (not recommended).
  // Defaults to "/data" inside the docker container
  directory: /data
wireguard:
  // The network interface name for wireguard
  // Optional
  interfaceName: wg0
  // The WireGuard PrivateKey
  // You can generate this value using "$ wg genkey"
  // If this value is empty then the server will use an in-memory
  // generated key
  privateKey: ""
  // ExternalAddress is the address that clients
  // use to connect to the wireguard interface
  // By default, this will be empty and the web ui
  // will use the current page's origin i.e. window.location.origin
  // Optional
  externalHost: ""
  // The WireGuard ListenPort
  // Optional
  port: 51820
} `yaml:"wireguard"`
vpn:
  // CIDR configures a network address space
  // that client (WireGuard peers) will be allocated
  // an IP address from.
  // Optional
  cidr: "10.44.0.0/24"
  // GatewayInterface will be used in iptable forwarding
  // rules that send VPN traffic from clients to this interface
  // Most use-cases will want this interface to have access
  // to the outside internet
  // If not configured then the server will select the default
  // network interface e.g. eth0
  // Optional
  gatewayInterface: ""
dns:
  // upstream DNS servers.
  // that the server-side DNS proxy will forward requests to.
  // By default /etc/resolv.conf will be used to find upstream
  // DNS servers.
  // Optional
  upstream:
    - "1.1.1.1"
// Auth configures optional authentication backends
// to controll access to the web ui.
// Devices will be managed on a per-user basis if any
// auth backends are configured.
// If no authentication backends are configured then
// the server will not require any authentication.
// It's recommended to make use of basic authentication
// or use an upstream HTTP proxy that enforces authentication
// Optional
auth:
  // HTTP Basic Authentication
  basic:
    // Users is a list of htpasswd encoded username:password pairs
    // supports BCrypt, Sha, Ssha, Md5
    // You can create a user using "htpasswd -nB <username>"
    users: []
  oidc:
    name: ""
    issuer: ""
    clientID: ""
    clientSecret: ""
    scopes: ""
    redirectURL: ""
  gitlab:
    name: ""
    baseURL: ""
    clientID: ""
    clientSecret: ""
    redirectURL: ""
```

## Screenshots

![Connect iOS](./screenshots/connect-ios.png)

![Connect MacOS](./screenshots/connect-macos.png)

![Devices](./screenshots/devices.png)

![Sign In](./screenshots/signin.png)

## Roadmap

- [ ] Implement administration features
  - administration of all devices
  - see when a device last connected
  - see owns the device
- [ ] VPN network client isolation
- [ ] ??? PRs, feedback, suggestions welcome

## Development

The software is made up a Golang Server and React App.

Here's how I develop locally:

2. run `cd website && npm install && npm start` to get the frontend running on `:3000`
3. run `sudo go run ./main.go` to get the server running on `:8000`

Here are some notes about the development configuration:

- sudo is required because the server uses iptables/ip to configure the VPN networking
- you'll access the website on `:3000` and it'll proxy API requests to `:8000` thanks to webpack
- in-memory storage and generated wireguard keys will be used

GRPC codegeneration:

The client communicates with the server via gRPC-Web. You can edit the API specification
in `./proto/*.proto`.

After changing a service or message definition you'll want to re-generate server and client
code using: `./codegen.sh`.
