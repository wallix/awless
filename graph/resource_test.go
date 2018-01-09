/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graph

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"runtime/pprof"
	"sort"
	"testing"
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
)

func TestEqualResources(t *testing.T) {
	i1 := &Resource{id: "inst_1", kind: "instance"}
	i2 := &Resource{id: "inst_2", kind: "instance"}
	i3 := &Resource{id: "toto", kind: "instance"}
	s1 := &Resource{id: "subnet_1", kind: "subnet"}
	s2 := &Resource{id: "subnet_1", kind: "subnet"}
	s3 := &Resource{id: "toto", kind: "subnet"}
	empty := &Resource{}
	tcases := []struct {
		from, to cloud.Resource
		exp      bool
	}{
		{from: i1, to: i1, exp: true},
		{from: i1, to: i2, exp: false},
		{from: i1, to: i3, exp: false},
		{from: i1, to: s1, exp: false},
		{from: s1, to: s2, exp: true},
		{from: s2, to: s1, exp: true},
		{from: s1, to: s3, exp: false},
		{from: i3, to: s3, exp: false},
		{from: empty, to: empty, exp: true},
		{from: empty, to: nil, exp: false},
		{from: empty, to: i1, exp: false},
		{from: i1, to: empty, exp: false},
	}

	for _, tcase := range tcases {
		if tcase.from.Same(tcase.to) != tcase.exp {
			t.Fatalf("expected %t, from %+v, to %+v", tcase.exp, tcase.from, tcase.to)
		}
	}
}

func TestPrintResource(t *testing.T) {
	tcases := []struct {
		res *Resource
		exp string
	}{
		{res: &Resource{id: "inst_1", kind: "instance"}, exp: "inst_1[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", properties: map[string]interface{}{"ID": "notthis"}}, exp: "inst_1[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", properties: map[string]interface{}{"ID": "notthis", "Name": "to-display"}}, exp: "@to-display[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", properties: map[string]interface{}{"Name": ""}}, exp: "inst_1[instance]"},
		{res: &Resource{kind: "instance", properties: map[string]interface{}{"ID": "notthis", "Name": "to-display"}}, exp: "@to-display[instance]"},
		{res: &Resource{}, exp: "[<none>]"},
		{res: nil, exp: "[<none>]"},
	}
	for _, tcase := range tcases {
		if got, want := tcase.res.String(), tcase.exp; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func TestFormatResource(t *testing.T) {
	tcases := []struct {
		res         *Resource
		layout, exp string
	}{
		{res: &Resource{id: "inst_1", kind: "instance"}, layout: "%i[%t]", exp: "inst_1[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", properties: map[string]interface{}{"Name": "to-display"}}, layout: "%n[%t]", exp: "@to-display[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", properties: map[string]interface{}{"Name": "to-display"}}, layout: "@%[Name]p[%t]", exp: "@to-display[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", properties: map[string]interface{}{"Test": "my-test-prop"}}, layout: "%i:%t:%[Test]p:%[Missing]p:", exp: "inst_1:instance:my-test-prop::"},
		{res: &Resource{id: "prop%n", kind: "instance", properties: map[string]interface{}{"Name": "to-display"}}, layout: "%i", exp: "prop%n"},
		{res: &Resource{id: "inst_1", kind: "instance"}, layout: "", exp: ""},
		{res: &Resource{id: "inst_1", kind: ""}, layout: "%i:%t", exp: "inst_1:<none>"},
		{res: &Resource{}, layout: "%i:%t", exp: "<none>:<none>"},
	}
	for i, tcase := range tcases {
		if got, want := tcase.res.Format(tcase.layout), tcase.exp; got != want {
			t.Fatalf("%d: got %s, want %s", i+1, got, want)
		}
	}
}

