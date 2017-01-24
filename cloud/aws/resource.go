package aws

import (
	"fmt"
	"reflect"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/graph"
)

func NewResource(source interface{}) (*graph.Resource, error) {
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
			res.Properties()[prop] = val
		}
	}

	return res, nil
}

func addCloudResourceToGraph(g *graph.Graph, cloudResource interface{}) error {
	res, err := NewResource(cloudResource)
	if err != nil {
		return err
	}
	triples, err := res.MarshalToTriples()
	if err != nil {
		return err
	}
	g.Add(triples...)
	return nil
}
