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
	"github.com/wallix/awless/cloud/properties"
)

var DefaultsColumnDefinitions = map[string][]ColumnDefinition{
	//EC2
	cloud.Instance: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.AvailabilityZone, Friendly: "Zone"},
		StringColumnDefinition{Prop: properties.Name},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.State},
			ColoredValues:          map[string]color.Attribute{"running": color.FgGreen, "stopped": color.FgRed},
		},
		StringColumnDefinition{Prop: properties.Type},
		StringColumnDefinition{Prop: properties.PublicIP, Friendly: "Public IP"},
		StringColumnDefinition{Prop: properties.PrivateIP, Friendly: "Private IP"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Launched, Friendly: "Uptime"}},
		StringColumnDefinition{Prop: properties.SSHKey, Friendly: "Access Key"},
	},
	cloud.Vpc: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.Default, Friendly: "Default"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgGreen},
		},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.CIDR},
	},
	cloud.Subnet: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.CIDR},
		StringColumnDefinition{Prop: properties.AvailabilityZone, Friendly: "Zone"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.Default, Friendly: "Default"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgGreen},
		},
		StringColumnDefinition{Prop: properties.Vpc},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.Public},
			ColoredValues:          map[string]color.Attribute{"true": color.FgYellow}},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.State},
			ColoredValues:          map[string]color.Attribute{"available": color.FgGreen}},
	},
	cloud.SecurityGroup: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Vpc},
		FirewallRulesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.InboundRules, Friendly: "Inbound"}},
		FirewallRulesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.OutboundRules, Friendly: "Outbound"}},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Description, DisableTruncate: true},
	},
	cloud.InternetGateway: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Vpcs, DisableTruncate: true},
	},
	cloud.RouteTable: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Vpc},
		StringColumnDefinition{Prop: properties.Main},
		RoutesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Routes}},
	},
	cloud.Keypair: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Fingerprint, DisableTruncate: true},
	},
	cloud.Volume: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Type},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.Size, Friendly: "Size (Gb)"},
		StringColumnDefinition{Prop: properties.Encrypted},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created, Friendly: "Created"}},
		StringColumnDefinition{Prop: properties.AvailabilityZone, Friendly: "Zone"},
	},
	cloud.AvailabilityZone: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.Region},
		StringColumnDefinition{Prop: properties.Messages},
	},
	// Loadbalancer
	cloud.LoadBalancer: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Vpc},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.PublicDNS},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created, Friendly: "Created"}},
		StringColumnDefinition{Prop: properties.Scheme},
	},
	cloud.TargetGroup: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Vpc},
		StringColumnDefinition{Prop: properties.CheckHTTPCode},
		StringColumnDefinition{Prop: properties.Port},
		StringColumnDefinition{Prop: properties.Protocol},
		StringColumnDefinition{Prop: properties.CheckInterval, Friendly: "HCInterval"},
		StringColumnDefinition{Prop: properties.CheckPath, Friendly: "HCPath"},
		StringColumnDefinition{Prop: properties.CheckPort, Friendly: "HCPort"},
		StringColumnDefinition{Prop: properties.CheckProtocol, Friendly: "HCProtocol"},
	},
	cloud.Listener: {
		StringColumnDefinition{Prop: properties.ID, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Actions},
		StringColumnDefinition{Prop: properties.LoadBalancer},
		StringColumnDefinition{Prop: properties.Port},
		StringColumnDefinition{Prop: properties.Protocol},
		StringColumnDefinition{Prop: properties.CipherSuite},
	},
	// Database
	cloud.Database: {
		StringColumnDefinition{Prop: properties.ID, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.AvailabilityZone, Friendly: "Zone"},
		StringColumnDefinition{Prop: properties.Class},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.State},
			ColoredValues:          map[string]color.Attribute{"available": color.FgGreen}},
		StringColumnDefinition{Prop: properties.Storage, Friendly: "Size(Gb)"},
		StringColumnDefinition{Prop: properties.Port},
		StringColumnDefinition{Prop: properties.Username},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.Public, Friendly: "Public"},
			ColoredValues:          map[string]color.Attribute{"true": color.FgYellow}},
		StringColumnDefinition{Prop: properties.Engine},
		StringColumnDefinition{Prop: properties.EngineVersion, Friendly: "Version"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created, Friendly: "Created"}},
	},
	cloud.DbSubnetGroup: {
		StringColumnDefinition{Prop: properties.ID, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.State, Friendly: "Status"},
		StringColumnDefinition{Prop: properties.Vpc},
		StringColumnDefinition{Prop: properties.Subnets, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Description},
	},
	//IAM
	cloud.User: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.PasswordLastUsed, Friendly: "PasswordLastUsed"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.Role: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.Policy: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Updated, Friendly: "Updated"}},
	},
	cloud.Group: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	// S3
	cloud.Bucket: {
		StringColumnDefinition{Prop: properties.ID, DisableTruncate: true},
		GrantsColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Grants}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.Object: {
		StringColumnDefinition{Prop: properties.ID, TruncateRight: true},
		StringColumnDefinition{Prop: properties.Bucket, Friendly: "Bucket"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Modified, Friendly: "Modified"}},
		StringColumnDefinition{Prop: properties.Owner, TruncateRight: true},
		StringColumnDefinition{Prop: properties.Size},
		StringColumnDefinition{Prop: properties.Class},
	},
	//Notification
	cloud.Subscription: {
		StringColumnDefinition{Prop: properties.Arn},
		StringColumnDefinition{Prop: properties.Topic},
		StringColumnDefinition{Prop: properties.Endpoint, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Protocol},
		StringColumnDefinition{Prop: properties.Owner},
	},
	cloud.Topic: {
		StringColumnDefinition{Prop: properties.Topic, DisableTruncate: true},
	},
	//Queue
	cloud.Queue: {
		StringColumnDefinition{Prop: properties.ID, Friendly: "URL", DisableTruncate: true},
		StringColumnDefinition{Prop: properties.ApproximateMessageCount, Friendly: "~NbMsg"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Modified, Friendly: "LastModif"}},
		StringColumnDefinition{Prop: properties.Delay, Friendly: "Delay(s)"},
	},
	// DNS
	cloud.Zone: {
		StringColumnDefinition{Prop: properties.ID, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Comment},
		StringColumnDefinition{Prop: properties.Private, Friendly: "Private"},
		StringColumnDefinition{Prop: properties.RecordCount, Friendly: "Nb Records"},
		StringColumnDefinition{Prop: properties.CallerReference, DisableTruncate: true},
	},
	cloud.Record: {
		StringColumnDefinition{Prop: properties.ID, Friendly: "AwlessId", DisableTruncate: true},
		StringColumnDefinition{Prop: properties.Type},
		StringColumnDefinition{Prop: properties.Name, DisableTruncate: true},
		SliceColumnDefinition{StringColumnDefinition{Prop: properties.Records}},
		StringColumnDefinition{Prop: properties.TTL},
	},
}
