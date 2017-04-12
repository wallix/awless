package template

import (
	"errors"
	"strings"
	"testing"

	"github.com/wallix/awless/template/internal/ast"
)

func TestRevertTemplate(t *testing.T) {
	t.Run("Simple template", func(t *testing.T) {
		tpl := MustParse("create instance type=t2.micro")
		for _, cmd := range tpl.CommandNodesIterator() {
			cmd.CmdResult = "i-54321"
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := "delete instance id=i-54321"
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert an instance creation that is not the first command", func(t *testing.T) {
		tpl := MustParse("create subnet\ncreate instance type=t2.micro")
		for i, cmd := range tpl.CommandNodesIterator() {
			if i == 0 {
				cmd.CmdResult = "i-12345"
			}
			if i == 1 {
				cmd.CmdResult = "i-54321"
			}
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := "delete instance id=i-54321\ncheck instance id=i-54321 state=terminated timeout=180\ndelete subnet id=i-12345"
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("More advanced template", func(t *testing.T) {
		tpl := MustParse("attach policy arn=stuff user=mrT\ncreate vpc\ncreate subnet\nstart instance id=i-54g3hj\ncreate tag key=Key resource=myinst value=Value\ncreate instance")
		for i, cmd := range tpl.CommandNodesIterator() {
			if i == 1 {
				cmd.CmdResult = "vpc-12345"
			}
			if i == 2 {
				cmd.CmdResult = "sub-12345"
			}
			if i == 3 {
				cmd.CmdResult = "i-12345"
			}
			if i == 5 {
				cmd.CmdErr = errors.New("cannot create instance")
			}
		}

		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := "delete tag key=Key resource=myinst value=Value\nstop instance id=i-54g3hj\ndelete subnet id=sub-12345\ndelete vpc id=vpc-12345\ndetach policy arn=stuff user=mrT"
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert load balancer creation", func(t *testing.T) {
		tpl := MustParse(`
loadbalancerfw = create securitygroup
update securitygroup cidr=0.0.0.0/0 id=$loadbalancerfw inbound=authorize portrange=80 protocol=tcp
lb = create loadbalancer groups=$loadbalancerfw name=loadbalancer
create listener actiontype=forward loadbalancer=$lb
inst1 = create instance
`)
		for i, cmd := range tpl.CommandNodesIterator() {
			if i == 0 {
				cmd.CmdResult = "secgroup-1"
			}
			if i == 2 {
				cmd.CmdResult = "lb-1"
			}
			if i == 3 {
				cmd.CmdResult = "list-1"
			}
			if i == 4 {
				cmd.CmdResult = "i-1"
			}
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}
		exp := `delete instance id=i-1
check instance id=i-1 state=terminated timeout=180
delete listener id=list-1
delete loadbalancer id=lb-1
check loadbalancer id=lb-1 state=not-found timeout=180
check securitygroup id=secgroup-1 state=unused timeout=180
delete securitygroup id=secgroup-1`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Delete one securitygroup", func(t *testing.T) {
		tpl := MustParse("create securitygroup")
		for _, cmd := range tpl.CommandNodesIterator() {
			cmd.CmdResult = "sg-54321"
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `check securitygroup id=sg-54321 state=unused timeout=180
delete securitygroup id=sg-54321`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

}

func TestCmdNodeIsRevertible(t *testing.T) {
	tcases := []struct {
		line, result string
		err          error
		revertible   bool
	}{
		{line: "update vpc", result: "any", revertible: false},
		{line: "delete vpc", result: "any", revertible: false},
		{line: "create vpc", result: "any", err: errors.New("any"), revertible: false},
		{line: "create vpc", revertible: false},
		{line: "start instance", revertible: false},
		{line: "create vpc", result: "any", revertible: true},
		{line: "stop instance", result: "any", revertible: true},
		{line: "attach policy", revertible: true},
		{line: "detach policy", revertible: true},
		{line: "create record", revertible: true},
		{line: "delete record", revertible: true},
	}

	for _, tc := range tcases {
		splits := strings.SplitN(tc.line, " ", 2)
		action, entity := splits[0], splits[1]
		cmd := &ast.CommandNode{Action: action, Entity: entity, CmdResult: tc.result, CmdErr: tc.err}
		if tc.revertible != isRevertible(cmd) {
			t.Fatalf("expected '%s' to have revertible=%t", cmd, tc.revertible)
		}
	}
}
