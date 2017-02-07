package definitions

import "github.com/wallix/awless/graph"

type service struct {
	Name     string
	Api      string
	Fetchers []fetcher
}

type fetcher struct {
	ResourceType                                string
	AWSType                                     string
	ApiMethod, Input                            string
	Output, OutputsContainers, OutputsExtractor string
	ManualFetcher                               bool
}

var Services = []service{
	{
		Name: "infra",
		Api:  "ec2",
		Fetchers: []fetcher{
			{ResourceType: graph.Instance.String(), AWSType: "Instance", ApiMethod: "DescribeInstances", Input: "DescribeInstancesInput{}", Output: "DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations"},
			{ResourceType: graph.Subnet.String(), AWSType: "Subnet", ApiMethod: "DescribeSubnets", Input: "DescribeSubnetsInput{}", Output: "DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{ResourceType: graph.Vpc.String(), AWSType: "Vpc", ApiMethod: "DescribeVpcs", Input: "DescribeVpcsInput{}", Output: "DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{ResourceType: graph.Keypair.String(), AWSType: "KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "DescribeKeyPairsInput{}", Output: "DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{ResourceType: graph.SecurityGroup.String(), AWSType: "SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "DescribeSecurityGroupsInput{}", Output: "DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{ResourceType: graph.Volume.String(), AWSType: "Volume", ApiMethod: "DescribeVolumes", Input: "DescribeVolumesInput{}", Output: "DescribeVolumesOutput", OutputsExtractor: "Volumes"},
			{ResourceType: graph.InternetGateway.String(), AWSType: "InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "DescribeInternetGatewaysInput{}", Output: "DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{ResourceType: graph.RouteTable.String(), AWSType: "RouteTable", ApiMethod: "DescribeRouteTables", Input: "DescribeRouteTablesInput{}", Output: "DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
		},
	},
	{
		Name: "access",
		Api:  "iam",
		Fetchers: []fetcher{
			{ResourceType: graph.User.String(), AWSType: "UserDetail", ManualFetcher: true},
			{ResourceType: graph.Group.String(), AWSType: "GroupDetail", ApiMethod: "GetAccountAuthorizationDetails", Input: "GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeGroup)}}", Output: "GetAccountAuthorizationDetailsOutput", OutputsExtractor: "GroupDetailList"},
			{ResourceType: graph.Role.String(), AWSType: "RoleDetail", ApiMethod: "GetAccountAuthorizationDetails", Input: "GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeRole)}}", Output: "GetAccountAuthorizationDetailsOutput", OutputsExtractor: "RoleDetailList"},
			{ResourceType: graph.Policy.String(), AWSType: "Policy", ApiMethod: "ListPolicies", Input: "ListPoliciesInput{OnlyAttached: awssdk.Bool(true)}", Output: "ListPoliciesOutput", OutputsExtractor: "Policies"},
		},
	},
	{
		Name: "storage",
		Api:  "s3",
		Fetchers: []fetcher{
			{ResourceType: "bucket", AWSType: "Bucket", ManualFetcher: true},
			{ResourceType: graph.Object.String(), AWSType: "Object", ManualFetcher: true},
		},
	},
}
