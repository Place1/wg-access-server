# wg-access-server

wg-access-server is a single binary that provides a WireGuard
VPN server and device management web ui. We support user authentication,
_1 click_ device registration that works with Mac, Linux, Windows, Ios and Android
including QR codes. You can configure different network isolation modes for
better control and more.

This project aims to deliver a simple VPN solution for developers,
homelab enthusiasts and anyone else feeling adventurous.

wg-access-server is a functional but young project. Contributions are welcome!

## Documentation

[See our documentation website](https://place1.github.io/wg-access-server/)

Quick Links:

- [Configuration Overview](https://place1.github.io/wg-access-server/2-configuration/)
- [Deploy With Docker](https://place1.github.io/wg-access-server/deployment/1-docker/)
- [Deploy With Helm](https://place1.github.io/wg-access-server/deployment/3-kubernetes/)
- [Deploy With Docker-Compose](https://place1.github.io/wg-access-server/deployment/2-docker-compose/)

## Running with Docker

Here's a quick command to run the server to try it out.

```bash
export WG_ADMIN_PASSWORD="example"
export WG_WIREGUARD_PRIVATE_KEY="$(wg genkey)"

docker run \
  -it \
  --rm \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  -v wg-access-server-data:/data \
  -e "WG_ADMIN_PASSWORD=$WG_ADMIN_PASSWORD" \
  -e "WG_WIREGUARD_PRIVATE_KEY=$WG_WIREGUARD_PRIVATE_KEY" \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  place1/wg-access-server
```

If you open your browser using your LAN ip address you can even connect your
phone to try it out: for example, i'll open my browser at http://192.168.0.XX:8000
using the local LAN IP address.

You can connect to the web server on the local machine browser at http://localhost:8000

## Running on Kubernetes via Helm

wg-access-server ships a Helm chart to make it easy to get started on
Kubernetes.

Here's a quick start, but you can read more at the [Helm Chart Deployment Docs](https://place1.github.io/wg-access-server/deployment/3-kubernetes/)

```bash
# deploy
helm install my-release --repo https://place1.github.io/wg-access-server wg-access-server

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

![Devices](https://github.com/Place1/wg-access-server/raw/master/screenshots/devices.png)

![Connect iOS](https://github.com/Place1/wg-access-server/raw/master/screenshots/connect-mobile.png)

![Connect MacOS](https://github.com/Place1/wg-access-server/raw/master/screenshots/connect-desktop.png)

![Sign In](https://github.com/Place1/wg-access-server/raw/master/screenshots/signin.png)

## Changelog

See the [CHANGELOG.md](https://github.com/Place1/wg-access-server/blob/master/CHANGELOG.md) file

## Development

The software is made up a Golang Server and React App.

Here's how I develop locally:

1. run `cd website && npm install && npm start` to get the frontend running on `:3000`
2. run `sudo go run ./main.go` to get the server running on `:8000`

Here are some notes about the development configuration:

- sudo is required because the server uses iptables/ip to configure the VPN networking
- you'll access the website on `:3000` and it'll proxy API requests to `:8000` thanks to webpack
- in-memory storage and generated wireguard keys will be used

GRPC codegeneration:

The client communicates with the server via gRPC-Web. You can edit the API specification
in `./proto/*.proto`.

After changing a service or message definition you'll want to re-generate server and client
code using: `./codegen.sh`.
