package config

import (
	"fmt"

	"github.com/wallix/awless/database"
)

var Config *config

type config struct {
	Defaults map[string]interface{}
}

func LoadConfig() error {
	db, err, dbclose := database.Current()
	if err != nil {
		return fmt.Errorf("load config: %s", err)
	}
	defer dbclose()

	defaults, err := db.GetDefaults()
	if err != nil {
		return fmt.Errorf("config: load defaults: %s", err)
	}

	Config = &config{defaults}

	return nil
}
