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

package console

import (
	"github.com/fatih/color"
	"github.com/wallix/awless/cloud"
)

var DefaultsColumnDefinitions = map[string][]ColumnDefinition{
	//EC2
	cloud.Instance: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "SubnetId"},
		StringColumnDefinition{Prop: "Name"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "State"},
			ColoredValues:          map[string]color.Attribute{"running": color.FgGreen, "stopped": color.FgRed},
		},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "KeyName", Friendly: "Access Key"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "LaunchTime"}},
	},
	cloud.Vpc: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "IsDefault", Friendly: "Default"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgGreen},
		},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "CidrBlock"},
	},
	cloud.Subnet: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "CidrBlock"},
		StringColumnDefinition{Prop: "AvailabilityZone", Friendly: "Zone"},
		StringColumnDefinition{Prop: "VpcId"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "MapPublicIpOnLaunch", Friendly: "Public VMs"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgYellow}},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "State"},
			ColoredValues:          map[string]color.Attribute{"available": color.FgGreen}},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "DefaultForAz", Friendly: "ZoneDefault"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgGreen},
		},
	},
	cloud.SecurityGroup: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "VpcId"},
		FirewallRulesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "InboundRules", Friendly: "Inbound"}},
		FirewallRulesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "OutboundRules", Friendly: "Outbound"}},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "Description", DisableTruncate: true},
	},
	cloud.InternetGateway: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "Vpcs", DisableTruncate: true},
	},
	cloud.RouteTable: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "VpcId"},
		StringColumnDefinition{Prop: "Main"},
		RoutesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "Routes"}},
	},
	cloud.Keypair: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "KeyFingerprint", DisableTruncate: true},
	},
	cloud.Volume: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "VolumeType"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "Size", Friendly: "Size (Gb)"},
		StringColumnDefinition{Prop: "Encrypted"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateTime"}},
		StringColumnDefinition{Prop: "AvailabilityZone"},
	},
	cloud.AvailabilityZone: {
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "Region"},
		StringColumnDefinition{Prop: "Messages"},
	},
	// Loadbalancer
	cloud.LoadBalancer: {
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "VpcId"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "DNSName"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateTime"}},
		StringColumnDefinition{Prop: "Scheme"},
	},
	cloud.TargetGroup: {
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "VpcId"},
		StringColumnDefinition{Prop: "Matcher"},
		StringColumnDefinition{Prop: "Port"},
		StringColumnDefinition{Prop: "Protocol"},
		StringColumnDefinition{Prop: "HealthCheckIntervalSeconds", Friendly: "HCInterval"},
		StringColumnDefinition{Prop: "HealthCheckPath", Friendly: "HCPath"},
		StringColumnDefinition{Prop: "HealthCheckPort", Friendly: "HCPort"},
		StringColumnDefinition{Prop: "HealthCheckProtocol", Friendly: "HCProtocol"},
	},
	cloud.Listener: {
		StringColumnDefinition{Prop: "Id", DisableTruncate: true},
		StringColumnDefinition{Prop: "Actions"},
		StringColumnDefinition{Prop: "LoadBalancer"},
		StringColumnDefinition{Prop: "Port"},
		StringColumnDefinition{Prop: "Protocol"},
		StringColumnDefinition{Prop: "SslPolicy"},
	},
	//IAM
	cloud.User: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "PasswordLastUsedDate", Friendly: "PasswordLastUsed"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
	},
	cloud.Role: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
	},
	cloud.Policy: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "UpdateDate"}},
	},
	cloud.Group: {
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
	},
	// S3
	cloud.Bucket: {
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		GrantsColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "Grants"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreateDate"}},
	},
	cloud.Object: {
		StringColumnDefinition{Prop: "Key", TruncateRight: true},
		StringColumnDefinition{Prop: "BucketName"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "ModifiedDate"}},
		StringColumnDefinition{Prop: "OwnerId", TruncateRight: true},
		StringColumnDefinition{Prop: "Size"},
		StringColumnDefinition{Prop: "Class"},
	},
	//Notification
	cloud.Subscription: {
		StringColumnDefinition{Prop: "SubscriptionArn"},
		StringColumnDefinition{Prop: "TopicArn"},
		StringColumnDefinition{Prop: "Endpoint", DisableTruncate: true},
		StringColumnDefinition{Prop: "Protocol"},
		StringColumnDefinition{Prop: "Owner"},
	},
	cloud.Topic: {
		StringColumnDefinition{Prop: "TopicArn", DisableTruncate: true},
	},
	//Queue
	cloud.Queue: {
		StringColumnDefinition{Prop: "Id", Friendly: "URL", DisableTruncate: true},
		StringColumnDefinition{Prop: "ApproximateNumberOfMessages", Friendly: "~NbMsg"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "CreatedTimestamp", Friendly: "Created"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: "LastModifiedTimestamp", Friendly: "LastModif"}},
		StringColumnDefinition{Prop: "DelaySeconds", Friendly: "Delay(s)"},
	},
	// DNS
	cloud.Zone: {
		StringColumnDefinition{Prop: "Id", DisableTruncate: true},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		StringColumnDefinition{Prop: "Comment"},
		StringColumnDefinition{Prop: "IsPrivateZone"},
		StringColumnDefinition{Prop: "ResourceRecordSetCount"},
		StringColumnDefinition{Prop: "CallerReference", DisableTruncate: true},
	},
	cloud.Record: {
		StringColumnDefinition{Prop: "Id", Friendly: "AwlessId", DisableTruncate: true},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "Name", DisableTruncate: true},
		SliceColumnDefinition{StringColumnDefinition{Prop: "Records"}},
		StringColumnDefinition{Prop: "TTL"},
	},
}
