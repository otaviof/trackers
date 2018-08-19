package trackers

import "log"

type Harvest struct {
	storage storageInterface
	client  clientInterface
}

// Execute runs throgh trackers coming from torrent-client and ...
func (h *Harvest) Execute() error {
	var trackers []*Tracker
	var tracker *Tracker
	var storedTrackers []*Tracker
	var storedTracker *Tracker
	var notStoredTrackers []*Tracker
	var err error

	log.Print("Harvesting running trackers...")

	if trackers, err = h.client.ListTrackers(); err != nil {
		log.Fatalf("Error on listing trackers: '%s'", err)
		return err
	}
	log.Printf("Found '%d' trackers in torrents", len(trackers))
	if storedTrackers, err = h.storage.Read(); err != nil {
		log.Fatalf("Error on reading stored trackers: '%s'", err)
		return err
	}
	log.Printf("Read '%d' trackers from storage", len(storedTrackers))

	// disanbiguanting results, checking which ones already exist
	for _, tracker = range trackers {
		var exists bool

		for _, storedTracker = range storedTrackers {
			if tracker.Announce == storedTracker.Announce {
				log.Printf("[In-Store] Tracker: '%s'", tracker.Announce)
				exists = true
				continue
			}
		}

		if !exists {
			notStoredTrackers = append(notStoredTrackers, tracker)
			log.Printf("[New-Tracker] Announce URL: '%s'", tracker.Announce)
		}
	}

	if len(notStoredTrackers) > 0 {
		log.Printf("Saving '%d' trackers...", len(notStoredTrackers))
		if err = h.storage.Write(notStoredTrackers); err != nil {
			return err
		}
	}

	return nil
}

// NewHarvest instantiate harvest, to look at the running torrents and extract and store a list of
// trackers.
func NewHarvest(storage storageInterface, client clientInterface) *Harvest {
	return &Harvest{storage: storage, client: client}
}
