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
	client   *dns.Client
	cache    *cache.Cache
	upstream []string
}

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

func (d *DNSProxy) Lookup(m *dns.Msg) (*dns.Msg, error) {
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
		} else if err == nil {
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

	return response, nil
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
