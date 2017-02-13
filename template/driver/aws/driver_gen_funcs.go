// DO NOT EDIT
// This file was automatically generated with go generate

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
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
)

// This function was auto generated
func (d *AwsDriver) Create_Vpc_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVpcInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["cidr"], input, "CidrBlock")

	_, err := d.ec2.CreateVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("vpc")
			d.logger.Verbose("full dry run: create vpc ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create vpc error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVpcInput{}

	// Required params
	setField(params["cidr"], input, "CidrBlock")

	output, err := d.ec2.CreateVpc(input)
	if err != nil {
		d.logger.Errorf("create vpc error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.Vpc.VpcId)
	d.logger.Verbosef("create vpc '%s' done", id)
	return aws.StringValue(output.Vpc.VpcId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Vpc_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVpcInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "VpcId")

	_, err := d.ec2.DeleteVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("vpc")
			d.logger.Verbose("full dry run: delete vpc ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete vpc error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVpcInput{}

	// Required params
	setField(params["id"], input, "VpcId")

	output, err := d.ec2.DeleteVpc(input)
	if err != nil {
		d.logger.Errorf("delete vpc error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete vpc done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSubnetInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["cidr"], input, "CidrBlock")
	setField(params["vpc"], input, "VpcId")

	// Extra params
	if _, ok := params["zone"]; ok {
		setField(params["zone"], input, "AvailabilityZone")
	}

	_, err := d.ec2.CreateSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("subnet")
			d.logger.Verbose("full dry run: create subnet ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create subnet error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSubnetInput{}

	// Required params
	setField(params["cidr"], input, "CidrBlock")
	setField(params["vpc"], input, "VpcId")

	// Extra params
	if _, ok := params["zone"]; ok {
		setField(params["zone"], input, "AvailabilityZone")
	}

	output, err := d.ec2.CreateSubnet(input)
	if err != nil {
		d.logger.Errorf("create subnet error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.Subnet.SubnetId)
	d.logger.Verbosef("create subnet '%s' done", id)
	return aws.StringValue(output.Subnet.SubnetId), nil
}

// This function was auto generated
func (d *AwsDriver) Update_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("update subnet: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: update subnet ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Update_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ModifySubnetAttributeInput{}

	// Required params
	setField(params["id"], input, "SubnetId")

	// Extra params
	if _, ok := params["public-vms"]; ok {
		setField(params["public-vms"], input, "MapPublicIpOnLaunch")
	}

	output, err := d.ec2.ModifySubnetAttribute(input)
	if err != nil {
		d.logger.Errorf("update subnet error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("update subnet done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSubnetInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "SubnetId")

	_, err := d.ec2.DeleteSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("subnet")
			d.logger.Verbose("full dry run: delete subnet ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete subnet error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSubnetInput{}

	// Required params
	setField(params["id"], input, "SubnetId")

	output, err := d.ec2.DeleteSubnet(input)
	if err != nil {
		d.logger.Errorf("delete subnet error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete subnet done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.RunInstancesInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["image"], input, "ImageId")
	setField(params["type"], input, "InstanceType")
	setField(params["count"], input, "MaxCount")
	setField(params["count"], input, "MinCount")
	setField(params["subnet"], input, "SubnetId")

	// Extra params
	if _, ok := params["lock"]; ok {
		setField(params["lock"], input, "DisableApiTermination")
	}
	if _, ok := params["key"]; ok {
		setField(params["key"], input, "KeyName")
	}
	if _, ok := params["ip"]; ok {
		setField(params["ip"], input, "PrivateIpAddress")
	}
	if _, ok := params["group"]; ok {
		setField(params["group"], input, "SecurityGroupIds")
	}
	if _, ok := params["userdata"]; ok {
		setField(params["userdata"], input, "UserData")
	}

	_, err := d.ec2.RunInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("instance")
			tagsParams := map[string]interface{}{"resource": id}
			if v, ok := params["name"]; ok {
				tagsParams["Name"] = v
			}
			if len(tagsParams) > 1 {
				d.Create_Tags_DryRun(tagsParams)
			}
			d.logger.Verbose("full dry run: create instance ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.RunInstancesInput{}

	// Required params
	setField(params["image"], input, "ImageId")
	setField(params["type"], input, "InstanceType")
	setField(params["count"], input, "MaxCount")
	setField(params["count"], input, "MinCount")
	setField(params["subnet"], input, "SubnetId")

	// Extra params
	if _, ok := params["lock"]; ok {
		setField(params["lock"], input, "DisableApiTermination")
	}
	if _, ok := params["key"]; ok {
		setField(params["key"], input, "KeyName")
	}
	if _, ok := params["ip"]; ok {
		setField(params["ip"], input, "PrivateIpAddress")
	}
	if _, ok := params["group"]; ok {
		setField(params["group"], input, "SecurityGroupIds")
	}
	if _, ok := params["userdata"]; ok {
		setField(params["userdata"], input, "UserData")
	}

	output, err := d.ec2.RunInstances(input)
	if err != nil {
		d.logger.Errorf("create instance error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.Instances[0].InstanceId)
	tagsParams := map[string]interface{}{"resource": id}
	if v, ok := params["name"]; ok {
		tagsParams["Name"] = v
	}
	if len(tagsParams) > 1 {
		d.Create_Tags(tagsParams)
	}
	d.logger.Verbosef("create instance '%s' done", id)
	return aws.StringValue(output.Instances[0].InstanceId), nil
}

// This function was auto generated
func (d *AwsDriver) Update_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ModifyInstanceAttributeInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "InstanceId")

	// Extra params
	if _, ok := params["lock"]; ok {
		setField(params["lock"], input, "DisableApiTermination")
	}
	if _, ok := params["group"]; ok {
		setField(params["group"], input, "Groups")
	}
	if _, ok := params["type"]; ok {
		setField(params["type"], input, "InstanceType")
	}

	_, err := d.ec2.ModifyInstanceAttribute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("instance")
			d.logger.Verbose("full dry run: update instance ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: update instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Update_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ModifyInstanceAttributeInput{}

	// Required params
	setField(params["id"], input, "InstanceId")

	// Extra params
	if _, ok := params["lock"]; ok {
		setField(params["lock"], input, "DisableApiTermination")
	}
	if _, ok := params["group"]; ok {
		setField(params["group"], input, "Groups")
	}
	if _, ok := params["type"]; ok {
		setField(params["type"], input, "InstanceType")
	}

	output, err := d.ec2.ModifyInstanceAttribute(input)
	if err != nil {
		d.logger.Errorf("update instance error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("update instance done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.TerminateInstancesInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "InstanceIds")

	_, err := d.ec2.TerminateInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("instance")
			d.logger.Verbose("full dry run: delete instance ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.TerminateInstancesInput{}

	// Required params
	setField(params["id"], input, "InstanceIds")

	output, err := d.ec2.TerminateInstances(input)
	if err != nil {
		d.logger.Errorf("delete instance error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete instance done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Start_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StartInstancesInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "InstanceIds")

	_, err := d.ec2.StartInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("instance")
			d.logger.Verbose("full dry run: start instance ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: start instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Start_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StartInstancesInput{}

	// Required params
	setField(params["id"], input, "InstanceIds")

	output, err := d.ec2.StartInstances(input)
	if err != nil {
		d.logger.Errorf("start instance error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.StartingInstances[0].InstanceId)
	d.logger.Verbosef("start instance '%s' done", id)
	return aws.StringValue(output.StartingInstances[0].InstanceId), nil
}

// This function was auto generated
func (d *AwsDriver) Stop_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StopInstancesInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "InstanceIds")

	_, err := d.ec2.StopInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("instance")
			d.logger.Verbose("full dry run: stop instance ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: stop instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Stop_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StopInstancesInput{}

	// Required params
	setField(params["id"], input, "InstanceIds")

	output, err := d.ec2.StopInstances(input)
	if err != nil {
		d.logger.Errorf("stop instance error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.StoppingInstances[0].InstanceId)
	d.logger.Verbosef("stop instance '%s' done", id)
	return aws.StringValue(output.StoppingInstances[0].InstanceId), nil
}

// This function was auto generated
func (d *AwsDriver) Create_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSecurityGroupInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["description"], input, "Description")
	setField(params["name"], input, "GroupName")
	setField(params["vpc"], input, "VpcId")

	_, err := d.ec2.CreateSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("securitygroup")
			d.logger.Verbose("full dry run: create securitygroup ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create securitygroup error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Securitygroup(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSecurityGroupInput{}

	// Required params
	setField(params["description"], input, "Description")
	setField(params["name"], input, "GroupName")
	setField(params["vpc"], input, "VpcId")

	output, err := d.ec2.CreateSecurityGroup(input)
	if err != nil {
		d.logger.Errorf("create securitygroup error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.GroupId)
	d.logger.Verbosef("create securitygroup '%s' done", id)
	return aws.StringValue(output.GroupId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSecurityGroupInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "GroupId")

	_, err := d.ec2.DeleteSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("securitygroup")
			d.logger.Verbose("full dry run: delete securitygroup ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete securitygroup error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Securitygroup(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSecurityGroupInput{}

	// Required params
	setField(params["id"], input, "GroupId")

	output, err := d.ec2.DeleteSecurityGroup(input)
	if err != nil {
		d.logger.Errorf("delete securitygroup error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete securitygroup done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVolumeInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["zone"], input, "AvailabilityZone")
	setField(params["size"], input, "Size")

	_, err := d.ec2.CreateVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("volume")
			d.logger.Verbose("full dry run: create volume ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create volume error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Volume(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVolumeInput{}

	// Required params
	setField(params["zone"], input, "AvailabilityZone")
	setField(params["size"], input, "Size")

	output, err := d.ec2.CreateVolume(input)
	if err != nil {
		d.logger.Errorf("create volume error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.VolumeId)
	d.logger.Verbosef("create volume '%s' done", id)
	return aws.StringValue(output.VolumeId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVolumeInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "VolumeId")

	_, err := d.ec2.DeleteVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("volume")
			d.logger.Verbose("full dry run: delete volume ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete volume error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Volume(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVolumeInput{}

	// Required params
	setField(params["id"], input, "VolumeId")

	output, err := d.ec2.DeleteVolume(input)
	if err != nil {
		d.logger.Errorf("delete volume error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete volume done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Attach_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AttachVolumeInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["device"], input, "Device")
	setField(params["instance"], input, "InstanceId")
	setField(params["id"], input, "VolumeId")

	_, err := d.ec2.AttachVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("volume")
			d.logger.Verbose("full dry run: attach volume ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: attach volume error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Attach_Volume(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AttachVolumeInput{}

	// Required params
	setField(params["device"], input, "Device")
	setField(params["instance"], input, "InstanceId")
	setField(params["id"], input, "VolumeId")

	output, err := d.ec2.AttachVolume(input)
	if err != nil {
		d.logger.Errorf("attach volume error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.VolumeId)
	d.logger.Verbosef("attach volume '%s' done", id)
	return aws.StringValue(output.VolumeId), nil
}

// This function was auto generated
func (d *AwsDriver) Create_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateInternetGatewayInput{}
	input.DryRun = aws.Bool(true)

	_, err := d.ec2.CreateInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("internetgateway")
			d.logger.Verbose("full dry run: create internetgateway ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create internetgateway error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Internetgateway(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateInternetGatewayInput{}

	output, err := d.ec2.CreateInternetGateway(input)
	if err != nil {
		d.logger.Errorf("create internetgateway error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.InternetGateway.InternetGatewayId)
	d.logger.Verbosef("create internetgateway '%s' done", id)
	return aws.StringValue(output.InternetGateway.InternetGatewayId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteInternetGatewayInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "InternetGatewayId")

	_, err := d.ec2.DeleteInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("internetgateway")
			d.logger.Verbose("full dry run: delete internetgateway ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete internetgateway error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Internetgateway(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteInternetGatewayInput{}

	// Required params
	setField(params["id"], input, "InternetGatewayId")

	output, err := d.ec2.DeleteInternetGateway(input)
	if err != nil {
		d.logger.Errorf("delete internetgateway error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Attach_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AttachInternetGatewayInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "InternetGatewayId")
	setField(params["vpc"], input, "VpcId")

	_, err := d.ec2.AttachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("internetgateway")
			d.logger.Verbose("full dry run: attach internetgateway ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: attach internetgateway error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Attach_Internetgateway(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AttachInternetGatewayInput{}

	// Required params
	setField(params["id"], input, "InternetGatewayId")
	setField(params["vpc"], input, "VpcId")

	output, err := d.ec2.AttachInternetGateway(input)
	if err != nil {
		d.logger.Errorf("attach internetgateway error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("attach internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Detach_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DetachInternetGatewayInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "InternetGatewayId")
	setField(params["vpc"], input, "VpcId")

	_, err := d.ec2.DetachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("internetgateway")
			d.logger.Verbose("full dry run: detach internetgateway ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: detach internetgateway error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Detach_Internetgateway(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DetachInternetGatewayInput{}

	// Required params
	setField(params["id"], input, "InternetGatewayId")
	setField(params["vpc"], input, "VpcId")

	output, err := d.ec2.DetachInternetGateway(input)
	if err != nil {
		d.logger.Errorf("detach internetgateway error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("detach internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateRouteTableInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["vpc"], input, "VpcId")

	_, err := d.ec2.CreateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("routetable")
			d.logger.Verbose("full dry run: create routetable ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create routetable error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Routetable(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateRouteTableInput{}

	// Required params
	setField(params["vpc"], input, "VpcId")

	output, err := d.ec2.CreateRouteTable(input)
	if err != nil {
		d.logger.Errorf("create routetable error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.RouteTable.RouteTableId)
	d.logger.Verbosef("create routetable '%s' done", id)
	return aws.StringValue(output.RouteTable.RouteTableId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteRouteTableInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "RouteTableId")

	_, err := d.ec2.DeleteRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("routetable")
			d.logger.Verbose("full dry run: delete routetable ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete routetable error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Routetable(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteRouteTableInput{}

	// Required params
	setField(params["id"], input, "RouteTableId")

	output, err := d.ec2.DeleteRouteTable(input)
	if err != nil {
		d.logger.Errorf("delete routetable error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete routetable done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Attach_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AssociateRouteTableInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "RouteTableId")
	setField(params["subnet"], input, "SubnetId")

	_, err := d.ec2.AssociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("routetable")
			d.logger.Verbose("full dry run: attach routetable ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: attach routetable error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Attach_Routetable(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AssociateRouteTableInput{}

	// Required params
	setField(params["id"], input, "RouteTableId")
	setField(params["subnet"], input, "SubnetId")

	output, err := d.ec2.AssociateRouteTable(input)
	if err != nil {
		d.logger.Errorf("attach routetable error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.AssociationId)
	d.logger.Verbosef("attach routetable '%s' done", id)
	return aws.StringValue(output.AssociationId), nil
}

// This function was auto generated
func (d *AwsDriver) Detach_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DisassociateRouteTableInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["association"], input, "AssociationId")

	_, err := d.ec2.DisassociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("routetable")
			d.logger.Verbose("full dry run: detach routetable ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: detach routetable error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Detach_Routetable(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DisassociateRouteTableInput{}

	// Required params
	setField(params["association"], input, "AssociationId")

	output, err := d.ec2.DisassociateRouteTable(input)
	if err != nil {
		d.logger.Errorf("detach routetable error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("detach routetable done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Route_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateRouteInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["cidr"], input, "DestinationCidrBlock")
	setField(params["gateway"], input, "GatewayId")
	setField(params["table"], input, "RouteTableId")

	_, err := d.ec2.CreateRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("route")
			d.logger.Verbose("full dry run: create route ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create route error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Route(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateRouteInput{}

	// Required params
	setField(params["cidr"], input, "DestinationCidrBlock")
	setField(params["gateway"], input, "GatewayId")
	setField(params["table"], input, "RouteTableId")

	output, err := d.ec2.CreateRoute(input)
	if err != nil {
		d.logger.Errorf("create route error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("create route done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Route_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteRouteInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["cidr"], input, "DestinationCidrBlock")
	setField(params["table"], input, "RouteTableId")

	_, err := d.ec2.DeleteRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("route")
			d.logger.Verbose("full dry run: delete route ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete route error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Route(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteRouteInput{}

	// Required params
	setField(params["cidr"], input, "DestinationCidrBlock")
	setField(params["table"], input, "RouteTableId")

	output, err := d.ec2.DeleteRoute(input)
	if err != nil {
		d.logger.Errorf("delete route error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete route done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Keypair_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteKeyPairInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	setField(params["id"], input, "KeyName")

	_, err := d.ec2.DeleteKeyPair(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("keypair")
			d.logger.Verbose("full dry run: delete keypair ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: delete keypair error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Keypair(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteKeyPairInput{}

	// Required params
	setField(params["id"], input, "KeyName")

	output, err := d.ec2.DeleteKeyPair(input)
	if err != nil {
		d.logger.Errorf("delete keypair error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete keypair done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create user: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create user ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Create_User(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateUserInput{}

	// Required params
	setField(params["name"], input, "UserName")

	output, err := d.iam.CreateUser(input)
	if err != nil {
		d.logger.Errorf("create user error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.User.UserId)
	d.logger.Verbosef("create user '%s' done", id)
	return aws.StringValue(output.User.UserId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete user: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete user ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_User(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteUserInput{}

	// Required params
	setField(params["name"], input, "UserName")

	output, err := d.iam.DeleteUser(input)
	if err != nil {
		d.logger.Errorf("delete user error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete user done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Group_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create group: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create group ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Group(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateGroupInput{}

	// Required params
	setField(params["name"], input, "GroupName")

	output, err := d.iam.CreateGroup(input)
	if err != nil {
		d.logger.Errorf("create group error: %s", err)
		return nil, err
	}
	output = output
	id := aws.StringValue(output.Group.GroupId)
	d.logger.Verbosef("create group '%s' done", id)
	return aws.StringValue(output.Group.GroupId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Group_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete group: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete group ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Group(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteGroupInput{}

	// Required params
	setField(params["name"], input, "GroupName")

	output, err := d.iam.DeleteGroup(input)
	if err != nil {
		d.logger.Errorf("delete group error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete group done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Attach_Policy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("attach policy: missing required params 'arn'")
	}

	if _, ok := params["user"]; !ok {
		return nil, errors.New("attach policy: missing required params 'user'")
	}

	d.logger.Verbose("params dry run: attach policy ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Attach_Policy(params map[string]interface{}) (interface{}, error) {
	input := &iam.AttachUserPolicyInput{}

	// Required params
	setField(params["arn"], input, "PolicyArn")
	setField(params["user"], input, "UserName")

	output, err := d.iam.AttachUserPolicy(input)
	if err != nil {
		d.logger.Errorf("attach policy error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("attach policy done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Detach_Policy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("detach policy: missing required params 'arn'")
	}

	if _, ok := params["user"]; !ok {
		return nil, errors.New("detach policy: missing required params 'user'")
	}

	d.logger.Verbose("params dry run: detach policy ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Detach_Policy(params map[string]interface{}) (interface{}, error) {
	input := &iam.DetachUserPolicyInput{}

	// Required params
	setField(params["arn"], input, "PolicyArn")
	setField(params["user"], input, "UserName")

	output, err := d.iam.DetachUserPolicy(input)
	if err != nil {
		d.logger.Errorf("detach policy error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("detach policy done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Bucket_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create bucket: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create bucket ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Bucket(params map[string]interface{}) (interface{}, error) {
	input := &s3.CreateBucketInput{}

	// Required params
	setField(params["name"], input, "Bucket")

	output, err := d.s3.CreateBucket(input)
	if err != nil {
		d.logger.Errorf("create bucket error: %s", err)
		return nil, err
	}
	output = output
	id := params["name"]
	d.logger.Verbosef("create bucket '%s' done", id)
	return params["name"], nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Bucket_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete bucket: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete bucket ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Bucket(params map[string]interface{}) (interface{}, error) {
	input := &s3.DeleteBucketInput{}

	// Required params
	setField(params["name"], input, "Bucket")

	output, err := d.s3.DeleteBucket(input)
	if err != nil {
		d.logger.Errorf("delete bucket error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete bucket done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Storageobject_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["bucket"]; !ok {
		return nil, errors.New("delete storageobject: missing required params 'bucket'")
	}

	if _, ok := params["key"]; !ok {
		return nil, errors.New("delete storageobject: missing required params 'key'")
	}

	d.logger.Verbose("params dry run: delete storageobject ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Storageobject(params map[string]interface{}) (interface{}, error) {
	input := &s3.DeleteObjectInput{}

	// Required params
	setField(params["bucket"], input, "Bucket")
	setField(params["key"], input, "Key")

	output, err := d.s3.DeleteObject(input)
	if err != nil {
		d.logger.Errorf("delete storageobject error: %s", err)
		return nil, err
	}
	output = output
	d.logger.Verbose("delete storageobject done")
	return output, nil
}
