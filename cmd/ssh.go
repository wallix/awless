package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/shell"
)

func init() {
	RootCmd.AddCommand(sshCmd)
}

var sshCmd = &cobra.Command{
	Use:   "ssh [user@]instance-id",
	Short: "Launch a SSH (Secure Shell) session connecting to an instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("ssh: instance id required")
		}
		var instanceID string
		var user string
		if strings.Contains(args[0], "@") {
			user = strings.Split(args[0], "@")[0]
			instanceID = strings.Split(args[0], "@")[1]
		} else {
			instanceID = args[0]
		}
		instancesGraph, err := aws.InfraService.InstancesGraph()
		if err != nil {
			return err
		}
		cred, err := aws.InstanceCredentialsFromGraph(instancesGraph, instanceID)
		if err != nil {
			return err
		}
		if user == "" {
			cred.User = "ec2-user" //TODO find a way to fetch the default user
		} else {
			cred.User = user
		}

		client, err := shell.NewClient(config.KeysDir, cred)
		if err != nil {
			return err
		}

		if err := shell.InteractiveTerminal(client); err != nil {
			return err
		}

		return nil
	},
}
