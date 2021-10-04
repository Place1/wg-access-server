package devices

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/place1/wg-embed/pkg/wgembed"

	"github.com/pkg/errors"
	"github.com/place1/wg-access-server/internal/network"
	"github.com/place1/wg-access-server/internal/storage"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
	"github.com/sirupsen/logrus"
)

type DeviceManager struct {
	wg      wgembed.WireGuardInterface
	storage storage.Storage
	cidr    string
	cidrv6  string
}

func New(wg wgembed.WireGuardInterface, s storage.Storage, cidr, cidrv6 string) *DeviceManager {
	return &DeviceManager{wg, s, cidr, cidrv6}
}

func (d *DeviceManager) StartSync(disableMetadataCollection bool) error {
	// Start listening to the device add/remove events
	d.storage.OnAdd(func(device *storage.Device) {
		logrus.Debugf("storage event: device added: %s/%s", device.Owner, device.Name)
		if err := d.wg.AddPeer(device.PublicKey, network.SplitAddresses(device.Address)); err != nil {
			logrus.Error(errors.Wrap(err, "failed to add wireguard peer"))
		}
	})

	d.storage.OnDelete(func(device *storage.Device) {
		logrus.Debugf("storage event: device removed: %s/%s", device.Owner, device.Name)
		if err := d.wg.RemovePeer(device.PublicKey); err != nil {
			logrus.Error(errors.Wrap(err, "failed to remove wireguard peer"))
		}
	})

	d.storage.OnReconnect(func() {
		if err := d.sync(); err != nil {
			logrus.Error(errors.Wrap(err, "device sync after storage backend reconnect event failed"))
		}
	})

	// Do an initial sync of existing devices
	if err := d.sync(); err != nil {
		return errors.Wrap(err, "initial device sync from storage failed")
	}

	// start the metrics loop
	if !disableMetadataCollection {
		go metadataLoop(d)
	}

	return nil
}

func (d *DeviceManager) AddDevice(identity *authsession.Identity, name string, publicKey string) (*storage.Device, error) {
	if name == "" {
		return nil, errors.New("device name must not be empty")
	}

	clientAddr, err := d.nextClientAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate an ip address for device")
	}

	device := &storage.Device{
		Owner:         identity.Subject,
		OwnerName:     identity.Name,
		OwnerEmail:    identity.Email,
		OwnerProvider: identity.Provider,
		Name:          name,
		PublicKey:     publicKey,
		Address:       clientAddr,
		CreatedAt:     time.Now(),
	}

	if err := d.SaveDevice(device); err != nil {
		return nil, errors.Wrap(err, "failed to save the new device")
	}

	return device, nil
}

func (d *DeviceManager) SaveDevice(device *storage.Device) error {
	return d.storage.Save(device)
}

func (d *DeviceManager) sync() error {
	devices, err := d.ListAllDevices()
	if err != nil {
		return errors.Wrap(err, "failed to list devices")
	}

	peers, err := d.wg.ListPeers()
	if err != nil {
		return errors.Wrap(err, "failed to list peers")
	}

	// Remove any peers for devices that are no longer in storage
	for _, peer := range peers {
		if !deviceListContains(devices, peer.PublicKey.String()) {
			if err := d.wg.RemovePeer(peer.PublicKey.String()); err != nil {
				logrus.Error(errors.Wrapf(err, "failed to remove peer during sync: %s", peer.PublicKey.String()))
			}
		}
	}

	// Add peers for all devices in storage
	for _, device := range devices {
		if err := d.wg.AddPeer(device.PublicKey, network.SplitAddresses(device.Address)); err != nil {
			logrus.Warn(errors.Wrapf(err, "failed to add device during sync: %s", device.Name))
		}
	}

	return nil
}

func (d *DeviceManager) ListAllDevices() ([]*storage.Device, error) {
	return d.storage.List("")
}

func (d *DeviceManager) ListDevices(user string) ([]*storage.Device, error) {
	return d.storage.List(user)
}

func (d *DeviceManager) DeleteDevice(user string, name string) error {
	device, err := d.storage.Get(user, name)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve device")
	}

	if err := d.storage.Delete(device); err != nil {
		return err
	}

	return nil
}

func (d *DeviceManager) GetByPublicKey(publicKey string) (*storage.Device, error) {
	return d.storage.GetByPublicKey(publicKey)
}

var nextIPLock = sync.Mutex{}

func (d *DeviceManager) nextClientAddress() (string, error) {
	nextIPLock.Lock()
	defer nextIPLock.Unlock()

	devices, err := d.ListDevices("")
	if err != nil {
		return "", errors.Wrap(err, "failed to list devices")
	}

	// TODO: read up on better ways to allocate client's IP
	// addresses from a configurable CIDR

	var usedIPv4s []net.IP
	var usedIPv6s []net.IP

	// Check what IP addresses are already occupied
	for _, device := range devices {
		addresses := network.SplitAddresses(device.Address)
		for _, addr := range addresses {
			ip, _ := MustParseCIDR(strings.TrimSpace(addr))
			if as4 := ip.To4(); as4 != nil {
				usedIPv4s = append(usedIPv4s, as4)
			} else {
				usedIPv6s = append(usedIPv6s, ip)
			}
		}
	}

	var ipv4 string
	var ipv6 string

	if d.cidr != "" {
		vpnipv4, vpnsubnetv4 := MustParseCIDR(d.cidr)
		startIPv4 := vpnipv4.Mask(vpnsubnetv4.Mask)

		// Add the network address and the VPN server address to the list of occupied addresses
		usedIPv4s = append(usedIPv4s,
			startIPv4,         // x.x.x.0
			nextIP(startIPv4), // x.x.x.1
		)

		for ip := startIPv4; vpnsubnetv4.Contains(ip); ip = nextIP(ip) {
			if !contains(usedIPv4s, ip) {
				ipv4 = fmt.Sprintf("%s/32", ip.String())
				break
			}
		}
	}

	if d.cidrv6 != "" {
		vpnipv6, vpnsubnetv6 := MustParseCIDR(d.cidrv6)
		startIPv6 := vpnipv6.Mask(vpnsubnetv6.Mask)

		// Add the network address and the VPN server address to the list of occupied addresses
		usedIPv6s = append(usedIPv6s,
			startIPv6,         // ::0
			nextIP(startIPv6), // ::1
		)

		for ip := startIPv6; vpnsubnetv6.Contains(ip); ip = nextIP(ip) {
			if !contains(usedIPv6s, ip) {
				ipv6 = fmt.Sprintf("%s/128", ip.String())
				break
			}
		}
	}

	if ipv4 != "" {
		if ipv6 != "" {
			return fmt.Sprintf("%s, %s", ipv4, ipv6), nil
		} else if d.cidrv6 != "" {
			return "", fmt.Errorf("there are no free IP addresses in the vpn subnet: '%s'", d.cidrv6)
		} else {
			return ipv4, nil
		}
	} else if ipv6 != "" {
		if d.cidr != "" {
			return "", fmt.Errorf("there are no free IP addresses in the vpn subnet: '%s'", d.cidr)
		} else {
			return ipv6, nil
		}
	} else {
		return "", fmt.Errorf("there are no free IP addresses in the vpn subnets: '%s', '%s'", d.cidr, d.cidrv6)
	}
}

func contains(ips []net.IP, target net.IP) bool {
	for _, ip := range ips {
		if ip.Equal(target) {
			return true
		}
	}
	return false
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

func deviceListContains(devices []*storage.Device, publicKey string) bool {
	for _, device := range devices {
		if device.PublicKey == publicKey {
			return true
		}
	}
	return false
}
