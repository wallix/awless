package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/graph"
)

func init() {
	//Resources
	for _, resource := range []graph.ResourceType{graph.Instance, graph.Vpc, graph.Subnet, graph.SecurityGroup, graph.Keypair, graph.InternetGateway, graph.RouteTable} {
		showCmd.AddCommand(showInfraResourceCmd(resource))
	}
	for _, resource := range []graph.ResourceType{graph.User, graph.Role, graph.Policy, graph.Group} {
		showCmd.AddCommand(showAccessResourceCmd(resource))
	}

	RootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:                "show",
	Short:              "Show resource and their relations via a given id: users, groups, instances, vpcs, ...",
	PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,
}

var showInfraResourceCmd = func(resourceType graph.ResourceType) *cobra.Command {
	command := &cobra.Command{
		Use:   resourceType.String() + " id",
		Short: "Show the properties of a AWS EC2 " + resourceType.String(),

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("id required")
			}
			id := args[0]
			var g *graph.Graph
			var err error
			if localResources {
				g, err = config.LoadInfraGraph()

			} else {
				g, err = aws.InfraService.FetchByType(resourceType.String())
			}
			exitOn(err)
			printResource(g, resourceType, id)
			return nil
		},
	}

	command.PersistentFlags().BoolVar(&localResources, "local", false, "List locally sync resources")
	return command
}

var showAccessResourceCmd = func(resourceType graph.ResourceType) *cobra.Command {
	command := &cobra.Command{
		Use:   resourceType.String() + " id",
		Short: "Show the properties of a AWS IAM " + resourceType.String(),

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("id required")
			}
			id := args[0]
			var g *graph.Graph
			var err error
			if localResources {
				g, err = config.LoadAccessGraph()

			} else {
				g, err = aws.AccessService.FetchByType(resourceType.String())
			}
			exitOn(err)
			printResource(g, resourceType, id)
			return nil
		},
	}

	command.PersistentFlags().BoolVar(&localResources, "local", false, "List locally sync resources")
	return command
}

func printResource(g *graph.Graph, resourceType graph.ResourceType, id string) {
	a := graph.Alias(id)
	if aID, ok := a.ResolveToId(g, resourceType); ok {
		id = aID
	}

	resource, err := g.GetResource(resourceType, id)
	if err != nil {
		exitOn(err)
	}

	displayer := display.BuildOptions(
		display.WithHeaders(display.DefaultsColumnDefinitions[resourceType]),
		display.WithFormat(listingFormat),
	).SetSource(resource).Build()

	err = displayer.Print(os.Stderr)
	exitOn(err)
}
