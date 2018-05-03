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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/adler32"
	"net"
	"net/url"
	"reflect"
	"sync"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
)

func InitResource(source interface{}) (*graph.Resource, error) {
	var res *graph.Resource
	switch ss := source.(type) {
	// EC2
	case *ec2.Instance:
		res = graph.InitResource(cloud.Instance, awssdk.StringValue(ss.InstanceId))
	case *ec2.Vpc:
		res = graph.InitResource(cloud.Vpc, awssdk.StringValue(ss.VpcId))
	case *ec2.Subnet:
		res = graph.InitResource(cloud.Subnet, awssdk.StringValue(ss.SubnetId))
	case *ec2.SecurityGroup:
		res = graph.InitResource(cloud.SecurityGroup, awssdk.StringValue(ss.GroupId))
	case *ec2.KeyPairInfo:
		res = graph.InitResource(cloud.Keypair, awssdk.StringValue(ss.KeyName))
	case *ec2.Volume:
		res = graph.InitResource(cloud.Volume, awssdk.StringValue(ss.VolumeId))
	case *ec2.Image:
		res = graph.InitResource(cloud.Image, awssdk.StringValue(ss.ImageId))
	case *ec2.ImportImageTask:
		res = graph.InitResource(cloud.ImportImageTask, awssdk.StringValue(ss.ImportTaskId))
	case *ec2.InternetGateway:
		res = graph.InitResource(cloud.InternetGateway, awssdk.StringValue(ss.InternetGatewayId))
	case *ec2.NatGateway:
		res = graph.InitResource(cloud.NatGateway, awssdk.StringValue(ss.NatGatewayId))
	case *ec2.RouteTable:
		res = graph.InitResource(cloud.RouteTable, awssdk.StringValue(ss.RouteTableId))
	case *ec2.AvailabilityZone:
		res = graph.InitResource(cloud.AvailabilityZone, awssdk.StringValue(ss.ZoneName))
	case *ec2.Address:
		if awssdk.StringValue(ss.AllocationId) != "" {
			res = graph.InitResource(cloud.ElasticIP, awssdk.StringValue(ss.AllocationId))
		} else {
			res = graph.InitResource(cloud.ElasticIP, awssdk.StringValue(ss.PublicIp))
		}
	case *ec2.Snapshot:
		res = graph.InitResource(cloud.Snapshot, awssdk.StringValue(ss.SnapshotId))
	case *ec2.NetworkInterface:
		res = graph.InitResource(cloud.NetworkInterface, awssdk.StringValue(ss.NetworkInterfaceId))
	// Loadbalancer
	case *elb.LoadBalancerDescription:
		res = graph.InitResource(cloud.ClassicLoadBalancer, awssdk.StringValue(ss.LoadBalancerName))
	case *elbv2.LoadBalancer:
		res = graph.InitResource(cloud.LoadBalancer, awssdk.StringValue(ss.LoadBalancerArn))
	case *elbv2.TargetGroup:
		res = graph.InitResource(cloud.TargetGroup, awssdk.StringValue(ss.TargetGroupArn))
	case *elbv2.Listener:
		res = graph.InitResource(cloud.Listener, awssdk.StringValue(ss.ListenerArn))
		// Database
	case *rds.DBInstance:
		res = graph.InitResource(cloud.Database, awssdk.StringValue(ss.DBInstanceIdentifier))
	case *rds.DBSubnetGroup:
		res = graph.InitResource(cloud.DbSubnetGroup, awssdk.StringValue(ss.DBSubnetGroupArn))
		// Autoscaling
	case *autoscaling.LaunchConfiguration:
		res = graph.InitResource(cloud.LaunchConfiguration, awssdk.StringValue(ss.LaunchConfigurationARN))
	case *autoscaling.Group:
		res = graph.InitResource(cloud.ScalingGroup, awssdk.StringValue(ss.AutoScalingGroupARN))
	case *autoscaling.ScalingPolicy:
		res = graph.InitResource(cloud.ScalingPolicy, awssdk.StringValue(ss.PolicyARN))
	// Container
	case *ecr.Repository:
		res = graph.InitResource(cloud.Repository, awssdk.StringValue(ss.RepositoryArn))
	case *ecs.Cluster:
		res = graph.InitResource(cloud.ContainerCluster, awssdk.StringValue(ss.ClusterArn))
	case *ecs.TaskDefinition:
		res = graph.InitResource(cloud.ContainerTask, awssdk.StringValue(ss.TaskDefinitionArn))
	case *ecs.Container:
		res = graph.InitResource(cloud.Container, awssdk.StringValue(ss.ContainerArn))
	case *ecs.ContainerInstance:
		res = graph.InitResource(cloud.ContainerInstance, awssdk.StringValue(ss.ContainerInstanceArn))
		// ACM
	case *acm.CertificateSummary:
		res = graph.InitResource(cloud.Certificate, awssdk.StringValue(ss.CertificateArn))
	// IAM
	case *iam.User:
		res = graph.InitResource(cloud.User, awssdk.StringValue(ss.UserId))
	case *iam.UserDetail:
		res = graph.InitResource(cloud.User, awssdk.StringValue(ss.UserId))
	case *iam.RoleDetail:
		res = graph.InitResource(cloud.Role, awssdk.StringValue(ss.RoleId))
	case *iam.GroupDetail:
		res = graph.InitResource(cloud.Group, awssdk.StringValue(ss.GroupId))
	case *iam.Policy:
		res = graph.InitResource(cloud.Policy, awssdk.StringValue(ss.PolicyId))
	case *iam.ManagedPolicyDetail:
		res = graph.InitResource(cloud.Policy, awssdk.StringValue(ss.PolicyId))
	case *iam.AccessKeyMetadata:
		res = graph.InitResource(cloud.AccessKey, awssdk.StringValue(ss.AccessKeyId))
	case *iam.InstanceProfile:
		res = graph.InitResource(cloud.InstanceProfile, awssdk.StringValue(ss.InstanceProfileId))
	case *iam.VirtualMFADevice:
		res = graph.InitResource(cloud.MFADevice, awssdk.StringValue(ss.SerialNumber))
	// S3
	case *s3.Bucket:
		res = graph.InitResource(cloud.Bucket, awssdk.StringValue(ss.Name))
	case *s3.Object:
		res = graph.InitResource(cloud.S3Object, awssdk.StringValue(ss.Key))
	//SNS
	case *sns.Subscription:
		res = graph.InitResource(cloud.Subscription, awssdk.StringValue(ss.Endpoint))
	case *sns.Topic:
		res = graph.InitResource(cloud.Topic, awssdk.StringValue(ss.TopicArn))
		// DNS
	case *route53.HostedZone:
		res = graph.InitResource(cloud.Zone, awssdk.StringValue(ss.Id))
	case *route53.ResourceRecordSet:
		id := HashFields(awssdk.StringValue(ss.Name), awssdk.StringValue(ss.Type))
		res = graph.InitResource(cloud.Record, id)
		// Lambda
	case *lambda.FunctionConfiguration:
		res = graph.InitResource(cloud.Function, awssdk.StringValue(ss.FunctionArn))
		// Monitoring
	case *cloudwatch.Metric:
		id := HashFields(awssdk.StringValue(ss.Namespace), awssdk.StringValue(ss.MetricName))
		res = graph.InitResource(cloud.Metric, id)
	case *cloudwatch.MetricAlarm:
		res = graph.InitResource(cloud.Alarm, awssdk.StringValue(ss.AlarmArn))
		// cdn
	case *cloudfront.DistributionSummary:
		res = graph.InitResource(cloud.Distribution, awssdk.StringValue(ss.Id))
		// cloudformation
	case *cloudformation.Stack:
		res = graph.InitResource(cloud.Stack, awssdk.StringValue(ss.StackId))
	default:
		return nil, fmt.Errorf("Unknown type of resource %T", source)
	}
	return res, nil
}

