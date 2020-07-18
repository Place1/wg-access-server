package main

import (
	"fmt"
	"net/http"

	"github.com/place1/wg-access-server/internal/services"
	"github.com/place1/wg-access-server/internal/storage"
	"github.com/place1/wg-access-server/pkg/authnz"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"

	"github.com/gorilla/mux"
	"github.com/place1/wg-embed/pkg/wgembed"

	"github.com/pkg/errors"
	"github.com/place1/wg-access-server/internal/config"
	"github.com/place1/wg-access-server/internal/devices"
	"github.com/place1/wg-access-server/internal/dnsproxy"
	"github.com/place1/wg-access-server/internal/network"
	"github.com/sirupsen/logrus"
)

func main() {
	conf := config.Read()

	// The server's IP within the VPN virtual network
	vpnip := network.ServerVPNIP(conf.VPN.CIDR)

	// WireGuard Server
	wg, err := wgembed.New(conf.WireGuard.InterfaceName)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wireguard interface"))
	}
	defer wg.Close()

	logrus.Infof("starting wireguard server on 0.0.0.0:%d", conf.WireGuard.Port)

	wgconfig := &wgembed.ConfigFile{
		Interface: wgembed.IfaceConfig{
			PrivateKey: conf.WireGuard.PrivateKey,
			Address:    vpnip.String(),
			ListenPort: &conf.WireGuard.Port,
		},
	}

	if err := wg.LoadConfig(wgconfig); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to load wireguard config"))
	}

	logrus.Infof("wireguard VPN network is %s", conf.VPN.CIDR)

	if err := network.ConfigureForwarding(conf.WireGuard.InterfaceName, conf.VPN.GatewayInterface, conf.VPN.CIDR, conf.VPN.AllowedIPs); err != nil {
		logrus.Fatal(err)
	}

	// DNS Server
	if *conf.DNS.Enabled {
		dns, err := dnsproxy.New(dnsproxy.DNSServerOpts{
			Upstream: conf.DNS.Upstream,
		})
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to start dns server"))
		}
		defer dns.Close()
	}

	// Storage
	storageBackend, err := storage.NewStorage(conf.Storage)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create storage backend"))
	}
	if err := storageBackend.Open(); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to connect/open storage backend"))
	}
	defer storageBackend.Close()

	// Services
	deviceManager := devices.New(wg.Name(), storageBackend, conf.VPN.CIDR)
	if err := deviceManager.StartSync(conf.DisableMetadata); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to sync"))
	}

	router := mux.NewRouter()
	router.Use(services.TracesMiddleware)
	router.Use(services.RecoveryMiddleware)

	// Health check endpoint
	router.PathPrefix("/health").Handler(services.HealthEndpoint())

	// Authentication middleware
	if conf.Auth.IsEnabled() {
		router.Use(authnz.NewMiddleware(conf.Auth, claimsMiddleware(conf)))
	} else {
		logrus.Warn("[DEPRECATION NOTICE] using wg-access-server without an admin user is deprecated and will be removed in an upcoming minior release.")
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r.WithContext(authsession.SetIdentityCtx(r.Context(), &authsession.AuthSession{
					Identity: &authsession.Identity{
						Subject: "",
					},
				})))
			})
		})
	}

	// Subrouter for our site (web + api)
	site := router.PathPrefix("/").Subrouter()
	site.Use(authnz.RequireAuthentication)

	// Grpc api
	site.PathPrefix("/api").Handler(services.ApiRouter(&services.ApiServices{
		Config:        conf,
		DeviceManager: deviceManager,
	}))

	// Static website
	site.PathPrefix("/").Handler(services.WebsiteRouter())

	// publicRouter.NotFoundHandler = authMiddleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	if authsession.Authenticated(r.Context()) {
	// 		router.ServeHTTP(w, r)
	// 	} else {
	// 		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
	// 	}
	// }))
	publicRouter := router

	// Listen
	address := fmt.Sprintf("0.0.0.0:%d", conf.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: publicRouter,
	}

	// Start Web server
	logrus.Infof("web ui listening on %v", address)
	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatal(errors.Wrap(err, "unable to start http server"))
	}
}

func claimsMiddleware(conf *config.AppConfig) authsession.ClaimsMiddleware {
	return func(user *authsession.Identity) error {
		if user.Subject == conf.AdminSubject {
			user.Claims.Add("admin", "true")
		}
		return nil
	}
}
