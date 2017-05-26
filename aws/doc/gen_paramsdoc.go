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
package awsdoc

var generatedParamsDoc = map[string]map[string]string{
	"attachalarm": {},
	"attachelasticip": {
		"allow-reassociation": "[EC2-VPC] For a VPC in an EC2-Classic account, specify true to allow an Elastic IP address that is already associated with an instance or network interface to be reassociated with the specified instance or network interface",
		"id":               "[EC2-VPC] The allocation ID",
		"instance":         "The ID of the instance",
		"networkinterface": "[EC2-VPC] The ID of the network interface",
		"privateip":        "[EC2-VPC] The primary or secondary private IP address to associate with the Elastic IP address",
	},
	"attachinstance": {
		"targetgroup": "The Amazon Resource Name (ARN) of the target group",
	},
	"attachinternetgateway": {
		"id":  "The ID of the Internet gateway",
		"vpc": "The ID of the VPC",
	},
	"attachpolicy": {},
	"attachrole": {
		"instanceprofile": "The name of the instance profile to update",
		"name":            "The name of the role to add",
	},
	"attachroutetable": {
		"id":     "The ID of the route table",
		"subnet": "The ID of the subnet",
	},
	"attachsecuritygroup": {},
	"attachuser": {
		"group": "The name of the group to update",
		"name":  "The name of the user to add",
	},
	"attachvolume": {
		"device":   "The device name to expose to the instance (for example, /dev/sdh or xvdh)",
		"id":       "The ID of the EBS volume",
		"instance": "The ID of the instance",
	},
	"checkdistribution":  {},
	"checkinstance":      {},
	"checkloadbalancer":  {},
	"checkscalinggroup":  {},
	"checksecuritygroup": {},
	"copyimage": {
		"description":   "A description for the new AMI in the destination region",
		"encrypted":     "Specifies whether the destination snapshots of the copied image should be encrypted",
		"name":          "The name of the new AMI in the destination region",
		"source-id":     "The ID of the AMI to copy",
		"source-region": "The name of the region that contains the AMI to copy",
	},
	"copysnapshot": {
		"description":   "A description for the EBS snapshot",
		"encrypted":     "Specifies whether the destination snapshot should be encrypted",
		"source-id":     "The ID of the EBS snapshot to copy",
		"source-region": "The ID of the region that contains the snapshot to be copied",
	},
	"createaccesskey": {},
	"createalarm": {
		"alarm-actions":            "The actions to execute when this alarm transitions to the ALARM state from any other state",
		"description":              "The description for the alarm",
		"dimensions":               "The dimensions for the metric associated with the alarm",
		"enabled":                  "Indicates whether actions should be executed during any changes to the alarm state",
		"evaluation-periods":       "The number of periods over which data is compared to the specified threshold",
		"insufficientdata-actions": "The actions to execute when this alarm transitions to the INSUFFICIENT_DATA state from any other state",
		"metric":                   "The name for the metric associated with the alarm",
		"name":                     "The name for the alarm",
		"namespace":                "The namespace for the metric associated with the alarm",
		"ok-actions":               "The actions to execute when this alarm transitions to an OK state from any other state",
		"operator":                 " The arithmetic operation to use when comparing the specified statistic and threshold",
		"period":                   "The period, in seconds, over which the specified statistic is applied",
		"statistic-function":       "The statistic for the metric associated with the alarm, other than percentile",
		"threshold":                "The value against which the specified statistic is compared",
		"unit":                     "The unit of measure for the statistic",
	},
	"createbucket": {
		"acl":  "The canned ACL to apply to the bucket",
		"name": "",
	},
	"createdatabase":      {},
	"createdbsubnetgroup": {},
	"createdistribution":  {},
	"createelasticip": {
		"domain": "Set to vpc to allocate the address for use with instances in a VPC",
	},
	"createfunction": {
		"description": "A short, user-defined function description",
		"handler":     "The function within your code that Lambda calls to begin execution",
		"memory":      "The amount of memory, in MB, your Lambda function is given",
		"name":        "The name you want to assign to the function you are uploading",
		"publish":     "This boolean parameter can be used to request AWS Lambda to create the Lambda function and publish a version as an atomic operation",
		"role":        "The Amazon Resource Name (ARN) of the IAM role that Lambda assumes when it executes your function to access any other Amazon Web Services (AWS) resources",
		"runtime":     "The runtime environment for the Lambda function you are uploading",
		"timeout":     "The function execution time at which Lambda should terminate the function",
	},
	"creategroup": {
		"name": "The name of the group to create",
	},
	"createinstance": {
		"count":         "The minimum number of instances to launch",
		"image":         "The ID of the AMI, which you can get by calling DescribeImages",
		"ip":            "[EC2-VPC] The primary IPv4 address",
		"keypair":       "The name of the key pair",
		"lock":          "If you set this parameter to true, you can't terminate the instance using the Amazon EC2 console, CLI, or API; otherwise, you can",
		"securitygroup": "One or more security group IDs",
		"subnet":        "[EC2-VPC] The ID of the subnet to launch the instance into",
		"type":          "The instance type",
		"userdata":      "The user data to make available to the instance",
	},
	"createinstanceprofile": {
		"name": "The name of the instance profile to create",
	},
	"createinternetgateway":     {},
	"createkeypair":             {},
	"createlaunchconfiguration": {},
	"createlistener": {
		"loadbalancer": "The Amazon Resource Name (ARN) of the load balancer",
		"port":         "The port on which the load balancer is listening",
		"protocol":     "The protocol for connections from clients to the load balancer",
		"sslpolicy":    "The security policy that defines which ciphers and protocols are supported",
	},
	"createloadbalancer": {
		"iptype":         "The type of IP addresses used by the subnets for your load balancer",
		"name":           "The name of the load balancer",
		"scheme":         "The nodes of an Internet-facing load balancer have public IP addresses",
		"securitygroups": "The IDs of the security groups to assign to the load balancer",
		"subnets":        "The IDs of the subnets to attach to the load balancer",
	},
	"createloginprofile": {
		"password":       "The new password for the user",
		"password-reset": "Specifies whether the user is required to set a new password on next sign-in",
		"username":       "The name of the IAM user to create a password for",
	},
	"createpolicy": {},
	"createqueue": {
		"name": "The name of the new queue",
	},
	"createrecord": {},
	"createrole":   {},
	"createroute": {
		"cidr":    "The IPv4 CIDR address block used for the destination match",
		"gateway": "The ID of an Internet gateway or virtual private gateway attached to your VPC",
		"table":   "The ID of the route table for the route",
	},
	"createroutetable": {
		"vpc": "The ID of the VPC",
	},
	"creates3object":      {},
	"createscalinggroup":  {},
	"createscalingpolicy": {},
	"createsecuritygroup": {
		"description": "A description for the security group",
		"name":        "The name of the security group",
		"vpc":         "[EC2-VPC] The ID of the VPC",
	},
	"createsnapshot": {
		"description": "A description for the snapshot",
		"volume":      "The ID of the EBS volume",
	},
	"createstack": {},
	"createsubnet": {
		"availabilityzone": "The Availability Zone for the subnet",
		"cidr":             "The IPv4 network range for the subnet, in CIDR notation",
		"vpc":              "The ID of the VPC",
	},
	"createsubscription": {
		"endpoint": "The endpoint that you want to receive notifications",
		"protocol": "The protocol you want to use",
		"topic":    "The ARN of the topic you want to subscribe to",
	},
	"createtag": {},
	"createtargetgroup": {
		"healthcheckinterval": "The approximate amount of time, in seconds, between health checks of an individual target",
		"healthcheckpath":     "The ping path that is the destination on the targets for health checks",
		"healthcheckport":     "The port the load balancer uses when performing health checks on targets",
		"healthcheckprotocol": "The protocol the load balancer uses when performing health checks on targets",
		"healthchecktimeout":  "The amount of time, in seconds, during which no response from a target means a failed health check",
		"healthythreshold":    "The number of consecutive health checks successes required before considering an unhealthy target healthy",
		"name":                "The name of the target group",
		"port":                "The port on which the targets receive traffic",
		"protocol":            "The protocol to use for routing traffic to the targets",
		"unhealthythreshold":  "The number of consecutive health check failures required before considering a target unhealthy",
		"vpc":                 "The identifier of the virtual private cloud (VPC)",
	},
	"createtopic": {
		"name": "The name of the topic you want to create",
	},
	"createuser": {
		"name": "The name of the user to create",
	},
	"createvolume": {
		"availabilityzone": "The Availability Zone in which to create the volume",
		"size":             "The size of the volume, in GiBs",
	},
	"createvpc": {
		"cidr": "The IPv4 network range for the VPC, in CIDR notation",
	},
	"createzone": {
		"callerreference": "A unique string that identifies the request and that allows failed CreateHostedZone requests to be retried without the risk of executing the operation twice",
		"delegationsetid": "If you want to associate a reusable delegation set with this hosted zone, the ID that Amazon Route 53 assigned to the reusable delegation set when you created it",
		"name":            "The name of the domain",
	},
	"deleteaccesskey": {
		"id": "The access key ID for the access key ID and secret access key you want to delete",
	},
	"deletealarm": {
		"name": "The alarms to be deleted",
	},
	"deletebucket": {
		"name": "",
	},
	"deletedatabase":      {},
	"deletedbsubnetgroup": {},
	"deletedistribution":  {},
	"deleteelasticip": {
		"id": "[EC2-VPC] The allocation ID",
		"ip": "[EC2-Classic] The Elastic IP address",
	},
	"deletefunction": {
		"id":      "The Lambda function to delete",
		"version": "Using this optional parameter you can specify a function version (but not the $LATEST version) to direct AWS Lambda to delete a specific function version",
	},
	"deletegroup": {
		"name": "The name of the IAM group to delete",
	},
	"deleteimage": {},
	"deleteinstance": {
		"id": "One or more instance IDs",
	},
	"deleteinstanceprofile": {
		"name": "The name of the instance profile to delete",
	},
	"deleteinternetgateway": {
		"id": "The ID of the Internet gateway",
	},
	"deletekeypair": {
		"id": "The name of the key pair",
	},
	"deletelaunchconfiguration": {},
	"deletelistener": {
		"id": "The Amazon Resource Name (ARN) of the listener",
	},
	"deleteloadbalancer": {
		"id": "The Amazon Resource Name (ARN) of the load balancer",
	},
	"deleteloginprofile": {
		"username": "The name of the user whose password you want to delete",
	},
	"deletepolicy": {
		"arn": "The Amazon Resource Name (ARN) of the IAM policy you want to delete",
	},
	"deletequeue": {
		"url": "The URL of the Amazon SQS queue to delete",
	},
	"deleterecord": {},
	"deleterole":   {},
	"deleteroute": {
		"cidr":  "The IPv4 CIDR range for the route",
		"table": "The ID of the route table",
	},
	"deleteroutetable": {
		"id": "The ID of the route table",
	},
	"deletes3object": {
		"bucket": "",
		"name":   "",
	},
	"deletescalinggroup":  {},
	"deletescalingpolicy": {},
	"deletesecuritygroup": {
		"id": "The ID of the security group",
	},
	"deletesnapshot": {
		"id": "The ID of the EBS snapshot",
	},
	"deletestack": {},
	"deletesubnet": {
		"id": "The ID of the subnet",
	},
	"deletesubscription": {
		"id": "The ARN of the subscription to be deleted",
	},
	"deletetag": {},
	"deletetargetgroup": {
		"id": "The Amazon Resource Name (ARN) of the target group",
	},
	"deletetopic": {
		"id": "The ARN of the topic you want to delete",
	},
	"deleteuser": {
		"name": "The name of the user to delete",
	},
	"deletevolume": {
		"id": "The ID of the volume",
	},
	"deletevpc": {
		"id": "The ID of the VPC",
	},
	"deletezone": {
		"id": "The ID of the hosted zone you want to delete",
	},
	"detachalarm": {},
	"detachelasticip": {
		"association": "[EC2-VPC] The association ID",
	},
	"detachinstance": {
		"targetgroup": "The Amazon Resource Name (ARN) of the target group",
	},
	"detachinternetgateway": {
		"id":  "The ID of the Internet gateway",
		"vpc": "The ID of the VPC",
	},
	"detachpolicy": {},
	"detachrole": {
		"instanceprofile": "The name of the instance profile to update",
		"name":            "The name of the role to remove",
	},
	"detachroutetable": {
		"association": "The association ID representing the current association between the route table and subnet",
	},
	"detachsecuritygroup": {},
	"detachuser": {
		"group": "The name of the group to update",
		"name":  "The name of the user to remove",
	},
	"detachvolume": {
		"device":   "The device name",
		"force":    "Forces detachment if the previous detachment attempt did not occur cleanly (for example, logging into an instance, unmounting the volume, and detaching normally)",
		"id":       "The ID of the volume",
		"instance": "The ID of the instance",
	},
	"importimage": {
		"architecture": "The architecture of the virtual machine",
		"description":  "A description string for the import image task",
		"license":      "The license type to be used for the Amazon Machine Image (AMI) after importing",
		"platform":     "The operating system of the virtual machine",
		"role":         "The name of the role to use when not using the default role, 'vmimport'",
	},
	"startalarm": {
		"names": "The names of the alarms",
	},
	"startinstance": {
		"id": "One or more instance IDs",
	},
	"stopalarm": {
		"names": "The names of the alarms",
	},
	"stopinstance": {
		"id": "One or more instance IDs",
	},
	"updatebucket":       {},
	"updatedistribution": {},
	"updateinstance": {
		"id":   "The ID of the instance",
		"lock": "If the value is true, you can't terminate the instance using the Amazon EC2 console, CLI, or API; otherwise, you can",
		"type": "Changes the instance type to the specified value",
	},
	"updateloginprofile": {
		"password":       "The new password for the specified IAM user",
		"password-reset": "Allows this new password to be used only once by requiring the specified IAM user to set a new password on next sign-in",
		"username":       "The name of the user whose password you want to update",
	},
	"updates3object": {
		"acl":     "The canned ACL to apply to the object",
		"bucket":  "",
		"name":    "",
		"version": "VersionId used to reference a specific version of the object",
	},
	"updatescalinggroup":  {},
	"updatesecuritygroup": {},
	"updatesubnet": {
		"id":     "The ID of the subnet",
		"public": "Specify true to indicate that network interfaces created in the specified subnet should be assigned a public IPv4 address",
	},
}
