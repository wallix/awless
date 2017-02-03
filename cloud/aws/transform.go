package aws

import (
	"errors"
	"fmt"
	"net"
	"reflect"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/graph"
)

func newResource(source interface{}) (*graph.Resource, error) {
	value := reflect.ValueOf(source)
	if !value.IsValid() || value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not a valid pointer.", value)
	}
	nodeV := value.Elem()

	var res *graph.Resource
	switch ss := source.(type) {
	case *ec2.Instance:
		res = graph.InitResource(awssdk.StringValue(ss.InstanceId), graph.Instance)
	case *ec2.Vpc:
		res = graph.InitResource(awssdk.StringValue(ss.VpcId), graph.Vpc)
	case *ec2.Subnet:
		res = graph.InitResource(awssdk.StringValue(ss.SubnetId), graph.Subnet)
	case *ec2.SecurityGroup:
		res = graph.InitResource(awssdk.StringValue(ss.GroupId), graph.SecurityGroup)
	case *ec2.KeyPairInfo:
		res = graph.InitResource(awssdk.StringValue(ss.KeyName), graph.Keypair)
	case *ec2.Volume:
		res = graph.InitResource(awssdk.StringValue(ss.VolumeId), graph.Volume)
	case *ec2.InternetGateway:
		res = graph.InitResource(awssdk.StringValue(ss.InternetGatewayId), graph.InternetGateway)
	case *ec2.RouteTable:
		res = graph.InitResource(awssdk.StringValue(ss.RouteTableId), graph.RouteTable)
	case *iam.User:
		res = graph.InitResource(awssdk.StringValue(ss.UserId), graph.User)
	case *iam.UserDetail:
		res = graph.InitResource(awssdk.StringValue(ss.UserId), graph.User)
	case *iam.Role:
		res = graph.InitResource(awssdk.StringValue(ss.RoleId), graph.Role)
	case *iam.RoleDetail:
		res = graph.InitResource(awssdk.StringValue(ss.RoleId), graph.Role)
	case *iam.Group:
		res = graph.InitResource(awssdk.StringValue(ss.GroupId), graph.Group)
	case *iam.GroupDetail:
		res = graph.InitResource(awssdk.StringValue(ss.GroupId), graph.Group)
	case *iam.Policy:
		res = graph.InitResource(awssdk.StringValue(ss.PolicyId), graph.Policy)
	case *iam.ManagedPolicyDetail:
		res = graph.InitResource(awssdk.StringValue(ss.PolicyId), graph.Policy)
	default:
		return nil, fmt.Errorf("Unknown type of resource %T", source)
	}

	for prop, trans := range awsResourcesDef[res.Type()] {
		sourceField := nodeV.FieldByName(trans.name)
		if sourceField.IsValid() && !sourceField.IsNil() {
			val, err := trans.transform(sourceField.Interface())
			if err == ErrTagNotFound {
				continue
			}
			if err != nil {
				return res, err
			}
			res.Properties[prop] = val
		}
	}

	return res, nil
}

var ErrTagNotFound = errors.New("aws tag key not found")
var ErrFieldNotFound = errors.New("aws struct field not found")

type propertyTransform struct {
	name      string
	transform transformFn
}

type transformFn func(i interface{}) (interface{}, error)

var extractValueFn = func(i interface{}) (interface{}, error) {
	iv := reflect.ValueOf(i)
	if iv.Kind() == reflect.Ptr && !iv.IsNil() {
		return iv.Elem().Interface(), nil
	}
	return nil, fmt.Errorf("aws type unknown: %T", i)
}

var extractIpPermissionSliceFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*ec2.IpPermission); !ok {
		return nil, fmt.Errorf("aws type unknown: %T", i)
	}
	var rules []*graph.FirewallRule
	for _, ipPerm := range i.([]*ec2.IpPermission) {
		rule := &graph.FirewallRule{}

		protocol := awssdk.StringValue(ipPerm.IpProtocol)
		switch protocol {
		case "-1":
			rule.Protocol = "any"
			rule.PortRange = graph.PortRange{Any: true}
		case "tcp", "udp", "icmp", "58":
			rule.Protocol = protocol
			fromPort := awssdk.Int64Value(ipPerm.FromPort)
			toPort := awssdk.Int64Value(ipPerm.ToPort)
			if fromPort == -1 || toPort == -1 {
				rule.PortRange = graph.PortRange{Any: true}
			} else {
				rule.PortRange = graph.PortRange{FromPort: fromPort, ToPort: toPort}
			}

		default:
			rule.Protocol = protocol
			rule.PortRange = graph.PortRange{Any: true}
		}
		for _, r := range ipPerm.IpRanges {
			_, net, err := net.ParseCIDR(awssdk.StringValue(r.CidrIp))
			if err != nil {
				return rules, err
			}
			rule.IPRanges = append(rule.IPRanges, net)
		}
		for _, r := range ipPerm.Ipv6Ranges {
			_, net, err := net.ParseCIDR(awssdk.StringValue(r.CidrIpv6))
			if err != nil {
				return rules, err
			}
			rule.IPRanges = append(rule.IPRanges, net)
		}

		rules = append(rules, rule)
	}
	return rules, nil

}

var extractFieldFn = func(field string) transformFn {
	return func(i interface{}) (interface{}, error) {
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("aws type unknown: %T", i)
		}
		struc := value.Elem()
		if struc.Kind() != reflect.Struct {
			return nil, fmt.Errorf("aws type unknown: %T", i)
		}

		structField := struc.FieldByName(field)

		if !structField.IsValid() {
			return nil, ErrFieldNotFound
		}

		return extractValueFn(structField.Interface())
	}
}

