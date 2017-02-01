package aws

import "github.com/wallix/awless/graph"

type serviceDefinition struct {
	Name     string
	Api      string
	Fetchers []fetcherDefinition
}

type fetcherDefinition struct {
	ResourceType                                graph.ResourceType
	AWSType                                     string
	ApiMethod, Input                            string
	Output, OutputsContainers, OutputsExtractor string
}

var ServicesDefinitions = []serviceDefinition{
	{
		Name: "infra",
		Api:  "ec2",
		Fetchers: []fetcherDefinition{
			{ResourceType: graph.Instance, AWSType: "Instance", ApiMethod: "DescribeInstances", Input: "DescribeInstancesInput", Output: "DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations"},
			{ResourceType: graph.Subnet, AWSType: "Subnet", ApiMethod: "DescribeSubnets", Input: "DescribeSubnetsInput", Output: "DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{ResourceType: graph.Vpc, AWSType: "Vpc", ApiMethod: "DescribeVpcs", Input: "DescribeVpcsInput", Output: "DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{ResourceType: graph.Keypair, AWSType: "KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "DescribeKeyPairsInput", Output: "DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{ResourceType: graph.SecurityGroup, AWSType: "SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "DescribeSecurityGroupsInput", Output: "DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{ResourceType: graph.Volume, AWSType: "Volume", ApiMethod: "DescribeVolumes", Input: "DescribeVolumesInput", Output: "DescribeVolumesOutput", OutputsExtractor: "Volumes"},
			{ResourceType: graph.Region, AWSType: "Region", ApiMethod: "DescribeRegions", Input: "DescribeRegionsInput", Output: "DescribeRegionsOutput", OutputsExtractor: "Regions"},
			{ResourceType: graph.InternetGateway, AWSType: "InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "DescribeInternetGatewaysInput", Output: "DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{ResourceType: graph.RouteTable, AWSType: "RouteTable", ApiMethod: "DescribeRouteTables", Input: "DescribeRouteTablesInput", Output: "DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
			{ResourceType: graph.Image, AWSType: "Image", ApiMethod: "DescribeImages", Input: "DescribeImagesInput", Output: "DescribeImagesOutput", OutputsExtractor: "Images"},
		},
	},
}
