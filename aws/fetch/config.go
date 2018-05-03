package awsfetch

import (
	"reflect"

	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/wallix/awless/logger"
)

type AWSAPI struct {
	Iam                    iamiface.IAMAPI
	Ec2                    ec2iface.EC2API
	Elbv2                  elbv2iface.ELBV2API
	Elb                    elbiface.ELBAPI
	Rds                    rdsiface.RDSAPI
	Autoscaling            autoscalingiface.AutoScalingAPI
	Ecr                    ecriface.ECRAPI
	Ecs                    ecsiface.ECSAPI
	Applicationautoscaling applicationautoscalingiface.ApplicationAutoScalingAPI
	Sts                    stsiface.STSAPI
	S3                     s3iface.S3API
	Sns                    snsiface.SNSAPI
	Sqs                    sqsiface.SQSAPI
	Route53                route53iface.Route53API
	Lambda                 lambdaiface.LambdaAPI
	Cloudwatch             cloudwatchiface.CloudWatchAPI
	Cloudfront             cloudfrontiface.CloudFrontAPI
	Cloudformation         cloudformationiface.CloudFormationAPI
	Acm                    acmiface.ACMAPI
}

type Config struct {
	Log   *logger.Logger
	Extra map[string]interface{}
	APIs  *AWSAPI
}

func NewConfig(apis ...interface{}) *Config {
	c := &Config{
		Extra: make(map[string]interface{}),
		Log:   logger.DiscardLogger,
	}
	assignAPIs(c, apis...)
	return c
}

func (c *Config) getBoolDefaultTrue(key string) bool {
	if c.Extra == nil {
		return true
	}

	if b, ok := c.Extra[key].(bool); ok {
		return b
	}

	return true
}

func assignAPIs(c *Config, apis ...interface{}) {
	c.APIs = new(AWSAPI)
	val := reflect.ValueOf(c.APIs).Elem()
	stru := val.Type()

	for _, api := range apis {
		if !reflect.ValueOf(api).IsValid() {
			continue
		}

		apiType := reflect.TypeOf(api)
		for i := 0; i < stru.NumField(); i++ {
			fieldType := stru.Field(i).Type
			if apiType.AssignableTo(fieldType) {
				val.Field(i).Set(reflect.ValueOf(api))
				break
			}
		}
	}
}
