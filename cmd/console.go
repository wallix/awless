package cmd

import (
	"fmt"
	"os"
)

func exitOn(err error) {
	if err != nil {
		if db != nil {
			db.AddLog(err.Error())
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
