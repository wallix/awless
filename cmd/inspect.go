package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/inspect"
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
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,

	RunE: func(c *cobra.Command, args []string) error {
		inspector, ok := inspect.InspectorsRegister[inspectorFlag]
		if !ok {
			return fmt.Errorf("command needs a valid inspector: %s", allInspectors())
		}

		region := database.MustGetDefaultRegion()

		infra, err := aws.InfraService.FetchAwsInfra()
		exitOn(err)

		infrag, err := aws.BuildAwsInfraGraph(region, infra)
		exitOn(err)

		err = inspector.Inspect(infrag)
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
