package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wallix/awless/cmd"
)

var (
	configFilename = "config.yaml"
	configDir      = filepath.Join(os.Getenv("HOME"), ".awless", "aws")
	configPath     = filepath.Join(configDir, configFilename)
)

func main() {
	createAwlessDefaultConf()

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func createAwlessDefaultConf() {
	os.MkdirAll(configDir, 0700)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		ioutil.WriteFile(configPath, []byte("region: \"eu-west-1\"\n"), 0400)
		fmt.Printf("Creating default config file at %s\n", configPath)
	}
}
