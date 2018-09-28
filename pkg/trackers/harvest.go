package trackers

import "log"

// Harvest query torrent client for trackers, saving the new ones.
type Harvest struct {
	storage storageInterface
	client  clientInterface
}

// Execute runs throgh trackers coming from torrent-client and ...
func (h *Harvest) Execute(dryRun bool) error {
	var trackers []*Tracker
	var torrents []*Torrent
	var tracker *Tracker
	var storedTrackers []*Tracker
	var storedTracker *Tracker
	var notStoredTrackers []*Tracker
	var err error

	log.Print("Harvesting running trackers...")

	if torrents, err = h.client.ListTorrents([]int{}); err != nil {
		log.Fatalf("Error on listing torrents: '%s'", err)
		return err
	}

	trackers = h.client.ListTrackers(torrents)
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
		log.Printf("Saving '%d' trackers... (dry-run: %v)", len(notStoredTrackers), dryRun)
		if !dryRun {
			if err = h.storage.Write(notStoredTrackers); err != nil {
				return err
			}
		}
	}

	return nil
}

// NewHarvest instantiate harvest, to look at the running torrents and extract and store a list of
// trackers.
func NewHarvest(storage storageInterface, client clientInterface) *Harvest {
	return &Harvest{storage: storage, client: client}
}
