package trackers

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// Probe checks if a tracker is responding.
type Probe struct {
	tracker   *Tracker // tracker in probe
	timeout   int      // dial timeout in seconds
	addresses []string // addresses to probe
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

func (p *Probe) getTimeout() time.Duration {
	return time.Duration(p.timeout) * time.Second
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

// LookupIPs resolve tracker's hostname IPv4 addresses
func (p *Probe) LookupIPs() error {
	var IPs []net.IP
	var IP net.IP
	var err error

	if IPs, err = net.LookupIP(p.tracker.Hostname); err != nil {
		return err
	}

	for _, IP = range IPs {
		var ipv4 = IP.To4().String()

		// ignoring "<nil>", since it represents tha IPv5 address
		if ipv4 != "" && ipv4 != "<nil>" {
			p.addresses = append(p.addresses, ipv4)
		}
	}

	log.Printf("Tracker '%s' IPv4 addresses: '[%s]'",
		p.tracker.Hostname, strings.Join(p.addresses, ", "))
	return nil
}

// SetAddresses instead of using LookupIPs, you can set the addresses.
func (p *Probe) SetAddresses(addresses []string) {
	p.addresses = addresses
}

// NewProbe instantiate a probe object for a specific tracker.
func NewProbe(tracker *Tracker, timeout int) *Probe {
	return &Probe{tracker: tracker, timeout: timeout}
}