func NewResource(source interface{}) (*graph.Resource, error) {
	res, err := InitResource(source)
	if err != nil {
		return res, err
	}

	res.Properties()[properties.ID] = res.Id()

	value := reflect.ValueOf(source)
	if !value.IsValid() || value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not a valid pointer.", value)
	}
	nodeV := value.Elem()

	type keyValResult struct {
		key string
		val interface{}
	}

	resultc := make(chan keyValResult)
	errc := make(chan error)

	var wg sync.WaitGroup
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
						errc <- fmt.Errorf("type [%s]: prop '%v': %s", res.Type(), p, err)
					}
					resultc <- keyValResult{p, val}
				}
			}
			if t.fetch != nil {
				val, err := t.fetch(source)
				if err != nil {
					errc <- fmt.Errorf("type [%s]: prop '%v': %s", res.Type(), p, err)
				}
				resultc <- keyValResult{p, val}
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
		case keyVal, ok := <-resultc:
			if !ok {
				return res, nil
			}
			res.Properties()[keyVal.key] = keyVal.val
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
	switch ii := i.(type) {
	case []*string:
		return awssdk.StringValueSlice(ii), nil

	}
	return nil, fmt.Errorf("extract value: not a pointer but a %T", i)
}

var extractValueAsStringFn = func(i interface{}) (interface{}, error) {
	val, err := extractValueFn(i)
	return fmt.Sprint(val), err
}

