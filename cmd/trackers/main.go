package main

import "log"

// main calls Cobra execute method.
func main() {
	var err error

	if err = rootCmd.Execute(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
