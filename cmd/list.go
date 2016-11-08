package cmd

import "github.com/spf13/cobra"

var (
	displayFormat string
)

func init() {
	listCmd.PersistentFlags().StringVarP(&displayFormat, "format", "f", "line", "Display entities as raw in the console")

	// access
	listCmd.AddCommand(listUsersCmd)
	listCmd.AddCommand(listGroupsCmd)
	listCmd.AddCommand(listRolesCmd)
	listCmd.AddCommand(listPoliciesCmd)

	// infra
	listCmd.AddCommand(listRegionsCmd)
	listCmd.AddCommand(listInstancesCmd)

	RootCmd.AddCommand(listCmd)
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

var listInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List instances",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := infraApi.Instances()
		display(displayFormat, resp, err)
	},
}

var listRegionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List regions",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := infraApi.Regions()
		display(displayFormat, resp, err)
	},
}
