package trackers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var monitor *Monitor

func TestNewMonitor(t *testing.T) {
	var storage *Storage

	storage, _ = NewStorage()
	monitor = NewMonitor(storage)
}

func TestMonitorInspect(t *testing.T) {
	var err error

	err = monitor.Inspect()

	assert.Nil(t, err)
}
