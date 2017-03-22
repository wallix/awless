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
	"regexp"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
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
		me := &whoami{}

		resp, err := aws.SecuAPI.GetCallerIdentity(nil)
		exitOn(err)

		me.Account = awssdk.StringValue(resp.Account)
		me.Arn = awssdk.StringValue(resp.Arn)
		me.UserId = awssdk.StringValue(resp.UserId)

		username := me.GetUsername()

		policies, err := aws.AccessService.(*aws.Access).ListUserPolicies(&iam.ListUserPoliciesInput{
			UserName: awssdk.String(username),
		})
		if err != nil {
			logger.Error(err)
		} else {
			for _, name := range policies.PolicyNames {
				me.InlinedPolicies = append(me.InlinedPolicies, awssdk.StringValue(name))
			}
		}

		attached, err := aws.AccessService.(*aws.Access).ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
			UserName: awssdk.String(username),
		})
		if err != nil {
			logger.Error(err)
		} else {
			for _, pol := range attached.AttachedPolicies {
				me.AttachedPolicies = append(me.AttachedPolicies, policy{Arn: awssdk.StringValue(pol.PolicyArn), Name: awssdk.StringValue(pol.PolicyName)})
			}
		}

		fmt.Printf("Username: %s, Id: %s, Account: %s\n", username, me.UserId, me.Account)
		if len(me.AttachedPolicies) > 0 {
			fmt.Println("\nAttached policies (i.e. managed):")
			for _, p := range me.AttachedPolicies {
				fmt.Printf("\t- %s\n", p.Name)
			}
		} else {
			fmt.Println("\nAttached policies (i.e. managed): none")
		}
		if len(me.InlinedPolicies) > 0 {
			fmt.Println("\nInlined policies:")
			for _, p := range me.InlinedPolicies {
				fmt.Printf("\t- %s\n", p)
			}
		} else {
			fmt.Println("\nInlined policies: none")
		}
	},
}

type policy struct {
	Arn, Name string
}

type whoami struct {
	Account, Arn, UserId string
	AttachedPolicies     []policy
	InlinedPolicies      []string
}

func (w *whoami) GetUsername() string {
	matches := usernameRegex.FindStringSubmatch(w.Arn)
	if len(matches) != 2 {
		return ""
	}
	return matches[1]
}

var usernameRegex = regexp.MustCompile(`:user/([\w-.]*)$`)
