package main

import (
	"log"
	"strconv"
	"strings"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/**
 * Root Command
 */

var rootCmd = &cobra.Command{
	Use:   "trackers",
	Short: "Keep track of trackers.",
}

/**
 * Sub-Commands
 */

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new tracker.",
	Run: func(cmd *cobra.Command, args []string) {
		var add *trackers.Add
		var announce = viper.GetString("announce")
		var addresses = strings.Split(viper.GetString("addresses"), ",")
		var err error

		// validating required option
		if announce == "" {
			log.Fatalf("[ERROR] --announce is required!")
		}

		add = trackers.NewAdd(storageInstance(), viper.GetInt("timeout"))
		if err = add.Tracker(announce, addresses, viper.GetBool("dry-run")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	},
}

var harvestCmd = &cobra.Command{
	Use:   "harvest",
	Short: "Retreive trackers from torrent client.",
	Run: func(cmd *cobra.Command, args []string) {
		var client *trackers.Client
		var harvest *trackers.Harvest
		var err error

		if client, err = trackers.NewClient(
			viper.GetString("rpc-url"), viper.GetString("username"), viper.GetString("password"),
		); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}

		harvest = trackers.NewHarvest(storageInstance(), client)
		if err = harvest.Execute(viper.GetBool("dry-run")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Report about trackers.",
	Run: func(cmd *cobra.Command, args []string) {
		var list *trackers.List
		var status []int
		var statusStr string
		var statusInt int
		var err error

		// parsing status flags to a slice of integers
		for _, statusStr = range viper.GetStringSlice("status") {
			if statusInt, err = strconv.Atoi(statusStr); err != nil {
				log.Fatalf("[ERROR] Invalid --status flag: '%s'", err)
			}
			status = append(status, statusInt)
		}

		list = trackers.NewList(storageInstance())

		if viper.GetBool("etc-hosts") {
			if err = list.AsEtcHosts(status); err != nil {
				log.Fatalf("[ERROR] %s", err)
			}
		} else {
			if err = list.AsTable(status); err != nil {
				log.Fatalf("[ERROR] %s", err)
			}
		}
	},
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor trackers functional status.",
	Run: func(cmd *cobra.Command, args []string) {
		var monitor = trackers.NewMonitor(storageInstance(), viper.GetInt("timeout"))
		var err error

		if err = monitor.Inspect(viper.GetBool("dry-run")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update tracker hostname IPv4 addresses.",
	Run: func(cmd *cobra.Command, args []string) {
		var update *trackers.Update
		var addresses = strings.Split(viper.GetString("addresses"), ",")
		var hostname = viper.GetString("hostname")
		var err error

		// validating required options first
		if hostname == "" {
			log.Fatalf("[ERROR] option --hostname is required!")
		}
		if len(addresses) == 0 {
			log.Fatalf("[ERROR] option --addresses is required!")
		}

		update = trackers.NewUpdate(storageInstance(), viper.GetInt("timeout"))

		if err = update.HostnameAddress(hostname, addresses, viper.GetBool("dry-run")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	},
}

/**
 * Helpers
 */

// storageInstance returns a Storage instance or die in error.
func storageInstance() *trackers.Storage {
	var storage *trackers.Storage
	var err error

	if storage, err = trackers.NewStorage(viper.GetString("db-file")); err != nil {
		log.Fatalf("[ERROR] On instantiating Storage: '%s'", err)
	}

	return storage
}

/**
 * Flags
 */

// rootFlags set the root command flags, and flags that are shared by more than one sub-command.
func rootFlags() {
	var flagSet = rootCmd.PersistentFlags()

	flagSet.Bool("dry-run", false, "Dry-run mode, don't commit any data.")
	flagSet.String("db-file", "/var/lib/trackers/trackers.sqlite", "SQLite database file path.")
	flagSet.Int("timeout", 15, "Timeout probing trackers.")
	flagSet.String("addresses", "", "IPv4 addresses, comma-separated list.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

// addFlags set flags for add sub-command.
func addFlags() {
	var flagSet = addCmd.PersistentFlags()

	flagSet.String("announce", "", "Tracker announce URL.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

// harvestFlags set flags for harvest sub-command.
func harvestFlags() {
	var flagSet = harvestCmd.PersistentFlags()

	flagSet.String("rpc-url", "", "Torrent client RPC URL.")
	flagSet.String("username", "", "Torrent client username.")
	flagSet.String("password", "", "Torrent client password.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

// listFlags set flgas for list sub-command.
func listFlags() {
	var flagSet = listCmd.PersistentFlags()

	flagSet.Bool("etc-hosts", false, "Format data as '/etc/hosts'")
	flagSet.String("output", "", "Save output to file.")
	// using string since IntSlice is not present in viper
	flagSet.StringSlice("status", []string{"-1"}, "Filter trackers by status.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

// monitorFlags set flags for monitor sub-command.
func monitorFlags() {
	var flagSet = monitorCmd.PersistentFlags()

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

// updateFlags set flags for update sub-command.
func updateFlags() {
	var flagSet = updateCmd.PersistentFlags()

	flagSet.String("hostname", "", "Tracker's hostname.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

/**
 * Execution
 */

// init link Cobra commands with root.
func init() {
	addFlags()
	rootFlags()
	harvestFlags()
	listFlags()
	monitorFlags()
	updateFlags()

	rootCmd.AddCommand(addCmd, harvestCmd, listCmd, monitorCmd, updateCmd)
}

// main calls Cobra execute method.
func main() {
	var err error

	if err = rootCmd.Execute(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
