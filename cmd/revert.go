package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"time"
	"strings"

	"github.com/oklog/ulid"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
)

var (
	revertFromIdFlag string
	flushHistoryFlag bool
	listHistoryFlag  bool
)

func init() {
	RootCmd.AddCommand(revertCmd)

	revertCmd.Flags().StringVarP(&revertFromIdFlag, "id", "i", "", "Template id to revert operation from")
	revertCmd.Flags().BoolVarP(&listHistoryFlag, "list", "l", false, "List all entries from operations history")
	revertCmd.Flags().BoolVarP(&flushHistoryFlag, "flush", "f", false, "Remove all entries from operations history")
}

var revertCmd = &cobra.Command{
	Use:               "revert",
	Short:             "List the history of your template action and revert them from an ID",
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,

	RunE: func(c *cobra.Command, args []string) error {
		if flushHistoryFlag {
			db, dbclose := database.Current()
			err := db.DeleteTemplateExecutions()
			dbclose()
			exitOn(err)
			return nil
		}

		if listHistoryFlag {
			listHistory()
			return nil
		}

		return errors.New("no flags given")
	},
}

func listHistory() {
	db, dbclose := database.Current()
	all, err := db.GetTemplateExecutions()
	dbclose()
	exitOn(err)

	for _, templ := range all {
		var buff bytes.Buffer

		buff.WriteByte('\n')
		for _, done := range templ.Executed {
			line := fmt.Sprintf("\t%s", done.Line)
			if done.Err != "" {
				buff.WriteString(renderRedFn(line))
				buff.WriteByte('\n')
				errMsg := strings.Replace(done.Err, "\n", "\t", -1)
				buff.WriteString(fmt.Sprintf("\t\t%s", errMsg))
			} else {
				buff.WriteString(renderGreenFn(line))
			}
			buff.WriteByte('\n')
		}

		uid, err := ulid.Parse(templ.ID)
		exitOn(err)

		date := time.Unix(int64(uid.Time())/int64(1000), time.Nanosecond.Nanoseconds())

		fmt.Printf("Date: %s, Id: %s%s\n", date.Format(time.Stamp), templ.ID, buff.String())
	}
}
