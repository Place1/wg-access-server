#!/bin/bash
# This script will build the Dockerfile
# and then run it with a minimalistic set of
# docker run arguments
#
# note that "WIREGUARD_PRIVATE_KEY" used in
# this configuration is for the demo and clearly
# not secure, please don't copy-paste it
set -eou pipefail

docker build -t demo .

# read -p "Enter LAN ip address i.e. 192.168.0.2 : " external_address

# docker run \
#   -it \
#   --rm \
#   --name wg \
#   --cap-add NET_ADMIN \
#   --device /dev/net/tun:/dev/net/tun \
#   -v wgdata:/data \
#   -p 8000:8000/tcp \
#   -p 51820:51820/udp \
#   -e WIREGUARD_PRIVATE_KEY="kH4F1lldSzgEMB7wfQ1ccujAhZCCCCEeh2Kvhxf+XFw=" \
#   -e WEB_EXTERNAL_ADDRESS="$external_address" \
#   demo

docker run \
  -it \
  --rm \
  --name wg \
  --network host \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  -v wgdata:/data \
  -v "$(pwd)"/config-demo.yaml:/config-demo.yaml \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  demo /server --config /config-demo.yaml
