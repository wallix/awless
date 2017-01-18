package aws

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

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

func (prop *Property) tripleFromNode(subject *node.Node) (*triple.Triple, error) {
	propL, err := prop.ToLiteralObject()
	if err != nil {
		return nil, err
	}
	if propT, err := triple.New(subject, rdf.PropertyPredicate, propL); err != nil {
		return nil, err
	} else {
		return propT, nil
	}
}

func (prop *Property) ToLiteralObject() (*triple.Object, error) {
	json, err := json.Marshal(prop)
	if err != nil {
		return nil, err
	}
	var propL *literal.Literal
	if propL, err = literal.DefaultBuilder().Build(literal.Text, string(json)); err != nil {
		return nil, err
	}
	return triple.NewLiteralObject(propL), nil
}

type Properties map[string]interface{}

type Resource struct {
	kind       rdf.ResourceType
	id         string
	properties Properties
}

func InitResource(id string, kind rdf.ResourceType) *Resource {
	return &Resource{id: id, kind: kind, properties: make(Properties)}
}

func InitFromRdfNode(n *node.Node) *Resource {
	return InitResource(n.ID().String(), rdf.NewResourceType(n.Type()))
}

func NewResource(source interface{}) (*Resource, error) {
	value := reflect.ValueOf(source)
	if !value.IsValid() || value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not a valid pointer.", value)
	}
	nodeV := value.Elem()

	var res *Resource
	switch ss := source.(type) {
	case *ec2.Instance:
		res = InitResource(awssdk.StringValue(ss.InstanceId), rdf.Instance)
	case *ec2.Vpc:
		res = InitResource(awssdk.StringValue(ss.VpcId), rdf.Vpc)
	case *ec2.Subnet:
		res = InitResource(awssdk.StringValue(ss.SubnetId), rdf.Subnet)
	case *ec2.SecurityGroup:
		res = InitResource(awssdk.StringValue(ss.GroupId), rdf.SecurityGroup)
	case *iam.User:
		res = InitResource(awssdk.StringValue(ss.UserId), rdf.User)
	case *iam.Role:
		res = InitResource(awssdk.StringValue(ss.RoleId), rdf.Role)
	case *iam.Group:
		res = InitResource(awssdk.StringValue(ss.GroupId), rdf.Group)
	case *iam.Policy:
		res = InitResource(awssdk.StringValue(ss.PolicyId), rdf.Policy)
	default:
		return nil, fmt.Errorf("Unknown type of resource %T", source)
	}

	for prop, trans := range awsResourcesDef[res.kind] {
		sourceField := nodeV.FieldByName(trans.name)
		if sourceField.IsValid() && !sourceField.IsNil() {
			val, err := trans.transform(sourceField.Interface())
			if err == ErrTagNotFound {
				continue
			}
			if err != nil {
				return res, err
			}
			res.properties[prop] = val
		}
	}

	return res, nil
}

func (res *Resource) ExistsInGraph(g *rdf.Graph) bool {
	r, err := res.buildRdfSubject()
	if err != nil {
		return false
	}

	nodes, err := g.NodesForType(res.kind)
	if err != nil {
		return false
	}
	for _, n := range nodes {
		if n.UUID().String() == r.UUID().String() {
			return true
		}
	}
	return false
}

func (res *Resource) UnmarshalFromGraph(g *rdf.Graph) error {
	node, err := res.buildRdfSubject()
	if err != nil {
		return err
	}

	triples, err := g.TriplesForSubjectPredicate(node, rdf.PropertyPredicate)
	if err != nil {
		return err
	}

	for _, t := range triples {
		prop, err := NewPropertyFromTriple(t)
		if err != nil {
			return err
		}
		res.properties[prop.Key] = prop.Value
	}

	return nil
}

func (res *Resource) Properties() Properties {
	return res.properties
}

func (res *Resource) Type() rdf.ResourceType {
	return res.kind
}

func (res *Resource) Id() string {
	return res.id
}

func (res *Resource) MarshalToTriples() ([]*triple.Triple, error) {
	var triples []*triple.Triple
	n, err := res.buildRdfSubject()
	if err != nil {
		return triples, err
	}
	var lit *literal.Literal
	if lit, err = literal.DefaultBuilder().Build(literal.Text, res.kind.ToRDFString()); err != nil {
		return triples, err
	}
	t, err := triple.New(n, rdf.HasTypePredicate, triple.NewLiteralObject(lit))
	if err != nil {
		return triples, err
	}
	triples = append(triples, t)

	for propKey, propValue := range res.properties {
		prop := Property{Key: propKey, Value: propValue}
		if propT, err := prop.tripleFromNode(n); err != nil {
			return nil, err
		} else {
			triples = append(triples, propT)
		}
	}

	return triples, nil
}

func LoadResourcesFromGraph(g *rdf.Graph, t rdf.ResourceType) ([]*Resource, error) {
	var res []*Resource
	nodes, err := g.NodesForType(t)
	if err != nil {
		return res, err
	}

	for _, node := range nodes {
		r := InitResource(node.ID().String(), t)
		if err := r.UnmarshalFromGraph(g); err != nil {
			return res, err
		}
		res = append(res, r)
	}
	return res, nil
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

	switch {
	case strings.HasSuffix(strings.ToLower(prop.Key), "time"), strings.HasSuffix(strings.ToLower(prop.Key), "date"):
		t, err := time.Parse(time.RFC3339, fmt.Sprint(prop.Value))
		if err == nil {
			prop.Value = t
		}
	}

	return &prop, nil
}

func (res *Resource) buildRdfSubject() (*node.Node, error) {
	return node.NewNodeFromStrings(res.kind.ToRDFString(), res.id)
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
