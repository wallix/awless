// Auto generated implementation for the AWS cloud service

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

package awsservices

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/wallix/awless/cloud"
)

// DO NOT EDIT - This file was automatically generated with go generate

type mockEc2 struct {
	ec2iface.EC2API
	instances         []*ec2.Instance
	subnets           []*ec2.Subnet
	vpcs              []*ec2.Vpc
	keypairinfos      []*ec2.KeyPairInfo
	securitygroups    []*ec2.SecurityGroup
	volumes           []*ec2.Volume
	internetgateways  []*ec2.InternetGateway
	natgateways       []*ec2.NatGateway
	routetables       []*ec2.RouteTable
	availabilityzones []*ec2.AvailabilityZone
	images            []*ec2.Image
	importimagetasks  []*ec2.ImportImageTask
	addresss          []*ec2.Address
	snapshots         []*ec2.Snapshot
	networkinterfaces []*ec2.NetworkInterface
}

func (m *mockEc2) Name() string {
	return ""
}

func (m *mockEc2) Region() string {
	return ""
}

func (m *mockEc2) Profile() string {
	return ""
}

func (m *mockEc2) Provider() string {
	return ""
}

func (m *mockEc2) ProviderAPI() string {
	return ""
}

func (m *mockEc2) ResourceTypes() []string {
	return []string{}
}

func (m *mockEc2) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEc2) IsSyncDisabled() bool {
	return false
}

func (m *mockEc2) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEc2) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return &ec2.DescribeSubnetsOutput{Subnets: m.subnets}, nil
}

func (m *mockEc2) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return &ec2.DescribeVpcsOutput{Vpcs: m.vpcs}, nil
}

func (m *mockEc2) DescribeKeyPairs(input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return &ec2.DescribeKeyPairsOutput{KeyPairs: m.keypairinfos}, nil
}

func (m *mockEc2) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return &ec2.DescribeSecurityGroupsOutput{SecurityGroups: m.securitygroups}, nil
}

