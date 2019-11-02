# WireGuard Access Server

_i'm still thinking of a name..._

## What is this

This project aims to create a simple VPN solution for developers,
homelab enthusiasts and perhaps some adventurous small businesses.

This project offers a single docker container that provides a WireGuard
VPN server and device management web ui that's simple to use.

Today, this project allows you to deploy a WireGuard VPN using a single
docker container; use a web ui to add/connect your Linux/Mac/Windows/iOS/Android
device; and manage connected devices. The server will automatically
configure ip routes and iptables rules to ensure that client VPN traffic
can access the internet.

Soon I hope to add the following features

- [ ] headless mode
  * in this mode there'll be no web ui
  * you can add devices (i.e. WireGuard peers) via files, flags or the environment
  * intended for use by developers to easily deploy a one-shot style
    VPN into a network to get access to it on their local machine,
    i'm hoping to use this mode to VPN into a kubernetes cluster's
    overlay network including DNS and cluster service routing.
- [ ] singleuser mode
  * this is how the project currently works but I'll expand it to support authentication
- [ ] multiuser mode
  * support pluggable authentication backends including OAuth, OpenID Connect, LDAP, etc.
  * allow different users to manage thier own devices without seeing others
  * allow network isolation to be turned on or off allowing users to communicate or be isolated

## Running with Docker

```
# Note that the hard coded WIREGUARD_PRIVATE_KEY is clearly
# not secure, please don't accidently copy-paste it to a real
# system.
# You can generate a wireguard private key using:
#    umask 077; wg genkey > privatekey
docker run \
  -it \
  --rm \
  --name wg \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  --network host \
  -v wgdata:/data \
  -e WIREGUARD_PRIVATE_KEY="kH4F1lldSzgEMB7wfQ1ccujAhZCCCCEeh2Kvhxf+XFw=" \
  place1/wireguard-access-server:0.0.1
```

## Screenshots

![IOS Connection Dialog](./screenshots/get-connected-ios.png)

![Windows Connection Dialog](./screenshots/get-connected.png)

## Development

The software is made up a Golang server, React webapp and a WireGuard
implementation that must be provided by the system.

I'll add more instructions here soon :)
