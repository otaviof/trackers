package trackers

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Tracker defines a tracker
type Tracker struct {
	Announce  string   // tracker original url
	Hostname  string   // tracker hostname
	Port      int      // tracker port
	Protocol  string   // tracker protocol
	Addresses []string // resolved IP addresses
	Status    int      // tracker status
}

// NewTracker parses the Announce URL.
func NewTracker(announce string) (*Tracker, error) {
	var parsed *url.URL
	var hostname string
	var portSrt string
	var port int64
	var err error

	if parsed, err = url.Parse(announce); err != nil {
		return nil, err
	}
	if strings.Contains(parsed.Host, ":") {
		if hostname, portSrt, err = net.SplitHostPort(parsed.Host); err != nil {
			return nil, err
		}
	} else {
		hostname = parsed.Host
		portSrt = "80"
	}

	if portSrt == "" {
		port = 80
	} else {
		if port, err = strconv.ParseInt(portSrt, 10, 64); err != nil {
			return nil, err
		}
	}

	return &Tracker{
		Announce:  announce,
		Hostname:  hostname,
		Protocol:  parsed.Scheme,
		Port:      int(port),
		Addresses: []string{"0.0.0.0"},
		Status:    99,
	}, nil
}
