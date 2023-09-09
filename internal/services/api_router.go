package services

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/freifunkMUC/wg-embed/pkg/wgembed"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/freifunkMUC/wg-access-server/internal/config"
	"github.com/freifunkMUC/wg-access-server/internal/devices"
	"github.com/freifunkMUC/wg-access-server/internal/traces"
	"github.com/freifunkMUC/wg-access-server/proto/proto"
	"github.com/sirupsen/logrus"
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
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				// wrapped in anonymous func to get ctx
				return grpcLogrus.UnaryServerInterceptor(grpcLoggerWith(ctx))(ctx, req, info, handler)
			},
			grpcRecovery.UnaryServerInterceptor(
				grpcRecovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
					// add trace id to error message so it's visible for the client
					return status.Errorf(codes.Internal, "%v; trace = %s", p, traces.TraceID(ctx))
				}),
			),
		)),
	}...)

	// Register GRPC services
	proto.RegisterServerServer(server, &ServerService{
		Config: deps.Config,
		Wg:     deps.Wg,
	})
	proto.RegisterUsersServer(server, &UserService{
		DeviceManager: deps.DeviceManager,
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

// GRPC events have the nature of DEBUG logs but are logged with INFO level. To clean up the log stream starting from INFO log level we only log WARN events.
func grpcLoggerWith(ctx context.Context) *logrus.Entry {
	if (logrus.GetLevel() == logrus.TraceLevel || logrus.GetLevel() == logrus.DebugLevel) {
		return traces.Logger(ctx)
	} else {
		return traces.WarnLogger(ctx)
	}
}
