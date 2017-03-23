/*
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
import "github.com/wallix/awless/cloud/properties"

var awsResourcesDef = map[string]map[string]*propertyTransform{
	//EC2
	cloud.Instance: {
		properties.Name:              {name: "Tags", transform: extractTagFn("Name")},
		properties.Type:              {name: "InstanceType", transform: extractValueFn},
		properties.Subnet:            {name: "SubnetId", transform: extractValueFn},
		properties.Vpc:               {name: "VpcId", transform: extractValueFn},
		properties.PublicIP:          {name: "PublicIpAddress", transform: extractValueFn},
		properties.PrivateIP:         {name: "PrivateIpAddress", transform: extractValueFn},
		properties.Image:             {name: "ImageId", transform: extractValueFn},
		properties.Launched:          {name: "LaunchTime", transform: extractValueFn},
		properties.State:             {name: "State", transform: extractFieldFn("Name")},
		properties.SSHKey:            {name: "KeyName", transform: extractValueFn},
		properties.SecurityGroups:    {name: "SecurityGroups", transform: extractStringSliceValues("GroupId")},
		properties.Affinity:          {name: "Placement", transform: extractFieldFn("Affinity")},
		properties.AvailabilityZone:  {name: "Placement", transform: extractFieldFn("AvailabilityZone")},
		properties.PlacementGroup:    {name: "Placement", transform: extractFieldFn("GroupName")},
		properties.Host:              {name: "Placement", transform: extractFieldFn("HostId")},
		properties.Architecture:      {name: "Architecture", transform: extractValueFn},
		properties.Hypervisor:        {name: "Hypervisor", transform: extractValueFn},
		properties.Profile:           {name: "IamInstanceProfile", transform: extractFieldFn("Arn")},
		properties.Lifecycle:         {name: "InstanceLifecycle", transform: extractValueFn},
		properties.NetworkInterfaces: {name: "NetworkInterfaces", transform: extractStringSliceValues("NetworkInterfaceId")},
		properties.PublicDNS:         {name: "PublicDnsName", transform: extractValueFn},
		properties.RootDevice:        {name: "RootDeviceName", transform: extractValueFn},
		properties.RootDeviceType:    {name: "RootDeviceType", transform: extractValueFn},
	},
	cloud.Vpc: {
		properties.Name:    {name: "Tags", transform: extractTagFn("Name")},
		properties.Default: {name: "IsDefault", transform: extractValueFn},
		properties.State:   {name: "State", transform: extractValueFn},
		properties.CIDR:    {name: "CidrBlock", transform: extractValueFn},
	},
	cloud.Subnet: {
		properties.Name:             {name: "Tags", transform: extractTagFn("Name")},
		properties.Vpc:              {name: "VpcId", transform: extractValueFn},
		properties.Public:           {name: "MapPublicIpOnLaunch", transform: extractValueFn},
		properties.State:            {name: "State", transform: extractValueFn},
		properties.CIDR:             {name: "CidrBlock", transform: extractValueFn},
		properties.AvailabilityZone: {name: "AvailabilityZone", transform: extractValueFn},
		properties.Default:          {name: "DefaultForAz", transform: extractValueFn},
	},
	cloud.SecurityGroup: {
		properties.Name:          {name: "GroupName", transform: extractValueFn},
		properties.Description:   {name: "Description", transform: extractValueFn},
		properties.InboundRules:  {name: "IpPermissions", transform: extractIpPermissionSliceFn},
		properties.OutboundRules: {name: "IpPermissionsEgress", transform: extractIpPermissionSliceFn},
		properties.Owner:         {name: "OwnerId", transform: extractValueFn},
		properties.Vpc:           {name: "VpcId", transform: extractValueFn},
	},
	cloud.Keypair: {
		properties.Fingerprint: {name: "KeyFingerprint", transform: extractValueFn},
	},
	cloud.Volume: {
		properties.Name:             {name: "Tags", transform: extractTagFn("Name")},
		properties.Type:             {name: "VolumeType", transform: extractValueFn},
		properties.State:            {name: "State", transform: extractValueFn},
		properties.Size:             {name: "Size", transform: extractValueFn},
		properties.Encrypted:        {name: "Encrypted", transform: extractValueFn},
		properties.Created:          {name: "CreateTime", transform: extractTimeFn},
		properties.AvailabilityZone: {name: "AvailabilityZone", transform: extractValueFn},
	},
	cloud.InternetGateway: {
		properties.Name: {name: "Tags", transform: extractTagFn("Name")},
		properties.Vpcs: {name: "Attachments", transform: extractStringSliceValues("VpcId")},
	},
	cloud.RouteTable: {
		properties.Name:   {name: "Tags", transform: extractTagFn("Name")},
		properties.Vpc:    {name: "VpcId", transform: extractValueFn},
		properties.Routes: {name: "Routes", transform: extractRoutesSliceFn},
		properties.Main:   {name: "Associations", transform: extractHasATrueBoolInStructSliceFn("Main")},
	},
	cloud.AvailabilityZone: {
		properties.Name:     {name: "ZoneName", transform: extractValueFn},
		properties.State:    {name: "State", transform: extractValueFn},
		properties.Region:   {name: "RegionName", transform: extractValueFn},
		properties.Messages: {name: "Messages", transform: extractStringSliceValues("Message")},
	},
	// LoadBalancer
	cloud.LoadBalancer: {
		properties.Name:              {name: "LoadBalancerName", transform: extractValueFn},
		properties.AvailabilityZones: {name: "AvailabilityZones", transform: extractStringSliceValues("ZoneName")},
		properties.Subnets:           {name: "AvailabilityZones", transform: extractStringSliceValues("SubnetId")},
		properties.Zone:              {name: "CanonicalHostedZoneId", transform: extractValueFn},
		properties.Created:           {name: "CreatedTime", transform: extractTimeFn},
		properties.PublicDNS:         {name: "DNSName", transform: extractValueFn},
		properties.IPType:            {name: "IpAddressType", transform: extractValueFn},
		properties.Scheme:            {name: "Scheme", transform: extractValueFn},
		properties.State:             {name: "State", transform: extractFieldFn("Code")},
		properties.Type:              {name: "Type", transform: extractValueFn},
		properties.Vpc:               {name: "VpcId", transform: extractValueFn},
	},
	cloud.TargetGroup: {
		properties.Name:                    {name: "TargetGroupName", transform: extractValueFn},
		properties.CheckInterval:           {name: "HealthCheckIntervalSeconds", transform: extractValueFn},
		properties.CheckPath:               {name: "HealthCheckPath", transform: extractValueFn},
		properties.CheckPort:               {name: "HealthCheckPort", transform: extractValueFn},
		properties.CheckProtocol:           {name: "HealthCheckProtocol", transform: extractValueFn},
		properties.CheckTimeout:            {name: "HealthCheckTimeoutSeconds", transform: extractValueFn},
		properties.CheckHTTPCode:           {name: "Matcher", transform: extractFieldFn("HttpCode")},
		properties.HealthyThresholdCount:   {name: "HealthyThresholdCount", transform: extractValueFn},
		properties.UnhealthyThresholdCount: {name: "UnhealthyThresholdCount", transform: extractValueFn},
		properties.Port:                    {name: "Port", transform: extractValueFn},
		properties.Protocol:                {name: "Protocol", transform: extractValueFn},
		properties.Vpc:                     {name: "VpcId", transform: extractValueFn},
	},
	cloud.Listener: {
		properties.Certificates: {name: "Certificates", transform: extractStringSliceValues("CertificateArn")},
		properties.Actions:      {name: "DefaultActions", transform: extractStringSliceValues("Type")},
		properties.LoadBalancer: {name: "LoadBalancerArn", transform: extractValueFn},
		properties.Port:         {name: "Port", transform: extractValueFn},
		properties.Protocol:     {name: "Protocol", transform: extractValueFn},
		properties.CipherSuite:  {name: "SslPolicy", transform: extractValueFn},
	},
	//Database
	cloud.Database: {
		properties.Storage:                   {name: "AllocatedStorage", transform: extractValueFn},
		properties.AutoUpgrade:               {name: "AutoMinorVersionUpgrade", transform: extractValueFn},
		properties.AvailabilityZone:          {name: "AvailabilityZone", transform: extractValueFn},
		properties.BackupRetentionPeriod:     {name: "BackupRetentionPeriod", transform: extractValueFn},
		properties.CertificateAuthority:      {name: "CACertificateIdentifier", transform: extractValueFn},
		properties.Charset:                   {name: "CharacterSetName", transform: extractValueFn},
		properties.CopyTagsToSnapshot:        {name: "CopyTagsToSnapshot", transform: extractValueFn},
		properties.Cluster:                   {name: "DBClusterIdentifier", transform: extractValueFn},
		properties.Arn:                       {name: "DBInstanceArn", transform: extractValueFn},
		properties.Class:                     {name: "DBInstanceClass", transform: extractValueFn},
		properties.State:                     {name: "DBInstanceStatus", transform: extractValueFn},
		properties.Name:                      {name: "DBName", transform: extractValueFn},
		properties.ParameterGroups:           {name: "DBParameterGroups", transform: extractStringSliceValues("DBParameterGroupName")},
		properties.DBSecurityGroups:          {name: "DBSecurityGroups", transform: extractStringSliceValues("DBSecurityGroupName")},
		properties.DBSubnetGroup:             {name: "DBSubnetGroup", transform: extractFieldFn("DBSubnetGroupName")},
		properties.Port:                      {name: "DbInstancePort", transform: extractValueFn},
		properties.GlobalID:                  {name: "DbiResourceId", transform: extractValueFn},
		properties.PublicDNS:                 {name: "Endpoint", transform: extractFieldFn("Address")},
		properties.Zone:                      {name: "Endpoint", transform: extractFieldFn("HostedZoneId")},
		properties.Engine:                    {name: "Engine", transform: extractValueFn},
		properties.EngineVersion:             {name: "EngineVersion", transform: extractValueFn},
		properties.Created:                   {name: "InstanceCreateTime", transform: extractValueFn},
		properties.IOPS:                      {name: "Iops", transform: extractValueFn},
		properties.LatestRestorableTime:      {name: "LatestRestorableTime", transform: extractValueFn},
		properties.License:                   {name: "LicenseModel", transform: extractValueFn},
		properties.Username:                  {name: "MasterUsername", transform: extractValueFn},
		properties.MonitoringInterval:        {name: "MonitoringInterval", transform: extractValueFn},
		properties.MonitoringRole:            {name: "MonitoringRoleArn", transform: extractValueFn},
		properties.MultiAZ:                   {name: "MultiAZ", transform: extractValueFn},
		properties.OptionGroups:              {name: "OptionGroupMemberships", transform: extractStringSliceValues("OptionGroupName")},
		properties.PreferredBackupDate:       {name: "PreferredBackupWindow", transform: extractValueFn},
		properties.PreferredMaintenanceDate:  {name: "PreferredMaintenanceWindow", transform: extractValueFn},
		properties.Public:                    {name: "PubliclyAccessible", transform: extractValueFn},
		properties.SecondaryAvailabilityZone: {name: "SecondaryAvailabilityZone", transform: extractValueFn},
		properties.Encrypted:                 {name: "StorageEncrypted", transform: extractValueFn},
		properties.StorageType:               {name: "StorageType", transform: extractValueFn},
		properties.Timezone:                  {name: "Timezone", transform: extractValueFn},
		properties.SecurityGroups:            {name: "VpcSecurityGroups", transform: extractStringSliceValues("VpcSecurityGroupId")},
	},
	cloud.DbSubnetGroup: {
		properties.Name:        {name: "DBSubnetGroupName", transform: extractValueFn},
		properties.Arn:         {name: "DBSubnetGroupArn", transform: extractValueFn},
		properties.Description: {name: "DBSubnetGroupDescription", transform: extractValueFn},
		properties.State:       {name: "SubnetGroupStatus", transform: extractValueFn},
		properties.Subnets:     {name: "Subnets", transform: extractStringSliceValues("SubnetIdentifier")},
		properties.Vpc:         {name: "VpcId", transform: extractValueFn},
	},
	//IAM
	cloud.User: {
		properties.Name:             {name: "UserName", transform: extractValueFn},
		properties.Arn:              {name: "Arn", transform: extractValueFn},
		properties.Path:             {name: "Path", transform: extractValueFn},
		properties.Created:          {name: "CreateDate", transform: extractTimeFn},
		properties.PasswordLastUsed: {name: "PasswordLastUsed", transform: extractTimeFn},
		properties.InlinePolicies:   {name: "UserPolicyList", transform: extractStringSliceValues("PolicyName")},
	},
	cloud.Role: {
		properties.Name:           {name: "RoleName", transform: extractValueFn},
		properties.Arn:            {name: "Arn", transform: extractValueFn},
		properties.Created:        {name: "CreateDate", transform: extractTimeFn},
		properties.Path:           {name: "Path", transform: extractValueFn},
		properties.InlinePolicies: {name: "RolePolicyList", transform: extractStringSliceValues("PolicyName")},
	},
	cloud.Group: {
		properties.Name:           {name: "GroupName", transform: extractValueFn},
		properties.Arn:            {name: "Arn", transform: extractValueFn},
		properties.Created:        {name: "CreateDate", transform: extractTimeFn},
		properties.Path:           {name: "Path", transform: extractValueFn},
		properties.InlinePolicies: {name: "GroupPolicyList", transform: extractStringSliceValues("PolicyName")},
	},
	cloud.Policy: {
		properties.Name:        {name: "PolicyName", transform: extractValueFn},
		properties.Arn:         {name: "Arn", transform: extractValueFn},
		properties.Created:     {name: "CreateDate", transform: extractTimeFn},
		properties.Updated:     {name: "UpdateDate", transform: extractTimeFn},
		properties.Description: {name: "Description", transform: extractValueFn},
		properties.Attachable:  {name: "IsAttachable", transform: extractValueFn},
		properties.Path:        {name: "Path", transform: extractValueFn},
	},
	//S3
	cloud.Bucket: {
		properties.Created: {name: "CreationDate", transform: extractTimeFn},
		properties.Grants:  {fetch: fetchAndExtractGrantsFn},
	},
	cloud.Object: {
		properties.Key:      {name: "Key", transform: extractValueFn},
		properties.Modified: {name: "LastModified", transform: extractTimeFn},
		properties.Owner:    {name: "Owner", transform: extractFieldFn("ID")},
		properties.Size:     {name: "Size", transform: extractValueFn},
		properties.Class:    {name: "StorageClass", transform: extractValueFn},
	},
	//Notification
	cloud.Subscription: {
		properties.Endpoint: {name: "Endpoint", transform: extractValueFn},
		properties.Owner:    {name: "Owner", transform: extractValueFn},
		properties.Protocol: {name: "Protocol", transform: extractValueFn},
		properties.Arn:      {name: "SubscriptionArn", transform: extractValueFn},
		properties.Topic:    {name: "TopicArn", transform: extractValueFn},
	},
	cloud.Topic: {
		properties.Arn: {name: "TopicArn", transform: extractValueFn},
	},
	// DNS
	cloud.Zone: {
		properties.Name:            {name: "Name", transform: extractValueFn},
		properties.Comment:         {name: "Config", transform: extractFieldFn("Comment")},
		properties.Private:         {name: "Config", transform: extractFieldFn("PrivateZone")},
		properties.CallerReference: {name: "CallerReference", transform: extractValueFn},
		properties.RecordCount:     {name: "ResourceRecordSetCount", transform: extractValueFn},
	},
	cloud.Record: {
		properties.Name:                  {name: "Name", transform: extractValueFn},
		properties.Failover:              {name: "Failover", transform: extractValueFn},
		properties.Continent:             {name: "GeoLocation", transform: extractFieldFn("ContinentCode")},
		properties.Country:               {name: "GeoLocation", transform: extractFieldFn("CountryCode")},
		properties.HealthCheck:           {name: "HealthCheckId", transform: extractValueFn},
		properties.Region:                {name: "Region", transform: extractValueFn},
		properties.Records:               {name: "ResourceRecords", transform: extractStringSliceValues("Value")},
		properties.Set:                   {name: "SetIdentifier", transform: extractValueFn},
		properties.TTL:                   {name: "TTL", transform: extractValueFn},
		properties.TrafficPolicyInstance: {name: "TrafficPolicyInstanceId", transform: extractValueFn},
		properties.Type:                  {name: "Type", transform: extractValueFn},
		properties.Weight:                {name: "Weight", transform: extractValueFn},
	},
	//Queue
	cloud.Queue: {}, //Manually set
}
