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

import "github.com/wallix/awless/cloud"

type param struct {
	AwsField, AwsType string
	TemplateName      string
}

type driver struct {
	RequiredParams                            []param
	ExtraParams                               []param
	TagsMapping                               map[string]string
	Action, Entity                            string
	Input, Output, ApiMethod, OutputExtractor string
	DryRunUnsupported                         bool
	ManualFuncDefinition                      bool
}

type driversDef struct {
	Api          string
	ApiInterface string
	Drivers      []driver
}

var DriversDefs = []driversDef{
	{
		Api: "ec2",
		Drivers: []driver{
			// VPC
			{
				Action: "create", Entity: cloud.Vpc, Input: "CreateVpcInput", Output: "CreateVpcOutput", ApiMethod: "CreateVpc", OutputExtractor: "aws.StringValue(output.Vpc.VpcId)",
				RequiredParams: []param{
					{AwsField: "CidrBlock", TemplateName: "cidr", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Vpc, Input: "DeleteVpcInput", Output: "DeleteVpcOutput", ApiMethod: "DeleteVpc",
				RequiredParams: []param{
					{AwsField: "VpcId", TemplateName: "id", AwsType: "awsstr"},
				},
			},

			// SUBNET
			{
				Action: "create", Entity: cloud.Subnet, Input: "CreateSubnetInput", Output: "CreateSubnetOutput", ApiMethod: "CreateSubnet", OutputExtractor: "aws.StringValue(output.Subnet.SubnetId)",
				RequiredParams: []param{
					{AwsField: "CidrBlock", TemplateName: "cidr", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "AvailabilityZone", TemplateName: "zone", AwsType: "awsstr"},
				},
				TagsMapping: map[string]string{
					"Name": "name",
				},
			},
			{
				Action: "update", Entity: cloud.Subnet, Input: "ModifySubnetAttributeInput", Output: "ModifySubnetAttributeOutput", ApiMethod: "ModifySubnetAttribute", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "SubnetId", TemplateName: "id", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "MapPublicIpOnLaunch", TemplateName: "public", AwsType: "awsboolattribute"},
				},
			},
			{
				Action: "delete", Entity: cloud.Subnet, Input: "DeleteSubnetInput", Output: "DeleteSubnetOutput", ApiMethod: "DeleteSubnet",
				RequiredParams: []param{
					{AwsField: "SubnetId", TemplateName: "id", AwsType: "awsstr"},
				},
			},

			// INSTANCES
			{
				Action: "create", Entity: cloud.Instance, Input: "RunInstancesInput", Output: "Reservation", ApiMethod: "RunInstances", OutputExtractor: "aws.StringValue(output.Instances[0].InstanceId)",
				RequiredParams: []param{
					{AwsField: "ImageId", TemplateName: "image", AwsType: "awsstr"},
					{AwsField: "MaxCount", TemplateName: "count", AwsType: "awsint64"},
					{AwsField: "MinCount", TemplateName: "count", AwsType: "awsint64"},
					{AwsField: "InstanceType", TemplateName: "type", AwsType: "awsstr"},
					{AwsField: "SubnetId", TemplateName: "subnet", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "KeyName", TemplateName: "key", AwsType: "awsstr"},
					{AwsField: "PrivateIpAddress", TemplateName: "ip", AwsType: "awsstr"},
					{AwsField: "UserData", TemplateName: "userdata", AwsType: "awsstr"},
					{AwsField: "SecurityGroupIds", TemplateName: "group", AwsType: "awsstringslice"},
					{AwsField: "DisableApiTermination", TemplateName: "lock", AwsType: "awsboolattribute"},
				},
				TagsMapping: map[string]string{
					"Name": "name",
				},
			},
			{
				Action: "update", Entity: cloud.Instance, Input: "ModifyInstanceAttributeInput", Output: "ModifyInstanceAttributeOutput", ApiMethod: "ModifyInstanceAttribute",
				RequiredParams: []param{
					{AwsField: "InstanceId", TemplateName: "id", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "InstanceType", TemplateName: "type", AwsType: "awsstr"},
					{AwsField: "Groups", TemplateName: "group", AwsType: "awsstringslice"},
					{AwsField: "DisableApiTermination", TemplateName: "lock", AwsType: "awsboolattribute"},
				},
			},
			{
				Action: "delete", Entity: cloud.Instance, Input: "TerminateInstancesInput", Output: "TerminateInstancesOutput", ApiMethod: "TerminateInstances",
				RequiredParams: []param{
					{AwsField: "InstanceIds", TemplateName: "id", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "start", Entity: cloud.Instance, Input: "StartInstancesInput", Output: "StartInstancesOutput", ApiMethod: "StartInstances", OutputExtractor: "aws.StringValue(output.StartingInstances[0].InstanceId)",
				RequiredParams: []param{
					{AwsField: "InstanceIds", TemplateName: "id", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "stop", Entity: cloud.Instance, Input: "StopInstancesInput", Output: "StopInstancesOutput", ApiMethod: "StopInstances", OutputExtractor: "aws.StringValue(output.StoppingInstances[0].InstanceId)",
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
				Action: "create", Entity: cloud.SecurityGroup, Input: "CreateSecurityGroupInput", Output: "CreateSecurityGroupOutput", ApiMethod: "CreateSecurityGroup", OutputExtractor: "aws.StringValue(output.GroupId)",
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
				Action: "delete", Entity: cloud.SecurityGroup, Input: "DeleteSecurityGroupInput", Output: "DeleteSecurityGroupOutput", ApiMethod: "DeleteSecurityGroup",
				RequiredParams: []param{
					{AwsField: "GroupId", TemplateName: "id", AwsType: "awsstr"},
				},
			},

			// VOLUME
			{
				Action: "create", Entity: cloud.Volume, Input: "CreateVolumeInput", Output: "Volume", ApiMethod: "CreateVolume", OutputExtractor: "aws.StringValue(output.VolumeId)",
				RequiredParams: []param{
					{AwsField: "AvailabilityZone", TemplateName: "zone", AwsType: "awsstr"},
					{AwsField: "Size", TemplateName: "size", AwsType: "awsint64"},
				},
			},
			{
				Action: "delete", Entity: cloud.Volume, Input: "DeleteVolumeInput", Output: "DeleteVolumeOutput", ApiMethod: "DeleteVolume",
				RequiredParams: []param{
					{AwsField: "VolumeId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.Volume, Input: "AttachVolumeInput", Output: "VolumeAttachment", ApiMethod: "AttachVolume", OutputExtractor: "aws.StringValue(output.VolumeId)",
				RequiredParams: []param{
					{AwsField: "Device", TemplateName: "device", AwsType: "awsstr"},
					{AwsField: "VolumeId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "InstanceId", TemplateName: "instance", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: cloud.Volume, Input: "DetachVolumeInput", Output: "VolumeAttachment", ApiMethod: "DetachVolume", OutputExtractor: "aws.StringValue(output.VolumeId)",
				RequiredParams: []param{
					{AwsField: "Device", TemplateName: "device", AwsType: "awsstr"},
					{AwsField: "VolumeId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "InstanceId", TemplateName: "instance", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Force", TemplateName: "force", AwsType: "awsbool"},
				},
			},
			// INTERNET GATEWAYS
			{
				Action: "create", Entity: cloud.InternetGateway, Input: "CreateInternetGatewayInput", Output: "CreateInternetGatewayOutput", ApiMethod: "CreateInternetGateway", OutputExtractor: "aws.StringValue(output.InternetGateway.InternetGatewayId)",
			},
			{
				Action: "delete", Entity: cloud.InternetGateway, Input: "DeleteInternetGatewayInput", Output: "DeleteInternetGatewayOutput", ApiMethod: "DeleteInternetGateway",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.InternetGateway, Input: "AttachInternetGatewayInput", Output: "AttachInternetGatewayOutput", ApiMethod: "AttachInternetGateway",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: cloud.InternetGateway, Input: "DetachInternetGatewayInput", Output: "DetachInternetGatewayOutput", ApiMethod: "DetachInternetGateway",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
			},
			// ROUTE TABLES
			{
				Action: "create", Entity: cloud.RouteTable, Input: "CreateRouteTableInput", Output: "CreateRouteTableOutput", ApiMethod: "CreateRouteTable", OutputExtractor: "aws.StringValue(output.RouteTable.RouteTableId)",
				RequiredParams: []param{
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"}},
			},
			{
				Action: "delete", Entity: cloud.RouteTable, Input: "DeleteRouteTableInput", Output: "DeleteRouteTableOutput", ApiMethod: "DeleteRouteTable",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: cloud.RouteTable, Input: "AssociateRouteTableInput", Output: "AssociateRouteTableOutput", ApiMethod: "AssociateRouteTable", OutputExtractor: "aws.StringValue(output.AssociationId)",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "SubnetId", TemplateName: "subnet", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: cloud.RouteTable, Input: "DisassociateRouteTableInput", Output: "DisassociateRouteTableOutput", ApiMethod: "DisassociateRouteTable",
				RequiredParams: []param{
					{AwsField: "AssociationId", TemplateName: "association", AwsType: "awsstr"},
				},
			},
			// ROUTES
			{
				Action: "create", Entity: "route", Input: "CreateRouteInput", Output: "CreateRouteOutput", ApiMethod: "CreateRoute",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "table", AwsType: "awsstr"},
					{AwsField: "DestinationCidrBlock", TemplateName: "cidr", AwsType: "awsstr"},
					{AwsField: "GatewayId", TemplateName: "gateway", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: "route", Input: "DeleteRouteInput", Output: "DeleteRouteOutput", ApiMethod: "DeleteRoute",
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

			// Keypair
			{
				Action: "create", Entity: cloud.Keypair, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
				},
			},
			{
				Action: "delete", Entity: cloud.Keypair, Input: "DeleteKeyPairInput", Output: "DeleteKeyPairOutput", ApiMethod: "DeleteKeyPair",
				RequiredParams: []param{
					{AwsField: "KeyName", TemplateName: "id", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api: "elbv2",
		Drivers: []driver{
			// LoadBalancer
			{
				Action: "create", Entity: cloud.LoadBalancer, Input: "CreateLoadBalancerInput", Output: "CreateLoadBalancerOutput", ApiMethod: "CreateLoadBalancer", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.LoadBalancers[0].LoadBalancerArn)",
				RequiredParams: []param{
					{AwsField: "Name", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "Subnets", TemplateName: "subnets", AwsType: "awsstringslice"},
				},
				ExtraParams: []param{
					{AwsField: "IpAddressType", TemplateName: "iptype", AwsType: "awsstr"},
					{AwsField: "Scheme", TemplateName: "scheme", AwsType: "awsstr"},
					{AwsField: "SecurityGroups", TemplateName: "groups", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "delete", Entity: cloud.LoadBalancer, Input: "DeleteLoadBalancerInput", Output: "DeleteLoadBalancerOutput", ApiMethod: "DeleteLoadBalancer", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "LoadBalancerArn", TemplateName: "arn", AwsType: "awsstr"},
				},
			},
			// Listener
			{
				Action: "create", Entity: cloud.Listener, Input: "CreateListenerInput", Output: "CreateListenerOutput", ApiMethod: "CreateListener", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.Listeners[0].ListenerArn)",
				RequiredParams: []param{
					{AwsField: "DefaultActions[0]Type", TemplateName: "actiontype", AwsType: "awsslicestruct"}, //always forward
					{AwsField: "DefaultActions[0]TargetGroupArn", TemplateName: "target", AwsType: "awsslicestruct"},
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
				Action: "delete", Entity: cloud.Listener, Input: "DeleteListenerInput", Output: "DeleteListenerOutput", ApiMethod: "DeleteListener", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "ListenerArn", TemplateName: "arn", AwsType: "awsstr"},
				},
			},
			// Target group
			{
				Action: "create", Entity: cloud.TargetGroup, Input: "CreateTargetGroupInput", Output: "CreateTargetGroupOutput", ApiMethod: "CreateTargetGroup", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.TargetGroups[0].TargetGroupArn)",
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
				Action: "delete", Entity: cloud.TargetGroup, Input: "DeleteTargetGroupInput", Output: "DeleteTargetGroupOutput", ApiMethod: "DeleteTargetGroup", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "TargetGroupArn", TemplateName: "arn", AwsType: "awsstr"},
				},
			},
		},
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

			// POLICY
			{
				Action: "attach", Entity: cloud.Policy, ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "arn"},
				},
				ExtraParams: []param{
					{TemplateName: "user"},
					{TemplateName: "group"},
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
			},
			{
				Action: "delete", Entity: cloud.Bucket, DryRunUnsupported: true, Input: "DeleteBucketInput", Output: "DeleteBucketOutput", ApiMethod: "DeleteBucket",
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// OBJECT
			{
				Action: "create", Entity: cloud.Object, ManualFuncDefinition: true,
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "bucket", AwsType: "awsstr"},
					{AwsField: "Body", TemplateName: "file", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Key", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: cloud.Object, DryRunUnsupported: true, Input: "DeleteObjectInput", Output: "DeleteObjectOutput", ApiMethod: "DeleteObject",
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "bucket", AwsType: "awsstr"},
					{AwsField: "Key", TemplateName: "key", AwsType: "awsstr"},
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
					{AwsField: "TopicArn", TemplateName: "arn", AwsType: "awsstr"},
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
					{AwsField: "SubscriptionArn", TemplateName: "arn", AwsType: "awsstr"},
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
					{AwsField: "Attributes[MaximumMessageSize]", TemplateName: "maxMsgSize", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[MessageRetentionPeriod]", TemplateName: "retentionPeriod", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[Policy]", TemplateName: "policy", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[ReceiveMessageWaitTimeSeconds]", TemplateName: "msgWait", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[RedrivePolicy]", TemplateName: "redrivePolicy", AwsType: "awsstringpointermap"},
					{AwsField: "Attributes[VisibilityTimeout]", TemplateName: "visibilityTimeout", AwsType: "awsstringpointermap"},
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
}
