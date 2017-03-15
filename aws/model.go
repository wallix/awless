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

var awsResourcesDef = map[string]map[string]*propertyTransform{
	//EC2
	cloud.Instance: {
		"Name":           {name: "Tags", transform: extractTagFn("Name")},
		"Type":           {name: "InstanceType", transform: extractValueFn},
		"SubnetId":       {name: "SubnetId", transform: extractValueFn},
		"VpcId":          {name: "VpcId", transform: extractValueFn},
		"PublicIp":       {name: "PublicIpAddress", transform: extractValueFn},
		"PrivateIp":      {name: "PrivateIpAddress", transform: extractValueFn},
		"ImageId":        {name: "ImageId", transform: extractValueFn},
		"LaunchTime":     {name: "LaunchTime", transform: extractValueFn},
		"State":          {name: "State", transform: extractFieldFn("Name")},
		"KeyName":        {name: "KeyName", transform: extractValueFn},
		"SecurityGroups": {name: "SecurityGroups", transform: extractSliceValues("GroupId")},
	},
	cloud.Vpc: {
		"Name":      {name: "Tags", transform: extractTagFn("Name")},
		"IsDefault": {name: "IsDefault", transform: extractValueFn},
		"State":     {name: "State", transform: extractValueFn},
		"CidrBlock": {name: "CidrBlock", transform: extractValueFn},
	},
	cloud.Subnet: {
		"Name":                {name: "Tags", transform: extractTagFn("Name")},
		"VpcId":               {name: "VpcId", transform: extractValueFn},
		"MapPublicIpOnLaunch": {name: "MapPublicIpOnLaunch", transform: extractValueFn},
		"State":               {name: "State", transform: extractValueFn},
		"CidrBlock":           {name: "CidrBlock", transform: extractValueFn},
		"AvailabilityZone":    {name: "AvailabilityZone", transform: extractValueFn},
		"DefaultForAz":        {name: "DefaultForAz", transform: extractValueFn},
	},
	cloud.SecurityGroup: {
		"Name":          {name: "GroupName", transform: extractValueFn},
		"Description":   {name: "Description", transform: extractValueFn},
		"InboundRules":  {name: "IpPermissions", transform: extractIpPermissionSliceFn},
		"OutboundRules": {name: "IpPermissionsEgress", transform: extractIpPermissionSliceFn},
		"OwnerId":       {name: "OwnerId", transform: extractValueFn},
		"VpcId":         {name: "VpcId", transform: extractValueFn},
	},
	cloud.Keypair: {
		"Name":           {name: "KeyName", transform: extractValueFn},
		"KeyFingerprint": {name: "KeyFingerprint", transform: extractValueFn},
	},
	cloud.Volume: {
		"Name":             {name: "Tags", transform: extractTagFn("Name")},
		"VolumeType":       {name: "VolumeType", transform: extractValueFn},
		"State":            {name: "State", transform: extractValueFn},
		"Size":             {name: "Size", transform: extractValueFn},
		"Encrypted":        {name: "Encrypted", transform: extractValueFn},
		"CreateTime":       {name: "CreateTime", transform: extractTimeFn},
		"AvailabilityZone": {name: "AvailabilityZone", transform: extractValueFn},
	},
	cloud.InternetGateway: {
		"Name": {name: "Tags", transform: extractTagFn("Name")},
		"Vpcs": {name: "Attachments", transform: extractSliceValues("VpcId")},
	},
	cloud.RouteTable: {
		"Name":   {name: "Tags", transform: extractTagFn("Name")},
		"VpcId":  {name: "VpcId", transform: extractValueFn},
		"Routes": {name: "Routes", transform: extractRoutesSliceFn},
		"Main":   {name: "Associations", transform: extractHasATrueBoolInStructSliceFn("Main")},
	},
	cloud.AvailabilityZone: {
		"Name":     {name: "ZoneName", transform: extractValueFn},
		"State":    {name: "State", transform: extractValueFn},
		"Region":   {name: "RegionName", transform: extractValueFn},
		"Messages": {name: "Messages", transform: extractSliceValues("Message")},
	},
	// LoadBalancer
	cloud.LoadBalancer: {
		"Name":                  {name: "LoadBalancerName", transform: extractValueFn},
		"AvailabilityZones":     {name: "AvailabilityZones", transform: extractSliceValues("ZoneName")},
		"Subnets":               {name: "AvailabilityZones", transform: extractSliceValues("SubnetId")},
		"CanonicalHostedZoneId": {name: "CanonicalHostedZoneId", transform: extractValueFn},
		"CreateTime":            {name: "CreatedTime", transform: extractTimeFn},
		"DNSName":               {name: "DNSName", transform: extractValueFn},
		"IpAddressType":         {name: "IpAddressType", transform: extractValueFn},
		"Scheme":                {name: "Scheme", transform: extractValueFn},
		"State":                 {name: "State", transform: extractFieldFn("Code")},
		"Type":                  {name: "Type", transform: extractValueFn},
		"VpcId":                 {name: "VpcId", transform: extractValueFn},
	},
	cloud.TargetGroup: {
		"Name": {name: "TargetGroupName", transform: extractValueFn},
		"HealthCheckIntervalSeconds": {name: "HealthCheckIntervalSeconds", transform: extractValueFn},
		"HealthCheckPath":            {name: "HealthCheckPath", transform: extractValueFn},
		"HealthCheckPort":            {name: "HealthCheckPort", transform: extractValueFn},
		"HealthCheckProtocol":        {name: "HealthCheckProtocol", transform: extractValueFn},
		"HealthCheckTimeoutSeconds":  {name: "HealthCheckTimeoutSeconds", transform: extractValueFn},
		"HealthyThresholdCount":      {name: "HealthyThresholdCount", transform: extractValueFn},
		"Matcher":                    {name: "Matcher", transform: extractFieldFn("HttpCode")},
		"Port":                       {name: "Port", transform: extractValueFn},
		"Protocol":                   {name: "Protocol", transform: extractValueFn},
		"UnhealthyThresholdCount":    {name: "UnhealthyThresholdCount", transform: extractValueFn},
		"VpcId":                      {name: "VpcId", transform: extractValueFn},
	},
	cloud.Listener: {
		"Certificates": {name: "Certificates", transform: extractSliceValues("CertificateArn")},
		"Actions":      {name: "DefaultActions", transform: extractSliceValues("Type")},
		"LoadBalancer": {name: "LoadBalancerArn", transform: extractValueFn},
		"Port":         {name: "Port", transform: extractValueFn},
		"Protocol":     {name: "Protocol", transform: extractValueFn},
		"SslPolicy":    {name: "SslPolicy", transform: extractValueFn},
	},
	//IAM
	cloud.User: {
		"Name":                 {name: "UserName", transform: extractValueFn},
		"Arn":                  {name: "Arn", transform: extractValueFn},
		"Path":                 {name: "Path", transform: extractValueFn},
		"CreateDate":           {name: "CreateDate", transform: extractTimeFn},
		"PasswordLastUsedDate": {name: "PasswordLastUsed", transform: extractTimeFn},
		"InlinePolicies":       {name: "UserPolicyList", transform: extractSliceValues("PolicyName")},
	},
	cloud.Role: {
		"Name":           {name: "RoleName", transform: extractValueFn},
		"Arn":            {name: "Arn", transform: extractValueFn},
		"CreateDate":     {name: "CreateDate", transform: extractTimeFn},
		"Path":           {name: "Path", transform: extractValueFn},
		"InlinePolicies": {name: "RolePolicyList", transform: extractSliceValues("PolicyName")},
	},
	cloud.Group: {
		"Name":           {name: "GroupName", transform: extractValueFn},
		"Arn":            {name: "Arn", transform: extractValueFn},
		"CreateDate":     {name: "CreateDate", transform: extractTimeFn},
		"Path":           {name: "Path", transform: extractValueFn},
		"InlinePolicies": {name: "GroupPolicyList", transform: extractSliceValues("PolicyName")},
	},
	cloud.Policy: {
		"Name":         {name: "PolicyName", transform: extractValueFn},
		"Arn":          {name: "Arn", transform: extractValueFn},
		"CreateDate":   {name: "CreateDate", transform: extractTimeFn},
		"UpdateDate":   {name: "UpdateDate", transform: extractTimeFn},
		"Description":  {name: "Description", transform: extractValueFn},
		"IsAttachable": {name: "IsAttachable", transform: extractValueFn},
		"Path":         {name: "Path", transform: extractValueFn},
	},
	//S3
	cloud.Bucket: {
		"Name":       {name: "Name", transform: extractValueFn},
		"CreateDate": {name: "CreationDate", transform: extractTimeFn},
		"Grants":     {fetch: fetchAndExtractGrantsFn},
	},
	cloud.Object: {
		"Key":          {name: "Key", transform: extractValueFn},
		"ModifiedDate": {name: "LastModified", transform: extractTimeFn},
		"OwnerId":      {name: "Owner", transform: extractFieldFn("ID")},
		"Size":         {name: "Size", transform: extractValueFn},
		"Class":        {name: "StorageClass", transform: extractValueFn},
	},
	//Notification
	cloud.Subscription: {
		"Endpoint":        {name: "Endpoint", transform: extractValueFn},
		"Owner":           {name: "Owner", transform: extractValueFn},
		"Protocol":        {name: "Protocol", transform: extractValueFn},
		"SubscriptionArn": {name: "SubscriptionArn", transform: extractValueFn},
		"TopicArn":        {name: "TopicArn", transform: extractValueFn},
	},
	cloud.Topic: {
		"TopicArn": {name: "TopicArn", transform: extractValueFn},
	},
	// DNS
	cloud.Zone: {
		"Name":                   {name: "Name", transform: extractValueFn},
		"Comment":                {name: "Config", transform: extractFieldFn("Comment")},
		"IsPrivateZone":          {name: "Config", transform: extractFieldFn("PrivateZone")},
		"CallerReference":        {name: "CallerReference", transform: extractValueFn},
		"ResourceRecordSetCount": {name: "ResourceRecordSetCount", transform: extractValueFn},
	},
	cloud.Record: {
		"Name":          {name: "Name", transform: extractValueFn},
		"Failover":      {name: "Failover", transform: extractValueFn},
		"Continent":     {name: "GeoLocation", transform: extractFieldFn("ContinentCode")},
		"Country":       {name: "GeoLocation", transform: extractFieldFn("CountryCode")},
		"HealthCheckId": {name: "HealthCheckId", transform: extractValueFn},
		"Region":        {name: "Region", transform: extractValueFn},
		"Records":       {name: "ResourceRecords", transform: extractSliceValues("Value")},
		"SetIdentifier": {name: "SetIdentifier", transform: extractValueFn},
		"TTL":           {name: "TTL", transform: extractValueFn},
		"TrafficPolicyInstanceId": {name: "TrafficPolicyInstanceId", transform: extractValueFn},
		"Type":   {name: "Type", transform: extractValueFn},
		"Weight": {name: "Weight", transform: extractValueFn},
	},
	//Queue
	cloud.Queue: {}, //Manually set
}
