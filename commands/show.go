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
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

var (
	listAllSiblingsFlag bool
)

func init() {
	RootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVar(&listAllSiblingsFlag, "siblings", false, "List all the resource's siblings")
}

var showCmd = &cobra.Command{
	Use:   "show REFERENCE",
	Short: "Show a resource and its interrelations given a REFERENCE: id or name",
	Example: `  awless show i-8d43b21b            # show an instance via its id
  awless show AIDAJ3Z24GOKHTZO4OIX6 # show a user via its id
  awless show jsmith                # show a user via its name`,
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
	PersistentPostRun: applyHooks(saveHistoryHook, verifyNewVersionHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("REFERENCE required. See examples.")
		}

		ref := args[0]
		notFound := fmt.Sprintf("resource with reference %s not found", deprefix(ref))

		var resource *graph.Resource
		var gph *graph.Graph

		resource, gph = findResourceInLocalGraphs(ref)

		if resource == nil && localGlobalFlag {
			logger.Info(notFound)
			return nil
		} else if resource == nil {
			runFullSync()

			if resource, gph = findResourceInLocalGraphs(ref); resource == nil {
				logger.Info(notFound)
				return nil
			}
		}

		if !localGlobalFlag && config.GetAutosync() {
			srv, err := cloud.GetServiceForType(resource.Type())
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

	fmt.Println(renderCyanBoldFn("\nRelations:"))

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
	printResourceList(renderCyanBoldFn("Siblings"), siblings, "display all with flag --siblings")

	appliedOn, err := gph.ListResourcesAppliedOn(resource)
	exitOn(err)
	printResourceList(renderCyanBoldFn("Applied on"), appliedOn)

	dependingOn, err := gph.ListResourcesDependingOn(resource)
	exitOn(err)
	printResourceList(renderCyanBoldFn("Depending on"), dependingOn)
}

func runFullSync() {
	if !config.GetAutosync() {
		logger.Info("autosync disabled")
		return
	}

	logger.Info("cannot resolve resource - running full sync")

	var services []cloud.Service
	for _, srv := range cloud.ServiceRegistry {
		services = append(services, srv)
	}

	if _, err := sync.DefaultSyncer.Sync(services...); err != nil {
		logger.Verbose(err)
	}
}

func findResourceInLocalGraphs(ref string) (*graph.Resource, *graph.Graph) {
	resources := resolveResourceFromRef(ref)
	switch len(resources) {
	case 0:
		return nil, nil
	case 1:
		res := resources[0]
		return res, sync.LoadCurrentLocalGraph(aws.ServicePerResourceType[res.Type()])
	default:
		var all []string
		for _, res := range resources {
			all = append(all, fmt.Sprintf("%s[%s]", res.Id(), res.Type()))
		}
		logger.Infof("%d resources found with name '%s': %s", len(resources), deprefix(ref), strings.Join(all, ", "))
		logger.Info("Show them using the id:")
		for _, res := range resources {
			logger.Infof("\t`awless show %s` for the %s", res.Id(), res.Type())
		}

		os.Exit(0)
	}

	return nil, nil
}

func resolveResourceFromRef(ref string) []*graph.Resource {
	g, err := sync.LoadAllGraphs()
	exitOn(err)

	name := deprefix(ref)
	byName := &graph.ByProperty{"Name", name}

	if strings.HasPrefix(ref, "@") {
		logger.Verbosef("prefixed with @: forcing research by name '%s'", name)
		rs, err := g.ResolveResources(byName)
		exitOn(err)
		return rs
	} else {

		rs, err := g.ResolveResources(&graph.ById{name})
		exitOn(err)

		if len(rs) > 0 {
			return rs
		} else {
			rs, err := g.ResolveResources(
				byName,
				&graph.ByProperty{"Arn", name},
			)
			exitOn(err)

			return rs
		}
	}
}

func deprefix(s string) string {
	return strings.TrimPrefix(s, "@")
}

func printResourceList(title string, list []*graph.Resource, shortenListMsg ...string) {
	all := graph.Resources(list).Map(func(r *graph.Resource) string { return r.String() })
	count := len(all)
	max := 3
	if count > 0 {
		if !listAllSiblingsFlag && len(shortenListMsg) > 0 && count > max {
			fmt.Printf("\n%s: %s, ... (%s)\n", title, strings.Join(all[0:max], ", "), shortenListMsg[0])
		} else {
			fmt.Printf("\n%s: %s\n", title, strings.Join(all, ", "))
		}
	}
}
