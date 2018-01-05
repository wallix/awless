package commands

import (
	"reflect"
	"testing"

	"sort"

	p "github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestEnumCompletionFunc(t *testing.T) {
	list, _ := enumCompletionFunc([]string{"am", "amstram", "gram"}).Do([]rune{'a'}, 1)
	if got, want := list, toRune("m ", "mstram "); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestTypedParamCompletionFunc(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(resourcetest.Instance("1").Prop(p.Name, "broker_1").Build())
	g.AddResource(resourcetest.Instance("2").Prop(p.Name, "broker_2").Build())
	g.AddResource(resourcetest.Instance("3").Prop(p.Name, "redis").Build())
	g.AddResource(resourcetest.Instance("4").Prop(p.Records, []string{"rec1", "rec2"}).Build())

	list, _ := typedParamCompletionFunc(g, "instance", p.Name).Do([]rune{'b'}, 1)
	sort.Slice(list, func(i, j int) bool { return string(list[i]) <= string(list[j]) })
	if got, want := list, toRune("roker_1 ", "roker_2 "); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
	list, _ = typedParamCompletionFunc(g, "instance", p.Records).Do([]rune{'r', 'e'}, 2)
	sort.Slice(list, func(i, j int) bool { return string(list[i]) <= string(list[j]) })
	if got, want := list, toRune("c1 ", "c2 "); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAutoCompletion(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(resourcetest.Instance("1").Prop(p.Name, "broker_1").Prop(p.Type, "t2.micro").Prop(p.Subnet, "1").Prop(p.ActiveServicesCount, 42).Build())
	g.AddResource(resourcetest.Instance("2").Prop(p.Name, "broker_2").Prop(p.Type, "t3.medium").Prop(p.Subnet, "2").Prop(p.ActiveServicesCount, 24).Build())
	g.AddResource(resourcetest.Instance("3").Prop(p.Name, "kafka").Prop(p.Type, "t3.medium").Prop(p.Subnet, "2").Prop(p.ActiveServicesCount, 44).Build())
	g.AddResource(resourcetest.Instance("4").Prop(p.Name, "redis").Build())
	g.AddResource(resourcetest.SecurityGroup("sg-1").Prop(p.Name, "ssh").Build())
	g.AddResource(resourcetest.Subnet("s-5").Prop(p.Name, "subnet 1").Prop(p.Public, true).Prop(p.CIDR, "10.0.0.0/0").Build())
	g.AddResource(resourcetest.Subnet("s-6").Prop(p.Name, "subnet 2").Prop(p.Public, false).Prop(p.CIDR, "192.168.0.0/0").Build())
	g.AddResource(resourcetest.Alarm("1").Prop(p.Dimensions, []*graph.KeyValue{{"abc", "val1"}, {"abd", "val2"}}).Build())
	g.AddResource(resourcetest.Alarm("2").Prop(p.Dimensions, []*graph.KeyValue{{"def", "val3"}}).Build())

	t.Run("no matches", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.instance.id"}).Do([]rune{'a'}, 1)
		if len(list) != 0 {
			t.Fatalf("expected empty, got %q", list)
		}
	})

	t.Run("empty when property valid but undefined", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.subnet.vpc"}).Do([]rune{}, 0)
		if len(list) != 0 {
			t.Fatalf("expected empty, got %q", list)
		}
	})

	t.Run("fallbacks on entity when property invalid", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.subnet.wrong"}).Do([]rune{}, 0)
		if got, want := list, toRune("'@subnet 1'", "'@subnet 2'", "s-5", "s-6"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("match on entity type", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.subnet.instance"}).Do([]rune{'@', 'b'}, 2)
		if got, want := list, toRune("roker_1", "roker_2"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.subnet.instance"}).Do([]rune{'@', 'r'}, 2)
		if got, want := list, toRune("edis"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		var empty [][]rune
		list, _ = holeAutoCompletion(g, []string{"create.subnet.instance"}).Do([]rune{'@', 'r', 'e', 'd', 'i', 's'}, 6)
		if got, want := list, empty; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("match on entity regular property. Entity has a Name property", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.securitygroup.id"}).Do([]rune{}, 0)
		if got, want := list, toRune("@ssh", "sg-1"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("match on entity regular property", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.subnet.cidr"}).Do([]rune{}, 0)
		if got, want := list, toRune("10.0.0.0/0", "192.168.0.0/0"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.instance.type"}).Do([]rune{}, 0)
		if got, want := list, toRune("t2.micro", "t3.medium"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("match on non string property", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.instance.activeservicescount"}).Do([]rune{}, 0)
		if got, want := list, toRune("24", "42", "44"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.instance.activeservicescount"}).Do([]rune{'4'}, 1)
		if got, want := list, toRune("2", "4"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.subnet.public"}).Do([]rune{'t', 'r'}, 2)
		if got, want := list, toRune("ue"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("match on property with '-'", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.instance.active-services-count"}).Do([]rune{}, 0)
		if got, want := list, toRune("24", "42", "44"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("match on entity property which is itself an entity", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.instance.subnet"}).Do([]rune{}, 0)
		if got, want := list, toRune("'@subnet 1'", "'@subnet 2'", "s-5", "s-6"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("suggests list of entities for plural resources", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.instance.subnets"}).Do([]rune{}, 0)
		if got, want := list, toRune("'@subnet 1'", "'@subnet 2'", "s-5", "s-6"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.instance.subnets"}).Do([]rune{'s', '-', '5', ','}, 4)
		if got, want := list, toRune("'@subnet 1'", "'@subnet 2'", "s-5", "s-6"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.instance.subnets"}).Do([]rune{'s', '-', '5', ',', 's', '-'}, 6)
		if got, want := list, toRune("5", "6"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.instance.subnets"}).Do([]rune{'s', '-', '5', ',', 's', '-', '6', ','}, 8)
		if got, want := list, toRune("'@subnet 1'", "'@subnet 2'", "s-5", "s-6"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.instance.subnets"}).Do([]rune{'s', '-', '5', ',', '\'', '@', 's', 'u'}, 7)
		if got, want := list, toRune("bnet 1'", "bnet 2'"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("suggests multiple properties when containing a ','", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.alarm.dimensions"}).Do([]rune{}, 0)
		if got, want := list, toRune("abc:val1", "abd:val2", "def:val3"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.alarm.dimensions"}).Do([]rune{'a', 'b'}, 2)
		if got, want := list, toRune("c:val1", "d:val2"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.alarm.dimensions"}).Do([]rune{'a', 'b', 'c', ':', 'v', 'a', 'l', '1', ','}, 9)
		if got, want := list, toRune("abc:val1", "abd:val2", "def:val3"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("support [ for multiple properties ", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.alarm.dimensions"}).Do([]rune{'[', 'a', 'b'}, 3)
		if got, want := list, toRune("c:val1", "d:val2"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
		list, _ = holeAutoCompletion(g, []string{"create.alarm.dimensions"}).Do([]rune{'[', 'a', 'b', 'c', ':', 'v', 'a', 'l', '1', ','}, 10)
		if got, want := list, toRune("abc:val1", "abd:val2", "def:val3"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("with multiple param paths", func(t *testing.T) {
		list, _ := holeAutoCompletion(g, []string{"create.instance.subnet", "update.securitygroup.name"}).Do([]rune{'s'}, 1)
		if got, want := list, toRune("-5", "-6", "sh"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

func TestGuessEntityType(t *testing.T) {
	tcases := []struct {
		hole, prop string
		types      []string
	}{
		{hole: ""},
		{hole: "any"},
		{hole: "inst"},
		{hole: "gateway"},
		{hole: "gateway.", types: []string{"internetgateway", "natgateway"}},
		{hole: "instance", types: []string{"instance"}},
		{hole: "instance.ip", types: []string{"instance"}, prop: "ip"},
		{hole: "securitygroup.id", types: []string{"securitygroup"}, prop: "id"},
		{hole: "subnets", types: []string{"subnet"}},
		{hole: "subnet", types: []string{"subnet"}},
		{hole: "subnets.cidr", types: []string{"subnet"}, prop: "cidr"},
		{hole: "subnet.cidr", types: []string{"subnet"}, prop: "cidr"},
		{hole: "subnet.cidr.any", types: []string{"subnet"}},
		{hole: "vpc.instance", types: []string{"instance"}},
		{hole: "route.gateway", types: []string{"internetgateway", "natgateway"}},
		{hole: "route.table", types: []string{"routetable"}},

		{hole: "zone.1", types: []string{"zone"}, prop: "1"},
		{hole: "availabilityzone.1", types: []string{"availabilityzone"}, prop: "1"},

		{hole: "gateway.1", types: []string{"internetgateway", "natgateway"}},
		{hole: "gateway.in", types: []string{"internetgateway", "natgateway"}},
		{hole: "gateway.inst", types: []string{"instance", "containerinstance", "instanceprofile"}},

		{hole: "gateway.inst.any", types: []string{"instance", "containerinstance", "instanceprofile"}},
		{hole: "gateway.any.inst", types: []string{"instance", "containerinstance", "instanceprofile"}},
		{hole: "inst.any.any", types: []string{"instance", "containerinstance", "instanceprofile"}},
	}

	for _, tcase := range tcases {
		types, prop := guessEntityTypeFromHoleQuestion(tcase.hole)
		if got, want := types, tcase.types; !reflect.DeepEqual(got, want) {
			t.Fatalf("case '%s': got %v, want %v", tcase.hole, got, want)
		}
		if got, want := prop, tcase.prop; got != want {
			t.Fatalf("case '%s': property: got '%s', want '%s'", tcase.hole, got, want)
		}
	}
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

	sort.Slice(out, func(i, j int) bool { return string(out[i]) <= string(out[j]) })

	return out
}
