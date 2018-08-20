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
 * Commands
 */

var rootCmd = &cobra.Command{
	Use:   "trackers",
	Short: "Keep track of trackers.",
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new tracker.",
	Run: func(cmd *cobra.Command, args []string) {
		var storage *trackers.Storage
		var add *trackers.Add
		var announce = viper.GetString("announce")
		var addresses = strings.Split(viper.GetString("addresses"), ",")
		var err error

		if announce == "" {
			log.Fatalf("[ERROR] --announce is required!")
		}

		if storage, err = trackers.NewStorage(viper.GetString("db-file")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}

		add = trackers.NewAdd(storage, viper.GetInt("timeout"))
		if err = add.Tracker(announce, addresses, viper.GetBool("dry-run")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	},
}

var harvestCmd = &cobra.Command{
	Use:   "harvest",
	Short: "Retreive trackers from torrent client.",
	Run: func(cmd *cobra.Command, args []string) {
		var storage *trackers.Storage
		var client *trackers.Client
		var harvest *trackers.Harvest
		var err error

		if storage, err = trackers.NewStorage(viper.GetString("db-file")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}

		if client, err = trackers.NewClient(
			viper.GetString("rpc-url"), viper.GetString("username"), viper.GetString("password"),
		); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}

		harvest = trackers.NewHarvest(storage, client)
		if err = harvest.Execute(viper.GetBool("dry-run")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Report about trackers.",
	Run: func(cmd *cobra.Command, args []string) {
		var storage *trackers.Storage
		var list *trackers.List
		var status []int
		var statusStr string
		var statusInt int
		var err error

		for _, statusStr = range viper.GetStringSlice("status") {
			if statusInt, err = strconv.Atoi(statusStr); err != nil {
				log.Fatalf("[ERROR] %s", err)
			}
			status = append(status, statusInt)
		}

		if storage, err = trackers.NewStorage(viper.GetString("db-file")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}

		list = trackers.NewList(storage)

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
		var storage *trackers.Storage
		var monitor *trackers.Monitor
		var err error

		if storage, err = trackers.NewStorage(viper.GetString("db-file")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}

		monitor = trackers.NewMonitor(storage, viper.GetInt("timeout"))

		if err = monitor.Inspect(viper.GetBool("dry-run")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update tracker hostname IPv4 addresses.",
	Run: func(cmd *cobra.Command, args []string) {
		var storage *trackers.Storage
		var update *trackers.Update
		var addresses = strings.Split(viper.GetString("addresses"), ",")
		var hostname = viper.GetString("hostname")
		var err error

		if hostname == "" {
			log.Fatalf("[ERROR] --hostname is required!")
		}
		if len(addresses) == 0 {
			log.Fatalf("[ERROR] --addresses is required!")
		}

		if storage, err = trackers.NewStorage(viper.GetString("db-file")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}

		update = trackers.NewUpdate(storage, viper.GetInt("timeout"))

		if err = update.HostnameAddress(viper.GetString("hostname"), addresses, viper.GetBool("dry-run")); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	},
}

/**
 * Flags
 */

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

func addFlags() {
	var flagSet = addCmd.PersistentFlags()

	flagSet.String("announce", "", "Tracker announce URL.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

func harvestFlags() {
	var flagSet = harvestCmd.PersistentFlags()

	flagSet.String("rpc-url", "", "Torrent client RPC URL.")
	flagSet.String("username", "", "Torrent client username.")
	flagSet.String("password", "", "Torrent client password.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

func listFlags() {
	var flagSet = listCmd.PersistentFlags()

	flagSet.Bool("etc-hosts", false, "Format data as '/etc/hosts'")
	flagSet.String("output", "", "Save output to file.")
	flagSet.StringSlice("status", []string{"-1"}, "Filter trackers by status, '-1' shows all.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

func monitorFlags() {
	var flagSet = monitorCmd.PersistentFlags()

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

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

func init() {
	addFlags()
	rootFlags()
	harvestFlags()
	listFlags()
	monitorFlags()
	updateFlags()

	rootCmd.AddCommand(addCmd, harvestCmd, listCmd, monitorCmd, updateCmd)
}

func main() {
	var err error

	if err = rootCmd.Execute(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
