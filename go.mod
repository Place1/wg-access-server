module github.com/place1/wireguard-access-server

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v39.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.5 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/alexedwards/scs/v2 v2.2.0
	github.com/beevik/etree v1.1.0 // indirect
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/coreos/go-iptables v0.4.3
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dexidp/dex v2.13.0+incompatible
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/golang/protobuf v1.3.3
	github.com/gorilla/handlers v1.4.2 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/sessions v1.2.0
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/improbable-eng/grpc-web v0.12.0
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/miekg/dns v1.1.27
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/place1/wg-embed v0.0.0
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/rs/cors v1.7.0 // indirect
	github.com/russellhaering/goxmldsig v0.0.0-20180430223755-7acd5e4a6ef7 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/tg123/go-htpasswd v1.0.0
	github.com/vishvananda/netlink v1.0.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20191205174707-786493d6718c
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/genproto v0.0.0-20200210034751-acff78025515 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/square/go-jose.v2 v2.4.0 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/place1/wg-embed => ../wg-embed
