package cmd

import (
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/store"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Manage your local infrastructure",

	RunE: func(cmd *cobra.Command, args []string) error {
		vpcs, subnets, instances, err := infraApi.FetchInfra()
		if err != nil {
			return err
		}

		region := store.BuildRegionTree(viper.GetString("region"), vpcs, subnets, instances)

		if err := ioutil.WriteFile(filepath.Join(config.Dir, config.InfraFilename), region.Json(), 0700); err != nil {
			return err
		}

		return nil
	},
}
