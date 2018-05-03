/* Copyright 2017 WALLIX

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

// DO NOT EDIT
// This file was automatically generated with go generate
package rdf

import "github.com/wallix/awless/cloud/properties"

const (
	Account                           = "cloud:account"
	ACMCertificate                    = "cloud:acmCertificate"
	Actions                           = "cloud:actions"
	ActionsEnabled                    = "cloud:actionsEnabled"
	ActiveServicesCount               = "cloud:activeServicesCount"
	AdjustmentType                    = "cloud:adjustmentType"
	Affinity                          = "cloud:affinity"
	AgentConnected                    = "cloud:agentConnected"
	AgentState                        = "cloud:agentState"
	AgentVersion                      = "cloud:agentVersion"
	AlarmActions                      = "cloud:alarmActions"
	AlarmNames                        = "cloud:alarmNames"
	Alias                             = "cloud:alias"
	Aliases                           = "cloud:aliases"
	ApproximateMessageCount           = "cloud:approximateMessageCount"
	Architecture                      = "cloud:architecture"
	Arn                               = "cloud:arn"
	Association                       = "cloud:association"
	Associations                      = "cloud:associations"
	Attachable                        = "cloud:attachable"
	Attached                          = "cloud:attached"
	AttachedAt                        = "cloud:attachedAt"
	Attachment                        = "cloud:attachment"
	Attributes                        = "cloud:attributes"
	AutoUpgrade                       = "cloud:autoUpgrade"
	AvailabilityZone                  = "cloud:availabilityZone"
	AvailabilityZones                 = "cloud:availabilityZones"
	BackupRetentionPeriod             = "cloud:backupRetentionPeriod"
	Bucket                            = "cloud:bucketName"
	CallerReference                   = "cloud:callerReference"
	Capabilities                      = "cloud:capabilities"
	Certificate                       = "cloud:certificate"
	CertificateAuthority              = "cloud:certificateAuthority"
	Certificates                      = "cloud:certificates"
	ChangeSet                         = "cloud:changeSet"
	Charset                           = "cloud:charset"
	CheckHTTPCode                     = "cloud:checkHTTPCode"
	CheckInterval                     = "cloud:checkInterval"
	CheckPath                         = "cloud:checkPath"
	CheckPort                         = "cloud:checkPort"
	CheckProtocol                     = "cloud:checkProtocol"
	CheckTimeout                      = "cloud:checkTimeout"
	CIDR                              = "net:cidr"
	CIDRv6                            = "net:cidrv6"
	CipherSuite                       = "cloud:cipherSuite"
	Class                             = "cloud:class"
	Cluster                           = "cloud:cluster"
	Comment                           = "rdfs:comment"
	Config                            = "cloud:config"
	ContainerInstance                 = "cloud:containerInstance"
	ContainersImages                  = "cloud:containersImages"
	ContainerTask                     = "cloud:containerTask"
	Continent                         = "cloud:continent"
	Cooldown                          = "cloud:cooldown"
	CopyTagsToSnapshot                = "cloud:copyTagsToSnapshot"
	Country                           = "cloud:country"
	Created                           = "cloud:created"
	DBSecurityGroups                  = "cloud:dbSecurityGroups"
	DBSubnetGroup                     = "cloud:dbSubnetGroup"
	Default                           = "cloud:default"
	DefaultCooldown                   = "cloud:defaultCooldown"
	Delay                             = "cloud:delaySeconds"
	DeploymentName                    = "cloud:deploymentName"
	Deployments                       = "cloud:deployments"
	Description                       = "cloud:description"
	DesiredCapacity                   = "cloud:desiredCapacity"
	Dimensions                        = "cloud:dimensions"
	DisableRollback                   = "cloud:disableRollback"
	DockerVersion                     = "cloud:dockerVersion"
	Document                          = "cloud:document"
	Enabled                           = "cloud:enabled"
	Encrypted                         = "cloud:encrypted"
	Endpoint                          = "cloud:endpoint"
	Engine                            = "cloud:engine"
	EngineVersion                     = "cloud:engineVersion"
	ExitCode                          = "cloud:exitCode"
	Failover                          = "cloud:failover"
	Fingerprint                       = "cloud:fingerprint"
	GlobalID                          = "cloud:globalID"
	GranteeType                       = "cloud:granteeType"
	Grants                            = "cloud:grants"
	Handler                           = "cloud:handler"
	Hash                              = "cloud:hash"
	HealthCheck                       = "cloud:healthCheck"
	HealthCheckGracePeriod            = "cloud:healthCheckGracePeriod"
	HealthCheckType                   = "cloud:healthCheckType"
	HealthyThresholdCount             = "cloud:healthyThresholdCount"
	Host                              = "cloud:host"
	HTTPVersion                       = "cloud:httpVersion"
	Hypervisor                        = "cloud:hypervisor"
	ID                                = "cloud:id"
	Image                             = "cloud:image"
	InboundRules                      = "net:inboundRules"
	InlinePolicies                    = "cloud:inlinePolicies"
	Instance                          = "cloud:instance"
	InstanceOwner                     = "cloud:instanceOwner"
	Instances                         = "cloud:instances"
	InsufficientDataActions           = "cloud:insufficientDataActions"
	IOPS                              = "cloud:iops"
	IPType                            = "net:ipType"
	IPv6Addresses                     = "cloud:ipv6Addresses"
	IPv6Enabled                       = "cloud:ipv6Enabled"
	Key                               = "cloud:key"
	KeyName                           = "cloud:keyName"
	KeyPair                           = "cloud:keyPair"
	LatestRestorableTime              = "cloud:latestRestorableTime"
	LaunchConfigurationName           = "cloud:launchConfigurationName"
	Launched                          = "cloud:launched"
	License                           = "cloud:license"
	Lifecycle                         = "cloud:lifecycle"
	LoadBalancer                      = "cloud:loadBalancer"
	Location                          = "cloud:location"
	MACAddress                        = "cloud:macAddress"
	Main                              = "cloud:main"
	MaxSize                           = "cloud:maxSize"
	Memory                            = "cloud:memory"
	Messages                          = "cloud:messages"
	MetricName                        = "cloud:metricName"
	MinSize                           = "cloud:minSize"
	Modified                          = "cloud:modified"
	MonitoringInterval                = "cloud:monitoringInterval"
	MonitoringRole                    = "cloud:monitoringRole"
	MultiAZ                           = "cloud:multiAZ"
	Name                              = "cloud:name"
	Namespace                         = "cloud:namemespace"
	NetworkInterfaces                 = "cloud:networkInterfaces"
	NewInstancesProtected             = "cloud:newInstancesProtected"
	Notifications                     = "cloud:notifications"
	OKActions                         = "cloud:okActions"
	OptionGroups                      = "cloud:optionGroups"
	Origins                           = "cloud:origins"
	OutboundRules                     = "net:outboundRules"
	Outputs                           = "cloud:outputs"
	Owner                             = "cloud:owner"
	ParameterGroups                   = "cloud:parameterGroups"
	Parameters                        = "cloud:parameters"
	PasswordLastUsed                  = "cloud:passwordLastUsed"
	Path                              = "cloud:path"
	PathPrefix                        = "cloud:pathPrefix"
	PendingTasksCount                 = "cloud:pendingTasksCount"
	PlacementGroup                    = "cloud:placementGroup"
	Port                              = "net:port"
	Ports                             = "net:ports"
	PortRange                         = "net:portRange"
	PreferredBackupDate               = "cloud:preferredBackupDate"
	PreferredMaintenanceDate          = "cloud:preferredMaintenanceDate"
	PriceClass                        = "cloud:priceClass"
	Private                           = "cloud:private"
	PrivateDNS                        = "cloud:privateDNS"
	PrivateIP                         = "net:privateIP"
	Profile                           = "cloud:profile"
	Progress                          = "cloud:progress"
	Protocol                          = "net:protocol"
	Public                            = "cloud:public"
	PublicDNS                         = "cloud:publicDNS"
	PublicIP                          = "net:publicIP"
	RecordCount                       = "cloud:records"
	Records                           = "cloud:recordCount"
	Region                            = "cloud:region"
	RegisteredContainerInstancesCount = "cloud:registeredContainerInstancesCount"
	ReplicaOf                         = "cloud:replicaOf"
	Role                              = "cloud:role"
	Roles                             = "cloud:roles"
	RootDevice                        = "cloud:rootDevice"
	RootDeviceType                    = "cloud:rootDeviceType"
	Routes                            = "net:routes"
	RunningTasksCount                 = "cloud:runningTasksCount"
	Runtime                           = "cloud:runtime"
	ScalingAdjustment                 = "cloud:scalingAdjustment"
	ScalingGroupName                  = "cloud:scalingGroupName"
	Scheme                            = "net:scheme"
	SecondaryAvailabilityZone         = "cloud:secondaryAvailabilityZone"
	SecurityGroups                    = "cloud:securityGroups"
	Set                               = "cloud:set"
	Size                              = "cloud:size"
	Source                            = "cloud:source"
	SpotInstanceRequestId             = "cloud:spotInstanceRequestId"
	SpotPrice                         = "cloud:spotPrice"
	SSLSupportMethod                  = "cloud:sslSupportMethod"
	State                             = "cloud:state"
	StateMessage                      = "cloud:stateMessage"
	Stopped                           = "cloud:stopped"
	Storage                           = "cloud:storage"
	StorageType                       = "cloud:storageType"
	Subnet                            = "cloud:subnet"
	Subnets                           = "cloud:subnets"
	Tags                              = "cloud:tags"
	TargetGroups                      = "cloud:targetGroups"
	Timeout                           = "cloud:timezone"
	Timezone                          = "cloud:timeout"
	TLSVersionRequired                = "cloud:tlsVersionRequired"
	Topic                             = "cloud:topic"
	TrafficPolicyInstance             = "cloud:trafficPolicyInstance"
	TrustPolicy                       = "cloud:trustPolicy"
	TTL                               = "cloud:ttl"
	Type                              = "cloud:type"
	UnhealthyThresholdCount           = "cloud:unhealthyThresholdCount"
	Updated                           = "cloud:updated"
	URI                               = "cloud:uri"
	UserData                          = "cloud:userData"
	Username                          = "cloud:username"
	Value                             = "cloud:value"
	Version                           = "cloud:version"
	Virtualization                    = "cloud:virtualization"
	Volume                            = "cloud:volume"
	Vpc                               = "cloud:vpc"
	Vpcs                              = "cloud:vpcs"
	WebACL                            = "cloud:webACL"
	Weight                            = "cloud:weight"
	Zone                              = "cloud:zone"
)

func init() {
	Labels = map[string]string{
		properties.Account:                           Account,
		properties.ACMCertificate:                    ACMCertificate,
		properties.Actions:                           Actions,
		properties.ActionsEnabled:                    ActionsEnabled,
		properties.ActiveServicesCount:               ActiveServicesCount,
		properties.AdjustmentType:                    AdjustmentType,
		properties.Affinity:                          Affinity,
		properties.AgentConnected:                    AgentConnected,
		properties.AgentState:                        AgentState,
		properties.AgentVersion:                      AgentVersion,
		properties.AlarmActions:                      AlarmActions,
		properties.AlarmNames:                        AlarmNames,
		properties.Alias:                             Alias,
		properties.Aliases:                           Aliases,
		properties.ApproximateMessageCount:           ApproximateMessageCount,
		properties.Architecture:                      Architecture,
		properties.Arn:                               Arn,
		properties.Association:                       Association,
		properties.Associations:                      Associations,
		properties.Attachable:                        Attachable,
		properties.Attached:                          Attached,
		properties.AttachedAt:                        AttachedAt,
		properties.Attachment:                        Attachment,
		properties.Attributes:                        Attributes,
		properties.AutoUpgrade:                       AutoUpgrade,
		properties.AvailabilityZone:                  AvailabilityZone,
		properties.AvailabilityZones:                 AvailabilityZones,
		properties.BackupRetentionPeriod:             BackupRetentionPeriod,
		properties.Bucket:                            Bucket,
		properties.CallerReference:                   CallerReference,
		properties.Capabilities:                      Capabilities,
		properties.Certificate:                       Certificate,
		properties.CertificateAuthority:              CertificateAuthority,
		properties.Certificates:                      Certificates,
		properties.ChangeSet:                         ChangeSet,
		properties.Charset:                           Charset,
		properties.CheckHTTPCode:                     CheckHTTPCode,
		properties.CheckInterval:                     CheckInterval,
		properties.CheckPath:                         CheckPath,
		properties.CheckPort:                         CheckPort,
		properties.CheckProtocol:                     CheckProtocol,
		properties.CheckTimeout:                      CheckTimeout,
		properties.CIDR:                              CIDR,
		properties.CIDRv6:                            CIDRv6,
		properties.CipherSuite:                       CipherSuite,
		properties.Class:                             Class,
		properties.Cluster:                           Cluster,
		properties.Comment:                           Comment,
		properties.Config:                            Config,
		properties.ContainerInstance:                 ContainerInstance,
		properties.ContainersImages:                  ContainersImages,
		properties.ContainerTask:                     ContainerTask,
		properties.Continent:                         Continent,
		properties.Cooldown:                          Cooldown,
		properties.CopyTagsToSnapshot:                CopyTagsToSnapshot,
		properties.Country:                           Country,
		properties.Created:                           Created,
		properties.DBSecurityGroups:                  DBSecurityGroups,
		properties.DBSubnetGroup:                     DBSubnetGroup,
		properties.Default:                           Default,
		properties.DefaultCooldown:                   DefaultCooldown,
		properties.Delay:                             Delay,
		properties.DeploymentName:                    DeploymentName,
		properties.Deployments:                       Deployments,
		properties.Description:                       Description,
		properties.DesiredCapacity:                   DesiredCapacity,
		properties.Dimensions:                        Dimensions,
		properties.DisableRollback:                   DisableRollback,
		properties.DockerVersion:                     DockerVersion,
		properties.Document:                          Document,
		properties.Enabled:                           Enabled,
		properties.Encrypted:                         Encrypted,
		properties.Endpoint:                          Endpoint,
		properties.Engine:                            Engine,
		properties.EngineVersion:                     EngineVersion,
		properties.ExitCode:                          ExitCode,
		properties.Failover:                          Failover,
		properties.Fingerprint:                       Fingerprint,
		properties.GlobalID:                          GlobalID,
		properties.GranteeType:                       GranteeType,
		properties.Grants:                            Grants,
		properties.Handler:                           Handler,
		properties.Hash:                              Hash,
		properties.HealthCheck:                       HealthCheck,
		properties.HealthCheckGracePeriod:            HealthCheckGracePeriod,
		properties.HealthCheckType:                   HealthCheckType,
		properties.HealthyThresholdCount:             HealthyThresholdCount,
		properties.Host:                              Host,
		properties.HTTPVersion:                       HTTPVersion,
		properties.Hypervisor:                        Hypervisor,
		properties.ID:                                ID,
		properties.Image:                             Image,
		properties.InboundRules:                      InboundRules,
		properties.InlinePolicies:                    InlinePolicies,
		properties.Instance:                          Instance,
		properties.InstanceOwner:                     InstanceOwner,
		properties.Instances:                         Instances,
		properties.InsufficientDataActions:           InsufficientDataActions,
		properties.IOPS:                              IOPS,
		properties.IPType:                            IPType,
		properties.IPv6Addresses:                     IPv6Addresses,
		properties.IPv6Enabled:                       IPv6Enabled,
		properties.Key:                               Key,
		properties.KeyName:                           KeyName,
		properties.KeyPair:                           KeyPair,
		properties.LatestRestorableTime:              LatestRestorableTime,
		properties.LaunchConfigurationName:           LaunchConfigurationName,
		properties.Launched:                          Launched,
		properties.License:                           License,
		properties.Lifecycle:                         Lifecycle,
		properties.LoadBalancer:                      LoadBalancer,
		properties.Location:                          Location,
		properties.MACAddress:                        MACAddress,
		properties.Main:                              Main,
		properties.MaxSize:                           MaxSize,
		properties.Memory:                            Memory,
		properties.Messages:                          Messages,
		properties.MetricName:                        MetricName,
		properties.MinSize:                           MinSize,
		properties.Modified:                          Modified,
		properties.MonitoringInterval:                MonitoringInterval,
		properties.MonitoringRole:                    MonitoringRole,
		properties.MultiAZ:                           MultiAZ,
		properties.Name:                              Name,
		properties.Namespace:                         Namespace,
		properties.NetworkInterfaces:                 NetworkInterfaces,
		properties.NewInstancesProtected:             NewInstancesProtected,
		properties.Notifications:                     Notifications,
		properties.OKActions:                         OKActions,
		properties.OptionGroups:                      OptionGroups,
		properties.Origins:                           Origins,
		properties.OutboundRules:                     OutboundRules,
		properties.Outputs:                           Outputs,
		properties.Owner:                             Owner,
		properties.ParameterGroups:                   ParameterGroups,
		properties.Parameters:                        Parameters,
		properties.PasswordLastUsed:                  PasswordLastUsed,
		properties.Path:                              Path,
		properties.PathPrefix:                        PathPrefix,
		properties.PendingTasksCount:                 PendingTasksCount,
		properties.PlacementGroup:                    PlacementGroup,
		properties.Port:                              Port,
		properties.Ports:                             Ports,
		properties.PortRange:                         PortRange,
		properties.PreferredBackupDate:               PreferredBackupDate,
		properties.PreferredMaintenanceDate:          PreferredMaintenanceDate,
		properties.PriceClass:                        PriceClass,
		properties.Private:                           Private,
		properties.PrivateDNS:                        PrivateDNS,
		properties.PrivateIP:                         PrivateIP,
		properties.Profile:                           Profile,
		properties.Progress:                          Progress,
		properties.Protocol:                          Protocol,
		properties.Public:                            Public,
		properties.PublicDNS:                         PublicDNS,
		properties.PublicIP:                          PublicIP,
		properties.RecordCount:                       RecordCount,
		properties.Records:                           Records,
		properties.Region:                            Region,
		properties.RegisteredContainerInstancesCount: RegisteredContainerInstancesCount,
		properties.ReplicaOf:                         ReplicaOf,
		properties.Role:                              Role,
		properties.Roles:                             Roles,
		properties.RootDevice:                        RootDevice,
		properties.RootDeviceType:                    RootDeviceType,
		properties.Routes:                            Routes,
		properties.RunningTasksCount:                 RunningTasksCount,
		properties.Runtime:                           Runtime,
		properties.ScalingAdjustment:                 ScalingAdjustment,
		properties.ScalingGroupName:                  ScalingGroupName,
		properties.Scheme:                            Scheme,
		properties.SecondaryAvailabilityZone:         SecondaryAvailabilityZone,
		properties.SecurityGroups:                    SecurityGroups,
		properties.Set:                               Set,
		properties.Size:                              Size,
		properties.Source:                            Source,
		properties.SpotInstanceRequestId:             SpotInstanceRequestId,
		properties.SpotPrice:                         SpotPrice,
		properties.SSLSupportMethod:                  SSLSupportMethod,
		properties.State:                             State,
		properties.StateMessage:                      StateMessage,
		properties.Stopped:                           Stopped,
		properties.Storage:                           Storage,
		properties.StorageType:                       StorageType,
		properties.Subnet:                            Subnet,
		properties.Subnets:                           Subnets,
		properties.Tags:                              Tags,
		properties.TargetGroups:                      TargetGroups,
		properties.Timeout:                           Timeout,
		properties.Timezone:                          Timezone,
		properties.TLSVersionRequired:                TLSVersionRequired,
		properties.Topic:                             Topic,
		properties.TrafficPolicyInstance:             TrafficPolicyInstance,
		properties.TrustPolicy:                       TrustPolicy,
		properties.TTL:                               TTL,
		properties.Type:                              Type,
		properties.UnhealthyThresholdCount:           UnhealthyThresholdCount,
		properties.Updated:                           Updated,
		properties.URI:                               URI,
		properties.UserData:                          UserData,
		properties.Username:                          Username,
		properties.Value:                             Value,
		properties.Version:                           Version,
		properties.Virtualization:                    Virtualization,
		properties.Volume:                            Volume,
		properties.Vpc:                               Vpc,
		properties.Vpcs:                              Vpcs,
		properties.WebACL:                            WebACL,
		properties.Weight:                            Weight,
		properties.Zone:                              Zone,
	}
}

var Properties = RDFProperties{
	Account:                 {ID: Account, RdfType: "rdf:Property", RdfsLabel: "Account", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	ACMCertificate:          {ID: ACMCertificate, RdfType: "rdf:Property", RdfsLabel: "ACMCertificate", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Actions:                 {ID: Actions, RdfType: "rdf:Property", RdfsLabel: "Actions", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	ActionsEnabled:          {ID: ActionsEnabled, RdfType: "rdf:Property", RdfsLabel: "ActionsEnabled", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	ActiveServicesCount:     {ID: ActiveServicesCount, RdfType: "rdf:Property", RdfsLabel: "ActiveServicesCount", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	AdjustmentType:          {ID: AdjustmentType, RdfType: "rdf:Property", RdfsLabel: "AdjustmentType", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Affinity:                {ID: Affinity, RdfType: "rdf:Property", RdfsLabel: "Affinity", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	AgentConnected:          {ID: AgentConnected, RdfType: "rdf:Property", RdfsLabel: "AgentConnected", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	AgentState:              {ID: AgentState, RdfType: "rdf:Property", RdfsLabel: "AgentState", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	AgentVersion:            {ID: AgentVersion, RdfType: "rdf:Property", RdfsLabel: "AgentVersion", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	AlarmActions:            {ID: AlarmActions, RdfType: "rdf:Property", RdfsLabel: "AlarmActions", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	AlarmNames:              {ID: AlarmNames, RdfType: "rdf:Property", RdfsLabel: "AlarmNames", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	Alias:                   {ID: Alias, RdfType: "rdf:Property", RdfsLabel: "Alias", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Aliases:                 {ID: Aliases, RdfType: "rdf:Property", RdfsLabel: "Aliases", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	ApproximateMessageCount: {ID: ApproximateMessageCount, RdfType: "rdf:Property", RdfsLabel: "ApproximateMessageCount", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Architecture:            {ID: Architecture, RdfType: "rdf:Property", RdfsLabel: "Architecture", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Arn:                     {ID: Arn, RdfType: "rdf:Property", RdfsLabel: "Arn", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Association:             {ID: Association, RdfType: "rdf:Property", RdfsLabel: "Association", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Associations:            {ID: Associations, RdfType: "rdf:Property", RdfsLabel: "Associations", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:KeyValue"},
	Attachable:              {ID: Attachable, RdfType: "rdf:Property", RdfsLabel: "Attachable", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	Attached:                {ID: Attached, RdfType: "rdf:Property", RdfsLabel: "Attached", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	AttachedAt:              {ID: AttachedAt, RdfType: "rdf:Property", RdfsLabel: "AttachedAt", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:dateTime"},
	Attachment:              {ID: Attachment, RdfType: "rdf:Property", RdfsLabel: "Attachment", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Attributes:              {ID: Attributes, RdfType: "rdf:Property", RdfsLabel: "Attributes", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:KeyValue"},
	AutoUpgrade:             {ID: AutoUpgrade, RdfType: "rdf:Property", RdfsLabel: "AutoUpgrade", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	AvailabilityZone:        {ID: AvailabilityZone, RdfType: "rdf:Property", RdfsLabel: "AvailabilityZone", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	AvailabilityZones:       {ID: AvailabilityZones, RdfType: "rdf:Property", RdfsLabel: "AvailabilityZones", RdfsDefinedBy: "rdfs:list", RdfsDataType: "rdfs:Class"},
	BackupRetentionPeriod:   {ID: BackupRetentionPeriod, RdfType: "rdf:Property", RdfsLabel: "BackupRetentionPeriod", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:dateTime"},
	Bucket:                  {ID: Bucket, RdfType: "rdf:Property", RdfsLabel: "Bucket", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	CallerReference:         {ID: CallerReference, RdfType: "rdf:Property", RdfsLabel: "CallerReference", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Capabilities:            {ID: Capabilities, RdfType: "rdf:Property", RdfsLabel: "Capabilities", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	Certificate:             {ID: Certificate, RdfType: "rdf:Property", RdfsLabel: "Certificate", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	CertificateAuthority:    {ID: CertificateAuthority, RdfType: "rdf:Property", RdfsLabel: "CertificateAuthority", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Certificates:            {ID: Certificates, RdfType: "rdf:Property", RdfsLabel: "Certificates", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	ChangeSet:               {ID: ChangeSet, RdfType: "rdf:Property", RdfsLabel: "ChangeSet", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Charset:                 {ID: Charset, RdfType: "rdf:Property", RdfsLabel: "Charset", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	CheckHTTPCode:           {ID: CheckHTTPCode, RdfType: "rdf:Property", RdfsLabel: "CheckHTTPCode", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	CheckInterval:           {ID: CheckInterval, RdfType: "rdf:Property", RdfsLabel: "CheckInterval", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	CheckPath:               {ID: CheckPath, RdfType: "rdf:Property", RdfsLabel: "CheckPath", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	CheckPort:               {ID: CheckPort, RdfType: "rdf:Property", RdfsLabel: "CheckPort", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	CheckProtocol:           {ID: CheckProtocol, RdfType: "rdf:Property", RdfsLabel: "CheckProtocol", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	CheckTimeout:            {ID: CheckTimeout, RdfType: "rdf:Property", RdfsLabel: "CheckTimeout", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	CIDR:                    {ID: CIDR, RdfType: "rdf:Property", RdfsLabel: "CIDR", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	CIDRv6:                  {ID: CIDRv6, RdfType: "rdf:Property", RdfsLabel: "CIDRv6", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	CipherSuite:             {ID: CipherSuite, RdfType: "rdf:Property", RdfsLabel: "CipherSuite", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Class:                   {ID: Class, RdfType: "rdf:Property", RdfsLabel: "Class", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Cluster:                 {ID: Cluster, RdfType: "rdf:Property", RdfsLabel: "Cluster", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Comment:                 {ID: Comment, RdfType: "rdf:Property", RdfsLabel: "Comment", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Config:                  {ID: Config, RdfType: "rdf:Property", RdfsLabel: "Config", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	ContainerInstance:       {ID: ContainerInstance, RdfType: "rdf:Property", RdfsLabel: "ContainerInstance", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	ContainersImages:        {ID: ContainersImages, RdfType: "rdf:Property", RdfsLabel: "ContainersImages", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:KeyValue"},
	ContainerTask:           {ID: ContainerTask, RdfType: "rdf:Property", RdfsLabel: "ContainerTask", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	Continent:               {ID: Continent, RdfType: "rdf:Property", RdfsLabel: "Continent", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Cooldown:                {ID: Cooldown, RdfType: "rdf:Property", RdfsLabel: "Cooldown", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	CopyTagsToSnapshot:      {ID: CopyTagsToSnapshot, RdfType: "rdf:Property", RdfsLabel: "CopyTagsToSnapshot", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Country:                 {ID: Country, RdfType: "rdf:Property", RdfsLabel: "Country", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Created:                 {ID: Created, RdfType: "rdf:Property", RdfsLabel: "Created", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:dateTime"},
	DBSecurityGroups:        {ID: DBSecurityGroups, RdfType: "rdf:Property", RdfsLabel: "DBSecurityGroups", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	DBSubnetGroup:           {ID: DBSubnetGroup, RdfType: "rdf:Property", RdfsLabel: "DBSubnetGroup", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Default:                 {ID: Default, RdfType: "rdf:Property", RdfsLabel: "Default", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	DefaultCooldown:         {ID: DefaultCooldown, RdfType: "rdf:Property", RdfsLabel: "DefaultCooldown", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Delay:                   {ID: Delay, RdfType: "rdf:Property", RdfsLabel: "Delay", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	DeploymentName:          {ID: DeploymentName, RdfType: "rdf:Property", RdfsLabel: "DeploymentName", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Deployments:             {ID: Deployments, RdfType: "rdf:Property", RdfsLabel: "Deployments", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:KeyValue"},
	Description:             {ID: Description, RdfType: "rdf:Property", RdfsLabel: "Description", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	DesiredCapacity:         {ID: DesiredCapacity, RdfType: "rdf:Property", RdfsLabel: "DesiredCapacity", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Dimensions:              {ID: Dimensions, RdfType: "rdf:Property", RdfsLabel: "Dimensions", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:KeyValue"},
	DisableRollback:         {ID: DisableRollback, RdfType: "rdf:Property", RdfsLabel: "DisableRollback", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	DockerVersion:           {ID: DockerVersion, RdfType: "rdf:Property", RdfsLabel: "DockerVersion", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Document:                {ID: Document, RdfType: "rdf:Property", RdfsLabel: "Document", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Enabled:                 {ID: Enabled, RdfType: "rdf:Property", RdfsLabel: "Enabled", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	Encrypted:               {ID: Encrypted, RdfType: "rdf:Property", RdfsLabel: "Encrypted", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	Endpoint:                {ID: Endpoint, RdfType: "rdf:Property", RdfsLabel: "Endpoint", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Engine:                  {ID: Engine, RdfType: "rdf:Property", RdfsLabel: "Engine", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	EngineVersion:           {ID: EngineVersion, RdfType: "rdf:Property", RdfsLabel: "EngineVersion", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	ExitCode:                {ID: ExitCode, RdfType: "rdf:Property", RdfsLabel: "ExitCode", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Failover:                {ID: Failover, RdfType: "rdf:Property", RdfsLabel: "Failover", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Fingerprint:             {ID: Fingerprint, RdfType: "rdf:Property", RdfsLabel: "Fingerprint", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	GlobalID:                {ID: GlobalID, RdfType: "rdf:Property", RdfsLabel: "GlobalID", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	GranteeType:             {ID: GranteeType, RdfType: "rdf:Property", RdfsLabel: "GranteeType", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Grants:                  {ID: Grants, RdfType: "rdf:Property", RdfsLabel: "Grants", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:Grant"},
	Handler:                 {ID: Handler, RdfType: "rdf:Property", RdfsLabel: "Handler", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Hash:                    {ID: Hash, RdfType: "rdf:Property", RdfsLabel: "Hash", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	HealthCheck:             {ID: HealthCheck, RdfType: "rdf:Property", RdfsLabel: "HealthCheck", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	HealthCheckGracePeriod:  {ID: HealthCheckGracePeriod, RdfType: "rdf:Property", RdfsLabel: "HealthCheckGracePeriod", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	HealthCheckType:         {ID: HealthCheckType, RdfType: "rdf:Property", RdfsLabel: "HealthCheckType", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	HealthyThresholdCount:   {ID: HealthyThresholdCount, RdfType: "rdf:Property", RdfsLabel: "HealthyThresholdCount", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Host:                    {ID: Host, RdfType: "rdf:Property", RdfsLabel: "Host", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	HTTPVersion:             {ID: HTTPVersion, RdfType: "rdf:Property", RdfsLabel: "HTTPVersion", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Hypervisor:              {ID: Hypervisor, RdfType: "rdf:Property", RdfsLabel: "Hypervisor", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	ID:                      {ID: ID, RdfType: "rdf:Property", RdfsLabel: "ID", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Image:                   {ID: Image, RdfType: "rdf:Property", RdfsLabel: "Image", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	InboundRules:            {ID: InboundRules, RdfType: "rdf:Property", RdfsLabel: "InboundRules", RdfsDefinedBy: "rdfs:list", RdfsDataType: "net-owl:FirewallRule"},
	InlinePolicies:          {ID: InlinePolicies, RdfType: "rdf:Property", RdfsLabel: "InlinePolicies", RdfsDefinedBy: "rdfs:list", RdfsDataType: "rdfs:Class"},
	Instance:                {ID: Instance, RdfType: "rdf:Property", RdfsLabel: "Instance", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	InstanceOwner:           {ID: InstanceOwner, RdfType: "rdf:Property", RdfsLabel: "InstanceOwner", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Instances:               {ID: Instances, RdfType: "rdf:Property", RdfsLabel: "Instances", RdfsDefinedBy: "rdfs:list", RdfsDataType: "rdfs:Class"},
	InsufficientDataActions: {ID: InsufficientDataActions, RdfType: "rdf:Property", RdfsLabel: "InsufficientDataActions", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	IOPS:                     {ID: IOPS, RdfType: "rdf:Property", RdfsLabel: "IOPS", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	IPType:                   {ID: IPType, RdfType: "rdf:Property", RdfsLabel: "IPType", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	IPv6Addresses:            {ID: IPv6Addresses, RdfType: "rdf:Property", RdfsLabel: "IPv6Addresses", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	IPv6Enabled:              {ID: IPv6Enabled, RdfType: "rdf:Property", RdfsLabel: "IPv6Enabled", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	Key:                      {ID: Key, RdfType: "rdf:Property", RdfsLabel: "Key", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	KeyName:                  {ID: KeyName, RdfType: "rdf:Property", RdfsLabel: "KeyName", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	KeyPair:                  {ID: KeyPair, RdfType: "rdf:Property", RdfsLabel: "KeyPair", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	LatestRestorableTime:     {ID: LatestRestorableTime, RdfType: "rdf:Property", RdfsLabel: "LatestRestorableTime", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:dateTime"},
	LaunchConfigurationName:  {ID: LaunchConfigurationName, RdfType: "rdf:Property", RdfsLabel: "LaunchConfigurationName", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Launched:                 {ID: Launched, RdfType: "rdf:Property", RdfsLabel: "Launched", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:dateTime"},
	License:                  {ID: License, RdfType: "rdf:Property", RdfsLabel: "License", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Lifecycle:                {ID: Lifecycle, RdfType: "rdf:Property", RdfsLabel: "Lifecycle", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	LoadBalancer:             {ID: LoadBalancer, RdfType: "rdf:Property", RdfsLabel: "LoadBalancer", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	Location:                 {ID: Location, RdfType: "rdf:Property", RdfsLabel: "Location", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	MACAddress:               {ID: MACAddress, RdfType: "rdf:Property", RdfsLabel: "MACAddress", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Main:                     {ID: Main, RdfType: "rdf:Property", RdfsLabel: "Main", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	MaxSize:                  {ID: MaxSize, RdfType: "rdf:Property", RdfsLabel: "MaxSize", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Memory:                   {ID: Memory, RdfType: "rdf:Property", RdfsLabel: "Memory", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Messages:                 {ID: Messages, RdfType: "rdf:Property", RdfsLabel: "Messages", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	MetricName:               {ID: MetricName, RdfType: "rdf:Property", RdfsLabel: "MetricName", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	MinSize:                  {ID: MinSize, RdfType: "rdf:Property", RdfsLabel: "MinSize", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Modified:                 {ID: Modified, RdfType: "rdf:Property", RdfsLabel: "Modified", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:dateTime"},
	MonitoringInterval:       {ID: MonitoringInterval, RdfType: "rdf:Property", RdfsLabel: "MonitoringInterval", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	MonitoringRole:           {ID: MonitoringRole, RdfType: "rdf:Property", RdfsLabel: "MonitoringRole", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	MultiAZ:                  {ID: MultiAZ, RdfType: "rdf:Property", RdfsLabel: "MultiAZ", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Name:                     {ID: Name, RdfType: "rdf:Property", RdfsLabel: "Name", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Namespace:                {ID: Namespace, RdfType: "rdf:Property", RdfsLabel: "Namespace", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	NetworkInterfaces:        {ID: NetworkInterfaces, RdfType: "rdf:Property", RdfsLabel: "NetworkInterfaces", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	NewInstancesProtected:    {ID: NewInstancesProtected, RdfType: "rdf:Property", RdfsLabel: "NewInstancesProtected", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	Notifications:            {ID: Notifications, RdfType: "rdf:Property", RdfsLabel: "Notifications", RdfsDefinedBy: "rdfs:list", RdfsDataType: "rdfs:Class"},
	OKActions:                {ID: OKActions, RdfType: "rdf:Property", RdfsLabel: "OKActions", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	OptionGroups:             {ID: OptionGroups, RdfType: "rdf:Property", RdfsLabel: "OptionGroups", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	Origins:                  {ID: Origins, RdfType: "rdf:Property", RdfsLabel: "Origins", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:DistributionOrigin"},
	OutboundRules:            {ID: OutboundRules, RdfType: "rdf:Property", RdfsLabel: "OutboundRules", RdfsDefinedBy: "rdfs:list", RdfsDataType: "net-owl:FirewallRule"},
	Outputs:                  {ID: Outputs, RdfType: "rdf:Property", RdfsLabel: "Outputs", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:KeyValue"},
	Owner:                    {ID: Owner, RdfType: "rdf:Property", RdfsLabel: "Owner", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	ParameterGroups:          {ID: ParameterGroups, RdfType: "rdf:Property", RdfsLabel: "ParameterGroups", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	Parameters:               {ID: Parameters, RdfType: "rdf:Property", RdfsLabel: "Parameters", RdfsDefinedBy: "rdfs:list", RdfsDataType: "cloud-owl:KeyValue"},
	PasswordLastUsed:         {ID: PasswordLastUsed, RdfType: "rdf:Property", RdfsLabel: "PasswordLastUsed", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:dateTime"},
	Path:                     {ID: Path, RdfType: "rdf:Property", RdfsLabel: "Path", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	PathPrefix:               {ID: PathPrefix, RdfType: "rdf:Property", RdfsLabel: "PathPrefix", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	PendingTasksCount:        {ID: PendingTasksCount, RdfType: "rdf:Property", RdfsLabel: "PendingTasksCount", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	PlacementGroup:           {ID: PlacementGroup, RdfType: "rdf:Property", RdfsLabel: "PlacementGroup", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Port:                     {ID: Port, RdfType: "rdf:Property", RdfsLabel: "Port", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Ports:                    {ID: Ports, RdfType: "rdf:Property", RdfsLabel: "Ports", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	PortRange:                {ID: PortRange, RdfType: "rdfs:subPropertyOf", RdfsLabel: "PortRange", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	PreferredBackupDate:      {ID: PreferredBackupDate, RdfType: "rdf:Property", RdfsLabel: "PreferredBackupDate", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	PreferredMaintenanceDate: {ID: PreferredMaintenanceDate, RdfType: "rdf:Property", RdfsLabel: "PreferredMaintenanceDate", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	PriceClass:               {ID: PriceClass, RdfType: "rdf:Property", RdfsLabel: "PriceClass", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Private:                  {ID: Private, RdfType: "rdf:Property", RdfsLabel: "Private", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	PrivateDNS:               {ID: PrivateDNS, RdfType: "rdf:Property", RdfsLabel: "PrivateDNS", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	PrivateIP:                {ID: PrivateIP, RdfType: "rdf:Property", RdfsLabel: "PrivateIP", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Profile:                  {ID: Profile, RdfType: "rdf:Property", RdfsLabel: "Profile", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Progress:                 {ID: Progress, RdfType: "rdf:Property", RdfsLabel: "Progress", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Protocol:                 {ID: Protocol, RdfType: "rdf:Property", RdfsLabel: "Protocol", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Public:                   {ID: Public, RdfType: "rdf:Property", RdfsLabel: "Public", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:boolean"},
	PublicDNS:                {ID: PublicDNS, RdfType: "rdf:Property", RdfsLabel: "PublicDNS", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	PublicIP:                 {ID: PublicIP, RdfType: "rdf:Property", RdfsLabel: "PublicIP", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	RecordCount:              {ID: RecordCount, RdfType: "rdf:Property", RdfsLabel: "RecordCount", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Records:                  {ID: Records, RdfType: "rdf:Property", RdfsLabel: "Records", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	Region:                   {ID: Region, RdfType: "rdf:Property", RdfsLabel: "Region", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	RegisteredContainerInstancesCount: {ID: RegisteredContainerInstancesCount, RdfType: "rdf:Property", RdfsLabel: "RegisteredContainerInstancesCount", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	ReplicaOf:                         {ID: ReplicaOf, RdfType: "rdf:Property", RdfsLabel: "ReplicaOf", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Role:                              {ID: Role, RdfType: "rdf:Property", RdfsLabel: "Role", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	Roles:                             {ID: Roles, RdfType: "rdf:Property", RdfsLabel: "Roles", RdfsDefinedBy: "rdfs:list", RdfsDataType: "rdfs:Class"},
	RootDevice:                        {ID: RootDevice, RdfType: "rdf:Property", RdfsLabel: "RootDevice", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	RootDeviceType:                    {ID: RootDeviceType, RdfType: "rdf:Property", RdfsLabel: "RootDeviceType", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Routes:                            {ID: Routes, RdfType: "rdf:Property", RdfsLabel: "Routes", RdfsDefinedBy: "rdfs:list", RdfsDataType: "net-owl:Route"},
	RunningTasksCount:                 {ID: RunningTasksCount, RdfType: "rdf:Property", RdfsLabel: "RunningTasksCount", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Runtime:                           {ID: Runtime, RdfType: "rdf:Property", RdfsLabel: "Runtime", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	ScalingAdjustment:                 {ID: ScalingAdjustment, RdfType: "rdf:Property", RdfsLabel: "ScalingAdjustment", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	ScalingGroupName:                  {ID: ScalingGroupName, RdfType: "rdf:Property", RdfsLabel: "ScalingGroupName", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Scheme:                            {ID: Scheme, RdfType: "rdf:Property", RdfsLabel: "Scheme", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	SecondaryAvailabilityZone: {ID: SecondaryAvailabilityZone, RdfType: "rdf:Property", RdfsLabel: "SecondaryAvailabilityZone", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	SecurityGroups:            {ID: SecurityGroups, RdfType: "rdf:Property", RdfsLabel: "SecurityGroups", RdfsDefinedBy: "rdfs:list", RdfsDataType: "rdfs:Class"},
	Set:                       {ID: Set, RdfType: "rdf:Property", RdfsLabel: "Set", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Size:                      {ID: Size, RdfType: "rdf:Property", RdfsLabel: "Size", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Source:                    {ID: Source, RdfType: "rdf:Property", RdfsLabel: "Source", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	SpotInstanceRequestId:     {ID: SpotInstanceRequestId, RdfType: "rdf:Property", RdfsLabel: "SpotInstanceRequestId", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	SpotPrice:                 {ID: SpotPrice, RdfType: "rdf:Property", RdfsLabel: "SpotPrice", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	SSLSupportMethod:          {ID: SSLSupportMethod, RdfType: "rdf:Property", RdfsLabel: "SSLSupportMethod", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	State:                     {ID: State, RdfType: "rdf:Property", RdfsLabel: "State", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	StateMessage:              {ID: StateMessage, RdfType: "rdf:Property", RdfsLabel: "StateMessage", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Stopped:                   {ID: Stopped, RdfType: "rdf:Property", RdfsLabel: "Stopped", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:dateTime"},
	Storage:                   {ID: Storage, RdfType: "rdf:Property", RdfsLabel: "Storage", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	StorageType:               {ID: StorageType, RdfType: "rdf:Property", RdfsLabel: "StorageType", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Subnet:                    {ID: Subnet, RdfType: "rdf:Property", RdfsLabel: "Subnet", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	Subnets:                   {ID: Subnets, RdfType: "rdf:Property", RdfsLabel: "Subnets", RdfsDefinedBy: "rdfs:list", RdfsDataType: "rdfs:Class"},
	Tags:                      {ID: Tags, RdfType: "rdf:Property", RdfsLabel: "Tags", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	TargetGroups:              {ID: TargetGroups, RdfType: "rdf:Property", RdfsLabel: "TargetGroups", RdfsDefinedBy: "rdfs:list", RdfsDataType: "xsd:string"},
	Timeout:                   {ID: Timeout, RdfType: "rdf:Property", RdfsLabel: "Timeout", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Timezone:                  {ID: Timezone, RdfType: "rdf:Property", RdfsLabel: "Timezone", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	TLSVersionRequired:        {ID: TLSVersionRequired, RdfType: "rdf:Property", RdfsLabel: "TLSVersionRequired", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Topic:                     {ID: Topic, RdfType: "rdf:Property", RdfsLabel: "Topic", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	TrafficPolicyInstance: {ID: TrafficPolicyInstance, RdfType: "rdf:Property", RdfsLabel: "TrafficPolicyInstance", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	TrustPolicy:           {ID: TrustPolicy, RdfType: "rdf:Property", RdfsLabel: "TrustPolicy", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	TTL:                   {ID: TTL, RdfType: "rdf:Property", RdfsLabel: "TTL", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Type:                  {ID: Type, RdfType: "rdf:Property", RdfsLabel: "Type", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	UnhealthyThresholdCount: {ID: UnhealthyThresholdCount, RdfType: "rdf:Property", RdfsLabel: "UnhealthyThresholdCount", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:int"},
	Updated:                 {ID: Updated, RdfType: "rdf:Property", RdfsLabel: "Updated", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	URI:                     {ID: URI, RdfType: "rdf:Property", RdfsLabel: "URI", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	UserData:                {ID: UserData, RdfType: "rdf:Property", RdfsLabel: "UserData", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Username:                {ID: Username, RdfType: "rdf:Property", RdfsLabel: "Username", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Value:                   {ID: Value, RdfType: "rdf:Property", RdfsLabel: "Value", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Version:                 {ID: Version, RdfType: "rdf:Property", RdfsLabel: "Version", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Virtualization:          {ID: Virtualization, RdfType: "rdf:Property", RdfsLabel: "Virtualization", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Volume:                  {ID: Volume, RdfType: "rdf:Property", RdfsLabel: "Volume", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	Vpc:                     {ID: Vpc, RdfType: "rdf:Property", RdfsLabel: "Vpc", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
	Vpcs:                    {ID: Vpcs, RdfType: "rdf:Property", RdfsLabel: "Vpcs", RdfsDefinedBy: "rdfs:list", RdfsDataType: "rdfs:Class"},
	WebACL:                  {ID: WebACL, RdfType: "rdf:Property", RdfsLabel: "WebACL", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Weight:                  {ID: Weight, RdfType: "rdf:Property", RdfsLabel: "Weight", RdfsDefinedBy: "rdfs:Literal", RdfsDataType: "xsd:string"},
	Zone:                    {ID: Zone, RdfType: "rdf:Property", RdfsLabel: "Zone", RdfsDefinedBy: "rdfs:Class", RdfsDataType: "xsd:string"},
}
