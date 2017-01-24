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
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/revision/repo"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:               "sync",
	Short:             "Manage your local infrastructure",
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,

	RunE: func(cmd *cobra.Command, args []string) error {
		region := database.MustGetDefaultRegion()
		infrag, accessg, err := performSync(region)
		if err != nil {
			return err
		}

		if verboseFlag {
			printWithTabs := func(g *graph.Graph, n *node.Node, distance int) {
				var tabs bytes.Buffer
				for i := 0; i < distance; i++ {
					tabs.WriteByte('\t')
				}
				fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
			}

			root, err := node.NewNodeFromStrings("/region", region)
			if err != nil {
				return err
			}

			infrag.Visit(root, printWithTabs)
			accessg.Visit(root, printWithTabs)
		}

		return nil
	},
}

func performSync(region string) (*graph.Graph, *graph.Graph, error) {
	var awsInfra *aws.AwsInfra
	var awsAccess *aws.AwsAccess

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		infra, err := aws.InfraService.FetchAwsInfra()
		exitOn(err)
		awsInfra = infra
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		access, err := aws.AccessService.FetchAwsAccess()
		exitOn(err)
		awsAccess = access
	}()

	wg.Wait()

	infrag, err := aws.BuildAwsInfraGraph(region, awsInfra)

	tofile, err := infrag.Marshal()
	if err != nil {
		return nil, nil, err
	}
	if err = ioutil.WriteFile(filepath.Join(config.RepoDir, config.InfraFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	accessg, err := aws.BuildAwsAccessGraph(region, awsAccess)

	tofile, err = accessg.Marshal()
	if err != nil {
		return nil, nil, err
	}
	if err := ioutil.WriteFile(filepath.Join(config.RepoDir, config.AccessFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	r, err := repo.NewRepo()
	if err != nil {
		return nil, nil, err
	}
	if err := r.Commit(config.InfraFilename, config.AccessFilename); err != nil {
		return nil, nil, err
	}

	return infrag, accessg, nil
}
