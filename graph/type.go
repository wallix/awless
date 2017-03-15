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
)

type FirewallRule struct {
	PortRange PortRange
	Protocol  string
	IPRanges  []*net.IPNet // IPv4 or IPv6 range
}

func (r *FirewallRule) String() string {
	return fmt.Sprintf("PortRange:%+v; Protocol:%s; IPRanges:%+v", r.PortRange, r.Protocol, r.IPRanges)
}

type PortRange struct {
	FromPort, ToPort int64
	Any              bool
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
	return fmt.Sprintf("Type:%+v; Ref:%s", t.Type, t.Ref)
}

type Route struct {
	Destination             *net.IPNet
	DestinationIPv6         *net.IPNet
	DestinationPrefixListId string
	Targets                 []*RouteTarget
}

func (r *Route) String() string {
	return fmt.Sprintf("Destination:%+v; DestinationIPv6:%+v; Targets:%+v", r.Destination, r.DestinationIPv6, r.Targets)
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

type ResourceType string

func (r ResourceType) String() string {
	return string(r)
}
