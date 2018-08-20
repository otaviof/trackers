package trackers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTracker(t *testing.T) {
	var announce string
	var tracker *Tracker
	var err error

	// testing with and without port in the url
	for _, announce = range []string{
		"udp://tracker.debian.org:80/announce",
		"udp://tracker.debian.org/announce",
	} {
		tracker, err = NewTracker(announce)

		assert.Nil(t, err)
		assert.NotNil(t, tracker)

		assert.Equal(t, announce, tracker.Announce)
		assert.Equal(t, "tracker.debian.org", tracker.Hostname)
		assert.Equal(t, 80, tracker.Port)
		assert.Equal(t, "udp", tracker.Protocol)
	}
}
