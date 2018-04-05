package awsservices

import (
	"fmt"
	"strconv"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
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

/*func (m *mockIam) ListPoliciesPages(input *iam.ListPoliciesInput, fn func(p *iam.ListPoliciesOutput, lastPage bool) (shouldContinue bool)) error {
	if awssdk.BoolValue(input.OnlyAttached) == false {
		return nil
	}
	var policies []*iam.Policy
	for _, p := range m.managedpolicydetails {
		policy := &iam.Policy{PolicyId: p.PolicyId, PolicyName: p.PolicyName, Arn: p.Arn}
		policies = append(policies, policy)
	}
	fn(&iam.ListPoliciesOutput{Policies: policies}, true)
	return nil
}*/

func (m *mockIam) GetAccountAuthorizationDetailsPages(input *iam.GetAccountAuthorizationDetailsInput, fn func(p *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&iam.GetAccountAuthorizationDetailsOutput{
		GroupDetailList: m.groupdetails,
		Policies:        m.managedpolicydetails,
		RoleDetailList:  m.roledetails,
		UserDetailList:  m.userdetails,
	}, true)
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
func (m *mockS3) ListObjectsPages(input *s3.ListObjectsInput, fn func(*s3.ListObjectsOutput, bool) bool) error {
	fn(&s3.ListObjectsOutput{Contents: m.objects[awssdk.StringValue(input.Bucket)]}, true)
	return nil
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

func (m *mockCloudfront) ListDistributionsPages(input *cloudfront.ListDistributionsInput, fn func(p *cloudfront.ListDistributionsOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*cloudfront.DistributionSummary
	for i := 0; i < len(m.distributionsummarys); i += 2 {
		page := []*cloudfront.DistributionSummary{m.distributionsummarys[i]}
		if i+1 < len(m.distributionsummarys) {
			page = append(page, m.distributionsummarys[i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&cloudfront.ListDistributionsOutput{DistributionList: &cloudfront.DistributionList{Items: page, NextMarker: awssdk.String(strconv.Itoa(i + 1))}},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockEcs) DescribeClusters(input *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error) {
	var clusters []*ecs.Cluster
	for _, cluster := range m.clusters {
		for _, inputC := range input.Clusters {
			if awssdk.StringValue(cluster.ClusterArn) == awssdk.StringValue(inputC) {
				clusters = append(clusters, cluster)
			}
		}
	}
	return &ecs.DescribeClustersOutput{Clusters: clusters}, nil
}

func (m *mockEcs) DescribeTaskDefinition(input *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	for _, def := range m.taskdefinitions {
		family := awssdk.StringValue(def.Family)
		familyRevision := family + ":" + fmt.Sprint(awssdk.Int64Value(def.Revision))
		if family == awssdk.StringValue(input.TaskDefinition) || familyRevision == awssdk.StringValue(input.TaskDefinition) || awssdk.StringValue(def.TaskDefinitionArn) == awssdk.StringValue(input.TaskDefinition) {
			return &ecs.DescribeTaskDefinitionOutput{TaskDefinition: def}, nil
		}
	}
	return nil, fmt.Errorf("task definition not found")
}

func (m *mockEcs) ListTasksPages(input *ecs.ListTasksInput, fn func(p *ecs.ListTasksOutput, lastPage bool) (shouldContinue bool)) error {
	if awssdk.StringValue(input.DesiredStatus) == "STOPPED" {
		return nil
	}
	var pages [][]*string
	for i := 0; i < len(m.tasksNames[awssdk.StringValue(input.Cluster)]); i += 2 {
		page := []*string{m.tasksNames[awssdk.StringValue(input.Cluster)][i]}
		if i+1 < len(m.tasksNames[awssdk.StringValue(input.Cluster)]) {
			page = append(page, m.tasksNames[awssdk.StringValue(input.Cluster)][i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&ecs.ListTasksOutput{TaskArns: page, NextToken: awssdk.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockEcs) DescribeTasks(input *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	var tasks []*ecs.Task
	for _, task := range m.tasks[awssdk.StringValue(input.Cluster)] {
		for _, inputT := range input.Tasks {
			if awssdk.StringValue(task.TaskArn) == awssdk.StringValue(inputT) {
				tasks = append(tasks, task)
			}
		}
	}
	return &ecs.DescribeTasksOutput{Tasks: tasks}, nil
}

func (m *mockEcs) ListContainerInstancesPages(input *ecs.ListContainerInstancesInput, fn func(p *ecs.ListContainerInstancesOutput, lastPage bool) (shouldContinue bool)) error {
	var pages [][]*string
	for i := 0; i < len(m.containerinstancesNames[awssdk.StringValue(input.Cluster)]); i += 2 {
		page := []*string{m.containerinstancesNames[awssdk.StringValue(input.Cluster)][i]}
		if i+1 < len(m.containerinstancesNames[awssdk.StringValue(input.Cluster)]) {
			page = append(page, m.containerinstancesNames[awssdk.StringValue(input.Cluster)][i+1])
		}
		pages = append(pages, page)
	}
	for i, page := range pages {
		fn(&ecs.ListContainerInstancesOutput{ContainerInstanceArns: page, NextToken: awssdk.String(strconv.Itoa(i + 1))},
			i < len(pages),
		)
	}
	return nil
}

func (m *mockEcs) DescribeContainerInstances(input *ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error) {
	return &ecs.DescribeContainerInstancesOutput{ContainerInstances: m.containerinstances[awssdk.StringValue(input.Cluster)]}, nil
}
