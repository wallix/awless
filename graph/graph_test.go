package graph

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestGetResource(t *testing.T) {
	g := NewGraph()

	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Tags","Value":[{"Key":"Name","Value":"redis"}]}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Type","Value":"t2.micro"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"PublicIp","Value":"1.2.3.4"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"State","Value":{"Code": 16,"Name":"running"}}"^^type:text`))

	res, err := g.GetResource(Instance, "inst_1")
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
		{kind: Instance, id: "inst_1", Properties: Properties{"Id": "inst_1", "Name": "redis"}},
		{kind: Instance, id: "inst_2", Properties: Properties{"Id": "inst_2", "Name": "redis2"}},
		{kind: Instance, id: "inst_3", Properties: Properties{"Id": "inst_3", "Name": "redis3", "CreationDate": time}},
	}
	res, err := g.GetAllResources(Instance)
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
		{kind: SecurityGroup, id: "sg-1234", Properties: Properties{
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
	res, err := g.GetAllResources(SecurityGroup)
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
