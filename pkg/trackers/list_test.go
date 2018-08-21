package trackers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var list *List

func TestNewList(t *testing.T) {
	var storage *Storage
	storage, _ = NewStorage("/var/tmp/test.sqlite")

	list = NewList(storage)
}

func TestListAsEtcHosts(t *testing.T) {
	var err error

	err = list.AsEtcHosts([]int{-1})

	assert.Nil(t, err)
}

func TestListAsTable(t *testing.T) {
	var err error

	err = list.AsTable([]int{-1})

	assert.Nil(t, err)
}
