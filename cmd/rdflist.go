package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

var infraResourcesToDisplay = map[string][]string{
	"instance": []string{"Id", "Tags[].Name", "State.Name", "Type", "PublicIp", "PrivateIp"},
	"vpc":      []string{"Id", "IsDefault", "State", "CidrBlock"},
	"subnet":   []string{"Id", "MapPublicIpOnLaunch", "State", "CidrBlock"},
}

var accessResourcesToDisplay = map[string][]string{
	"user":   []string{"Id", "Name", "Arn", "Path", "PasswordLastUsed"},
	"role":   []string{"Id", "Name", "Arn", "CreateDate", "Path"},
	"policy": []string{"Id", "Name", "Arn", "Description", "isAttachable", "CreateDate", "UpdateDate", "Path"},
	"group":  []string{"Id", "Name", "Arn", "CreateDate", "Path"},
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

var rdfListInfraResourceCmd = func(resource string, properties []string) *cobra.Command {
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

var rdfListAccessResourceCmd = func(resource string, properties []string) *cobra.Command {
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

func ListCloudResource(cloudService interface{}, resources string, nodeType string, properties []string) {
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
