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
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type FirewallRules []*FirewallRule

func (rules FirewallRules) Sort() {
	for _, r := range rules {
		sort.Slice(r.IPRanges, func(i int, j int) bool {
			return r.IPRanges[i].String() < r.IPRanges[j].String()
		})
	}
	sort.Slice(rules, func(i int, j int) bool {
		return rules[i].String() < rules[j].String()
	})
}

type FirewallRule struct {
	PortRange PortRange    `predicate:"net:portRange"`
	Protocol  string       `predicate:"net:protocol"`
	IPRanges  []*net.IPNet `predicate:"net:cidr"` // IPv4 or IPv6 range
	Sources   []string     `predicate:"cloud:source"`
}

func (r *FirewallRule) Contains(ip string) bool {
	addr := net.ParseIP(ip)
	for _, n := range r.IPRanges {
		if n.Contains(addr) {
			return true
		}
	}
	return false
}

func (r *FirewallRule) String() string {
	return fmt.Sprintf("PortRange:%+v; Protocol:%s; IPRanges:%+v; Sources:%+v", r.PortRange, r.Protocol, r.IPRanges, r.Sources)
}

func (r *FirewallRule) marshalToTriples(id string) []tstore.Triple {
	var triples []tstore.Triple
	triples = append(triples, tstore.SubjPred(id, rdf.RdfType).Resource(rdf.NetFirewallRule))
	triples = append(triples, tstore.TriplesFromStruct(id, r)...)
	return triples
}

func (r *FirewallRule) unmarshalFromTriples(g tstore.RDFGraph, id string) error {
	portRangeTs := g.WithSubjPred(id, rdf.PortRange)
	ports, err := extractUniqueLiteralTextFromTriples(portRangeTs)
	if err != nil {
		return fmt.Errorf("unmarshal firewall rule: port range: %s", err)
	}
	pr, err := ParsePortRange(ports)
	if err != nil {
		return fmt.Errorf("unmarshal firewall rule: %s", err)
	}
	r.PortRange = pr

	protocolTs := g.WithSubjPred(id, rdf.Protocol)
	protocol, err := extractUniqueLiteralTextFromTriples(protocolTs)
	if err != nil {
		return fmt.Errorf("unmarshal firewall rule: protocol: %s", err)
	}
	r.Protocol = protocol

	cidrTs := g.WithSubjPred(id, rdf.CIDR)
	for _, cidrT := range cidrTs {
		cidrTxt, err := tstore.ParseString(cidrT.Object())
		if err != nil {
			return fmt.Errorf("unmarshal firewall rule: cidr: %s", err)
		}
		_, cidr, err := net.ParseCIDR(cidrTxt)
		if err != nil {
			return fmt.Errorf("unmarshal firewall rule: cidr: %s", err)
		}
		r.IPRanges = append(r.IPRanges, cidr)
	}

	sourceTs := g.WithSubjPred(id, rdf.Source)
	for _, sourceT := range sourceTs {
		source, err := tstore.ParseString(sourceT.Object())
		if err != nil {
			return fmt.Errorf("unmarshal firewall rule: source: %s", err)
		}
		r.Sources = append(r.Sources, source)
	}
	return nil
}

type PortRange struct {
	FromPort, ToPort int64
	Any              bool
}

func (p PortRange) Contains(port int64) bool {
	if p.Any {
		return true
	}

	from, to := p.FromPort, p.ToPort
	if from == port || to == port || (from < port && to > port) {
		return true
	}

	return false
}

func (p PortRange) String() string {
	switch {
	case p.Any:
		return ":"
	case p.FromPort == int64(-1):
		return fmt.Sprintf("%d:%[1]d", p.ToPort)
	case p.ToPort == int64(-1):
		return fmt.Sprintf("%d:%[1]d", p.FromPort)
	default:
		return fmt.Sprintf("%d:%d", p.FromPort, p.ToPort)
	}

}

func ParsePortRange(s string) (PortRange, error) {
	splits := strings.Split(s, ":")
	switch {
	case s == ":":
		return PortRange{Any: true}, nil
	case len(splits) == 2:
		from, err := strconv.Atoi(splits[0])
		if err != nil {
			return PortRange{}, err
		}
		to, err := strconv.Atoi(splits[1])
		if err != nil {
			return PortRange{}, err
		}
		return PortRange{FromPort: int64(from), ToPort: int64(to)}, nil
	default:
		return PortRange{}, fmt.Errorf("unexpected portrange: '%s'", s)
	}
}

type routeTargetType int

