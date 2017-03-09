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
	Api                                         string
}

var FetchersDefs = []fetchersDef{
	{
		Name: "infra",
		Api:  []string{"ec2", "elbv2"},
		Fetchers: []fetcher{
			{Api: "ec2", ResourceType: graph.Instance.String(), AWSType: "ec2.Instance", ApiMethod: "DescribeInstancesPages", Input: "ec2.DescribeInstancesInput{}", Output: "ec2.DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ec2", ResourceType: graph.Subnet.String(), AWSType: "ec2.Subnet", ApiMethod: "DescribeSubnets", Input: "ec2.DescribeSubnetsInput{}", Output: "ec2.DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{Api: "ec2", ResourceType: graph.Vpc.String(), AWSType: "ec2.Vpc", ApiMethod: "DescribeVpcs", Input: "ec2.DescribeVpcsInput{}", Output: "ec2.DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{Api: "ec2", ResourceType: graph.Keypair.String(), AWSType: "ec2.KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "ec2.DescribeKeyPairsInput{}", Output: "ec2.DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{Api: "ec2", ResourceType: graph.SecurityGroup.String(), AWSType: "ec2.SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "ec2.DescribeSecurityGroupsInput{}", Output: "ec2.DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{Api: "ec2", ResourceType: graph.Volume.String(), AWSType: "ec2.Volume", ApiMethod: "DescribeVolumesPages", Input: "ec2.DescribeVolumesInput{}", Output: "ec2.DescribeVolumesOutput", OutputsExtractor: "Volumes", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ec2", ResourceType: graph.InternetGateway.String(), AWSType: "ec2.InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "ec2.DescribeInternetGatewaysInput{}", Output: "ec2.DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{Api: "ec2", ResourceType: graph.RouteTable.String(), AWSType: "ec2.RouteTable", ApiMethod: "DescribeRouteTables", Input: "ec2.DescribeRouteTablesInput{}", Output: "ec2.DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
			{Api: "ec2", ResourceType: graph.AvailabilityZone.String(), AWSType: "ec2.AvailabilityZone", ApiMethod: "DescribeAvailabilityZones", Input: "ec2.DescribeAvailabilityZonesInput{}", Output: "ec2.DescribeAvailabilityZonesOutput", OutputsExtractor: "AvailabilityZones"},
			{Api: "ec2", ResourceType: graph.LoadBalancer.String(), AWSType: "elbv2.LoadBalancer", ApiMethod: "DescribeLoadBalancersPages", Input: "elbv2.DescribeLoadBalancersInput{}", Output: "elbv2.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancers", Multipage: true, NextPageMarker: "NextMarker"},
			{Api: "elbv2", ResourceType: graph.TargetGroup.String(), AWSType: "elbv2.TargetGroup", ApiMethod: "DescribeTargetGroups", Input: "elbv2.DescribeTargetGroupsInput{}", Output: "elbv2.DescribeTargetGroupsOutput", OutputsExtractor: "TargetGroups"},
			{Api: "elbv2", ResourceType: graph.Listener.String(), AWSType: "elbv2.Listener", ManualFetcher: true},
		},
	},
	{
		Name: "access",
		Api:  []string{"iam"},
		Fetchers: []fetcher{
			{Api: "iam", ResourceType: graph.User.String(), AWSType: "iam.UserDetail", ManualFetcher: true},
			{Api: "iam", ResourceType: graph.Group.String(), AWSType: "iam.GroupDetail", ApiMethod: "GetAccountAuthorizationDetailsPages", Input: "iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeGroup)}}", Output: "iam.GetAccountAuthorizationDetailsOutput", OutputsExtractor: "GroupDetailList", Multipage: true, NextPageMarker: "Marker"},
			{Api: "iam", ResourceType: graph.Role.String(), AWSType: "iam.RoleDetail", ApiMethod: "GetAccountAuthorizationDetailsPages", Input: "iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeRole)}}", Output: "iam.GetAccountAuthorizationDetailsOutput", OutputsExtractor: "RoleDetailList", Multipage: true, NextPageMarker: "Marker"},
			{Api: "iam", ResourceType: graph.Policy.String(), AWSType: "iam.Policy", ApiMethod: "ListPoliciesPages", Input: "iam.ListPoliciesInput{OnlyAttached: awssdk.Bool(true)}", Output: "iam.ListPoliciesOutput", OutputsExtractor: "Policies", Multipage: true, NextPageMarker: "Marker"},
		},
	},
	{
		Name: "storage",
		Api:  []string{"s3"},
		Fetchers: []fetcher{
			{Api: "s3", ResourceType: graph.Bucket.String(), AWSType: "s3.Bucket", ManualFetcher: true},
			{Api: "s3", ResourceType: graph.Object.String(), AWSType: "s3.Object", ManualFetcher: true},
		},
	},
	{
		Name: "notification",
		Api:  []string{"sns"},
		Fetchers: []fetcher{
			{Api: "sns", ResourceType: graph.Subscription.String(), AWSType: "sns.Subscription", ApiMethod: "ListSubscriptionsPages", Input: "sns.ListSubscriptionsInput{}", Output: "sns.ListSubscriptionsOutput", OutputsExtractor: "Subscriptions", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "sns", ResourceType: graph.Topic.String(), AWSType: "sns.Topic", ApiMethod: "ListTopicsPages", Input: "sns.ListTopicsInput{}", Output: "sns.ListTopicsOutput", OutputsExtractor: "Topics", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "queue",
		Api:  []string{"sqs"},
		Fetchers: []fetcher{
			{Api: "sqs", ResourceType: graph.Queue.String(), AWSType: "string", ManualFetcher: true},
		},
	},
}
