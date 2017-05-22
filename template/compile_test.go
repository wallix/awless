package template

import (
	"reflect"
	"strings"
	"testing"

	"github.com/wallix/awless/template/internal/ast"
)

func TestWholeCompilation(t *testing.T) {
	env := NewEnv()

	env.AddFillers(map[string]interface{}{
		"instance.type":  "t2.micro",
		"test.cidr":      "10.0.2.0/24",
		"instance.count": 42,
		"unused":         "filler",
	})
	env.AliasFunc = func(e, k, v string) string {
		vals := map[string]string{
			"vpc": "vpc-1234",
		}
		return vals[v]
	}
	env.DefLookupFunc = func(in string) (Definition, bool) {
		t, ok := DefsExample[in]
		return t, ok
	}

	tcases := []struct {
		tpl    string
		expect string
	}{
		{
			`subnetname = my-subnet
vpcref=@vpc
testsubnet = create subnet cidr={test.cidr} vpc=$vpcref name=$subnetname
update subnet id=$testsubnet public=true
instancecount = {instance.count}
create instance subnet=$testsubnet image=ami-12345 count=$instancecount name='my test instance'`,
			`testsubnet = create subnet cidr=10.0.2.0/24 name=my-subnet vpc=vpc-1234
update subnet id=$testsubnet public=true
create instance count=42 image=ami-12345 name='my test instance' subnet=$testsubnet type=t2.micro`,
		},
	}

	for i, tcase := range tcases {
		inTpl := MustParse(tcase.tpl)

		pass := newMultiPass(NormalCompileMode...)

		compiled, _, err := pass.compile(inTpl, env)
		if err != nil {
			t.Fatal(err)
		}
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}

		if got, want := compiled.String(), tcase.expect; got != want {
			t.Fatalf("%d: got\n%s\nwant\n%s", i+1, got, want)
		}

		expProcessedFillers := map[string]interface{}{"instance.type": "t2.micro", "subnet.cidr": "10.0.2.0/24", "instance.count": 42}
		if got, want := env.GetProcessedFillers(), expProcessedFillers; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestReplaceVariableWithValue(t *testing.T) {
	env := NewEnv()
	tcases := []struct {
		tpl      string
		expError string
		expTpl   string
	}{
		{"ip = 127.0.0.1\ncreate instance ip=$ip", "", "ip = 127.0.0.1\ncreate instance ip=127.0.0.1"},
		{"ip = 1.2.3.4\ncreate instance ip=$ip\ncreate subnet cidr=$ip", "", "ip = 1.2.3.4\ncreate instance ip=1.2.3.4\ncreate subnet cidr=1.2.3.4"},
	}

	for i, tcase := range tcases {
		inTpl := MustParse(tcase.tpl)

		resolvedTpl, _, err := replaceVariableValuePass(inTpl, env)
		if tcase.expError != "" {
			if err == nil {
				t.Fatalf("%d: expected error, got nil", i+1)
			}
			if got, want := err.Error(), tcase.expError; !strings.Contains(got, want) {
				t.Fatalf("%d: got %s, want %s", i+1, got, want)
			}
			continue
		}
		if got, want := resolvedTpl.String(), tcase.expTpl; got != want {
			t.Fatalf("%d: got\n%s\nwant\n%s", i+1, got, want)
		}
	}
}

func TestRemoveVariablesPass(t *testing.T) {
	env := NewEnv()
	tcases := []struct {
		tpl    string
		expTpl string
	}{
		{"ip = 127.0.0.1\ncreate instance ip=127.0.0.1", "create instance ip=127.0.0.1"},
		{"ip = 1.2.3.4\ncreate instance ip=1.2.3.4\ncreate subnet cidr=2.3.4.5\nsubnet=2.3.4.5", "create instance ip=1.2.3.4\ncreate subnet cidr=2.3.4.5"},
		{"ip = {elasticip}\ncreate instance ip=$ip", "ip = {elasticip}\ncreate instance ip=$ip"},
	}

	for i, tcase := range tcases {
		inTpl := MustParse(tcase.tpl)

		resolvedTpl, _, err := removeValueStatementsPass(inTpl, env)
		if err != nil {
			t.Fatalf("%d: %v", i+1, err)
		}
		if got, want := resolvedTpl.String(), tcase.expTpl; got != want {
			t.Fatalf("%d: got\n%s\nwant\n%s", i+1, got, want)
		}
	}
}

func TestDefaultEnvWithNilFunc(t *testing.T) {
	text := "create instance name={instance.name} subnet=@mysubnet"
	env := NewEnv()
	tpl := MustParse(text)

	pass := newMultiPass(resolveHolesPass, resolveMissingHolesPass, resolveAliasPass)

	compiled, _, err := pass.compile(tpl, env)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if got, want := compiled.String(), text; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestBailOnUnresolvedAliasOrHoles(t *testing.T) {
	env := NewEnv()
	tpl := MustParse("create subnet\ncreate instance subnet=@mysubnet name={instance.name}\ncreate instance")

	_, _, err := failOnUnresolvedAlias(tpl, env)
	if err == nil || !strings.Contains(err.Error(), "unresolved alias") {
		t.Fatalf("expected err unresolved alias. Got %s", err)
	}

	_, _, err = failOnUnresolvedHoles(tpl, env)
	if err == nil || !strings.Contains(err.Error(), "unresolved holes") {
		t.Fatalf("expected err unresolved holes. Got %s", err)
	}
}

func TestCheckReferencesDeclarationPass(t *testing.T) {
	env := NewEnv()
	tcases := []struct {
		tpl    string
		expErr string
	}{
		{"sub = create subnet\ninst = create instance subnet=$sub\nip = 127.0.0.1\ncreate instance subnet=$inst ip=$ip", ""},
		{"sub = create subnet\ninst = create instance subnet=$sub\ninst = create instance", "'inst' has already been assigned in template"},
		{"sub = create subnet\ninst = create instance subnet=$sub\ncreate instance subnet=$inst_2", "'inst_2' is undefined in template"},
		{"sub = create subnet\ncreate vpc cidr=10.0.0.0/4", "unused reference 'sub' in template"},
		{"create instance subnet=$sub\nsub = create subnet", "'sub' is undefined in template"},
		{"create instance\nip = 127.0.0.1", "unused reference 'ip'"},
		{"new_inst = create instance autoref=$new_inst\n", "'new_inst' is undefined in template"},
	}

	for i, tcase := range tcases {
		_, _, err := checkReferencesDeclaration(MustParse(tcase.tpl), env)
		if tcase.expErr == "" && err != nil {
			t.Fatalf("%d: %v", i+1, err)
		}
		if tcase.expErr != "" && (err == nil || !strings.Contains(err.Error(), tcase.expErr)) {
			t.Fatalf("%d: got %v, expected %s", i+1, err, tcase.expErr)
		}
	}
}

func TestResolveAgainstDefinitionsPass(t *testing.T) {
	env := NewEnv()
	env.DefLookupFunc = func(in string) (Definition, bool) {
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
		tpl := MustParse(`create none type=t2.micro`)

		_, _, err := resolveAgainstDefinitions(tpl, env)
		if err == nil || !strings.Contains(err.Error(), "createnone") {
			t.Fatalf("expected err with message containing 'createnone'")
		}
	})

	t.Run("Err on unexpected param key", func(t *testing.T) {
		tpl := MustParse(`create instance type=t2.micro
	                        create keypair name={key.name} type=wrong`)

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
	ip = {instance.elasticip}
	create instance subnet={instance.subnet} type={instance.type} name={redis.prod} ip=$ip
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
		case "instance.elasticip":
			return "1.2.3.4"
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

	if got, want := count, 4; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	assertVariableValues(t, tpl,
		"1.2.3.4",
	)
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
	if got, want := env.GetProcessedFillers(), map[string]interface{}{"instance.ami": "ami-12345"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
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
type holes map[string]string
type refs map[string]string

func assertVariableValues(t *testing.T, tpl *Template, exp ...interface{}) {
	for i, decl := range tpl.expressionNodesIterator() {
		if vn, ok := decl.(*ast.ValueNode); ok {
			if got, want := vn.Value, exp[i]; !reflect.DeepEqual(got, want) {
				t.Fatalf("variables value %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
			}
		}

	}
}

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
