package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"

	"github.com/vishvananda/netlink"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/alecthomas/kingpin.v2"
)

type AppConfig struct {
	Web struct {
		// ExternalAddress is the address that
		// clients should use to connect to this
		// server. It will be used in generated
		// VPN client connection configuration
		// ExternalAddress should not include any
		// port infomation.
		// The WireGuard port will be appended.
		ExternalAddress string
		// Port that the web server should listen on
		Port int
	}
	Storage struct {
		// Directory that VPN devices (WireGuard peers)
		// should be saved under.
		// If this value is empty then an InMemory storage
		// backend will be used (not recommended).
		Directory string
	}
	WireGuard struct {
		// UserspaceImplementation is a command (program on $PATH)
		// that implements the WireGuard protocol in userspace.
		// In our Docker image we make use of `boringtun` so that
		// users aren't required to setup kernel modules
		UserspaceImplementation string
		// The network interface name of the WireGuard
		// network device
		InterfaceName string
		// The WireGuard PrivateKey
		// If this value is lost then any existing
		// clients (WireGuard peers) will no longer
		// be able to connect.
		// Clients will either have to manually update
		// their connection configuration or setup
		// their VPN again using the web ui (easier for most people)
		PrivateKey string
		// The WireGuard ListenPort
		Port int
	}
	VPN struct {
		// SubnetCIDR configures a network address space
		// that client (WireGuard peers) will be allocated
		// an IP address from
		SubnetCIDR string
		// GatewayInterface will be used in iptable forwarding
		// rules that send VPN traffic from clients to this interface
		// Most use-cases will want this interface to have access
		// to the outside internet
		GatewayInterface netlink.Link
	}
}

var (
	app                              = kingpin.New("TODO: name-this-program", "An all-in-one WireGuard VPN solution")
	logLevel                         = app.Flag("loglevel", "Enable debug mode").Default("info").OverrideDefaultFromEnvar("LOG_LEVEL").String()
	webPort                          = app.Flag("web-port", "The web server port").Default("8000").OverrideDefaultFromEnvar("WEB_PORT").Int()
	webExternalAddress               = app.Flag("web-external-address", "The external address that the service is accessible from excluding any scheme or port (e.g. vpn.example.com). Defaults to the IP address of your default interface").OverrideDefaultFromEnvar("WEB_EXTERNAL_ADDRESS").String()
	storageDirectory                 = app.Flag("storage-directory", "The directory where vpn devices (i.e. peers) will be stored").OverrideDefaultFromEnvar("STORAGE_DIRECTORY").String()
	wireGuardUserspaceImplementation = app.Flag("wireguard-userspace-implementation", "The a userspace implementation of wireguard e.g. wireguard-go or boringtun").OverrideDefaultFromEnvar("WIREGUARD_USERSPACE_IMPLEMENTATION").String()
	wireGuardInterfaceName           = app.Flag("wireguard-interface-name", "The name of the WireGuard interface").Default("wg0").OverrideDefaultFromEnvar("WIREGUARD_INTERFACE_NAME").String()
	wireguardPort                    = app.Flag("wireguard-port", "The WireGuard ListenPort").Default("51820").OverrideDefaultFromEnvar("WIREGUARD_PORT").Int()
	wireguardPrivateKey              = app.Flag("wireguard-private-key", "The WireGuard private key").OverrideDefaultFromEnvar("WIREGUARD_PRIVATE_KEY").String()
	vpnGatewayInterfaceName          = app.Flag("vpn-gateway-interface-name", "The name of the network interface you want VPN client traffic to foward to").OverrideDefaultFromEnvar("VPN_GATEWAY_INTERFACE_NAME").String()
	vpnSubnetCIDR                    = app.Flag("vpn-subnet-cidr", "The subnet CIDR that clients should be networked within").Default("10.44.0.1/24").OverrideDefaultFromEnvar("VPN_SUBNET_CIDR").String()
)

func Read() *AppConfig {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	config := AppConfig{}

	config.Web.Port = *webPort
	config.Web.ExternalAddress = *webExternalAddress
	config.Storage.Directory = *storageDirectory
	config.WireGuard.UserspaceImplementation = *wireGuardUserspaceImplementation
	config.WireGuard.InterfaceName = *wireGuardInterfaceName
	config.WireGuard.Port = *wireguardPort
	config.WireGuard.PrivateKey = *wireguardPrivateKey
	config.VPN.SubnetCIDR = *vpnSubnetCIDR
	config.VPN.GatewayInterface = findGatewayLink(*vpnGatewayInterfaceName)
	if config.Web.ExternalAddress == "" && config.VPN.GatewayInterface != nil {
		if ip, err := linkIPAddr(config.VPN.GatewayInterface); err == nil {
			config.Web.ExternalAddress = ip.String()
			logrus.Infof("no external address was configured - using %s from the gateway interface", config.Web.ExternalAddress)
		}
	}

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "invalid log level - should be one of fatal, error, warn, info, debug, trace"))
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		},
	})

	if config.WireGuard.PrivateKey == "" {
		logrus.Warn("no private key has been configured! using an in-memory private key that will be lost when the process exits!")
		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to generate a server private key"))
		}
		config.WireGuard.PrivateKey = key.String()
	}

	if config.Storage.Directory == "" {
		logrus.Warn("storage directory not configured - using in-memory storage backend! wireguard devices will be lost when the process exits!")
	}

	return &config
}

func defaultInterface() string {
	links, err := netlink.LinkList()
	if err != nil {
		logrus.Warn(errors.Wrap(err, "failed to list network interfaces"))
		return ""
	}
	for _, link := range links {
		routes, err := netlink.RouteList(link, 4)
		if err != nil {
			logrus.Warn(errors.Wrapf(err, "failed to list routes for interface %s", link.Attrs().Name))
			return ""
		}
		for _, route := range routes {
			if route.Dst == nil {
				return link.Attrs().Name
			}
		}
	}
	return ""
}

func findGatewayLink(name string) netlink.Link {
	if name == "" {
		if name = defaultInterface(); name == "" {
			logrus.Warn("a gateway interface name was not configured - vpn forwarding rules will not be applied!")
			return nil
		} else {
			logrus.Infof("no gateway interface name was configured - using the system's default route's interface %s", name)
		}
	}
	if name == "" {
	}
	link, err := netlink.LinkByName(name)
	if err != nil {
		logrus.Warn(errors.Wrapf(err, "the gateway interface '%s' could not be found - vpn forwarding rules will not be applied!", name))
		return nil
	}
	return link
}

func linkIPAddr(link netlink.Link) (net.IP, error) {
	routes, err := netlink.RouteList(link, 4)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list routes for interface %s", link.Attrs().Name)
	}
	for _, route := range routes {
		if route.Src != nil {
			return route.Src, nil
		}
	}
	return nil, fmt.Errorf("no source IP found for interface %s", link.Attrs().Name)
}
