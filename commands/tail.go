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
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws/tailers"
)

var tailFollowFrequencyFlag time.Duration
var tailEnableFollowFlag bool
var tailNumberEventsFlag int

func init() {
	RootCmd.AddCommand(tailCmd)

	tailCmd.PersistentFlags().DurationVar(&tailFollowFrequencyFlag, "frequency", 10*time.Second, "Fetch refresh frequency")
	tailCmd.PersistentFlags().BoolVar(&tailEnableFollowFlag, "follow", false, "Periodically refresh and append new data to output")
	tailCmd.PersistentFlags().IntVarP(&tailNumberEventsFlag, "number", "n", 10, "Number of events to display")

	tailCmd.AddCommand(scalingActivitiesCmd)
}

var tailCmd = &cobra.Command{
	Use:               "tail",
	Hidden:            true,
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook),
	Short:             "Tail cloud events",
}

var scalingActivitiesCmd = &cobra.Command{
	Use:   "scaling-activities",
	Short: "Watch scaling-activities",

	Run: func(cmd *cobra.Command, args []string) {
		exitOn(awstailers.NewScalingActivitiesTailer(tailNumberEventsFlag, tailEnableFollowFlag, tailFollowFrequencyFlag).Tail(os.Stdout))
	},
}
