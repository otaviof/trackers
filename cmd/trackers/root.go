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
	Short: "Keep track of trackers.",
}

var config *trackers.Config
var storage *trackers.Storage
var dryRun bool

func init() {
	var flagSet = rootCmd.PersistentFlags()

	cobra.OnInitialize(load)

	flagSet.String("config", "/etc/trackers/trackers.yaml", "Trackers configuration file path.")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal(err)
	}
}

func load() {
	loadConfig()
	storageInstance()
}

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
