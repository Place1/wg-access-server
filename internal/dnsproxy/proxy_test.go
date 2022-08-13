package dnsproxy

import (
	"context"
	"net"
	"testing"
	"time"
)

var ffmucUpstreams, _ = net.LookupHost("dns.ffmuc.net")

func TestDNSProxy_ServeDNS(t *testing.T) {
	const listen = "[::1]:8053"

	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Second}
			return d.DialContext(ctx, network, listen)
		},
	}

	server, err := New(DNSServerOpts{
		Domain:     "",
		ListenAddr: []string{listen},
		Upstream:   ffmucUpstreams,
	})
	server.ListenAndServe()
	defer func() { _ = server.Close() }()

	if err != nil {
		t.Fatal(err)
	}

	t.Run("Reply over 1300 bytes", func(t *testing.T) {
		_, err := resolver.LookupTXT(context.Background(), "cloudflare.com.")
		if err != nil {
			t.Error(err)
			return
		}
	})
	t.Run("Reply over 1500 bytes", func(t *testing.T) {
		records, err := resolver.LookupTXT(context.Background(), "txtfill1500.go.dnscheck.tools.")
		if err != nil {
			t.Error(err)
			return
		}
		var containsBigRecord bool
		for _, r := range records {
			if len(r) >= 1500 {
				containsBigRecord = true
			}
		}
		if !containsBigRecord {
			t.Error("missing big TXT record, packet probably truncated")
		}
	})
}
