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
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestRegionsValid(t *testing.T) {
	if got, want := stringInSlice("eu-west-1", AllRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("us-east-1", AllRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("us-west-1", AllRegions()), true; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	if got, want := stringInSlice("eu-test-1", AllRegions()), false; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	for _, k := range AllRegions() {
		if got, want := IsValidRegion(k), true; got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
	if got, want := IsValidRegion("aa-test-10"), false; got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

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
			GroupList:               []*string{awssdk.String("group_1"), awssdk.String("group_2")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}, {PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_2"),
			GroupList:      []*string{awssdk.String("group_1")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}},
		},
		{
			UserId:                  awssdk.String("usr_3"),
			GroupList:               []*string{awssdk.String("group_1"), awssdk.String("group_4")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}, {PolicyName: awssdk.String("nmanaged_policy_2")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}, {PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_4"),
			GroupList:      []*string{awssdk.String("group_2")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_5"),
			GroupList:      []*string{awssdk.String("group_2")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:                  awssdk.String("usr_6"),
			GroupList:               []*string{awssdk.String("group_2")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_3")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_7"),
			GroupList:      []*string{awssdk.String("group_2"), awssdk.String("group_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}, {PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_8"),
			GroupList:      []*string{awssdk.String("group_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_9"),
			GroupList:      []*string{awssdk.String("group_4")},
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
			PasswordLastUsed: awssdk.Time(time.Unix(1486139077, 0)),
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

	result := g.MustMarshal()
	expectContent, err := ioutil.ReadFile(filepath.Join("testdata", "access.rdf"))
	if err != nil {
		t.Fatal(err)
	}

	if err := diffText(result, string(expectContent)); err != nil {
		t.Fatal(err)
	}
}

func TestBuildInfraRdfGraph(t *testing.T) {
	instances := []*ec2.Instance{
		{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}}},
		{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{InstanceId: awssdk.String("inst_4"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}, {GroupId: awssdk.String("secgroup_2")}}, KeyName: awssdk.String("my_key_pair")},
		{InstanceId: awssdk.String("inst_5"), SubnetId: nil, VpcId: nil}, // terminated instance (no vpc, subnet ids)
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
		{RouteTableId: awssdk.String("rt_1"), VpcId: awssdk.String("vpc_1"), Associations: []*ec2.RouteTableAssociation{{RouteTableId: awssdk.String("rt_1"), SubnetId: awssdk.String("subnet_1")}}},
	}

	//ELB
	lbPages := [][]*elbv2.LoadBalancer{
		{{LoadBalancerArn: awssdk.String("lb_1"), LoadBalancerName: awssdk.String("my_loadbalancer"), VpcId: awssdk.String("vpc_1")}, {LoadBalancerArn: awssdk.String("lb_2"), VpcId: awssdk.String("vpc_2")}},
		{{LoadBalancerArn: awssdk.String("lb_3"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*string{awssdk.String("secgroup_1"), awssdk.String("secgroup_2")}}},
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
	infra := Infra{EC2API: mock, ELBV2API: mockLb, region: "eu-west-1"}
	InfraService = &infra

	g, err := infra.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()

	expectContent, err := ioutil.ReadFile(filepath.Join("testdata", "infra.rdf"))
	if err != nil {
		t.Fatal(err)
	}

	if err := diffText(result, string(expectContent)); err != nil {
		t.Fatal(err)
	}
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
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{URI: awssdk.String("usr_2")}},
		},
		"bucket_eu_1": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{URI: awssdk.String("usr_2")}},
		},
		"bucket_eu_2": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{URI: awssdk.String("usr_1")}},
		},
	}

	mocks3 := &mockS3{bucketsPerRegion: buckets, objectsPerBucket: objects, bucketsACL: bucketsACL}
	StorageService = mocks3
	storage := Storage{S3API: mocks3, region: "eu-west-1"}

	g, err := storage.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()

	expectContent, err := ioutil.ReadFile(filepath.Join("testdata", "storage.rdf"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := result, string(expectContent); got != want {
		t.Fatalf("got\n[%s]\n\nwant\n[%s]", got, want)
	}
}

func TestBuildEmptyRdfGraphWhenNoData(t *testing.T) {
	expect := `/region<eu-west-1>	"has_type"@[]	"/region"^^type:text`
	access := Access{IAMAPI: &mockIam{}, region: "eu-west-1"}

	g, err := access.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}

	infra := Infra{EC2API: &mockEc2{}, ELBV2API: &mockELB{}, region: "eu-west-1"}

	g, err = infra.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result = g.MustMarshal()
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
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
			return fmt.Errorf("text diff: diff at line %d\ngot:\n%s\nwant:\n%s", i+1, actuals[i], expecteds[i])
		}
	}

	return nil
}
