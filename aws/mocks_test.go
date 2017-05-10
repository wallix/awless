package aws

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (m *mockEc2) DescribeInstancesPages(input *ec2.DescribeInstancesInput, fn func(p *ec2.DescribeInstancesOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{{Instances: m.instances}}}, true)
	return nil
}

func (m *mockElbv2) DescribeListenersPages(input *elbv2.DescribeListenersInput, fn func(p *elbv2.DescribeListenersOutput, lastPage bool) (shouldContinue bool)) error {
	listeners := make(map[string][]*elbv2.Listener)
	for _, l := range m.listeners {
		listeners[awssdk.StringValue(l.LoadBalancerArn)] = append(listeners[awssdk.StringValue(l.LoadBalancerArn)], l)
	}
	fn(&elbv2.DescribeListenersOutput{Listeners: listeners[awssdk.StringValue(input.LoadBalancerArn)]}, true)
	return nil
}

func (m *mockElbv2) DescribeTargetHealth(input *elbv2.DescribeTargetHealthInput) (*elbv2.DescribeTargetHealthOutput, error) {
	return &elbv2.DescribeTargetHealthOutput{TargetHealthDescriptions: m.targethealthdescriptions[awssdk.StringValue(input.TargetGroupArn)]}, nil
}

func (m *mockRoute53) ListResourceRecordSetsPages(input *route53.ListResourceRecordSetsInput, fn func(p *route53.ListResourceRecordSetsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&route53.ListResourceRecordSetsOutput{ResourceRecordSets: m.resourcerecordsets[awssdk.StringValue(input.HostedZoneId)]}, true)
	return nil
}

func (m *mockIam) ListUsersPages(input *iam.ListUsersInput, fn func(p *iam.ListUsersOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&iam.ListUsersOutput{Users: m.users}, true)
	return nil
}

func (m *mockIam) ListPoliciesPages(input *iam.ListPoliciesInput, fn func(p *iam.ListPoliciesOutput, lastPage bool) (shouldContinue bool)) error {
	var policies []*iam.Policy
	for _, p := range m.managedpolicydetails {
		policy := &iam.Policy{PolicyId: p.PolicyId, PolicyName: p.PolicyName}
		policies = append(policies, policy)
	}
	fn(&iam.ListPoliciesOutput{Policies: policies}, true)
	return nil
}

func (m *mockIam) GetAccountAuthorizationDetailsPages(input *iam.GetAccountAuthorizationDetailsInput, fn func(p *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&iam.GetAccountAuthorizationDetailsOutput{GroupDetailList: m.groupdetails, Policies: m.managedpolicydetails, RoleDetailList: m.roledetails, UserDetailList: m.userdetails}, true)
	return nil
}

func (m *mockS3) GetBucketAcl(input *s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error) {
	return &s3.GetBucketAclOutput{Grants: m.grants[awssdk.StringValue(input.Bucket)]}, nil
}

func (m *mockS3) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	var buckets []*s3.Bucket
	for _, b := range m.buckets {
		buckets = append(buckets, b...)
	}
	return &s3.ListBucketsOutput{Buckets: buckets}, nil
}
func (m *mockS3) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return &s3.ListObjectsOutput{Contents: m.objects[awssdk.StringValue(input.Bucket)]}, nil
}
func (m *mockS3) GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {
	for region, buckets := range m.buckets {
		for _, bucket := range buckets {
			if awssdk.StringValue(bucket.Name) == awssdk.StringValue(input.Bucket) {
				return &s3.GetBucketLocationOutput{LocationConstraint: awssdk.String(region)}, nil
			}
		}
	}
	return nil, fmt.Errorf("bucket location mock: bucket %s not found", awssdk.StringValue(input.Bucket))
}

func (m *mockSqs) GetQueueAttributes(input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	return &sqs.GetQueueAttributesOutput{Attributes: m.attributes[awssdk.StringValue(input.QueueUrl)]}, nil
}
