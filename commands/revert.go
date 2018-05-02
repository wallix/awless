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

	"github.com/spf13/cobra"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template"
)

func init() {
	RootCmd.AddCommand(revertCmd)
}

var revertCmd = &cobra.Command{
	Use:               "revert REVERTID",
	Short:             "Revert a template from a revert ID (see `awless log`). If deployment has changed there is no guarantee that it is still revertible.",
	Example:           "  awless revert 01BA7RV6ES86PZYCM3H28WM6KZ",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),

	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("REVERTID required (see `awless log` to list revert ids)")
		}

		revertID := args[0]

		var loaded *template.TemplateExecution
		exitOn(database.Execute(func(db *database.DB) (terr error) {
			loaded, terr = db.GetTemplate(revertID)
			return
		}))

		if loc := loaded.Locale; loc != "" && loc != config.GetAWSRegion() {
			logger.Errorf("This template was originally run in region %s", loc)
			logger.Infof("Revert with `awless revert %s -r %s -p %s`", revertID, loc, loaded.Profile)
			os.Exit(1)
		}

		if prof := loaded.Profile; prof != config.GetAWSProfile() {
			logger.Warningf("This template was originally run with profile %s", prof)
		}

		reverted, err := loaded.Template.Revert()
		exitOn(err)

		tplExec := &template.TemplateExecution{
			Template: reverted,
			Locale:   config.GetAWSRegion(),
			Profile:  config.GetAWSProfile(),
			Source:   reverted.String(),
		}
		tplExec.SetMessage(fmt.Sprintf("Revert %s: %s", loaded.ID, loaded.Message))

		exitOn(NewRunnerRequiredParamsOnly(tplExec.Template, tplExec.Message, tplExec.Path).Run())

		return nil
	},
}
