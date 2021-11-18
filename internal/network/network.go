package network

import (
	"net"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
)

// ServerVPNIPs returns two net.IPNet objects (for IPv4 + IPv6)
// with the IP attribute set to the server's IP addresses
// in these subnets, i.e. the first usable address
// The  return values are nil if the corresponding input is an empty string
func ServerVPNIPs(cidr, cidr6 string) (ipv4, ipv6 *net.IPNet, err error) {
	if cidr != "" {
		vpnip, vpnsubnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, nil, err
		}
		vpnsubnet.IP = nextIP(vpnip.Mask(vpnsubnet.Mask))
		ipv4 = vpnsubnet
	}
	if cidr6 != "" {
		vpnip, vpnsubnet, err := net.ParseCIDR(cidr6)
		if err != nil {
			return nil, nil, err
		}
		vpnsubnet.IP = nextIP(vpnip.Mask(vpnsubnet.Mask))
		ipv6 = vpnsubnet
	}
	return ipv4, ipv6, nil
}

// StringJoinIPNets joins the string representations of a and b using a comma
func StringJoinIPNets(a, b *net.IPNet) string {
	if a != nil && b != nil {
		return strings.Join([]string{a.String(), b.String()}, ", ")
	} else if a != nil {
		return a.String()
	} else if b != nil {
		return b.String()
	}
	return ""
}

// StringJoinIPs joins the string representations of the IPs of a and b using a comma
func StringJoinIPs(a, b *net.IPNet) string {
	if a != nil && b != nil {
		return strings.Join([]string{a.IP.String(), b.IP.String()}, ", ")
	} else if a != nil {
		return a.IP.String()
	} else if b != nil {
		return b.IP.String()
	}
	return ""
}

func ConfigureForwarding(gatewayIface string, cidr string, cidrv6 string, nat66 bool, allowedIPs []string) error {
	// Networking configuration (iptables) configuration
	// to ensure that traffic from clients of the wireguard interface
	// is sent to the provided network interface
	allowedIPv4s := make([]string, 0, len(allowedIPs)/2)
	allowedIPv6s := make([]string, 0, len(allowedIPs)/2)

	for _, allowedCIDR := range allowedIPs {
		parsedAddress, parsedNetwork, err := net.ParseCIDR(allowedCIDR)
		if err != nil {
			return errors.Wrap(err, "invalid cidr in AllowedIPs")
		}
		if as4 := parsedAddress.To4(); as4 != nil {
			// Handle IPv4-mapped IPv6 addresses, if they go into ip6tables they don't get hit
			// and go-iptables can't convert them (whereas commandline iptables can).
			parsedNetwork.IP = as4
			allowedIPv4s = append(allowedIPv4s, parsedNetwork.String())
		} else {
			allowedIPv6s = append(allowedIPv6s, parsedNetwork.String())
		}
	}

	if cidr != ""{
		if err := configureForwardingv4(gatewayIface, cidr, allowedIPv4s); err != nil {
			return err
		}
	}
	if cidrv6 != "" {
		if err := configureForwardingv6(gatewayIface, cidrv6, nat66, allowedIPv6s); err != nil {
			return err
		}
	}
	return nil
}

func configureForwardingv4(gatewayIface string, cidr string, allowedIPs []string) error {
	ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		return errors.Wrap(err, "failed to init iptables")
	}

	// Cleanup our chains first so that we don't leak
	// iptable rules when the network configuration changes.
	err = ipt.ClearChain("filter", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return errors.Wrap(err, "failed to clear filter chain")
	}
	err = ipt.ClearChain("nat", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return errors.Wrap(err, "failed to clear nat chain")
	}

	// Create our own chain for forwarding rules
	err = ipt.NewChain("filter", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return errors.Wrap(err, "failed to create filter chain")
	}
	err = ipt.AppendUnique("filter", "FORWARD", "-j", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return errors.Wrap(err, "failed to append FORWARD rule to filter chain")
	}

	// Create our own chain for postrouting rules
	err = ipt.NewChain("nat", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return errors.Wrap(err, "failed to create nat chain")
	}
	err = ipt.AppendUnique("nat", "POSTROUTING", "-j", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return errors.Wrap(err, "failed to append POSTROUTING rule to nat chain")
	}

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

func configureForwardingv6(gatewayIface string, cidrv6 string, nat66 bool, allowedIPs []string) error {
	ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
	if err != nil {
		return errors.Wrap(err, "failed to init ip6tables")
	}

	err = ipt.ClearChain("filter", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return errors.Wrap(err, "failed to clear filter chain")
	}
	err = ipt.ClearChain("nat", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return errors.Wrap(err, "failed to clear nat chain")
	}

	err = ipt.NewChain("filter", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return errors.Wrap(err, "failed to create filter chain")
	}
	err = ipt.AppendUnique("filter", "FORWARD", "-j", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return errors.Wrap(err, "failed to append FORWARD rule to filter chain")
	}

	err = ipt.NewChain("nat", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return errors.Wrap(err, "failed to create nat chain")
	}
	err = ipt.AppendUnique("nat", "POSTROUTING", "-j", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return errors.Wrap(err, "failed to append POSTROUTING rule to nat chain")
	}

	// Accept client traffic for given allowed ips
	for _, allowedCIDR := range allowedIPs {
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidrv6, "-d", allowedCIDR, "-j", "ACCEPT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}

	if gatewayIface != "" {
		if nat66 {
			if err := ipt.AppendUnique("nat", "WG_ACCESS_SERVER_POSTROUTING", "-s", cidrv6, "-o", gatewayIface, "-j", "MASQUERADE"); err != nil {
				return errors.Wrap(err, "failed to set ip tables rule")
			}
		}
	}

	if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", cidrv6, "-j", "REJECT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	return nil
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
