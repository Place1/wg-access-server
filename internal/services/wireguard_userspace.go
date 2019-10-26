package services

import (
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func ExecUserWireGuard(wgcommand string, ifaceName string) error {
	logrus.Infof("using userspace wireguard implementation %s", wgcommand)

	// create the command to exec
	// if it's "boringtun" we'll provide some non-standard
	// flags to better support running within a docker container
	var cmd *exec.Cmd
	if wgcommand == "boringtun" {
		cmd = exec.Command(
			wgcommand,
			ifaceName,
			"--disable-drop-privileges=root",
			"--foreground",
		)
	} else {
		cmd = exec.Command(
			wgcommand,
			"-f",
			ifaceName,
		)
	}

	entry := logrus.NewEntry(logrus.New()).WithField("process", wgcommand)
	cmd.Stdout = entry.Writer()
	cmd.Stderr = entry.Writer()
	logrus.Infof("starting %s", cmd.String())
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "userspace wireguard exitted")
	}

	return nil
}
