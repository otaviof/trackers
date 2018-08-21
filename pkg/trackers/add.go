package trackers

import (
	"log"
	"strings"
)

// Add handles adding a new tracker to storage.
type Add struct {
	storage storageInterface // storage interface
	timeout int              // probe timeout in seconds
}

// Tracker includes a new tracker in storage, when addresses are informed, it will execute Probe
// against those addresses, if empty, it will use regular DNS queries. When a non-reachable Tracker
// is added it will hold status 99, and 0.0.0.0 addresses, the default behaviour.
func (a *Add) Tracker(announce string, addresses []string, dryRun bool) error {
	var tracker *Tracker
	var probe *Probe
	var reacheableAddresses []string
	var err error

	log.Printf("Adding Tracker '%s' with addresses '[%s]'", announce, strings.Join(addresses, ", "))

	if tracker, err = NewTracker(announce); err != nil {
		return err
	}

	probe = NewProbe(tracker, a.timeout)
	if len(addresses) == 0 {
		if err = probe.LookupIPs(); err != nil {
			return err
		}
	} else {
		probe.SetAddresses(addresses)
	}

	if reacheableAddresses, err = probe.ReachableAddresses(); err != nil {
		return err
	}

	if len(reacheableAddresses) == 0 {
		log.Printf("[WARN] No reachable addresses found!")
		tracker.Addresses = []string{"0.0.0.0"}
		tracker.Status = 99
	} else {
		log.Printf("Reachable addresses: '[%s]'", strings.Join(reacheableAddresses, ", "))
		tracker.Addresses = reacheableAddresses
		tracker.Status = 3
	}

	return a.storage.Write([]*Tracker{tracker})
}

// NewAdd instantiate Add, using storage and global timeout.
func NewAdd(storage storageInterface, timeout int) *Add {
	return &Add{storage: storage, timeout: timeout}
}
