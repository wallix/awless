package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration",

	RunE: func(cmd *cobra.Command, args []string) error {
		if file, err := ioutil.ReadFile(viper.ConfigFileUsed()); err != nil {
			return fmt.Errorf("config: %s", err)
		} else {
			fmt.Printf("%s", file)
			fmt.Printf("(config at %s)\n", viper.ConfigFileUsed())
			return nil
		}
	},
}
