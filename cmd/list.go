package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/api"
)

var (
	accessApi *api.Access

	displayFormat string
)

func init() {
	listCmd.PersistentFlags().StringVarP(&displayFormat, "format", "f", "line", "Display entities as raw in the console")

	listCmd.AddCommand(listUsersCmd)
	listCmd.AddCommand(listGroupsCmd)
	listCmd.AddCommand(listRolesCmd)
	listCmd.AddCommand(listPoliciesCmd)

	RootCmd.AddCommand(listCmd)

	var err error

	if accessApi, err = api.NewAccess(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to init the api: %s\n", err)
		os.Exit(-1)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List various type of items: users, groups, instances, ...",
}

var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := accessApi.Users()
		display(displayFormat, resp, err)
	},
}

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := accessApi.Groups()
		display(displayFormat, resp, err)
	},
}

var listRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List roles",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := accessApi.Roles()
		display(displayFormat, resp, err)
	},
}

var listPoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List policies",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := accessApi.Policies()
		display(displayFormat, resp, err)
	},
}
