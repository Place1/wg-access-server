package devices

import (
	"github.com/pkg/errors"
	"github.com/place1/wg-access-server/internal/storage"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/place1/wg-embed/pkg/wgembed"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func metricsLoop(d *DeviceManager) {
	for {
		syncMetrics(d)
		time.Sleep(5 * time.Second)
	}
}

func syncMetrics(d *DeviceManager) {
	devices, err := d.ListAllDevices()
	if err != nil {
		logrus.Warn(errors.Wrap(err, "failed to list devices - metrics cannot be recorded"))
		return
	}
	peers, err := wgembed.ListPeers(d.iface)
	if err != nil {
		logrus.Warn(errors.Wrap(err, "failed to list peers - metrics cannot be recorded"))
		return
	}
	for _, peer := range peers {
		for _, device := range devices {
			if peer.PublicKey.String() == device.PublicKey {
				device.Connected = true
				device.Endpoint = peer.Endpoint.String()
				device.ReceiveBytes = peer.ReceiveBytes
				device.TransmitBytes = peer.TransmitBytes
				device.LastHandshakeTime = peer.LastHandshakeTime
			}
		}
	}
	for _, device := range devices {
		if !isConnected(peers, device) {
			device.Connected = false
			device.Endpoint = "-"
			device.ReceiveBytes = 0
			device.TransmitBytes = 0
		}
	}
}

func isConnected(currentPeers []wgtypes.Peer, device *storage.Device) bool {
	for _, peer := range currentPeers {
		if peer.PublicKey.String() == device.PublicKey {
			return true
		}
	}
	return false
}
