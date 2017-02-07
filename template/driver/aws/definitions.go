package aws

import "github.com/wallix/awless/graph"

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
		Action: "create", Entity: graph.Vpc.String(), Api: "ec2", Input: "CreateVpcInput", ApiMethod: "CreateVpc", OutputExtractor: "aws.StringValue(output.Vpc.VpcId)",
		RequiredParams: map[string]string{
			"CidrBlock": "cidr",
		},
	},
	{
		Action: "delete", Entity: graph.Vpc.String(), Api: "ec2", Input: "DeleteVpcInput", ApiMethod: "DeleteVpc",
		RequiredParams: map[string]string{
			"VpcId": "id",
		},
	},

	// SUBNET
	{
		Action: "create", Entity: graph.Subnet.String(), Api: "ec2", Input: "CreateSubnetInput", ApiMethod: "CreateSubnet", OutputExtractor: "aws.StringValue(output.Subnet.SubnetId)",
		RequiredParams: map[string]string{
			"CidrBlock": "cidr",
			"VpcId":     "vpc",
		},
		ExtraParams: map[string]string{
			"AvailabilityZone": "zone",
		},
	},
	{
		Action: "update", Entity: graph.Subnet.String(), Api: "ec2", Input: "ModifySubnetAttributeInput", ApiMethod: "ModifySubnetAttribute", DryRunUnsupported: true,
		RequiredParams: map[string]string{
			"SubnetId": "id",
		},
		ExtraParams: map[string]string{
			"MapPublicIpOnLaunch": "public-vms",
		},
	},
	{
		Action: "delete", Entity: graph.Subnet.String(), Api: "ec2", Input: "DeleteSubnetInput", ApiMethod: "DeleteSubnet",
		RequiredParams: map[string]string{
			"SubnetId": "id",
		},
	},

	// INSTANCES
	{
		Action: "create", Entity: graph.Instance.String(), Api: "ec2", Input: "RunInstancesInput", ApiMethod: "RunInstances", OutputExtractor: "aws.StringValue(output.Instances[0].InstanceId)",
		RequiredParams: map[string]string{
			"ImageId":      "image",
			"MaxCount":     "count",
			"MinCount":     "count",
			"InstanceType": "type",
			"SubnetId":     "subnet",
		},
		ExtraParams: map[string]string{
			"KeyName":               "key",
			"PrivateIpAddress":      "ip",
			"UserData":              "userdata",
			"SecurityGroupIds":      "group",
			"DisableApiTermination": "lock",
		},
		TagsMapping: map[string]string{
			"Name": "name",
		},
	},
	{
		Action: "update", Entity: graph.Instance.String(), Api: "ec2", Input: "ModifyInstanceAttributeInput", ApiMethod: "ModifyInstanceAttribute",
		RequiredParams: map[string]string{
			"InstanceId": "id",
		},
		ExtraParams: map[string]string{
			"InstanceType":          "type",
			"Groups":                "group",
			"DisableApiTermination": "lock",
		},
	},
	{
		Action: "delete", Entity: graph.Instance.String(), Api: "ec2", Input: "TerminateInstancesInput", ApiMethod: "TerminateInstances",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},
	{
		Action: "start", Entity: graph.Instance.String(), Api: "ec2", Input: "StartInstancesInput", ApiMethod: "StartInstances", OutputExtractor: "aws.StringValue(output.StartingInstances[0].InstanceId)",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},
	{
		Action: "stop", Entity: graph.Instance.String(), Api: "ec2", Input: "StopInstancesInput", ApiMethod: "StopInstances", OutputExtractor: "aws.StringValue(output.StoppingInstances[0].InstanceId)",
		RequiredParams: map[string]string{
			"InstanceIds": "id",
		},
	},
	{
		Action: "check", Entity: graph.Instance.String(), Api: "ec2", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"Id":      "id",
			"State":   "state",
			"Timeout": "timeout",
		},
	},

	// Security Group
	{
		Action: "create", Entity: graph.SecurityGroup.String(), Api: "ec2", Input: "CreateSecurityGroupInput", ApiMethod: "CreateSecurityGroup", OutputExtractor: "aws.StringValue(output.GroupId)",
		RequiredParams: map[string]string{
			"GroupName":   "name",
			"VpcId":       "vpc",
			"Description": "description",
		},
	},
	{
		Action: "update", Entity: graph.SecurityGroup.String(), Api: "ec2", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"GroupId":    "id",
			"CidrIp":     "cidr",
			"IpProtocol": "protocol",
			// + either inbound or outbound = either authorize or revoke
		},
		//ExtraParams : portrange
	},
	{
		Action: "delete", Entity: graph.SecurityGroup.String(), Api: "ec2", Input: "DeleteSecurityGroupInput", ApiMethod: "DeleteSecurityGroup",
		RequiredParams: map[string]string{
			"GroupId": "id",
		},
	},

	// VOLUME
	{
		Action: "create", Entity: graph.Volume.String(), Api: "ec2", Input: "CreateVolumeInput", ApiMethod: "CreateVolume", OutputExtractor: "aws.StringValue(output.VolumeId)",
		RequiredParams: map[string]string{
			"AvailabilityZone": "zone",
			"Size":             "size",
		},
	},
	{
		Action: "delete", Entity: graph.Volume.String(), Api: "ec2", Input: "DeleteVolumeInput", ApiMethod: "DeleteVolume",
		RequiredParams: map[string]string{
			"VolumeId": "id",
		},
	},
	{
		Action: "attach", Entity: graph.Volume.String(), Api: "ec2", Input: "AttachVolumeInput", ApiMethod: "AttachVolume", OutputExtractor: "aws.StringValue(output.VolumeId)",
		RequiredParams: map[string]string{
			"Device":     "device",
			"VolumeId":   "id",
			"InstanceId": "instance",
		},
	},
	// INTERNET GATEWAYS
	{
		Action: "create", Entity: graph.InternetGateway.String(), Api: "ec2", Input: "CreateInternetGatewayInput", ApiMethod: "CreateInternetGateway", OutputExtractor: "aws.StringValue(output.InternetGateway.InternetGatewayId)",
		RequiredParams: map[string]string{},
	},
	{
		Action: "delete", Entity: graph.InternetGateway.String(), Api: "ec2", Input: "DeleteInternetGatewayInput", ApiMethod: "DeleteInternetGateway",
		RequiredParams: map[string]string{
			"InternetGatewayId": "id",
		},
	},
	{
		Action: "attach", Entity: graph.InternetGateway.String(), Api: "ec2", Input: "AttachInternetGatewayInput", ApiMethod: "AttachInternetGateway",
		RequiredParams: map[string]string{
			"InternetGatewayId": "id",
			"VpcId":             "vpc",
		},
	},
	{
		Action: "detach", Entity: graph.InternetGateway.String(), Api: "ec2", Input: "DetachInternetGatewayInput", ApiMethod: "DetachInternetGateway",
		RequiredParams: map[string]string{
			"InternetGatewayId": "id",
			"VpcId":             "vpc",
		},
	},
	// ROUTE TABLES
	{
		Action: "create", Entity: graph.RouteTable.String(), Api: "ec2", Input: "CreateRouteTableInput", ApiMethod: "CreateRouteTable", OutputExtractor: "aws.StringValue(output.RouteTable.RouteTableId)",
		RequiredParams: map[string]string{
			"VpcId": "vpc"},
	},
	{
		Action: "delete", Entity: graph.RouteTable.String(), Api: "ec2", Input: "DeleteRouteTableInput", ApiMethod: "DeleteRouteTable",
		RequiredParams: map[string]string{
			"RouteTableId": "id",
		},
	},
	{
		Action: "attach", Entity: graph.RouteTable.String(), Api: "ec2", Input: "AssociateRouteTableInput", ApiMethod: "AssociateRouteTable", OutputExtractor: "aws.StringValue(output.AssociationId)",
		RequiredParams: map[string]string{
			"RouteTableId": "id",
			"SubnetId":     "subnet",
		},
	},
	{
		Action: "detach", Entity: graph.RouteTable.String(), Api: "ec2", Input: "DisassociateRouteTableInput", ApiMethod: "DisassociateRouteTable",
		RequiredParams: map[string]string{
			"AssociationId": "association",
		},
	},
	// ROUTES
	{
		Action: "create", Entity: "route", Api: "ec2", Input: "CreateRouteInput", ApiMethod: "CreateRoute",
		RequiredParams: map[string]string{
			"RouteTableId":         "table",
			"DestinationCidrBlock": "cidr",
			"GatewayId":            "gateway",
		},
	},
	{
		Action: "delete", Entity: "route", Api: "ec2", Input: "DeleteRouteInput", ApiMethod: "DeleteRoute",
		RequiredParams: map[string]string{
			"RouteTableId":         "table",
			"DestinationCidrBlock": "cidr",
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
		Action: "create", Entity: graph.Keypair.String(), Api: "ec2", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"KeyName": "name",
		},
	},
	{
		Action: "delete", Entity: graph.Keypair.String(), Api: "ec2", Input: "DeleteKeyPairInput", ApiMethod: "DeleteKeyPair",
		RequiredParams: map[string]string{
			"KeyName": "id",
		},
	},

	//// IAM

	// USER
	{
		Action: "create", Entity: graph.User.String(), Api: "iam", DryRunUnsupported: true, Input: "CreateUserInput", ApiMethod: "CreateUser", OutputExtractor: "aws.StringValue(output.User.UserId)",
		RequiredParams: map[string]string{
			"UserName": "name",
		},
	},
	{
		Action: "delete", Entity: graph.User.String(), Api: "iam", DryRunUnsupported: true, Input: "DeleteUserInput", ApiMethod: "DeleteUser",
		RequiredParams: map[string]string{
			"UserName": "name",
		},
	},

	// GROUP
	{
		Action: "create", Entity: graph.Group.String(), Api: "iam", DryRunUnsupported: true, Input: "CreateGroupInput", ApiMethod: "CreateGroup", OutputExtractor: "aws.StringValue(output.Group.GroupId)",
		RequiredParams: map[string]string{
			"GroupName": "name",
		},
	},
	{
		Action: "delete", Entity: graph.Group.String(), Api: "iam", DryRunUnsupported: true, Input: "DeleteGroupInput", ApiMethod: "DeleteGroup",
		RequiredParams: map[string]string{
			"GroupName": "name",
		},
	},

	// POLICY
	{
		Action: "attach", Entity: graph.Policy.String(), Api: "iam", DryRunUnsupported: true, Input: "AttachUserPolicyInput", ApiMethod: "AttachUserPolicy",
		RequiredParams: map[string]string{
			"PolicyArn": "arn",
			"UserName":  "user",
		},
	},
	{
		Action: "detach", Entity: "policy", Api: "iam", DryRunUnsupported: true, Input: "DetachUserPolicyInput", ApiMethod: "DetachUserPolicy",
		RequiredParams: map[string]string{
			"PolicyArn": "arn",
			"UserName":  "user",
		},
	},

	//// S3

	// BUCKET
	{
		Action: "create", Entity: graph.Bucket.String(), Api: "s3", DryRunUnsupported: true, Input: "CreateBucketInput", ApiMethod: "CreateBucket", OutputExtractor: "params[\"name\"]",
		RequiredParams: map[string]string{
			"Bucket": "name",
		},
	},
	{
		Action: "delete", Entity: graph.Bucket.String(), Api: "s3", DryRunUnsupported: true, Input: "DeleteBucketInput", ApiMethod: "DeleteBucket",
		RequiredParams: map[string]string{
			"Bucket": "name",
		},
	},

	// OBJECT
	{
		Action: "create", Entity: graph.Object.String(), Api: "s3", ManualFuncDefinition: true,
		RequiredParams: map[string]string{
			"Bucket": "bucket",
			"Body":   "file",
		},
		ExtraParams: map[string]string{
			"Key": "name",
		},
	},
	{
		Action: "delete", Entity: graph.Object.String(), Api: "s3", DryRunUnsupported: true, Input: "DeleteObjectInput", ApiMethod: "DeleteObject",
		RequiredParams: map[string]string{
			"Bucket": "bucket",
			"Key":    "key",
		},
	},
}
