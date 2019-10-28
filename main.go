package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gorilla/mux"

	"github.com/pkg/errors"
	"github.com/place1/wireguard-access-server/internal/auth"
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
		config.WireGuard.ExternalAddress,
	)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wgserver"))
	}
	defer wgserver.Close()
	logrus.Infof("wireguard server public key is %s", wgserver.PublicKey())
	logrus.Infof("wireguard endpoint is %s", wgserver.Endpoint())

	// Networking configuration
	if err := services.ConfigureRouting(config.WireGuard.InterfaceName, config.VPN.CIDR); err != nil {
		logrus.Fatal(err)
	}
	if config.VPN.GatewayInterface != "" {
		logrus.Infof("vpn gateway interface is %s", config.VPN.GatewayInterface)
		if err := services.ConfigureForwarding(config.WireGuard.InterfaceName, config.VPN.GatewayInterface, config.VPN.CIDR); err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Warn("VPN.GatewayInterface is not configured - vpn clients will not have access to the internet")
	}

	// Storage
	var storageDriver storage.Storage
	if config.Storage.Directory != "" {
		storageDriver = storage.NewDiskStorage(config.Storage.Directory)
	} else {
		storageDriver = storage.NewMemoryStorage()
	}

	// Services
	deviceManager := services.NewDeviceManager(wgserver, storageDriver, config.VPN.CIDR)
	if err := deviceManager.Sync(); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to sync"))
	}

	// Http sessions
	session := scs.New()
	session.Store = memstore.New()

	// Router
	router := mux.NewRouter()
	if dex := dexIntegration(config, session); dex != nil {
		router.PathPrefix("/auth").Handler(dex)
	}
	secureRouter := router.PathPrefix("/").Subrouter()
	secureRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if session.GetString(r.Context(), "auth/subject") == "" {
				http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})
	secureRouter.HandleFunc("/api/devices", web.AddDevice(deviceManager)).Methods("POST")
	secureRouter.HandleFunc("/api/devices", web.ListDevices(deviceManager)).Methods("GET")
	secureRouter.HandleFunc("/api/devices/{name}", web.DeleteDevice(deviceManager)).Methods("DELETE")
	secureRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("website/build")))

	// Listen
	address := fmt.Sprintf("0.0.0.0:%d", config.Web.Port)
	logrus.Infof("website external address is %s", config.Web.ExternalAddress)
	logrus.Infof("website listening on %s", address)
	if err := http.ListenAndServe(address, session.LoadAndSave(router)); err != nil {
		logrus.Fatal(errors.Wrap(err, "server exited"))
	}
}

func dexIntegration(config *config.AppConfig, session *scs.SessionManager) *mux.Router {
	authBackends := []auth.AuthConnector{}
	if config.Auth.OIDC != nil {
		logrus.Infof("adding oidc auth backend %s", config.Auth.OIDC.Name)
		authBackends = append(authBackends, config.Auth.OIDC)
	}
	if config.Auth.Gitlab != nil {
		logrus.Infof("adding gitlab auth backend %s", config.Auth.Gitlab.Name)
		authBackends = append(authBackends, config.Auth.Gitlab)
	}
	c := auth.Config{}
	if len(authBackends) > 0 {
		c.Connectors = authBackends
	}
	if config.Auth.StaticUser != nil {
		c.StaticUsers = []auth.StaticUser{*config.Auth.StaticUser}
	}
	if c.Connectors != nil || c.StaticUsers != nil {
		dex, err := auth.NewDexServer(session, config.Web.ExternalAddress, config.Web.Port, c)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to initialize auth system"))
		}
		return dex.Router()
	}
	return nil
}
