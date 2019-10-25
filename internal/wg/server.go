package wg

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Server struct {
	client       *wgctrl.Client
	iface        string
	externalName string
	port         int
	publicKey    wgtypes.Key
	lock         sync.Mutex
}

func New(iface string, privateKey string, port int, externalName string) (*Server, error) {
	// wgctrl.New() will search for a kernel implementation
	// of wireguard, then user implementations
	// user implementations are found in /var/run/wireguard/<iface>.sock
	// this unix socket likely requires root to access
	client, err := wgctrl.New()
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wgctrl"))
	}
	key, err := wgtypes.ParseKey(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "bad private key format")
	}
	server := &Server{
		client:       client,
		iface:        iface,
		port:         port,
		externalName: externalName,
		publicKey:    key.PublicKey(),
	}
	err = server.configure(func(config *wgtypes.Config) error {
		config.PrivateKey = &key
		config.ListenPort = &port
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure wireguard - is wireguard running?")
	}
	return server, nil
}

func (s *Server) AddPeer(publicKey string, addressCIDR string) error {
	logrus.
		WithField("publicKey", publicKey).
		WithField("address", addressCIDR).
		Debugf("adding peer")
	key, err := wgtypes.ParseKey(publicKey)
	if err != nil {
		return errors.Wrapf(err, "bad public key %v", publicKey)
	}
	_, allowedIPs, err := net.ParseCIDR(addressCIDR)
	if err != nil || allowedIPs == nil {
		return errors.Wrap(err, "bad CIDR value for AllowedIPs")
	}
	if s.HasPeer(key.String()) {
		s.RemovePeer(key.String())
	}
	return s.configure(func(config *wgtypes.Config) error {
		config.ReplacePeers = false
		config.Peers = []wgtypes.PeerConfig{
			wgtypes.PeerConfig{
				PublicKey:  key,
				AllowedIPs: []net.IPNet{*allowedIPs},
			},
		}
		return nil
	})
}

func (s *Server) ListPeers() ([]wgtypes.Peer, error) {
	d, err := s.Device()
	if err != nil {
		return nil, err
	}
	return d.Peers, nil
}

func (s *Server) Peer(publicKey string) (*wgtypes.Peer, error) {
	peers, err := s.ListPeers()
	if err != nil {
		return nil, err
	}
	for _, peer := range peers {
		if peer.PublicKey.String() == publicKey {
			return &peer, nil
		}
	}
	return nil, fmt.Errorf("peer with public key '%s' not found", publicKey)
}

func (s *Server) HasPeer(publicKey string) bool {
	peers, err := s.ListPeers()
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to list peers"))
		return false
	}
	for _, peer := range peers {
		if peer.PublicKey.String() == publicKey {
			return true
		}
	}
	return false
}

func (s *Server) RemovePeer(publicKey string) error {
	logrus.WithField("publicKey", publicKey).Debug("removing peer")
	key, err := wgtypes.ParseKey(publicKey)
	if err != nil {
		return errors.Wrap(err, "bad public key")
	}
	return s.configure(func(config *wgtypes.Config) error {
		config.ReplacePeers = false
		config.Peers = []wgtypes.PeerConfig{
			wgtypes.PeerConfig{
				Remove:    true,
				PublicKey: key,
			},
		}
		return nil
	})
}

func (s *Server) PublicKey() string {
	return s.publicKey.String()
}

func (s *Server) Endpoint() string {
	return fmt.Sprintf("%s:%d", s.externalName, s.port)
}

func (s *Server) DNS() string {
	return "1.1.1.1, 8.8.8.8" // TODO: dns stuff
}

func (s *Server) Device() (*wgtypes.Device, error) {
	return s.client.Device(s.iface)
}

func (s *Server) Close() error {
	return s.client.Close()
}

func (s *Server) configure(cb func(*wgtypes.Config) error) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	next := wgtypes.Config{}
	if err := cb(&next); err != nil {
		return errors.Wrap(err, "failed to get next wireguard config")
	} else {
		return s.client.ConfigureDevice(s.iface, next)
	}
}

func trimLines(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	output := make([]string, len(lines))
	for index, line := range lines {
		output[index] = strings.TrimSpace(line)
	}
	return strings.Join(output, "\n")
}
