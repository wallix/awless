package template

import (
	"reflect"
	"strings"
	"testing"

	"github.com/wallix/awless/template/internal/ast"
)

func TestWholeCompilation(t *testing.T) {
	tcases := []struct {
		tpl                  string
		expect               string
		expProcessedFillers  map[string]interface{}
		expResolvedVariables map[string]interface{}
	}{
		{
			tpl: `subnetname = my-subnet
vpcref=@vpc
testsubnet = create subnet cidr={test.cidr} vpc=$vpcref name=$subnetname
update subnet id=$testsubnet public=true
instancecount = {instance.count}
create instance subnet=$testsubnet image=ami-12345 count=$instancecount name='my test instance'`,
			expect: `testsubnet = create subnet cidr=10.0.2.0/24 name=my-subnet vpc=vpc-1234
update subnet id=$testsubnet public=true
create instance count=42 image=ami-12345 name='my test instance' subnet=$testsubnet type=t2.micro`,
			expProcessedFillers:  map[string]interface{}{"instance.type": "t2.micro", "test.cidr": "10.0.2.0/24", "instance.count": 42},
			expResolvedVariables: map[string]interface{}{"subnetname": "my-subnet", "vpcref": "vpc-1234", "instancecount": 42},
		},
		{
			tpl: `
create loadbalancer subnets=[sub-1234, sub-2345,@subalias,@subalias] name=mylb
sub1 = create subnet cidr={test.cidr} vpc=@vpc name=subnet1
sub2 = create subnet cidr=10.0.3.0/24 vpc=@vpc name=subnet2
create loadbalancer subnets=[$sub1, $sub2, sub-3456,{backup-subnet}] name=mylb2
`,
			expect: `create loadbalancer name=mylb subnets=[sub-1234,sub-2345,sub-1111,sub-1111]
sub1 = create subnet cidr=10.0.2.0/24 name=subnet1 vpc=vpc-1234
sub2 = create subnet cidr=10.0.3.0/24 name=subnet2 vpc=vpc-1234
create loadbalancer name=mylb2 subnets=[$sub1,$sub2,sub-3456,sub-0987]`,
			expProcessedFillers:  map[string]interface{}{"test.cidr": "10.0.2.0/24", "backup-subnet": "sub-0987"},
			expResolvedVariables: map[string]interface{}{},
		},
		{
			tpl: `
lb0 = create loadbalancer subnets=[sub-1234, sub-2345,@subalias,@subalias] name=mylb
sub1 = create subnet cidr={test.cidr} vpc=@vpc name=subnet1
sub2 = create subnet cidr=10.0.3.0/24 vpc=@vpc name=subnet2
lb1 = create loadbalancer subnets=[$sub1, $sub2, sub-3456,{backup-subnet}] name=mylb2
`,
			expect: `lb0 = create loadbalancer name=mylb subnets=[sub-1234,sub-2345,sub-1111,sub-1111]
sub1 = create subnet cidr=10.0.2.0/24 name=subnet1 vpc=vpc-1234
sub2 = create subnet cidr=10.0.3.0/24 name=subnet2 vpc=vpc-1234
lb1 = create loadbalancer name=mylb2 subnets=[$sub1,$sub2,sub-3456,sub-0987]`,
			expProcessedFillers:  map[string]interface{}{"test.cidr": "10.0.2.0/24", "backup-subnet": "sub-0987"},
			expResolvedVariables: map[string]interface{}{},
		},
		{
			tpl: `
			a = "mysubnet-1"
b = $a
c = {mysubnet2.hole}
d = [$b,$c,{mysubnet3.hole},mysubnet-4]
create loadbalancer subnets=$d name=lb1
e=$b
secondlb = create loadbalancer subnets=[$e,mysubnet-4,{mysubnet5.hole}] name=lb2
`,
			expect: `create loadbalancer name=lb1 subnets=[mysubnet-1,mysubnet-2,mysubnet-3,mysubnet-4]
secondlb = create loadbalancer name=lb2 subnets=[mysubnet-1,mysubnet-4,mysubnet-5]`,
			expProcessedFillers:  map[string]interface{}{"mysubnet2.hole": "mysubnet-2", "mysubnet3.hole": "mysubnet-3", "mysubnet5.hole": "mysubnet-5"},
			expResolvedVariables: map[string]interface{}{"a": "mysubnet-1", "b": "mysubnet-1", "e": "mysubnet-1", "c": "mysubnet-2", "d": []interface{}{"mysubnet-1", "mysubnet-2", "mysubnet-3", "mysubnet-4"}},
		},
	}

	for i, tcase := range tcases {
		env := NewEnv()

		env.AddFillers(map[string]interface{}{
			"instance.type":  "t2.micro",
			"test.cidr":      "10.0.2.0/24",
			"instance.count": 42,
			"unused":         "filler",
			"backup-subnet":  "sub-0987",
			"mysubnet2.hole": "mysubnet-2",
			"mysubnet3.hole": "mysubnet-3",
			"mysubnet5.hole": "mysubnet-5",
		})
		env.AliasFunc = func(e, k, v string) string {
			vals := map[string]string{
				"vpc":      "vpc-1234",
				"subalias": "sub-1111",
			}
			return vals[v]
		}
		env.DefLookupFunc = func(in string) (Definition, bool) {
			t, ok := DefsExample[in]
			return t, ok
		}

		inTpl := MustParse(tcase.tpl)

		pass := newMultiPass(NormalCompileMode...)

		compiled, _, err := pass.compile(inTpl, env)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}

		if got, want := compiled.String(), tcase.expect; got != want {
			t.Fatalf("%d: got\n%s\nwant\n%s", i+1, got, want)
		}

		if got, want := env.GetProcessedFillers(), tcase.expProcessedFillers; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %v, want %v", i+1, got, want)
		}

		if got, want := env.ResolvedVariables, tcase.expResolvedVariables; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestInlineVariableWithValue(t *testing.T) {
	env := NewEnv()
	tcases := []struct {
		tpl      string
		expError string
		expTpl   string
	}{
		{"ip = 127.0.0.1\ncreate instance ip=$ip", "", "create instance ip=127.0.0.1"},
		{"ip = 1.2.3.4\ncreate instance ip=$ip\ncreate subnet cidr=$ip", "", "create instance ip=1.2.3.4\ncreate subnet cidr=1.2.3.4"},
	}

	for i, tcase := range tcases {
		inTpl := MustParse(tcase.tpl)

		resolvedTpl, _, err := inlineVariableValuePass(inTpl, env)
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
	tcases := []struct {
		tpl         string
		expAliasErr string
		expHolesErr string
	}{
		{tpl: "create subnet\ncreate instance subnet=@mysubnet name={instance.name}\ncreate instance", expAliasErr: "unresolved alias", expHolesErr: "unresolved holes"},
		{tpl: "create subnet\ncreate instance subnet=@mysubnet\ncreate instance", expAliasErr: "unresolved alias: [mysubnet]"},
		{tpl: "create subnet hole=@myhole\ncreate instance subnet=@mysubnet\ncreate instance", expAliasErr: "unresolved alias: [myhole mysubnet]"},
		{tpl: "create subnet name=subnet\nname=@myinstance\ncreate instance name=$myinstance\ncreate instance", expAliasErr: "unresolved alias: [myinstance]"},
		{tpl: "create subnet\ncreate instance name={instance.name}\ncreate instance", expHolesErr: "unresolved holes: [instance.name]"},
		{tpl: "create subnet\ncreate instance name={instance.name}\ncreate instance\ncreate subnet name={subnet.name}", expHolesErr: "unresolved holes: [instance.name subnet.name]"},
		{tpl: "subnetname = {subnet.name} create subnet name=$subnetname\ncreate instance name=instancename\ncreate instance", expHolesErr: "unresolved holes: [subnet.name]"},
		{tpl: "create subnet\ncreate instance name=instancename\ncreate instance\ncreate subnet subnet=name"},
	}

	for i, tcase := range tcases {
		tpl := MustParse(tcase.tpl)
		_, _, err := failOnUnresolvedAlias(tpl, env)
		if err == nil && tcase.expAliasErr != "" {
			t.Fatalf("%d: unresolved aliases: got nil error, expect '%s'", i+1, tcase.expAliasErr)
		} else if err != nil && tcase.expAliasErr == "" {
			t.Fatalf("%d: unresolved aliases: got '%s' error, expect nil", i+1, err.Error())
		} else if got, want := err, tcase.expAliasErr; got != nil && want != "" && !strings.Contains(err.Error(), want) {
			t.Fatalf("%d: unresolved aliases: got '%s', want '%s'", i+1, got.Error(), want)
		}

		_, _, err = failOnUnresolvedHoles(tpl, env)
		if err == nil && tcase.expHolesErr != "" {
			t.Fatalf("%d: unresolved holes: got nil error, expect '%s'", i+1, tcase.expHolesErr)
		} else if err != nil && tcase.expHolesErr == "" {
			t.Fatalf("%d: unresolved holes: got '%s' error, expect nil", i+1, err.Error())
		} else if got, want := err, tcase.expHolesErr; got != nil && want != "" && !strings.Contains(err.Error(), want) {
			t.Fatalf("%d: unresolved holes: got '%s', want '%s'", i+1, got.Error(), want)
		}
	}
}

func TestCheckInvalidReferencesDeclarationPass(t *testing.T) {
	env := NewEnv()
	tcases := []struct {
		tpl    string
		expErr string
	}{
		{"sub = create subnet\ninst = create instance subnet=$sub\nip = 127.0.0.1\ncreate instance subnet=$inst ip=$ip", ""},
		{"sub = create subnet\ninst = create instance subnet=$sub\ninst = create instance", "'inst' has already been assigned in template"},
		{"sub = create subnet\ninst = create instance subnet=$sub\ncreate instance subnet=$inst_2", "'inst_2' is undefined in template"},
		{"sub = create subnet\ncreate vpc cidr=10.0.0.0/4", ""},
		{"create instance subnet=$sub\nsub = create subnet", "'sub' is undefined in template"},
		{"create instance\nip = 127.0.0.1", ""},
		{"new_inst = create instance autoref=$new_inst\n", "'new_inst' is undefined in template"},
		{"a = $test", "'test' is undefined in template"},
		{"b = [test1,$test2,{test4}]", "'test2' is undefined in template"},
	}

	for i, tcase := range tcases {
		_, _, err := checkInvalidReferenceDeclarations(MustParse(tcase.tpl), env)
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

		assertCmdHoles(t, tpl, map[string][]string{
			"subnet": {"instance.subnet"},
			"image":  {"instance.image"},
		})
		assertCmdAliases(t, tpl, map[string][]string{
			"type": {"custom_type"},
		})
		assertCmdRefs(t, tpl, map[string][]string{
			"count": {"inst_num"},
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

	assertCmdHoles(t, tpl, map[string][]string{})
	assertCmdParams(t, tpl, map[string]interface{}{"type": "t2.micro", "count": 3})
}

type params map[string]interface{}
type holes map[string][]string
type refs map[string][]string
type aliases map[string][]string

func assertVariableValues(t *testing.T, tpl *Template, exp ...interface{}) {
	for i, decl := range tpl.expressionNodesIterator() {
		if vn, ok := decl.(*ast.ValueNode); ok {
			if got, want := vn.Value.Value(), exp[i]; !reflect.DeepEqual(got, want) {
				t.Fatalf("variables value %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
			}
		}

	}
}

func assertCmdParams(t *testing.T, tpl *Template, exp ...params) {
	for i, cmd := range tpl.CommandNodesIterator() {
		if got, want := params(cmd.ToDriverParams()), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("params: cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
		}
	}
}

func assertCmdHoles(t *testing.T, tpl *Template, exp ...holes) {
	for i, cmd := range tpl.CommandNodesIterator() {
		h := make(map[string][]string)
		for k, p := range cmd.Params {
			if withHoles, ok := p.(ast.WithHoles); ok && len(withHoles.GetHoles()) > 0 {
				h[k] = withHoles.GetHoles()
			}
		}
		if got, want := holes(h), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("holes: cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
		}
	}
}

func assertCmdRefs(t *testing.T, tpl *Template, exp ...refs) {
	for i, cmd := range tpl.CommandNodesIterator() {
		r := make(map[string][]string)
		for k, p := range cmd.Params {
			if withRefs, ok := p.(ast.WithRefs); ok && len(withRefs.GetRefs()) > 0 {
				r[k] = withRefs.GetRefs()
			}
		}
		if got, want := refs(r), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("refs: cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
		}
	}
}

func assertCmdAliases(t *testing.T, tpl *Template, exp ...aliases) {
	for i, cmd := range tpl.CommandNodesIterator() {
		r := make(map[string][]string)
		for k, p := range cmd.Params {
			if withAliases, ok := p.(ast.WithAlias); ok && len(withAliases.GetAliases()) > 0 {
				r[k] = withAliases.GetAliases()
			}
		}
		if got, want := aliases(r), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("refs: cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
		}
	}
}
