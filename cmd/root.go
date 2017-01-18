package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/rdf"
	"github.com/wallix/awless/stats"
)

var (
	verboseFlag bool
)

var RootCmd = &cobra.Command{
	Use:   "awless",
	Short: "Manage your cloud",
	Long:  "Awless is a powerful command line tool to inspect, sync and manage your infrastructure",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db, close := database.Current()
		defer close()
		db.AddHistoryCommand(append(strings.Split(cmd.CommandPath(), " "), args...))
	},
	BashCompletionFunction: bash_completion_func,
}

func InitCli() {
	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Turn on verbose mode for all commands")
	db, dbclose := database.Current()
	statsToSend := stats.CheckStatsToSend(db)
	dbclose()
	if statsToSend {
		go loadRdfsAndSendStats()
	}
}

func ExecuteRoot() error {
	err := RootCmd.Execute()
	db, dbclose := database.Current()
	defer dbclose()
	if err != nil && db != nil {
		db.AddLog(err.Error())
	}

	return err
}

func loadRdfsAndSendStats() {
	var err error
	localInfra, localAccess := rdf.NewGraph(), rdf.NewGraph()
	if !config.AwlessFirstSync {
		localInfra, err = rdf.NewGraphFromFile(filepath.Join(config.RepoDir, config.InfraFilename))
		if err != nil {
			db, dbclose := database.Current()
			db.AddLog(err.Error())
			dbclose()
		}
		localAccess, err = rdf.NewGraphFromFile(filepath.Join(config.RepoDir, config.AccessFilename))
		if err != nil {
			db, dbclose := database.Current()
			db.AddLog(err.Error())
			dbclose()
		}
	}
	db, dbclose := database.Current()
	if err := stats.SendStats(db, localInfra, localAccess); err != nil {
		db.AddLog(err.Error())
	}
	dbclose()
}

const (
	bash_completion_func = `
__awless_get_all_ids()
{
		local all_ids_output
		if all_ids_output=$(awless list all --local --ids --infra --access 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${all_ids_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_conf_keys()
{
		local all_keys_output
		if all_keys_output=$(awless config list --keys 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${all_keys_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_groups_ids()
{
		local ids_output
		if ids_output=$(awless list groups --local --ids 2>/dev/null); then
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
__awless_get_policies_ids()
{
		local ids_output
		if ids_output=$(awless list policies --local --ids 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${ids_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_roles_ids()
{
		local ids_output
		if ids_output=$(awless list roles --local --ids 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${ids_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_subnets_ids()
{
		local ids_output
		if ids_output=$(awless list subnets --local --ids 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${ids_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_users_ids()
{
		local ids_output
		if ids_output=$(awless list users --local --ids 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${ids_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_vpcs_ids()
{
		local ids_output
		if ids_output=$(awless list vpcs --local --ids 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${ids_output[*]}" -- "$cur" ) )
		fi
}
__custom_func() {
    case ${last_command} in
				awless_ssh )
            __awless_get_instances_ids
            return
            ;;
				awless_show_group )
            __awless_get_groups_ids
            return
            ;;
				awless_show_instance )
            __awless_get_instances_ids
            return
            ;;
				awless_show_policy )
            __awless_get_policies_ids
            return
            ;;
				awless_show_role )
            __awless_get_roles_ids
            return
            ;;
				awless_show_subnet )
						__awless_get_subnets_ids
						return
						;;
				awless_show_user )
						__awless_get_users_ids
						return
						;;
				awless_show_vpc )
						__awless_get_vpcs_ids
						return
						;;
				awless_config_set )
						__awless_get_conf_keys
						return
						;;
				awless_config_get )
						__awless_get_conf_keys
						return
						;;
				awless_config_unset )
						__awless_get_conf_keys
						return
						;;
        *)
            ;;
    esac
}`
)
