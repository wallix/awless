package awsdoc

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

	"create.launchconfiguration.distro": {"amazonlinux", "canonical", "redhat", "debian", "suselinux", "windows"},
	"create.launchconfiguration.type":   {"t2.nano", "t2.micro", "t2.small", "t2.medium", "t2.large", "t2.xlarge", "t2.2xlarge", "m4.large", "m4.xlarge", "c4.large", "c4.xlarge"},
}
