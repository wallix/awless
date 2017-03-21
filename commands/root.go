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
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
)

var (
	verboseGlobalFlag      bool
	extraVerboseGlobalFlag bool
	localGlobalFlag        bool
	forceGlobalFlag        bool
	versionGlobalFlag      bool
	awsRegionGlobalFlag    string
	awsProfileGlobalFlag   string

	renderGreenFn    = color.New(color.FgGreen).SprintFunc()
	renderRedFn      = color.New(color.FgRed).SprintFunc()
	renderCyanBoldFn = color.New(color.FgCyan, color.Bold).SprintFunc()
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verboseGlobalFlag, "verbose", "v", false, "Turn on verbose mode for all commands")
	RootCmd.PersistentFlags().BoolVarP(&extraVerboseGlobalFlag, "extra-verbose", "e", false, "Turn on extra verbose mode (including regular verbose) for all commands")
	RootCmd.PersistentFlags().BoolVarP(&localGlobalFlag, "local", "l", false, "Work offline only with synced/local resources")
	RootCmd.PersistentFlags().BoolVarP(&forceGlobalFlag, "force", "f", false, "Force the command and bypass any confirmation prompt")
	RootCmd.PersistentFlags().StringVar(&awsRegionGlobalFlag, "aws-region", "", "Overwrite AWS region")
	RootCmd.PersistentFlags().StringVar(&awsProfileGlobalFlag, "aws-profile", "", "Overwrite AWS profile")
	RootCmd.Flags().BoolVar(&versionGlobalFlag, "version", false, "Print awless version")

	cobra.AddTemplateFunc("IsCmdAnnotatedOneliner", IsCmdAnnotatedOneliner)
	cobra.AddTemplateFunc("HasCmdOnelinerChilds", HasCmdOnelinerChilds)

	RootCmd.SetUsageTemplate(customRootUsage)
}

var RootCmd = &cobra.Command{
	Use:   "awless COMMAND",
	Short: "Manage  and explore your cloud",
	Long:  "awless is a powerful CLI to explore, sync and manage your infrastructure",
	BashCompletionFunction: bash_completion_func,
	RunE: func(c *cobra.Command, args []string) error {
		if versionGlobalFlag {
			printVersion(c, args)
			return nil
		}
		return c.Usage()
	},
}

func ExecuteRoot() error {
	err := RootCmd.Execute()

	if err != nil {
		db, dberr, dbclose := database.Current()
		if dberr == nil && db != nil {
			db.AddLog(err.Error())
			dbclose()
		}
	}

	return err
}

const customRootUsage = `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{ if not (IsCmdAnnotatedOneliner .Annotations)}}{{if .IsAvailableCommand }}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if HasCmdOnelinerChilds .}}

One-liner Template Commands:{{range .Commands}}{{ if IsCmdAnnotatedOneliner .Annotations}}{{if .IsAvailableCommand }}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func IsCmdAnnotatedOneliner(annot map[string]string) bool {
	if annot == nil {
		return false
	}
	_, ok := annot["one-liner"]
	return ok
}

func HasCmdOnelinerChilds(cmd *cobra.Command) bool {
	for _, child := range cmd.Commands() {
		if IsCmdAnnotatedOneliner(child.Annotations) {
			return true
		}
	}

	return false
}

const (
	bash_completion_func = `
__awless_get_all_ids()
{
		local all_ids_output
		if all_ids_output=$(awless list infra --local --ids 2>/dev/null; awless list access --local --ids 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${all_ids_output[*]}" -- "$cur" ) )
		fi
}
__awless_get_instances_ids()
{
		local all_ids_output
		if all_ids_output=$(awless list instances --local --ids 2>/dev/null); then
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

__custom_func() {
    case ${last_command} in
				awless_ssh )
            __awless_get_instances_ids
            return
            ;;
				awless_show )
            __awless_get_all_ids
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
