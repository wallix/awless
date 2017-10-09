package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

func TestInstanceprofile(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create instanceprofile name=my-instance-profile").Mock(&iamMock{
			CreateInstanceProfileFunc: func(input *iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
				return nil, nil
			}}).
			ExpectInput("CreateInstanceProfile", &iam.CreateInstanceProfileInput{
				InstanceProfileName: String("my-instance-profile"),
			}).ExpectCalls("CreateInstanceProfile").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete instanceprofile name=my-instance-profile").Mock(&iamMock{
			DeleteInstanceProfileFunc: func(input *iam.DeleteInstanceProfileInput) (*iam.DeleteInstanceProfileOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteInstanceProfile", &iam.DeleteInstanceProfileInput{
				InstanceProfileName: String("my-instance-profile"),
			}).ExpectCalls("DeleteInstanceProfile").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		t.Run("without previous associations", func(t *testing.T) {
			Template("attach instanceprofile name=my-instance-profile instance=i-12345 replace=true").Mock(&ec2Mock{
				DescribeIamInstanceProfileAssociationsFunc: func(input *ec2.DescribeIamInstanceProfileAssociationsInput) (*ec2.DescribeIamInstanceProfileAssociationsOutput, error) {
					return &ec2.DescribeIamInstanceProfileAssociationsOutput{IamInstanceProfileAssociations: []*ec2.IamInstanceProfileAssociation{}}, nil
				},
				AssociateIamInstanceProfileFunc: func(input *ec2.AssociateIamInstanceProfileInput) (*ec2.AssociateIamInstanceProfileOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DescribeIamInstanceProfileAssociations", &ec2.DescribeIamInstanceProfileAssociationsInput{
				Filters: []*ec2.Filter{
					{Name: String("instance-id"), Values: []*string{String("i-12345")}},
					{Name: String("state"), Values: []*string{String("associated")}},
				},
			}).ExpectInput("AssociateIamInstanceProfile", &ec2.AssociateIamInstanceProfileInput{
				InstanceId: String("i-12345"),
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: String("my-instance-profile"),
				},
			}).ExpectCalls("DescribeIamInstanceProfileAssociations", "AssociateIamInstanceProfile").Run(t)
		})
		t.Run("with existing associations", func(t *testing.T) {
			Template("attach instanceprofile name=my-instance-profile instance=i-12345 replace=true").Mock(&ec2Mock{
				DescribeIamInstanceProfileAssociationsFunc: func(input *ec2.DescribeIamInstanceProfileAssociationsInput) (*ec2.DescribeIamInstanceProfileAssociationsOutput, error) {
					return &ec2.DescribeIamInstanceProfileAssociationsOutput{
						IamInstanceProfileAssociations: []*ec2.IamInstanceProfileAssociation{
							{AssociationId: String("my-assoc-1"), InstanceId: String("i-12345"), IamInstanceProfile: &ec2.IamInstanceProfile{Arn: String("arn:of:old:profile")}},
						},
					}, nil
				},
				ReplaceIamInstanceProfileAssociationFunc: func(input *ec2.ReplaceIamInstanceProfileAssociationInput) (*ec2.ReplaceIamInstanceProfileAssociationOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DescribeIamInstanceProfileAssociations", &ec2.DescribeIamInstanceProfileAssociationsInput{
				Filters: []*ec2.Filter{
					{Name: String("instance-id"), Values: []*string{String("i-12345")}},
					{Name: String("state"), Values: []*string{String("associated")}},
				},
			}).ExpectInput("ReplaceIamInstanceProfileAssociation", &ec2.ReplaceIamInstanceProfileAssociationInput{
				AssociationId: String("my-assoc-1"),
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: String("my-instance-profile"),
				},
			}).ExpectCalls("DescribeIamInstanceProfileAssociations", "ReplaceIamInstanceProfileAssociation").Run(t)
		})
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach instanceprofile name=my-instance-profile instance=i-12345").Mock(&ec2Mock{
			DescribeIamInstanceProfileAssociationsFunc: func(input *ec2.DescribeIamInstanceProfileAssociationsInput) (*ec2.DescribeIamInstanceProfileAssociationsOutput, error) {
				return &ec2.DescribeIamInstanceProfileAssociationsOutput{
					IamInstanceProfileAssociations: []*ec2.IamInstanceProfileAssociation{
						{AssociationId: String("my-assoc-1"), InstanceId: String("i-12345"), IamInstanceProfile: &ec2.IamInstanceProfile{Arn: String("arn:of:my-instance-profile")}},
					},
				}, nil
			},
			DisassociateIamInstanceProfileFunc: func(input *ec2.DisassociateIamInstanceProfileInput) (*ec2.DisassociateIamInstanceProfileOutput, error) {
				return &ec2.DisassociateIamInstanceProfileOutput{
					IamInstanceProfileAssociation: &ec2.IamInstanceProfileAssociation{
						IamInstanceProfile: &ec2.IamInstanceProfile{Id: String("assoc-id")},
					},
				}, nil
			},
		}).ExpectInput("DescribeIamInstanceProfileAssociations", &ec2.DescribeIamInstanceProfileAssociationsInput{
			Filters: []*ec2.Filter{
				{Name: String("instance-id"), Values: []*string{String("i-12345")}},
			},
		}).ExpectInput("DisassociateIamInstanceProfile", &ec2.DisassociateIamInstanceProfileInput{
			AssociationId: String("my-assoc-1"),
		}).ExpectCommandResult("assoc-id").ExpectCalls("DescribeIamInstanceProfileAssociations", "DisassociateIamInstanceProfile").Run(t)
	})
}
