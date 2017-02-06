package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

func init() {
	RootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:                "show",
	Short:              "Show resource and their relations via a given id: users, groups, instances, vpcs, ...",
	PersistentPreRun:   applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("id required")
		}

		id := args[0]
		notFound := fmt.Sprintf("resource with id %s not found", id)
		var gph *graph.Graph
		var graphs map[string]*graph.Graph

		resType := resolveResourceType(id)

		if resType == "" && localFlag {
			logger.Info(notFound)
			return nil
		} else if resType == "" {
			graphs = runFullSync()

			if resType = resolveResourceType(id); resType == "" {
				logger.Info(notFound)
				return nil
			}
		}

		srv, err := cloud.GetServiceForType(resType)
		exitOn(err)

		if graphs == nil {
			logger.Verbosef("syncing service for %s type", resType)
			graphs, err = sync.DefaultSyncer.Sync(srv)
			exitOn(err)
		}

		gph = graphs[srv.Name()]

		res, err := gph.FindResource(args[0])
		exitOn(err)

		if res != nil {
			displayer := console.BuildOptions(
				console.WithHeaders(console.DefaultsColumnDefinitions[res.Type()]),
				console.WithFormat(listingFormat),
			).SetSource(res).Build()

			exitOn(displayer.Print(os.Stderr))
		}

		return nil
	},
}

func runFullSync() map[string]*graph.Graph {
	logger.Info("cannot resolve resource - running full sync")

	var services []cloud.Service
	for _, srv := range cloud.ServiceRegistry {
		services = append(services, srv)
	}

	graphs, err := sync.DefaultSyncer.Sync(services...)
	exitOn(err)

	return graphs
}

func resolveResourceType(id string) (resType string) {
	for _, name := range aws.ServiceNames {
		g := sync.LoadCurrentLocalGraph(name)
		localRes, err := g.FindResource(id)
		exitOn(err)
		if localRes != nil {
			resType = localRes.Type().String()
		}
	}
	return
}
