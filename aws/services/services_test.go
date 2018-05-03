/*
Copyright 2017 WALLIX

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

package awsservices

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"sort"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wallix/awless/aws/fetch"
	"github.com/wallix/awless/cloud"
	p "github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestBuildAccessRdfGraph(t *testing.T) {
	policyDoc := `{"Version":"2012-10-17","Statement":[{"Sid":"Stmt1486739000000","Effect":"Allow","Action":["ec2:*"],"Resource":["arn:aws:ec2:::vpc/vpc-123456","arn:aws:ec2:::subnet/*","arn:aws:ec2:::instance/*"]}]}`
	assumeRoleDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":"sts:AssumeRole","Condition":{"Bool":{"aws:MultiFactorAuthPresent":"true"}}}]}`
	managedPolicies := []*iam.ManagedPolicyDetail{
		{PolicyId: awssdk.String("managed_policy_1"), PolicyName: awssdk.String("nmanaged_policy_1"), AttachmentCount: awssdk.Int64(3)},
		{PolicyId: awssdk.String("managed_policy_2"), PolicyName: awssdk.String("nmanaged_policy_2"), AttachmentCount: awssdk.Int64(0), PolicyVersionList: []*iam.PolicyVersion{
			{Document: awssdk.String("this policy will be ignored")},
			{IsDefaultVersion: awssdk.Bool(true), Document: awssdk.String(url.QueryEscape(policyDoc))},
		}},
		{PolicyId: awssdk.String("managed_policy_3"), PolicyName: awssdk.String("nmanaged_policy_3"), Arn: awssdk.String("arn:aws:iam::aws:policy/managed_policy_3"), AttachmentCount: awssdk.Int64(1)},
	}

	groups := []*iam.GroupDetail{
		{GroupId: awssdk.String("group_1"), GroupName: awssdk.String("ngroup_1"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}}},
		{GroupId: awssdk.String("group_2"), GroupName: awssdk.String("ngroup_2"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_2")}}},
		{GroupId: awssdk.String("group_3"), GroupName: awssdk.String("ngroup_3"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_3")}}},
		{GroupId: awssdk.String("group_4"), GroupName: awssdk.String("ngroup_4"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}}},
	}

	roles := []*iam.RoleDetail{
		{RoleId: awssdk.String("role_1"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}}, AssumeRolePolicyDocument: awssdk.String(url.QueryEscape(assumeRoleDoc))},
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

	users := []*iam.User{
		{
			UserId:           awssdk.String("usr_1"),
			PasswordLastUsed: awssdk.Time(time.Unix(1486139077, 0).UTC()),
		},
		{
			UserId: awssdk.String("usr_2"),
		},
		{
			UserId: awssdk.String("usr_3"),
		},
		{
			UserId: awssdk.String("usr_4"),
		},
		{
			UserId: awssdk.String("usr_5"),
		},
		{
			UserId: awssdk.String("usr_6"),
		},
		{
			UserId: awssdk.String("usr_7"),
		},
		{
			UserId: awssdk.String("usr_8"),
		},
		{
			UserId: awssdk.String("usr_9"),
		},
		{
			UserId: awssdk.String("usr_10"), //users not in any groups
		},
		{
			UserId: awssdk.String("usr_11"),
		},
	}
	now := time.Now().UTC()
	mfaDevices := []*iam.VirtualMFADevice{
		{EnableDate: awssdk.Time(now), SerialNumber: awssdk.String("mfa-device-1"), User: &iam.User{UserId: awssdk.String("usr_1")}},
		{SerialNumber: awssdk.String("mfa-device-2")},
	}

	mock := &mockIam{groupdetails: groups, userdetails: usersDetails, roledetails: roles, managedpolicydetails: managedPolicies, users: users, virtualmfadevices: mfaDevices}
	access := Access{
		IAMAPI:  mock,
		region:  "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildAccessFetchFuncs(awsfetch.NewConfig(mock))),
	}

	g, err := access.Fetch(context.Background())
	if err != nil {
		fmt.Printf("%#v", err)
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery("policy", "group", "role", "user", cloud.MFADevice))
	if err != nil {
		t.Fatal(err)
	}

	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties()[p.InlinePolicies].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]cloud.Resource{
		"managed_policy_1": resourcetest.Policy("managed_policy_1").Prop(p.Name, "nmanaged_policy_1").Prop(p.Type, "Customer Managed").Prop(p.Attached, true).Build(),
		"managed_policy_2": resourcetest.Policy("managed_policy_2").Prop(p.Name, "nmanaged_policy_2").Prop(p.Type, "Customer Managed").Prop(p.Attached, false).Prop(p.Document, policyDoc).Build(),
		"managed_policy_3": resourcetest.Policy("managed_policy_3").Prop(p.Name, "nmanaged_policy_3").Prop(p.Arn, "arn:aws:iam::aws:policy/managed_policy_3").Prop(p.Type, "AWS Managed").Prop(p.Attached, true).Build(),
		"group_1":          resourcetest.Group("group_1").Prop(p.Name, "ngroup_1").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"group_2":          resourcetest.Group("group_2").Prop(p.Name, "ngroup_2").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"group_3":          resourcetest.Group("group_3").Prop(p.Name, "ngroup_3").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"group_4":          resourcetest.Group("group_4").Prop(p.Name, "ngroup_4").Prop(p.InlinePolicies, []string{"npolicy_4"}).Build(),
		"role_1":           resourcetest.Role("role_1").Prop(p.InlinePolicies, []string{"npolicy_1"}).Prop(p.TrustPolicy, assumeRoleDoc).Build(),
		"role_2":           resourcetest.Role("role_2").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"role_3":           resourcetest.Role("role_3").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"role_4":           resourcetest.Role("role_4").Prop(p.InlinePolicies, []string{"npolicy_4"}).Build(),
		"usr_1":            resourcetest.User("usr_1").Prop(p.InlinePolicies, []string{"npolicy_1", "npolicy_2"}).Prop(p.PasswordLastUsed, time.Unix(1486139077, 0).UTC()).Build(),
		"usr_2":            resourcetest.User("usr_2").Prop(p.InlinePolicies, []string{"npolicy_1"}).Build(),
		"usr_3":            resourcetest.User("usr_3").Prop(p.InlinePolicies, []string{"npolicy_1", "npolicy_4"}).Build(),
		"usr_4":            resourcetest.User("usr_4").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"usr_5":            resourcetest.User("usr_5").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"usr_6":            resourcetest.User("usr_6").Prop(p.InlinePolicies, []string{"npolicy_2"}).Build(),
		"usr_7":            resourcetest.User("usr_7").Prop(p.InlinePolicies, []string{"npolicy_2", "npolicy_4"}).Build(),
		"usr_8":            resourcetest.User("usr_8").Prop(p.InlinePolicies, []string{"npolicy_4"}).Build(),
		"usr_9":            resourcetest.User("usr_9").Prop(p.InlinePolicies, []string{"npolicy_4"}).Build(),
		"usr_10":           resourcetest.User("usr_10").Build(),
		"usr_11":           resourcetest.User("usr_11").Build(),
		"mfa-device-1":     resourcetest.MfaDevice("mfa-device-1").Prop(p.AttachedAt, now).Build(),
		"mfa-device-2":     resourcetest.MfaDevice("mfa-device-2").Build(),
	}

	expectedChildren := map[string][]string{}

	expectedAppliedOn := map[string][]string{
		"group_1":          {"usr_1", "usr_2", "usr_3"},
		"group_2":          {"usr_1", "usr_4", "usr_5", "usr_6", "usr_7"},
		"group_4":          {"usr_3", "usr_7", "usr_8", "usr_9"},
		"managed_policy_1": {"group_1", "role_1", "usr_1", "usr_3"},
		"managed_policy_2": {"group_2", "role_3", "usr_3"},
		"managed_policy_3": {"group_3", "usr_6"},
		"mfa-device-1":     {"usr_1"},
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildInfraRdfGraph(t *testing.T) {
	now := time.Now().UTC()
	instances := []*ec2.Instance{
		{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("securitygroup_1")}}},
		{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{InstanceId: awssdk.String("inst_4"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("securitygroup_1")}, {GroupId: awssdk.String("securitygroup_2")}}, KeyName: awssdk.String("my_key")},
		{InstanceId: awssdk.String("inst_5"), SubnetId: nil, VpcId: nil, KeyName: awssdk.String("unexisting_key")}, // terminated instance (no vpc, subnet ids)
		{
			InstanceId:         awssdk.String("inst_6"),
			Tags:               []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("inst_6_name")}},
			InstanceType:       awssdk.String("t2.micro"),
			SubnetId:           awssdk.String("sub_3"),
			VpcId:              awssdk.String("vpc_2"),
			PublicIpAddress:    awssdk.String("1.2.3.4"),
			PrivateIpAddress:   awssdk.String("10.0.0.1"),
			ImageId:            awssdk.String("ami-1234"),
			LaunchTime:         awssdk.Time(now),
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

	vpcs := []*ec2.Vpc{
		{VpcId: awssdk.String("vpc_1")},
		{VpcId: awssdk.String("vpc_2")},
	}

	securityGroups := []*ec2.SecurityGroup{
		{
			GroupId:   awssdk.String("securitygroup_1"),
			GroupName: awssdk.String("my_securitygroup"),
			VpcId:     awssdk.String("vpc_1"),
			IpPermissions: []*ec2.IpPermission{
				{FromPort: awssdk.Int64(22), ToPort: awssdk.Int64(80), IpProtocol: awssdk.String("tcp"), UserIdGroupPairs: []*ec2.UserIdGroupPair{{GroupId: awssdk.String("group_1")}, {GroupId: awssdk.String("group_2")}}},
			},
			IpPermissionsEgress: []*ec2.IpPermission{
				{FromPort: awssdk.Int64(0), ToPort: awssdk.Int64(65535), IpProtocol: awssdk.String("tcp"), IpRanges: []*ec2.IpRange{{CidrIp: awssdk.String("10.20.0.0/16")}}},
			},
		},
		{GroupId: awssdk.String("securitygroup_2"), VpcId: awssdk.String("vpc_1")},
	}

	subnets := []*ec2.Subnet{
		{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
		{SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
		{SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{SubnetId: awssdk.String("sub_4"), VpcId: nil}, // edge case subnet with no vpc id
	}

	keypairs := []*ec2.KeyPairInfo{
		{KeyName: awssdk.String("my_key")},
	}

	igws := []*ec2.InternetGateway{
		{InternetGatewayId: awssdk.String("igw_1"), Attachments: []*ec2.InternetGatewayAttachment{{VpcId: awssdk.String("vpc_2")}}},
	}

	natgws := []*ec2.NatGateway{
		{NatGatewayId: awssdk.String("natgw_1"), VpcId: awssdk.String("vpc_1"), SubnetId: awssdk.String("sub_1")},
	}

	routeTables := []*ec2.RouteTable{
		{
			RouteTableId: awssdk.String("rt_1"),
			VpcId:        awssdk.String("vpc_1"),
			Associations: []*ec2.RouteTableAssociation{
				{RouteTableId: awssdk.String("rt_1"), SubnetId: awssdk.String("sub_1"), RouteTableAssociationId: awssdk.String("assoc_1")},
				{RouteTableId: awssdk.String("rt_1"), SubnetId: awssdk.String("sub_2"), RouteTableAssociationId: awssdk.String("assoc_2"), Main: awssdk.Bool(true)},
			},
		},
	}

	images := []*ec2.Image{
		{ImageId: awssdk.String("img_1")},
		{ImageId: awssdk.String("img_2"), Name: awssdk.String("img_2_name"), Architecture: awssdk.String("img_2_arch"), Hypervisor: awssdk.String("img_2_hyper"), CreationDate: awssdk.String("2010-04-01T12:05:01.000Z")},
	}

	networkInterfaces := []*ec2.NetworkInterface{
		{
			Association:        &ec2.NetworkInterfaceAssociation{PublicIp: awssdk.String("1.2.3.4"), PublicDnsName: awssdk.String("my.ip.dns.name")},
			Attachment:         &ec2.NetworkInterfaceAttachment{AttachmentId: awssdk.String("eni-attach-12345"), InstanceId: awssdk.String("inst_1"), InstanceOwnerId: awssdk.String("12345678")},
			AvailabilityZone:   awssdk.String("us-west-1b"),
			Description:        awssdk.String("my network interface description"),
			Groups:             []*ec2.GroupIdentifier{{GroupId: awssdk.String("securitygroup_1")}, {GroupId: awssdk.String("securitygroup_2")}},
			InterfaceType:      awssdk.String("type"),
			Ipv6Addresses:      []*ec2.NetworkInterfaceIpv6Address{{Ipv6Address: awssdk.String("ab:cd:ef::")}, {Ipv6Address: awssdk.String("cd:ef:ab::")}},
			MacAddress:         awssdk.String("01:23:34:56:78:9a"),
			NetworkInterfaceId: awssdk.String("eni-1"),
			OwnerId:            awssdk.String("12345678"),
			PrivateDnsName:     awssdk.String("my.private.dns.name"),
			PrivateIpAddress:   awssdk.String("10.10.20.12"),
			Status:             awssdk.String("in-use"),
			SubnetId:           awssdk.String("sub_1"),
			VpcId:              awssdk.String("vpc_1"),
		},
		{
			NetworkInterfaceId: awssdk.String("eni-2"),
			SubnetId:           awssdk.String("sub_3"),
			VpcId:              awssdk.String("vpc_2")},
	}

	availabilityZones := []*ec2.AvailabilityZone{
		{ZoneName: awssdk.String("us-west-1a"), State: awssdk.String("available"), RegionName: awssdk.String("us-west-1"), Messages: []*ec2.AvailabilityZoneMessage{{Message: awssdk.String("msg 1")}, {Message: awssdk.String("msg 2")}}},
		{ZoneName: awssdk.String("us-west-1b")},
	}
	//ELBV2
	lbPages := []*elbv2.LoadBalancer{
		{LoadBalancerArn: awssdk.String("lb_1"), LoadBalancerName: awssdk.String("my_loadbalancer"), VpcId: awssdk.String("vpc_1")},
		{LoadBalancerArn: awssdk.String("lb_2"), VpcId: awssdk.String("vpc_2")},
		{LoadBalancerArn: awssdk.String("lb_3"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*string{awssdk.String("securitygroup_1"), awssdk.String("securitygroup_2")}},
	}
	//ELB
	classicLbPages := []*elb.LoadBalancerDescription{
		{LoadBalancerName: awssdk.String("my_classic_loadbalancer_1"), VPCId: awssdk.String("vpc_1"), ListenerDescriptions: []*elb.ListenerDescription{{Listener: &elb.Listener{LoadBalancerPort: awssdk.Int64(443), Protocol: awssdk.String("HTTPS"), InstancePort: awssdk.Int64(8080), InstanceProtocol: awssdk.String("HTTP")}}}},
		{LoadBalancerName: awssdk.String("my_classic_loadbalancer_2"), VPCId: awssdk.String("vpc_2")},
		{LoadBalancerName: awssdk.String("my_classic_loadbalancer_3"), VPCId: awssdk.String("vpc_1"), SecurityGroups: []*string{awssdk.String("securitygroup_1"), awssdk.String("securitygroup_2")}},
	}

	targetGroups := []*elbv2.TargetGroup{
		{TargetGroupArn: awssdk.String("tg_1"), VpcId: awssdk.String("vpc_1"), LoadBalancerArns: []*string{awssdk.String("lb_1"), awssdk.String("lb_3")}},
		{TargetGroupArn: awssdk.String("tg_2"), VpcId: awssdk.String("vpc_2"), LoadBalancerArns: []*string{awssdk.String("lb_2")}},
	}
	listeners := []*elbv2.Listener{
		{ListenerArn: awssdk.String("list_1"), LoadBalancerArn: awssdk.String("lb_1")}, {ListenerArn: awssdk.String("list_1.2"), LoadBalancerArn: awssdk.String("lb_1")},
		{ListenerArn: awssdk.String("list_2"), LoadBalancerArn: awssdk.String("lb_2")},
		{ListenerArn: awssdk.String("list_3"), LoadBalancerArn: awssdk.String("lb_3")},
	}
	targetHealths := map[string][]*elbv2.TargetHealthDescription{
		"tg_1": {{HealthCheckPort: awssdk.String("80"), Target: &elbv2.TargetDescription{Id: awssdk.String("inst_1"), Port: awssdk.Int64(443)}}},
		"tg_2": {{Target: &elbv2.TargetDescription{Id: awssdk.String("inst_2"), Port: awssdk.Int64(80)}}, {Target: &elbv2.TargetDescription{Id: awssdk.String("inst_3"), Port: awssdk.Int64(80)}}},
	}

	//Autoscaling
	launchConfigs := []*autoscaling.LaunchConfiguration{
		{LaunchConfigurationARN: awssdk.String("launchconfig_arn"), LaunchConfigurationName: awssdk.String("launchconfig_name"), KeyName: awssdk.String("my_key")},
	}
	scalingGroups := []*autoscaling.Group{
		{AutoScalingGroupARN: awssdk.String("asg_arn_1"), AutoScalingGroupName: awssdk.String("asg_name_1"), Instances: []*autoscaling.Instance{{InstanceId: awssdk.String("inst_1")}, {InstanceId: awssdk.String("inst_3")}}, VPCZoneIdentifier: awssdk.String("sub_1,sub_2"), LaunchConfigurationName: awssdk.String("launchconfig_name")},
		{AutoScalingGroupARN: awssdk.String("asg_arn_2"), AutoScalingGroupName: awssdk.String("asg_name_2"), LaunchConfigurationName: awssdk.String("launchconfig_name"), TargetGroupARNs: []*string{awssdk.String("tg_1"), awssdk.String("tg_2")}},
	}

	//ECR
	repositories := []*ecr.Repository{
		{CreatedAt: awssdk.Time(now), RegistryId: awssdk.String("account_id"), RepositoryArn: awssdk.String("repo_1"), RepositoryName: awssdk.String("repo_name_1"), RepositoryUri: awssdk.String("http://my.repository.url")},
		{RepositoryArn: awssdk.String("repo_2")},
		{RepositoryArn: awssdk.String("repo_3")},
	}

	//ECS
	clusterNames := []*string{awssdk.String("clust_1"), awssdk.String("clust_2"), awssdk.String("clust_3")}
	clusters := []*ecs.Cluster{
		{ActiveServicesCount: awssdk.Int64(3), ClusterArn: awssdk.String("clust_1"), ClusterName: awssdk.String("my_cust_1"), PendingTasksCount: awssdk.Int64(1), RegisteredContainerInstancesCount: awssdk.Int64(3), RunningTasksCount: awssdk.Int64(2), Status: awssdk.String("ACTIVE")},
		{ClusterArn: awssdk.String("clust_2")},
		{ClusterArn: awssdk.String("clust_3"), ClusterName: awssdk.String("my_cust_3")},
	}
	defNames := []*string{awssdk.String("cs_1:1"), awssdk.String("cs_2:1"), awssdk.String("cs_2:2"), awssdk.String("cs_3:1")}
	tasksDef := []*ecs.TaskDefinition{
		{
			ContainerDefinitions: []*ecs.ContainerDefinition{
				{Name: awssdk.String("cont_name_1"), Image: awssdk.String("image_1")},
				{Name: awssdk.String("cont_name_2"), Image: awssdk.String("image_2")},
				{Name: awssdk.String("cont_name_3"), Image: awssdk.String("image_3")},
			},
			Family:            awssdk.String("cs_1"),
			Revision:          awssdk.Int64(1),
			Status:            awssdk.String("ENABLED"),
			TaskDefinitionArn: awssdk.String("cs_1:1"),
			TaskRoleArn:       awssdk.String("role:arn"),
		},
		{
			ContainerDefinitions: []*ecs.ContainerDefinition{},
			Family:               awssdk.String("cs_2"),
			Revision:             awssdk.Int64(1),
			TaskDefinitionArn:    awssdk.String("cs_2:1"),
		},
		{
			ContainerDefinitions: []*ecs.ContainerDefinition{},
			Family:               awssdk.String("cs_2"),
			Revision:             awssdk.Int64(2),
			TaskDefinitionArn:    awssdk.String("cs_2:2"),
		},
		{
			ContainerDefinitions: []*ecs.ContainerDefinition{},
			Family:               awssdk.String("cs_3"),
			TaskDefinitionArn:    awssdk.String("cs_3:1"),
			Status:               awssdk.String("ACTIVE"),
		},
	}
	tasksNames := map[string][]*string{
		"clust_1": {awssdk.String("task_1")},
		"clust_2": {awssdk.String("task_2"), awssdk.String("task_3")},
	}
	tasks := map[string][]*ecs.Task{
		"clust_1": {
			{
				ClusterArn:           awssdk.String("clust_1"),
				ContainerInstanceArn: awssdk.String("cont_inst_1"),
				LastStatus:           awssdk.String("running"),
				Containers: []*ecs.Container{
					{
						ContainerArn: awssdk.String("container_1"),
						ExitCode:     awssdk.Int64(-1),
						LastStatus:   awssdk.String("running"),
						Name:         awssdk.String("my_container_1"),
						Reason:       awssdk.String("no reason"),
					}, {
						ContainerArn: awssdk.String("container_2"),
						Name:         awssdk.String("my_container_2"),
					}, {
						ContainerArn: awssdk.String("container_3"),
					},
				},
				Group:             awssdk.String("service:container-service-1"),
				CreatedAt:         awssdk.Time(now.Add(-2 * time.Hour)),
				StartedAt:         awssdk.Time(now.Add(-1 * time.Hour)),
				StoppedAt:         awssdk.Time(now),
				TaskArn:           awssdk.String("task_1"),
				TaskDefinitionArn: awssdk.String("cs_2:1"),
			},
		},
		"clust_2": {
			{
				ClusterArn:           awssdk.String("clust_2"),
				ContainerInstanceArn: awssdk.String("cont_inst_2"),
				LastStatus:           awssdk.String("stopped"),
				Containers: []*ecs.Container{
					{
						ContainerArn: awssdk.String("container_4"),
						ExitCode:     awssdk.Int64(0),
						LastStatus:   awssdk.String("stopped"),
						Name:         awssdk.String("my_container_4"),
					},
				},
				Group:             awssdk.String("service:container-service-2"),
				TaskArn:           awssdk.String("task_2"),
				TaskDefinitionArn: awssdk.String("cs_2:2"),
			},
			{
				ClusterArn:           awssdk.String("clust_2"),
				ContainerInstanceArn: awssdk.String("cont_inst_3"),
				Group:                awssdk.String("family:cs_1"),
				Containers: []*ecs.Container{
					{
						ContainerArn: awssdk.String("container_5"),
						Name:         awssdk.String("my_container_5"),
					},
				},
				LastStatus:        awssdk.String("running"),
				TaskArn:           awssdk.String("task_3"),
				TaskDefinitionArn: awssdk.String("cs_1:1"),
			},
		},
	}

	containerInstancesNames := map[string][]*string{
		"clust_1": {awssdk.String("cont_inst_1"), awssdk.String("cont_inst_2")},
		"clust_2": {awssdk.String("cont_inst_3")},
	}
	containerInstances := map[string][]*ecs.ContainerInstance{
		"clust_1": {
			{
				AgentConnected:    awssdk.Bool(true),
				AgentUpdateStatus: awssdk.String("AgentRunning"),
				Attributes: []*ecs.Attribute{
					{
						Name:  awssdk.String("attr_1"),
						Value: awssdk.String("val1"),
					},
					{
						Name:  awssdk.String("attr_2"),
						Value: awssdk.String("val2"),
					},
				},
				ContainerInstanceArn: awssdk.String("cont_inst_1"),
				Ec2InstanceId:        awssdk.String("inst_2"),
				PendingTasksCount:    awssdk.Int64(4),
				RegisteredAt:         awssdk.Time(now.Add(-2 * time.Hour)),
				RunningTasksCount:    awssdk.Int64(2),
				Status:               awssdk.String("ACTIVE"),
				Version:              awssdk.Int64(2),
				VersionInfo:          &ecs.VersionInfo{AgentVersion: awssdk.String("0.0.5"), DockerVersion: awssdk.String("v1.0.12")},
			},
			{
				ContainerInstanceArn: awssdk.String("cont_inst_2"),
				Ec2InstanceId:        awssdk.String("inst_3"),
			},
		},
		"clust_2": {
			{
				ContainerInstanceArn: awssdk.String("cont_inst_3"),
				Ec2InstanceId:        awssdk.String("inst_1"),
			},
		},
	}

	//ACM
	certificates := []*acm.CertificateSummary{
		{CertificateArn: awssdk.String("arn:certif_1234"), DomainName: awssdk.String("domain-name.1")},
		{CertificateArn: awssdk.String("arn:certif_2345"), DomainName: awssdk.String("domain-name.2")},
		{CertificateArn: awssdk.String("arn:certif_3456"), DomainName: awssdk.String("domain-name.3")},
	}

	mock := &mockEc2{vpcs: vpcs, securitygroups: securityGroups, subnets: subnets, instances: instances, keypairinfos: keypairs, internetgateways: igws, routetables: routeTables, images: images, availabilityzones: availabilityZones, natgateways: natgws, networkinterfaces: networkInterfaces}
	mockLb := &mockElbv2{loadbalancers: lbPages, targetgroups: targetGroups, listeners: listeners, targethealthdescriptions: targetHealths}
	mockClassicLb := &mockElb{loadbalancerdescriptions: classicLbPages}
	mockEcr := &mockEcr{repositorys: repositories}
	mockEcs := &mockEcs{clusterNames: clusterNames, clusters: clusters, taskdefinitionNames: defNames, taskdefinitions: tasksDef, tasksNames: tasksNames, tasks: tasks, containerinstancesNames: containerInstancesNames, containerinstances: containerInstances}
	mockRds := &mockRds{}
	mockAcm := &mockAcm{certificatesummarys: certificates}
	mockAutoscaling := &mockAutoscaling{launchconfigurations: launchConfigs, groups: scalingGroups}
	InfraService = &Infra{
		EC2API:         mock,
		ECRAPI:         mockEcr,
		ECSAPI:         mockEcs,
		ELBAPI:         mockClassicLb,
		ELBV2API:       mockLb,
		RDSAPI:         mockRds,
		ACMAPI:         mockAcm,
		AutoScalingAPI: mockAutoscaling,
		region:         "eu-west-1",
		fetcher:        fetch.NewFetcher(awsfetch.BuildInfraFetchFuncs(awsfetch.NewConfig(mock, mockEcr, mockEcs, mockClassicLb, mockLb, mockRds, mockAutoscaling, mockAcm))),
	}
	g, err := InfraService.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resources, err := g.Find(cloud.NewQuery("region", "instance", "vpc", "securitygroup", "subnet", "keypair", "internetgateway", cloud.NatGateway, "routetable", "classicloadbalancer", "loadbalancer", "targetgroup", "listener", "launchconfiguration", "scalinggroup", "image", "availabilityzone", "repository", cloud.ContainerCluster, cloud.ContainerTask, cloud.Container, cloud.ContainerInstance, cloud.NetworkInterface, cloud.Certificate))
	if err != nil {
		t.Fatal(err)
	}

	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties()[p.SecurityGroups].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.Vpcs].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.Messages].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.ContainersImages].([]*graph.KeyValue); ok {
			sort.Slice(p, func(i, j int) bool {
				if p[i].KeyName == p[j].KeyName {
					return p[i].Value < p[j].Value
				}
				return p[i].KeyName < p[j].KeyName
			})
		}
		if p, ok := res.Properties()[p.Attributes].([]*graph.KeyValue); ok {
			sort.Slice(p, func(i, j int) bool {
				if p[i].KeyName == p[j].KeyName {
					return p[i].Value < p[j].Value
				}
				return p[i].KeyName < p[j].KeyName
			})
		}
		if p, ok := res.Properties()[p.Associations].([]*graph.KeyValue); ok {
			sort.Slice(p, func(i, j int) bool {
				if p[i].KeyName == p[j].KeyName {
					return p[i].Value < p[j].Value
				}
				return p[i].KeyName < p[j].KeyName
			})
		}
		if p, ok := res.Properties()[p.InboundRules].([]*graph.FirewallRule); ok {
			for _, r := range p {
				sort.Strings(r.Sources)
			}
		}
		if p, ok := res.Properties()[p.IPv6Addresses].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]cloud.Resource{
		"eu-west-1": resourcetest.Region("eu-west-1").Build(),
		"inst_1":    resourcetest.Instance("inst_1").Prop(p.Subnet, "sub_1").Prop(p.Vpc, "vpc_1").Prop(p.Name, "instance1-name").Prop(p.Tags, []string{"Name=instance1-name"}).Build(),
		"inst_2":    resourcetest.Instance("inst_2").Prop(p.Subnet, "sub_2").Prop(p.Vpc, "vpc_1").Prop(p.SecurityGroups, []string{"securitygroup_1"}).Build(),
		"inst_3":    resourcetest.Instance("inst_3").Prop(p.Subnet, "sub_3").Prop(p.Vpc, "vpc_2").Build(),
		"inst_4":    resourcetest.Instance("inst_4").Prop(p.Subnet, "sub_3").Prop(p.Vpc, "vpc_2").Prop(p.SecurityGroups, []string{"securitygroup_1", "securitygroup_2"}).Prop(p.KeyPair, "my_key").Build(),
		"inst_5":    resourcetest.Instance("inst_5").Prop(p.KeyPair, "unexisting_key").Build(),
		"inst_6": resourcetest.Instance("inst_6").Prop(p.Name, "inst_6_name").Prop(p.Tags, []string{"Name=inst_6_name"}).Prop(p.Type, "t2.micro").Prop(p.Subnet, "sub_3").Prop(p.Vpc, "vpc_2").Prop(p.PublicIP, "1.2.3.4").Prop(p.PrivateIP, "10.0.0.1").
			Prop(p.Image, "ami-1234").Prop(p.Launched, now).Prop(p.State, "running").Prop(p.KeyPair, "my_key").Prop(p.SecurityGroups, []string{"securitygroup_1"}).Prop(p.Affinity, "inst_affinity").
			Prop(p.AvailabilityZone, "inst_az").Prop(p.PlacementGroup, "inst_group").Prop(p.Host, "inst_host").Prop(p.Architecture, "x86").Prop(p.Hypervisor, "xen").Prop(p.Profile, "arn:instance:profile").
			Prop(p.Lifecycle, "lifecycle").Prop(p.NetworkInterfaces, []string{"my-network-interface"}).Prop(p.PublicDNS, "my-instance.dns").Prop(p.RootDevice, "/dev/xvda").Prop(p.RootDeviceType, "ebs").Build(),
		"vpc_1": resourcetest.VPC("vpc_1").Build(),
		"vpc_2": resourcetest.VPC("vpc_2").Build(),
		"securitygroup_1": resourcetest.SecurityGroup("securitygroup_1").Prop(p.Name, "my_securitygroup").Prop(p.Vpc, "vpc_1").
			Prop(p.InboundRules, []*graph.FirewallRule{{PortRange: graph.PortRange{FromPort: 22, ToPort: 80, Any: false}, Protocol: "tcp", Sources: []string{"group_1", "group_2"}}}).
			Prop(p.OutboundRules, []*graph.FirewallRule{{PortRange: graph.PortRange{FromPort: 0, ToPort: 65535, Any: false}, Protocol: "tcp", IPRanges: []*net.IPNet{{IP: net.IP{0xa, 0x14, 0x0, 0x0}, Mask: net.CIDRMask(16, 32)}}}}).Build(),
		"securitygroup_2": resourcetest.SecurityGroup("securitygroup_2").Prop(p.Vpc, "vpc_1").Build(),
		"sub_1":           resourcetest.Subnet("sub_1").Prop(p.Vpc, "vpc_1").Build(),
		"sub_2":           resourcetest.Subnet("sub_2").Prop(p.Vpc, "vpc_1").Build(),
		"sub_3":           resourcetest.Subnet("sub_3").Prop(p.Vpc, "vpc_2").Build(),
		"sub_4":           resourcetest.Subnet("sub_4").Build(),
		"us-west-1a":      resourcetest.AvailabilityZone("us-west-1a").Prop(p.Name, "us-west-1a").Prop(p.State, "available").Prop(p.Region, "us-west-1").Prop(p.Messages, []string{"msg 1", "msg 2"}).Build(),
		"us-west-1b":      resourcetest.AvailabilityZone("us-west-1b").Prop(p.Name, "us-west-1b").Build(),
		"my_key":          resourcetest.KeyPair("my_key").Build(),
		"igw_1":           resourcetest.InternetGw("igw_1").Prop(p.Vpcs, []string{"vpc_2"}).Build(),
		"natgw_1":         resourcetest.NatGw("natgw_1").Prop(p.Vpc, "vpc_1").Prop(p.Subnet, "sub_1").Build(),
		"rt_1":            resourcetest.RouteTable("rt_1").Prop(p.Vpc, "vpc_1").Prop(p.Default, true).Prop(p.Associations, []*graph.KeyValue{{KeyName: "assoc_1", Value: "sub_1"}, {KeyName: "assoc_2", Value: "sub_2"}}).Build(),
		"lb_1":            resourcetest.LoadBalancer("lb_1").Prop(p.Arn, "lb_1").Prop(p.Name, "my_loadbalancer").Prop(p.Vpc, "vpc_1").Build(),
		"lb_2":            resourcetest.LoadBalancer("lb_2").Prop(p.Arn, "lb_2").Prop(p.Vpc, "vpc_2").Build(),
		"lb_3":            resourcetest.LoadBalancer("lb_3").Prop(p.Arn, "lb_3").Prop(p.Vpc, "vpc_1").Build(),
		"my_classic_loadbalancer_1": resourcetest.ClassicLoadBalancer("my_classic_loadbalancer_1").Prop(p.Name, "my_classic_loadbalancer_1").Prop(p.Vpc, "vpc_1").Prop(p.Ports, []string{"HTTPS:443:HTTP:8080"}).Build(),
		"my_classic_loadbalancer_2": resourcetest.ClassicLoadBalancer("my_classic_loadbalancer_2").Prop(p.Name, "my_classic_loadbalancer_2").Prop(p.Vpc, "vpc_2").Build(),
		"my_classic_loadbalancer_3": resourcetest.ClassicLoadBalancer("my_classic_loadbalancer_3").Prop(p.Name, "my_classic_loadbalancer_3").Prop(p.Vpc, "vpc_1").Build(),
		"tg_1":             resourcetest.TargetGroup("tg_1").Prop(p.Arn, "tg_1").Prop(p.Vpc, "vpc_1").Build(),
		"tg_2":             resourcetest.TargetGroup("tg_2").Prop(p.Arn, "tg_2").Prop(p.Vpc, "vpc_2").Build(),
		"list_1":           resourcetest.Listener("list_1").Prop(p.Arn, "list_1").Prop(p.LoadBalancer, "lb_1").Build(),
		"list_1.2":         resourcetest.Listener("list_1.2").Prop(p.Arn, "list_1.2").Prop(p.LoadBalancer, "lb_1").Build(),
		"list_2":           resourcetest.Listener("list_2").Prop(p.Arn, "list_2").Prop(p.LoadBalancer, "lb_2").Build(),
		"list_3":           resourcetest.Listener("list_3").Prop(p.Arn, "list_3").Prop(p.LoadBalancer, "lb_3").Build(),
		"launchconfig_arn": resourcetest.LaunchConfig("launchconfig_arn").Prop(p.Arn, "launchconfig_arn").Prop(p.Name, "launchconfig_name").Prop(p.KeyPair, "my_key").Build(),
		"asg_arn_1":        resourcetest.ScalingGroup("asg_arn_1").Prop(p.Arn, "asg_arn_1").Prop(p.Name, "asg_name_1").Prop(p.LaunchConfigurationName, "launchconfig_name").Build(),
		"asg_arn_2":        resourcetest.ScalingGroup("asg_arn_2").Prop(p.Arn, "asg_arn_2").Prop(p.Name, "asg_name_2").Prop(p.LaunchConfigurationName, "launchconfig_name").Build(),
		"img_1":            resourcetest.Image("img_1").Build(),
		"img_2":            resourcetest.Image("img_2").Prop(p.Name, "img_2_name").Prop(p.Architecture, "img_2_arch").Prop(p.Hypervisor, "img_2_hyper").Prop(p.Created, time.Unix(1270123501, 0).UTC()).Build(),
		"repo_1":           resourcetest.Repository("repo_1").Prop(p.Created, now).Prop(p.Arn, "repo_1").Prop(p.Account, "account_id").Prop(p.Name, "repo_name_1").Prop(p.URI, "http://my.repository.url").Build(),
		"repo_2":           resourcetest.Repository("repo_2").Prop(p.Arn, "repo_2").Build(),
		"repo_3":           resourcetest.Repository("repo_3").Prop(p.Arn, "repo_3").Build(),
		"clust_1":          resourcetest.ContainerCluster("clust_1").Prop(p.Arn, "clust_1").Prop(p.Name, "my_cust_1").Prop(p.PendingTasksCount, 1).Prop(p.ActiveServicesCount, 3).Prop(p.RegisteredContainerInstancesCount, 3).Prop(p.RunningTasksCount, 2).Prop(p.State, "ACTIVE").Build(),
		"clust_2":          resourcetest.ContainerCluster("clust_2").Prop(p.Arn, "clust_2").Build(),
		"clust_3":          resourcetest.ContainerCluster("clust_3").Prop(p.Arn, "clust_3").Prop(p.Name, "my_cust_3").Build(),
		"cs_1:1": resourcetest.ContainerTask("cs_1:1").Prop(p.Arn, "cs_1:1").Prop(p.ContainersImages, []*graph.KeyValue{{"cont_name_1", "image_1"}, {"cont_name_2", "image_2"}, {"cont_name_3", "image_3"}}).Prop(p.Name, "cs_1").Prop(p.Version, "1").
			Prop(p.State, "1 task running").Prop(p.Role, "role:arn").Prop(p.Deployments, []*graph.KeyValue{{"clust_2", "cs_1 (running task)"}}).Build(),
		"cs_2:1": resourcetest.ContainerTask("cs_2:1").Prop(p.Arn, "cs_2:1").Prop(p.Name, "cs_2").Prop(p.State, "1 service running").Prop(p.Version, "1").Prop(p.Deployments, []*graph.KeyValue{{"clust_1", "container-service-1 (running service)"}}).Build(),
		"cs_2:2": resourcetest.ContainerTask("cs_2:2").Prop(p.Arn, "cs_2:2").Prop(p.Name, "cs_2").Prop(p.State, "1 service stopped").Prop(p.Version, "2").Prop(p.Deployments, []*graph.KeyValue{{"clust_2", "container-service-2 (stopped service)"}}).Build(),
		"cs_3:1": resourcetest.ContainerTask("cs_3:1").Prop(p.Arn, "cs_3:1").Prop(p.Name, "cs_3").Prop(p.State, "ready").Build(),
		"container_1": resourcetest.Container("container_1").Prop(p.Arn, "container_1").Prop(p.ExitCode, -1).Prop(p.State, "running").Prop(p.Name, "my_container_1").Prop(p.StateMessage, "no reason").Prop(p.Cluster, "clust_1").
			Prop(p.ContainerInstance, "cont_inst_1").Prop(p.Created, now.Add(-2*time.Hour)).Prop(p.Launched, now.Add(-1*time.Hour)).Prop(p.Stopped, now).Prop(p.DeploymentName, "service:container-service-1").Prop(p.ContainerTask, "cs_2:1").Build(),
		"container_2": resourcetest.Container("container_2").Prop(p.Arn, "container_2").Prop(p.Name, "my_container_2").Prop(p.Cluster, "clust_1").Prop(p.ContainerInstance, "cont_inst_1").Prop(p.Created, now.Add(-2*time.Hour)).Prop(p.Launched, now.Add(-1*time.Hour)).Prop(p.Stopped, now).Prop(p.DeploymentName, "service:container-service-1").Prop(p.ContainerTask, "cs_2:1").Build(),
		"container_3": resourcetest.Container("container_3").Prop(p.Arn, "container_3").Prop(p.Cluster, "clust_1").Prop(p.ContainerInstance, "cont_inst_1").Prop(p.Created, now.Add(-2*time.Hour)).Prop(p.Launched, now.Add(-1*time.Hour)).Prop(p.Stopped, now).Prop(p.ContainerTask, "cs_2:1").Prop(p.DeploymentName, "service:container-service-1").Build(),
		"container_4": resourcetest.Container("container_4").Prop(p.Arn, "container_4").Prop(p.Name, "my_container_4").Prop(p.ExitCode, 0).Prop(p.State, "stopped").Prop(p.Cluster, "clust_2").Prop(p.ContainerInstance, "cont_inst_2").Prop(p.ContainerTask, "cs_2:2").Prop(p.DeploymentName, "service:container-service-2").Build(),
		"container_5": resourcetest.Container("container_5").Prop(p.Arn, "container_5").Prop(p.Name, "my_container_5").Prop(p.Cluster, "clust_2").Prop(p.ContainerInstance, "cont_inst_3").Prop(p.ContainerTask, "cs_1:1").Prop(p.DeploymentName, "family:cs_1").Build(),
		"cont_inst_1": resourcetest.ContainerInstance("cont_inst_1").Prop(p.Arn, "cont_inst_1").Prop(p.AgentConnected, true).Prop(p.AgentState, "AgentRunning").Prop(p.Attributes, []*graph.KeyValue{{"attr_1", "val1"}, {"attr_2", "val2"}}).
			Prop(p.Instance, "inst_2").Prop(p.PendingTasksCount, 4).Prop(p.Created, now.Add(-2*time.Hour)).Prop(p.RunningTasksCount, 2).Prop(p.State, "ACTIVE").Prop(p.Version, "2").Prop(p.AgentVersion, "0.0.5").Prop(p.DockerVersion, "v1.0.12").Prop(p.Cluster, "clust_1").Build(),
		"cont_inst_2": resourcetest.ContainerInstance("cont_inst_2").Prop(p.Arn, "cont_inst_2").Prop(p.Instance, "inst_3").Prop(p.Cluster, "clust_1").Build(),
		"cont_inst_3": resourcetest.ContainerInstance("cont_inst_3").Prop(p.Arn, "cont_inst_3").Prop(p.Instance, "inst_1").Prop(p.Cluster, "clust_2").Build(),
		"eni-1": resourcetest.NetworkInterface("eni-1").Prop(p.PublicIP, "1.2.3.4").Prop(p.PublicDNS, "my.ip.dns.name").Prop(p.Attachment, "eni-attach-12345").Prop(p.Instance, "inst_1").Prop(p.InstanceOwner, "12345678").
			Prop(p.AvailabilityZone, "us-west-1b").Prop(p.Description, "my network interface description").Prop(p.SecurityGroups, []string{"securitygroup_1", "securitygroup_2"}).Prop(p.Type, "type").Prop(p.IPv6Addresses, []string{"ab:cd:ef::", "cd:ef:ab::"}).Prop(p.MACAddress, "01:23:34:56:78:9a").
			Prop(p.Owner, "12345678").Prop(p.PrivateDNS, "my.private.dns.name").Prop(p.PrivateIP, "10.10.20.12").Prop(p.State, "in-use").Prop(p.Subnet, "sub_1").Prop(p.Vpc, "vpc_1").Build(),
		"eni-2":           resourcetest.NetworkInterface("eni-2").Prop(p.Subnet, "sub_3").Prop(p.Vpc, "vpc_2").Build(),
		"arn:certif_1234": resourcetest.Certificate("arn:certif_1234").Prop(p.Arn, "arn:certif_1234").Prop(p.Name, "domain-name.1").Build(),
		"arn:certif_2345": resourcetest.Certificate("arn:certif_2345").Prop(p.Arn, "arn:certif_2345").Prop(p.Name, "domain-name.2").Build(),
		"arn:certif_3456": resourcetest.Certificate("arn:certif_3456").Prop(p.Arn, "arn:certif_3456").Prop(p.Name, "domain-name.3").Build(),
	}

	expectedChildren := map[string][]string{
		"eu-west-1": {"arn:certif_1234", "arn:certif_2345", "arn:certif_3456", "asg_arn_1", "asg_arn_2", "clust_1", "clust_2", "clust_3", "cs_1:1", "cs_2:1", "cs_2:2", "cs_3:1", "igw_1", "img_1", "img_2", "launchconfig_arn", "my_key", "natgw_1", "repo_1", "repo_2", "repo_3", "us-west-1a", "us-west-1b", "vpc_1", "vpc_2"},
		"lb_1":      {"list_1", "list_1.2"},
		"lb_2":      {"list_2"},
		"lb_3":      {"list_3"},
		"sub_1":     {"eni-1", "inst_1"},
		"sub_2":     {"inst_2"},
		"sub_3":     {"eni-2", "inst_3", "inst_4", "inst_6"},
		"vpc_1":     {"lb_1", "lb_3", "my_classic_loadbalancer_1", "my_classic_loadbalancer_3", "natgw_1", "rt_1", "securitygroup_1", "securitygroup_2", "sub_1", "sub_2", "tg_1"},
		"vpc_2":     {"lb_2", "my_classic_loadbalancer_2", "sub_3", "tg_2"},
		"clust_1":   {"cont_inst_1", "cont_inst_2", "container_1", "container_2", "container_3"},
		"clust_2":   {"cont_inst_3", "container_4", "container_5"},
	}

	expectedAppliedOn := map[string][]string{
		"igw_1":           {"vpc_2"},
		"lb_1":            {"tg_1"},
		"lb_2":            {"tg_2"},
		"lb_3":            {"tg_1"},
		"my_key":          {"inst_4", "inst_6", "launchconfig_arn"},
		"natgw_1":         {"sub_1"},
		"rt_1":            {"sub_1", "sub_2"},
		"securitygroup_1": {"eni-1", "inst_2", "inst_4", "inst_6", "lb_3", "my_classic_loadbalancer_3"},
		"securitygroup_2": {"eni-1", "inst_4", "lb_3", "my_classic_loadbalancer_3"},
		"tg_1":            {"inst_1"},
		"tg_2":            {"inst_2", "inst_3"},
		"asg_arn_1":       {"inst_1", "inst_3", "sub_1", "sub_2"},
		"asg_arn_2":       {"tg_1", "tg_2"},
		"cs_1:1":          {"container_5"},
		"cs_2:1":          {"container_1", "container_2", "container_3"},
		"cs_2:2":          {"container_4"},
		"inst_1":          {"cont_inst_3"},
		"inst_2":          {"cont_inst_1"},
		"inst_3":          {"cont_inst_2"},
		"cont_inst_1":     {"container_1", "container_2", "container_3"},
		"cont_inst_2":     {"container_4"},
		"cont_inst_3":     {"container_5"},
		"eni-1":           {"inst_1"},
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildStorageRdfGraph(t *testing.T) {
	buckets := map[string][]*s3.Bucket{
		"us-west-1": {
			{Name: awssdk.String("bucket_us_1")},
			{Name: awssdk.String("bucket_us_2")},
			{Name: awssdk.String("bucket_us_3")},
		},
		"eu-west-1": {
			{Name: awssdk.String("bucket_eu_1")},
			{Name: awssdk.String("bucket_eu_2")},
		},
	}
	objects := map[string][]*s3.Object{
		"bucket_us_1": {
			{Key: awssdk.String("obj_1")},
			{Key: awssdk.String("obj_2")},
		},
		"bucket_us_2": {},
		"bucket_us_3": {
			{Key: awssdk.String("obj_3")},
		},
		"bucket_eu_1": {
			{Key: awssdk.String("obj_4")},
		},
		"bucket_eu_2": {
			{Key: awssdk.String("obj_5")},
			{Key: awssdk.String("obj_6")},
		},
	}
	bucketsACL := map[string][]*s3.Grant{
		"bucket_us_1": {
			{Permission: awssdk.String("Read"), Grantee: &s3.Grantee{ID: awssdk.String("usr_1")}},
		},
		"bucket_us_3": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_2")}},
		},
		"bucket_eu_1": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_2")}},
		},
		"bucket_eu_2": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_1")}},
		},
	}

	mocks3 := &mockS3{buckets: buckets, objects: objects, grants: bucketsACL}
	StorageService = mocks3
	storage := Storage{
		S3API:   mocks3,
		region:  "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildStorageFetchFuncs(awsfetch.NewConfig(mocks3))),
	}

	g, err := storage.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	resources, err := g.Find(cloud.NewQuery("region", "bucket"))
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]cloud.Resource{
		"eu-west-1":   resourcetest.Region("eu-west-1").Build(),
		"bucket_eu_1": resourcetest.Bucket("bucket_eu_1").Prop(p.Grants, []*graph.Grant{{Grantee: graph.Grantee{GranteeID: "usr_2"}, Permission: "Write"}}).Build(),
		"bucket_eu_2": resourcetest.Bucket("bucket_eu_2").Prop(p.Grants, []*graph.Grant{{Grantee: graph.Grantee{GranteeID: "usr_1"}, Permission: "Write"}}).Build(),
	}
	expectedChildren := map[string][]string{
		"eu-west-1":   {"bucket_eu_1", "bucket_eu_2"},
		"bucket_eu_1": {"obj_4"},
		"bucket_eu_2": {"obj_5", "obj_6"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildDnsRdfGraph(t *testing.T) {
	zonePages := []*route53.HostedZone{
		{Id: awssdk.String("/hostedzone/12345"), Name: awssdk.String("my.first.domain")},
		{Id: awssdk.String("/hostedzone/23456"), Name: awssdk.String("my.second.domain")},
		{Id: awssdk.String("/hostedzone/34567"), Name: awssdk.String("my.third.domain")},
	}
	recordPages := map[string][]*route53.ResourceRecordSet{
		"/hostedzone/12345": {
			{Type: awssdk.String("A"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain1.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("1.2.3.4")}, {Value: awssdk.String("2.3.4.5")}}},
			{Type: awssdk.String("A"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain2.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("3.4.5.6")}}},
			{Type: awssdk.String("CNAME"), TTL: awssdk.Int64(60), Name: awssdk.String("subdomain3.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("4.5.6.7")}}},
		},
		"/hostedzone/23456": {
			{Type: awssdk.String("A"), TTL: awssdk.Int64(30), Name: awssdk.String("subdomain1.my.second.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("5.6.7.8")}}},
			{Type: awssdk.String("CNAME"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain3.my.second.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("6.7.8.9")}}},
		},
	}
	mockRoute53 := &mockRoute53{hostedzones: zonePages, resourcerecordsets: recordPages}

	dns := Dns{
		Route53API: mockRoute53, region: "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildDnsFetchFuncs(awsfetch.NewConfig(mockRoute53))),
	}

	g, err := dns.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery("zone", "record"))
	if err != nil {
		t.Fatal(err)
	}
	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties()[p.Records].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]cloud.Resource{
		"/hostedzone/12345": resourcetest.Zone("/hostedzone/12345").Prop(p.Name, "my.first.domain").Build(),
		"/hostedzone/23456": resourcetest.Zone("/hostedzone/23456").Prop(p.Name, "my.second.domain").Build(),
		"/hostedzone/34567": resourcetest.Zone("/hostedzone/34567").Prop(p.Name, "my.third.domain").Build(),
		"awls-91fa0a45":     resourcetest.Record("awls-91fa0a45").Prop(p.Name, "subdomain1.my.first.domain").Prop(p.Zone, "my.first.domain").Prop(p.Type, "A").Prop(p.TTL, 10).Prop(p.Records, []string{"1.2.3.4", "2.3.4.5"}).Build(),
		"awls-920c0a46":     resourcetest.Record("awls-920c0a46").Prop(p.Name, "subdomain2.my.first.domain").Prop(p.Zone, "my.first.domain").Prop(p.Type, "A").Prop(p.TTL, 10).Prop(p.Records, []string{"3.4.5.6"}).Build(),
		"awls-be1e0b6a":     resourcetest.Record("awls-be1e0b6a").Prop(p.Name, "subdomain3.my.first.domain").Prop(p.Zone, "my.first.domain").Prop(p.Type, "CNAME").Prop(p.TTL, 60).Prop(p.Records, []string{"4.5.6.7"}).Build(),
		"awls-9c420a99":     resourcetest.Record("awls-9c420a99").Prop(p.Name, "subdomain1.my.second.domain").Prop(p.Zone, "my.second.domain").Prop(p.Type, "A").Prop(p.TTL, 30).Prop(p.Records, []string{"5.6.7.8"}).Build(),
		"awls-c9b80bbe":     resourcetest.Record("awls-c9b80bbe").Prop(p.Name, "subdomain3.my.second.domain").Prop(p.Zone, "my.second.domain").Prop(p.Type, "CNAME").Prop(p.TTL, 10).Prop(p.Records, []string{"6.7.8.9"}).Build(),
	}
	expectedChildren := map[string][]string{
		"/hostedzone/12345": {"awls-91fa0a45", "awls-920c0a46", "awls-be1e0b6a"},
		"/hostedzone/23456": {"awls-9c420a99", "awls-c9b80bbe"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildNotificationGraph(t *testing.T) {
	topics := []*sns.Topic{
		{TopicArn: awssdk.String("topic_arn_1")},
		{TopicArn: awssdk.String("topic_arn_2")},
		{TopicArn: awssdk.String("topic_arn_3")},
	}

	subscriptions := []*sns.Subscription{
		{Endpoint: awssdk.String("endpoint_1")},
		{Endpoint: awssdk.String("endpoint_2"), Owner: awssdk.String("subscr_owner"), Protocol: awssdk.String("subscr_prot"), SubscriptionArn: awssdk.String("subscr_arn"), TopicArn: awssdk.String("topic_arn_2")},
		{Endpoint: awssdk.String("endpoint_3"), TopicArn: awssdk.String("topic_arn_2")},
	}
	queues := []*string{awssdk.String("queue_1"), awssdk.String("queue_2"), awssdk.String("queue_3")}
	attributes := map[string]map[string]*string{
		"queue_2": {
			"ApproximateNumberOfMessages": awssdk.String("4"),
			"CreatedTimestamp":            awssdk.String("1494419259"),
			"LastModifiedTimestamp":       awssdk.String("1494332859"),
			"QueueArn":                    awssdk.String("queue_2_arn"),
			"DelaySeconds":                awssdk.String("15"),
		},
		"queue_3": {
			"ApproximateNumberOfMessages": awssdk.String("12"),
		},
	}

	sqs := &mockSqs{strings: queues, attributes: attributes}
	sns := &mockSns{subscriptions: subscriptions, topics: topics}

	service := Messaging{
		SNSAPI: sns, SQSAPI: sqs, region: "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildMessagingFetchFuncs(awsfetch.NewConfig(sqs, sns))),
	}

	g, err := service.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery("subscription", "topic"))
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]cloud.Resource{
		"endpoint_1":  resourcetest.Subscription("endpoint_1").Prop(p.Endpoint, "endpoint_1").Build(),
		"endpoint_2":  resourcetest.Subscription("endpoint_2").Prop(p.Endpoint, "endpoint_2").Prop(p.Owner, "subscr_owner").Prop(p.Protocol, "subscr_prot").Prop(p.Arn, "subscr_arn").Prop(p.Topic, "topic_arn_2").Build(),
		"endpoint_3":  resourcetest.Subscription("endpoint_3").Prop(p.Endpoint, "endpoint_3").Prop(p.Topic, "topic_arn_2").Build(),
		"topic_arn_1": resourcetest.Topic("topic_arn_1").Prop(p.Arn, "topic_arn_1").Build(),
		"topic_arn_2": resourcetest.Topic("topic_arn_2").Prop(p.Arn, "topic_arn_2").Build(),
		"topic_arn_3": resourcetest.Topic("topic_arn_3").Prop(p.Arn, "topic_arn_3").Build(),
	}
	expectedChildren := map[string][]string{
		"eu-west-1":   {"topic_arn_1", "topic_arn_2", "topic_arn_3"},
		"topic_arn_2": {"endpoint_2", "endpoint_3"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)

	resources, err = g.Find(cloud.NewQuery("queue"))
	if err != nil {
		t.Fatal(err)
	}

	expected = map[string]cloud.Resource{
		"queue_1": resourcetest.Queue("queue_1").Build(),
		"queue_2": resourcetest.Queue("queue_2").Prop(p.ApproximateMessageCount, 4).Prop(p.Created, time.Unix(1494419259, 0).UTC()).Prop(p.Modified, time.Unix(1494332859, 0).UTC()).Prop(p.Arn, "queue_2_arn").Prop(p.Delay, 15).Build(),
		"queue_3": resourcetest.Queue("queue_3").Prop(p.ApproximateMessageCount, 12).Build(),
	}
	expectedChildren = map[string][]string{}
	expectedAppliedOn = map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)

}

func TestBuildLambdaGraph(t *testing.T) {
	functions := []*lambda.FunctionConfiguration{
		{FunctionArn: awssdk.String("func_1_arn")},
		{
			FunctionArn:  awssdk.String("func_2_arn"),
			FunctionName: awssdk.String("func_2_name"),
			CodeSha256:   awssdk.String("abcdef123456789"),
			CodeSize:     awssdk.Int64(1234),
			Description:  awssdk.String("my function desc"),
			Handler:      awssdk.String("handl"),
			LastModified: awssdk.String("2006-01-02T15:04:05.000+0000"),
			MemorySize:   awssdk.Int64(1234),
			Role:         awssdk.String("role"),
			Runtime:      awssdk.String("runtime"),
			Timeout:      awssdk.Int64(60),
			Version:      awssdk.String("v2"),
		},
		{FunctionArn: awssdk.String("func_3_arn")},
	}

	mock := &mockLambda{functionconfigurations: functions}

	service := Lambda{
		LambdaAPI: mock, region: "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildLambdaFetchFuncs(awsfetch.NewConfig(mock))),
	}

	g, err := service.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery("function"))
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]cloud.Resource{
		"func_1_arn": resourcetest.Function("func_1_arn").Prop(p.Arn, "func_1_arn").Build(),
		"func_2_arn": resourcetest.Function("func_2_arn").Prop(p.Arn, "func_2_arn").Prop(p.Name, "func_2_name").Prop(p.Hash, "abcdef123456789").Prop(p.Size, 1234).
			Prop(p.Description, "my function desc").Prop(p.Handler, "handl").Prop(p.Modified, time.Unix(1136214245, 0).UTC()).Prop(p.Memory, 1234).Prop(p.Role, "role").
			Prop(p.Runtime, "runtime").Prop(p.Timeout, 60).Prop(p.Version, "v2").Build(),
		"func_3_arn": resourcetest.Function("func_3_arn").Prop(p.Arn, "func_3_arn").Build(),
	}

	expectedChildren := map[string][]string{
		"eu-west-1": {"func_1_arn", "func_2_arn", "func_3_arn"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildMonitoringGraph(t *testing.T) {
	now := time.Now().UTC()
	metrics := []*cloudwatch.Metric{
		{Namespace: awssdk.String("namespace_1"), MetricName: awssdk.String("metric_1")},
		{Namespace: awssdk.String("namespace_1"), MetricName: awssdk.String("metric_2"), Dimensions: []*cloudwatch.Dimension{{Name: awssdk.String("first"), Value: awssdk.String("dimension")}, {Name: awssdk.String("second"), Value: awssdk.String("dimension")}}},
		{Namespace: awssdk.String("namespace_2"), MetricName: awssdk.String("metric_1")},
		{Namespace: awssdk.String("namespace_2"), MetricName: awssdk.String("metric_2")},
	}
	alarms := []*cloudwatch.MetricAlarm{
		{AlarmArn: awssdk.String("alarm_1")},
		{AlarmArn: awssdk.String("alarm_2")},
		{
			AlarmArn:                awssdk.String("alarm_3"),
			AlarmName:               awssdk.String("my_alarm"),
			ActionsEnabled:          awssdk.Bool(true),
			AlarmActions:            []*string{awssdk.String("action_arn_1"), awssdk.String("action_arn_2"), awssdk.String("action_arn_3")},
			InsufficientDataActions: []*string{awssdk.String("action_arn_1"), awssdk.String("action_arn_3")},
			OKActions:               []*string{awssdk.String("action_arn_2")},
			AlarmDescription:        awssdk.String("my alarm description"),
			Dimensions:              []*cloudwatch.Dimension{{Name: awssdk.String("first"), Value: awssdk.String("dimension")}, {Name: awssdk.String("second"), Value: awssdk.String("dimension")}},
			MetricName:              awssdk.String("metric_2"),
			Namespace:               awssdk.String("namespace_2"),
			StateUpdatedTimestamp:   awssdk.Time(now),
			StateValue:              awssdk.String("OK"),
		},
	}

	mock := &mockCloudwatch{metrics: metrics, metricalarms: alarms}

	service := Monitoring{
		CloudWatchAPI: mock, region: "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildMonitoringFetchFuncs(awsfetch.NewConfig(mock))),
	}

	g, err := service.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery("metric", "alarm"))
	if err != nil {
		t.Fatal(err)
	}
	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties()[p.OKActions].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.AlarmActions].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.InsufficientDataActions].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.Dimensions].([]*graph.KeyValue); ok {
			sort.Slice(p, func(i, j int) bool {
				if p[i].KeyName != p[j].KeyName {
					return p[i].KeyName < p[j].KeyName
				}
				return p[i].Value <= p[j].Value
			})
		}
	}

	expected := map[string]cloud.Resource{
		"awls-4ba90752": resourcetest.Metric("awls-4ba90752").Prop(p.Name, "metric_1").Prop(p.Namespace, "namespace_1").Build(),
		"awls-4baa0753": resourcetest.Metric("awls-4baa0753").Prop(p.Name, "metric_2").Prop(p.Namespace, "namespace_1").Prop(p.Dimensions, []*graph.KeyValue{{KeyName: "first", Value: "dimension"}, {KeyName: "second", Value: "dimension"}}).Build(),
		"awls-4bb20753": resourcetest.Metric("awls-4bb20753").Prop(p.Name, "metric_1").Prop(p.Namespace, "namespace_2").Build(),
		"awls-4bb30754": resourcetest.Metric("awls-4bb30754").Prop(p.Name, "metric_2").Prop(p.Namespace, "namespace_2").Build(),
		"alarm_1":       resourcetest.Alarm("alarm_1").Prop(p.Arn, "alarm_1").Build(),
		"alarm_2":       resourcetest.Alarm("alarm_2").Prop(p.Arn, "alarm_2").Build(),
		"alarm_3": resourcetest.Alarm("alarm_3").Prop(p.Arn, "alarm_3").Prop(p.Name, "my_alarm").Prop(p.ActionsEnabled, true).Prop(p.AlarmActions, []string{"action_arn_1", "action_arn_2", "action_arn_3"}).Prop(p.InsufficientDataActions, []string{"action_arn_1", "action_arn_3"}).
			Prop(p.OKActions, []string{"action_arn_2"}).Prop(p.Description, "my alarm description").Prop(p.Dimensions, []*graph.KeyValue{{KeyName: "first", Value: "dimension"}, {KeyName: "second", Value: "dimension"}}).Prop(p.MetricName, "metric_2").
			Prop(p.Namespace, "namespace_2").Prop(p.Updated, now).Prop(p.State, "OK").Build(),
	}

	expectedChildren := map[string][]string{
		"eu-west-1": {"awls-4ba90752", "awls-4baa0753", "awls-4bb20753", "awls-4bb30754", "alarm_1", "alarm_2", "alarm_3"},
	}
	expectedAppliedOn := map[string][]string{
		"alarm_3": {"awls-4bb30754"},
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildCdnGraph(t *testing.T) {
	now := time.Now().UTC()
	distributions := []*cloudfront.DistributionSummary{
		{
			ARN:              awssdk.String("ds_1_arn"),
			Aliases:          &cloudfront.Aliases{Items: []*string{awssdk.String("cname1.domain.name"), awssdk.String("cname2.domain.name")}, Quantity: awssdk.Int64(2)},
			Comment:          awssdk.String("my cdn distribution"),
			DomainName:       awssdk.String("domain.name"),
			Enabled:          awssdk.Bool(true),
			HttpVersion:      awssdk.String("http/2"),
			Id:               awssdk.String("ds_1"),
			IsIPV6Enabled:    awssdk.Bool(true),
			LastModifiedTime: awssdk.Time(now),
			Origins: &cloudfront.Origins{
				Quantity: awssdk.Int64(2),
				Items: []*cloudfront.Origin{
					{
						DomainName:     awssdk.String("domain.name"),
						Id:             awssdk.String("origin_1"),
						OriginPath:     awssdk.String("my/s3/path"),
						S3OriginConfig: &cloudfront.S3OriginConfig{OriginAccessIdentity: awssdk.String("origin-access-identity/CloudFront/ID-of-origin-access-identity")},
					},
					{
						DomainName: awssdk.String("domain2.name"),
						Id:         awssdk.String("origin_2"),
						OriginPath: awssdk.String("my/other/path"),
					},
				},
			},
			PriceClass: awssdk.String("expensive"),
			Status:     awssdk.String("running"),
			ViewerCertificate: &cloudfront.ViewerCertificate{
				ACMCertificateArn: awssdk.String("acm-certificate"),
				Certificate:       awssdk.String("<ViewerProtocolPolicy>https-only<ViewerProtocolPolicy>"),
				//IAMCertificateId:             awssdk.String("iam-certificate"),
				MinimumProtocolVersion: awssdk.String("TLSv1"),
				SSLSupportMethod:       awssdk.String("sni-only"),
			},
			WebACLId: awssdk.String("id"),
		},
		{
			ARN:        awssdk.String("ds_2_arn"),
			DomainName: awssdk.String("other.domain.name"),
			Id:         awssdk.String("ds_2"),
		},
		{
			Id: awssdk.String("ds_3"),
		},
	}

	mock := &mockCloudfront{distributionsummarys: distributions}

	service := Cdn{
		CloudFrontAPI: mock, region: "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildCdnFetchFuncs(awsfetch.NewConfig(mock))),
	}

	g, err := service.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery("distribution"))
	if err != nil {
		t.Fatal(err)
	}

	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties()[p.Aliases].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.Origins].([]*graph.DistributionOrigin); ok {
			sort.Slice(p, func(i, j int) bool {
				return p[i].ID <= p[j].ID
			})
		}
	}

	expected := map[string]cloud.Resource{
		"ds_1": resourcetest.Distribution("ds_1").
			Prop(p.Arn, "ds_1_arn").
			Prop(p.Aliases, []string{"cname1.domain.name", "cname2.domain.name"}).
			Prop(p.Comment, "my cdn distribution").
			Prop(p.PublicDNS, "domain.name").
			Prop(p.Enabled, true).
			Prop(p.HTTPVersion, "http/2").
			Prop(p.IPv6Enabled, true).
			Prop(p.Modified, now).
			Prop(p.PriceClass, "expensive").
			Prop(p.State, "running").
			Prop(p.WebACL, "id").
			Prop(p.ACMCertificate, "acm-certificate").
			Prop(p.Certificate, "<ViewerProtocolPolicy>https-only<ViewerProtocolPolicy>").
			Prop(p.TLSVersionRequired, "TLSv1").
			Prop(p.SSLSupportMethod, "sni-only").
			Prop(p.Origins, []*graph.DistributionOrigin{
				{ID: "origin_1", PublicDNS: "domain.name", OriginType: "s3", PathPrefix: "my/s3/path", Config: "origin-access-identity/CloudFront/ID-of-origin-access-identity"},
				{ID: "origin_2", PublicDNS: "domain2.name", PathPrefix: "my/other/path"},
			}).
			Build(),
		"ds_2": resourcetest.Distribution("ds_2").Prop(p.Arn, "ds_2_arn").Prop(p.PublicDNS, "other.domain.name").Build(),
		"ds_3": resourcetest.Distribution("ds_3").Build(),
	}

	expectedChildren := map[string][]string{}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildCloudFormationGraph(t *testing.T) {
	now := time.Now().UTC()
	stacks := []*cloudformation.Stack{
		{
			Capabilities:      []*string{awssdk.String("cap_1"), awssdk.String("cap_2"), awssdk.String("cap_3")},
			ChangeSetId:       awssdk.String("changeset"),
			CreationTime:      awssdk.Time(now.Add(-2 * time.Hour)),
			Description:       awssdk.String("my cf stack"),
			DisableRollback:   awssdk.Bool(true),
			LastUpdatedTime:   awssdk.Time(now),
			NotificationARNs:  []*string{awssdk.String("notif_1"), awssdk.String("notif_2")},
			Outputs:           []*cloudformation.Output{{OutputKey: awssdk.String("output1"), OutputValue: awssdk.String("myoutput1")}, {OutputKey: awssdk.String("output2"), OutputValue: awssdk.String("myoutput2")}},
			Parameters:        []*cloudformation.Parameter{{ParameterKey: awssdk.String("key1"), ParameterValue: awssdk.String("val1")}, {ParameterKey: awssdk.String("key2"), ParameterValue: awssdk.String("val2")}},
			RoleARN:           awssdk.String("role_arn"),
			StackId:           awssdk.String("id_1"),
			StackName:         awssdk.String("name_1"),
			StackStatus:       awssdk.String("deployed"),
			StackStatusReason: awssdk.String("evrything ok"),
		},
		{
			StackId:   awssdk.String("id_2"),
			StackName: awssdk.String("name_2"),
		},
		{
			StackId: awssdk.String("id_3"),
		},
	}

	mock := &mockCloudformation{stacks: stacks}

	service := Cloudformation{
		CloudFormationAPI: mock, region: "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildCloudformationFetchFuncs(awsfetch.NewConfig(mock))),
	}

	g, err := service.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery("stack"))
	if err != nil {
		t.Fatal(err)
	}
	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties()[p.Capabilities].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.Notifications].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties()[p.Parameters].([]*graph.KeyValue); ok {
			sort.Slice(p, func(i, j int) bool {
				if p[i].KeyName != p[j].KeyName {
					return p[i].KeyName < p[j].KeyName
				}
				return p[i].Value <= p[j].Value
			})
		}
		if p, ok := res.Properties()[p.Outputs].([]*graph.KeyValue); ok {
			sort.Slice(p, func(i, j int) bool {
				if p[i].KeyName != p[j].KeyName {
					return p[i].KeyName < p[j].KeyName
				}
				return p[i].Value <= p[j].Value
			})
		}
	}

	expected := map[string]cloud.Resource{
		"id_1": resourcetest.Stack("id_1").
			Prop(p.Name, "name_1").
			Prop(p.Capabilities, []string{"cap_1", "cap_2", "cap_3"}).
			Prop(p.ChangeSet, "changeset").
			Prop(p.Created, now.Add(-2*time.Hour)).
			Prop(p.Description, "my cf stack").
			Prop(p.DisableRollback, true).
			Prop(p.Modified, now).
			Prop(p.Notifications, []string{"notif_1", "notif_2"}).
			Prop(p.Outputs, []*graph.KeyValue{{KeyName: "output1", Value: "myoutput1"}, {KeyName: "output2", Value: "myoutput2"}}).
			Prop(p.Parameters, []*graph.KeyValue{{KeyName: "key1", Value: "val1"}, {KeyName: "key2", Value: "val2"}}).
			Prop(p.Role, "role_arn").
			Prop(p.State, "deployed").
			Prop(p.StateMessage, "evrything ok").
			Build(),
		"id_2": resourcetest.Stack("id_2").Prop(p.Name, "name_2").Build(),
		"id_3": resourcetest.Stack("id_3").Build(),
	}

	expectedChildren := map[string][]string{
		"eu-west-1": {"id_1", "id_2", "id_3"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildEmptyRdfGraphWhenNoData(t *testing.T) {
	expectG := graph.NewGraph()
	expectG.AddResource(resourcetest.Region("eu-west-1").Build())

	mock := &mockIam{}

	access := Access{
		IAMAPI: mock, region: "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildAccessFetchFuncs(awsfetch.NewConfig(mock))),
	}

	g, err := access.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]cloud.Resource{
		"eu-west-1": resourcetest.Region("eu-west-1").Build(),
	}
	expectedChildren, expectedAppliedOn := map[string][]string{}, map[string][]string{}
	resources, err := g.Find(cloud.NewQuery("region"))
	if err != nil {
		t.Fatal(err)
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)

	infra := Infra{
		EC2API:         &mockEc2{},
		ELBAPI:         &mockElb{},
		ELBV2API:       &mockElbv2{},
		RDSAPI:         &mockRds{},
		AutoScalingAPI: &mockAutoscaling{},
		ECRAPI:         &mockEcr{},
		ECSAPI:         &mockEcs{},
		ACMAPI:         &mockAcm{},
		region:         "eu-west-1",
		fetcher: fetch.NewFetcher(awsfetch.BuildInfraFetchFuncs(awsfetch.NewConfig(
			&mockEc2{}, &mockElb{}, &mockElbv2{}, &mockRds{}, &mockEcr{}, &mockEcs{}, &mockAutoscaling{}, &mockAcm{},
		))),
	}

	g, err = infra.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	resources, err = g.Find(cloud.NewQuery("region"))
	if err != nil {
		t.Fatal(err)
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func mustGetChildrenId(g cloud.GraphAPI, res cloud.Resource) []string {
	var collect []string
	children, err := g.ResourceRelations(res, rdf.ChildrenOfRel, false)
	if err != nil {
		panic(err)
	}
	for _, child := range children {
		collect = append(collect, child.Id())
	}
	return collect
}

func mustGetAppliedOnId(g cloud.GraphAPI, res cloud.Resource) []string {
	var collect []string
	children, err := g.ResourceRelations(res, rdf.ApplyOn, false)
	if err != nil {
		panic(err)
	}
	for _, child := range children {
		collect = append(collect, child.Id())
	}
	return collect
}

func compareResources(t *testing.T, g cloud.GraphAPI, resources []cloud.Resource, expected map[string]cloud.Resource, expectedChildren, expectedAppliedOn map[string][]string) {
	t.Helper()
	if got, want := len(resources), len(expected); got != want {
		t.Errorf("got %d, want %d", got, want)
		//t.Fatalf("got %#v\nwant %#v\n", resources, expected)
	}
	for _, got := range resources {
		want := expected[got.Id()]
		if !reflect.DeepEqual(got, want) {
			//fmt.Println("got:")
			//pretty.Print(got)
			//fmt.Println("\nwant:")
			//pretty.Print(want)
			t.Errorf("got \n%#v\nwant\n%#v", got, want)
		}
		children := mustGetChildrenId(g, got)
		sort.Strings(children)
		if g, w := children, expectedChildren[got.Id()]; !reflect.DeepEqual(g, w) {
			t.Errorf("'%s' children: got %v, want %v", got.Id(), g, w)
		}
		appliedOn := mustGetAppliedOnId(g, got)
		sort.Strings(appliedOn)
		if g, w := appliedOn, expectedAppliedOn[got.Id()]; !reflect.DeepEqual(g, w) {
			t.Errorf("'%s' appliedOn: got %v, want %v", got.Id(), g, w)
		}
	}
}

func TestSliceOfSlice(t *testing.T) {
	var empty [][]*string
	tcases := []struct {
		in        []*string
		maxlength int
		out       [][]*string
	}{
		{in: []*string{awssdk.String("1"), awssdk.String("2"), awssdk.String("3")}, maxlength: 2, out: [][]*string{{awssdk.String("1"), awssdk.String("2")}, {awssdk.String("3")}}},
		{in: []*string{awssdk.String("1"), awssdk.String("2"), awssdk.String("3")}, maxlength: 1, out: [][]*string{{awssdk.String("1")}, {awssdk.String("2")}, {awssdk.String("3")}}},
		{in: []*string{awssdk.String("1"), awssdk.String("2"), awssdk.String("3")}, maxlength: 3, out: [][]*string{{awssdk.String("1"), awssdk.String("2"), awssdk.String("3")}}},
		{in: []*string{awssdk.String("1"), awssdk.String("2"), awssdk.String("3")}, maxlength: 5, out: [][]*string{{awssdk.String("1"), awssdk.String("2"), awssdk.String("3")}}},
		{in: []*string{awssdk.String("1"), awssdk.String("2"), awssdk.String("3")}, maxlength: 0, out: empty},
		{in: []*string{}, maxlength: 2, out: empty},
		{in: []*string{awssdk.String("1"), awssdk.String("2"), awssdk.String("3"), awssdk.String("4")}, maxlength: 2, out: [][]*string{{awssdk.String("1"), awssdk.String("2")}, {awssdk.String("3"), awssdk.String("4")}}},
	}
	for i, tcase := range tcases {
		if got, want := sliceOfSlice(tcase.in, tcase.maxlength), tcase.out; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %+v, want %+v", i+1, got, want)
		}
	}
}

func sliceOfSlice(in []*string, maxLength int) (res [][]*string) {
	if maxLength <= 0 {
		return
	}
	if len(in) == 0 {
		return
	}
	for i := 0; i < len(in); i += maxLength {
		if i+maxLength < len(in) {
			res = append(res, in[i:i+maxLength])
		} else {
			res = append(res, in[i:])
		}
	}

	return
}