func TestReduceResources(t *testing.T) {
	res := Resources{{id: "1"}, {id: "2"}, {id: "3"}}
	if got, want := res.Map(func(r *Resource) string { return r.Format("%i") }), []string{"1", "2", "3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCompareProperties(t *testing.T) {
	props1 := map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": "3",
		"four":  4,
	}
	props2 := map[string]interface{}{
		"zero":  0,
		"two":   2,
		"three": "3",
		"four":  "4",
		"five":  "5",
	}

	exp := map[string]interface{}{"one": 1, "four": 4}
	if got, want := Subtract(props1, props2), exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	exp = map[string]interface{}{"zero": 0, "four": "4", "five": "5"}
	if got, want := Subtract(props2, props1), exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestMarshalUnmarshalFullRdf(t *testing.T) {
	res := []*Resource{
		instResource("inst1").prop(properties.ID, "inst1").prop(properties.Name, "inst1_name").prop(properties.Subnet, "sub1").prop(properties.Vpc, "vpc1").prop(properties.Launched, time.Now().UTC()).build(),
		subResource("sub1").prop(properties.ID, "sub1").prop(properties.Vpc, "vpc1").prop(properties.Default, true).build(),
		vpcResource("vpc1").prop(properties.ID, "vpc1").build(),
	}
	for _, r := range res {
		g := NewGraph()
		triples, err := r.marshalFullRDF()
		if err != nil {
			t.Fatal(err)
		}
		g.store.Add(triples...)
		rawRes := InitResource(r.Type(), r.Id())
		err = rawRes.unmarshalFullRdf(g.store.Snapshot())
		if err != nil {
			t.Fatal(err)
		}
		if got, want := rawRes, r; !reflect.DeepEqual(got, want) {
			t.Fatalf("got\n%#v\nwant\n%#v\n", got, want)
		}
	}
}

func TestResourceMarshalWithoutProps(t *testing.T) {
	res := InitResource("mytype", "myid")
	if got, want := res.Properties(), map[string]interface{}{properties.ID: "myid"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%#v\nwant\n%#v\n", got, want)
	}
}

func TestResourceMarshalUnmarshalHandlesNilValues(t *testing.T) {
	r := instResource("inst1").prop(properties.ID, "inst1").prop(properties.Host, nil).prop(properties.Name, "any").build()
	g := NewGraph()
	triples, err := r.marshalFullRDF()
	if err != nil {
		t.Fatal(err)
	}
	g.store.Add(triples...)

	rawRes := InitResource(r.Type(), r.Id())
	err = rawRes.unmarshalFullRdf(g.store.Snapshot())
	if err != nil {
		t.Fatal(err)
	}
	expectedNoNil := instResource("inst1").prop(properties.ID, "inst1").prop(properties.Name, "any").build()
	if got, want := rawRes, expectedNoNil; !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%#v\nwant\n%#v\n", got, want)
	}
}

func TestMarshalUnmarshalList(t *testing.T) {
	r := instResource("inst1").prop(properties.ID, "inst1").prop(properties.SecurityGroups, []string{"sgroup1", "sgroup2", "sgroup3"}).prop(properties.Actions, []string{"start", "stop", "delete"}).build()
	g := NewGraph()
	triples, err := r.marshalFullRDF()
	if err != nil {
		t.Fatal(err)
	}
	g.store.Add(triples...)
	rawRes := InitResource(r.Type(), r.Id())
	err = rawRes.unmarshalFullRdf(g.store.Snapshot())
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(r.Properties()["SecurityGroups"].([]string))
	sort.Strings(rawRes.Properties()["SecurityGroups"].([]string))

	sort.Strings(r.Properties()["Actions"].([]string))
	sort.Strings(rawRes.Properties()["Actions"].([]string))

	if got, want := rawRes, r; !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%#v\nwant\n%#v\n", got, want)
	}
}

func TestMarshalUnmarshalFirewallRules(t *testing.T) {
	_, localhost, _ := net.ParseCIDR("127.0.0.1/32")
	_, subnetcidr, _ := net.ParseCIDR("10.192.24.0/24")
	r := sGrpResource("sgroup1").prop(properties.ID, "sgroup1").prop(
		"InboundRules", []*FirewallRule{
			{PortRange: PortRange{FromPort: 80, ToPort: 80}, Protocol: "tcp"},
			{PortRange: PortRange{FromPort: 1, ToPort: 1024}, Protocol: "udp", IPRanges: []*net.IPNet{subnetcidr}},
		}).prop(
		"OutboundRules", []*FirewallRule{
			{PortRange: PortRange{Any: true}, Protocol: "icmp", IPRanges: []*net.IPNet{localhost, {IP: net.ParseIP("::1"), Mask: net.CIDRMask(128, 128)}}},
		}).build()
	g := NewGraph()
	triples, err := r.marshalFullRDF()
	if err != nil {
		t.Fatal(err)
	}
	g.store.Add(triples...)
	rawRes := InitResource(r.Type(), r.Id())
	err = rawRes.unmarshalFullRdf(g.store.Snapshot())
	if err != nil {
		t.Fatal(err)
	}
	FirewallRules(rawRes.Properties()["InboundRules"].([]*FirewallRule)).Sort()
	FirewallRules(rawRes.Properties()["OutboundRules"].([]*FirewallRule)).Sort()
	FirewallRules(r.Properties()["InboundRules"].([]*FirewallRule)).Sort()
	FirewallRules(r.Properties()["OutboundRules"].([]*FirewallRule)).Sort()

	if got, want := rawRes, r; !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%#v\nwant\n%#v\n", got, want)
	}
}

