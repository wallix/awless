package commands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
)

func init() {
	RootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:                "log",
	Short:              "Logs all your awless template executions",
	PersistentPreRun:   applyHooks(initAwlessEnvHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(c *cobra.Command, args []string) error {
		db, err, dbclose := database.Current()
		exitOn(err)
		all, err := db.ListTemplateExecutions()
		dbclose()
		exitOn(err)

		for _, templ := range all {
			var buff bytes.Buffer

			for _, done := range templ.Executed {
				line := fmt.Sprintf("\t%s", done.Line)
				if done.Err != "" {
					buff.WriteString(renderRedFn(line))
					buff.WriteByte('\n')
					buff.WriteString(formatMultiLineErrMsg(done.Err))
				} else {
					buff.WriteString(renderGreenFn(line))
				}
				buff.WriteByte('\n')
			}

			uid, err := ulid.Parse(templ.ID)
			exitOn(err)

			date := time.Unix(int64(uid.Time())/int64(1000), time.Nanosecond.Nanoseconds())

			fmt.Printf("Date: %s\n", date.Format(time.Stamp))
			if templ.IsRevertible() {
				fmt.Printf("Revert id: %s\n", templ.ID)
			} else {
				fmt.Println("Revert id: <not revertible>")
			}
			fmt.Println(buff.String())
		}

		return nil
	},
}

func formatMultiLineErrMsg(msg string) string {
	notabs := strings.Replace(msg, "\t", "", -1)
	var indented []string
	for _, line := range strings.Split(notabs, "\n") {
		indented = append(indented, fmt.Sprintf("\t\t%s", line))
	}
	return strings.Join(indented, "\n")
}
