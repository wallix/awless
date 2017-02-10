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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/stats"
	"github.com/wallix/awless/sync"
)

func applyHooks(funcs ...func(*cobra.Command, []string) error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		for _, fn := range funcs {
			if err := fn(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "command hook failed: %s\n", err)
				os.Exit(1)
			}
		}
	}
}

func initAwlessEnvHook(cmd *cobra.Command, args []string) error {
	if err := config.InitAwlessEnv(); err != nil {
		return fmt.Errorf("cannot init awless environment: %s", err)
	}
	return nil
}

func initCloudServicesHook(cmd *cobra.Command, args []string) error {
	region := os.Getenv("__AWLESS_REGION")
	if region == "" {
		return errors.New("region should be in env")
	}

	if err := aws.InitServices(region); err != nil {
		return err
	}

	if err := database.InitDB(); err != nil {
		db, err, closing := database.Current()
		if err == nil && db != nil {
			db.AddLog(fmt.Sprintf("cannot init database: %s", err))
		}
		closing()
	}

	return nil
}

func initSyncerHook(cmd *cobra.Command, args []string) error {
	sync.DefaultSyncer = sync.NewSyncer(dryRunSyncFlag)
	sync.DefaultSyncer.SetLogger(logger.DefaultLogger)
	return nil
}

func initLoggerHook(cmd *cobra.Command, args []string) error {
	logger.DefaultLogger.SetVerbose(verboseFlag)
	return nil
}

func saveHistoryHook(cmd *cobra.Command, args []string) error {
	db, err, close := database.Current()
	if err == nil && db != nil {
		db.AddHistoryCommand(append(strings.Split(cmd.CommandPath(), " "), args...))
		defer close()
	}
	return nil
}

func checkStatsHook(cmd *cobra.Command, args []string) error {
	db, err, dbclose := database.Current()
	if err != nil {
		return nil
	}
	statsToSend := stats.CheckStatsToSend(db)
	dbclose()

	if statsToSend {
		go func() {
			localInfra := sync.LoadCurrentLocalGraph(aws.InfraService.Name())
			localAccess := sync.LoadCurrentLocalGraph(aws.AccessService.Name())

			db, dberr, dbclose := database.Current()
			if dberr != nil {
				return
			}
			if err := stats.SendStats(db, localInfra, localAccess); err != nil {
				db.AddLog(err.Error())
			}
			dbclose()
		}()
	}

	return nil
}
