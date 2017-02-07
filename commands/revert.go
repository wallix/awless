package commands

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

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
	Use:                "revert",
	Short:              "List the history of your template action and revert them from an ID",
	PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(c *cobra.Command, args []string) error {
		if listHistoryFlag {
			listHistory()
			return nil
		}

		if flushHistoryFlag {
			fmt.Print("Are you sure you want to delete all your template executions history? (y/n): ")
			var yesorno string
			_, err := fmt.Scanln(&yesorno)
			exitOn(err)

			if strings.TrimSpace(yesorno) == "y" {
				db, dbclose := database.Current()
				err := db.DeleteTemplateExecutions()
				dbclose()
				exitOn(err)
			}

			return nil
		}

		if revertFromIdFlag != "" {
			db, dbclose := database.Current()
			tplExec, err := db.GetTemplateExecution(revertFromIdFlag)
			dbclose()
			exitOn(err)

			reverted, err := tplExec.Revert()
			exitOn(err)

			exitOn(runTemplate(reverted, getCurrentDefaults()))

			return nil
		}

		return errors.New("no flags given")
	},
}

func listHistory() {
	db, dbclose := database.Current()
	all, err := db.ListTemplateExecutions()
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

		if templ.IsRevertible() {
			fmt.Printf("Date: %s, Revert Id: %s%s\n", date.Format(time.Stamp), templ.ID, buff.String())
		} else {
			fmt.Printf("Date: %s%s\n", date.Format(time.Stamp), buff.String())
		}
	}
}
