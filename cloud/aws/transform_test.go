package aws

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/wallix/awless/graph"

	"github.com/aws/aws-sdk-go/service/ec2"
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
				IpRanges:   []*ec2.IpRange{&ec2.IpRange{CidrIp: awssdk.String("10.192.24.0/24")}},
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

		val, err := extractSliceValues("GroupId")(slice)
		if err != nil {
			t.Fatal(err)
		}
		expectedI := []interface{}{"MyGroup1", "MyGroup2"}
		if got, want := val, expectedI; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		val, err = extractSliceValues("GroupName")(slice)
		if err != nil {
			t.Fatal(err)
		}
		expectedI = []interface{}{"MyGroupName1", "MyGroupName2"}
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
				DestinationIpv6CidrBlock: awssdk.String("0.0.0.0/0"),
				InstanceId:               awssdk.String("test_id_3"),
			},
			{
				DestinationIpv6CidrBlock: awssdk.String("0.0.0.0/0"),
				NatGatewayId:             awssdk.String("test_id_4"),
			},
			{
				DestinationIpv6CidrBlock: awssdk.String("0.0.0.0/0"),
				NetworkInterfaceId:       awssdk.String("test_id_5"),
			},
			{
				DestinationIpv6CidrBlock: awssdk.String("0.0.0.0/0"),
				VpcPeeringConnectionId:   awssdk.String("test_id_6"),
			},
		}

		expected := []*graph.Route{
			{
				Destination: &net.IPNet{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(24, 32)},
				TargetType:  graph.EgressOnlyInternetGatewayTarget,
				Target:      "test_id_1",
			},
			{
				Destination: &net.IPNet{IP: net.ParseIP("fd34:fe56:7891:2f3a::"), Mask: net.CIDRMask(64, 128)},
				TargetType:  graph.GatewayTarget,
				Target:      "test_id_2",
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				TargetType:  graph.InstanceTarget,
				Target:      "test_id_3",
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				TargetType:  graph.NatTarget,
				Target:      "test_id_4",
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				TargetType:  graph.NetworkInterfaceTarget,
				Target:      "test_id_5",
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				TargetType:  graph.VpcPeeringConnectionTarget,
				Target:      "test_id_6",
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
}
