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
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	dryRunOperation = "DryRunOperation"
	notFound        = "NotFound"
)

// This function was auto generated
func (d *Ec2Driver) Create_Vpc_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVpcInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["cidr"], input, "CidrBlock", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params

	_, err = d.CreateVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("vpc")
			// Extra param as tag
			if v, ok := params["name"]; ok {
				_, err = d.Create_Tag_DryRun(map[string]interface{}{"key": "Name", "value": v, "resource": id})
				if err != nil {
					return nil, fmt.Errorf("dry run: create vpc: adding tags: %s", err)
				}
			}
			d.logger.Verbose("dry run: create vpc ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create vpc: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVpcInput{}
	var err error

	// Required params
	err = setFieldWithType(params["cidr"], input, "CidrBlock", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params

	start := time.Now()
	var output *ec2.CreateVpcOutput
	output, err = d.CreateVpc(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create vpc: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateVpc call took %s", time.Since(start))
	id := aws.StringValue(output.Vpc.VpcId)
	// Extra param as tag
	if v, ok := params["name"]; ok {
		_, err = d.Create_Tag(map[string]interface{}{"key": "Name", "value": v, "resource": id})
		if err != nil {
			return nil, fmt.Errorf("create vpc: adding tags: %s", err)
		}
	}

	d.logger.Verbosef("create vpc '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Vpc_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVpcInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("vpc")
			d.logger.Verbose("dry run: delete vpc ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete vpc: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVpcInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteVpcOutput
	output, err = d.DeleteVpc(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete vpc: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteVpc call took %s", time.Since(start))
	d.logger.Verbose("delete vpc done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.CreateSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("subnet")
			// Extra param as tag
			if v, ok := params["name"]; ok {
				_, err = d.Create_Tag_DryRun(map[string]interface{}{"key": "Name", "value": v, "resource": id})
				if err != nil {
					return nil, fmt.Errorf("dry run: create subnet: adding tags: %s", err)
				}
			}
			d.logger.Verbose("dry run: create subnet ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create subnet: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Subnet(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.CreateSubnet(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create subnet: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateSubnet call took %s", time.Since(start))
	id := aws.StringValue(output.Subnet.SubnetId)
	// Extra param as tag
	if v, ok := params["name"]; ok {
		_, err = d.Create_Tag(map[string]interface{}{"key": "Name", "value": v, "resource": id})
		if err != nil {
			return nil, fmt.Errorf("create subnet: adding tags: %s", err)
		}
	}

	d.logger.Verbosef("create subnet '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Update_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("update subnet: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: update subnet ok")
	return nil, nil
}

// This function was auto generated
func (d *Ec2Driver) Update_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ModifySubnetAttributeInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["public"]; ok {
		err = setFieldWithType(params["public"], input, "MapPublicIpOnLaunch", awsboolattribute)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.ModifySubnetAttributeOutput
	output, err = d.ModifySubnetAttribute(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("update subnet: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.ModifySubnetAttribute call took %s", time.Since(start))
	d.logger.Verbose("update subnet done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSubnetInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("subnet")
			d.logger.Verbose("dry run: delete subnet ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete subnet: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSubnetInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteSubnetOutput
	output, err = d.DeleteSubnet(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete subnet: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteSubnet call took %s", time.Since(start))
	d.logger.Verbose("delete subnet done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.RunInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("instance")
			// Required param as tag
			_, err = d.Create_Tag_DryRun(map[string]interface{}{"key": "Name", "value": params["name"], "resource": id})
			if err != nil {
				return nil, fmt.Errorf("dry run: create instance: adding tags: %s", err)
			}
			d.logger.Verbose("dry run: create instance ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create instance: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Instance(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.RunInstances(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create instance: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.RunInstances call took %s", time.Since(start))
	id := aws.StringValue(output.Instances[0].InstanceId)
	// Required param as tag
	_, err = d.Create_Tag(map[string]interface{}{"key": "Name", "value": params["name"], "resource": id})
	if err != nil {
		return nil, fmt.Errorf("create instance: adding tags: %s", err)
	}

	d.logger.Verbosef("create instance '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Update_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.ModifyInstanceAttribute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("instance")
			d.logger.Verbose("dry run: update instance ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: update instance: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Update_Instance(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.ModifyInstanceAttribute(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("update instance: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.ModifyInstanceAttribute call took %s", time.Since(start))
	d.logger.Verbose("update instance done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.TerminateInstancesInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	_, err = d.TerminateInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("instance")
			d.logger.Verbose("dry run: delete instance ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete instance: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.TerminateInstancesInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.TerminateInstancesOutput
	output, err = d.TerminateInstances(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete instance: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.TerminateInstances call took %s", time.Since(start))
	d.logger.Verbose("delete instance done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Start_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StartInstancesInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	_, err = d.StartInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("instance")
			d.logger.Verbose("dry run: start instance ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: start instance: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Start_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StartInstancesInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.StartInstancesOutput
	output, err = d.StartInstances(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("start instance: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.StartInstances call took %s", time.Since(start))
	id := aws.StringValue(output.StartingInstances[0].InstanceId)

	d.logger.Verbosef("start instance '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Stop_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StopInstancesInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	_, err = d.StopInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("instance")
			d.logger.Verbose("dry run: stop instance ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: stop instance: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Stop_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StopInstancesInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.StopInstancesOutput
	output, err = d.StopInstances(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("stop instance: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.StopInstances call took %s", time.Since(start))
	id := aws.StringValue(output.StoppingInstances[0].InstanceId)

	d.logger.Verbosef("stop instance '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.CreateSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("securitygroup")
			d.logger.Verbose("dry run: create securitygroup ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create securitygroup: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Securitygroup(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.CreateSecurityGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create securitygroup: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateSecurityGroup call took %s", time.Since(start))
	id := aws.StringValue(output.GroupId)

	d.logger.Verbosef("create securitygroup '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSecurityGroupInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "GroupId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("securitygroup")
			d.logger.Verbose("dry run: delete securitygroup ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete securitygroup: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Securitygroup(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSecurityGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "GroupId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteSecurityGroupOutput
	output, err = d.DeleteSecurityGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete securitygroup: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteSecurityGroup call took %s", time.Since(start))
	d.logger.Verbose("delete securitygroup done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.CreateVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("volume")
			d.logger.Verbose("dry run: create volume ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create volume: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Volume(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.CreateVolume(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create volume: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateVolume call took %s", time.Since(start))
	id := aws.StringValue(output.VolumeId)

	d.logger.Verbosef("create volume '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVolumeInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "VolumeId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("volume")
			d.logger.Verbose("dry run: delete volume ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete volume: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Volume(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVolumeInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "VolumeId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteVolumeOutput
	output, err = d.DeleteVolume(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete volume: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteVolume call took %s", time.Since(start))
	d.logger.Verbose("delete volume done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Attach_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.AttachVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("volume")
			d.logger.Verbose("dry run: attach volume ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: attach volume: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Attach_Volume(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.AttachVolume(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("attach volume: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.AttachVolume call took %s", time.Since(start))
	id := aws.StringValue(output.VolumeId)

	d.logger.Verbosef("attach volume '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Detach_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DetachVolumeInput{}
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

	// Extra params
	if _, ok := params["force"]; ok {
		err = setFieldWithType(params["force"], input, "Force", awsbool)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.DetachVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("volume")
			d.logger.Verbose("dry run: detach volume ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: detach volume: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Detach_Volume(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DetachVolumeInput{}
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

	// Extra params
	if _, ok := params["force"]; ok {
		err = setFieldWithType(params["force"], input, "Force", awsbool)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.VolumeAttachment
	output, err = d.DetachVolume(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("detach volume: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DetachVolume call took %s", time.Since(start))
	id := aws.StringValue(output.VolumeId)

	d.logger.Verbosef("detach volume '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	var err error

	_, err = d.CreateInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("internetgateway")
			d.logger.Verbose("dry run: create internetgateway ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create internetgateway: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Internetgateway(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateInternetGatewayInput{}
	var err error

	start := time.Now()
	var output *ec2.CreateInternetGatewayOutput
	output, err = d.CreateInternetGateway(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create internetgateway: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateInternetGateway call took %s", time.Since(start))
	id := aws.StringValue(output.InternetGateway.InternetGatewayId)

	d.logger.Verbosef("create internetgateway '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InternetGatewayId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("internetgateway")
			d.logger.Verbose("dry run: delete internetgateway ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete internetgateway: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Internetgateway(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteInternetGatewayInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "InternetGatewayId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteInternetGatewayOutput
	output, err = d.DeleteInternetGateway(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete internetgateway: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteInternetGateway call took %s", time.Since(start))
	d.logger.Verbose("delete internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Attach_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.AttachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("internetgateway")
			d.logger.Verbose("dry run: attach internetgateway ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: attach internetgateway: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Attach_Internetgateway(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.AttachInternetGateway(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("attach internetgateway: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.AttachInternetGateway call took %s", time.Since(start))
	d.logger.Verbose("attach internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Detach_Internetgateway_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.DetachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("internetgateway")
			d.logger.Verbose("dry run: detach internetgateway ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: detach internetgateway: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Detach_Internetgateway(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.DetachInternetGateway(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("detach internetgateway: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DetachInternetGateway call took %s", time.Since(start))
	d.logger.Verbose("detach internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateRouteTableInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.CreateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("routetable")
			d.logger.Verbose("dry run: create routetable ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create routetable: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Routetable(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateRouteTableInput{}
	var err error

	// Required params
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.CreateRouteTableOutput
	output, err = d.CreateRouteTable(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create routetable: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateRouteTable call took %s", time.Since(start))
	id := aws.StringValue(output.RouteTable.RouteTableId)

	d.logger.Verbosef("create routetable '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteRouteTableInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("routetable")
			d.logger.Verbose("dry run: delete routetable ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete routetable: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Routetable(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteRouteTableInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "RouteTableId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteRouteTableOutput
	output, err = d.DeleteRouteTable(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete routetable: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteRouteTable call took %s", time.Since(start))
	d.logger.Verbose("delete routetable done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Attach_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.AssociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("routetable")
			d.logger.Verbose("dry run: attach routetable ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: attach routetable: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Attach_Routetable(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.AssociateRouteTable(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("attach routetable: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.AssociateRouteTable call took %s", time.Since(start))
	id := aws.StringValue(output.AssociationId)

	d.logger.Verbosef("attach routetable '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Detach_Routetable_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DisassociateRouteTableInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["association"], input, "AssociationId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DisassociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("routetable")
			d.logger.Verbose("dry run: detach routetable ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: detach routetable: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Detach_Routetable(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DisassociateRouteTableInput{}
	var err error

	// Required params
	err = setFieldWithType(params["association"], input, "AssociationId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DisassociateRouteTableOutput
	output, err = d.DisassociateRouteTable(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("detach routetable: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DisassociateRouteTable call took %s", time.Since(start))
	d.logger.Verbose("detach routetable done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Route_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.CreateRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("route")
			d.logger.Verbose("dry run: create route ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create route: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Route(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.CreateRoute(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create route: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateRoute call took %s", time.Since(start))
	d.logger.Verbose("create route done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Route_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.DeleteRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("route")
			d.logger.Verbose("dry run: delete route ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete route: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Route(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.DeleteRoute(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete route: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteRoute call took %s", time.Since(start))
	d.logger.Verbose("delete route done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Keypair_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteKeyPairInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "KeyName", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteKeyPair(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("keypair")
			d.logger.Verbose("dry run: delete keypair ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete keypair: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Keypair(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteKeyPairInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "KeyName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteKeyPairOutput
	output, err = d.DeleteKeyPair(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete keypair: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteKeyPair call took %s", time.Since(start))
	d.logger.Verbose("delete keypair done")
	return output, nil
}

// This function was auto generated
func (d *Elbv2Driver) Create_Loadbalancer_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create loadbalancer: missing required params 'name'")
	}

	if _, ok := params["subnets"]; !ok {
		return nil, errors.New("create loadbalancer: missing required params 'subnets'")
	}

	d.logger.Verbose("params dry run: create loadbalancer ok")
	return nil, nil
}

// This function was auto generated
func (d *Elbv2Driver) Create_Loadbalancer(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.CreateLoadBalancerInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Name", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["subnets"], input, "Subnets", awsstringslice)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["iptype"]; ok {
		err = setFieldWithType(params["iptype"], input, "IpAddressType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["scheme"]; ok {
		err = setFieldWithType(params["scheme"], input, "Scheme", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["groups"]; ok {
		err = setFieldWithType(params["groups"], input, "SecurityGroups", awsstringslice)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *elbv2.CreateLoadBalancerOutput
	output, err = d.CreateLoadBalancer(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create loadbalancer: %s", err)
	}
	d.logger.ExtraVerbosef("elbv2.CreateLoadBalancer call took %s", time.Since(start))
	id := aws.StringValue(output.LoadBalancers[0].LoadBalancerArn)

	d.logger.Verbosef("create loadbalancer '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Loadbalancer_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("delete loadbalancer: missing required params 'arn'")
	}

	d.logger.Verbose("params dry run: delete loadbalancer ok")
	return nil, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Loadbalancer(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.DeleteLoadBalancerInput{}
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "LoadBalancerArn", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *elbv2.DeleteLoadBalancerOutput
	output, err = d.DeleteLoadBalancer(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete loadbalancer: %s", err)
	}
	d.logger.ExtraVerbosef("elbv2.DeleteLoadBalancer call took %s", time.Since(start))
	d.logger.Verbose("delete loadbalancer done")
	return output, nil
}

// This function was auto generated
func (d *Elbv2Driver) Create_Listener_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["actiontype"]; !ok {
		return nil, errors.New("create listener: missing required params 'actiontype'")
	}

	if _, ok := params["target"]; !ok {
		return nil, errors.New("create listener: missing required params 'target'")
	}

	if _, ok := params["loadbalancer"]; !ok {
		return nil, errors.New("create listener: missing required params 'loadbalancer'")
	}

	if _, ok := params["port"]; !ok {
		return nil, errors.New("create listener: missing required params 'port'")
	}

	if _, ok := params["protocol"]; !ok {
		return nil, errors.New("create listener: missing required params 'protocol'")
	}

	d.logger.Verbose("params dry run: create listener ok")
	return nil, nil
}

// This function was auto generated
func (d *Elbv2Driver) Create_Listener(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.CreateListenerInput{}
	var err error

	// Required params
	err = setFieldWithType(params["actiontype"], input, "DefaultActions[0]Type", awsslicestruct)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["target"], input, "DefaultActions[0]TargetGroupArn", awsslicestruct)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["loadbalancer"], input, "LoadBalancerArn", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["port"], input, "Port", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["protocol"], input, "Protocol", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["certificate"]; ok {
		err = setFieldWithType(params["certificate"], input, "Certificates[0]CertificateArn", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["sslpolicy"]; ok {
		err = setFieldWithType(params["sslpolicy"], input, "SslPolicy", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *elbv2.CreateListenerOutput
	output, err = d.CreateListener(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create listener: %s", err)
	}
	d.logger.ExtraVerbosef("elbv2.CreateListener call took %s", time.Since(start))
	id := aws.StringValue(output.Listeners[0].ListenerArn)

	d.logger.Verbosef("create listener '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Listener_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("delete listener: missing required params 'arn'")
	}

	d.logger.Verbose("params dry run: delete listener ok")
	return nil, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Listener(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.DeleteListenerInput{}
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "ListenerArn", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *elbv2.DeleteListenerOutput
	output, err = d.DeleteListener(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete listener: %s", err)
	}
	d.logger.ExtraVerbosef("elbv2.DeleteListener call took %s", time.Since(start))
	d.logger.Verbose("delete listener done")
	return output, nil
}

// This function was auto generated
func (d *Elbv2Driver) Create_Targetgroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create targetgroup: missing required params 'name'")
	}

	if _, ok := params["port"]; !ok {
		return nil, errors.New("create targetgroup: missing required params 'port'")
	}

	if _, ok := params["protocol"]; !ok {
		return nil, errors.New("create targetgroup: missing required params 'protocol'")
	}

	if _, ok := params["vpc"]; !ok {
		return nil, errors.New("create targetgroup: missing required params 'vpc'")
	}

	d.logger.Verbose("params dry run: create targetgroup ok")
	return nil, nil
}

// This function was auto generated
func (d *Elbv2Driver) Create_Targetgroup(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.CreateTargetGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Name", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["port"], input, "Port", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["protocol"], input, "Protocol", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["vpc"], input, "VpcId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["healthcheckinterval"]; ok {
		err = setFieldWithType(params["healthcheckinterval"], input, "HealthCheckIntervalSeconds", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthcheckpath"]; ok {
		err = setFieldWithType(params["healthcheckpath"], input, "HealthCheckPath", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthcheckport"]; ok {
		err = setFieldWithType(params["healthcheckport"], input, "HealthCheckPort", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthcheckprotocol"]; ok {
		err = setFieldWithType(params["healthcheckprotocol"], input, "HealthCheckProtocol", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthchecktimeout"]; ok {
		err = setFieldWithType(params["healthchecktimeout"], input, "HealthCheckTimeoutSeconds", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthythreshold"]; ok {
		err = setFieldWithType(params["healthythreshold"], input, "HealthyThresholdCount", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["unhealthythreshold"]; ok {
		err = setFieldWithType(params["unhealthythreshold"], input, "UnhealthyThresholdCount", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["matcher"]; ok {
		err = setFieldWithType(params["matcher"], input, "Matcher.HttpCode", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *elbv2.CreateTargetGroupOutput
	output, err = d.CreateTargetGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create targetgroup: %s", err)
	}
	d.logger.ExtraVerbosef("elbv2.CreateTargetGroup call took %s", time.Since(start))
	id := aws.StringValue(output.TargetGroups[0].TargetGroupArn)

	d.logger.Verbosef("create targetgroup '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Targetgroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("delete targetgroup: missing required params 'arn'")
	}

	d.logger.Verbose("params dry run: delete targetgroup ok")
	return nil, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Targetgroup(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.DeleteTargetGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "TargetGroupArn", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *elbv2.DeleteTargetGroupOutput
	output, err = d.DeleteTargetGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete targetgroup: %s", err)
	}
	d.logger.ExtraVerbosef("elbv2.DeleteTargetGroup call took %s", time.Since(start))
	d.logger.Verbose("delete targetgroup done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Create_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create user: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create user ok")
	return nil, nil
}

// This function was auto generated
func (d *IamDriver) Create_User(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateUserInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.CreateUserOutput
	output, err = d.CreateUser(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create user: %s", err)
	}
	d.logger.ExtraVerbosef("iam.CreateUser call took %s", time.Since(start))
	id := aws.StringValue(output.User.UserId)

	d.logger.Verbosef("create user '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *IamDriver) Delete_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete user: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete user ok")
	return nil, nil
}

// This function was auto generated
func (d *IamDriver) Delete_User(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteUserInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.DeleteUserOutput
	output, err = d.DeleteUser(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete user: %s", err)
	}
	d.logger.ExtraVerbosef("iam.DeleteUser call took %s", time.Since(start))
	d.logger.Verbose("delete user done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Attach_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["group"]; !ok {
		return nil, errors.New("attach user: missing required params 'group'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("attach user: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: attach user ok")
	return nil, nil
}

// This function was auto generated
func (d *IamDriver) Attach_User(params map[string]interface{}) (interface{}, error) {
	input := &iam.AddUserToGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["group"], input, "GroupName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.AddUserToGroupOutput
	output, err = d.AddUserToGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("attach user: %s", err)
	}
	d.logger.ExtraVerbosef("iam.AddUserToGroup call took %s", time.Since(start))
	d.logger.Verbose("attach user done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Detach_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["group"]; !ok {
		return nil, errors.New("detach user: missing required params 'group'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("detach user: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: detach user ok")
	return nil, nil
}

// This function was auto generated
func (d *IamDriver) Detach_User(params map[string]interface{}) (interface{}, error) {
	input := &iam.RemoveUserFromGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["group"], input, "GroupName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.RemoveUserFromGroupOutput
	output, err = d.RemoveUserFromGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("detach user: %s", err)
	}
	d.logger.ExtraVerbosef("iam.RemoveUserFromGroup call took %s", time.Since(start))
	d.logger.Verbose("detach user done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Create_Group_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create group: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create group ok")
	return nil, nil
}

// This function was auto generated
func (d *IamDriver) Create_Group(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "GroupName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.CreateGroupOutput
	output, err = d.CreateGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create group: %s", err)
	}
	d.logger.ExtraVerbosef("iam.CreateGroup call took %s", time.Since(start))
	id := aws.StringValue(output.Group.GroupId)

	d.logger.Verbosef("create group '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *IamDriver) Delete_Group_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete group: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete group ok")
	return nil, nil
}

// This function was auto generated
func (d *IamDriver) Delete_Group(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "GroupName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.DeleteGroupOutput
	output, err = d.DeleteGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete group: %s", err)
	}
	d.logger.ExtraVerbosef("iam.DeleteGroup call took %s", time.Since(start))
	d.logger.Verbose("delete group done")
	return output, nil
}

// This function was auto generated
func (d *S3Driver) Create_Bucket_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create bucket: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create bucket ok")
	return nil, nil
}

// This function was auto generated
func (d *S3Driver) Create_Bucket(params map[string]interface{}) (interface{}, error) {
	input := &s3.CreateBucketInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *s3.CreateBucketOutput
	output, err = d.CreateBucket(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create bucket: %s", err)
	}
	d.logger.ExtraVerbosef("s3.CreateBucket call took %s", time.Since(start))
	id := params["name"]

	d.logger.Verbosef("create bucket '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *S3Driver) Delete_Bucket_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete bucket: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete bucket ok")
	return nil, nil
}

// This function was auto generated
func (d *S3Driver) Delete_Bucket(params map[string]interface{}) (interface{}, error) {
	input := &s3.DeleteBucketInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *s3.DeleteBucketOutput
	output, err = d.DeleteBucket(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete bucket: %s", err)
	}
	d.logger.ExtraVerbosef("s3.DeleteBucket call took %s", time.Since(start))
	d.logger.Verbose("delete bucket done")
	return output, nil
}

// This function was auto generated
func (d *S3Driver) Delete_Storageobject_DryRun(params map[string]interface{}) (interface{}, error) {
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
func (d *S3Driver) Delete_Storageobject(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.DeleteObject(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete storageobject: %s", err)
	}
	d.logger.ExtraVerbosef("s3.DeleteObject call took %s", time.Since(start))
	d.logger.Verbose("delete storageobject done")
	return output, nil
}

// This function was auto generated
func (d *SnsDriver) Create_Topic_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create topic: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create topic ok")
	return nil, nil
}

// This function was auto generated
func (d *SnsDriver) Create_Topic(params map[string]interface{}) (interface{}, error) {
	input := &sns.CreateTopicInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Name", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sns.CreateTopicOutput
	output, err = d.CreateTopic(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create topic: %s", err)
	}
	d.logger.ExtraVerbosef("sns.CreateTopic call took %s", time.Since(start))
	id := aws.StringValue(output.TopicArn)

	d.logger.Verbosef("create topic '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *SnsDriver) Delete_Topic_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("delete topic: missing required params 'arn'")
	}

	d.logger.Verbose("params dry run: delete topic ok")
	return nil, nil
}

// This function was auto generated
func (d *SnsDriver) Delete_Topic(params map[string]interface{}) (interface{}, error) {
	input := &sns.DeleteTopicInput{}
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "TopicArn", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sns.DeleteTopicOutput
	output, err = d.DeleteTopic(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete topic: %s", err)
	}
	d.logger.ExtraVerbosef("sns.DeleteTopic call took %s", time.Since(start))
	d.logger.Verbose("delete topic done")
	return output, nil
}

// This function was auto generated
func (d *SnsDriver) Create_Subscription_DryRun(params map[string]interface{}) (interface{}, error) {
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
func (d *SnsDriver) Create_Subscription(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.Subscribe(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create subscription: %s", err)
	}
	d.logger.ExtraVerbosef("sns.Subscribe call took %s", time.Since(start))
	id := aws.StringValue(output.SubscriptionArn)

	d.logger.Verbosef("create subscription '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *SnsDriver) Delete_Subscription_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("delete subscription: missing required params 'arn'")
	}

	d.logger.Verbose("params dry run: delete subscription ok")
	return nil, nil
}

// This function was auto generated
func (d *SnsDriver) Delete_Subscription(params map[string]interface{}) (interface{}, error) {
	input := &sns.UnsubscribeInput{}
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "SubscriptionArn", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sns.UnsubscribeOutput
	output, err = d.Unsubscribe(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete subscription: %s", err)
	}
	d.logger.ExtraVerbosef("sns.Unsubscribe call took %s", time.Since(start))
	d.logger.Verbose("delete subscription done")
	return output, nil
}

// This function was auto generated
func (d *SqsDriver) Create_Queue_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create queue: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create queue ok")
	return nil, nil
}

// This function was auto generated
func (d *SqsDriver) Create_Queue(params map[string]interface{}) (interface{}, error) {
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
	output, err = d.CreateQueue(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create queue: %s", err)
	}
	d.logger.ExtraVerbosef("sqs.CreateQueue call took %s", time.Since(start))
	id := aws.StringValue(output.QueueUrl)

	d.logger.Verbosef("create queue '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *SqsDriver) Delete_Queue_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["url"]; !ok {
		return nil, errors.New("delete queue: missing required params 'url'")
	}

	d.logger.Verbose("params dry run: delete queue ok")
	return nil, nil
}

// This function was auto generated
func (d *SqsDriver) Delete_Queue(params map[string]interface{}) (interface{}, error) {
	input := &sqs.DeleteQueueInput{}
	var err error

	// Required params
	err = setFieldWithType(params["url"], input, "QueueUrl", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *sqs.DeleteQueueOutput
	output, err = d.DeleteQueue(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete queue: %s", err)
	}
	d.logger.ExtraVerbosef("sqs.DeleteQueue call took %s", time.Since(start))
	d.logger.Verbose("delete queue done")
	return output, nil
}

// This function was auto generated
func (d *Route53Driver) Create_Zone_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["callerreference"]; !ok {
		return nil, errors.New("create zone: missing required params 'callerreference'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("create zone: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create zone ok")
	return nil, nil
}

// This function was auto generated
func (d *Route53Driver) Create_Zone(params map[string]interface{}) (interface{}, error) {
	input := &route53.CreateHostedZoneInput{}
	var err error

	// Required params
	err = setFieldWithType(params["callerreference"], input, "CallerReference", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "Name", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["delegationsetid"]; ok {
		err = setFieldWithType(params["delegationsetid"], input, "DelegationSetId", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["comment"]; ok {
		err = setFieldWithType(params["comment"], input, "HostedZoneConfig.Comment", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["isprivate"]; ok {
		err = setFieldWithType(params["isprivate"], input, "HostedZoneConfig.PrivateZone", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["vpcid"]; ok {
		err = setFieldWithType(params["vpcid"], input, "VPC.VPCId", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["vpcregion"]; ok {
		err = setFieldWithType(params["vpcregion"], input, "VPC.VPCRegion", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *route53.CreateHostedZoneOutput
	output, err = d.CreateHostedZone(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create zone: %s", err)
	}
	d.logger.ExtraVerbosef("route53.CreateHostedZone call took %s", time.Since(start))
	id := aws.StringValue(output.HostedZone.Id)

	d.logger.Verbosef("create zone '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Route53Driver) Delete_Zone_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete zone: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete zone ok")
	return nil, nil
}

// This function was auto generated
func (d *Route53Driver) Delete_Zone(params map[string]interface{}) (interface{}, error) {
	input := &route53.DeleteHostedZoneInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "Id", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *route53.DeleteHostedZoneOutput
	output, err = d.DeleteHostedZone(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete zone: %s", err)
	}
	d.logger.ExtraVerbosef("route53.DeleteHostedZone call took %s", time.Since(start))
	d.logger.Verbose("delete zone done")
	return output, nil
}
