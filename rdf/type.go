package rdf

import "strings"

const (
	REGION   = "/region"
	VPC      = "/vpc"
	SUBNET   = "/subnet"
	INSTANCE = "/instance"
	USER     = "/user"
	ROLE     = "/role"
	GROUP    = "/group"
	POLICY   = "/policy"
)

func ToRDFType(resourceType string) string {
	return "/" + strings.ToLower(resourceType)
}

func ToResourceType(rdfType string) string {
	if rdfType[0] == '/' {
		return rdfType[1:]
	}
	return rdfType
}
