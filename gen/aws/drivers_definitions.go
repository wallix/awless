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

import "github.com/wallix/awless/graph"

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
				Action: "create", Entity: graph.Vpc.String(), Input: "CreateVpcInput", Output: "CreateVpcOutput", ApiMethod: "CreateVpc", OutputExtractor: "aws.StringValue(output.Vpc.VpcId)",
				RequiredParams: []param{
					{AwsField: "CidrBlock", TemplateName: "cidr", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: graph.Vpc.String(), Input: "DeleteVpcInput", Output: "DeleteVpcOutput", ApiMethod: "DeleteVpc",
				RequiredParams: []param{
					{AwsField: "VpcId", TemplateName: "id", AwsType: "awsstr"},
				},
			},

			// SUBNET
			{
				Action: "create", Entity: graph.Subnet.String(), Input: "CreateSubnetInput", Output: "CreateSubnetOutput", ApiMethod: "CreateSubnet", OutputExtractor: "aws.StringValue(output.Subnet.SubnetId)",
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
				Action: "update", Entity: graph.Subnet.String(), Input: "ModifySubnetAttributeInput", Output: "ModifySubnetAttributeOutput", ApiMethod: "ModifySubnetAttribute", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "SubnetId", TemplateName: "id", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "MapPublicIpOnLaunch", TemplateName: "public", AwsType: "awsboolattribute"},
				},
			},
			{
				Action: "delete", Entity: graph.Subnet.String(), Input: "DeleteSubnetInput", Output: "DeleteSubnetOutput", ApiMethod: "DeleteSubnet",
				RequiredParams: []param{
					{AwsField: "SubnetId", TemplateName: "id", AwsType: "awsstr"},
				},
			},

			// INSTANCES
			{
				Action: "create", Entity: graph.Instance.String(), Input: "RunInstancesInput", Output: "Reservation", ApiMethod: "RunInstances", OutputExtractor: "aws.StringValue(output.Instances[0].InstanceId)",
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
				Action: "update", Entity: graph.Instance.String(), Input: "ModifyInstanceAttributeInput", Output: "ModifyInstanceAttributeOutput", ApiMethod: "ModifyInstanceAttribute",
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
				Action: "delete", Entity: graph.Instance.String(), Input: "TerminateInstancesInput", Output: "TerminateInstancesOutput", ApiMethod: "TerminateInstances",
				RequiredParams: []param{
					{AwsField: "InstanceIds", TemplateName: "id", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "start", Entity: graph.Instance.String(), Input: "StartInstancesInput", Output: "StartInstancesOutput", ApiMethod: "StartInstances", OutputExtractor: "aws.StringValue(output.StartingInstances[0].InstanceId)",
				RequiredParams: []param{
					{AwsField: "InstanceIds", TemplateName: "id", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "stop", Entity: graph.Instance.String(), Input: "StopInstancesInput", Output: "StopInstancesOutput", ApiMethod: "StopInstances", OutputExtractor: "aws.StringValue(output.StoppingInstances[0].InstanceId)",
				RequiredParams: []param{
					{AwsField: "InstanceIds", TemplateName: "id", AwsType: "awsstringslice"},
				},
			},
			{
				Action: "check", Entity: graph.Instance.String(), ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "id"},
					{TemplateName: "state"},
					{TemplateName: "timeout"},
				},
			},

			// Security Group
			{
				Action: "create", Entity: graph.SecurityGroup.String(), Input: "CreateSecurityGroupInput", Output: "CreateSecurityGroupOutput", ApiMethod: "CreateSecurityGroup", OutputExtractor: "aws.StringValue(output.GroupId)",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "name", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
					{AwsField: "Description", TemplateName: "description", AwsType: "awsstr"},
				},
			},
			{
				Action: "update", Entity: graph.SecurityGroup.String(), ManualFuncDefinition: true,
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
				Action: "delete", Entity: graph.SecurityGroup.String(), Input: "DeleteSecurityGroupInput", Output: "DeleteSecurityGroupOutput", ApiMethod: "DeleteSecurityGroup",
				RequiredParams: []param{
					{AwsField: "GroupId", TemplateName: "id", AwsType: "awsstr"},
				},
			},

			// VOLUME
			{
				Action: "create", Entity: graph.Volume.String(), Input: "CreateVolumeInput", Output: "Volume", ApiMethod: "CreateVolume", OutputExtractor: "aws.StringValue(output.VolumeId)",
				RequiredParams: []param{
					{AwsField: "AvailabilityZone", TemplateName: "zone", AwsType: "awsstr"},
					{AwsField: "Size", TemplateName: "size", AwsType: "awsint64"},
				},
			},
			{
				Action: "delete", Entity: graph.Volume.String(), Input: "DeleteVolumeInput", Output: "DeleteVolumeOutput", ApiMethod: "DeleteVolume",
				RequiredParams: []param{
					{AwsField: "VolumeId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: graph.Volume.String(), Input: "AttachVolumeInput", Output: "VolumeAttachment", ApiMethod: "AttachVolume", OutputExtractor: "aws.StringValue(output.VolumeId)",
				RequiredParams: []param{
					{AwsField: "Device", TemplateName: "device", AwsType: "awsstr"},
					{AwsField: "VolumeId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "InstanceId", TemplateName: "instance", AwsType: "awsstr"},
				},
			},
			// INTERNET GATEWAYS
			{
				Action: "create", Entity: graph.InternetGateway.String(), Input: "CreateInternetGatewayInput", Output: "CreateInternetGatewayOutput", ApiMethod: "CreateInternetGateway", OutputExtractor: "aws.StringValue(output.InternetGateway.InternetGatewayId)",
			},
			{
				Action: "delete", Entity: graph.InternetGateway.String(), Input: "DeleteInternetGatewayInput", Output: "DeleteInternetGatewayOutput", ApiMethod: "DeleteInternetGateway",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: graph.InternetGateway.String(), Input: "AttachInternetGatewayInput", Output: "AttachInternetGatewayOutput", ApiMethod: "AttachInternetGateway",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: graph.InternetGateway.String(), Input: "DetachInternetGatewayInput", Output: "DetachInternetGatewayOutput", ApiMethod: "DetachInternetGateway",
				RequiredParams: []param{
					{AwsField: "InternetGatewayId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"},
				},
			},
			// ROUTE TABLES
			{
				Action: "create", Entity: graph.RouteTable.String(), Input: "CreateRouteTableInput", Output: "CreateRouteTableOutput", ApiMethod: "CreateRouteTable", OutputExtractor: "aws.StringValue(output.RouteTable.RouteTableId)",
				RequiredParams: []param{
					{AwsField: "VpcId", TemplateName: "vpc", AwsType: "awsstr"}},
			},
			{
				Action: "delete", Entity: graph.RouteTable.String(), Input: "DeleteRouteTableInput", Output: "DeleteRouteTableOutput", ApiMethod: "DeleteRouteTable",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "id", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: graph.RouteTable.String(), Input: "AssociateRouteTableInput", Output: "AssociateRouteTableOutput", ApiMethod: "AssociateRouteTable", OutputExtractor: "aws.StringValue(output.AssociationId)",
				RequiredParams: []param{
					{AwsField: "RouteTableId", TemplateName: "id", AwsType: "awsstr"},
					{AwsField: "SubnetId", TemplateName: "subnet", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: graph.RouteTable.String(), Input: "DisassociateRouteTableInput", Output: "DisassociateRouteTableOutput", ApiMethod: "DisassociateRouteTable",
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
				Action: "create", Entity: graph.Keypair.String(), ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "name"},
				},
			},
			{
				Action: "delete", Entity: graph.Keypair.String(), Input: "DeleteKeyPairInput", Output: "DeleteKeyPairOutput", ApiMethod: "DeleteKeyPair",
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
				Action: "create", Entity: graph.LoadBalancer.String(), Input: "CreateLoadBalancerInput", Output: "CreateLoadBalancerOutput", ApiMethod: "CreateLoadBalancer", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.LoadBalancers[0].LoadBalancerArn)",
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
				Action: "delete", Entity: graph.LoadBalancer.String(), Input: "DeleteLoadBalancerInput", Output: "DeleteLoadBalancerOutput", ApiMethod: "DeleteLoadBalancer", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "LoadBalancerArn", TemplateName: "arn", AwsType: "awsstr"},
				},
			},
			// Listener
			{
				Action: "create", Entity: graph.Listener.String(), Input: "CreateListenerInput", Output: "CreateListenerOutput", ApiMethod: "CreateListener", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.Listeners[0].ListenerArn)",
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
				Action: "delete", Entity: graph.Listener.String(), Input: "DeleteListenerInput", Output: "DeleteListenerOutput", ApiMethod: "DeleteListener", DryRunUnsupported: true,
				RequiredParams: []param{
					{AwsField: "ListenerArn", TemplateName: "arn", AwsType: "awsstr"},
				},
			},
			// Target group
			{
				Action: "create", Entity: graph.TargetGroup.String(), Input: "CreateTargetGroupInput", Output: "CreateTargetGroupOutput", ApiMethod: "CreateTargetGroup", DryRunUnsupported: true, OutputExtractor: "aws.StringValue(output.TargetGroups[0].TargetGroupArn)",
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
				Action: "delete", Entity: graph.TargetGroup.String(), Input: "DeleteTargetGroupInput", Output: "DeleteTargetGroupOutput", ApiMethod: "DeleteTargetGroup", DryRunUnsupported: true,
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
				Action: "create", Entity: graph.User.String(), DryRunUnsupported: true, Input: "CreateUserInput", Output: "CreateUserOutput", ApiMethod: "CreateUser", OutputExtractor: "aws.StringValue(output.User.UserId)",
				RequiredParams: []param{
					{AwsField: "UserName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: graph.User.String(), DryRunUnsupported: true, Input: "DeleteUserInput", Output: "DeleteUserOutput", ApiMethod: "DeleteUser",
				RequiredParams: []param{
					{AwsField: "UserName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "attach", Entity: graph.User.String(), DryRunUnsupported: true, Input: "AddUserToGroupInput", Output: "AddUserToGroupOutput", ApiMethod: "AddUserToGroup",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "group", AwsType: "awsstr"},
					{AwsField: "UserName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "detach", Entity: graph.User.String(), DryRunUnsupported: true, Input: "RemoveUserFromGroupInput", Output: "RemoveUserFromGroupOutput", ApiMethod: "RemoveUserFromGroup",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "group", AwsType: "awsstr"},
					{AwsField: "UserName", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// GROUP
			{
				Action: "create", Entity: graph.Group.String(), DryRunUnsupported: true, Input: "CreateGroupInput", Output: "CreateGroupOutput", ApiMethod: "CreateGroup", OutputExtractor: "aws.StringValue(output.Group.GroupId)",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: graph.Group.String(), DryRunUnsupported: true, Input: "DeleteGroupInput", Output: "DeleteGroupOutput", ApiMethod: "DeleteGroup",
				RequiredParams: []param{
					{AwsField: "GroupName", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// POLICY
			{
				Action: "attach", Entity: graph.Policy.String(), ManualFuncDefinition: true,
				RequiredParams: []param{
					{TemplateName: "arn"},
				},
				ExtraParams: []param{
					{TemplateName: "user"},
					{TemplateName: "group"},
				},
			},
			{
				Action: "detach", Entity: graph.Policy.String(), ManualFuncDefinition: true,
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
				Action: "create", Entity: graph.Bucket.String(), DryRunUnsupported: true, Input: "CreateBucketInput", Output: "CreateBucketOutput", ApiMethod: "CreateBucket", OutputExtractor: "params[\"name\"]",
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: graph.Bucket.String(), DryRunUnsupported: true, Input: "DeleteBucketInput", Output: "DeleteBucketOutput", ApiMethod: "DeleteBucket",
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "name", AwsType: "awsstr"},
				},
			},

			// OBJECT
			{
				Action: "create", Entity: graph.Object.String(), ManualFuncDefinition: true,
				RequiredParams: []param{
					{AwsField: "Bucket", TemplateName: "bucket", AwsType: "awsstr"},
					{AwsField: "Body", TemplateName: "file", AwsType: "awsstr"},
				},
				ExtraParams: []param{
					{AwsField: "Key", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: graph.Object.String(), DryRunUnsupported: true, Input: "DeleteObjectInput", Output: "DeleteObjectOutput", ApiMethod: "DeleteObject",
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
				Action: "create", Entity: graph.Topic.String(), DryRunUnsupported: true, Input: "CreateTopicInput", Output: "CreateTopicOutput", ApiMethod: "CreateTopic", OutputExtractor: "aws.StringValue(output.TopicArn)",
				RequiredParams: []param{
					{AwsField: "Name", TemplateName: "name", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: graph.Topic.String(), DryRunUnsupported: true, Input: "DeleteTopicInput", Output: "DeleteTopicOutput", ApiMethod: "DeleteTopic",
				RequiredParams: []param{
					{AwsField: "TopicArn", TemplateName: "arn", AwsType: "awsstr"},
				},
			},
			//Subscription
			{
				Action: "create", Entity: graph.Subscription.String(), DryRunUnsupported: true, Input: "SubscribeInput", Output: "SubscribeOutput", ApiMethod: "Subscribe", OutputExtractor: "aws.StringValue(output.SubscriptionArn)",
				RequiredParams: []param{
					{AwsField: "TopicArn", TemplateName: "topic", AwsType: "awsstr"},
					{AwsField: "Endpoint", TemplateName: "endpoint", AwsType: "awsstr"},
					{AwsField: "Protocol", TemplateName: "protocol", AwsType: "awsstr"},
				},
			},
			{
				Action: "delete", Entity: graph.Subscription.String(), DryRunUnsupported: true, Input: "UnsubscribeInput", Output: "UnsubscribeOutput", ApiMethod: "Unsubscribe",
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
				Action: "create", Entity: graph.Queue.String(), DryRunUnsupported: true, Input: "CreateQueueInput", Output: "CreateQueueOutput", ApiMethod: "CreateQueue", OutputExtractor: "aws.StringValue(output.QueueUrl)",
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
				Action: "delete", Entity: graph.Queue.String(), DryRunUnsupported: true, Input: "DeleteQueueInput", Output: "DeleteQueueOutput", ApiMethod: "DeleteQueue",
				RequiredParams: []param{
					{AwsField: "QueueUrl", TemplateName: "url", AwsType: "awsstr"},
				},
			},
		},
	},
	{
		Api:          "route53",
		ApiInterface: "Route53API",
		Drivers:      []driver{},
	},
}
