package trackers

import (
	trans "github.com/odwrtw/transmission"
)

type clientInterface interface {
	ListTrackers() ([]*Tracker, error)
	ListTorrents() ([]*Torrent, error)
}

// Client tranmission torrent client.
type Client struct {
	t *trans.Client
}

// ListTrackers based on torrents, extracts a unique list.
func (c *Client) ListTrackers() ([]*Tracker, error) {
	var torrent *Torrent
	var torrents []*Torrent
	var tracker *Tracker
	var trackers []*Tracker
	var trackerMap = make(map[string]*Tracker)
	var exists bool
	var err error

	if torrents, err = c.ListTorrents(); err != nil {
		return nil, err
	}

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

	return trackers, nil
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

// ListTorrents returns a list of torrents.
func (c *Client) ListTorrents() ([]*Torrent, error) {
	var torrents []*trans.Torrent
	var torrent *trans.Torrent
	var extractedTorrents []*Torrent
	var err error

	if torrents, err = c.t.GetTorrents(); err != nil {
		return nil, err
	}

	for _, torrent = range torrents {
		var extractedTorrent *Torrent

		if extractedTorrent, err = c.extractTorrent(torrent); err != nil {
			return nil, err
		}

		extractedTorrents = append(extractedTorrents, extractedTorrent)
	}

	return extractedTorrents, nil
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
