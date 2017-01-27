// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"strings"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
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
			d.logger.Println("full dry run: create vpc ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: create vpc error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Create_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateVpcInput{}
	
	// Required params
	setField(params["cidr"], input, "CidrBlock")
	

	output, err := d.ec2.CreateVpc(input)
	if err != nil {
		d.logger.Printf("create vpc error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.Vpc.VpcId)
	d.logger.Printf("create vpc '%s' done", id)
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
			d.logger.Println("full dry run: delete vpc ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: delete vpc error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Vpc(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVpcInput{}
	
	// Required params
	setField(params["id"], input, "VpcId")
	

	output, err := d.ec2.DeleteVpc(input)
	if err != nil {
		d.logger.Printf("delete vpc error: %s", err)
		return nil, err
	}
	d.logger.Println("delete vpc done")
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
			d.logger.Println("full dry run: create subnet ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: create subnet error: %s", err)
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
		d.logger.Printf("create subnet error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.Subnet.SubnetId)
	d.logger.Printf("create subnet '%s' done", id)
	return aws.StringValue(output.Subnet.SubnetId), nil
}

// This function was auto generated
func (d *AwsDriver) Update_Subnet_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("update subnet: missing required params 'id'")
	}
	
	d.logger.Println("params dry run: update subnet ok")
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
		d.logger.Printf("update subnet error: %s", err)
		return nil, err
	}
	d.logger.Println("update subnet done")
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
			d.logger.Println("full dry run: delete subnet ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: delete subnet error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Subnet(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSubnetInput{}
	
	// Required params
	setField(params["id"], input, "SubnetId")
	

	output, err := d.ec2.DeleteSubnet(input)
	if err != nil {
		d.logger.Printf("delete subnet error: %s", err)
		return nil, err
	}
	d.logger.Println("delete subnet done")
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
			d.logger.Println("full dry run: create instance ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: create instance error: %s", err)
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
		d.logger.Printf("create instance error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.Instances[0].InstanceId)
	tagsParams := map[string]interface{}{"resource": id}
	if v, ok := params["name"]; ok {
		tagsParams["Name"] = v
	}
	if len(tagsParams) > 1 {
		d.Create_Tags(tagsParams)
	}
	d.logger.Printf("create instance '%s' done", id)
	return aws.StringValue(output.Instances[0].InstanceId), nil
}

// This function was auto generated
func (d *AwsDriver) Update_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ModifyInstanceAttributeInput{}
	input.DryRun = aws.Bool(true)
	
	// Required params
	setField(params["id"], input, "InstanceId")
	
	// Extra params
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
			d.logger.Println("full dry run: update instance ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: update instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Update_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ModifyInstanceAttributeInput{}
	
	// Required params
	setField(params["id"], input, "InstanceId")
	
	// Extra params
	if _, ok := params["group"]; ok {
		setField(params["group"], input, "Groups")
	}
	if _, ok := params["type"]; ok {
		setField(params["type"], input, "InstanceType")
	}

	output, err := d.ec2.ModifyInstanceAttribute(input)
	if err != nil {
		d.logger.Printf("update instance error: %s", err)
		return nil, err
	}
	d.logger.Println("update instance done")
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
			d.logger.Println("full dry run: delete instance ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: delete instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.TerminateInstancesInput{}
	
	// Required params
	setField(params["id"], input, "InstanceIds")
	

	output, err := d.ec2.TerminateInstances(input)
	if err != nil {
		d.logger.Printf("delete instance error: %s", err)
		return nil, err
	}
	d.logger.Println("delete instance done")
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
			d.logger.Println("full dry run: start instance ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: start instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Start_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StartInstancesInput{}
	
	// Required params
	setField(params["id"], input, "InstanceIds")
	

	output, err := d.ec2.StartInstances(input)
	if err != nil {
		d.logger.Printf("start instance error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.StartingInstances[0].InstanceId)
	d.logger.Printf("start instance '%s' done", id)
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
			d.logger.Println("full dry run: stop instance ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: stop instance error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Stop_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.StopInstancesInput{}
	
	// Required params
	setField(params["id"], input, "InstanceIds")
	

	output, err := d.ec2.StopInstances(input)
	if err != nil {
		d.logger.Printf("stop instance error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.StoppingInstances[0].InstanceId)
	d.logger.Printf("stop instance '%s' done", id)
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
			d.logger.Println("full dry run: create securitygroup ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: create securitygroup error: %s", err)
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
		d.logger.Printf("create securitygroup error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.GroupId)
	d.logger.Printf("create securitygroup '%s' done", id)
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
			d.logger.Println("full dry run: delete securitygroup ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: delete securitygroup error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Securitygroup(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteSecurityGroupInput{}
	
	// Required params
	setField(params["id"], input, "GroupId")
	

	output, err := d.ec2.DeleteSecurityGroup(input)
	if err != nil {
		d.logger.Printf("delete securitygroup error: %s", err)
		return nil, err
	}
	d.logger.Println("delete securitygroup done")
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
			d.logger.Println("full dry run: create volume ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: create volume error: %s", err)
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
		d.logger.Printf("create volume error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.VolumeId)
	d.logger.Printf("create volume '%s' done", id)
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
			d.logger.Println("full dry run: delete volume ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: delete volume error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Volume(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteVolumeInput{}
	
	// Required params
	setField(params["id"], input, "VolumeId")
	

	output, err := d.ec2.DeleteVolume(input)
	if err != nil {
		d.logger.Printf("delete volume error: %s", err)
		return nil, err
	}
	d.logger.Println("delete volume done")
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
			d.logger.Println("full dry run: attach volume ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: attach volume error: %s", err)
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
		d.logger.Printf("attach volume error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.VolumeId)
	d.logger.Printf("attach volume '%s' done", id)
	return aws.StringValue(output.VolumeId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Keypair_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteKeyPairInput{}
	input.DryRun = aws.Bool(true)
	
	// Required params
	setField(params["name"], input, "KeyName")
	

	_, err := d.ec2.DeleteKeyPair(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fakeDryRunId("keypair")
			d.logger.Println("full dry run: delete keypair ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: delete keypair error: %s", err)
	return nil, err
}

// This function was auto generated
func (d *AwsDriver) Delete_Keypair(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteKeyPairInput{}
	
	// Required params
	setField(params["name"], input, "KeyName")
	

	output, err := d.ec2.DeleteKeyPair(input)
	if err != nil {
		d.logger.Printf("delete keypair error: %s", err)
		return nil, err
	}
	d.logger.Println("delete keypair done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create user: missing required params 'name'")
	}
	
	d.logger.Println("params dry run: create user ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Create_User(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateUserInput{}
	
	// Required params
	setField(params["name"], input, "UserName")
	

	output, err := d.iam.CreateUser(input)
	if err != nil {
		d.logger.Printf("create user error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.User.UserId)
	d.logger.Printf("create user '%s' done", id)
	return aws.StringValue(output.User.UserId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_User_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete user: missing required params 'name'")
	}
	
	d.logger.Println("params dry run: delete user ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_User(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteUserInput{}
	
	// Required params
	setField(params["name"], input, "UserName")
	

	output, err := d.iam.DeleteUser(input)
	if err != nil {
		d.logger.Printf("delete user error: %s", err)
		return nil, err
	}
	d.logger.Println("delete user done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Group_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("create group: missing required params 'name'")
	}
	
	d.logger.Println("params dry run: create group ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Group(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateGroupInput{}
	
	// Required params
	setField(params["name"], input, "GroupName")
	

	output, err := d.iam.CreateGroup(input)
	if err != nil {
		d.logger.Printf("create group error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.Group.GroupId)
	d.logger.Printf("create group '%s' done", id)
	return aws.StringValue(output.Group.GroupId), nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Group_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete group: missing required params 'name'")
	}
	
	d.logger.Println("params dry run: delete group ok")
	return nil, nil
}

// This function was auto generated
func (d *AwsDriver) Delete_Group(params map[string]interface{}) (interface{}, error) {
	input := &iam.DeleteGroupInput{}
	
	// Required params
	setField(params["name"], input, "GroupName")
	

	output, err := d.iam.DeleteGroup(input)
	if err != nil {
		d.logger.Printf("delete group error: %s", err)
		return nil, err
	}
	d.logger.Println("delete group done")
	return output, nil
}
