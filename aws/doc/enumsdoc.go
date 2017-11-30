package awsdoc

import (
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
)

var EnumDoc = map[string][]string{
	"update.securitygroup.inbound":  {"revoke", "authorize"},
	"update.securitygroup.outbound": {"revoke", "authorize"},
	"update.securitygroup.protocol": {"tcp", "udp", "icmp", "any"},

	"attach.policy.access":  {"readonly", "full"},
	"attach.policy.service": {"iam", "ec2", "s3", "route53", "elbv2", "rds", "autoscaling", "lambda", "sns", "sqs", "cloudwatch", "cloudfront", "ecr", "ecs", "applicationautoscaling", "acm", "sts", "cloudformation"},

	"create.subnet.public": {"true", "false"},
	"update.subnet.public": {"true", "false"},

	"create.instance.distro": {"amazonlinux", "canonical", "redhat", "debian", "suselinux", "windows"},
	"create.instance.type":   {"t2.nano", "t2.micro", "t2.small", "t2.medium", "t2.large", "t2.xlarge", "t2.2xlarge", "m4.large", "m4.xlarge", "c4.large", "c4.xlarge"},
	"create.instance.lock":   {"true", "false"},

	"update.image.operation": {"add", "remove"},

	"create.database.engine": {"mysql", "mariadb", "postgres", "aurora", "oracle-se1", "oracle-se2", "oracle-se", "oracle-ee", "sqlserver-ee", "sqlserver-se", "sqlserver-ex", "sqlserver-web"},

	"create.launchconfiguration.distro": {"amazonlinux", "canonical", "redhat", "debian", "suselinux", "windows"},
	"create.launchconfiguration.type":   {"t2.nano", "t2.micro", "t2.small", "t2.medium", "t2.large", "t2.xlarge", "t2.2xlarge", "m4.large", "m4.xlarge", "c4.large", "c4.xlarge"},

	"create.policy.action":   {""},
	"create.policy.effect":   {"Allow", "Deny"},
	"create.policy.resource": {"*"},

	"check.database.state":   {"available", "backing-up", "creating", "deleting", "failed", "maintenance", "modifying", "rebooting", "renaming", "resetting-master-credentials", "restore-error", "storage-full", "upgrading", "not-found"},
	"check.database.timeout": {"10", "60", "180", "300", "600", "900"},
}

type ParamType struct {
	ResourceType, PropertyName string
}

var ParamTypeDoc = map[string]*ParamType{
	"create.accesskey.user":     {ResourceType: cloud.User, PropertyName: properties.Name},
	"update.securitygroup.cidr": {ResourceType: cloud.Subnet, PropertyName: properties.CIDR},
}