var extractTagFn = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		tags, ok := i.([]*ec2.Tag)
		if !ok {
			return nil, fmt.Errorf("aws model: unexpected type %T", i)
		}
		for _, t := range tags {
			if key == awssdk.StringValue(t.Key) {
				return awssdk.StringValue(t.Value), nil
			}
		}

		return nil, ErrTagNotFound
	}
}

var extractSliceValues = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		var res []interface{}
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Slice {
			return nil, fmt.Errorf("aws type invalid: %T", i)
		}
		for i := 0; i < value.Len(); i++ {
			e, err := extractFieldFn(key)(value.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			res = append(res, e)
		}

		return res, nil
	}
}

var extractRoutesSliceFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*ec2.Route); !ok {
		return nil, fmt.Errorf("aws type unknown: %T", i)
	}
	var routes []*graph.Route
	for _, r := range i.([]*ec2.Route) {
		route := &graph.Route{}
		var err error
		switch {
		case awssdk.StringValue(r.DestinationCidrBlock) != "" && awssdk.StringValue(r.DestinationIpv6CidrBlock) != "":
			return nil, fmt.Errorf("extract values: both IPv4 and IPv6 destination in route %v", r)
		case awssdk.StringValue(r.DestinationCidrBlock) != "":
			_, route.Destination, err = net.ParseCIDR(awssdk.StringValue(r.DestinationCidrBlock))
			if err != nil {
				return nil, err
			}
		case awssdk.StringValue(r.DestinationIpv6CidrBlock) != "":
			_, route.Destination, err = net.ParseCIDR(awssdk.StringValue(r.DestinationIpv6CidrBlock))
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("extract values: no IPv4 nor IPv6 destination in route %v", r)
		}
		switch {
		case notEmpty(r.EgressOnlyInternetGatewayId):
			if notEmpty(r.GatewayId) || notEmpty(r.InstanceId) || notEmpty(r.NatGatewayId) || notEmpty(r.NetworkInterfaceId) || notEmpty(r.VpcPeeringConnectionId) {
				return nil, fmt.Errorf("extract values: multiple non-empty target type in route %v", r)
			}
			route.TargetType = graph.EgressOnlyInternetGatewayTarget
			route.Target = awssdk.StringValue(r.EgressOnlyInternetGatewayId)
		case notEmpty(r.GatewayId):
			if notEmpty(r.EgressOnlyInternetGatewayId) || notEmpty(r.InstanceId) || notEmpty(r.NatGatewayId) || notEmpty(r.NetworkInterfaceId) || notEmpty(r.VpcPeeringConnectionId) {
				return nil, fmt.Errorf("extract values: multiple non-empty target type in route %v", r)
			}
			route.TargetType = graph.GatewayTarget
			route.Target = awssdk.StringValue(r.GatewayId)
		case notEmpty(r.InstanceId):
			if notEmpty(r.EgressOnlyInternetGatewayId) || notEmpty(r.GatewayId) || notEmpty(r.NatGatewayId) || notEmpty(r.NetworkInterfaceId) || notEmpty(r.VpcPeeringConnectionId) {
				return nil, fmt.Errorf("extract values: multiple non-empty target type in route %v", r)
			}
			route.TargetType = graph.InstanceTarget
			route.Target = awssdk.StringValue(r.InstanceId)
		case notEmpty(r.NatGatewayId):
			if notEmpty(r.EgressOnlyInternetGatewayId) || notEmpty(r.GatewayId) || notEmpty(r.InstanceId) || notEmpty(r.NetworkInterfaceId) || notEmpty(r.VpcPeeringConnectionId) {
				return nil, fmt.Errorf("extract values: multiple non-empty target type in route %v", r)
			}
			route.TargetType = graph.NatTarget
			route.Target = awssdk.StringValue(r.NatGatewayId)
		case notEmpty(r.NetworkInterfaceId):
			if notEmpty(r.EgressOnlyInternetGatewayId) || notEmpty(r.GatewayId) || notEmpty(r.InstanceId) || notEmpty(r.NatGatewayId) || notEmpty(r.VpcPeeringConnectionId) {
				return nil, fmt.Errorf("extract values: multiple non-empty target type in route %v", r)
			}
			route.TargetType = graph.NetworkInterfaceTarget
			route.Target = awssdk.StringValue(r.NetworkInterfaceId)
		case notEmpty(r.VpcPeeringConnectionId):
			if notEmpty(r.EgressOnlyInternetGatewayId) || notEmpty(r.GatewayId) || notEmpty(r.InstanceId) || notEmpty(r.NatGatewayId) || notEmpty(r.NetworkInterfaceId) {
				return nil, fmt.Errorf("extract values: multiple non-empty target type in route %v", r)
			}
			route.TargetType = graph.VpcPeeringConnectionTarget
			route.Target = awssdk.StringValue(r.VpcPeeringConnectionId)

		default:
			return nil, fmt.Errorf("extract values: no non-empty target type in route %v", r)
		}

		routes = append(routes, route)
	}
	return routes, nil
}

var extractHasATrueBoolInStructSliceFn = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		var res bool
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Slice {
			return nil, fmt.Errorf("aws type invalid: %T", i)
		}
		for i := 0; i < value.Len(); i++ {
			e, err := extractFieldFn(key)(value.Index(i).Interface())
			if err != nil {
				continue //Empty field, we do not need to throw the error
			}
			b, ok := e.(bool)
			if !ok {
				return nil, fmt.Errorf("this field is not a boolean, but has type: %T", e)
			}
			if b {
				res = true
			}
		}

		return res, nil
	}
}

func notEmpty(str *string) bool {
	return awssdk.StringValue(str) != ""
}
