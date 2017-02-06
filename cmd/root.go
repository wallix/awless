package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
)

var (
	verboseFlag bool
	localFlag   bool
	versionFlag bool
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Turn on verbose mode for all commands")
	RootCmd.PersistentFlags().BoolVar(&localFlag, "local", false, "Work offline only with synced/local resources")
	RootCmd.Flags().BoolVar(&versionFlag, "version", false, "Print awless version")
}

var RootCmd = &cobra.Command{
	Use:   "awless",
	Short: "Manage your cloud",
	Long:  "Awless is a powerful command line tool to inspect, sync and manage your infrastructure",
	BashCompletionFunction: bash_completion_func,
	RunE: func(c *cobra.Command, args []string) error {
		if versionFlag {
			printVersion(c, args)
			return nil
		}
		return c.Usage()
	},
}

func ExecuteRoot() error {
	err := RootCmd.Execute()

	if err != nil {
		db, dbclose := database.Current()
		defer dbclose()
		if db != nil {
			db.AddLog(err.Error())
		}
	}

	return err
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
