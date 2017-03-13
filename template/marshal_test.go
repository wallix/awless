package template

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/wallix/awless/template/ast"
)

func TestUnmarshalFromJSON(t *testing.T) {
	tpl := &Template{}
	err := tpl.UnmarshalJSON([]byte(`{"id": "123456", "commands": [
	  {"errors": ["first error"], "results": ["vpc-12345"], "line": "create vpc cidr=10.0.0.0/24"},
	   {"line": "create subnet"},
	   {"errors": ["third error"], "results": ["i-12345"], "line": "create instance type=t2.micro count=4"}
	  ]
	 }`))
	if err != nil {
		t.Fatal(err)
	}

	var cmds []*ast.CommandNode
	for _, cmd := range tpl.CommandNodesIterator() {
		cmds = append(cmds, cmd)
	}

	if got, want := tpl.ID, "123456"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	if got, want := cmds[0].CmdResult, "vpc-12345"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := cmds[0].Action, "create"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cmds[0].Entity, "vpc"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	exp := map[string]interface{}{"cidr": "10.0.0.0/24"}
	if got, want := cmds[0].Params, exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cmds[0].CmdErr.Error(), "first error"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	if got, want := cmds[1].Action, "create"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cmds[1].Entity, "subnet"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	exp = make(map[string]interface{})
	if got, want := len(cmds[1].Params), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if cmds[1].CmdErr != nil {
		t.Fatal("expected nil error for cmd")
	}

	if got, want := cmds[2].CmdResult, "i-12345"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cmds[2].Action, "create"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cmds[2].Entity, "instance"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	exp = map[string]interface{}{"type": "t2.micro", "count": 4}
	if got, want := cmds[2].Params, exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := cmds[2].CmdErr.Error(), "third error"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestMarshalToJSON(t *testing.T) {
	tmplWithErrors := MustParse("create vpc\ncreate subnet\ncreate instance")
	tmplWithErrors.ID = "12345"
	for i, cmd := range tmplWithErrors.CommandNodesIterator() {
		if i == 0 {
			cmd.CmdErr = errors.New("first error")
			cmd.CmdResult = "first result"
		}
		if i == 2 {
			cmd.CmdErr = errors.New("third error")
			cmd.CmdResult = "third result"
		}
	}

	tcases := []struct {
		templ *Template
		out   string
	}{
		{
			tmplWithErrors,
			`{
			  "id": "12345",
			  "commands": [
			  {"errors": ["first error"], "results": ["first result"], "line": "create vpc"},
			   {"line": "create subnet"},
			   {"errors": ["third error"], "results": ["third result"], "line": "create instance"}
			  ]
		         }`,
		},
		{
			MustParse("create subnet cidr=10.0.0.0/24\ncreate instance name=@myinst"),
			`{ 
			  "id": "",
			  "commands": [
			    {"line": "create subnet cidr=10.0.0.0/24"},
			    {"line": "create instance name=@myinst"}
			  ]
			}`,
		},
	}

	for _, c := range tcases {
		actual, err := c.templ.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := identJSON(actual), identJSON([]byte(c.out)); got != want {
			t.Fatalf("\ngot\n\n%q\nwant\n\n%q\n", got, want)
		}
	}
}

func identJSON(content []byte) string {
	var v interface{}
	err := json.Unmarshal(content, &v)
	if err != nil {
		panic(err)
	}
	ident, err := json.MarshalIndent(v, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(ident)
}
