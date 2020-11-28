package services

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/place1/wg-access-server/internal/network"

	"github.com/place1/wg-access-server/internal/config"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
	"github.com/place1/wg-access-server/proto/proto"
	"github.com/place1/wg-embed/pkg/wgembed"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerService struct {
	Config *config.AppConfig
	Wg     wgembed.WireGuardInterface
}

func (s *ServerService) Info(ctx context.Context, req *proto.InfoReq) (*proto.InfoRes, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	publicKey, err := s.Wg.PublicKey()
	if err != nil {
		ctxlogrus.Extract(ctx).Error(err)
		return nil, status.Errorf(codes.Internal, "failed to get public key")
	}

	return &proto.InfoRes{
		Host:            stringValue(&s.Config.ExternalHost),
		PublicKey:       publicKey,
		Port:            int32(s.Config.WireGuard.Port),
		HostVpnIp:       network.ServerVPNIP(s.Config.VPN.CIDR).IP.String(),
		MetadataEnabled: !s.Config.DisableMetadata,
		IsAdmin:         user.Claims.Contains("admin"),
		AllowedIps:      allowedIPs(s.Config),
		DnsEnabled:      s.Config.DNS.Enabled,
		DnsAddress:      network.ServerVPNIP(s.Config.VPN.CIDR).IP.String(),
	}, nil
}

func allowedIPs(config *config.AppConfig) string {
	return strings.Join(config.VPN.AllowedIPs, ", ")
}
