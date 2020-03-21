package network

import (
	"fmt"
	"net"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func ServerVPNIP(cidr string) *net.IPNet {
	vpnip, vpnsubnet := MustParseCIDR(cidr)
	vpnsubnet.IP = nextIP(vpnip.Mask(vpnsubnet.Mask))
	return vpnsubnet
}

func ConfigureRouting(wgIface string, cidr string) error {
	// Networking configuration (ip links and route tables)
	// to ensure that network traffic in the VPN subnet
	// moves through the wireguard interface
	link, err := netlink.LinkByName(wgIface)
	if err != nil {
		return errors.Wrap(err, "failed to find wireguard interface")
	}
	vpnip := ServerVPNIP(cidr)
	logrus.Infof("server VPN subnet IP is %s", vpnip.String())
	addr, err := netlink.ParseAddr(vpnip.String())
	if err != nil {
		return errors.Wrap(err, "failed to parse subnet address")
	}
	if err := netlink.AddrAdd(link, addr); err != nil {
		logrus.Warn(errors.Wrap(err, "failed to add subnet to wireguard interface"))
	}
	if err := netlink.LinkSetUp(link); err != nil {
		logrus.Warn(errors.Wrap(err, "failed to bring wireguard interface up"))
	}
	return nil
}

type NetworkRules struct {
	// AllowVPNLAN enables routing between VPN clients
	// i.e. allows the VPN to work like a LAN.
	AllowVPNLAN bool
	// AllowServerLAN enables routing to private IPv4
	// address ranges. Enabling this will allow VPN clients
	// to access networks on the server's LAN.
	AllowServerLAN bool
	// AllowInternet enables routing of all traffic
	// to the public internet.
	AllowInternet bool
	// AllowedNetworks allows you to whitelist a partcular
	// network CIDR. This is useful if you want to block
	// access to the Server's LAN but allow access to a few
	// specific IPs or a small range.
	// e.g. "192.0.2.0/24" or "192.0.2.10/32".
	AllowedNetworks []string
}

func ConfigureForwarding(wgIface string, gatewayIface string, cidr string, rules NetworkRules) error {
	// Networking configuration (iptables) configuration
	// to ensure that traffic from clients the wireguard interface
	// is sent to the provided network interface
	ipt, err := iptables.New()
	if err != nil {
		return errors.Wrap(err, "failed to init iptables")
	}

	// Cleanup our chains first so that we don't leak
	// iptable rules when the network configuration changes.
	ipt.ClearChain("filter", "WG_ACCESS_SERVER_FORWARD")
	ipt.ClearChain("nat", "WG_ACCESS_SERVER_POSTROUTING")

	// Create our own chain for forwarding rules
	ipt.NewChain("filter", "WG_ACCESS_SERVER_FORWARD")
	ipt.AppendUnique("filter", "FORWARD", "-j", "WG_ACCESS_SERVER_FORWARD")

	// Create our own chain for postrouting rules
	ipt.NewChain("nat", "WG_ACCESS_SERVER_POSTROUTING")
	ipt.AppendUnique("nat", "POSTROUTING", "-j", "WG_ACCESS_SERVER_POSTROUTING")

	if err := ConfigureRouting(wgIface, cidr); err != nil {
		logrus.Error(errors.Wrap(err, "failed to configure interface"))
	}

	// https://simple.wikipedia.org/wiki/Private_network
	privateCIDRs := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

	// White listed networks
	if len(rules.AllowedNetworks) != 0 {
		for _, subnet := range rules.AllowedNetworks {
			if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-d", subnet, "-j", "ACCEPT"); err != nil {
				return errors.Wrap(err, "failed to set ip tables rule")
			}
		}
	}

	// VPN LAN
	if rules.AllowVPNLAN {
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-d", cidr, "-j", "ACCEPT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	} else {
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-d", cidr, "-j", "REJECT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}

	// Server LAN
	for _, privateCIDR := range privateCIDRs {
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-d", privateCIDR, "-j", boolToRule(rules.AllowServerLAN)); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}

	// Internet
	if rules.AllowInternet && gatewayIface != "" {
		// TODO: do we actually need to specify a gateway interface?
		// I suppose i neet to refresh my knowledge of nat.
		// if you're reading this please open a Github issue and help teach me nat and iptables :P
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-i", gatewayIface, "-o", wgIface, "-j", "ACCEPT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-i", wgIface, "-o", gatewayIface, "-j", "ACCEPT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
		if err := ipt.AppendUnique("nat", "WG_ACCESS_SERVER_POSTROUTING", "-s", cidr, "-o", gatewayIface, "-j", "MASQUERADE"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	} else {
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-j", "REJECT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}

	return nil
}

func MustParseCIDR(cidr string) (net.IP, *net.IPNet) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return ip, ipnet
}

func MustParseIP(ip string) net.IP {
	netip, _ := MustParseCIDR(fmt.Sprintf("%s/32", ip))
	return netip
}

func nextIP(ip net.IP) net.IP {
	next := make([]byte, len(ip))
	copy(next, ip)
	for j := len(next) - 1; j >= 0; j-- {
		next[j]++
		if next[j] > 0 {
			break
		}
	}
	return next
}

func boolToRule(accept bool) string {
	if accept {
		return "ACCEPT"
	}
	return "REJECT"
}
