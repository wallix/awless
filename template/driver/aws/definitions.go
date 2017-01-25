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
	Action, Entity                    string
	Input, ApiMethod, OutputExtractor string
	ManualFuncDefinition              bool
}{
	// VPC
	{
		Action: "create", Entity: "vpc", Input: "CreateVpcInput", ApiMethod: "CreateVpc", OutputExtractor: "Vpc.VpcId",
		RequiredParams: map[string]string{
			"CidrBlock": "cidr",
		},
	},
	{
		Action: "delete", Entity: "vpc", Input: "DeleteVpcInput", ApiMethod: "DeleteVpc",
		RequiredParams: map[string]string{
			"VpcId": "id",
		},
	},

	// SUBNET
	{
		Action: "create", Entity: "subnet", Input: "CreateSubnetInput", ApiMethod: "CreateSubnet", OutputExtractor: "Subnet.SubnetId",
		RequiredParams: map[string]string{
			"CidrBlock": "cidr",
			"VpcId":     "vpc",
		},
	},
	{
		Action: "delete", Entity: "subnet", Input: "DeleteSubnetInput", ApiMethod: "DeleteSubnet",
		RequiredParams: map[string]string{
			"SubnetId": "id",
		},
	},

	// INSTANCES
	{
		Action: "create", Entity: "instance", Input: "RunInstancesInput", ApiMethod: "RunInstances", OutputExtractor: "Instances[0].InstanceId",
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
		Action: "delete", Entity: "instance", Input: "TerminateInstancesInput", ApiMethod: "TerminateInstances",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},
	{
		Action: "start", Entity: "instance", Input: "StartInstancesInput", ApiMethod: "StartInstances", OutputExtractor: "StartingInstances[0].InstanceId",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},
	{
		Action: "stop", Entity: "instance", Input: "StopInstancesInput", ApiMethod: "StopInstances", OutputExtractor: "StoppingInstances[0].InstanceId",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},

	// TAG
	{
		Action: "create", Entity: "tags", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"Resources": "resource",
		},
	},
}
