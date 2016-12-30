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

var infraResourcesToDisplay = map[string][]*display.PropertyDisplayer{
	"instance": []*display.PropertyDisplayer{
		{Property: "Id"},
		{Property: "Tags[].Name", Label: "Name"},
		{Property: "State.Name", Label: "State", ColoredValues: map[string]string{"running": "green", "stopped": "red"}},
		{Property: "Type"},
		{Property: "PublicIp", Label: "Public IP"},
		{Property: "PrivateIp", Label: "Private IP"},
	},
	"vpc": []*display.PropertyDisplayer{
		{Property: "Id"},
		{Property: "IsDefault", Label: "Default", ColoredValues: map[string]string{"true": "green"}},
		{Property: "State"}, {Property: "CidrBlock"},
	},
	"subnet": []*display.PropertyDisplayer{
		{Property: "Id"},
		{Property: "MapPublicIpOnLaunch", Label: "Public VMs", ColoredValues: map[string]string{"true": "red"}},
		{Property: "State", ColoredValues: map[string]string{"available": "green"}},
		{Property: "CidrBlock"},
	},
}

var accessResourcesToDisplay = map[string][]*display.PropertyDisplayer{
	"user": []*display.PropertyDisplayer{
		{Property: "Id"},
		{Property: "Name"},
		{Property: "Arn"},
		{Property: "Path"},
		{Property: "PasswordLastUsed"},
	},
	"role": []*display.PropertyDisplayer{
		{Property: "Id"},
		{Property: "Name"},
		{Property: "Arn"},
		{Property: "CreateDate"},
		{Property: "Path"},
	},
	"policy": []*display.PropertyDisplayer{
		{Property: "Id"},
		{Property: "Name"},
		{Property: "Arn"},
		{Property: "Description"},
		{Property: "isAttachable"},
		{Property: "CreateDate"},
		{Property: "UpdateDate"},
		{Property: "Path"},
	},
	"group": []*display.PropertyDisplayer{
		{Property: "Id"},
		{Property: "Name"},
		{Property: "Arn"},
		{Property: "CreateDate"},
		{Property: "Path"},
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	for resource, properties := range infraResourcesToDisplay {
		listCmd.AddCommand(listInfraResourceCmd(resource, properties))
	}
	for resource, properties := range accessResourcesToDisplay {
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

var listInfraResourceCmd = func(resource string, properties []*display.PropertyDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := "/" + resource
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS EC2 " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			if localResources {
				localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
				display.ResourceOfGraph(localInfra, nodeType, properties, sortBy, listOnlyIDs, err)
			} else {
				listRemoteCloudResource(aws.InfraService, resources, nodeType, properties)
			}
		},
	}
}

var listAccessResourceCmd = func(resource string, properties []*display.PropertyDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := "/" + resource
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS IAM " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			if localResources {
				localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
				display.ResourceOfGraph(localAccess, nodeType, properties, sortBy, listOnlyIDs, err)
			} else {
				listRemoteCloudResource(aws.AccessService, resources, nodeType, properties)
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
			display.SeveralResourcesOfGraph(localInfra, infraResourcesToDisplay, listOnlyIDs, err)
		}
		if listAllAccess {
			if !listOnlyIDs {
				fmt.Println("Access")
			}
			localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
			display.SeveralResourcesOfGraph(localAccess, accessResourcesToDisplay, listOnlyIDs, err)
		}
	},
}

func listRemoteCloudResource(cloudService interface{}, resources string, nodeType string, properties []*display.PropertyDisplayer) {
	fnName := fmt.Sprintf("%sGraph", humanize(resources))
	method := reflect.ValueOf(cloudService).MethodByName(fnName)
	if method.IsValid() && !method.IsNil() {
		methodI := method.Interface()
		if graphFn, ok := methodI.(func() (*rdf.Graph, error)); ok {
			graph, err := graphFn()
			display.ResourceOfGraph(graph, nodeType, properties, sortBy, listOnlyIDs, err)
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
