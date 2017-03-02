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

package aws

import "github.com/wallix/awless/graph"

type fetchersDef struct {
	Name     string
	Api      []string
	Fetchers []fetcher
}

type fetcher struct {
	ResourceType                                string
	AWSType                                     string
	ApiMethod, Input                            string
	Output, OutputsContainers, OutputsExtractor string
	ManualFetcher                               bool
	Multipage                                   bool
	NextPageMarker                              string
}

var FetchersDefs = []fetchersDef{
	{
		Name: "infra",
		Api:  []string{"ec2", "elbv2"},
		Fetchers: []fetcher{
			{ResourceType: graph.Instance.String(), AWSType: "ec2.Instance", ApiMethod: "DescribeInstancesPages", Input: "ec2.DescribeInstancesInput{}", Output: "ec2.DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations", Multipage: true, NextPageMarker: "NextToken"},
			{ResourceType: graph.Subnet.String(), AWSType: "ec2.Subnet", ApiMethod: "DescribeSubnets", Input: "ec2.DescribeSubnetsInput{}", Output: "ec2.DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{ResourceType: graph.Vpc.String(), AWSType: "ec2.Vpc", ApiMethod: "DescribeVpcs", Input: "ec2.DescribeVpcsInput{}", Output: "ec2.DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{ResourceType: graph.Keypair.String(), AWSType: "ec2.KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "ec2.DescribeKeyPairsInput{}", Output: "ec2.DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{ResourceType: graph.SecurityGroup.String(), AWSType: "ec2.SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "ec2.DescribeSecurityGroupsInput{}", Output: "ec2.DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{ResourceType: graph.Volume.String(), AWSType: "ec2.Volume", ApiMethod: "DescribeVolumesPages", Input: "ec2.DescribeVolumesInput{}", Output: "ec2.DescribeVolumesOutput", OutputsExtractor: "Volumes", Multipage: true, NextPageMarker: "NextToken"},
			{ResourceType: graph.InternetGateway.String(), AWSType: "ec2.InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "ec2.DescribeInternetGatewaysInput{}", Output: "ec2.DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{ResourceType: graph.RouteTable.String(), AWSType: "ec2.RouteTable", ApiMethod: "DescribeRouteTables", Input: "ec2.DescribeRouteTablesInput{}", Output: "ec2.DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
			{ResourceType: graph.AvailabilityZone.String(), AWSType: "ec2.AvailabilityZone", ApiMethod: "DescribeAvailabilityZones", Input: "ec2.DescribeAvailabilityZonesInput{}", Output: "ec2.DescribeAvailabilityZonesOutput", OutputsExtractor: "AvailabilityZones"},
			{ResourceType: graph.LoadBalancer.String(), AWSType: "elbv2.LoadBalancer", ApiMethod: "DescribeLoadBalancersPages", Input: "elbv2.DescribeLoadBalancersInput{}", Output: "elbv2.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancers", Multipage: true, NextPageMarker: "NextMarker"},
		},
	},
	{
		Name: "access",
		Api:  []string{"iam"},
		Fetchers: []fetcher{
			{ResourceType: graph.User.String(), AWSType: "iam.UserDetail", ManualFetcher: true},
			{ResourceType: graph.Group.String(), AWSType: "iam.GroupDetail", ApiMethod: "GetAccountAuthorizationDetails", Input: "iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeGroup)}}", Output: "iam.GetAccountAuthorizationDetailsOutput", OutputsExtractor: "GroupDetailList"},
			{ResourceType: graph.Role.String(), AWSType: "iam.RoleDetail", ApiMethod: "GetAccountAuthorizationDetails", Input: "iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeRole)}}", Output: "iam.GetAccountAuthorizationDetailsOutput", OutputsExtractor: "RoleDetailList"},
			{ResourceType: graph.Policy.String(), AWSType: "iam.Policy", ApiMethod: "ListPolicies", Input: "iam.ListPoliciesInput{OnlyAttached: awssdk.Bool(true)}", Output: "iam.ListPoliciesOutput", OutputsExtractor: "Policies"},
		},
	},
	{
		Name: "storage",
		Api:  []string{"s3"},
		Fetchers: []fetcher{
			{ResourceType: graph.Bucket.String(), AWSType: "s3.Bucket", ManualFetcher: true},
			{ResourceType: graph.Object.String(), AWSType: "s3.Object", ManualFetcher: true},
		},
	},
	{
		Name: "notification",
		Api:  []string{"sns"},
		Fetchers: []fetcher{
			{ResourceType: graph.Subscription.String(), AWSType: "sns.Subscription", ApiMethod: "ListSubscriptionsPages", Input: "sns.ListSubscriptionsInput{}", Output: "sns.ListSubscriptionsOutput", OutputsExtractor: "Subscriptions", Multipage: true, NextPageMarker: "NextToken"},
			{ResourceType: graph.Topic.String(), AWSType: "sns.Topic", ApiMethod: "ListTopicsPages", Input: "sns.ListTopicsInput{}", Output: "sns.ListTopicsOutput", OutputsExtractor: "Topics", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "queue",
		Api:  []string{"sqs"},
		Fetchers: []fetcher{
			{ResourceType: graph.Queue.String(), AWSType: "string", ManualFetcher: true},
		},
	},
}
