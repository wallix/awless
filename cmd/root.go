package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
	"github.com/wallix/awless/stats"
)

var (
	statsDB *stats.DB

	verboseFlag bool
)

const (
	bash_completion_func = `
__awless_get_all_ids()
{
		local all_ids_output
		if all_ids_output=$(awless list all --local --ids --infra --access 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${all_ids_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_alias_ids()
{
		local ids_output
		if ids_output=$(awless list aliases --local --ids 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${ids_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_instances_ids()
{
		local ids_output
		if ids_output=$(awless list instances --local --ids 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${ids_output[*]}" -- "$cur" ) )
		fi
}
__custom_func() {
    case ${last_command} in
        awless_create_alias )
            __awless_get_all_ids
            return
            ;;
				awless_ssh )
            __awless_get_instances_ids
            return
            ;;
				awless_delete_alias )
            __awless_get_alias_ids
            return
            ;;
        *)
            ;;
    esac
}`
)

var RootCmd = &cobra.Command{
	Use:   "awless",
	Short: "Manage your cloud",
	Long:  "Awless is a powerful command line tool to inspect, sync and manage your infrastructure",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if statsDB != nil {
			defer statsDB.Close()
			statsDB.AddHistoryCommand(append(strings.Split(cmd.CommandPath(), " "), args...))
		}
	},
	BashCompletionFunction: bash_completion_func,
}

func InitCli() {
	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Turn on verbose mode for all commands")

	var err error

	statsDB, err = stats.OpenDB(config.DatabasePath)
	if err != nil {
		if statsDB != nil {
			statsDB.AddLog("can not save history: " + err.Error())
		}
	} else if statsDB.CheckStatsToSend(config.StatsExpirationDuration) {
		publicKey, err := config.LoadPublicKey()
		if err != nil {
			statsDB.AddLog(err.Error())
		} else {
			if !config.AwlessFirstSync {
				go func() {
					localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
					if err != nil {
						statsDB.AddLog(err.Error())
					}
					localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
					if err != nil {
						statsDB.AddLog(err.Error())
					}

					if err := statsDB.SendStats(config.StatsServerUrl, *publicKey, localInfra, localAccess); err != nil {
						statsDB.AddLog(err.Error())
					}
				}()
			}
		}
	}
}

func ExecuteRoot() error {
	err := RootCmd.Execute()
	if err != nil && statsDB != nil {
		statsDB.AddLog(err.Error())
	}

	return err
}
