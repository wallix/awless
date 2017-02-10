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

package aws

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

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

func TestDriver(t *testing.T) {
	awsMock := &mockEc2{}
	driv := NewDriver(awsMock, &mockIam{}, &mockS3{})

	t.Run("Create vpc", func(t *testing.T) {
		cidr := "10.0.0.0/16"

		awsMock.verifyVpcInput = func(input *ec2.CreateVpcInput) error {
			if got, want := aws.StringValue(input.CidrBlock), cidr; got != want {
				return fmt.Errorf("got '%s', want '%s'", got, want)
			}
			return nil
		}

		id, err := driv.Create_Vpc(map[string]interface{}{"cidr": cidr})
		if err != nil {
			t.Fatal(err)
		}
		if got, want := id.(string), "mynewvpc"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("Create subnet", func(t *testing.T) {
		cidr, vpc := "10.0.0.0/16", "anyvpc"

		awsMock.verifySubnetInput = func(input *ec2.CreateSubnetInput) error {
			if got, want := aws.StringValue(input.CidrBlock), cidr; got != want {
				return fmt.Errorf("got '%s', want '%s'", got, want)
			}
			if got, want := aws.StringValue(input.VpcId), vpc; got != want {
				return fmt.Errorf("got %s, want %s", got, want)
			}
			return nil
		}

		id, err := driv.Create_Subnet(map[string]interface{}{"cidr": cidr, "vpc": vpc})
		if err != nil {
			t.Fatal(err)
		}
		if got, want := id.(string), "mynewsubnet"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("Create instance", func(t *testing.T) {
		countInt := 2
		image, typ, subnet, count, name := "ami-12", "t2.medium", "anysubnet", strconv.Itoa(countInt), "my_instance_name"
		var tagNameCreated bool

		awsMock.verifyInstanceInput = func(input *ec2.RunInstancesInput) error {
			if got, want := aws.StringValue(input.SubnetId), subnet; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			if got, want := aws.StringValue(input.ImageId), image; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			if got, want := aws.StringValue(input.InstanceType), typ; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			if got, want := aws.Int64Value(input.MinCount), int64(countInt); got != want {
				t.Fatalf("got %d, want %d", got, want)
			}
			if got, want := aws.Int64Value(input.MaxCount), int64(countInt); got != want {
				t.Fatalf("got %d, want %d", got, want)
			}
			return nil
		}

		awsMock.verifyTagInput = func(input *ec2.CreateTagsInput) error {
			if got, want := len(input.Tags), 1; got != want {
				t.Fatalf("got %d, want %d", got, want)
			}
			tagNameCreated = true
			if got, want := aws.StringValue(input.Tags[0].Key), "Name"; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			if got, want := aws.StringValue(input.Tags[0].Value), name; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			return nil
		}

		id, err := driv.Create_Instance(map[string]interface{}{"image": image, "type": typ, "subnet": subnet, "count": count, "name": name})
		if err != nil {
			t.Fatal(err)
		}
		if got, want := id.(string), "mynewinstance"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		if got, want := tagNameCreated, true; got != want {
			t.Fatalf("got %t, want %t", got, want)
		}
	})
}

func TestBuildIpPermissionsFromParams(t *testing.T) {
	params := map[string]interface{}{
		"protocol":  "tcp",
		"cidr":      "192.168.1.10/24",
		"portrange": 80,
	}
	expected := []*ec2.IpPermission{
		{
			IpProtocol: aws.String("tcp"),
			IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("192.168.1.10/24")}},
			FromPort:   aws.Int64(int64(80)),
			ToPort:     aws.Int64(int64(80)),
		},
	}
	ipPermissions, err := buildIpPermissionsFromParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := ipPermissions, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	params = map[string]interface{}{
		"protocol": "any",
		"cidr":     "192.168.1.18/32",
	}
	expected = []*ec2.IpPermission{
		{
			IpProtocol: aws.String("-1"),
			IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("192.168.1.18/32")}},
			FromPort:   aws.Int64(int64(-1)),
			ToPort:     aws.Int64(int64(-1)),
		},
	}
	ipPermissions, err = buildIpPermissionsFromParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := ipPermissions, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	params = map[string]interface{}{
		"protocol":  "udp",
		"cidr":      "0.0.0.0/0",
		"portrange": "22-23",
	}
	expected = []*ec2.IpPermission{
		{
			IpProtocol: aws.String("udp"),
			IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("0.0.0.0/0")}},
			FromPort:   aws.Int64(int64(22)),
			ToPort:     aws.Int64(int64(23)),
		},
	}
	ipPermissions, err = buildIpPermissionsFromParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := ipPermissions, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	params = map[string]interface{}{
		"protocol":  "icmp",
		"cidr":      "10.0.0.0/16",
		"portrange": "any",
	}
	expected = []*ec2.IpPermission{
		{
			IpProtocol: aws.String("icmp"),
			IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("10.0.0.0/16")}},
			FromPort:   aws.Int64(int64(-1)),
			ToPort:     aws.Int64(int64(-1)),
		},
	}
	ipPermissions, err = buildIpPermissionsFromParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := ipPermissions, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

type mockIam struct {
	iamiface.IAMAPI
}

type mockS3 struct {
	s3iface.S3API
}

type mockEc2 struct {
	ec2iface.EC2API
	verifyVpcInput      func(*ec2.CreateVpcInput) error
	verifySubnetInput   func(*ec2.CreateSubnetInput) error
	verifyInstanceInput func(*ec2.RunInstancesInput) error
	verifyTagInput      func(*ec2.CreateTagsInput) error
}

func (m *mockEc2) CreateVpc(input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	if err := m.verifyVpcInput(input); err != nil {
		return nil, err
	}
	return &ec2.CreateVpcOutput{Vpc: &ec2.Vpc{VpcId: aws.String("mynewvpc")}}, nil
}

func (m *mockEc2) CreateSubnet(input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	if err := m.verifySubnetInput(input); err != nil {
		return nil, err
	}
	return &ec2.CreateSubnetOutput{Subnet: &ec2.Subnet{SubnetId: aws.String("mynewsubnet")}}, nil
}

func (m *mockEc2) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if err := m.verifyInstanceInput(input); err != nil {
		return nil, err
	}
	return &ec2.Reservation{Instances: []*ec2.Instance{&ec2.Instance{InstanceId: aws.String("mynewinstance")}}}, nil
}

func (m *mockEc2) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	if err := m.verifyTagInput(input); err != nil {
		return nil, err
	}
	return &ec2.CreateTagsOutput{}, nil
}
