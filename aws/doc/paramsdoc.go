package awsdoc

import "strings"

func TemplateParamsDoc(action, entity, param string) (string, bool) {
	if doc, ok := manualParamsDoc[action+"."+entity][param]; ok {
		return doc, ok
	}
	doc, ok := generatedParamsDoc[action+"."+entity][param]
	return doc, ok
}

func TemplateParamsDocWithEnums(action, entity, param string) (string, bool) {
	var suffix string
	if enum, ok := EnumDoc[action+"."+entity+"."+param]; ok && len(enum) > 0 && strings.TrimSpace(enum[0]) != "" {
		suffix = " (" + strings.Join(enum, " | ") + ")"
	}
	doc, ok := TemplateParamsDoc(action, entity, param)
	return doc + suffix, ok
}

var manualParamsDoc = map[string]map[string]string{
	"attach.alarm": {
		"name":       "The Name of the Alarm to update",
		"action-arn": "The Amazon Resource Name (ARN) of the action to execute when this alarm transitions to the ALARM state from any other state",
	},
	"attach.classicloadbalancer": {
		"name":     "The name of the Classic load balancer",
		"instance": "The ID of the instance",
	},
	"attach.containertask": {
		"container-name":    "The name of a container",
		"name":              "The name of the new or existing task containing the container to attach",
		"image":             "The image used to start a container. Images in the Docker Hub registry are available by default. Other repositories are specified with repository-url/image:tag",
		"memory-hard-limit": "The hard limit (in MiB) of memory to present to the container. If your container attempts to exceed the memory specified here, the container is killed",
		"command":           "The command that is passed to the container",
		"env":               "The environment variables to pass to a container using this format: [key1:val1,key2:val2,...]",
		"privileged":        "When this parameter is true, the container is given elevated privileges on the host container instance",
		"workdir":           "The working directory in which to run commands inside the container",
		"ports":             "The list of port mappings for the container. Port mappings allow containers to access ports on the host container instance to send or receive traffic (format [host-port:]container-port[/protocol][,[host-port:]container-port[/protocol]])",
	},
	"attach.elasticip": {
		"allow-reassociation": "Specify false to ensure the operation fails if the Elastic IP address is already associated with another resource",
	},
	"attach.instance": {
		"id":   "The ID of the Instance",
		"port": "The port on which the Instance is listenning",
	},
	"attach.instanceprofile": {
		"instance": "The ID of the Instance",
		"name":     "The name of the InstanceProfile to associate to the Instance",
		"replace":  "If 'true' will replace existing instance profile with provided one",
	},
	"attach.listener": {
		"certificate": "The awless alias of the certificate's name (ex: @www.mysite.com), or the full certificate's ARN",
	},
	"attach.mfadevice": {
		"no-prompt": "Use 'true' to disable the prompt that asks to append the mfadevice to ~/.aws/config file",
	},
	"attach.policy": {
		"access":  "Type of access to retrieve an AWS policy",
		"service": "Service string to retrieve an AWS policy",
		"arn":     "The Amazon Resource Name (ARN) of the IAM policy you want to attach",
		"user":    "The name (friendly name, not ARN) of the IAM user to attach the policy to",
		"group":   "The name (friendly name, not ARN) of the IAM group to attach the policy to",
		"role":    "The name (friendly name, not ARN) of the IAM role to attach the policy to",
	},
	"attach.securitygroup": {
		"id":       "The ID of the Security Group to add to the instance",
		"instance": "The ID of the Instance",
	},
	"authenticate.registry": {
		"accounts":        "A list of AWS account IDs that are associated with the registries for which to authenticate",
		"no-confirm":      "Do not ask confirmation before effectively running `docker login` command",
		"no-docker-login": "Set to 'true' to disable the prompt and automatic execution of `docker login` command",
	},
	"check.certificate": {
		"arn":     "The Amazon Resource Name (ARN) of the certificate to check",
		"state":   "The state of the certificate to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.database": {
		"id":      "The ID of the RDS Database to check",
		"state":   "The state of the RDS Database to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.distribution": {
		"id":      "The ID of the CloudFront Distribution to check",
		"state":   "The state of the CloudFront Distribution to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.instance": {
		"id":      "The ID of the EC2 Instance to check",
		"state":   "The state of the EC2 Instance to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.loadbalancer": {
		"id":      "The ID of the ELBv2 Loadbalancer to check",
		"state":   "The state of the ELBv2 Loadbalancer to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.natgateway": {
		"id":      "The ID of the NAT Gateway to check",
		"state":   "The state of the NAT Gateway to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.networkinterface": {
		"id":      "The ID of the Network Interface to check",
		"state":   "The state of the Network Interface to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.scalinggroup": {
		"name":    "The name of the AutoScaling Group to check",
		"count":   "The number of Instances + Loadbalancers + TargetGroups in the AutoScaling Group to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.securitygroup": {
		"id":      "The ID of the EC2 Security Group to check",
		"state":   "The state of the EC2 Security Group to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"check.volume": {
		"id":      "The ID of the EC2 Volume to check",
		"state":   "The state of the EC2 Volume to reach",
		"timeout": "The time (in seconds) after which the check is failed",
	},
	"create.accesskey": {
		"user":      "The name of the user for which the access key will be generated",
		"save":      "Use 'true' to save the access key in ~/.aws/credentials under 'user' profile; use 'false' to disable the prompt",
		"no-prompt": "Deprecated - use the save param",
	},
	"create.alarm": {
		"operator":           "The arithmetic operation to use when comparing the specified statistic and threshold",
		"statistic-function": "The statistic for the metric associated with the alarm, other than percentile",
		"unit":               "The unit of measure for the statistic",
	},
	"create.appscalingtarget": {
		"dimension":         "The scalable dimension associated with the scalable target",
		"resource":          "The identifier of the resource associated with the scalable target (eg. for ECS: service/cluster-name/service-deployment-name, for EC2 spot-fleet: spot-fleet-request/sfr-73fbd2ce-aa30-494c-8788-1cee4EXAMPLE, for EMR cluster: instancegroup/j-2EEZNYKUA1NTV/ig-1791Y4E1L8YI0, for AppStream 2.0 fleet: fleet/sample-fleet, for DynamoDB table: table/my-table, for DynamoDB global secondary index: table/my-table/index/my-table-index)",
		"service-namespace": "The namespace of the AWS service",
	},
	"create.appscalingpolicy": {
		"dimension":         "The scalable dimension associated with the scalable target",
		"resource":          "The identifier of the resource associated with the scalable target (eg. for ECS: service/cluster-name/service-deployment-name, for EC2 spot-fleet: spot-fleet-request/sfr-73fbd2ce-aa30-494c-8788-1cee4EXAMPLE, for EMR cluster: instancegroup/j-2EEZNYKUA1NTV/ig-1791Y4E1L8YI0, for AppStream 2.0 fleet: fleet/sample-fleet, for DynamoDB table: table/my-table, for DynamoDB global secondary index: table/my-table/index/my-table-index)",
		"service-namespace": "The namespace of the AWS service",
		"type":              "The policy type",
		"stepscaling-adjustment-type":          "The scalable dimension",
		"stepscaling-adjustments":              "A set of adjustments that enable you to scale based on the size of the alarm breach using this format: [[from]:[to]:scaling-adjustment[,[from]:[to]:scaling-adjustment[,...]]]",
		"stepscaling-cooldown":                 "The amount of time, in seconds, after a scaling activity completes where previous trigger-related scaling activities can influence future scaling events",
		"stepscaling-aggregation-type":         "The aggregation type for the CloudWatch metrics",
		"stepscaling-min-adjustment-magnitude": "The minimum number to adjust your scalable dimension as a result of a scaling activity",
	},
	"create.bucket": {
		"acl":  "The canned ACL to apply to the bucket",
		"name": "The name of bucket to create",
	},
	"create.certificate": {
		"domains":            "Main and Additional Fully qualified domain names (FQDNs) to be included in the Certificate name and Subject Alternative Name of the ACM Certificate",
		"validation-domains": "The domain name that you want ACM to use to send you validation emails. This domain name is the suffix of the email addresses that you want ACM to use. This must be the same as the DomainName value or a superdomain of the domain value",
	},
	"create.classicloadbalancer": {
		"name":             "The name of the Classic load balancer",
		"zones":            "At least 1 availability zone from the same region as the load balancer",
		"listeners":        "The list of listeners. Format for a listener: LOADB_PROTO:LOADB_PORT:INST_PROTO:INST_PORT . Ex: [HTTPS:443:HTTP:8080,TCP:53:TCP:4567]",
		"healthcheck-path": "The healthcheck ping path only. Otherwise default to ping port 80 using TCP. Ex: home.html (protocol and port are derived from the first listener)",
		"tags":             "A list of tags to assign to the load balancer. Example: tags=Env:Prod,Country:US",
	},
	"create.database": {
		"autoupgrade":        "Set to true to indicate that minor version patches are applied automatically",
		"availabilityzone":   "Specifies the name of the Availability Zone the DB instance is located in",
		"backupretention":    "Specifies the number of days for which automatic DB snapshots are retained",
		"backupwindow":       "Specifies the daily time range during which automated backups are created if automated backups are enabled, as determined by the BackupRetentionPeriod (format hh24:mi-hh24:mi)",
		"cluster":            "If the DB instance is a member of a DB cluster, contains the name of the DB cluster that the DB instance is a member of",
		"copytagstosnapshot": "True to copy all tags from the READ replica to snapshots of the READ Replica. Default is false",
		"dbname":             "The name of the database to create when the DB instance is created",
		"dbsecuritygroups":   "A list of DB security groups to associate with this DB instance",
		"domain":             "Specify the Active Directory Domain to create the instance in",
		"encrypted":          "Specifies whether the DB instance is encrypted",
		"engine":             "The name of the database engine to be used for this DB instance (not every engine is available for every region)",
		"iamrole":            "Specify the name of the IAM role to be used when making API calls to the Directory Service",
		"id":                 "Contains a user-supplied database identifier",
		"iops":               "Specifies the Provisioned IOPS (I/O operations per second) value",
		"license":            "License model information for this DB instance",
		"maintenancewindow":  "Specifies the weekly time range during which system maintenance can occur, in Universal Coordinated Time (UTC)",
		"multiaz":            "Specifies if the DB instance is a Multi-AZ deployment",
		"optiongroup":        "Indicates that the DB instance should be associated with the specified option group",
		"password":           "The password for the master database user",
		"public":             "'true' specifies an Internet-facing instance with a publicly resolvable DNS name, which resolves to a public IP address. 'false' specifies an internal instance with a DNS name that resolves to a private IP address",
		"parametergroup":     "The name of the DB parameter group to associate with this DB instance",
		"port":               "The port number on which the database accepts connections",
		"replica":            "The DB instance identifier of the READ replica",
		"replica-source":     "The identifier of the DB instance that will act as the source for the READ replica (each DB instance can have up to 5 Read replicas). Use the Amazon Resource Name (ARN) of the database if it is not in the same region",
		"size":               "Specifies the allocated storage size specified in gigabytes",
		"storagetype":        "Specifies the storage type associated with DB instance",
		"subnetgroup":        "A DB subnet group to associate with this DB instance",
		"timezone":           "The time zone of the DB instance",
		"type":               "Contains the name of the compute and memory capacity class of the DB instance",
		"username":           "Contains the master username for the DB instance",
		"version":            "Indicates the database engine version",
		"vpcsecuritygroups":  "A list of EC2 VPC security groups to associate with this DB instance",
	},
	"create.dbsubnetgroup": {
		"description": "The description for the DB subnet group",
		"name":        "The name for the DB subnet group",
		"subnets":     "The EC2 Subnet IDs for the DB subnet group",
	},
	"create.distribution": {
		"origin-domain":   "The DNS name of the Amazon S3 bucket from which you want CloudFront to get objects for this origin, for example, myawsbucket.s3.amazonaws.com",
		"certificate":     "The Amazon Resource Name (ARN) of the AWS Certificate Manager (ACM) certificate you want to use for TSL connection",
		"comment":         "Any comments you want to include about the distribution",
		"default-file":    "The object that you want CloudFront to request from your origin when a viewer requests the root URL for your distribution (http://www.example.com)",
		"domain-aliases":  "A list of CNAMEs (alternate domain names), if any, for this distribution",
		"enable":          "From this field, you can enable or disable the selected distribution",
		"forward-cookies": "Specifies which cookies to forward to the origin for this cache behavior",
		"forward-queries": "Indicates whether you want CloudFront to forward query strings to the origin that is associated with this cache behavior and cache based on the query string parameters",
		"https-behaviour": "The protocol (HTTP or HTTPS) that viewers can use to access the files",
		"origin-path":     "An optional element that causes CloudFront to request your content from a directory in your Amazon S3 bucket or your custom origin. When you include this element, specify the directory name, beginning with a /",
		"price-class":     "The price class that corresponds with the maximum price that you want to pay for CloudFront service. If you specify PriceClass_All, CloudFront responds to requests for your objects from all CloudFront edge locations",
		"min-ttl":         "The minimum amount of time that you want objects to stay in CloudFront caches before CloudFront forwards another request to your origin to determine whether the object has been updated",
	},
	"create.elasticip": {
		"domain": "Set to vpc to allocate the address for use with instances in a VPC else the address is for use with instances in EC2-Classic",
	},
	"create.function": {
		"bucket":        "Amazon S3 bucket name where the .zip file containing your deployment package is stored. This bucket must reside in the same AWS region where you are creating the Lambda function",
		"object":        "The Amazon S3 object (the deployment package) key name you want to upload",
		"objectversion": "The Amazon S3 object (the deployment package) version you want to upload",
		"runtime":       "The runtime environment for the Lambda function you are uploading",
		"zipfile":       "The path toward the zip file containing your deployment package",
	},
	"create.group": {
		"name": "The name of the group to create",
	},
	"create.instance": {
		"count":  "The number of instances to launch",
		"name":   "The name of the instance to launch",
		"role":   "The name of the instance profile (role) to launch the instance with",
		"image":  "The ID of an AMI for the instance to be launched",
		"distro": "The distro query to resolve official community free bare distro AMI from current region. See above description from this help for specific queries. Default choices:",
	},
	"create.image": {
		"reboot": "True to shut down and reboot the instance before creating the image, otherwise no reboot and file system integrity on the created image cannot be guaranteed",
	},
	"create.keypair": {
		"name":      "The name of the keypair to create (it will also be the name of the file stored in ~/.awless/keys)",
		"encrypted": "Set to 'true' if you want to encrypt the keypair"},
	"create.launchconfiguration": {
		"distro": "The distro query to resolve official community bare distro AMI from current region. See `awless search images -h`",
		"public": "Used for groups that launch instances into a virtual private cloud (VPC). Specifies whether to assign a public IP address to each instance",
	},
	"create.listener": {
		"actiontype":  "The type of action",
		"targetgroup": "The Amazon Resource Name (ARN) of the target group",
		"certificate": "The Amazon Resource Name (ARN) of the certificate",
		"protocol":    "The protocol for connections from clients to the load balancer",
		"sslpolicy":   "The security policy that defines which ciphers and protocols are supported",
	},
	"create.mfadevice": {
		"name": "The name of the virtual MFA device",
	},
	"create.policy": {
		"name":        "The friendly name of the policy",
		"description": "A friendly description of the policy",
		"effect":      "The Effect element is required and specifies whether the policy will result in an allow or an explicit deny",
		"action":      "The Action elements describing the actions that will be allowed or denied. You specify a value using a namespace that identifies a service followed by the name of the action to allow or deny (eg. sqs:SendMessage, s3:*). Use a list for multiple actions",
		"resource":    "The Amazon Resource Name (ARN) of the Resource element which specifies the object or objects that the policy covers",
		"conditions":  "List of conditions necessary for the policy to be in effect (e.g. [aws:UserAgent!=My user agent,s3:prefix=~home/,aws:CurrentTime>=2013-06-30T00:00:00Z,aws:SourceIp!=203.0.113.0/24,aws:SourceArn==arn:aws:sns:eu-west-1:*:*])",
	},
	"create.queue": {
		"delay":              "The length of time, in seconds, for which the delivery of all messages in the queue is delayed. Valid values: An integer from 0 to 900 seconds (15 minutes). The default is 0",
		"max-msg-size":       "The limit of how many bytes a message can contain before Amazon SQS rejects it. Valid values: An integer from 1024 bytes (1 KiB) to 262144 bytes (256 KiB). The default is 262144 (256 KiB)",
		"retention-period":   "The length of time, in seconds, for which Amazon SQS retains a message. Valid values: An integer from 60 seconds (1 minute) to 1209600 seconds (14 days). The default is 345600 (4 days)",
		"policy":             "The queue's policy",
		"msg-wait":           "The length of time, in seconds, for which a ReceiveMessage action waits for a message to arrive. Valid values: An integer from 0 to 20 (seconds). The default is 0",
		"redrive-policy":     "The parameters for the dead letter queue functionality of the source queue",
		"visibility-timeout": "The visibility timeout for the queue. Valid values: An integer from 0 to 43200 (12 hours). The default is 30",
	},
	"create.record": {
		"zone":    "The ID of the hosted zone that contains the resource record sets that you want to change",
		"name":    "The name of the domain you want to perform the action on. Enter a fully qualified domain name, for example, www.example.com. You can optionally include a trailing dot",
		"type":    "The DNS record type",
		"value":   "The new DNS record value",
		"values":  "The new DNS record value(s)",
		"ttl":     "The resource record cache time to live (TTL), in seconds",
		"comment": "Any comments you want to include about a change batch request",
	},
	"create.role": {
		"conditions":        "List of conditions necessary for the policy to be in effect (e.g. [aws:UserAgent!=My user agent,s3:prefix=~home/,aws:CurrentTime>=2013-06-30T00:00:00Z,aws:SourceIp!=203.0.113.0/24,aws:SourceArn==arn:aws:sns:eu-west-1:*:*])",
		"name":              "The name of the role to create",
		"principal-account": "The ID of the account that can perform actions and access resources of the role (you can know your account ID with `awless whoami`)",
		"principal-user":    "The Amazon Resource Name (ARN) of the user that can perform actions and access resources of the role",
		"principal-service": "The AWS Service that can assume this role to perform actions and access resources of the role (e.g. 'ec2.amazonaws.com')",
		"sleep-after":       "The amount of time in seconds you want to wait after creating the role (usually used to be sure that the role creation has been propagated)",
	},
	"create.s3object": {
		"bucket": "Name of the bucket to which object will be added",
		"file":   "The path toward to file to upload",
		"name":   "The name of the Object to create (by default the file name is used)",
		"acl":    "The canned ACL to apply to the object",
	},
	"create.scalinggroup": {
		"healthcheck-type": "The service to use for the health checks",
	},
	"create.scalingpolicy": {
		"adjustment-type":    "The adjustment type",
		"adjustment-scaling": "The amount by which to scale, based on the specified adjustment type (e.g. '-2', '3')",
	},
	"create.stack": {
		"capabilities":            "A list of values that you must specify before AWS CloudFormation can create certain stacks",
		"on-failure":              "Determines what action will be taken if stack creation fails",
		"parameters":              "A list of Parameters that specify input parameters for the stack given using this format: [key1:val1,key2:val2,...]",
		"policy-file":             "The path to the file containing the stack policy body",
		"template-file":           "The path to the file containing the template body with a minimum size of 1 byte and a maximum size of 51,200 bytes",
		"stack-file":              "The path to the file containing Parameters/Tags/StackPolices definition (http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/continuous-delivery-codepipeline-cfn-artifacts.html#w2ab2c13c15c15). Values passed via CLI has higher priority than ones defined in StackFile",
		"rollback-triggers":       "List of CloudWatch Alarm ARNs to monitor during and after creation",
		"rollback-monitoring-min": "Time to monitor rollback-triggers during and after creation",
	},
	"create.subnet": {
		"name":   "The 'Name' Tag for the subnet to create",
		"public": "A value (true) to indicate that network interfaces created in this subnet should be assigned a public IPv4 address (instances, etc.)",
	},
	"create.subscription": {
		"endpoint": "The endpoint that you want to receive notifications. Endpoints vary by protocol: For the http or https protocol, the endpoint is a URL beginning with 'http://' or 'https://', for the email or email-json protocol, the endpoint is an email address, for the sms protocol, the endpoint is a phone number of an SMS-enabled, for the sqs protocol, the endpoint is the ARN of an Amazon SQS queue, for the application protocol, the endpoint is the EndpointArn of a mobile app and device, for the lambda protocol, the endpoint is the ARN of an AWS Lambda function",
		"protocol": "The protocol you want to use",
		"topic":    "The ARN of the topic you want to subscribe to",
	},
	"create.tag": {
		"resource": "The ID of the resource on which you want to add a tag",
		"key":      "The Tag key",
		"value":    "The Tag value",
	},
	"create.targetgroup": {
		"matcher": "The HTTP codes to use when checking for a successful response from a target",
	},
	"create.vpc": {
		"name": "The 'Name' Tag for the VPC to create",
	},
	"create.zone": {
		"comment":   "Any comments that you want to include about the hosted zone",
		"isprivate": "A value that indicates whether this is a private hosted zone",
		"vpcid":     "(Private hosted zones only) The ID of an Amazon VPC",
		"vpcregion": "(Private hosted zones only) The region in which you created an Amazon VPC",
	},
	"delete.accesskey": {
		"id": "The ID of the access key and secret access key you want to delete",
	},
	"delete.alarm": {
		"name": "The name of the alarm(s) to be deleted",
	},
	"delete.bucket": {
		"name": "The name of the bucket to be deleted",
	},
	"delete.containertask": {
		"name":         "The name of the containertask to be deleted",
		"all-versions": "Set to 'true' to delete all existing versions of the containertask to be deleted",
	},
	"delete.database": {
		"id":            "The ID of the database to be deleted",
		"skip-snapshot": "Determines whether a final DB snapshot is created before the DB instance is deleted. If true is specified, no DBSnapshot is created. If false is specified, a DB snapshot is created before the DB instance is deleted",
		"snapshot":      "The ID of the new DBSnapshot created when skip-snapshot=false",
	},
	"delete.dbsubnetgroup": {
		"name": "The name of the database subnet group to be deleted",
	},
	"delete.distribution": {
		"id": "The ID of the distribution to be deleted",
	},
	"delete.function": {
		"id": "The ID of the Lambda function to be deleted",
	},
	"delete.image": {
		"id":               "The ID of the AMI to be deleted",
		"delete-snapshots": "Set to 'true' to also delete the snapshots created from this image",
	},
	"delete.instance": {
		"ids": "The ID(s) of the instance(s) to be deleted",
		"id":  "The ID of the instance(s) to be deleted",
	},
	"delete.internetgateway": {
		"id": "The ID of the Internet gateway to be deleted",
	},
	"delete.keypair": {
		"name": "The name of the key pair to be deleted",
	},
	"delete.launchconfiguration": {
		"name": "The name of the launch configuration to be deleted",
	},
	"delete.classicloadbalancer": {
		"name": "The name of the Classic load balancer",
	},
	"delete.policy": {
		"all-versions": "Set to 'true' to delete all existing versions of the policy to be deleted",
	},
	"delete.record": {
		"id":     "The awless id (cf `awless list records`) of the record to delete",
		"zone":   "The ID of the hosted zone that contains the resource record sets that you want to delete",
		"name":   "The name of the domain you want to perform the action on. Enter a fully qualified domain name, for example, www.example.com. You can optionally include a trailing dot",
		"type":   "The DNS record type",
		"value":  "The DNS record value to delete",
		"values": "The DNS record value(s) to delete",
		"ttl":    "The resource record cache time to live (TTL), in seconds",
	},
	"delete.role": {
		"name": "The name of the role to be deleted",
	},
	"delete.s3object": {
		"bucket": "The name of the bucket containing the object to be deleted",
		"name":   "The name (i.e. key) of the object to be deleted",
	},
	"delete.tag": {
		"resource": "The ID of the resource on which you want to remove a tag",
		"key":      "The Tag key",
		"value":    "The Tag value",
	},
	"detach.alarm": {
		"name":       "The name of the alarm",
		"action-arn": "The Amazon Resource Name (ARN) to be detached of the ALARM actions",
	},
	"detach.classicloadbalancer": {
		"name":     "The name of the Classic load balancer",
		"instance": "The ID of the instance",
	},
	"detach.containertask": {
		"container-name": "The name of the container to detach",
		"name":           "The name of the existing container task containing the container to detach",
	},
	"detach.instance": {
		"id": "The ID of the instance to be detached from target group",
	},
	"detach.instanceprofile": {
		"instance": "The ID of the Instance",
		"name":     "The name of the InstanceProfile to detach from the Instance",
	},
	"detach.networkinterface": {
		"attachment": "The ID of the attachment",
		"force":      "Specifies whether to force a detachment",
		"id":         "The ID of the network interface",
		"instance":   "The ID of the instance this network interface is attached to",
	},

	"detach.policy": {
		"access":  "Type of access to retrieve an AWS policy",
		"service": "Service string to retrieve an AWS policy",
		"arn":     "The Amazon Resource Name (ARN) of the IAM policy you want to detach",
		"user":    "The name (friendly name, not ARN) of the IAM user to detach the policy to",
		"group":   "The name (friendly name, not ARN) of the IAM group to detach the policy to",
		"role":    "The name (friendly name, not ARN) of the IAM role to detach the policy to",
	},
	"detach.securitygroup": {
		"id":       "The ID of the security group",
		"instance": "The ID of the instance to be detached",
	},
	"import.image": {
		"architecture": "The architecture of the virtual machine",
		"url":          "The URL to the Amazon S3-based disk image being imported. The URL can either be a https URL (https://..) or an Amazon S3 URL (s3://..)",
		"snapshot":     "The ID of the EBS snapshot to be used for importing the snapshot",
		"bucket":       "The name of the S3 bucket where the disk image is located",
		"s3object":     "The name of the S3 object where the disk image is located",
		"license":      "The license type to be used for the Amazon Machine Image (AMI) after importing",
		"platform":     "The operating system of the virtual machine",
	},
	"restart.instance": {
		"id": "The ID of the instance to be restarted",
	},
	"restart.database": {
		"with-failover": "When true, the reboot is conducted through a MultiAZ failover",
	},
	"start.containertask": {
		"cluster":                     "The short name or full Amazon Resource Name (ARN) of the cluster on which to run your task",
		"type":                        "The type of task to launch",
		"desired-count":               "The number of instantiations of the specified service to place and keep running on your cluster",
		"loadbalancer.container-name": "The name of the container (as it appears in a container definition) to associate with the load balancer",
		"loadbalancer.container-port": "The port on the container to associate with the load balancer",
		"loadbalancer.targetgroup":    "The full Amazon Resource Name (ARN) of the Elastic Load Balancing target group associated with a service",
		"name":            "The name of the container task to start",
		"deployment-name": "The deployment name of the service (e.g. prod, staging...)",
		"role":            "The name or full Amazon Resource Name (ARN) of the IAM role that allows Amazon ECS to make calls to your load balancer on your behalf",
	},
	"start.instance": {
		"id": "The ID of the instance to be started",
	},
	"stop.containertask": {
		"cluster":         "The short name or full Amazon Resource Name (ARN) of the cluster on which to run your task",
		"type":            "The type of task to launch",
		"deployment-name": "The deployment name of the service (e.g. prod, staging...)",
		"run-arn":         "The ID or full Amazon Resource Name (ARN) entry of the run of the task to stop",
	},
	"stop.instance": {
		"id": "The ID of the instance to be stopped",
	},
	"update.bucket": {
		"name":              "The name of the bucket to update",
		"acl":               "The canned ACL to apply to the bucket",
		"public-website":    "Set to 'true' if you want to publish the content of the bucket as a public HTTP website",
		"redirect-hostname": "Hostname where HTTP requests will be redirected when publishing website",
		"index-suffix":      "A suffix that is appended to a request that is for a directory on the website endpoint",
		"enforce-https":     "Use HTTPS rather than HTTP when redirecting requests",
	},
	"update.classicloadbalancer": {
		"health-interval":     "The approximate interval, in seconds, between health checks of an individual instance",
		"health-timeout":      "The amount of time, in seconds, during which no response means a failed health check. This value must be less than the Interval value",
		"healthy-threshold":   "The number of consecutive health checks successes required before moving the instance to the Healthy state",
		"unhealthy-threshold": "The number of consecutive health check failures required before moving the instance to the Unhealthy state",
		"health-target":       "String with format PROTOCOL:PORT[PING_PATH]. For HTTP/HTTPS, you must include a ping path. Protocols: TCP, HTTP, HTTPS, or SSL. Ex: TCP:5000, or HTTP:80/weather/us/wa/seattle",
	},
	"update.distribution": {
		"id":              "The ID of the distribution to update",
		"origin-domain":   "The DNS name of the Amazon S3 bucket from which you want CloudFront to get objects for this origin, for example, myawsbucket.s3.amazonaws.com",
		"certificate":     "The Amazon Resource Name (ARN) of the AWS Certificate Manager (ACM) certificate you want to use for TSL connection",
		"comment":         "Any comments you want to include about the distribution",
		"default-file":    "The object that you want CloudFront to request from your origin when a viewer requests the root URL for your distribution (http://www.example.com)",
		"domain-aliases":  "A list of CNAMEs (alternate domain names), if any, for this distribution",
		"forward-cookies": "Specifies which cookies to forward to the origin for this cache behavior",
		"forward-queries": "Indicates whether you want CloudFront to forward query strings to the origin that is associated with this cache behavior and cache based on the query string parameters (true | false)",
		"https-behaviour": "The protocol (HTTP or HTTPS) that viewers can use to access the files",
		"origin-path":     "An optional element that causes CloudFront to request your content from a directory in your Amazon S3 bucket or your custom origin. When you include this element, specify the directory name, beginning with a /",
		"price-class":     "The price class that corresponds with the maximum price that you want to pay for CloudFront service. If you specify PriceClass_All, CloudFront responds to requests for your objects from all CloudFront edge locations",
		"min-ttl":         "The minimum amount of time that you want objects to stay in CloudFront caches before CloudFront forwards another request to your origin to determine whether the object has been updated",
		"enable":          "Enable/Disable the distribution",
	},
	"update.image": {
		"accounts":      "List (one or more) AWS account IDs",
		"description":   "A new description for the AMI",
		"id":            "The ID of the AMI",
		"groups":        "List (one or more) user groups",
		"operation":     "The operation type for launch permissions",
		"product-codes": "One or more DevPay product codes. After adding a product code, it cannot be removed",
	},
	"update.instance": {
		"type": "Changes the instance type to the specified value",
	},
	"update.policy": {
		"arn":        "The Amazon Resource Name (ARN) of the IAM policy you want to attach",
		"effect":     "The Effect element is required and specifies whether the policy will result in an allow or an explicit deny",
		"action":     "The Action elements describing the actions that will be allowed or denied. You specify a value using a namespace that identifies a service followed by the name of the action to allow or deny (eg. sqs:SendMessage, s3:*). Use a list for multiple actions",
		"resource":   "The Amazon Resource Name (ARN) of the Resource element which specifies the object or objects that the policy covers",
		"conditions": "List of conditions necessary for the policy to be in effect (e.g. [aws:UserAgent!=My user agent,s3:prefix=~home/,aws:CurrentTime>=2013-06-30T00:00:00Z,aws:SourceIp!=203.0.113.0/24,aws:SourceArn==arn:aws:sns:eu-west-1:*:*])",
	},
	"update.record": {
		"zone":    "The ID of the hosted zone that contains the resource record sets that you want to change",
		"name":    "The name of the domain you want to perform the action on. Enter a fully qualified domain name, for example, www.example.com. You can optionally include a trailing dot",
		"type":    "The DNS record type",
		"value":   "The current or new DNS record value",
		"values":  "The current or new DNS record value(s)",
		"ttl":     "The resource record cache time to live (TTL), in seconds",
		"comment": "Any comments you want to include about a change batch request",
	},
	"update.s3object": {
		"acl":     "The canned ACL to apply to the bucket",
		"bucket":  "The name of the bucket containing the object to be updated",
		"name":    "The name of the object to be updated",
		"version": "Used to reference a specific version of the object",
	},
	"update.securitygroup": {
		"id":            "The ID of the security group to be updated",
		"cidr":          "The CIDR IPv4 address range",
		"securitygroup": "The ID of the source security group. Cannot be used when using cidr param",
		"protocol":      "The IP protocol name or number",
		"inbound":       "Set inbound to either authorize or revoke, to update the security group ingress rules",
		"outbound":      "Set outbound to either authorize or revoke, to update the security group egress rules",
		"portrange":     "The portrange for the rule to update: any, 80, 22-23...",
	},
	"update.stack": {
		"capabilities":            "A list of values that you must specify before AWS CloudFormation can update certain stacks",
		"parameters":              "A list of Parameters that specify input parameters for the stack given using this format: [key1:val1,key2:val2,...]",
		"policy-file":             "The path to the file containing the stack policy body",
		"policy-update-file":      "The path to the file containing the temporary overriding stack policy",
		"template-file":           "The path to the file containing the template body with a minimum size of 1 byte and a maximum size of 51,200 bytes",
		"stack-file":              "The path to the file containing Parameters/Tags/StackPolices definition (http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/continuous-delivery-codepipeline-cfn-artifacts.html#w2ab2c13c15c15). Values passed via CLI has higher priority than ones defined in StackFile",
		"rollback-triggers":       "List of CloudWatch Alarm ARNs to monitor during and after update",
		"rollback-monitoring-min": "Time to monitor rollback-triggers during and after update",
	},
	"update.targetgroup": {
		"id": "The Amazon Resource Name (ARN) of the target group",
		"deregistrationdelay": "The amount time for Elastic Load Balancing to wait before changing the state of a deregistering target from draining to unused. The range is 0-3600 seconds. The default value is 300 seconds",
		"healthcheckinterval": "The approximate amount of time, in seconds, between health checks of an individual target",
		"healthcheckpath":     "The ping path that is the destination on the targets for health checks",
		"healthcheckport":     "The port the load balancer uses when performing health checks on targets",
		"healthcheckprotocol": "The protocol the load balancer uses when performing health checks on targets",
		"healthchecktimeout":  "The amount of time, in seconds, during which no response from a target means a failed health check",
		"healthythreshold":    "The number of consecutive health checks successes required before considering an unhealthy target healthy",
		"matcher":             "The HTTP codes to use when checking for a successful response from a target",
		"name":                "The name of the target group",
		"port":                "The port on which the targets receive traffic",
		"protocol":            "The protocol to use for routing traffic to the targets",
		"unhealthythreshold":  "The number of consecutive health check failures required before considering a target unhealthy",
		"stickiness":          "Indicates whether sticky sessions (of type load balancer cookies) are enabled",
		"stickinessduration":  "The time period, in seconds, during which requests from a client should be routed to the same target. After this time period expires, the load balancer-generated cookie is considered stale. The range is 1 second to 1 week (604800 seconds). The default value is 1 day (86400 seconds)",
	},
}
