package dnsproxy

import (
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DNSServerOpts struct {
	Upstream []string
}

type DNSServer struct {
	server   *dns.Server
	client   *dns.Client
	cache    *cache.Cache
	upstream []string
}

func New(opts DNSServerOpts) (*DNSServer, error) {
	if len(opts.Upstream) == 0 {
		return nil, errors.New("at least 1 upstream dns server is required for the dns proxy server to function")
	}

	addr := ":53"
	logrus.Infof("starting dns server on %s with upstreams: %s", addr, strings.Join(opts.Upstream, ", "))

	dnsServer := &DNSServer{
		server: &dns.Server{
			Addr: addr,
			Net:  "udp",
		},
		client: &dns.Client{
			SingleInflight: true,
			Timeout:        5 * time.Second,
		},
		cache:    cache.New(10*time.Minute, 10*time.Minute),
		upstream: opts.Upstream,
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
			HandleFailed(w, r)
			return
		}
		m.SetReply(r)
		err = w.WriteMsg(m)
		if err != nil {
			logrus.Errorf("failed write response for client with error: %s\n%s", err.Error(), r)
			return
		}
	default:
		m := &dns.Msg{}
		m.SetReply(r)
		err := w.WriteMsg(m)
		if err != nil {
			logrus.Errorf("failed write response for client with error: %s\n%s", err.Error(), r)
			return
		}
	}

}

// HandleFailed returns a HandlerFunc that returns SERVFAIL for every request it gets.
func HandleFailed(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetRcode(r, dns.RcodeServerFailure)
	// does not matter if this write fails
	_ = w.WriteMsg(m)
}

func (d *DNSServer) Lookup(m *dns.Msg) (*dns.Msg, error) {
	key := makekey(m)

	// check the cache first
	if item, found := d.cache.Get(key); found {
		logrus.Debugf("dns cache hit %s", prettyPrintMsg(m))
		return item.(*dns.Msg), nil
	}

	// fallback to upstream exchange
	// TODO disable upstream after certain amount of failures?
	var response *dns.Msg
	var firstErr error
	for _, upstream := range d.upstream {
		resp, _, err := d.client.Exchange(m, net.JoinHostPort(upstream, "53"))
		if err != nil && firstErr == nil {
			logrus.Warnf(errors.Wrap(err, fmt.Sprintf("DNS lookup failed for upstream %s", upstream)).Error())
			firstErr = err
		}
		response = resp
	}
	if response == nil {
		return nil, firstErr
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
