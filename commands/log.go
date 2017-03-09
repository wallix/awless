/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/template"
)

var (
	logPorcelainFlag bool
	deleteLogsFlag   bool
)

func init() {
	RootCmd.AddCommand(logCmd)

	logCmd.Flags().BoolVarP(&logPorcelainFlag, "porcelain", "p", false, "Format for machine consumption")
	logCmd.Flags().BoolVarP(&deleteLogsFlag, "delete", "d", false, "Delete all logs from local db")
}

var logCmd = &cobra.Command{
	Use:               "log",
	Short:             "Logs all executions done against your cloud",
	PersistentPreRun:  applyHooks(initAwlessEnvHook),
	PersistentPostRun: applyHooks(saveHistoryHook, verifyNewVersionHook),

	RunE: func(c *cobra.Command, args []string) error {
		db, err, dbclose := database.Current()
		exitOn(err)

		if deleteLogsFlag {
			db.DeleteTemplateExecutions()
			return nil
		}

		all, err := db.ListTemplateExecutions()
		dbclose()
		exitOn(err)

		for _, templ := range all {
			var buff bytes.Buffer

			if logPorcelainFlag {
				formatForMachine(&buff, templ)
			} else {
				formatForHuman(&buff, templ)
			}

			fmt.Println(buff.String())
		}

		return nil
	},
}

func formatForMachine(buff *bytes.Buffer, templ *template.TemplateExecution) {
	sep := '\t'

	buff.WriteString(parseULIDDate(templ.ID))
	buff.WriteRune(sep)
	if templ.IsRevertible() {
		buff.WriteString(templ.ID)
	} else {
		buff.WriteString("<not revertible>")
	}
	buff.WriteByte('\n')
	for _, done := range templ.Executed {
		if done.Err != "" {
			buff.WriteString("KO")
		} else {
			buff.WriteString("OK")
		}
		buff.WriteRune(sep)
		buff.WriteString(done.Result)
		buff.WriteRune(sep)
		buff.WriteString(done.Line)
		buff.WriteByte('\n')
	}
}

func formatForHuman(buff *bytes.Buffer, templ *template.TemplateExecution) {
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

	fmt.Printf("Date: %s\n", parseULIDDate(templ.ID))
	if templ.IsRevertible() {
		fmt.Printf("Revert id: %s\n", templ.ID)
	} else {
		fmt.Println("Revert id: <not revertible>")
	}
}

func parseULIDDate(uid string) string {
	parsed, err := ulid.Parse(uid)
	exitOn(err)

	date := time.Unix(int64(parsed.Time())/int64(1000), time.Nanosecond.Nanoseconds())

	return date.Format(time.Stamp)
}

func formatMultiLineErrMsg(msg string) string {
	notabs := strings.Replace(msg, "\t", "", -1)
	var indented []string
	for _, line := range strings.Split(notabs, "\n") {
		indented = append(indented, fmt.Sprintf("\t\t%s", line))
	}
	return strings.Join(indented, "\n")
}
