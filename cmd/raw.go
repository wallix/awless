package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
)

func init() {
	// access
	listRawCmd.AddCommand(listRawUsersCmd)
	listRawCmd.AddCommand(listRawGroupsCmd)
	listRawCmd.AddCommand(listRawRolesCmd)
	listRawCmd.AddCommand(listRawPoliciesCmd)

	// infra
	listRawCmd.AddCommand(listRawRegionsCmd)
	listRawCmd.AddCommand(listRawVpcsCmd)
	listRawCmd.AddCommand(listRawSubnetsCmd)
	listRawCmd.AddCommand(listRawInstancesCmd)
	listRawCmd.AddCommand(listRawImagesCmd)

	listCmd.AddCommand(listRawCmd)
}

var listRawCmd = &cobra.Command{
	Use:   "raw",
	Short: "List raw JSONs from AWS.",
}

// access

var listRawUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.AccessService.Users()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}

var listRawGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.AccessService.Groups()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}

var listRawRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List roles",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.AccessService.Roles()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}

var listRawPoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List policies",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.AccessService.LocalPolicies()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}

// infra

var listRawRegionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List regions",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Regions()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}

var listRawVpcsCmd = &cobra.Command{
	Use:   "vpcs",
	Short: "List vpcs",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Vpcs()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}

var listRawSubnetsCmd = &cobra.Command{
	Use:   "subnets",
	Short: "List subnets",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Subnets()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}

var listRawInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List instances",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Instances()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}

var listRawImagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List images",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.InfraService.Images()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(resp)
	},
}
