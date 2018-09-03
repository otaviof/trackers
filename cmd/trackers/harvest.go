package main

import (
	"log"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
)

var harvestCmd = &cobra.Command{
	Use:   "harvest",
	Run:   runHarvestCmd,
	Short: "Reaper new trackers from torrent client.",
	Long: `
Harvest trackers from torrent client, adding the new entries in the database.
	`,
}

// init bind command-line flags and sub-command in main command.
func init() {
	var flagSet = harvestCmd.PersistentFlags()

	flagSet.BoolVar(&dryRun, "dry-run", false, "Dry-run mode.")

	rootCmd.AddCommand(harvestCmd)
}

// runHarvestCmd execute harvest sub-command.
func runHarvestCmd(cmd *cobra.Command, args []string) {
	var client *trackers.Client
	var harvest *trackers.Harvest
	var err error

	if client, err = trackers.NewClient(config); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	harvest = trackers.NewHarvest(storage, client)

	if err = harvest.Execute(dryRun); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
