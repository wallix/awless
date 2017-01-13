package rdf

import "strings"

type ResourceType int

const (
	Region ResourceType = iota
	Vpc
	Subnet
	Instance
	User
	Role
	Group
	Policy
)

func NewResourceTypeFromRdfType(str string) ResourceType {
	switch str {
	case "/region":
		return Region
	case "/vpc":
		return Vpc
	case "/subnet":
		return Subnet
	case "/instance":
		return Instance
	case "/user":
		return User
	case "/role":
		return Role
	case "/group":
		return Group
	case "/policy":
		return Policy
	default:
		panic("invalid resource type:" + str)
	}
}

func (r ResourceType) String() string {
	switch r {
	case Region:
		return "region"
	case Vpc:
		return "vpc"
	case Subnet:
		return "subnet"
	case Instance:
		return "instance"
	case User:
		return "user"
	case Role:
		return "role"
	case Group:
		return "group"
	case Policy:
		return "policy"
	default:
		panic("invalid resource type")
	}
}

func (r ResourceType) ToRDFType() string {
	return "/" + r.String()
}

func (r ResourceType) PluralString() string {
	return pluralize(r.String())
}

func pluralize(singular string) string {
	if strings.HasSuffix(singular, "y") {
		return strings.TrimSuffix(singular, "y") + "ies"
	}
	return singular + "s"
}
