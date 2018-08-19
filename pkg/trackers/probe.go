package trackers

import (
	"fmt"
	"log"
	"net"
)

type Probe struct {
	tracker *Tracker
	ipv4s   []string
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

// ReachableIPv4s check connectivity with tracker service IPv4 addresses, return the ones working.
func (p *Probe) ReachableIPv4s() ([]string, error) {
	var conn net.Conn
	var ipv4 string
	var proto string
	var reachable []string
	var err error

	if proto, err = p.protocol(); err != nil {
		return nil, err
	}

	log.Printf("Probing '%s' service...", p.tracker.Announce)
	for _, ipv4 = range p.ipv4s {
		var serviceAndPort = fmt.Sprintf("%s:%d", ipv4, p.tracker.Port)

		if conn, err = net.Dial(proto, serviceAndPort); err != nil {
			log.Printf("[%s] Service is NOT available at: '%s'", proto, serviceAndPort)
		} else {
			log.Printf("[%s] Service is available at: '%s'", proto, serviceAndPort)
			reachable = append(reachable, ipv4)
		}

		defer conn.Close()
	}

	log.Printf("Tracker is reachable over: '%s'", reachable)
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
			p.ipv4s = append(p.ipv4s, ipv4)
		}
	}

	log.Printf("Tracker '%s' IPv4 addresses: '%v'", p.tracker.Hostname, p.ipv4s)
	return nil
}

func NewProbe(tracker *Tracker) *Probe {
	return &Probe{tracker: tracker}
}
