package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"net/http"
	"runtime/debug"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/place1/wireguard-access-server/proto/proto"

	"github.com/gorilla/mux"
	"github.com/place1/wg-embed/pkg/wgembed"

	"github.com/pkg/errors"
	"github.com/place1/wireguard-access-server/internal/auth"
	"github.com/place1/wireguard-access-server/internal/config"
	"github.com/place1/wireguard-access-server/internal/devices"
	"github.com/place1/wireguard-access-server/internal/dnsproxy"
	"github.com/place1/wireguard-access-server/internal/services"
	"github.com/place1/wireguard-access-server/internal/storage"
	"github.com/sirupsen/logrus"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

func main() {
	conf := config.Read()

	// The server's IP within the VPN virtual network
	vpnip := services.ServerVPNIP(conf.VPN.CIDR)

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
	if err := services.ConfigureRouting(conf.WireGuard.InterfaceName, conf.VPN.CIDR); err != nil {
		logrus.Fatal(err)
	}
	if conf.VPN.GatewayInterface != "" {
		logrus.Infof("vpn gateway interface is %s", conf.VPN.GatewayInterface)
		if err := services.ConfigureForwarding(conf.WireGuard.InterfaceName, conf.VPN.GatewayInterface, conf.VPN.CIDR); err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Warn("VPN.GatewayInterface is not configured - vpn clients will not have access to the internet")
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
		storageDriver = storage.NewDiskStorage(conf.Storage.Directory)
	} else {
		storageDriver = storage.NewMemoryStorage()
	}

	// Services
	deviceManager := devices.New(wg.Name(), storageDriver, conf.VPN.CIDR)
	if err := deviceManager.Sync(); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to sync"))
	}

	// Router
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("website/build")))

	// GRPC Server
	server := grpc.NewServer([]grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(1 * math.Pow(2, 20))), // 1MB
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		)),
	}...)
	proto.RegisterServerServer(server, &services.ServerService{
		Config: conf,
	})
	proto.RegisterDevicesServer(server, &services.DeviceService{
		DeviceManager: deviceManager,
	})
	grpcServer := grpcweb.WrapServer(server)

	var handler http.Handler = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithField("stack", string(debug.Stack())).Error(err)
			}
		}()
		if grpcServer.IsGrpcWebRequest(req) {
			grpcServer.ServeHTTP(resp, req)
		} else {
			router.ServeHTTP(resp, req)
		}
	})

	if conf.Auth != nil {
		handler = auth.New(conf.Auth).Wrap(handler)
	}

	// Listen
	address := fmt.Sprintf("0.0.0.0:%d", conf.Web.Port)
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
