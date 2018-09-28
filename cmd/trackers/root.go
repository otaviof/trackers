package main

import (
	"io/ioutil"
	"log"

	trackers "github.com/otaviof/trackers/pkg/trackers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "trackers",
	Short: "Track and monitor tracker services.",
	Long: `
Book-keeping application focused in Torrent Tracker services. Keep a up-to-date database (SQLite3)
with the trackers URI, hostname and service status. Using "trackers" you can harvest new instances,
and monitor their functional status over time. Please consider the sub-commands to see the other
actions that you can execute.
	`,
}

var config *trackers.Config   // global configuration
var storage *trackers.Storage // storage interface
var statuses []string         // status filter, used in filter and update sub-commands
var dryRun bool               // dry-run flag, used in a number of sub-commands

// init instantiate command, with loading setps and command-line flags.
func init() {
	var flagSet = rootCmd.PersistentFlags()

	cobra.OnInitialize(load)

	flagSet.String("config", "/etc/trackers/trackers.yaml", "Trackers configuration file path.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

// load execute load methods.
func load() {
	loadConfig()
	storageInstance()
}

// loadConfig read configuratoin file from informed location, export a instance of config.
func loadConfig() {
	var configPath = viper.GetString("config")
	var bytes []byte
	var err error

	if bytes, err = ioutil.ReadFile(configPath); err != nil {
		log.Fatal(err)
	}

	if err = yaml.Unmarshal(bytes, &config); err != nil {
		log.Fatal(err)
	}
}

// storageInstance returns a Storage instance or die in error.
func storageInstance() {
	var err error

	if storage, err = trackers.NewStorage(config); err != nil {
		log.Fatalf("[ERROR] On instantiating Storage: '%s'", err)
	}
}
