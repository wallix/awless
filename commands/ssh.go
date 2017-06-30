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
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/ssh"
)

var keyPathFlag string
var sshPortFlag int
var printSSHConfigFlag bool
var printSSHCLIFlag bool
var privateIPFlag bool
var disableStrictHostKeyCheckingFlag bool

func init() {
	RootCmd.AddCommand(sshCmd)
	sshCmd.Flags().StringVarP(&keyPathFlag, "identity", "i", "", "Set path or name toward the identity (key file) to use to connect through SSH")
	sshCmd.Flags().IntVar(&sshPortFlag, "port", 22, "Set SSH port")
	sshCmd.Flags().BoolVar(&printSSHConfigFlag, "print-config", false, "Print SSH configuration for ~/.ssh/config file.")
	sshCmd.Flags().BoolVar(&printSSHCLIFlag, "print-cli", false, "Print the CLI one-liner to connect with SSH. (/usr/bin/ssh user@ip -i ...)")
	sshCmd.Flags().BoolVar(&privateIPFlag, "private", false, "Use private ip to connect to host")
	sshCmd.Flags().BoolVar(&disableStrictHostKeyCheckingFlag, "disable-strict-host-keychecking", false, "Disable the remote host key check from ~/.ssh/known_hosts or ~/.awless/known_hosts file")
}

var sshCmd = &cobra.Command{
	Use:   "ssh [USER@]INSTANCE",
	Short: "Launch a SSH (Secure Shell) session to an instance given an id or alias",
	Example: `  awless ssh i-8d43b21b                       # using the instance id
  awless ssh ec2-user@redis-prod              # using the instance name and specify a user
  awless ssh redis-prod -i keyname # using a key stored in ~/.ssh/keyname.pem
  awless ssh redis-prod -i ./path/toward/key # with a keyfile`,
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("instance required")
		}

		connectionCtx, err := initInstanceConnectionContext(args[0], keyPathFlag)
		exitOn(err)

		client, err := ssh.InitClient(connectionCtx.keypath, config.KeysDir, filepath.Join(os.Getenv("HOME"), ".ssh"))
		client.Port = sshPortFlag

		if err != nil && strings.Contains(err.Error(), "cannot find SSH key") && keyPathFlag == "" {
			logger.Info("you may want to specify a key filepath with `-i /path/to/key.pem`")
		}
		exitOn(err)

		client.SetLogger(logger.DefaultLogger)
		client.SetStrictHostKeyChecking(!disableStrictHostKeyCheckingFlag)
		client.InteractiveTerminalFunc = console.InteractiveTerminal
		client.IP = connectionCtx.ip

		if connectionCtx.user != "" {
			client, err = client.DialWithUsers(connectionCtx.user)
		} else {
			client, err = client.DialWithUsers(awsconfig.DefaultAMIUsers...)
		}

		if isConnectionRefusedErr(err) {
			logger.Warning("cannot connect to this instance, maybe the system is still booting?")
			exitOn(err)
			return nil
		}

		if err != nil {
			if e := connectionCtx.checkInstanceAccessible(); e != nil {
				logger.Error(e.Error())
			}
			exitOn(err)
		}

		if printSSHConfigFlag {
			fmt.Println(client.SSHConfigString(connectionCtx.instanceName))
			return nil
		}

		if printSSHCLIFlag {
			fmt.Println(client.ConnectString())
			return nil
		}

		exitOn(client.Connect())
		return nil
	},
}

func instanceCredentialsFromGraph(g *graph.Graph, inst *graph.Resource, keyFlag string) (keypath string, ip string, err error) {
	if ip, err = getIP(inst); err != nil {
		return
	}

	if keyFlag != "" {
		keypath = keyFlag
	} else {
		keypair, ok := inst.Properties[properties.KeyPair]
		if !ok {
			return
		}
		keypath = fmt.Sprint(keypair)
	}

	return
}

func isConnectionRefusedErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "connection refused")
}

func getIP(inst *graph.Resource) (string, error) {
	var ipKeyType string

	if privateIPFlag {
		ipKeyType = properties.PrivateIP
	} else {
		ipKeyType = properties.PublicIP
	}

	ip, ok := inst.Properties[ipKeyType]

	if !ok {
		return "", fmt.Errorf("no IP address for instance %s", inst.Id())
	}

	return fmt.Sprint(ip), nil
}

type instanceConnectionContext struct {
	user           string
	ip             string
	keypath        string
	instanceName   string
	instance       *graph.Resource
	resourcesGraph *graph.Graph
	myip           net.IP
}

