package graph

type rBuilder struct {
	id, typ string
	props   map[string]interface{}
}

func testResource(id, typ string) *rBuilder {
	return &rBuilder{id: id, typ: typ, props: map[string]interface{}{"ID": id}}
}

func instResource(id string) *rBuilder {
	return testResource(id, "instance")
}

func subResource(id string) *rBuilder {
	return testResource(id, "subnet")
}

func vpcResource(id string) *rBuilder {
	return testResource(id, "vpc")
}

func sGrpResource(id string) *rBuilder {
	return testResource(id, "securitygroup")
}

func (b *rBuilder) prop(key string, value interface{}) *rBuilder {
	b.props[key] = value
	return b
}

func (b *rBuilder) build() *Resource {
	return &Resource{id: b.id, kind: b.typ, properties: b.props, meta: make(map[string]interface{}), relations: make(map[string][]*Resource)}
}
