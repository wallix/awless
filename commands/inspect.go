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
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/inspect"
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
	Use: "inspect",
	Short: fmt.Sprintf(
		"Inspecting your infrastructure using available inspectors below: %s", allInspectors(),
	),
	PersistentPreRun:   applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, verifyNewVersionHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(c *cobra.Command, args []string) error {
		inspector, ok := inspect.InspectorsRegister[inspectorFlag]
		if !ok {
			return fmt.Errorf("command needs a valid inspector: %s", allInspectors())
		}

		var graphs []*graph.Graph
		if localFlag {
			for _, name := range inspector.Services() {
				graphs = append(graphs, sync.LoadCurrentLocalGraph(name))
			}
		} else {
			var err error
			services := []cloud.Service{}
			for _, name := range inspector.Services() {
				srv, ok := cloud.ServiceRegistry[name]
				if !ok {
					return fmt.Errorf("unknown service %s for inspector %s", name, inspector.Name())
				}
				services = append(services, srv)
			}

			graphPerService, err := sync.DefaultSyncer.Sync(services...)
			exitOn(err)
			for _, g := range graphPerService {
				graphs = append(graphs, g)
			}
		}

		err := inspector.Inspect(graphs...)
		exitOn(err)

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
