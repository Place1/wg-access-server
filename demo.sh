#!/bin/bash
# This script will build the Dockerfile
# and then run it with a minimalistic set of
# docker run arguments
#
# note that "WIREGUARD_PRIVATE_KEY" used in
# this configuration is for the demo and clearly
# not secure, please don't copy-paste it
set -eo pipefail

if [[ -z $1 ]]; then
  echo "USAGE: $0 <path-to-config-file>"
  exit 1
fi

CONFIG_FILE="$1"

docker build -t place1/wireguard-access-server .

docker run \
  -it \
  --rm \
  --name wg \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  -v "$CONFIG_FILE:/config.yaml" \
  -v demo-data:/data \
  -e "LOG_LEVEL=Debug" \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  -p 53:53/udp \
  place1/wireguard-access-server /server --config /config.yaml
