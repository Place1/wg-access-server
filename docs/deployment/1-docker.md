# Docker

## TL;DR;

Here's a one-liner to run wg-access-server:

```bash
docker run --rm -it \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  place1/wg-access-server
```
