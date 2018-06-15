package awsdoc

import (
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
)

var (
	timeouts      = []string{"10", "60", "180", "300", "600", "900"}
	boolean       = []string{"true", "false"}
	services      = []string{"iam", "ec2", "s3", "route53", "elbv2", "rds", "autoscaling", "lambda", "sns", "sqs", "cloudwatch", "cloudfront", "ecr", "ecs", "applicationautoscaling", "acm", "sts", "cloudformation"}
	instanceTypes = []string{"t2.nano", "t2.micro", "t2.small", "t2.medium", "t2.large", "t2.xlarge", "t2.2xlarge", "m4.large", "m4.xlarge", "c4.large", "c4.xlarge"}
	s3ACLs        = []string{"private", "public-read", "public-read-write", "aws-exec-read", "authenticated-read", "bucket-owner-read", "bucket-owner-full-control", "log-delivery-write"}
	distros       = []string{"amazonlinux", "canonical:ubuntu", "redhat:rhel", "debian:debian", "centos:centos", "coreos:coreos", "suselinux", "windows:server"}
	regions       = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1", "ca-central-1", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1"}
)

var EnumDoc = map[string][]string{

	"attach.mfadevice.mfa-code-1": {""},
	"attach.mfadevice.mfa-code-2": {""},
	"attach.mfadevice.no-prompt":  boolean,

	"attach.policy.access":  {"readonly", "full"},
	"attach.policy.service": services,

	"check.database.state":   {"available", "backing-up", "creating", "deleting", "failed", "maintenance", "modifying", "rebooting", "renaming", "resetting-master-credentials", "restore-error", "storage-full", "upgrading", "not-found"},
	"check.database.timeout": timeouts,

	"check.certificate.state":   {"issued", "pending_validation", "not-found"},
	"check.certificate.timeout": timeouts,

	"check.distribution.state":   {"Deployed", "InProgress", "not-found"},
	"check.distribution.timeout": timeouts,

	"check.instance.state":   {"pending", "running", "shutting-down", "terminated", "stopping", "stopped", "not-found"},
	"check.instance.timeout": timeouts,

	"check.loadbalancer.state":   {"provisioning", "active", "failed", "not-found"},
	"check.loadbalancer.timeout": timeouts,

	"check.natgateway.state":   {"pending", "failed", "available", "deleting", "deleted", "not-found"},
	"check.natgateway.timeout": timeouts,

	"check.networkinterface.state":   {"available", "attaching", "detaching", "in-use", "not-found"},
	"check.networkinterface.timeout": timeouts,

	"check.scalinggroup.count":   {"0"},
	"check.scalinggroup.timeout": timeouts,

	"check.securitygroup.state":   {"unused"},
	"check.securitygroup.timeout": timeouts,

	"check.volume.state":   {"available", "in-use", "not-found"},
	"check.volume.timeout": timeouts,

	"create.accesskey.save": boolean,

	"create.alarm.operator":           {"GreaterThanThreshold", "LessThanThreshold", "LessThanOrEqualToThreshold", "GreaterThanOrEqualToThreshold"},
	"create.alarm.statistic-function": {"Minimum", "Maximum", "Sum", "Average", "SampleCount", "pNN.NN"},
	"create.alarm.unit":               {"Seconds", "Microseconds", "Milliseconds", "Bytes", "Kilobytes", "Megabytes", "Gigabytes", "Terabytes", "Bits", "Kilobits", "Megabits", "Gigabits", "Terabits", "Percent", "Count", "Bytes/Second", "Kilobytes/Second", "Megabytes/Second", "Gigabytes/Second", "Terabytes/Second", "Bits/Second", "Kilobits/Second", "Megabits/Second", "Gigabits/Second", "Terabits/Second", "Count/Second", "None"},

	"create.appscalingtarget.dimension":         {"ecs:service:DesiredCount", "ec2:spot-fleet-request:TargetCapacity", "elasticmapreduce:instancegroup:InstanceCount", "appstream:fleet:DesiredCapacity", "dynamodb:table:ReadCapacityUnits", "dynamodb:table:WriteCapacityUnits", "dynamodb:index:ReadCapacityUnits", "dynamodb:index:WriteCapacityUnits"},
	"create.appscalingtarget.service-namespace": {"ecs", "ec2", "elasticmapreduce", "appstream", "dynamodb"},

	"create.appscalingpolicy.dimension":                    {"ecs:service:DesiredCount", "ec2:spot-fleet-request:TargetCapacity", "elasticmapreduce:instancegroup:InstanceCount", "appstream:fleet:DesiredCapacity", "dynamodb:table:ReadCapacityUnits", "dynamodb:table:WriteCapacityUnits", "dynamodb:index:ReadCapacityUnits", "dynamodb:index:WriteCapacityUnits"},
	"create.appscalingpolicy.service-namespace":            {"ecs", "ec2", "elasticmapreduce", "appstream", "dynamodb"},
	"create.appscalingpolicy.type":                         {"StepScaling"},
	"create.appscalingpolicy.stepscaling-adjustment-type":  {"ChangeInCapacity", "ExactCapacity", "PercentChangeInCapacity"},
	"create.appscalingpolicy.stepscaling-adjustments":      {"0::+1", ":0:-1", "75::+1"},
	"create.appscalingpolicy.stepscaling-aggregation-type": {"Minimum", "Maximum", "Average"},

	"create.bucket.acl": s3ACLs,

	"create.database.engine":             {"mysql", "mariadb", "postgres", "aurora", "oracle-se1", "oracle-se2", "oracle-se", "oracle-ee", "sqlserver-ee", "sqlserver-se", "sqlserver-ex", "sqlserver-web"},
	"create.database.copytagstosnapshot": boolean,
	"create.database.encrypted":          boolean,
	"create.database.license":            {"license-included", "bring-your-own-license", "general-public-license"},
	"create.database.multiaz":            boolean,
	"create.database.public":             boolean,
	"create.database.storagetype":        {"standard", "gp2", "io1"},
	"create.database.type":               {"db.t1.micro", "db.m1.small", "db.m1.medium", "db.m1.large", "db.m1.xlarge", "db.m2.xlarge |db.m2.2xlarge", "db.m2.4xlarge", "db.m3.medium", "db.m3.large", "db.m3.xlarge", "db.m3.2xlarge", "db.m4.large", "db.m4.xlarge", "db.m4.2xlarge", "db.m4.4xlarge", "db.m4.10xlarge", "db.r3.large", "db.r3.xlarge", "db.r3.2xlarge", "db.r3.4xlarge", "db.r3.8xlarge", "db.t2.micro", "db.t2.small", "db.t2.medium", "db.t2.large"},

	"create.distribution.default-file":    {"index.html"},
	"create.distribution.enable":          boolean,
	"create.distribution.forward-cookies": {"all", "none", "whitelist"},
	"create.distribution.forward-queries": boolean,
	"create.distribution.https-behaviour": {"allow-all", "redirect-to-https", "https-only"},
	"create.distribution.price-class":     {"PriceClass_All", "PriceClass_100", "PriceClass_200"},

	"create.elasticip.domain": {"vpc", "ec2-classic"},

	"create.function.runtime": {"nodejs", "nodejs4.3", "nodejs6.10", "java8", "python2.7", "python3.6", "dotnetcore1.0", "nodejs4.3-edge"},

	"create.instance.distro":   distros,
	"create.instance.type":     instanceTypes,
	"create.instance.lock":     boolean,
	"create.instance.userdata": {""},

	"create.image.reboot": boolean,

	"create.keypair.encrypted": boolean,

	"create.launchconfiguration.distro":   distros,
	"create.launchconfiguration.type":     {"t2.nano", "t2.micro", "t2.small", "t2.medium", "t2.large", "t2.xlarge", "t2.2xlarge", "m4.large", "m4.xlarge", "c4.large", "c4.xlarge"},
	"create.launchconfiguration.userdata": {""},
	"create.launchconfiguration.public":   boolean,

	"create.listener.actiontype": {"forward"},
	"create.listener.protocol":   {"HTTP", "HTTPS"},
	"create.listener.sslpolicy":  {"ELBSecurityPolicy-2016-08", "ELBSecurityPolicy-TLS-1-2-2017-01", "ELBSecurityPolicy-TLS-1-1-2017-01", "ELBSecurityPolicy-2015-05", "ELBSecurityPolicy-TLS-1-0-2015-04"},

	"create.policy.action":   {""},
	"create.policy.effect":   {"Allow", "Deny"},
	"create.policy.resource": {"*"},

	"create.record.type": {"A", "AAAA", "CNAME", "MX", "NAPTR", "NS", "PTR", "SOA", "SPF", "SRV", "TXT"},

	"create.s3object.acl": s3ACLs,

	"create.scalinggroup.healthcheck-type": {"EC2", "ELB"},

	"create.scalingpolicy.adjustment-type": {"ChangeInCapacity", "ExactCapacity", "PercentChangeInCapacity"},

	"create.stack.capabilities": {"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"},
	"create.stack.on-failure":   {"DO_NOTHING", "ROLLBACK", "DELETE"},

	"create.subnet.public": boolean,

	"create.subscription.protocol": {"http", "https", "email", "email-json", "sms", "sqs", "lambda"},

	"create.zone.isprivate": boolean,

	"copy.image.source-id":     {""},
	"copy.image.source-region": regions,

	"delete.containertask.all-versions": boolean,

	"delete.database.skip-snapshot": boolean,

	"delete.image.delete-snapshots": boolean,

	"delete.policy.all-versions": boolean,

	"delete.record.type": {"A", "AAAA", "CNAME", "MX", "NAPTR", "NS", "PTR", "SOA", "SPF", "SRV", "TXT"},

	"detach.networkinterface.force": boolean,

	"detach.policy.access":  {"readonly", "full"},
	"detach.policy.service": services,

	"import.image.architecture": {"i386", "x86_64"},
	"import.image.license":      {"AWS", "BYOL"},
	"import.image.platform":     {"Windows", "Linux"},

	"restart.database.with-failover": boolean,

	"start.containertask.type": {"task", "service"},

	"stop.containertask.type": {"task", "service"},

	"update.bucket.acl":            {"private", "public-read", "public-read-write", "aws-exec-read", "authenticated-read", "bucket-owner-read", "bucket-owner-full-control", "log-delivery-write"},
	"update.bucket.public-website": boolean,
	"update.bucket.index-suffix":   {"index.html"},

	"update.distribution.default-file":    {"index.html"},
	"update.distribution.forward-cookies": {"all", "none", "whitelist"},
	"update.distribution.forward-queries": boolean,
	"update.distribution.https-behaviour": {"allow-all", "redirect-to-https", "https-only"},
	"update.distribution.price-class":     {"PriceClass_All", "PriceClass_100", "PriceClass_200"},
	"update.distribution.enable":          boolean,

	"update.image.operation": {"add", "remove"},

	"update.instance.type": instanceTypes,

	"update.policy.effect": {"Allow", "Deny"},

	"update.s3object.acl": s3ACLs,

	"update.securitygroup.inbound":   {"revoke", "authorize"},
	"update.securitygroup.outbound":  {"revoke", "authorize"},
	"update.securitygroup.protocol":  {"tcp", "udp", "icmp", "any"},
	"update.securitygroup.portrange": {""},

	"update.stack.capabilities": {"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"},

	"update.targetgroup.stickiness": boolean,

	"update.subnet.public": boolean,

	"update.record.type": {"A", "AAAA", "CNAME", "MX", "NAPTR", "NS", "PTR", "SOA", "SPF", "SRV", "TXT"},
}

