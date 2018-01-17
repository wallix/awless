/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package template

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/wallix/awless/template/internal/ast"
)

func TestParseTemplatesWithList(t *testing.T) {
	tcases := []struct {
		text string
	}{
		{"create loadbalancer subnets=[subnet1,subnet2,subnet3]"},
		{"lb = create loadbalancer subnets=[subnet1,subnet2,subnet3]"},
		{"create loadbalancer subnets=[$subnet1,$subnet2,$subnet3]"},
		{"lb = create loadbalancer subnets=[$subnet1,$subnet2,$subnet3]"},
		{"create loadbalancer subnets=[{subnet1},{subnet2},{subnet3}]"},
		{"lb = create loadbalancer subnets=[{subnet1},{subnet2},{subnet3}]"},
		{"create loadbalancer name=mylb subnets=[sub-1234,sub-2345]"},
		{"lb = create loadbalancer name=mylb subnets=[sub-1234,sub-2345]"},
		{"create loadbalancer name=mylb subnets=[sub-1234,$subnet2,{subnet3}]"},
		{"lb = create loadbalancer name=mylb subnets=[sub-1234,$subnet2,{subnet3}]"},
		{"create loadbalancer name=mylb subnets=[@mysubnet,$subnet2,{subnet3}]"},
		{"lb = create loadbalancer name=mylb subnets=[@mysubnet,$subnet2,{subnet3}]"},
	}

	for i, tcase := range tcases {
		tpl, err := Parse(tcase.text)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := tpl.String(), tcase.text; !strings.HasSuffix(got, want) {
			t.Fatalf("%d: parsing [%s]\ngot  [%s]\nwant [%s]\n", i+1, tcase.text, got, want)
		}
	}
}

func TestParseVariousTemplatesCorrectly(t *testing.T) {
	tcases := []struct {
		desc string
		text string
		exp  string
	}{
		{"keep quote in output", "create policy description=\"my desc\"", "create policy description='my desc'"},
		{"support wildcard", "create policy action=ec2:Get*", ""},
		{"support wildcard in quote", "create policy action=\"ec2:Get*\"", "create policy action=ec2:Get*"},
		{"support single wildcard", "create policy resource=*", ""},
		{"support parameter value beginning with number", "create keypair name=123test", ""},
		{"support prefix/suffixes around holes (values)", "name = prefix-{instance.name}-{instance.version}-suffix", "name = 'prefix-'+{instance.name}+'-'+{instance.version}+'-suffix'"},
		{"support prefix/suffixes around holes (params)", "instance = create instance name=prefix-{instance.name}-{instance.version}-suffix", "instance = create instance name='prefix-'+{instance.name}+'-'+{instance.version}+'-suffix'"},
		{"support suffix after holes", "instance = create instance name={instance.name}-suffix", "instance = create instance name={instance.name}+'-suffix'"},
		{"retro-compatibility with old lists without []", "create loadbalancer subnets=subnet-1,subnet-2", "create loadbalancer subnets=[subnet-1,subnet-2]"},
		{"support concatenation with '+' of quoted string and holes", "instance = create instance name='prefix-'+{instance.name}+{instance.version}+'-suffix'", ""},
		{"support concatenation with '+' of quoted string and holes", "instance = create instance name='pre${}fix-' + {instance.name}+'middle-' +{instance.version}+ '-suffix'", "instance = create instance name='pre${}fix-'+{instance.name}+'middle-'+{instance.version}+'-suffix'"},
		{"support concatenation with '+' of quoted string and holes with a hole as prefix", "instance = create instance name={instance.name}+'midl${}fix-'+'midle2${}fix-'+{instance.version}+'-suffix'", ""},
	}

	for _, tcase := range tcases {
		tpl, err := Parse(tcase.text)
		if err != nil {
			t.Fatal(err)
		}

		exp := tcase.text
		if tcase.exp != "" {
			exp = tcase.exp
		}

		if got, want := tpl.String(), exp; got != want {
			t.Fatalf("%s: parsing [%s]\ngot  [%s]\nwant [%s]\n", tcase.desc, tcase.text, got, want)
		}
	}
}

