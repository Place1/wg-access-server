package network

import (
	"fmt"
	"net"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
)

func ServerVPNIP(cidr string) *net.IPNet {
	vpnip, vpnsubnet := MustParseCIDR(cidr)
	vpnsubnet.IP = nextIP(vpnip.Mask(vpnsubnet.Mask))
	return vpnsubnet
}

func ConfigureForwarding(wgIface string, gatewayIface string, cidr string, allowedIPs []string) error {
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

	// Accept client traffic for given allowed ips
	for _, allowedCIDR := range allowedIPs {
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-d", allowedCIDR, "-j", "ACCEPT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}

	if gatewayIface != "" {
		if err := ipt.AppendUnique("nat", "WG_ACCESS_SERVER_POSTROUTING", "-s", cidr, "-o", gatewayIface, "-j", "MASQUERADE"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}

	if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidr, "-j", "REJECT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
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
