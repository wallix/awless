package main

import (
	"github.com/wallix/awless/cmd"
	"github.com/wallix/awless/config"
)

func main() {
	config.CreateDefaultConf()
	cmd.RootCmd.Execute()
}
