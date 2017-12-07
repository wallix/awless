package sync_test

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/sync"
)

func BenchmarkSync(b *testing.B) {
	dir, err := ioutil.TempDir("", "synctest")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	os.Setenv("__AWLESS_HOME", dir)

	awsservices.InfraService = &awsservices.Infra{
		EC2API:         &ec2Mock{},
		ELBV2API:       &elbMock{},
		RDSAPI:         &rdsMock{},
		ECRAPI:         &ecrMock{},
		AutoScalingAPI: &autoscalingStub{},
	}
	awsservices.AccessService = &awsservices.Access{
		IAMAPI: &iamMock{},
		STSAPI: &stsMock{},
	}
	awsservices.StorageService = &awsservices.Storage{S3API: &s3Mock{}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sync.NewSyncer().Sync(awsservices.InfraService, awsservices.AccessService, awsservices.StorageService)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type ec2Mock struct {
	ec2iface.EC2API
}

func (*ec2Mock) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	vpcs := []*ec2.Vpc{
		{VpcId: awssdk.String("vpc_1")},
		{VpcId: awssdk.String("vpc_2")},
	}
	return &ec2.DescribeVpcsOutput{Vpcs: vpcs}, nil
}

func (*ec2Mock) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	subnets := []*ec2.Subnet{
		{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
		{SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
		{SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{SubnetId: awssdk.String("sub_4"), VpcId: nil},
	}
	return &ec2.DescribeSubnetsOutput{Subnets: subnets}, nil
}

func (m *ec2Mock) DescribeInstancesPages(input *ec2.DescribeInstancesInput, fn func(p *ec2.DescribeInstancesOutput, lastPage bool) (shouldContinue bool)) error {
	instances := []*ec2.Instance{
		{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("securitygroup_1")}}},
		{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{InstanceId: awssdk.String("inst_4"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("securitygroup_1")}, {GroupId: awssdk.String("securitygroup_2")}}, KeyName: awssdk.String("my_key")},
		{InstanceId: awssdk.String("inst_5"), SubnetId: nil, VpcId: nil, KeyName: awssdk.String("unexisting_key")},
		{
			InstanceId:         awssdk.String("inst_6"),
			Tags:               []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("inst_6_name")}},
			InstanceType:       awssdk.String("t2.micro"),
			SubnetId:           awssdk.String("sub_3"),
			VpcId:              awssdk.String("vpc_2"),
			PublicIpAddress:    awssdk.String("1.2.3.4"),
			PrivateIpAddress:   awssdk.String("10.0.0.1"),
			ImageId:            awssdk.String("ami-1234"),
			LaunchTime:         awssdk.Time(time.Now()),
			State:              &ec2.InstanceState{Name: awssdk.String("running")},
			KeyName:            awssdk.String("my_key"),
			SecurityGroups:     []*ec2.GroupIdentifier{{GroupId: awssdk.String("securitygroup_1")}},
			Placement:          &ec2.Placement{Affinity: awssdk.String("inst_affinity"), AvailabilityZone: awssdk.String("inst_az"), GroupName: awssdk.String("inst_group"), HostId: awssdk.String("inst_host")},
			Architecture:       awssdk.String("x86"),
			Hypervisor:         awssdk.String("xen"),
			IamInstanceProfile: &ec2.IamInstanceProfile{Arn: awssdk.String("arn:instance:profile")},
			InstanceLifecycle:  awssdk.String("lifecycle"),
			NetworkInterfaces:  []*ec2.InstanceNetworkInterface{{NetworkInterfaceId: awssdk.String("my-network-interface")}},
			PublicDnsName:      awssdk.String("my-instance.dns"),
			RootDeviceName:     awssdk.String("/dev/xvda"),
			RootDeviceType:     awssdk.String("ebs"),
		},
	}
	fn(&ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{{Instances: instances}}}, true)
	return nil
}

func (*ec2Mock) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	sgroups := []*ec2.SecurityGroup{
		{GroupId: awssdk.String("securitygroup_1"), GroupName: awssdk.String("my_securitygroup"), VpcId: awssdk.String("vpc_1")},
		{GroupId: awssdk.String("securitygroup_2"), VpcId: awssdk.String("vpc_1")},
	}
	return &ec2.DescribeSecurityGroupsOutput{SecurityGroups: sgroups}, nil
}

func (*ec2Mock) DescribeImportImageTasks(input *ec2.DescribeImportImageTasksInput) (*ec2.DescribeImportImageTasksOutput, error) {
	return &ec2.DescribeImportImageTasksOutput{ImportImageTasks: []*ec2.ImportImageTask{}}, nil
}

