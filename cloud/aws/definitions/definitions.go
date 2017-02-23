/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	Multipage                                   bool
}

var Services = []service{
	{
		Name: "infra",
		Api:  "ec2",
		Fetchers: []fetcher{
			{ResourceType: graph.Instance.String(), AWSType: "ec2.Instance", ApiMethod: "DescribeInstancesPages", Input: "DescribeInstancesInput{}", Output: "DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations", Multipage: true},
			{ResourceType: graph.Subnet.String(), AWSType: "ec2.Subnet", ApiMethod: "DescribeSubnets", Input: "DescribeSubnetsInput{}", Output: "DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{ResourceType: graph.Vpc.String(), AWSType: "ec2.Vpc", ApiMethod: "DescribeVpcs", Input: "DescribeVpcsInput{}", Output: "DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{ResourceType: graph.Keypair.String(), AWSType: "ec2.KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "DescribeKeyPairsInput{}", Output: "DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{ResourceType: graph.SecurityGroup.String(), AWSType: "ec2.SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "DescribeSecurityGroupsInput{}", Output: "DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{ResourceType: graph.Volume.String(), AWSType: "ec2.Volume", ApiMethod: "DescribeVolumesPages", Input: "DescribeVolumesInput{}", Output: "DescribeVolumesOutput", OutputsExtractor: "Volumes", Multipage: true},
			{ResourceType: graph.InternetGateway.String(), AWSType: "ec2.InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "DescribeInternetGatewaysInput{}", Output: "DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{ResourceType: graph.RouteTable.String(), AWSType: "ec2.RouteTable", ApiMethod: "DescribeRouteTables", Input: "DescribeRouteTablesInput{}", Output: "DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
		},
	},
	{
		Name: "access",
		Api:  "iam",
		Fetchers: []fetcher{
			{ResourceType: graph.User.String(), AWSType: "iam.UserDetail", ManualFetcher: true},
			{ResourceType: graph.Group.String(), AWSType: "iam.GroupDetail", ApiMethod: "GetAccountAuthorizationDetails", Input: "GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeGroup)}}", Output: "GetAccountAuthorizationDetailsOutput", OutputsExtractor: "GroupDetailList"},
			{ResourceType: graph.Role.String(), AWSType: "iam.RoleDetail", ApiMethod: "GetAccountAuthorizationDetails", Input: "GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeRole)}}", Output: "GetAccountAuthorizationDetailsOutput", OutputsExtractor: "RoleDetailList"},
			{ResourceType: graph.Policy.String(), AWSType: "iam.Policy", ApiMethod: "ListPolicies", Input: "ListPoliciesInput{OnlyAttached: awssdk.Bool(true)}", Output: "ListPoliciesOutput", OutputsExtractor: "Policies"},
		},
	},
	{
		Name: "storage",
		Api:  "s3",
		Fetchers: []fetcher{
			{ResourceType: graph.Bucket.String(), AWSType: "s3.Bucket", ManualFetcher: true},
			{ResourceType: graph.Object.String(), AWSType: "s3.Object", ManualFetcher: true},
		},
	},
	{
		Name: "notification",
		Api:  "sns",
		Fetchers: []fetcher{
			{ResourceType: graph.Subscription.String(), AWSType: "sns.Subscription", ApiMethod: "ListSubscriptionsPages", Input: "ListSubscriptionsInput{}", Output: "ListSubscriptionsOutput", OutputsExtractor: "Subscriptions", Multipage: true},
			{ResourceType: graph.Topic.String(), AWSType: "sns.Topic", ApiMethod: "ListTopicsPages", Input: "ListTopicsInput{}", Output: "ListTopicsOutput", OutputsExtractor: "Topics", Multipage: true},
		},
	},
	{
		Name: "queue",
		Api:  "sqs",
		Fetchers: []fetcher{
			{ResourceType: graph.Queue.String(), AWSType: "string", ManualFetcher: true},
		},
	},
}
