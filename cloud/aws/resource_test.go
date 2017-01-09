package aws

import (
	"encoding/json"
	"reflect"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
	"github.com/wallix/awless/rdf"
)

func TestUnmarshalResource(t *testing.T) {
	res := Resource{id: "inst_1", kind: rdf.INSTANCE}

	g := rdf.NewGraph()
	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Tags","Value":[{"Key":"Name","Value":"redis"}]}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Type","Value":"t2.micro"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"PublicIp","Value":"1.2.3.4"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"State","Value":{"Code": 16,"Name":"running"}}"^^type:text`))

	res.UnmarshalFromGraph(g)

	expected := Properties{"Id": "inst_1", "Type": "t2.micro", "PublicIp": "1.2.3.4",
		"State": map[string]interface{}{"Code": float64(16), "Name": "running"},
		"Tags": []interface{}{
			map[string]interface{}{"Key": "Name", "Value": "redis"},
		},
	}

	if got, want := res.properties, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got \n%#v\n\nwant \n%#v\n", got, want)
	}
}

func TestLoadPropertiesTriples(t *testing.T) {
	g := rdf.NewGraph()

	aLiteral, err := literal.DefaultBuilder().Build(literal.Text, mustJsonMarshal(Property{Key: "prop1", Value: "val1"}))
	if err != nil {
		t.Fatal(err)
	}
	bLiteral, err := literal.DefaultBuilder().Build(literal.Text, mustJsonMarshal(Property{Key: "prop2", Value: "val2"}))
	if err != nil {
		t.Fatal(err)
	}
	cLiteral, err := literal.DefaultBuilder().Build(literal.Text, mustJsonMarshal(Property{Key: "prop3", Value: "val3"}))
	if err != nil {
		t.Fatal(err)
	}
	dLiteral, err := literal.DefaultBuilder().Build(literal.Text, mustJsonMarshal(Property{Key: "prop4", Value: "val4"}))
	if err != nil {
		t.Fatal(err)
	}

	one, _ := node.NewNodeFromStrings("/one", "1")
	g.Add(noErrLiteralTriple(one, rdf.PropertyPredicate, aLiteral))
	g.Add(noErrLiteralTriple(one, rdf.PropertyPredicate, bLiteral))
	g.Add(noErrLiteralTriple(one, rdf.PropertyPredicate, cLiteral))
	two, _ := node.NewNodeFromStrings("/two", "2")
	g.Add(noErrLiteralTriple(two, rdf.PropertyPredicate, dLiteral))

	properties, err := LoadPropertiesFromGraph(g, one)
	if err != nil {
		t.Fatal(err)
	}
	expected := Properties{
		"prop1": "val1",
		"prop2": "val2",
		"prop3": "val3",
	}

	if got, want := properties, expected; !reflect.DeepEqual(properties, expected) {
		t.Fatalf("got %s, want %s", got, want)
	}

	properties, err = LoadPropertiesFromGraph(g, two)
	expected = Properties{
		"prop4": "val4",
	}

	if got, want := properties, expected; !reflect.DeepEqual(properties, expected) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestResourceName(t *testing.T) {
	properties := Properties{}
	if got, want := NameFromProperties(properties), ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	properties = Properties{
		"Id":    "my_id",
		"Name":  "my_name",
		"Other": "my_other",
	}
	if got, want := NameFromProperties(properties), "my_name"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	properties = Properties{
		"Id":    "my_id",
		"Other": "my_other",
	}
	if got, want := NameFromProperties(properties), ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	awsInfra := &AwsInfra{}

	awsInfra.Instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		&ec2.Instance{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
	}

	awsInfra.Vpcs = []*ec2.Vpc{
		&ec2.Vpc{VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Other"), Value: awssdk.String("Tag")}}},
		&ec2.Vpc{VpcId: awssdk.String("vpc_2"), Tags: []*ec2.Tag{}},
	}

	awsInfra.Subnets = []*ec2.Subnet{
		&ec2.Subnet{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
	}

	g, err := BuildAwsInfraGraph("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}
	properties, err = LoadPropertiesFromGraph(g, newNode(rdf.INSTANCE, "inst_1"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := NameFromProperties(properties), "instance1-name"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	properties, err = LoadPropertiesFromGraph(g, newNode(rdf.INSTANCE, "inst_2"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := NameFromProperties(properties), ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	properties, err = LoadPropertiesFromGraph(g, newNode(rdf.INSTANCE, "vpc_1"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := NameFromProperties(properties), ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func mustJsonMarshal(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func noErrLiteralTriple(s *node.Node, p *predicate.Predicate, l *literal.Literal) *triple.Triple {
	tri, err := triple.New(s, p, triple.NewLiteralObject(l))
	if err != nil {
		panic(err)
	}
	return tri
}

func newNode(t, id string) *node.Node {
	nodeT, err := node.NewType(t)
	if err != nil {
		panic(err)
	}
	nodeID, err := node.NewID(id)
	if err != nil {
		panic(err)
	}
	return node.NewNode(nodeT, nodeID)
}

func parseTriple(s string) *triple.Triple {
	t, err := triple.Parse(s, literal.DefaultBuilder())
	if err != nil {
		panic(err)
	}

	return t
}