func TestMarshalUnmarshalRouteTables(t *testing.T) {
	_, subnet1cidr, _ := net.ParseCIDR("10.192.24.0/24")
	_, subnet2cidr, _ := net.ParseCIDR("10.20.24.0/24")
	_, subnet2ipv6, _ := net.ParseCIDR("2001:db8::/110")
	r := testResource("rt1", "routetable").prop(properties.ID, "rt1").prop(
		"Routes", []*Route{
			{Destination: subnet1cidr, DestinationPrefixListId: "toto", Targets: []*RouteTarget{{Type: InstanceTarget, Ref: "ref_1", Owner: "me"}, {Type: GatewayTarget, Ref: "ref_2"}}},
			{Destination: subnet2cidr, DestinationIPv6: subnet2ipv6, DestinationPrefixListId: "tata", Targets: []*RouteTarget{{Type: NetworkInterfaceTarget, Ref: "ref_3"}}},
		}).build()
	g := NewGraph()
	triples, err := r.marshalFullRDF()
	if err != nil {
		t.Fatal(err)
	}
	g.store.Add(triples...)
	rawRes := InitResource(r.Type(), r.Id())
	err = rawRes.unmarshalFullRdf(g.store.Snapshot())
	if err != nil {
		t.Fatal(err)
	}

	Routes(r.Properties()["Routes"].([]*Route)).Sort()
	Routes(rawRes.Properties()["Routes"].([]*Route)).Sort()

	if got, want := rawRes, r; !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%#v\nwant\n%#v\n", got, want)
	}
}

func TestMarshalUnmarshalGrants(t *testing.T) {
	r := testResource("bck1", "bucket").prop(properties.ID, "bck1").prop(
		"Grants", []*Grant{
			{Permission: "denied"},
			{Permission: "granted", Grantee: Grantee{GranteeID: "123", GranteeDisplayName: "John Smith", GranteeType: "user"}},
			{Permission: "other", Grantee: Grantee{GranteeID: "myid"}},
		}).build()
	g := NewGraph()
	triples, err := r.marshalFullRDF()
	if err != nil {
		t.Fatal(err)
	}
	g.store.Add(triples...)
	rawRes := InitResource(r.Type(), r.Id())

	snap := g.store.Snapshot()
	if err := rawRes.unmarshalFullRdf(snap); err != nil {
		t.Fatal(err)
	}

	Grants(r.Properties()["Grants"].([]*Grant)).Sort()
	Grants(rawRes.Properties()["Grants"].([]*Grant)).Sort()

	if got, want := rawRes, r; !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%#v\nwant\n%#v\n", got, want)
	}
}

func TestMarshalUnmarshalIPPermissions(t *testing.T) {
	expected := []*Resource{
		testResource("sg-1234", "securitygroup").prop("InboundRules", []*FirewallRule{
			{
				PortRange: PortRange{FromPort: int64(22), ToPort: int64(22), Any: false},
				Protocol:  "tcp",
				IPRanges:  []*net.IPNet{{IP: net.IPv4(10, 10, 0, 0), Mask: net.CIDRMask(16, 32)}},
			},
			{
				PortRange: PortRange{FromPort: int64(443), ToPort: int64(443), Any: false},
				Protocol:  "tcp",
				IPRanges:  []*net.IPNet{{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)}},
			},
		}).prop("OutboundRules", []*FirewallRule{
			{
				PortRange: PortRange{Any: true},
				Protocol:  "any",
				IPRanges:  []*net.IPNet{{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)}},
			},
		}).build(),
	}

	g := NewGraph()
	for _, r := range expected {
		triples, err := r.marshalFullRDF()
		if err != nil {
			t.Fatal(err)
		}
		g.store.Add(triples...)
	}

	res, err := g.GetAllResources("securitygroup")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), len(expected); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	if got, want := res[0].Id(), expected[0].Id(); got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	if got, want := res[0].Type(), expected[0].Type(); got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	if got, want := len(res[0].Properties()), len(expected[0].Properties()); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for k := range expected[0].Properties() {
		if got, want := fmt.Sprintf("%T", res[0].Properties()[k]), fmt.Sprintf("%T", expected[0].Properties()[k]); got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	}
}

