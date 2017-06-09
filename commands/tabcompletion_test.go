package commands

import (
	"reflect"
	"testing"

	p "github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestAutoCompletion(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(resourcetest.Instance("1").Prop(p.Name, "broker_1").Prop(p.Type, "t2.micro").Prop(p.Subnet, "1").Build())
	g.AddResource(resourcetest.Instance("2").Prop(p.Name, "broker_2").Prop(p.Type, "t3.medium").Prop(p.Subnet, "2").Build())
	g.AddResource(resourcetest.Instance("3").Prop(p.Name, "kafka").Prop(p.Type, "t3.medium").Prop(p.Subnet, "2").Build())
	g.AddResource(resourcetest.Instance("4").Prop(p.Name, "redis").Build())
	g.AddResource(resourcetest.Subnet("5").Prop(p.Name, "subnet 1").Prop(p.CIDR, "10.0.0.0/0").Build())
	g.AddResource(resourcetest.Subnet("6").Prop(p.Name, "subnet 2").Prop(p.CIDR, "192.168.0.0/0").Build())

	t.Run("no matches", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, "instance").Do([]rune{'a'}, 1)
		if len(list) != 0 {
			t.Fatalf("expected empty, got %q", list)
		}
		list, _ = holeAutoCompletion(g, "instance.ip").Do([]rune{}, 0)
		if len(list) != 0 {
			t.Fatalf("expected empty, got %q", list)
		}

	})

	t.Run("match on entity type", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, "instance").Do([]rune{'@', 'b'}, 2)
		if got, want := list, toRune("roker_1 ", "roker_2 "); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, "instance").Do([]rune{'@', 'r'}, 2)
		if got, want := list, toRune("edis "); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("match on entity regular property", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, "subnet.cidr").Do([]rune{}, 0)
		if got, want := list, toRune("10.0.0.0/0 ", "192.168.0.0/0 "); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, "instance.type").Do([]rune{}, 0)
		if got, want := list, toRune("t2.micro ", "t3.medium "); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("match on entity property which is itself an entity", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, "instance.subnet").Do([]rune{}, 0)
		if got, want := list, toRune("'@subnet 1' ", "'@subnet 2' ", "5 ", "6 "); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

func toRune(arr ...string) [][]rune {
	out := make([][]rune, len(arr))
	for i, s := range arr {
		if out[i] == nil {
			out[i] = make([]rune, 0)
		}
		for _, r := range s {
			out[i] = append(out[i], r)
		}
	}
	return out
}
