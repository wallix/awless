package config

import (
	"fmt"
	"os"

	"github.com/wallix/awless/database"
)

func GetDefaultRegion() string {
	db, close := database.Current()
	defer close()
	region, ok := db.GetDefaultString(RegionKey)
	if !ok {
		fmt.Fprintf(os.Stderr, "config: missing region. Set it with `awless config set region`")
		os.Exit(-1)
	}
	return region
}
