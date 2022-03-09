## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install my-release --repo https://freie-netze.org/wg-access-server wg-access-server
```

The command deploys wg-access-server on the Kubernetes cluster in the default configuration. The configuration section lists the parameters that can be configured during installation.

A wireguard private key needs to be set in order for the pod to start successfully. Use `wg genkey` and append `--set wireguard.config.privateKey="<wg-private-key>"` to the command above.

Per default persistence is disable and devices will not persist. To enable persistence, set `persistence.enabled`.

Because IPv6 on Kubernetes is disabled by default in most clusters and can't be enabled on a per-pod basis, the default `values.yaml` disables it for the VPN as well. If you have a cluster with working IPv6, set `config: {}` in your `values.yaml` or specify a custom VPN-internal prefix under `config.vpn.cidrv6`.

If no admin password is set, the Chart generates a random one. You can retrieve it using `kubectl get secret ...` as prompted by helm after installing the Chart.

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Example values.yaml

```yaml
config:
  externalHost: "<loadbalancer-ip>"

# wg access server is an http server without TLS. Exposing it via a loadbalancer is NOT secure!
# Uncomment the following section only if you are running on private network or simple testing.
# A much better option would be TLS terminating ingress controller or reverse-proxy.
# web:
#   service:
#     type: "LoadBalancer"
#     loadBalancerIP: "<loadbalancer-ip>"

wireguard:
  config:
    privateKey: "<wireguard-private-key>"
  service:
    type: "LoadBalancer"
    loadBalancerIP: "<loadbalancer-ip>"

persistence:
  enabled: true

ingress:
  enabled: true
  hosts: ["vpn.example.com"]
  tls:
    - hosts: ["vpn.example.com"]
      secretName: "tls-wg-access-server"
```



## All Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| config | object | `{}` | inline wg-access-server config (config.yaml) |
| web.config.adminUsername | string | `"admin"` |  |
| web.config.adminPassword | string | `""` | If omitted a random password will be generated and stored in the secret |
| web.service.annotations | object | `{}` |  |
| web.service.externalTrafficPolicy | string | `""` |  |
| web.service.type | string | `"ClusterIP"` |  |
| web.service.loadBalancerIP | string | `""` |  |
| wireguard.config.privateKey | string | `""` | REQUIRED - A wireguard private key. You can generate one using `$ wg genkey` |
| wireguard.service.annotations | object | `{}` |  |
| wireguard.service.type | string | `"ClusterIP"` |  |
| wireguard.service.sessionAffinity | string | `"ClientIP"` |  |
| wireguard.service.externalTrafficPolicy | string | `""` |  |
| wireguard.service.ipFamilyPolicy | string | `"SingleStack"` |  |
| wireguard.service.loadBalancerIP | string | `""` |  |
| wireguard.service.port | int | `51820` |  |
| wireguard.service.nodePort | int | `""` | Use available port from range 30000-32768 |
| persistence.enabled | bool | `false` |  |
| persistence.existingClaim | string | `""` | Use existing PVC claim for persistence instead |
| persistence.annotations | object | `{}` |  |
| persistence.accessModes[0] | string | `"ReadWriteOnce"` |  |
| persistence.storageClass | string | `""` |  |
| persistence.size | string | `"100Mi"` |  |
| ingress.enabled | bool | `false` |  |
| ingress.annotations | object | `{}` |  |
| ingress.ingressClassName | string | `""` |  |
| ingress.hosts | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| imagePullSecrets | list | `[]` |  |
| image.repository | string | `"ghcr.io/freifunkmuc/wg-access-server"` |  |
| image.tag | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| replicas | int | `1` |  |
| strategy.type | string | `""` | `Recreate` if `persistence.enabled` true or `RollingUpdate` if false |
| resources | object | `{}` | pod cpu/memory resource requests and limits |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
