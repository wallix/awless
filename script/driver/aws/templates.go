// DO NOT EDIT
// This file was automatically generated with go generate
package aws

var AWSTemplates = map[string]string{
	"createvpc":      "create vpc cidr={ vpc_cidr }",
	"deletevpc":      "delete vpc id={ vpc_id }",
	"createsubnet":   "create subnet cidr={ subnet_cidr } vpc={ subnet_vpc }",
	"deletesubnet":   "delete subnet id={ subnet_id }",
	"createinstance": "create instance base={ instance_base } type={ instance_type } count={ instance_count } count={ instance_count } subnet={ instance_subnet }",
	"deleteinstance": "delete instance id={ instance_id }",
}
