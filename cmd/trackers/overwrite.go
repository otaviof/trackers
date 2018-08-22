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
	Run:   runOverwriteCmd,
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
	var overwrite *trackers.Overwrite
	var err error

	if hostname == "" {
		log.Fatal("[ERROR] Parameter --hostname is required!")
	}

	overwrite = trackers.NewOverwrite(storage, config)

	if err = overwrite.HostnameAddress(hostname, addresses, viper.GetBool("dry-run")); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

}
