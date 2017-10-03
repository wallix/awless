package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/wallix/awless/template/internal/ast"
)

func TestTemplateIsOneLiner(t *testing.T) {
	tplExec := &TemplateExecution{}
	if err := tplExec.UnmarshalJSON([]byte(`{"commands":[{"line":"create vpc"},{"line": "detach policy"}]}`)); err != nil {
		t.Fatal(err)
	}
	if tplExec.IsOneLiner() {
		t.Fatal("not expecting one-liner")
	}

	if err := tplExec.UnmarshalJSON([]byte(`{"commands":[{"line":"create vpc"}]}`)); err != nil {
		t.Fatal(err)
	}
	if !tplExec.IsOneLiner() {
		t.Fatal("expecting one-liner")
	}
}
func TestSetMessageTruncatingSizeWhenNeeded(t *testing.T) {
	valid := strings.Repeat("a", 140)

	te := &TemplateExecution{}
	te.SetMessage(valid + "aa")
	if got, want := te.Message, valid[:len(valid)-3]+"..."; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	te.SetMessage(valid)
	if got, want := te.Message, valid; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
func TestGetStatsFromTemplateExecution(t *testing.T) {
	tplExec := &TemplateExecution{}
	if err := tplExec.UnmarshalJSON([]byte(`{
		"commands": [
			{"line": "create vpc"},
			{"line": "create subnet"},
			{"line": "create instance"},
			{"line": "create instance"},
			{"line": "create instance"},
			{"line": "create subnet"},
			{"line": "attach policy"},
			{"line": "stop instance"},
			{"line": "detach policy", "errors": ["any"]}
		]}`)); err != nil {
		t.Fatal(err)
	}
	expected := map[string]int{
		"create instance": 3, "create vpc": 1, "create subnet": 2, "attach policy": 1, "stop instance": 1, "detach policy": 1,
	}
	stats := tplExec.Stats()
	if got, want := stats.ActionEntityCount, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := stats.OKCount, 8; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := stats.KOCount, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := stats.CmdCount, 9; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := stats.Oneliner, ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	if err := tplExec.UnmarshalJSON([]byte(`{"commands": [{"line": "create vpc"}]}`)); err != nil {
		t.Fatal(err)
	}
	if got, want := tplExec.Stats().Oneliner, "create vpc"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTemplateExecutionUnmarshalFromJSON(t *testing.T) {
	tplExec := &TemplateExecution{}
	err := tplExec.UnmarshalJSON([]byte(`{
		"source": "create stuff",
		"locale": "eu-west-2",
		"message": "Make the CLI great again",
		"profile": "admin",
		"path": "http://gist.com/mytemplate.aws",
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
	if got, want := tplExec.Profile, "admin"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := tplExec.Message, "Make the CLI great again"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := tplExec.Path, "http://gist.com/mytemplate.aws"; got != want {
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
	if got, want := cmds[0].ToDriverParams(), exp; !reflect.DeepEqual(got, want) {
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
	if got, want := len(cmds[1].ToDriverParams()), 0; got != want {
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
	if got, want := cmds[2].ToDriverParams(), exp; !reflect.DeepEqual(got, want) {
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
		profile, message, path string
		templ                  *Template
		fillers                map[string]interface{}
		out                    string
	}{
		{
			"create vpc\ncreate subnet\ncreate instance",
			"us-west-2", "michael", "admin", "bijour", "http://gist.com",
			tmplWithErrors,
			nil,
			`{"source": "create vpc\ncreate subnet\ncreate instance",
			  "locale": "us-west-2",
			  "profile": "admin",
			  "message": "bijour",
			  "path": "http://gist.com",
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
			"us-west-1", "", "admin", "bijour", "http://gist.com",
			MustParse("create subnet cidr=10.0.0.0/24\ncreate instance name=@myinst"),
			map[string]interface{}{"two": "2"},
			`{"source": "create subnet cidr=10.0.0.0/24\ncreate instance name=@myinst",
			  "locale": "us-west-1",
			  "profile": "admin",
			  "message": "bijour",
			  "path": "http://gist.com",
			  "fillers": {"two": "2"},
			  "id": "",
			  "commands": [
			    {"line": "create subnet cidr=10.0.0.0/24"},
			    {"line": "create instance name=@myinst"}
			  ]
			}`,
		},
		{
			"create instance name='my instance'",
			"eu-central-2", "michael", "admin", "bijour", "http://gist.com",
			MustParse("create instance name='my instance'"),
			map[string]interface{}{"three": "3"},
			`{"source": "create instance name='my instance'",
			  "locale": "eu-central-2",
			  "profile": "admin",
			  "message": "bijour",
			  "path": "http://gist.com",
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
			"eu-central-1", "michael", "", "", "",
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
		{
			"create loadbalancer subnets={private.subnets}",
			"eu-central-1", "michael", "", "bijour", "",
			MustParse("create loadbalancer subnets=[subnet-1234,subnet-2345]"),
			map[string]interface{}{"private.subnets": []interface{}{"subnet-1234", "subnet-2345"}},
			`{"source": "create loadbalancer subnets={private.subnets}",
			  "fillers": {"private.subnets": ["subnet-1234","subnet-2345"]},
			  "author": "michael",
			  "message": "bijour",
			  "locale": "eu-central-1",
				"id": "",
				"commands": [
					{"line": "create loadbalancer subnets=[subnet-1234,subnet-2345]"}
				]
			}`,
		},
	}

	for _, c := range tcases {
		tplExec := TemplateExecution{Template: c.templ, Source: c.source, Author: c.author, Locale: c.locale, Profile: c.profile, Message: c.message, Path: c.path, Fillers: c.fillers}
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
