package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	configFilename   = "config.yaml"
	databaseFilename = "database.db"
	Dir              = filepath.Join(os.Getenv("HOME"), ".awless", "aws")
	Path             = filepath.Join(Dir, configFilename)
	DatabasePath     = filepath.Join(os.Getenv("HOME"), ".awless", databaseFilename)

	InfraFilename  = "infra.rdf"
	AccessFilename = "access.rdf"
)

func CreateDefaultConf() {
	os.MkdirAll(Dir, 0700)

	if _, err := os.Stat(Path); os.IsNotExist(err) {
		ioutil.WriteFile(Path, []byte("region: \"eu-west-1\"\n"), 0400)
		fmt.Printf("Creating default config file at %s\n", Path)
	}
}
