package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/api"
)

var accessApi *api.Access

func init() {
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
		fmt.Println(accessApi.Users())
	},
}

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(accessApi.Groups())
	},
}

var listRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List roles",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(accessApi.Roles())
	},
}

var listPoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List policies",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(accessApi.Policies())
	},
}
