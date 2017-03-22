package template

import (
	"reflect"
	"strings"
	"testing"
)

func TestCheckReferencesDeclarationPass(t *testing.T) {
	env := NewEnv()
	tpl := MustParse(`
	sub = create subnet
	inst = create instance subnet=$sub
	create instance subnet=$inst_2
	`)

	_, _, err := checkReferencesDeclaration(tpl, env)
	if err == nil ||
		!strings.Contains(err.Error(), "undefined") ||
		!strings.Contains(err.Error(), "inst_2") {
		t.Fatalf("expected err with specific words. Got %s", err)
	}

	tpl = MustParse(`
	sub = create subnet
	create vpc cidr=10.0.0.0/4
	`)

	_, _, err = checkReferencesDeclaration(tpl, env)
	if err == nil ||
		!strings.Contains(err.Error(), "unused") ||
		!strings.Contains(err.Error(), "sub") {
		t.Fatalf("expected err with specific words. Got %s", err)
	}
}

func TestResolveAgainstDefinitionsPass(t *testing.T) {
	env := NewEnv()
	env.DefLookupFunc = func(in string) (TemplateDefinition, bool) {
		t, ok := DefsExample[in]
		return t, ok
	}

	t.Run("Put definition required param in holes", func(t *testing.T) {
		tpl := MustParse(`create instance type=@custom_type count=$inst_num`)

		resolveAgainstDefinitions(tpl, env)

		assertCmdHoles(t, tpl, map[string]string{
			"subnet": "instance.subnet",
			"image":  "instance.image",
		})
		assertCmdParams(t, tpl, map[string]interface{}{
			"type": "@custom_type",
		})
		assertCmdRefs(t, tpl, map[string]string{
			"count": "inst_num",
		})
	})

	t.Run("Err on unexisting templ def", func(t *testing.T) {
		tpl := MustParse(`create subnet type=t2.micro`)

		_, _, err := resolveAgainstDefinitions(tpl, env)
		if err == nil || !strings.Contains(err.Error(), "createsubnet") {
			t.Fatalf("expected err with message containing 'createsubnet'")
		}
	})

	t.Run("Err on unexpected param key", func(t *testing.T) {
		tpl := MustParse(`create instance type=t2.micro
	                        create keypair name={keypair.name} type=wrong`)

		_, _, err := resolveAgainstDefinitions(tpl, env)
		if err == nil || !strings.Contains(err.Error(), "type") {
			t.Fatalf("expected err with message containing 'type'")
		}
	})

	t.Run("Err on unexpected ref key", func(t *testing.T) {
		tpl := MustParse(`create instance type=t2.micro
		create tag stuff=$any`)

		_, _, err := resolveAgainstDefinitions(tpl, env)
		if err == nil || !strings.Contains(err.Error(), "stuff") {
			t.Fatalf("expected err with message containing 'stuff'")
		}
	})

	t.Run("Err on unexpected hole key", func(t *testing.T) {
		tpl := MustParse(`create instance type=t2.micro
		create tag stuff={stuff.any}`)

		_, _, err := resolveAgainstDefinitions(tpl, env)
		if err == nil || !strings.Contains(err.Error(), "stuff") {
			t.Fatalf("expected err with message containing 'stuff'")
		}
	})
}

func TestResolveMissingHolesPass(t *testing.T) {
	tpl := MustParse(`
	create instance subnet={instance.subnet} type={instance.type} name={redis.prod}
	create vpc cidr={vpc.cidr}
	create instance name={redis.prod} id={redis.prod} count=3`)

	var count int
	env := NewEnv()
	env.MissingHolesFunc = func(in string) interface{} {
		count++
		switch in {
		case "instance.subnet":
			return "sub-98765"
		case "redis.prod":
			return "redis-124.32.34.54"
		case "vpc.cidr":
			return "10.0.0.0/24"
		default:
			return ""
		}
	}
	env.AddFillers(map[string]interface{}{"instance.type": "t2.micro"})

	pass := newMultiPass(resolveHolesPass, resolveMissingHolesPass)

	tpl, _, err := pass.compile(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := count, 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	assertCmdParams(t, tpl,
		map[string]interface{}{"type": "t2.micro", "name": "redis-124.32.34.54", "subnet": "sub-98765"},
		map[string]interface{}{"cidr": "10.0.0.0/24"},
		map[string]interface{}{"id": "redis-124.32.34.54", "name": "redis-124.32.34.54", "count": 3},
	)
}

func TestResolveAliasPass(t *testing.T) {
	tpl := MustParse("create instance subnet=@my-subnet ami={instance.ami} count=3")

	env := NewEnv()
	env.AliasFunc = func(e, k, v string) string {
		vals := map[string]string{
			"my-ami":    "ami-12345",
			"my-subnet": "sub-12345",
		}
		return vals[v]
	}
	env.AddFillers(map[string]interface{}{"instance.ami": "@my-ami"})

	pass := newMultiPass(resolveHolesPass, resolveAliasPass)

	tpl, _, err := pass.compile(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	assertCmdParams(t, tpl, map[string]interface{}{"subnet": "sub-12345", "ami": "ami-12345", "count": 3})
}

func TestResolveHolesPass(t *testing.T) {
	tpl := MustParse("create instance count={instance.count} type={instance.type}")

	env := NewEnv()
	env.AddFillers(map[string]interface{}{
		"instance.count": 3,
		"instance.type":  "t2.micro",
	})

	tpl, _, err := resolveHolesPass(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	assertCmdHoles(t, tpl, map[string]string{})
	assertCmdParams(t, tpl, map[string]interface{}{"type": "t2.micro", "count": 3})
}

type params map[string]interface{}
type paramsPerCommand []params

type holes map[string]string
type holesPerCommand []holes

type refs map[string]string
type refsPerCommand []refs

func assertCmdParams(t *testing.T, tpl *Template, exp ...params) {
	for i, cmd := range tpl.CommandNodesIterator() {
		if got, want := params(cmd.Params), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("params: cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
		}
	}
}

func assertCmdHoles(t *testing.T, tpl *Template, exp ...holes) {
	for i, cmd := range tpl.CommandNodesIterator() {
		if got, want := holes(cmd.Holes), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("holes: cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
		}
	}
}

func assertCmdRefs(t *testing.T, tpl *Template, exp ...refs) {
	for i, cmd := range tpl.CommandNodesIterator() {
		if got, want := refs(cmd.Refs), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("refs: cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
		}
	}
}
