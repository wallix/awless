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

import (
	"github.com/wallix/awless/cloud"
)

type fetchersDef struct {
	Name          string
	Api           []string
	ApiInterfaces map[string]string
	Fetchers      []fetcher
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
		Name:          "infra",
		Api:           []string{"ec2", "elbv2", "rds", "autoscaling"},
		ApiInterfaces: map[string]string{"autoscaling": "AutoScalingAPI"},
		Fetchers: []fetcher{
			{Api: "ec2", ResourceType: cloud.Instance, AWSType: "ec2.Instance", ApiMethod: "DescribeInstancesPages", Input: "ec2.DescribeInstancesInput{}", Output: "ec2.DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ec2", ResourceType: cloud.Subnet, AWSType: "ec2.Subnet", ApiMethod: "DescribeSubnets", Input: "ec2.DescribeSubnetsInput{}", Output: "ec2.DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{Api: "ec2", ResourceType: cloud.Vpc, AWSType: "ec2.Vpc", ApiMethod: "DescribeVpcs", Input: "ec2.DescribeVpcsInput{}", Output: "ec2.DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{Api: "ec2", ResourceType: cloud.Keypair, AWSType: "ec2.KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "ec2.DescribeKeyPairsInput{}", Output: "ec2.DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{Api: "ec2", ResourceType: cloud.SecurityGroup, AWSType: "ec2.SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "ec2.DescribeSecurityGroupsInput{}", Output: "ec2.DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{Api: "ec2", ResourceType: cloud.Volume, AWSType: "ec2.Volume", ApiMethod: "DescribeVolumesPages", Input: "ec2.DescribeVolumesInput{}", Output: "ec2.DescribeVolumesOutput", OutputsExtractor: "Volumes", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ec2", ResourceType: cloud.InternetGateway, AWSType: "ec2.InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "ec2.DescribeInternetGatewaysInput{}", Output: "ec2.DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{Api: "ec2", ResourceType: cloud.RouteTable, AWSType: "ec2.RouteTable", ApiMethod: "DescribeRouteTables", Input: "ec2.DescribeRouteTablesInput{}", Output: "ec2.DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
			{Api: "ec2", ResourceType: cloud.AvailabilityZone, AWSType: "ec2.AvailabilityZone", ApiMethod: "DescribeAvailabilityZones", Input: "ec2.DescribeAvailabilityZonesInput{}", Output: "ec2.DescribeAvailabilityZonesOutput", OutputsExtractor: "AvailabilityZones"},
			{Api: "ec2", ResourceType: cloud.Image, AWSType: "ec2.Image", ApiMethod: "DescribeImages", Input: "ec2.DescribeImagesInput{Owners: []*string{awssdk.String(\"self\")}}", Output: "ec2.DescribeImagesOutput", OutputsExtractor: "Images"},
			{Api: "ec2", ResourceType: cloud.ElasticIP, AWSType: "ec2.Address", ApiMethod: "DescribeAddresses", Input: "ec2.DescribeAddressesInput{}", Output: "ec2.DescribeAddressesOutput", OutputsExtractor: "Addresses"},
			{Api: "elbv2", ResourceType: cloud.LoadBalancer, AWSType: "elbv2.LoadBalancer", ApiMethod: "DescribeLoadBalancersPages", Input: "elbv2.DescribeLoadBalancersInput{}", Output: "elbv2.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancers", Multipage: true, NextPageMarker: "NextMarker"},
			{Api: "elbv2", ResourceType: cloud.TargetGroup, AWSType: "elbv2.TargetGroup", ApiMethod: "DescribeTargetGroups", Input: "elbv2.DescribeTargetGroupsInput{}", Output: "elbv2.DescribeTargetGroupsOutput", OutputsExtractor: "TargetGroups"},
			{Api: "elbv2", ResourceType: cloud.Listener, AWSType: "elbv2.Listener", ManualFetcher: true},
			{Api: "rds", ResourceType: cloud.Database, AWSType: "rds.DBInstance", ApiMethod: "DescribeDBInstancesPages", Input: "rds.DescribeDBInstancesInput{}", Output: "rds.DescribeDBInstancesOutput", OutputsExtractor: "DBInstances", Multipage: true, NextPageMarker: "Marker"},
			{Api: "rds", ResourceType: cloud.DbSubnetGroup, AWSType: "rds.DBSubnetGroup", ApiMethod: "DescribeDBSubnetGroupsPages", Input: "rds.DescribeDBSubnetGroupsInput{}", Output: "rds.DescribeDBSubnetGroupsOutput", OutputsExtractor: "DBSubnetGroups", Multipage: true, NextPageMarker: "Marker"},
			{Api: "autoscaling", ResourceType: cloud.LaunchConfiguration, AWSType: "autoscaling.LaunchConfiguration", ApiMethod: "DescribeLaunchConfigurationsPages", Input: "autoscaling.DescribeLaunchConfigurationsInput{}", Output: "autoscaling.DescribeLaunchConfigurationsOutput", OutputsExtractor: "LaunchConfigurations", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "autoscaling", ResourceType: cloud.ScalingGroup, AWSType: "autoscaling.Group", ApiMethod: "DescribeAutoScalingGroupsPages", Input: "autoscaling.DescribeAutoScalingGroupsInput{}", Output: "autoscaling.DescribeAutoScalingGroupsOutput", OutputsExtractor: "AutoScalingGroups", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "autoscaling", ResourceType: cloud.ScalingPolicy, AWSType: "autoscaling.ScalingPolicy", ApiMethod: "DescribePoliciesPages", Input: "autoscaling.DescribePoliciesInput{}", Output: "autoscaling.DescribePoliciesOutput", OutputsExtractor: "ScalingPolicies", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "access",
		Api:  []string{"iam", "sts"},
		Fetchers: []fetcher{
			{Api: "iam", ResourceType: cloud.User, AWSType: "iam.UserDetail", ManualFetcher: true},
			{Api: "iam", ResourceType: cloud.Group, AWSType: "iam.GroupDetail", ApiMethod: "GetAccountAuthorizationDetailsPages", Input: "iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeGroup)}}", Output: "iam.GetAccountAuthorizationDetailsOutput", OutputsExtractor: "GroupDetailList", Multipage: true, NextPageMarker: "Marker"},
			{Api: "iam", ResourceType: cloud.Role, AWSType: "iam.RoleDetail", ApiMethod: "GetAccountAuthorizationDetailsPages", Input: "iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeRole)}}", Output: "iam.GetAccountAuthorizationDetailsOutput", OutputsExtractor: "RoleDetailList", Multipage: true, NextPageMarker: "Marker"},
			{Api: "iam", ResourceType: cloud.Policy, AWSType: "iam.Policy", ApiMethod: "ListPoliciesPages", Input: "iam.ListPoliciesInput{OnlyAttached: awssdk.Bool(true)}", Output: "iam.ListPoliciesOutput", OutputsExtractor: "Policies", Multipage: true, NextPageMarker: "Marker"},
			{Api: "iam", ResourceType: cloud.AccessKey, AWSType: "iam.AccessKeyMetadata", ApiMethod: "ListAccessKeysPages", Input: "iam.ListAccessKeysInput{}", Output: "iam.ListAccessKeysOutput", OutputsExtractor: "AccessKeyMetadata", Multipage: true, NextPageMarker: "Marker"},
		},
	},
	{
		Name: "storage",
		Api:  []string{"s3"},
		Fetchers: []fetcher{
			{Api: "s3", ResourceType: cloud.Bucket, AWSType: "s3.Bucket", ManualFetcher: true},
			{Api: "s3", ResourceType: cloud.S3Object, AWSType: "s3.Object", ManualFetcher: true},
		},
	},
	{
		Name: "notification",
		Api:  []string{"sns"},
		Fetchers: []fetcher{
			{Api: "sns", ResourceType: cloud.Subscription, AWSType: "sns.Subscription", ApiMethod: "ListSubscriptionsPages", Input: "sns.ListSubscriptionsInput{}", Output: "sns.ListSubscriptionsOutput", OutputsExtractor: "Subscriptions", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "sns", ResourceType: cloud.Topic, AWSType: "sns.Topic", ApiMethod: "ListTopicsPages", Input: "sns.ListTopicsInput{}", Output: "sns.ListTopicsOutput", OutputsExtractor: "Topics", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "queue",
		Api:  []string{"sqs"},
		Fetchers: []fetcher{
			{Api: "sqs", ResourceType: cloud.Queue, AWSType: "string", ManualFetcher: true},
		},
	},
	{
		Name:          "dns",
		Api:           []string{"route53"},
		ApiInterfaces: map[string]string{"route53": "Route53API"},
		Fetchers: []fetcher{
			{Api: "route53", ResourceType: cloud.Zone, AWSType: "route53.HostedZone", ApiMethod: "ListHostedZonesPages", Input: "route53.ListHostedZonesInput{}", Output: "route53.ListHostedZonesOutput", OutputsExtractor: "HostedZones", Multipage: true, NextPageMarker: "NextMarker"},
			{Api: "route53", ResourceType: cloud.Record, AWSType: "route53.ResourceRecordSet", ManualFetcher: true},
		},
	},

	{
		Name:          "lambda",
		Api:           []string{"lambda"},
		ApiInterfaces: map[string]string{"lambda": "LambdaAPI"},
		Fetchers: []fetcher{
			{Api: "lambda", ResourceType: cloud.Function, AWSType: "lambda.FunctionConfiguration", ApiMethod: "ListFunctionsPages", Input: "lambda.ListFunctionsInput{}", Output: "lambda.ListFunctionsOutput", OutputsExtractor: "Functions", Multipage: true, NextPageMarker: "NextMarker"},
		},
	},
	{
		Name:          "monitoring",
		Api:           []string{"cloudwatch"},
		ApiInterfaces: map[string]string{"cloudwatch": "CloudWatchAPI"},
		Fetchers: []fetcher{
			{Api: "cloudwatch", ResourceType: cloud.Metric, AWSType: "cloudwatch.Metric", ApiMethod: "ListMetricsPages", Input: "cloudwatch.ListMetricsInput{}", Output: "cloudwatch.ListMetricsOutput", OutputsExtractor: "Metrics", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "cloudwatch", ResourceType: cloud.Alarm, AWSType: "cloudwatch.MetricAlarm", ApiMethod: "DescribeAlarmsPages", Input: "cloudwatch.DescribeAlarmsInput{}", Output: "cloudwatch.DescribeAlarmsOutput", OutputsExtractor: "MetricAlarms", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
}
