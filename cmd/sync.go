package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/api"
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
		infrag, accessg, err := performSync()
		if err != nil {
			return err
		}

		printWithTabs := func(g *rdf.Graph, n *node.Node, distance int) {
			var tabs bytes.Buffer
			for i := 0; i < distance; i++ {
				tabs.WriteByte('\t')
			}
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
		}

		root, err := node.NewNodeFromStrings("/region", viper.GetString("region"))
		if err != nil {
			return err
		}

		infrag.VisitDepthFirst(root, printWithTabs)
		accessg.VisitDepthFirst(root, printWithTabs)

		return nil
	},
}

func performSync() (*rdf.Graph, *rdf.Graph, error) {
	var awsInfra *api.AwsInfra
	var awsAccess *api.AwsAccess

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		infra, err := infraApi.FetchAwsInfra()
		exitOn(err)
		awsInfra = infra
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		access, err := accessApi.FetchAwsAccess()
		exitOn(err)
		awsAccess = access
	}()

	wg.Wait()

	infrag, err := rdf.BuildAwsInfraGraph(viper.GetString("region"), awsInfra)

	tofile, err := infrag.Marshal()
	if err != nil {
		return nil, nil, err
	}
	if err = ioutil.WriteFile(filepath.Join(config.Dir, config.InfraFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	accessg, err := rdf.BuildAwsAccessGraph(viper.GetString("region"), awsAccess)

	tofile, err = accessg.Marshal()
	if err != nil {
		return nil, nil, err
	}
	if err := ioutil.WriteFile(filepath.Join(config.Dir, config.AccessFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	return infrag, accessg, nil
}
