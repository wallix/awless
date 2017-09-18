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

package awsdriver

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/wallix/awless/template/driver"
)

func TestLookupFn(t *testing.T) {
	driv := NewEc2Driver(&mockEc2{})

	t.Run("Known function", func(t *testing.T) {
		driverFn, err := driv.Lookup("create", "vpc")
		if err != nil {
			t.Fatal(err)
		}
		if got, want := reflect.ValueOf(driverFn).Pointer(), reflect.ValueOf(driv.(*Ec2Driver).Create_Vpc).Pointer(); got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
		driv.SetDryRun(true)
		driverFn, err = driv.Lookup("create", "vpc")
		if err != nil {
			t.Fatal(err)
		}
		if got, want := reflect.ValueOf(driverFn).Pointer(), reflect.ValueOf(driv.(*Ec2Driver).Create_Vpc_DryRun).Pointer(); got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("unKnown function", func(t *testing.T) {
		_, err := driv.Lookup("unknown", "function")
		if got, want := err, driver.ErrDriverFnNotFound; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})
}

func TestDriver(t *testing.T) {
	awsMock := &mockEc2{}
	driv := NewEc2Driver(awsMock).(*Ec2Driver)

	t.Run("Create vpc", func(t *testing.T) {
		cidr := "10.0.0.0/16"

		awsMock.verifyVpcInput = func(input *ec2.CreateVpcInput) error {
			if got, want := aws.StringValue(input.CidrBlock), cidr; got != want {
				return fmt.Errorf("got '%s', want '%s'", got, want)
			}
			return nil
		}

		id, err := driv.Create_Vpc(driver.EmptyContext, map[string]interface{}{"cidr": cidr})
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

		id, err := driv.Create_Subnet(driver.EmptyContext, map[string]interface{}{"cidr": cidr, "vpc": vpc})
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

		id, err := driv.Create_Instance(driver.EmptyContext, map[string]interface{}{"image": image, "type": typ, "subnet": subnet, "count": count, "name": name})
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

	t.Run("Create policy", func(t *testing.T) {
		awsMockIam := &mockIam{}
		iamDriv := NewIamDriver(awsMockIam).(*IamDriver)
		policyName := "AwlessInfraReadonlyPolicy"
		policyDesc := "Readonly access to infra resources"
		expectedPolicyDocument := `{
 "Version": "2012-10-17",
 "Statement": [
  {
   "Effect": "Allow",
   "Action": [
    "ec2:Describe*",
    "autoscaling:Describe*",
    "elasticloadbalancing:Describe*"
   ],
   "Resource": "*"
  }
 ]
}`

		awsMockIam.verifyCreatePolicyInput = func(input *iam.CreatePolicyInput) error {
			if got, want := aws.StringValue(input.PolicyName), policyName; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			if got, want := aws.StringValue(input.Description), policyDesc; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			if got, want := aws.StringValue(input.PolicyDocument), expectedPolicyDocument; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			return nil
		}
		id, err := iamDriv.Create_Policy(driver.EmptyContext, map[string]interface{}{
			"name":        policyName,
			"effect":      "Allow",
			"action":      []interface{}{"ec2:Describe*", "autoscaling:Describe*", "elasticloadbalancing:Describe*"},
			"resource":    "*",
			"description": policyDesc,
		})
		if err != nil {
			t.Fatal(err)
		}
		if got, want := id.(string), "mynewpolicy"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("Update policy", func(t *testing.T) {
		policyArn := "arn:aws:iam::0123456789:policy/AwlessInfraReadonlyPolicy"
		previousPolicyDocument := `{
 "Version": "2012-10-17",
 "Statement": [
  {
   "Effect": "Allow",
   "Action": [
    "ec2:AttachVolume",
    "ec2:DetachVolume"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:instance/*",
   "Condition": {
    "StringEquals": {
     "ec2:ResourceTag/department": "dev"
    }
   }
  },
  {
   "Effect": "Allow",
   "Action": [
    "ec2:AttachVolume",
    "ec2:DetachVolume"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:volume/*",
   "Condition": {
    "StringEquals": {
     "ec2:ResourceTag/volume_user": "${aws:username}"}
    }
  }
 ]
}`

		awsMockIam := &mockIam{policyVersions: []*iam.PolicyVersion{
			{
				Document:         aws.String("not this one"),
				IsDefaultVersion: aws.Bool(false),
				VersionId:        aws.String("v1"),
			},
			{
				Document:         aws.String(url.QueryEscape(previousPolicyDocument)),
				IsDefaultVersion: aws.Bool(true),
				VersionId:        aws.String("v2"),
			},
		}}
		iamDriv := NewIamDriver(awsMockIam).(*IamDriver)

		expectedNewVersionPolicyDocument := `{
 "Version": "2012-10-17",
 "Statement": [
  {
   "Effect": "Allow",
   "Action": [
    "ec2:AttachVolume",
    "ec2:DetachVolume"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:instance/*",
   "Condition": {
    "StringEquals": {
     "ec2:ResourceTag/department": "dev"
    }
   }
  },
  {
   "Effect": "Allow",
   "Action": [
    "ec2:AttachVolume",
    "ec2:DetachVolume"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:volume/*",
   "Condition": {
    "StringEquals": {
     "ec2:ResourceTag/volume_user": "${aws:username}"
    }
   }
  },
  {
   "Effect": "Deny",
   "Action": [
    "ec2:AttachVolume",
    "DescribeVolumeAttribute"
   ],
   "Resource": "arn:aws:ec2:eu-west-1:0123456789:volume/*"
  }
 ]
}`

		awsMockIam.verifyCreatePolicyVersionInput = func(input *iam.CreatePolicyVersionInput) error {
			if got, want := aws.StringValue(input.PolicyArn), policyArn; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			if got, want := aws.StringValue(input.PolicyDocument), expectedNewVersionPolicyDocument; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
			if got, want := aws.BoolValue(input.SetAsDefault), true; got != want {
				t.Fatalf("got %t, want %t", got, want)
			}
			return nil
		}
		out, err := iamDriv.Update_Policy(driver.EmptyContext, map[string]interface{}{
			"arn":      policyArn,
			"effect":   "Deny",
			"action":   []interface{}{"ec2:AttachVolume", "DescribeVolumeAttribute"},
			"resource": "arn:aws:ec2:eu-west-1:0123456789:volume/*",
		})
		if err != nil {
			t.Fatal(err)
		}
		if out != nil {
			t.Fatalf("got %#v, expected nil output", out)
		}
	})

	t.Run("Request certificate", func(t *testing.T) {
		awsMockAcm := &mockACM{}
		acmDriv := NewAcmDriver(awsMockAcm).(*AcmDriver)

		tcases := []struct {
			domains                    []interface{}
			validationDomains          []interface{}
			expDomainName              string
			expSubjectAlternativeNames []string
			expDomainValidationOptions []*acm.DomainValidationOption
			expErr                     error
		}{
			{
				domains:                    []interface{}{"my.domain.1", "my.domain.2", "my.domain.3"},
				validationDomains:          []interface{}{"domain.1", "2"},
				expDomainName:              "my.domain.1",
				expSubjectAlternativeNames: []string{"my.domain.2", "my.domain.3"},
				expDomainValidationOptions: []*acm.DomainValidationOption{
					{DomainName: aws.String("my.domain.1"), ValidationDomain: aws.String("domain.1")},
					{DomainName: aws.String("my.domain.2"), ValidationDomain: aws.String("2")},
				},
			},
			{
				domains:                    []interface{}{"my.domain.1", "my.domain.2", "my.domain.3"},
				validationDomains:          []interface{}{},
				expDomainName:              "my.domain.1",
				expSubjectAlternativeNames: []string{"my.domain.2", "my.domain.3"},
			},
			{
				domains:           []interface{}{"my.domain.1"},
				validationDomains: []interface{}{"my.domain.1", "my.domain.2", "my.domain.3"},
				expErr:            fmt.Errorf("there is more validation-domains than certificate domains: [my.domain.1 my.domain.2 my.domain.3]"),
			},

			{
				domains:                    []interface{}{"my.domain.1", "my.domain.2"},
				validationDomains:          []interface{}{"domain.1", "domain.2"},
				expDomainName:              "my.domain.1",
				expSubjectAlternativeNames: []string{"my.domain.2"},
				expDomainValidationOptions: []*acm.DomainValidationOption{
					{DomainName: aws.String("my.domain.1"), ValidationDomain: aws.String("domain.1")},
					{DomainName: aws.String("my.domain.2"), ValidationDomain: aws.String("domain.2")},
				},
			},
		}

		for i, tcase := range tcases {
			awsMockAcm.verifyRequestCertificateInput = func(input *acm.RequestCertificateInput) error {
				if got, want := aws.StringValue(input.DomainName), tcase.expDomainName; got != want {
					t.Fatalf("%d: got %s, want %s", i+1, got, want)
				}
				if got, want := aws.StringValueSlice(input.SubjectAlternativeNames), tcase.expSubjectAlternativeNames; !reflect.DeepEqual(got, want) {
					t.Fatalf("%d: got %#v, want %#v", i+1, got, want)
				}
				if got, want := input.DomainValidationOptions, tcase.expDomainValidationOptions; !reflect.DeepEqual(got, want) {
					t.Fatalf("%d: got %#v, want %#v", i+1, got, want)
				}
				return nil
			}
			id, err := acmDriv.Create_Certificate(driver.EmptyContext, map[string]interface{}{
				"domains":            tcase.domains,
				"validation-domains": tcase.validationDomains,
			})
			if got, want := err, tcase.expErr; !reflect.DeepEqual(got, want) {
				t.Fatalf("%d: got %#v, want %#v", i+1, got, want)
			}
			if err != nil {
				continue
			}
			if got, want := id.(string), "mynewcertificate"; got != want {
				t.Fatalf("%d: got %s, want %s", i+1, got, want)
			}
		}

	})
}

func TestBuildIpPermissionsFromParams(t *testing.T) {
	tcases := []struct {
		params   map[string]interface{}
		expected []*ec2.IpPermission
	}{
		{
			params: map[string]interface{}{
				"protocol":  "tcp",
				"cidr":      "192.168.1.10/24",
				"portrange": 80,
			},
			expected: []*ec2.IpPermission{
				{
					IpProtocol: aws.String("tcp"),
					IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("192.168.1.10/24")}},
					FromPort:   aws.Int64(int64(80)),
					ToPort:     aws.Int64(int64(80)),
				},
			},
		},
		{
			params: map[string]interface{}{
				"protocol": "any",
				"cidr":     "192.168.1.18/32",
			},
			expected: []*ec2.IpPermission{
				{
					IpProtocol: aws.String("-1"),
					IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("192.168.1.18/32")}},
					FromPort:   aws.Int64(int64(-1)),
					ToPort:     aws.Int64(int64(-1)),
				},
			},
		},
		{
			params: map[string]interface{}{
				"protocol":  "udp",
				"cidr":      "0.0.0.0/0",
				"portrange": "22-23",
			},
			expected: []*ec2.IpPermission{
				{
					IpProtocol: aws.String("udp"),
					IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("0.0.0.0/0")}},
					FromPort:   aws.Int64(int64(22)),
					ToPort:     aws.Int64(int64(23)),
				},
			},
		},
		{
			params: map[string]interface{}{
				"protocol":  "icmp",
				"cidr":      "10.0.0.0/16",
				"portrange": "any",
			},
			expected: []*ec2.IpPermission{
				{
					IpProtocol: aws.String("icmp"),
					IpRanges:   []*ec2.IpRange{{CidrIp: aws.String("10.0.0.0/16")}},
					FromPort:   aws.Int64(int64(-1)),
					ToPort:     aws.Int64(int64(-1)),
				},
			},
		},
		{
			params: map[string]interface{}{
				"protocol":      "icmp",
				"securitygroup": "sg-12345",
				"portrange":     "any",
			},
			expected: []*ec2.IpPermission{
				{
					IpProtocol:       aws.String("icmp"),
					UserIdGroupPairs: []*ec2.UserIdGroupPair{{GroupId: aws.String("sg-12345")}},
					FromPort:         aws.Int64(int64(-1)),
					ToPort:           aws.Int64(int64(-1)),
				},
			},
		},

		{
			params: map[string]interface{}{
				"protocol":      "tcp",
				"securitygroup": "sg-23456",
				"portrange":     80,
			},
			expected: []*ec2.IpPermission{
				{
					IpProtocol:       aws.String("tcp"),
					UserIdGroupPairs: []*ec2.UserIdGroupPair{{GroupId: aws.String("sg-23456")}},
					FromPort:         aws.Int64(int64(80)),
					ToPort:           aws.Int64(int64(80)),
				},
			},
		},
	}

	for i, tcase := range tcases {
		ipPermissions, err := buildIpPermissionsFromParams(tcase.params)
		if err != nil {
			t.Fatal(i+1, ":", err)
		}
		if got, want := ipPermissions, tcase.expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %+v, want %+v", i+1, got, want)
		}
	}
}

