package template

import (
	"reflect"
	"testing"
)

func TestMergeExternalParamsPass(t *testing.T) {
	extTpl := MustParse(`create instance subnet=@my-subnet count=4`)
	tpl := MustParse(`create instance ami=r45ty3`)

	env := &Env{}
	env.AddExternalParams(extTpl.GetParams())

	tpl, _, err := mergeExternalParamsPass(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	assertAllParams(t, tpl, map[string]interface{}{
		"subnet": "@my-subnet",
		"ami":    "r45ty3",
		"count":  4,
	})
}

func TestResolveMissingHolesPass(t *testing.T) {
	tpl := MustParse(`
	create instance subnet={instance.subnet} type={instance.type} name={redis.prod}
	create instance name={redis.prod} count=3`)

	var count int
	env := &Env{
		MissingHolesFunc: func(in string) interface{} {
			count++
			switch in {
			case "instance.subnet":
				return "sub-98765"
			case "redis.prod":
				return "redis-124.32.34.54"
			default:
				return ""
			}
		},
	}
	env.AddFillers(map[string]interface{}{"instance.type": "t2.micro"})

	pass := newMultiPass(resolveHolesPass, resolveMissingHolesPass)

	tpl, _, err := pass.compile(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := count, 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	assertAllParams(t, tpl, map[string]interface{}{
		"type":   "t2.micro",
		"name":   "redis-124.32.34.54",
		"subnet": "sub-98765",
		"count":  3,
	})
}

func TestResolveAliasPass(t *testing.T) {
	tpl := MustParse("create instance subnet=@my-subnet ami={instance.ami} count=3")

	env := &Env{
		AliasFunc: func(k, v string) string {
			vals := map[string]string{
				"my-ami":    "ami-12345",
				"my-subnet": "sub-12345",
			}
			return vals[v]
		},
	}
	env.AddFillers(map[string]interface{}{"instance.ami": "@my-ami"})

	pass := newMultiPass(resolveHolesPass, resolveAliasPass)

	tpl, _, err := pass.compile(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	assertAllParams(t, tpl, map[string]interface{}{
		"subnet": "sub-12345",
		"ami":    "ami-12345",
		"count":  3,
	})
}

func TestResolveHolesPass(t *testing.T) {
	tpl := MustParse("create instance count={instance.count} type={instance.type}")

	env := &Env{}
	env.AddFillers(map[string]interface{}{
		"instance.count": 3,
		"instance.type":  "t2.micro",
	})

	tpl, _, err := resolveHolesPass(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	assertAllParams(t, tpl, map[string]interface{}{
		"type":  "t2.micro",
		"count": 3,
	})
}

func assertAllParams(t *testing.T, tpl *Template, exp map[string]interface{}) {
	all := make(map[string]interface{})
	for _, cmd := range tpl.CommandNodesIterator() {
		for k, v := range cmd.Params {
			all[k] = v
		}
	}
	if got, want := all, exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot\n%v\n\nwant\n%v\n", got, want)
	}
}
