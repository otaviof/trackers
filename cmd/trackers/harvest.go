package main

import (
	"log"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
)

var harvestCmd = &cobra.Command{
	Use:   "harvest",
	Short: "Retreive trackers from torrent client.",
	Run:   runHarvestCmd,
}

func init() {
	var flagSet = harvestCmd.PersistentFlags()

	flagSet.BoolVar(&dryRun, "dry-run", false, "Dry-run mode.")

	rootCmd.AddCommand(harvestCmd)
}

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
