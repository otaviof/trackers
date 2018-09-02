package trackers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var monitor *Monitor

func TestNewMonitor(t *testing.T) {
	var config = &Config{
		Probe:       ProbeConfig{Timeout: 5},
		Persistence: PersistenceConfig{DbPath: "/var/tmp/test.sqlite"},
		Nameservers: []string{"1.1.1.1:853", "1.0.0.1:853"},
	}
	var storage *Storage

	storage, _ = NewStorage(config)
	monitor = NewMonitor(storage, config, 1)
}

func TestMonitorInspect(t *testing.T) {
	var err error

	err = monitor.Inspect(false)

	assert.Nil(t, err)
}