func TestStringWithDigitValues(t *testing.T) {
	tcases := []struct {
		text      string
		expParams map[string]interface{}
	}{
		{"create keypair name=1test", map[string]interface{}{"name": "1test"}},
		{"create keypair name=11test", map[string]interface{}{"name": "11test"}},
		{"create keypair name=123test", map[string]interface{}{"name": "123test"}},
		{"create keypair name=0test", map[string]interface{}{"name": "0test"}},
		{"create keypair name=110", map[string]interface{}{"name": 110}},
		{"create keypair name=1/test", map[string]interface{}{"name": "1/test"}},
		{"create keypair name=123456789", map[string]interface{}{"name": 123456789}},
		{"create keypair name=0.5", map[string]interface{}{"name": 0.5}},
		{"create keypair name=0.5:0.6:+1", map[string]interface{}{"name": "0.5:0.6:+1"}},
	}

	for i, tcase := range tcases {
		tpl, err := Parse(tcase.text)
		if err != nil {
			t.Fatalf("%d. %s", i+1, err)
		}

		if n, ok := tpl.Statements[0].Node.(*ast.CommandNode); ok {
			if got, want := n.ToDriverParams(), tcase.expParams; !reflect.DeepEqual(got, want) {
				t.Fatalf("%d. got %#v, want %#v", i+1, got, want)
			}
		} else {
			t.Fatalf("expected command node, was %T", n)
		}

	}
}

func TestParseDoubleQuotedString(t *testing.T) {
	tcases := []struct {
		text, exp string
	}{
		{"create instance data=\"\"", ""},
		{"create instance data=\"hello\"", "hello"},
		{"create instance data=\"hello.\"", "hello."},
		{"create instance data=\"just jack\"", "just jack"},
		{"create instance data=\" just  jack \"", " just  jack "},

		{"create instance data=\"\t\tjust\t \tjack\t\"", "\t\tjust\t \tjack\t"},

		{"create instance data=\"just jack\n\"", "just jack\n"},
		{"create instance data=\"just jack\r\"", "just jack\r"},
		{"create instance data=\"just jack\n\r\"", "just jack\n\r"},
		{"create instance data=\"just jack\r\n\"", "just jack\r\n"},

		{"create instance data=\"\njust jack\"", "\njust jack"},
		{"create instance data=\"\rjust jack\"", "\rjust jack"},
		{"create instance data=\"\n\rjust jack\"", "\n\rjust jack"},
		{"create instance data=\"\r\njust jack\"", "\r\njust jack"},

		{"create instance data=\"!£$%^*()/{}_-+=:;@'~#,.?/<>\"", "!£$%^*()/{}_-+=:;@'~#,.?/<>"},
		{"create instance data=\"#!/bin/bash;touch /home/ubuntu/stuff.txt\"", "#!/bin/bash;touch /home/ubuntu/stuff.txt"},
	}

	for i, tcase := range tcases {
		tpl, err := Parse(tcase.text)
		if err != nil {
			t.Fatalf("%d. %s", i+1, err)
		}

		if n, ok := tpl.Statements[0].Node.(*ast.CommandNode); ok {
			if got, want := n.ParamNodes["data"].(ast.InterfaceNode).Value(), tcase.exp; got != want {
				t.Fatalf("%d. got %s, want %s", i+1, got, want)
			}
		} else {
			t.Fatalf("expected command node, was %T", n)
		}

	}
}