type ParamType struct {
	ResourceType, PropertyName string
}

var ParamTypeDoc = map[string]*ParamType{
	"attach.mfadevice.user": {ResourceType: cloud.User, PropertyName: properties.Name},

	"attach.policy.group": {ResourceType: cloud.Group, PropertyName: properties.Name},
	"attach.policy.role":  {ResourceType: cloud.Role, PropertyName: properties.Name},
	"attach.policy.user":  {ResourceType: cloud.User, PropertyName: properties.Name},
	"attach.policy.arn":   {ResourceType: cloud.Policy, PropertyName: properties.Arn},

	"attach.user.name":  {ResourceType: cloud.User, PropertyName: properties.Name},
	"attach.user.group": {ResourceType: cloud.Group, PropertyName: properties.Name},

	"attach.role.instanceprofile": {ResourceType: cloud.InstanceProfile, PropertyName: properties.Name},

	"create.accesskey.user": {ResourceType: cloud.User, PropertyName: properties.Name},

	"create.instance.role": {ResourceType: cloud.Role, PropertyName: properties.Name},

	"create.record.values": {ResourceType: cloud.Record, PropertyName: properties.Records},

	"delete.policy.arn":   {ResourceType: cloud.Policy, PropertyName: properties.Arn},
	"detach.policy.arn":   {ResourceType: cloud.Policy, PropertyName: properties.Arn},
	"detach.policy.group": {ResourceType: cloud.Group, PropertyName: properties.Name},
	"detach.policy.role":  {ResourceType: cloud.Role, PropertyName: properties.Name},
	"detach.policy.user":  {ResourceType: cloud.User, PropertyName: properties.Name},

	"detach.role.instanceprofile": {ResourceType: cloud.InstanceProfile, PropertyName: properties.Name},

	"update.policy.arn": {ResourceType: cloud.Policy, PropertyName: properties.Arn},

	"update.securitygroup.cidr": {ResourceType: cloud.Subnet, PropertyName: properties.CIDR},
}