func initInstanceConnectionContext(userhost, keypath string) (*instanceConnectionContext, error) {
	ctx := &instanceConnectionContext{}

	if strings.Contains(userhost, "@") {
		ctx.user = strings.Split(userhost, "@")[0]
		ctx.instanceName = strings.Split(userhost, "@")[1]
	} else {
		ctx.instanceName = userhost
	}

	ctx.fetchConnectionInfo()

	instanceResolvers := []graph.Resolver{&graph.ByProperty{Key: "Name", Value: ctx.instanceName}, &graph.ByType{Typ: cloud.Instance}}
	resources, err := ctx.resourcesGraph.ResolveResources(&graph.And{Resolvers: instanceResolvers})
	exitOn(err)
	switch len(resources) {
	case 0:
		// No instance with that name, use the id
		ctx.instance, err = findResource(ctx.resourcesGraph, ctx.instanceName, cloud.Instance)
		exitOn(err)
	case 1:
		ctx.instance = resources[0]
	default:
		idStatus := graph.Resources(resources).Map(func(r *graph.Resource) string {
			return fmt.Sprintf("%s (%s)", r.Id(), r.Properties[properties.State])
		})
		logger.Infof("Found %d resources with name '%s': %s", len(resources), ctx.instanceName, strings.Join(idStatus, ", "))

		var running []*graph.Resource
		running, err = ctx.resourcesGraph.ResolveResources(&graph.And{Resolvers: append(instanceResolvers, &graph.ByProperty{Key: properties.State, Value: "running"})})
		exitOn(err)

		switch len(running) {
		case 0:
			logger.Warning("None of them is running, cannot connect through SSH")
			return ctx, errors.New("non running instances")
		case 1:
			logger.Infof("Found only one instance running: %s. Will connect to this instance.", running[0].Id())
			ctx.instance = running[0]
		default:
			logger.Warning("Connect through the running ones using their id:")
			for _, res := range running {
				var up string
				if uptime, ok := res.Properties[properties.Launched].(time.Time); ok {
					up = fmt.Sprintf("\t\t(uptime: %s)", console.HumanizeTime(uptime))
				}
				logger.Warningf("\t`awless ssh %s`%s", res.Id(), up)
			}
			return ctx, errors.New("use instances ids")
		}
	}

	keypath, IP, err := instanceCredentialsFromGraph(ctx.resourcesGraph, ctx.instance, keypath)
	if err != nil {
		return nil, err
	}

	ctx.keypath = keypath
	ctx.ip = IP

	return ctx, nil
}

func (ctx *instanceConnectionContext) fetchConnectionInfo() {
	var resourcesGraph, sgroupsGraph *graph.Graph
	var myip net.IP
	var wg sync.WaitGroup
	var errc = make(chan error)

	wg.Add(1)
	go func() {
		var err error
		defer wg.Done()
		resourcesGraph, err = aws.InfraService.FetchByType(cloud.Instance)
		if err != nil {
			errc <- err
		}
	}()

	wg.Add(1)
	go func() {
		var err error
		defer wg.Done()
		sgroupsGraph, err = aws.InfraService.FetchByType(cloud.SecurityGroup)
		if err != nil {
			errc <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		myip = getMyIP()
	}()
	go func() {
		wg.Wait()
		close(errc)
	}()
	for err := range errc {
		if err != nil {
			exitOn(err)
		}
	}
	resourcesGraph.AddGraph(sgroupsGraph)

	ctx.resourcesGraph = resourcesGraph
	ctx.myip = myip
	return
}

func (ctx *instanceConnectionContext) checkInstanceAccessible() (err error) {
	state, ok := ctx.instance.Properties[properties.State]
	if st := fmt.Sprint(state); ok && st != "running" {
		logger.Warningf("this instance is '%s' (cannot ssh to a non running state)", st)
		if st == "stopped" {
			logger.Warningf("you can start it with `awless -f start instance id=%s`", ctx.instance.Id())
		}
		return errors.New("instance not accessible")
	}

	sgroups, ok := ctx.instance.Properties[properties.SecurityGroups].([]string)
	if ok {
		var sshPortOpen, myIPAllowed bool
		for _, id := range sgroups {
			var sgroup *graph.Resource
			sgroup, err = findResource(ctx.resourcesGraph, id, cloud.SecurityGroup)
			if err != nil {
				logger.Errorf("cannot get securitygroup '%s' for instance '%s': %s", id, ctx.instance.Id(), err)
				break
			}

			rules, ok := sgroup.Properties[properties.InboundRules].([]*graph.FirewallRule)
			if ok {
				for _, r := range rules {
					if r.PortRange.Contains(22) {
						sshPortOpen = true
					}
					if ctx.myip != nil && r.Contains(ctx.myip.String()) {
						myIPAllowed = true
					}
				}
			}
		}

		if !sshPortOpen {
			logger.Warning("port 22 is not open on this instance")
			return errors.New("instance not accessible")
		}

		if !myIPAllowed && ctx.myip != nil {
			logger.Warningf("your ip %s is not authorized for this instance. You might want to update the securitygroup with:", ctx.myip)
			var group = "mygroup"
			if len(sgroups) == 1 {
				group = sgroups[0]
			}
			logger.Warningf("`awless update securitygroup id=%s inbound=authorize protocol=tcp cidr=%s/32 portrange=22`", group, ctx.myip)
			return errors.New("instance not accessible")
		}
	}

	return nil
}

func findResource(g *graph.Graph, id, typ string) (*graph.Resource, error) {
	if found, err := g.FindResource(id); found == nil || err != nil {
		return nil, fmt.Errorf("instance '%s' not found", id)
	}

	return g.GetResource(typ, id)
}