func TestParseSingleQuotedString(t *testing.T) {
	tcases := []struct {
		text, exp string
	}{
		{"create instance data=''", ""},
		{"create instance data='hello'", "hello"},
		{"create instance data='hello.'", "hello."},
		{"create instance data='just jack'", "just jack"},
		{"create instance data=' just  jack '", " just  jack "},

		{"create instance data='\t\tjust\t \tjack\t'", "\t\tjust\t \tjack\t"},

		{"create instance data='just jack\n'", "just jack\n"},
		{"create instance data='just jack\r'", "just jack\r"},
		{"create instance data='just jack\n\r'", "just jack\n\r"},
		{"create instance data='just jack\r\n'", "just jack\r\n"},

		{"create instance data='\njust jack'", "\njust jack"},
		{"create instance data='\rjust jack'", "\rjust jack"},
		{"create instance data='\n\rjust jack'", "\n\rjust jack"},
		{"create instance data='\r\njust jack'", "\r\njust jack"},

		{"create instance data='!£$%^*()/{}_-+=:;@\"~#,.?/<>'", "!£$%^*()/{}_-+=:;@\"~#,.?/<>"},
		{"create instance data='#!/bin/bash;touch /home/ubuntu/stuff.txt'", "#!/bin/bash;touch /home/ubuntu/stuff.txt"},
	}

	for i, tcase := range tcases {
		tpl, err := Parse(tcase.text)
		if err != nil {
			t.Fatalf("%d. %s", i+1, err)
		}

		if n, ok := tpl.Statements[0].Node.(*ast.CommandNode); ok {
			if got, want := n.ParamNodes["data"].(ast.InterfaceNode).Value(), tcase.exp; got != want {
				t.Fatalf("%d. got %s, want %s", i+1, got, want)
			}
		} else {
			t.Fatalf("expected command node, was %T", n)
		}

	}
}

func TestParsingInvalidActionAndEntities(t *testing.T) {
	_, err := Parse(`creat instance`)
	if err == nil || !strings.Contains(err.Error(), "action 'creat'") {
		t.Fatalf("expected error with specific message, got: %s", err)
	}

	_, err = Parse(`create intance`)
	if err == nil || !strings.Contains(err.Error(), "entity 'intance'") {
		t.Fatalf("expected error with specific message, got: %s", err)
	}
}

func TestParsingEmptyTemplate(t *testing.T) {
	_, err := Parse(``)
	if err == nil || err.Error() != "empty template" {
		t.Fatalf("expected error with specific message, got: %s", err)
	}
}

func TestWrapPegParseError(t *testing.T) {
	t.Run("Display better error message", func(t *testing.T) {
		text := "create subnet\ncreate instance type= wrong=\ncreate vpc"
		_, err := Parse(text)

		perr, _ := err.(*parseError)

		if got, want := perr.line, 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := perr.start, 23; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := perr.end, 28; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		exp := "error parsing template at line 2 (char 23):\n\t   create subnet\n\t-> create instance type= w\n\t   create vpc"
		if got, want := err.Error(), exp; got != want {
			t.Fatalf("got\n\n%s\n\nwant\n\n%s\n", got, want)
		}
	})

	t.Run("Fallback on peg msg if cannot find contextual info", func(t *testing.T) {
		orig := "rubbish"
		err := newParseError("create vpc\ncreate vpc", orig)
		if got, want := err.Error(), orig; got != want {
			t.Fatalf("got\n\n%q\n\nwant\n\n%q\n", got, want)
		}
	})

	t.Run("Fallback on peg msg if out of bounds indexes info", func(t *testing.T) {
		orig := "line 1 symbol 21 - line 1 symbol 22"
		err := newParseError("create vpc\ncreate vpc", orig)
		if got, want := err.Error(), orig; got != want {
			t.Fatalf("got\n\n%q\n\nwant\n\n%q\n", got, want)
		}
	})
}

func TestParamsOnlyParsing(t *testing.T) {
	tcases := []struct {
		input string
		exp   map[string]interface{}
	}{
		{input: "type=t2.micro subnet=@my-subnet count=4", exp: map[string]interface{}{"type": "t2.micro", "subnet": ast.NewAliasNode("my-subnet"), "count": 4}},
		{input: "subnet=[sub-1234,sub-2345]", exp: map[string]interface{}{"subnet": ast.NewListNode([]interface{}{"sub-1234", "sub-2345"})}},
	}
	for i, tcase := range tcases {
		params, err := ParseParams(tcase.input)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}
		if got, want := params, tcase.exp; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got\n%#v\n\nwant\n%#v\n", i+1, got, want)
		}
	}

}

