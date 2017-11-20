package template

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/wallix/awless/template/internal/ast"
)

type mockCommand struct{ id string }

func (c *mockCommand) ValidateCommand(map[string]interface{}, []string) []error {
	return []error{errors.New(c.id)}
}
func (c *mockCommand) Run(ctx, params map[string]interface{}) (interface{}, error)    { return nil, nil }
func (c *mockCommand) DryRun(ctx, params map[string]interface{}) (interface{}, error) { return nil, nil }
func (c *mockCommand) ValidateParams(p []string) ([]string, error) {
	switch c.id {
	case "1", "2":
		return []string{c.id}, nil
	case "3":
		return []string{c.id}, errors.New("unexpected")
	}
	panic("wooot")
}

func (c *mockCommand) ConvertParams() ([]string, func(values map[string]interface{}) (map[string]interface{}, error)) {
	return []string{"param1", "param2"},
		func(values map[string]interface{}) (map[string]interface{}, error) {
			_, hasParam1 := values["param1"]
			_, hasParam2 := values["param2"]
			if hasParam1 && hasParam2 {
				return map[string]interface{}{"new": fmt.Sprint(values["param1"], values["param2"])}, nil
			}
			return values, nil
		}
}

func TestCommandsPasses(t *testing.T) {
	cmd1, cmd2, cmd3 := &mockCommand{"1"}, &mockCommand{"2"}, &mockCommand{"3"}
	var count int
	env := NewEnv()
	env.Lookuper = func(...string) interface{} {
		count++
		switch count {
		case 1:
			return cmd1
		case 2:
			return cmd2
		case 3:
			return cmd3
		default:
			panic("whaat")
		}
	}

	t.Run("verify commands exist", func(t *testing.T) {
		tpl := MustParse("create instance\nsub = create subnet\ncreate instance")
		count = 0
		_, _, err := verifyCommandsDefinedPass(tpl, env)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("convert params", func(t *testing.T) {
		tpl := MustParse("create instance\nsub = create subnet param1=anything param2=other\ncreate instance param1=anything")
		count = 0
		compiled, _, err := convertParamsPass(tpl, env)
		if err != nil {
			t.Fatal(err)
		}
		exp := map[string]interface{}{"new": "anythingother"}
		if got, want := compiled.CommandNodesIterator()[1].ToDriverParams(), exp; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		exp = map[string]interface{}{"param1": "anything"}
		if got, want := compiled.CommandNodesIterator()[2].ToDriverParams(), exp; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("validate commands params", func(t *testing.T) {
		tpl := MustParse("create instance\nsub = create subnet\ncreate instance")
		count = 0
		_, _, err := validateCommandsParamsPass(tpl, env)
		if err == nil {
			t.Fatal("expected err got none")
		}
		if got, want := err.Error(), "unexpected"; !strings.Contains(got, want) {
			t.Fatalf("%s should contain %s", got, want)
		}
	})

	t.Run("normalize missing required params as hole", func(t *testing.T) {
		tpl := MustParse("create instance\nsub = create subnet")
		count = 0
		compiled, _, err := normalizeMissingRequiredParamsAsHolePass(tpl, env)
		if err != nil {
			t.Fatal(err)
		}

		for i, cmd := range compiled.CommandNodesIterator() {
			if got, want := len(cmd.GetHoles()), 1; got != want {
				t.Fatalf("%d. got %d, want %d", i+1, got, want)
			}
			var first string
			for h := range cmd.GetHoles() {
				first = h
			}
			if got, want := first, fmt.Sprintf("%s.%d", cmd.Entity, i+1); got != want {
				t.Fatalf("%d. got %s, want %s", i+1, got, want)
			}
		}
	})

	t.Run("validate commands", func(t *testing.T) {
		tpl := MustParse("create instance\nsub = create subnet\ncreate instance")
		count = 0
		_, _, err := validateCommandsPass(tpl, env)
		if err == nil {
			t.Fatal("expected err got none")
		}

		checkContainsAll(t, err.Error(), "123")
	})

	t.Run("inject command", func(t *testing.T) {
		count = 0
		tpl := MustParse("create instance\nsub = create subnet\ncreate instance")

		compiled, _, err := injectCommandsPass(tpl, env)
		if err != nil {
			t.Fatal(err)
		}

		cmds := compiled.CommandNodesIterator()

		sameObject := func(got, want interface{}) bool {
			return reflect.ValueOf(got).Pointer() == reflect.ValueOf(want).Pointer()
		}
		if got, want := cmds[0].Command, cmd1; !sameObject(got, want) {
			t.Fatalf("different object: got %#v, want %#v", got, want)
		}
		if got, want := cmds[1].Command, cmd2; !sameObject(got, want) {
			t.Fatalf("different object: got %#v, want %#v", got, want)
		}
		if got, want := cmds[2].Command, cmd3; !sameObject(got, want) {
			t.Fatalf("different object: got %#v, want %#v", got, want)
		}
	})
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
		_, _, err := failOnUnresolvedAliasPass(tpl, env)
		if err == nil && tcase.expAliasErr != "" {
			t.Fatalf("%d: unresolved aliases: got nil error, expect '%s'", i+1, tcase.expAliasErr)
		} else if err != nil && tcase.expAliasErr == "" {
			t.Fatalf("%d: unresolved aliases: got '%s' error, expect nil", i+1, err.Error())
		} else if got, want := err, tcase.expAliasErr; got != nil && want != "" && !strings.Contains(err.Error(), want) {
			t.Fatalf("%d: unresolved aliases: got '%s', want '%s'", i+1, got.Error(), want)
		}

		_, _, err = failOnUnresolvedHolesPass(tpl, env)
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
		_, _, err := checkInvalidReferenceDeclarationsPass(MustParse(tcase.tpl), env)
		if tcase.expErr == "" && err != nil {
			t.Fatalf("%d: %v", i+1, err)
		}
		if tcase.expErr != "" && (err == nil || !strings.Contains(err.Error(), tcase.expErr)) {
			t.Fatalf("%d: got %v, expected %s", i+1, err, tcase.expErr)
		}
	}
}

type mockCommandWithResult struct{ id string }

func (c *mockCommandWithResult) ValidateCommand(map[string]interface{}, []string) []error { return nil }
func (c *mockCommandWithResult) Run(ctx, params map[string]interface{}) (interface{}, error) {
	return nil, nil
}
func (c *mockCommandWithResult) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return nil, nil
}
func (c *mockCommandWithResult) ValidateParams(p []string) ([]string, error) { return nil, nil }
func (c *mockCommandWithResult) ExtractResult(i interface{}) string          { return "" }

func TestFailOnDeclarationWithNoResultPass(t *testing.T) {
	env := NewEnv()
	env.Lookuper = func(tokens ...string) interface{} {
		switch strings.Join(tokens, "") {
		case "createinstance":
			return &mockCommandWithResult{"create instance"}
		case "createsubnet":
			return &mockCommandWithResult{"create subnet"}
		case "attachpolicy":
			return &mockCommand{"attach policy"}
		default:
			panic("whaat")
		}
	}
	tcases := []struct {
		tpl    string
		expErr string
	}{
		{"sub = create subnet\ncreate instance subnet=$sub", ""},
		{"sub = attach policy\ncreate instance subnet=$sub", "cannot assign"},
	}

	for i, tcase := range tcases {
		_, _, err := failOnDeclarationWithNoResultPass(MustParse(tcase.tpl), env)
		if tcase.expErr == "" && err != nil {
			t.Fatalf("%d: %v", i+1, err)
		}
		if tcase.expErr != "" && (err == nil || !strings.Contains(err.Error(), tcase.expErr)) {
			t.Fatalf("%d: got %v, expected %s", i+1, err, tcase.expErr)
		}
	}
}

func TestResolveMissingHolesPass(t *testing.T) {
	tpl := MustParse(`
	ip = {instance.elasticip}
	create instance subnet={instance.subnet} type={instance.type} name={redis.prod} ip=$ip
	create vpc cidr={vpc.cidr}
	create instance name={redis.prod} id={redis.prod} count=3`)

	var count int
	env := NewEnv()
	env.MissingHolesFunc = func(in string, paramPaths []string) interface{} {
		count++
		switch in {
		case "instance.subnet":
			if got, want := paramPaths, []string{"create.instance.subnet"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("%s: got %v, want %v", in, got, want)
			}
			return "sub-98765"
		case "redis.prod":
			sort.Strings(paramPaths)
			if got, want := paramPaths, []string{"create.instance.id", "create.instance.name"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("%s: got %v, want %v", in, got, want)
			}
			return "redis-124.32.34.54"
		case "vpc.cidr":
			if got, want := paramPaths, []string{"create.vpc.cidr"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("%s: got %v, want %v", in, got, want)
			}
			return "10.0.0.0/24"
		case "instance.elasticip":
			if got, want := len(paramPaths), 0; got != want {
				t.Fatalf("%s: got %#v, want no element", in, paramPaths)
			}
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
	env.AddFillers(map[string]interface{}{"instance.ami": ast.NewAliasValue("my-ami")})

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

func TestCmdErr(t *testing.T) {
	tcases := []struct {
		cmd    *ast.CommandNode
		err    interface{}
		ifaces []interface{}
		expErr error
	}{
		{&ast.CommandNode{Action: "create", Entity: "instance"}, nil, nil, nil},
		{&ast.CommandNode{Action: "create", Entity: "instance"}, "my error", nil, errors.New("create instance: my error")},
		{&ast.CommandNode{Action: "create", Entity: "instance"}, errors.New("my error"), nil, errors.New("create instance: my error")},
		{nil, "my error", nil, errors.New("my error")},
		{&ast.CommandNode{Action: "create", Entity: "instance"}, "my error with %s %d", []interface{}{"Donald", 1}, errors.New("create instance: my error with Donald 1")},
	}
	for i, tcase := range tcases {
		if got, want := cmdErr(tcase.cmd, tcase.err, tcase.ifaces...), tcase.expErr; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %#v, want %#v", i+1, got, want)
		}
	}
}

type params map[string]interface{}
type holesKeys map[string][]string
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

func assertCmdHoles(t *testing.T, tpl *Template, exp ...holesKeys) {
	for i, cmd := range tpl.CommandNodesIterator() {
		h := make(map[string][]string)
		for k, p := range cmd.Params {
			if withHoles, ok := p.(ast.WithHoles); ok && len(withHoles.GetHoles()) > 0 {
				for key := range withHoles.GetHoles() {
					h[k] = append(h[k], key)
				}
			}
		}
		if got, want := holesKeys(h), exp[i]; !reflect.DeepEqual(got, want) {
			t.Fatalf("holes keys: cmd %d: \ngot\n%v\n\nwant\n%v\n", i+1, got, want)
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

func checkContainsAll(t *testing.T, s, chars string) {
	for _, e := range chars {
		if !strings.ContainsRune(s, e) {
			t.Fatalf("%s does not contain '%q'", s, e)
		}
	}
}
