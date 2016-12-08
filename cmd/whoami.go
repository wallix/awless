package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wallix/awless/api"
)

func init() {
	RootCmd.AddCommand(whoamiCmd)
}

var whoamiCmd = &cobra.Command{
	Use:     "whoami",
	Aliases: []string{"who"},
	Short:   "Show the caller identity",

	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.AccessService.CallerIdentity()
		display(resp, err)
		return nil
	},
}
