package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
)

func init() {
	createCmd.AddCommand(createInstanceCmd)

	RootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create various type of resources by id: users, groups, instances, vpcs, ...",
}

var createInstanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"inst", "i"},
	Short:   "Create an instance",

	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := aws.InfraService.CreateInstance("ami-9398d3e0")
		display(resp, err)

		return nil
	},
}
