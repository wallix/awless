package template

import (
	"errors"
	"strings"
	"testing"

	"github.com/wallix/awless/template/internal/ast"
)

func TestRevertOneliner(t *testing.T) {
	tcases := []struct {
		in, exp string
	}{
		{in: "create instanceprofile name=stuff", exp: "delete instanceprofile name=stuff"},
		{in: "delete instanceprofile name=stuff", exp: "create instanceprofile name=stuff"},
		{in: "create appscalingtarget dimension=dim max-capacity=10 min-capacity=4 resource=res role=role service-namespace=ecs", exp: "delete appscalingtarget dimension=dim resource=res service-namespace=ecs"},
		{in: "create appscalingpolicy dimension=my-dim name=my-name resource=my-res service-namespace=my-ns stepscaling-adjustment-type=my-sat stepscaling-adjustments=0:0.25:-1,0.25:0.75:0,0.75::+1 type=my-type", exp: "delete appscalingpolicy dimension=my-dim name=my-name resource=my-res service-namespace=my-ns"},
	}

	for _, tcase := range tcases {
		reverted, err := MustParse(tcase.in).Revert()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := reverted.String(), tcase.exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	}
}

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

		exp := "delete tag key=Key resource=myinst value=Value\ncheck instance id=i-54g3hj state=running timeout=180\nstop instance id=i-54g3hj\ndelete subnet id=sub-12345\ndelete vpc id=vpc-12345\ndetach policy arn=stuff user=mrT"
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
				cmd.CmdResult = "securitygroup-1"
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
check securitygroup id=securitygroup-1 state=unused timeout=180
delete securitygroup id=securitygroup-1`
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

	t.Run("Delete the copy of an image", func(t *testing.T) {
		tpl := MustParse("copy image")
		for _, cmd := range tpl.CommandNodesIterator() {
			cmd.CmdResult = "ami-12345678"
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete image delete-snapshots=true id=ami-12345678`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert detach a volume removes the force param", func(t *testing.T) {
		tpl := MustParse("detach volume device=/dev/sdh force=true id=vol-12345 instance=i-12345")
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `attach volume device=/dev/sdh id=vol-12345 instance=i-12345`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert attach a volume waits that the volume is available", func(t *testing.T) {
		tpl := MustParse("detach volume device=/dev/sdh id=vol-12345 instance=i-12345\nattach volume device=/dev/sdh id=vol-12345 instance=i-12345")
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `detach volume device=/dev/sdh id=vol-12345 instance=i-12345
check volume id=vol-12345 state=available timeout=180
attach volume device=/dev/sdh id=vol-12345 instance=i-12345`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert create route", func(t *testing.T) {
		tpl := MustParse("create route cidr=0.0.0.0/0 gateway=igw-12345 table=rtb-12345")
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete route cidr=0.0.0.0/0 table=rtb-12345`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert attach instance", func(t *testing.T) {
		tpl := MustParse("attach instance id=i-123456 port=80 targetgroup=mytargetgrouparn")
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `detach instance id=i-123456 targetgroup=mytargetgrouparn`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Delete scaling group", func(t *testing.T) {
		tpl := MustParse("create scalinggroup")
		for _, cmd := range tpl.CommandNodesIterator() {
			cmd.CmdResult = "my-scalinggroup"
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `update scalinggroup max-size=0 min-size=0 name=my-scalinggroup
check scalinggroup count=0 name=my-scalinggroup timeout=180
delete scalinggroup force=true name=my-scalinggroup`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert create accesskey", func(t *testing.T) {
		tpl := MustParse("create accesskey user=myuser")
		for _, cmd := range tpl.CommandNodesIterator() {
			cmd.CmdResult = "my-accesskey"
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete accesskey id=my-accesskey user=myuser`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert create queue", func(t *testing.T) {
		tpl := MustParse("create queue name=my-queue")
		for _, cmd := range tpl.CommandNodesIterator() {
			cmd.CmdResult = "my-queue-url"
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete queue url=my-queue-url`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Param with space is quoted", func(t *testing.T) {
		tpl := MustParse("create queue name=my-queue")
		for _, cmd := range tpl.CommandNodesIterator() {
			cmd.CmdResult = "my queue url"
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete queue url='my queue url'`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert create record", func(t *testing.T) {
		tpl := MustParse("create record comment='my test record' name=test.awlesstest.io. ttl=60 type=A value=1.2.3.4 zone=/hostedzone/Z29L20HGD4CX07")
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete record name=test.awlesstest.io. ttl=60 type=A value=1.2.3.4 zone=/hostedzone/Z29L20HGD4CX07`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert create database", func(t *testing.T) {
		tpl := MustParse("dbsubgroup = create dbsubnetgroup\ncreate database subnetgroup=$dbsubgroup")
		for i, cmd := range tpl.CommandNodesIterator() {
			if i == 0 {
				cmd.CmdResult = "my-dbsubgroup"
			}
			if i == 1 {
				cmd.CmdResult = "my-database"
			}
		}
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete database id=my-database skip-snapshot=true
check database id=my-database state=not-found timeout=600
delete dbsubnetgroup name=my-dbsubgroup`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert start containerservice", func(t *testing.T) {
		tpl := MustParse("start containerservice name=awless-test-service cluster=awless-cluster deployment-name=awless-test")
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `update containerservice cluster=awless-cluster deployment-name=awless-test desired-count=0
stop containerservice cluster=awless-cluster deployment-name=awless-test`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert attach containertask", func(t *testing.T) {
		tpl := MustParse("attach containertask image=toto memory-hard-limit=64 container-name=test-container name=test-service")
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `detach containertask container-name=test-container name=test-service`
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
		{line: "copy image", result: "any", revertible: true},
		{line: "detach routetable", revertible: false},
		{line: "start alarm", revertible: true},
		{line: "stop alarm", revertible: true},
		{line: "start containerservice", revertible: true},
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
