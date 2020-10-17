package serve

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/place1/wg-access-server/internal/services"
	"github.com/place1/wg-access-server/internal/storage"
	"github.com/place1/wg-access-server/pkg/authnz"
	"github.com/place1/wg-access-server/pkg/authnz/authconfig"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
	"github.com/vishvananda/netlink"
	"golang.org/x/crypto/bcrypt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
	"github.com/place1/wg-embed/pkg/wgembed"

	"github.com/pkg/errors"
	"github.com/place1/wg-access-server/internal/config"
	"github.com/place1/wg-access-server/internal/devices"
	"github.com/place1/wg-access-server/internal/dnsproxy"
	"github.com/place1/wg-access-server/internal/network"
	"github.com/sirupsen/logrus"
)

func RegisterCommand(app *kingpin.Application) *servecmd {
	cmd := &servecmd{}
	cli := app.Command(cmd.Name(), "Run the server")
	cli.Flag("config", "Path to a config file").Envar("CONFIG").StringVar(&cmd.configPath)
	cli.Flag("web-port", "The port that the web ui server will listen on").Envar("WEB_PORT").Default("8000").IntVar(&cmd.webPort)
	cli.Flag("wireguard-port", "The port that the Wireguard server will listen on").Envar("WIREGUARD_PORT").Default("51820").IntVar(&cmd.wireguardPort)
	cli.Flag("storage", "The storage backend connection string").Envar("STORAGE").Default("memory://").StringVar(&cmd.storage)
	cli.Flag("wireguard-private-key", "Wireguard private key").Envar("WIREGUARD_PRIVATE_KEY").StringVar(&cmd.privateKey)
	cli.Flag("disable-metadata", "Disable metadata collection (i.e. metrics)").Envar("DISABLE_METADATA").Default("false").BoolVar(&cmd.disableMetadata)
	cli.Flag("admin-username", "Admin username (defaults to admin)").Envar("ADMIN_USERNAME").Default("admin").StringVar(&cmd.adminUsername)
	cli.Flag("admin-password", "Admin password (provide plaintext, stored in-memory only)").Envar("ADMIN_PASSWORD").StringVar(&cmd.adminPassword)
	cli.Flag("upstream-dns", "An upstream DNS server to proxy DNS traffic to").Envar("UPSTREAM_DNS").StringVar(&cmd.upstreamDNS)
	return cmd
}

type servecmd struct {
	configPath      string
	webPort         int
	wireguardPort   int
	storage         string
	privateKey      string
	disableMetadata bool
	adminUsername   string
	adminPassword   string
	upstreamDNS     string
}

func (cmd *servecmd) Name() string {
	return "serve"
}

func (cmd *servecmd) Run() {
	conf := cmd.ReadConfig()

	// The server's IP within the VPN virtual network
	vpnip := network.ServerVPNIP(conf.VPN.CIDR)

	// WireGuard Server
	wg := wgembed.NewNoOpInterface()
	if conf.WireGuard.Enabled {
		wgimpl, err := wgembed.New(conf.WireGuard.InterfaceName)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to create wireguard interface"))
		}
		defer wgimpl.Close()
		wg = wgimpl

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
	}

	// DNS Server
	if conf.DNS.Enabled {
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
	deviceManager := devices.New(wg, storageBackend, conf.VPN.CIDR)
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
		logrus.Warn("[DEPRECATION NOTICE] using wg-access-server without an admin user is deprecated and will be removed in an upcoming minor release.")
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
		Wg:            wg,
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

func (cmd *servecmd) ReadConfig() *config.AppConfig {
	// here we're filling out the config struct
	// with values from our flags/defaults.
	config := &config.AppConfig{}
	config.Port = cmd.webPort
	config.WireGuard.InterfaceName = "wg0"
	config.WireGuard.Port = cmd.wireguardPort
	config.VPN.CIDR = "10.44.0.0/24"
	config.DisableMetadata = cmd.disableMetadata
	config.WireGuard.Enabled = true
	config.WireGuard.PrivateKey = cmd.privateKey
	config.Storage = cmd.storage
	config.VPN.AllowedIPs = []string{"0.0.0.0/0"}
	config.DNS.Enabled = true
	config.AdminPassword = cmd.adminPassword
	config.AdminSubject = cmd.adminUsername

	if cmd.upstreamDNS != "" {
		config.DNS.Upstream = []string{cmd.upstreamDNS}
	}

	if cmd.configPath != "" {
		if b, err := ioutil.ReadFile(cmd.configPath); err == nil {
			if err := yaml.Unmarshal(b, &config); err != nil {
				logrus.Fatal(errors.Wrap(err, "failed to bind configuration file"))
			}
		}
	}

	if config.LogLevel != "" {
		level, err := logrus.ParseLevel(config.LogLevel)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "invalid log level - should be one of fatal, error, warn, info, debug, trace"))
		}
		logrus.SetLevel(level)
	}

	if config.DisableMetadata {
		logrus.Info("Metadata collection has been disabled. No metrics or device connectivity information will be recorded or shown")
	}

	if config.VPN.GatewayInterface == "" {
		iface, err := defaultInterface()
		if err != nil {
			logrus.Warn(errors.Wrap(err, "failed to set default value for VPN.GatewayInterface"))
		} else {
			config.VPN.GatewayInterface = iface
		}
	}

	if config.WireGuard.PrivateKey == "" {
		if !strings.HasPrefix(config.Storage, "memory://") {
			logrus.Fatal(missingPrivateKey)
		}
		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to generate a server private key"))
		}
		config.WireGuard.PrivateKey = key.String()
	}

	if config.AdminPassword != "" && config.AdminSubject != "" {
		if config.Auth.Basic == nil {
			config.Auth.Basic = &authconfig.BasicAuthConfig{}
		}
		// htpasswd.AcceptBcrypt(config.AdminPassword)
		pw, err := bcrypt.GenerateFromPassword([]byte(config.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to generate a bcrypt hash for the provided admin password"))
		}
		config.Auth.Basic.Users = append(config.Auth.Basic.Users, fmt.Sprintf("%s:%s", config.AdminSubject, string(pw)))
	}

	return config
}

func claimsMiddleware(conf *config.AppConfig) authsession.ClaimsMiddleware {
	return func(user *authsession.Identity) error {
		if user.Subject == conf.AdminSubject {
			user.Claims.Add("admin", "true")
		}
		return nil
	}
}

func defaultInterface() (string, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return "", errors.Wrap(err, "failed to list network interfaces")
	}
	for _, link := range links {
		routes, err := netlink.RouteList(link, 4)
		if err != nil {
			return "", errors.Wrapf(err, "failed to list routes for interface %s", link.Attrs().Name)
		}
		for _, route := range routes {
			if route.Dst == nil {
				return link.Attrs().Name, nil
			}
		}
	}
	return "", errors.New("could not determine the default network interface name")
}

var missingPrivateKey = `missing wireguard private key:

    create a key:

        $ wg genkey

    configure via environment variable:

        $ export WIREGUARD_PRIVATE_KEY="<private-key>"

    or configure via flag:

        $ wg-access-server serve --wireguard-private-key="<private-key>"

    or configure via file:

      wireguard:
        privateKey: "<private-key>"

`
