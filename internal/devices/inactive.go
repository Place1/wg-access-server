package devices

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func inactiveLoop(d *DeviceManager, inactiveDuration time.Duration) {
	for {
		checkAndRemove(d, inactiveDuration)
		time.Sleep(30 * time.Second)
	}
}

func checkAndRemove(d *DeviceManager, inactiveDuration time.Duration) {
	logrus.Debug("inactive check executing")

	devices, err := d.ListAllDevices()
	if err != nil {
		logrus.Warn(errors.Wrap(err, "failed to list devices - inactive devices cannot be deleted"))
		return
	}

	for _, dev := range devices {
		logrus.Debugf("checking inactive device: %s/%s", dev.Owner, dev.Name)

		var elapsed time.Duration
		if dev.LastHandshakeTime == nil {
			// Never connected
			elapsed = time.Since(dev.CreatedAt)
		} else {
			elapsed = time.Since(*dev.LastHandshakeTime)
		}

		if elapsed > inactiveDuration {
			logrus.Debug("deleting inactive device")
			err := d.DeleteDevice(dev.Owner, dev.Name)
			if err != nil {
				logrus.Error(errors.Wrap(err, "failed to delete device"))
				continue
			}
		}
	}
}
