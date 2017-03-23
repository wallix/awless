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
	"sort"
	"strconv"
	"strings"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	cloudrdf "github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/graph/internal/rdf"
)

type triplesMarshaler interface {
	marshalToTriples(id string) []*triple.Triple
	unmarshalFromTriples(gph *rdf.Graph, node *node.Node) error
}

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
	PortRange PortRange
	Protocol  string
	IPRanges  []*net.IPNet // IPv4 or IPv6 range
}

func (r *FirewallRule) String() string {
	return fmt.Sprintf("PortRange:%+v; Protocol:%s; IPRanges:%+v", r.PortRange, r.Protocol, r.IPRanges)
}

func (r *FirewallRule) marshalToTriples(id string) []*triple.Triple {
	var triples []*triple.Triple

	triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.PortRange).Literal(r.PortRange.String()))
	triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.Protocol).Literal(r.Protocol))
	for _, cidr := range r.IPRanges {
		triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.CIDR).Literal(cidr.String()))
	}
	return triples
}

func (r *FirewallRule) unmarshalFromTriples(g *rdf.Graph, node *node.Node) error {
	portRangeTs, err := g.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.PortRange))
	if err != nil {
		return err
	}
	ports, err := extractUniqueLiteralTextFromTriples(portRangeTs)
	if err != nil {
		return fmt.Errorf("unmarshal firewall rule: port range: %s", err)
	}
	pr, err := ParsePortRange(ports)
	if err != nil {
		return fmt.Errorf("unmarshal firewall rule: %s", err)
	}
	r.PortRange = pr

	protocolTs, err := g.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.Protocol))
	if err != nil {
		return err
	}
	protocol, err := extractUniqueLiteralTextFromTriples(protocolTs)
	if err != nil {
		return fmt.Errorf("unmarshal firewall rule: protocol: %s", err)
	}
	r.Protocol = protocol

	cidrTs, err := g.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.CIDR))
	for _, cidrT := range cidrTs {
		lit, err := cidrT.Object().Literal()
		if err != nil {
			return fmt.Errorf("unmarshal firewall rule: cidr: %s", err)
		}
		cidrTxt, err := lit.Text()
		if err != nil {
			return fmt.Errorf("unmarshal firewall rule: cidr: %s", err)
		}
		_, cidr, err := net.ParseCIDR(cidrTxt)
		if err != nil {
			return fmt.Errorf("unmarshal firewall rule: cidr: %s", err)
		}
		r.IPRanges = append(r.IPRanges, cidr)
	}
	return nil
}

type PortRange struct {
	FromPort, ToPort int64
	Any              bool
}

func (p PortRange) String() string {
	switch {
	case p.Any:
		return ":"
	case p.FromPort == int64(0):
		return fmt.Sprintf("%d:%[1]d", p.ToPort)
	case p.ToPort == int64(0):
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
	Destination             *net.IPNet
	DestinationIPv6         *net.IPNet
	DestinationPrefixListId string
	Targets                 []*RouteTarget
}

func (r *Route) String() string {
	return fmt.Sprintf("Destination:%+v; DestinationIPv6:%+v; DestinationPrefixListId:%s; Targets:%+v", r.Destination, r.DestinationIPv6, r.DestinationPrefixListId, r.Targets)
}

func (r *Route) marshalToTriples(id string) []*triple.Triple {
	var triples []*triple.Triple

	if r.Destination != nil {
		triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.CIDR).Literal(r.Destination.String()))
	}
	if r.DestinationIPv6 != nil {
		triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.CIDRv6).Literal(r.DestinationIPv6.String()))
	}
	triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.NetDestinationPrefixList).Literal(r.DestinationPrefixListId))
	for _, t := range r.Targets {
		triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.NetRouteTargets).Literal(t.String()))
	}
	return triples
}

