package main

import (
	"fmt"
	"os"

	"github.com/wallix/awless/cmd"
	"github.com/wallix/awless/config"
)

func main() {
	config.CreateDefaultConf()

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
