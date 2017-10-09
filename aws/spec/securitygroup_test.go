/* Copyright 2017 WALLIX

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

package awsspec

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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
		cmd := &UpdateSecuritygroup{}
		cmd.inject(tcase.params)
		ipPermissions, err := cmd.buildIpPermissions()
		if err != nil {
			t.Fatal(i+1, ":", err)
		}
		if got, want := ipPermissions, tcase.expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %+v, want %+v", i+1, got, want)
		}
	}
}
