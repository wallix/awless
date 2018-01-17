package template_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/internal/ast"
)

func TestDryRun(t *testing.T) {
	env := template.NewEnv().WithLookupCommandFunc(func(tokens ...string) interface{} {
		return awsspec.MockAWSSessionFactory.Build(strings.Join(tokens, ""))()
	}).Build()

	t.Run("return error", func(t *testing.T) {
		tpl := template.MustParse("create instance userdata=/invalid-file count=1 image=ami-123456 name=any subnet=any type=t2.micro")
		_, _, err := template.Compile(tpl, env, template.NewRunnerCompileMode)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tpl.DryRun(template.NewRunEnv(env)); err == nil {
			t.Fatal("expected error got none")
		}
	})
}

func TestParamsProcessing(t *testing.T) {
	env := template.NewEnv().WithLookupCommandFunc(func(tokens ...string) interface{} {
		return awsspec.MockAWSSessionFactory.Build(strings.Join(tokens, ""))()
	}).Build()

	t.Run("unexpected param", func(t *testing.T) {
		tpl := template.MustParse("create instance invalid=any")
		_, _, err := template.Compile(tpl, env, template.NewRunnerCompileMode)
		if err == nil {
			t.Fatal("expected err got none")
		}
		if got, want := err.Error(), "create instance: unexpected param(s): invalid"; !strings.Contains(got, want) {
			t.Fatalf("%s should contain %s", got, want)
		}
	})

	t.Run("normalizing missing required params as holes", func(t *testing.T) {
		tpl := template.MustParse("create instance image=ami-123456")
		compiled, _, _ := template.Compile(tpl, env, template.NewRunnerCompileMode)
		if got, want := compiled.String(), "create instance count={instance.count} image=ami-123456 name={instance.name} subnet={instance.subnet} type={instance.type}"; got != want {
			t.Fatalf("%s should contain %s", got, want)
		}
	})

	t.Run("format validation", func(t *testing.T) {
		tpl := template.MustParse("check instance state=woot id=i-45678 timeout=180")
		_, _, err := template.Compile(tpl, env, template.NewRunnerCompileMode)
		if err == nil {
			t.Fatal("expected err got none")
		}
		if got, want := err.Error(), "expected any of"; !strings.Contains(got, want) {
			t.Fatalf("%s should contain %s", got, want)
		}
	})
}

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
		{
			tpl: `
name = instance-{instance.name}-{version}
name2 = my-test-{hole}
create instance image=ami-1234 name=$name subnet=subnet-{version}
create instance image=ami-1234 name=$name2 subnet=sub1234
`,
			expect: `create instance count=42 image=ami-1234 name=instance-myinstance-10 subnet=subnet-10 type=t2.micro
create instance count=42 image=ami-1234 name=my-test-sub-2345 subnet=sub1234 type=t2.micro`,
			expProcessedFillers:  map[string]interface{}{"instance.name": "myinstance", "version": 10, "instance.type": "t2.micro", "instance.count": 42, "hole": "@sub"},
			expResolvedVariables: map[string]interface{}{"name": "instance-myinstance-10", "name2": "my-test-sub-2345"},
		},
		{
			tpl: `
name = "ins$\ta{nce}-"+{instance.name}+{version}
name2 = {hole}+{hole}+"text-with $Special {char-s"
create instance image=ami-1234 name=$name subnet=subnet-{version}
create instance image=ami-1234 name=$name2 subnet=sub1234
`,
			expect: `create instance count=42 image=ami-1234 name='ins$\ta{nce}-myinstance10' subnet=subnet-10 type=t2.micro
create instance count=42 image=ami-1234 name='sub-2345sub-2345text-with $Special {char-s' subnet=sub1234 type=t2.micro`,
			expProcessedFillers:  map[string]interface{}{"instance.name": "myinstance", "version": 10, "instance.type": "t2.micro", "instance.count": 42, "hole": "@sub"},
			expResolvedVariables: map[string]interface{}{"name": "ins$\\ta{nce}-myinstance10", "name2": "sub-2345sub-2345text-with $Special {char-s"},
		},
		{
			tpl: `
create loadbalancer name=mylb subnets={private.subnets}
`,
			expect:               `create loadbalancer name=mylb subnets=[sub-1234,sub-2345]`,
			expProcessedFillers:  map[string]interface{}{"private.subnets": []interface{}{"sub-1234", "sub-2345"}},
			expResolvedVariables: map[string]interface{}{},
		},
		{
			tpl: `
create loadbalancer name=mylb subnets=subnet-1, subnet-2
`,
			expect:               `create loadbalancer name=mylb subnets=[subnet-1,subnet-2]`,
			expProcessedFillers:  map[string]interface{}{},
			expResolvedVariables: map[string]interface{}{},
		}, //retro-compatibility with old list style, without brackets
	}

	for i, tcase := range tcases {
		cenv := template.NewEnv().WithAliasFunc(func(p, v string) string {
			vals := map[string]string{
				"vpc":      "vpc-1234",
				"subalias": "sub-1111",
				"sub":      "sub-2345",
			}
			return vals[v]
		}).WithLookupCommandFunc(func(tokens ...string) interface{} {
			return awsspec.MockAWSSessionFactory.Build(strings.Join(tokens, ""))()
		}).Build()

		cenv.Push(env.FILLERS, map[string]interface{}{
			"instance.type":   "t2.micro",
			"test.cidr":       "10.0.2.0/24",
			"instance.count":  42,
			"unused":          "filler",
			"backup-subnet":   "sub-0987",
			"mysubnet2.hole":  "mysubnet-2",
			"mysubnet3.hole":  "mysubnet-3",
			"mysubnet5.hole":  "mysubnet-5",
			"version":         10,
			"instance.name":   "myinstance",
			"hole":            ast.NewAliasNode("sub"),
			"private.subnets": ast.NewListNode([]interface{}{"sub-1234", "sub-2345"}),
		})

		inTpl := template.MustParse(tcase.tpl)

		compiled, _, err := template.Compile(inTpl, cenv, template.NewRunnerCompileMode)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}

		if got, want := compiled.String(), tcase.expect; got != want {
			t.Fatalf("%d: got\n%s\nwant\n%s", i+1, got, want)
		}

		if got, want := cenv.Get(env.PROCESSED_FILLERS), tcase.expProcessedFillers; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %v, want %v", i+1, got, want)
		}

		if got, want := cenv.Get(env.RESOLVED_VARS), tcase.expResolvedVariables; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got\n%#v\nwant\n%#v\n", i+1, got, want)
		}
	}
}

