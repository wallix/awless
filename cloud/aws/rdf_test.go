package aws

import (
	"fmt"
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

	g, err := BuildAwsAccessGraph("eu-west-1", awsAccess)
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()
	expect := `/group<group_1>	"has_type"@[]	"/group"^^type:text
/group<group_1>	"parent_of"@[]	/user<usr_1>
/group<group_1>	"parent_of"@[]	/user<usr_2>
/group<group_1>	"parent_of"@[]	/user<usr_3>
/group<group_1>	"property"@[]	"{"Key":"Id","Value":"group_1"}"^^type:text
/group<group_1>	"property"@[]	"{"Key":"Name","Value":"ngroup_1"}"^^type:text
/group<group_2>	"has_type"@[]	"/group"^^type:text
/group<group_2>	"parent_of"@[]	/user<usr_1>
/group<group_2>	"parent_of"@[]	/user<usr_4>
/group<group_2>	"parent_of"@[]	/user<usr_5>
/group<group_2>	"parent_of"@[]	/user<usr_6>
/group<group_2>	"parent_of"@[]	/user<usr_7>
/group<group_2>	"property"@[]	"{"Key":"Id","Value":"group_2"}"^^type:text
/group<group_2>	"property"@[]	"{"Key":"Name","Value":"ngroup_2"}"^^type:text
/group<group_3>	"has_type"@[]	"/group"^^type:text
/group<group_3>	"property"@[]	"{"Key":"Id","Value":"group_3"}"^^type:text
/group<group_3>	"property"@[]	"{"Key":"Name","Value":"ngroup_3"}"^^type:text
/group<group_4>	"has_type"@[]	"/group"^^type:text
/group<group_4>	"parent_of"@[]	/user<usr_3>
/group<group_4>	"parent_of"@[]	/user<usr_7>
/group<group_4>	"parent_of"@[]	/user<usr_8>
/group<group_4>	"parent_of"@[]	/user<usr_9>
/group<group_4>	"property"@[]	"{"Key":"Id","Value":"group_4"}"^^type:text
/group<group_4>	"property"@[]	"{"Key":"Name","Value":"ngroup_4"}"^^type:text
/policy<policy_1>	"has_type"@[]	"/policy"^^type:text
/policy<policy_1>	"parent_of"@[]	/group<group_1>
/policy<policy_1>	"parent_of"@[]	/group<group_2>
/policy<policy_1>	"parent_of"@[]	/role<role_1>
/policy<policy_1>	"parent_of"@[]	/role<role_2>
/policy<policy_1>	"parent_of"@[]	/user<usr_1>
/policy<policy_1>	"parent_of"@[]	/user<usr_2>
/policy<policy_1>	"parent_of"@[]	/user<usr_3>
/policy<policy_1>	"property"@[]	"{"Key":"Id","Value":"policy_1"}"^^type:text
/policy<policy_1>	"property"@[]	"{"Key":"Name","Value":"npolicy_1"}"^^type:text
/policy<policy_2>	"has_type"@[]	"/policy"^^type:text
/policy<policy_2>	"parent_of"@[]	/group<group_3>
/policy<policy_2>	"parent_of"@[]	/role<role_3>
/policy<policy_2>	"parent_of"@[]	/user<usr_1>
/policy<policy_2>	"parent_of"@[]	/user<usr_4>
/policy<policy_2>	"parent_of"@[]	/user<usr_5>
/policy<policy_2>	"parent_of"@[]	/user<usr_6>
/policy<policy_2>	"parent_of"@[]	/user<usr_7>
/policy<policy_2>	"property"@[]	"{"Key":"Id","Value":"policy_2"}"^^type:text
/policy<policy_2>	"property"@[]	"{"Key":"Name","Value":"npolicy_2"}"^^type:text
/policy<policy_3>	"has_type"@[]	"/policy"^^type:text
/policy<policy_3>	"property"@[]	"{"Key":"Id","Value":"policy_3"}"^^type:text
/policy<policy_3>	"property"@[]	"{"Key":"Name","Value":"npolicy_3"}"^^type:text
/policy<policy_4>	"has_type"@[]	"/policy"^^type:text
/policy<policy_4>	"parent_of"@[]	/group<group_4>
/policy<policy_4>	"parent_of"@[]	/role<role_4>
/policy<policy_4>	"parent_of"@[]	/user<usr_3>
/policy<policy_4>	"parent_of"@[]	/user<usr_7>
/policy<policy_4>	"parent_of"@[]	/user<usr_8>
/policy<policy_4>	"parent_of"@[]	/user<usr_9>
/policy<policy_4>	"property"@[]	"{"Key":"Id","Value":"policy_4"}"^^type:text
/policy<policy_4>	"property"@[]	"{"Key":"Name","Value":"npolicy_4"}"^^type:text
/region<eu-west-1>	"has_type"@[]	"/region"^^type:text
/region<eu-west-1>	"parent_of"@[]	/group<group_1>
/region<eu-west-1>	"parent_of"@[]	/group<group_2>
/region<eu-west-1>	"parent_of"@[]	/group<group_3>
/region<eu-west-1>	"parent_of"@[]	/group<group_4>
/region<eu-west-1>	"parent_of"@[]	/policy<policy_1>
/region<eu-west-1>	"parent_of"@[]	/policy<policy_2>
/region<eu-west-1>	"parent_of"@[]	/policy<policy_3>
/region<eu-west-1>	"parent_of"@[]	/policy<policy_4>
/region<eu-west-1>	"parent_of"@[]	/role<role_1>
/region<eu-west-1>	"parent_of"@[]	/role<role_2>
/region<eu-west-1>	"parent_of"@[]	/role<role_3>
/region<eu-west-1>	"parent_of"@[]	/role<role_4>
/region<eu-west-1>	"parent_of"@[]	/user<usr_10>
/region<eu-west-1>	"parent_of"@[]	/user<usr_11>
/region<eu-west-1>	"parent_of"@[]	/user<usr_1>
/region<eu-west-1>	"parent_of"@[]	/user<usr_2>
/region<eu-west-1>	"parent_of"@[]	/user<usr_3>
/region<eu-west-1>	"parent_of"@[]	/user<usr_4>
/region<eu-west-1>	"parent_of"@[]	/user<usr_5>
/region<eu-west-1>	"parent_of"@[]	/user<usr_6>
/region<eu-west-1>	"parent_of"@[]	/user<usr_7>
/region<eu-west-1>	"parent_of"@[]	/user<usr_8>
/region<eu-west-1>	"parent_of"@[]	/user<usr_9>
/role<role_1>	"has_type"@[]	"/role"^^type:text
/role<role_1>	"property"@[]	"{"Key":"Id","Value":"role_1"}"^^type:text
/role<role_2>	"has_type"@[]	"/role"^^type:text
/role<role_2>	"property"@[]	"{"Key":"Id","Value":"role_2"}"^^type:text
/role<role_3>	"has_type"@[]	"/role"^^type:text
/role<role_3>	"property"@[]	"{"Key":"Id","Value":"role_3"}"^^type:text
/role<role_4>	"has_type"@[]	"/role"^^type:text
/role<role_4>	"property"@[]	"{"Key":"Id","Value":"role_4"}"^^type:text
/user<usr_10>	"has_type"@[]	"/user"^^type:text
/user<usr_10>	"property"@[]	"{"Key":"Id","Value":"usr_10"}"^^type:text
/user<usr_11>	"has_type"@[]	"/user"^^type:text
/user<usr_11>	"property"@[]	"{"Key":"Id","Value":"usr_11"}"^^type:text
/user<usr_1>	"has_type"@[]	"/user"^^type:text
/user<usr_1>	"property"@[]	"{"Key":"Id","Value":"usr_1"}"^^type:text
/user<usr_2>	"has_type"@[]	"/user"^^type:text
/user<usr_2>	"property"@[]	"{"Key":"Id","Value":"usr_2"}"^^type:text
/user<usr_3>	"has_type"@[]	"/user"^^type:text
/user<usr_3>	"property"@[]	"{"Key":"Id","Value":"usr_3"}"^^type:text
/user<usr_4>	"has_type"@[]	"/user"^^type:text
/user<usr_4>	"property"@[]	"{"Key":"Id","Value":"usr_4"}"^^type:text
/user<usr_5>	"has_type"@[]	"/user"^^type:text
/user<usr_5>	"property"@[]	"{"Key":"Id","Value":"usr_5"}"^^type:text
/user<usr_6>	"has_type"@[]	"/user"^^type:text
/user<usr_6>	"property"@[]	"{"Key":"Id","Value":"usr_6"}"^^type:text
/user<usr_7>	"has_type"@[]	"/user"^^type:text
/user<usr_7>	"property"@[]	"{"Key":"Id","Value":"usr_7"}"^^type:text
/user<usr_8>	"has_type"@[]	"/user"^^type:text
/user<usr_8>	"property"@[]	"{"Key":"Id","Value":"usr_8"}"^^type:text
/user<usr_9>	"has_type"@[]	"/user"^^type:text
/user<usr_9>	"property"@[]	"{"Key":"Id","Value":"usr_9"}"^^type:text`
	if result != expect {
		t.Fatalf("got\n[%s]\n\nwant\n[%s]", result, expect)
	}

}

