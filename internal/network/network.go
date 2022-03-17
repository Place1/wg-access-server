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

// SplitAddresses splits multiple comma-separated addresses into a slice of address strings
func SplitAddresses(addresses string) []string {
	split := strings.Split(addresses, ",")
	for i, addr := range split {
		split[i] = strings.TrimSpace(addr)
	}
	return split
}

// ForwardingOptions contains all options used for configuring the firewall rules
type ForwardingOptions struct {
	GatewayIface    string
	CIDR, CIDRv6    string
	NAT44, NAT66    bool
	ClientIsolation bool
	AllowedIPs      []string
	allowedIPv4s    []string
	allowedIPv6s    []string
}

func ConfigureForwarding(options ForwardingOptions) error {
	// Networking configuration (iptables) configuration
	// to ensure that traffic from clients of the wireguard interface
	// is sent to the provided network interface
	allowedIPv4s := make([]string, 0, len(options.AllowedIPs)/2)
	allowedIPv6s := make([]string, 0, len(options.AllowedIPs)/2)

	for _, allowedCIDR := range options.AllowedIPs {
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
	options.allowedIPv4s = allowedIPv4s
	options.allowedIPv6s = allowedIPv6s

	if options.CIDR != "" {
		if err := configureForwardingv4(options); err != nil {
			return err
		}
	}
	if options.CIDRv6 != "" {
		if err := configureForwardingv6(options); err != nil {
			return err
		}
	}
	return nil
}

func configureForwardingv4(options ForwardingOptions) error {
	ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		return errors.Wrap(err, "failed to init iptables")
	}

	// Cleanup our chains first so that we don't leak
	// iptable rules when the network configuration changes.
	err = clearOrCreateChain(ipt, "filter", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return err
	}

	err = clearOrCreateChain(ipt, "nat", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return err
	}

	err = ipt.AppendUnique("filter", "FORWARD", "-j", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return errors.Wrap(err, "failed to append FORWARD rule to filter chain")
	}

	err = ipt.AppendUnique("nat", "POSTROUTING", "-j", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return errors.Wrap(err, "failed to append POSTROUTING rule to nat chain")
	}

	if options.ClientIsolation {
		// Reject inter-device traffic
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", options.CIDR, "-d", options.CIDR, "-j", "REJECT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}
	// Accept client traffic for given allowed ips
	for _, allowedCIDR := range options.allowedIPv4s {
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", options.CIDR, "-d", allowedCIDR, "-j", "ACCEPT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}
	// And reject everything else
	if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", options.CIDR, "-j", "REJECT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}

	if options.GatewayIface != "" {
		if options.NAT44 {
			if err := ipt.AppendUnique("nat", "WG_ACCESS_SERVER_POSTROUTING", "-s", options.CIDR, "-o", options.GatewayIface, "-j", "MASQUERADE"); err != nil {
				return errors.Wrap(err, "failed to set ip tables rule")
			}
		}
	}
	return nil
}

func configureForwardingv6(options ForwardingOptions) error {
	ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
	if err != nil {
		return errors.Wrap(err, "failed to init ip6tables")
	}

	err = clearOrCreateChain(ipt, "filter", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return err
	}

	err = clearOrCreateChain(ipt, "nat", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return err
	}

	err = ipt.AppendUnique("filter", "FORWARD", "-j", "WG_ACCESS_SERVER_FORWARD")
	if err != nil {
		return errors.Wrap(err, "failed to append FORWARD rule to filter chain")
	}

	err = ipt.AppendUnique("nat", "POSTROUTING", "-j", "WG_ACCESS_SERVER_POSTROUTING")
	if err != nil {
		return errors.Wrap(err, "failed to append POSTROUTING rule to nat chain")
	}

	if options.ClientIsolation {
		// Reject inter-device traffic
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", options.CIDRv6, "-d", options.CIDRv6, "-j", "REJECT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}
	// Accept client traffic for given allowed ips
	for _, allowedCIDR := range options.allowedIPv6s {
		if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", options.CIDRv6, "-d", allowedCIDR, "-j", "ACCEPT"); err != nil {
			return errors.Wrap(err, "failed to set ip tables rule")
		}
	}
	// And reject everything else
	if err := ipt.AppendUnique("filter", "WG_ACCESS_SERVER_FORWARD", "-s", options.CIDRv6, "-j", "REJECT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}

	if options.GatewayIface != "" {
		if options.NAT66 {
			if err := ipt.AppendUnique("nat", "WG_ACCESS_SERVER_POSTROUTING", "-s", options.CIDRv6, "-o", options.GatewayIface, "-j", "MASQUERADE"); err != nil {
				return errors.Wrap(err, "failed to set ip tables rule")
			}
		}
	}
	return nil
}

func clearOrCreateChain(ipt *iptables.IPTables, table, chain string) error {
	exists, err := ipt.ChainExists(table, chain)
	if err != nil {
		return errors.Wrapf(err, "failed to read table %s", table)
	}
	if exists {
		err = ipt.ClearChain(table, chain)
		if err != nil {
			return errors.Wrapf(err, "failed to clear chain %s in table %s", chain, table)
		}
	} else {
		// Create our own chain for forwarding rules
		err = ipt.NewChain(table, chain)
		if err != nil {
			return errors.Wrapf(err, "failed to create chain %s in table %s", chain, table)
		}
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
