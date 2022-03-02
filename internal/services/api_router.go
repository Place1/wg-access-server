package services

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/place1/wg-embed/pkg/wgembed"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/freifunkMUC/wg-access-server/internal/config"
	"github.com/freifunkMUC/wg-access-server/internal/devices"
	"github.com/freifunkMUC/wg-access-server/internal/traces"
	"github.com/freifunkMUC/wg-access-server/proto/proto"
	"google.golang.org/grpc"
)

type ApiServices struct {
	Config        *config.AppConfig
	DeviceManager *devices.DeviceManager
	Wg            wgembed.WireGuardInterface
}

func ApiRouter(deps *ApiServices) http.Handler {
	// Native GRPC server
	server := grpc.NewServer([]grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(1 * math.Pow(2, 20))), // 1MB
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return grpc_logrus.UnaryServerInterceptor(traces.Logger(ctx))(ctx, req, info, handler)
		}),
	}...)

	// Register GRPC services
	proto.RegisterServerServer(server, &ServerService{
		Config: deps.Config,
		Wg:     deps.Wg,
	})
	proto.RegisterDevicesServer(server, &DeviceService{
		DeviceManager: deps.DeviceManager,
	})

	// Grpc Web in process proxy (wrapper)
	grpcServer := grpcweb.WrapServer(server,
		grpcweb.WithAllowNonRootResource(true),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if grpcServer.IsGrpcWebRequest(r) {
			grpcServer.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(400)
		fmt.Fprintln(w, "expected grpc request")
	})
}
