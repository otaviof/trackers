package trackers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var harvest *Harvest

func TestNewHarvest(t *testing.T) {
	var client *Client
	var storage *Storage
	var config = &Config{
		Persistence: PersistenceConfig{
			DbPath: "/var/tmp/test.sqlite",
		},
		Transmission: TransmissionConfig{
			URL:      os.Getenv("TRANSMISSION_RPC_URL"),
			Username: os.Getenv("TRANSMISSION_RPC_USERNAME"),
			Password: os.Getenv("TRANSMISSION_RPC_PASSWORD"),
		},
	}

	client, _ = NewClient(config)
	storage, _ = NewStorage(config)

	harvest = NewHarvest(storage, client)

	assert.NotNil(t, harvest)
}

func TestExecute(t *testing.T) {
	var err error

	err = harvest.Execute(false)

	assert.Nil(t, err)
}
