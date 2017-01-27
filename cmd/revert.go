package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
)

var (
	revertFromIdFlag string
)

func init() {
	RootCmd.AddCommand(revertCmd)

	revertCmd.Flags().StringVarP(&revertFromIdFlag, "id", "i", "", "Template id to revert action from")
}

var revertCmd = &cobra.Command{
	Use:               "revert",
	Short:             "List the history of yoru template action and revert them from an ID",
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,

	RunE: func(c *cobra.Command, args []string) error {
		db, dbclose := database.Current()
		all, err := db.GetTemplateOperations()
		dbclose()
		exitOn(err)

		for _, a := range all {
			fmt.Println(a)
		}

		return nil
	},
}
