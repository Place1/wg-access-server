#!/bin/bash
# This script will build the Dockerfile
# and then run it with a minimalistic set of
# docker run arguments
#
# note that "WIREGUARD_PRIVATE_KEY" used in
# this configuration is for the demo and clearly
# not secure, please don't copy-paste it
set -eou pipefail

docker build -t place1/wireguard-access-server .

docker run \
  -it \
  --rm \
  --name wg \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  place1/wireguard-access-server
