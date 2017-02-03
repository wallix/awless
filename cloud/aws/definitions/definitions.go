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
	ManualFetcher                               bool
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
		},
	},
	{
		Name:              "access",
		Api:               "iam",
		ManualGlobalFetch: true,
		Fetchers: []fetcher{
			{ResourceType: "user", AWSType: "User", ManualFetcher: true},
			{ResourceType: "group", AWSType: "GroupDetail", ApiMethod: "GetAccountAuthorizationDetails", Input: "GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeUser),awssdk.String(iam.EntityTypeRole),awssdk.String(iam.EntityTypeGroup),awssdk.String(iam.EntityTypeLocalManagedPolicy),awssdk.String(iam.EntityTypeAwsmanagedPolicy)}}", Output: "GetAccountAuthorizationDetailsOutput", OutputsExtractor: "GroupDetailList"},
			{ResourceType: "role", AWSType: "RoleDetail", ApiMethod: "GetAccountAuthorizationDetails", Input: "GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeUser),awssdk.String(iam.EntityTypeRole),awssdk.String(iam.EntityTypeGroup),awssdk.String(iam.EntityTypeLocalManagedPolicy),awssdk.String(iam.EntityTypeAwsmanagedPolicy)}}", Output: "GetAccountAuthorizationDetailsOutput", OutputsExtractor: "RoleDetailList"},
			{ResourceType: "policy", AWSType: "Policy", ApiMethod: "ListPolicies", Input: "ListPoliciesInput{OnlyAttached: awssdk.Bool(true)}", Output: "ListPoliciesOutput", OutputsExtractor: "Policies"},
		},
	},
}
