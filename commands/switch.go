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

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/config"
)

func init() {
	RootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
	Use:     "switch [REGION] [PROFILE]",
	Aliases: []string{"sw"},
	Short:   "Quick way to switch awless config to given profile and/or region",
	Example: `  awless switch eu-west-2           # now using region eu-west-2'
  awless switch mfa                 # now using profile mfa (with mfa a valid profile in ~/.aws/{config,credentials})
  awless switch default us-west-1   # now using region us-west-1 and the default profile
  awless sw eu-west-3 admin         # now using profile admin in region eu-west-3`,
	PersistentPreRun:  applyHooks(initAwlessEnvHook, initLoggerHook),
	PersistentPostRun: applyHooks(includeHookIf(&config.TriggerSyncOnConfigUpdate, initCloudServicesHook)),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("REGION and/or PROFILE required. See examples.")
		}
		if len(args) > 2 {
			return errors.New("too many arguments provided, expected REGION and/or PROFILE. See examples.")
		}
		for _, arg := range args {
			if awsconfig.IsValidRegion(arg) {
				exitOn(config.Set(config.RegionConfigKey, arg))
				continue
			}
			if awsconfig.IsValidProfile(arg) {
				exitOn(config.Set(config.ProfileConfigKey, arg))
				continue
			}
			exitOn(fmt.Errorf("could not find profile: '%s' in $HOME/.aws/{credentials,config}", arg))
		}
		return nil
	},
}
