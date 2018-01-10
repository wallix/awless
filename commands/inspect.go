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
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/inspect"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

var (
	inspectorFlag string
)

func init() {
	RootCmd.AddCommand(inspectCmd)

	inspectCmd.Flags().StringVarP(&inspectorFlag, "inspector", "i", "", "Indicates which inspector to run")
}

var inspectCmd = &cobra.Command{
	Use:               "inspect",
	Short:             "Analyze your infrastructure through inspectors",
	Long:              fmt.Sprintf("Basic proof of concept inspectors to analyze your infrastructure: %s", allInspectors()),
	Example:           "  awless inspect -i bucket_sizer\n  awless inspect -i pricer\n  awless inspect -i port_scanner",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),

	RunE: func(c *cobra.Command, args []string) error {
		inspector, ok := inspect.InspectorsRegister[inspectorFlag]
		if !ok {
			return fmt.Errorf("command needs a valid inspector: %s", allInspectors())
		}

		if !localGlobalFlag {
			logger.Info("Running full sync before inspection (disable it with --local flag)\n")
			var services []cloud.Service
			for _, srv := range cloud.ServiceRegistry {
				services = append(services, srv)
			}

			if _, err := sync.DefaultSyncer.Sync(services...); err != nil {
				logger.Verbose(err)
			}
		}

		g, err := sync.LoadLocalGraphs(config.GetAWSProfile(), config.GetAWSRegion())
		exitOn(err)

		exitOn(inspector.Inspect(g))

		inspector.Print(os.Stdout)

		return nil
	},
}

func allInspectors() string {
	var all []string
	for name := range inspect.InspectorsRegister {
		all = append(all, name)
	}
	return strings.Join(all, ", ")
}
