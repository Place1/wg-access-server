# wg-access-server

wg-access-server is a single binary file that contains a WireGuard
VPN server and a web user interface for device management. We support user authentication,
_1-click_ device enrollment that works with macOS, Linux, Windows, iOS/iPadOS and Android
including QR codes. Furthermore, you can choose from different network isolation modes for a
better control over connected devices. Generally speaking you can customize the project
to your use-case with relative ease.

This project aims to provide a simple VPN solution for developers,
homelab enthusiasts, and anyone else who is adventurous.

**This is a fork of the original work of place1, maintained by [Freifunk Munich](https://ffmuc.net/).
Since the upstream is currently unmaintained, we try to add new features and keep the project up to date and in a working state.**

This fork supports IPv6. The VPN can run in dual-stack, IPv6-only or IPv4-only mode.
NAT can be disabled separately for IPv4 and IPv6.

**Contributions are always welcome so that we can offer new bug fixes, features and improvements to the users of this project**.

## Features

- Pluggable authentication using OpenID Connect
- Authentication using GitLab
- IPv6 support in tunnel
- Caching DNS proxy (stub resolver)
- WireGuard client configuration QR codes
- PostgreSQL, MySQL or SQLite3 storage backend

## Documentation

[See our documentation website](https://www.freie-netze.org/wg-access-server/)

Quick Links:

- [Configuration Overview](https://www.freie-netze.org/wg-access-server/2-configuration)
- [Deploy With Docker](https://www.freie-netze.org/wg-access-server/deployment/1-docker)
- [Deploy With Helm](https://www.freie-netze.org/wg-access-server/deployment/2-docker-compose)
- [Deploy With Docker-Compose](https://www.freie-netze.org/wg-access-server/deployment/2-docker-compose)

## Running with Docker

Here is a quick command to start the wg-access-server for the first time and try it out.

```bash
export WG_ADMIN_PASSWORD="example"
export WG_WIREGUARD_PRIVATE_KEY="$(wg genkey)"

docker run \
  -it \
  --rm \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  --sysctl net.ipv6.conf.all.disable_ipv6=0 \
  --sysctl net.ipv6.conf.all.forwarding=1 \
  -v wg-access-server-data:/data \
  -v /lib/modules:/lib/modules \
  -e "WG_ADMIN_PASSWORD=$WG_ADMIN_PASSWORD" \
  -e "WG_WIREGUARD_PRIVATE_KEY=$WG_WIREGUARD_PRIVATE_KEY" \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  ghcr.io/freifunkmuc/wg-access-server:latest
```

If the wg-access-server is accessible via LAN or a network you are in, you can directly connect your phone to the VPN. You have to call the webfrontent of the project for this. Normally, this is done via the IP address of the device or server on which the wg-access-server is running followed by the standard port 8000, via which the web interface can be reached. For most deployments something like this should work: http://192.168.0.XX:8000

If the project is running locally on the computer, you can easily connect to the web interface by connecting to http://localhost:8000 in the browser.

## Running on Kubernetes via Helm

wg-access-server ships a Helm chart to make it easy to get started on
Kubernetes.

Here's a quick start, but you can read more at the [Helm Chart Deployment Docs](https://freifunkMUC.github.io/wg-access-server/deployment/3-kubernetes/)

```bash
# deploy
helm install my-release --repo https://freifunkMUC.github.io/wg-access-server wg-access-server

# cleanup
helm delete my-release
```

## Running with Docker-Compose

Download the the docker-compose.yml file from the repo and run the following command.

```bash
export WG_ADMIN_PASSWORD="example"
export WG_WIREGUARD_PRIVATE_KEY="$(wg genkey)"

docker-compose up
```

You can connect to the web server on the local machine browser at http://localhost:8000

If you open your browser to your machine's LAN IP address you'll be able
to connect your phone using the UI and QR code!

## Screenshots

![Devices](https://github.com/freifunkMUC/wg-access-server/raw/master/screenshots/devices.png)

![Connect iOS](https://github.com/freifunkMUC/wg-access-server/raw/master/screenshots/connect-mobile.png)

![Connect MacOS](https://github.com/freifunkMUC/wg-access-server/raw/master/screenshots/connect-desktop.png)

![Sign In](https://github.com/freifunkMUC/wg-access-server/raw/master/screenshots/signin.png)

## Changelog

See the [CHANGELOG.md](https://github.com/freifunkMUC/wg-access-server/blob/master/CHANGELOG.md) file

## Development

The software consists of a Golang server and a React app.

If you want to make changes to the project locally, you can do so relatively easily with the following steps.

1. Run `cd website && npm install && npm start` to get the frontend running on `:3000`.
2. Run `sudo go run ./main.go` to get the server running on `:8000`.

Here are some notes on development configuration:

- sudo is required because the server uses iptables/ip to configure the VPN network
- access to the website is on `:3000` and API requests are redirected to `:8000` thanks to webpack
- in-memory storage and generated WireGuard keys are used

gRPC code generation:

The client communicates with the server via gRPC web. You can edit the API specification in `./proto/*.proto`.

After changing a service or message definition, you must regenerate the server and client code using: `./codegen.sh`.
