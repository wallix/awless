package template

var DefsExample = map[string]Definition{
	"createinstance": {
		Action:         "create",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []string{"image", "count", "count", "type", "subnet"},
		ExtraParams:    []string{"key", "ip", "userdata", "secgroup", "lock"},
	},
	"createkey": {
		Action:         "create",
		Entity:         "key",
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
}
