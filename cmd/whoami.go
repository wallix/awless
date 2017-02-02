package cmd

import (
	"fmt"

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

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.SecuService.CallerIdentity()
		exitOn(err)
		fmt.Println(resp)
	},
}
