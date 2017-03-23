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
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	p "github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestBuildAccessRdfGraph(t *testing.T) {
	managedPolicies := []*iam.ManagedPolicyDetail{
		{PolicyId: awssdk.String("managed_policy_1"), PolicyName: awssdk.String("nmanaged_policy_1")},
		{PolicyId: awssdk.String("managed_policy_2"), PolicyName: awssdk.String("nmanaged_policy_2")},
		{PolicyId: awssdk.String("managed_policy_3"), PolicyName: awssdk.String("nmanaged_policy_3")},
	}

	groups := []*iam.GroupDetail{
		{GroupId: awssdk.String("group_1"), GroupName: awssdk.String("ngroup_1"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}}},
		{GroupId: awssdk.String("group_2"), GroupName: awssdk.String("ngroup_2"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_2")}}},
		{GroupId: awssdk.String("group_3"), GroupName: awssdk.String("ngroup_3"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_3")}}},
		{GroupId: awssdk.String("group_4"), GroupName: awssdk.String("ngroup_4"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}}},
	}

	roles := []*iam.RoleDetail{
		{RoleId: awssdk.String("role_1"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}}},
		{RoleId: awssdk.String("role_2"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}},
		{RoleId: awssdk.String("role_3"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_2")}}},
		{RoleId: awssdk.String("role_4"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}}},
	}

	usersDetails := []*iam.UserDetail{
		{
			UserId:                  awssdk.String("usr_1"),
			GroupList:               []*string{awssdk.String("ngroup_1"), awssdk.String("ngroup_2")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}, {PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_2"),
			GroupList:      []*string{awssdk.String("ngroup_1")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}},
		},
		{
			UserId:                  awssdk.String("usr_3"),
			GroupList:               []*string{awssdk.String("ngroup_1"), awssdk.String("ngroup_4")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}, {PolicyName: awssdk.String("nmanaged_policy_2")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}, {PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_4"),
			GroupList:      []*string{awssdk.String("ngroup_2")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_5"),
			GroupList:      []*string{awssdk.String("ngroup_2")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:                  awssdk.String("usr_6"),
			GroupList:               []*string{awssdk.String("ngroup_2")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_3")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_7"),
			GroupList:      []*string{awssdk.String("ngroup_2"), awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}, {PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_8"),
			GroupList:      []*string{awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_9"),
			GroupList:      []*string{awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId: awssdk.String("usr_10"), //users not in any groups
		},
		{
			UserId: awssdk.String("usr_11"),
		},
	}

	users := []*iam.User{
		{
			UserId:           awssdk.String("usr_1"),
			PasswordLastUsed: awssdk.Time(time.Unix(1486139077, 0).UTC()),
		},
		{
			UserId: awssdk.String("usr_2"),
		},
		{
			UserId: awssdk.String("usr_3"),
		},
		{
			UserId: awssdk.String("usr_4"),
		},
		{
			UserId: awssdk.String("usr_5"),
		},
		{
			UserId: awssdk.String("usr_6"),
		},
		{
			UserId: awssdk.String("usr_7"),
		},
		{
			UserId: awssdk.String("usr_8"),
		},
		{
			UserId: awssdk.String("usr_9"),
		},
		{
			UserId: awssdk.String("usr_10"), //users not in any groups
		},
		{
			UserId: awssdk.String("usr_11"),
		},
	}

	mock := &mockIam{groups: groups, usersDetails: usersDetails, roles: roles, managedPolicies: managedPolicies, users: users}
	access := Access{IAMAPI: mock, region: "eu-west-1"}

	g, err := access.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.GetAllResources("policy", "group", "role", "user")
	if err != nil {
		t.Fatal(err)
	}

	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties[p.InlinePolicies].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]*graph.Resource{
		"managed_policy_1": resourcetest.Policy("managed_policy_1").Prop(p.Name, "nmanaged_policy_1").Build(),
		"managed_policy_2": resourcetest.Policy("managed_policy_2").Prop(p.Name, "nmanaged_policy_2").Build(),
		"managed_policy_3": resourcetest.Policy("managed_policy_3").Prop(p.Name, "nmanaged_policy_3").Build(),
		"group_1":          resourcetest.Group("group_1").Prop(p.Name, "ngroup_1").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"group_2":          resourcetest.Group("group_2").Prop(p.Name, "ngroup_2").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"group_3":          resourcetest.Group("group_3").Prop(p.Name, "ngroup_3").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"group_4":          resourcetest.Group("group_4").Prop(p.Name, "ngroup_4").Prop(p.InlinePolicies, []string{"npolicy_4"}).Build(),
		"role_1":           resourcetest.Role("role_1").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"role_2":           resourcetest.Role("role_2").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"role_3":           resourcetest.Role("role_3").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"role_4":           resourcetest.Role("role_4").Prop(p.InlinePolicies, []string{"npolicy_4"}).Build(),
		"usr_1":            resourcetest.User("usr_1").Prop(p.InlinePolicies, []string{"npolicy_1", "npolicy_2"}).Prop(p.PasswordLastUsed, time.Unix(1486139077, 0).UTC()).Build(),
		"usr_2":            resourcetest.User("usr_2").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"usr_3":            resourcetest.User("usr_3").Prop(p.InlinePolicies, []string{"npolicy_1", "npolicy_4"}).Build(),
		"usr_4":            resourcetest.User("usr_4").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"usr_5":            resourcetest.User("usr_5").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"usr_6":            resourcetest.User("usr_6").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"usr_7":            resourcetest.User("usr_7").Prop(p.InlinePolicies, []string{"npolicy_2", "npolicy_4"}).Build(),
		"usr_8":            resourcetest.User("usr_8").Prop(p.InlinePolicies, []string{"npolicy_4"}).Build(),
		"usr_9":            resourcetest.User("usr_9").Prop(p.InlinePolicies, []string{"npolicy_4"}).Build(),
		"usr_10":           resourcetest.User("usr_10").Build(),
		"usr_11":           resourcetest.User("usr_11").Build(),
	}

	expectedChildren := map[string][]string{}

	expectedAppliedOn := map[string][]string{
		"group_1":          {"usr_1", "usr_2", "usr_3"},
		"group_2":          {"usr_1", "usr_4", "usr_5", "usr_6", "usr_7"},
		"group_4":          {"usr_3", "usr_7", "usr_8", "usr_9"},
		"managed_policy_1": {"group_1", "role_1", "usr_1", "usr_3"},
		"managed_policy_2": {"group_2", "role_3", "usr_3"},
		"managed_policy_3": {"group_3", "usr_6"},
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildInfraRdfGraph(t *testing.T) {
	instances := []*ec2.Instance{
		{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}}},
		{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{InstanceId: awssdk.String("inst_4"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}, {GroupId: awssdk.String("secgroup_2")}}, KeyName: awssdk.String("my_key_pair")},
		{InstanceId: awssdk.String("inst_5"), SubnetId: nil, VpcId: nil, KeyName: awssdk.String("unexisting_keypair")}, // terminated instance (no vpc, subnet ids)
	}

	vpcs := []*ec2.Vpc{
		{VpcId: awssdk.String("vpc_1")},
		{VpcId: awssdk.String("vpc_2")},
	}

	securityGroups := []*ec2.SecurityGroup{
		{GroupId: awssdk.String("secgroup_1"), GroupName: awssdk.String("my_secgroup"), VpcId: awssdk.String("vpc_1")},
		{GroupId: awssdk.String("secgroup_2"), VpcId: awssdk.String("vpc_1")},
	}

	subnets := []*ec2.Subnet{
		{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
		{SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
		{SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{SubnetId: awssdk.String("sub_4"), VpcId: nil}, // edge case subnet with no vpc id
	}

	keypairs := []*ec2.KeyPairInfo{
		{KeyName: awssdk.String("my_key_pair")},
	}

	igws := []*ec2.InternetGateway{
		{InternetGatewayId: awssdk.String("igw_1"), Attachments: []*ec2.InternetGatewayAttachment{{VpcId: awssdk.String("vpc_2")}}},
	}

	routeTables := []*ec2.RouteTable{
		{RouteTableId: awssdk.String("rt_1"), VpcId: awssdk.String("vpc_1"), Associations: []*ec2.RouteTableAssociation{{RouteTableId: awssdk.String("rt_1"), SubnetId: awssdk.String("sub_1")}}},
	}

	//ELB
	lbPages := [][]*elbv2.LoadBalancer{
		{
			{LoadBalancerArn: awssdk.String("lb_1"), LoadBalancerName: awssdk.String("my_loadbalancer"), VpcId: awssdk.String("vpc_1")},
			{LoadBalancerArn: awssdk.String("lb_2"), VpcId: awssdk.String("vpc_2")},
		},
		{
			{LoadBalancerArn: awssdk.String("lb_3"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*string{awssdk.String("secgroup_1"), awssdk.String("secgroup_2")}},
		},
	}
	targetGroups := []*elbv2.TargetGroup{
		{TargetGroupArn: awssdk.String("tg_1"), VpcId: awssdk.String("vpc_1"), LoadBalancerArns: []*string{awssdk.String("lb_1"), awssdk.String("lb_3")}},
		{TargetGroupArn: awssdk.String("tg_2"), VpcId: awssdk.String("vpc_2"), LoadBalancerArns: []*string{awssdk.String("lb_2")}},
	}
	listeners := map[string][]*elbv2.Listener{
		"lb_1": {{ListenerArn: awssdk.String("list_1"), LoadBalancerArn: awssdk.String("lb_1")}, {ListenerArn: awssdk.String("list_1.2"), LoadBalancerArn: awssdk.String("lb_1")}},
		"lb_2": {{ListenerArn: awssdk.String("list_2"), LoadBalancerArn: awssdk.String("lb_2")}},
		"lb_3": {{ListenerArn: awssdk.String("list_3"), LoadBalancerArn: awssdk.String("lb_3")}},
	}
	targetHealths := map[string][]*elbv2.TargetHealthDescription{
		"tg_1": {{HealthCheckPort: awssdk.String("80"), Target: &elbv2.TargetDescription{Id: awssdk.String("inst_1"), Port: awssdk.Int64(443)}}},
		"tg_2": {{Target: &elbv2.TargetDescription{Id: awssdk.String("inst_2"), Port: awssdk.Int64(80)}}, {Target: &elbv2.TargetDescription{Id: awssdk.String("inst_3"), Port: awssdk.Int64(80)}}},
	}

	mock := &mockEc2{vpcs: vpcs, securityGroups: securityGroups, subnets: subnets, instances: instances, keyPairs: keypairs, internetGateways: igws, routeTables: routeTables}
	mockLb := &mockELB{loadBalancerPages: lbPages, targetGroups: targetGroups, listeners: listeners, targetHealths: targetHealths}
	infra := Infra{EC2API: mock, ELBV2API: mockLb, RDSAPI: &mockRDS{}, region: "eu-west-1"}
	InfraService = &infra

	g, err := infra.FetchResources()
	if err != nil {
		t.Fatal(err)
	}
	resources, err := g.GetAllResources("region", "instance", "vpc", "securitygroup", "subnet", "keypair", "internetgateway", "routetable", "loadbalancer", "targetgroup", "listener")
	if err != nil {
		t.Fatal(err)
	}

	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties[p.SecurityGroups].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties[p.Vpcs].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]*graph.Resource{
		"eu-west-1":   resourcetest.Region("eu-west-1").Build(),
		"inst_1":      resourcetest.Instance("inst_1").Prop(p.Subnet, "sub_1").Prop(p.Vpc, "vpc_1").Prop(p.Name, "instance1-name").Build(),
		"inst_2":      resourcetest.Instance("inst_2").Prop(p.Subnet, "sub_2").Prop(p.Vpc, "vpc_1").Prop(p.SecurityGroups, []string{"secgroup_1"}).Build(),
		"inst_3":      resourcetest.Instance("inst_3").Prop(p.Subnet, "sub_3").Prop(p.Vpc, "vpc_2").Build(),
		"inst_4":      resourcetest.Instance("inst_4").Prop(p.Subnet, "sub_3").Prop(p.Vpc, "vpc_2").Prop(p.SecurityGroups, []string{"secgroup_1", "secgroup_2"}).Prop(p.SSHKey, "my_key_pair").Build(),
		"inst_5":      resourcetest.Instance("inst_5").Prop(p.SSHKey, "unexisting_keypair").Build(),
		"vpc_1":       resourcetest.VPC("vpc_1").Build(),
		"vpc_2":       resourcetest.VPC("vpc_2").Build(),
		"secgroup_1":  resourcetest.SecGroup("secgroup_1").Prop(p.Name, "my_secgroup").Prop(p.Vpc, "vpc_1").Build(),
		"secgroup_2":  resourcetest.SecGroup("secgroup_2").Prop(p.Vpc, "vpc_1").Build(),
		"sub_1":       resourcetest.Subnet("sub_1").Prop(p.Vpc, "vpc_1").Build(),
		"sub_2":       resourcetest.Subnet("sub_2").Prop(p.Vpc, "vpc_1").Build(),
		"sub_3":       resourcetest.Subnet("sub_3").Prop(p.Vpc, "vpc_2").Build(),
		"sub_4":       resourcetest.Subnet("sub_4").Build(),
		"my_key_pair": resourcetest.Keypair("my_key_pair").Build(),
		"igw_1":       resourcetest.InternetGw("igw_1").Prop(p.Vpcs, []string{"vpc_2"}).Build(),
		"rt_1":        resourcetest.RouteTable("rt_1").Prop(p.Vpc, "vpc_1").Prop(p.Main, false).Build(),
		"lb_1":        resourcetest.LoadBalancer("lb_1").Prop(p.Name, "my_loadbalancer").Prop(p.Vpc, "vpc_1").Build(),
		"lb_2":        resourcetest.LoadBalancer("lb_2").Prop(p.Vpc, "vpc_2").Build(),
		"lb_3":        resourcetest.LoadBalancer("lb_3").Prop(p.Vpc, "vpc_1").Build(),
		"tg_1":        resourcetest.TargetGroup("tg_1").Prop(p.Vpc, "vpc_1").Build(),
		"tg_2":        resourcetest.TargetGroup("tg_2").Prop(p.Vpc, "vpc_2").Build(),
		"list_1":      resourcetest.Listener("list_1").Prop(p.LoadBalancer, "lb_1").Build(),
		"list_1.2":    resourcetest.Listener("list_1.2").Prop(p.LoadBalancer, "lb_1").Build(),
		"list_2":      resourcetest.Listener("list_2").Prop(p.LoadBalancer, "lb_2").Build(),
		"list_3":      resourcetest.Listener("list_3").Prop(p.LoadBalancer, "lb_3").Build(),
	}

	expectedChildren := map[string][]string{
		"eu-west-1": {"igw_1", "my_key_pair", "vpc_1", "vpc_2"},
		"lb_1":      {"list_1", "list_1.2"},
		"lb_2":      {"list_2"},
		"lb_3":      {"list_3"},
		"sub_1":     {"inst_1"},
		"sub_2":     {"inst_2"},
		"sub_3":     {"inst_3", "inst_4"},
		"vpc_1":     {"lb_1", "lb_3", "rt_1", "secgroup_1", "secgroup_2", "sub_1", "sub_2", "tg_1"},
		"vpc_2":     {"lb_2", "sub_3", "tg_2"},
	}

	expectedAppliedOn := map[string][]string{
		"igw_1":       {"vpc_2"},
		"lb_1":        {"tg_1"},
		"lb_2":        {"tg_2"},
		"lb_3":        {"tg_1"},
		"my_key_pair": {"inst_4"},
		"rt_1":        {"sub_1"},
		"secgroup_1":  {"inst_2", "inst_4", "lb_3"},
		"secgroup_2":  {"inst_4", "lb_3"},
		"tg_1":        {"inst_1"},
		"tg_2":        {"inst_2", "inst_3"},
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildStorageRdfGraph(t *testing.T) {
	buckets := map[string][]*s3.Bucket{
		"us-west-1": {
			{Name: awssdk.String("bucket_us_1")},
			{Name: awssdk.String("bucket_us_2")},
			{Name: awssdk.String("bucket_us_3")},
		},
		"eu-west-1": {
			{Name: awssdk.String("bucket_eu_1")},
			{Name: awssdk.String("bucket_eu_2")},
		},
	}
	objects := map[string][]*s3.Object{
		"bucket_us_1": {
			{Key: awssdk.String("obj_1")},
			{Key: awssdk.String("obj_2")},
		},
		"bucket_us_2": {},
		"bucket_us_3": {
			{Key: awssdk.String("obj_3")},
		},
		"bucket_eu_1": {
			{Key: awssdk.String("obj_4")},
		},
		"bucket_eu_2": {
			{Key: awssdk.String("obj_5")},
			{Key: awssdk.String("obj_6")},
		},
	}
	bucketsACL := map[string][]*s3.Grant{
		"bucket_us_1": {
			{Permission: awssdk.String("Read"), Grantee: &s3.Grantee{ID: awssdk.String("usr_1")}},
		},
		"bucket_us_3": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_2")}},
		},
		"bucket_eu_1": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_2")}},
		},
		"bucket_eu_2": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_1")}},
		},
	}

	mocks3 := &mockS3{bucketsPerRegion: buckets, objectsPerBucket: objects, bucketsACL: bucketsACL}
	StorageService = mocks3
	storage := Storage{S3API: mocks3, region: "eu-west-1"}

	g, err := storage.FetchResources()
	if err != nil {
		t.Fatal(err)
	}
	resources, err := g.GetAllResources("region", "bucket")
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]*graph.Resource{
		"eu-west-1":   resourcetest.Region("eu-west-1").Build(),
		"bucket_eu_1": resourcetest.Bucket("bucket_eu_1").Prop(p.Grants, []*graph.Grant{{GranteeID: "usr_2", Permission: "Write"}}).Build(),
		"bucket_eu_2": resourcetest.Bucket("bucket_eu_2").Prop(p.Grants, []*graph.Grant{{GranteeID: "usr_1", Permission: "Write"}}).Build(),
	}
	expectedChildren := map[string][]string{
		"eu-west-1":   {"bucket_eu_1", "bucket_eu_2"},
		"bucket_eu_1": {"obj_4"},
		"bucket_eu_2": {"obj_5", "obj_6"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildDnsRdfGraph(t *testing.T) {
	zonePages := [][]*route53.HostedZone{
		{
			{Id: awssdk.String("/hostedzone/12345"), Name: awssdk.String("my.first.domain")},
			{Id: awssdk.String("/hostedzone/23456"), Name: awssdk.String("my.second.domain")},
		},
		{{Id: awssdk.String("/hostedzone/34567"), Name: awssdk.String("my.third.domain")}},
	}
	recordPages := map[string][][]*route53.ResourceRecordSet{
		"/hostedzone/12345": {
			{
				{Type: awssdk.String("A"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain1.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("1.2.3.4")}, {Value: awssdk.String("2.3.4.5")}}},
				{Type: awssdk.String("A"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain2.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("3.4.5.6")}}},
			},
			{
				{Type: awssdk.String("CNAME"), TTL: awssdk.Int64(60), Name: awssdk.String("subdomain3.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("4.5.6.7")}}},
			},
		},
		"/hostedzone/23456": {
			{
				{Type: awssdk.String("A"), TTL: awssdk.Int64(30), Name: awssdk.String("subdomain1.my.second.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("5.6.7.8")}}},
				{Type: awssdk.String("CNAME"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain3.my.second.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("6.7.8.9")}}},
			},
		},
	}
	mockRoute53 := &mockRoute53{zonePages: zonePages, recordPages: recordPages}

	dns := Dns{Route53API: mockRoute53, region: "eu-west-1"}

	g, err := dns.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.GetAllResources("zone", "record")
	if err != nil {
		t.Fatal(err)
	}
	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties[p.Records].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]*graph.Resource{
		"/hostedzone/12345": resourcetest.Zone("/hostedzone/12345").Prop(p.Name, "my.first.domain").Build(),
		"/hostedzone/23456": resourcetest.Zone("/hostedzone/23456").Prop(p.Name, "my.second.domain").Build(),
		"/hostedzone/34567": resourcetest.Zone("/hostedzone/34567").Prop(p.Name, "my.third.domain").Build(),
		"awls-91fa0a45":     resourcetest.Record("awls-91fa0a45").Prop(p.Name, "subdomain1.my.first.domain").Prop(p.Type, "A").Prop(p.TTL, 10).Prop(p.Records, []string{"1.2.3.4", "2.3.4.5"}).Build(),
		"awls-920c0a46":     resourcetest.Record("awls-920c0a46").Prop(p.Name, "subdomain2.my.first.domain").Prop(p.Type, "A").Prop(p.TTL, 10).Prop(p.Records, []string{"3.4.5.6"}).Build(),
		"awls-be1e0b6a":     resourcetest.Record("awls-be1e0b6a").Prop(p.Name, "subdomain3.my.first.domain").Prop(p.Type, "CNAME").Prop(p.TTL, 60).Prop(p.Records, []string{"4.5.6.7"}).Build(),
		"awls-9c420a99":     resourcetest.Record("awls-9c420a99").Prop(p.Name, "subdomain1.my.second.domain").Prop(p.Type, "A").Prop(p.TTL, 30).Prop(p.Records, []string{"5.6.7.8"}).Build(),
		"awls-c9b80bbe":     resourcetest.Record("awls-c9b80bbe").Prop(p.Name, "subdomain3.my.second.domain").Prop(p.Type, "CNAME").Prop(p.TTL, 10).Prop(p.Records, []string{"6.7.8.9"}).Build(),
	}
	expectedChildren := map[string][]string{
		"/hostedzone/12345": {"awls-91fa0a45", "awls-920c0a46", "awls-be1e0b6a"},
		"/hostedzone/23456": {"awls-9c420a99", "awls-c9b80bbe"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildEmptyRdfGraphWhenNoData(t *testing.T) {

	expectG := graph.NewGraph()
	expectG.AddResource(resourcetest.Region("eu-west-1").Build())

	access := Access{IAMAPI: &mockIam{}, region: "eu-west-1"}

	g, err := access.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()
	if result != expectG.MustMarshal() {
		t.Fatalf("got [%s]\nwant [%s]", result, expectG.MustMarshal())
	}

	infra := Infra{EC2API: &mockEc2{}, ELBV2API: &mockELB{}, RDSAPI: &mockRDS{}, region: "eu-west-1"}

	g, err = infra.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result = g.MustMarshal()
	if result != expectG.MustMarshal() {
		t.Fatalf("got [%s]\nwant [%s]", result, expectG.MustMarshal())
	}
}

func diffText(actual, expected string) error {
	actuals := strings.Split(actual, "\n")
	expecteds := strings.Split(expected, "\n")

	if len(actuals) != len(expecteds) {
		return fmt.Errorf("text diff: not same number of lines:\ngot \n%s\n\nwant\n%s\n", actual, expected)
	}

	for i := 0; i < len(actuals); i++ {
		if actuals[i] != expecteds[i] {
			return fmt.Errorf("text diff: diff at line %d\ngot:\n%q\nwant:\n%q", i+1, actuals[i], expecteds[i])
		}
	}

	return nil
}

func mustGetChildrenId(g *graph.Graph, res *graph.Resource) []string {
	var collect []string

	err := g.Accept(&graph.ChildrenVisitor{From: res, IncludeFrom: false, Each: func(res *graph.Resource, depth int) error {
		if depth == 1 {
			collect = append(collect, res.Id())
		}
		return nil
	}})
	if err != nil {
		panic(err)
	}
	return collect
}

func mustGetAppliedOnId(g *graph.Graph, res *graph.Resource) []string {
	resources, err := g.ListResourcesAppliedOn(res)
	if err != nil {
		panic(err)
	}
	var ids []string
	for _, r := range resources {
		ids = append(ids, r.Id())
	}
	return ids
}

func compareResources(t *testing.T, g *graph.Graph, resources []*graph.Resource, expected map[string]*graph.Resource, expectedChildren, expectedAppliedOn map[string][]string) {
	if got, want := len(resources), len(expected); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	for _, got := range resources {
		want := expected[got.Id()]
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got \n%#v\nwant\n%#v", got, want)
		}
		children := mustGetChildrenId(g, got)
		sort.Strings(children)
		if g, w := children, expectedChildren[got.Id()]; !reflect.DeepEqual(g, w) {
			t.Fatalf("'%s' children: got %v, want %v", got.Id(), g, w)
		}
		appliedOn := mustGetAppliedOnId(g, got)
		sort.Strings(appliedOn)
		if g, w := appliedOn, expectedAppliedOn[got.Id()]; !reflect.DeepEqual(g, w) {
			t.Fatalf("'%s' appliedOn: got %v, want %v", got.Id(), g, w)
		}
	}
}
