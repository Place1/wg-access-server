package services

import (
	"context"

	"github.com/place1/wg-access-server/internal/auth/authsession"
	"github.com/place1/wg-access-server/internal/config"
	"github.com/place1/wg-access-server/proto/proto"
	"github.com/place1/wg-embed/pkg/wgembed"
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

	return &proto.InfoRes{
		Host:            stringValue(s.Config.WireGuard.ExternalHost),
		PublicKey:       publicKey,
		Port:            int32(s.Config.WireGuard.Port),
		HostVpnIp:       ServerVPNIP(s.Config.VPN.CIDR).IP.String(),
		MetadataEnabled: s.Config.MetadataEnabled,
	}, nil
}
