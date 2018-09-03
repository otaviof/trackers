package main

import (
	"log"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Run:   runAddCmd,
	Short: "Add a new tracker.",
	Long: `
Adds a new tracker in the database. On adding it will check if informed address(es) are responding,
accordingly with the announce URL protocol and port.
	`,
}

var announce string    // announce URL
var addresses []string // tracker addresses

// init bind commnad-line flags and sub-command in root command.
func init() {
	var flagSet = addCmd.PersistentFlags()

	flagSet.StringVar(&announce, "announce", "", "Tracker announce URL.")
	flagSet.StringSliceVar(&addresses, "addresses", []string{}, "Tracker addresses, comma separated list.")
	flagSet.BoolVar(&dryRun, "dry-run", false, "Dry-run mode.")

	addCmd.MarkFlagRequired("announce")

	rootCmd.AddCommand(addCmd)
}

// runAddCmd executes Add sub-command.
func runAddCmd(cmd *cobra.Command, args []string) {
	var add *trackers.Add
	var err error

	if announce == "" {
		log.Fatal("[ERROR] Parameter --announce is required!")
	}

	add = trackers.NewAdd(storage, config)

	if err = add.Tracker(announce, addresses, dryRun); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
