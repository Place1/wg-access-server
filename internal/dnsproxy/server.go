package dnsproxy

import (
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/docker/libnetwork/resolvconf"
	"github.com/docker/libnetwork/types"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DNSServer struct {
	server   *dns.Server
	client   *dns.Client
	cache    *cache.Cache
	upstream []string
}

func New() (*DNSServer, error) {

	upstream := []string{}

	if r, err := resolvconf.Get(); err == nil {
		upstream = resolvconf.GetNameservers(r.Content, types.IPv4)
	}

	if len(upstream) == 0 {
		logrus.Warn("failed to get nameservers from /etc/resolv.conf defaulting to 1.1.1.1 for DNS instead")
		upstream = append(upstream, "1.1.1.1")
	}

	logrus.Infof("starting dns server with upstreams: %s", strings.Join(upstream, ", "))

	dnsServer := &DNSServer{
		server: &dns.Server{
			Addr: "0.0.0.0:53",
			Net:  "udp",
		},
		client: &dns.Client{
			SingleInflight: true,
			Timeout:        5 * time.Second,
		},
		cache:    cache.New(10*time.Minute, 10*time.Minute),
		upstream: upstream,
	}
	dnsServer.server.Handler = dnsServer

	go func() {
		if err := dnsServer.server.ListenAndServe(); err != nil {
			logrus.Error(errors.Wrap(err, "failed to start dns server"))
		}
	}()

	return dnsServer, nil
}

func (d *DNSServer) Close() error {
	return d.server.Shutdown()
}

func (d *DNSServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("dns server panic handled: %v\n%s", err, string(debug.Stack()))
			dns.HandleFailed(w, r)
		}
	}()

	logrus.Debugf("dns query: %s", prettyPrintMsg(r))

	switch r.Opcode {
	case dns.OpcodeQuery:
		m, err := d.Lookup(r)
		if err != nil {
			logrus.Errorf("failed lookup record with error: %s\n%s", err.Error(), r)
			dns.HandleFailed(w, r)
			return
		}
		m.SetReply(r)
		w.WriteMsg(m)
	default:
		m := &dns.Msg{}
		m.SetReply(r)
		w.WriteMsg(m)
	}

}

func (d *DNSServer) Lookup(m *dns.Msg) (*dns.Msg, error) {
	key := makekey(m)

	// check the cache first
	if item, found := d.cache.Get(key); found {
		logrus.Debugf("dns cache hit %s", prettyPrintMsg(m))
		return item.(*dns.Msg), nil
	}

	// fallback to upstream exchange
	response, _, err := d.client.Exchange(m, net.JoinHostPort(d.upstream[0], "53"))
	if err != nil {
		return nil, err
	}

	if len(response.Answer) > 0 {
		ttl := time.Duration(response.Answer[0].Header().Ttl) * time.Second
		logrus.Debugf("caching dns response for %s for %v seconds", prettyPrintMsg(m), ttl)
		d.cache.Set(key, response, ttl)
	}

	return response, nil
}

func makekey(m *dns.Msg) string {
	q := m.Question[0]
	return fmt.Sprintf("%s:%d:%d", q.Name, q.Qtype, q.Qclass)
}

func prettyPrintMsg(m *dns.Msg) string {
	if len(m.Question) > 0 {
		return fmt.Sprintf("dns query for: %s", makekey(m))
	}
	return m.String()
}
