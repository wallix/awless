package commands

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
)

func init() {
	RootCmd.AddCommand(revertCmd)
}

var revertCmd = &cobra.Command{
	Use:                "revert",
	Short:              "Revert an template action given an revert ID (see `awless log` to list revert ids)",
	PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("revert id required (see `awless log` to list revert ids)")
		}

		revertId := args[0]

		db, err, dbclose := database.Current()
		exitOn(err)
		tplExec, err := db.GetTemplateExecution(revertId)
		dbclose()
		exitOn(err)

		reverted, err := tplExec.Revert()
		exitOn(err)

		exitOn(runTemplate(reverted, getCurrentDefaults()))

		return nil
	},
}