const (
	EgressOnlyInternetGatewayTarget routeTargetType = iota
	GatewayTarget
	InstanceTarget
	NatTarget
	NetworkInterfaceTarget
	VpcPeeringConnectionTarget
)

type RouteTarget struct {
	Type  routeTargetType
	Ref   string
	Owner string
}

func (t *RouteTarget) String() string {
	return fmt.Sprintf("%d|%s|%s", t.Type, t.Ref, t.Owner)
}

func ParseRouteTarget(s string) (*RouteTarget, error) {
	splits := strings.Split(s, "|")
	if len(splits) != 3 {
		return &RouteTarget{}, fmt.Errorf("unexpected route target: '%s'", s)
	}
	typ, err := strconv.Atoi(splits[0])
	if err != nil {
		return &RouteTarget{}, err
	}
	return &RouteTarget{Type: routeTargetType(typ), Ref: splits[1], Owner: splits[2]}, nil
}

type Routes []*Route

func (routes Routes) Sort() {
	for _, r := range routes {
		sort.Slice(r.Targets, func(i int, j int) bool {
			return r.Targets[i].String() < r.Targets[j].String()
		})
	}
	sort.Slice(routes, func(i int, j int) bool {
		return routes[i].String() < routes[j].String()
	})
}

type Route struct {
	Destination             *net.IPNet     `predicate:"net:cidr"`
	DestinationIPv6         *net.IPNet     `predicate:"net:cidrv6"`
	DestinationPrefixListId string         `predicate:"net:routeDestinationPrefixList"`
	Targets                 []*RouteTarget `predicate:"net:routeTargets"`
}

func (r *Route) String() string {
	return fmt.Sprintf("Destination:%+v; DestinationIPv6:%+v; DestinationPrefixListId:%s; Targets:%+v", r.Destination, r.DestinationIPv6, r.DestinationPrefixListId, r.Targets)
}

func (r *Route) marshalToTriples(id string) []tstore.Triple {
	var triples []tstore.Triple
	triples = append(triples, tstore.SubjPred(id, rdf.RdfType).Resource(rdf.NetRoute))
	triples = append(triples, tstore.TriplesFromStruct(id, r)...)
	return triples
}

func (r *Route) unmarshalFromTriples(g tstore.RDFGraph, id string) error {
	routeDestTs := g.WithSubjPred(id, rdf.CIDR)
	if len(routeDestTs) > 0 {
		dest, err := extractUniqueLiteralTextFromTriples(routeDestTs)
		if err != nil {
			return fmt.Errorf("unmarshal route: destination: %s", err)
		}
		_, r.Destination, err = net.ParseCIDR(dest)
		if err != nil {
			return fmt.Errorf("unmarshal route: destination: %s", err)
		}
	}

	routeDestv6Ts := g.WithSubjPred(id, rdf.CIDRv6)

	if len(routeDestv6Ts) > 0 {
		destv6, err := extractUniqueLiteralTextFromTriples(routeDestv6Ts)
		if err != nil {
			return fmt.Errorf("unmarshal route: destinationV6: %s", err)
		}
		_, r.DestinationIPv6, err = net.ParseCIDR(destv6)
		if err != nil {
			return fmt.Errorf("unmarshal route: destinationV6: %s", err)
		}
	}

	destPrefixTs := g.WithSubjPred(id, rdf.NetDestinationPrefixList)
	if len(destPrefixTs) > 0 {
		var err error
		r.DestinationPrefixListId, err = extractUniqueLiteralTextFromTriples(destPrefixTs)
		if err != nil {
			return fmt.Errorf("unmarshal route: destination prefix: %s", err)
		}
	}

	targetTs := g.WithSubjPred(id, rdf.NetRouteTargets)
	for _, targetT := range targetTs {
		litText, err := tstore.ParseString(targetT.Object())
		if err != nil {
			return err
		}
		target, err := ParseRouteTarget(litText)
		if err != nil {
			return fmt.Errorf("unmarshal route target: %s", err)
		}
		r.Targets = append(r.Targets, target)
	}
	return nil
}

type Grants []*Grant

func (grants Grants) Sort() {
	sort.Slice(grants, func(i int, j int) bool {
		return grants[i].String() < grants[j].String()
	})
}

type Grant struct {
	Permission string  `predicate:"cloud:permission"`
	Grantee    Grantee `predicate:"cloud:grantee" bnode:""`
}

type Grantee struct {
	GranteeID          string `predicate:"cloud:id"`
	GranteeDisplayName string `predicate:"cloud:name"`
	GranteeType        string `predicate:"cloud:granteeType"`
}

