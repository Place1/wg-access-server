package devices

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/place1/wg-embed/pkg/wgembed"

	"github.com/pkg/errors"
	"github.com/place1/wg-access-server/internal/storage"
	"github.com/sirupsen/logrus"
)

type DeviceManager struct {
	iface   string
	storage storage.Storage
	cidr    string
}

func New(iface string, s storage.Storage, cidr string) *DeviceManager {
	return &DeviceManager{iface, s, cidr}
}

func (d *DeviceManager) StartSync(disableMetadataCollection bool) error {
	// sync devices from storage once
	devices, err := d.ListDevices("")
	if err != nil {
		return errors.Wrap(err, "failed to list devices")
	}
	for _, device := range devices {
		if err := wgembed.AddPeer(d.iface, device.PublicKey, device.Address); err != nil {
			logrus.Warn(errors.Wrapf(err, "failed to sync device '%s' (ignoring)", device.Name))
		}
	}

	// start the metrics loop
	if !disableMetadataCollection {
		go metadataLoop(d)
	}

	return nil
}

func (d *DeviceManager) AddDevice(user string, name string, publicKey string) (*storage.Device, error) {
	if name == "" {
		return nil, errors.New("device name must not be empty")
	}

	clientAddr, err := d.nextClientAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate an ip address for device")
	}

	device := &storage.Device{
		Owner:     user,
		Name:      name,
		PublicKey: publicKey,
		Address:   clientAddr,
		CreatedAt: time.Now(),
	}

	if err := d.SaveDevice(device); err != nil {
		return nil, errors.Wrap(err, "failed to save the new device")
	}

	if err := wgembed.AddPeer(d.iface, publicKey, clientAddr); err != nil {
		return nil, errors.Wrap(err, "unable to provision peer")
	}

	return device, nil
}

func (d *DeviceManager) SaveDevice(device *storage.Device) error {
	return d.storage.Save(key(device.Owner, device.Name), device)
}

func (d *DeviceManager) ListAllDevices() ([]*storage.Device, error) {
	return d.storage.List("")
}

func (d *DeviceManager) ListDevices(user string) ([]*storage.Device, error) {
	prefix := ""
	if user != "" {
		prefix = user + string(os.PathSeparator)
	}
	return d.storage.List(prefix)
}

func (d *DeviceManager) DeleteDevice(user string, name string) error {
	device, err := d.storage.Get(key(user, name))
	if err != nil {
		return errors.Wrap(err, "failed to retrieve device")
	}
	if err := d.storage.Delete(key(user, name)); err != nil {
		return err
	}
	if err := wgembed.RemovePeer(d.iface, device.PublicKey); err != nil {
		return errors.Wrap(err, "device was removed from storage but failed to be removed from the wireguard interface")
	}
	return nil
}

func key(user string, device string) string {
	return filepath.Join(user, device)
}

var nextIPLock = sync.Mutex{}

func (d *DeviceManager) nextClientAddress() (string, error) {
	nextIPLock.Lock()
	defer nextIPLock.Unlock()

	devices, err := d.ListDevices("")
	if err != nil {
		return "", errors.Wrap(err, "failed to list devices")
	}

	vpnip, vpnsubnet := MustParseCIDR(d.cidr)
	ip := vpnip.Mask(vpnsubnet.Mask)

	// TODO: read up on better ways to allocate client's IP
	// addresses from a configurable CIDR
	usedIPs := []net.IP{
		ip,         // x.x.x.0
		nextIP(ip), // x.x.x.1
	}
	for _, device := range devices {
		ip, _ := MustParseCIDR(device.Address)
		usedIPs = append(usedIPs, ip)
	}

	for ip := ip; vpnsubnet.Contains(ip); ip = nextIP(ip) {
		if !contains(usedIPs, ip) {
			return fmt.Sprintf("%s/32", ip.String()), nil
		}
	}

	return "", fmt.Errorf("there are no free IP addresses in the vpn subnet: '%s'", vpnsubnet)
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
