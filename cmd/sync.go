package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/badwolf/triple/node"
	"github.com/libgit2/git2go"
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

		if verboseFlag {
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
		}

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
	if err = ioutil.WriteFile(filepath.Join(config.GitDir, config.InfraFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	accessg, err := rdf.BuildAwsAccessGraph(viper.GetString("region"), awsAccess)

	tofile, err = accessg.Marshal()
	if err != nil {
		return nil, nil, err
	}
	if err := ioutil.WriteFile(filepath.Join(config.GitDir, config.AccessFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	if err := saveSyncRevision(); err != nil {
		return nil, nil, err
	}

	return infrag, accessg, nil
}

func saveSyncRevision() error {
	if _, err := os.Stat(filepath.Join(config.GitDir, ".git")); os.IsNotExist(err) {
		if _, err := git.InitRepository(config.GitDir, false); err != nil {
			return err
		}
	}

	repo, err := git.OpenRepository(config.GitDir)
	if err != nil {
		return err
	}

	idx, err := repo.Index()
	if err != nil {
		return err
	}

	if err := idx.AddByPath(config.InfraFilename); err != nil {
		return err
	}
	if err := idx.AddByPath(config.AccessFilename); err != nil {
		return err
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		return err
	}

	if err := idx.Write(); err != nil {
		return err
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return err
	}

	var parents []*git.Commit

	head, err := repo.Head()
	if err == nil {
		headCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}
		parents = append(parents, headCommit)
	}

	sig := &git.Signature{Name: "awless", Email: "git@awless.io"}

	if _, err = repo.CreateCommit("HEAD", sig, sig, "new sync", tree, parents...); err != nil {
		return err
	}

	return nil
}
