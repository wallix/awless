package main

import (
	"github.com/wallix/awless/api"
	"github.com/wallix/awless/cmd"
	"github.com/wallix/awless/config"
)

func main() {
	config.InitAwlessEnv()
	config.InitCloudSession()

	api.InitServices()

	cmd.InitCli()
	cmd.ExecuteRoot()
}
