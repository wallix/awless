package cmd

import (
	"bytes"
	"context"
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
		infra, err := infraApi.FetchAwsInfra()
		if err != nil {
			return err
		}

		triples, err := rdf.BuildInfraRdfTriples(viper.GetString("region"), infra)
		if err != nil {
			return err
		}

		if err = ioutil.WriteFile(filepath.Join(config.Dir, config.InfraFilename), []byte(rdf.MarshalTriples(triples)), 0600); err != nil {
			return err
		}

		infrag, err := rdf.NewMemGraph("infra")
		if err != nil {
			return err
		}

		infrag.AddTriples(context.Background(), triples)
		root, err := node.NewNodeFromStrings("/region", viper.GetString("region"))
		if err != nil {
			return err
		}

		rdf.VisitDepthFirst(infrag, root, func(n *node.Node, distance int) {
			var tabs bytes.Buffer
			for i := 0; i < distance; i++ {
				tabs.WriteByte('\t')
			}
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
		})

		access, err := accessApi.FetchAwsAccess()
		if err != nil {
			return err
		}

		triples, err = rdf.BuildAccessRdfTriples(viper.GetString("region"), access)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath.Join(config.Dir, config.AccessFilename), []byte(rdf.MarshalTriples(triples)), 0600); err != nil {
			return err
		}

		accessg, err := rdf.NewMemGraph("access")
		if err != nil {
			return err
		}

		accessg.AddTriples(context.Background(), triples)
		root, err = node.NewNodeFromStrings("/region", viper.GetString("region"))
		if err != nil {
			return err
		}

		rdf.VisitDepthFirst(accessg, root, func(n *node.Node, distance int) {
			var tabs bytes.Buffer
			for i := 0; i < distance; i++ {
				tabs.WriteByte('\t')
			}
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
		})

		return nil
	},
}
