package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	historyCmd.AddCommand(historyFlushCmd)
	historyCmd.AddCommand(historyShowCmd)

	RootCmd.AddCommand(historyCmd)
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Display or delete the history of the command lines entered in awless",
}

var historyFlushCmd = &cobra.Command{
	Use:   "delete",
	Short: "Empty the history",

	RunE: func(cmd *cobra.Command, args []string) error {
		if db != nil {
			err := db.EmptyHistory()
			if err != nil {
				return err
			}
			return nil
		} else {
			fmt.Fprintln(os.Stderr, "The history database has not been loaded")
			return nil
		}
	},
}

var historyShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the history",

	RunE: func(cmd *cobra.Command, args []string) error {
		if db != nil {
			lines, err := db.GetHistory(0)
			if err != nil {
				return err
			}
			if len(lines) == 0 {
				fmt.Println("There is no line in the awless history")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
			for _, line := range lines {
				fmt.Fprintln(w, line.Time.Format(time.RFC822), "\t", strings.Join(line.Command, " "))
			}
			w.Flush()

			return nil
		} else {
			fmt.Fprintln(os.Stderr, "The history database has not been loaded")
			return nil
		}
	},
}
