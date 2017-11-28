package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestSecuritygroup(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template(`create securitygroup name=my-sg-name vpc=my-vpc-id description="security group description"`).Mock(&ec2Mock{
			CreateSecurityGroupFunc: func(input *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
				return &ec2.CreateSecurityGroupOutput{GroupId: String("new-secgroup-id")}, nil
			}}).
			ExpectInput("CreateSecurityGroup", &ec2.CreateSecurityGroupInput{
				GroupName:   String("my-sg-name"),
				VpcId:       String("my-vpc-id"),
				Description: String("security group description"),
			}).ExpectCommandResult("new-secgroup-id").ExpectCalls("CreateSecurityGroup").Run(t)
	})

	t.Run("update", func(t *testing.T) {
		t.Run("inbound authorize with another secgroup", func(t *testing.T) {
			Template("update securitygroup id=my-secgroup-id inbound=authorize protocol=tcp securitygroup=any-secgroup-id portrange=8080").Mock(&ec2Mock{
				AuthorizeSecurityGroupIngressFunc: func(input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
					return nil, nil
				}}).
				ExpectInput("AuthorizeSecurityGroupIngress", &ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: String("my-secgroup-id"),
					IpPermissions: []*ec2.IpPermission{
						{
							UserIdGroupPairs: []*ec2.UserIdGroupPair{{GroupId: String("any-secgroup-id")}},
							IpProtocol:       String("tcp"),
							FromPort:         Int64(8080),
							ToPort:           Int64(8080),
						},
					},
				}).ExpectCalls("AuthorizeSecurityGroupIngress").Run(t)
		})

		t.Run("inbound authorize", func(t *testing.T) {
			Template("update securitygroup id=my-secgroup-id inbound=authorize protocol=tcp cidr=10.10.10.0/24 portrange=10-22").Mock(&ec2Mock{
				AuthorizeSecurityGroupIngressFunc: func(input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
					return nil, nil
				}}).
				ExpectInput("AuthorizeSecurityGroupIngress", &ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: String("my-secgroup-id"),
					IpPermissions: []*ec2.IpPermission{
						{
							IpProtocol: String("tcp"),
							IpRanges:   []*ec2.IpRange{{CidrIp: String("10.10.10.0/24")}},
							FromPort:   Int64(10),
							ToPort:     Int64(22),
						},
					},
				}).ExpectCalls("AuthorizeSecurityGroupIngress").Run(t)
		})
		t.Run("inbound revoke", func(t *testing.T) {
			Template("update securitygroup id=my-secgroup-id inbound=revoke protocol=tcp cidr=10.10.10.0/24 portrange=10-22").Mock(&ec2Mock{
				RevokeSecurityGroupIngressFunc: func(input *ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error) {
					return nil, nil
				}}).
				ExpectInput("RevokeSecurityGroupIngress", &ec2.RevokeSecurityGroupIngressInput{
					GroupId: String("my-secgroup-id"),
					IpPermissions: []*ec2.IpPermission{
						{
							IpProtocol: String("tcp"),
							IpRanges:   []*ec2.IpRange{{CidrIp: String("10.10.10.0/24")}},
							FromPort:   Int64(10),
							ToPort:     Int64(22),
						},
					},
				}).ExpectCalls("RevokeSecurityGroupIngress").Run(t)
		})
		t.Run("outbound authorize", func(t *testing.T) {
			Template("update securitygroup id=my-secgroup-id outbound=authorize protocol=tcp cidr=10.10.10.0/24 portrange=10-22").Mock(&ec2Mock{
				AuthorizeSecurityGroupEgressFunc: func(input *ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
					return nil, nil
				}}).
				ExpectInput("AuthorizeSecurityGroupEgress", &ec2.AuthorizeSecurityGroupEgressInput{
					GroupId: String("my-secgroup-id"),
					IpPermissions: []*ec2.IpPermission{
						{
							IpProtocol: String("tcp"),
							IpRanges:   []*ec2.IpRange{{CidrIp: String("10.10.10.0/24")}},
							FromPort:   Int64(10),
							ToPort:     Int64(22),
						},
					},
				}).ExpectCalls("AuthorizeSecurityGroupEgress").Run(t)
		})
		t.Run("outbound revoke", func(t *testing.T) {
			Template("update securitygroup id=my-secgroup-id outbound=revoke protocol=tcp cidr=10.10.10.0/24 portrange=10-22").Mock(&ec2Mock{
				RevokeSecurityGroupEgressFunc: func(input *ec2.RevokeSecurityGroupEgressInput) (*ec2.RevokeSecurityGroupEgressOutput, error) {
					return nil, nil
				}}).
				ExpectInput("RevokeSecurityGroupEgress", &ec2.RevokeSecurityGroupEgressInput{
					GroupId: String("my-secgroup-id"),
					IpPermissions: []*ec2.IpPermission{
						{
							IpProtocol: String("tcp"),
							IpRanges:   []*ec2.IpRange{{CidrIp: String("10.10.10.0/24")}},
							FromPort:   Int64(10),
							ToPort:     Int64(22),
						},
					},
				}).ExpectCalls("RevokeSecurityGroupEgress").Run(t)
		})
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete securitygroup id=my-secgroup-id").Mock(&ec2Mock{
			DeleteSecurityGroupFunc: func(input *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteSecurityGroup", &ec2.DeleteSecurityGroupInput{
				GroupId: String("my-secgroup-id"),
			}).ExpectCalls("DeleteSecurityGroup").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach securitygroup id=my-secgroup-id instance=secgroup-instance-id").Mock(&ec2Mock{
			DescribeInstanceAttributeFunc: func(input *ec2.DescribeInstanceAttributeInput) (*ec2.DescribeInstanceAttributeOutput, error) {
				return &ec2.DescribeInstanceAttributeOutput{Groups: []*ec2.GroupIdentifier{
					{GroupId: String("secgroup-1")}, {GroupId: String("secgroup-2")},
				}}, nil
			},
			ModifyInstanceAttributeFunc: func(input *ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DescribeInstanceAttribute", &ec2.DescribeInstanceAttributeInput{
				Attribute:  String("groupSet"),
				InstanceId: String("secgroup-instance-id"),
			}).
			ExpectInput("ModifyInstanceAttribute", &ec2.ModifyInstanceAttributeInput{
				InstanceId: String("secgroup-instance-id"),
				Groups:     []*string{String("secgroup-1"), String("secgroup-2"), String("my-secgroup-id")},
			}).ExpectCalls("DescribeInstanceAttribute", "ModifyInstanceAttribute").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach securitygroup id=my-secgroup-id instance=secgroup-instance-id").Mock(&ec2Mock{
			DescribeInstanceAttributeFunc: func(input *ec2.DescribeInstanceAttributeInput) (*ec2.DescribeInstanceAttributeOutput, error) {
				return &ec2.DescribeInstanceAttributeOutput{Groups: []*ec2.GroupIdentifier{
					{GroupId: String("secgroup-1")}, {GroupId: String("my-secgroup-id")},
				}}, nil
			},
			ModifyInstanceAttributeFunc: func(input *ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DescribeInstanceAttribute", &ec2.DescribeInstanceAttributeInput{
				Attribute:  String("groupSet"),
				InstanceId: String("secgroup-instance-id"),
			}).
			ExpectInput("ModifyInstanceAttribute", &ec2.ModifyInstanceAttributeInput{
				InstanceId: String("secgroup-instance-id"),
				Groups:     []*string{String("secgroup-1")},
			}).ExpectCalls("DescribeInstanceAttribute", "ModifyInstanceAttribute").Run(t)
	})

	t.Run("check", func(t *testing.T) {
		Template("check securitygroup id=my-secgroup-id state=unused timeout=2").Mock(&ec2Mock{
			DescribeNetworkInterfacesFunc: func(input *ec2.DescribeNetworkInterfacesInput) (*ec2.DescribeNetworkInterfacesOutput, error) {
				return &ec2.DescribeNetworkInterfacesOutput{
					NetworkInterfaces: []*ec2.NetworkInterface{},
				}, nil
			}}).ExpectInput("DescribeNetworkInterfaces", &ec2.DescribeNetworkInterfacesInput{
			Filters: []*ec2.Filter{{Name: String("group-id"), Values: []*string{String("my-secgroup-id")}}},
		}).ExpectCalls("DescribeNetworkInterfaces").Run(t)
	})
}
