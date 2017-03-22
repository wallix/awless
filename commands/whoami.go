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

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/logger"
)

func init() {
	RootCmd.AddCommand(whoamiCmd)
}

var whoamiCmd = &cobra.Command{
	Use:               "whoami",
	Aliases:           []string{"who"},
	PersistentPreRun:  applyHooks(initAwlessEnvHook, initLoggerHook, initCloudServicesHook),
	PersistentPostRun: applyHooks(saveHistoryHook, verifyNewVersionHook),
	Short:             "Show your account, attached (i.e. managed) and inlined policies",

	Run: func(cmd *cobra.Command, args []string) {
		me, err := aws.AccessService.(*aws.Access).GetIdentity()
		exitOn(err)

		fmt.Printf("Username: %s, Id: %s, Account: %s\n", me.Username, me.UserId, me.Account)

		policies, err := aws.AccessService.(*aws.Access).GetUserPolicies(me.Username)
		if err != nil {
			logger.Error(err)
			return
		} else {
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
		}
	},
}
