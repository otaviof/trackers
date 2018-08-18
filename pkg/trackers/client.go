package trackers

import (
	trans "github.com/odwrtw/transmission"
)

type ClientInterface interface {
}

type Client struct {
	t *trans.Client
}

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

func (c *Client) List() ([]*Torrent, error) {
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

func NewClient(address string, user string, pass string) (*Client, error) {
	var client = &Client{}
	var err error

	if client.t, err = trans.New(trans.Config{
		Address:  address,
		User:     user,
		Password: pass,
	}); err != nil {
		return nil, err
	}

	return client, nil
}
