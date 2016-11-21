package store

import (
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
)

func TestBuildInfraRdfTriples(t *testing.T) {
	var instances []*ec2.Instance
	var vpcs []*ec2.Vpc
	var subnets []*ec2.Subnet

	triples, err := BuildInfraRdfTriples("eu-west-1", vpcs, subnets, instances)
	if err != nil {
		t.Fatalf("error while building triples : %s", err)
	}

	result := MarshalTriples(triples)
	expect := ``
	if result != expect {
		t.Fatalf("got %s\nwant %s", result, expect)
	}

	instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: aws.String("inst_1"), SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_2"), SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_3"), SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2")},
		&ec2.Instance{InstanceId: aws.String("inst_4"), SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2")},
		&ec2.Instance{InstanceId: aws.String("inst_5"), SubnetId: nil, VpcId: nil}, // terminated instance (no vpc, subnet ids)
	}

	vpcs = []*ec2.Vpc{
		&ec2.Vpc{VpcId: aws.String("vpc_1")},
		&ec2.Vpc{VpcId: aws.String("vpc_2")},
	}

	subnets = []*ec2.Subnet{
		&ec2.Subnet{SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Subnet{SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Subnet{SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2")},
		&ec2.Subnet{SubnetId: aws.String("sub_3"), VpcId: nil}, // edge case subnet with no vpc id
	}

	triples, err = BuildInfraRdfTriples("eu-west-1", vpcs, subnets, instances)
	if err != nil {
		t.Fatalf("error while building triples : %s", err)
	}

	result = SortLines(MarshalTriples(triples))
	expect = `/region<eu-west-1>	"parent_of"@[]	/vpc<vpc_1>
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
