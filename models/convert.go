package models

import (
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Convert(source interface{}) interface{} {
	awsObject := reflect.ValueOf(source).Elem()

	switch awsType := source.(type) {
	case *ec2.Instance:
		instance := &Instance{}
		return populate(awsObject, instance)
	case *ec2.Vpc:
		vpc := &Vpc{}
		return populate(awsObject, vpc)
	case *ec2.Subnet:
		subnet := &Subnet{}
		return populate(awsObject, subnet)
	default:
		fmt.Fprintf(os.Stderr, "struct conversion: aws %T unsupported\n", awsType)
		os.Exit(-1)
	}

	return nil
}

func populate(source reflect.Value, dest interface{}) interface{} {
	destStruct := reflect.ValueOf(dest).Elem()

	for i := 0; i < destStruct.NumField(); i++ {
		destField := destStruct.Field(i)

		if awsTag, ok := destStruct.Type().Field(i).Tag.Lookup("aws"); ok {
			switch destField.Kind() {
			case reflect.String:
				sourceField := source.FieldByName(awsTag)
				if sourceField.IsValid() {
					destField.SetString(aws.StringValue(sourceField.Interface().(*string)))
				}
			default:
				fmt.Fprintf(os.Stderr, "struct conversion: type %s unsupported\n", destField.Kind())
				os.Exit(-1)
			}
		}
	}

	return dest
}
