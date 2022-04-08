# Docker

Load the `ip_tables`, `ip6_tables` and `wireguard` kernel modules on the host.

```bash
modprobe ip_tables && modprobe ip6_tables && modprobe wireguard
# Load modules on boot
echo ip_tables >> /etc/modules
echo ip6_tables >> /etc/modules
echo wireguard >> /etc/modules
```

```bash
docker run \
  -it \
  --rm \
  --cap-add NET_ADMIN \
  --device /dev/net/tun:/dev/net/tun \
  --sysctl net.ipv6.conf.all.disable_ipv6=0 \
  --sysctl net.ipv6.conf.all.forwarding=1 \
  -v wg-access-server-data:/data \
  -e "WG_ADMIN_PASSWORD=$WG_ADMIN_PASSWORD" \
  -e "WG_WIREGUARD_PRIVATE_KEY=$WG_WIREGUARD_PRIVATE_KEY" \
  -p 8000:8000/tcp \
  -p 51820:51820/udp \
  ghcr.io/freifunkmuc/wg-access-server:latest
```

## Modules

If you are unable to load the `iptables` kernel modules, you can add the `SYS_MODULE` capability instead: `--cap-add SYS_MODULE`. You must also add the following mount: `-v /lib/modules:/lib/modules:ro`.

This is not recommended as it essentially gives the container root privileges over the host system and an attacker could easily break out of the container.

The WireGuard module should be loaded automatically, even without `SYS_MODULE` capability or `/lib/modules` mount.
If it still fails to load, the server automatically falls back to the userspace implementation. 

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
