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
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

var (
	listAllSiblingsFlag          bool
	showPropertiesValuesOnlyFlag []string
)

func init() {
	RootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVar(&listAllSiblingsFlag, "siblings", false, "List all the resource's siblings")
	showCmd.Flags().StringSliceVar(&showPropertiesValuesOnlyFlag, "values-for", []string{}, "Output values only for given properties keys")
}

var showCmd = &cobra.Command{
	Use:   "show REFERENCE",
	Short: "Show a resource and its interrelations given a REFERENCE: id or name",
	Example: `  awless show i-8d43b21b            # show an instance via its ref
  awless show AIDAJ3Z24GOKHTZO4OIX6 # show a user via its ref
  awless show jsmith                # show a user via its ref,
  awless show @jsmith               # forcing search by name`,
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade),

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
				logger.Verbose(err)
			}
			resource, gph = findResourceInLocalGraphs(ref)
		}

		if resource != nil {
			if len(showPropertiesValuesOnlyFlag) > 0 {
				showResourceValuesOnlyFor(resource, showPropertiesValuesOnlyFlag)
				return nil
			}
			showResource(resource, gph)
		}

		return nil
	},
}

func showResourceValuesOnlyFor(resource *graph.Resource, propKeys []string) {
	var normalized []string
	for _, p := range propKeys {
		normalized = append(normalized, strings.ToLower(strings.Replace(p, " ", "", -1)))
	}

	valuesForKeys := map[string]string{}
	isIncluded := func(s string) (bool, string) {
		for _, n := range normalized {
			if n == strings.ToLower(s) {
				return true, n
			}
		}
		return false, ""
	}
	for k, v := range resource.Properties {
		if ok, p := isIncluded(k); ok {
			valuesForKeys[p] = fmt.Sprint(v)
		}
	}

	var values []string
	for _, n := range normalized {
		if v, ok := valuesForKeys[n]; ok {
			values = append(values, v)
		}
	}

	fmt.Println(strings.Join(values, ","))
}

func showResource(resource *graph.Resource, gph *graph.Graph) {
	displayer, err := console.BuildOptions(
		console.WithHeaders(console.DefaultsColumnDefinitions[resource.Type()]),
		console.WithFormat(listingFormat),
		console.WithMaxWidth(console.GetTerminalWidth()),
	).SetSource(resource).Build()
	exitOn(err)

	exitOn(displayer.Print(os.Stdout))

	var parents []*graph.Resource
	err = gph.Accept(&graph.ParentsVisitor{From: resource, Each: graph.VisitorCollectFunc(&parents)})
	exitOn(err)

	var parentsW bytes.Buffer
	var count int
	for i := len(parents) - 1; i >= 0; i-- {
		if count == 0 {
			fmt.Fprintf(&parentsW, "%s\n", parents[i])
		} else {
			fmt.Fprintf(&parentsW, "%s↳ %s\n", strings.Repeat("\t", count), parents[i])
		}
		count++
	}

	var childrenW bytes.Buffer
	var hasChildren bool
	printWithTabs := func(r *graph.Resource, distance int) error {
		var tabs bytes.Buffer
		tabs.WriteString(strings.Repeat("\t", count))
		for i := 0; i < distance; i++ {
			tabs.WriteByte('\t')
		}

		display := r.String()
		if r.Same(resource) {
			display = renderGreenFn(resource.String())
		} else {
			hasChildren = true
		}
		fmt.Fprintf(&childrenW, "%s↳ %s\n", tabs.String(), display)
		return nil
	}
	err = gph.Accept(&graph.ChildrenVisitor{From: resource, Each: printWithTabs, IncludeFrom: true})
	exitOn(err)

	if len(parents) > 0 || hasChildren {
		fmt.Println(renderCyanBoldFn("\n# Relations:"))
		fmt.Printf(parentsW.String())
		fmt.Printf(childrenW.String())
	}

	appliedOn, err := gph.ListResourcesAppliedOn(resource)
	exitOn(err)
	printResourceList(renderCyanBoldFn("Applied on"), appliedOn)

	dependingOn, err := gph.ListResourcesDependingOn(resource)
	exitOn(err)
	printResourceList(renderCyanBoldFn("Depending on"), dependingOn)

	var siblings []*graph.Resource
	err = gph.Accept(&graph.SiblingsVisitor{From: resource, Each: graph.VisitorCollectFunc(&siblings)})
	exitOn(err)
	printResourceList(renderCyanBoldFn("Siblings"), siblings, "display all with flag --siblings")
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
	g, resources := resolveResourceFromRef(ref)
	switch len(resources) {
	case 0:
		return nil, nil
	case 1:
		return resources[0], g
	default:
		logger.Infof("%d resources found with name '%s'. Show a specific resource with:", len(resources), deprefix(ref))
		for _, res := range resources {
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("\t`awless show %s` to show the %s", res.Id(), res.Type()))
			if state, ok := res.Properties["State"].(string); ok {
				buf.WriteString(fmt.Sprintf(" (state: '%s')", state))
			}
			logger.Infof(buf.String())
		}

		os.Exit(0)
	}

	return nil, nil
}

func resolveResourceFromRef(ref string) (*graph.Graph, []*graph.Resource) {
	g, err := sync.LoadAllGraphs()
	exitOn(err)

	name := deprefix(ref)
	byName := &graph.ByProperty{Key: "Name", Value: name}

	if strings.HasPrefix(ref, "@") {
		logger.Verbosef("prefixed with @: forcing research by name '%s'", name)
		rs, err := g.ResolveResources(byName)
		exitOn(err)
		return g, rs
	} else {
		rs, err := g.ResolveResources(&graph.ById{Id: name})
		exitOn(err)

		if len(rs) > 0 {
			return g, rs
		} else {
			rs, err := g.ResolveResources(
				byName,
				&graph.ByProperty{Key: "Arn", Value: name},
			)
			exitOn(err)

			return g, rs
		}
	}
}

func deprefix(s string) string {
	return strings.TrimPrefix(s, "@")
}

func printResourceList(title string, list []*graph.Resource, shortenListMsg ...string) {
	sort.Sort(byTypeAndString{list})
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

type byTypeAndString struct {
	res []*graph.Resource
}

func (b byTypeAndString) Len() int { return len(b.res) }
func (b byTypeAndString) Swap(i, j int) {
	b.res[i], b.res[j] = b.res[j], b.res[i]
}
func (b byTypeAndString) Less(i, j int) bool {
	if b.res[i].Type() != b.res[j].Type() {
		return b.res[i].Type() < b.res[j].Type()
	}
	return b.res[i].String() <= b.res[j].String()
}
