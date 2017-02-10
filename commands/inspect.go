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
	PersistentPreRun:   applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
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
	for name, _ := range inspect.InspectorsRegister {
		all = append(all, name)
	}
	return strings.Join(all, ", ")
}
