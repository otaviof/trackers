package trackers

// Update holds the interfaces needed to updated trackers in a torrent.
type Update struct {
	storage storageInterface
	client  clientInterface
}

// Execute contains the steps to list torrents and update trackers, accordingly to informed
// parameters, tracker and torrent statuses.
func (u *Update) Execute(trackerStatuses []int, torrentStatues []int, dryRun bool) error {
	var torrents []*Torrent
	var trackersToAdd []*Tracker
	var trackers []*Tracker
	var tracker *Tracker
	var err error

	if torrents, err = u.client.ListTorrents(torrentStatues); err != nil {
		return err
	}

	if trackers, err = u.storage.Read(); err != nil {
		return err
	}

	for _, tracker = range trackers {
		// skipping trackers that don't match desired statuses
		if !intSliceEq(trackerStatuses, []int{-1}) &&
			!intSliceContains(trackerStatuses, tracker.Status) {
			continue
		}
		trackersToAdd = append(trackersToAdd, tracker)
	}

	if err = u.client.UpdateTrackers(torrents, trackersToAdd, dryRun); err != nil {
		return err
	}

	return nil
}

// NewUpdate instantiate Update.
func NewUpdate(storage storageInterface, client clientInterface) *Update {
	return &Update{storage: storage, client: client}
}
