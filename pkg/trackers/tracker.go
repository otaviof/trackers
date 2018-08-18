package trackers

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Tracker defines a tracker
type Tracker struct {
	Announce string        // tracker original url
	Hostname string        // tracker hostname
	Port     int           // tracker port
	Protocol string        // tracker protocol
	Status   TrackerStatus // tracker status
}

// TrackerStatus exposes a given tracker funcitonal status
type TrackerStatus struct {
	Reachable      bool // tracker is reachable
	HistoryAddress bool // tracker address is coming from history, not an updated entry
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
		Announce: announce,
		Hostname: hostname,
		Protocol: parsed.Scheme,
		Port:     int(port),
	}, nil
}