func TestTemplateParsing(t *testing.T) {
	t.Run("Parse special characters", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(tpl *Template) error
		}{
			{
				input: "attach policy arn=arn:aws:iam::aws:policy/AmazonS3FullAccess",
				verifyFn: func(tpl *Template) error {
					if err := isCommandNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
			{
				input: "create instance name=a2zR_-+:;@~./<>",
				verifyFn: func(tpl *Template) error {
					if err := isCommandNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return assertParams(tpl.Statements[0].Node, map[string]interface{}{"name": "a2zR_-+:;@~./<>"})
				},
			},
			{
				input: "attach policy arn=@arn:aws:iam::aws:policy/AmazonS3FullAccess",
				verifyFn: func(tpl *Template) error {
					if err := isCommandNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return assertAliases(tpl.Statements[0].Node, map[string]string{"arn": "arn:aws:iam::aws:policy/AmazonS3FullAccess"})
				},
			},
			{
				input: "attach instance id=@\"my vm name\"",
				verifyFn: func(tpl *Template) error {
					if err := isCommandNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return assertAliases(tpl.Statements[0].Node, map[string]string{"id": "my vm name"})
				},
			},
			{
				input: "attach instance id=@'my f$!=€&g vm name'",
				verifyFn: func(tpl *Template) error {
					if err := isCommandNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return assertAliases(tpl.Statements[0].Node, map[string]string{"id": "my f$!=€&g vm name"})
				},
			},
			{
				input: "vpc_1 = create vpc",
				verifyFn: func(tpl *Template) error {
					if err := isDeclarationNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
			{
				input: "launch = create launchconfiguration spotprice=0.01",
				verifyFn: func(tpl *Template) error {
					if err := isDeclarationNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
			{
				input: "launch = create launchconfiguration spotprice=100.",
				verifyFn: func(tpl *Template) error {
					if err := isDeclarationNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
		}

		for _, tcase := range tcases {
			node, err := Parse(tcase.input)
			if err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}

			if err := tcase.verifyFn(node); err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}
		}
	})

	t.Run("Allow and ignore comments", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(tpl *Template) error
		}{
			{
				input: "create vpc\n#my comment\ncreate subnet",
				verifyFn: func(tpl *Template) error {
					if got, want := len(tpl.Statements), 2; got != want {
						t.Fatalf("got %d, want %d", got, want)
					}
					if err := isCommandNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					if err := isCommandNode(tpl.Statements[1].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
			{
				input: "create vpc \n//my comment\ncreate subnet",
				verifyFn: func(tpl *Template) error {
					if got, want := len(tpl.Statements), 2; got != want {
						t.Fatalf("got %d, want %d", got, want)
					}
					if err := isCommandNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					if err := isCommandNode(tpl.Statements[1].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
		}

		for _, tcase := range tcases {
			node, err := Parse(tcase.input)
			if err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}

			if err := tcase.verifyFn(node); err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}
		}
	})

	t.Run("Onliner statement", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(n ast.Node) error
		}{
			{
				input:    `create vpc`,
				verifyFn: func(n ast.Node) error { return assertParams(n, make(map[string]interface{})) },
			},
			{
				input:    `create vpc`,
				verifyFn: func(n ast.Node) error { return isCommandNode(n) },
			},
			{
				input:    `mysubnet = create subnet`,
				verifyFn: func(n ast.Node) error { return isDeclarationNode(n) },
			},
			{
				input: `create vpc cidr=10.0.0.0/24 num=3 ip=127.0.0.1 name=bousin`,
				verifyFn: func(n ast.Node) error {
					return assertParams(n, map[string]interface{}{"cidr": "10.0.0.0/24", "num": 3, "ip": "127.0.0.1", "name": "bousin"})
				},
			},
			{
				input: `create vpc cidr="10.0.0.0/24" num="3" ip="127.0.0.1" name="bousin"`,
				verifyFn: func(n ast.Node) error {
					return assertParams(n, map[string]interface{}{"cidr": "10.0.0.0/24", "num": "3", "ip": "127.0.0.1", "name": "bousin"})
				},
			},
			{
				input: `create vpc cidr='10.0.0.0/24' num='3' ip='127.0.0.1' name='bousin'`,
				verifyFn: func(n ast.Node) error {
					return assertParams(n, map[string]interface{}{"cidr": "10.0.0.0/24", "num": "3", "ip": "127.0.0.1", "name": "bousin"})
				},
			},
			{
				input: `create subnet vpc=$myvpc`,
				verifyFn: func(n ast.Node) error {
					return assertRefs(n, map[string]string{"vpc": "myvpc"})
				},
			},
			{
				input: `create instance subnet=@my-subnet`,
				verifyFn: func(n ast.Node) error {
					return assertAliases(n, map[string]string{"subnet": "my-subnet"})
				},
			},
			{
				input: `delete vpc id={my-vpc-id}`,
				verifyFn: func(n ast.Node) error {
					return assertHoleKeys(n, map[string]string{"id": "my-vpc-id"})
				},
			},
			{
				input: `create securitygroup port=20-80`,
				verifyFn: func(n ast.Node) error {
					if err := assertParams(n, map[string]interface{}{"port": "20-80"}); err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `create securitygroup port="20-80"`,
				verifyFn: func(n ast.Node) error {
					if err := assertParams(n, map[string]interface{}{"port": "20-80"}); err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `create securitygroup port='20-80'`,
				verifyFn: func(n ast.Node) error {
					if err := assertParams(n, map[string]interface{}{"port": "20-80"}); err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `create vpc array=[test1,test2, 20 , my-array-elem4] ip=127.0.0.1`,
				verifyFn: func(n ast.Node) error {
					return assertCmdNodeParams(n, map[string]interface{}{"array": []interface{}{"test1", "test2", 20, "my-array-elem4"}, "ip": "127.0.0.1"})
				},
			},
			{
				input: `create vpc array=["test1","test2", "20" , "my-array-elem4"] ip="127.0.0.1"`,
				verifyFn: func(n ast.Node) error {
					return assertCmdNodeParams(n, map[string]interface{}{"array": []interface{}{"test1", "test2", "20", "my-array-elem4"}, "ip": "127.0.0.1"})
				},
			},
			{
				input: `create vpc array=['test1','test2', '20' , 'my-array-elem4'] ip='127.0.0.1'`,
				verifyFn: func(n ast.Node) error {
					return assertCmdNodeParams(n, map[string]interface{}{"array": []interface{}{"test1", "test2", "20", "my-array-elem4"}, "ip": "127.0.0.1"})
				},
			},
			{
				input: `myinstance = create instance type={instance.type} cidr=10.0.0.0/25 subnet=@default-subnet vpc=$myvpc`,
				verifyFn: func(n ast.Node) error {
					if err := assertParams(n, map[string]interface{}{"cidr": "10.0.0.0/25"}); err != nil {
						return err
					}
					if err := assertHoleKeys(n, map[string]string{"type": "instance.type"}); err != nil {
						return err
					}
					if err := assertRefs(n, map[string]string{"vpc": "myvpc"}); err != nil {
						return err
					}
					if err := assertAliases(n, map[string]string{"subnet": "default-subnet"}); err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `create policy name=policyName effect=Allow action=[ec2:Describe*,autoscaling:Describe*,elasticloadbalancing:Describe*] resource=["arn:aws:iam::0123456789:mfa/${aws:username}", "arn:aws:iam::0123456789:user/${aws:username}"] conditions=["aws:MultiFactorAuthPresent==true", "aws:TokenIssueTime!=Null"]`,
				verifyFn: func(n ast.Node) error {
					if err := assertCmdNodeParams(n, map[string]interface{}{
						"name":       "policyName",
						"effect":     "Allow",
						"action":     []interface{}{"ec2:Describe*", "autoscaling:Describe*", "elasticloadbalancing:Describe*"},
						"resource":   []interface{}{"arn:aws:iam::0123456789:mfa/${aws:username}", "arn:aws:iam::0123456789:user/${aws:username}"},
						"conditions": []interface{}{"aws:MultiFactorAuthPresent==true", "aws:TokenIssueTime!=Null"}}); err != nil {
						return err
					}
					if err := assertHoleKeys(n, map[string]string{}); err != nil {
						return err
					}
					return nil
				},
			},
		}

		for _, tcase := range tcases {
			node, err := parseStatement(tcase.input)
			if err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}

			if err := tcase.verifyFn(node); err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}
		}
	})

	t.Run("Multiline parsing", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(s *Template) error
		}{
			{
				input: `create vpc
create subnet`,
				verifyFn: func(s *Template) error {
					if err := assertCommandNode(s.Statements[0].Node, "create", "vpc",
						make(map[string]string), make(map[string]interface{}), make(map[string]string), make(map[string]string),
					); err != nil {
						return err
					}
					if err := assertCommandNode(s.Statements[1].Node, "create", "subnet",
						make(map[string]string), make(map[string]interface{}), make(map[string]string), make(map[string]string),
					); err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `


				create vpc

create subnet


`,
				verifyFn: func(s *Template) error {
					if err := assertCommandNode(s.Statements[0].Node, "create", "vpc",
						make(map[string]string), make(map[string]interface{}), make(map[string]string), make(map[string]string),
					); err != nil {
						return err
					}
					if err := assertCommandNode(s.Statements[1].Node, "create", "subnet",
						make(map[string]string), make(map[string]interface{}), make(map[string]string), make(map[string]string),
					); err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `# first comment. Next line contains whitespaces on purpose
       
				# second comment

				create vpc  # inlined comment
create subnet
# third statement

`,
				verifyFn: func(s *Template) error {
					if err := assertCommandNode(s.Statements[0].Node, "create", "vpc",
						make(map[string]string), make(map[string]interface{}), make(map[string]string), make(map[string]string),
					); err != nil {
						return err
					}
					if err := assertCommandNode(s.Statements[1].Node, "create", "subnet",
						make(map[string]string), make(map[string]interface{}), make(map[string]string), make(map[string]string),
					); err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `
			myvpc  =   create   vpc  cidr=10.0.0.0/24 num=3
mysubnet = delete subnet vpc=$myvpc name={ the_name } cidr=10.0.0.0/25
create instance count=1 instance.type=t2.micro subnet=$mysubnet image=ami-9398d3e0 ip=127.0.0.1
                       `,

				verifyFn: func(s *Template) error {
					err := assertDeclarationNode(s.Statements[0].Node, "myvpc", "create", "vpc",
						map[string]string{},
						map[string]interface{}{"cidr": "10.0.0.0/24", "num": 3},
						map[string]string{},
						map[string]string{},
					)
					if err != nil {
						return err
					}

					err = assertDeclarationNode(s.Statements[1].Node, "mysubnet", "delete", "subnet",
						map[string]string{"vpc": "myvpc"},
						map[string]interface{}{"cidr": "10.0.0.0/25"},
						map[string]string{"name": "the_name"},
						map[string]string{},
					)
					if err != nil {
						return err
					}

					err = assertCommandNode(s.Statements[2].Node, "create", "instance",
						map[string]string{"subnet": "mysubnet"},
						map[string]interface{}{"count": 1, "instance.type": "t2.micro", "ip": "127.0.0.1", "image": "ami-9398d3e0"},
						map[string]string{},
						map[string]string{},
					)

					return err
				},
			},
			{
				input: `
			myname  =   "my var-value"
mysubnet = create subnet vpc=$myvpc name=$myname cidr=10.0.0.0/25
mysecondvar = {var-hole}
                       `,

				verifyFn: func(s *Template) error {
					err := assertVariableDeclarationNode(s.Statements[0].Node, "myname", "my var-value", "")
					if err != nil {
						return err
					}

					err = assertDeclarationNode(s.Statements[1].Node, "mysubnet", "create", "subnet",
						map[string]string{"vpc": "myvpc", "name": "myname"},
						map[string]interface{}{"cidr": "10.0.0.0/25"},
						map[string]string{},
						map[string]string{},
					)
					if err != nil {
						return err
					}

					err = assertVariableDeclarationNode(s.Statements[2].Node, "mysecondvar", nil, "{var-hole}")
					if err != nil {
						return err
					}

					return err
				},
			},
		}

		for i, tcase := range tcases {
			templ, err := Parse(tcase.input)
			if err != nil {
				t.Fatalf("\n%d: input: [%s]\nError: %s\n", i, tcase.input, err)
			}

			if err := tcase.verifyFn(templ); err != nil {
				t.Fatalf("\n%d: input: [%s]\nError: %s\n", i, tcase.input, err)
			}
		}
	})

	t.Run("Create s3 object", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(s *Template) error
		}{
			{
				input: `create s3object bucket=my-existing-bucket file=./todolist.txt`,
				verifyFn: func(s *Template) error {
					if err := assertCommandNode(s.Statements[0].Node, "create", "s3object",
						map[string]string{}, map[string]interface{}{"bucket": "my-existing-bucket", "file": "./todolist.txt"}, map[string]string{}, map[string]string{},
					); err != nil {
						return err
					}
					return nil
				},
			},
		}

		for _, tcase := range tcases {
			templ, err := Parse(tcase.input)
			if err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}

			if err := tcase.verifyFn(templ); err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}
		}
	})
}

func assertParams(n ast.Node, expected map[string]interface{}) error {
	compare := func(got, want map[string]interface{}) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("params: got %#v, want %#v", got, want)
		}
		return nil
	}

	cmd := extractCommandNode(n)
	return compare(cmd.ToDriverParams(), expected)
}

func assertAliases(n ast.Node, expected map[string]string) error {
	compare := func(got, want map[string]string) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("aliases: got %#v, want %#v", got, want)
		}
		return nil
	}

	cmd := extractCommandNode(n)
	aliases := make(map[string]string)
	for k, param := range cmd.ParamNodes {
		if alias, ok := param.(ast.AliasNode); ok {
			aliases[k] = alias.Alias()
		}
	}
	return compare(aliases, expected)
}

