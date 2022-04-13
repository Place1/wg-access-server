package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/freifunkMUC/wg-access-server/cmd"
	"github.com/freifunkMUC/wg-access-server/cmd/migrate"
	"github.com/freifunkMUC/wg-access-server/cmd/serve"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("wg-access-server", "An all-in-one WireGuard Access Server & VPN solution")
	logLevel = app.Flag("log-level", "Log level: trace, debug, info, error, fatal").Envar("WG_LOG_LEVEL").Default("info").String()
)

func main() {
	// all the subcommands for wg-access-server
	commands := []cmd.Command{
		serve.Register(app),
		migrate.Register(app),
	}

	// parse CLI arguments
	clicmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	// set global log level
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

	for _, c := range commands {
		if clicmd == c.Name() {
			c.Run()
			return
		}
	}
}