func buildBenchmarkData() (data []*Resource) {
	_, localhost, _ := net.ParseCIDR("127.0.0.1/32")
	_, subnetcidr, _ := net.ParseCIDR("10.192.24.0/24")
	_, subnet2cidr, _ := net.ParseCIDR("10.20.24.0/24")
	_, subnet2ipv6, _ := net.ParseCIDR("2001:db8::/110")
	routes := []*Route{
		{Destination: subnetcidr, DestinationPrefixListId: "toto", Targets: []*RouteTarget{{Type: InstanceTarget, Ref: "ref_1", Owner: "me"}, {Type: GatewayTarget, Ref: "ref_2"}}},
		{Destination: subnet2cidr, DestinationIPv6: subnet2ipv6, DestinationPrefixListId: "tata", Targets: []*RouteTarget{{Type: NetworkInterfaceTarget, Ref: "ref_3"}}},
	}
	rules := []*FirewallRule{
		{PortRange: PortRange{FromPort: 80, ToPort: 80}, Protocol: "tcp", IPRanges: []*net.IPNet{localhost, subnetcidr}},
		{PortRange: PortRange{FromPort: 1, ToPort: 1024}, Protocol: "udp", IPRanges: []*net.IPNet{subnetcidr}},
	}
	for i := 0; i < 10; i++ {
		vpcId := fmt.Sprintf("vpc%d", i)
		data = append(data, vpcResource(vpcId).prop(properties.ID, vpcId).build())

		routeId := fmt.Sprintf("%s_route", vpcId)
		data = append(data, testResource(routeId, "routetable").prop(properties.ID, routeId).prop(properties.Vpc, vpcId).prop("Routes", routes).build())

		for j := 0; j < 10; j++ {
			subId := fmt.Sprintf("%ssub%d", vpcId, j)
			data = append(data, subResource(subId).prop(properties.ID, subId).prop(properties.Vpc, vpcId).prop(properties.Default, true).build())
			for k := 0; k < 10; k++ {
				instId := fmt.Sprintf("%sinst%d", subId, k)

				secGroup1Id := fmt.Sprintf("%s_securitygroup1", instId)
				data = append(data, sGrpResource(secGroup1Id).prop(properties.ID, secGroup1Id).prop("InboundRules", rules).prop("OutboundRules", rules).prop(properties.Vpc, vpcId).prop(properties.Launched, time.Now()).build())
				secGroup2Id := fmt.Sprintf("%s_securitygroup2", instId)
				data = append(data, sGrpResource(secGroup2Id).prop(properties.ID, secGroup2Id).prop("InboundRules", rules).prop("OutboundRules", rules).prop(properties.Vpc, vpcId).prop(properties.Launched, time.Now()).build())

				data = append(data, instResource(instId).prop(properties.ID, instId).prop("Name", instId+"name").prop(properties.Subnet, subId).prop(properties.Vpc, vpcId).prop(properties.Launched, time.Now()).prop(properties.SecurityGroups, []string{secGroup1Id, secGroup2Id}).build())
			}
		}
	}
	return
}

func BenchmarkRdfMarshaling(b *testing.B) {
	resources := buildBenchmarkData()
	b.ResetTimer()
	b.Run("full RDF", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, res := range resources {
				if _, err := res.marshalFullRDF(); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func buildBenchmarkFullRdfTriples(resources []*Resource) *Graph {
	graph := NewGraph()
	for _, res := range resources {
		if triples, err := res.marshalFullRDF(); err != nil {
			panic(err)
		} else {
			graph.store.Add(triples...)
		}
	}
	return graph
}

func BenchmarkRdfUnmarshaling(b *testing.B) {
	resources := buildBenchmarkData()
	graph := buildBenchmarkFullRdfTriples(resources)
	b.ResetTimer()
	f, _ := os.Create("./fullcpu.out")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	b.Run("full RDF", func(b *testing.B) {
		snap := graph.store.Snapshot()
		for i := 0; i < b.N; i++ {
			for _, r := range resources {
				rawRes := InitResource(r.Type(), r.Id())
				if err := rawRes.unmarshalFullRdf(snap); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
