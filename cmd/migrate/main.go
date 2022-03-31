package migrate

import (
	"github.com/freifunkMUC/wg-access-server/internal/storage"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func Register(app *kingpin.Application) *migratecmd {
	cmd := &migratecmd{}
	cli := app.Command(cmd.Name(), "Migrate your wg-access-server devices between storage backends. This tool is provided on a best effort bases.")
	cli.Arg("source", "The source storage URI").Required().StringVar(&cmd.src)
	cli.Arg("destination", "The destination storage URI").Required().StringVar(&cmd.dest)
	return cmd
}

type migratecmd struct {
	src  string
	dest string
}

func (cmd *migratecmd) Name() string {
	return "migrate"
}

func (cmd *migratecmd) Run() {
	srcBackend, err := storage.NewStorage(cmd.src)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create src storage backend"))
	}
	if err := srcBackend.Open(); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to connect/open src storage backend"))
	}
	defer srcBackend.Close()

	destBackend, err := storage.NewStorage(cmd.dest)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create destination storage backend"))
	}
	if err := destBackend.Open(); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to connect/open destination storage backend"))
	}
	defer destBackend.Close()

	srcDevices, err := srcBackend.List("")
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to list all devices from source storage backend"))
	}

	logrus.Infof("copying %v devices from source --> destination backend", len(srcDevices))

	for _, device := range srcDevices {
		if err := destBackend.Save(device); err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to write device to destination storage backend"))
		}
	}
}
