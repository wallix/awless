package commands

import (
	"fmt"
	"os"

	"github.com/wallix/awless/database"
)

func exitOn(err error) {
	if err != nil {
		db, close := database.Current()
		defer close()
		if db != nil {
			db.AddLog(err.Error())
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
