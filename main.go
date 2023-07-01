package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/freifunkMUC/wg-access-server/cmd"
	"github.com/freifunkMUC/wg-access-server/cmd/migrate"
	"github.com/freifunkMUC/wg-access-server/cmd/serve"
)

var (
	app      = kingpin.New("wg-access-server", "An all-in-one WireGuard Access Server & VPN solution")
	logLevel = app.Flag("log-level", "Log level: trace, debug, info, error, fatal").Envar("WG_LOG_LEVEL").Default("info").String()
)

func main() {
	// All the subcommands for wg-access-server
	commands := []cmd.Command{
		serve.Register(app),
		migrate.Register(app),
	}

	// Parse CLI arguments
	clicmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Set global log level
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "Invalid log level - should be one of fatal, error, warn, info, debug, trace"))
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		},
	})

	// Hooks
	logrus.AddHook(&GrpcInfoLogDemotionHook{})

	for _, c := range commands {
		if clicmd == c.Name() {
			c.Run()
			return
		}
	}
}

// Logrus hook for downgrading these 'finished unary call...' GRPC info logs to debug.
type GrpcInfoLogDemotionHook struct {
}

func (h *GrpcInfoLogDemotionHook) Levels() []logrus.Level {
	// Only concerns info entries.
	return []logrus.Level{logrus.InfoLevel}
}

func (h *GrpcInfoLogDemotionHook) Fire(e *logrus.Entry) error {
	// Demotes info log lines like `INFO[0010]options.go:220 finished unary call with code OK grpc.code=OK grpc.method=ListDevices...` to debug level.
	if e.Data != nil && e.Data["system"] == "grpc" && e.Data["grpc.code"] == "OK" {
		e.Level = logrus.DebugLevel
	}
	return nil
}
