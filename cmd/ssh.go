package cmd

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/shell"
)

func init() {
	RootCmd.AddCommand(sshCmd)
}

var sshCmd = &cobra.Command{
	Use:               "ssh [user@]instance",
	Short:             "Launch a SSH (Secure Shell) session connecting to an instance",
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("instance required")
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
		exitOn(err)

		a := graph.Alias(instanceID)
		if id, ok := a.ResolveToId(instancesGraph, graph.Instance); ok {
			instanceID = id
		}

		cred, err := aws.InstanceCredentialsFromGraph(instancesGraph, instanceID)
		exitOn(err)
		var client *ssh.Client
		if user != "" {
			cred.User = user
			client, err = shell.NewClient(config.KeysDir, cred)
			exitOn(err)
			if verboseFlag {
				log.Printf("Login as '%s' on '%s', using key '%s'", user, cred.IP, cred.KeyName)
			}
			if err = shell.InteractiveTerminal(client); err != nil {
				exitOn(err)
			}
			return nil
		}
		for _, user := range aws.DefaultAMIUsers {
			cred.User = user
			client, err = shell.NewClient(config.KeysDir, cred)
			if err != nil && strings.Contains(err.Error(), "unable to authenticate") {
				continue
			}
			exitOn(err)
			log.Printf("Login as '%s' on '%s', using key '%s'", user, cred.IP, cred.KeyName)
			if err = shell.InteractiveTerminal(client); err != nil {
				exitOn(err)
			}
			return nil
		}
		return err
	},
}
