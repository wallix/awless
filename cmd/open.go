package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open your AWS console in your default browser",

	RunE: func(c *cobra.Command, args []string) error {
		console := fmt.Sprintf("https://%s.console.aws.amazon.com/console/home", viper.GetString("region"))

		var verb string
		switch runtime.GOOS {
		case "darwin":
			verb = "open"
		default:
			verb = "xdg-open"
		}

		cmd := exec.Command(verb, console)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	},
}
