package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
)

func init() {
	RootCmd.AddCommand(whoamiCmd)
}

var whoamiCmd = &cobra.Command{
	Use:     "whoami",
	Aliases: []string{"who"},
	Short:   "Show the caller identity",

	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := aws.AccessService.CallerIdentity()
		display(resp, err)
		return nil
	},
}
