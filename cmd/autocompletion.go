package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(genautocompleteCmd)
}

var genautocompleteCmd = &cobra.Command{
	Use:   "genautocomplete filepath",
	Short: "Generate shell autocompletion script for Awless",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("path must be provided")
		}

		return RootCmd.GenBashCompletionFile(args[0])
	},
}
