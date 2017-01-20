package rdf

import (
	"strings"

	"github.com/google/badwolf/triple/node"
)

type ResourceType string

const (
	Region        ResourceType = "region"
	Vpc           ResourceType = "vpc"
	Subnet        ResourceType = "subnet"
	SecurityGroup ResourceType = "securitygroup"
	Instance      ResourceType = "instance"
	User          ResourceType = "user"
	Role          ResourceType = "role"
	Group         ResourceType = "group"
	Policy        ResourceType = "policy"
)

func NewResourceType(t *node.Type) ResourceType {
	if !strings.HasPrefix(t.String(), "/") {
		panic("invalid resource type:" + t.String())
	}
	return ResourceType(strings.Split(t.String(), "/")[1])
}

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
	if strings.HasSuffix(singular, "y") {
		return strings.TrimSuffix(singular, "y") + "ies"
	}
	return singular + "s"
}
