package config

import (
	"github.com/place1/wg-access-server/pkg/authnz/authconfig"
)

type AppConfig struct {
	// Set the log level.
	// Defaults to "info" (fatal, error, warn, info, debug, trace)
	LogLevel string `yaml:"loglevel"`
	// Set the superadmin username
	// Defaults to "admin"
	AdminUsername string `yaml:"adminUsername"`
	// Set the superadmin password (required)
	AdminPassword string `yaml:"adminPassword"`
	// Port sets the port that the web UI will listen on.
	// Defaults to 8000
	Port int `yaml:"port"`
	// ExternalHost is the address that clients
	// use to connect to the wireguard interface
	// By default, this will be empty and the web ui
	// will use the current page's origin.
	ExternalHost string `yaml:"externalHost"`
	// The storage backend where device configuration will
	// be persisted.
	// Supports memory:// postgresql:// mysql:// sqlite3://
	// Defaults to memory://
	Storage string `yaml:"storage"`
	// DisableMetadata allows you to turn off collection of device
	// metadata including last handshake time & rx/tx bytes
	DisableMetadata bool `yaml:"disableMetadata"`
	// Configure WireGuard related settings
	WireGuard struct {
		// Set this to false to disable the embedded wireguard
		// server. This is useful for development environments
		// on mac and windows where we don't currently support
		// the OS's network stack.
		Enabled bool `yaml:"enabled"`
		// The network interface name of the WireGuard
		// network device.
		// Defaults to wg0
		Interface string `yaml:"interface"`
		// The WireGuard PrivateKey
		// If this value is lost then any existing
		// clients (WireGuard peers) will no longer
		// be able to connect.
		// Clients will either have to manually update
		// their connection configuration or setup
		// their VPN again using the web ui (easier for most people)
		PrivateKey string `yaml:"privateKey"`
		// The WireGuard ListenPort
		// Defaults to 51820
		Port int `yaml:"port"`
	} `yaml:"wireguard"`
	// Configure VPN related settings (networking)
	VPN struct {
		// CIDR configures a network address space
		// that client (WireGuard peers) will be allocated
		// an IP address from
		// defaults to 10.44.0.0/24
		CIDR string `yaml:"cidr"`
		// CIDRv6 configures an IPv6 network address space
		// that client (WireGuard peers) will be allocated
		// an IP address from
		// defaults to none
		CIDRv6 string `yaml:"cidrv6"`
		// NAT66 configures whether IPv6 traffic leaving
		// through the GatewayInterface should be
		// masqueraded like IPv4 traffic
		// defaults to true
		NAT66 bool `yaml:"nat66"`
		// GatewayInterface will be used in iptable forwarding
		// rules that send VPN traffic from clients to this interface
		// Most use-cases will want this interface to have access
		// to the outside internet
		GatewayInterface string `yaml:"gatewayInterface"`
		// The "AllowedIPs" for VPN clients.
		// This value will be included in client config
		// files and in server-side iptable rules
		// to enforce network access.
		// defaults to ["0.0.0.0/0", "::/0"]
		AllowedIPs []string `yaml:"allowedIPs"`
	} `yaml:"vpn"`
	// Configure the embeded DNS server
	DNS struct {
		// Enabled allows you to turn on/off
		// the VPN DNS proxy feature.
		// DNS Proxying is enabled by default.
		Enabled bool `yaml:"enabled"`
		// Upstream configures the addresses of upstream
		// DNS servers to which client DNS requests will be sent to.
		// Defaults the host's upstream DNS servers (via resolveconf)
		// or 1.1.1.1 if resolveconf cannot be used.
		// NOTE: currently wg-access-server will only use the first upstream.
		Upstream []string `yaml:"upstream"`
	} `yaml:"dns"`
	// Auth configures optional authentication backends
	// to controll access to the web ui.
	// Devices will be managed on a per-user basis if any
	// auth backends are configured.
	// If no authentication backends are configured then
	// the server will not require any authentication.
	Auth authconfig.AuthConfig `yaml:"auth"`
}
