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
func (l *List) AsEtcHosts(status int) error {
	var trackers []*Tracker
	var selectedTrackers []*Tracker
	var tracker *Tracker
	var hostname string
	var addresses []string
	var err error

	if trackers, err = l.storage.Read(); err != nil {
		return err
	}

	for _, tracker = range trackers {
		var matched bool

		if status != -1 && tracker.Status != status {
			continue
		}
		// skipping ip addresses stored as domains
		if matched, _ = regexp.MatchString(`^\d+\.\d+\.\d+\.\d+$`, tracker.Hostname); matched {
			continue
		}

		selectedTrackers = append(selectedTrackers, tracker)
	}

	for hostname, addresses = range l.groupByHostname(selectedTrackers) {
		var address string

		for _, address = range addresses {
			fmt.Printf("%s %s\n", address, hostname)
		}
	}

	return nil
}

// groupByHostname from a list of trackers group information by hostname, deplicating results and
// making sure that 0.0.0.0 is not mapped for working trackers.
func (l *List) groupByHostname(trackers []*Tracker) map[string][]string {
	var tracker *Tracker
	var groupBy = make(map[string][]string)
	var hostname string
	var addresses []string

	// collecting unique entries per hostname
	for _, tracker = range trackers {
		var hostname = tracker.Hostname
		var address string

		for _, address = range tracker.Addresses {
			var exists bool

			if _, exists = groupBy[hostname]; !exists {
				groupBy[hostname] = []string{address}
			} else {
				// only saving the non-duplicates
				if !l.sliceContains(groupBy[hostname], address) {
					groupBy[hostname] = append(groupBy[hostname], address)
				}
			}
		}
	}

	// cleaning up results
	for hostname, addresses = range groupBy {
		// when a 0.0.0.0 entry is on the list, it must not contain other addresses
		if len(addresses) > 1 && l.sliceContains(addresses, "0.0.0.0") {
			groupBy[hostname] = l.sliceRemove(addresses, "0.0.0.0")
		}
	}

	return groupBy
}

// sliceContains checks if a slice contiains a string.
func (l *List) sliceContains(slice []string, str string) bool {
	var sliceStr string

	for _, sliceStr = range slice {
		if str == sliceStr {
			return true
		}
	}

	return false
}

// sliceRemove removes a string from a slice.
func (l *List) sliceRemove(slice []string, remove string) []string {
	var str string
	var cleanSlice []string

	for _, str = range slice {
		if remove == str {
			continue
		}
		cleanSlice = append(cleanSlice, str)
	}

	return cleanSlice
}

// AsTable show contents of storage object as a ascii table.
func (l *List) AsTable(status int) error {
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
		if status != -1 && tracker.Status != status {
			continue
		}

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
