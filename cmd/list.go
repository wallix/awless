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

var listInfraResourceCmd = func(resource string, displayer *display.ResourceDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := rdf.ToRDFType(resource)
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS EC2 " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			if localResources {
				localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
				display.ResourceOfGraph(localInfra, nodeType, displayer, sortBy, listOnlyIDs, err)
			} else {
				listRemoteCloudResource(aws.InfraService, resources, nodeType, displayer)
			}
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
			if localResources {
				localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
				display.ResourceOfGraph(localAccess, nodeType, displayer, sortBy, listOnlyIDs, err)
			} else {
				listRemoteCloudResource(aws.AccessService, resources, nodeType, displayer)
			}
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
			display.SeveralResourcesOfGraph(localInfra, display.PropertiesDisplayer.Services[aws.InfraServiceName], listOnlyIDs, err)
		}
		if listAllAccess {
			if !listOnlyIDs {
				fmt.Println("Access")
			}
			localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
			display.SeveralResourcesOfGraph(localAccess, display.PropertiesDisplayer.Services[aws.AccessServiceName], listOnlyIDs, err)
		}
	},
}

func listRemoteCloudResource(cloudService interface{}, resources string, nodeType string, displayer *display.ResourceDisplayer) {
	fnName := fmt.Sprintf("%sGraph", humanize(resources))
	method := reflect.ValueOf(cloudService).MethodByName(fnName)
	if method.IsValid() && !method.IsNil() {
		methodI := method.Interface()
		if graphFn, ok := methodI.(func() (*rdf.Graph, error)); ok {
			graph, err := graphFn()
			display.ResourceOfGraph(graph, nodeType, displayer, sortBy, listOnlyIDs, err)
			return
		}
	}
	fmt.Println(fmt.Errorf("Unknown type of resource: %s", resources))
	return
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
