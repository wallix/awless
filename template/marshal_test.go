package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/wallix/awless/template/internal/ast"
)

func TestTemplateExecutionUnmarshalFromJSON(t *testing.T) {
	tplExec := &TemplateExecution{}
	err := tplExec.UnmarshalJSON([]byte(`{
		"source": "create stuff",
		"locale": "eu-west-2",
		"fillers": {
			"mykey": "myvalue",
			"mysecondkey": "mysecondvalue"
		},
		"id": "123456", "author": "michael", "commands": [
		{"errors": ["first error"], "results": ["vpc-12345"], "line": "create vpc cidr=10.0.0.0/24"},
		{"line": "create subnet"},
		{"errors": ["third error"], "results": ["i-12345"], "line": "create instance type=t2.micro count=4"}
		]
	}`))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := tplExec.Source, "create stuff"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := tplExec.Fillers, map[string]interface{}{"mykey": "myvalue", "mysecondkey": "mysecondvalue"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	var cmds []*ast.CommandNode
	for _, cmd := range tplExec.CommandNodesIterator() {
		cmds = append(cmds, cmd)
	}

	if got, want := tplExec.ID, "123456"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := tplExec.Author, "michael"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := tplExec.Locale, "eu-west-2"; got != want {
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

func TestTemplateExecutionMarshalToJSON(t *testing.T) {
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
		source, locale, author string
		templ                  *Template
		fillers                map[string]interface{}
		out                    string
	}{
		{
			"create vpc\ncreate subnet\ncreate instance",
			"us-west-2", "michael",
			tmplWithErrors,
			nil,
			`{"source": "create vpc\ncreate subnet\ncreate instance",
			  "locale": "us-west-2",
			  "fillers": {}, 
				"id": "12345",
				"author": "michael",
				"commands": [
					{"errors": ["first error"], "results": ["first result"], "line": "create vpc"},
					{"line": "create subnet"},
					{"errors": ["third error"], "results": ["third result"], "line": "create instance"}
				]
		     }`,
		},
		{
			"create subnet cidr=10.0.0.0/24\ncreate instance name=@myinst",
			"us-west-1", "michael",
			MustParse("create subnet cidr=10.0.0.0/24\ncreate instance name=@myinst"),
			map[string]interface{}{"two": "2"},
			`{"source": "create subnet cidr=10.0.0.0/24\ncreate instance name=@myinst",
			  "locale": "us-west-1",
			  "fillers": {"two": "2"},
			  "author": "michael",
			  "id": "",
			  "commands": [
			    {"line": "create subnet cidr=10.0.0.0/24"},
			    {"line": "create instance name=@myinst"}
			  ]
			}`,
		},
		{
			"create instance name='my instance'",
			"eu-central-2", "michael",
			MustParse("create instance name='my instance'"),
			map[string]interface{}{"three": "3"},
			`{"source": "create instance name='my instance'",
			  "locale": "eu-central-2",
			  "author": "michael",
			  "fillers": {"three": "3"},
				"id": "",
				"commands": [
					{"line": "create instance name='my instance'"}
				]
			}`,
		},
		{
			"create instance name=\"my instance '$&\\ special) chars\"",
			"eu-central-1", "michael",
			MustParse("create instance name=\"my instance '$&\\ special) chars\""),
			map[string]interface{}{"four": 4},
			`{"source": "create instance name=\"my instance '$&\\ special) chars\"",
			  "fillers": {"four": 4},
			  "author": "michael",
			  "locale": "eu-central-1",
				"id": "",
				"commands": [
					{"line": "create instance name=\"my instance '$&\\ special) chars\""}
				]
			}`,
		},
	}

	for _, c := range tcases {
		tplExec := TemplateExecution{Template: c.templ, Source: c.source, Author: c.author, Locale: c.locale, Fillers: c.fillers}
		actual, err := tplExec.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := identJSON(actual), identJSON([]byte(c.out)); got != want {
			t.Fatalf("\ngot\n\n%s\nwant\n\n%s\n", got, want)
		}
	}
}

func identJSON(content []byte) string {
	var v interface{}
	err := json.Unmarshal(content, &v)
	if err != nil {
		fmt.Println("", string(content))
		panic(err)
	}
	ident, err := json.MarshalIndent(v, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(ident)
}
