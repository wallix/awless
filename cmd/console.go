package cmd

import (
	"fmt"
	"os"

	"github.com/wallix/awless/database"
)

func exitOn(err error) {
	if err != nil {
		if database.Current != nil {
			database.Current.AddLog(err.Error())
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
