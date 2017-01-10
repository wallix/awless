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
	res := Resource{id: "inst_1", kind: rdf.Instance}

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

func TestLoadResources(t *testing.T) {
	g := rdf.NewGraph()
	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
	/instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Name","Value":"redis2"}"^^type:text
	/instance<inst_3>  "has_type"@[] "/instance"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Id","Value":"inst_3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Name","Value":"redis3"}"^^type:text
	/instance<subnet>  "has_type"@[] "/subnet"^^type:text
  /instance<subnet>  "property"@[] "{"Key":"Id","Value":"my subnet"}"^^type:text`))

	expected := []*Resource{
		{kind: rdf.Instance, id: "inst_1", properties: Properties{"Id": "inst_1", "Name": "redis"}},
		{kind: rdf.Instance, id: "inst_2", properties: Properties{"Id": "inst_2", "Name": "redis2"}},
		{kind: rdf.Instance, id: "inst_3", properties: Properties{"Id": "inst_3", "Name": "redis3"}},
	}
	res, err := LoadResourcesFromGraph(g, rdf.Instance)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), len(expected); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for _, r := range expected {
		found := false
		for _, r2 := range res {
			if r2.kind == r.kind && r2.id == r.id && reflect.DeepEqual(r2.properties, r.properties) {
				found = true
			}
		}
		if !found {
			t.Fatalf("%+v not found", r)
		}
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
	properties, err = LoadPropertiesFromGraph(g, newNode(rdf.Instance.ToRDFType(), "inst_1"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := NameFromProperties(properties), "instance1-name"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	properties, err = LoadPropertiesFromGraph(g, newNode(rdf.Instance.ToRDFType(), "inst_2"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := NameFromProperties(properties), ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	properties, err = LoadPropertiesFromGraph(g, newNode(rdf.Instance.ToRDFType(), "vpc_1"))
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
