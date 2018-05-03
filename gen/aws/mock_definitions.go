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

import "strings"

type mockDef struct {
	Api, Name string
	Funcs     []*mockFuncDef
}

type mockFuncDef struct {
	FuncType, AWSType, ApiMethod, Input, Output, OutputsExtractor, OutputsContainers string
	Manual                                                                           bool
	Multipage                                                                        bool
	NextPageMarker                                                                   string
	MockField, MockFieldType                                                         string
}

var mocksDefs = []*mockDef{
	{
		Api: "ec2",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "ec2.Instance", ApiMethod: "DescribeInstancesPages", Input: "ec2.DescribeInstancesInput", Output: "ec2.DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations", Multipage: true, NextPageMarker: "NextToken", Manual: true},
			{FuncType: "list", AWSType: "ec2.Subnet", ApiMethod: "DescribeSubnets", Input: "ec2.DescribeSubnetsInput", Output: "ec2.DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{FuncType: "list", AWSType: "ec2.Vpc", ApiMethod: "DescribeVpcs", Input: "ec2.DescribeVpcsInput", Output: "ec2.DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{FuncType: "list", AWSType: "ec2.KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "ec2.DescribeKeyPairsInput", Output: "ec2.DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{FuncType: "list", AWSType: "ec2.SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "ec2.DescribeSecurityGroupsInput", Output: "ec2.DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{FuncType: "list", AWSType: "ec2.Volume", ApiMethod: "DescribeVolumesPages", Input: "ec2.DescribeVolumesInput", Output: "ec2.DescribeVolumesOutput", OutputsExtractor: "Volumes", Multipage: true, NextPageMarker: "NextToken"},
			{FuncType: "list", AWSType: "ec2.InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "ec2.DescribeInternetGatewaysInput", Output: "ec2.DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{FuncType: "list", AWSType: "ec2.NatGateway", ApiMethod: "DescribeNatGateways", Input: "ec2.DescribeNatGatewaysInput", Output: "ec2.DescribeNatGatewaysOutput", OutputsExtractor: "NatGateways"},
			{FuncType: "list", AWSType: "ec2.RouteTable", ApiMethod: "DescribeRouteTables", Input: "ec2.DescribeRouteTablesInput", Output: "ec2.DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
			{FuncType: "list", AWSType: "ec2.AvailabilityZone", ApiMethod: "DescribeAvailabilityZones", Input: "ec2.DescribeAvailabilityZonesInput", Output: "ec2.DescribeAvailabilityZonesOutput", OutputsExtractor: "AvailabilityZones"},
			{FuncType: "list", AWSType: "ec2.Image", ApiMethod: "DescribeImages", Input: "ec2.DescribeImagesInput", Output: "ec2.DescribeImagesOutput", OutputsExtractor: "Images"},
			{FuncType: "list", AWSType: "ec2.ImportImageTask", ApiMethod: "DescribeImportImageTasks", Input: "ec2.DescribeImportImageTasksInput", Output: "ec2.DescribeImportImageTasksOutput", OutputsExtractor: "ImportImageTasks"},
			{FuncType: "list", AWSType: "ec2.Address", ApiMethod: "DescribeAddresses", Input: "ec2.DescribeAddressesInput", Output: "ec2.DescribeAddressesOutput", OutputsExtractor: "Addresses"},
			{FuncType: "list", AWSType: "ec2.Snapshot", ApiMethod: "DescribeSnapshotsPages", Input: "ec2.DescribeSnapshotsInput", Output: "ec2.DescribeSnapshotsOutput", OutputsExtractor: "Snapshots", Multipage: true, NextPageMarker: "NextToken"},
			{FuncType: "list", AWSType: "ec2.NetworkInterface", ApiMethod: "DescribeNetworkInterfaces", Input: "ec2.DescribeNetworkInterfacesInput", Output: "ec2.DescribeNetworkInterfacesOutput", OutputsExtractor: "NetworkInterfaces"},
		},
	},
	{
		Api: "elbv2",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "elbv2.LoadBalancer", ApiMethod: "DescribeLoadBalancersPages", Input: "elbv2.DescribeLoadBalancersInput", Output: "elbv2.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancers", Multipage: true, NextPageMarker: "NextMarker"},
			{FuncType: "list", AWSType: "elbv2.TargetGroup", ApiMethod: "DescribeTargetGroups", Input: "elbv2.DescribeTargetGroupsInput", Output: "elbv2.DescribeTargetGroupsOutput", OutputsExtractor: "TargetGroups"},
			{FuncType: "list", AWSType: "elbv2.Listener", Manual: true},
			{FuncType: "list", AWSType: "elbv2.TargetHealthDescription", Manual: true, MockFieldType: "mapslice"},
		},
	},
	{
		Api: "elb",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "elb.LoadBalancerDescription", ApiMethod: "DescribeLoadBalancersPages", Input: "elb.DescribeLoadBalancersInput", Output: "elb.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancerDescriptions", Multipage: true, NextPageMarker: "NextMarker"},
		},
	},
	{
		Api: "rds",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "rds.DBInstance", ApiMethod: "DescribeDBInstancesPages", Input: "rds.DescribeDBInstancesInput", Output: "rds.DescribeDBInstancesOutput", OutputsExtractor: "DBInstances", Multipage: true, NextPageMarker: "Marker"},
			{FuncType: "list", AWSType: "rds.DBSubnetGroup", ApiMethod: "DescribeDBSubnetGroupsPages", Input: "rds.DescribeDBSubnetGroupsInput", Output: "rds.DescribeDBSubnetGroupsOutput", OutputsExtractor: "DBSubnetGroups", Multipage: true, NextPageMarker: "Marker"},
		},
	},
	{
		Api: "autoscaling",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "autoscaling.LaunchConfiguration", ApiMethod: "DescribeLaunchConfigurationsPages", Input: "autoscaling.DescribeLaunchConfigurationsInput", Output: "autoscaling.DescribeLaunchConfigurationsOutput", OutputsExtractor: "LaunchConfigurations", Multipage: true, NextPageMarker: "NextToken"},
			{FuncType: "list", AWSType: "autoscaling.Group", ApiMethod: "DescribeAutoScalingGroupsPages", Input: "autoscaling.DescribeAutoScalingGroupsInput", Output: "autoscaling.DescribeAutoScalingGroupsOutput", OutputsExtractor: "AutoScalingGroups", Multipage: true, NextPageMarker: "NextToken"},
			{FuncType: "list", AWSType: "autoscaling.ScalingPolicy", ApiMethod: "DescribePoliciesPages", Input: "autoscaling.DescribePoliciesInput", Output: "autoscaling.DescribePoliciesOutput", OutputsExtractor: "ScalingPolicies", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Api: "acm",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "acm.CertificateSummary", ApiMethod: "ListCertificatesPages", Input: "acm.ListCertificatesInput", Output: "acm.ListCertificatesOutput", OutputsExtractor: "CertificateSummaryList", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Api: "iam",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "iam.UserDetail", Manual: true},
			{FuncType: "list", AWSType: "iam.GroupDetail", ApiMethod: "GetAccountAuthorizationDetailsPages", Input: "iam.GetAccountAuthorizationDetailsInput", Output: "iam.GetAccountAuthorizationDetailsOutput", OutputsExtractor: "GroupDetailList", Multipage: true, NextPageMarker: "Marker", Manual: true},
			{FuncType: "list", AWSType: "iam.RoleDetail", ApiMethod: "GetAccountAuthorizationDetailsPages", Input: "iam.GetAccountAuthorizationDetailsInput", Output: "iam.GetAccountAuthorizationDetailsOutput", OutputsExtractor: "RoleDetailList", Multipage: true, NextPageMarker: "Marker", Manual: true},
			{FuncType: "list", AWSType: "iam.Policy", ApiMethod: "ListPoliciesPages", Input: "iam.ListPoliciesInput", Output: "iam.ListPoliciesOutput", OutputsExtractor: "Policies", Multipage: true, NextPageMarker: "Marker", Manual: true},
			{FuncType: "list", AWSType: "iam.AccessKeyMetadata", ApiMethod: "ListAccessKeysPages", Input: "iam.ListAccessKeysInput", Output: "iam.ListAccessKeysOutput", OutputsExtractor: "AccessKeyMetadata", Multipage: true, NextPageMarker: "Marker"},
			{FuncType: "list", AWSType: "iam.InstanceProfile", ApiMethod: "ListInstanceProfilesPages", Input: "iam.ListInstanceProfilesInput", Output: "iam.ListInstanceProfilesOutput", OutputsExtractor: "InstanceProfiles", Multipage: true, NextPageMarker: "Marker"},
			{FuncType: "list", AWSType: "iam.ManagedPolicyDetail", Manual: true},
			{FuncType: "list", AWSType: "iam.User", Manual: true},
			{FuncType: "list", AWSType: "iam.VirtualMFADevice", ApiMethod: "ListVirtualMFADevicesPages", Input: "iam.ListVirtualMFADevicesInput", Output: "iam.ListVirtualMFADevicesOutput", OutputsExtractor: "VirtualMFADevices", Multipage: true, NextPageMarker: "Marker"},
		},
	},
	{
		Api: "s3",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "s3.Bucket", Manual: true, MockFieldType: "mapslice"},
			{FuncType: "list", AWSType: "s3.Object", Manual: true, MockFieldType: "mapslice"},
			{FuncType: "list", AWSType: "s3.Grant", Manual: true, MockFieldType: "mapslice"},
		},
	},
	{
		Api: "sns",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "sns.Subscription", ApiMethod: "ListSubscriptionsPages", Input: "sns.ListSubscriptionsInput", Output: "sns.ListSubscriptionsOutput", OutputsExtractor: "Subscriptions", Multipage: true, NextPageMarker: "NextToken"},
			{FuncType: "list", AWSType: "sns.Topic", ApiMethod: "ListTopicsPages", Input: "sns.ListTopicsInput", Output: "sns.ListTopicsOutput", OutputsExtractor: "Topics", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Api: "sqs",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "string", ApiMethod: "ListQueues", Input: "sqs.ListQueuesInput", Output: "sqs.ListQueuesOutput", OutputsExtractor: "QueueUrls"},
			{FuncType: "list", AWSType: "map[string]*string", Manual: true, MockFieldType: "map"},
		},
	},
	{
		Api: "route53",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "route53.HostedZone", ApiMethod: "ListHostedZonesPages", Input: "route53.ListHostedZonesInput", Output: "route53.ListHostedZonesOutput", OutputsExtractor: "HostedZones", Multipage: true, NextPageMarker: "NextMarker"},
			{FuncType: "list", AWSType: "route53.ResourceRecordSet", Manual: true, MockFieldType: "mapslice"},
		},
	},
	{
		Api: "lambda",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "lambda.FunctionConfiguration", ApiMethod: "ListFunctionsPages", Input: "lambda.ListFunctionsInput", Output: "lambda.ListFunctionsOutput", OutputsExtractor: "Functions", Multipage: true, NextPageMarker: "NextMarker"},
		},
	},
	{
		Api: "cloudwatch",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "cloudwatch.Metric", ApiMethod: "ListMetricsPages", Input: "cloudwatch.ListMetricsInput", Output: "cloudwatch.ListMetricsOutput", OutputsExtractor: "Metrics", Multipage: true, NextPageMarker: "NextToken"},
			{FuncType: "list", AWSType: "cloudwatch.MetricAlarm", ApiMethod: "DescribeAlarmsPages", Input: "cloudwatch.DescribeAlarmsInput", Output: "cloudwatch.DescribeAlarmsOutput", OutputsExtractor: "MetricAlarms", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Api: "cloudfront",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "cloudfront.DistributionSummary", Manual: true},
		},
	},
	{
		Api: "cloudformation",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "cloudformation.Stack", ApiMethod: "DescribeStacksPages", Input: "cloudformation.DescribeStacksInput", Output: "cloudformation.DescribeStacksOutput", OutputsExtractor: "Stacks", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Api: "ecr",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "ecr.Repository", ApiMethod: "DescribeRepositoriesPages", Input: "ecr.DescribeRepositoriesInput", Output: "ecr.DescribeRepositoriesOutput", OutputsExtractor: "Repositories", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Api: "ecs",
		Funcs: []*mockFuncDef{
			{FuncType: "list", AWSType: "ecs.Cluster", Manual: true},
			{FuncType: "list", MockField: "clusterNames", AWSType: "string", ApiMethod: "ListClustersPages", Input: "ecs.ListClustersInput", Output: "ecs.ListClustersOutput", OutputsExtractor: "ClusterArns", Multipage: true, NextPageMarker: "NextToken"},
			{FuncType: "list", AWSType: "ecs.TaskDefinition", Manual: true},
			{FuncType: "list", MockField: "taskdefinitionNames", AWSType: "string", ApiMethod: "ListTaskDefinitionsPages", Input: "ecs.ListTaskDefinitionsInput", Output: "ecs.ListTaskDefinitionsOutput", OutputsExtractor: "TaskDefinitionArns", Multipage: true, NextPageMarker: "NextToken"},
			{FuncType: "list", MockFieldType: "mapslice", AWSType: "ecs.Task", Manual: true},
			{FuncType: "list", MockFieldType: "mapslice", MockField: "tasksNames", AWSType: "string", Manual: true},
			{FuncType: "list", MockFieldType: "mapslice", MockField: "containerinstancesNames", AWSType: "string", Manual: true},
			{FuncType: "list", MockFieldType: "mapslice", AWSType: "ecs.ContainerInstance", Manual: true},
		},
	},
}

func Mocks() []*mockDef {
	for _, def := range mocksDefs {
		def.Name = "mock" + strings.Title(def.Api)
		for _, f := range def.Funcs {
			if f.MockField == "" {
				f.MockField = nameFromAwsType(f.AWSType)
			}
		}
	}
	return mocksDefs
}

func nameFromAwsType(awstype string) string {
	if awstype == "map[string]*string" {
		return "attributes"
	}
	splits := strings.Split(awstype, ".")
	return strings.ToLower(splits[len(splits)-1]) + "s"
}