func assertCmdNodeParams(n ast.Node, expected map[string]interface{}) error {
	cmd := extractCmdNode(n)
	if got, want := len(cmd.ParamNodes), len(expected); got != want {
		return fmt.Errorf("got %d params (%#v), want %d params (%#v)", got, cmd.ParamNodes, want, expected)
	}
	for key, expVal := range expected {
		node, ok := cmd.ParamNodes[key]
		if !ok {
			return fmt.Errorf("param '%s' missing in action params.", key)
		}
		var value interface{}
		if i, ok := node.(ast.InterfaceNode); ok {
			value = i.Value()
		}
		if l, ok := node.(ast.ListNode); ok {
			var arr []interface{}
			for _, e := range l.Elems() {
				switch ev := e.(type) {
				case ast.InterfaceNode:
					arr = append(arr, ev.Value())
				default:
					arr = append(arr, ev)
				}
			}
			value = arr
		}
		if got, want := value, expVal; !reflect.DeepEqual(got, want) {
			return fmt.Errorf("param '%s': got %#v, want %#v", key, got, want)
		}
	}
	return nil
}

func assertRefs(n ast.Node, expected map[string]string) error {
	compare := func(got, want map[string]string) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("refs: got %#v, want %#v", got, want)
		}
		return nil
	}

	cmd := extractCommandNode(n)
	refs := make(map[string]string)
	for k, p := range cmd.ParamNodes {
		if ref, ok := p.(ast.RefNode); ok {
			refs[k] = ref.Ref()
		}
	}
	return compare(refs, expected)
}

