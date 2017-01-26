// DO NOT EDIT
// This file was automatically generated with go generate
package aws

var AWSTemplatesDefinitions = map[string]string{
	"createvpc": "create vpc cidr={ vpc.cidr } ",
	"deletevpc": "delete vpc id={ vpc.id } ",
	"createsubnet": "create subnet cidr={ subnet.cidr } vpc={ subnet.vpc } ",
	"updatesubnet": "update subnet id={ subnet.id } ",
	"deletesubnet": "delete subnet id={ subnet.id } ",
	"createinstance": "create instance image={ instance.image } type={ instance.type } count={ instance.count } count={ instance.count } subnet={ instance.subnet }  name={ instance.name }",
	"deleteinstance": "delete instance id={ instance.id } ",
	"startinstance": "start instance id={ instance.id } ",
	"stopinstance": "stop instance id={ instance.id } ",
	"createtags": "create tags resource={ tags.resource } ",
	"createkeypair": "create keypair name={ keypair.name } ",
	"deletekeypair": "delete keypair name={ keypair.name } ",
	"createuser": "create user name={ user.name } ",
	"deleteuser": "delete user name={ user.name } ",
	"creategroup": "create group name={ group.name } ",
	"deletegroup": "delete group name={ group.name } ",
}