func (m *mockEc2) DescribeVolumesPages(input *ec2.DescribeVolumesInput, fn func(p *ec2.DescribeVolumesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*ec2.Volume
	for i := 0; i < len(m.volumes); i += 2 {
		page := []*ec2.Volume{m.volumes[i]}
		if i+1 < len(m.volumes) {
			page = append(page, m.volumes[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&ec2.DescribeVolumesOutput{Volumes: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockEc2) DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	return &ec2.DescribeInternetGatewaysOutput{InternetGateways: m.internetgateways}, nil
}

func (m *mockEc2) DescribeNatGateways(input *ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
	return &ec2.DescribeNatGatewaysOutput{NatGateways: m.natgateways}, nil
}

func (m *mockEc2) DescribeRouteTables(input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	return &ec2.DescribeRouteTablesOutput{RouteTables: m.routetables}, nil
}

func (m *mockEc2) DescribeAvailabilityZones(input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	return &ec2.DescribeAvailabilityZonesOutput{AvailabilityZones: m.availabilityzones}, nil
}

func (m *mockEc2) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	return &ec2.DescribeImagesOutput{Images: m.images}, nil
}

func (m *mockEc2) DescribeImportImageTasks(input *ec2.DescribeImportImageTasksInput) (*ec2.DescribeImportImageTasksOutput, error) {
	return &ec2.DescribeImportImageTasksOutput{ImportImageTasks: m.importimagetasks}, nil
}

func (m *mockEc2) DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	return &ec2.DescribeAddressesOutput{Addresses: m.addresss}, nil
}

func (m *mockEc2) DescribeSnapshotsPages(input *ec2.DescribeSnapshotsInput, fn func(p *ec2.DescribeSnapshotsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*ec2.Snapshot
	for i := 0; i < len(m.snapshots); i += 2 {
		page := []*ec2.Snapshot{m.snapshots[i]}
		if i+1 < len(m.snapshots) {
			page = append(page, m.snapshots[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&ec2.DescribeSnapshotsOutput{Snapshots: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockEc2) DescribeNetworkInterfaces(input *ec2.DescribeNetworkInterfacesInput) (*ec2.DescribeNetworkInterfacesOutput, error) {
	return &ec2.DescribeNetworkInterfacesOutput{NetworkInterfaces: m.networkinterfaces}, nil
}

type mockElbv2 struct {
	elbv2iface.ELBV2API
	loadbalancers            []*elbv2.LoadBalancer
	targetgroups             []*elbv2.TargetGroup
	listeners                []*elbv2.Listener
	targethealthdescriptions map[string][]*elbv2.TargetHealthDescription
}

func (m *mockElbv2) Name() string {
	return ""
}

func (m *mockElbv2) Region() string {
	return ""
}

func (m *mockElbv2) Profile() string {
	return ""
}

func (m *mockElbv2) Provider() string {
	return ""
}

func (m *mockElbv2) ProviderAPI() string {
	return ""
}

func (m *mockElbv2) ResourceTypes() []string {
	return []string{}
}

func (m *mockElbv2) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockElbv2) IsSyncDisabled() bool {
	return false
}

func (m *mockElbv2) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockElbv2) DescribeLoadBalancersPages(input *elbv2.DescribeLoadBalancersInput, fn func(p *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*elbv2.LoadBalancer
	for i := 0; i < len(m.loadbalancers); i += 2 {
		page := []*elbv2.LoadBalancer{m.loadbalancers[i]}
		if i+1 < len(m.loadbalancers) {
			page = append(page, m.loadbalancers[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&elbv2.DescribeLoadBalancersOutput{LoadBalancers: page, NextMarker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockElbv2) DescribeTargetGroups(input *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	return &elbv2.DescribeTargetGroupsOutput{TargetGroups: m.targetgroups}, nil
}

type mockElb struct {
	elbiface.ELBAPI
	loadbalancerdescriptions []*elb.LoadBalancerDescription
}

func (m *mockElb) Name() string {
	return ""
}

func (m *mockElb) Region() string {
	return ""
}

func (m *mockElb) Profile() string {
	return ""
}

func (m *mockElb) Provider() string {
	return ""
}

func (m *mockElb) ProviderAPI() string {
	return ""
}

func (m *mockElb) ResourceTypes() []string {
	return []string{}
}

func (m *mockElb) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockElb) IsSyncDisabled() bool {
	return false
}

func (m *mockElb) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockElb) DescribeLoadBalancersPages(input *elb.DescribeLoadBalancersInput, fn func(p *elb.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*elb.LoadBalancerDescription
	for i := 0; i < len(m.loadbalancerdescriptions); i += 2 {
		page := []*elb.LoadBalancerDescription{m.loadbalancerdescriptions[i]}
		if i+1 < len(m.loadbalancerdescriptions) {
			page = append(page, m.loadbalancerdescriptions[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&elb.DescribeLoadBalancersOutput{LoadBalancerDescriptions: page, NextMarker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockRds struct {
	rdsiface.RDSAPI
	dbinstances    []*rds.DBInstance
	dbsubnetgroups []*rds.DBSubnetGroup
}

func (m *mockRds) Name() string {
	return ""
}

func (m *mockRds) Region() string {
	return ""
}

func (m *mockRds) Profile() string {
	return ""
}

func (m *mockRds) Provider() string {
	return ""
}

func (m *mockRds) ProviderAPI() string {
	return ""
}

func (m *mockRds) ResourceTypes() []string {
	return []string{}
}

func (m *mockRds) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockRds) IsSyncDisabled() bool {
	return false
}

func (m *mockRds) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockRds) DescribeDBInstancesPages(input *rds.DescribeDBInstancesInput, fn func(p *rds.DescribeDBInstancesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*rds.DBInstance
	for i := 0; i < len(m.dbinstances); i += 2 {
		page := []*rds.DBInstance{m.dbinstances[i]}
		if i+1 < len(m.dbinstances) {
			page = append(page, m.dbinstances[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&rds.DescribeDBInstancesOutput{DBInstances: page, Marker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockRds) DescribeDBSubnetGroupsPages(input *rds.DescribeDBSubnetGroupsInput, fn func(p *rds.DescribeDBSubnetGroupsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*rds.DBSubnetGroup
	for i := 0; i < len(m.dbsubnetgroups); i += 2 {
		page := []*rds.DBSubnetGroup{m.dbsubnetgroups[i]}
		if i+1 < len(m.dbsubnetgroups) {
			page = append(page, m.dbsubnetgroups[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&rds.DescribeDBSubnetGroupsOutput{DBSubnetGroups: page, Marker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockAutoscaling struct {
	autoscalingiface.AutoScalingAPI
	launchconfigurations []*autoscaling.LaunchConfiguration
	groups               []*autoscaling.Group
	scalingpolicys       []*autoscaling.ScalingPolicy
}

func (m *mockAutoscaling) Name() string {
	return ""
}

func (m *mockAutoscaling) Region() string {
	return ""
}

func (m *mockAutoscaling) Profile() string {
	return ""
}

func (m *mockAutoscaling) Provider() string {
	return ""
}

func (m *mockAutoscaling) ProviderAPI() string {
	return ""
}

func (m *mockAutoscaling) ResourceTypes() []string {
	return []string{}
}

func (m *mockAutoscaling) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockAutoscaling) IsSyncDisabled() bool {
	return false
}

func (m *mockAutoscaling) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockAutoscaling) DescribeLaunchConfigurationsPages(input *autoscaling.DescribeLaunchConfigurationsInput, fn func(p *autoscaling.DescribeLaunchConfigurationsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*autoscaling.LaunchConfiguration
	for i := 0; i < len(m.launchconfigurations); i += 2 {
		page := []*autoscaling.LaunchConfiguration{m.launchconfigurations[i]}
		if i+1 < len(m.launchconfigurations) {
			page = append(page, m.launchconfigurations[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&autoscaling.DescribeLaunchConfigurationsOutput{LaunchConfigurations: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockAutoscaling) DescribeAutoScalingGroupsPages(input *autoscaling.DescribeAutoScalingGroupsInput, fn func(p *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*autoscaling.Group
	for i := 0; i < len(m.groups); i += 2 {
		page := []*autoscaling.Group{m.groups[i]}
		if i+1 < len(m.groups) {
			page = append(page, m.groups[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockAutoscaling) DescribePoliciesPages(input *autoscaling.DescribePoliciesInput, fn func(p *autoscaling.DescribePoliciesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*autoscaling.ScalingPolicy
	for i := 0; i < len(m.scalingpolicys); i += 2 {
		page := []*autoscaling.ScalingPolicy{m.scalingpolicys[i]}
		if i+1 < len(m.scalingpolicys) {
			page = append(page, m.scalingpolicys[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&autoscaling.DescribePoliciesOutput{ScalingPolicies: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockAcm struct {
	acmiface.ACMAPI
	certificatesummarys []*acm.CertificateSummary
}

func (m *mockAcm) Name() string {
	return ""
}

func (m *mockAcm) Region() string {
	return ""
}

func (m *mockAcm) Profile() string {
	return ""
}

func (m *mockAcm) Provider() string {
	return ""
}

func (m *mockAcm) ProviderAPI() string {
	return ""
}

func (m *mockAcm) ResourceTypes() []string {
	return []string{}
}

func (m *mockAcm) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockAcm) IsSyncDisabled() bool {
	return false
}

func (m *mockAcm) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockAcm) ListCertificatesPages(input *acm.ListCertificatesInput, fn func(p *acm.ListCertificatesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*acm.CertificateSummary
	for i := 0; i < len(m.certificatesummarys); i += 2 {
		page := []*acm.CertificateSummary{m.certificatesummarys[i]}
		if i+1 < len(m.certificatesummarys) {
			page = append(page, m.certificatesummarys[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&acm.ListCertificatesOutput{CertificateSummaryList: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockIam struct {
	iamiface.IAMAPI
	userdetails          []*iam.UserDetail
	groupdetails         []*iam.GroupDetail
	roledetails          []*iam.RoleDetail
	policys              []*iam.Policy
	accesskeymetadatas   []*iam.AccessKeyMetadata
	instanceprofiles     []*iam.InstanceProfile
	managedpolicydetails []*iam.ManagedPolicyDetail
	users                []*iam.User
	virtualmfadevices    []*iam.VirtualMFADevice
}

func (m *mockIam) Name() string {
	return ""
}

func (m *mockIam) Region() string {
	return ""
}

func (m *mockIam) Profile() string {
	return ""
}

func (m *mockIam) Provider() string {
	return ""
}

func (m *mockIam) ProviderAPI() string {
	return ""
}

func (m *mockIam) ResourceTypes() []string {
	return []string{}
}

func (m *mockIam) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockIam) IsSyncDisabled() bool {
	return false
}

func (m *mockIam) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockIam) ListAccessKeysPages(input *iam.ListAccessKeysInput, fn func(p *iam.ListAccessKeysOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*iam.AccessKeyMetadata
	for i := 0; i < len(m.accesskeymetadatas); i += 2 {
		page := []*iam.AccessKeyMetadata{m.accesskeymetadatas[i]}
		if i+1 < len(m.accesskeymetadatas) {
			page = append(page, m.accesskeymetadatas[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&iam.ListAccessKeysOutput{AccessKeyMetadata: page, Marker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockIam) ListInstanceProfilesPages(input *iam.ListInstanceProfilesInput, fn func(p *iam.ListInstanceProfilesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*iam.InstanceProfile
	for i := 0; i < len(m.instanceprofiles); i += 2 {
		page := []*iam.InstanceProfile{m.instanceprofiles[i]}
		if i+1 < len(m.instanceprofiles) {
			page = append(page, m.instanceprofiles[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&iam.ListInstanceProfilesOutput{InstanceProfiles: page, Marker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockIam) ListVirtualMFADevicesPages(input *iam.ListVirtualMFADevicesInput, fn func(p *iam.ListVirtualMFADevicesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*iam.VirtualMFADevice
	for i := 0; i < len(m.virtualmfadevices); i += 2 {
		page := []*iam.VirtualMFADevice{m.virtualmfadevices[i]}
		if i+1 < len(m.virtualmfadevices) {
			page = append(page, m.virtualmfadevices[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&iam.ListVirtualMFADevicesOutput{VirtualMFADevices: page, Marker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockS3 struct {
	s3iface.S3API
	buckets map[string][]*s3.Bucket
	objects map[string][]*s3.Object
	grants  map[string][]*s3.Grant
}

func (m *mockS3) Name() string {
	return ""
}

func (m *mockS3) Region() string {
	return ""
}

func (m *mockS3) Profile() string {
	return ""
}

func (m *mockS3) Provider() string {
	return ""
}

func (m *mockS3) ProviderAPI() string {
	return ""
}

func (m *mockS3) ResourceTypes() []string {
	return []string{}
}

func (m *mockS3) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockS3) IsSyncDisabled() bool {
	return false
}

func (m *mockS3) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

type mockSns struct {
	snsiface.SNSAPI
	subscriptions []*sns.Subscription
	topics        []*sns.Topic
}

func (m *mockSns) Name() string {
	return ""
}

func (m *mockSns) Region() string {
	return ""
}

func (m *mockSns) Profile() string {
	return ""
}

func (m *mockSns) Provider() string {
	return ""
}

func (m *mockSns) ProviderAPI() string {
	return ""
}

func (m *mockSns) ResourceTypes() []string {
	return []string{}
}

func (m *mockSns) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockSns) IsSyncDisabled() bool {
	return false
}

func (m *mockSns) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockSns) ListSubscriptionsPages(input *sns.ListSubscriptionsInput, fn func(p *sns.ListSubscriptionsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*sns.Subscription
	for i := 0; i < len(m.subscriptions); i += 2 {
		page := []*sns.Subscription{m.subscriptions[i]}
		if i+1 < len(m.subscriptions) {
			page = append(page, m.subscriptions[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&sns.ListSubscriptionsOutput{Subscriptions: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockSns) ListTopicsPages(input *sns.ListTopicsInput, fn func(p *sns.ListTopicsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*sns.Topic
	for i := 0; i < len(m.topics); i += 2 {
		page := []*sns.Topic{m.topics[i]}
		if i+1 < len(m.topics) {
			page = append(page, m.topics[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&sns.ListTopicsOutput{Topics: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockSqs struct {
	sqsiface.SQSAPI
	strings    []*string
	attributes map[string]map[string]*string
}

func (m *mockSqs) Name() string {
	return ""
}

func (m *mockSqs) Region() string {
	return ""
}

func (m *mockSqs) Profile() string {
	return ""
}

func (m *mockSqs) Provider() string {
	return ""
}

func (m *mockSqs) ProviderAPI() string {
	return ""
}

func (m *mockSqs) ResourceTypes() []string {
	return []string{}
}

func (m *mockSqs) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockSqs) IsSyncDisabled() bool {
	return false
}

func (m *mockSqs) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockSqs) ListQueues(input *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	return &sqs.ListQueuesOutput{QueueUrls: m.strings}, nil
}

type mockRoute53 struct {
	route53iface.Route53API
	hostedzones        []*route53.HostedZone
	resourcerecordsets map[string][]*route53.ResourceRecordSet
}

func (m *mockRoute53) Name() string {
	return ""
}

func (m *mockRoute53) Region() string {
	return ""
}

func (m *mockRoute53) Profile() string {
	return ""
}

func (m *mockRoute53) Provider() string {
	return ""
}

func (m *mockRoute53) ProviderAPI() string {
	return ""
}

func (m *mockRoute53) ResourceTypes() []string {
	return []string{}
}

func (m *mockRoute53) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockRoute53) IsSyncDisabled() bool {
	return false
}

func (m *mockRoute53) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockRoute53) ListHostedZonesPages(input *route53.ListHostedZonesInput, fn func(p *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*route53.HostedZone
	for i := 0; i < len(m.hostedzones); i += 2 {
		page := []*route53.HostedZone{m.hostedzones[i]}
		if i+1 < len(m.hostedzones) {
			page = append(page, m.hostedzones[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&route53.ListHostedZonesOutput{HostedZones: page, NextMarker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockLambda struct {
	lambdaiface.LambdaAPI
	functionconfigurations []*lambda.FunctionConfiguration
}

func (m *mockLambda) Name() string {
	return ""
}

func (m *mockLambda) Region() string {
	return ""
}

func (m *mockLambda) Profile() string {
	return ""
}

func (m *mockLambda) Provider() string {
	return ""
}

func (m *mockLambda) ProviderAPI() string {
	return ""
}

func (m *mockLambda) ResourceTypes() []string {
	return []string{}
}

func (m *mockLambda) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockLambda) IsSyncDisabled() bool {
	return false
}

func (m *mockLambda) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockLambda) ListFunctionsPages(input *lambda.ListFunctionsInput, fn func(p *lambda.ListFunctionsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*lambda.FunctionConfiguration
	for i := 0; i < len(m.functionconfigurations); i += 2 {
		page := []*lambda.FunctionConfiguration{m.functionconfigurations[i]}
		if i+1 < len(m.functionconfigurations) {
			page = append(page, m.functionconfigurations[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&lambda.ListFunctionsOutput{Functions: page, NextMarker: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockCloudwatch struct {
	cloudwatchiface.CloudWatchAPI
	metrics      []*cloudwatch.Metric
	metricalarms []*cloudwatch.MetricAlarm
}

func (m *mockCloudwatch) Name() string {
	return ""
}

func (m *mockCloudwatch) Region() string {
	return ""
}

func (m *mockCloudwatch) Profile() string {
	return ""
}

func (m *mockCloudwatch) Provider() string {
	return ""
}

func (m *mockCloudwatch) ProviderAPI() string {
	return ""
}

func (m *mockCloudwatch) ResourceTypes() []string {
	return []string{}
}

func (m *mockCloudwatch) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudwatch) IsSyncDisabled() bool {
	return false
}

func (m *mockCloudwatch) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudwatch) ListMetricsPages(input *cloudwatch.ListMetricsInput, fn func(p *cloudwatch.ListMetricsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*cloudwatch.Metric
	for i := 0; i < len(m.metrics); i += 2 {
		page := []*cloudwatch.Metric{m.metrics[i]}
		if i+1 < len(m.metrics) {
			page = append(page, m.metrics[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&cloudwatch.ListMetricsOutput{Metrics: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockCloudwatch) DescribeAlarmsPages(input *cloudwatch.DescribeAlarmsInput, fn func(p *cloudwatch.DescribeAlarmsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*cloudwatch.MetricAlarm
	for i := 0; i < len(m.metricalarms); i += 2 {
		page := []*cloudwatch.MetricAlarm{m.metricalarms[i]}
		if i+1 < len(m.metricalarms) {
			page = append(page, m.metricalarms[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&cloudwatch.DescribeAlarmsOutput{MetricAlarms: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockCloudfront struct {
	cloudfrontiface.CloudFrontAPI
	distributionsummarys []*cloudfront.DistributionSummary
}

func (m *mockCloudfront) Name() string {
	return ""
}

func (m *mockCloudfront) Region() string {
	return ""
}

func (m *mockCloudfront) Profile() string {
	return ""
}

func (m *mockCloudfront) Provider() string {
	return ""
}

func (m *mockCloudfront) ProviderAPI() string {
	return ""
}

func (m *mockCloudfront) ResourceTypes() []string {
	return []string{}
}

func (m *mockCloudfront) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudfront) IsSyncDisabled() bool {
	return false
}

func (m *mockCloudfront) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

type mockCloudformation struct {
	cloudformationiface.CloudFormationAPI
	stacks []*cloudformation.Stack
}

func (m *mockCloudformation) Name() string {
	return ""
}

func (m *mockCloudformation) Region() string {
	return ""
}

func (m *mockCloudformation) Profile() string {
	return ""
}

func (m *mockCloudformation) Provider() string {
	return ""
}

func (m *mockCloudformation) ProviderAPI() string {
	return ""
}

func (m *mockCloudformation) ResourceTypes() []string {
	return []string{}
}

func (m *mockCloudformation) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudformation) IsSyncDisabled() bool {
	return false
}

func (m *mockCloudformation) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudformation) DescribeStacksPages(input *cloudformation.DescribeStacksInput, fn func(p *cloudformation.DescribeStacksOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*cloudformation.Stack
	for i := 0; i < len(m.stacks); i += 2 {
		page := []*cloudformation.Stack{m.stacks[i]}
		if i+1 < len(m.stacks) {
			page = append(page, m.stacks[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&cloudformation.DescribeStacksOutput{Stacks: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockEcr struct {
	ecriface.ECRAPI
	repositorys []*ecr.Repository
}

func (m *mockEcr) Name() string {
	return ""
}

func (m *mockEcr) Region() string {
	return ""
}

func (m *mockEcr) Profile() string {
	return ""
}

func (m *mockEcr) Provider() string {
	return ""
}

func (m *mockEcr) ProviderAPI() string {
	return ""
}

func (m *mockEcr) ResourceTypes() []string {
	return []string{}
}

func (m *mockEcr) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEcr) IsSyncDisabled() bool {
	return false
}

func (m *mockEcr) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEcr) DescribeRepositoriesPages(input *ecr.DescribeRepositoriesInput, fn func(p *ecr.DescribeRepositoriesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*ecr.Repository
	for i := 0; i < len(m.repositorys); i += 2 {
		page := []*ecr.Repository{m.repositorys[i]}
		if i+1 < len(m.repositorys) {
			page = append(page, m.repositorys[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&ecr.DescribeRepositoriesOutput{Repositories: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

type mockEcs struct {
	ecsiface.ECSAPI
	clusters                []*ecs.Cluster
	clusterNames            []*string
	taskdefinitions         []*ecs.TaskDefinition
	taskdefinitionNames     []*string
	tasks                   map[string][]*ecs.Task
	tasksNames              map[string][]*string
	containerinstancesNames map[string][]*string
	containerinstances      map[string][]*ecs.ContainerInstance
}

func (m *mockEcs) Name() string {
	return ""
}

func (m *mockEcs) Region() string {
	return ""
}

func (m *mockEcs) Profile() string {
	return ""
}

func (m *mockEcs) Provider() string {
	return ""
}

func (m *mockEcs) ProviderAPI() string {
	return ""
}

func (m *mockEcs) ResourceTypes() []string {
	return []string{}
}

func (m *mockEcs) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEcs) IsSyncDisabled() bool {
	return false
}

func (m *mockEcs) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEcs) ListClustersPages(input *ecs.ListClustersInput, fn func(p *ecs.ListClustersOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*string
	for i := 0; i < len(m.clusterNames); i += 2 {
		page := []*string{m.clusterNames[i]}
		if i+1 < len(m.clusterNames) {
			page = append(page, m.clusterNames[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&ecs.ListClustersOutput{ClusterArns: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockEcs) ListTaskDefinitionsPages(input *ecs.ListTaskDefinitionsInput, fn func(p *ecs.ListTaskDefinitionsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*string
	for i := 0; i < len(m.taskdefinitionNames); i += 2 {
		page := []*string{m.taskdefinitionNames[i]}
		if i+1 < len(m.taskdefinitionNames) {
			page = append(page, m.taskdefinitionNames[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&ecs.ListTaskDefinitionsOutput{TaskDefinitionArns: page, NextToken: aws.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}
