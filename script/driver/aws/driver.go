package aws

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/script/driver"
)

type AwsDriver struct {
	api    ec2iface.EC2API
	logger *log.Logger
}

func NewDriver(api ec2iface.EC2API) *AwsDriver {
	return &AwsDriver{
		api:    api,
		logger: log.New(ioutil.Discard, "", 0),
	}
}

func (d *AwsDriver) SetLogger(l *log.Logger) {
	d.logger = l
}

func (d *AwsDriver) Lookup(lookups ...string) driver.DriverFn {
	if len(lookups) < 2 {
		panic("need at least 2 string to lookup driver method")
	}

	fnName := fmt.Sprintf("%s_%s", humanize(lookups[0]), humanize(lookups[1]))
	method := reflect.ValueOf(d).MethodByName(fnName).Interface()

	driverFn, converted := method.(func(map[string]interface{}) (interface{}, error))
	if !converted {
		panic(fmt.Sprintf("method '%s' found on '%T' is not a driver function", fnName, d))
	}

	return driverFn
}

func (d *AwsDriver) Create_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVpcInput{}

	setField(params["cidr"], input, "CidrBlock")

	output, err := d.api.CreateVpc(input)
	if err != nil {
		d.logger.Printf("error creating vpc\n%s\n", err)
		return nil, err
	}
	d.logger.Println("vpc created")

	return aws.StringValue(output.Vpc.VpcId), nil
}

func (d *AwsDriver) Create_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSubnetInput{}

	setField(params["cidr"], input, "CidrBlock")
	setField(params["vpc"], input, "VpcId")

	output, err := d.api.CreateSubnet(input)
	if err != nil {
		d.logger.Printf("error creating subnet\n%s\n", err)
		return nil, err
	}
	d.logger.Println("subnet created")

	return aws.StringValue(output.Subnet.SubnetId), nil
}

func (d *AwsDriver) Create_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-9398d3e0"),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		InstanceType: aws.String("t2.micro"),
	}

	setField(params["base"], input, "ImageId")
	setField(params["type"], input, "InstanceType")
	setField(params["count"], input, "MaxCount")
	setField(params["count"], input, "MinCount")
	setField(params["subnet"], input, "SubnetId")

	output, err := d.api.RunInstances(input)
	if err != nil {
		d.logger.Printf("error creating instance\n%s\n", err)
		return nil, err
	}
	d.logger.Println("instance created")

	return aws.StringValue(output.Instances[0].InstanceId), nil
}

func humanize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	return strings.ToUpper(s)
}
