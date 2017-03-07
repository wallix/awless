/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

func init() {
	RootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show REFERENCE",
	Short: "Show a resource and its interrelations given a REFERENCE: id or name",
	Example: `  awless show i-8d43b21b            # show an instance via its id
  awless show AIDAJ3Z24GOKHTZO4OIX6 # show a user via its id
  awless show jsmith                # show a user via its name`,
	PersistentPreRun:   applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, verifyNewVersionHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("REFERENCE required. See examples.")
		}

		id := strings.TrimPrefix(args[0], "@") // in case users still use @ or its captured in copy/paste
		notFound := fmt.Sprintf("resource with id %s not found", id)

		var resource *graph.Resource
		var gph *graph.Graph

		resource, gph = findResourceInLocalGraphs(id)

		if resource == nil && localGlobalFlag {
			logger.Info(notFound)
			return nil
		} else if resource == nil {
			runFullSync()

			if resource, gph = findResourceInLocalGraphs(id); resource == nil {
				logger.Info(notFound)
				return nil
			}
		}

		if !localGlobalFlag {
			srv, err := cloud.GetServiceForType(resource.Type().String())
			exitOn(err)
			logger.Verbosef("syncing service for %s type", resource.Type())
			if _, err = sync.DefaultSyncer.Sync(srv); err != nil {
				logger.Error(err)
			}
		}

		if resource != nil {
			showResource(resource, gph)
		}

		return nil
	},
}

func showResource(resource *graph.Resource, gph *graph.Graph) {
	displayer := console.BuildOptions(
		console.WithHeaders(console.DefaultsColumnDefinitions[resource.Type()]),
		console.WithFormat(listingFormat),
	).SetSource(resource).Build()

	exitOn(displayer.Print(os.Stderr))

	var parents []*graph.Resource
	err := gph.Accept(&graph.ParentsVisitor{From: resource, Each: graph.VisitorCollectFunc(&parents)})
	exitOn(err)

	fmt.Println("\nRelations:")

	var count int
	for i := len(parents) - 1; i >= 0; i-- {
		if count == 0 {
			fmt.Printf("%s\n", parents[i])
		} else {
			fmt.Printf("%s↳ %s\n", strings.Repeat("\t", count), parents[i])
		}
		count++
	}

	printWithTabs := func(r *graph.Resource, distance int) error {
		var tabs bytes.Buffer
		tabs.WriteString(strings.Repeat("\t", count))
		for i := 0; i < distance; i++ {
			tabs.WriteByte('\t')
		}

		display := r.String()
		if r.Same(resource) {
			display = renderGreenFn(resource.String())
		}
		fmt.Printf("%s↳ %s\n", tabs.String(), display)

		return nil
	}

	err = gph.Accept(&graph.ChildrenVisitor{From: resource, Each: printWithTabs, IncludeFrom: true})
	exitOn(err)

	var siblings []*graph.Resource
	err = gph.Accept(&graph.SiblingsVisitor{From: resource, Each: graph.VisitorCollectFunc(&siblings)})
	exitOn(err)
	printResourceList("Siblings", siblings)

	appliedOn, err := gph.ListResourcesAppliedOn(resource)
	exitOn(err)
	printResourceList("Applied on", appliedOn)

	dependingOn, err := gph.ListResourcesDependingOn(resource)
	exitOn(err)
	printResourceList("Depending on", dependingOn)
}

func runFullSync() map[string]*graph.Graph {
	logger.Info("cannot resolve resource - running full sync")

	var services []cloud.Service
	for _, srv := range cloud.ServiceRegistry {
		services = append(services, srv)
	}

	graphs, err := sync.DefaultSyncer.Sync(services...)
	logger.Verbose(err)

	return graphs
}

func findResourceInLocalGraphs(ref string) (*graph.Resource, *graph.Graph) {
	resources := resolveResourceFromRef(ref)
	switch len(resources) {
	case 0:
		return nil, nil
	case 1:
		res := resources[0]
		return res, sync.LoadCurrentLocalGraph(aws.ServicePerResourceType[res.Type().String()])
	default:
		var all []string
		for _, res := range resources {
			all = append(all, fmt.Sprintf("%s[%s]", res.Id(), res.Type()))
		}
		logger.Infof("%d resources found with the name '%s': %s", len(resources), ref, strings.Join(all, ", "))
		logger.Info("Show them using their id")
		os.Exit(0)
	}

	return nil, nil
}

func resolveResourceFromRef(ref string) []*graph.Resource {
	g := graph.NewGraph()
	for _, s := range aws.ServiceNames {
		g.AddGraph(sync.LoadCurrentLocalGraph(s))
	}

	byId, err := g.FindResource(ref)
	exitOn(err)

	var resources []*graph.Resource
	if byId != nil {
		resources = append(resources, byId)
	} else {
		rs, err := g.ResolveResources(
			&graph.ByProperty{"Name", ref},
			&graph.ByProperty{"Arn", ref},
		)
		exitOn(err)
		resources = append(resources, rs...)
	}

	return resources
}

func printResourceList(title string, list []*graph.Resource) {
	all := graph.Resources(list).Map(func(r *graph.Resource) string { return r.String() })
	if len(all) > 0 {
		fmt.Printf("\n%s: %s\n", title, strings.Join(all, ", "))
	}
}
