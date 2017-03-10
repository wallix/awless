package template

import (
	"reflect"
	"testing"

	"github.com/wallix/awless/logger"
)

func TestCompile(t *testing.T) {
	tpl := MustParse(`create keypair name={keypair}
	create instance name={instance.name}`)

	env := NewEnv()
	logger.DefaultLogger.SetVerbose(logger.ExtraVerboseF)
	env.Log = logger.DefaultLogger

	env.MissingHolesFunc = func(in string) interface{} {
		switch in {
		case "keypair":
			return "ewofh2o424"
		case "instance.name":
			return "redis"
		default:
			return ""
		}
	}

	tpl, _, err := Compile(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	assertAllParams(t, tpl,
		map[string]interface{}{"name": "ewofh2o424"},
		map[string]interface{}{"name": "redis"},
	)
}

func TestMergeExternalParamsPass(t *testing.T) {
	extTpl := MustParse(`create instance subnet=@my-subnet count=4`)
	tpl := MustParse(`create instance ami=r45ty3`)

	env := NewEnv()
	env.AddExternalParams(extTpl.GetNormalizedParams())

	tpl, _, err := mergeExternalParamsPass(tpl, env)
	if err != nil {
		t.Fatal(err)
	}

	assertAllParams(t, tpl, map[string]interface{}{"instance.subnet": "@my-subnet", "ami": "r45ty3", "instance.count": 4})
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
	assertAllParams(t, tpl,
		map[string]interface{}{"type": "t2.micro", "name": "redis-124.32.34.54", "subnet": "sub-98765"},
		map[string]interface{}{"cidr": "10.0.0.0/24"},
		map[string]interface{}{"id": "redis-124.32.34.54", "name": "redis-124.32.34.54", "count": 3},
	)
}

func TestResolveAliasPass(t *testing.T) {
	tpl := MustParse("create instance subnet=@my-subnet ami={instance.ami} count=3")

	env := NewEnv()
	env.AliasFunc = func(k, v string) string {
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

	assertAllParams(t, tpl, map[string]interface{}{"subnet": "sub-12345", "ami": "ami-12345", "count": 3})
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

	assertAllParams(t, tpl, map[string]interface{}{"type": "t2.micro", "count": 3})
}

type params map[string]interface{}
type paramsPerCommand []params

func assertAllParams(t *testing.T, tpl *Template, exp ...params) {
	for i, cmd := range tpl.CommandNodesIterator() {
		if got, want := params(cmd.Params), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
		}
	}
}
