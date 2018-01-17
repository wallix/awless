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
	template     string
	cmdResult    *string
	expectCalls  map[string]int
	expectInput  map[string]interface{}
	ignoredInput map[string]struct{}
	fillers      map[string]string
	expectRevert string
	mock         mock
	graph        *graph.Graph
}

func Template(template string) *ATBuilder {
	return &ATBuilder{template: template,
		expectCalls:  make(map[string]int),
		expectInput:  make(map[string]interface{}),
		ignoredInput: make(map[string]struct{}),
	}
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

func (b *ATBuilder) IgnoreInput(calls ...string) *ATBuilder {
	for _, call := range calls {
		b.ignoredInput[call] = struct{}{}
	}
	return b
}

func (b *ATBuilder) Graph(g *graph.Graph) *ATBuilder {
	b.graph = g
	return b
}

func (b *ATBuilder) Mock(i mock) *ATBuilder {
	b.mock = i
	return b
}

func (b *ATBuilder) Fillers(fillers map[string]string) *ATBuilder {
	b.fillers = fillers
	return b
}

func (b *ATBuilder) ExpectRevert(revert string) *ATBuilder {
	b.expectRevert = revert
	return b
}

func (b *ATBuilder) Run(t *testing.T, l ...*logger.Logger) {
	t.Helper()
	b.mock.SetInputs(b.expectInput)
	b.mock.SetIgnored(b.ignoredInput)
	b.mock.SetTesting(t)

	tpl, err := template.Parse(b.template)
	if err != nil {
		t.Fatal(err)
	}
	if b.graph == nil {
		b.graph = graph.NewGraph()
	}
	awsspec.CommandFactory = NewAcceptanceFactory(b.mock, b.graph, l...)

	cenv := template.NewEnv().WithLookupCommandFunc(func(tokens ...string) interface{} {
		return awsspec.CommandFactory.Build(strings.Join(tokens, ""))()
	}).WithMissingHolesFunc(func(key string, paramPaths []string, isOptional bool) string {
		return b.fillers[key]
	}).Build()
	compiled, cenv, err := template.Compile(tpl, cenv, template.NewRunnerCompileMode)
	if err != nil {
		t.Fatal(err)
	}

	ran, err := compiled.Run(template.NewRunEnv(cenv))
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
	if b.expectRevert != "" {
		revert, err := ran.Revert()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := revert.String(), b.expectRevert; got != want {
			t.Fatalf("got\n%s\nwant\n%s", got, want)
		}
	}
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
