package services

import (
	"context"

	"github.com/place1/wg-embed/pkg/wgembed"
	"github.com/place1/wireguard-access-server/internal/auth/authsession"
	"github.com/place1/wireguard-access-server/internal/config"
	"github.com/place1/wireguard-access-server/proto/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerService struct {
	Config *config.AppConfig
}

func (s *ServerService) Info(ctx context.Context, req *proto.InfoReq) (*proto.InfoRes, error) {
	if _, err := authsession.CurrentUser(ctx); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	publicKey, err := wgembed.PublicKey(s.Config.WireGuard.InterfaceName)
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "failed to get public key")
	}

	port, err := wgembed.Port(s.Config.WireGuard.InterfaceName)
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "failed to get port")
	}

	return &proto.InfoRes{
		Host:      stringValue(s.Config.WireGuard.ExternalAddress),
		PublicKey: publicKey,
		Port:      int32(port),
	}, nil
}
