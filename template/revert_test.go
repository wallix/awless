package template

import (
	"errors"
	"strings"
	"testing"

	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/template/internal/ast"
)

func TestRevertOneliner(t *testing.T) {
	tcases := []struct {
		in, exp   string
		cmdResult interface{}
	}{
		{in: "create instanceprofile name=stuff", exp: "delete instanceprofile name=stuff"},
		{in: "delete instanceprofile name=stuff", exp: "create instanceprofile name=stuff"},
		{in: "create appscalingtarget dimension=dim max-capacity=10 min-capacity=4 resource=res role=role service-namespace=ecs", exp: "delete appscalingtarget dimension=dim resource=res service-namespace=ecs"},
		{in: "create appscalingpolicy dimension=my-dim name=my-name resource=my-res service-namespace=my-ns stepscaling-adjustment-type=my-sat stepscaling-adjustments=[0:0.25:-1,0.25:0.75:0,0.75::+1] type=my-type", exp: "delete appscalingpolicy dimension=my-dim name=my-name resource=my-res service-namespace=my-ns"},
		{in: "attach containertask image=toto memory-hard-limit=64 container-name=test-container name=test-service", exp: "detach containertask container-name=test-container name=test-service"},
		{in: "update securitygroup cidr=0.0.0.0/0 id=sg-12345 inbound=authorize portrange=443 protocol=tcp", exp: "update securitygroup cidr=0.0.0.0/0 id=sg-12345 inbound=revoke portrange=443 protocol=tcp"},
		{in: "update securitygroup cidr=0.0.0.0/0 id=sg-12345 outbound=revoke portrange=443 protocol=tcp", exp: "update securitygroup cidr=0.0.0.0/0 id=sg-12345 outbound=authorize portrange=443 protocol=tcp"},
		{in: "attach mfadevice id=my-mfa-device-id user=toto mfa-code-1=1234 mfa-code-2=2345", exp: "detach mfadevice id=my-mfa-device-id user=toto"},
		{in: "detach mfadevice id=my-mfa-device-id user=toto", exp: "attach mfadevice id=my-mfa-device-id user=toto"},

		{in: "stop instance ids=inst-id-1", exp: "check instance id=inst-id-1 state=stopped timeout=180\nstart instance ids=inst-id-1", cmdResult: "inst-id-1"},
		{in: "start instance ids=inst-id-1", exp: "check instance id=inst-id-1 state=running timeout=180\nstop instance ids=inst-id-1", cmdResult: "inst-id-1"},

		{in: "stop instance ids=inst-id-1,inst-id-2", exp: "check instance id=inst-id-1 state=stopped timeout=180\ncheck instance id=inst-id-2 state=stopped timeout=180\nstart instance ids=[inst-id-1,inst-id-2]", cmdResult: "inst-id-1"},
		{in: "start instance ids=inst-id-1,inst-id-2", exp: "check instance id=inst-id-1 state=running timeout=180\ncheck instance id=inst-id-2 state=running timeout=180\nstop instance ids=[inst-id-1,inst-id-2]", cmdResult: "inst-id-1"},

		{in: "stop database id=my-db-id", exp: "start database id=my-db-id"},
		{in: "start database id=my-db-id", exp: "stop database id=my-db-id"},

		{in: "create instanceprofile name='my funny name with spaces'", exp: "delete instanceprofile name='my funny name with spaces'"},
		{in: "create appscalingtarget dimension=dim max-capacity=10 min-capacity=4 resource=['one res','two','three','4', 5, '4.3', 5.1] role=role service-namespace=ecs", exp: "delete appscalingtarget dimension=dim resource=['one res',two,three,'4',5,'4.3',5.1] service-namespace=ecs"},

		{in: "create classicloadbalancer name=my-classic-loadb", exp: "delete classicloadbalancer name=my-classic-loadb", cmdResult: "my-classic-loadb"},
	}

	for _, tcase := range tcases {
		parsed := MustParse(tcase.in)
		if tcase.cmdResult != nil {
			parsed.CommandNodesIterator()[0].CmdResult = tcase.cmdResult
		}
		reverted, err := parsed.Revert()
		if err != nil {
			t.Fatalf("case '%s': %s", tcase.in, err)
		}
		if got, want := reverted.String(), tcase.exp; got != want {
			t.Fatalf("got\n%q\n\nwant\n%q\n", got, want)
		}
	}
}

