package main

import (
	"fmt"
	"os"

	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/cmd"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
)

func main() {
	err := config.InitAwlessEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot init environment: %s\n", err)
		os.Exit(1)
	}

	sess, err := aws.InitSession()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	aws.InitServices(sess)

	err = database.InitDB(config.AwlessFirstInstall)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot init database: %s\n", err)
		os.Exit(1)
	}

	cmd.InitCli()
	cmd.ExecuteRoot()
}
