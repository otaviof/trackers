package main

import (
	"log"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor trackers functional status.",
	Run:   runMonitorCmd,
}

func init() {
	var flagSet = monitorCmd.PersistentFlags()

	flagSet.BoolVar(&dryRun, "dry-run", false, "Dry-run mode.")

	rootCmd.AddCommand(monitorCmd)
}

func runMonitorCmd(cmd *cobra.Command, args []string) {
	var monitor = trackers.NewMonitor(storage, config)
	var err error

	if err = monitor.Inspect(dryRun); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
