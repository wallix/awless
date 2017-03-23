package resourcetest

import (
	"fmt"
	"strings"

	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
)

type rBuilder struct {
	id, typ string
	props   map[string]interface{}
}

func new(typ, id string) *rBuilder {
	return &rBuilder{id: id, typ: typ, props: make(map[string]interface{})}
}

func Region(id string) *rBuilder {
	return new("region", id)
}

func Instance(id string) *rBuilder {
	return new("instance", id).Prop(properties.ID, id)
}

func Subnet(id string) *rBuilder {
	return new("subnet", id).Prop(properties.ID, id)
}

func VPC(id string) *rBuilder {
	return new("vpc", id).Prop(properties.ID, id)
}

func SecGroup(id string) *rBuilder {
	return new("securitygroup", id).Prop(properties.ID, id)
}

func Keypair(id string) *rBuilder {
	return new("keypair", id).Prop(properties.ID, id)
}

func InternetGw(id string) *rBuilder {
	return new("internetgateway", id).Prop(properties.ID, id)
}

func RouteTable(id string) *rBuilder {
	return new("routetable", id).Prop(properties.ID, id)
}

func LoadBalancer(id string) *rBuilder {
	return new("loadbalancer", id).Prop(properties.ID, id)
}

func TargetGroup(id string) *rBuilder {
	return new("targetgroup", id).Prop(properties.ID, id)
}

func Policy(id string) *rBuilder {
	return new("policy", id).Prop(properties.ID, id)
}

func Group(id string) *rBuilder {
	return new("group", id).Prop(properties.ID, id)
}

func Role(id string) *rBuilder {
	return new("role", id).Prop(properties.ID, id)
}

func User(id string) *rBuilder {
	return new("user", id).Prop(properties.ID, id)
}

func Listener(id string) *rBuilder {
	return new("listener", id).Prop(properties.ID, id)
}

func Bucket(id string) *rBuilder {
	return new("bucket", id).Prop(properties.ID, id)
}

func Zone(id string) *rBuilder {
	return new("zone", id).Prop(properties.ID, id)
}

func Record(id string) *rBuilder {
	return new("record", id).Prop(properties.ID, id)
}

func (b *rBuilder) Prop(key string, value interface{}) *rBuilder {
	b.props[key] = value
	return b
}

func (b *rBuilder) Build() *graph.Resource {
	res := graph.InitResource(b.typ, b.id)
	res.Properties = b.props
	return res
}

func AddParents(g *graph.Graph, relations ...string) {
	for _, rel := range relations {
		splits := strings.Split(rel, "->")
		if len(splits) != 2 {
			panic(fmt.Sprintf("invalid relation '%s'", rel))
		}
		r1 := graph.InitResource("", strings.TrimSpace(splits[0]))
		r2 := graph.InitResource("", strings.TrimSpace(splits[1]))
		err := g.AddParentRelation(r1, r2)
		if err != nil {
			panic(err)
		}
	}
}