// Extract time forcing timezone to UTC (friendlier when running test in different timezones i.e. travis)
var extractTimeFn = func(i interface{}) (interface{}, error) {
	t, ok := i.(*time.Time)
	if ok {
		return t.UTC(), nil
	}
	s, ok := i.(*string)
	if ok {
		t, err := time.Parse("2006-01-02T15:04:05.000+0000", awssdk.StringValue(s))
		if err != nil {
			return nil, err
		}
		return t.UTC(), nil
	}
	return nil, fmt.Errorf("extract time: expected time pointer, got: %T", i)
}

// Extract time that have a Z directly after the time without a space which means UTC
// (https://en.wikipedia.org/wiki/ISO_8601#UTC)
var extractTimeWithZSuffixFn = func(i interface{}) (interface{}, error) {
	t, ok := i.(*time.Time)
	if ok {
		return t.UTC(), nil
	}
	s, ok := i.(*string)
	if ok {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", awssdk.StringValue(s))
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	return nil, fmt.Errorf("extract time: expected time pointer, got: %T", i)
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
		for _, group := range ipPerm.UserIdGroupPairs {
			rule.Sources = append(rule.Sources, awssdk.StringValue(group.GroupId))
		}

		rules = append(rules, rule)
	}
	return rules, nil

}

var extractNameValueFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*cloudwatch.Dimension); !ok {
		return nil, fmt.Errorf("extract ip namevalue: not a dimension slice but a %T", i)
	}
	var nameValues []*graph.KeyValue
	for _, dimension := range i.([]*cloudwatch.Dimension) {
		keyval := &graph.KeyValue{KeyName: awssdk.StringValue(dimension.Name), Value: awssdk.StringValue(dimension.Value)}

		nameValues = append(nameValues, keyval)
	}
	return nameValues, nil
}

var extractECSAttributesFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*ecs.Attribute); !ok {
		return nil, fmt.Errorf("extract ECS attributes: not an attribute slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, attribute := range i.([]*ecs.Attribute) {
		keyval := &graph.KeyValue{KeyName: awssdk.StringValue(attribute.Name), Value: awssdk.StringValue(attribute.Value)}

		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

var extractRouteTableAssociationsFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*ec2.RouteTableAssociation); !ok {
		return nil, fmt.Errorf("extract route table associations: not an association slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, assoc := range i.([]*ec2.RouteTableAssociation) {
		keyval := &graph.KeyValue{KeyName: awssdk.StringValue(assoc.RouteTableAssociationId), Value: awssdk.StringValue(assoc.SubnetId)}
		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

var extractFieldFn = func(field string) transformFn {
	return func(i interface{}) (interface{}, error) {
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("extract field '%s': not a pointer but a %T", field, i)
		}
		struc := value.Elem()
		if struc.Kind() != reflect.Struct {
			return nil, fmt.Errorf("extract field '%s': not a struct pointer but a %T", field, i)
		}

		structField := struc.FieldByName(field)

		if !structField.IsValid() {
			return nil, fmt.Errorf("extract field: field not found: %s", field)
		}

		return extractValueFn(structField.Interface())
	}
}

var extractTagsFn = func(i interface{}) (interface{}, error) {
	var out []string
	switch tags := i.(type) {
	case []*ec2.Tag:
		for _, t := range tags {
			out = append(out, fmt.Sprintf("%s=%s", awssdk.StringValue(t.Key), awssdk.StringValue(t.Value)))
		}
	case []*autoscaling.TagDescription:
		for _, t := range tags {
			out = append(out, fmt.Sprintf("%s=%s", awssdk.StringValue(t.Key), awssdk.StringValue(t.Value)))
		}
	default:
		return nil, fmt.Errorf("extract tags: not a tag slice, but a %T", i)
	}

	return out, nil
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

var extractStringPointerSliceValues = func(i interface{}) (interface{}, error) {
	pointers, ok := i.([]*string)
	if !ok {
		return nil, fmt.Errorf("extract string pointer: not a string pointer slice but a %T", i)
	}
	return awssdk.StringValueSlice(pointers), nil
}

var extractStringSliceValues = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		var res []string
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Slice {
			return nil, fmt.Errorf("extract slice: not a slice but a %T", i)
		}
		for i := 0; i < value.Len(); i++ {
			e, err := extractFieldFn(key)(value.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			str, ok := e.(string)
			if !ok {
				return nil, fmt.Errorf("extract string slice: not a string but a %T", e)
			}
			res = append(res, str)
		}

		return res, nil
	}
}

var extractClassicLoadbListenerDescriptionsFn = func(i interface{}) (interface{}, error) {
	listeners, ok := i.([]*elb.ListenerDescription)
	if !ok {
		return nil, fmt.Errorf("extract classic loadb listener descriptions: unexpected type %T", i)
	}
	var out []string
	for _, d := range listeners {
		if list := d.Listener; list != nil {
			out = append(out, fmt.Sprintf("%s:%d:%s:%d", *list.Protocol, *list.LoadBalancerPort, *list.InstanceProtocol, *list.InstancePort))
		}
	}
	return out, nil
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
		if notEmpty(r.DestinationPrefixListId) {
			route.DestinationPrefixListId = awssdk.StringValue(r.DestinationPrefixListId)
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

var extractDistributionOriginFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.(*cloudfront.Origins); !ok {
		return nil, fmt.Errorf("extract origins: not a origins pointer but a %T", i)
	}
	var origins []*graph.DistributionOrigin
	for _, o := range i.(*cloudfront.Origins).Items {
		origin := &graph.DistributionOrigin{
			ID:         awssdk.StringValue(o.Id),
			PublicDNS:  awssdk.StringValue(o.DomainName),
			PathPrefix: awssdk.StringValue(o.OriginPath),
		}
		if o.S3OriginConfig != nil && awssdk.StringValue(o.S3OriginConfig.OriginAccessIdentity) != "" {
			origin.OriginType = "s3"
			origin.Config = awssdk.StringValue(o.S3OriginConfig.OriginAccessIdentity)
		}

		origins = append(origins, origin)
	}
	return origins, nil
}

var extractStackOutputsFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*cloudformation.Output); !ok {
		return nil, fmt.Errorf("extract ouutputs not an output slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, out := range i.([]*cloudformation.Output) {
		keyval := &graph.KeyValue{KeyName: awssdk.StringValue(out.OutputKey), Value: awssdk.StringValue(out.OutputValue)}

		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

var extractStackParametersFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*cloudformation.Parameter); !ok {
		return nil, fmt.Errorf("extract parameters not a parameter slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, out := range i.([]*cloudformation.Parameter) {
		keyval := &graph.KeyValue{KeyName: awssdk.StringValue(out.ParameterKey), Value: awssdk.StringValue(out.ParameterValue)}

		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

var extractContainersImagesFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]*ecs.ContainerDefinition); !ok {
		return nil, fmt.Errorf("extract containers images, not a container definition slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, out := range i.([]*ecs.ContainerDefinition) {
		keyval := &graph.KeyValue{KeyName: awssdk.StringValue(out.Name), Value: awssdk.StringValue(out.Image)}

		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

func extractDocumentDefaultVersion(i interface{}) (interface{}, error) {
	if _, ok := i.([]*iam.PolicyVersion); !ok {
		return nil, fmt.Errorf("extract default version of document, not a policy version slice but a %T", i)
	}
	for _, version := range i.([]*iam.PolicyVersion) {
		if awssdk.BoolValue(version.IsDefaultVersion) {
			docStr := awssdk.StringValue(version.Document)
			if str, err := url.QueryUnescape(docStr); err == nil {
				var buff bytes.Buffer
				err = json.Compact(&buff, []byte(str))
				return buff.String(), err
			}
			return docStr, nil
		}
	}
	return "", nil
}

func extractURLEncodedJson(i interface{}) (interface{}, error) {
	if _, ok := i.(*string); !ok {
		return nil, fmt.Errorf("extract URL-encoded JSON, not a *string but a %T", i)
	}
	docStr := awssdk.StringValue(i.(*string))
	if str, err := url.QueryUnescape(docStr); err == nil {
		var buff bytes.Buffer
		err = json.Compact(&buff, []byte(str))
		return buff.String(), err
	}
	return docStr, nil
}

func notEmpty(str *string) bool {
	return awssdk.StringValue(str) != ""
}

func HashFields(fields ...interface{}) string {
	var buf bytes.Buffer
	for _, field := range fields {
		buf.WriteString(fmt.Sprint(field))
	}
	h := adler32.New()
	buf.WriteTo(h)
	return "awls-" + hex.EncodeToString(h.Sum(nil))
}
