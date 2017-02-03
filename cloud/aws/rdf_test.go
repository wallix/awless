package aws

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

func TestBuildAccessRdfGraph(t *testing.T) {
	awsAccess := &AwsAccess{}

	awsAccess.GroupsDetail = []*iam.GroupDetail{
		&iam.GroupDetail{GroupId: awssdk.String("group_1"), GroupName: awssdk.String("ngroup_1")},
		&iam.GroupDetail{GroupId: awssdk.String("group_2"), GroupName: awssdk.String("ngroup_2")},
		&iam.GroupDetail{GroupId: awssdk.String("group_3"), GroupName: awssdk.String("ngroup_3")},
		&iam.GroupDetail{GroupId: awssdk.String("group_4"), GroupName: awssdk.String("ngroup_4")},
	}

	awsAccess.Policies = []*iam.ManagedPolicyDetail{
		&iam.ManagedPolicyDetail{PolicyId: awssdk.String("policy_1"), PolicyName: awssdk.String("npolicy_1")},
		&iam.ManagedPolicyDetail{PolicyId: awssdk.String("policy_2"), PolicyName: awssdk.String("npolicy_2")},
		&iam.ManagedPolicyDetail{PolicyId: awssdk.String("policy_3"), PolicyName: awssdk.String("npolicy_3")},
		&iam.ManagedPolicyDetail{PolicyId: awssdk.String("policy_4"), PolicyName: awssdk.String("npolicy_4")},
	}

	awsAccess.RolesDetail = []*iam.RoleDetail{
		&iam.RoleDetail{RoleId: awssdk.String("role_1")},
		&iam.RoleDetail{RoleId: awssdk.String("role_2")},
		&iam.RoleDetail{RoleId: awssdk.String("role_3")},
		&iam.RoleDetail{RoleId: awssdk.String("role_4")},
	}

	awsAccess.UsersDetail = []*iam.UserDetail{
		&iam.UserDetail{UserId: awssdk.String("usr_1")},
		&iam.UserDetail{UserId: awssdk.String("usr_2")},
		&iam.UserDetail{UserId: awssdk.String("usr_3")},
		&iam.UserDetail{UserId: awssdk.String("usr_4")},
		&iam.UserDetail{UserId: awssdk.String("usr_5")},
		&iam.UserDetail{UserId: awssdk.String("usr_6")},
		&iam.UserDetail{UserId: awssdk.String("usr_7")},
		&iam.UserDetail{UserId: awssdk.String("usr_8")},
		&iam.UserDetail{UserId: awssdk.String("usr_9")},
		&iam.UserDetail{UserId: awssdk.String("usr_10")}, //users not in any groups
		&iam.UserDetail{UserId: awssdk.String("usr_11")},
	}

	awsAccess.Users = []*iam.User{
		&iam.User{UserId: awssdk.String("usr_1")},
		&iam.User{UserId: awssdk.String("usr_2")},
		&iam.User{UserId: awssdk.String("usr_3")},
		&iam.User{UserId: awssdk.String("usr_4")},
		&iam.User{UserId: awssdk.String("usr_5")},
		&iam.User{UserId: awssdk.String("usr_6")},
		&iam.User{UserId: awssdk.String("usr_7")},
		&iam.User{UserId: awssdk.String("usr_8")},
		&iam.User{UserId: awssdk.String("usr_9")},
		&iam.User{UserId: awssdk.String("usr_10")}, //users not in any groups
		&iam.User{UserId: awssdk.String("usr_11")},
	}

	awsAccess.UserGroups = map[string][]string{
		"usr_1": []string{"group_1", "group_2"},
		"usr_2": []string{"group_1"},
		"usr_3": []string{"group_1", "group_4"},
		"usr_4": []string{"group_2"},
		"usr_5": []string{"group_2"},
		"usr_6": []string{"group_2"},
		"usr_7": []string{"group_2", "group_4"},
		"usr_8": []string{"group_4"},
		"usr_9": []string{"group_4"},
	}

	awsAccess.UserPolicies = map[string][]string{
		"usr_1": []string{"npolicy_1", "npolicy_2"},
		"usr_2": []string{"npolicy_1"},
		"usr_3": []string{"npolicy_1", "npolicy_4"},
		"usr_4": []string{"npolicy_2"},
		"usr_5": []string{"npolicy_2"},
		"usr_6": []string{"npolicy_2"},
		"usr_7": []string{"npolicy_2", "npolicy_4"},
		"usr_8": []string{"npolicy_4"},
		"usr_9": []string{"npolicy_4"},
	}

	awsAccess.RolePolicies = map[string][]string{
		"role_1": []string{"npolicy_1"},
		"role_2": []string{"npolicy_1"},
		"role_3": []string{"npolicy_2"},
		"role_4": []string{"npolicy_4"},
	}

	awsAccess.GroupPolicies = map[string][]string{
		"group_1": []string{"npolicy_1"},
		"group_2": []string{"npolicy_1"},
		"group_3": []string{"npolicy_2"},
		"group_4": []string{"npolicy_4"},
	}

	g, err := buildAccessGraph("eu-west-1", awsAccess)
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
	awsInfra := &AwsInfra{}

	awsInfra.instanceList = []*ec2.Instance{
		&ec2.Instance{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		&ec2.Instance{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}}},
		&ec2.Instance{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		&ec2.Instance{InstanceId: awssdk.String("inst_4"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}, {GroupId: awssdk.String("secgroup_2")}}},
		&ec2.Instance{InstanceId: awssdk.String("inst_5"), SubnetId: nil, VpcId: nil}, // terminated instance (no vpc, subnet ids)
	}

	awsInfra.vpcList = []*ec2.Vpc{
		&ec2.Vpc{VpcId: awssdk.String("vpc_1")},
		&ec2.Vpc{VpcId: awssdk.String("vpc_2")},
	}

	awsInfra.securitygroupList = []*ec2.SecurityGroup{
		&ec2.SecurityGroup{GroupId: awssdk.String("secgroup_1"), GroupName: awssdk.String("my_secgroup"), VpcId: awssdk.String("vpc_1")},
		&ec2.SecurityGroup{GroupId: awssdk.String("secgroup_2"), VpcId: awssdk.String("vpc_1")},
	}

	awsInfra.subnetList = []*ec2.Subnet{
		&ec2.Subnet{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_4"), VpcId: nil}, // edge case subnet with no vpc id
	}

	g, err := buildInfraGraph("eu-west-1", awsInfra)
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
	g, err := buildAccessGraph("eu-west-1", NewAwsAccess())
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}

	g, err = buildInfraGraph("eu-west-1", &AwsInfra{})
	if err != nil {
		t.Fatal(err)
	}

	result = g.MustMarshal()
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}
}

func BenchmarkBuildInfraRdfGraph(b *testing.B) {
	awsInfra := &AwsInfra{}

	for i := 0; i < 10; i++ {
		vpcId := fmt.Sprintf("vpc_%d", i+1)
		vpc := &ec2.Vpc{VpcId: awssdk.String(vpcId)}
		awsInfra.vpcList = append(awsInfra.vpcList, vpc)
		for j := 0; j < 10; j++ {
			subnetId := fmt.Sprintf("%s_sub_%d", vpcId, j+1)
			subnet := &ec2.Subnet{SubnetId: awssdk.String(subnetId), VpcId: awssdk.String(vpcId)}
			awsInfra.subnetList = append(awsInfra.subnetList, subnet)
			for k := 0; k < 1000; k++ {
				inst := &ec2.Instance{InstanceId: awssdk.String(fmt.Sprintf("%s_inst_%d", subnetId, k)), SubnetId: awssdk.String(subnetId), VpcId: awssdk.String(vpcId), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String(fmt.Sprintf("instance_%d_name", k))}}}
				awsInfra.instanceList = append(awsInfra.instanceList, inst)
			}
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buildInfraGraph("eu-west-1", awsInfra)
		if err != nil {
			b.Fatal(err)
		}
	}
}
