package main

import (
	"log"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var overwriteCmd = &cobra.Command{
	Use:   "overwrite",
	Short: "Overwrite tracker hostname's IPv4 addresses.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var hostname string

func init() {
	var flagSet = overwriteCmd.PersistentFlags()

	flagSet.StringVar(&hostname, "hostname", "", "Tracker's hostname.")
	flagSet.StringSliceVar(&addresses, "addresses", []string{}, "Tracker addresses, comma separated list.")

	overwriteCmd.MarkFlagRequired("hostname")
	overwriteCmd.MarkFlagRequired("addresses")

	rootCmd.AddCommand(overwriteCmd)
}

func runOverwriteCmd(cmd *cobra.Command, args []string) {
	var update *trackers.Update
	var hostname = viper.GetString("hostname")
	var err error

	update = trackers.NewUpdate(storage, config)

	if err = update.HostnameAddress(hostname, addresses, viper.GetBool("dry-run")); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

}
