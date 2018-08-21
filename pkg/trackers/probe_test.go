package trackers

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var probe *Probe

func TestNewProbe(t *testing.T) {
	var tracker *Tracker

	tracker, _ = NewTracker("udp://tracker.debian.org/announce")

	probe = NewProbe(tracker, 5)

	assert.NotNil(t, probe)
}

func TestProbeLookupAddresses(t *testing.T) {
	var err error

	err = probe.LookupAddresses()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(probe.addresses))
	assert.Regexp(t, regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`), probe.addresses[0])
}

func TestProbeReachableAddresses(t *testing.T) {
	var err error
	var reachable []string

	reachable, err = probe.ReachableAddresses()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(reachable))
	assert.Regexp(t, regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`), reachable[0])
}
