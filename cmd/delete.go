package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.AddCommand(deleteAliasCmd)

	RootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete various type of resources by id",
}

var deleteAliasCmd = &cobra.Command{
	Use:   "alias [name]...",
	Short: "Delete alias",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Not enough args, need at least one alias id\n")
		}
		return statsDB.DeleteAlias(args...)
	},
}