func (g *Grant) String() string {
	return fmt.Sprintf("Permission:%s; GranteeID:%s; GranteeDisplayName:%s; GranteeType:%s", g.Permission, g.Grantee.GranteeID, g.Grantee.GranteeDisplayName, g.Grantee.GranteeType)
}

func (g *Grant) marshalToTriples(id string) []tstore.Triple {
	var triples []tstore.Triple

	triples = append(triples, tstore.SubjPred(id, rdf.RdfType).Resource(rdf.Grant))
	triples = append(triples, tstore.TriplesFromStruct(id, g)...)

	return triples
}

func (g *Grant) unmarshalFromTriples(gph tstore.RDFGraph, id string) error {
	permissionTs := gph.WithSubjPred(id, rdf.Permission)
	var err error
	g.Permission, err = extractUniqueLiteralTextFromTriples(permissionTs)
	if err != nil {
		return fmt.Errorf("unmarshal grant: permission: %s", err)
	}
	granteeTs := gph.WithSubjPred(id, rdf.Grantee)
	if len(granteeTs) != 1 {
		return fmt.Errorf("unmarshal grant: expect 1 grantee got: %d", len(granteeTs))
	}
	granteeNode, ok := granteeTs[0].Object().Bnode()
	if !ok {
		return fmt.Errorf("unmarshal grant: grantee does not contain a resource identifier")
	}
	granteeIdTs := gph.WithSubjPred(granteeNode, rdf.ID)
	if len(granteeIdTs) > 0 {
		g.Grantee.GranteeID, err = extractUniqueLiteralTextFromTriples(granteeIdTs)
		if err != nil {
			return fmt.Errorf("unmarshal grant: grantee id: %s", err)
		}
	}
	granteeNameTs := gph.WithSubjPred(granteeNode, rdf.Name)
	if len(granteeNameTs) > 0 {
		g.Grantee.GranteeDisplayName, err = extractUniqueLiteralTextFromTriples(granteeNameTs)
		if err != nil {
			return fmt.Errorf("unmarshal grant: grantee name: %s", err)
		}
	}

	granteeTypeTs := gph.WithSubjPred(granteeNode, rdf.GranteeType)
	if len(granteeTypeTs) > 0 {
		g.Grantee.GranteeType, err = extractUniqueLiteralTextFromTriples(granteeTypeTs)
		if err != nil {
			return fmt.Errorf("unmarshal grant: grantee type: %s", err)
		}
	}

	return nil
}

type KeyValue struct {
	KeyName string `predicate:"cloud:keyName"`
	Value   string `predicate:"cloud:value"`
}

func (kv *KeyValue) String() string {
	return fmt.Sprintf("[Key:%s,Value:%s]", kv.KeyName, kv.Value)
}

func (kv *KeyValue) marshalToTriples(id string) []tstore.Triple {
	var triples []tstore.Triple

	triples = append(triples, tstore.SubjPred(id, rdf.RdfType).Resource(rdf.KeyValue))
	triples = append(triples, tstore.TriplesFromStruct(id, kv)...)

	return triples
}

func (kv *KeyValue) unmarshalFromTriples(gph tstore.RDFGraph, id string) error {
	var err error
	kv.KeyName, err = extractUniqueLiteralTextFromGraph(gph, id, rdf.KeyName)
	if err != nil {
		return fmt.Errorf("unmarshal keyvalue: key name: %s", err)
	}
	kv.Value, err = extractUniqueLiteralTextFromGraph(gph, id, rdf.Value)
	if err != nil {
		return fmt.Errorf("unmarshal keyvalue: val name: %s", err)
	}
	return nil
}

type DistributionOrigin struct {
	ID         string `predicate:"cloud:id"`
	PublicDNS  string `predicate:"cloud:publicDNS"`
	PathPrefix string `predicate:"cloud:pathPrefix"`
	OriginType string `predicate:"cloud:type"`
	Config     string `predicate:"cloud:config"`
}

func (o *DistributionOrigin) String() string {
	var elems []string
	elems = append(elems, "ID:"+o.ID)
	if o.PublicDNS != "" {
		elems = append(elems, "PublicDNS:"+o.PublicDNS)
	}
	if o.PathPrefix != "" {
		elems = append(elems, "PathPrefix:"+o.PathPrefix)
	}
	if o.OriginType != "" {
		elems = append(elems, "Type:"+o.OriginType)
	}
	if o.Config != "" {
		elems = append(elems, "Config:"+o.Config)
	}
	return fmt.Sprintf("[%s]", strings.Join(elems, ","))
}

func (o *DistributionOrigin) marshalToTriples(id string) []tstore.Triple {
	var triples []tstore.Triple

	triples = append(triples, tstore.SubjPred(id, rdf.RdfType).Resource(rdf.DistributionOrigin))
	triples = append(triples, tstore.TriplesFromStruct(id, o)...)

	return triples
}