func (*ec2Mock) DescribeTargetGroups(input *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	tgroups := []*elbv2.TargetGroup{
		{TargetGroupArn: awssdk.String("tg_1"), VpcId: awssdk.String("vpc_1"), LoadBalancerArns: []*string{awssdk.String("lb_1"), awssdk.String("lb_3")}},
		{TargetGroupArn: awssdk.String("tg_2"), VpcId: awssdk.String("vpc_2"), LoadBalancerArns: []*string{awssdk.String("lb_2")}},
	}
	return &elbv2.DescribeTargetGroupsOutput{TargetGroups: tgroups}, nil
}

func (*ec2Mock) DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	return &ec2.DescribeAddressesOutput{Addresses: []*ec2.Address{}}, nil
}

func (*ec2Mock) DescribeKeyPairs(input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	keypairs := []*ec2.KeyPairInfo{
		{KeyName: awssdk.String("my_key")},
	}
	return &ec2.DescribeKeyPairsOutput{KeyPairs: keypairs}, nil
}

func (*ec2Mock) DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	igws := []*ec2.InternetGateway{
		{InternetGatewayId: awssdk.String("igw_1"), Attachments: []*ec2.InternetGatewayAttachment{{VpcId: awssdk.String("vpc_2")}}},
	}
	return &ec2.DescribeInternetGatewaysOutput{InternetGateways: igws}, nil
}

func (*ec2Mock) DescribeRouteTables(input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	rTables := []*ec2.RouteTable{
		{RouteTableId: awssdk.String("rt_1"), VpcId: awssdk.String("vpc_1"), Associations: []*ec2.RouteTableAssociation{{RouteTableId: awssdk.String("rt_1"), SubnetId: awssdk.String("sub_1")}}},
	}
	return &ec2.DescribeRouteTablesOutput{RouteTables: rTables}, nil
}

func (*ec2Mock) DescribeAvailabilityZones(input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	zones := []*ec2.AvailabilityZone{
		{ZoneName: awssdk.String("us-west-1a"), State: awssdk.String("available"), RegionName: awssdk.String("us-west-1"), Messages: []*ec2.AvailabilityZoneMessage{{Message: awssdk.String("msg 1")}, {Message: awssdk.String("msg 2")}}},
		{ZoneName: awssdk.String("us-west-1b")},
	}
	return &ec2.DescribeAvailabilityZonesOutput{AvailabilityZones: zones}, nil
}

func (*ec2Mock) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	images := []*ec2.Image{
		{ImageId: awssdk.String("img_1")},
		{ImageId: awssdk.String("img_2"), Name: awssdk.String("img_2_name"), Architecture: awssdk.String("img_2_arch"), Hypervisor: awssdk.String("img_2_hyper"), CreationDate: awssdk.String("2010-04-01T12:05:01.000Z")},
	}
	return &ec2.DescribeImagesOutput{Images: images}, nil
}

func (*ec2Mock) DescribeSnapshotsPages(input *ec2.DescribeSnapshotsInput, fn func(p *ec2.DescribeSnapshotsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&ec2.DescribeSnapshotsOutput{Snapshots: []*ec2.Snapshot{}}, true)
	return nil
}

func (*ec2Mock) DescribeVolumesPages(input *ec2.DescribeVolumesInput, fn func(p *ec2.DescribeVolumesOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{}}, true)
	return nil
}

type autoscalingStub struct {
	autoscalingiface.AutoScalingAPI
}

func (*autoscalingStub) DescribeLaunchConfigurationsPages(input *autoscaling.DescribeLaunchConfigurationsInput, fn func(p *autoscaling.DescribeLaunchConfigurationsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&autoscaling.DescribeLaunchConfigurationsOutput{LaunchConfigurations: []*autoscaling.LaunchConfiguration{}}, true)
	return nil
}

func (*autoscalingStub) DescribeAutoScalingGroupsPages(input *autoscaling.DescribeAutoScalingGroupsInput, fn func(p *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: []*autoscaling.Group{}}, true)
	return nil
}

func (*autoscalingStub) DescribePoliciesPages(input *autoscaling.DescribePoliciesInput, fn func(p *autoscaling.DescribePoliciesOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&autoscaling.DescribePoliciesOutput{ScalingPolicies: []*autoscaling.ScalingPolicy{}}, true)
	return nil
}

type ecrMock struct {
	ecriface.ECRAPI
}

