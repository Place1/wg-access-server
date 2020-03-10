package devices

import (
	"time"

	"github.com/pkg/errors"
	"github.com/place1/wg-embed/pkg/wgembed"
	"github.com/sirupsen/logrus"
)

func metadataLoop(d *DeviceManager) {
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
				device.ReceiveBytes = peer.ReceiveBytes
				device.TransmitBytes = peer.TransmitBytes
				if !peer.LastHandshakeTime.IsZero() {
					device.LastHandshakeTime = &peer.LastHandshakeTime
				}
				if peer.Endpoint != nil {
					device.Endpoint = peer.Endpoint.IP.String()
				}
				if err := d.SaveDevice(device); err != nil {
					logrus.Debug(errors.Wrap(err, "failed to update device metadata"))
				}
			}
		}
	}
}
