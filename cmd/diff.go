package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/config"
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

		local, err := store.UnmarshalTriples(string(content))
		if err != nil {
			return err
		}

		vpcs, subnets, instances, err := infraApi.FetchInfra()
		if err != nil {
			return err
		}

		remote, err := store.BuildInfraRdfTriples(viper.GetString("region"), vpcs, subnets, instances)
		if err != nil {
			return err
		}
		extras, missings, err := store.Compare(viper.GetString("region"), local, remote)
		if err != nil {
			return err
		}

		fmt.Println("Extras:")
		fmt.Printf("\t%s\n", store.MarshalTriples(extras))
		fmt.Println("Missings:")
		fmt.Printf("\t%s\n", store.MarshalTriples(missings))

		return nil
	},
}
