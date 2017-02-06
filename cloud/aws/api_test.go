package aws

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
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

	if got, want := result, string(expectContent); got != want {
		t.Fatalf("got\n[%s]\n\nwant\n[%s]", got, want)
	}
}

func TestBuildInfraRdfGraph(t *testing.T) {
	instances := []*ec2.Instance{
		{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}}},
		{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{InstanceId: awssdk.String("inst_4"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}, {GroupId: awssdk.String("secgroup_2")}}},
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

	mock := &mockEc2{vpcs: vpcs, securityGroups: securityGroups, subnets: subnets, instances: instances, keyPairs: keypairs, internetGateways: igws, routeTables: routeTables}
	infra := Infra{EC2API: mock, region: "eu-west-1"}

	g, err := infra.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()
	expectContent, err := ioutil.ReadFile(filepath.Join("testdata", "infra.rdf"))
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

	infra := Infra{EC2API: &mockEc2{}, region: "eu-west-1"}

	g, err = infra.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result = g.MustMarshal()
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}
}

type mockEc2 struct {
	ec2iface.EC2API
	vpcs             []*ec2.Vpc
	subnets          []*ec2.Subnet
	instances        []*ec2.Instance
	securityGroups   []*ec2.SecurityGroup
	keyPairs         []*ec2.KeyPairInfo
	internetGateways []*ec2.InternetGateway
	routeTables      []*ec2.RouteTable
}

func (m *mockEc2) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return &ec2.DescribeVpcsOutput{Vpcs: m.vpcs}, nil
}

func (m *mockEc2) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return &ec2.DescribeSubnetsOutput{Subnets: m.subnets}, nil
}

func (m *mockEc2) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{{Instances: m.instances}}}, nil
}

func (m *mockEc2) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return &ec2.DescribeSecurityGroupsOutput{SecurityGroups: m.securityGroups}, nil
}

func (m *mockEc2) DescribeKeyPairs(input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return &ec2.DescribeKeyPairsOutput{KeyPairs: m.keyPairs}, nil
}

func (m *mockEc2) DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	return &ec2.DescribeInternetGatewaysOutput{InternetGateways: m.internetGateways}, nil
}

func (m *mockEc2) DescribeRouteTables(input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	return &ec2.DescribeRouteTablesOutput{RouteTables: m.routeTables}, nil
}

// Not tested
func (m *mockEc2) DescribeVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	return &ec2.DescribeVolumesOutput{}, nil
}

type mockIam struct {
	iamiface.IAMAPI
	groups          []*iam.GroupDetail
	managedPolicies []*iam.ManagedPolicyDetail
	policies        []*iam.Policy
	roles           []*iam.RoleDetail
	users           []*iam.User
	usersDetails    []*iam.UserDetail
}

func (m *mockIam) ListUsers(input *iam.ListUsersInput) (*iam.ListUsersOutput, error) {
	return &iam.ListUsersOutput{Users: m.users}, nil
}

func (m *mockIam) ListPolicies(input *iam.ListPoliciesInput) (*iam.ListPoliciesOutput, error) {
	var policies []*iam.Policy
	for _, p := range m.managedPolicies {
		policy := &iam.Policy{PolicyId: p.PolicyId, PolicyName: p.PolicyName}
		policies = append(policies, policy)
	}
	return &iam.ListPoliciesOutput{Policies: policies}, nil
}

func (m *mockIam) GetAccountAuthorizationDetails(input *iam.GetAccountAuthorizationDetailsInput) (*iam.GetAccountAuthorizationDetailsOutput, error) {
	return &iam.GetAccountAuthorizationDetailsOutput{GroupDetailList: m.groups, Policies: m.managedPolicies, RoleDetailList: m.roles, UserDetailList: m.usersDetails}, nil
}

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