func assertHoleKeys(n ast.Node, expected map[string]string) error {
	compare := func(got, want map[string]string) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("holes: got %#v, want %#v", got, want)
		}
		return nil
	}

	cmd := extractCommandNode(n)
	holes := make(map[string]string)
	for k, p := range cmd.ParamNodes {
		if hole, ok := p.(ast.HoleNode); ok {
			holes[k] = hole.Hole()
		}
	}
	return compare(holes, expected)
}

func assertVariableDeclarationNode(n ast.Node, expIdent string, value interface{}, expHole string) error {
	if err := isDeclarationNode(n); err != nil {
		return err
	}

	decl := n.(*ast.DeclarationNode)
	if got, want := decl.Ident, expIdent; got != want {
		return fmt.Errorf("ident: got '%s' want '%s'", got, want)
	}
	if err := isRightExpressionNode(decl.Expr); err != nil {
		return err
	}
	val := decl.Expr.(*ast.RightExpressionNode)
	if got, want := val.Result(), value; got != want {
		return fmt.Errorf("value: got '%s' want '%s'", got, want)
	}
	if expHole != "" {
		hole, ok := val.Node().(ast.HoleNode)
		if !ok {
			return fmt.Errorf("hole value: expect '%#v': got no hole (%#v)", expHole, val.Node())
		}
		if got, want := hole.String(), expHole; !reflect.DeepEqual(got, want) {
			return fmt.Errorf("hole value: got '%#v' want '%#v'", got, want)
		}
	}
	return nil
}

