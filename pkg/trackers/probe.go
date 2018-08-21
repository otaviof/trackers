package trackers

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	"github.com/miekg/dns"
)

// Probe checks if a tracker is responding.
type Probe struct {
	tracker     *Tracker // tracker in probe
	timeout     int64    // dial timeout in seconds
	addresses   []string // addresses to probe
	nameservers []string // nameservers to use on lookup
}

// protocol defines which protocol (TCP or UDP) should be used against tracker service.
func (p *Probe) protocol() (string, error) {
	var proto string

	switch proto = p.tracker.Protocol; proto {
	case "tcp":
		proto = "tcp"
	case "udp":
		proto = "udp"
	case "http":
		proto = "tcp"
	case "https":
		proto = "tcp"
	default:
		return "", fmt.Errorf("can't define protocol probe for '%s'", proto)
	}

	return proto, nil
}

// ReachableAddresses check connectivity with tracker service IPv4 addresses, return the ones working.
func (p *Probe) ReachableAddresses() ([]string, error) {
	var conn net.Conn
	var ipv4 string
	var proto string
	var reachable []string
	var err error

	if proto, err = p.protocol(); err != nil {
		return nil, err
	}

	for _, ipv4 = range p.addresses {
		var serviceAndPort = fmt.Sprintf("%s:%d", ipv4, p.tracker.Port)

		if ipv4 == "127.0.0.1" || ipv4 == "0.0.0.0" {
			log.Printf("Skipping address '%s'", ipv4)
			continue
		}

		// trying to dial tracker service, with timeout
		if conn, err = net.DialTimeout(proto, serviceAndPort, p.getTimeout()); err != nil {
			log.Printf("[%s] Service is NOT available at: '%s' (%s)", proto, serviceAndPort, err)
		} else {
			// if no error returned assuming it's reachable
			log.Printf("[%s] Service is available at: '%s'", proto, serviceAndPort)
			reachable = append(reachable, ipv4)
			defer conn.Close()
		}
	}

	return reachable, nil
}

// getTimeout as time.Duration
func (p *Probe) getTimeout() time.Duration {
	return time.Duration(p.timeout) * time.Second
}

// dnsQuery executes the actual dns query against nameservers.
func (p *Probe) dnsQuery(client *dns.Client) ([]string, error) {
	var msg *dns.Msg
	var resp *dns.Msg
	var addresses []string
	var ns string
	var err error

	msg = &dns.Msg{
		MsgHdr:   dns.MsgHdr{RecursionDesired: true},
		Question: make([]dns.Question, 1),
	}
	msg.SetQuestion(dns.Fqdn(p.tracker.Hostname), dns.TypeA)

	for _, ns = range p.nameservers {
		if resp, _, err = client.Exchange(msg, ns); err != nil {
			return nil, err
		}

		for _, answer := range resp.Answer {
			switch t := answer.(type) {
			case *dns.A:
				addresses = append(addresses, t.A.String())
			}
		}
	}

	return addresses, nil
}

// LookupAddresses uses dns-query results, disanbiguate the entries and only keep ipv4 addresses.
func (p *Probe) LookupAddresses() error {
	var client *dns.Client
	var addresses []string
	var address string
	var err error

	client = &dns.Client{
		Net:       "tcp-tls",
		Timeout:   p.getTimeout(),
		TLSConfig: &tls.Config{InsecureSkipVerify: false},
	}

	if addresses, err = p.dnsQuery(client); err != nil {
		return err
	}

	for _, address = range addresses {
		var match bool

		if match, _ = regexp.MatchString(`^\d+\.\d+\.\d+\.\d+$`, address); !match {
			continue
		}
		if !stringSliceContains(p.addresses, address) {
			p.addresses = append(p.addresses, address)
		}
	}

	return nil
}

// SetAddresses instead of using LookupIPs, you can set the addresses.
func (p *Probe) SetAddresses(addresses []string) {
	p.addresses = addresses
}

// NewProbe instantiate a probe object for a specific tracker.
func NewProbe(tracker *Tracker, timeout int64) *Probe {
	var nameservers = []string{"1.1.1.1:853", "1.0.0.1:853"}
	return &Probe{tracker: tracker, timeout: timeout, nameservers: nameservers}
}
