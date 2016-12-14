package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
)

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
	listCmd.AddCommand(listVpcsCmd)
	listCmd.AddCommand(listSubnetsCmd)
	listCmd.AddCommand(listInstancesCmd)
	listCmd.AddCommand(listImagesCmd)

	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List various type of items: users, groups, instances, ...",
}

// access

var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.AccessService.Users()
		display(resp, err, displayFormat)
	},
}

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.AccessService.Groups()
		display(resp, err, displayFormat)
	},
}

var listRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List roles",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.AccessService.Roles()
		display(resp, err, displayFormat)
	},
}

var listPoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List policies",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.AccessService.LocalPolicies()
		display(resp, err, displayFormat)
	},
}

// infra

var listRegionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List regions",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Regions()
		display(resp, err, displayFormat)
	},
}

var listVpcsCmd = &cobra.Command{
	Use:   "vpcs",
	Short: "List vpcs",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Vpcs()
		display(resp, err, displayFormat)
	},
}

var listSubnetsCmd = &cobra.Command{
	Use:   "subnets",
	Short: "List subnets",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Subnets()
		display(resp, err, displayFormat)
	},
}

var listInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List instances",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Instances()
		display(resp, err, displayFormat)
	},
}

var listImagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List images",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Images()
		display(resp, err, displayFormat)
	},
}
