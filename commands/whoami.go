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
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/logger"
)

var onlyMyIPFlag, onlyMyNameFlag, onlyMyTypeFlag, onlyMyIDFlag, onlyMyAccountFlag, onlyMyResourcePathFlag bool

func init() {
	RootCmd.AddCommand(whoamiCmd)

	whoamiCmd.Flags().BoolVar(&onlyMyIPFlag, "ip-only", false, "Only returns your IP address as seen by AWS")
	whoamiCmd.Flags().BoolVar(&onlyMyNameFlag, "name-only", false, "Only returns the name of the resource (ex: username for a user)")
	whoamiCmd.Flags().BoolVar(&onlyMyTypeFlag, "type-only", false, "Only returns the type of the resource (ex: user for a user)")
	whoamiCmd.Flags().BoolVar(&onlyMyIDFlag, "id-only", false, "Only returns the ID of the resource")
	whoamiCmd.Flags().BoolVar(&onlyMyAccountFlag, "account-only", false, "Only returns the AWS account number")
	whoamiCmd.Flags().BoolVar(&onlyMyResourcePathFlag, "resource-only", false, "Only returns the AWS ARN resource path suffix (ex: user/jsmith)")
}

var whoamiCmd = &cobra.Command{
	Use:               "whoami",
	Aliases:           []string{"who"},
	PersistentPreRun:  applyHooks(initAwlessEnvHook, initLoggerHook, initCloudServicesHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),
	Short:             "Show your account, attached (i.e. managed) and inlined policies",

	Run: func(cmd *cobra.Command, args []string) {
		if onlyMyIPFlag {
			fmt.Println(getMyIP())
			return
		}

		if localGlobalFlag {
			logger.Warning("`--local` flag prevent the command from fetching remote information")
			return
		}

		me, err := awsservices.AccessService.(*awsservices.Access).GetIdentity()
		exitOn(err)

		if me.IsRoot() {
			logger.Warning("You are currently root")
			logger.Warning("Best practices suggest to create a new user and affecting it roles of access")
			logger.Warning("awless official templates might help https://github.com/wallix/awless-templates\n")
		}

		switch {
		case onlyMyAccountFlag:
			fmt.Println(me.Account)
			return
		case onlyMyIDFlag:
			fmt.Println(me.UserId)
			return
		case onlyMyNameFlag:
			fmt.Println(me.Resource)
			return
		case onlyMyTypeFlag:
			fmt.Println(me.ResourceType)
			return
		case onlyMyResourcePathFlag:
			fmt.Println(me.ResourcePath)
			return
		}

		if !me.IsUserType() {
			fmt.Printf("ResourceType: %s, Resource: %s, Id: %s, Account: %s\n", me.ResourceType, me.Resource, me.UserId, me.Account)
			return
		}

		fmt.Printf("Username: %s, Id: %s, Account: %s\n", me.Resource, me.UserId, me.Account)

		policies, err := awsservices.AccessService.(*awsservices.Access).GetUserPolicies(me.Resource)
		if err != nil {
			if aerr, ok := err.(awserr.RequestFailure); ok && aerr.Code() == "AccessDenied" {
				logger.Warningf("user '%s' is not authorized to list its policies", me.Resource)
			} else {
				logger.Error(err)
			}
			return
		}

		if attached := policies.Attached; len(attached) > 0 {
			fmt.Println("\nAttached policies (i.e. managed):")
			for _, name := range attached {
				fmt.Printf("\t- %s\n", name)
			}
		} else {
			fmt.Println("\nAttached policies (i.e. managed): none")
		}
		if inlined := policies.Inlined; len(inlined) > 0 {
			fmt.Println("\nInlined policies:")
			for _, name := range inlined {
				fmt.Printf("\t- %s\n", name)
			}
		} else {
			fmt.Println("\nInlined policies: none")
		}
		if byGroup := policies.ByGroup; len(byGroup) > 0 {
			for g, pol := range byGroup {
				fmt.Printf("\nPolicies from group '%s': %s\n", g, strings.Join(pol, ", "))
			}
		}
	},
}

func getMyIP() net.IP {
	client := &http.Client{Timeout: 3 * time.Second}
	if resp, err := client.Get("http://checkip.amazonaws.com/"); err == nil {
		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return net.ParseIP(strings.TrimSpace(string(b)))
	}
	return nil
}
