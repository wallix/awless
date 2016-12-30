package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete various type of resources by id",
}
