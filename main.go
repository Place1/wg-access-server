package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/place1/wg-access-server/proto/proto"

	"github.com/gorilla/mux"
	"github.com/place1/wg-embed/pkg/wgembed"

	"github.com/pkg/errors"
	"github.com/place1/wg-access-server/internal/config"
	"github.com/place1/wg-access-server/internal/devices"
	"github.com/place1/wg-access-server/internal/dnsproxy"
	"github.com/place1/wg-access-server/internal/network"
	"github.com/place1/wg-access-server/internal/services"
	"github.com/place1/wg-access-server/internal/storage"
	"github.com/place1/wg-access-server/pkg/authnz"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
	"github.com/sirupsen/logrus"

	"net/http/httputil"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
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

	wg.LoadConfig(&wgembed.ConfigFile{
		Interface: wgembed.IfaceConfig{
			PrivateKey: conf.WireGuard.PrivateKey,
			Address:    vpnip.IP.String(),
			ListenPort: &conf.WireGuard.Port,
		},
	})

	// Networking configuration
	rules := network.NetworkRules{
		AllowVPNLAN:    true,
		AllowServerLAN: true,
		AllowInternet:  true,
	}
	if err := network.ConfigureForwarding(conf.WireGuard.InterfaceName, conf.VPN.GatewayInterface, conf.VPN.CIDR, rules); err != nil {
		logrus.Fatal(err)
	}

	// DNS Server
	dns, err := dnsproxy.New(conf.DNS.Upstream)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to start dns server"))
	}
	defer dns.Close()

	// Storage
	var storageDriver storage.Storage
	if conf.Storage.Directory != "" {
		logrus.Infof("storing data in %s", conf.Storage.Directory)
		storageDriver = storage.NewDiskStorage(conf.Storage.Directory)
	} else {
		storageDriver = storage.NewMemoryStorage()
	}

	// Services
	deviceManager := devices.New(wg.Name(), storageDriver, conf.VPN.CIDR)
	if err := deviceManager.StartSync(conf.DisableMetadata); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to sync"))
	}

	// Router
	router := mux.NewRouter()

	// Healthcheck endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	// if the built website exists, serve that
	// otherwise proxy to a local webpack development server
	if _, err := os.Stat("website/build"); os.IsNotExist(err) {
		u, _ := url.Parse("http://localhost:3000")
		router.NotFoundHandler = httputil.NewSingleHostReverseProxy(u)
	} else {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("website/build")))
	}

	// GRPC Server
	server := grpc.NewServer([]grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(1 * math.Pow(2, 20))), // 1MB
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}...)
	proto.RegisterServerServer(server, &services.ServerService{
		Config: conf,
	})
	proto.RegisterDevicesServer(server, &services.DeviceService{
		DeviceManager: deviceManager,
	})
	grpcServer := grpcweb.WrapServer(server)

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithField("stack", string(debug.Stack())).Error(err)
			}
		}()
		if grpcServer.IsGrpcWebRequest(r) {
			grpcServer.ServeHTTP(w, r)
		} else {
			if authsession.Authenticated(r.Context()) {
				router.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			}
		}
	})

	if conf.Auth.IsEnabled() {
		handler = authnz.New(conf.Auth, func(user *authsession.Identity) error {
			if user.Subject == conf.AdminSubject {
				user.Claims.Add("admin", "true")
			}
			return nil
		}).Wrap(handler)
	} else {
		base := handler
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			base.ServeHTTP(w, r.WithContext(authsession.SetIdentityCtx(r.Context(), &authsession.AuthSession{
				Identity: &authsession.Identity{
					Subject: "",
				},
			})))
		})
	}

	// Listen
	address := "0.0.0.0:8000"
	srv := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	// Start Web server
	logrus.Infof("listening on %v", address)
	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatal(errors.Wrap(err, "unable to start http server"))
	}
}

func randomBytes(size int) []byte {
	blk := make([]byte, size)
	_, err := rand.Read(blk)
	if err != nil {
		logrus.Fatal(err)
	}
	return blk
}
