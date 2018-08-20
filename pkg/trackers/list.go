package trackers

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// List shows data held in storage.
type List struct {
	storage storageInterface
}

// AsEtcHosts print storage contents following /etc/hosts format.
func (l *List) AsEtcHosts() error {
	var trackers []*Tracker
	var tracker *Tracker
	var err error

	if trackers, err = l.storage.Read(); err != nil {
		return err
	}

	for _, tracker = range trackers {
		var address string
		var matched bool

		// skipping ipaddresses
		if matched, _ = regexp.MatchString(`^\d+\.\d+\.\d+\.\d+$`, tracker.Hostname); matched {
			continue
		}

		for _, address = range tracker.Addresses {
			fmt.Printf("%s %s # announce='%s'\n", address, tracker.Hostname, tracker.Announce)
		}
	}

	return nil
}

// AsTable show contents of storage object as a ascii table.
func (l *List) AsTable() error {
	var trackers []*Tracker
	var tracker *Tracker
	var tableWriter *tablewriter.Table
	var err error

	if trackers, err = l.storage.Read(); err != nil {
		return err
	}

	tableWriter = tablewriter.NewWriter(os.Stdout)
	tableWriter.SetHeader([]string{"Hostname", "Announce", "Addresses", "Status"})

	for _, tracker = range trackers {
		tableWriter.Append([]string{
			tracker.Hostname,
			tracker.Announce,
			strings.Join(tracker.Addresses, ", "),
			strconv.Itoa(tracker.Status),
		})
	}

	tableWriter.Render()
	return nil
}

// NewList instantiate a List object with a storage instance.
func NewList(storage storageInterface) *List {
	return &List{storage: storage}
}
