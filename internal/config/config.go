package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"

	"github.com/place1/wireguard-access-server/internal/auth/authconfig"
	"github.com/vishvananda/netlink"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/alecthomas/kingpin.v2"
)

type AppConfig struct {
	LogLevel string `yaml:"loglevel"`
	Web      struct {
		// ExternalAddress is that users access the web ui
		// using. This value is required for using auth backends
		// This value should include the scheme.
		// The port should be included if non-standard.
		// e.g. http://192.168.0.2:8000
		// or https://myvpn.example.com
		ExternalAddress string `yaml:"externalAddress"`
		// Port that the web server should listen on
		Port int `yaml:"port"`
	} `yaml:"web"`
	Storage struct {
		// Directory that VPN devices (WireGuard peers)
		// should be saved under.
		// If this value is empty then an InMemory storage
		// backend will be used (not recommended).
		Directory string `yaml:"directory"`
	} `yaml:"storage"`
	WireGuard struct {
		// The network interface name of the WireGuard
		// network device.
		// Defaults to wg0
		InterfaceName string `yaml:"interfaceName"`
		// The WireGuard PrivateKey
		// If this value is lost then any existing
		// clients (WireGuard peers) will no longer
		// be able to connect.
		// Clients will either have to manually update
		// their connection configuration or setup
		// their VPN again using the web ui (easier for most people)
		PrivateKey string `yaml:"privateKey"`
		// ExternalAddress is the address that users
		// use to connect to the wireguard interface
		// By default, this will use the Web.ExternalAddress
		// domain with the WireGuard.Port
		ExternalAddress *string `yaml:"externalAddress"`
		// The WireGuard ListenPort
		Port int `yaml:"port"`
	} `yaml:"wireguard"`
	VPN struct {
		// CIDR configures a network address space
		// that client (WireGuard peers) will be allocated
		// an IP address from
		CIDR string `yaml:"cidr"`
		// GatewayInterface will be used in iptable forwarding
		// rules that send VPN traffic from clients to this interface
		// Most use-cases will want this interface to have access
		// to the outside internet
		GatewayInterface string `yaml:"gatewayInterface"`
	}
	DNS struct {
		// TODO: docs
		Upstream []string `yaml:"upstream"`
	} `yaml:"dns"`
	Auth *authconfig.AuthConfig `yaml:"auth"`
}

var (
	app        = kingpin.New("wg-access-server", "An all-in-one WireGuard Access Server & VPN solution")
	configPath = app.Flag("config", "Path to a config file").OverrideDefaultFromEnvar("CONFIG").String()
)

func Read() *AppConfig {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	config := AppConfig{}
	config.LogLevel = "info"
	config.Web.Port = 8000
	config.WireGuard.InterfaceName = "wg0"
	config.WireGuard.Port = 51820
	config.VPN.CIDR = "10.44.0.0/24"

	if *configPath != "" {
		b, err := ioutil.ReadFile(*configPath)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to read the configuration file"))
		}
		if err := yaml.Unmarshal(b, &config); err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to bind configuration file"))
		}
	}

	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		config.LogLevel = v
	}

	if v, ok := os.LookupEnv("STORAGE_DIRECTORY"); ok {
		config.Storage.Directory = v
	}

	if v, ok := os.LookupEnv("WIREGUARD_PRIVATE_KEY"); ok {
		config.WireGuard.PrivateKey = v
	}

	level, err := logrus.ParseLevel(config.LogLevel)
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

	if config.VPN.GatewayInterface == "" {
		iface, err := defaultInterface()
		if err != nil {
			logrus.Warn(errors.Wrap(err, "failed to set default value for VPN.GatewayInterface"))
		} else {
			config.VPN.GatewayInterface = iface
		}
	}

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
	} else {
		config.Storage.Directory, err = filepath.Abs(config.Storage.Directory)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to get absolute path to storage directory"))
		}
	}

	return &config
}

func (config *AppConfig) IsAuthEnabled() bool {
	return config.Auth != nil
}

func defaultInterface() (string, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return "", errors.Wrap(err, "failed to list network interfaces")
	}
	for _, link := range links {
		routes, err := netlink.RouteList(link, 4)
		if err != nil {
			return "", errors.Wrapf(err, "failed to list routes for interface %s", link.Attrs().Name)
		}
		for _, route := range routes {
			if route.Dst == nil {
				return link.Attrs().Name, nil
			}
		}
	}
	return "", errors.New("could not determine the default network interface name")
}

func linkIPAddr(name string) (net.IP, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find network interface %s", name)
	}
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
