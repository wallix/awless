// DO NOT EDIT
// This file was automatically generated with go generate
package aws

var AWSTemplates = map[string]string{
	"createvpc":      "create vpc cidr={ vpc.cidr }",
	"deletevpc":      "delete vpc id={ vpc.id }",
	"createsubnet":   "create subnet cidr={ subnet.cidr } vpc={ subnet.vpc }",
	"deletesubnet":   "delete subnet id={ subnet.id }",
	"createinstance": "create instance image={ instance.image } type={ instance.type } count={ instance.count } count={ instance.count } subnet={ instance.subnet }",
	"deleteinstance": "delete instance id={ instance.id }",
}
