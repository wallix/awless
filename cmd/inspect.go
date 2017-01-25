package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/inspect"
)

func init() {
	RootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:               "inspect",
	Short:             "Experimental! Inspecting your infrastructure through any inspector",
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,

	RunE: func(c *cobra.Command, args []string) error {
		region := database.MustGetDefaultRegion()

		infra, err := aws.InfraService.FetchAwsInfra()
		exitOn(err)

		infrag, err := aws.BuildAwsInfraGraph(region, infra)
		exitOn(err)

		pricer := &inspect.Pricer{}
		err = pricer.Inspect(infrag)
		exitOn(err)

		pricer.Print(os.Stdout)

		return nil
	},
}
