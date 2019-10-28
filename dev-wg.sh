#!/bin/bash
set -eou pipefail

docker build -t dev-wg --target boringtun .

docker run \
  --rm \
  -it \
  --network host \
  --device /dev/net/tun:/dev/net/tun \
  --cap-add NET_ADMIN \
  -v /var/run/wireguard:/var/run/wireguard \
  dev-wg \
    boringtun wg0 --disable-drop-privileges=root --foreground --verbosity=debug
