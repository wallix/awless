package rdf

import (
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
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

	triples, err := BuildAccessRdfTriples("eu-west-1", awsAccess)
	if err != nil {
		t.Fatal(err)
	}

	result := SortLines(MarshalTriples(triples))
	expect := `/group<group_1>	"parent_of"@[]	/user<usr_1>
/group<group_1>	"parent_of"@[]	/user<usr_2>
/group<group_1>	"parent_of"@[]	/user<usr_3>
/group<group_2>	"parent_of"@[]	/user<usr_1>
/group<group_2>	"parent_of"@[]	/user<usr_4>
/group<group_2>	"parent_of"@[]	/user<usr_5>
/group<group_2>	"parent_of"@[]	/user<usr_6>
/group<group_2>	"parent_of"@[]	/user<usr_7>
/group<group_4>	"parent_of"@[]	/user<usr_3>
/group<group_4>	"parent_of"@[]	/user<usr_7>
/group<group_4>	"parent_of"@[]	/user<usr_8>
/group<group_4>	"parent_of"@[]	/user<usr_9>
/region<eu-west-1>	"parent_of"@[]	/group<group_1>
/region<eu-west-1>	"parent_of"@[]	/group<group_2>
/region<eu-west-1>	"parent_of"@[]	/group<group_3>
/region<eu-west-1>	"parent_of"@[]	/group<group_4>
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
/region<eu-west-1>	"parent_of"@[]	/user<usr_9>`
	if result != expect {
		t.Fatalf("got\n%s\n\nwant\n%s", result, expect)
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
		&ec2.Subnet{SubnetId: aws.String("sub_3"), VpcId: nil}, // edge case subnet with no vpc id
	}

	triples, err := BuildInfraRdfTriples("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}

	result := SortLines(MarshalTriples(triples))
	expect := `/region<eu-west-1>	"parent_of"@[]	/vpc<vpc_1>
/region<eu-west-1>	"parent_of"@[]	/vpc<vpc_2>
/subnet<sub_1>	"parent_of"@[]	/instance<inst_1>
/subnet<sub_2>	"parent_of"@[]	/instance<inst_2>
/subnet<sub_3>	"parent_of"@[]	/instance<inst_3>
/subnet<sub_3>	"parent_of"@[]	/instance<inst_4>
/vpc<vpc_1>	"parent_of"@[]	/subnet<sub_1>
/vpc<vpc_1>	"parent_of"@[]	/subnet<sub_2>
/vpc<vpc_2>	"parent_of"@[]	/subnet<sub_3>`
	if result != expect {
		t.Fatalf("got %s\nwant %s", result, expect)
	}
}

func TestBuildEmptyRdfTriplesWhenNoData(t *testing.T) {
	triples, err := BuildAccessRdfTriples("any", api.NewAwsAccess())
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(triples), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	triples, err = BuildInfraRdfTriples("any", &api.AwsInfra{})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(triples), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestIntersectTriples(t *testing.T) {
	var a, b, expect []*triple.Triple

	a = append(a, parseTriple("/a<1>	\"to\"@[]	/b<1>"))
	a = append(a, parseTriple("/a<2>	\"to\"@[]	/b<2>"))
	a = append(a, parseTriple("/a<3>	\"to\"@[]	/b<3>"))
	a = append(a, parseTriple("/a<4>	\"to\"@[]	/b<4>"))

	b = append(b, parseTriple("/a<0>	\"to\"@[]	/b<0>"))
	b = append(b, parseTriple("/a<2>	\"to\"@[]	/b<2>"))
	b = append(b, parseTriple("/a<3>	\"to\"@[]	/b<3>"))
	b = append(b, parseTriple("/a<5>	\"to\"@[]	/b<5>"))
	b = append(b, parseTriple("/a<6>	\"to\"@[]	/b<6>"))

	result := IntersectTriples(a, b)
	expect = append(expect, parseTriple("/a<2>	\"to\"@[]	/b<2>"))
	expect = append(expect, parseTriple("/a<3>	\"to\"@[]	/b<3>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestSubstractTriples(t *testing.T) {
	var a, b, expect []*triple.Triple

	a = append(a, parseTriple("/a<1>	\"to\"@[]	/b<1>"))
	a = append(a, parseTriple("/a<2>	\"to\"@[]	/b<2>"))
	a = append(a, parseTriple("/a<3>	\"to\"@[]	/b<3>"))
	a = append(a, parseTriple("/a<4>	\"to\"@[]	/b<4>"))

	b = append(b, parseTriple("/a<0>	\"to\"@[]	/b<0>"))
	b = append(b, parseTriple("/a<2>	\"to\"@[]	/b<2>"))
	b = append(b, parseTriple("/a<3>	\"to\"@[]	/b<3>"))
	b = append(b, parseTriple("/a<5>	\"to\"@[]	/b<5>"))
	b = append(b, parseTriple("/a<6>	\"to\"@[]	/b<6>"))

	result := SubstractTriples(a, b)
	expect = append(expect, parseTriple("/a<1>	\"to\"@[]	/b<1>"))
	expect = append(expect, parseTriple("/a<4>	\"to\"@[]	/b<4>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	result = SubstractTriples(b, a)
	expect = []*triple.Triple{}
	expect = append(expect, parseTriple("/a<0>	\"to\"@[]	/b<0>"))
	expect = append(expect, parseTriple("/a<5>	\"to\"@[]	/b<5>"))
	expect = append(expect, parseTriple("/a<6>	\"to\"@[]	/b<6>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func parseTriple(s string) *triple.Triple {
	t, err := triple.Parse(s, literal.DefaultBuilder())
	if err != nil {
		panic(err)
	}

	return t
}

func SortLines(lines string) string {
	linesToSort := strings.Split(lines, "\n")
	sort.Strings(linesToSort)
	return strings.Join(linesToSort, "\n")
}
