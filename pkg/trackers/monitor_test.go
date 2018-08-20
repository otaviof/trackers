package trackers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var monitor *Monitor

func TestNewMonitor(t *testing.T) {
	var storage *Storage

	storage, _ = NewStorage("/var/tmp/test.sqlite")
	monitor = NewMonitor(storage, 5)
}

func TestMonitorInspect(t *testing.T) {
	var err error

	err = monitor.Inspect(false)

	assert.Nil(t, err)
}
