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
		TagsMapping: map[string]string{
			"Name": "name",
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
	// TAG
	{
		Action: "create", Entity: "tags", Api: "ec2", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"Resources": "resource",
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
