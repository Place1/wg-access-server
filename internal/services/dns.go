package services

import (
	"fmt"
	"net"
	"runtime/debug"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DNSServer struct {
	server   *dns.Server
	client   *dns.Client
	upstream []string
}

func NewDNSServer(upstream []string) (*DNSServer, error) {
	logrus.Infof("starting dns server")
	server := &dns.Server{Addr: "0.0.0.0:53", Net: "udp"}
	client := &dns.Client{
		SingleInflight: true,
	}
	dnsServer := &DNSServer{server, client, upstream}
	server.Handler = dnsServer
	go func() {
		if err := server.ListenAndServe(); err != nil {
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
	switch r.Opcode {
	case dns.OpcodeQuery:
		if logrus.GetLevel() == logrus.DebugLevel {
			// log behind a condition to ensure we don't call prettyPrintMsg
			// when the log level would filter out the message anyway
			logrus.Debugf("dns query: %s", prettyPrintMsg(r))
		}
		m, err := d.Lookup(r)
		if err != nil {
			logrus.Errorf("failed lookup record %s with error: %s\n", r, err.Error())
			m.SetReply(r)
			w.WriteMsg(m)
			return
		}
		m.SetReply(r)
		w.WriteMsg(m)
		return
	}
}

func (d *DNSServer) Lookup(m *dns.Msg) (*dns.Msg, error) {
	// TODO: add support for caching
	response, _, err := d.client.Exchange(m, net.JoinHostPort(d.upstream[0], "53"))
	if err != nil {
		return nil, err
	}
	return response, nil
}

func prettyPrintMsg(m *dns.Msg) string {
	if len(m.Question) > 0 {
		return fmt.Sprintf("dns query for: %s", m.Question[0].Name)
	}
	return m.String()
}
