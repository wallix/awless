package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/models"
	"github.com/wallix/awless/store"
)

func init() {
	RootCmd.AddCommand(diffCmd)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show diff between your local and remote infra",

	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := ioutil.ReadFile(filepath.Join(config.Dir, config.InfraFilename))
		if err != nil {
			return err
		}

		localRegion := &models.Region{}
		if err := json.Unmarshal(content, localRegion); err != nil {
			return err
		}

		vpcs, subnets, instances, err := infraApi.FetchInfra()
		if err != nil {
			return err
		}

		remoteRegion := store.BuildRegionTree(viper.GetString("region"), vpcs, subnets, instances)

		fmt.Printf(string(models.Compare(localRegion, remoteRegion).Json()))

		return nil
	},
}
