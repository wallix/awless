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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wallix/awless-scheduler/client"
	"github.com/wallix/awless-scheduler/model"
	"github.com/wallix/awless/config"
)

var (
	listSchedulerTasksFlag    bool
	listSchedulerFailuresFlag bool
)

func init() {
	RootCmd.AddCommand(schedulerCmd)

	schedulerCmd.Flags().BoolVar(&listSchedulerTasksFlag, "list-tasks", false, "List scheduler tasks")
	schedulerCmd.Flags().BoolVar(&listSchedulerFailuresFlag, "list-failures", false, "List scheduler failures")
}

var schedulerCmd = &cobra.Command{
	Use:               "scheduler",
	PersistentPreRun:  applyHooks(initAwlessEnvHook, initLoggerHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade),
	Hidden:            true,
	Short:             "Accessing the scheduler API (when installed). To schedule templates runs/reverts use `awless run --schedule`",

	Run: func(cmd *cobra.Command, args []string) {
		if config.GetSchedulerURL() == "" {
			exitOn(errors.New("no scheduler URL in configuration. Set it with `awless config set scheduler.url`"))
		}
		cli, err := client.New(config.GetSchedulerURL())
		exitOn(err)

		if listSchedulerTasksFlag {
			tasks, err := cli.ListTasks()
			exitOn(err)
			printTasks(tasks)
			return
		}

		if listSchedulerFailuresFlag {
			tasks, err := cli.ListFailures()
			exitOn(err)
			printTasks(tasks)
			return
		}

		info := cli.ServiceInfo()
		fmt.Printf("Scheduler up!\nAddress: '%s'\nTickerFrequency: %s\nUptime: %s\n", info.ServiceAddr, info.TickerFrequency, info.Uptime)
	},
}

func printTasks(tasks []*model.Task) {
	for _, t := range tasks {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("Region: %s", t.Region))
		buf.WriteString(fmt.Sprintf(", RunAt: %s", t.RunAt))
		if !t.RevertAt.IsZero() {
			buf.WriteString(fmt.Sprintf(", RevertAt: %s", t.RevertAt))
		}
		buf.WriteString(fmt.Sprintf("\nContent: %s\n", t.Content))
		fmt.Println(buf.String())
	}
}
