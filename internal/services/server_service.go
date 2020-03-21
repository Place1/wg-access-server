package services

import (
	"github.com/place1/wg-access-server/internal/network"
	"context"

	"github.com/place1/wg-access-server/internal/config"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
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
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
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
		HostVpnIp:       network.ServerVPNIP(s.Config.VPN.CIDR).IP.String(),
		MetadataEnabled: !s.Config.DisableMetadata,
		IsAdmin:         user.Claims.Contains("admin"),
	}, nil
}
