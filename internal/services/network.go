package services

import (
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

func ConfigureForwarding(wgIface string, gatewayIface string, cidr string) error {
	// Networking configuration (iptables) configuration
	// to ensure that traffic from clients the wireguard interface
	// is sent to the provided network interface
	ipt, err := iptables.New()
	if err != nil {
		return errors.Wrap(err, "failed to init iptables")
	}
	logrus.Infof("iptables rule - accept forwarding traffic from %s to interface %s", gatewayIface, wgIface)
	if err := ipt.AppendUnique("filter", "FORWARD", "-i", gatewayIface, "-o", wgIface, "-j", "ACCEPT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	logrus.Infof("iptables rule - accept forwarding traffic from %s to interface %s", wgIface, gatewayIface)
	if err := ipt.AppendUnique("filter", "FORWARD", "-i", wgIface, "-o", gatewayIface, "-j", "ACCEPT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	logrus.Infof("iptables rule - masquerade traffic from %s to interface %s", cidr, gatewayIface)
	if err := ipt.AppendUnique("nat", "POSTROUTING", "-s", cidr, "-o", gatewayIface, "-j", "MASQUERADE"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	return nil
}
