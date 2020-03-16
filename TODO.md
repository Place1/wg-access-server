## Docs
- mkdocs
- about
- deploying
  - simple docker 1 liner
  - docker-compose
  - kubernetes quickstart
  - helm
- configuring
  - general
  - config file/flag/env
- how-to-guides
  - docker + docker-compose
  - kubernetes + nginx ingress
  - raspberry-pi + pihole dns

## Features
- ARM docker image for raspberry-pi
- admin
  - list all devices
  - remove device
- networking
  - isolate clients
  - forward to internet only (isolate LAN/WAN)
  - allowed networks (configure forwarding to specific CIDRs)
    - also limit which CIDRs clients forward
    - i.e. only forward to specific server-side LAN and not all internet traffic
