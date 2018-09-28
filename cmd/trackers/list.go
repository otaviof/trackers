package main

import (
	"log"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Run:   runListCmd,
	Short: "List trackers.",
	Long: `
Read information from database and display table for "/etc/hosts" style, table format is handy for
ad-hoc usage, while "/etc/hosts" can be used to cache Tracker hostnames locally.

Trackers can also be filtered by functional status:
  - 0: service is reachable and responding;
  - 1: Can't resolv tracker's hostname;
  - 2: service does not respond;
  - 3: tracker was overwritten by "trackers overwrite";
	`,
	Example: `
trackers list --status "0,1"
trackers list --status 0 --etc-hosts`,
}

var etcHosts bool // show data as etc-hosts format

// init link command line arguments and join sub-command on main command.
func init() {
	var flagSet = listCmd.PersistentFlags()

	flagSet.BoolVar(&etcHosts, "etc-hosts", false, "Format output as '/etc/hosts' style.")
	flagSet.StringSliceVar(&statuses, "status", []string{"-1"}, "Comma-separated list of status.")

	rootCmd.AddCommand(listCmd)
}

// runListCmd execute the List sub-command.
func runListCmd(cmd *cobra.Command, args []string) {
	var list *trackers.List
	var statusIntSlice []int
	var err error

	if statusIntSlice, err = trackers.StringSliceToInt(statuses); err != nil {
		log.Fatalf("[ERROR] Invalid --status flag: '%s'", err)
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
