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

package awsconv

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/wallix/awless/graph"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
)

func TestTransformFunctions(t *testing.T) {
	t.Parallel()
	t.Run("extractTag", func(t *testing.T) {
		t.Parallel()
		tag := []*ec2.Tag{
			{Key: awssdk.String("Name"), Value: awssdk.String("instance-name")},
			{Key: awssdk.String("Created with"), Value: awssdk.String("awless")},
		}

		val, _ := extractTagFn("Name")(tag)
		if got, want := fmt.Sprint(val), "instance-name"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		val, _ = extractTagFn("Created with")(tag)
		if got, want := fmt.Sprint(val), "awless"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("extractValue", func(t *testing.T) {
		t.Parallel()
		val, _ := extractValueFn(awssdk.String("any"))
		if got, want := val.(string), "any"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}

		val, _ = extractValueFn(awssdk.Int(2))
		if got, want := val.(int), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		val, _ = extractValueFn(awssdk.Int64(4))
		if got, want := val.(int64), int64(4); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		val, _ = extractValueFn(awssdk.Bool(true))
		if got, want := val.(bool), true; got != want {
			t.Fatalf("got %t, want %t", got, want)
		}
	})

	t.Run("extractField", func(t *testing.T) {
		t.Parallel()
		data := &ec2.InstanceState{Code: awssdk.Int64(12), Name: awssdk.String("running")}

		val, _ := extractFieldFn("Code")(data)
		if got, want := val.(int64), int64(12); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		val, _ = extractFieldFn("Name")(data)
		if got, want := val.(string), "running"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("extractClassicLoadbListenerDescriptions", func(t *testing.T) {
		t.Parallel()
		descriptions := []*elb.ListenerDescription{
			{Listener: &elb.Listener{LoadBalancerPort: awssdk.Int64(443), Protocol: awssdk.String("HTTPS"), InstancePort: awssdk.Int64(8080), InstanceProtocol: awssdk.String("HTTP")}},
			{Listener: &elb.Listener{LoadBalancerPort: awssdk.Int64(5000), Protocol: awssdk.String("SSL"), InstancePort: awssdk.Int64(3000), InstanceProtocol: awssdk.String("TCP")}},
		}

		expected := []string{"HTTPS:443:HTTP:8080", "SSL:5000:TCP:3000"}

		i, err := extractClassicLoadbListenerDescriptionsFn(descriptions)
		if err != nil {
			t.Fatal(err)
		}
		res := i.([]string)
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		for i := range expected {
			if got, want := res[i], expected[i]; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}

	})

	t.Run("extractIpPermissions", func(t *testing.T) {
		t.Parallel()
		ipPermissions := []*ec2.IpPermission{
			{FromPort: awssdk.Int64(70),
				ToPort:     awssdk.Int64(85),
				IpProtocol: awssdk.String("udp"),
				IpRanges:   []*ec2.IpRange{},
			},
			{FromPort: awssdk.Int64(-1),
				ToPort:     awssdk.Int64(-1),
				IpProtocol: awssdk.String("-1"),
				IpRanges:   []*ec2.IpRange{{CidrIp: awssdk.String("10.192.24.0/24")}},
			},
			{FromPort: awssdk.Int64(12),
				ToPort:     awssdk.Int64(12),
				IpProtocol: awssdk.String("27"),
				IpRanges: []*ec2.IpRange{
					{CidrIp: awssdk.String("1.2.3.4/32")},
					{CidrIp: awssdk.String("2.3.0.0/16")},
				},
			},
			{FromPort: awssdk.Int64(22),
				ToPort:     awssdk.Int64(22),
				IpProtocol: awssdk.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{CidrIp: awssdk.String("1.2.3.4/32")},
				},
				Ipv6Ranges: []*ec2.Ipv6Range{
					{CidrIpv6: awssdk.String("fd34:fe56:7891:2f3a::/64")},
					{CidrIpv6: awssdk.String("2001:db8::/110")},
				},
			},
			{FromPort: awssdk.Int64(0),
				ToPort:     awssdk.Int64(65535),
				IpProtocol: awssdk.String("tcp"),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{GroupId: awssdk.String("group_1")},
					{GroupId: awssdk.String("group_2")},
				},
			},
		}

		expected := []*graph.FirewallRule{
			{
				PortRange: graph.PortRange{FromPort: int64(70), ToPort: int64(85), Any: false},
				Protocol:  "udp",
				IPRanges:  []*net.IPNet{},
			},
			{
				PortRange: graph.PortRange{Any: true},
				Protocol:  "any",
				IPRanges:  []*net.IPNet{{IP: net.IPv4(10, 192, 24, 0), Mask: net.CIDRMask(24, 32)}},
			},
			{
				PortRange: graph.PortRange{Any: true},
				Protocol:  "27",
				IPRanges: []*net.IPNet{
					{IP: net.IPv4(1, 2, 3, 4), Mask: net.CIDRMask(32, 32)},
					{IP: net.IPv4(2, 3, 0, 0), Mask: net.CIDRMask(16, 32)},
				},
			},
			{
				PortRange: graph.PortRange{FromPort: int64(22), ToPort: int64(22), Any: false},
				Protocol:  "tcp",
				IPRanges: []*net.IPNet{
					{IP: net.IPv4(1, 2, 3, 4), Mask: net.CIDRMask(32, 32)},
					{IP: net.ParseIP("fd34:fe56:7891:2f3a::"), Mask: net.CIDRMask(64, 128)},
					{IP: net.ParseIP("2001:db8::"), Mask: net.CIDRMask(110, 128)},
				},
			},
			{
				PortRange: graph.PortRange{FromPort: int64(0), ToPort: int64(65535), Any: false},
				Protocol:  "tcp",
				Sources: []string{
					"group_1", "group_2",
				},
			},
		}

		i, err := extractIpPermissionSliceFn(ipPermissions)
		if err != nil {
			t.Fatal(err)
		}
		res := i.([]*graph.FirewallRule)
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		for i := range expected {
			if got, want := res[i].String(), expected[i].String(); got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}
	})

	t.Run("extractSliceValues", func(t *testing.T) {
		t.Parallel()
		slice := []*ec2.GroupIdentifier{
			{GroupId: awssdk.String("MyGroup1"), GroupName: awssdk.String("MyGroupName1")},
			{GroupId: awssdk.String("MyGroup2"), GroupName: awssdk.String("MyGroupName2")},
		}

		val, err := extractStringSliceValues("GroupId")(slice)
		if err != nil {
			t.Fatal(err)
		}
		expectedI := []string{"MyGroup1", "MyGroup2"}
		if got, want := val, expectedI; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		val, err = extractStringSliceValues("GroupName")(slice)
		if err != nil {
			t.Fatal(err)
		}
		expectedI = []string{"MyGroupName1", "MyGroupName2"}
		if got, want := val, expectedI; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("extractRoutesSlice", func(t *testing.T) {
		t.Parallel()
		routes := []*ec2.Route{
			{
				DestinationCidrBlock:        awssdk.String("10.0.0.0/24"),
				EgressOnlyInternetGatewayId: awssdk.String("test_id_1"),
			},
			{
				DestinationIpv6CidrBlock: awssdk.String("fd34:fe56:7891:2f3a::/64"),
				GatewayId:                awssdk.String("test_id_2"),
			},
			{
				DestinationCidrBlock: awssdk.String("0.0.0.0/0"),
				InstanceId:           awssdk.String("test_id_3"),
			},
			{
				DestinationCidrBlock: awssdk.String("0.0.0.0/0"),
				NatGatewayId:         awssdk.String("test_id_4"),
			},
			{
				DestinationCidrBlock: awssdk.String("0.0.0.0/0"),
				NetworkInterfaceId:   awssdk.String("test_id_5"),
			},
			{
				DestinationCidrBlock:   awssdk.String("0.0.0.0/0"),
				VpcPeeringConnectionId: awssdk.String("test_id_6"),
			},
			{
				DestinationCidrBlock:     awssdk.String("10.0.0.0/24"),
				DestinationIpv6CidrBlock: awssdk.String("fd34:fe56:7891:2f3a::/64"),
				VpcPeeringConnectionId:   awssdk.String("test_id_7"),
			},
			{
				DestinationPrefixListId: awssdk.String("pl-0123456"),
				GatewayId:               awssdk.String("test_id_8"),
			},
			{
				DestinationCidrBlock:    awssdk.String("0.0.0.0/0"),
				DestinationPrefixListId: awssdk.String("pl-0123456"),
				InstanceId:              awssdk.String("test_id_9"),
				InstanceOwnerId:         awssdk.String("owner"),
				NetworkInterfaceId:      awssdk.String("eni-123456"),
			},
		}

		expected := []*graph.Route{
			{
				Destination: &net.IPNet{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(24, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.EgressOnlyInternetGatewayTarget, Ref: "test_id_1"},
				},
			},
			{
				DestinationIPv6: &net.IPNet{IP: net.ParseIP("fd34:fe56:7891:2f3a::"), Mask: net.CIDRMask(64, 128)},
				Targets: []*graph.RouteTarget{
					{Type: graph.GatewayTarget, Ref: "test_id_2"},
				},
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.InstanceTarget, Ref: "test_id_3"},
				},
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.NatTarget, Ref: "test_id_4"},
				},
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.NetworkInterfaceTarget, Ref: "test_id_5"},
				},
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.VpcPeeringConnectionTarget, Ref: "test_id_6"},
				},
			},
			{
				Destination:     &net.IPNet{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(24, 32)},
				DestinationIPv6: &net.IPNet{IP: net.ParseIP("fd34:fe56:7891:2f3a::"), Mask: net.CIDRMask(64, 128)},
				Targets: []*graph.RouteTarget{
					{Type: graph.VpcPeeringConnectionTarget, Ref: "test_id_7"},
				},
			},
			{
				DestinationPrefixListId: "pl-0123456",
				Targets: []*graph.RouteTarget{
					{Type: graph.GatewayTarget, Ref: "test_id_8"},
				},
			},
			{
				Destination:             &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				DestinationPrefixListId: "pl-0123456",
				Targets: []*graph.RouteTarget{
					{Type: graph.InstanceTarget, Ref: "test_id_9", Owner: "owner"},
					{Type: graph.NetworkInterfaceTarget, Ref: "eni-123456"},
				},
			},
		}

		i, err := extractRoutesSliceFn(routes)
		if err != nil {
			t.Fatal(err)
		}
		res := i.([]*graph.Route)
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		for i := range expected {
			if got, want := res[i].String(), expected[i].String(); got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}
	})

	t.Run("extractHasATrueBoolInStructSlice", func(t *testing.T) {
		t.Parallel()
		slice := []*ec2.RouteTableAssociation{
			{Main: awssdk.Bool(false), RouteTableAssociationId: awssdk.String("test")},
			{RouteTableId: awssdk.String("test2"), Main: awssdk.Bool(true)},
		}

		val, err := extractHasATrueBoolInStructSliceFn("Main")(slice)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(bool), true; got != want {
			t.Fatalf("got %t, want %t", got, want)
		}

		slice = []*ec2.RouteTableAssociation{
			{Main: awssdk.Bool(false), RouteTableAssociationId: awssdk.String("test")},
			{RouteTableId: awssdk.String("test2"), Main: awssdk.Bool(false)},
		}

		val, err = extractHasATrueBoolInStructSliceFn("Main")(slice)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(bool), false; got != want {
			t.Fatalf("got %t, want %t", got, want)
		}
	})
}
