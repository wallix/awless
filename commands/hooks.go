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
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

func applyHooks(funcs ...func(*cobra.Command, []string) error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		for _, fn := range funcs {
			if err := fn(cmd, args); err != nil {
				exitOn(err)
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
	if localFlag {
		return nil
	}
	db, err, dbclose := database.Current()
	if err != nil {
		return fmt.Errorf("init cloud service: database error: %s", err)
	}
	profile, _ := db.GetDefaultString(database.ProfileKey)
	region := db.MustGetDefaultRegion()
	dbclose()

	if err := aws.InitServices(region, profile); err != nil {
		return err
	}

	return nil
}

func initConfigStruct(cmd *cobra.Command, args []string) error {
	return config.LoadConfig()
}

func initSyncerHook(cmd *cobra.Command, args []string) error {
	sync.DefaultSyncer = sync.NewSyncer()
	sync.DefaultSyncer.SetLogger(logger.DefaultLogger)
	return nil
}

func initLoggerHook(cmd *cobra.Command, args []string) error {
	var flag int
	if verboseFlag {
		flag = logger.VerboseF
	}
	if extraVerboseFlag {
		flag = flag | logger.ExtraVerboseF
	}

	logger.DefaultLogger.SetVerbose(flag)
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

func verifyNewVersionHook(cmd *cobra.Command, args []string) error {
	config.VerifyNewVersionAvailable("https://updates.awless.io", os.Stderr)
	return nil
}
