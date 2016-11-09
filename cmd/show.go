package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	showCmd.AddCommand(showVpcCmd)

	RootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show various type of items by id: users, groups, instances, vpcs, ...",
}

var showVpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Show a vpc from a given id",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("show vpc: id required")
		}
		resp, err := infraApi.Vpc(args[0])
		display(resp, err)
		return nil
	},
}
