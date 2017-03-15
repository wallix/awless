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
	"reflect"
	"testing"
	"time"
)

func TestAddGraphRelation(t *testing.T) {

	t.Run("Add parent", func(t *testing.T) {
		g := NewGraph()
		g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text`))

		res, err := g.GetResource("instance", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		g.AddParentRelation(InitResource("subnet_1", "subnet"), res)

		exp := `/instance<inst_1>	"has_type"@[]	"/instance"^^type:text
/subnet<subnet_1>	"parent_of"@[]	/instance<inst_1>`

		if got, want := g.MustMarshal(), exp; got != want {
			t.Fatalf("got\n%q\nwant\n%q\n", got, want)
		}
	})

	t.Run("Add applies on", func(t *testing.T) {
		g := NewGraph()
		g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text`))

		res, err := g.GetResource("instance", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		g.AddAppliesOnRelation(InitResource("subnet_1", "subnet"), res)

		exp := `/instance<inst_1>	"has_type"@[]	"/instance"^^type:text
/subnet<subnet_1>	"applies_on"@[]	/instance<inst_1>`

		if got, want := g.MustMarshal(), exp; got != want {
			t.Fatalf("got\n%q\nwant\n%q\n", got, want)
		}
	})
}

func TestGetResource(t *testing.T) {
	g := NewGraph()

	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Tags","Value":[{"Key":"Name","Value":"redis"}]}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Type","Value":"t2.micro"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"PublicIp","Value":"1.2.3.4"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"State","Value":{"Code": 16,"Name":"running"}}"^^type:text`))

	res, err := g.GetResource("instance", "inst_1")
	if err != nil {
		t.Fatal(err)
	}

	expected := Properties{"Id": "inst_1", "Type": "t2.micro", "PublicIp": "1.2.3.4",
		"State": map[string]interface{}{"Code": float64(16), "Name": "running"},
		"Tags": []interface{}{
			map[string]interface{}{"Key": "Name", "Value": "redis"},
		},
	}

	if got, want := res.Properties, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got \n%#v\n\nwant \n%#v\n", got, want)
	}
}

func TestFindResources(t *testing.T) {
	t.Parallel()
	g := NewGraph()

	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
  /instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /subnet<sub_1>  "has_type"@[] "/subnet"^^type:text
  /subnet<sub_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text`))

	t.Run("FindResource", func(t *testing.T) {
		t.Parallel()
		res, err := g.FindResource("inst_1")
		if err != nil {
			t.Fatal(err)
		}
		if got, want := res.Properties["Name"], "redis"; got != want {
			t.Fatalf("got %s want %s", got, want)
		}

		if res, err = g.FindResource("none"); err != nil {
			t.Fatal(err)
		}
		if res != nil {
			t.Fatalf("expected nil got %v", res)
		}

		if res, err = g.FindResource("sub_1"); err != nil {
			t.Fatal(err)
		}
		if got, want := res.Type(), "subnet"; got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	})
	t.Run("FindResourcesByProperty", func(t *testing.T) {
		t.Parallel()
		res, err := g.FindResourcesByProperty("Id", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		expected := []*Resource{
			{id: "inst_1", kind: "instance", Properties: map[string]interface{}{"Id": interface{}("inst_1"), "Name": interface{}("redis")}, Meta: make(Properties)},
		}
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		if got, want := res[0], expected[0]; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %+v want %+v", got, want)
		}
		res, err = g.FindResourcesByProperty("Name", "redis")
		if err != nil {
			t.Fatal(err)
		}
		expected = []*Resource{
			{id: "inst_1", kind: "instance", Properties: map[string]interface{}{"Id": "inst_1", "Name": "redis"}, Meta: make(Properties)},
			{id: "sub_1", kind: "subnet", Properties: map[string]interface{}{"Name": "redis"}, Meta: make(Properties)},
		}
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		if res[0].Id() == expected[0].Id() {
			if got, want := res[0], expected[0]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
			if got, want := res[1], expected[1]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
		} else {
			if got, want := res[0], expected[1]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
			if got, want := res[1], expected[0]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
		}
	})
}

func TestGetAllResources(t *testing.T) {
	g := NewGraph()

	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
  /instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Name","Value":"redis2"}"^^type:text
  /instance<inst_3>  "has_type"@[] "/instance"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Id","Value":"inst_3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Name","Value":"redis3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"CreationDate","Value":"2017-01-10T16:47:18Z"}"^^type:text
  /instance<subnet>  "has_type"@[] "/subnet"^^type:text
  /instance<subnet>  "property"@[] "{"Key":"Id","Value":"my subnet"}"^^type:text`))

	time, _ := time.Parse(time.RFC3339, "2017-01-10T16:47:18Z")

	expected := []*Resource{
		{kind: "instance", id: "inst_1", Properties: Properties{"Id": "inst_1", "Name": "redis"}},
		{kind: "instance", id: "inst_2", Properties: Properties{"Id": "inst_2", "Name": "redis2"}},
		{kind: "instance", id: "inst_3", Properties: Properties{"Id": "inst_3", "Name": "redis3", "CreationDate": time}},
	}
	res, err := g.GetAllResources("instance")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), len(expected); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for _, r := range expected {
		found := false
		for _, r2 := range res {
			if r2.kind == r.kind && r2.id == r.id && reflect.DeepEqual(r2.Properties, r.Properties) {
				found = true
			}
		}
		if !found {
			t.Fatalf("%+v not found", r)
		}
	}
}

func TestLoadIpPermissions(t *testing.T) {
	g := NewGraph()
	g.Unmarshal([]byte(`/securitygroup<sg-1234>	"has_type"@[]	"/securitygroup"^^type:text
/securitygroup<sg-1234>	"property"@[]	"{"Key":"Id","Value":"sg-1234"}"^^type:text
/securitygroup<sg-1234>	"property"@[]	"{"Key":"InboundRules","Value":[{"PortRange":{"FromPort":22,"ToPort":22,"Any":false},"Protocol":"tcp","IPRanges":[{"IP":"10.10.0.0","Mask":"//8AAA=="}]},{"PortRange":{"FromPort":443,"ToPort":443,"Any":false},"Protocol":"tcp","IPRanges":[{"IP":"0.0.0.0","Mask":"AAAAAA=="}]}]}"^^type:text
/securitygroup<sg-1234>	"property"@[]	"{"Key":"OutboundRules","Value":[{"PortRange":{"FromPort":0,"ToPort":0,"Any":true},"Protocol":"any","IPRanges":[{"IP":"0.0.0.0","Mask":"AAAAAA=="}]}]}"^^type:text`))
	expected := []*Resource{
		{kind: "securitygroup", id: "sg-1234", Properties: Properties{
			"Id": "sg-1234",
			"InboundRules": []*FirewallRule{
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
			},
			"OutboundRules": []*FirewallRule{
				{
					PortRange: PortRange{Any: true},
					Protocol:  "any",
					IPRanges:  []*net.IPNet{{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)}},
				},
			},
		},
		},
	}
	res, err := g.GetAllResources("securitygroup")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), len(expected); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	if got, want := res[0].id, expected[0].id; got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	if got, want := res[0].kind, expected[0].kind; got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	if got, want := len(res[0].Properties), len(expected[0].Properties); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for k := range expected[0].Properties {
		if got, want := fmt.Sprintf("%T", res[0].Properties[k]), fmt.Sprintf("%T", expected[0].Properties[k]); got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	}
}