func (o *DistributionOrigin) unmarshalFromTriples(gph tstore.RDFGraph, id string) error {
	var err error
	o.ID, err = extractUniqueLiteralTextFromGraph(gph, id, rdf.ID)
	if err != nil {
		return fmt.Errorf("unmarshal DistributionOrigin: extract id: %s", err)
	}
	o.PublicDNS, err = extractUniqueLiteralTextFromGraph(gph, id, rdf.PublicDNS)
	if err != nil {
		return fmt.Errorf("unmarshal DistributionOrigin: extract PublicDNS: %s", err)
	}
	o.PathPrefix, err = extractUniqueLiteralTextFromGraph(gph, id, rdf.PathPrefix)
	if err != nil {
		return fmt.Errorf("unmarshal DistributionOrigin: extract PathPrefix: %s", err)
	}
	o.OriginType, err = extractUniqueLiteralTextFromGraph(gph, id, rdf.Type)
	if err != nil {
		return fmt.Errorf("unmarshal DistributionOrigin: extract Type: %s", err)
	}
	o.Config, err = extractUniqueLiteralTextFromGraph(gph, id, rdf.Config)
	if err != nil {
		return fmt.Errorf("unmarshal DistributionOrigin: extract Config: %s", err)
	}
	return nil
}

type Policy struct {
	Version    string             `json:",omitempty"`
	ID         string             `json:"Id,omitempty"`
	Statements compositeStatement `json:"Statement,omitempty"`
}

type PolicyStatement struct {
	ID           string              `json:"Sid,omitempty"`
	Principal    *StatementPrincipal `json:",omitempty"`
	NotPrincipal *StatementPrincipal `json:",omitempty"`
	Effect       string              `json:",omitempty"`
	Actions      compositeString     `json:"Action,omitempty"`
	NotActions   compositeString     `json:"NotAction,omitempty"`
	Resources    compositeString     `json:"Resource,omitempty"`
	NotResources compositeString     `json:"NotResource,omitempty"`
	Condition    interface{}         `json:",omitempty"`
}

type StatementPrincipal struct {
	AWS       compositeString `json:",omitempty"`
	Service   compositeString `json:",omitempty"`
	Federated compositeString `json:",omitempty"`
}

// To support AWS JSON for Policy in which Principal, Action,... can be either string or slice of string
type compositeString []string

func (c *compositeString) UnmarshalJSON(data []byte) (err error) {
	var str string
	if err = json.Unmarshal(data, &str); err == nil {
		*c = []string{str}
		return
	}

	var slice []string
	if err = json.Unmarshal(data, &slice); err != nil {
		return
	}
	*c = slice
	return
}

// To support AWS JSON for Policy in which Statement can be either Statement or slice of Statement
type compositeStatement []*PolicyStatement

func (c *compositeStatement) UnmarshalJSON(data []byte) (err error) {
	var statement *PolicyStatement
	if err = json.Unmarshal(data, &statement); err == nil {
		*c = []*PolicyStatement{statement}
		return
	}

	var slice []*PolicyStatement
	if err = json.Unmarshal(data, &slice); err != nil {
		return
	}
	*c = slice
	return
}

// To support AWS JSON for Policy in which a principal can be either a JSON object, either "*"
func (c *StatementPrincipal) UnmarshalJSON(data []byte) (err error) {
	var wildCardString string
	if err = json.Unmarshal(data, &wildCardString); err == nil {
		if wildCardString == "*" {
			c.AWS = []string{"*"} // according to doc, "Principal":"*" is equivalent to "Principal":{"AWS":"*"}
			return
		} else {
			return fmt.Errorf("unmarshaling policy: a principal string can only contain '*', but got %s", wildCardString)
		}
	}

	type aliasPrincipal struct {
		AWS       compositeString `json:",omitempty"`
		Service   compositeString `json:",omitempty"`
		Federated compositeString `json:",omitempty"`
	}
	var principal *aliasPrincipal
	if err = json.Unmarshal(data, &principal); err != nil {
		return
	}
	*c = StatementPrincipal(*principal)
	return
}

func extractUniqueLiteralTextFromGraph(gph tstore.RDFGraph, subj, pred string) (string, error) {
	ts := gph.WithSubjPred(subj, pred)
	if len(ts) != 1 {
		return "", fmt.Errorf("%s,%s: expect 1 triple got: %d", subj, pred, len(ts))
	}
	return extractUniqueLiteralTextFromTriples(ts)
}
