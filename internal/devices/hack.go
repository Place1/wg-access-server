package devices

import (
	"fmt"
	"github.com/place1/wg-embed/pkg/wgembed"
)

func x(d *DeviceManager) {
	devices, _ := d.ListAllDevices()
	peers, _ := wgembed.ListPeers("wg0")
	for _, peer := range peers {
		for _, device := range devices {
			if peer.PublicKey.String() == device.PublicKey {
				device.Endpoint = peer.Endpoint.String()
				device.ReceiveBytes = peer.ReceiveBytes
				device.TransmitBytes = peer.TransmitBytes
				device.LastHandshakeTime = peer.LastHandshakeTime

				device.LifetimeReceivedBytes = -1
				device.LifetimeTransmitBytes = -1
			}
		}
	}
}
