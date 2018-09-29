package trackers

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client *Client

func TestNewClient(t *testing.T) {
	var config = &Config{
		Transmission: TransmissionConfig{
			URL:      os.Getenv("TRANSMISSION_RPC_URL"),
			Username: os.Getenv("TRANSMISSION_RPC_USERNAME"),
			Password: os.Getenv("TRANSMISSION_RPC_PASSWORD"),
		},
	}
	var err error

	client, err = NewClient(config)

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestListTorrents(t *testing.T) {
	var torrents []*Torrent
	var err error

	torrents, err = client.ListTorrents([]int{})

	log.Printf("[TestListTorrrents] torrents='%#v'", torrents[0])

	assert.Nil(t, err)
}

func TestListTrackers(t *testing.T) {
	var torrents []*Torrent
	var trackers []*Tracker

	torrents, _ = client.ListTorrents([]int{})
	trackers = client.ListTrackers(torrents)

	log.Printf("[TestListTrackers] trackers='%#v'", trackers[0])

	assert.True(t, len(trackers) > 0)
}
