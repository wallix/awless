/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	awsconfig "github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
)

var keyPathFlag string

func init() {
	RootCmd.AddCommand(sshCmd)
	sshCmd.Flags().StringVarP(&keyPathFlag, "identity", "i", "", "Set path toward the identity (key file) to use to connect through SSH")
}

var sshCmd = &cobra.Command{
	Use:   "ssh [USER@]INSTANCE",
	Short: "Launch a SSH (Secure Shell) session to an instance given an id or alias",
	Example: `  awless ssh i-8d43b21b                       # using the instance id
  awless ssh ec2-user@redis-prod              # using the instance name and specify a user
  awless ssh @redis-prod -i ./path/toward/key # with a keyfile`,
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook),
	PersistentPostRun: applyHooks(saveHistoryHook, verifyNewVersionHook),

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

		instancesGraph, err := aws.InfraService.FetchByType(cloud.Instance)
		exitOn(err)

		a := graph.Alias(instanceID)
		if id, ok := a.ResolveToId(instancesGraph, cloud.Instance); ok {
			instanceID = id
		}

		cred, err := instanceCredentialsFromGraph(instancesGraph, instanceID, keyPathFlag)
		exitOn(err)

		var client *ssh.Client
		if user != "" {
			cred.User = user
			client, err = console.NewSSHClient(cred)
			exitOn(err)
			exitOn(sshConnect(client, cred))
			return nil
		}
		for _, user := range awsconfig.DefaultAMIUsers {
			cred.User = user
			client, err = console.NewSSHClient(cred)
			if err != nil && strings.Contains(err.Error(), "unable to authenticate") {
				continue
			}
			exitOn(err)
			exitOn(sshConnect(client, cred))
			return nil
		}
		return err
	},
}

func instanceCredentialsFromGraph(g *graph.Graph, instanceID, keyPathFlag string) (*console.Credentials, error) {
	inst, err := g.GetResource(cloud.Instance, instanceID)
	if err != nil {
		return nil, err
	}

	ip, ok := inst.Properties["PublicIP"]
	if !ok {
		return nil, fmt.Errorf("no public IP address for instance %s", instanceID)
	}
	var keyPath string
	if keyPathFlag != "" {
		keyPath = keyPathFlag
	} else {
		key, ok := inst.Properties["SSHKey"]
		if !ok {
			return nil, fmt.Errorf("no access key set for instance %s", instanceID)
		}
		keyPath = path.Join(config.KeysDir, fmt.Sprint(key))
	}
	return &console.Credentials{IP: fmt.Sprint(ip), User: "", KeyPath: keyPath}, nil
}

func sshConnect(sshClient *ssh.Client, cred *console.Credentials) error {
	sshPath, sshErr := exec.LookPath("ssh")
	if sshErr == nil {
		logger.Infof("Login as '%s' on '%s', using key '%s' with ssh client at '%s'", cred.User, cred.IP, cred.KeyPath, sshPath)
		args := []string{"ssh", "-i", cred.KeyPath, fmt.Sprintf("%s@%s", cred.User, cred.IP)}
		return syscall.Exec(sshPath, args, os.Environ())
	} else { // Fallback SSH
		logger.Infof("No SSH. Fallback on builtin client. Login as '%s' on '%s', using key '%s'", cred.User, cred.IP, cred.KeyPath)
		return console.InteractiveTerminal(sshClient)
	}
}
