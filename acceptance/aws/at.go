package awsat

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template"
)

type ATBuilder struct {
	template    string
	cmdResult   *string
	expectCalls map[string]int
	expectInput map[string]interface{}
	mock        mock
	graph       *graph.Graph
}

func Template(template string) *ATBuilder {
	return &ATBuilder{template: template, expectCalls: make(map[string]int), expectInput: make(map[string]interface{})}
}

func (b *ATBuilder) ExpectCommandResult(key string) *ATBuilder {
	b.cmdResult = &key
	return b
}

func (b *ATBuilder) ExpectCalls(expects ...string) *ATBuilder {
	for _, expect := range expects {
		b.expectCalls[expect]++
	}
	return b
}

func (b *ATBuilder) ExpectInput(call string, input interface{}) *ATBuilder {
	b.expectInput[call] = input
	return b
}

func (b *ATBuilder) Graph(g *graph.Graph) *ATBuilder {
	b.graph = g
	return b
}

func (b *ATBuilder) Run(t *testing.T, l ...*logger.Logger) {
	t.Helper()
	b.mock.SetInputs(b.expectInput)
	b.mock.SetTesting(t)

	tpl, err := template.Parse(b.template)
	if err != nil {
		t.Fatal(err)
	}
	if b.graph == nil {
		b.graph = graph.NewGraph()
	}
	awsspec.CommandFactory = NewAcceptanceFactory(b.mock, b.graph, l...)

	env := template.NewEnv()
	env.Lookuper = func(tokens ...string) interface{} {
		return awsspec.CommandFactory.Build(strings.Join(tokens, ""))()
	}
	compiled, env, err := template.Compile(tpl, env, template.NewRunnerCompileMode)
	if err != nil {
		t.Fatal(err)
	}

	ran, err := compiled.Run(env)
	if err != nil {
		t.Fatal(err)
	}
	if ran.HasErrors() {
		for _, cmd := range ran.CommandNodesIterator() {
			if cmd.Err() != nil {
				t.Fatal(cmd.Err())
			}
		}
	}
	if len(b.expectCalls) > 0 {
		if got, want := b.mock.Calls(), b.expectCalls; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	}
	if b.cmdResult != nil {
		if got, want := fmt.Sprint(ran.CommandNodesIterator()[0].Result()), StringValue(b.cmdResult); got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func (b *ATBuilder) Mock(i mock) *ATBuilder {
	b.mock = i
	return b
}

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func String(v string) *string {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

func Int64AsIntValue(v *int64) int {
	if v != nil {
		return int(*v)
	}
	return 0
}

func Bool(v bool) *bool {
	return &v
}

func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}
