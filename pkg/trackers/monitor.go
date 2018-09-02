package trackers

import (
	"log"
	"strings"
	"sync"
)

// Monitor monitor instance
type Monitor struct {
	storage storageInterface
	config  *Config
	workers int
}

var wg sync.WaitGroup

// probedTracker execute the probes and name resolver against informed tracker, it returns a new
// tracker instance with up-to-date fields regarding its functional status.
func (m *Monitor) probeTracker(tracker *Tracker) (*Tracker, error) {
	var probe = NewProbe(tracker, m.config.Probe.Timeout)
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

// worker combine in and out channels to probe trackers and send back their funcitnoal status.
func (m *Monitor) worker(id int, in <-chan *Tracker, out chan<- *Tracker) {
	var tracker *Tracker
	var err error

	for tracker = range in {
		var probed *Tracker
		log.Printf("[worker: %d] Probing tracker '%s'", id, tracker.Announce)
		if probed, err = m.probeTracker(tracker); err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		out <- probed
	}
}

// readTrackers read trackers from storage, skipping trackers on status 3.
func (m *Monitor) readTrackers() ([]*Tracker, error) {
	var trackers []*Tracker
	var tracker *Tracker
	var toProbe []*Tracker
	var err error

	if trackers, err = m.storage.Read(); err != nil {
		return nil, err
	}

	for _, tracker = range trackers {
		if tracker.Status == 3 {
			continue
		}
		toProbe = append(toProbe, tracker)
	}

	return toProbe, nil
}

// Inspect load trackers from storage and start to probe if services are up, the results are stored
// back via storage interface.
func (m *Monitor) Inspect(dryRun bool) error {
	var trackers []*Tracker
	var tracker *Tracker
	var in chan *Tracker
	var out chan *Tracker
	var i int
	var err error

	if trackers, err = m.readTrackers(); err != nil {
		return err
	}

	in = make(chan *Tracker, len(trackers)-1)
	out = make(chan *Tracker, len(trackers)-1)

	for i = 0; i < m.workers; i++ {
		go m.worker(i, in, out)
	}

	for _, tracker = range trackers {
		in <- tracker
	}
	close(in)
	wg.Wait()

	for i = 0; i < len(trackers); i++ {
		tracker = <-out
		log.Printf("[result] Tracker: announce='%s', status='%d', addresses='[%s]'",
			tracker.Announce, tracker.Status, strings.Join(tracker.Addresses, ", "))

		if !dryRun {
			if err = m.storage.Update(tracker); err != nil {
				return err
			}
		}
	}
	defer close(out)

	return nil
}

// NewMonitor instantiate a monitor object, requires storage interface.
func NewMonitor(storage storageInterface, config *Config, workers int) *Monitor {
	return &Monitor{storage: storage, config: config, workers: workers}
}