type mockIam struct {
	iamiface.IAMAPI
	policyVersions                 []*iam.PolicyVersion
	verifyCreatePolicyInput        func(*iam.CreatePolicyInput) error
	verifyCreatePolicyVersionInput func(*iam.CreatePolicyVersionInput) error
}

func (m *mockIam) CreatePolicy(input *iam.CreatePolicyInput) (*iam.CreatePolicyOutput, error) {
	if err := m.verifyCreatePolicyInput(input); err != nil {
		return nil, err
	}
	return &iam.CreatePolicyOutput{Policy: &iam.Policy{Arn: aws.String("mynewpolicy")}}, nil
}

func (m *mockIam) CreatePolicyVersion(input *iam.CreatePolicyVersionInput) (*iam.CreatePolicyVersionOutput, error) {
	if err := m.verifyCreatePolicyVersionInput(input); err != nil {
		return nil, err
	}
	return &iam.CreatePolicyVersionOutput{PolicyVersion: &iam.PolicyVersion{VersionId: aws.String("mynewpolicyversion")}}, nil
}

func (m *mockIam) ListPolicyVersions(input *iam.ListPolicyVersionsInput) (*iam.ListPolicyVersionsOutput, error) {
	return &iam.ListPolicyVersionsOutput{Versions: m.policyVersions}, nil
}

