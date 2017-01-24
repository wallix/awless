package aws

import (
	"fmt"
	"reflect"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/rdf"
)

func NewResource(source interface{}) (*cloud.Resource, error) {
	value := reflect.ValueOf(source)
	if !value.IsValid() || value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not a valid pointer.", value)
	}
	nodeV := value.Elem()

	var res *cloud.Resource
	switch ss := source.(type) {
	case *ec2.Instance:
		res = cloud.InitResource(awssdk.StringValue(ss.InstanceId), rdf.Instance)
	case *ec2.Vpc:
		res = cloud.InitResource(awssdk.StringValue(ss.VpcId), rdf.Vpc)
	case *ec2.Subnet:
		res = cloud.InitResource(awssdk.StringValue(ss.SubnetId), rdf.Subnet)
	case *ec2.SecurityGroup:
		res = cloud.InitResource(awssdk.StringValue(ss.GroupId), rdf.SecurityGroup)
	case *iam.User:
		res = cloud.InitResource(awssdk.StringValue(ss.UserId), rdf.User)
	case *iam.UserDetail:
		res = cloud.InitResource(awssdk.StringValue(ss.UserId), rdf.User)
	case *iam.Role:
		res = cloud.InitResource(awssdk.StringValue(ss.RoleId), rdf.Role)
	case *iam.RoleDetail:
		res = cloud.InitResource(awssdk.StringValue(ss.RoleId), rdf.Role)
	case *iam.Group:
		res = cloud.InitResource(awssdk.StringValue(ss.GroupId), rdf.Group)
	case *iam.GroupDetail:
		res = cloud.InitResource(awssdk.StringValue(ss.GroupId), rdf.Group)
	case *iam.Policy:
		res = cloud.InitResource(awssdk.StringValue(ss.PolicyId), rdf.Policy)
	case *iam.ManagedPolicyDetail:
		res = cloud.InitResource(awssdk.StringValue(ss.PolicyId), rdf.Policy)
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

func addCloudResourceToGraph(g *rdf.Graph, cloudResource interface{}) error {
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
