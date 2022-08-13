package dnsproxy

import (
	"fmt"
	"net"
	"runtime/debug"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DNSProxy struct {
	udpClient *dns.Client
	tcpClient *dns.Client
	cache     *cache.Cache
	upstream  []string
}

// ServeDNS is called by the mux from the listening servers.
func (d *DNSProxy) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("dns server panic handled: %v\n%s", err, string(debug.Stack()))
			dns.HandleFailed(w, r)
		}
	}()

	logrus.Debugf("dns query: %s", prettyPrintMsg(r))

	switch r.Opcode {
	case dns.OpcodeQuery:
		// Remove EDNS0 Client Subnet information as we don't handle them in the cache
		purgeECS(r)
		outQuery := r.Copy()
		// Set EDNS BufSize for forwarding to upstream
		ensureEDNS0BufSize(outQuery)
		m, err := d.Lookup(outQuery)
		if err != nil {
			logrus.Errorf("failed lookup record with error: %s\n%s", err.Error(), r)
			HandleFailed(w, r)
			return
		}
		m.SetReply(r)
		truncateIfRequired(m, r, w.RemoteAddr().Network())
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

// Lookup first checks the cache for a matching response, and if unsuccessful queries the upstream resolvers.
func (d *DNSProxy) Lookup(m *dns.Msg) (*dns.Msg, error) {
	key := makekey(m)

	// check the cache first
	if item, found := d.cache.Get(key); found {
		logrus.Debugf("dns cache hit %s", prettyPrintMsg(m))
		return item.(*dns.Msg).Copy(), nil
	}

	// fallback to upstream exchange
	// TODO disable upstream after certain amount of failures?
	var response *dns.Msg
	var firstErr error
	for _, upstream := range d.upstream {
		target := net.JoinHostPort(upstream, "53")
		resp, _, err := d.udpClient.Exchange(m, target)
		if err != nil && firstErr == nil {
			logrus.Warnf(errors.Wrap(err, fmt.Sprintf("DNS lookup failed for upstream %s", upstream)).Error())
			firstErr = err
		} else if err == nil {
			// Retry truncated responses over TCP
			if resp.Truncated {
				resp, _, err = d.tcpClient.Exchange(m, target)
				if err != nil && firstErr == nil {
					logrus.Warnf(errors.Wrap(err, fmt.Sprintf("DNS lookup failed over TCP for upstream %s", upstream)).Error())
					firstErr = err
					continue
				}
			}
			response = resp
			break
		}
	}
	if response == nil {
		return nil, firstErr
	}

	if len(response.Answer) > 0 {
		ttl := time.Duration(response.Answer[0].Header().Ttl) * time.Second
		logrus.Debugf("caching dns response for %s for %v seconds", prettyPrintMsg(m), ttl)
		d.cache.Set(key, response, ttl)
	}

	return response.Copy(), nil
}

func purgeECS(m *dns.Msg) {
	if opt := m.IsEdns0(); opt != nil {
		for i, option := range opt.Option {
			if option.Option() == dns.EDNS0SUBNET {
				opt.Option = append(opt.Option[:i], opt.Option[i+1:]...)
			}
		}
	}
}

func ensureEDNS0BufSize(m *dns.Msg) {
	if opt := m.IsEdns0(); opt != nil {
		opt.SetUDPSize(1232)
	} else {
		m.SetEdns0(1232, false)
	}
}

func truncateIfRequired(response *dns.Msg, original *dns.Msg, transport string) {
	size := dns.MinMsgSize
	if transport == "tcp" {
		size = dns.MaxMsgSize
	} else if opt := original.IsEdns0(); opt != nil {
		size = int(opt.UDPSize())
	}
	logrus.Debugf("truncating to %d", size)
	response.Truncate(size)
}
