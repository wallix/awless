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
	listCmd.AddCommand(listVpcsCmd)
	listCmd.AddCommand(listSubnetsCmd)

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
		display(resp, err, displayFormat)
	},
}

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := accessApi.Groups()
		display(resp, err, displayFormat)
	},
}

var listRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List roles",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := accessApi.Roles()
		display(resp, err, displayFormat)
	},
}

var listPoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List policies",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := accessApi.Policies()
		display(resp, err, displayFormat)
	},
}

var listRegionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List regions",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := infraApi.Regions()
		display(resp, err)
	},
}

var listInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List instances",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := infraApi.Instances()
		display(resp, err)
	},
}

var listVpcsCmd = &cobra.Command{
	Use:   "vpcs",
	Short: "List vpcs",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := infraApi.Vpcs()
		display(resp, err)
	},
}

var listSubnetsCmd = &cobra.Command{
	Use:   "subnets",
	Short: "List subnets",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := infraApi.Subnets()
		display(resp, err)
	},
}
