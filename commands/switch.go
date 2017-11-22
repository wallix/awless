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
	"github.com/wallix/awless/logger"
)

func init() {
	RootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
	Use:     "switch REGION/PROFILE",
	Aliases: []string{"sw"},
	Short:   "Quick way to switch to profiles and regions",
	Example: `  awless switch eu-west-2         # equivalent to 'awless config set aws.region eu-west-2'
  awless switch mfa               # equivalent to 'awless config set aws.profile mfa', if mfa is a valid profile in ~/.aws/{config,credentials}
  awless sw default us-west-1     # switch in region 'us-west-1', with profile 'default'`,
	PersistentPreRun:  applyHooks(initAwlessEnvHook, initLoggerHook),
	PersistentPostRun: applyHooks(includeHookIf(&config.TriggerSyncOnConfigUpdate, initCloudServicesHook)),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("REGION or PROFILE required. See examples.")
		}
		if len(args) > 2 {
			return errors.New("two many arguments provided, expected REGION or PROFILE. See examples.")
		}
		for _, arg := range args {
			if awsconfig.IsValidRegion(arg) {
				logger.Infof("Switching to region '%s'", arg)
				exitOn(config.Set(config.RegionConfigKey, arg))
				continue
			}
			if awsconfig.IsValidProfile(arg) {
				logger.Infof("Switching to profile '%s'", arg)
				exitOn(config.Set(config.ProfileConfigKey, arg))
				continue
			}
			return fmt.Errorf("invalid region or profile: '%s'", arg)
		}
		return nil
	},
}
