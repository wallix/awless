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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/template"
)

var (
	deleteLogsFlag bool
)

func init() {
	RootCmd.AddCommand(logCmd)

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
			db.DeleteTemplates()
			return nil
		}

		all, err := db.ListTemplates()
		dbclose()
		exitOn(err)

		for _, templ := range all {
			printer := template.NewLogPrinter(os.Stdout)
			printer.RenderKO = renderRedFn
			printer.RenderOK = renderGreenFn

			printer.Print(templ)
			fmt.Println()
		}

		return nil
	},
}
