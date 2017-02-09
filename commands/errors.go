package commands

import (
	"fmt"
	"os"

	"github.com/wallix/awless/database"
)

func exitOn(err error) {
	if err != nil {
		db, dberr, close := database.Current()
		if dberr == nil && db != nil {
			defer close()
			db.AddLog(err.Error())
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
