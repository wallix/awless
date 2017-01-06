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
	api        ec2iface.EC2API
	references map[string]map[string]string
	logger     *log.Logger
}

func NewDriver(api ec2iface.EC2API) *AwsDriver {
	return &AwsDriver{
		api:        api,
		references: make(map[string]map[string]string),
		logger:     log.New(ioutil.Discard, "", 0),
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

	driverFn, converted := method.(func(map[string]interface{}) error)
	if !converted {
		panic(fmt.Sprintf("method '%s' found on '%T' is not a driver function", fnName, d))
	}

	return driverFn
}

func (d *AwsDriver) Create_Vpc(params map[string]interface{}) error {
	input := &ec2.CreateVpcInput{}

	setField(d.lookupValue(driver.CIDR, params), input, "CidrBlock")

	output, err := d.api.CreateVpc(input)
	if err != nil {
		d.logger.Printf("error creating vpc\n%s\n", err)
		return err
	}
	d.logger.Println("vpc created")

	if refname, ok := params[driver.REF]; ok {
		d.logger.Printf("vpc referenced with '%s'", refname)
		d.addReference(driver.VPC, refname, aws.StringValue(output.Vpc.VpcId))
	}

	return nil
}

func (d *AwsDriver) Create_Subnet(params map[string]interface{}) error {
	input := &ec2.CreateSubnetInput{}

	setField(d.lookupValue(driver.CIDR, params), input, "CidrBlock")
	setField(d.lookupValue(driver.VPC, params), input, "VpcId")

	output, err := d.api.CreateSubnet(input)
	if err != nil {
		d.logger.Printf("error creating subnet\n%s\n", err)
		return err
	}
	d.logger.Println("subnet created")

	if refname, ok := params[driver.REF]; ok {
		d.logger.Printf("subnet referenced with '%s'", refname)
		d.addReference(driver.SUBNET, refname, aws.StringValue(output.Subnet.SubnetId))
	}

	return nil
}

func (d *AwsDriver) Create_Instance(params map[string]interface{}) error {
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-9398d3e0"),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		InstanceType: aws.String("t2.micro"),
	}

	setField(d.lookupValue(driver.BASE, params), input, "ImageId")
	setField(d.lookupValue(driver.TYPE, params), input, "InstanceType")
	setField(d.lookupValue(driver.COUNT, params), input, "MaxCount")
	setField(d.lookupValue(driver.COUNT, params), input, "MinCount")
	setField(d.lookupValue(driver.SUBNET, params), input, "SubnetId")

	output, err := d.api.RunInstances(input)
	if err != nil {
		d.logger.Printf("error creating instance\n%s\n", err)
		return err
	}
	d.logger.Println("instance created")

	if refname, ok := params[driver.REF]; ok {
		d.logger.Printf("instance referenced with '%s'", refname)
		d.addReference(driver.INSTANCE, refname, aws.StringValue(output.Instances[0].InstanceId))
	}

	return nil
}

func (d *AwsDriver) lookupValue(tok string, params map[string]interface{}) interface{} {
	if backref, ok := params[driver.REFERENCES]; ok {
		if refs, ok := d.references[tok]; ok {
			return refs[backref.(string)]
		}
	}

	return params[tok]
}

func (d *AwsDriver) addReference(t string, refName interface{}, resourceId string) {
	if d.references[t] == nil {
		d.references[t] = make(map[string]string)
	}

	d.references[t][refName.(string)] = resourceId
}

func humanize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	return strings.ToUpper(s)
}
