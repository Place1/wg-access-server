## Docs
- [x] mkdocs
- [ ] about
- [x] deploying
  - [x] simple docker 1 liner
  - [x] docker-compose
  - [x] kubernetes quickstart
  - [x] helm
- [x] configuring
  - [x] general
  - [x] config file/flag/env
- [ ] how-to-guides
  - [ ] docker + docker-compose
  - [ ] kubernetes + nginx ingress
  - [ ] raspberry-pi + pihole dns

## Features
- [ ] ARM docker image for raspberry-pi
- [ ] admin
  - [x] list all devices
  - [ ] remove device
- [x] networking
  - [x] isolate clients
  - [x] forward to internet only (isolate LAN/WAN)
  - [x] allowed networks (configure forwarding to specific CIDRs)
    - [x] also limit which CIDRs clients forward
    - [x] i.e. only forward to specific server-side LAN and not all internet traffic
