package awsdoc

func TemplateParamsDoc(templateDef, param string) (string, bool) {
	if doc, ok := manualParamsDoc[templateDef][param]; ok {
		return doc, ok
	}
	doc, ok := generatedParamsDoc[templateDef][param]
	return doc, ok
}

var manualParamsDoc = map[string]map[string]string{
	"attachalarm": {
		"name":       "The Name of the Alarm to update",
		"action-arn": "The Amazon Resource Name (ARN) of the action to execute when this alarm transitions to the ALARM state from any other state",
	},
	"attachelasticip": {
		"allow-reassociation": "Specify false to ensure the operation fails if the Elastic IP address is already associated with another resource",
	},
	"attachinstance": {
		"id":   "The ID of the Instance",
		"port": "The port on which the Instance is listenning",
	},
	"attachpolicy": {
		"arn":   "The Amazon Resource Name (ARN) of the IAM policy you want to attach",
		"user":  "The name (friendly name, not ARN) of the IAM user to attach the policy to",
		"group": "The name (friendly name, not ARN) of the IAM group to attach the policy to",
		"role":  "The name (friendly name, not ARN) of the IAM role to attach the policy to",
	},
	"attachsecuritygroup": {
		"id":       "The ID of the Security Group to add to the instance",
		"instance": "The ID of the Instance",
	},
	"checkdistribution": {
		"id":      "The ID of the CloudFront Distribution to check",
		"state":   "The state of the CloudFront Distribution to reach (Deployed | InProgress | not-found)",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"checkinstance": {
		"id":      "The ID of the EC2 Instance to check",
		"state":   "The state of the EC2 Instance to reach (pending | running | shutting-down | terminated | stopping | stopped | not-found)",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"checkloadbalancer": {
		"id":      "The ID of the ELBv2 Loadbalancer to check",
		"state":   "The state of the ELBv2 Loadbalancer to reach (provisioning | active | failed | not-found)",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"checkscalinggroup": {
		"id":      "The ID of the AutoScaling Group to check",
		"count":   "The number of Instances + Loadbalancers + TargetGroups in the AutoScaling Group to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"checksecuritygroup": {
		"id":      "The ID of the EC2 Security Group to check",
		"state":   "The state of the EC2 Security Group to reach (unused)",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"createaccesskey": {
		"user": "The name of the user for which the access key will be generated",
	},
	"createalarm": {
		"operator":           "The arithmetic operation to use when comparing the specified statistic and threshold (GreaterThanThreshold | LessThanThreshold | LessThanOrEqualToThreshold | GreaterThanOrEqualToThreshold)",
		"statistic-function": "The statistic for the metric associated with the alarm, other than percentile (Minimum | Maximum | Sum | Average | SampleCount | pNN.NN)",
	},
	"createbucket": {
		"acl":  "The canned ACL to apply to the bucket (private | public-read | public-read-write | aws-exec-read | authenticated-read | bucket-owner-read | bucket-owner-full-control | log-delivery-write)",
		"name": "The name of bucket to create",
	},
	"createdatabase": {
		"autoupgrade":       "Set to true to indicate that minor version patches are applied automatically",
		"dbname":            "The name of the database to create when the DB instance is created",
		"dbsecuritygroups":  "A list of DB security groups to associate with this DB instance",
		"domain":            "Specify the Active Directory Domain to create the instance in.",
		"engine":            "Provides the name of the database engine to be used for this DB instance (mysql | mariadb | oracle-se1 | oracle-se2 | oracle-se | oracle-ee | sqlserver-ee | sqlserver-se | sqlserver-ex | sqlserver-web | postgres | aurora)",
		"iamrole":           "Specify the name of the IAM role to be used when making API calls to the Directory Service",
		"license":           "License model information for this DB instance (license-included | bring-your-own-license | general-public-license)",
		"optiongroup":       "Indicates that the DB instance should be associated with the specified option group",
		"password":          "The password for the master database user",
		"public":            "'true' specifies an Internet-facing instance with a publicly resolvable DNS name, which resolves to a public IP address. 'false' specifies an internal instance with a DNS name that resolves to a private IP address",
		"parametergroup":    "The name of the DB parameter group to associate with this DB instance",
		"port":              "The port number on which the database accepts connections",
		"storagetype":       "Specifies the storage type associated with DB instance (standard | gp2 | io1)",
		"subnetgroup":       "A DB subnet group to associate with this DB instance",
		"vpcsecuritygroups": "A list of EC2 VPC security groups to associate with this DB instance",
	},
	"createdbsubnetgroup": {
		"description": "The description for the DB subnet group",
		"name":        "The name for the DB subnet group",
		"subnets":     "The EC2 Subnet IDs for the DB subnet group",
	},
	"createdistribution": {
		"origin-domain":   "The DNS name of the Amazon S3 bucket from which you want CloudFront to get objects for this origin, for example, myawsbucket.s3.amazonaws.com",
		"certificate":     "The Amazon Resource Name (ARN) of the AWS Certificate Manager (ACM) certificate you want to use for TSL connection",
		"comment":         "Any comments you want to include about the distribution",
		"default-file":    "The object that you want CloudFront to request from your origin (for example, index.html) when a viewer requests the root URL for your distribution (http://www.example.com)",
		"domain-aliases":  "A list of CNAMEs (alternate domain names), if any, for this distribution",
		"enable":          "From this field, you can enable or disable the selected distribution",
		"forward-cookies": "Specifies which cookies to forward to the origin for this cache behavior (all | none | whitelist)",
		"forward-queries": "Indicates whether you want CloudFront to forward query strings to the origin that is associated with this cache behavior and cache based on the query string parameters (true | false)",
		"https-behaviour": "The protocol (HTTP or HTTPS) that viewers can use to access the files (allow-all | redirect-to-https | https-only)",
		"origin-path":     "An optional element that causes CloudFront to request your content from a directory in your Amazon S3 bucket or your custom origin. When you include this element, specify the directory name, beginning with a /",
		"price-class":     "The price class that corresponds with the maximum price that you want to pay for CloudFront service. If you specify PriceClass_All, CloudFront responds to requests for your objects from all CloudFront edge locations",
		"min-ttl":         "The minimum amount of time that you want objects to stay in CloudFront caches before CloudFront forwards another request to your origin to determine whether the object has been updated",
	},
	"createelasticip": {
		"domain": "Set to vpc to allocate the address for use with instances in a VPC else the address is for use with instances in EC2-Classic (vpc | ec2-classic)",
	},
	"createfunction": {
		"bucket":        "Amazon S3 bucket name where the .zip file containing your deployment package is stored. This bucket must reside in the same AWS region where you are creating the Lambda function",
		"object":        "The Amazon S3 object (the deployment package) key name you want to upload",
		"objectversion": "The Amazon S3 object (the deployment package) version you want to upload",
		"zipfile":       "The path toward the zip file containing your deployment package",
	},
	"creategroup": {
		"name": "The name of the group to create",
	},
	"createinstance": {
		"count": "The number of instances to launch",
		"name":  "The name of the instance to launch",
		"role":  "The name of the instance profile (role) to launch the instance with",
		"image": "The ID of the AMI of the instance to launch, which you can get by using `awless search images`",
	},
	"createkeypair": {
		"name":      "The name of the keypair to create (it will also be the name of the file stored in ~/.awless/keys)",
		"encrypted": "Set to 'true' if you want to encrypt the keypair"},
	"createlaunchconfiguration": {
		"public": "Used for groups that launch instances into a virtual private cloud (VPC). Specifies whether to assign a public IP address to each instance",
	},
	"createlistener": {
		"actiontype":  "The type of action (forward)",
		"targetgroup": "The Amazon Resource Name (ARN) of the target group",
		"certificate": "The Amazon Resource Name (ARN) of the certificate",
		"protocol":    "The protocol for connections from clients to the load balancer (TCP | HTTP | HTTPS)",
	},
	"createloadbalancer": {
		"scheme": "The routing range of the loadbalancer (Internet-facing | internal)",
	},
	"createpolicy": {
		"name":        "The friendly name of the policy",
		"description": "A friendly description of the policy",
		"effect":      "The Effect element is required and specifies whether the policy will result in an allow or an explicit deny (Allow | Deny)",
		"action":      "The Action element describes the specific action or actions that will be allowed or denied. You specify a value using a namespace that identifies a service followed by the name of the action to allow or deny (eg. sqs:SendMessage, s3:*)",
		"resource":    "The Amazon Resource Name (ARN) of the Resource element which specifies the object or objects that the policy covers",
	},
	"createqueue": {
		"delay":              "The length of time, in seconds, for which the delivery of all messages in the queue is delayed. Valid values: An integer from 0 to 900 seconds (15 minutes). The default is 0",
		"max-msg-size":       "The limit of how many bytes a message can contain before Amazon SQS rejects it. Valid values: An integer from 1,024 bytes (1 KiB) to 262,144 bytes (256 KiB). The default is 262,144 (256 KiB)",
		"retention-period":   "The length of time, in seconds, for which Amazon SQS retains a message. Valid values: An integer from 60 seconds (1 minute) to 1,209,600 seconds (14 days). The default is 345,600 (4 days)",
		"policy":             "The queue's policy",
		"msg-wait":           "The length of time, in seconds, for which a ReceiveMessage action waits for a message to arrive. Valid values: An integer from 0 to 20 (seconds). The default is 0",
		"redrive-policy":     "The parameters for the dead letter queue functionality of the source queue",
		"visibility-timeout": "The visibility timeout for the queue. Valid values: An integer from 0 to 43,200 (12 hours). The default is 30",
	},
	"createrecord": {
		"zone":    "The ID of the hosted zone that contains the resource record sets that you want to change",
		"name":    "The name of the domain you want to perform the action on. Enter a fully qualified domain name, for example, www.example.com. You can optionally include a trailing dot",
		"type":    "The DNS record type. (A | AAAA | CNAME | MX | NAPTR | NS | PTR | SOA | SPF | SRV | TXT)",
		"value":   "The current or new DNS record value",
		"ttl":     "The resource record cache time to live (TTL), in seconds",
		"comment": "Any comments you want to include about a change batch request",
	},
	"createrole": {
		"name":              "The name of the role to create",
		"principal-account": "The account entity that can perform actions and access resources",
		"principal-user":    "The user entity that can perform actions and access resources",
		"principal-service": "The service entity that can perform actions and access resources",
		"sleep-after":       "The amount of time in seconds you want to wait after creating the role (usually used to be sure that the role creation has been propagated)",
	},
	"creates3object": {
		"bucket": "Name of the bucket to which object will be added",
		"file":   "The path toward to file to upload",
		"name":   "The name of the Object to create (by default the file name is used)",
		"acl":    "The canned ACL to apply to the object (private | public-read | public-read-write | aws-exec-read | authenticated-read | bucket-owner-read | bucket-owner-full-control | log-delivery-write)",
	},
	"createscalingpolicy": {
		"adjustment-type": "The adjustment type (ChangeInCapacity | ExactCapacity | PercentChangeInCapacity)",
	},
	"createstack": {
		"capabilities":  "A list of values that you must specify before AWS CloudFormation can create certain stacks (CAPABILITY_IAM | CAPABILITY_NAMED_IAM)",
		"on-failure":    "Determines what action will be taken if stack creation fails (DO_NOTHING | ROLLBACK | DELETE)",
		"parameters":    "A list of Parameters that specify input parameters for the stack given using this format: key1:val1,key2:val2,...",
		"policy-file":   "The path to the file containing the stack policy body",
		"template-file": "The path to the file containing the template body with a minimum size of 1 byte and a maximum size of 51,200 bytes",
	},
	"createsubnet": {
		"name": "The 'Name' Tag for the subnet to create",
	},
	"createsubscription": {
		"endpoint": "The endpoint that you want to receive notifications. Endpoints vary by protocol: For the http or https protocol, the endpoint is a URL beginning with 'http://' or 'https://', for the email or email-json protocol, the endpoint is an email address, for the sms protocol, the endpoint is a phone number of an SMS-enabled, for the sqs protocol, the endpoint is the ARN of an Amazon SQS queue, for the application protocol, the endpoint is the EndpointArn of a mobile app and device, for the lambda protocol, the endpoint is the ARN of an AWS Lambda function.",
		"protocol": "The protocol you want to use (http | https | email | email-json | sms | sqs | lambda)",
		"topic":    "The ARN of the topic you want to subscribe to",
	},
	"createtag": {
		"resource": "The ID of the resource on which you want to add a tag",
		"key":      "The Tag key",
		"value":    "The Tag value",
	},
	"createtargetgroup": {
		"matcher": "The HTTP codes to use when checking for a successful response from a target",
	},
	"createvpc": {
		"name": "The 'Name' Tag for the VPC to create",
	},
	"createzone": {
		"comment":   "Any comments that you want to include about the hosted zone",
		"isprivate": "A value that indicates whether this is a private hosted zone",
		"vpcid":     "(Private hosted zones only) The ID of an Amazon VPC",
		"vpcregion": "(Private hosted zones only) The region in which you created an Amazon VPC",
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
	"detachalarm": {},
	"detachelasticip": {
		"association": "The association ID",
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
}
