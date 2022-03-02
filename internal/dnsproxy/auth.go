package dnsproxy

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const authoritativeTTL = 300

type ZoneKey struct{ Owner, Name string }
type Zone map[ZoneKey][]net.IP

type DNSAuth struct {
	Domain string
	// zone is a map of users to a map of device names to IP addresses
	// Lock zoneLock before accessing
	zone     Zone
	zoneLock *sync.RWMutex
}

func (d *DNSAuth) PushZone(zone Zone) {
	logrus.Debugln("pushing new auth zone")
	d.zoneLock.Lock()
	d.zone = zone
	d.zoneLock.Unlock()
}

func (d *DNSAuth) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	logrus.Debugf("auth dns query: %s", prettyPrintMsg(r))

	switch r.Opcode {
	case dns.OpcodeQuery:
		m, err := d.Lookup(r)
		if err != nil {
			logrus.Errorf("failed lookup record with error: %s\n%s", err.Error(), r)
			HandleFailed(w, r)
			return
		}
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

func (d *DNSAuth) Lookup(m *dns.Msg) (*dns.Msg, error) {
	if len(m.Question) != 1 {
		return nil, errors.New("only one question per packet allowed")
	}
	question := m.Question[0]
	qname := question.Name
	if question.Qclass != dns.ClassINET {
		return nil, errors.New("only class INET allowed")
	}

	deviceAndOwner := strings.TrimSuffix(qname, d.Domain)
	parts := dns.SplitDomainName(deviceAndOwner)

	response := new(dns.Msg)
	response.Authoritative = true
	var addresses []net.IP

	if parts == nil {
		// Query for the search domain itself, return server address
		addresses = d.getDevice("", "")
	} else if len(parts) < 2 {
		// Do not send NXDOMAIN because the device owner could exist (RFC 8020)
		return response.SetReply(m), nil
	} else {
		device, owner := parts[len(parts)-2], parts[len(parts)-1]

		if len(parts) > 2 {
			// Insert a CNAME to <device>.<owner>.<domain>
			target := fmt.Sprintf("%s.%s.%s", device, owner, d.Domain)
			rr, err := newRR(qname, question.Qclass, dns.TypeCNAME, target)
			if err != nil {
				return nil, err
			}
			response.Answer = append(response.Answer, rr)
			qname = target
		}

		addresses = d.getDevice(owner, device)
		if len(addresses) == 0 {
			// The requested device does not exist
			// The RCODE is always based on the final name in an CNAME chain (RFC 6604)
			return response.SetRcode(m, dns.RcodeNameError), nil
		}
	}

	// Figure out which addresses to send
	for _, addr := range addresses {
		if question.Qtype == dns.TypeAAAA || question.Qtype == dns.TypeANY {
			if addr.To4() == nil {
				rr, err := newRR(qname, question.Qclass, dns.TypeAAAA, addr.To16().String())
				if err == nil {
					response.Answer = append(response.Answer, rr)
				}
			}
		}
		if question.Qtype == dns.TypeA || question.Qtype == dns.TypeANY {
			if addr.To4() != nil {
				rr, err := newRR(qname, question.Qclass, dns.TypeA, addr.To4().String())
				if err == nil {
					response.Answer = append(response.Answer, rr)
				}
			}
		}
	}

	response.SetReply(m)
	return response, nil
}

func (d *DNSAuth) getDevice(owner, device string) []net.IP {
	d.zoneLock.RLock()
	defer d.zoneLock.RUnlock()
	return d.zone[ZoneKey{owner, device}]
}

// newRR creates a new resource record from the arguments
func newRR(qname string, qclass uint16, qtype uint16, data string) (dns.RR, error) {
	return dns.NewRR(fmt.Sprintf("%s %d %s %s %s", qname, authoritativeTTL, dns.ClassToString[qclass], dns.TypeToString[qtype], data))
}
