package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/cmd"
	"github.com/wallix/awless/config"
)

func main() {
	config.InitAwlessEnv()

	sess, err := aws.InitSession(viper.GetString("region"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	aws.InitServices(sess)

	cmd.InitCli()
	cmd.ExecuteRoot()
}
