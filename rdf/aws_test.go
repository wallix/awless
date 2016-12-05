package rdf

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/api"
)

func TestBuildAccessRdfTriples(t *testing.T) {
	awsAccess := &api.AwsAccess{}

	awsAccess.Groups = []*iam.Group{
		&iam.Group{GroupId: aws.String("group_1"), GroupName: aws.String("ngroup_1")},
		&iam.Group{GroupId: aws.String("group_2"), GroupName: aws.String("ngroup_2")},
		&iam.Group{GroupId: aws.String("group_3"), GroupName: aws.String("ngroup_3")},
		&iam.Group{GroupId: aws.String("group_4"), GroupName: aws.String("ngroup_4")},
	}

	awsAccess.LocalPolicies = []*iam.Policy{
		&iam.Policy{PolicyId: aws.String("policy_1"), PolicyName: aws.String("npolicy_1")},
		&iam.Policy{PolicyId: aws.String("policy_2"), PolicyName: aws.String("npolicy_2")},
		&iam.Policy{PolicyId: aws.String("policy_3"), PolicyName: aws.String("npolicy_3")},
		&iam.Policy{PolicyId: aws.String("policy_4"), PolicyName: aws.String("npolicy_4")},
	}

	awsAccess.Roles = []*iam.Role{
		&iam.Role{RoleId: aws.String("role_1")},
		&iam.Role{RoleId: aws.String("role_2")},
		&iam.Role{RoleId: aws.String("role_3")},
		&iam.Role{RoleId: aws.String("role_4")},
	}

	awsAccess.Users = []*iam.User{
		&iam.User{UserId: aws.String("usr_1")},
		&iam.User{UserId: aws.String("usr_2")},
		&iam.User{UserId: aws.String("usr_3")},
		&iam.User{UserId: aws.String("usr_4")},
		&iam.User{UserId: aws.String("usr_5")},
		&iam.User{UserId: aws.String("usr_6")},
		&iam.User{UserId: aws.String("usr_7")},
		&iam.User{UserId: aws.String("usr_8")},
		&iam.User{UserId: aws.String("usr_9")},
		&iam.User{UserId: aws.String("usr_10")}, //users not in any groups
		&iam.User{UserId: aws.String("usr_11")},
	}

	awsAccess.UsersByGroup = map[string][]string{
		"group_1": []string{"usr_1", "usr_2", "usr_3"},
		"group_2": []string{"usr_1", "usr_4", "usr_5", "usr_6", "usr_7"},
		"group_4": []string{"usr_3", "usr_8", "usr_9", "usr_7"},
	}

	awsAccess.UsersByLocalPolicies = map[string][]string{
		"policy_1": []string{"usr_1", "usr_2", "usr_3"},
		"policy_2": []string{"usr_1", "usr_4", "usr_5", "usr_6", "usr_7"},
		"policy_4": []string{"usr_3", "usr_8", "usr_9", "usr_7"},
	}

	awsAccess.RolesByLocalPolicies = map[string][]string{
		"policy_1": []string{"role_1", "role_2"},
		"policy_2": []string{"role_3"},
		"policy_4": []string{"role_4"},
	}

	awsAccess.GroupsByLocalPolicies = map[string][]string{
		"policy_1": []string{"group_1", "group_2"},
		"policy_2": []string{"group_3"},
		"policy_4": []string{"group_4"},
	}

	triples, err := buildAccessRdfTriples("eu-west-1", awsAccess)
	if err != nil {
		t.Fatal(err)
	}

	result := marshalTriples(triples)
	expect := `/group<group_1>	"has_type"@[]	"/group"^^type:text
/group<group_1>	"parent_of"@[]	/user<usr_1>
/group<group_1>	"parent_of"@[]	/user<usr_2>
/group<group_1>	"parent_of"@[]	/user<usr_3>
/group<group_2>	"has_type"@[]	"/group"^^type:text
/group<group_2>	"parent_of"@[]	/user<usr_1>
/group<group_2>	"parent_of"@[]	/user<usr_4>
/group<group_2>	"parent_of"@[]	/user<usr_5>
/group<group_2>	"parent_of"@[]	/user<usr_6>
/group<group_2>	"parent_of"@[]	/user<usr_7>
/group<group_3>	"has_type"@[]	"/group"^^type:text
/group<group_4>	"has_type"@[]	"/group"^^type:text
/group<group_4>	"parent_of"@[]	/user<usr_3>
/group<group_4>	"parent_of"@[]	/user<usr_7>
/group<group_4>	"parent_of"@[]	/user<usr_8>
/group<group_4>	"parent_of"@[]	/user<usr_9>
/policy<policy_1>	"has_type"@[]	"/policy"^^type:text
/policy<policy_1>	"parent_of"@[]	/group<group_1>
/policy<policy_1>	"parent_of"@[]	/group<group_2>
/policy<policy_1>	"parent_of"@[]	/role<role_1>
/policy<policy_1>	"parent_of"@[]	/role<role_2>
/policy<policy_1>	"parent_of"@[]	/user<usr_1>
/policy<policy_1>	"parent_of"@[]	/user<usr_2>
/policy<policy_1>	"parent_of"@[]	/user<usr_3>
/policy<policy_2>	"has_type"@[]	"/policy"^^type:text
/policy<policy_2>	"parent_of"@[]	/group<group_3>
/policy<policy_2>	"parent_of"@[]	/role<role_3>
/policy<policy_2>	"parent_of"@[]	/user<usr_1>
/policy<policy_2>	"parent_of"@[]	/user<usr_4>
/policy<policy_2>	"parent_of"@[]	/user<usr_5>
/policy<policy_2>	"parent_of"@[]	/user<usr_6>
/policy<policy_2>	"parent_of"@[]	/user<usr_7>
/policy<policy_3>	"has_type"@[]	"/policy"^^type:text
/policy<policy_4>	"has_type"@[]	"/policy"^^type:text
/policy<policy_4>	"parent_of"@[]	/group<group_4>
/policy<policy_4>	"parent_of"@[]	/role<role_4>
/policy<policy_4>	"parent_of"@[]	/user<usr_3>
/policy<policy_4>	"parent_of"@[]	/user<usr_7>
/policy<policy_4>	"parent_of"@[]	/user<usr_8>
/policy<policy_4>	"parent_of"@[]	/user<usr_9>
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
/role<role_2>	"has_type"@[]	"/role"^^type:text
/role<role_3>	"has_type"@[]	"/role"^^type:text
/role<role_4>	"has_type"@[]	"/role"^^type:text
/user<usr_10>	"has_type"@[]	"/user"^^type:text
/user<usr_11>	"has_type"@[]	"/user"^^type:text
/user<usr_1>	"has_type"@[]	"/user"^^type:text
/user<usr_2>	"has_type"@[]	"/user"^^type:text
/user<usr_3>	"has_type"@[]	"/user"^^type:text
/user<usr_4>	"has_type"@[]	"/user"^^type:text
/user<usr_5>	"has_type"@[]	"/user"^^type:text
/user<usr_6>	"has_type"@[]	"/user"^^type:text
/user<usr_7>	"has_type"@[]	"/user"^^type:text
/user<usr_8>	"has_type"@[]	"/user"^^type:text
/user<usr_9>	"has_type"@[]	"/user"^^type:text`
	if result != expect {
		t.Fatalf("got\n[%s]\n\nwant\n[%s]", result, expect)
	}

}

