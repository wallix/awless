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
package awsdriver

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/rds"
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("create vpc '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("delete vpc done")
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
	if _, ok := params["availabilityzone"]; ok {
		err = setFieldWithType(params["availabilityzone"], input, "AvailabilityZone", awsstr)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.CreateSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	if _, ok := params["availabilityzone"]; ok {
		err = setFieldWithType(params["availabilityzone"], input, "AvailabilityZone", awsstr)
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

	d.logger.Infof("create subnet '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Update_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("update subnet: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: update subnet ok")
	return fakeDryRunId("subnet"), nil
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
	d.logger.Info("update subnet done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("delete subnet done")
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
	if _, ok := params["keypair"]; ok {
		err = setFieldWithType(params["keypair"], input, "KeyName", awsstr)
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
		err = setFieldWithType(params["userdata"], input, "UserData", awsfiletobase64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["securitygroup"]; ok {
		err = setFieldWithType(params["securitygroup"], input, "SecurityGroupIds", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["lock"]; ok {
		err = setFieldWithType(params["lock"], input, "DisableApiTermination", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["role"]; ok {
		err = setFieldWithType(params["role"], input, "IamInstanceProfile.Name", awsstr)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.RunInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	if _, ok := params["keypair"]; ok {
		err = setFieldWithType(params["keypair"], input, "KeyName", awsstr)
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
		err = setFieldWithType(params["userdata"], input, "UserData", awsfiletobase64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["securitygroup"]; ok {
		err = setFieldWithType(params["securitygroup"], input, "SecurityGroupIds", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["lock"]; ok {
		err = setFieldWithType(params["lock"], input, "DisableApiTermination", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["role"]; ok {
		err = setFieldWithType(params["role"], input, "IamInstanceProfile.Name", awsstr)
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

	d.logger.Infof("create instance '%s' done", id)
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
		err = setFieldWithType(params["type"], input, "InstanceType.Value", awsstr)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
		err = setFieldWithType(params["type"], input, "InstanceType.Value", awsstr)
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
	d.logger.Info("update instance done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("delete instance done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("start instance '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("stop instance '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("create securitygroup '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("delete securitygroup done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Copy_Image_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CopyImageInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Name", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["source-id"], input, "SourceImageId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["source-region"], input, "SourceRegion", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["encrypted"]; ok {
		err = setFieldWithType(params["encrypted"], input, "Encrypted", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.CopyImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("image")
			d.logger.Verbose("dry run: copy image ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: copy image: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Copy_Image(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CopyImageInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "Name", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["source-id"], input, "SourceImageId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["source-region"], input, "SourceRegion", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["encrypted"]; ok {
		err = setFieldWithType(params["encrypted"], input, "Encrypted", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.CopyImageOutput
	output, err = d.CopyImage(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("copy image: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CopyImage call took %s", time.Since(start))
	id := aws.StringValue(output.ImageId)

	d.logger.Infof("copy image '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Import_Image_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ImportImageInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Extra params
	if _, ok := params["architecture"]; ok {
		err = setFieldWithType(params["architecture"], input, "Architecture", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["license"]; ok {
		err = setFieldWithType(params["license"], input, "LicenseType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["platform"]; ok {
		err = setFieldWithType(params["platform"], input, "Platform", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["role"]; ok {
		err = setFieldWithType(params["role"], input, "RoleName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["snapshot"]; ok {
		err = setFieldWithType(params["snapshot"], input, "DiskContainers[0]SnapshotId", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["url"]; ok {
		err = setFieldWithType(params["url"], input, "DiskContainers[0]Url", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["bucket"]; ok {
		err = setFieldWithType(params["bucket"], input, "DiskContainers[0]UserBucket.S3Bucket", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["s3object"]; ok {
		err = setFieldWithType(params["s3object"], input, "DiskContainers[0]UserBucket.S3Key", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.ImportImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("image")
			d.logger.Verbose("dry run: import image ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: import image: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Import_Image(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ImportImageInput{}
	var err error

	// Extra params
	if _, ok := params["architecture"]; ok {
		err = setFieldWithType(params["architecture"], input, "Architecture", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["license"]; ok {
		err = setFieldWithType(params["license"], input, "LicenseType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["platform"]; ok {
		err = setFieldWithType(params["platform"], input, "Platform", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["role"]; ok {
		err = setFieldWithType(params["role"], input, "RoleName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["snapshot"]; ok {
		err = setFieldWithType(params["snapshot"], input, "DiskContainers[0]SnapshotId", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["url"]; ok {
		err = setFieldWithType(params["url"], input, "DiskContainers[0]Url", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["bucket"]; ok {
		err = setFieldWithType(params["bucket"], input, "DiskContainers[0]UserBucket.S3Bucket", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["s3object"]; ok {
		err = setFieldWithType(params["s3object"], input, "DiskContainers[0]UserBucket.S3Key", awsslicestruct)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.ImportImageOutput
	output, err = d.ImportImage(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("import image: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.ImportImage call took %s", time.Since(start))
	id := aws.StringValue(output.ImportTaskId)

	d.logger.Infof("import image '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Volume_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVolumeInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["availabilityzone"], input, "AvailabilityZone", awsstr)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	err = setFieldWithType(params["availabilityzone"], input, "AvailabilityZone", awsstr)
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

	d.logger.Infof("create volume '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("delete volume done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("attach volume '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("detach volume '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Snapshot_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSnapshotInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["volume"], input, "VolumeId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.CreateSnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("snapshot")
			d.logger.Verbose("dry run: create snapshot ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create snapshot: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Snapshot(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateSnapshotInput{}
	var err error

	// Required params
	err = setFieldWithType(params["volume"], input, "VolumeId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.Snapshot
	output, err = d.CreateSnapshot(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create snapshot: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateSnapshot call took %s", time.Since(start))
	id := aws.StringValue(output.SnapshotId)

	d.logger.Infof("create snapshot '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Snapshot_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSnapshotInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SnapshotId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteSnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("snapshot")
			d.logger.Verbose("dry run: delete snapshot ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete snapshot: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Snapshot(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSnapshotInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SnapshotId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteSnapshotOutput
	output, err = d.DeleteSnapshot(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete snapshot: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteSnapshot call took %s", time.Since(start))
	d.logger.Info("delete snapshot done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Copy_Snapshot_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CopySnapshotInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["source-id"], input, "SourceSnapshotId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["source-region"], input, "SourceRegion", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["encrypted"]; ok {
		err = setFieldWithType(params["encrypted"], input, "Encrypted", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.CopySnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("snapshot")
			d.logger.Verbose("dry run: copy snapshot ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: copy snapshot: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Copy_Snapshot(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CopySnapshotInput{}
	var err error

	// Required params
	err = setFieldWithType(params["source-id"], input, "SourceSnapshotId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["source-region"], input, "SourceRegion", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["encrypted"]; ok {
		err = setFieldWithType(params["encrypted"], input, "Encrypted", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.CopySnapshotOutput
	output, err = d.CopySnapshot(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("copy snapshot: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CopySnapshot call took %s", time.Since(start))
	id := aws.StringValue(output.SnapshotId)

	d.logger.Infof("copy snapshot '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("create internetgateway '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("delete internetgateway done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("attach internetgateway done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("detach internetgateway done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Natgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["elasticip-id"]; !ok {
		return nil, errors.New("create natgateway: missing required params 'elasticip-id'")
	}

	if _, ok := params["subnet"]; !ok {
		return nil, errors.New("create natgateway: missing required params 'subnet'")
	}

	d.logger.Verbose("params dry run: create natgateway ok")
	return fakeDryRunId("natgateway"), nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Natgateway(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateNatGatewayInput{}
	var err error

	// Required params
	err = setFieldWithType(params["elasticip-id"], input, "AllocationId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["subnet"], input, "SubnetId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.CreateNatGatewayOutput
	output, err = d.CreateNatGateway(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create natgateway: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateNatGateway call took %s", time.Since(start))
	id := aws.StringValue(output.NatGateway.NatGatewayId)

	d.logger.Infof("create natgateway '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Natgateway_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete natgateway: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete natgateway ok")
	return fakeDryRunId("natgateway"), nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Natgateway(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteNatGatewayInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "NatGatewayId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DeleteNatGatewayOutput
	output, err = d.DeleteNatGateway(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete natgateway: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteNatGateway call took %s", time.Since(start))
	d.logger.Info("delete natgateway done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("create routetable '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("delete routetable done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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

	d.logger.Infof("attach routetable '%s' done", id)
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("detach routetable done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("create route done")
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
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	d.logger.Info("delete route done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Keypair_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteKeyPairInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "KeyName", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DeleteKeyPair(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
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
	err = setFieldWithType(params["name"], input, "KeyName", awsstr)
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
	d.logger.Info("delete keypair done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Create_Elasticip_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AllocateAddressInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["domain"], input, "Domain", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.AllocateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("elasticip")
			d.logger.Verbose("dry run: create elasticip ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create elasticip: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Create_Elasticip(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AllocateAddressInput{}
	var err error

	// Required params
	err = setFieldWithType(params["domain"], input, "Domain", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.AllocateAddressOutput
	output, err = d.AllocateAddress(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create elasticip: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.AllocateAddress call took %s", time.Since(start))
	id := aws.StringValue(output.AllocationId)

	d.logger.Infof("create elasticip '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Delete_Elasticip_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ReleaseAddressInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Extra params
	if _, ok := params["id"]; ok {
		err = setFieldWithType(params["id"], input, "AllocationId", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["ip"]; ok {
		err = setFieldWithType(params["ip"], input, "PublicIp", awsstr)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.ReleaseAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("elasticip")
			d.logger.Verbose("dry run: delete elasticip ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete elasticip: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Delete_Elasticip(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ReleaseAddressInput{}
	var err error

	// Extra params
	if _, ok := params["id"]; ok {
		err = setFieldWithType(params["id"], input, "AllocationId", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["ip"]; ok {
		err = setFieldWithType(params["ip"], input, "PublicIp", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.ReleaseAddressOutput
	output, err = d.ReleaseAddress(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete elasticip: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.ReleaseAddress call took %s", time.Since(start))
	d.logger.Info("delete elasticip done")
	return output, nil
}

// This function was auto generated
func (d *Ec2Driver) Attach_Elasticip_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AssociateAddressInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "AllocationId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["instance"]; ok {
		err = setFieldWithType(params["instance"], input, "InstanceId", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["networkinterface"]; ok {
		err = setFieldWithType(params["networkinterface"], input, "NetworkInterfaceId", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["privateip"]; ok {
		err = setFieldWithType(params["privateip"], input, "PrivateIpAddress", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["allow-reassociation"]; ok {
		err = setFieldWithType(params["allow-reassociation"], input, "AllowReassociation", awsbool)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.AssociateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("elasticip")
			d.logger.Verbose("dry run: attach elasticip ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: attach elasticip: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Attach_Elasticip(params map[string]interface{}) (interface{}, error) {
	input := &ec2.AssociateAddressInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "AllocationId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["instance"]; ok {
		err = setFieldWithType(params["instance"], input, "InstanceId", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["networkinterface"]; ok {
		err = setFieldWithType(params["networkinterface"], input, "NetworkInterfaceId", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["privateip"]; ok {
		err = setFieldWithType(params["privateip"], input, "PrivateIpAddress", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["allow-reassociation"]; ok {
		err = setFieldWithType(params["allow-reassociation"], input, "AllowReassociation", awsbool)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.AssociateAddressOutput
	output, err = d.AssociateAddress(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("attach elasticip: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.AssociateAddress call took %s", time.Since(start))
	id := aws.StringValue(output.AssociationId)

	d.logger.Infof("attach elasticip '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Ec2Driver) Detach_Elasticip_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DisassociateAddressInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["association"], input, "AssociationId", awsstr)
	if err != nil {
		return nil, err
	}

	_, err = d.DisassociateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("elasticip")
			d.logger.Verbose("dry run: detach elasticip ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: detach elasticip: %s", err)
}

// This function was auto generated
func (d *Ec2Driver) Detach_Elasticip(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DisassociateAddressInput{}
	var err error

	// Required params
	err = setFieldWithType(params["association"], input, "AssociationId", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ec2.DisassociateAddressOutput
	output, err = d.DisassociateAddress(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("detach elasticip: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DisassociateAddress call took %s", time.Since(start))
	d.logger.Info("detach elasticip done")
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
	return fakeDryRunId("loadbalancer"), nil
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
	if _, ok := params["securitygroups"]; ok {
		err = setFieldWithType(params["securitygroups"], input, "SecurityGroups", awsstringslice)
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

	d.logger.Infof("create loadbalancer '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Loadbalancer_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete loadbalancer: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete loadbalancer ok")
	return fakeDryRunId("loadbalancer"), nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Loadbalancer(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.DeleteLoadBalancerInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "LoadBalancerArn", awsstr)
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
	d.logger.Info("delete loadbalancer done")
	return output, nil
}

// This function was auto generated
func (d *Elbv2Driver) Create_Listener_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["actiontype"]; !ok {
		return nil, errors.New("create listener: missing required params 'actiontype'")
	}

	if _, ok := params["targetgroup"]; !ok {
		return nil, errors.New("create listener: missing required params 'targetgroup'")
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
	return fakeDryRunId("listener"), nil
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
	err = setFieldWithType(params["targetgroup"], input, "DefaultActions[0]TargetGroupArn", awsslicestruct)
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

	d.logger.Infof("create listener '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Listener_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete listener: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete listener ok")
	return fakeDryRunId("listener"), nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Listener(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.DeleteListenerInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "ListenerArn", awsstr)
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
	d.logger.Info("delete listener done")
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
	return fakeDryRunId("targetgroup"), nil
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

	d.logger.Infof("create targetgroup '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Targetgroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete targetgroup: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete targetgroup ok")
	return fakeDryRunId("targetgroup"), nil
}

// This function was auto generated
func (d *Elbv2Driver) Delete_Targetgroup(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.DeleteTargetGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "TargetGroupArn", awsstr)
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
	d.logger.Info("delete targetgroup done")
	return output, nil
}

// This function was auto generated
func (d *Elbv2Driver) Attach_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["targetgroup"]; !ok {
		return nil, errors.New("attach instance: missing required params 'targetgroup'")
	}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("attach instance: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: attach instance ok")
	return fakeDryRunId("instance"), nil
}

// This function was auto generated
func (d *Elbv2Driver) Attach_Instance(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.RegisterTargetsInput{}
	var err error

	// Required params
	err = setFieldWithType(params["targetgroup"], input, "TargetGroupArn", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["id"], input, "Targets[0]Id", awsslicestruct)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["port"]; ok {
		err = setFieldWithType(params["port"], input, "Targets[0]Port", awsslicestructint64)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *elbv2.RegisterTargetsOutput
	output, err = d.RegisterTargets(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("attach instance: %s", err)
	}
	d.logger.ExtraVerbosef("elbv2.RegisterTargets call took %s", time.Since(start))
	d.logger.Info("attach instance done")
	return output, nil
}

// This function was auto generated
func (d *Elbv2Driver) Detach_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["targetgroup"]; !ok {
		return nil, errors.New("detach instance: missing required params 'targetgroup'")
	}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("detach instance: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: detach instance ok")
	return fakeDryRunId("instance"), nil
}

// This function was auto generated
func (d *Elbv2Driver) Detach_Instance(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.DeregisterTargetsInput{}
	var err error

	// Required params
	err = setFieldWithType(params["targetgroup"], input, "TargetGroupArn", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["id"], input, "Targets[0]Id", awsslicestruct)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *elbv2.DeregisterTargetsOutput
	output, err = d.DeregisterTargets(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("detach instance: %s", err)
	}
	d.logger.ExtraVerbosef("elbv2.DeregisterTargets call took %s", time.Since(start))
	d.logger.Info("detach instance done")
	return output, nil
}

// This function was auto generated
func (d *AutoscalingDriver) Create_Launchconfiguration_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["image"]; !ok {
		return nil, errors.New("create launchconfiguration: missing required params 'image'")
	}

	if _, ok := params["type"]; !ok {
		return nil, errors.New("create launchconfiguration: missing required params 'type'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("create launchconfiguration: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create launchconfiguration ok")
	return fakeDryRunId("launchconfiguration"), nil
}

// This function was auto generated
func (d *AutoscalingDriver) Create_Launchconfiguration(params map[string]interface{}) (interface{}, error) {
	input := &autoscaling.CreateLaunchConfigurationInput{}
	var err error

	// Required params
	err = setFieldWithType(params["image"], input, "ImageId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["type"], input, "InstanceType", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "LaunchConfigurationName", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["public"]; ok {
		err = setFieldWithType(params["public"], input, "AssociatePublicIpAddress", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["keypair"]; ok {
		err = setFieldWithType(params["keypair"], input, "KeyName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["userdata"]; ok {
		err = setFieldWithType(params["userdata"], input, "UserData", awsfiletobase64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["securitygroups"]; ok {
		err = setFieldWithType(params["securitygroups"], input, "SecurityGroups", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["role"]; ok {
		err = setFieldWithType(params["role"], input, "IamInstanceProfile", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["spotprice"]; ok {
		err = setFieldWithType(params["spotprice"], input, "SpotPrice", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *autoscaling.CreateLaunchConfigurationOutput
	output, err = d.CreateLaunchConfiguration(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create launchconfiguration: %s", err)
	}
	d.logger.ExtraVerbosef("autoscaling.CreateLaunchConfiguration call took %s", time.Since(start))
	id := params["name"]

	d.logger.Infof("create launchconfiguration '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *AutoscalingDriver) Delete_Launchconfiguration_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete launchconfiguration: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete launchconfiguration ok")
	return fakeDryRunId("launchconfiguration"), nil
}

// This function was auto generated
func (d *AutoscalingDriver) Delete_Launchconfiguration(params map[string]interface{}) (interface{}, error) {
	input := &autoscaling.DeleteLaunchConfigurationInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "LaunchConfigurationName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *autoscaling.DeleteLaunchConfigurationOutput
	output, err = d.DeleteLaunchConfiguration(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete launchconfiguration: %s", err)
	}
	d.logger.ExtraVerbosef("autoscaling.DeleteLaunchConfiguration call took %s", time.Since(start))
	d.logger.Info("delete launchconfiguration done")
	return output, nil
}

// This function was auto generated
func (d *AutoscalingDriver) Create_Scalinggroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create scalinggroup: missing required params 'name'")
	}

	if _, ok := params["launchconfiguration"]; !ok {
		return nil, errors.New("create scalinggroup: missing required params 'launchconfiguration'")
	}

	if _, ok := params["max-size"]; !ok {
		return nil, errors.New("create scalinggroup: missing required params 'max-size'")
	}

	if _, ok := params["min-size"]; !ok {
		return nil, errors.New("create scalinggroup: missing required params 'min-size'")
	}

	if _, ok := params["subnets"]; !ok {
		return nil, errors.New("create scalinggroup: missing required params 'subnets'")
	}

	d.logger.Verbose("params dry run: create scalinggroup ok")
	return fakeDryRunId("scalinggroup"), nil
}

// This function was auto generated
func (d *AutoscalingDriver) Create_Scalinggroup(params map[string]interface{}) (interface{}, error) {
	input := &autoscaling.CreateAutoScalingGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "AutoScalingGroupName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["launchconfiguration"], input, "LaunchConfigurationName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["max-size"], input, "MaxSize", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["min-size"], input, "MinSize", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["subnets"], input, "VPCZoneIdentifier", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["cooldown"]; ok {
		err = setFieldWithType(params["cooldown"], input, "DefaultCooldown", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["desired-capacity"]; ok {
		err = setFieldWithType(params["desired-capacity"], input, "DesiredCapacity", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthcheck-grace-period"]; ok {
		err = setFieldWithType(params["healthcheck-grace-period"], input, "HealthCheckGracePeriod", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthcheck-type"]; ok {
		err = setFieldWithType(params["healthcheck-type"], input, "HealthCheckType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["new-instances-protected"]; ok {
		err = setFieldWithType(params["new-instances-protected"], input, "NewInstancesProtectedFromScaleIn", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["targetgroups"]; ok {
		err = setFieldWithType(params["targetgroups"], input, "TargetGroupARNs", awsstringslice)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *autoscaling.CreateAutoScalingGroupOutput
	output, err = d.CreateAutoScalingGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create scalinggroup: %s", err)
	}
	d.logger.ExtraVerbosef("autoscaling.CreateAutoScalingGroup call took %s", time.Since(start))
	id := params["name"]

	d.logger.Infof("create scalinggroup '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *AutoscalingDriver) Update_Scalinggroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("update scalinggroup: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: update scalinggroup ok")
	return fakeDryRunId("scalinggroup"), nil
}

// This function was auto generated
func (d *AutoscalingDriver) Update_Scalinggroup(params map[string]interface{}) (interface{}, error) {
	input := &autoscaling.UpdateAutoScalingGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "AutoScalingGroupName", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["cooldown"]; ok {
		err = setFieldWithType(params["cooldown"], input, "DefaultCooldown", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["desired-capacity"]; ok {
		err = setFieldWithType(params["desired-capacity"], input, "DesiredCapacity", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthcheck-grace-period"]; ok {
		err = setFieldWithType(params["healthcheck-grace-period"], input, "HealthCheckGracePeriod", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["healthcheck-type"]; ok {
		err = setFieldWithType(params["healthcheck-type"], input, "HealthCheckType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["launchconfiguration"]; ok {
		err = setFieldWithType(params["launchconfiguration"], input, "LaunchConfigurationName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["max-size"]; ok {
		err = setFieldWithType(params["max-size"], input, "MaxSize", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["min-size"]; ok {
		err = setFieldWithType(params["min-size"], input, "MinSize", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["new-instances-protected"]; ok {
		err = setFieldWithType(params["new-instances-protected"], input, "NewInstancesProtectedFromScaleIn", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["subnets"]; ok {
		err = setFieldWithType(params["subnets"], input, "VPCZoneIdentifier", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *autoscaling.UpdateAutoScalingGroupOutput
	output, err = d.UpdateAutoScalingGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("update scalinggroup: %s", err)
	}
	d.logger.ExtraVerbosef("autoscaling.UpdateAutoScalingGroup call took %s", time.Since(start))
	d.logger.Info("update scalinggroup done")
	return output, nil
}

// This function was auto generated
func (d *AutoscalingDriver) Delete_Scalinggroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete scalinggroup: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete scalinggroup ok")
	return fakeDryRunId("scalinggroup"), nil
}

// This function was auto generated
func (d *AutoscalingDriver) Delete_Scalinggroup(params map[string]interface{}) (interface{}, error) {
	input := &autoscaling.DeleteAutoScalingGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "AutoScalingGroupName", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["force"]; ok {
		err = setFieldWithType(params["force"], input, "ForceDelete", awsbool)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *autoscaling.DeleteAutoScalingGroupOutput
	output, err = d.DeleteAutoScalingGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete scalinggroup: %s", err)
	}
	d.logger.ExtraVerbosef("autoscaling.DeleteAutoScalingGroup call took %s", time.Since(start))
	d.logger.Info("delete scalinggroup done")
	return output, nil
}

// This function was auto generated
func (d *AutoscalingDriver) Create_Scalingpolicy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["adjustment-type"]; !ok {
		return nil, errors.New("create scalingpolicy: missing required params 'adjustment-type'")
	}

	if _, ok := params["scalinggroup"]; !ok {
		return nil, errors.New("create scalingpolicy: missing required params 'scalinggroup'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("create scalingpolicy: missing required params 'name'")
	}

	if _, ok := params["adjustment-scaling"]; !ok {
		return nil, errors.New("create scalingpolicy: missing required params 'adjustment-scaling'")
	}

	d.logger.Verbose("params dry run: create scalingpolicy ok")
	return fakeDryRunId("scalingpolicy"), nil
}

// This function was auto generated
func (d *AutoscalingDriver) Create_Scalingpolicy(params map[string]interface{}) (interface{}, error) {
	input := &autoscaling.PutScalingPolicyInput{}
	var err error

	// Required params
	err = setFieldWithType(params["adjustment-type"], input, "AdjustmentType", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["scalinggroup"], input, "AutoScalingGroupName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "PolicyName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["adjustment-scaling"], input, "ScalingAdjustment", awsint64)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["cooldown"]; ok {
		err = setFieldWithType(params["cooldown"], input, "Cooldown", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["adjustment-magnitude"]; ok {
		err = setFieldWithType(params["adjustment-magnitude"], input, "MinAdjustmentMagnitude", awsint64)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *autoscaling.PutScalingPolicyOutput
	output, err = d.PutScalingPolicy(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create scalingpolicy: %s", err)
	}
	d.logger.ExtraVerbosef("autoscaling.PutScalingPolicy call took %s", time.Since(start))
	id := aws.StringValue(output.PolicyARN)

	d.logger.Infof("create scalingpolicy '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *AutoscalingDriver) Delete_Scalingpolicy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete scalingpolicy: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete scalingpolicy ok")
	return fakeDryRunId("scalingpolicy"), nil
}

// This function was auto generated
func (d *AutoscalingDriver) Delete_Scalingpolicy(params map[string]interface{}) (interface{}, error) {
	input := &autoscaling.DeletePolicyInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "PolicyName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *autoscaling.DeletePolicyOutput
	output, err = d.DeletePolicy(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete scalingpolicy: %s", err)
	}
	d.logger.ExtraVerbosef("autoscaling.DeletePolicy call took %s", time.Since(start))
	d.logger.Info("delete scalingpolicy done")
	return output, nil
}

// This function was auto generated
func (d *RdsDriver) Create_Database_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["type"]; !ok {
		return nil, errors.New("create database: missing required params 'type'")
	}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("create database: missing required params 'id'")
	}

	if _, ok := params["engine"]; !ok {
		return nil, errors.New("create database: missing required params 'engine'")
	}

	if _, ok := params["password"]; !ok {
		return nil, errors.New("create database: missing required params 'password'")
	}

	if _, ok := params["username"]; !ok {
		return nil, errors.New("create database: missing required params 'username'")
	}

	if _, ok := params["size"]; !ok {
		return nil, errors.New("create database: missing required params 'size'")
	}

	d.logger.Verbose("params dry run: create database ok")
	return fakeDryRunId("database"), nil
}

// This function was auto generated
func (d *RdsDriver) Create_Database(params map[string]interface{}) (interface{}, error) {
	input := &rds.CreateDBInstanceInput{}
	var err error

	// Required params
	err = setFieldWithType(params["type"], input, "DBInstanceClass", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["id"], input, "DBInstanceIdentifier", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["engine"], input, "Engine", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["password"], input, "MasterUserPassword", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["username"], input, "MasterUsername", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["size"], input, "AllocatedStorage", awsint64)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["autoupgrade"]; ok {
		err = setFieldWithType(params["autoupgrade"], input, "AutoMinorVersionUpgrade", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["availabilityzone"]; ok {
		err = setFieldWithType(params["availabilityzone"], input, "AvailabilityZone", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["backupretention"]; ok {
		err = setFieldWithType(params["backupretention"], input, "BackupRetentionPeriod", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["cluster"]; ok {
		err = setFieldWithType(params["cluster"], input, "DBClusterIdentifier", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["dbname"]; ok {
		err = setFieldWithType(params["dbname"], input, "DBName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["parametergroup"]; ok {
		err = setFieldWithType(params["parametergroup"], input, "DBParameterGroupName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["dbsecuritygroups"]; ok {
		err = setFieldWithType(params["dbsecuritygroups"], input, "DBSecurityGroups", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["subnetgroup"]; ok {
		err = setFieldWithType(params["subnetgroup"], input, "DBSubnetGroupName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["domain"]; ok {
		err = setFieldWithType(params["domain"], input, "Domain", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["iamrole"]; ok {
		err = setFieldWithType(params["iamrole"], input, "DomainIAMRoleName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["version"]; ok {
		err = setFieldWithType(params["version"], input, "EngineVersion", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["iops"]; ok {
		err = setFieldWithType(params["iops"], input, "Iops", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["license"]; ok {
		err = setFieldWithType(params["license"], input, "LicenseModel", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["multiaz"]; ok {
		err = setFieldWithType(params["multiaz"], input, "MultiAZ", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["optiongroup"]; ok {
		err = setFieldWithType(params["optiongroup"], input, "OptionGroupName", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["port"]; ok {
		err = setFieldWithType(params["port"], input, "Port", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["backupwindow"]; ok {
		err = setFieldWithType(params["backupwindow"], input, "PreferredBackupWindow", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["maintenancewindow"]; ok {
		err = setFieldWithType(params["maintenancewindow"], input, "PreferredMaintenanceWindow", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["public"]; ok {
		err = setFieldWithType(params["public"], input, "PubliclyAccessible", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["encrypted"]; ok {
		err = setFieldWithType(params["encrypted"], input, "StorageEncrypted", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["storagetype"]; ok {
		err = setFieldWithType(params["storagetype"], input, "StorageType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["timezone"]; ok {
		err = setFieldWithType(params["timezone"], input, "Timezone", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["vpcsecuritygroups"]; ok {
		err = setFieldWithType(params["vpcsecuritygroups"], input, "VpcSecurityGroupIds", awsstringslice)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *rds.CreateDBInstanceOutput
	output, err = d.CreateDBInstance(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create database: %s", err)
	}
	d.logger.ExtraVerbosef("rds.CreateDBInstance call took %s", time.Since(start))
	id := aws.StringValue(output.DBInstance.DBInstanceIdentifier)

	d.logger.Infof("create database '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *RdsDriver) Delete_Database_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete database: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete database ok")
	return fakeDryRunId("database"), nil
}

// This function was auto generated
func (d *RdsDriver) Delete_Database(params map[string]interface{}) (interface{}, error) {
	input := &rds.DeleteDBInstanceInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "DBInstanceIdentifier", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["skip-snapshot"]; ok {
		err = setFieldWithType(params["skip-snapshot"], input, "SkipFinalSnapshot", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["snapshot"]; ok {
		err = setFieldWithType(params["snapshot"], input, "FinalDBSnapshotIdentifier", awsbool)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *rds.DeleteDBInstanceOutput
	output, err = d.DeleteDBInstance(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete database: %s", err)
	}
	d.logger.ExtraVerbosef("rds.DeleteDBInstance call took %s", time.Since(start))
	d.logger.Info("delete database done")
	return output, nil
}

// This function was auto generated
func (d *RdsDriver) Create_Dbsubnetgroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["description"]; !ok {
		return nil, errors.New("create dbsubnetgroup: missing required params 'description'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("create dbsubnetgroup: missing required params 'name'")
	}

	if _, ok := params["subnets"]; !ok {
		return nil, errors.New("create dbsubnetgroup: missing required params 'subnets'")
	}

	d.logger.Verbose("params dry run: create dbsubnetgroup ok")
	return fakeDryRunId("dbsubnetgroup"), nil
}

// This function was auto generated
func (d *RdsDriver) Create_Dbsubnetgroup(params map[string]interface{}) (interface{}, error) {
	input := &rds.CreateDBSubnetGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["description"], input, "DBSubnetGroupDescription", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "DBSubnetGroupName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["subnets"], input, "SubnetIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *rds.CreateDBSubnetGroupOutput
	output, err = d.CreateDBSubnetGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create dbsubnetgroup: %s", err)
	}
	d.logger.ExtraVerbosef("rds.CreateDBSubnetGroup call took %s", time.Since(start))
	id := aws.StringValue(output.DBSubnetGroup.DBSubnetGroupName)

	d.logger.Infof("create dbsubnetgroup '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *RdsDriver) Delete_Dbsubnetgroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete dbsubnetgroup: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete dbsubnetgroup ok")
	return fakeDryRunId("dbsubnetgroup"), nil
}

// This function was auto generated
func (d *RdsDriver) Delete_Dbsubnetgroup(params map[string]interface{}) (interface{}, error) {
	input := &rds.DeleteDBSubnetGroupInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "DBSubnetGroupName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *rds.DeleteDBSubnetGroupOutput
	output, err = d.DeleteDBSubnetGroup(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete dbsubnetgroup: %s", err)
	}
	d.logger.ExtraVerbosef("rds.DeleteDBSubnetGroup call took %s", time.Since(start))
	d.logger.Info("delete dbsubnetgroup done")
	return output, nil
}

// This function was auto generated
func (d *EcrDriver) Create_Repository_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create repository: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create repository ok")
	return fakeDryRunId("repository"), nil
}

// This function was auto generated
func (d *EcrDriver) Create_Repository(params map[string]interface{}) (interface{}, error) {
	input := &ecr.CreateRepositoryInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "RepositoryName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ecr.CreateRepositoryOutput
	output, err = d.CreateRepository(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create repository: %s", err)
	}
	d.logger.ExtraVerbosef("ecr.CreateRepository call took %s", time.Since(start))
	id := aws.StringValue(output.Repository.RepositoryArn)

	d.logger.Infof("create repository '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *EcrDriver) Delete_Repository_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete repository: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete repository ok")
	return fakeDryRunId("repository"), nil
}

// This function was auto generated
func (d *EcrDriver) Delete_Repository(params map[string]interface{}) (interface{}, error) {
	input := &ecr.DeleteRepositoryInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "RepositoryName", awsstr)
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
	if _, ok := params["account"]; ok {
		err = setFieldWithType(params["account"], input, "RegistryId", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ecr.DeleteRepositoryOutput
	output, err = d.DeleteRepository(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete repository: %s", err)
	}
	d.logger.ExtraVerbosef("ecr.DeleteRepository call took %s", time.Since(start))
	d.logger.Info("delete repository done")
	return output, nil
}

// This function was auto generated
func (d *EcsDriver) Create_Containercluster_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create containercluster: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create containercluster ok")
	return fakeDryRunId("containercluster"), nil
}

// This function was auto generated
func (d *EcsDriver) Create_Containercluster(params map[string]interface{}) (interface{}, error) {
	input := &ecs.CreateClusterInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "ClusterName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ecs.CreateClusterOutput
	output, err = d.CreateCluster(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create containercluster: %s", err)
	}
	d.logger.ExtraVerbosef("ecs.CreateCluster call took %s", time.Since(start))
	id := aws.StringValue(output.Cluster.ClusterArn)

	d.logger.Infof("create containercluster '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *EcsDriver) Delete_Containercluster_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete containercluster: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete containercluster ok")
	return fakeDryRunId("containercluster"), nil
}

// This function was auto generated
func (d *EcsDriver) Delete_Containercluster(params map[string]interface{}) (interface{}, error) {
	input := &ecs.DeleteClusterInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "Cluster", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *ecs.DeleteClusterOutput
	output, err = d.DeleteCluster(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete containercluster: %s", err)
	}
	d.logger.ExtraVerbosef("ecs.DeleteCluster call took %s", time.Since(start))
	d.logger.Info("delete containercluster done")
	return output, nil
}

// This function was auto generated
func (d *EcsDriver) Update_Containertask_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["cluster"]; !ok {
		return nil, errors.New("update containertask: missing required params 'cluster'")
	}

	if _, ok := params["deployment-name"]; !ok {
		return nil, errors.New("update containertask: missing required params 'deployment-name'")
	}

	d.logger.Verbose("params dry run: update containertask ok")
	return fakeDryRunId("containertask"), nil
}

// This function was auto generated
func (d *EcsDriver) Update_Containertask(params map[string]interface{}) (interface{}, error) {
	input := &ecs.UpdateServiceInput{}
	var err error

	// Required params
	err = setFieldWithType(params["cluster"], input, "Cluster", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["deployment-name"], input, "Service", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["desired-count"]; ok {
		err = setFieldWithType(params["desired-count"], input, "DesiredCount", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["name"]; ok {
		err = setFieldWithType(params["name"], input, "TaskDefinition", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ecs.UpdateServiceOutput
	output, err = d.UpdateService(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("update containertask: %s", err)
	}
	d.logger.ExtraVerbosef("ecs.UpdateService call took %s", time.Since(start))
	d.logger.Info("update containertask done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Create_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create user: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create user ok")
	return fakeDryRunId("user"), nil
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

	d.logger.Infof("create user '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *IamDriver) Delete_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete user: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete user ok")
	return fakeDryRunId("user"), nil
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
	d.logger.Info("delete user done")
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
	return fakeDryRunId("user"), nil
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
	d.logger.Info("attach user done")
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
	return fakeDryRunId("user"), nil
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
	d.logger.Info("detach user done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Delete_Accesskey_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete accesskey: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete accesskey ok")
	return fakeDryRunId("accesskey"), nil
}

// This function was auto generated
func (d *IamDriver) Delete_Accesskey(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteAccessKeyInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "AccessKeyId", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["user"]; ok {
		err = setFieldWithType(params["user"], input, "UserName", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *iam.DeleteAccessKeyOutput
	output, err = d.DeleteAccessKey(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete accesskey: %s", err)
	}
	d.logger.ExtraVerbosef("iam.DeleteAccessKey call took %s", time.Since(start))
	d.logger.Info("delete accesskey done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Create_Loginprofile_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["username"]; !ok {
		return nil, errors.New("create loginprofile: missing required params 'username'")
	}

	if _, ok := params["password"]; !ok {
		return nil, errors.New("create loginprofile: missing required params 'password'")
	}

	d.logger.Verbose("params dry run: create loginprofile ok")
	return fakeDryRunId("loginprofile"), nil
}

// This function was auto generated
func (d *IamDriver) Create_Loginprofile(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateLoginProfileInput{}
	var err error

	// Required params
	err = setFieldWithType(params["username"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["password"], input, "Password", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["password-reset"]; ok {
		err = setFieldWithType(params["password-reset"], input, "PasswordResetRequired", awsbool)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *iam.CreateLoginProfileOutput
	output, err = d.CreateLoginProfile(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create loginprofile: %s", err)
	}
	d.logger.ExtraVerbosef("iam.CreateLoginProfile call took %s", time.Since(start))
	id := aws.StringValue(output.LoginProfile.UserName)

	d.logger.Infof("create loginprofile '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *IamDriver) Update_Loginprofile_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["username"]; !ok {
		return nil, errors.New("update loginprofile: missing required params 'username'")
	}

	if _, ok := params["password"]; !ok {
		return nil, errors.New("update loginprofile: missing required params 'password'")
	}

	d.logger.Verbose("params dry run: update loginprofile ok")
	return fakeDryRunId("loginprofile"), nil
}

// This function was auto generated
func (d *IamDriver) Update_Loginprofile(params map[string]interface{}) (interface{}, error) {
	input := &iam.UpdateLoginProfileInput{}
	var err error

	// Required params
	err = setFieldWithType(params["username"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["password"], input, "Password", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["password-reset"]; ok {
		err = setFieldWithType(params["password-reset"], input, "PasswordResetRequired", awsbool)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *iam.UpdateLoginProfileOutput
	output, err = d.UpdateLoginProfile(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("update loginprofile: %s", err)
	}
	d.logger.ExtraVerbosef("iam.UpdateLoginProfile call took %s", time.Since(start))
	d.logger.Info("update loginprofile done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Delete_Loginprofile_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["username"]; !ok {
		return nil, errors.New("delete loginprofile: missing required params 'username'")
	}

	d.logger.Verbose("params dry run: delete loginprofile ok")
	return fakeDryRunId("loginprofile"), nil
}

// This function was auto generated
func (d *IamDriver) Delete_Loginprofile(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteLoginProfileInput{}
	var err error

	// Required params
	err = setFieldWithType(params["username"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.DeleteLoginProfileOutput
	output, err = d.DeleteLoginProfile(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete loginprofile: %s", err)
	}
	d.logger.ExtraVerbosef("iam.DeleteLoginProfile call took %s", time.Since(start))
	d.logger.Info("delete loginprofile done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Create_Group_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create group: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create group ok")
	return fakeDryRunId("group"), nil
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

	d.logger.Infof("create group '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *IamDriver) Delete_Group_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete group: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete group ok")
	return fakeDryRunId("group"), nil
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
	d.logger.Info("delete group done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Attach_Role_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["instanceprofile"]; !ok {
		return nil, errors.New("attach role: missing required params 'instanceprofile'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("attach role: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: attach role ok")
	return fakeDryRunId("role"), nil
}

// This function was auto generated
func (d *IamDriver) Attach_Role(params map[string]interface{}) (interface{}, error) {
	input := &iam.AddRoleToInstanceProfileInput{}
	var err error

	// Required params
	err = setFieldWithType(params["instanceprofile"], input, "InstanceProfileName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "RoleName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.AddRoleToInstanceProfileOutput
	output, err = d.AddRoleToInstanceProfile(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("attach role: %s", err)
	}
	d.logger.ExtraVerbosef("iam.AddRoleToInstanceProfile call took %s", time.Since(start))
	d.logger.Info("attach role done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Detach_Role_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["instanceprofile"]; !ok {
		return nil, errors.New("detach role: missing required params 'instanceprofile'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("detach role: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: detach role ok")
	return fakeDryRunId("role"), nil
}

// This function was auto generated
func (d *IamDriver) Detach_Role(params map[string]interface{}) (interface{}, error) {
	input := &iam.RemoveRoleFromInstanceProfileInput{}
	var err error

	// Required params
	err = setFieldWithType(params["instanceprofile"], input, "InstanceProfileName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "RoleName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.RemoveRoleFromInstanceProfileOutput
	output, err = d.RemoveRoleFromInstanceProfile(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("detach role: %s", err)
	}
	d.logger.ExtraVerbosef("iam.RemoveRoleFromInstanceProfile call took %s", time.Since(start))
	d.logger.Info("detach role done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Create_Instanceprofile_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create instanceprofile: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create instanceprofile ok")
	return fakeDryRunId("instanceprofile"), nil
}

// This function was auto generated
func (d *IamDriver) Create_Instanceprofile(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateInstanceProfileInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "InstanceProfileName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.CreateInstanceProfileOutput
	output, err = d.CreateInstanceProfile(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create instanceprofile: %s", err)
	}
	d.logger.ExtraVerbosef("iam.CreateInstanceProfile call took %s", time.Since(start))
	d.logger.Info("create instanceprofile done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Delete_Instanceprofile_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete instanceprofile: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete instanceprofile ok")
	return fakeDryRunId("instanceprofile"), nil
}

// This function was auto generated
func (d *IamDriver) Delete_Instanceprofile(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteInstanceProfileInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "InstanceProfileName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.DeleteInstanceProfileOutput
	output, err = d.DeleteInstanceProfile(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete instanceprofile: %s", err)
	}
	d.logger.ExtraVerbosef("iam.DeleteInstanceProfile call took %s", time.Since(start))
	d.logger.Info("delete instanceprofile done")
	return output, nil
}

// This function was auto generated
func (d *IamDriver) Delete_Policy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("delete policy: missing required params 'arn'")
	}

	d.logger.Verbose("params dry run: delete policy ok")
	return fakeDryRunId("policy"), nil
}

// This function was auto generated
func (d *IamDriver) Delete_Policy(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeletePolicyInput{}
	var err error

	// Required params
	err = setFieldWithType(params["arn"], input, "PolicyArn", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.DeletePolicyOutput
	output, err = d.DeletePolicy(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete policy: %s", err)
	}
	d.logger.ExtraVerbosef("iam.DeletePolicy call took %s", time.Since(start))
	d.logger.Info("delete policy done")
	return output, nil
}

// This function was auto generated
func (d *S3Driver) Create_Bucket_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create bucket: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create bucket ok")
	return fakeDryRunId("bucket"), nil
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

	// Extra params
	if _, ok := params["acl"]; ok {
		err = setFieldWithType(params["acl"], input, "ACL", awsstr)
		if err != nil {
			return nil, err
		}
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

	d.logger.Infof("create bucket '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *S3Driver) Delete_Bucket_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete bucket: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete bucket ok")
	return fakeDryRunId("bucket"), nil
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
	d.logger.Info("delete bucket done")
	return output, nil
}

// This function was auto generated
func (d *S3Driver) Update_S3object_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["bucket"]; !ok {
		return nil, errors.New("update s3object: missing required params 'bucket'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("update s3object: missing required params 'name'")
	}

	if _, ok := params["acl"]; !ok {
		return nil, errors.New("update s3object: missing required params 'acl'")
	}

	d.logger.Verbose("params dry run: update s3object ok")
	return fakeDryRunId("s3object"), nil
}

// This function was auto generated
func (d *S3Driver) Update_S3object(params map[string]interface{}) (interface{}, error) {
	input := &s3.PutObjectAclInput{}
	var err error

	// Required params
	err = setFieldWithType(params["bucket"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "Key", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["acl"], input, "ACL", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["version"]; ok {
		err = setFieldWithType(params["version"], input, "VersionId", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *s3.PutObjectAclOutput
	output, err = d.PutObjectAcl(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("update s3object: %s", err)
	}
	d.logger.ExtraVerbosef("s3.PutObjectAcl call took %s", time.Since(start))
	d.logger.Info("update s3object done")
	return output, nil
}

// This function was auto generated
func (d *S3Driver) Delete_S3object_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["bucket"]; !ok {
		return nil, errors.New("delete s3object: missing required params 'bucket'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete s3object: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete s3object ok")
	return fakeDryRunId("s3object"), nil
}

// This function was auto generated
func (d *S3Driver) Delete_S3object(params map[string]interface{}) (interface{}, error) {
	input := &s3.DeleteObjectInput{}
	var err error

	// Required params
	err = setFieldWithType(params["bucket"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], input, "Key", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *s3.DeleteObjectOutput
	output, err = d.DeleteObject(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete s3object: %s", err)
	}
	d.logger.ExtraVerbosef("s3.DeleteObject call took %s", time.Since(start))
	d.logger.Info("delete s3object done")
	return output, nil
}

// This function was auto generated
func (d *SnsDriver) Create_Topic_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create topic: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create topic ok")
	return fakeDryRunId("topic"), nil
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

	d.logger.Infof("create topic '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *SnsDriver) Delete_Topic_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete topic: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete topic ok")
	return fakeDryRunId("topic"), nil
}

// This function was auto generated
func (d *SnsDriver) Delete_Topic(params map[string]interface{}) (interface{}, error) {
	input := &sns.DeleteTopicInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "TopicArn", awsstr)
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
	d.logger.Info("delete topic done")
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
	return fakeDryRunId("subscription"), nil
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

	d.logger.Infof("create subscription '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *SnsDriver) Delete_Subscription_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete subscription: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete subscription ok")
	return fakeDryRunId("subscription"), nil
}

// This function was auto generated
func (d *SnsDriver) Delete_Subscription(params map[string]interface{}) (interface{}, error) {
	input := &sns.UnsubscribeInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "SubscriptionArn", awsstr)
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
	d.logger.Info("delete subscription done")
	return output, nil
}

// This function was auto generated
func (d *SqsDriver) Create_Queue_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create queue: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: create queue ok")
	return fakeDryRunId("queue"), nil
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
	if _, ok := params["max-msg-size"]; ok {
		err = setFieldWithType(params["max-msg-size"], input, "Attributes[MaximumMessageSize]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["retention-period"]; ok {
		err = setFieldWithType(params["retention-period"], input, "Attributes[MessageRetentionPeriod]", awsstringpointermap)
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
	if _, ok := params["msg-wait"]; ok {
		err = setFieldWithType(params["msg-wait"], input, "Attributes[ReceiveMessageWaitTimeSeconds]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["redrive-policy"]; ok {
		err = setFieldWithType(params["redrive-policy"], input, "Attributes[RedrivePolicy]", awsstringpointermap)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["visibility-timeout"]; ok {
		err = setFieldWithType(params["visibility-timeout"], input, "Attributes[VisibilityTimeout]", awsstringpointermap)
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

	d.logger.Infof("create queue '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *SqsDriver) Delete_Queue_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["url"]; !ok {
		return nil, errors.New("delete queue: missing required params 'url'")
	}

	d.logger.Verbose("params dry run: delete queue ok")
	return fakeDryRunId("queue"), nil
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
	d.logger.Info("delete queue done")
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
	return fakeDryRunId("zone"), nil
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

	d.logger.Infof("create zone '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *Route53Driver) Delete_Zone_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete zone: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete zone ok")
	return fakeDryRunId("zone"), nil
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
	d.logger.Info("delete zone done")
	return output, nil
}

// This function was auto generated
func (d *LambdaDriver) Create_Function_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create function: missing required params 'name'")
	}

	if _, ok := params["handler"]; !ok {
		return nil, errors.New("create function: missing required params 'handler'")
	}

	if _, ok := params["role"]; !ok {
		return nil, errors.New("create function: missing required params 'role'")
	}

	if _, ok := params["runtime"]; !ok {
		return nil, errors.New("create function: missing required params 'runtime'")
	}

	d.logger.Verbose("params dry run: create function ok")
	return fakeDryRunId("function"), nil
}

// This function was auto generated
func (d *LambdaDriver) Create_Function(params map[string]interface{}) (interface{}, error) {
	input := &lambda.CreateFunctionInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "FunctionName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["handler"], input, "Handler", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["role"], input, "Role", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["runtime"], input, "Runtime", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["bucket"]; ok {
		err = setFieldWithType(params["bucket"], input, "Code.S3Bucket", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["object"]; ok {
		err = setFieldWithType(params["object"], input, "Code.S3Key", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["objectversion"]; ok {
		err = setFieldWithType(params["objectversion"], input, "Code.S3ObjectVersion", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["zipfile"]; ok {
		err = setFieldWithType(params["zipfile"], input, "Code.ZipFile", awsfiletobyteslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "Description", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["memory"]; ok {
		err = setFieldWithType(params["memory"], input, "MemorySize", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["publish"]; ok {
		err = setFieldWithType(params["publish"], input, "Publish", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["timeout"]; ok {
		err = setFieldWithType(params["timeout"], input, "Timeout", awsint64)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *lambda.FunctionConfiguration
	output, err = d.CreateFunction(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create function: %s", err)
	}
	d.logger.ExtraVerbosef("lambda.CreateFunction call took %s", time.Since(start))
	id := aws.StringValue(output.FunctionArn)

	d.logger.Infof("create function '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *LambdaDriver) Delete_Function_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("delete function: missing required params 'id'")
	}

	d.logger.Verbose("params dry run: delete function ok")
	return fakeDryRunId("function"), nil
}

// This function was auto generated
func (d *LambdaDriver) Delete_Function(params map[string]interface{}) (interface{}, error) {
	input := &lambda.DeleteFunctionInput{}
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "FunctionName", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["version"]; ok {
		err = setFieldWithType(params["version"], input, "Qualifier", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *lambda.DeleteFunctionOutput
	output, err = d.DeleteFunction(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete function: %s", err)
	}
	d.logger.ExtraVerbosef("lambda.DeleteFunction call took %s", time.Since(start))
	d.logger.Info("delete function done")
	return output, nil
}

// This function was auto generated
func (d *CloudwatchDriver) Create_Alarm_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create alarm: missing required params 'name'")
	}

	if _, ok := params["operator"]; !ok {
		return nil, errors.New("create alarm: missing required params 'operator'")
	}

	if _, ok := params["metric"]; !ok {
		return nil, errors.New("create alarm: missing required params 'metric'")
	}

	if _, ok := params["namespace"]; !ok {
		return nil, errors.New("create alarm: missing required params 'namespace'")
	}

	if _, ok := params["evaluation-periods"]; !ok {
		return nil, errors.New("create alarm: missing required params 'evaluation-periods'")
	}

	if _, ok := params["period"]; !ok {
		return nil, errors.New("create alarm: missing required params 'period'")
	}

	if _, ok := params["statistic-function"]; !ok {
		return nil, errors.New("create alarm: missing required params 'statistic-function'")
	}

	if _, ok := params["threshold"]; !ok {
		return nil, errors.New("create alarm: missing required params 'threshold'")
	}

	d.logger.Verbose("params dry run: create alarm ok")
	return fakeDryRunId("alarm"), nil
}

// This function was auto generated
func (d *CloudwatchDriver) Create_Alarm(params map[string]interface{}) (interface{}, error) {
	input := &cloudwatch.PutMetricAlarmInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "AlarmName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["operator"], input, "ComparisonOperator", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["metric"], input, "MetricName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["namespace"], input, "Namespace", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["evaluation-periods"], input, "EvaluationPeriods", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["period"], input, "Period", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["statistic-function"], input, "Statistic", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["threshold"], input, "Threshold", awsfloat)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["enabled"]; ok {
		err = setFieldWithType(params["enabled"], input, "ActionsEnabled", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["alarm-actions"]; ok {
		err = setFieldWithType(params["alarm-actions"], input, "AlarmActions", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["insufficientdata-actions"]; ok {
		err = setFieldWithType(params["insufficientdata-actions"], input, "InsufficientDataActions", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["ok-actions"]; ok {
		err = setFieldWithType(params["ok-actions"], input, "OKActions", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["description"]; ok {
		err = setFieldWithType(params["description"], input, "AlarmDescription", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["dimensions"]; ok {
		err = setFieldWithType(params["dimensions"], input, "Dimensions", awsdimensionslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["unit"]; ok {
		err = setFieldWithType(params["unit"], input, "Unit", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *cloudwatch.PutMetricAlarmOutput
	output, err = d.PutMetricAlarm(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create alarm: %s", err)
	}
	d.logger.ExtraVerbosef("cloudwatch.PutMetricAlarm call took %s", time.Since(start))
	id := params["name"]

	d.logger.Infof("create alarm '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *CloudwatchDriver) Delete_Alarm_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete alarm: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete alarm ok")
	return fakeDryRunId("alarm"), nil
}

// This function was auto generated
func (d *CloudwatchDriver) Delete_Alarm(params map[string]interface{}) (interface{}, error) {
	input := &cloudwatch.DeleteAlarmsInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "AlarmNames", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *cloudwatch.DeleteAlarmsOutput
	output, err = d.DeleteAlarms(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete alarm: %s", err)
	}
	d.logger.ExtraVerbosef("cloudwatch.DeleteAlarms call took %s", time.Since(start))
	d.logger.Info("delete alarm done")
	return output, nil
}

// This function was auto generated
func (d *CloudwatchDriver) Start_Alarm_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["names"]; !ok {
		return nil, errors.New("start alarm: missing required params 'names'")
	}

	d.logger.Verbose("params dry run: start alarm ok")
	return fakeDryRunId("alarm"), nil
}

// This function was auto generated
func (d *CloudwatchDriver) Start_Alarm(params map[string]interface{}) (interface{}, error) {
	input := &cloudwatch.EnableAlarmActionsInput{}
	var err error

	// Required params
	err = setFieldWithType(params["names"], input, "AlarmNames", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *cloudwatch.EnableAlarmActionsOutput
	output, err = d.EnableAlarmActions(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("start alarm: %s", err)
	}
	d.logger.ExtraVerbosef("cloudwatch.EnableAlarmActions call took %s", time.Since(start))
	d.logger.Info("start alarm done")
	return output, nil
}

// This function was auto generated
func (d *CloudwatchDriver) Stop_Alarm_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["names"]; !ok {
		return nil, errors.New("stop alarm: missing required params 'names'")
	}

	d.logger.Verbose("params dry run: stop alarm ok")
	return fakeDryRunId("alarm"), nil
}

// This function was auto generated
func (d *CloudwatchDriver) Stop_Alarm(params map[string]interface{}) (interface{}, error) {
	input := &cloudwatch.DisableAlarmActionsInput{}
	var err error

	// Required params
	err = setFieldWithType(params["names"], input, "AlarmNames", awsstringslice)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *cloudwatch.DisableAlarmActionsOutput
	output, err = d.DisableAlarmActions(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("stop alarm: %s", err)
	}
	d.logger.ExtraVerbosef("cloudwatch.DisableAlarmActions call took %s", time.Since(start))
	d.logger.Info("stop alarm done")
	return output, nil
}

// This function was auto generated
func (d *CloudformationDriver) Create_Stack_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create stack: missing required params 'name'")
	}

	if _, ok := params["template-file"]; !ok {
		return nil, errors.New("create stack: missing required params 'template-file'")
	}

	d.logger.Verbose("params dry run: create stack ok")
	return fakeDryRunId("stack"), nil
}

// This function was auto generated
func (d *CloudformationDriver) Create_Stack(params map[string]interface{}) (interface{}, error) {
	input := &cloudformation.CreateStackInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "StackName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["template-file"], input, "TemplateBody", awsfiletostring)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["capabilities"]; ok {
		err = setFieldWithType(params["capabilities"], input, "Capabilities", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["disable-rollback"]; ok {
		err = setFieldWithType(params["disable-rollback"], input, "DisableRollback", awsbool)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["notifications"]; ok {
		err = setFieldWithType(params["notifications"], input, "NotificationARNs", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["on-failure"]; ok {
		err = setFieldWithType(params["on-failure"], input, "OnFailure", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["parameters"]; ok {
		err = setFieldWithType(params["parameters"], input, "Parameters", awsparameterslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["resource-types"]; ok {
		err = setFieldWithType(params["resource-types"], input, "ResourceTypes", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["role"]; ok {
		err = setFieldWithType(params["role"], input, "RoleARN", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["policy-file"]; ok {
		err = setFieldWithType(params["policy-file"], input, "StackPolicyBody", awsfiletostring)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["timeout"]; ok {
		err = setFieldWithType(params["timeout"], input, "TimeoutInMinutes", awsint64)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *cloudformation.CreateStackOutput
	output, err = d.CreateStack(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create stack: %s", err)
	}
	d.logger.ExtraVerbosef("cloudformation.CreateStack call took %s", time.Since(start))
	id := aws.StringValue(output.StackId)

	d.logger.Infof("create stack '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *CloudformationDriver) Update_Stack_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("update stack: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: update stack ok")
	return fakeDryRunId("stack"), nil
}

// This function was auto generated
func (d *CloudformationDriver) Update_Stack(params map[string]interface{}) (interface{}, error) {
	input := &cloudformation.UpdateStackInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "StackName", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["capabilities"]; ok {
		err = setFieldWithType(params["capabilities"], input, "Capabilities", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["notifications"]; ok {
		err = setFieldWithType(params["notifications"], input, "NotificationARNs", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["parameters"]; ok {
		err = setFieldWithType(params["parameters"], input, "Parameters", awsparameterslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["resource-types"]; ok {
		err = setFieldWithType(params["resource-types"], input, "ResourceTypes", awsstringslice)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["role"]; ok {
		err = setFieldWithType(params["role"], input, "RoleARN", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["policy-file"]; ok {
		err = setFieldWithType(params["policy-file"], input, "StackPolicyBody", awsfiletostring)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["policy-update-file"]; ok {
		err = setFieldWithType(params["policy-update-file"], input, "StackPolicyDuringUpdateBody", awsfiletostring)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["template-file"]; ok {
		err = setFieldWithType(params["template-file"], input, "TemplateBody", awsfiletostring)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["use-previous-template"]; ok {
		err = setFieldWithType(params["use-previous-template"], input, "UsePreviousTemplate", awsbool)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *cloudformation.UpdateStackOutput
	output, err = d.UpdateStack(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("update stack: %s", err)
	}
	d.logger.ExtraVerbosef("cloudformation.UpdateStack call took %s", time.Since(start))
	id := aws.StringValue(output.StackId)

	d.logger.Infof("update stack '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *CloudformationDriver) Delete_Stack_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete stack: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete stack ok")
	return fakeDryRunId("stack"), nil
}

// This function was auto generated
func (d *CloudformationDriver) Delete_Stack(params map[string]interface{}) (interface{}, error) {
	input := &cloudformation.DeleteStackInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "StackName", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["retain-resources"]; ok {
		err = setFieldWithType(params["retain-resources"], input, "RetainResources", awsstringslice)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *cloudformation.DeleteStackOutput
	output, err = d.DeleteStack(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete stack: %s", err)
	}
	d.logger.ExtraVerbosef("cloudformation.DeleteStack call took %s", time.Since(start))
	d.logger.Info("delete stack done")
	return output, nil
}

// This function was auto generated
func (d *ApplicationautoscalingDriver) Create_Appscalingtarget_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["max-capacity"]; !ok {
		return nil, errors.New("create appscalingtarget: missing required params 'max-capacity'")
	}

	if _, ok := params["min-capacity"]; !ok {
		return nil, errors.New("create appscalingtarget: missing required params 'min-capacity'")
	}

	if _, ok := params["resource"]; !ok {
		return nil, errors.New("create appscalingtarget: missing required params 'resource'")
	}

	if _, ok := params["role"]; !ok {
		return nil, errors.New("create appscalingtarget: missing required params 'role'")
	}

	if _, ok := params["dimension"]; !ok {
		return nil, errors.New("create appscalingtarget: missing required params 'dimension'")
	}

	if _, ok := params["service-namespace"]; !ok {
		return nil, errors.New("create appscalingtarget: missing required params 'service-namespace'")
	}

	d.logger.Verbose("params dry run: create appscalingtarget ok")
	return fakeDryRunId("appscalingtarget"), nil
}

// This function was auto generated
func (d *ApplicationautoscalingDriver) Create_Appscalingtarget(params map[string]interface{}) (interface{}, error) {
	input := &applicationautoscaling.RegisterScalableTargetInput{}
	var err error

	// Required params
	err = setFieldWithType(params["max-capacity"], input, "MaxCapacity", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["min-capacity"], input, "MinCapacity", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["resource"], input, "ResourceId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["role"], input, "RoleARN", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["dimension"], input, "ScalableDimension", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["service-namespace"], input, "ServiceNamespace", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *applicationautoscaling.RegisterScalableTargetOutput
	output, err = d.RegisterScalableTarget(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create appscalingtarget: %s", err)
	}
	d.logger.ExtraVerbosef("applicationautoscaling.RegisterScalableTarget call took %s", time.Since(start))
	d.logger.Info("create appscalingtarget done")
	return output, nil
}

// This function was auto generated
func (d *ApplicationautoscalingDriver) Delete_Appscalingtarget_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["resource"]; !ok {
		return nil, errors.New("delete appscalingtarget: missing required params 'resource'")
	}

	if _, ok := params["dimension"]; !ok {
		return nil, errors.New("delete appscalingtarget: missing required params 'dimension'")
	}

	if _, ok := params["service-namespace"]; !ok {
		return nil, errors.New("delete appscalingtarget: missing required params 'service-namespace'")
	}

	d.logger.Verbose("params dry run: delete appscalingtarget ok")
	return fakeDryRunId("appscalingtarget"), nil
}

// This function was auto generated
func (d *ApplicationautoscalingDriver) Delete_Appscalingtarget(params map[string]interface{}) (interface{}, error) {
	input := &applicationautoscaling.DeregisterScalableTargetInput{}
	var err error

	// Required params
	err = setFieldWithType(params["resource"], input, "ResourceId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["dimension"], input, "ScalableDimension", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["service-namespace"], input, "ServiceNamespace", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *applicationautoscaling.DeregisterScalableTargetOutput
	output, err = d.DeregisterScalableTarget(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete appscalingtarget: %s", err)
	}
	d.logger.ExtraVerbosef("applicationautoscaling.DeregisterScalableTarget call took %s", time.Since(start))
	d.logger.Info("delete appscalingtarget done")
	return output, nil
}

// This function was auto generated
func (d *ApplicationautoscalingDriver) Create_Appscalingpolicy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create appscalingpolicy: missing required params 'name'")
	}

	if _, ok := params["type"]; !ok {
		return nil, errors.New("create appscalingpolicy: missing required params 'type'")
	}

	if _, ok := params["resource"]; !ok {
		return nil, errors.New("create appscalingpolicy: missing required params 'resource'")
	}

	if _, ok := params["dimension"]; !ok {
		return nil, errors.New("create appscalingpolicy: missing required params 'dimension'")
	}

	if _, ok := params["service-namespace"]; !ok {
		return nil, errors.New("create appscalingpolicy: missing required params 'service-namespace'")
	}

	if _, ok := params["stepscaling-adjustment-type"]; !ok {
		return nil, errors.New("create appscalingpolicy: missing required params 'stepscaling-adjustment-type'")
	}

	if _, ok := params["stepscaling-adjustments"]; !ok {
		return nil, errors.New("create appscalingpolicy: missing required params 'stepscaling-adjustments'")
	}

	d.logger.Verbose("params dry run: create appscalingpolicy ok")
	return fakeDryRunId("appscalingpolicy"), nil
}

// This function was auto generated
func (d *ApplicationautoscalingDriver) Create_Appscalingpolicy(params map[string]interface{}) (interface{}, error) {
	input := &applicationautoscaling.PutScalingPolicyInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "PolicyName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["type"], input, "PolicyType", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["resource"], input, "ResourceId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["dimension"], input, "ScalableDimension", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["service-namespace"], input, "ServiceNamespace", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["stepscaling-adjustment-type"], input, "StepScalingPolicyConfiguration.AdjustmentType", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["stepscaling-adjustments"], input, "StepScalingPolicyConfiguration.StepAdjustments", awsstepadjustments)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["stepscaling-cooldown"]; ok {
		err = setFieldWithType(params["stepscaling-cooldown"], input, "StepScalingPolicyConfiguration.Cooldown", awsint64)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["stepscaling-aggregation-type"]; ok {
		err = setFieldWithType(params["stepscaling-aggregation-type"], input, "StepScalingPolicyConfiguration.MetricAggregationType", awsstr)
		if err != nil {
			return nil, err
		}
	}
	if _, ok := params["stepscaling-min-adjustment-magnitude"]; ok {
		err = setFieldWithType(params["stepscaling-min-adjustment-magnitude"], input, "StepScalingPolicyConfiguration.MinAdjustmentMagnitude", awsint64)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *applicationautoscaling.PutScalingPolicyOutput
	output, err = d.PutScalingPolicy(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("create appscalingpolicy: %s", err)
	}
	d.logger.ExtraVerbosef("applicationautoscaling.PutScalingPolicy call took %s", time.Since(start))
	id := aws.StringValue(output.PolicyARN)

	d.logger.Infof("create appscalingpolicy '%s' done", id)
	return id, nil
}

// This function was auto generated
func (d *ApplicationautoscalingDriver) Delete_Appscalingpolicy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete appscalingpolicy: missing required params 'name'")
	}

	if _, ok := params["resource"]; !ok {
		return nil, errors.New("delete appscalingpolicy: missing required params 'resource'")
	}

	if _, ok := params["dimension"]; !ok {
		return nil, errors.New("delete appscalingpolicy: missing required params 'dimension'")
	}

	if _, ok := params["service-namespace"]; !ok {
		return nil, errors.New("delete appscalingpolicy: missing required params 'service-namespace'")
	}

	d.logger.Verbose("params dry run: delete appscalingpolicy ok")
	return fakeDryRunId("appscalingpolicy"), nil
}

// This function was auto generated
func (d *ApplicationautoscalingDriver) Delete_Appscalingpolicy(params map[string]interface{}) (interface{}, error) {
	input := &applicationautoscaling.DeleteScalingPolicyInput{}
	var err error

	// Required params
	err = setFieldWithType(params["name"], input, "PolicyName", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["resource"], input, "ResourceId", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["dimension"], input, "ScalableDimension", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["service-namespace"], input, "ServiceNamespace", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *applicationautoscaling.DeleteScalingPolicyOutput
	output, err = d.DeleteScalingPolicy(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete appscalingpolicy: %s", err)
	}
	d.logger.ExtraVerbosef("applicationautoscaling.DeleteScalingPolicy call took %s", time.Since(start))
	d.logger.Info("delete appscalingpolicy done")
	return output, nil
}
