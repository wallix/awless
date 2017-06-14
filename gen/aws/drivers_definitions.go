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
	"sort"

	"github.com/wallix/awless/cloud"
)

type param struct {
	AwsField, AwsType string
	TemplateName      string
	AsAwsTag          bool
}

type driver struct {
	RequiredParams                            []param
	ExtraParams                               []param
	Action, Entity                            string
	Input, Output, ApiMethod, OutputExtractor string
	DryRunUnsupported                         bool
	ManualFuncDefinition                      bool
}

func (d *driver) RequiredKeys() []string {
	var keys []string
	for _, p := range d.RequiredParams {
		keys = append(keys, p.TemplateName)
	}

	return sortUnique(keys)
}

func (d *driver) ExtraKeys() []string {
	var keys []string
	for _, p := range d.ExtraParams {
		keys = append(keys, p.TemplateName)
	}

	return sortUnique(keys)
}

type driversDef struct {
	Api          string
	ApiInterface string
	Drivers      []driver
}

func sortUnique(arr []string) (sorted []string) {
	unique := make(map[string]struct{})
	for _, s := range arr {
		unique[s] = struct{}{}
	}

	for k := range unique {
		sorted = append(sorted, k)
	}

	sort.Strings(sorted)
	return
}

var DriversDefs = []driversDef{
	{
		Api: "ec2",
		Drivers: []driver{
			// VPC
			{
				Action: "create", Entity: cloud.Vpc, ApiMethod: "CreateVpc", Input: "CreateVpcInput", Output: "CreateVpcOutput", OutputExtractor: "aws.StringValue(output.Vpc.VpcId)",
				RequiredParams: []param{
					{AwsField: "CidrBlock", TemplateName: "cidr", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Name", TemplateName: "name", AsAwsTag: true},
				},
			},
			{
				Action: "delete", Entity: cloud.Vpc, ApiMethod: "DeleteVpc", Input: "DeleteVpcInput", Output: "DeleteVpcOutput",
				RequiredParams: []param{
					{AwsField: "VpcId", TemplateName: "id", AwsType: "awsstr"},
				},
			},

			// SUBNET
			{
				Action: "create", Entity: cloud.Subnet, ApiMethod: "CreateSubnet", Input: "CreateSubnetInput", Output: "CreateSubnetOutput", OutputExtractor: "aws.StringValue(output.Subnet.SubnetId)",
				RequiredParams: []param{
					{AwsField: "CidrBlock", TemplateName: "cidr", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "AvailabilityZone", TemplateName: "availabilityzone", AwsType: "awsstr"},
					{AwsField: "Name", TemplateName: "name", AsAwsTag: true},
				},
			},
			{
				Action: "update", Entity: cloud.Subnet, ApiMethod: "ModifySubnetAttribute", Input: "ModifySubnetAttributeInput", Output: "ModifySubnetAttributeOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "SubnetId", TemplateName: "id", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "MapPublicIpOnLaunch", TemplateName: "public", AwsType: "awsboolattribute"},
				},
			},
			{
				Action: "delete", Entity: cloud.Subnet, ApiMethod: "DeleteSubnet", Input: "DeleteSubnetInput", Output: "DeleteSubnetOutput",
				RequiredParams: []param{
					{AwsField: "SubnetId", TemplateName: "id", AwsType: "awsstr"},
				},
			},

			// INSTANCES
			{
				Action: "create", Entity: cloud.Instance, ApiMethod: "RunInstances", Input: "RunInstancesInput", Output: "Reservation", OutputExtractor: "aws.StringValue(output.Instances[0].InstanceId)",
				RequiredParams: []param{
					{AwsField: "ImageId", TemplateName: "image", AwsType: "awsstr"},
					{AwsField: "MaxCount", TemplateName: "count", AwsType: "awsint64"},
					{AwsField: "MinCount", TemplateName: "count", AwsType: "awsint64"},
					{AwsField: "InstanceType", TemplateName: "type", AwsType: "awsstr"},
					{AwsField: "SubnetId", TemplateName: "subnet", AwsType: "awsstr"},
					{AwsField: "Name", TemplateName: "name", AsAwsTag: true},
				},
				ExtraParams: []param{
					{AwsField: "KeyName", TemplateName: "keypair", AwsType: "awsstr"},
					{AwsField: "PrivateIpAddress", TemplateName: "ip", AwsType: "awsstr"},
					{AwsField: "UserData", TemplateName: "userdata", AwsType: "awsfiletobase64"},
					{AwsField: "SecurityGroupIds", TemplateName: "securitygroup", AwsType: "awsstringslice"},
					{AwsField: "DisableApiTermination", TemplateName: "lock", AwsType: "awsbool"},
					{AwsField: "IamInstanceProfile.Name", TemplateName: "role", AwsType: "awsstr"},
				},
			},
			{
				Action: "update", Entity: cloud.Instance, ApiMethod: "ModifyInstanceAttribute", Input: "ModifyInstanceAttributeInput", Output: "ModifyInstanceAttributeOutput",
				RequiredParams: []param{
					{AwsField: "InstanceId", TemplateName: "id", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "InstanceType.Value", TemplateName: "type", AwsType: "awsstr"},
					{AwsField: "DisableApiTermination", TemplateName: "lock", AwsType: "awsboolattribute"},
				},
			},
			{
				Action: "delete", Entity: cloud.Instance, ApiMethod: "TerminateInstances", Input: "TerminateInstancesInput", Output: "TerminateInstancesOutput",
				RequiredParams: []param{
					{AwsField: "InstanceIds", TemplateName: "id", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "start", Entity: cloud.Instance, ApiMethod: "StartInstances", Input: "StartInstancesInput", Output: "StartInstancesOutput", OutputExtractor: "aws.StringValue(output.StartingInstances[0].InstanceId)",
				RequiredParams: []param{
					{AwsField: "InstanceIds", TemplateName: "id", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "stop", Entity: cloud.Instance, ApiMethod: "StopInstances", Input: "StopInstancesInput", Output: "StopInstancesOutput", OutputExtractor: "aws.StringValue(output.StoppingInstances[0].InstanceId)",
				RequiredParams: []param{
					{AwsField: "InstanceIds", TemplateName: "id", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "check", Entity: cloud.Instance, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "state"},
					{TemplateName: "timeout"},
				},
			},
			// Security Group
			{
				Action: "create", Entity: cloud.SecurityGroup, ApiMethod: "CreateSecurityGroup", Input: "CreateSecurityGroupInput", Output: "CreateSecurityGroupOutput", OutputExtractor: "aws.StringValue(output.GroupId)",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
					{AwsField: "Description", TemplateName: "description", AwsType: "awsstr"},
				},
			},
			{
				Action: "update", Entity: cloud.SecurityGroup, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "cidr"},
					{TemplateName: "protocol"},
				},
				ExtraParams: []param{
					{TemplateName: "inbound"}, // either inbound or outbound = either authorize or revoke
					{TemplateName: "outbound"},
					{TemplateName: "portrange"},
				},
			},
			{
				Action: "delete", Entity: cloud.SecurityGroup, ApiMethod: "DeleteSecurityGroup", Input: "DeleteSecurityGroupInput", Output: "DeleteSecurityGroupOutput",
				RequiredParams: []param{
					{AwsField: "GroupId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "check", Entity: cloud.SecurityGroup, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "state"},
					{TemplateName: "timeout"},
				},
			},
			{
				Action: "attach", Entity: cloud.SecurityGroup, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
				},
				ExtraParams: []param{
					{TemplateName: "instance"},
				},
			},
			{
				Action: "detach", Entity: cloud.SecurityGroup, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
				},
				ExtraParams: []param{
					{TemplateName: "instance"},
				},
			},
			// IMAGES
			{
				Action: "copy", Entity: cloud.Image, ApiMethod: "CopyImage", Input: "CopyImageInput", Output: "CopyImageOutput", OutputExtractor: "aws.StringValue(output.ImageId)",
				RequiredParams: []param{
					{AwsField: "Name", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "SourceImageId", TemplateName: "source-id", AwsType: "awsstr"},
					{AwsField: "SourceRegion", TemplateName: "source-region", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Encrypted", TemplateName: "encrypted", AwsType: "awsbool"},
					{AwsField: "Description", TemplateName: "description", AwsType: "awsstr"},
				},
			},
			{
				Action: "import", Entity: cloud.Image, ApiMethod: "ImportImage", Input: "ImportImageInput", Output: "ImportImageOutput", OutputExtractor: "aws.StringValue(output.ImportTaskId)",
				RequiredParams: []param{},
				ExtraParams: []param{
					{AwsField: "Architecture", TemplateName: "architecture", AwsType: "awsstr"},
					{AwsField: "Description", TemplateName: "description", AwsType: "awsstr"},
					{AwsField: "LicenseType", TemplateName: "license", AwsType: "awsstr"},
					{AwsField: "Platform", TemplateName: "platform", AwsType: "awsstr"},
					{AwsField: "RoleName", TemplateName: "role", AwsType: "awsstr"},
					{AwsField: "DiskContainers[0]SnapshotId", TemplateName: "snapshot", AwsType: "awsslicestruct"},
					{AwsField: "DiskContainers[0]Url", TemplateName: "url", AwsType: "awsslicestruct"},
					{AwsField: "DiskContainers[0]UserBucket.S3Bucket", TemplateName: "bucket", AwsType: "awsslicestruct"},
					{AwsField: "DiskContainers[0]UserBucket.S3Key", TemplateName: "s3object", AwsType: "awsslicestruct"},
				},
			},
			{
				Action: "delete", Entity: cloud.Image, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "delete-snapshots"},
				},
			},

			// VOLUME
			{
				Action: "create", Entity: cloud.Volume, ApiMethod: "CreateVolume", Input: "CreateVolumeInput", Output: "Volume", OutputExtractor: "aws.StringValue(output.VolumeId)",
				RequiredParams: []param{
					{AwsField: "AvailabilityZone", TemplateName: "availabilityzone", AwsType: "awsstr"},
					{AwsField: "Size", TemplateName: "size", AwsType: "awsint64"},
				},
			},
			{
				Action: "check", Entity: cloud.Volume, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "state"},
					{TemplateName: "timeout"},
				},
			},
			{
				Action: "delete", Entity: cloud.Volume, ApiMethod: "DeleteVolume", Input: "DeleteVolumeInput", Output: "DeleteVolumeOutput",
				RequiredParams: []param{
					{AwsField: "VolumeId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.Volume, ApiMethod: "AttachVolume", Input: "AttachVolumeInput", Output: "VolumeAttachment", OutputExtractor: "aws.StringValue(output.VolumeId)",
				RequiredParams: []param{
					{AwsField: "Device", TemplateName: "device", AwsType: "awsstr"},
					{AwsField: "VolumeId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "InstanceId", TemplateName: "instance", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: cloud.Volume, ApiMethod: "DetachVolume", Input: "DetachVolumeInput", Output: "VolumeAttachment", OutputExtractor: "aws.StringValue(output.VolumeId)",
				RequiredParams: []param{
					{AwsField: "Device", TemplateName: "device", AwsType: "awsstr"},
					{AwsField: "VolumeId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "InstanceId", TemplateName: "instance", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Force", TemplateName: "force", AwsType: "awsbool"},
				},
			},
			// Snapshot
			{
				Action: "create", Entity: cloud.Snapshot, ApiMethod: "CreateSnapshot", Input: "CreateSnapshotInput", Output: "Snapshot", OutputExtractor: "aws.StringValue(output.SnapshotId)",
				RequiredParams: []param{
					{AwsField: "VolumeId", TemplateName: "volume", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Description", TemplateName: "description", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Snapshot, ApiMethod: "DeleteSnapshot", Input: "DeleteSnapshotInput", Output: "DeleteSnapshotOutput",
				RequiredParams: []param{
					{AwsField: "SnapshotId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "copy", Entity: cloud.Snapshot, ApiMethod: "CopySnapshot", Input: "CopySnapshotInput", Output: "CopySnapshotOutput", OutputExtractor: "aws.StringValue(output.SnapshotId)",
				RequiredParams: []param{
					{AwsField: "SourceSnapshotId", TemplateName: "source-id", AwsType: "awsstr"},
					{AwsField: "SourceRegion", TemplateName: "source-region", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Encrypted", TemplateName: "encrypted", AwsType: "awsbool"},
					{AwsField: "Description", TemplateName: "description", AwsType: "awsstr"},
				},
			},
			// INTERNET GATEWAYS
			{
				Action: "create", Entity: cloud.InternetGateway, ApiMethod: "CreateInternetGateway", Input: "CreateInternetGatewayInput", Output: "CreateInternetGatewayOutput", OutputExtractor: "aws.StringValue(output.InternetGateway.InternetGatewayId)",
			},
			{
				Action: "delete", Entity: cloud.InternetGateway, ApiMethod: "DeleteInternetGateway", Input: "DeleteInternetGatewayInput", Output: "DeleteInternetGatewayOutput",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.InternetGateway, ApiMethod: "AttachInternetGateway", Input: "AttachInternetGatewayInput", Output: "AttachInternetGatewayOutput",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: cloud.InternetGateway, ApiMethod: "DetachInternetGateway", Input: "DetachInternetGatewayInput", Output: "DetachInternetGatewayOutput",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
			},
			// ROUTE TABLES
			{
				Action: "create", Entity: cloud.RouteTable, ApiMethod: "CreateRouteTable", Input: "CreateRouteTableInput", Output: "CreateRouteTableOutput", OutputExtractor: "aws.StringValue(output.RouteTable.RouteTableId)",
				RequiredParams: []param{
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"}},
			},
			{
				Action: "delete", Entity: cloud.RouteTable, ApiMethod: "DeleteRouteTable", Input: "DeleteRouteTableInput", Output: "DeleteRouteTableOutput",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.RouteTable, ApiMethod: "AssociateRouteTable", Input: "AssociateRouteTableInput", Output: "AssociateRouteTableOutput", OutputExtractor: "aws.StringValue(output.AssociationId)",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "SubnetId", TemplateName: "subnet", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: cloud.RouteTable, ApiMethod: "DisassociateRouteTable", Input: "DisassociateRouteTableInput", Output: "DisassociateRouteTableOutput",
				RequiredParams: []param{
					{AwsField: "AssociationId", TemplateName: "association", AwsType: "awsstr"},
				},
			},
			// ROUTES
			{
				Action: "create", Entity: "route", ApiMethod: "CreateRoute", Input: "CreateRouteInput", Output: "CreateRouteOutput",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "table", AwsType: "awsstr"},
					{AwsField: "DestinationCidrBlock", TemplateName: "cidr", AwsType: "awsstr"},
					{AwsField: "GatewayId", TemplateName: "gateway", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: "route", ApiMethod: "DeleteRoute", Input: "DeleteRouteInput", Output: "DeleteRouteOutput",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "table", AwsType: "awsstr"},
					{AwsField: "DestinationCidrBlock", TemplateName: "cidr", AwsType: "awsstr"},
				},
			},
			// TAG
			{
				Action: "create", Entity: "tag", ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "resource"},
					{TemplateName: "key"},
					{TemplateName: "value"},
				},
			},
			{
				Action: "delete", Entity: "tag", ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "resource"},
					{TemplateName: "key"},
					{TemplateName: "value"},
				},
			},

			// Key
			{
				Action: "create", Entity: cloud.Keypair, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
				},
				ExtraParams: []param{
					{TemplateName: "encrypted"},
				},
			},
			{
				Action: "delete", Entity: cloud.Keypair, ApiMethod: "DeleteKeyPair", Input: "DeleteKeyPairInput", Output: "DeleteKeyPairOutput",
				RequiredParams: []param{
					{AwsField: "KeyName", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// Allocate address
			{
				Action: "create", Entity: cloud.ElasticIP, ApiMethod: "AllocateAddress", Input: "AllocateAddressInput", Output: "AllocateAddressOutput", OutputExtractor: "aws.StringValue(output.AllocationId)", // should return PublicIp if params["domain"] == "standard"
				RequiredParams: []param{
					{AwsField: "Domain", TemplateName: "domain", AwsType: "awsstr"},
				},
				ExtraParams: []param{},
			},
			{
				Action: "delete", Entity: cloud.ElasticIP, ApiMethod: "ReleaseAddress", Input: "ReleaseAddressInput", Output: "ReleaseAddressOutput",
				ExtraParams: []param{
					{AwsField: "AllocationId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "PublicIp", TemplateName: "ip", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.ElasticIP, ApiMethod: "AssociateAddress", Input: "AssociateAddressInput", Output: "AssociateAddressOutput", OutputExtractor: "aws.StringValue(output.AssociationId)",
				RequiredParams: []param{
					{AwsField: "AllocationId", TemplateName: "id", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "InstanceId", TemplateName: "instance", AwsType: "awsstr"},
					{AwsField: "NetworkInterfaceId", TemplateName: "networkinterface", AwsType: "awsstr"},
					{AwsField: "PrivateIpAddress", TemplateName: "privateip", AwsType: "awsstr"},
					{AwsField: "AllowReassociation", TemplateName: "allow-reassociation", AwsType: "awsbool"},
				},
			},
			{
				Action: "detach", Entity: cloud.ElasticIP, ApiMethod: "DisassociateAddress", Input: "DisassociateAddressInput", Output: "DisassociateAddressOutput",
				RequiredParams: []param{
					{AwsField: "AssociationId", TemplateName: "association", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api: "elbv2",
		Drivers: []driver{
			// LoadBalancer
			{
				Action: "create", Entity: cloud.LoadBalancer, ApiMethod: "CreateLoadBalancer", Input: "CreateLoadBalancerInput", Output: "CreateLoadBalancerOutput", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.LoadBalancers[0].LoadBalancerArn)",
				RequiredParams: []param{
					{AwsField: "Name", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "Subnets", TemplateName: "subnets", AwsType: "awsstringslice"},
				},
				ExtraParams: []param{
					{AwsField: "IpAddressType", TemplateName: "iptype", AwsType: "awsstr"},
					{AwsField: "Scheme", TemplateName: "scheme", AwsType: "awsstr"},
					{AwsField: "SecurityGroups", TemplateName: "securitygroups", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "delete", Entity: cloud.LoadBalancer, ApiMethod: "DeleteLoadBalancer", Input: "DeleteLoadBalancerInput", Output: "DeleteLoadBalancerOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "LoadBalancerArn", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "check", Entity: cloud.LoadBalancer, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "state"},
					{TemplateName: "timeout"},
				},
			},
			// Listener
			{
				Action: "create", Entity: cloud.Listener, ApiMethod: "CreateListener", Input: "CreateListenerInput", Output: "CreateListenerOutput", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.Listeners[0].ListenerArn)",
				RequiredParams: []param{
					{AwsField: "DefaultActions[0]Type", TemplateName: "actiontype", AwsType: "awsslicestruct"}, //always forward
					{AwsField: "DefaultActions[0]TargetGroupArn", TemplateName: "targetgroup", AwsType: "awsslicestruct"},
					{AwsField: "LoadBalancerArn", TemplateName: "loadbalancer", AwsType: "awsstr"},
					{AwsField: "Port", TemplateName: "port", AwsType: "awsint64"},
					{AwsField: "Protocol", TemplateName: "protocol", AwsType: "awsstr"}, // TCP, HTTP, HTTPS
				},
				ExtraParams: []param{
					{AwsField: "Certificates[0]CertificateArn", TemplateName: "certificate", AwsType: "awsslicestruct"},
					{AwsField: "SslPolicy", TemplateName: "sslpolicy", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Listener, ApiMethod: "DeleteListener", Input: "DeleteListenerInput", Output: "DeleteListenerOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "ListenerArn", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			// Target group
			{
				Action: "create", Entity: cloud.TargetGroup, ApiMethod: "CreateTargetGroup", Input: "CreateTargetGroupInput", Output: "CreateTargetGroupOutput", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.TargetGroups[0].TargetGroupArn)",
				RequiredParams: []param{
					{AwsField: "Name", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "Port", TemplateName: "port", AwsType: "awsint64"},
					{AwsField: "Protocol", TemplateName: "protocol", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "HealthCheckIntervalSeconds", TemplateName: "healthcheckinterval", AwsType: "awsint64"},
					{AwsField: "HealthCheckPath", TemplateName: "healthcheckpath", AwsType: "awsstr"},
					{AwsField: "HealthCheckPort", TemplateName: "healthcheckport", AwsType: "awsstr"},
					{AwsField: "HealthCheckProtocol", TemplateName: "healthcheckprotocol", AwsType: "awsstr"},
					{AwsField: "HealthCheckTimeoutSeconds", TemplateName: "healthchecktimeout", AwsType: "awsint64"},
					{AwsField: "HealthyThresholdCount", TemplateName: "healthythreshold", AwsType: "awsint64"},
					{AwsField: "UnhealthyThresholdCount", TemplateName: "unhealthythreshold", AwsType: "awsint64"},
					{AwsField: "Matcher.HttpCode", TemplateName: "matcher", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.TargetGroup, ApiMethod: "DeleteTargetGroup", Input: "DeleteTargetGroupInput", Output: "DeleteTargetGroupOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "TargetGroupArn", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.Instance, ApiMethod: "RegisterTargets", Input: "RegisterTargetsInput", Output: "RegisterTargetsOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "TargetGroupArn", TemplateName: "targetgroup", AwsType: "awsstr"},
					{AwsField: "Targets[0]Id", TemplateName: "id", AwsType: "awsslicestruct"},
				},
				ExtraParams: []param{
					{AwsField: "Targets[0]Port", TemplateName: "port", AwsType: "awsslicestructint64"},
				},
			},
			{
				Action: "detach", Entity: cloud.Instance, ApiMethod: "DeregisterTargets", Input: "DeregisterTargetsInput", Output: "DeregisterTargetsOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "TargetGroupArn", TemplateName: "targetgroup", AwsType: "awsstr"},
					{AwsField: "Targets[0]Id", TemplateName: "id", AwsType: "awsslicestruct"},
				},
			},
		},
	},
	{
		Api:          "autoscaling",
		ApiInterface: "AutoScalingAPI",
		Drivers: []driver{
			{
				Action: "create", Entity: cloud.LaunchConfiguration, ApiMethod: "CreateLaunchConfiguration", Input: "CreateLaunchConfigurationInput", Output: "CreateLaunchConfigurationOutput", DryRunUnsupported: true, OutputExtractor: "params[\"name\"]",
				RequiredParams: []param{
					{AwsField: "ImageId", TemplateName: "image", AwsType: "awsstr"},
					{AwsField: "InstanceType", TemplateName: "type", AwsType: "awsstr"},
					{AwsField: "LaunchConfigurationName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "AssociatePublicIpAddress", TemplateName: "public", AwsType: "awsbool"},
					{AwsField: "KeyName", TemplateName: "keypair", AwsType: "awsstr"},
					{AwsField: "UserData", TemplateName: "userdata", AwsType: "awsfiletobase64"},
					{AwsField: "SecurityGroups", TemplateName: "securitygroups", AwsType: "awsstringslice"},
					{AwsField: "IamInstanceProfile", TemplateName: "role", AwsType: "awsstr"},
					{AwsField: "SpotPrice", TemplateName: "spotprice", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.LaunchConfiguration, ApiMethod: "DeleteLaunchConfiguration", Input: "DeleteLaunchConfigurationInput", Output: "DeleteLaunchConfigurationOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "LaunchConfigurationName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "create", Entity: cloud.ScalingGroup, ApiMethod: "CreateAutoScalingGroup", Input: "CreateAutoScalingGroupInput", Output: "CreateAutoScalingGroupOutput", DryRunUnsupported: true, OutputExtractor: "params[\"name\"]",
				RequiredParams: []param{
					{AwsField: "AutoScalingGroupName", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "LaunchConfigurationName", TemplateName: "launchconfiguration", AwsType: "awsstr"},
					{AwsField: "MaxSize", TemplateName: "max-size", AwsType: "awsint64"},
					{AwsField: "MinSize", TemplateName: "min-size", AwsType: "awsint64"},
					{AwsField: "VPCZoneIdentifier", TemplateName: "subnets", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "DefaultCooldown", TemplateName: "cooldown", AwsType: "awsint64"},
					{AwsField: "DesiredCapacity", TemplateName: "desired-capacity", AwsType: "awsint64"},
					{AwsField: "HealthCheckGracePeriod", TemplateName: "healthcheck-grace-period", AwsType: "awsint64"},
					{AwsField: "HealthCheckType", TemplateName: "healthcheck-type", AwsType: "awsstr"},
					{AwsField: "NewInstancesProtectedFromScaleIn", TemplateName: "new-instances-protected", AwsType: "awsbool"},
					{AwsField: "TargetGroupARNs", TemplateName: "targetgroups", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "update", Entity: cloud.ScalingGroup, ApiMethod: "UpdateAutoScalingGroup", Input: "UpdateAutoScalingGroupInput", Output: "UpdateAutoScalingGroupOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "AutoScalingGroupName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "DefaultCooldown", TemplateName: "cooldown", AwsType: "awsint64"},
					{AwsField: "DesiredCapacity", TemplateName: "desired-capacity", AwsType: "awsint64"},
					{AwsField: "HealthCheckGracePeriod", TemplateName: "healthcheck-grace-period", AwsType: "awsint64"},
					{AwsField: "HealthCheckType", TemplateName: "healthcheck-type", AwsType: "awsstr"},
					{AwsField: "LaunchConfigurationName", TemplateName: "launchconfiguration", AwsType: "awsstr"},
					{AwsField: "MaxSize", TemplateName: "max-size", AwsType: "awsint64"},
					{AwsField: "MinSize", TemplateName: "min-size", AwsType: "awsint64"},
					{AwsField: "NewInstancesProtectedFromScaleIn", TemplateName: "new-instances-protected", AwsType: "awsbool"},
					{AwsField: "VPCZoneIdentifier", TemplateName: "subnets", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.ScalingGroup, ApiMethod: "DeleteAutoScalingGroup", Input: "DeleteAutoScalingGroupInput", Output: "DeleteAutoScalingGroupOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "AutoScalingGroupName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "ForceDelete", TemplateName: "force", AwsType: "awsbool"},
				},
			},
			{
				Action: "check", Entity: cloud.ScalingGroup, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
					{TemplateName: "count"},
					{TemplateName: "timeout"},
				},
			},
			{
				Action: "create", Entity: cloud.ScalingPolicy, ApiMethod: "PutScalingPolicy", Input: "PutScalingPolicyInput", Output: "PutScalingPolicyOutput", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.PolicyARN)",
				RequiredParams: []param{
					{AwsField: "AdjustmentType", TemplateName: "adjustment-type", AwsType: "awsstr"},
					{AwsField: "AutoScalingGroupName", TemplateName: "scalinggroup", AwsType: "awsstr"},
					{AwsField: "PolicyName", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "ScalingAdjustment", TemplateName: "adjustment-scaling", AwsType: "awsint64"},
				},
				ExtraParams: []param{
					{AwsField: "Cooldown", TemplateName: "cooldown", AwsType: "awsint64"},
					{AwsField: "MinAdjustmentMagnitude", TemplateName: "adjustment-magnitude", AwsType: "awsint64"},
				},
			},
			{
				Action: "delete", Entity: cloud.ScalingPolicy, ApiMethod: "DeletePolicy", Input: "DeletePolicyInput", Output: "DeletePolicyOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "PolicyName", TemplateName: "id", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api: "rds",
		Drivers: []driver{
			// Database
			{
				Action: "create", Entity: cloud.Database, ApiMethod: "CreateDBInstance", Input: "CreateDBInstanceInput", Output: "CreateDBInstanceOutput", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.DBInstance.DBInstanceIdentifier)",
				RequiredParams: []param{
					{AwsField: "DBInstanceClass", TemplateName: "type", AwsType: "awsstr"},
					{AwsField: "DBInstanceIdentifier", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "Engine", TemplateName: "engine", AwsType: "awsstr"},
					{AwsField: "MasterUserPassword", TemplateName: "password", AwsType: "awsstr"},
					{AwsField: "MasterUsername", TemplateName: "username", AwsType: "awsstr"},
					{AwsField: "AllocatedStorage", TemplateName: "size", AwsType: "awsint64"},
				},
				ExtraParams: []param{
					{AwsField: "AutoMinorVersionUpgrade", TemplateName: "autoupgrade", AwsType: "awsbool"},
					{AwsField: "AvailabilityZone", TemplateName: "availabilityzone", AwsType: "awsstr"},
					{AwsField: "BackupRetentionPeriod", TemplateName: "backupretention", AwsType: "awsint64"},
					{AwsField: "DBClusterIdentifier", TemplateName: "cluster", AwsType: "awsstr"},
					{AwsField: "DBName", TemplateName: "dbname", AwsType: "awsstr"},
					{AwsField: "DBParameterGroupName", TemplateName: "parametergroup", AwsType: "awsstr"},
					{AwsField: "DBSecurityGroups", TemplateName: "dbsecuritygroups", AwsType: "awsstringslice"},
					{AwsField: "DBSubnetGroupName", TemplateName: "subnetgroup", AwsType: "awsstr"},
					{AwsField: "Domain", TemplateName: "domain", AwsType: "awsstr"},
					{AwsField: "DomainIAMRoleName", TemplateName: "iamrole", AwsType: "awsstr"},
					{AwsField: "EngineVersion", TemplateName: "version", AwsType: "awsstr"},
					{AwsField: "Iops", TemplateName: "iops", AwsType: "awsint64"},
					{AwsField: "LicenseModel", TemplateName: "license", AwsType: "awsstr"}, // license-included | bring-your-own-license | general-public-license
					{AwsField: "MultiAZ", TemplateName: "multiaz", AwsType: "awsbool"},
					{AwsField: "OptionGroupName", TemplateName: "optiongroup", AwsType: "awsstr"},
					{AwsField: "Port", TemplateName: "port", AwsType: "awsint64"},
					{AwsField: "PreferredBackupWindow", TemplateName: "backupwindow", AwsType: "awsstr"},
					{AwsField: "PreferredMaintenanceWindow", TemplateName: "maintenancewindow", AwsType: "awsstr"},
					{AwsField: "PubliclyAccessible", TemplateName: "public", AwsType: "awsbool"},
					{AwsField: "StorageEncrypted", TemplateName: "encrypted", AwsType: "awsbool"},
					{AwsField: "StorageType", TemplateName: "storagetype", AwsType: "awsstr"},
					{AwsField: "Timezone", TemplateName: "timezone", AwsType: "awsstr"},
					{AwsField: "VpcSecurityGroupIds", TemplateName: "vpcsecuritygroups", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "delete", Entity: cloud.Database, ApiMethod: "DeleteDBInstance", Input: "DeleteDBInstanceInput", Output: "DeleteDBInstanceOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "DBInstanceIdentifier", TemplateName: "id", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "SkipFinalSnapshot", TemplateName: "skip-snapshot", AwsType: "awsbool"},
					{AwsField: "FinalDBSnapshotIdentifier", TemplateName: "snapshot", AwsType: "awsbool"},
				},
			},
			{
				Action: "check", Entity: cloud.Database, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "state"},
					{TemplateName: "timeout"},
				},
			},
			{
				Action: "create", Entity: cloud.DbSubnetGroup, ApiMethod: "CreateDBSubnetGroup", Input: "CreateDBSubnetGroupInput", Output: "CreateDBSubnetGroupOutput", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.DBSubnetGroup.DBSubnetGroupName)",
				RequiredParams: []param{
					{AwsField: "DBSubnetGroupDescription", TemplateName: "description", AwsType: "awsstr"},
					{AwsField: "DBSubnetGroupName", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "SubnetIds", TemplateName: "subnets", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "delete", Entity: cloud.DbSubnetGroup, ApiMethod: "DeleteDBSubnetGroup", Input: "DeleteDBSubnetGroupInput", Output: "DeleteDBSubnetGroupOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "DBSubnetGroupName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api: "ecr",
		Drivers: []driver{
			// Repository
			{
				Action: "create", Entity: cloud.Repository, ApiMethod: "CreateRepository", Input: "CreateRepositoryInput", Output: "CreateRepositoryOutput", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.Repository.RepositoryArn)",
				RequiredParams: []param{
					{AwsField: "RepositoryName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Repository, ApiMethod: "DeleteRepository", Input: "DeleteRepositoryInput", Output: "DeleteRepositoryOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "RepositoryName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Force", TemplateName: "force", AwsType: "awsbool"},
					{AwsField: "RegistryId", TemplateName: "account", AwsType: "awsstr"},
				},
			},
			// Registry
			{
				Action: "authenticate", Entity: cloud.Registry, ManualFuncDefinition: true,
				RequiredParams: []param{},
				ExtraParams: []param{
					{TemplateName: "accounts"},
					{TemplateName: "no-confirm"},
				},
			},
		},
	},
	{
		Api: "ecs",
		Drivers: []driver{
			//Cluster
			{
				Action: "create", Entity: cloud.ContainerCluster, ApiMethod: "CreateCluster", Input: "CreateClusterInput", Output: "CreateClusterOutput", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.Cluster.ClusterArn)",
				RequiredParams: []param{
					{AwsField: "ClusterName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.ContainerCluster, ApiMethod: "DeleteCluster", Input: "DeleteClusterInput", Output: "DeleteClusterOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "Cluster", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "start", Entity: cloud.ContainerService, ApiMethod: "CreateService", Input: "CreateServiceInput", Output: "CreateServiceOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "Cluster", TemplateName: "cluster", AwsType: "awsstr"},
					{AwsField: "ServiceName", TemplateName: "deployment-name", AwsType: "awsstr"},
					{AwsField: "DesiredCount", TemplateName: "desired-count", AwsType: "awsint64"},
					{AwsField: "TaskDefinition", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Role", TemplateName: "role", AwsType: "awsstr"},
				},
			},
			{
				Action: "stop", Entity: cloud.ContainerService, ApiMethod: "DeleteService", Input: "DeleteServiceInput", Output: "DeleteServiceOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "Cluster", TemplateName: "cluster", AwsType: "awsstr"},
					{AwsField: "Service", TemplateName: "deployment-name", AwsType: "awsstr"},
				},
			},
			{
				Action: "update", Entity: cloud.ContainerService, ApiMethod: "UpdateService", Input: "UpdateServiceInput", Output: "UpdateServiceOutput", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "Cluster", TemplateName: "cluster", AwsType: "awsstr"},
					{AwsField: "Service", TemplateName: "deployment-name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "DesiredCount", TemplateName: "desired-count", AwsType: "awsint64"},
					{AwsField: "TaskDefinition", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			//Container
			{
				Action: "create", Entity: cloud.Container, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
					{TemplateName: "service"},
					{TemplateName: "image"},
					{TemplateName: "memory-hard-limit"},
				},
				ExtraParams: []param{
					{TemplateName: "command"},
					{TemplateName: "env"},
					{TemplateName: "privileged"},
					{TemplateName: "workdir"},
				},
			},
			{
				Action: "delete", Entity: cloud.Container, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
					{TemplateName: "service"},
				},
			},
		},
	},
	{
		Api:     "sts",
		Drivers: []driver{},
	},
	{
		Api: "iam",
		Drivers: []driver{
			// USER
			{
				Action: "create", Entity: cloud.User, DryRunUnsupported: true, Input: "CreateUserInput", Output: "CreateUserOutput", ApiMethod: "CreateUser", OutputExtractor: "aws.StringValue(output.User.UserId)",
				RequiredParams: []param{
					{AwsField: "UserName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.User, DryRunUnsupported: true, Input: "DeleteUserInput", Output: "DeleteUserOutput", ApiMethod: "DeleteUser",
				RequiredParams: []param{
					{AwsField: "UserName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.User, DryRunUnsupported: true, Input: "AddUserToGroupInput", Output: "AddUserToGroupOutput", ApiMethod: "AddUserToGroup",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "group", AwsType: "awsstr"},
					{AwsField: "UserName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: cloud.User, DryRunUnsupported: true, Input: "RemoveUserFromGroupInput", Output: "RemoveUserFromGroupOutput", ApiMethod: "RemoveUserFromGroup",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "group", AwsType: "awsstr"},
					{AwsField: "UserName", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// Access key
			{
				Action: "create", Entity: cloud.AccessKey, DryRunUnsupported: true, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "user"},
				},
			},
			{
				Action: "delete", Entity: cloud.AccessKey, DryRunUnsupported: true, ApiMethod: "DeleteAccessKey", Input: "DeleteAccessKeyInput", Output: "DeleteAccessKeyOutput",
				RequiredParams: []param{
					{AwsField: "AccessKeyId", TemplateName: "id", AwsType: "awsstr"},
				}, ExtraParams: []param{
					{AwsField: "UserName", TemplateName: "user", AwsType: "awsstr"},
				},
			},
			// Access key
			{
				Action: "create", Entity: cloud.LoginProfile, DryRunUnsupported: true, Input: "CreateLoginProfileInput", Output: "CreateLoginProfileOutput", ApiMethod: "CreateLoginProfile", OutputExtractor: "aws.StringValue(output.LoginProfile.UserName)",
				RequiredParams: []param{
					{AwsField: "UserName", TemplateName: "username", AwsType: "awsstr"},
					{AwsField: "Password", TemplateName: "password", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "PasswordResetRequired", TemplateName: "password-reset", AwsType: "awsbool"},
				},
			},
			{
				Action: "update", Entity: cloud.LoginProfile, DryRunUnsupported: true, Input: "UpdateLoginProfileInput", Output: "UpdateLoginProfileOutput", ApiMethod: "UpdateLoginProfile",
				RequiredParams: []param{
					{AwsField: "UserName", TemplateName: "username", AwsType: "awsstr"},
					{AwsField: "Password", TemplateName: "password", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "PasswordResetRequired", TemplateName: "password-reset", AwsType: "awsbool"},
				},
			},
			{
				Action: "delete", Entity: cloud.LoginProfile, DryRunUnsupported: true, ApiMethod: "DeleteLoginProfile", Input: "DeleteLoginProfileInput", Output: "DeleteLoginProfileOutput",
				RequiredParams: []param{
					{AwsField: "UserName", TemplateName: "username", AwsType: "awsstr"},
				},
			},

			// GROUP
			{
				Action: "create", Entity: cloud.Group, DryRunUnsupported: true, Input: "CreateGroupInput", Output: "CreateGroupOutput", ApiMethod: "CreateGroup", OutputExtractor: "aws.StringValue(output.Group.GroupId)",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Group, DryRunUnsupported: true, Input: "DeleteGroupInput", Output: "DeleteGroupOutput", ApiMethod: "DeleteGroup",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// ROLE
			{
				Action: "create", Entity: cloud.Role, ManualFuncDefinition: true,
				RequiredParams: []param{
					{AwsField: "RoleName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{TemplateName: "principal-account"},
					{TemplateName: "principal-user"},
					{TemplateName: "principal-service"},
					{TemplateName: "sleep-after"},
				},
			},
			{
				Action: "delete", Entity: cloud.Role, ManualFuncDefinition: true,
				RequiredParams: []param{
					{AwsField: "RoleName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.Role, DryRunUnsupported: true, Input: "AddRoleToInstanceProfileInput", Output: "AddRoleToInstanceProfileOutput", ApiMethod: "AddRoleToInstanceProfile",
				RequiredParams: []param{
					{AwsField: "InstanceProfileName", TemplateName: "instanceprofile", AwsType: "awsstr"},
					{AwsField: "RoleName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: cloud.Role, DryRunUnsupported: true, Input: "RemoveRoleFromInstanceProfileInput", Output: "RemoveRoleFromInstanceProfileOutput", ApiMethod: "RemoveRoleFromInstanceProfile",
				RequiredParams: []param{
					{AwsField: "InstanceProfileName", TemplateName: "instanceprofile", AwsType: "awsstr"},
					{AwsField: "RoleName", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// INSTANCE PROFILE
			{
				Action: "create", Entity: cloud.InstanceProfile, DryRunUnsupported: true, Input: "CreateInstanceProfileInput", Output: "CreateInstanceProfileOutput", ApiMethod: "CreateInstanceProfile",
				RequiredParams: []param{
					{AwsField: "InstanceProfileName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.InstanceProfile, DryRunUnsupported: true, Input: "DeleteInstanceProfileInput", Output: "DeleteInstanceProfileOutput", ApiMethod: "DeleteInstanceProfile",
				RequiredParams: []param{
					{AwsField: "InstanceProfileName", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// POLICY
			{
				Action: "create", Entity: cloud.Policy, ManualFuncDefinition: true,
				RequiredParams: []param{
					{AwsField: "PolicyName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Description", TemplateName: "description", AwsType: "awsstr"},
					{TemplateName: "effect"},
					{TemplateName: "action"},
					{TemplateName: "resource"},
				},
			},
			{
				Action: "delete", Entity: cloud.Policy, DryRunUnsupported: true, Input: "DeletePolicyInput", Output: "DeletePolicyOutput", ApiMethod: "DeletePolicy",
				RequiredParams: []param{
					{AwsField: "PolicyArn", TemplateName: "arn", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.Policy, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "arn"},
				},
				ExtraParams: []param{
					{TemplateName: "user"},
					{TemplateName: "group"},
					{TemplateName: "role"},
				},
			},
			{
				Action: "detach", Entity: cloud.Policy, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "arn"},
				},
				ExtraParams: []param{
					{TemplateName: "user"},
					{TemplateName: "group"},
					{TemplateName: "role"},
				},
			},
		},
	},
	{
		Api: "s3",
		Drivers: []driver{
			// BUCKET
			{
				Action: "create", Entity: cloud.Bucket, DryRunUnsupported: true, Input: "CreateBucketInput", Output: "CreateBucketOutput", ApiMethod: "CreateBucket", OutputExtractor: "params[\"name\"]",
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "ACL", TemplateName: "acl", AwsType: "awsstr"},
				},
			},
			{
				Action: "update", Entity: cloud.Bucket, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
				},
				ExtraParams: []param{
					{TemplateName: "acl"},
					{TemplateName: "public-website"},
					{TemplateName: "redirect-hostname"},
					{TemplateName: "index-suffix"},
					{TemplateName: "enforce-https"},
				},
			},
			{
				Action: "delete", Entity: cloud.Bucket, DryRunUnsupported: true, Input: "DeleteBucketInput", Output: "DeleteBucketOutput", ApiMethod: "DeleteBucket",
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// OBJECT
			{
				Action: "create", Entity: cloud.S3Object, ManualFuncDefinition: true,
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "bucket", AwsType: "awsstr"},
					{AwsField: "Body", TemplateName: "file", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Key", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "ACL", TemplateName: "acl", AwsType: "awsstr"},
				},
			},
			{
				Action: "update", Entity: cloud.S3Object, DryRunUnsupported: true, Input: "PutObjectAclInput", Output: "PutObjectAclOutput", ApiMethod: "PutObjectAcl",
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "bucket", AwsType: "awsstr"},
					{AwsField: "Key", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "ACL", TemplateName: "acl", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "VersionId", TemplateName: "version", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.S3Object, DryRunUnsupported: true, Input: "DeleteObjectInput", Output: "DeleteObjectOutput", ApiMethod: "DeleteObject",
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "bucket", AwsType: "awsstr"},
					{AwsField: "Key", TemplateName: "name", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api: "sns",
		Drivers: []driver{
			// TOPIC
			{
				Action: "create", Entity: cloud.Topic, DryRunUnsupported: true, Input: "CreateTopicInput", Output: "CreateTopicOutput", ApiMethod: "CreateTopic", OutputExtractor: "aws.StringValue(output.TopicArn)",
				RequiredParams: []param{
					{AwsField: "Name", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Topic, DryRunUnsupported: true, Input: "DeleteTopicInput", Output: "DeleteTopicOutput", ApiMethod: "DeleteTopic",
				RequiredParams: []param{
					{AwsField: "TopicArn", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			//Subscription
			{
				Action: "create", Entity: cloud.Subscription, DryRunUnsupported: true, Input: "SubscribeInput", Output: "SubscribeOutput", ApiMethod: "Subscribe", OutputExtractor: "aws.StringValue(output.SubscriptionArn)",
				RequiredParams: []param{
					{AwsField: "TopicArn", TemplateName: "topic", AwsType: "awsstr"},
					{AwsField: "Endpoint", TemplateName: "endpoint", AwsType: "awsstr"},
					{AwsField: "Protocol", TemplateName: "protocol", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Subscription, DryRunUnsupported: true, Input: "UnsubscribeInput", Output: "UnsubscribeOutput", ApiMethod: "Unsubscribe",
				RequiredParams: []param{
					{AwsField: "SubscriptionArn", TemplateName: "id", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api: "sqs",
		Drivers: []driver{
			// QUEUE
			{
				Action: "create", Entity: cloud.Queue, DryRunUnsupported: true, Input: "CreateQueueInput", Output: "CreateQueueOutput", ApiMethod: "CreateQueue", OutputExtractor: "aws.StringValue(output.QueueUrl)",
				RequiredParams: []param{
					{AwsField: "QueueName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Attributes[DelaySeconds]", TemplateName: "delay", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[MaximumMessageSize]", TemplateName: "max-msg-size", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[MessageRetentionPeriod]", TemplateName: "retention-period", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[Policy]", TemplateName: "policy", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[ReceiveMessageWaitTimeSeconds]", TemplateName: "msg-wait", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[RedrivePolicy]", TemplateName: "redrive-policy", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[VisibilityTimeout]", TemplateName: "visibility-timeout", AwsType: "awsstringpointermap"},
				},
			},
			{
				Action: "delete", Entity: cloud.Queue, DryRunUnsupported: true, Input: "DeleteQueueInput", Output: "DeleteQueueOutput", ApiMethod: "DeleteQueue",
				RequiredParams: []param{
					{AwsField: "QueueUrl", TemplateName: "url", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api:          "route53",
		ApiInterface: "Route53API",
		Drivers: []driver{
			{
				Action: "create", Entity: cloud.Zone, DryRunUnsupported: true, Input: "CreateHostedZoneInput", Output: "CreateHostedZoneOutput", ApiMethod: "CreateHostedZone", OutputExtractor: "aws.StringValue(output.HostedZone.Id)",
				RequiredParams: []param{
					{AwsField: "CallerReference", TemplateName: "callerreference", AwsType: "awsstr"}, // unique string (random/date/timestamp)
					{AwsField: "Name", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "DelegationSetId", TemplateName: "delegationsetid", AwsType: "awsstr"},
					{AwsField: "HostedZoneConfig.Comment", TemplateName: "comment", AwsType: "awsstr"},
					{AwsField: "HostedZoneConfig.PrivateZone", TemplateName: "isprivate", AwsType: "awsbool"},
					{AwsField: "VPC.VPCId", TemplateName: "vpcid", AwsType: "awsstr"},
					{AwsField: "VPC.VPCRegion", TemplateName: "vpcregion", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Zone, DryRunUnsupported: true, Input: "DeleteHostedZoneInput", Output: "DeleteHostedZoneOutput", ApiMethod: "DeleteHostedZone",
				RequiredParams: []param{
					{AwsField: "Id", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "create", Entity: cloud.Record, DryRunUnsupported: true, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "zone"},
					{TemplateName: "name"},
					{TemplateName: "type"},
					{TemplateName: "value"},
					{TemplateName: "ttl"},
				},
				ExtraParams: []param{
					{TemplateName: "comment"},
				},
			},
			{
				Action: "delete", Entity: cloud.Record, DryRunUnsupported: true, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "zone"},
					{TemplateName: "name"},
					{TemplateName: "type"},
					{TemplateName: "value"},
					{TemplateName: "ttl"},
				},
			},
		},
	},
	{
		Api:          "lambda",
		ApiInterface: "LambdaAPI",
		Drivers: []driver{
			{
				Action: "create", Entity: cloud.Function, DryRunUnsupported: true, ApiMethod: "CreateFunction", Input: "CreateFunctionInput", Output: "FunctionConfiguration", OutputExtractor: "aws.StringValue(output.FunctionArn)",
				RequiredParams: []param{
					{AwsField: "FunctionName", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "Handler", TemplateName: "handler", AwsType: "awsstr"},
					{AwsField: "Role", TemplateName: "role", AwsType: "awsstr"},
					{AwsField: "Runtime", TemplateName: "runtime", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Code.S3Bucket", TemplateName: "bucket", AwsType: "awsstr"},
					{AwsField: "Code.S3Key", TemplateName: "object", AwsType: "awsstr"},
					{AwsField: "Code.S3ObjectVersion", TemplateName: "objectversion", AwsType: "awsstr"},
					{AwsField: "Code.ZipFile", TemplateName: "zipfile", AwsType: "awsfiletobyteslice"},
					{AwsField: "Description", TemplateName: "description", AwsType: "awsstr"},
					{AwsField: "MemorySize", TemplateName: "memory", AwsType: "awsint64"},
					{AwsField: "Publish", TemplateName: "publish", AwsType: "awsbool"},
					{AwsField: "Timeout", TemplateName: "timeout", AwsType: "awsint64"},
				},
			},
			{
				Action: "delete", Entity: cloud.Function, DryRunUnsupported: true, ApiMethod: "DeleteFunction", Input: "DeleteFunctionInput", Output: "DeleteFunctionOutput",
				RequiredParams: []param{
					{AwsField: "FunctionName", TemplateName: "id", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Qualifier", TemplateName: "version", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api:          "cloudwatch",
		ApiInterface: "CloudWatchAPI",
		Drivers: []driver{
			{
				Action: "create", Entity: cloud.Alarm, DryRunUnsupported: true, ApiMethod: "PutMetricAlarm", Input: "PutMetricAlarmInput", Output: "PutMetricAlarmOutput", OutputExtractor: "params[\"name\"]",
				RequiredParams: []param{
					{AwsField: "AlarmName", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "ComparisonOperator", TemplateName: "operator", AwsType: "awsstr"}, // [GreaterThanThreshold, LessThanThreshold, LessThanOrEqualToThreshold, GreaterThanOrEqualToThreshold]
					{AwsField: "MetricName", TemplateName: "metric", AwsType: "awsstr"},
					{AwsField: "Namespace", TemplateName: "namespace", AwsType: "awsstr"},
					{AwsField: "EvaluationPeriods", TemplateName: "evaluation-periods", AwsType: "awsint64"},
					{AwsField: "Period", TemplateName: "period", AwsType: "awsint64"},
					{AwsField: "Statistic", TemplateName: "statistic-function", AwsType: "awsstr"}, // Minimum, Maximum, Sum, Average, SampleCount, pNN.NN
					{AwsField: "Threshold", TemplateName: "threshold", AwsType: "awsfloat"},
				},
				ExtraParams: []param{
					{AwsField: "ActionsEnabled", TemplateName: "enabled", AwsType: "awsbool"},
					{AwsField: "AlarmActions", TemplateName: "alarm-actions", AwsType: "awsstringslice"},
					{AwsField: "InsufficientDataActions", TemplateName: "insufficientdata-actions", AwsType: "awsstringslice"},
					{AwsField: "OKActions", TemplateName: "ok-actions", AwsType: "awsstringslice"},
					{AwsField: "AlarmDescription", TemplateName: "description", AwsType: "awsstr"},
					{AwsField: "Dimensions", TemplateName: "dimensions", AwsType: "awsdimensionslice"},
					{AwsField: "Unit", TemplateName: "unit", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Alarm, DryRunUnsupported: true, ApiMethod: "DeleteAlarms", Input: "DeleteAlarmsInput", Output: "DeleteAlarmsOutput",
				RequiredParams: []param{
					{AwsField: "AlarmNames", TemplateName: "name", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "start", Entity: cloud.Alarm, DryRunUnsupported: true, ApiMethod: "EnableAlarmActions", Input: "EnableAlarmActionsInput", Output: "EnableAlarmActionsOutput",
				RequiredParams: []param{
					{AwsField: "AlarmNames", TemplateName: "names", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "stop", Entity: cloud.Alarm, DryRunUnsupported: true, ApiMethod: "DisableAlarmActions", Input: "DisableAlarmActionsInput", Output: "DisableAlarmActionsOutput",
				RequiredParams: []param{
					{AwsField: "AlarmNames", TemplateName: "names", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "attach", Entity: cloud.Alarm, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
					{TemplateName: "action-arn"},
				},
			},
			{
				Action: "detach", Entity: cloud.Alarm, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
					{TemplateName: "action-arn"},
				},
			},
		},
	},
	{
		Api:          "cloudfront",
		ApiInterface: "CloudFrontAPI",
		Drivers: []driver{
			{
				Action: "create", Entity: cloud.Distribution, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "origin-domain"},
				},
				ExtraParams: []param{
					{TemplateName: "certificate"},
					{TemplateName: "comment"},
					{TemplateName: "default-file"},
					{TemplateName: "domain-aliases"},
					{TemplateName: "enable"},
					{TemplateName: "forward-cookies"},
					{TemplateName: "forward-queries"},
					{TemplateName: "https-behaviour"},
					{TemplateName: "origin-path"},
					{TemplateName: "price-class"},
					{TemplateName: "min-ttl"},
				},
			},
			{
				Action: "check", Entity: cloud.Distribution, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "state"},
					{TemplateName: "timeout"},
				},
			},
			{
				Action: "update", Entity: cloud.Distribution, ManualFuncDefinition: true,
				RequiredParams: []param{
					{AwsField: "Id", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "DistributionConfig.Enabled", TemplateName: "enable", AwsType: "awsbool"},
				},
			},
			{
				Action: "delete", Entity: cloud.Distribution, ManualFuncDefinition: true,
				RequiredParams: []param{
					{AwsField: "Id", TemplateName: "id", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api:          "cloudformation",
		ApiInterface: "CloudFormationAPI",
		Drivers: []driver{
			{
				Action: "create", Entity: cloud.Stack, DryRunUnsupported: true, ApiMethod: "CreateStack", Input: "CreateStackInput", Output: "CreateStackOutput", OutputExtractor: "aws.StringValue(output.StackId)",
				RequiredParams: []param{
					{AwsField: "StackName", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "TemplateBody", TemplateName: "template-file", AwsType: "awsfiletostring"},
				},
				ExtraParams: []param{
					{AwsField: "Capabilities", TemplateName: "capabilities", AwsType: "awsstringslice"}, //CAPABILITY_IAM and CAPABILITY_NAMED_IAM
					{AwsField: "DisableRollback", TemplateName: "disable-rollback", AwsType: "awsbool"},
					{AwsField: "NotificationARNs", TemplateName: "notifications", AwsType: "awsstringslice"},
					{AwsField: "OnFailure", TemplateName: "on-failure", AwsType: "awsstr"},                 //DO_NOTHING, ROLLBACK, or DELETE
					{AwsField: "Parameters", TemplateName: "parameters", AwsType: "awsparameterslice"},     //Format, key1:val1,key2:val2,...
					{AwsField: "ResourceTypes", TemplateName: "resource-types", AwsType: "awsstringslice"}, //AWS::EC2::Instance, AWS::EC2::*, or Custom::MyCustomInstance or Custom::*
					{AwsField: "RoleARN", TemplateName: "role", AwsType: "awsstr"},
					{AwsField: "StackPolicyBody", TemplateName: "policy-file", AwsType: "awsfiletostring"},
					{AwsField: "TimeoutInMinutes", TemplateName: "timeout", AwsType: "awsint64"},
				},
			},
			{
				Action: "update", Entity: cloud.Stack, DryRunUnsupported: true, ApiMethod: "UpdateStack", Input: "UpdateStackInput", Output: "UpdateStackOutput", OutputExtractor: "aws.StringValue(output.StackId)",
				RequiredParams: []param{
					{AwsField: "StackName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Capabilities", TemplateName: "capabilities", AwsType: "awsstringslice"}, //CAPABILITY_IAM and CAPABILITY_NAMED_IAM
					{AwsField: "NotificationARNs", TemplateName: "notifications", AwsType: "awsstringslice"},
					{AwsField: "Parameters", TemplateName: "parameters", AwsType: "awsparameterslice"},     //Format, key1:val1,key2:val2,...
					{AwsField: "ResourceTypes", TemplateName: "resource-types", AwsType: "awsstringslice"}, //AWS::EC2::Instance, AWS::EC2::*, or Custom::MyCustomInstance or Custom::*
					{AwsField: "RoleARN", TemplateName: "role", AwsType: "awsstr"},
					{AwsField: "StackPolicyBody", TemplateName: "policy-file", AwsType: "awsfiletostring"},
					{AwsField: "StackPolicyDuringUpdateBody", TemplateName: "policy-update-file", AwsType: "awsfiletostring"},
					{AwsField: "TemplateBody", TemplateName: "template-file", AwsType: "awsfiletostring"},
					{AwsField: "UsePreviousTemplate", TemplateName: "use-previous-template", AwsType: "awsbool"},
				},
			},
			{
				Action: "delete", Entity: cloud.Stack, DryRunUnsupported: true, ApiMethod: "DeleteStack", Input: "DeleteStackInput", Output: "DeleteStackOutput",
				RequiredParams: []param{
					{AwsField: "StackName", TemplateName: "name", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "RetainResources", TemplateName: "retain-resources", AwsType: "awsstringslice"},
				},
			},
		},
	},
}
