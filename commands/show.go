package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

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

		var resource *graph.Resource
		var gph *graph.Graph

		resource, gph = findResourceInLocalGraphs(id)

		if resource == nil && localFlag {
			logger.Info(notFound)
			return nil
		} else if resource == nil {
			runFullSync()

			if resource, gph = findResourceInLocalGraphs(id); resource == nil {
				logger.Info(notFound)
				return nil
			}
		}

		if !localFlag {
			srv, err := cloud.GetServiceForType(resource.Type().String())
			exitOn(err)
			logger.Verbosef("syncing service for %s type", resource.Type())
			_, err = sync.DefaultSyncer.Sync(srv)
			exitOn(err)
		}

		if resource != nil {
			displayer := console.BuildOptions(
				console.WithHeaders(console.DefaultsColumnDefinitions[resource.Type()]),
				console.WithFormat(listingFormat),
			).SetSource(resource).Build()

			exitOn(displayer.Print(os.Stderr))

			printWithTabs := func(r *graph.Resource, distance int) {
				var tabs bytes.Buffer
				for i := 0; i < distance; i++ {
					tabs.WriteByte('\t')
				}
				if resource.Id() != r.Id() {
					fmt.Fprintf(os.Stdout, "%s%s[%s]\n", tabs.String(), r.Type(), r.Id())
				}
			}

			fmt.Println("\nParents:")
			gph.VisitParents(resource, printWithTabs)
			fmt.Println("\nChildrens:")
			gph.VisitChildren(resource, printWithTabs)

			var siblings []string
			gph.VisitSiblings(resource, func(r *graph.Resource, distance int) {
				siblings = append(siblings, fmt.Sprintf("%s[%s]", r.Type(), r.Id()))
			})
			fmt.Printf("\nSiblings: %s\n", strings.Join(siblings, ", "))
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

func findResourceInLocalGraphs(id string) (*graph.Resource, *graph.Graph) {
	for _, name := range aws.ServiceNames {
		g := sync.LoadCurrentLocalGraph(name)
		localRes, err := g.FindResource(id)
		exitOn(err)
		if localRes != nil {
			return localRes, g
		}
	}
	return nil, nil
}
