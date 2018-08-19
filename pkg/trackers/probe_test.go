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

	probe = NewProbe(tracker)

	assert.NotNil(t, probe)
}

func TestProbeLookupIPs(t *testing.T) {
	var err error

	err = probe.LookupIPs()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(probe.ipv4s))
	assert.Regexp(t, regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`), probe.ipv4s[0])
}

func TestProbeReachableIPv4s(t *testing.T) {
	var err error
	var reachable []string

	reachable, err = probe.ReachableIPv4s()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(reachable))
	assert.Regexp(t, regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`), reachable[0])
}
