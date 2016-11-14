package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/store"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Manage your local infrastructure",

	RunE: func(cmd *cobra.Command, args []string) error {
		var instances []*ec2.Instance
		var vpcs []*ec2.Vpc
		var subnets []*ec2.Subnet

		type fetchFn func() (interface{}, error)

		allFetch := []fetchFn{infraApi.Instances, infraApi.Subnets, infraApi.Vpcs}
		resultc := make(chan interface{})
		errc := make(chan error)

		for _, fetch := range allFetch {
			go func(fn fetchFn) {
				if r, err := fn(); err != nil {
					errc <- err
				} else {
					resultc <- r
				}
			}(fetch)
		}

		for range allFetch {
			select {
			case r := <-resultc:
				switch r.(type) {
				case *ec2.DescribeVpcsOutput:
					vpcs = r.(*ec2.DescribeVpcsOutput).Vpcs
				case *ec2.DescribeSubnetsOutput:
					subnets = r.(*ec2.DescribeSubnetsOutput).Subnets
				case *ec2.DescribeInstancesOutput:
					instances = r.(*ec2.DescribeInstancesOutput).Reservations[0].Instances
				}
			case err := <-errc:
				return err
			}
		}

		tree := store.BuildRegionTree(viper.GetString("region"), vpcs, subnets, instances)

		fmt.Println(tree.Json())
		return nil
	},
}
