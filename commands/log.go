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
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template"
)

var (
	deleteAllLogsFlag    bool
	deleteFromIdLogsFlag string
	logsAsRawJSONFlag    bool
)

func init() {
	RootCmd.AddCommand(logCmd)

	logCmd.Flags().BoolVar(&deleteAllLogsFlag, "delete-all", false, "Delete all logs from local db")
	logCmd.Flags().StringVar(&deleteFromIdLogsFlag, "delete", "", "Delete a specifc log entry given its id")
	logCmd.Flags().BoolVar(&logsAsRawJSONFlag, "raw-json", false, "Display logs as raw json")
}

var logCmd = &cobra.Command{
	Use:               "log",
	Short:             "Shows the cloud infrastructure changes log",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook),

	RunE: func(c *cobra.Command, args []string) error {
		if deleteAllLogsFlag {
			exitOn(database.Execute(func(db *database.DB) error {
				return db.DeleteTemplates()
			}))
			return nil
		}

		if tid := deleteFromIdLogsFlag; tid != "" {
			exitOn(database.Execute(func(db *database.DB) error {
				return db.DeleteTemplate(tid)
			}))
			return nil
		}

		var all []*database.LoadedTemplate
		err := database.Execute(func(db *database.DB) (dberr error) {
			all, dberr = db.ListTemplates()
			return
		})
		exitOn(err)

		if logsAsRawJSONFlag {
			for _, loaded := range all {
				if loaded.Err != nil {
					logger.Errorf("Template '%s' in error: %s", string(loaded.Key), loaded.Err)
					logger.Verbosef("Template raw content\n%s", loaded.Raw)
					fmt.Println()
					continue
				}

				b, err := loaded.TplExec.MarshalJSON()
				if err != nil {
					logger.Errorf("cannot display template as raw json: %s", err)
				}
				fmt.Println(string(b))
				fmt.Println()
			}
			return nil
		}

		printer := template.NewLogPrinter(os.Stdout)
		printer.RenderKO = renderRedFn
		printer.RenderOK = renderGreenFn

		for _, loaded := range all {
			if loaded.Err != nil {
				logger.Errorf("Template '%s' in error: %s", string(loaded.Key), loaded.Err)
				logger.Verbosef("Template raw content\n%s", loaded.Raw)
				fmt.Println()
				continue
			}

			printer.Print(loaded.TplExec)
			fmt.Println()
		}

		return nil
	},
}