func assertDeclarationNode(n ast.Node, expIdent, expAction, expEntity string, refs map[string]string, params map[string]interface{}, holes, aliases map[string]string) error {
	if err := isDeclarationNode(n); err != nil {
		return err
	}

	decl := n.(*ast.DeclarationNode)

	if err := verifyCommandNode(decl.Expr, expAction, expEntity, refs, params, holes, aliases); err != nil {
		return err
	}

	return nil
}

func assertCommandNode(n ast.Node, expAction, expEntity string, refs map[string]string, params map[string]interface{}, holes, aliases map[string]string) error {
	return verifyCommandNode(n, expAction, expEntity, refs, params, holes, aliases)
}

func verifyCommandNode(n ast.Node, expAction, expEntity string, refs map[string]string, params map[string]interface{}, holes, aliases map[string]string) error {
	if err := isCommandNode(n); err != nil {
		return err
	}

	expr := n.(*ast.CommandNode)

	if got, want := expr.Action, expAction; got != want {
		return fmt.Errorf("action: got '%s' want '%s'", got, want)
	}
	if got, want := expr.Entity, expEntity; got != want {
		return fmt.Errorf("entity: got '%s' want '%s'", got, want)
	}

	if err := assertParams(n, params); err != nil {
		return err
	}

	if err := assertAliases(n, aliases); err != nil {
		return err
	}

	if err := assertRefs(n, refs); err != nil {
		return err
	}

	if err := assertHoleKeys(n, holes); err != nil {
		return err
	}

	return nil
}

