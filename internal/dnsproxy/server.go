package dnsproxy

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DNSServerOpts struct {
	Domain     string
	ListenAddr []string
	Upstream   []string
}

type DNSServer struct {
	servers []*dns.Server
	proxy   *DNSProxy
	auth    *DNSAuth
}

// New returns a pointer to a DNSServer configured using opts DNSServerOpts.
// The returned server needs to be started using DNSServer.ListenAndServe()
func New(opts DNSServerOpts) (*DNSServer, error) {
	if len(opts.Upstream) == 0 {
		return nil, errors.New("At least 1 upstream dns server is required for the dns proxy server to function")
	}

	dnsServer := &DNSServer{
		servers: []*dns.Server{},
		proxy: &DNSProxy{
			udpClient: &dns.Client{
				SingleInflight: true,
				Timeout:        5 * time.Second,
			},
			tcpClient: &dns.Client{
				Net:            "tcp",
				SingleInflight: true,
				Timeout:        5 * time.Second,
			},
			cache:    cache.New(10*time.Minute, 10*time.Minute),
			upstream: opts.Upstream,
		},
		auth: &DNSAuth{
			Domain:   dns.Fqdn(opts.Domain),
			zoneLock: new(sync.RWMutex),
		},
	}

	// Send queries for VPN search domain to the authoritative server and everything else to the proxy
	serveMux := dns.NewServeMux()
	if opts.Domain != "" {
		serveMux.Handle(dnsServer.auth.Domain, dnsServer.auth)
	}
	serveMux.Handle(".", dnsServer.proxy)

	// Create one UDP and one TCP server per listen address
	for _, addr := range opts.ListenAddr {
		udpServer := &dns.Server{
			Addr: addr,
			Net:  "udp",
			// https://dnsflagday.net/2020/
			UDPSize: 1232,
			Handler: serveMux,
		}
		tcpServer := &dns.Server{
			Addr:    addr,
			Net:     "tcp",
			Handler: serveMux,
		}
		dnsServer.servers = append(dnsServer.servers, udpServer)
		dnsServer.servers = append(dnsServer.servers, tcpServer)
	}

	return dnsServer, nil
}

// ListenAndServe starts the DNSServer and waits until all listeners are up.
func (d *DNSServer) ListenAndServe() {
	var sb strings.Builder
	for i, s := range d.servers {
		sb.WriteString(s.Addr)
		sb.WriteString("/")
		sb.WriteString(s.Net)
		if i < len(d.servers)-1 {
			sb.WriteString(", ")
		}
	}

	logrus.Infof("Starting DNS server on %s with upstreams: %s", sb.String(), strings.Join(d.proxy.upstream, ", "))

	var wg sync.WaitGroup

	for _, server := range d.servers {
		wg.Add(1)
		server.NotifyStartedFunc = func() {
			wg.Done()
		}
		go func(server *dns.Server) {
			if err := server.ListenAndServe(); err != nil {
				logrus.Error(errors.Errorf("Failed to start DNS server on %s/%s: %s", server.Addr, server.Net, err))
				wg.Done()
			}
		}(server)
	}

	wg.Wait()
}

func (d *DNSServer) Close() error {
	var firstErr error
	for _, server := range d.servers {
		err := server.Shutdown()
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		return errors.Wrap(firstErr, "DNS server shutdown failed")
	}
	return nil
}

func (d *DNSServer) PushAuthZone(zone Zone) {
	d.auth.PushZone(zone)
}

// HandleFailed is a HandlerFunc that returns SERVFAIL for every request it gets.
func HandleFailed(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetRcode(r, dns.RcodeServerFailure)
	// does not matter if this write fails
	_ = w.WriteMsg(m)
}

func makekey(m *dns.Msg) string {
	q := m.Question[0]
	// Definitely not standard compliant, but better than nothing
	// A proper security-aware caching stub resolver always sets the DO bit to upstream
	// and customizes the client response based on the query, leaving out DNSSEC RRs if DO wasn't set
	var flags uint8
	if m.AuthenticatedData {
		flags |= 1 << 0
	}
	if m.CheckingDisabled {
		flags |= 1 << 1
	}
	if opt := m.IsEdns0(); opt != nil && opt.Do() {
		flags |= 1 << 2
	}
	return fmt.Sprintf("%s:%d:%d:%d", q.Name, q.Qtype, q.Qclass, flags)
}

func prettyPrintMsg(m *dns.Msg) string {
	if len(m.Question) > 0 {
		return fmt.Sprintf("DNS query for: %s", makekey(m))
	}
	return m.String()
}
