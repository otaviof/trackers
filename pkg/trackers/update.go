package trackers

import (
	"log"
	"strings"
)

// Update object to execute updates in storage, validating addresses.
type Update struct {
	storage storageInterface // interface with storage
	config  *Config          // configuration
}

// probeAddresses creates a new tracker based on original and its new addresses, executing Probe
// against new addresses, when successful returns a new Tracker object.
func (u *Update) probeAddresses(tracker *Tracker, addresses []string) (*Tracker, error) {
	var updatedTracker *Tracker
	var recheableAddresses []string
	var probe *Probe
	var err error

	if updatedTracker, err = NewTracker(tracker.Announce); err != nil {
		return nil, err
	}
	// marking tracker as updated
	updatedTracker.Status = 3

	probe = NewProbe(updatedTracker, u.config.Probe.Timeout)
	probe.SetAddresses(addresses)

	if recheableAddresses, err = probe.ReachableAddresses(); err != nil {
		return nil, err
	}

	if len(recheableAddresses) == 0 {
		log.Printf("[WARN] New addresses are not reachable for '%s'", updatedTracker.Announce)
		return nil, nil
	}

	log.Printf("[INFO] Reachable addresses: '[%s]'", strings.Join(recheableAddresses, ", "))
	// saving reachable addresses in tracker
	updatedTracker.Addresses = recheableAddresses

	return updatedTracker, nil
}

// HostnameAddress handle updates for all trackers that are matching hostname.
func (u *Update) HostnameAddress(hostname string, addresses []string, dryRun bool) error {
	var trackers []*Tracker
	var tracker *Tracker
	var err error

	if trackers, err = u.storage.Read(); err != nil {
		return err
	}

	for _, tracker = range trackers {
		var updatedTracker *Tracker

		// skipping trackers that don't match the hostname
		if hostname != tracker.Hostname {
			continue
		}

		log.Printf("Updating tracker: '%s', probing addresses: [%s]", tracker.Announce,
			strings.Join(addresses, ", "))

		if updatedTracker, err = u.probeAddresses(tracker, addresses); err != nil {
			return err
		}

		if updatedTracker != nil && !dryRun {
			if err = u.storage.Update(updatedTracker); err != nil {
				return err
			}
		}
	}

	return nil
}

// NewUpdate returns a update instnace.
func NewUpdate(storage storageInterface, config *Config) *Update {
	return &Update{storage: storage, config: config}
}
