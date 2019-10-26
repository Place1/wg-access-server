package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/pkg/errors"
	"github.com/place1/wireguard-access-server/internal/config"
	"github.com/place1/wireguard-access-server/internal/services"
	"github.com/place1/wireguard-access-server/internal/storage"
	"github.com/place1/wireguard-access-server/internal/web"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.Read()

	// Userspace wireguard command
	if config.WireGuard.UserspaceImplementation != "" {
		go func() {
			// execute the userspace wireguard implementation
			// if it exists/crashes for some reason then we'll also crash
			if err := services.ExecUserWireGuard(config.WireGuard.UserspaceImplementation, config.WireGuard.InterfaceName); err != nil {
				logrus.Fatal(err)
			}
		}()
		// Wait for the userspace wireguard process to
		// startup and create the wg0 interface
		// Super sorry if this just caused a race
		// condition for you :(
		time.Sleep(1 * time.Second)
	}

	// WireGuard
	wgserver, err := services.NewWireGuard(
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

	// Networking configuration
	if err := services.ConfigureRouting(config.WireGuard.InterfaceName); err != nil {
		logrus.Fatal(err)
	}
	if config.VPN.GatewayInterface != nil {
		if err := services.ConfigureForwarding(config.WireGuard.InterfaceName, config.VPN.GatewayInterface.Attrs().Name); err != nil {
			logrus.Fatal(err)
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
