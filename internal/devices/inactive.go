package devices

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func inactiveLoop(d *DeviceManager, inactiveDeviceGracePeriod time.Duration) {
	for {
		checkAndRemove(d, inactiveDeviceGracePeriod)
		time.Sleep(30 * time.Second)
	}
}

func checkAndRemove(d *DeviceManager, inactiveDeviceGracePeriod time.Duration) {
	logrus.Debug("Inactive check executing")

	devices, err := d.ListAllDevices()
	if err != nil {
		logrus.Warn(errors.Wrap(err, "Failed to list devices - inactive devices cannot be deleted"))
		return
	}

	for _, dev := range devices {
		logrus.Debugf("Checking inactive device: %s/%s", dev.Owner, dev.Name)

		var elapsed time.Duration
		if dev.LastHandshakeTime == nil {
			// Never connected
			elapsed = time.Since(dev.CreatedAt)
		} else {
			elapsed = time.Since(*dev.LastHandshakeTime)
		}

		if elapsed > inactiveDeviceGracePeriod {
			logrus.Warnf("Deleting inactive device: %s/%s", dev.Owner, dev.Name)
			err := d.DeleteDevice(dev.Owner, dev.Name)
			if err != nil {
				logrus.Error(errors.Wrap(err, fmt.Sprintf("Failed to delete device: %s/%s", dev.Owner, dev.Name)))
				continue
			}
		}
	}
}