func (m *mockIam) GetPolicyVersion(input *iam.GetPolicyVersionInput) (*iam.GetPolicyVersionOutput, error) {
	for _, v := range m.policyVersions {
		if v.VersionId == input.VersionId {
			return &iam.GetPolicyVersionOutput{PolicyVersion: v}, nil
		}
	}
	return nil, nil
}

type mockS3 struct {
	s3iface.S3API
}

type mockSNS struct {
	snsiface.SNSAPI
}

type mockSQS struct {
	sqsiface.SQSAPI
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
	return &ec2.Reservation{Instances: []*ec2.Instance{{InstanceId: aws.String("mynewinstance")}}}, nil
}

func (m *mockEc2) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	if err := m.verifyTagInput(input); err != nil {
		return nil, err
	}
	return &ec2.CreateTagsOutput{}, nil
}

func (m *mockEc2) CreateTagsRequest(input *ec2.CreateTagsInput) (*request.Request, *ec2.CreateTagsOutput) {
	return &request.Request{Error: m.verifyTagInput(input)}, &ec2.CreateTagsOutput{}
}

type mockACM struct {
	acmiface.ACMAPI
	verifyRequestCertificateInput func(*acm.RequestCertificateInput) error
}

func (m *mockACM) RequestCertificate(input *acm.RequestCertificateInput) (*acm.RequestCertificateOutput, error) {
	if err := m.verifyRequestCertificateInput(input); err != nil {
		return nil, err
	}
	return &acm.RequestCertificateOutput{CertificateArn: aws.String("mynewcertificate")}, nil
}
