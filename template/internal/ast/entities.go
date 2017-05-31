package ast

type Entity string

var entities = map[Entity]struct{}{
	"none": {},

	"accesskey":           {},
	"alarm":               {},
	"scalinggroup":        {},
	"bucket":              {},
	"database":            {},
	"distribution":        {},
	"dbsubnetgroup":       {},
	"elasticip":           {},
	"function":            {},
	"group":               {},
	"instance":            {},
	"image":               {},
	"internetgateway":     {},
	"instanceprofile":     {},
	"keypair":             {},
	"launchconfiguration": {},
	"listener":            {},
	"loadbalancer":        {},
	"loginprofile":        {},
	"policy":              {},
	"queue":               {},
	"record":              {},
	"role":                {},
	"route":               {},
	"routetable":          {},
	"s3object":            {},
	"scalingpolicy":       {},
	"securitygroup":       {},
	"snapshot":            {},
	"stack":               {},
	"subnet":              {},
	"subscription":        {},
	"tag":                 {},
	"targetgroup":         {},
	"topic":               {},
	"user":                {},
	"volume":              {},
	"vpc":                 {},
	"zone":                {},
}

func IsInvalidEntity(s string) bool {
	_, ok := entities[Entity(s)]
	return !ok
}