func (*ecrMock) DescribeRepositoriesPages(input *ecr.DescribeRepositoriesInput, fn func(p *ecr.DescribeRepositoriesOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&ecr.DescribeRepositoriesOutput{Repositories: []*ecr.Repository{}}, true)
	return nil
}

type elbMock struct {
	elbv2iface.ELBV2API
}

func (*elbMock) DescribeLoadBalancersPages(input *elbv2.DescribeLoadBalancersInput, fn func(p *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool)) error {
	lbPages := [][]*elbv2.LoadBalancer{{
		{LoadBalancerArn: awssdk.String("lb_1"), LoadBalancerName: awssdk.String("my_loadbalancer"), VpcId: awssdk.String("vpc_1")},
		{LoadBalancerArn: awssdk.String("lb_2"), VpcId: awssdk.String("vpc_2")},
		{LoadBalancerArn: awssdk.String("lb_3"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*string{awssdk.String("securitygroup_1"), awssdk.String("securitygroup_2")}},
	}}

	for i, page := range lbPages {
		fn(&elbv2.DescribeLoadBalancersOutput{LoadBalancers: page, NextMarker: awssdk.String(strconv.Itoa(i + 1))},
			i < len(lbPages),
		)
	}
	return nil
}

func (*elbMock) DescribeListenersPages(input *elbv2.DescribeListenersInput, fn func(p *elbv2.DescribeListenersOutput, lastPage bool) (shouldContinue bool)) error {
	listeners := []*elbv2.Listener{
		{ListenerArn: awssdk.String("list_1"), LoadBalancerArn: awssdk.String("lb_1")}, {ListenerArn: awssdk.String("list_1.2"), LoadBalancerArn: awssdk.String("lb_1")},
		{ListenerArn: awssdk.String("list_2"), LoadBalancerArn: awssdk.String("lb_2")},
		{ListenerArn: awssdk.String("list_3"), LoadBalancerArn: awssdk.String("lb_3")},
	}
	fn(&elbv2.DescribeListenersOutput{Listeners: listeners}, true)
	return nil
}

func (*elbMock) DescribeTargetGroups(input *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	targetGroups := []*elbv2.TargetGroup{
		{TargetGroupArn: awssdk.String("tg_1"), VpcId: awssdk.String("vpc_1"), LoadBalancerArns: []*string{awssdk.String("lb_1"), awssdk.String("lb_3")}},
		{TargetGroupArn: awssdk.String("tg_2"), VpcId: awssdk.String("vpc_2"), LoadBalancerArns: []*string{awssdk.String("lb_2")}},
	}
	return &elbv2.DescribeTargetGroupsOutput{TargetGroups: targetGroups}, nil
}

func (*elbMock) DescribeTargetHealth(input *elbv2.DescribeTargetHealthInput) (*elbv2.DescribeTargetHealthOutput, error) {
	return &elbv2.DescribeTargetHealthOutput{TargetHealthDescriptions: []*elbv2.TargetHealthDescription{}}, nil
}

type rdsMock struct {
	rdsiface.RDSAPI
}

func (*rdsMock) DescribeDBInstancesPages(input *rds.DescribeDBInstancesInput, fn func(p *rds.DescribeDBInstancesOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&rds.DescribeDBInstancesOutput{DBInstances: []*rds.DBInstance{}}, true)
	return nil
}

func (*rdsMock) DescribeDBSubnetGroupsPages(input *rds.DescribeDBSubnetGroupsInput, fn func(p *rds.DescribeDBSubnetGroupsOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&rds.DescribeDBSubnetGroupsOutput{DBSubnetGroups: []*rds.DBSubnetGroup{}}, true)
	return nil
}

type iamMock struct {
	iamiface.IAMAPI
}

func (*iamMock) ListUsersPages(input *iam.ListUsersInput, fn func(p *iam.ListUsersOutput, lastPage bool) (shouldContinue bool)) error {
	users := []*iam.User{
		{UserId: awssdk.String("usr_1"), PasswordLastUsed: awssdk.Time(time.Unix(1486139077, 0).UTC())},
		{UserId: awssdk.String("usr_2")},
		{UserId: awssdk.String("usr_3")},
		{UserId: awssdk.String("usr_4")},
		{UserId: awssdk.String("usr_5")},
		{UserId: awssdk.String("usr_6")},
		{UserId: awssdk.String("usr_7")},
		{UserId: awssdk.String("usr_8")},
		{UserId: awssdk.String("usr_9")},
		{UserId: awssdk.String("usr_10")},
		{UserId: awssdk.String("usr_11")},
	}
	fn(&iam.ListUsersOutput{Users: users}, true)
	return nil
}

func (*iamMock) GetAccountAuthorizationDetailsPages(input *iam.GetAccountAuthorizationDetailsInput, fn func(p *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool)) error {
	managedPolicies := []*iam.ManagedPolicyDetail{
		{PolicyId: awssdk.String("managed_policy_1"), PolicyName: awssdk.String("nmanaged_policy_1")},
		{PolicyId: awssdk.String("managed_policy_2"), PolicyName: awssdk.String("nmanaged_policy_2")},
		{PolicyId: awssdk.String("managed_policy_3"), PolicyName: awssdk.String("nmanaged_policy_3")},
	}

	groups := []*iam.GroupDetail{
		{GroupId: awssdk.String("group_1"), GroupName: awssdk.String("ngroup_1"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}}},
		{GroupId: awssdk.String("group_2"), GroupName: awssdk.String("ngroup_2"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_2")}}},
		{GroupId: awssdk.String("group_3"), GroupName: awssdk.String("ngroup_3"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_3")}}},
		{GroupId: awssdk.String("group_4"), GroupName: awssdk.String("ngroup_4"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}}},
	}

	roles := []*iam.RoleDetail{
		{RoleId: awssdk.String("role_1"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}}},
		{RoleId: awssdk.String("role_2"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}},
		{RoleId: awssdk.String("role_3"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_2")}}},
		{RoleId: awssdk.String("role_4"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}}},
	}

	usersDetails := []*iam.UserDetail{
		{
			UserId:                  awssdk.String("usr_1"),
			GroupList:               []*string{awssdk.String("ngroup_1"), awssdk.String("ngroup_2")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}, {PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_2"),
			GroupList:      []*string{awssdk.String("ngroup_1")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}},
		},
		{
			UserId:                  awssdk.String("usr_3"),
			GroupList:               []*string{awssdk.String("ngroup_1"), awssdk.String("ngroup_4")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}, {PolicyName: awssdk.String("nmanaged_policy_2")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}, {PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_4"),
			GroupList:      []*string{awssdk.String("ngroup_2")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_5"),
			GroupList:      []*string{awssdk.String("ngroup_2")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:                  awssdk.String("usr_6"),
			GroupList:               []*string{awssdk.String("ngroup_2")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_3")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_7"),
			GroupList:      []*string{awssdk.String("ngroup_2"), awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}, {PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_8"),
			GroupList:      []*string{awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_9"),
			GroupList:      []*string{awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId: awssdk.String("usr_10"), //users not in any groups
		},
		{
			UserId: awssdk.String("usr_11"),
		},
	}

	fn(&iam.GetAccountAuthorizationDetailsOutput{
		GroupDetailList: groups,
		Policies:        managedPolicies,
		RoleDetailList:  roles,
		UserDetailList:  usersDetails}, true)
	return nil
}

func (*iamMock) ListAccessKeysPages(input *iam.ListAccessKeysInput, fn func(p *iam.ListAccessKeysOutput, lastPage bool) (shouldContinue bool)) error {
	return nil
}

func (*iamMock) ListPoliciesPages(input *iam.ListPoliciesInput, fn func(p *iam.ListPoliciesOutput, lastPage bool) (shouldContinue bool)) error {
	policies := []*iam.Policy{
		{PolicyId: awssdk.String("managed_policy_1"), PolicyName: awssdk.String("nmanaged_policy_1")},
		{PolicyId: awssdk.String("managed_policy_2"), PolicyName: awssdk.String("nmanaged_policy_2")},
		{PolicyId: awssdk.String("managed_policy_3"), PolicyName: awssdk.String("nmanaged_policy_3")},
	}
	fn(&iam.ListPoliciesOutput{Policies: policies}, true)
	return nil
}

type stsMock struct {
	stsiface.STSAPI
}

type s3Mock struct {
	s3iface.S3API
}

func (*s3Mock) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	buckets := []*s3.Bucket{
		{Name: awssdk.String("bucket_1")},
		{Name: awssdk.String("bucket_2")},
		{Name: awssdk.String("bucket_3")},
		{Name: awssdk.String("bucket_4")},
		{Name: awssdk.String("bucket_5")},
	}
	return &s3.ListBucketsOutput{Buckets: buckets}, nil
}

func (*s3Mock) GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {
	return &s3.GetBucketLocationOutput{LocationConstraint: awssdk.String("us-west-1")}, nil
}
