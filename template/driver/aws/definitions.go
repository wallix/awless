package aws

func DriverSupportedActions() map[string][]string {
	supported := make(map[string][]string)
	for _, def := range DriverDefinitions {
		supported[def.Action] = append(supported[def.Action], def.Entity)
	}
	return supported
}

var DriverDefinitions = []struct {
	RequiredParams                    map[string]string
	ExtraParams                       map[string]string
	TagsMapping                       map[string]string
	Api                               string
	Action, Entity                    string
	Input, ApiMethod, OutputExtractor string
	DryRunUnsupported                 bool
	ManualFuncDefinition              bool
}{

	//// EC2

	// VPC
	{
		Action: "create", Entity: "vpc", Api: "ec2", Input: "CreateVpcInput", ApiMethod: "CreateVpc", OutputExtractor: "Vpc.VpcId",
		RequiredParams: map[string]string{
			"CidrBlock": "cidr",
		},
	},
	{
		Action: "delete", Entity: "vpc", Api: "ec2", Input: "DeleteVpcInput", ApiMethod: "DeleteVpc",
		RequiredParams: map[string]string{
			"VpcId": "id",
		},
	},

	// SUBNET
	{
		Action: "create", Entity: "subnet", Api: "ec2", Input: "CreateSubnetInput", ApiMethod: "CreateSubnet", OutputExtractor: "Subnet.SubnetId",
		RequiredParams: map[string]string{
			"CidrBlock": "cidr",
			"VpcId":     "vpc",
		},
		ExtraParams: map[string]string{
			"AvailabilityZone": "zone",
		},
	},
	{
		Action: "update", Entity: "subnet", Api: "ec2", Input: "ModifySubnetAttributeInput", ApiMethod: "ModifySubnetAttribute", DryRunUnsupported: true,
		RequiredParams: map[string]string{
			"SubnetId": "id",
		},
		ExtraParams: map[string]string{
			"MapPublicIpOnLaunch": "public-vms",
		},
	},
	{
		Action: "delete", Entity: "subnet", Api: "ec2", Input: "DeleteSubnetInput", ApiMethod: "DeleteSubnet",
		RequiredParams: map[string]string{
			"SubnetId": "id",
		},
	},

	// INSTANCES
	{
		Action: "create", Entity: "instance", Api: "ec2", Input: "RunInstancesInput", ApiMethod: "RunInstances", OutputExtractor: "Instances[0].InstanceId",
		RequiredParams: map[string]string{
			"ImageId":      "image",
			"MaxCount":     "count",
			"MinCount":     "count",
			"InstanceType": "type",
			"SubnetId":     "subnet",
		},
		ExtraParams: map[string]string{
			"KeyName":          "key",
			"PrivateIpAddress": "ip",
			"UserData":         "userdata",
			"SecurityGroupIds": "group",
		},
		TagsMapping: map[string]string{
			"Name": "name",
		},
	},
	{
		Action: "update", Entity: "instance", Api: "ec2", Input: "ModifyInstanceAttributeInput", ApiMethod: "ModifyInstanceAttribute",
		RequiredParams: map[string]string{
			"InstanceId": "id",
		},
		ExtraParams: map[string]string{
			"InstanceType": "type",
			"Groups":       "group",
		},
	},
	{
		Action: "delete", Entity: "instance", Api: "ec2", Input: "TerminateInstancesInput", ApiMethod: "TerminateInstances",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},
	{
		Action: "start", Entity: "instance", Api: "ec2", Input: "StartInstancesInput", ApiMethod: "StartInstances", OutputExtractor: "StartingInstances[0].InstanceId",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},
	{
		Action: "stop", Entity: "instance", Api: "ec2", Input: "StopInstancesInput", ApiMethod: "StopInstances", OutputExtractor: "StoppingInstances[0].InstanceId",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},
	{
		Action: "check", Entity: "instance", Api: "ec2", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"Id":      "id",
			"State":   "state",
			"Timeout": "timeout",
		},
	},

	// Security Group
	{
		Action: "create", Entity: "securitygroup", Api: "ec2", Input: "CreateSecurityGroupInput", ApiMethod: "CreateSecurityGroup", OutputExtractor: "GroupId",
		RequiredParams: map[string]string{
			"GroupName":   "name",
			"VpcId":       "vpc",
			"Description": "description",
		},
	},
	{
		Action: "update", Entity: "securitygroup", Api: "ec2", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"GroupId":    "id",
			"CidrIp":     "cidr",
			"IpProtocol": "protocol",
			// + either inbound or outbound = either authorize or revoke
		},
		//ExtraParams : portrange
	},
	{
		Action: "delete", Entity: "securitygroup", Api: "ec2", Input: "DeleteSecurityGroupInput", ApiMethod: "DeleteSecurityGroup",
		RequiredParams: map[string]string{
			"GroupId": "id",
		},
	},

	// VOLUME
	{
		Action: "create", Entity: "volume", Api: "ec2", Input: "CreateVolumeInput", ApiMethod: "CreateVolume", OutputExtractor: "VolumeId",
		RequiredParams: map[string]string{
			"AvailabilityZone": "zone",
			"Size":             "size",
		},
	},
	{
		Action: "delete", Entity: "volume", Api: "ec2", Input: "DeleteVolumeInput", ApiMethod: "DeleteVolume",
		RequiredParams: map[string]string{
			"VolumeId": "id",
		},
	},
	{
		Action: "attach", Entity: "volume", Api: "ec2", Input: "AttachVolumeInput", ApiMethod: "AttachVolume", OutputExtractor: "VolumeId",
		RequiredParams: map[string]string{
			"Device":     "device",
			"VolumeId":   "id",
			"InstanceId": "instance",
		},
	},
	// TAG
	{
		Action: "create", Entity: "tags", Api: "ec2", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"Resources": "resource",
		},
	},

	// Keypair
	{
		Action: "create", Entity: "keypair", Api: "ec2", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"KeyName": "name",
		},
	},
	{
		Action: "delete", Entity: "keypair", Api: "ec2", Input: "DeleteKeyPairInput", ApiMethod: "DeleteKeyPair",
		RequiredParams: map[string]string{
			"KeyName": "name",
		},
	},

	//// IAM

	// USER
	{
		Action: "create", Entity: "user", Api: "iam", DryRunUnsupported: true, Input: "CreateUserInput", ApiMethod: "CreateUser", OutputExtractor: "User.UserId",
		RequiredParams: map[string]string{
			"UserName": "name",
		},
	},
	{
		Action: "delete", Entity: "user", Api: "iam", DryRunUnsupported: true, Input: "DeleteUserInput", ApiMethod: "DeleteUser",
		RequiredParams: map[string]string{
			"UserName": "name",
		},
	},

	// GROUP
	{
		Action: "create", Entity: "group", Api: "iam", DryRunUnsupported: true, Input: "CreateGroupInput", ApiMethod: "CreateGroup", OutputExtractor: "Group.GroupId",
		RequiredParams: map[string]string{
			"GroupName": "name",
		},
	},
	{
		Action: "delete", Entity: "group", Api: "iam", DryRunUnsupported: true, Input: "DeleteGroupInput", ApiMethod: "DeleteGroup",
		RequiredParams: map[string]string{
			"GroupName": "name",
		},
	},
}
