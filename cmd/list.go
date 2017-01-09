package cmd

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/rdf"
)

var (
	listingFormat string

	listOnlyIDs    bool
	listAllInfra   bool
	listAllAccess  bool
	localResources bool
	sortBy         []string
)

func init() {
	RootCmd.AddCommand(listCmd)
	for resource, properties := range display.PropertiesDisplayer.Services[aws.InfraServiceName].Resources {
		listCmd.AddCommand(listInfraResourceCmd(resource, properties))
	}
	for resource, properties := range display.PropertiesDisplayer.Services[aws.AccessServiceName].Resources {
		listCmd.AddCommand(listAccessResourceCmd(resource, properties))
	}
	listCmd.AddCommand(listAllCmd)

	listCmd.AddCommand(listInstancesCmd)
	listCmd.Flags().StringVar(&listingFormat, "format", "csv", "Format for the display of resources: csv, table, ...")

	listCmd.PersistentFlags().BoolVar(&listOnlyIDs, "ids", false, "List only ids")
	listCmd.PersistentFlags().BoolVar(&localResources, "local", false, "List locally sync resources")
	listCmd.PersistentFlags().StringSliceVar(&sortBy, "sort-by", []string{"Id"}, "Sort tables by column(s) name(s)")

	listAllCmd.PersistentFlags().BoolVar(&listAllInfra, "infra", false, "List infrastructure resources")
	listAllCmd.PersistentFlags().BoolVar(&listAllAccess, "access", false, "List access resources")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List various type of items: instances, vpc, subnet ...",
}

var listInstancesCmd = &cobra.Command{
	Use:   "csvinstances",
	Short: "List aws instances",

	Run: func(cmd *cobra.Command, args []string) {
		g, err := aws.InfraService.InstancesGraph()

		exitOn(err)

		displayer := display.BuildDisplayer(display.Options{
			RdfType: rdf.INSTANCE, Format: listingFormat,
		})
		displayer.SetGraph(g)
		displayer.SetHeaders([]display.Header{
			display.StringHeader{Prop: "Id"},
			display.StringHeader{Prop: "Name"},
			display.StringHeader{Prop: "State"},
			display.StringHeader{Prop: "Type"},
			display.StringHeader{Prop: "KeyName", Friendly: "Access Key"},
			display.StringHeader{Prop: "PublicIp", Friendly: "Public IP"},
		})

		fmt.Println(displayer.Print())
	},
}

var listInfraResourceCmd = func(resource string, displayer *display.ResourceDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := rdf.ToRDFType(resource)
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS EC2 " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			var g *rdf.Graph
			var err error
			if localResources {
				g, err = rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))

			} else {
				g, err = remoteResourceGraph(aws.InfraService, resources)
			}
			exitOn(err)
			display.ResourcesOfGraph(g, nodeType, displayer, sortBy, listOnlyIDs)
		},
	}
}

var listAccessResourceCmd = func(resource string, displayer *display.ResourceDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := rdf.ToRDFType(resource)
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS IAM " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			var g *rdf.Graph
			var err error
			if localResources {
				g, err = rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
			} else {
				g, err = remoteResourceGraph(aws.AccessService, resources)
			}
			exitOn(err)
			display.ResourcesOfGraph(g, nodeType, displayer, sortBy, listOnlyIDs)
		},
	}
}

var listAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all local resources",

	Run: func(cmd *cobra.Command, args []string) {
		if !listAllInfra && !listAllAccess {
			listAllInfra = true //By default, print only infra
		}
		if listAllInfra {
			if !listOnlyIDs {
				fmt.Println("Infrastructure")
			}
			localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
			exitOn(err)
			display.SeveralResourcesOfGraph(localInfra, display.PropertiesDisplayer.Services[aws.InfraServiceName], listOnlyIDs)
		}
		if listAllAccess {
			if !listOnlyIDs {
				fmt.Println("Access")
			}
			localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
			exitOn(err)
			display.SeveralResourcesOfGraph(localAccess, display.PropertiesDisplayer.Services[aws.AccessServiceName], listOnlyIDs)
		}
	},
}

func remoteResourceGraph(cloudService interface{}, resources string) (*rdf.Graph, error) {
	fnName := fmt.Sprintf("%sGraph", humanize(resources))
	method := reflect.ValueOf(cloudService).MethodByName(fnName)
	if method.IsValid() && !method.IsNil() {
		methodI := method.Interface()
		if graphFn, ok := methodI.(func() (*rdf.Graph, error)); ok {
			return graphFn()
		}
	}
	return nil, (fmt.Errorf("Unknown type of resource: %s", resources))
}

func pluralize(singular string) string {
	if strings.HasSuffix(singular, "y") {
		return strings.TrimSuffix(singular, "y") + "ies"
	}
	return singular + "s"
}

func humanize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}
	return strings.ToUpper(s)
}
