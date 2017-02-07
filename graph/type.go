package graph

import (
	"fmt"
	"net"
)

const (
	Region ResourceType = "region"
	//infra
	Vpc             ResourceType = "vpc"
	Subnet          ResourceType = "subnet"
	Image           ResourceType = "image"
	SecurityGroup   ResourceType = "securitygroup"
	Keypair         ResourceType = "keypair"
	Volume          ResourceType = "volume"
	Instance        ResourceType = "instance"
	InternetGateway ResourceType = "internetgateway"
	RouteTable      ResourceType = "routetable"

	//access
	User   ResourceType = "user"
	Role   ResourceType = "role"
	Group  ResourceType = "group"
	Policy ResourceType = "policy"

	//s3
	Bucket ResourceType = "bucket"
	Object ResourceType = "storageobject"
	Acl    ResourceType = "storageacl"
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

func (r ResourceType) ToRDFString() string {
	return "/" + r.String()
}
