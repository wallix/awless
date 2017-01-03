package cmd

import (
	"fmt"
	"os"
)

func exitOn(err error) {
	if err != nil {
		if statsDB != nil {
			statsDB.AddLog(err.Error())
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
