package definitions

type service struct {
	Name              string
	Api               string
	Fetchers          []fetcher
	ManualGlobalFetch bool
}

type fetcher struct {
	ResourceType                                string
	AWSType                                     string
	ApiMethod, Input                            string
	Output, OutputsContainers, OutputsExtractor string
}

var Services = []service{
	{
		Name: "infra",
		Api:  "ec2",
		Fetchers: []fetcher{
			{ResourceType: "instance", AWSType: "Instance", ApiMethod: "DescribeInstances", Input: "DescribeInstancesInput{}", Output: "DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations"},
			{ResourceType: "subnet", AWSType: "Subnet", ApiMethod: "DescribeSubnets", Input: "DescribeSubnetsInput{}", Output: "DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{ResourceType: "vpc", AWSType: "Vpc", ApiMethod: "DescribeVpcs", Input: "DescribeVpcsInput{}", Output: "DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{ResourceType: "keypair", AWSType: "KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "DescribeKeyPairsInput{}", Output: "DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{ResourceType: "securitygroup", AWSType: "SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "DescribeSecurityGroupsInput{}", Output: "DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{ResourceType: "volume", AWSType: "Volume", ApiMethod: "DescribeVolumes", Input: "DescribeVolumesInput{}", Output: "DescribeVolumesOutput", OutputsExtractor: "Volumes"},
			{ResourceType: "region", AWSType: "Region", ApiMethod: "DescribeRegions", Input: "DescribeRegionsInput{}", Output: "DescribeRegionsOutput", OutputsExtractor: "Regions"},
			{ResourceType: "internetgateway", AWSType: "InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "DescribeInternetGatewaysInput{}", Output: "DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{ResourceType: "routetable", AWSType: "RouteTable", ApiMethod: "DescribeRouteTables", Input: "DescribeRouteTablesInput{}", Output: "DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
			{ResourceType: "image", AWSType: "Image", ApiMethod: "DescribeImages", Input: "DescribeImagesInput{}", Output: "DescribeImagesOutput", OutputsExtractor: "Images"},
		},
	},
	{
		Name:              "access",
		Api:               "iam",
		ManualGlobalFetch: true,
		Fetchers: []fetcher{
			{ResourceType: "user", AWSType: "User", ApiMethod: "ListUsers", Input: "ListUsersInput{}", Output: "ListUsersOutput", OutputsExtractor: "Users"},
			{ResourceType: "group", AWSType: "Group", ApiMethod: "ListGroups", Input: "ListGroupsInput{}", Output: "ListGroupsOutput", OutputsExtractor: "Groups"},
			{ResourceType: "role", AWSType: "Role", ApiMethod: "ListRoles", Input: "ListRolesInput{}", Output: "ListRolesOutput", OutputsExtractor: "Roles"},
			{ResourceType: "policy", AWSType: "Policy", ApiMethod: "ListPolicies", Input: "ListPoliciesInput{Scope: awssdk.String(iam.PolicyScopeTypeLocal)}", Output: "ListPoliciesOutput", OutputsExtractor: "Policies"},
		},
	},
}
