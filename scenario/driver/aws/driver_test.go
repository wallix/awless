package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/scenario"
	"github.com/wallix/awless/scenario/driver"
)

func TestRunSimpleScenario(t *testing.T) {
	raw := `CREATE VPC CIDR 10.0.0.0/16 REF vpc_1
CREATE SUBNET REFERENCES vpc_1 REF subnet_1
CREATE INSTANCE COUNT 2 BASE linux REFERENCES subnet_1
`

	verifyVpc := func(input *ec2.CreateVpcInput) {
		if got, want := aws.StringValue(input.CidrBlock), "10.0.0.0/16"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}

	verifySubnet := func(input *ec2.CreateSubnetInput) {
		if got, want := aws.StringValue(input.VpcId), "vpc_1"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}

	verifyInstance := func(input *ec2.RunInstancesInput) {
		if got, want := aws.StringValue(input.SubnetId), "subnet_1"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		if got, want := aws.StringValue(input.ImageId), "linux"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		if got, want := aws.StringValue(input.InstanceType), "t2.micro"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		if got, want := aws.Int64Value(input.MinCount), int64(2); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := aws.Int64Value(input.MaxCount), int64(2); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	}

	lex := &scenario.Lexer{}
	scen := lex.ParseScenario(raw)

	mock := &mockEc2{
		verifyInstance: verifyInstance,
		verifySubnet:   verifySubnet,
		verifyVpc:      verifyVpc,
	}

	runner := &driver.Runner{NewAwsDriver(mock)}

	if err := runner.Run(scen); err != nil {
		t.Fatal(err)
	}
}

type mockEc2 struct {
	ec2iface.EC2API
	verifyInstance func(*ec2.RunInstancesInput)
	verifySubnet   func(*ec2.CreateSubnetInput)
	verifyVpc      func(*ec2.CreateVpcInput)
}

func (m *mockEc2) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	m.verifyInstance(input)
	return nil, nil
}

func (m *mockEc2) CreateSubnet(input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	m.verifySubnet(input)
	return &ec2.CreateSubnetOutput{Subnet: &ec2.Subnet{SubnetId: aws.String("subnet_1")}}, nil
}

func (m *mockEc2) CreateVpc(input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	m.verifyVpc(input)
	return &ec2.CreateVpcOutput{Vpc: &ec2.Vpc{VpcId: aws.String("vpc_1")}}, nil
}

func TestHumanizeString(t *testing.T) {
	if got, want := humanize(""), ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanize("s"), "S"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanize("STUFF"), "Stuff"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := humanize("stuff"), "Stuff"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
