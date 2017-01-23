package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
)

func init() {
	RootCmd.AddCommand(whoamiCmd)
}

var whoamiCmd = &cobra.Command{
	Use:               "whoami",
	Aliases:           []string{"who"},
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,
	Short:             "Show the caller identity",

	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := aws.AccessService.CallerIdentity()
		displayItem(resp, err)
		return nil
	},
}