func TestRevertTemplate(t *testing.T) {
	env := NewEnv().WithLookupCommandFunc(func(tokens ...string) interface{} {
		return awsspec.MockAWSSessionFactory.Build(strings.Join(tokens, ""))()
	}).Build()
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
		compiled, _, err := Compile(MustParse("attach policy arn=stuff user=mrT\ncreate vpc cidr=10.0.0.0/16\ncreate subnet vpc=vpc-1234 cidr=10.0.0.0/24\nstart instance ids=i-54g3hj\ncreate tag key=Key resource=myinst value=Value\ncreate instance count=1 image=ami-1234 name=myinstance subnet=sub-1234 type=t2.nano"), env, NewRunnerCompileMode)
		if err != nil {
			t.Fatal(err)
		}
		for i, cmd := range compiled.CommandNodesIterator() {
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

		reverted, err := compiled.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := "delete tag key=Key resource=myinst value=Value\ncheck instance id=i-54g3hj state=running timeout=180\nstop instance ids=i-54g3hj\ndelete subnet id=sub-12345\ndelete vpc id=vpc-12345\ndetach policy arn=stuff user=mrT"
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: \n%s\n\nwant:\n%s\n", got, want)
		}
	})

	t.Run("Revert load balancer creation", func(t *testing.T) {
		tpl := MustParse(`
loadbalancerfw = create securitygroup
lb = create loadbalancer groups=$loadbalancerfw name=loadbalancer
create listener actiontype=forward loadbalancer=$lb
inst1 = create instance
`)
		for i, cmd := range tpl.CommandNodesIterator() {
			if i == 0 {
				cmd.CmdResult = "securitygroup-1"
			}
			if i == 1 {
				cmd.CmdResult = "lb-1"
			}
			if i == 2 {
				cmd.CmdResult = "list-1"
			}
			if i == 3 {
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
check securitygroup id=securitygroup-1 state=unused timeout=300
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

		exp := `check securitygroup id=sg-54321 state=unused timeout=300
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
check scalinggroup count=0 name=my-scalinggroup timeout=600
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
check database id=my-database state=not-found timeout=900
delete dbsubnetgroup name=my-dbsubgroup`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert start containertask type=service", func(t *testing.T) {
		tpl := MustParse("start containertask cluster=cl desired-count=2 name=taskname deployment-name=dpname type=service")
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `update containertask cluster=cl deployment-name=dpname desired-count=0
stop containertask cluster=cl deployment-name=dpname type=service`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert start containertask type=task", func(t *testing.T) {
		tpl := MustParse("start containertask type=task cluster=cl desired-count=2 name=taskname")
		tpl.CommandNodesIterator()[0].CmdResult = "my-task-arn"
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `stop containertask cluster=cl run-arn=my-task-arn type=task`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert create certificate", func(t *testing.T) {
		tpl := MustParse("create certificate domain=test.awless.io")
		tpl.CommandNodesIterator()[0].CmdResult = "my-certificate-arn"
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete certificate arn=my-certificate-arn`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})

	t.Run("Revert create policy", func(t *testing.T) {
		tpl := MustParse("create policy name=mypolicy")
		tpl.CommandNodesIterator()[0].CmdResult = "my-policy-arn"
		reverted, err := tpl.Revert()
		if err != nil {
			t.Fatal(err)
		}

		exp := `delete policy all-versions=true arn=my-policy-arn`
		if got, want := reverted.String(), exp; got != want {
			t.Fatalf("got: %s\nwant: %s\n", got, want)
		}
	})
}

func TestCmdNodeIsRevertible(t *testing.T) {
	tcases := []struct {
		line, result string
		params       map[string]interface{}
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
		{line: "start containertask", params: map[string]interface{}{"type": "service"}, revertible: true},
		{line: "start containertask", params: map[string]interface{}{"type": "task"}, revertible: true},
	}

	for _, tc := range tcases {
		splits := strings.SplitN(tc.line, " ", 2)
		action, entity := splits[0], splits[1]
		cmd := &ast.CommandNode{Action: action, Entity: entity, CmdResult: tc.result, CmdErr: tc.err}
		if tc.params != nil {
			cmd.ParamNodes = tc.params
		}
		if tc.revertible != isRevertible(cmd) {
			t.Fatalf("expected '%s' to have revertible=%t", cmd, tc.revertible)
		}
	}
}
