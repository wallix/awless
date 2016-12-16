package aws

import (
	"encoding/json"
	"fmt"
	"reflect"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/rdf"
)

type Property struct {
	Key   string
	Value interface{}
}

type Properties map[string]interface{}

type Resource struct {
	kind       string
	id         string
	source     interface{}
	properties Properties
}

func NewResource(source interface{}) (*Resource, error) {
	res := Resource{}
	switch ss := source.(type) {
	case *ec2.Instance:
		res.kind = rdf.INSTANCE
		res.id = awssdk.StringValue(ss.InstanceId)
	case *ec2.Vpc:
		res.kind = rdf.VPC
		res.id = awssdk.StringValue(ss.VpcId)
	case *ec2.Subnet:
		res.kind = rdf.SUBNET
		res.id = awssdk.StringValue(ss.SubnetId)
	case *iam.User:
		res.kind = rdf.USER
		res.id = awssdk.StringValue(ss.UserId)
	case *iam.Role:
		res.kind = rdf.ROLE
		res.id = awssdk.StringValue(ss.RoleId)
	case *iam.Group:
		res.kind = rdf.GROUP
		res.id = awssdk.StringValue(ss.GroupId)
	case *iam.Policy:
		res.kind = rdf.POLICY
		res.id = awssdk.StringValue(ss.PolicyId)
	default:
		return nil, fmt.Errorf("Unknown type of resource %T", source)
	}
	res.source = source

	value := reflect.ValueOf(source)
	if !value.IsValid() || value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not a valid pointer.", value)
	}

	nodeV := value.Elem()
	res.properties = make(Properties)
	for propertyId, cloudId := range awsResourcesProperties[res.kind] {
		sourceField := nodeV.FieldByName(cloudId)
		if sourceField.IsValid() && !sourceField.IsNil() {
			res.properties[propertyId] = sourceField.Interface()
		}
	}
	return &res, nil
}

func (res *Resource) MarshalToTriples() ([]*triple.Triple, error) {
	var triples []*triple.Triple
	n, err := res.buildRdfSubject()
	if err != nil {
		return triples, err
	}
	var lit *literal.Literal
	if lit, err = literal.DefaultBuilder().Build(literal.Text, res.kind); err != nil {
		return triples, err
	}
	t, err := triple.New(n, rdf.HasTypePredicate, triple.NewLiteralObject(lit))
	if err != nil {
		return triples, err
	}
	triples = append(triples, t)

	for propKey, propValue := range res.properties {
		if propT, err := NewPropertyTriple(n, propKey, propValue); err != nil {
			return nil, err
		} else {
			triples = append(triples, propT)
		}
	}

	return triples, nil
}

func (res *Resource) buildRdfSubject() (*node.Node, error) {
	return node.NewNodeFromStrings(res.kind, res.id)
}

func NewPropertyTriple(subject *node.Node, propertyKey string, propertyValue interface{}) (*triple.Triple, error) {
	prop := Property{Key: propertyKey, Value: propertyValue}
	json, err := json.Marshal(prop)
	if err != nil {
		return nil, err
	}
	var propL *literal.Literal
	if propL, err = literal.DefaultBuilder().Build(literal.Text, string(json)); err != nil {
		return nil, err
	}
	if propT, err := triple.New(subject, rdf.PropertyPredicate, triple.NewLiteralObject(propL)); err != nil {
		return nil, err
	} else {
		return propT, nil
	}
}

func LoadPropertiesFromGraph(g *rdf.Graph, node *node.Node) (Properties, error) {
	triples, err := g.TriplesForSubjectPredicate(node, rdf.PropertyPredicate)
	if err != nil {
		return nil, err
	}
	properties := make(Properties)
	for _, t := range triples {
		oL, e := t.Object().Literal()
		if e != nil {
			return properties, e
		}
		propStr, e := oL.Text()
		if e != nil {
			return properties, e
		}
		var prop Property
		e = json.Unmarshal([]byte(propStr), &prop)
		if e != nil {
			return properties, e
		}
		properties[prop.Key] = prop.Value
	}
	return properties, nil
}
