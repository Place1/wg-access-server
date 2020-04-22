package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/place1/wg-access-server/internal/network"

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
		AllowedIps:      allowedIPs(s.Config),
		DnsEnabled:      *s.Config.DNS.Enabled,
		DnsAddress:      fmt.Sprintf("%s:%d", network.ServerVPNIP(s.Config.VPN.CIDR).IP.String(), s.Config.DNS.Port),
	}, nil
}

func allowedIPs(config *config.AppConfig) string {
	if config.VPN.Rules == nil {
		return "0.0.0.0/1, 128.0.0.0/1, ::/0"
	}

	allowed := []string{}

	if *config.DNS.Enabled {
		allowed = append(allowed, network.ServerVPNIP(config.VPN.CIDR).IP.String())
	}

	if config.VPN.Rules.AllowVPNLAN {
		allowed = append(allowed, config.VPN.CIDR)
	}

	if config.VPN.Rules.AllowServerLAN {
		allowed = append(allowed, network.ServerLANSubnets...)
	}

	if config.VPN.Rules.AllowInternet {
		allowed = append(allowed, network.PublicInternetSubnets...)
	}

	allowed = append(allowed, config.VPN.Rules.AllowedNetworks...)

	return strings.Join(allowed, ", ")
}
