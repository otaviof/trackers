package main

import (
	"log"
	"strconv"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Report about trackers.",
	Run:   executeListCmd,
}

var etcHosts bool
var status []string

func init() {
	var flagSet = listCmd.PersistentFlags()

	flagSet.BoolVar(&etcHosts, "etc-hosts", false, "Format data as '/etc/hosts'")
	flagSet.StringSliceVar(&status, "status", []string{"-1"}, "Filter trackers by status.")

	rootCmd.AddCommand(listCmd)
}

func executeListCmd(cmd *cobra.Command, args []string) {
	var list *trackers.List
	var statusIntSlice []int
	var i int
	var err error

	// parsing status flags to a slice of integers
	for _, s := range status {
		if i, err = strconv.Atoi(s); err != nil {
			log.Fatalf("[ERROR] Invalid --status flag: '%s'", err)
		}
		statusIntSlice = append(statusIntSlice, i)
	}

	list = trackers.NewList(storage)

	if etcHosts {
		if err = list.AsEtcHosts(statusIntSlice); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	} else {
		if err = list.AsTable(statusIntSlice); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	}
}