func TestBuildInfraRdfGraph(t *testing.T) {
	awsInfra := &AwsInfra{}

	awsInfra.Instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		&ec2.Instance{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}}},
		&ec2.Instance{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		&ec2.Instance{InstanceId: awssdk.String("inst_4"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}, {GroupId: awssdk.String("secgroup_2")}}},
		&ec2.Instance{InstanceId: awssdk.String("inst_5"), SubnetId: nil, VpcId: nil}, // terminated instance (no vpc, subnet ids)
	}

	awsInfra.Vpcs = []*ec2.Vpc{
		&ec2.Vpc{VpcId: awssdk.String("vpc_1")},
		&ec2.Vpc{VpcId: awssdk.String("vpc_2")},
	}

	awsInfra.SecurityGroups = []*ec2.SecurityGroup{
		&ec2.SecurityGroup{GroupId: awssdk.String("secgroup_1"), GroupName: awssdk.String("my_secgroup"), VpcId: awssdk.String("vpc_1")},
		&ec2.SecurityGroup{GroupId: awssdk.String("secgroup_2"), VpcId: awssdk.String("vpc_1")},
	}

	awsInfra.Subnets = []*ec2.Subnet{
		&ec2.Subnet{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_4"), VpcId: nil}, // edge case subnet with no vpc id
	}

	g, err := BuildAwsInfraGraph("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()
	expect := `/instance<inst_1>	"has_type"@[]	"/instance"^^type:text
/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"inst_1"}"^^type:text
/instance<inst_1>	"property"@[]	"{"Key":"Name","Value":"instance1-name"}"^^type:text
/instance<inst_1>	"property"@[]	"{"Key":"SubnetId","Value":"sub_1"}"^^type:text
/instance<inst_1>	"property"@[]	"{"Key":"VpcId","Value":"vpc_1"}"^^type:text
/instance<inst_2>	"has_type"@[]	"/instance"^^type:text
/instance<inst_2>	"property"@[]	"{"Key":"Id","Value":"inst_2"}"^^type:text
/instance<inst_2>	"property"@[]	"{"Key":"SubnetId","Value":"sub_2"}"^^type:text
/instance<inst_2>	"property"@[]	"{"Key":"VpcId","Value":"vpc_1"}"^^type:text
/instance<inst_3>	"has_type"@[]	"/instance"^^type:text
/instance<inst_3>	"property"@[]	"{"Key":"Id","Value":"inst_3"}"^^type:text
/instance<inst_3>	"property"@[]	"{"Key":"SubnetId","Value":"sub_3"}"^^type:text
/instance<inst_3>	"property"@[]	"{"Key":"VpcId","Value":"vpc_2"}"^^type:text
/instance<inst_4>	"has_type"@[]	"/instance"^^type:text
/instance<inst_4>	"property"@[]	"{"Key":"Id","Value":"inst_4"}"^^type:text
/instance<inst_4>	"property"@[]	"{"Key":"SubnetId","Value":"sub_3"}"^^type:text
/instance<inst_4>	"property"@[]	"{"Key":"VpcId","Value":"vpc_2"}"^^type:text
/instance<inst_5>	"has_type"@[]	"/instance"^^type:text
/instance<inst_5>	"property"@[]	"{"Key":"Id","Value":"inst_5"}"^^type:text
/region<eu-west-1>	"has_type"@[]	"/region"^^type:text
/region<eu-west-1>	"parent_of"@[]	/vpc<vpc_1>
/region<eu-west-1>	"parent_of"@[]	/vpc<vpc_2>
/securitygroup<secgroup_1>	"has_type"@[]	"/securitygroup"^^type:text
/securitygroup<secgroup_1>	"parent_of"@[]	/instance<inst_2>
/securitygroup<secgroup_1>	"parent_of"@[]	/instance<inst_4>
/securitygroup<secgroup_1>	"property"@[]	"{"Key":"Id","Value":"secgroup_1"}"^^type:text
/securitygroup<secgroup_1>	"property"@[]	"{"Key":"Name","Value":"my_secgroup"}"^^type:text
/securitygroup<secgroup_1>	"property"@[]	"{"Key":"VpcId","Value":"vpc_1"}"^^type:text
/securitygroup<secgroup_2>	"has_type"@[]	"/securitygroup"^^type:text
/securitygroup<secgroup_2>	"parent_of"@[]	/instance<inst_4>
/securitygroup<secgroup_2>	"property"@[]	"{"Key":"Id","Value":"secgroup_2"}"^^type:text
/securitygroup<secgroup_2>	"property"@[]	"{"Key":"VpcId","Value":"vpc_1"}"^^type:text
/subnet<sub_1>	"has_type"@[]	"/subnet"^^type:text
/subnet<sub_1>	"parent_of"@[]	/instance<inst_1>
/subnet<sub_1>	"property"@[]	"{"Key":"Id","Value":"sub_1"}"^^type:text
/subnet<sub_1>	"property"@[]	"{"Key":"VpcId","Value":"vpc_1"}"^^type:text
/subnet<sub_2>	"has_type"@[]	"/subnet"^^type:text
/subnet<sub_2>	"parent_of"@[]	/instance<inst_2>
/subnet<sub_2>	"property"@[]	"{"Key":"Id","Value":"sub_2"}"^^type:text
/subnet<sub_2>	"property"@[]	"{"Key":"VpcId","Value":"vpc_1"}"^^type:text
/subnet<sub_3>	"has_type"@[]	"/subnet"^^type:text
/subnet<sub_3>	"parent_of"@[]	/instance<inst_3>
/subnet<sub_3>	"parent_of"@[]	/instance<inst_4>
/subnet<sub_3>	"property"@[]	"{"Key":"Id","Value":"sub_3"}"^^type:text
/subnet<sub_3>	"property"@[]	"{"Key":"VpcId","Value":"vpc_2"}"^^type:text
/subnet<sub_4>	"has_type"@[]	"/subnet"^^type:text
/subnet<sub_4>	"property"@[]	"{"Key":"Id","Value":"sub_4"}"^^type:text
/vpc<vpc_1>	"has_type"@[]	"/vpc"^^type:text
/vpc<vpc_1>	"parent_of"@[]	/securitygroup<secgroup_1>
/vpc<vpc_1>	"parent_of"@[]	/securitygroup<secgroup_2>
/vpc<vpc_1>	"parent_of"@[]	/subnet<sub_1>
/vpc<vpc_1>	"parent_of"@[]	/subnet<sub_2>
/vpc<vpc_1>	"property"@[]	"{"Key":"Id","Value":"vpc_1"}"^^type:text
/vpc<vpc_2>	"has_type"@[]	"/vpc"^^type:text
/vpc<vpc_2>	"parent_of"@[]	/subnet<sub_3>
/vpc<vpc_2>	"property"@[]	"{"Key":"Id","Value":"vpc_2"}"^^type:text`
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}
}

