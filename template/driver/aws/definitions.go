package aws

func DriverSupportedActions() map[string][]string {
	supported := make(map[string][]string)
	for _, def := range DriverDefinitions {
		supported[def.Action] = append(supported[def.Action], def.Entity)
	}
	return supported
}

var DriverDefinitions = []struct {
	ParamsMapping                     map[string]string
	TagsMapping                       map[string]string
	Action, Entity                    string
	Input, ApiMethod, OutputExtractor string
	ManualFuncDefinition              bool
}{
	// VPC
	{
		Action: "create", Entity: "vpc", Input: "CreateVpcInput", ApiMethod: "CreateVpc", OutputExtractor: "Vpc.VpcId",
		ParamsMapping: map[string]string{
			"CidrBlock": "cidr",
		},
	},
	{
		Action: "delete", Entity: "vpc", Input: "DeleteVpcInput", ApiMethod: "DeleteVpc",
		ParamsMapping: map[string]string{
			"VpcId": "id",
		},
	},

	// SUBNET
	{
		Action: "create", Entity: "subnet", Input: "CreateSubnetInput", ApiMethod: "CreateSubnet", OutputExtractor: "Subnet.SubnetId",
		ParamsMapping: map[string]string{
			"CidrBlock": "cidr",
			"VpcId":     "vpc",
		},
	},
	{
		Action: "delete", Entity: "subnet", Input: "DeleteSubnetInput", ApiMethod: "DeleteSubnet",
		ParamsMapping: map[string]string{
			"SubnetId": "id",
		},
	},

	// INSTANCES
	{
		Action: "create", Entity: "instance", Input: "RunInstancesInput", ApiMethod: "RunInstances", OutputExtractor: "Instances[0].InstanceId",
		ParamsMapping: map[string]string{
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
		ParamsMapping: map[string]string{
			"InstanceIds": "id",
		},
	},

	// TAG
	{
		Action: "create", Entity: "tags", ManualFuncDefinition: true,
		ParamsMapping: map[string]string{
			"Resources": "resource",
		},
	},
}
