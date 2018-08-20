package trackers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var storage *Storage

func TestNewStorage(t *testing.T) {
	var err error

	storage, err = NewStorage("/var/tmp/test.sqlite")

	assert.Nil(t, err)
}

func TestStorageInitDB(t *testing.T) {
	var err error

	err = storage.InitDB()

	assert.Nil(t, err)
}

func TestStorageWrite(t *testing.T) {
	var tracker *Tracker
	// var err error

	tracker, _ = NewTracker("udp://tracker.debian.org:80/announce")
	_ = storage.Write([]*Tracker{tracker})

	// assert.Nil(t, err)
}
