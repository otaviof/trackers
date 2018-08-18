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
		"udp://tracker.publicbt.com:80/announce",
		"udp://tracker.publicbt.com/announce",
	} {
		tracker, err = NewTracker(announce)

		assert.Nil(t, err)
		assert.NotNil(t, tracker)

		assert.Equal(t, announce, tracker.Announce)
		assert.Equal(t, "tracker.publicbt.com", tracker.Hostname)
		assert.Equal(t, 80, tracker.Port)
		assert.Equal(t, "udp", tracker.Protocol)
	}
}
