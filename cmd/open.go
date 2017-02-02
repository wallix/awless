package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
)

func init() {
	RootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open your AWS console in your default browser",

	RunE: func(c *cobra.Command, args []string) error {
		console := fmt.Sprintf("https://%s.console.aws.amazon.com/console/home", database.MustGetDefaultRegion())

		var verb string
		switch runtime.GOOS {
		case "darwin":
			verb = "open"
		default:
			verb = "xdg-open"
		}

		var stderr bytes.Buffer
		cmd := exec.Command(verb, console)
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil || stderr.String() != "" {
			return fmt.Errorf("%s:%s", err, stderr.String())
		}
		return nil
	},
}
