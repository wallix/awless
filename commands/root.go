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
)

var (
	verboseGlobalFlag      bool
	extraVerboseGlobalFlag bool
	silentGlobalFlag       bool
	localGlobalFlag        bool
	noSyncGlobalFlag       bool
	forceGlobalFlag        bool
	versionGlobalFlag      bool
	awsRegionGlobalFlag    string
	awsProfileGlobalFlag   string
	awsColorGlobalFlag     string
	networkMonitorFlag     bool

	renderGreenFn    = color.New(color.FgGreen).SprintFunc()
	renderRedFn      = color.New(color.FgRed).SprintFunc()
	renderYellowFn   = color.New(color.FgYellow).SprintFunc()
	renderBlueFn     = color.New(color.FgBlue).SprintFunc()
	renderCyanBoldFn = color.New(color.FgCyan, color.Bold).SprintFunc()
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verboseGlobalFlag, "verbose", "v", false, "Turn on verbose mode for all commands")
	RootCmd.PersistentFlags().BoolVarP(&extraVerboseGlobalFlag, "extra-verbose", "e", false, "Turn on extra verbose mode (including regular verbose) for all commands")
	RootCmd.PersistentFlags().BoolVar(&silentGlobalFlag, "silent", false, "Turn on silent mode for all commands: disable logging, etc...")
	RootCmd.PersistentFlags().BoolVarP(&localGlobalFlag, "local", "l", false, "Work offline only using locally synced resources")
	RootCmd.PersistentFlags().BoolVarP(&forceGlobalFlag, "force", "f", false, "Force the command and bypass confirmation prompts")
	RootCmd.PersistentFlags().BoolVar(&noSyncGlobalFlag, "no-sync", false, "Do not run any sync on command")
	RootCmd.PersistentFlags().StringVarP(&awsRegionGlobalFlag, "aws-region", "r", "", "Override AWS region temporarily for the current command")
	RootCmd.PersistentFlags().SetAnnotation("aws-region", cobra.BashCompCustom, []string{"__awless_region_list"})
	RootCmd.PersistentFlags().StringVarP(&awsProfileGlobalFlag, "aws-profile", "p", "", "Override AWS profile temporarily for the current command")
	RootCmd.PersistentFlags().SetAnnotation("aws-profile", cobra.BashCompCustom, []string{"__awless_profile_list"})
	RootCmd.PersistentFlags().StringVar(&awsColorGlobalFlag, "color", "auto", "Force enabling/disabling colors in display (auto, never, always)")
	RootCmd.PersistentFlags().BoolVar(&networkMonitorFlag, "network-monitor", false, "Debug requests with network monitor")
	RootCmd.PersistentFlags().MarkHidden("network-monitor")

	RootCmd.Flags().BoolVar(&versionGlobalFlag, "version", false, "Print awless version")

	cobra.AddTemplateFunc("IsCmdAnnotatedOneliner", IsCmdAnnotatedOneliner)
	cobra.AddTemplateFunc("HasCmdOnelinerChilds", HasCmdOnelinerChilds)

	RootCmd.SetUsageTemplate(customRootUsage)

	cobra.OnInitialize(func() {
		switch awsColorGlobalFlag {
		case "never":
			color.NoColor = true
		case "always":
			color.NoColor = false
		}
	})
}

var RootCmd = &cobra.Command{
	Use:   "awless COMMAND",
	Short: "Manage  and explore your cloud",
	Long:  "awless is a powerful CLI to explore, sync and manage your cloud infrastructure",
	BashCompletionFunction: bash_completion_func,
	RunE: func(c *cobra.Command, args []string) error {
		if versionGlobalFlag {
			printVersion(c, args)
			return nil
		}
		return c.Usage()
	},
}

const customRootUsage = `USAGE:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if gt .Aliases 0}}

ALIASES:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

EXAMPLES:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

COMMANDS:{{range .Commands}}{{ if not (IsCmdAnnotatedOneliner .Annotations)}}{{if .IsAvailableCommand }}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if HasCmdOnelinerChilds .}}

ONE-LINER TEMPLATE COMMANDS:{{range .Commands}}{{ if IsCmdAnnotatedOneliner .Annotations}}{{if .IsAvailableCommand }}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

FLAGS:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}

GLOBAL FLAGS:
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
				awless_switch )
						__awless_profile_region_list
						return
						;;
        *)
            ;;
    esac
}

__awless_region_list()
{
    cur="${COMP_WORDS[COMP_CWORD]#*=}"
    regions="us-east-1 us-east-2 us-west-1 us-west-2 ca-central-1 eu-west-1 eu-central-1 eu-west-2 eu-west-3 ap-northeast-1 ap-northeast-2 ap-southeast-1 ap-southeast-2 ap-south-1 sa-east-1"
    COMPREPLY=( $(compgen -W "${regions}" -- ${cur}) )
}

__awless_profile_list()
{
    cur="${COMP_WORDS[COMP_CWORD]#*=}"
    profiles="$((egrep '^\[ *[a-zA-Z0-9_-]+ *\]$' ~/.aws/credentials 2>/dev/null; grep '\[profile' ~/.aws/config 2>/dev/null | sed 's|\[profile ||g') | tr -d '[]' | sort | uniq)"
    COMPREPLY=( $(compgen -W "${profiles}" -- ${cur}) )
}

__awless_profile_region_list()
{
    cur="${COMP_WORDS[COMP_CWORD]#*=}"
		regions="us-east-1 us-east-2 us-west-1 us-west-2 ca-central-1 eu-west-1 eu-central-1 eu-west-2 eu-west-3 ap-northeast-1 ap-northeast-2 ap-southeast-1 ap-southeast-2 ap-south-1 sa-east-1"
    profiles="$((egrep '^\[ *[a-zA-Z0-9_-]+ *\]$' ~/.aws/credentials 2>/dev/null; grep '\[profile' ~/.aws/config 2>/dev/null | sed 's|\[profile ||g') | tr -d '[]' | sort | uniq)"
    COMPREPLY=( $(compgen -W "${profiles} ${regions}" -- ${cur}) )
}

`
)
