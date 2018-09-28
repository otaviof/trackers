package main

import (
	"log"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var overwriteCmd = &cobra.Command{
	Use:   "overwrite",
	Run:   runOverwriteCmd,
	Short: "Overwrite tracker hostname's IPv4 addresses.",
	Long: `
Overwrite command allows you set a given IPv4 address to a given tracker's hostname. Keep in mind
that the IPv4 address belongs to the hostname, and therefore it may affect more than a single
entry in the database.

Upon overwrite the IPv4 addresses will be probed, they must be responsive accordingly to the announce
URL that will be relying in that hostname.
	`,
}

var hostname string // tracker hostname

func init() {
	var flagSet = overwriteCmd.PersistentFlags()

	flagSet.StringVar(&hostname, "hostname", "", "Tracker's hostname.")
	flagSet.StringSliceVar(&addresses, "addresses", []string{}, "Tracker addresses, comma separated list.")

	overwriteCmd.MarkFlagRequired("hostname")
	overwriteCmd.MarkFlagRequired("addresses")

	rootCmd.AddCommand(overwriteCmd)
}

// runOverwriteCmd executes overwrite sub-command.
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
