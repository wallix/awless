package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	configFilename = "config.yaml"
	dir            = filepath.Join(os.Getenv("HOME"), ".awless", "aws")
	Path           = filepath.Join(dir, configFilename)
)

func CreateDefaultConf() {
	os.MkdirAll(dir, 0700)

	if _, err := os.Stat(Path); os.IsNotExist(err) {
		ioutil.WriteFile(Path, []byte("region: \"eu-west-1\"\n"), 0400)
		fmt.Printf("Creating default config file at %s\n", Path)
	}
}
