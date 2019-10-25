package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/coreos/go-iptables/iptables"

	"github.com/vishvananda/netlink"

	"github.com/gorilla/mux"

	"github.com/pkg/errors"
	"github.com/place1/wireguard-access-server/internal/config"
	"github.com/place1/wireguard-access-server/internal/services"
	"github.com/place1/wireguard-access-server/internal/storage"
	"github.com/place1/wireguard-access-server/internal/web"
	"github.com/place1/wireguard-access-server/internal/wg"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.Read()

	// Userspace wireguard command
	if config.WireGuard.UserspaceImplementation != "" {
		go func() {
			logrus.Infof("using userspace wireguard implementation %s", config.WireGuard.UserspaceImplementation)
			var command *exec.Cmd
			if config.WireGuard.UserspaceImplementation == "boringtun" {
				command = exec.Command(
					config.WireGuard.UserspaceImplementation,
					config.WireGuard.InterfaceName,
					"--disable-drop-privileges=root",
					"--foreground",
				)
			} else {
				command = exec.Command(
					config.WireGuard.UserspaceImplementation,
					"-f",
					config.WireGuard.InterfaceName,
				)
			}
			entry := logrus.NewEntry(logrus.New()).WithField("process", config.WireGuard.UserspaceImplementation)
			command.Stdout = entry.Writer()
			command.Stderr = entry.Writer()
			logrus.Infof("starting %s", command.String())
			if err := command.Run(); err != nil {
				logrus.Fatal(errors.Wrap(err, "userspace wireguard exitted"))
			}
		}()

		// Wait for the userspace wireguard process to
		// startup and create the wg0 interface
		// Super sorry if this just caused a race
		// condition for you :(
		time.Sleep(1 * time.Second)
	}

	// WireGuard
	wgserver, err := wg.New(
		config.WireGuard.InterfaceName,
		config.WireGuard.PrivateKey,
		config.WireGuard.Port,
		config.Web.ExternalAddress,
	)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wgserver"))
	}
	defer wgserver.Close()
	logrus.Infof("wireguard server public key is %s", wgserver.PublicKey())
	logrus.Infof("wireguard endpoint is %s", wgserver.Endpoint())

	// Networking configuration (ip links and route tables)
	link, err := netlink.LinkByName(config.WireGuard.InterfaceName)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to find wireguard interface"))
	}
	addr, err := netlink.ParseAddr("10.0.0.1/24")
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to parse subnet address"))
	}
	if err := netlink.AddrAdd(link, addr); err != nil {
		logrus.Warn(errors.Wrap(err, "failed to add subnet to wireguard interface"))
	}
	if err := netlink.LinkSetUp(link); err != nil {
		logrus.Warn(errors.Wrap(err, "failed to bring wireguard interface up"))
	}

	// Networking configuration (iptables)
	if config.VPN.GatewayInterface != nil {
		ipt, err := iptables.New()
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to init iptables"))
		}
		if err := ipt.AppendUnique("filter", "FORWARD", "-s", "10.0.0.1/24", "-o", config.WireGuard.InterfaceName, "-j", "ACCEPT"); err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to set ip tables rule"))
		}
		if err := ipt.AppendUnique("filter", "FORWARD", "-s", "10.0.0.1/24", "-i", config.WireGuard.InterfaceName, "-j", "ACCEPT"); err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to set ip tables rule"))
		}
		if err := ipt.AppendUnique("nat", "POSTROUTING", "-s", "10.0.0.1/24", "-o", config.VPN.GatewayInterface.Attrs().Name, "-j", "MASQUERADE"); err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to set ip tables rule"))
		}
	}

	// Storage
	var storageDriver storage.Storage
	if config.Storage.Directory != "" {
		storageDriver = storage.NewDiskStorage(config.Storage.Directory)
	} else {
		storageDriver = storage.NewMemoryStorage()
	}

	// Services
	deviceManager := services.NewDeviceManager(wgserver, storageDriver)
	if err := deviceManager.Sync(); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to sync"))
	}

	// Router
	router := mux.NewRouter()
	router.HandleFunc("/api/devices", web.AddDevice(deviceManager)).Methods("POST")
	router.HandleFunc("/api/devices", web.ListDevices(deviceManager)).Methods("GET")
	router.HandleFunc("/api/devices/{name}", web.DeleteDevice(deviceManager)).Methods("DELETE")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("website/build")))

	// Listen
	address := fmt.Sprintf("0.0.0.0:%d", config.Web.Port)
	logrus.Infof("website listening on %s", address)
	if err := http.ListenAndServe(address, router); err != nil {
		logrus.Fatal(errors.Wrap(err, "server exited"))
	}
}
