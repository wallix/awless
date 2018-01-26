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
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/config"
)

var keysOnly bool

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVar(&keysOnly, "keys", false, "list only config keys")
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
}

var configCmd = &cobra.Command{
	Use:                "config",
	Short:              "get, set, unset configuration values",
	Example:            "  awless config        # list all your config\n  awless config set aws.region eu-west-1\n  awless config unset instance.count",
	PersistentPreRunE:  initAwlessEnvHook,
	PersistentPostRunE: notifyOnRegionOrProfilePrecedenceHook,

	Run: func(cmd *cobra.Command, args []string) {
		if keysOnly {
			for k := range config.Config {
				fmt.Println(k)
			}
			for k := range config.Defaults {
				fmt.Println(k)
			}
		} else {
			fmt.Println(config.DisplayConfig())
		}
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Get a configuration value",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("not enough parameters")
		}

		d, ok := config.Get(args[0])
		if !ok {
			fmt.Println("this parameter has not been set")
		} else {
			fmt.Printf("%v\n", d)
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:               "set KEY [VALUE]",
	Short:             "Set or update a configuration value",
	PersistentPreRun:  applyHooks(initAwlessEnvHook),
	PersistentPostRun: applyHooks(includeHookIf(&config.TriggerSyncOnConfigUpdate, initCloudServicesHook)),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("not enough parameters")
		}
		switch len(args) {
		case 0:
			return fmt.Errorf("not enough parameters")
		case 1:
			exitOn(config.InteractiveSet(strings.TrimSpace(args[0])))
		default:
			exitOn(config.Set(strings.TrimSpace(args[0]), strings.TrimSpace(args[1])))
		}

		return nil
	},
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset KEY",
	Short: "Unset a configuration value",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("not enough parameters")
		}
		_, ok := config.Get(args[0])
		if !ok {
			fmt.Println("this parameter has not been set")
		} else {

			exitOn(config.Unset(args[0]))
		}
		return nil
	},
}