func TestBuildInfraRdfTriples(t *testing.T) {
	awsInfra := &api.AwsInfra{}

	awsInfra.Instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: aws.String("inst_1"), SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_2"), SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_3"), SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2")},
		&ec2.Instance{InstanceId: aws.String("inst_4"), SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2")},
		&ec2.Instance{InstanceId: aws.String("inst_5"), SubnetId: nil, VpcId: nil}, // terminated instance (no vpc, subnet ids)
	}

	awsInfra.Vpcs = []*ec2.Vpc{
		&ec2.Vpc{VpcId: aws.String("vpc_1")},
		&ec2.Vpc{VpcId: aws.String("vpc_2")},
	}

	awsInfra.Subnets = []*ec2.Subnet{
		&ec2.Subnet{SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Subnet{SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Subnet{SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2")},
		&ec2.Subnet{SubnetId: aws.String("sub_4"), VpcId: nil}, // edge case subnet with no vpc id
	}

	infraTriples, propertiesTriples, err := buildInfraRdfTriples("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}

	result := marshalTriples(infraTriples)
	expect := `/instance<inst_1>	"has_type"@[]	"/instance"^^type:text
/instance<inst_2>	"has_type"@[]	"/instance"^^type:text
/instance<inst_3>	"has_type"@[]	"/instance"^^type:text
/instance<inst_4>	"has_type"@[]	"/instance"^^type:text
/instance<inst_5>	"has_type"@[]	"/instance"^^type:text
/region<eu-west-1>	"has_type"@[]	"/region"^^type:text
/region<eu-west-1>	"parent_of"@[]	/vpc<vpc_1>
/region<eu-west-1>	"parent_of"@[]	/vpc<vpc_2>
/subnet<sub_1>	"has_type"@[]	"/subnet"^^type:text
/subnet<sub_1>	"parent_of"@[]	/instance<inst_1>
/subnet<sub_2>	"has_type"@[]	"/subnet"^^type:text
/subnet<sub_2>	"parent_of"@[]	/instance<inst_2>
/subnet<sub_3>	"has_type"@[]	"/subnet"^^type:text
/subnet<sub_3>	"parent_of"@[]	/instance<inst_3>
/subnet<sub_3>	"parent_of"@[]	/instance<inst_4>
/subnet<sub_4>	"has_type"@[]	"/subnet"^^type:text
/vpc<vpc_1>	"has_type"@[]	"/vpc"^^type:text
/vpc<vpc_1>	"parent_of"@[]	/subnet<sub_1>
/vpc<vpc_1>	"parent_of"@[]	/subnet<sub_2>
/vpc<vpc_2>	"has_type"@[]	"/vpc"^^type:text
/vpc<vpc_2>	"parent_of"@[]	/subnet<sub_3>`
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}

	result = marshalTriples(propertiesTriples)
	expect = `/instance<inst_1>	"Id"@[]	"inst_1"^^type:text
/instance<inst_1>	"SubnetId"@[]	"sub_1"^^type:text
/instance<inst_1>	"VpcId"@[]	"vpc_1"^^type:text
/instance<inst_2>	"Id"@[]	"inst_2"^^type:text
/instance<inst_2>	"SubnetId"@[]	"sub_2"^^type:text
/instance<inst_2>	"VpcId"@[]	"vpc_1"^^type:text
/instance<inst_3>	"Id"@[]	"inst_3"^^type:text
/instance<inst_3>	"SubnetId"@[]	"sub_3"^^type:text
/instance<inst_3>	"VpcId"@[]	"vpc_2"^^type:text
/instance<inst_4>	"Id"@[]	"inst_4"^^type:text
/instance<inst_4>	"SubnetId"@[]	"sub_3"^^type:text
/instance<inst_4>	"VpcId"@[]	"vpc_2"^^type:text
/instance<inst_5>	"Id"@[]	"inst_5"^^type:text
/subnet<sub_1>	"Id"@[]	"sub_1"^^type:text
/subnet<sub_1>	"VpcId"@[]	"vpc_1"^^type:text
/subnet<sub_2>	"Id"@[]	"sub_2"^^type:text
/subnet<sub_2>	"VpcId"@[]	"vpc_1"^^type:text
/subnet<sub_3>	"Id"@[]	"sub_3"^^type:text
/subnet<sub_3>	"VpcId"@[]	"vpc_2"^^type:text
/subnet<sub_4>	"Id"@[]	"sub_4"^^type:text
/vpc<vpc_1>	"Id"@[]	"vpc_1"^^type:text
/vpc<vpc_2>	"Id"@[]	"vpc_2"^^type:text`
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}
}

func TestBuildEmptyRdfTriplesWhenNoData(t *testing.T) {
	expect := `/region<eu-west-1>	"has_type"@[]	"/region"^^type:text`
	triples, err := buildAccessRdfTriples("eu-west-1", api.NewAwsAccess())
	if err != nil {
		t.Fatal(err)
	}

	result := marshalTriples(triples)
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}

	infraTriples, propertiesTriples, err := buildInfraRdfTriples("eu-west-1", &api.AwsInfra{})
	if err != nil {
		t.Fatal(err)
	}

	result = marshalTriples(infraTriples)
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}

	result = marshalTriples(propertiesTriples)
	if got, want := len(propertiesTriples), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}