func (r *Route) unmarshalFromTriples(g *rdf.Graph, node *node.Node) error {
	routeDestTs, err := g.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.CIDR))
	if err != nil {
		return err
	}
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

	routeDestv6Ts, err := g.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.CIDRv6))
	if err != nil {
		return err
	}
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

	destPrefixTs, err := g.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.NetDestinationPrefixList))
	if err != nil {
		return err
	}
	if len(destPrefixTs) > 0 {
		r.DestinationPrefixListId, err = extractUniqueLiteralTextFromTriples(destPrefixTs)
		if err != nil {
			return fmt.Errorf("unmarshal route: destination prefix: %s", err)
		}
	}

	targetTs, err := g.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.NetRouteTargets))
	if err != nil {
		return err
	}
	for _, targetT := range targetTs {
		lit, err := targetT.Object().Literal()
		if err != nil {
			return err
		}
		litText, err := lit.Text()
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
	Permission,
	GranteeID,
	GranteeDisplayName,
	GranteeType string
}

func (g *Grant) String() string {
	return fmt.Sprintf("Permission:%s; GranteeID:%s; GranteeDisplayName:%s; GranteeType:%s", g.Permission, g.GranteeID, g.GranteeDisplayName, g.GranteeType)
}

func (g *Grant) marshalToTriples(id string) []*triple.Triple {
	var triples []*triple.Triple

	triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.Permission).Literal(g.Permission))

	granteeId := randomRdfId()
	triples = append(triples, rdf.Subject(id).Predicate(cloudrdf.Grantee).Object(granteeId))
	if g.GranteeID != "" {
		triples = append(triples, rdf.Subject(granteeId).Predicate(cloudrdf.ID).Literal(g.GranteeID))
	}
	if g.GranteeDisplayName != "" {
		triples = append(triples, rdf.Subject(granteeId).Predicate(cloudrdf.Name).Literal(g.GranteeDisplayName))
	}
	if g.GranteeType != "" {
		triples = append(triples, rdf.Subject(granteeId).Predicate(cloudrdf.RdfType).Literal(g.GranteeType))
	}
	return triples
}

func (g *Grant) unmarshalFromTriples(gph *rdf.Graph, node *node.Node) error {
	permissionTs, err := gph.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.Permission))
	if err != nil {
		return err
	}
	g.Permission, err = extractUniqueLiteralTextFromTriples(permissionTs)
	if err != nil {
		return fmt.Errorf("unmarshal grant: permission: %s", err)
	}
	granteeTs, err := gph.TriplesForSubjectPredicate(node, rdf.MustBuildPredicate(cloudrdf.Grantee))
	if err != nil {
		return err
	}
	if len(granteeTs) != 1 {
		return fmt.Errorf("unmarshal grant: expect 1 grantee got: %d", len(granteeTs))
	}
	granteeNode, err := granteeTs[0].Object().Node()
	if err != nil {
		return err
	}
	granteeIdTs, err := gph.TriplesForSubjectPredicate(granteeNode, rdf.MustBuildPredicate(cloudrdf.ID))
	if err != nil {
		return err
	}
	if len(granteeIdTs) > 0 {
		g.GranteeID, err = extractUniqueLiteralTextFromTriples(granteeIdTs)
		if err != nil {
			return fmt.Errorf("unmarshal grant: grantee id: %s", err)
		}
	}
	granteeNameTs, err := gph.TriplesForSubjectPredicate(granteeNode, rdf.MustBuildPredicate(cloudrdf.Name))
	if err != nil {
		return err
	}
	if len(granteeNameTs) > 0 {
		g.GranteeDisplayName, err = extractUniqueLiteralTextFromTriples(granteeNameTs)
		if err != nil {
			return fmt.Errorf("unmarshal grant: grantee name: %s", err)
		}
	}

	granteeTypeTs, err := gph.TriplesForSubjectPredicate(granteeNode, rdf.MustBuildPredicate(cloudrdf.RdfType))
	if err != nil {
		return err
	}
	if len(granteeTypeTs) > 0 {
		g.GranteeType, err = extractUniqueLiteralTextFromTriples(granteeTypeTs)
		if err != nil {
			return fmt.Errorf("unmarshal grant: grantee type: %s", err)
		}
	}

	return nil
}
