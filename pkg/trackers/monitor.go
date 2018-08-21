package trackers

import "log"

// Monitor monitor instance
type Monitor struct {
	storage storageInterface
	timeout int64
}

// probedTracker execute the probes and name resolver against informed tracker, it returns a new
// tracker instance with up-to-date fields regarding its functional status.
func (m *Monitor) probeTracker(tracker *Tracker) (*Tracker, error) {
	var probe = NewProbe(tracker, m.timeout)
	var probedTracker *Tracker
	var addresses []string
	var err error

	if probedTracker, err = NewTracker(tracker.Announce); err != nil {
		return nil, err
	}
	probedTracker.Status = 2
	probedTracker.Addresses = []string{"0.0.0.0"}

	// resolving tracker ip addresses
	if err = probe.LookupAddresses(); err != nil {
		log.Printf("[DNS] Can't rsolve tracker's address: '%s'", err)
		return probedTracker, nil
	}

	// executing probe, returning which addreses are reachable
	if addresses, err = probe.ReachableAddresses(); err != nil {
		return nil, err
	}

	if len(addresses) > 0 {
		probedTracker.Status = 0
		probedTracker.Addresses = addresses
	} else {
		probedTracker.Status = 1
	}

	return probedTracker, nil
}

// Inspect load trackers from storage and start to probe if services are up, the results are stored
// back via storage interface.
func (m *Monitor) Inspect(dryRun bool) error {
	var trackers []*Tracker
	var tracker *Tracker
	var err error

	if trackers, err = m.storage.Read(); err != nil {
		return err
	}

	for _, tracker = range trackers {
		var probedTracker *Tracker

		log.Printf("Tracker: '%s' (hostname: '%s')", tracker.Announce, tracker.Hostname)
		if tracker.Status == 3 {
			log.Printf("Skipping tracker in status '%d'", tracker.Status)
			continue
		}

		if probedTracker, err = m.probeTracker(tracker); err != nil {
			return err
		}

		log.Printf("Tracker: status='%d', addresses='%v'",
			probedTracker.Status, probedTracker.Addresses)

		if !dryRun {
			if err = m.storage.Update(probedTracker); err != nil {
				return err
			}
		}
	}

	return nil
}

// NewMonitor instantiate a monitor object, requires storage interface.
func NewMonitor(storage storageInterface, timeout int64) *Monitor {
	return &Monitor{storage: storage, timeout: timeout}
}
