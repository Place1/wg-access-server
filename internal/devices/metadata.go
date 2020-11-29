package devices

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func metadataLoop(d *DeviceManager) {
	for {
		syncMetrics(d)
		time.Sleep(30 * time.Second)
	}
}

func syncMetrics(d *DeviceManager) {
	logrus.Debug("metadata sync executing")

	peers, err := d.wg.ListPeers()
	if err != nil {
		logrus.Warn(errors.Wrap(err, "failed to list peers - metrics cannot be recorded"))
		return
	}

	for _, peer := range peers {
		// if the peer is connected we can update their metrics
		// importantly, we'll ignore peers that we know about
		// but aren't connected at the moment.
		// they may actually be connected to another replica.
		if peer.Endpoint != nil {
			if device, err := d.GetByPublicKey(peer.PublicKey.String()); err == nil {
				device.Endpoint = peer.Endpoint.IP.String()
				device.ReceiveBytes = peer.ReceiveBytes
				device.TransmitBytes = peer.TransmitBytes
				if !peer.LastHandshakeTime.IsZero() {
					device.LastHandshakeTime = &peer.LastHandshakeTime
				}
				if err := d.SaveDevice(device); err != nil {
					logrus.Error(errors.Wrap(err, "failed to save device during metadata sync"))
				}
			}
		}
	}
}
