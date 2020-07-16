## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install my-release --repo https://place1.github.io/wg-access-server wg-access-server
```

The command deploys wg-access-server on the Kubernetes cluster in the default configuration. The configuration section lists the parameters that can be configured during installation.

By default an in-memory wireguard private key will be generated and devices will not persist
between pod restarts.

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Example values.yaml

```yaml
config:
  wireguard:
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
| web.service.type | string | `"ClusterIP"` |  |
| wireguard.config.privateKey | string | "" | A wireguard private key. You can generate one using `$ wg genkey` |
| wireguard.service.type | string | `"ClusterIP"` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts | string | `nil` |  |
| ingress.tls | list | `[]` |  |
| ingress.annotations | object | `{}` |  |
| persistence.enabled | bool | `false` |  |
| persistence.existingClaim | string | `""` | Use existing PVC claim for persistence instead |
| persistence.size | string | `"100Mi"` |  |
| persistence.subPath | string | `""` |  |
| persistence.annotations | object | `{}` |  |
| persistence.accessModes[0] | string | `"ReadWriteOnce"` |  |
| strategy.type | string | `"Recreate"` |  |
| resources | object | `{}` | pod cpu/memory resource requests and limits |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| affinity | object | `{}` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"place1/wg-access-server"` |  |
| imagePullSecrets | list | `[]` |  |