func extractCommandNode(n ast.Node) *ast.CommandNode {
	msg := func(i interface{}) string {
		return fmt.Sprintf("extracting node: want CommandNode, got %T", i)
	}
	switch n.(type) {
	case *ast.CommandNode:
		return n.(*ast.CommandNode)
	case *ast.DeclarationNode:
		expr := n.(*ast.DeclarationNode).Expr
		switch expr.(type) {
		case *ast.CommandNode:
			return expr.(*ast.CommandNode)
		default:
			panic(msg(expr))
		}
	default:
		panic(msg(n))
	}
}

func extractCmdNode(n ast.Node) *ast.CommandNode {
	msg := func(i interface{}) string {
		return fmt.Sprintf("extracting node: want ActionNode, got %T", i)
	}
	switch n.(type) {
	case *ast.CommandNode:
		return n.(*ast.CommandNode)
	case *ast.DeclarationNode:
		expr := n.(*ast.DeclarationNode).Expr
		switch expr.(type) {
		case *ast.CommandNode:
			return expr.(*ast.CommandNode)
		default:
			panic(msg(expr))
		}
	default:
		panic(msg(n))
	}
}

func isCommandNode(n ast.Node) error {
	switch n.(type) {
	case *ast.CommandNode:
	default:
		return errors.New("expected expression node")
	}
	return nil
}

func isDeclarationNode(n ast.Node) error {
	switch n.(type) {
	case *ast.DeclarationNode:
	default:
		return errors.New("expected declaration node")
	}
	return nil
}

func isRightExpressionNode(n ast.Node) error {
	switch n.(type) {
	case *ast.RightExpressionNode:
	default:
		return fmt.Errorf("expected right expression node, got %#v", n)
	}
	return nil
}
