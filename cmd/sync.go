package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Manage your local infrastructure",

	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := node.NewNodeFromStrings("/region", viper.GetString("region"))
		if err != nil {
			return err
		}

		printWithTabs := func(n *node.Node, distance int) {
			var tabs bytes.Buffer
			for i := 0; i < distance; i++ {
				tabs.WriteByte('\t')
			}
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
		}

		infra, err := infraApi.FetchAwsInfra()
		if err != nil {
			return err
		}

		infrag, err := rdf.BuildAwsInfraGraph("infra", viper.GetString("region"), infra)

		tofile, err := infrag.Marshal()
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(filepath.Join(config.Dir, config.InfraFilename), tofile, 0600); err != nil {
			return err
		}

		infrag.VisitDepthFirst(root, printWithTabs)

		access, err := accessApi.FetchAwsAccess()
		if err != nil {
			return err
		}

		accessg, err := rdf.BuildAwsAccessGraph("access", viper.GetString("region"), access)

		tofile, err = accessg.Marshal()
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(config.Dir, config.AccessFilename), tofile, 0600); err != nil {
			return err
		}

		accessg.VisitDepthFirst(root, printWithTabs)

		return nil
	},
}
