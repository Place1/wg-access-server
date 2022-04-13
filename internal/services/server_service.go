package services

import (
	"context"
	"strings"

	"github.com/freifunkMUC/wg-access-server/internal/config"
	"github.com/freifunkMUC/wg-access-server/internal/network"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"
	"github.com/freifunkMUC/wg-access-server/proto/proto"

	"github.com/freifunkMUC/wg-embed/pkg/wgembed"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerService struct {
	proto.UnimplementedServerServer
	Config *config.AppConfig
	Wg     wgembed.WireGuardInterface
}

func (s *ServerService) Info(ctx context.Context, req *proto.InfoReq) (*proto.InfoRes, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	host := s.Config.ExternalHost
	if strings.Contains(host, ":") {
		if !strings.HasPrefix(host, "[") {
			host = "[" + host
		}
		if !strings.HasSuffix(host, "]") {
			host = host + "]"
		}
	}

	publicKey, err := s.Wg.PublicKey()
	if err != nil {
		ctxlogrus.Extract(ctx).Error(err)
		return nil, status.Errorf(codes.Internal, "failed to get public key")
	}

	vpnip, vpnipv6, err := network.ServerVPNIPs(s.Config.VPN.CIDR, s.Config.VPN.CIDRv6)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get server IPs")
	}
	dnsAddress := network.StringJoinIPs(vpnip, vpnipv6)

	var hostVPNIP string
	if vpnip != nil {
		hostVPNIP = vpnip.IP.String()
	} else {
		hostVPNIP = ""
	}

	return &proto.InfoRes{
		Host:      stringValue(&host),
		PublicKey: publicKey,
		Port:      int32(s.Config.WireGuard.Port),
		// TODO IPv6 what is HostVpnIp used for, do we need HostVpnIpv6 as well?
		HostVpnIp:       hostVPNIP,
		MetadataEnabled: !s.Config.DisableMetadata,
		IsAdmin:         user.Claims.Contains("admin"),
		AllowedIps:      allowedIPs(s.Config),
		DnsEnabled:      s.Config.DNS.Enabled,
		DnsAddress:      dnsAddress,
	}, nil
}

func allowedIPs(config *config.AppConfig) string {
	return strings.Join(config.VPN.AllowedIPs, ", ")
}
