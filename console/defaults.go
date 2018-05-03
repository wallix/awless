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

var ColumnsInListing = map[string][]string{
	cloud.Instance:            {properties.ID, properties.AvailabilityZone, properties.Name, properties.State, properties.Type, properties.PublicIP, properties.PrivateIP, properties.Launched, properties.KeyPair},
	cloud.Vpc:                 {properties.ID, properties.Name, properties.Default, properties.State, properties.CIDR},
	cloud.Subnet:              {properties.ID, properties.Name, properties.CIDR, properties.AvailabilityZone, properties.Default, properties.Vpc, properties.Public, properties.State},
	cloud.SecurityGroup:       {properties.ID, properties.Vpc, properties.InboundRules, properties.OutboundRules, properties.Name, properties.Description},
	cloud.InternetGateway:     {properties.ID, properties.Name, properties.Vpcs},
	cloud.NatGateway:          {properties.ID, properties.State, properties.Vpc, properties.Subnet, properties.Created},
	cloud.RouteTable:          {properties.ID, properties.Name, properties.Vpc, properties.Default, properties.Routes, properties.Associations},
	cloud.Keypair:             {properties.ID, properties.Fingerprint},
	cloud.Image:               {properties.ID, properties.Name, properties.State, properties.Location, properties.Public, properties.Type, properties.Created, properties.Architecture, properties.Hypervisor, properties.Virtualization},
	cloud.ImportImageTask:     {properties.ID, properties.Description, properties.Image, properties.Progress, properties.State, properties.StateMessage},
	cloud.Volume:              {properties.ID, properties.Name, properties.Type, properties.State, properties.Size, properties.Encrypted, properties.Created, properties.AvailabilityZone, properties.Instances},
	cloud.AvailabilityZone:    {properties.Name, properties.State, properties.Region, properties.Messages},
	cloud.ElasticIP:           {properties.ID, properties.PublicIP, properties.PrivateIP, properties.Association},
	cloud.Snapshot:            {properties.ID, properties.Volume, properties.Encrypted, properties.Owner, properties.State, properties.Progress, properties.Created, properties.Size},
	cloud.NetworkInterface:    {properties.ID, properties.Vpc, properties.Subnet, properties.State, properties.Instance, properties.PrivateIP, properties.PublicIP, properties.Description},
	cloud.LoadBalancer:        {properties.Name, properties.Vpc, properties.State, properties.PublicDNS, properties.Created, properties.Scheme},
	cloud.ClassicLoadBalancer: {properties.Name, properties.Vpc, properties.PublicDNS, properties.Instances, properties.Ports, properties.Created, properties.Scheme},
	cloud.TargetGroup:         {properties.Name, properties.Vpc, properties.CheckHTTPCode, properties.Port, properties.Protocol, properties.CheckInterval, properties.CheckPath, properties.CheckPort, properties.CheckProtocol},
	cloud.Listener:            {properties.ID, properties.Protocol, properties.Port, properties.LoadBalancer, properties.TargetGroups, properties.AlarmActions},
	cloud.Database:            {properties.ID, properties.Name, properties.AvailabilityZone, properties.Class, properties.State, properties.Storage, properties.Port, properties.Username, properties.Public, properties.ReplicaOf, properties.Engine, properties.EngineVersion, properties.Created},
	cloud.DbSubnetGroup:       {properties.ID, properties.State, properties.Vpc, properties.Subnets, properties.Description},
	cloud.LaunchConfiguration: {properties.Name, properties.Type, properties.Created, properties.KeyPair},
	cloud.ScalingGroup:        {properties.Name, properties.LaunchConfigurationName, properties.DesiredCapacity, properties.State, properties.Created, properties.NewInstancesProtected},
	cloud.ScalingPolicy:       {properties.Name, properties.Type, properties.ScalingGroupName, properties.AlarmNames, properties.AdjustmentType, properties.ScalingAdjustment},
	cloud.Repository:          {properties.Name, properties.URI, properties.Created, properties.Account, properties.Arn},
	cloud.ContainerCluster:    {properties.Name, properties.State, properties.ActiveServicesCount, properties.PendingTasksCount, properties.RegisteredContainerInstancesCount, properties.RunningTasksCount},
	cloud.ContainerTask:       {properties.Name, properties.Version, properties.State, properties.ContainersImages, properties.Deployments},
	cloud.Container:           {properties.Name, properties.DeploymentName, properties.State, properties.Created, properties.Launched, properties.Stopped, properties.Cluster, properties.ContainerTask},
	cloud.ContainerInstance:   {properties.ID, properties.Instance, properties.Cluster, properties.State, properties.RunningTasksCount, properties.PendingTasksCount, properties.Created, properties.AgentConnected},
	cloud.Certificate:         {properties.Arn, properties.Name},
	cloud.User:                {properties.ID, properties.Name, properties.PasswordLastUsed, properties.Created},
	cloud.Role:                {properties.ID, properties.Name, properties.Created},
	cloud.InstanceProfile:     {properties.ID, properties.Name, properties.Path, properties.Created},
	cloud.Policy:              {properties.ID, properties.Name, properties.Type, properties.Created, properties.Updated, properties.Attached},
	cloud.Group:               {properties.ID, properties.Name, properties.Created},
	cloud.AccessKey:           {properties.ID, properties.State, properties.Username, properties.Created},
	cloud.MFADevice:           {properties.ID, properties.AttachedAt},
	cloud.Bucket:              {properties.ID, properties.Grants, properties.Created},
	cloud.S3Object:            {properties.ID, properties.Bucket, properties.Modified, properties.Owner, properties.Size, properties.Class},
	cloud.Subscription:        {properties.Arn, properties.Topic, properties.Endpoint, properties.Protocol, properties.Owner},
	cloud.Topic:               {properties.ID},
	cloud.Queue:               {properties.ID, properties.ApproximateMessageCount, properties.Created, properties.Modified, properties.Delay},
	cloud.Zone:                {properties.ID, properties.Name, properties.Comment, properties.Private, properties.RecordCount, properties.CallerReference},
	cloud.Record:              {properties.ID, properties.Type, properties.Name, properties.Records, properties.Zone, properties.Alias, properties.TTL},
	cloud.Function:            {properties.Name, properties.Size, properties.Memory, properties.Runtime, properties.Version, properties.Modified, properties.Description},
	cloud.Metric:              {properties.ID, properties.Name, properties.Namespace, properties.Dimensions},
	cloud.Alarm:               {properties.Name, properties.Namespace, properties.MetricName, properties.Description, properties.State, properties.Updated, properties.Dimensions},
	cloud.Distribution:        {properties.ID, properties.PublicDNS, properties.Enabled, properties.State, properties.Modified, properties.Aliases, properties.SSLSupportMethod, properties.Origins},
	cloud.Stack:               {properties.ID, properties.Name, properties.State, properties.Created, properties.Modified},
}

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
		StringColumnDefinition{Prop: properties.KeyPair},
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
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Description},
	},
	cloud.InternetGateway: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Vpcs},
	},
	cloud.NatGateway: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.Vpc},
		StringColumnDefinition{Prop: properties.Subnet},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created, Friendly: "Created"}},
	},
	cloud.RouteTable: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Vpc},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.Default},
			ColoredValues:          map[string]color.Attribute{"true": color.FgGreen},
		},
		RoutesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Routes}},
		KeyValuesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Associations}},
	},
	cloud.Keypair: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Fingerprint},
	},
	cloud.Image: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.Location},
		StringColumnDefinition{Prop: properties.Public},
		StringColumnDefinition{Prop: properties.Type},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created, Friendly: "Created"}},
		StringColumnDefinition{Prop: properties.Architecture, Friendly: "Arch"},
		StringColumnDefinition{Prop: properties.Hypervisor, Friendly: "Hyperv"},
		StringColumnDefinition{Prop: properties.Virtualization, Friendly: "Virt"},
	},
	cloud.ImportImageTask: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Description},
		StringColumnDefinition{Prop: properties.Image},
		StringColumnDefinition{Prop: properties.Progress},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.StateMessage},
	},
	cloud.Volume: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Type},
		StringColumnDefinition{Prop: properties.State},
		StorageColumnDefinition{Unit: gb, StringColumnDefinition: StringColumnDefinition{Prop: properties.Size}},
		StringColumnDefinition{Prop: properties.Encrypted},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		StringColumnDefinition{Prop: properties.AvailabilityZone, Friendly: "Zone"},
		StringColumnDefinition{Prop: properties.Instances},
	},
	cloud.AvailabilityZone: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.Region},
		StringColumnDefinition{Prop: properties.Messages},
	},
	cloud.ElasticIP: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.PublicIP},
		StringColumnDefinition{Prop: properties.PrivateIP},
		StringColumnDefinition{Prop: properties.Association},
	},
	cloud.Snapshot: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Volume},
		StringColumnDefinition{Prop: properties.Encrypted},
		StringColumnDefinition{Prop: properties.Owner},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.Progress},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		StorageColumnDefinition{Unit: gb, StringColumnDefinition: StringColumnDefinition{Prop: properties.Size}},
	},
	cloud.NetworkInterface: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Vpc},
		StringColumnDefinition{Prop: properties.Subnet},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.Instance},
		StringColumnDefinition{Prop: properties.PrivateIP},
		StringColumnDefinition{Prop: properties.PublicIP},
		StringColumnDefinition{Prop: properties.Description},
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
	cloud.ClassicLoadBalancer: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Vpc},
		StringColumnDefinition{Prop: properties.PublicDNS},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created, Friendly: "Created"}},
		StringColumnDefinition{Prop: properties.Scheme},
		SliceColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Instances}},
		SliceColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Ports}},
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
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Port},
		StringColumnDefinition{Prop: properties.Protocol},
		StringColumnDefinition{Prop: properties.LoadBalancer},
		SliceColumnDefinition{ForEach: ToShortArn, StringColumnDefinition: StringColumnDefinition{Prop: properties.TargetGroups}},
		StringColumnDefinition{Prop: properties.AlarmActions},
	},
	// Database
	cloud.Database: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.AvailabilityZone, Friendly: "Zone"},
		StringColumnDefinition{Prop: properties.Zone, Friendly: "HostedZone"},
		StringColumnDefinition{Prop: properties.Class},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.State},
			ColoredValues:          map[string]color.Attribute{"available": color.FgGreen}},
		StorageColumnDefinition{Unit: gb, StringColumnDefinition: StringColumnDefinition{Prop: properties.Storage}},
		StringColumnDefinition{Prop: properties.Port},
		StringColumnDefinition{Prop: properties.Username},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.Public},
			ColoredValues:          map[string]color.Attribute{"true": color.FgYellow}},
		StringColumnDefinition{Prop: properties.ReplicaOf, Friendly: "Replica Of"},
		StringColumnDefinition{Prop: properties.Engine},
		StringColumnDefinition{Prop: properties.EngineVersion, Friendly: "Version"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created, Friendly: "Created"}},
	},
	cloud.DbSubnetGroup: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.State, Friendly: "Status"},
		StringColumnDefinition{Prop: properties.Vpc},
		StringColumnDefinition{Prop: properties.Subnets},
		StringColumnDefinition{Prop: properties.Description},
	},
	//Autoscaling
	cloud.LaunchConfiguration: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Type},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		StringColumnDefinition{Prop: properties.KeyPair},
	},
	cloud.ScalingGroup: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.LaunchConfigurationName, Friendly: "LaunchConfiguration"},
		StringColumnDefinition{Prop: properties.DesiredCapacity},
		StringColumnDefinition{Prop: properties.State},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		StringColumnDefinition{Prop: properties.NewInstancesProtected},
	},
	cloud.ScalingPolicy: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Type},
		StringColumnDefinition{Prop: properties.ScalingGroupName},
		StringColumnDefinition{Prop: properties.AlarmNames},
		StringColumnDefinition{Prop: properties.AdjustmentType},
		StringColumnDefinition{Prop: properties.ScalingAdjustment},
	},
	//Containers
	cloud.Repository: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.URI},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		StringColumnDefinition{Prop: properties.Account},
		StringColumnDefinition{Prop: properties.Arn},
	},
	cloud.ContainerCluster: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.ActiveServicesCount, Friendly: "ActiveServices"},
		StringColumnDefinition{Prop: properties.PendingTasksCount, Friendly: "PendingTasks"},
		StringColumnDefinition{Prop: properties.RegisteredContainerInstancesCount, Friendly: "RegisteredContainerInstances"},
		StringColumnDefinition{Prop: properties.RunningTasksCount, Friendly: "RunningTasks"},
	},
	cloud.ContainerTask: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Version},
		StringColumnDefinition{Prop: properties.State},
		KeyValuesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.ContainersImages}},
		KeyValuesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Deployments}},
	},
	cloud.Container: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.DeploymentName},
		StringColumnDefinition{Prop: properties.State},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Launched}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Stopped}},
		ARNLastValueColumnDefinition{Separator: "/", StringColumnDefinition: StringColumnDefinition{Prop: properties.Cluster}},
		ARNLastValueColumnDefinition{Separator: "/", StringColumnDefinition: StringColumnDefinition{Prop: properties.ContainerTask}},
	},
	cloud.ContainerInstance: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Instance},
		ARNLastValueColumnDefinition{Separator: "/", StringColumnDefinition: StringColumnDefinition{Prop: properties.Cluster}},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.RunningTasksCount, Friendly: "RunningTasks"},
		StringColumnDefinition{Prop: properties.PendingTasksCount, Friendly: "PendingTasks"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		StringColumnDefinition{Prop: properties.AgentConnected},
	},
	//ACM
	cloud.Certificate: {
		StringColumnDefinition{Prop: properties.Arn},
		StringColumnDefinition{Prop: properties.Name},
	},
	//IAM
	cloud.User: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.PasswordLastUsed, Friendly: "PasswordLastUsed"}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.Role: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.InstanceProfile: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Path},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.Policy: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Type},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Updated, Friendly: "Updated"}},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: properties.Attached},
			ColoredValues:          map[string]color.Attribute{"false": color.FgYellow},
		},
	},
	cloud.Group: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.AccessKey: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.State},
		StringColumnDefinition{Prop: properties.Username},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.MFADevice: {
		StringColumnDefinition{Prop: properties.ID},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.AttachedAt}},
	},
	// S3
	cloud.Bucket: {
		StringColumnDefinition{Prop: properties.ID},
		GrantsColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Grants}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
	},
	cloud.S3Object: {
		StringColumnDefinition{Prop: properties.ID, Friendly: "Name"},
		StringColumnDefinition{Prop: properties.Bucket, Friendly: "Bucket"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Modified, Friendly: "Modified"}},
		StringColumnDefinition{Prop: properties.Owner},
		StorageColumnDefinition{Unit: b, StringColumnDefinition: StringColumnDefinition{Prop: properties.Size}},
		StringColumnDefinition{Prop: properties.Class},
	},
	//Notification
	cloud.Subscription: {
		StringColumnDefinition{Prop: properties.Arn},
		StringColumnDefinition{Prop: properties.Topic},
		StringColumnDefinition{Prop: properties.Endpoint},
		StringColumnDefinition{Prop: properties.Protocol},
		StringColumnDefinition{Prop: properties.Owner},
	},
	cloud.Topic: {
		StringColumnDefinition{Prop: properties.ID},
	},
	//Queue
	cloud.Queue: {
		StringColumnDefinition{Prop: properties.ID, Friendly: "URL"},
		StringColumnDefinition{Prop: properties.ApproximateMessageCount, Friendly: "~NbMsg"},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Modified, Friendly: "LastModif"}},
		StringColumnDefinition{Prop: properties.Delay, Friendly: "Delay(s)"},
	},
	// DNS
	cloud.Zone: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Comment},
		StringColumnDefinition{Prop: properties.Private, Friendly: "Private"},
		StringColumnDefinition{Prop: properties.RecordCount, Friendly: "Nb Records"},
		StringColumnDefinition{Prop: properties.CallerReference},
	},
	cloud.Record: {
		StringColumnDefinition{Prop: properties.ID, Friendly: "AwlessId"},
		StringColumnDefinition{Prop: properties.Type},
		StringColumnDefinition{Prop: properties.Name},
		SliceColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Records}},
		StringColumnDefinition{Prop: properties.Zone},
		StringColumnDefinition{Prop: properties.Alias},
		StringColumnDefinition{Prop: properties.TTL},
	},
	// Lamba
	cloud.Function: {
		StringColumnDefinition{Prop: properties.Name},
		StorageColumnDefinition{Unit: b, StringColumnDefinition: StringColumnDefinition{Prop: properties.Size}},
		StorageColumnDefinition{Unit: mb, StringColumnDefinition: StringColumnDefinition{Prop: properties.Memory}},
		StringColumnDefinition{Prop: properties.Runtime},
		StringColumnDefinition{Prop: properties.Version},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Modified}},
		StringColumnDefinition{Prop: properties.Description},
	},
	//Monitoring
	cloud.Metric: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Namespace},
		KeyValuesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Dimensions}},
	},
	cloud.Alarm: {
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.Namespace},
		StringColumnDefinition{Prop: properties.MetricName},
		StringColumnDefinition{Prop: properties.Description},
		StringColumnDefinition{Prop: properties.State},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Updated}},
		KeyValuesColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Dimensions}},
	},
	//CDN
	cloud.Distribution: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.PublicDNS},
		StringColumnDefinition{Prop: properties.Enabled},
		StringColumnDefinition{Prop: properties.State},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Modified}},
		SliceColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Aliases}},
		StringColumnDefinition{Prop: properties.SSLSupportMethod},
		SliceColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Origins}},
	},
	//Cloudformation
	cloud.Stack: {
		StringColumnDefinition{Prop: properties.ID},
		StringColumnDefinition{Prop: properties.Name},
		StringColumnDefinition{Prop: properties.State},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Created}},
		TimeColumnDefinition{StringColumnDefinition: StringColumnDefinition{Prop: properties.Modified}},
	},
}
