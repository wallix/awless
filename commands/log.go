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
)

var (
	deleteAllLogsFlag             bool
	deleteFromIdLogsFlag          string
	limitLogCountFlag             int
	rawJSONLogFlag, idOnlyLogFlag bool
	fullLogFlag, shortLogFlag     bool
)

func init() {
	RootCmd.AddCommand(logCmd)

	logCmd.Flags().BoolVar(&deleteAllLogsFlag, "delete-all", false, "Delete all logs from local db")
	logCmd.Flags().StringVar(&deleteFromIdLogsFlag, "delete", "", "Delete a specifc log entry given its id")
	logCmd.Flags().IntVarP(&limitLogCountFlag, "number", "n", 0, "Limit log output to the last n logs")
	logCmd.Flags().BoolVar(&rawJSONLogFlag, "raw", false, "Display logs as raw json with template context info, usually for debug")
	logCmd.Flags().BoolVar(&shortLogFlag, "short", false, "Display one or more template log with less info")
	logCmd.Flags().BoolVar(&fullLogFlag, "full", false, "Display template logs with full info")
	logCmd.Flags().BoolVar(&idOnlyLogFlag, "id-only", false, "Show only log template IDs (i.e. revert IDs)")
}

var logCmd = &cobra.Command{
	Use:               "log [REVERTID]",
	Short:             "Show all awless template actions against your cloud infrastructure",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade),

	RunE: func(c *cobra.Command, args []string) error {
		var all []*database.LoadedTemplate

		printer := getPrinter(args)

		if len(args) > 0 {
			exitOn(database.Execute(func(db *database.DB) error {
				single, err := db.GetLoadedTemplate(args[0])
				if err != nil {
					return err
				}
				all = append(all, single)
				return nil
			}))
			print(all, printer)
			return nil
		}

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

		exitOn(database.Execute(func(db *database.DB) (dberr error) {
			all, dberr = db.ListTemplates()
			return
		}))

		print(all, printer)
		return nil
	},
}

func print(all []*database.LoadedTemplate, printer logPrinter) {
	if limitLogCountFlag > 0 && limitLogCountFlag < len(all) {
		all = all[len(all)-limitLogCountFlag:]
	}

	for i, loaded := range all {
		if loaded.Err != nil {
			logger.Errorf("Template '%s' in error: %s", string(loaded.Key), loaded.Err)
			logger.Verbosef("Template raw content\n%s", loaded.Raw)
			fmt.Println()
			continue
		}

		if err := printer.print(loaded.TplExec); err != nil {
			logger.Error(err.Error())
		}

		if i < len(all)-1 {
			fmt.Println()
		}
	}

	if shortLogFlag {
		fmt.Println()
	}
}

func getPrinter(args []string) logPrinter {
	var defaultPrinter logPrinter
	if len(args) > 0 {
		defaultPrinter = &fullLogPrinter{os.Stdout}
	} else {
		defaultPrinter = &statLogPrinter{os.Stdout}
	}

	switch {
	case rawJSONLogFlag:
		return &rawJSONPrinter{os.Stdout}
	case idOnlyLogFlag:
		return &idOnlyPrinter{os.Stdout}
	case shortLogFlag:
		return &shortLogPrinter{os.Stdout}
	case fullLogFlag:
		return &fullLogPrinter{os.Stdout}
	default:
		return defaultPrinter
	}
}
