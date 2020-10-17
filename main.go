package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/place1/wg-access-server/cmd/migrate"
	"github.com/place1/wg-access-server/cmd/serve"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("wg-access-server", "An all-in-one WireGuard Access Server & VPN solution")
	logLevel = app.Flag("log-level", "Log level (debug, info, error)").Envar("LOG_LEVEL").Default("info").String()

	servecmd   = serve.RegisterCommand(app)
	migratecmd = migrate.RegisterCommand(app)
)

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "invalid log level - should be one of fatal, error, warn, info, debug, trace"))
	}

	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		},
	})

	switch cmd {
	case servecmd.Name():
		servecmd.Run()
	case migratecmd.Name():
		migratecmd.Run()
	default:
		logrus.Fatal(fmt.Errorf("unknown command: %s", cmd))
	}
}
