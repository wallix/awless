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

	if err := database.InitDB(config.AwlessFirstInstall); err != nil {
		return fmt.Errorf("cannot init database: %s", err)
	}

	return nil
}

func initSyncerHook(cmd *cobra.Command, args []string) error {
	sync.DefaultSyncer = sync.NewSyncer()
	return nil
}

func initLoggerHook(cmd *cobra.Command, args []string) error {
	logger.DefaultLogger.SetVerbose(verboseFlag)
	return nil
}

func saveHistoryHook(cmd *cobra.Command, args []string) error {
	db, close := database.Current()
	defer close()
	db.AddHistoryCommand(append(strings.Split(cmd.CommandPath(), " "), args...))
	return nil
}

func checkStatsHook(cmd *cobra.Command, args []string) error {
	db, dbclose := database.Current()
	statsToSend := stats.CheckStatsToSend(db)
	dbclose()

	if statsToSend {
		go func() {
			localInfra := sync.LoadCurrentLocalGraph(aws.InfraService.Name())
			localAccess := sync.LoadCurrentLocalGraph(aws.AccessService.Name())

			db, dbclose := database.Current()
			if err := stats.SendStats(db, localInfra, localAccess); err != nil {
				db.AddLog(err.Error())
			}
			dbclose()
		}()
	}

	return nil
}
