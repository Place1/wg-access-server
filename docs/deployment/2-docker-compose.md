# Docker Compose

You can run wg-access-server using the following example Docker Compose file.

Checkout the [configuration docs](../2-configuration.md) to learn how wg-access-server can be configured.

Please also read the [Docker instructions](1-docker.md) for general information regarding Docker deployments.

```yaml
{!../docker-compose.yml!}
```

## IPv4-only (without IPv6)

```yaml
version: "3.0"
services:
  wg-access-server:
    image: ghcr.io/freifunkmuc/wg-access-server:latest
    container_name: wg-access-server
    cap_add:
      - NET_ADMIN
    volumes:
      - "wg-access-server-data:/data"
    environment:
      - "WG_ADMIN_PASSWORD=${WG_ADMIN_PASSWORD:?\n\nplease set the WG_ADMIN_PASSWORD environment variable:\n    export WG_ADMIN_PASSWORD=example\n}"
      - "WG_WIREGUARD_PRIVATE_KEY=${WG_WIREGUARD_PRIVATE_KEY:?\n\nplease set the WG_WIREGUARD_PRIVATE_KEY environment variable:\n    export WG_WIREGUARD_PRIVATE_KEY=$(wg genkey)\n}"
      - "WG_VPN_CIDRV6=0" # to disable IPv6
    ports:
      - "8000:8000/tcp"
      - "51820:51820/udp"
    devices:
      - "/dev/net/tun:/dev/net/tun"

volumes:
  wg-access-server-data:
    driver: local
```
