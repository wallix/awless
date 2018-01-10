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
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

var (
	listingFormat              string
	listingFiltersFlag         []string
	listingTagFiltersFlag      []string
	listingTagKeyFiltersFlag   []string
	listingTagValueFiltersFlag []string
	listingColumnsFlag         []string
	listOnlyIDs                bool
	noHeadersFlag              bool
	sortBy                     []string
	reverseFlag                bool
)

func init() {
	RootCmd.AddCommand(listCmd)

	cobra.EnableCommandSorting = false

	for _, srvName := range awsservices.ServiceNames {
		listCmd.AddCommand(listAllResourceInServiceCmd(srvName))
	}

	for _, name := range awsservices.ServiceNames {
		var resources []string
		for _, resType := range awsservices.ResourceTypesPerServiceName()[name] {
			resources = append(resources, resType)
		}
		sort.Strings(resources)
		for _, resType := range resources {
			listCmd.AddCommand(listSpecificResourceCmd(resType))
		}
	}

	listCmd.PersistentFlags().StringVar(&listingFormat, "format", "table", "Output format: table, csv, tsv, json (default to table)")
	listCmd.PersistentFlags().StringSliceVar(&listingFiltersFlag, "filter", []string{}, "Filter resources given key/values fields (case insensitive). Ex: --filter type=t2.micro")
	listCmd.PersistentFlags().StringSliceVar(&listingTagFiltersFlag, "tag", []string{}, "Filter EC2 resources given tags (case sensitive!). Ex: --tag Env=Production")
	listCmd.PersistentFlags().StringSliceVar(&listingTagKeyFiltersFlag, "tag-key", []string{}, "Filter EC2 resources given a tag key only (case sensitive!). Ex: --tag-key Env")
	listCmd.PersistentFlags().StringSliceVar(&listingTagValueFiltersFlag, "tag-value", []string{}, "Filter EC2 resources given a tag value only (case sensitive!). Ex: --tag-value Staging")
	listCmd.PersistentFlags().StringSliceVar(&listingColumnsFlag, "columns", []string{}, "Select the properties to display in the columns. Ex: --columns id,name,cidr")
	listCmd.PersistentFlags().BoolVar(&listOnlyIDs, "ids", false, "List only ids")
	listCmd.PersistentFlags().BoolVar(&noHeadersFlag, "no-headers", false, "Do not display headers")
	listCmd.PersistentFlags().BoolVar(&reverseFlag, "reverse", false, "Use in conjunction with --sort to reverse sort")
	listCmd.PersistentFlags().StringSliceVar(&sortBy, "sort", []string{"Id"}, "Sort tables by column(s) name(s)")
}

var listCmd = &cobra.Command{
	Use:               "list",
	Aliases:           []string{"ls"},
	Example:           "  awless list instances --sort uptime\n  awless list users --format csv\n  awless list volumes --filter state=use --filter type=gp2\n  awless list volumes --tag-value Purchased\n  awless list vpcs --tag-key Dept --tag-key Internal\n  awless list instances --tag Env=Production,Dept=Marketing\n  awless list instances --filter state=running,type=micro\n  awless list s3objects --filter bucket=pdf-bucket ",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),
	Short:             "List resources: sorting, filtering via tag/properties, output formatting, etc...",
}

var listSpecificResourceCmd = func(resType string) *cobra.Command {
	return &cobra.Command{
		Use:   cloud.PluralizeResource(resType),
		Short: fmt.Sprintf("[%s] List %s %s", awsservices.ServicePerResourceType[resType], strings.ToUpper(awsservices.APIPerResourceType[resType]), cloud.PluralizeResource(resType)),

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				var plural string
				if len(args) > 1 {
					plural = "s"
				}
				logger.Errorf("invalid parameter%s '%s'", plural, strings.Join(args, " "))
				if strings.Contains(args[0], "=") {
					if !promptConfirmDefaultYes("Did you mean `awless list %s --filter %s`? ", cloud.PluralizeResource(resType), strings.Join(args, " ")) {
						os.Exit(1)
					}
					listingFiltersFlag = append(listingFiltersFlag, args...)
				} else {
					os.Exit(1)
				}
			}
			var g cloud.GraphAPI

			if localGlobalFlag {
				if srvName, ok := awsservices.ServicePerResourceType[resType]; ok {
					g = sync.LoadLocalGraphForService(srvName, config.GetAWSProfile(), config.GetAWSRegion())
				} else {
					exitOn(fmt.Errorf("cannot find service for resource type %s", resType))
				}
			} else {
				srv, err := cloud.GetServiceForType(resType)
				exitOn(err)
				g, err = srv.FetchByType(context.WithValue(context.Background(), "force", true), resType)
				exitOn(err)
			}

			printResources(g, resType)
		},
	}
}

var listAllResourceInServiceCmd = func(srvName string) *cobra.Command {
	return &cobra.Command{
		Use:    srvName,
		Short:  fmt.Sprintf("List all %s resources", srvName),
		Hidden: true,

		Run: func(cmd *cobra.Command, args []string) {
			g := sync.LoadLocalGraphForService(srvName, config.GetAWSProfile(), config.GetAWSRegion())
			displayer, err := console.BuildOptions(
				console.WithFormat(listingFormat),
				console.WithMaxWidth(console.GetTerminalWidth()),
				console.WithIDsOnly(listOnlyIDs),
			).SetSource(g).Build()
			exitOn(err)
			exitOn(displayer.Print(os.Stdout))
		},
	}
}

func printResources(g cloud.GraphAPI, resType string) {
	displayer, err := console.BuildOptions(
		console.WithRdfType(resType),
		console.WithColumns(listingColumnsFlag),
		console.WithFilters(listingFiltersFlag),
		console.WithTagFilters(listingTagFiltersFlag),
		console.WithTagKeyFilters(listingTagKeyFiltersFlag),
		console.WithTagValueFilters(listingTagValueFiltersFlag),
		console.WithMaxWidth(console.GetTerminalWidth()),
		console.WithFormat(listingFormat),
		console.WithIDsOnly(listOnlyIDs),
		console.WithSortBy(sortBy...),
		console.WithReverseSort(reverseFlag),
		console.WithNoHeaders(noHeadersFlag),
	).SetSource(g).Build()
	exitOn(err)

	exitOn(displayer.Print(os.Stdout))
}
