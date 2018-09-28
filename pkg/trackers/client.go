package trackers

import (
	"log"

	trans "github.com/odwrtw/transmission"
)

type clientInterface interface {
	ListTorrents(torrentStatus []int) ([]*Torrent, error)
	ListTrackers(torrents []*Torrent) []*Tracker
	UpdateTrackers(torrents []*Torrent, trackers []*Tracker, dryRun bool) error
}

// Client tranmission torrent client.
type Client struct {
	t *trans.Client
}

// ListTrackers based on torrents, extracts a unique list.
func (c *Client) ListTrackers(torrents []*Torrent) []*Tracker {
	var torrent *Torrent
	var tracker *Tracker
	var trackers []*Tracker
	var trackerMap = make(map[string]*Tracker)
	var exists bool

	for _, torrent = range torrents {
		for _, tracker = range torrent.Trackers {
			if _, exists = trackerMap[tracker.Announce]; !exists {
				trackerMap[tracker.Announce] = tracker
			}
		}
	}

	for _, tracker = range trackerMap {
		trackers = append(trackers, tracker)
	}

	return trackers
}

// extractedTorrent transforms a upstream transmission torrent in a local torrent object.
func (c *Client) extractTorrent(torrent *trans.Torrent) (*Torrent, error) {
	var trackerStats trans.TrackerStats
	var extractedTrackers []*Tracker
	var extractedTracker *Tracker
	var err error

	for _, trackerStats = range *torrent.TrackerStats {
		if extractedTracker, err = NewTracker(trackerStats.Announce); err != nil {
			return nil, err
		}
		extractedTrackers = append(extractedTrackers, extractedTracker)
	}

	return &Torrent{ID: torrent.ID, Name: torrent.Name, Trackers: extractedTrackers}, nil
}

// ListTorrents returns a list of torrents, based in the torrent status informed.
func (c *Client) ListTorrents(torrentStatus []int) ([]*Torrent, error) {
	var torrents []*trans.Torrent
	var torrent *trans.Torrent
	var extractedTorrents []*Torrent
	var err error

	if torrents, err = c.t.GetTorrents(); err != nil {
		return nil, err
	}
	log.Printf("Found '%d' torrents in client.", len(torrents))

	for _, torrent = range torrents {
		var extractedTorrent *Torrent

		// skipping torrents that don't match desired statuses, if option informed
		if len(torrentStatus) > 0 && !intSliceContains(torrentStatus, torrent.Status) {
			log.Printf("Skipping torrent '%s' in status '%d', does not match '%v'",
				torrent.Name, torrent.Status, torrentStatus)
			continue
		}

		if extractedTorrent, err = c.extractTorrent(torrent); err != nil {
			return nil, err
		}

		extractedTorrents = append(extractedTorrents, extractedTorrent)
	}

	return extractedTorrents, nil
}

// setTrackers execute thet Set API calls to update trackers in informed torrent, when dry-run it
// will only print usual log messages.
func (c *Client) setTrackers(torrent *trans.Torrent, t *Torrent, trackers []*Tracker, dryRun bool) error {
	var trackersInTorrent = c.ListTrackers([]*Torrent{t})
	var trackerInTorrent *Tracker
	var tracker *Tracker
	var announcesInTorrent []string
	var announcesToAdd []string
	var err error

	for _, trackerInTorrent = range trackersInTorrent {
		announcesInTorrent = append(announcesInTorrent, trackerInTorrent.Announce)
	}

	for _, tracker = range trackers {
		// checking if tracker is already present in torrent
		if stringSliceContains(announcesInTorrent, tracker.Announce) {
			continue
		}
		announcesToAdd = append(announcesToAdd, tracker.Announce)
	}

	log.Printf("Adding '%d' trackers in torrent '%s' (dry-run: '%v')",
		len(announcesToAdd), t.Name, dryRun)

	if !dryRun {
		if err = torrent.Set(trans.SetTorrentArg{TrackerAdd: announcesToAdd}); err != nil {
			return err
		}
	}

	return nil
}

// UpdateTrackers executes the api calls to update the trackers in the informed list of torrents,
// when running with dry-run mode, it will pass the boolean ahead.
func (c *Client) UpdateTrackers(torrents []*Torrent, trackers []*Tracker, dryRun bool) error {
	var torrentsInClient []*trans.Torrent
	var torrentInClient *trans.Torrent
	var torrent *Torrent
	var IDs []int
	var err error

	// picking up the IDs of torrents to be updated
	for _, torrent = range torrents {
		IDs = append(IDs, torrent.ID)
	}

	// loading torrents in client
	if torrentsInClient, err = c.t.GetTorrents(); err != nil {
		return err
	}

	for _, torrentInClient = range torrentsInClient {
		log.Printf("Updating torrent: '%s' (ID: '%d')", torrentInClient.Name, torrentInClient.ID)
		if intSliceContains(IDs, torrentInClient.ID) {
			log.Printf("Skpping update on torrent '%s' (id '%d')...",
				torrentInClient.Name, torrentInClient.ID)
			continue
		}
		if torrent, err = c.extractTorrent(torrentInClient); err != nil {
			return err
		}

		if err = c.setTrackers(torrentInClient, torrent, trackers, dryRun); err != nil {
			return err
		}
		log.Printf("Done updating torrent '%s'!", torrentInClient.Name)
	}

	return nil
}

// NewClient instantiate a transmission API client.
func NewClient(config *Config) (*Client, error) {
	var client = &Client{}
	var err error

	if client.t, err = trans.New(trans.Config{
		Address:  config.Transmission.URL,
		User:     config.Transmission.Username,
		Password: config.Transmission.Password,
	}); err != nil {
		return nil, err
	}

	return client, nil
}
