package graph

import (
	"fmt"
	"net"
	"strings"
)

const (
	Region          ResourceType = "region"
	Vpc             ResourceType = "vpc"
	Subnet          ResourceType = "subnet"
	Image           ResourceType = "image"
	SecurityGroup   ResourceType = "securitygroup"
	Keypair         ResourceType = "keypair"
	Volume          ResourceType = "volume"
	Instance        ResourceType = "instance"
	InternetGateway ResourceType = "internetgateway"
	RouteTable      ResourceType = "routetable"
	User            ResourceType = "user"
	Role            ResourceType = "role"
	Group           ResourceType = "group"
	Policy          ResourceType = "policy"
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

type Route struct {
	Destination *net.IPNet
	TargetType  routeTargetType
	Target      string
}

func (r *Route) String() string {
	return fmt.Sprintf("Destination:%+v; TargetType:%v; Target:%s", r.Destination, r.TargetType, r.Target)
}

type ResourceType string

func (r ResourceType) String() string {
	return string(r)
}

func (r ResourceType) ToRDFString() string {
	return "/" + r.String()
}

func (r ResourceType) PluralString() string {
	return pluralize(r.String())
}

func pluralize(singular string) string {
	if strings.HasSuffix(singular, "cy") {
		return strings.TrimSuffix(singular, "y") + "ies"
	}
	return singular + "s"
}
