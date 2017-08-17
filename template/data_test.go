package template

var DefsExample = map[string]Definition{
	"createsubnet": {
		Action:         "create",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []string{"cidr", "vpc"},
		ExtraParams:    []string{"availabilityzone", "name"},
	},
	"updatesubnet": {
		Action:         "update",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{"public"},
	},
	"createinstance": {
		Action:         "create",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []string{"image", "count", "count", "type", "subnet"},
		ExtraParams:    []string{"keypair", "ip", "userdata", "securitygroup", "lock", "name"},
	},
	"createkeypair": {
		Action:         "create",
		Entity:         "keypair",
		Api:            "ec2",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"createtag": {
		Action:         "create",
		Entity:         "tag",
		Api:            "ec2",
		RequiredParams: []string{"resource", "key", "value"},
		ExtraParams:    []string{},
	},
	"createloadbalancer": {
		Action:         "create",
		Entity:         "loadbalancer",
		Api:            "elbv2",
		RequiredParams: []string{"name", "subnets"},
		ExtraParams:    []string{"iptype", "scheme", "securitygroups"},
	},
}
