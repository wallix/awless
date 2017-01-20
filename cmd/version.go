package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/config"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show awless version",

	Run: printVersion,
}

func printVersion(*cobra.Command, []string) {
	fmt.Println("awless", config.CurrentBuildInfo.String())
}
