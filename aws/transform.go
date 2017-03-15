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
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/adler32"
	"net"
	"reflect"
	"sync"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

func initResource(source interface{}) (*graph.Resource, error) {
	var res *graph.Resource
	switch ss := source.(type) {
	// EC2
	case *ec2.Instance:
		res = graph.InitResource(awssdk.StringValue(ss.InstanceId), cloud.Instance)
	case *ec2.Vpc:
		res = graph.InitResource(awssdk.StringValue(ss.VpcId), cloud.Vpc)
	case *ec2.Subnet:
		res = graph.InitResource(awssdk.StringValue(ss.SubnetId), cloud.Subnet)
	case *ec2.SecurityGroup:
		res = graph.InitResource(awssdk.StringValue(ss.GroupId), cloud.SecurityGroup)
	case *ec2.KeyPairInfo:
		res = graph.InitResource(awssdk.StringValue(ss.KeyName), cloud.Keypair)
	case *ec2.Volume:
		res = graph.InitResource(awssdk.StringValue(ss.VolumeId), cloud.Volume)
	case *ec2.InternetGateway:
		res = graph.InitResource(awssdk.StringValue(ss.InternetGatewayId), cloud.InternetGateway)
	case *ec2.RouteTable:
		res = graph.InitResource(awssdk.StringValue(ss.RouteTableId), cloud.RouteTable)
	case *ec2.AvailabilityZone:
		res = graph.InitResource(awssdk.StringValue(ss.ZoneName), cloud.AvailabilityZone)
	// Loadbalancer
	case *elbv2.LoadBalancer:
		res = graph.InitResource(awssdk.StringValue(ss.LoadBalancerArn), cloud.LoadBalancer)
	case *elbv2.TargetGroup:
		res = graph.InitResource(awssdk.StringValue(ss.TargetGroupArn), cloud.TargetGroup)
	case *elbv2.Listener:
		res = graph.InitResource(awssdk.StringValue(ss.ListenerArn), cloud.Listener)
	// IAM
	case *iam.User:
		res = graph.InitResource(awssdk.StringValue(ss.UserId), cloud.User)
	case *iam.UserDetail:
		res = graph.InitResource(awssdk.StringValue(ss.UserId), cloud.User)
	case *iam.RoleDetail:
		res = graph.InitResource(awssdk.StringValue(ss.RoleId), cloud.Role)
	case *iam.GroupDetail:
		res = graph.InitResource(awssdk.StringValue(ss.GroupId), cloud.Group)
	case *iam.Policy:
		res = graph.InitResource(awssdk.StringValue(ss.PolicyId), cloud.Policy)
	case *iam.ManagedPolicyDetail:
		res = graph.InitResource(awssdk.StringValue(ss.PolicyId), cloud.Policy)
	// S3
	case *s3.Bucket:
		res = graph.InitResource(awssdk.StringValue(ss.Name), cloud.Bucket)
	case *s3.Object:
		res = graph.InitResource(awssdk.StringValue(ss.Key), cloud.Object)
	//SNS
	case *sns.Subscription:
		res = graph.InitResource(awssdk.StringValue(ss.Endpoint), cloud.Subscription)
	case *sns.Topic:
		res = graph.InitResource(awssdk.StringValue(ss.TopicArn), cloud.Topic)
		// DNS
	case *route53.HostedZone:
		res = graph.InitResource(awssdk.StringValue(ss.Id), cloud.Zone)
	case *route53.ResourceRecordSet:
		id := hashFields(awssdk.StringValue(ss.Name), awssdk.StringValue(ss.Type))
		res = graph.InitResource(id, cloud.Record)
	default:
		return nil, fmt.Errorf("Unknown type of resource %T", source)
	}
	return res, nil
}

func newResource(source interface{}) (*graph.Resource, error) {
	res, err := initResource(source)
	if err != nil {
		return res, err
	}

	value := reflect.ValueOf(source)
	if !value.IsValid() || value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not a valid pointer.", value)
	}
	nodeV := value.Elem()

	resultc := make(chan graph.Property)
	errc := make(chan error)
	var wg sync.WaitGroup
	res.Properties["Id"] = res.Id()

	for prop, trans := range awsResourcesDef[res.Type()] {
		wg.Add(1)
		go func(p string, t *propertyTransform) {
			defer wg.Done()
			if t.transform != nil {
				sourceField := nodeV.FieldByName(t.name)
				if sourceField.IsValid() && !sourceField.IsNil() {
					val, err := t.transform(sourceField.Interface())
					if err == ErrTagNotFound {
						return
					}
					if err != nil {
						errc <- err
					}
					p := graph.Property{Key: p, Value: val}
					resultc <- p
				}
			}
			if t.fetch != nil {
				val, err := t.fetch(source)
				if err != nil {
					errc <- err
				}
				p := graph.Property{Key: p, Value: val}
				resultc <- p
			}
		}(prop, trans)
	}

	go func() {
		wg.Wait()
		close(errc)
		close(resultc)
	}()

	for {
		select {
		case e := <-errc:
			if e != nil {
				return res, e
			}
		case p, ok := <-resultc:
			if !ok {
				return res, nil
			}
			res.Properties[p.Key] = p.Value
		}
	}

}

var ErrTagNotFound = errors.New("aws tag key not found")

type propertyTransform struct {
	name      string
	transform transformFn
	fetch     fetchFn
}

type transformFn func(i interface{}) (interface{}, error)
type fetchFn func(i interface{}) (interface{}, error)