func TestBuildEmptyRdfGraphWhenNoData(t *testing.T) {
	expect := `/region<eu-west-1>	"has_type"@[]	"/region"^^type:text`
	g, err := BuildAwsAccessGraph("eu-west-1", NewAwsAccess())
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}

	g, err = BuildAwsInfraGraph("eu-west-1", &AwsInfra{})
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
		awsInfra.Vpcs = append(awsInfra.Vpcs, vpc)
		for j := 0; j < 10; j++ {
			subnetId := fmt.Sprintf("%s_sub_%d", vpcId, j+1)
			subnet := &ec2.Subnet{SubnetId: awssdk.String(subnetId), VpcId: awssdk.String(vpcId)}
			awsInfra.Subnets = append(awsInfra.Subnets, subnet)
			for k := 0; k < 1000; k++ {
				inst := &ec2.Instance{InstanceId: awssdk.String(fmt.Sprintf("%s_inst_%d", subnetId, k)), SubnetId: awssdk.String(subnetId), VpcId: awssdk.String(vpcId), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String(fmt.Sprintf("instance_%d_name", k))}}}
				awsInfra.Instances = append(awsInfra.Instances, inst)
			}
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := BuildAwsInfraGraph("eu-west-1", awsInfra)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestInstanceCredentialsFromName(t *testing.T) {
	awsInfra := &AwsInfra{}

	awsInfra.Instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: awssdk.String("inst_1"),
			SubnetId:        awssdk.String("sub_1"),
			VpcId:           awssdk.String("vpc_1"),
			Tags:            []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}},
			PublicIpAddress: awssdk.String("1.2.3.4"),
			KeyName:         awssdk.String("my-key-name"),
		},
		&ec2.Instance{
			InstanceId:      awssdk.String("inst_2"),
			SubnetId:        awssdk.String("sub_1"),
			VpcId:           awssdk.String("vpc_1"),
			PublicIpAddress: awssdk.String("2.3.4.5"),
		},
		&ec2.Instance{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
	}

	awsInfra.Vpcs = []*ec2.Vpc{
		&ec2.Vpc{VpcId: awssdk.String("vpc_1")},
	}

	awsInfra.Subnets = []*ec2.Subnet{
		&ec2.Subnet{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
	}

	g, err := BuildAwsInfraGraph("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}
	cred, err := InstanceCredentialsFromGraph(g, "inst_1")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := cred.IP, "1.2.3.4"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cred.KeyName, "my-key-name"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cred.User, ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	_, err = InstanceCredentialsFromGraph(g, "inst_12")
	if got, want := err, ErrInstanceNotFound; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	_, err = InstanceCredentialsFromGraph(g, "inst_3")
	if got, want := err, ErrNoPublicIP; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	_, err = InstanceCredentialsFromGraph(g, "inst_2")
	if got, want := err, ErrNoAccessKey; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
