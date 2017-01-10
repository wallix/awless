// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (d *AwsDriver) Create_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVpcInput{}

	setField(params["cidr"], input, "CidrBlock")

	output, err := d.api.CreateVpc(input)
	if err != nil {
		d.logger.Printf("create vpc error: %s", err)
		return nil, err
	}
	d.logger.Println("create vpc done")

	return aws.StringValue(output.Vpc.VpcId), nil
}

func (d *AwsDriver) Delete_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVpcInput{}

	setField(params["id"], input, "VpcId")

	output, err := d.api.DeleteVpc(input)
	if err != nil {
		d.logger.Printf("delete vpc error: %s", err)
		return nil, err
	}
	d.logger.Println("delete vpc done")

	return output, nil
}

func (d *AwsDriver) Create_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSubnetInput{}

	setField(params["cidr"], input, "CidrBlock")
	setField(params["vpc"], input, "VpcId")

	output, err := d.api.CreateSubnet(input)
	if err != nil {
		d.logger.Printf("create subnet error: %s", err)
		return nil, err
	}
	d.logger.Println("create subnet done")

	return aws.StringValue(output.Subnet.SubnetId), nil
}

func (d *AwsDriver) Delete_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSubnetInput{}

	setField(params["id"], input, "SubnetId")

	output, err := d.api.DeleteSubnet(input)
	if err != nil {
		d.logger.Printf("delete subnet error: %s", err)
		return nil, err
	}
	d.logger.Println("delete subnet done")

	return output, nil
}

func (d *AwsDriver) Create_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.RunInstancesInput{}

	setField(params["base"], input, "ImageId")
	setField(params["type"], input, "InstanceType")
	setField(params["count"], input, "MaxCount")
	setField(params["count"], input, "MinCount")
	setField(params["subnet"], input, "SubnetId")

	output, err := d.api.RunInstances(input)
	if err != nil {
		d.logger.Printf("create instance error: %s", err)
		return nil, err
	}
	d.logger.Println("create instance done")

	return aws.StringValue(output.Instances[0].InstanceId), nil
}

func (d *AwsDriver) Delete_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.TerminateInstancesInput{}

	setField(params["id"], input, "InstanceIds")

	output, err := d.api.TerminateInstances(input)
	if err != nil {
		d.logger.Printf("delete instance error: %s", err)
		return nil, err
	}
	d.logger.Println("delete instance done")

	return output, nil
}
