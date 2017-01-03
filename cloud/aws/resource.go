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
		prop, err := NewPropertyFromTriple(t)
		if err != nil {
			return properties, err
		}
		properties[prop.Key] = prop.Value
	}
	return properties, nil
}

func NewPropertyFromTriple(t *triple.Triple) (*Property, error) {
	if t.Predicate().String() != rdf.PropertyPredicate.String() {
		return nil, fmt.Errorf("This triple has not a property prediate: %s", t.Predicate().String())
	}
	oL, err := t.Object().Literal()
	if err != nil {
		return nil, err
	}
	propStr, err := oL.Text()
	if err != nil {
		return nil, err
	}
	var prop Property
	err = json.Unmarshal([]byte(propStr), &prop)
	if err != nil {
		return nil, err
	}
	return &prop, nil
}

func NameFromProperties(p Properties) string {
	if n, ok := p["Name"]; ok {
		return fmt.Sprint(n)
	}
	if t, ok := p["Tags"]; ok {
		switch tt := t.(type) {
		case []interface{}:
			for _, tag := range tt {
				//map [key: result]
				if m, ok := tag.(map[string]interface{}); ok && m["Name"] != nil {
					return fmt.Sprint(m["Name"])
				}

				//map["Key": key, "Value": result]
				if m, ok := tag.(map[string]interface{}); ok && m["Key"] == "Name" {
					return fmt.Sprint(m["Value"])
				}
			}
		}

		return fmt.Sprint(t)
	}
	return ""
}

func (res *Resource) buildRdfSubject() (*node.Node, error) {
	return node.NewNodeFromStrings(res.kind, res.id)
}

func addCloudResourceToGraph(g *rdf.Graph, cloudResource interface{}) error {
	res, err := NewResource(cloudResource)
	if err != nil {
		return err
	}
	triples, err := res.MarshalToTriples()
	if err != nil {
		return err
	}
	g.Add(triples...)
	return nil
}
