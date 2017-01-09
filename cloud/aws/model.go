package aws

import (
	"errors"
	"fmt"
	"reflect"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/wallix/awless/rdf"
)

type propertyTransform struct {
	name      string
	transform transformFn
}

type transformFn func(i interface{}) (interface{}, error)

var ErrTagNotFound = errors.New("aws tag key not found")
var ErrFieldNotFound = errors.New("aws struct field not found")
var ErrUnknownType = errors.New("aws type unknown")

var extractValueFn = func(i interface{}) (interface{}, error) {
	switch ii := i.(type) {
	case *string:
		return awssdk.StringValue(ii), nil
	case *int:
		return awssdk.IntValue(ii), nil
	case *int64:
		return awssdk.Int64Value(ii), nil
	default:
		return nil, ErrUnknownType
	}
}

var extractFieldFn = func(field string) transformFn {
	return func(i interface{}) (interface{}, error) {
		value := reflect.ValueOf(i)
		struc := value.Elem()

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

var instanceDef = map[string]*propertyTransform{
	"Id":        {name: "InstanceId", transform: extractValueFn},
	"Name":      {name: "Tags", transform: extractTagFn("Name")},
	"Type":      {name: "InstanceType", transform: extractValueFn},
	"SubnetId":  {name: "SubnetId", transform: extractValueFn},
	"VpcId":     {name: "VpcId", transform: extractValueFn},
	"PublicIp":  {name: "PublicIpAddress", transform: extractValueFn},
	"PrivateIp": {name: "PrivateIpAddress", transform: extractValueFn},
	"ImageId":   {name: "ImageId", transform: extractValueFn},
	"State":     {name: "State", transform: extractFieldFn("Name")},
	"KeyName":   {name: "KeyName", transform: extractValueFn},
}

var awsResourcesProperties = map[string]map[string]string{
	rdf.VPC: {
		"Id":        "VpcId",
		"IsDefault": "IsDefault",
		"State":     "State",
		"CidrBlock": "CidrBlock",
	},
	rdf.SUBNET: {
		"Id":                  "SubnetId",
		"VpcId":               "VpcId",
		"MapPublicIpOnLaunch": "MapPublicIpOnLaunch",
		"State":               "State",
		"CidrBlock":           "CidrBlock",
	},
	rdf.USER: {
		"Id":               "UserId",
		"Name":             "UserName",
		"Arn":              "Arn",
		"Path":             "Path",
		"PasswordLastUsed": "PasswordLastUsed",
	},
	rdf.ROLE: {
		"Id":         "RoleId",
		"Name":       "RoleName",
		"Arn":        "Arn",
		"CreateDate": "CreateDate",
		"Path":       "Path",
	},
	rdf.GROUP: {
		"Id":         "GroupId",
		"Name":       "GroupName",
		"Arn":        "Arn",
		"CreateDate": "CreateDate",
		"Path":       "Path",
	},
	rdf.POLICY: {
		"Id":           "PolicyId",
		"Name":         "PolicyName",
		"Arn":          "Arn",
		"CreateDate":   "CreateDate",
		"UpdateDate":   "UpdateDate",
		"Description":  "Description",
		"IsAttachable": "IsAttachable",
		"Path":         "Path",
	},
}