func TestExternallyProvidedParams(t *testing.T) {
	tcases := []struct {
		template            string
		externalParams      string
		expect              string
		expProcessedFillers map[string]interface{}
	}{
		{
			template:            `create instance count=1 image=ami-123 name=test subnet={hole.name} type=t2.micro`,
			externalParams:      "hole.name=subnet-2345",
			expect:              `create instance count=1 image=ami-123 name=test subnet=subnet-2345 type=t2.micro`,
			expProcessedFillers: map[string]interface{}{"hole.name": "subnet-2345"},
		},
		{
			template:            `create instance count=1 image=ami-123 name=test subnet={hole.name} type={instance.type}`,
			externalParams:      "instance.type=t2.nano hole.name=@subalias",
			expect:              `create instance count=1 image=ami-123 name=test subnet=subnet-111 type=t2.nano`,
			expProcessedFillers: map[string]interface{}{"hole.name": "@subalias", "instance.type": "t2.nano"},
		},
		{
			template:            `create loadbalancer name=elbv2 subnets={my.subnets}`,
			externalParams:      "my.subnets=[@sub1, @sub2]",
			expect:              `create loadbalancer name=elbv2 subnets=[subnet-123,subnet-234]`,
			expProcessedFillers: map[string]interface{}{"my.subnets": []interface{}{"@sub1", "@sub2"}},
		},
		{
			template:            `create loadbalancer name={my.name} subnets={my.subnets}`,
			externalParams:      "my.subnets=sub1, sub2 my.name=loadbalancername",
			expect:              `create loadbalancer name=loadbalancername subnets=[sub1,sub2]`,
			expProcessedFillers: map[string]interface{}{"my.name": "loadbalancername", "my.subnets": []interface{}{"sub1", "sub2"}},
		}, //retro-compatibility with old list style, without brackets
	}
	for i, tcase := range tcases {
		externalFillters, err := template.ParseParams(tcase.externalParams)
		if err != nil {
			t.Fatal(err)
		}
		cenv := template.NewEnv().WithLookupCommandFunc(func(tokens ...string) interface{} {
			return awsspec.MockAWSSessionFactory.Build(strings.Join(tokens, ""))()
		}).WithAliasFunc(func(p, v string) string {
			vals := map[string]string{
				"subalias": "subnet-111",
				"sub1":     "subnet-123",
				"sub2":     "subnet-234",
			}
			return vals[v]
		}).Build()

		cenv.Push(env.FILLERS, externalFillters)

		inTpl := template.MustParse(tcase.template)

		compiled, _, err := template.Compile(inTpl, cenv, template.NewRunnerCompileMode)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}

		if got, want := compiled.String(), tcase.expect; got != want {
			t.Fatalf("%d: got\n%s\nwant\n%s", i+1, got, want)
		}

		if got, want := cenv.Get(env.PROCESSED_FILLERS), tcase.expProcessedFillers; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %#v, want %#v", i+1, got, want)
		}
	}
}
