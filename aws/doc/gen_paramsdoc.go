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
	"attachalarm":         {},
	"attachcontainertask": {},
	"attachelasticip": {
		"allow-reassociation": "For a VPC in an EC2-Classic account, specify true to allow an Elastic IP address that is already associated with an instance or network interface to be reassociated with the specified instance or network interface",
		"id":               "The allocation ID",
		"instance":         "The ID of the instance",
		"networkinterface": "The ID of the network interface",
		"privateip":        "The primary or secondary private IP address to associate with the Elastic IP address",
	},
	"attachinstance": {
		"targetgroup": "The Amazon Resource Name (ARN) of the target group",
	},
	"attachinstanceprofile": {},
	"attachinternetgateway": {
		"id":  "The ID of the Internet gateway",
		"vpc": "The ID of the VPC",
	},
	"attachmfadevice": {
		"id":         "The serial number that uniquely identifies the MFA device",
		"mfa-code-1": "An authentication code emitted by the device",
		"mfa-code-2": "A subsequent authentication code emitted by the device",
		"user":       "The name of the IAM user for whom you want to enable the MFA device",
	},
	"attachnetworkinterface": {
		"device-index": "The index of the device for the network interface attachment",
		"id":           "The ID of the network interface",
		"instance":     "The ID of the instance",
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
	"authenticateregistry":  {},
	"checkcertificate":      {},
	"checkdatabase":         {},
	"checkdistribution":     {},
	"checkinstance":         {},
	"checkloadbalancer":     {},
	"checknatgateway":       {},
	"checknetworkinterface": {},
	"checkscalinggroup":     {},
	"checksecuritygroup":    {},
	"checkvolume":           {},
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
	"createaccesskey": {
		"user": "The name of the IAM user that the new key will belong to",
	},
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
		"operator":                 "The arithmetic operation to use when comparing the specified statistic and threshold",
		"period":                   "The period, in seconds, over which the specified statistic is applied",
		"statistic-function":       "The statistic for the metric associated with the alarm, other than percentile",
		"threshold":                "The value against which the specified statistic is compared",
		"unit":                     "The unit of measure for the statistic",
	},
	"createappscalingpolicy": {
		"dimension":         "The scalable dimension",
		"name":              "The name of the scaling policy",
		"resource":          "The identifier of the resource associated with the scaling policy",
		"service-namespace": "The namespace of the AWS service",
		"type":              "The policy type",
	},
	"createappscalingtarget": {
		"dimension":         "The scalable dimension associated with the scalable target",
		"max-capacity":      "The maximum value to scale to in response to a scale out event",
		"min-capacity":      "The minimum value to scale to in response to a scale in event",
		"resource":          "The identifier of the resource associated with the scalable target",
		"role":              "The ARN of an IAM role that allows Application Auto Scaling to modify the scalable target on your behalf",
		"service-namespace": "The namespace of the AWS service",
	},
	"createbucket": {
		"acl":  "The canned ACL to apply to the bucket",
		"name": "",
	},
	"createcertificate": {},
	"createcontainercluster": {
		"name": "The name of your cluster",
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
		"image":         "The ID of the AMI, which you can get by calling DescribeImages",
		"ip":            "The primary IPv4 address",
		"keypair":       "The name of the key pair",
		"lock":          "If you set this parameter to true, you can't terminate the instance using the Amazon EC2 console, CLI, or API; otherwise, you can",
		"securitygroup": "One or more security group IDs",
		"subnet":        "The ID of the subnet to launch the instance into",
		"type":          "The instance type",
		"userdata":      "The user data to make available to the instance",
	},
	"createinstanceprofile": {
		"name": "The name of the instance profile to create",
	},
	"createinternetgateway": {},
	"createkeypair": {
		"name": "A unique name for the key pair",
	},
	"createlaunchconfiguration": {
		"image":          "The ID of the Amazon Machine Image (AMI) to use to launch your EC2 instances",
		"keypair":        "The name of the key pair",
		"name":           "The name of the launch configuration",
		"public":         "Used for groups that launch instances into a virtual private cloud (VPC)",
		"role":           "The name or the Amazon Resource Name (ARN) of the instance profile associated with the IAM role for the instance",
		"securitygroups": "One or more security groups with which to associate the instances",
		"spotprice":      "The maximum hourly price to be paid for any Spot Instance launched to fulfill the request",
		"type":           "The instance type of the EC2 instance",
		"userdata":       "The user data to make available to the launched EC2 instances",
	},
	"createlistener": {
		"loadbalancer": "The Amazon Resource Name (ARN) of the load balancer",
		"port":         "The port on which the load balancer is listening",
		"protocol":     "The protocol for connections from clients to the load balancer",
		"sslpolicy":    "[HTTPS listeners] The security policy that defines which ciphers and protocols are supported",
	},
	"createloadbalancer": {
		"iptype":          "[Application Load Balancers] The type of IP addresses used by the subnets for your load balancer",
		"name":            "The name of the load balancer",
		"scheme":          "The nodes of an Internet-facing load balancer have public IP addresses",
		"securitygroups":  "[Application Load Balancers] The IDs of the security groups to assign to the load balancer",
		"subnet-mappings": "The IDs of the subnets to attach to the load balancer",
		"subnets":         "The IDs of the subnets to attach to the load balancer",
		"type":            "The type of load balancer to create",
	},
	"createloginprofile": {
		"password":       "The new password for the user",
		"password-reset": "Specifies whether the user is required to set a new password on next sign-in",
		"username":       "The name of the IAM user to create a password for",
	},
	"createmfadevice": {},
	"createnatgateway": {
		"elasticip-id": "The allocation ID of an Elastic IP address to associate with the NAT gateway",
		"subnet":       "The subnet in which to create the NAT gateway",
	},
	"createnetworkinterface": {
		"description":    "A description for the network interface",
		"privateip":      "The primary private IPv4 address of the network interface",
		"securitygroups": "The IDs of one or more security groups",
		"subnet":         "The ID of the subnet to associate with the network interface",
	},
	"createpolicy": {
		"description": "A friendly description of the policy",
		"name":        "The friendly name of the policy",
	},
	"createqueue": {
		"name": "The name of the new queue",
	},
	"createrecord": {},
	"createrepository": {
		"name": "The name to use for the repository",
	},
	"createrole": {},
	"createroute": {
		"cidr":    "The IPv4 CIDR address block used for the destination match",
		"gateway": "The ID of an Internet gateway or virtual private gateway attached to your VPC",
		"table":   "The ID of the route table for the route",
	},
	"createroutetable": {
		"vpc": "The ID of the VPC",
	},
	"creates3object": {},
	"createscalinggroup": {
		"cooldown":                 "The amount of time, in seconds, after a scaling activity completes before another scaling activity can start",
		"desired-capacity":         "The number of EC2 instances that should be running in the group",
		"healthcheck-grace-period": "The amount of time, in seconds, that Auto Scaling waits before checking the health status of an EC2 instance that has come into service",
		"healthcheck-type":         "The service to use for the health checks",
		"launchconfiguration":      "The name of the launch configuration",
		"max-size":                 "The maximum size of the group",
		"min-size":                 "The minimum size of the group",
		"name":                     "The name of the group",
		"new-instances-protected": "Indicates whether newly launched instances are protected from termination by Auto Scaling when scaling in",
		"subnets":                 "A comma-separated list of subnet identifiers for your virtual private cloud (VPC)",
		"targetgroups":            "The Amazon Resource Names (ARN) of the target groups",
	},
	"createscalingpolicy": {
		"adjustment-magnitude": "The minimum number of instances to scale",
		"adjustment-scaling":   "The amount by which to scale, based on the specified adjustment type",
		"adjustment-type":      "The adjustment type",
		"cooldown":             "The amount of time, in seconds, after a scaling activity completes and before the next scaling activity can start",
		"name":                 "The name of the scaling policy",
		"scalinggroup":         "The name or ARN of the group",
	},
	"createsecuritygroup": {
		"description": "A description for the security group",
		"name":        "The name of the security group",
		"vpc":         "The ID of the VPC",
	},
	"createsnapshot": {
		"description": "A description for the snapshot",
		"volume":      "The ID of the EBS volume",
	},
	"createstack": {
		"capabilities":     "A list of values that you must specify before AWS CloudFormation can create certain stacks",
		"disable-rollback": "Set to true to disable rollback of the stack if stack creation failed",
		"name":             "The name that is associated with the stack",
		"notifications":    "The Simple Notification Service (SNS) topic ARNs to publish stack related events",
		"on-failure":       "Determines what action will be taken if stack creation fails",
		"parameters":       "A list of Parameter structures that specify input parameters for the stack",
		"policy-file":      "Structure containing the stack policy body",
		"resource-types":   "The template resource types that you have permissions to work with for this create stack action, such as AWS::EC2::Instance, AWS::EC2::*, or Custom::MyCustomInstance",
		"role":             "The Amazon Resource Name (ARN) of an AWS Identity and Access Management (IAM) role that AWS CloudFormation assumes to create the stack",
		"template-file":    "Structure containing the template body with a minimum length of 1 byte and a maximum length of 51,200 bytes",
		"timeout":          "The amount of time that can pass before the stack status becomes CREATE_FAILED; if DisableRollback is not set or is set to false, the stack will be rolled back",
	},
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
		"healthcheckpath":     "[HTTP/HTTPS health checks] The ping path that is the destination on the targets for health checks",
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
		"id":   "The access key ID for the access key ID and secret access key you want to delete",
		"user": "The name of the user whose access key pair you want to delete",
	},
	"deletealarm": {
		"name": "The alarms to be deleted",
	},
	"deleteappscalingpolicy": {
		"dimension":         "The scalable dimension",
		"name":              "The name of the scaling policy",
		"resource":          "The identifier of the resource associated with the scalable target",
		"service-namespace": "The namespace of the AWS service",
	},
	"deleteappscalingtarget": {
		"dimension":         "The scalable dimension associated with the scalable target",
		"resource":          "The identifier of the resource associated with the scalable target",
		"service-namespace": "The namespace of the AWS service",
	},
	"deletebucket": {
		"name": "",
	},
	"deletecertificate": {
		"arn": "String that contains the ARN of the ACM Certificate to be deleted",
	},
	"deletecontainercluster": {
		"id": "The short name or full Amazon Resource Name (ARN) of the cluster to delete",
	},
	"deletecontainertask": {},
	"deletedatabase": {
		"id": "Contains a user-supplied database identifier",
	},
	"deletedbsubnetgroup": {},
	"deletedistribution":  {},
	"deleteelasticip": {
		"id": "The allocation ID",
		"ip": "The Elastic IP address",
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
		"ids": "One or more instance IDs",
	},
	"deleteinstanceprofile": {
		"name": "The name of the instance profile to delete",
	},
	"deleteinternetgateway": {
		"id": "The ID of the Internet gateway",
	},
	"deletekeypair": {
		"name": "The name of the key pair",
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
	"deletemfadevice": {
		"id": "The serial number that uniquely identifies the MFA device",
	},
	"deletenatgateway": {
		"id": "The ID of the NAT gateway",
	},
	"deletenetworkinterface": {
		"id": "The ID of the network interface",
	},
	"deletepolicy": {
		"arn": "The Amazon Resource Name (ARN) of the IAM policy you want to delete",
	},
	"deletequeue": {
		"url": "The URL of the Amazon SQS queue to delete",
	},
	"deleterecord": {},
	"deleterepository": {
		"account": "The AWS account ID associated with the registry that contains the repository to delete",
		"force":   "Force the deletion of the repository if it contains images",
		"name":    "The name of the repository to delete",
	},
	"deleterole": {},
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
	"deletescalinggroup": {
		"force": "Specifies that the group will be deleted along with all instances associated with the group, without waiting for all instances to be terminated",
		"name":  "The name of the group to delete",
	},
	"deletescalingpolicy": {
		"id": "The name or Amazon Resource Name (ARN) of the policy",
	},
	"deletesecuritygroup": {
		"id": "The ID of the security group",
	},
	"deletesnapshot": {
		"id": "The ID of the EBS snapshot",
	},
	"deletestack": {
		"name":             "The name or the unique stack ID that is associated with the stack",
		"retain-resources": "For stacks in the DELETE_FAILED state, a list of resource logical IDs that are associated with the resources you want to retain",
	},
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
	"detachalarm":         {},
	"detachcontainertask": {},
	"detachelasticip": {
		"association": "The association ID",
	},
	"detachinstance": {
		"targetgroup": "The Amazon Resource Name (ARN) of the target group",
	},
	"detachinstanceprofile": {},
	"detachinternetgateway": {
		"id":  "The ID of the Internet gateway",
		"vpc": "The ID of the VPC",
	},
	"detachmfadevice": {
		"id":   "The serial number that uniquely identifies the MFA device",
		"user": "The name of the user whose MFA device you want to deactivate",
	},
	"detachnetworkinterface": {},
	"detachpolicy":           {},
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
	"startcontainertask": {},
	"startinstance": {
		"id": "One or more instance IDs",
	},
	"stopalarm": {
		"names": "The names of the alarms",
	},
	"stopcontainertask": {},
	"stopinstance": {
		"id": "One or more instance IDs",
	},
	"updatebucket": {},
	"updatecontainertask": {
		"cluster":         "The short name or full Amazon Resource Name (ARN) of the cluster that your service is running on",
		"deployment-name": "The name of the service to update",
		"desired-count":   "The number of instantiations of the task to place and keep running in your service",
		"name":            "The family and revision (family:revision) or full Amazon Resource Name (ARN) of the task definition to run in your service",
	},
	"updatedistribution": {},
	"updateinstance": {
		"id":   "The ID of the instance",
		"lock": "If the value is true, you can't terminate the instance using the Amazon EC2 console, CLI, or API; otherwise, you can",
	},
	"updateloginprofile": {
		"password":       "The new password for the specified IAM user",
		"password-reset": "Allows this new password to be used only once by requiring the specified IAM user to set a new password on next sign-in",
		"username":       "The name of the user whose password you want to update",
	},
	"updatepolicy": {
		"arn": "The Amazon Resource Name (ARN) of the IAM policy to which you want to add a new version",
	},
	"updaterecord": {},
	"updates3object": {
		"acl":     "The canned ACL to apply to the object",
		"bucket":  "",
		"name":    "",
		"version": "VersionId used to reference a specific version of the object",
	},
	"updatescalinggroup": {
		"cooldown":                 "The amount of time, in seconds, after a scaling activity completes before another scaling activity can start",
		"desired-capacity":         "The number of EC2 instances that should be running in the Auto Scaling group",
		"healthcheck-grace-period": "The amount of time, in seconds, that Auto Scaling waits before checking the health status of an EC2 instance that has come into service",
		"healthcheck-type":         "The service to use for the health checks",
		"launchconfiguration":      "The name of the launch configuration",
		"max-size":                 "The maximum size of the Auto Scaling group",
		"min-size":                 "The minimum size of the Auto Scaling group",
		"name":                     "The name of the Auto Scaling group",
		"new-instances-protected": "Indicates whether newly launched instances are protected from termination by Auto Scaling when scaling in",
		"subnets":                 "The ID of the subnet, if you are launching into a VPC",
	},
	"updatesecuritygroup": {},
	"updatestack": {
		"capabilities":          "A list of values that you must specify before AWS CloudFormation can update certain stacks",
		"name":                  "The name or unique stack ID of the stack to update",
		"notifications":         "Amazon Simple Notification Service topic Amazon Resource Names (ARNs) that AWS CloudFormation associates with the stack",
		"parameters":            "A list of Parameter structures that specify input parameters for the stack",
		"policy-file":           "Structure containing a new stack policy body",
		"policy-update-file":    "Structure containing the temporary overriding stack policy body",
		"resource-types":        "The template resource types that you have permissions to work with for this update stack action, such as AWS::EC2::Instance, AWS::EC2::*, or Custom::MyCustomInstance",
		"role":                  "The Amazon Resource Name (ARN) of an AWS Identity and Access Management (IAM) role that AWS CloudFormation assumes to update the stack",
		"template-file":         "Structure containing the template body with a minimum length of 1 byte and a maximum length of 51,200 bytes",
		"use-previous-template": "Reuse the existing template that is associated with the stack that you are updating",
	},
	"updatesubnet": {
		"id":     "The ID of the subnet",
		"public": "Specify true to indicate that network interfaces created in the specified subnet should be assigned a public IPv4 address",
	},
	"updatetargetgroup": {},
}