var extractValueFn = func(i interface{}) (interface{}, error) {
	iv := reflect.ValueOf(i)
	if iv.Kind() == reflect.Ptr {
		if iv.IsNil() {
			return nil, nil
		}
		return iv.Elem().Interface(), nil
	}
	return nil, fmt.Errorf("extract value: not a pointer but a %T", i)
}

// Extract time forcing timezone to UTC (friendlier when running test in different timezones i.e. travis)
var extractTimeFn = func(i interface{}) (interface{}, error) {
	t, ok := i.(*time.Time)
	if !ok {
		return nil, fmt.Errorf("extract time: expected time pointer, got: %T", i)
	}
	return t.UTC(), nil
}

var extractIpPermissionSliceFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*ec2.IpPermission); !ok {
		return nil, fmt.Errorf("extract ip permission: not a permission slice but a %T", i)
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
			return nil, fmt.Errorf("extract field: not a pointer but a %T", i)
		}
		struc := value.Elem()
		if struc.Kind() != reflect.Struct {
			return nil, fmt.Errorf("extract field: not a struct pointer but a %T", i)
		}

		structField := struc.FieldByName(field)

		if !structField.IsValid() {
			return nil, fmt.Errorf("extract field: field not found: %s", field)
		}

		return extractValueFn(structField.Interface())
	}
}

var extractTagFn = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		tags, ok := i.([]*ec2.Tag)
		if !ok {
			return nil, fmt.Errorf("extract tag: not a tag slice, but a %T", i)
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
			return nil, fmt.Errorf("extract slice: not a slice but a %T", i)
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
		return nil, fmt.Errorf("extract route: not a route slice but a %T", i)
	}
	var routes []*graph.Route
	for _, r := range i.([]*ec2.Route) {
		route := &graph.Route{}
		var err error
		if notEmpty(r.DestinationCidrBlock) {
			if _, route.Destination, err = net.ParseCIDR(awssdk.StringValue(r.DestinationCidrBlock)); err != nil {
				return nil, err
			}
		}
		if notEmpty(r.DestinationIpv6CidrBlock) {
			if _, route.DestinationIPv6, err = net.ParseCIDR(awssdk.StringValue(r.DestinationIpv6CidrBlock)); err != nil {
				return nil, err
			}
		}
		if notEmpty(r.EgressOnlyInternetGatewayId) {
			routeTarget := &graph.RouteTarget{Type: graph.EgressOnlyInternetGatewayTarget, Ref: awssdk.StringValue(r.EgressOnlyInternetGatewayId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmpty(r.GatewayId) {
			routeTarget := &graph.RouteTarget{Type: graph.GatewayTarget, Ref: awssdk.StringValue(r.GatewayId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmpty(r.InstanceId) {
			routeTarget := &graph.RouteTarget{Type: graph.InstanceTarget, Ref: awssdk.StringValue(r.InstanceId), Owner: awssdk.StringValue(r.InstanceOwnerId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmpty(r.NatGatewayId) {
			routeTarget := &graph.RouteTarget{Type: graph.NatTarget, Ref: awssdk.StringValue(r.NatGatewayId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmpty(r.NetworkInterfaceId) {
			routeTarget := &graph.RouteTarget{Type: graph.NetworkInterfaceTarget, Ref: awssdk.StringValue(r.NetworkInterfaceId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmpty(r.VpcPeeringConnectionId) {
			routeTarget := &graph.RouteTarget{Type: graph.VpcPeeringConnectionTarget, Ref: awssdk.StringValue(r.VpcPeeringConnectionId)}
			route.Targets = append(route.Targets, routeTarget)
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
			return nil, fmt.Errorf("extract true bool: not a slice but a %T", i)
		}
		for i := 0; i < value.Len(); i++ {
			e, err := extractFieldFn(key)(value.Index(i).Interface())
			if err != nil {
				return res, err
			}
			if e == nil {
				continue //Empty field
			}
			b, ok := e.(bool)
			if !ok {
				return nil, fmt.Errorf("extract true bool: the field %s is not a boolean, but has type: %T", key, e)
			}
			if b {
				res = true
			}
		}

		return res, nil
	}
}

var fetchAndExtractGrantsFn = func(i interface{}) (interface{}, error) {
	b, ok := i.(*s3.Bucket)
	if !ok {
		return nil, fmt.Errorf("fetch grants: not a bucket but a %T", i)
	}

	acls, err := StorageService.(s3iface.S3API).GetBucketAcl(&s3.GetBucketAclInput{Bucket: b.Name})
	if err != nil {
		return nil, err
	}
	var grants []*graph.Grant
	for _, acl := range acls.Grants {
		grant := &graph.Grant{
			Permission:         awssdk.StringValue(acl.Permission),
			GranteeID:          awssdk.StringValue(acl.Grantee.ID),
			GranteeType:        awssdk.StringValue(acl.Grantee.Type),
			GranteeDisplayName: awssdk.StringValue(acl.Grantee.DisplayName),
		}
		if awssdk.StringValue(acl.Grantee.EmailAddress) != "" {
			grant.GranteeDisplayName += "<" + awssdk.StringValue(acl.Grantee.EmailAddress) + ">"
		}
		if grant.GranteeType == "Group" {
			grant.GranteeID += awssdk.StringValue(acl.Grantee.URI)
		}
		grants = append(grants, grant)
	}
	return grants, nil
}

func notEmpty(str *string) bool {
	return awssdk.StringValue(str) != ""
}

func hashFields(fields ...interface{}) string {
	var buf bytes.Buffer
	for _, field := range fields {
		buf.WriteString(fmt.Sprint(field))
	}
	h := adler32.New()
	buf.WriteTo(h)
	return "awls-" + hex.EncodeToString(h.Sum(nil))
}
