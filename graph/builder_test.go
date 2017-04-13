package graph

type rBuilder struct {
	id, typ string
	props   map[string]interface{}
}

func testResource(id, typ string) *rBuilder {
	return &rBuilder{id: id, typ: typ, props: make(map[string]interface{})}
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
	return &Resource{id: b.id, kind: b.typ, Properties: b.props, Meta: make(map[string]interface{})}
}
