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
	r := &rBuilder{id: id, typ: typ, props: make(map[string]interface{})}
	return r.Prop(properties.ID, id)
}

func Region(id string) *rBuilder {
	return new("region", id)
}

func Instance(id string) *rBuilder {
	return new("instance", id)
}

func Subnet(id string) *rBuilder {
	return new("subnet", id)
}

func VPC(id string) *rBuilder {
	return new("vpc", id)
}

func SecurityGroup(id string) *rBuilder {
	return new("securitygroup", id)
}

func KeyPair(id string) *rBuilder {
	return new("keypair", id)
}

func InternetGw(id string) *rBuilder {
	return new("internetgateway", id)
}

func NatGw(id string) *rBuilder {
	return new("natgateway", id)
}

func RouteTable(id string) *rBuilder {
	return new("routetable", id)
}

func LoadBalancer(id string) *rBuilder {
	return new("loadbalancer", id)
}

func ClassicLoadBalancer(id string) *rBuilder {
	return new("classicloadbalancer", id)
}

func AvailabilityZone(id string) *rBuilder {
	return new("availabilityzone", id)
}

func TargetGroup(id string) *rBuilder {
	return new("targetgroup", id)
}

func Policy(id string) *rBuilder {
	return new("policy", id)
}

func Group(id string) *rBuilder {
	return new("group", id)
}

func Role(id string) *rBuilder {
	return new("role", id)
}

func User(id string) *rBuilder {
	return new("user", id)
}

func MfaDevice(id string) *rBuilder {
	return new("mfadevice", id)
}

func Listener(id string) *rBuilder {
	return new("listener", id)
}

func Bucket(id string) *rBuilder {
	return new("bucket", id)
}

func Zone(id string) *rBuilder {
	return new("zone", id)
}

func Record(id string) *rBuilder {
	return new("record", id)
}

func ScalingGroup(id string) *rBuilder {
	return new("scalinggroup", id)
}

func LaunchConfig(id string) *rBuilder {
	return new("launchconfiguration", id)
}

func Subscription(id string) *rBuilder {
	return new("subscription", id)
}

func Topic(id string) *rBuilder {
	return new("topic", id)
}

func Queue(id string) *rBuilder {
	return new("queue", id)
}

func Function(id string) *rBuilder {
	return new("function", id)
}

func Alarm(id string) *rBuilder {
	return new("alarm", id)
}

func Metric(id string) *rBuilder {
	return new("metric", id)
}

func Image(id string) *rBuilder {
	return new("image", id)
}

func Distribution(id string) *rBuilder {
	return new("distribution", id)
}

func Stack(id string) *rBuilder {
	return new("stack", id)
}

func Repository(id string) *rBuilder {
	return new("repository", id)
}

func ContainerCluster(id string) *rBuilder {
	return new("containercluster", id)
}

func ContainerTask(id string) *rBuilder {
	return new("containertask", id)
}

func Container(id string) *rBuilder {
	return new("container", id)
}

func ContainerInstance(id string) *rBuilder {
	return new("containerinstance", id)
}

func NetworkInterface(id string) *rBuilder {
	return new("networkinterface", id)
}

func Certificate(id string) *rBuilder {
	return new("certificate", id)
}

func AccessKey(id string) *rBuilder {
	return new("accesskey", id)
}

func (b *rBuilder) Prop(key string, value interface{}) *rBuilder {
	b.props[key] = value
	return b
}

func (b *rBuilder) Build() *graph.Resource {
	res := graph.InitResource(b.typ, b.id)
	for k, v := range b.props {
		res.Properties()[k] = v
	}

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
