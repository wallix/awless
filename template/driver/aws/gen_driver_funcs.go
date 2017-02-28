/* Copyright 2017 WALLIX

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

// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	dryRunOperation = "DryRunOperation"
	notFound        = "NotFound"
)

// This function was auto generated
func (d *AwsDriver) Create_Vpc_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVpcInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["cidr"], input, "CidrBlock", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.CreateVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["cidr"], input, "CidrBlock", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.CreateVpcOutput
	output, err = d.ec2.CreateVpc(input)
	output = output
	if err != nil {
		d.logger.Errorf("create vpc error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.CreateVpc call took %s", time.Since(start))
	id := aws.StringValue(output.Vpc.VpcId)
	d.logger.Verbosef("create vpc '%s' done", id)
	return aws.StringValue(output.Vpc.VpcId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Vpc_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVpcInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DeleteVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteVpcOutput
	output, err = d.ec2.DeleteVpc(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete vpc error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DeleteVpc call took %s", time.Since(start))
	d.logger.Verbose("delete vpc done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSubnetInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["cidr"], input, "CidrBlock", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["zone"]; ok {
		err = setFieldWithType(params["zone"], input, "AvailabilityZone", awsstr)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.ec2.CreateSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["cidr"], input, "CidrBlock", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["zone"]; ok {
		err = setFieldWithType(params["zone"], input, "AvailabilityZone", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.CreateSubnetOutput
	output, err = d.ec2.CreateSubnet(input)
	output = output
	if err != nil {
		d.logger.Errorf("create subnet error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.CreateSubnet call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["public"]; ok {
		err = setFieldWithType(params["public"], input, "MapPublicIpOnLaunch", awsbool)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.ModifySubnetAttributeOutput
	output, err = d.ec2.ModifySubnetAttribute(input)
	output = output
	if err != nil {
		d.logger.Errorf("update subnet error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.ModifySubnetAttribute call took %s", time.Since(start))
	d.logger.Verbose("update subnet done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSubnetInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DeleteSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteSubnetOutput
	output, err = d.ec2.DeleteSubnet(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete subnet error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DeleteSubnet call took %s", time.Since(start))
	d.logger.Verbose("delete subnet done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.RunInstancesInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["image"], input, "ImageId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["count"], input, "MaxCount", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["count"], input, "MinCount", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["type"], input, "InstanceType", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["subnet"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["key"]; ok {
		err = setFieldWithType(params["key"], input, "KeyName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["ip"]; ok {
		err = setFieldWithType(params["ip"], input, "PrivateIpAddress", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["userdata"]; ok {
		err = setFieldWithType(params["userdata"], input, "UserData", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["group"]; ok {
		err = setFieldWithType(params["group"], input, "SecurityGroupIds", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["lock"]; ok {
		err = setFieldWithType(params["lock"], input, "DisableApiTermination", awsboolattribute)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.ec2.RunInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("instance")
			tagsParams := map[string]interface{}{"resource": id}
			if v, ok := params["name"]; ok {
				tagsParams["Name"] = v
			}
			if len(tagsParams) > 1 {
				_, err = d.Create_Tags_DryRun(tagsParams)
				if err != nil {
					d.logger.Errorf("create instance: adding tags: error: %s", err)
					return nil, err
				}
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
	var err error

	// Required params
	err = setFieldWithType(params["image"], input, "ImageId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["count"], input, "MaxCount", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["count"], input, "MinCount", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["type"], input, "InstanceType", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["subnet"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["key"]; ok {
		err = setFieldWithType(params["key"], input, "KeyName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["ip"]; ok {
		err = setFieldWithType(params["ip"], input, "PrivateIpAddress", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["userdata"]; ok {
		err = setFieldWithType(params["userdata"], input, "UserData", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["group"]; ok {
		err = setFieldWithType(params["group"], input, "SecurityGroupIds", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["lock"]; ok {
		err = setFieldWithType(params["lock"], input, "DisableApiTermination", awsboolattribute)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.Reservation
	output, err = d.ec2.RunInstances(input)
	output = output
	if err != nil {
		d.logger.Errorf("create instance error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.RunInstances call took %s", time.Since(start))
	id := aws.StringValue(output.Instances[0].InstanceId)
	tagsParams := map[string]interface{}{"resource": id}
	if v, ok := params["name"]; ok {
		tagsParams["Name"] = v
	}
	if len(tagsParams) > 1 {
		_, err := d.Create_Tags(tagsParams)
		if err != nil {
			d.logger.Errorf("create instance: adding tags: error: %s", err)
			return nil, err
		}
	}
	d.logger.Verbosef("create instance '%s' done", id)
	return aws.StringValue(output.Instances[0].InstanceId), nil
}

// This function was auto generated
func (d *AwsDriver) Update_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ModifyInstanceAttributeInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["type"]; ok {
		err = setFieldWithType(params["type"], input, "InstanceType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["group"]; ok {
		err = setFieldWithType(params["group"], input, "Groups", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["lock"]; ok {
		err = setFieldWithType(params["lock"], input, "DisableApiTermination", awsboolattribute)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.ec2.ModifyInstanceAttribute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["type"]; ok {
		err = setFieldWithType(params["type"], input, "InstanceType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["group"]; ok {
		err = setFieldWithType(params["group"], input, "Groups", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["lock"]; ok {
		err = setFieldWithType(params["lock"], input, "DisableApiTermination", awsboolattribute)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.ModifyInstanceAttributeOutput
	output, err = d.ec2.ModifyInstanceAttribute(input)
	output = output
	if err != nil {
		d.logger.Errorf("update instance error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.ModifyInstanceAttribute call took %s", time.Since(start))
	d.logger.Verbose("update instance done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.TerminateInstancesInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.TerminateInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.TerminateInstancesOutput
	output, err = d.ec2.TerminateInstances(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete instance error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.TerminateInstances call took %s", time.Since(start))
	d.logger.Verbose("delete instance done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Start_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StartInstancesInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.StartInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.StartInstancesOutput
	output, err = d.ec2.StartInstances(input)
	output = output
	if err != nil {
		d.logger.Errorf("start instance error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.StartInstances call took %s", time.Since(start))
	id := aws.StringValue(output.StartingInstances[0].InstanceId)
	d.logger.Verbosef("start instance '%s' done", id)
	return aws.StringValue(output.StartingInstances[0].InstanceId), nil
}

// This function was auto generated
func (d *AwsDriver) Stop_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StopInstancesInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.StopInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.StopInstancesOutput
	output, err = d.ec2.StopInstances(input)
	output = output
	if err != nil {
		d.logger.Errorf("stop instance error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.StopInstances call took %s", time.Since(start))
	id := aws.StringValue(output.StoppingInstances[0].InstanceId)
	d.logger.Verbosef("stop instance '%s' done", id)
	return aws.StringValue(output.StoppingInstances[0].InstanceId), nil
}

// This function was auto generated
func (d *AwsDriver) Create_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSecurityGroupInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "GroupName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["description"], input, "Description", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.CreateSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "GroupName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["description"], input, "Description", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.CreateSecurityGroupOutput
	output, err = d.ec2.CreateSecurityGroup(input)
	output = output
	if err != nil {
		d.logger.Errorf("create securitygroup error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.CreateSecurityGroup call took %s", time.Since(start))
	id := aws.StringValue(output.GroupId)
	d.logger.Verbosef("create securitygroup '%s' done", id)
	return aws.StringValue(output.GroupId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSecurityGroupInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "GroupId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DeleteSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "GroupId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteSecurityGroupOutput
	output, err = d.ec2.DeleteSecurityGroup(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete securitygroup error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DeleteSecurityGroup call took %s", time.Since(start))
	d.logger.Verbose("delete securitygroup done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVolumeInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["zone"], input, "AvailabilityZone", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["size"], input, "Size", awsint64)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.CreateVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["zone"], input, "AvailabilityZone", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["size"], input, "Size", awsint64)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.Volume
	output, err = d.ec2.CreateVolume(input)
	output = output
	if err != nil {
		d.logger.Errorf("create volume error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.CreateVolume call took %s", time.Since(start))
	id := aws.StringValue(output.VolumeId)
	d.logger.Verbosef("create volume '%s' done", id)
	return aws.StringValue(output.VolumeId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVolumeInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "VolumeId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DeleteVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "VolumeId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteVolumeOutput
	output, err = d.ec2.DeleteVolume(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete volume error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DeleteVolume call took %s", time.Since(start))
	d.logger.Verbose("delete volume done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Attach_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AttachVolumeInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["device"], input, "Device", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["id"], input, "VolumeId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["instance"], input, "InstanceId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.AttachVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["device"], input, "Device", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["id"], input, "VolumeId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["instance"], input, "InstanceId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.VolumeAttachment
	output, err = d.ec2.AttachVolume(input)
	output = output
	if err != nil {
		d.logger.Errorf("attach volume error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.AttachVolume call took %s", time.Since(start))
	id := aws.StringValue(output.VolumeId)
	d.logger.Verbosef("attach volume '%s' done", id)
	return aws.StringValue(output.VolumeId), nil
}

// This function was auto generated
func (d *AwsDriver) Create_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	var err error

	_, err = d.ec2.CreateInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	start := time.Now()
	var output *ec2.CreateInternetGatewayOutput
	output, err = d.ec2.CreateInternetGateway(input)
	output = output
	if err != nil {
		d.logger.Errorf("create internetgateway error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.CreateInternetGateway call took %s", time.Since(start))
	id := aws.StringValue(output.InternetGateway.InternetGatewayId)
	d.logger.Verbosef("create internetgateway '%s' done", id)
	return aws.StringValue(output.InternetGateway.InternetGatewayId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InternetGatewayId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DeleteInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InternetGatewayId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteInternetGatewayOutput
	output, err = d.ec2.DeleteInternetGateway(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete internetgateway error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DeleteInternetGateway call took %s", time.Since(start))
	d.logger.Verbose("delete internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Attach_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AttachInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InternetGatewayId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.AttachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InternetGatewayId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.AttachInternetGatewayOutput
	output, err = d.ec2.AttachInternetGateway(input)
	output = output
	if err != nil {
		d.logger.Errorf("attach internetgateway error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.AttachInternetGateway call took %s", time.Since(start))
	d.logger.Verbose("attach internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Detach_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DetachInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InternetGatewayId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DetachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InternetGatewayId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DetachInternetGatewayOutput
	output, err = d.ec2.DetachInternetGateway(input)
	output = output
	if err != nil {
		d.logger.Errorf("detach internetgateway error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DetachInternetGateway call took %s", time.Since(start))
	d.logger.Verbose("detach internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateRouteTableInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.CreateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.CreateRouteTableOutput
	output, err = d.ec2.CreateRouteTable(input)
	output = output
	if err != nil {
		d.logger.Errorf("create routetable error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.CreateRouteTable call took %s", time.Since(start))
	id := aws.StringValue(output.RouteTable.RouteTableId)
	d.logger.Verbosef("create routetable '%s' done", id)
	return aws.StringValue(output.RouteTable.RouteTableId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteRouteTableInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DeleteRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteRouteTableOutput
	output, err = d.ec2.DeleteRouteTable(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete routetable error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DeleteRouteTable call took %s", time.Since(start))
	d.logger.Verbose("delete routetable done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Attach_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AssociateRouteTableInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["subnet"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.AssociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["subnet"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.AssociateRouteTableOutput
	output, err = d.ec2.AssociateRouteTable(input)
	output = output
	if err != nil {
		d.logger.Errorf("attach routetable error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.AssociateRouteTable call took %s", time.Since(start))
	id := aws.StringValue(output.AssociationId)
	d.logger.Verbosef("attach routetable '%s' done", id)
	return aws.StringValue(output.AssociationId), nil
}

// This function was auto generated
func (d *AwsDriver) Detach_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DisassociateRouteTableInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["association"], input, "AssociationId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DisassociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["association"], input, "AssociationId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DisassociateRouteTableOutput
	output, err = d.ec2.DisassociateRouteTable(input)
	output = output
	if err != nil {
		d.logger.Errorf("detach routetable error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DisassociateRouteTable call took %s", time.Since(start))
	d.logger.Verbose("detach routetable done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Route_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateRouteInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["table"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["cidr"], input, "DestinationCidrBlock", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["gateway"], input, "GatewayId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.CreateRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["table"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["cidr"], input, "DestinationCidrBlock", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["gateway"], input, "GatewayId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.CreateRouteOutput
	output, err = d.ec2.CreateRoute(input)
	output = output
	if err != nil {
		d.logger.Errorf("create route error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.CreateRoute call took %s", time.Since(start))
	d.logger.Verbose("create route done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Route_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteRouteInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["table"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["cidr"], input, "DestinationCidrBlock", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DeleteRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["table"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["cidr"], input, "DestinationCidrBlock", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteRouteOutput
	output, err = d.ec2.DeleteRoute(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete route error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DeleteRoute call took %s", time.Since(start))
	d.logger.Verbose("delete route done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Keypair_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteKeyPairInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "KeyName", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.ec2.DeleteKeyPair(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
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
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "KeyName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteKeyPairOutput
	output, err = d.ec2.DeleteKeyPair(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete keypair error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.DeleteKeyPair call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.CreateUserOutput
	output, err = d.iam.CreateUser(input)
	output = output
	if err != nil {
		d.logger.Errorf("create user error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("iam.CreateUser call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.DeleteUserOutput
	output, err = d.iam.DeleteUser(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete user error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("iam.DeleteUser call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "GroupName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.CreateGroupOutput
	output, err = d.iam.CreateGroup(input)
	output = output
	if err != nil {
		d.logger.Errorf("create group error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("iam.CreateGroup call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "GroupName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.DeleteGroupOutput
	output, err = d.iam.DeleteGroup(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete group error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("iam.DeleteGroup call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "PolicyArn", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["user"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.AttachUserPolicyOutput
	output, err = d.iam.AttachUserPolicy(input)
	output = output
	if err != nil {
		d.logger.Errorf("attach policy error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("iam.AttachUserPolicy call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "PolicyArn", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["user"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.DetachUserPolicyOutput
	output, err = d.iam.DetachUserPolicy(input)
	output = output
	if err != nil {
		d.logger.Errorf("detach policy error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("iam.DetachUserPolicy call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *s3.CreateBucketOutput
	output, err = d.s3.CreateBucket(input)
	output = output
	if err != nil {
		d.logger.Errorf("create bucket error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("s3.CreateBucket call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *s3.DeleteBucketOutput
	output, err = d.s3.DeleteBucket(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete bucket error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("s3.DeleteBucket call took %s", time.Since(start))
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
	var err error

	// Required params
	err = setFieldWithType(params["bucket"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["key"], input, "Key", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *s3.DeleteObjectOutput
	output, err = d.s3.DeleteObject(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete storageobject error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("s3.DeleteObject call took %s", time.Since(start))
	d.logger.Verbose("delete storageobject done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Topic_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create topic: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create topic ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Topic(params map[string]interface{}) (interface{}, error) {
	input := &sns.CreateTopicInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Name", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sns.CreateTopicOutput
	output, err = d.sns.CreateTopic(input)
	output = output
	if err != nil {
		d.logger.Errorf("create topic error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("sns.CreateTopic call took %s", time.Since(start))
	id := aws.StringValue(output.TopicArn)
	d.logger.Verbosef("create topic '%s' done", id)
	return aws.StringValue(output.TopicArn), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Topic_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("delete topic: missing required params 'arn'")
	}

	d.logger.Verbose("params dry run: delete topic ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Topic(params map[string]interface{}) (interface{}, error) {
	input := &sns.DeleteTopicInput{}
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "TopicArn", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sns.DeleteTopicOutput
	output, err = d.sns.DeleteTopic(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete topic error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("sns.DeleteTopic call took %s", time.Since(start))
	d.logger.Verbose("delete topic done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Subscription_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["topic"]; !ok {
		return nil, errors.New("create subscription: missing required params 'topic'")
	}

	if _, ok := params["endpoint"]; !ok {
		return nil, errors.New("create subscription: missing required params 'endpoint'")
	}

	if _, ok := params["protocol"]; !ok {
		return nil, errors.New("create subscription: missing required params 'protocol'")
	}

	d.logger.Verbose("params dry run: create subscription ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Subscription(params map[string]interface{}) (interface{}, error) {
	input := &sns.SubscribeInput{}
	var err error

	// Required params
	err = setFieldWithType(params["topic"], input, "TopicArn", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["endpoint"], input, "Endpoint", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["protocol"], input, "Protocol", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sns.SubscribeOutput
	output, err = d.sns.Subscribe(input)
	output = output
	if err != nil {
		d.logger.Errorf("create subscription error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("sns.Subscribe call took %s", time.Since(start))
	id := aws.StringValue(output.SubscriptionArn)
	d.logger.Verbosef("create subscription '%s' done", id)
	return aws.StringValue(output.SubscriptionArn), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Subscription_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("delete subscription: missing required params 'arn'")
	}

	d.logger.Verbose("params dry run: delete subscription ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Subscription(params map[string]interface{}) (interface{}, error) {
	input := &sns.UnsubscribeInput{}
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "SubscriptionArn", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sns.UnsubscribeOutput
	output, err = d.sns.Unsubscribe(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete subscription error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("sns.Unsubscribe call took %s", time.Since(start))
	d.logger.Verbose("delete subscription done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Queue_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create queue: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create queue ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Queue(params map[string]interface{}) (interface{}, error) {
	input := &sqs.CreateQueueInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "QueueName", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["delay"]; ok {
		err = setFieldWithType(params["delay"], input, "Attributes[DelaySeconds]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["maxMsgSize"]; ok {
		err = setFieldWithType(params["maxMsgSize"], input, "Attributes[MaximumMessageSize]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["retentionPeriod"]; ok {
		err = setFieldWithType(params["retentionPeriod"], input, "Attributes[MessageRetentionPeriod]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["policy"]; ok {
		err = setFieldWithType(params["policy"], input, "Attributes[Policy]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["msgWait"]; ok {
		err = setFieldWithType(params["msgWait"], input, "Attributes[ReceiveMessageWaitTimeSeconds]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["redrivePolicy"]; ok {
		err = setFieldWithType(params["redrivePolicy"], input, "Attributes[RedrivePolicy]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["visibilityTimeout"]; ok {
		err = setFieldWithType(params["visibilityTimeout"], input, "Attributes[VisibilityTimeout]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *sqs.CreateQueueOutput
	output, err = d.sqs.CreateQueue(input)
	output = output
	if err != nil {
		d.logger.Errorf("create queue error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("sqs.CreateQueue call took %s", time.Since(start))
	id := aws.StringValue(output.QueueUrl)
	d.logger.Verbosef("create queue '%s' done", id)
	return aws.StringValue(output.QueueUrl), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Queue_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["url"]; !ok {
		return nil, errors.New("delete queue: missing required params 'url'")
	}

	d.logger.Verbose("params dry run: delete queue ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Queue(params map[string]interface{}) (interface{}, error) {
	input := &sqs.DeleteQueueInput{}
	var err error

	// Required params
	err = setFieldWithType(params["url"], input, "QueueUrl", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sqs.DeleteQueueOutput
	output, err = d.sqs.DeleteQueue(input)
	output = output
	if err != nil {
		d.logger.Errorf("delete queue error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("sqs.DeleteQueue call took %s", time.Since(start))
	d.logger.Verbose("delete queue done")
	return output, nil
}
