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
	printOnlyID    bool
	localResources bool
)

var infraResourcesToDisplay = map[string][]*display.PropertyDisplayer{
	"instance": []*display.PropertyDisplayer{{Property: "Id"}, {Property: "Tags[].Name", Label: "Name"}, {Property: "State.Name", Label: "State", ColoredValues: map[string]string{"running": "green", "stopped": "red"}}, {Property: "Type"}, {Property: "PublicIp", Label: "Public IP"}, {Property: "PrivateIp", Label: "Private IP"}},
	"vpc":      []*display.PropertyDisplayer{{Property: "Id"}, {Property: "IsDefault", Label: "Default", ColoredValues: map[string]string{"true": "green"}}, {Property: "State"}, {Property: "CidrBlock"}},
	"subnet":   []*display.PropertyDisplayer{{Property: "Id"}, {Property: "MapPublicIpOnLaunch", Label: "Public VMs", ColoredValues: map[string]string{"true": "red"}}, {Property: "State", ColoredValues: map[string]string{"available": "green"}}, {Property: "CidrBlock"}},
}

var accessResourcesToDisplay = map[string][]*display.PropertyDisplayer{
	"user":   []*display.PropertyDisplayer{{Property: "Id"}, {Property: "Name"}, {Property: "Arn"}, {Property: "Path"}, {Property: "PasswordLastUsed"}},
	"role":   []*display.PropertyDisplayer{{Property: "Id"}, {Property: "Name"}, {Property: "Arn"}, {Property: "CreateDate"}, {Property: "Path"}},
	"policy": []*display.PropertyDisplayer{{Property: "Id"}, {Property: "Name"}, {Property: "Arn"}, {Property: "Description"}, {Property: "isAttachable"}, {Property: "CreateDate"}, {Property: "UpdateDate"}, {Property: "Path"}},
	"group":  []*display.PropertyDisplayer{{Property: "Id"}, {Property: "Name"}, {Property: "Arn"}, {Property: "CreateDate"}, {Property: "Path"}},
}

func init() {
	RootCmd.AddCommand(rdfListCmd)
	for resource, properties := range infraResourcesToDisplay {
		rdfListCmd.AddCommand(rdfListInfraResourceCmd(resource, properties))
	}
	for resource, properties := range accessResourcesToDisplay {
		rdfListCmd.AddCommand(rdfListAccessResourceCmd(resource, properties))
	}
	rdfListCmd.AddCommand(rdfListAliasesCmd)
	rdfListCmd.AddCommand(rdfListAllCmd)

	rdfListCmd.PersistentFlags().BoolVar(&printOnlyID, "ids", false, "List only ids")
	rdfListCmd.PersistentFlags().BoolVar(&localResources, "local", false, "List locally sync resources")
}

var rdfListCmd = &cobra.Command{
	Use:   "rdflist",
	Short: "List various type of items: instances, vpc, subnet ...",
}

var rdfListAliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "List aliases",

	Run: func(cmd *cobra.Command, args []string) {
		displayAliases(statsDB.GetAliases())
	},
}

var rdfListInfraResourceCmd = func(resource string, properties []*display.PropertyDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := "/" + resource
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS EC2 " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			if localResources {
				localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
				display.ResourceOfGraph(localInfra, nodeType, properties, printOnlyID, err)
			} else {
				listRemoteCloudResource(aws.InfraService, resources, nodeType, properties)
			}
		},
	}
}

var rdfListAccessResourceCmd = func(resource string, properties []*display.PropertyDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := "/" + resource
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS IAM " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			if localResources {
				localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
				display.ResourceOfGraph(localAccess, nodeType, properties, printOnlyID, err)
			} else {
				listRemoteCloudResource(aws.AccessService, resources, nodeType, properties)
			}
		},
	}
}

var rdfListAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all local resources",

	Run: func(cmd *cobra.Command, args []string) {
		if !printOnlyID {
			fmt.Println("Infrastructure")
		}
		localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
		display.SeveralResourcesOfGraph(localInfra, infraResourcesToDisplay, printOnlyID, err)

		if !printOnlyID {
			fmt.Println("Access")
		}
		localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
		display.SeveralResourcesOfGraph(localAccess, accessResourcesToDisplay, printOnlyID, err)
	},
}

func listRemoteCloudResource(cloudService interface{}, resources string, nodeType string, properties []*display.PropertyDisplayer) {
	fnName := fmt.Sprintf("%sGraph", humanize(resources))
	method := reflect.ValueOf(cloudService).MethodByName(fnName)
	if method.IsValid() && !method.IsNil() {
		methodI := method.Interface()
		if graphFn, ok := methodI.(func() (*rdf.Graph, error)); ok {
			graph, err := graphFn()
			display.ResourceOfGraph(graph, nodeType, properties, printOnlyID, err)
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
