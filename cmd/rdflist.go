package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

var infraResourcesToDisplay = map[string][]PropertyDisplayer{
	"instance": []PropertyDisplayer{{Property: "Id"}, {Property: "Tags[].Name", Label: "Name"}, {Property: "State.Name", Label: "State", ColoredValues: map[string]string{"running": "green", "stopped": "red"}}, {Property: "Type"}, {Property: "PublicIp", Label: "Public IP"}, {Property: "PrivateIp", Label: "Private IP"}},
	"vpc":      []PropertyDisplayer{{Property: "Id"}, {Property: "IsDefault", Label: "Default", ColoredValues: map[string]string{"true": "green"}}, {Property: "State"}, {Property: "CidrBlock"}},
	"subnet":   []PropertyDisplayer{{Property: "Id"}, {Property: "MapPublicIpOnLaunch", Label: "Public VMs", ColoredValues: map[string]string{"true": "red"}}, {Property: "State", ColoredValues: map[string]string{"available": "green"}}, {Property: "CidrBlock"}},
}

var accessResourcesToDisplay = map[string][]PropertyDisplayer{
	"user":   []PropertyDisplayer{{Property: "Id"}, {Property: "Name"}, {Property: "Arn"}, {Property: "Path"}, {Property: "PasswordLastUsed"}},
	"role":   []PropertyDisplayer{{Property: "Id"}, {Property: "Name"}, {Property: "Arn"}, {Property: "CreateDate"}, {Property: "Path"}},
	"policy": []PropertyDisplayer{{Property: "Id"}, {Property: "Name"}, {Property: "Arn"}, {Property: "Description"}, {Property: "isAttachable"}, {Property: "CreateDate"}, {Property: "UpdateDate"}, {Property: "Path"}},
	"group":  []PropertyDisplayer{{Property: "Id"}, {Property: "Name"}, {Property: "Arn"}, {Property: "CreateDate"}, {Property: "Path"}},
}

func init() {
	RootCmd.AddCommand(rdfListCmd)
	for resource, properties := range infraResourcesToDisplay {
		rdfListCmd.AddCommand(rdfListInfraResourceCmd(resource, properties))
	}
	for resource, properties := range accessResourcesToDisplay {
		rdfListCmd.AddCommand(rdfListAccessResourceCmd(resource, properties))
	}
}

var rdfListCmd = &cobra.Command{
	Use:   "rdflist",
	Short: "List various type of items: instances, vpc, subnet ...",
}

var rdfListInfraResourceCmd = func(resource string, properties []PropertyDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := "/" + resource
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS EC2 " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			ListCloudResource(aws.InfraService, resources, nodeType, properties)
		},
	}
}

var rdfListAccessResourceCmd = func(resource string, properties []PropertyDisplayer) *cobra.Command {
	resources := pluralize(resource)
	nodeType := "/" + resource
	return &cobra.Command{
		Use:   resources,
		Short: "List AWS IAM " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			ListCloudResource(aws.AccessService, resources, nodeType, properties)
		},
	}
}

func ListCloudResource(cloudService interface{}, resources string, nodeType string, properties []PropertyDisplayer) {
	fnName := fmt.Sprintf("%sGraph", humanize(resources))
	method := reflect.ValueOf(cloudService).MethodByName(fnName)
	if method.IsValid() && !method.IsNil() {
		methodI := method.Interface()
		if graphFn, ok := methodI.(func() (*rdf.Graph, error)); ok {
			graph, err := graphFn()
			displayGraph(graph, nodeType, properties, err)
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
