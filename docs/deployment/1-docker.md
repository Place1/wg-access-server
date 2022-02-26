# Docker

## TL;DR;

Here's a one-liner to run wg-access-server:

```bash
docker run \
  -it \
  --rm \
  --cap-add NET_ADMIN \
  --cap-add SYS_MODULE \
  --device /dev/net/tun:/dev/net/tun \
  --sysctl net.ipv6.conf.all.disable_ipv6=0 \
  --sysctl net.ipv6.conf.all.forwarding=1 \
  -v wg-access-server-data:/data \
  -v /lib/modules:/lib/modules:ro \
  -e "WG_ADMIN_PASSWORD=$WG_ADMIN_PASSWORD" \
  -e "WG_WIREGUARD_PRIVATE_KEY=$WG_WIREGUARD_PRIVATE_KEY" \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  ghcr.io/freifunkmuc/wg-access-server:latest
```

## Modules

If you load the kernel modules `ip_tables` and `ip6_tables` on the host,
you can drop the `SYS_MODULE` capability and remove the `/lib/modules` mount:
```bash
modprobe ip_tables && modprobe ip6_tables
# Load modules on boot
echo ip_tables >> /etc/modules
echo ip6_tables >> /etc/modules
```
This is highly recommended, as a container with CAP_SYS_MODULE essentially has root privileges
over the host system and attacker could easily break out of the container.

## IPv4-only (without IPv6)

If you don't want IPv6 inside the VPN network, set `WG_VPN_CIDRV6=0`.
In this case you can also get rid of the sysctls:

```bash
docker run \
  -it \
  --rm \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  -v wg-access-server-data:/data \
  -e "WG_ADMIN_PASSWORD=$WG_ADMIN_PASSWORD" \
  -e "WG_WIREGUARD_PRIVATE_KEY=$WG_WIREGUARD_PRIVATE_KEY" \
  -e "WG_VPN_CIDRV6=0"
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  ghcr.io/freifunkmuc/wg-access-server:latest
```

Likewise you can disable IPv4 by setting `WG_VPN_CIDR=0`.
