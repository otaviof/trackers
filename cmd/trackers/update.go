package main

import (
	"log"

	"github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Run:   runUpdateCmd,
	Short: "Update list of trackers in torrents",
	Long: `
Update running torrents with local trackers database entries. By running update, Trackers will
inspect running torrents, and let you add trackers based on tracker and torrents status.

The torrent status may be:
 - 0: stopped;
 - 1: check pending;
 - 2: checking;
 - 3: download pending;
 - 4: downloading;
 - 5: seed pending;
 - 6: seeding;

Trackers status can be:
  - 0: service is reachable and responding;
  - 1: Can't resolv tracker's hostname;
  - 2: service does not respond;
  - 3: tracker was overwritten by "trackers overwrite";`,
}

var torrentStatuses []string // torrent status filter

func init() {
	var flagSet = updateCmd.PersistentFlags()

	flagSet.StringSliceVar(
		&statuses,
		"status",
		[]string{"0", "3"},
		"Comma-separated list of status.",
	)
	flagSet.StringSliceVar(
		&torrentStatuses,
		"torrent-status",
		[]string{"4"},
		"Comma-separated list of torrent status.",
	)
	flagSet.BoolVar(&dryRun, "dry-run", false, "Dry-run mode.")

	rootCmd.AddCommand(updateCmd)
}

// runUpdateCmd execute the steps to update the trackers in torrents, based in tracker and torrent
// status.
func runUpdateCmd(cmd *cobra.Command, args []string) {
	var client *trackers.Client
	var update *trackers.Update
	var statusesInt []int
	var torrentStatusesInt []int
	var err error

	if client, err = trackers.NewClient(config); err != nil {
		log.Fatal(err)
	}
	update = trackers.NewUpdate(storage, client)

	if statusesInt, err = trackers.StringSliceToInt(statuses); err != nil {
		log.Fatalf("Error parsing --status flag: '%s'", err)
	}
	if torrentStatusesInt, err = trackers.StringSliceToInt(torrentStatuses); err != nil {
		log.Fatalf("Error parsing --torrent-status flag: '%s'", err)
	}

	if err = update.Execute(statusesInt, torrentStatusesInt, dryRun); err != nil {
		log.Fatal(err)
	}
}
