package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestSubnet(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create subnet name=my-subnet cidr=10.10.10.0/24 vpc=any-vpc-id availabilityzone=eu-west-1a").Mock(&ec2Mock{
			CreateSubnetFunc: func(input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
				return &ec2.CreateSubnetOutput{Subnet: &ec2.Subnet{SubnetId: String("new-subnet-id")}}, nil
			},
			CreateTagsRequestFunc: func(input *ec2.CreateTagsInput) (req *request.Request, output *ec2.CreateTagsOutput) {
				output = &ec2.CreateTagsOutput{}
				req = request.New(aws.Config{}, metadata.ClientInfo{}, request.Handlers{}, nil, &request.Operation{}, input, output)
				return
			}}).
			ExpectInput("CreateSubnet", &ec2.CreateSubnetInput{
				AvailabilityZone: String("eu-west-1a"),
				CidrBlock:        String("10.10.10.0/24"),
				VpcId:            String("any-vpc-id"),
			}).
			ExpectInput("CreateTagsRequest", &ec2.CreateTagsInput{
				Resources: []*string{String("new-subnet-id")},
				Tags:      []*ec2.Tag{{Key: String("Name"), Value: String("my-subnet")}},
			}).ExpectCommandResult("new-subnet-id").ExpectCalls("CreateSubnet", "CreateTagsRequest").Run(t)
	})

	t.Run("create public", func(t *testing.T) {
		Template("create subnet public=true name=my-subnet cidr=10.10.10.0/24 vpc=any-vpc-id availabilityzone=eu-west-1a").Mock(&ec2Mock{
			CreateSubnetFunc: func(input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
				return &ec2.CreateSubnetOutput{Subnet: &ec2.Subnet{SubnetId: String("new-subnet-id")}}, nil
			},
			CreateTagsRequestFunc: func(input *ec2.CreateTagsInput) (req *request.Request, output *ec2.CreateTagsOutput) {
				output = &ec2.CreateTagsOutput{}
				req = request.New(aws.Config{}, metadata.ClientInfo{}, request.Handlers{}, nil, &request.Operation{}, input, output)
				return
			}, ModifySubnetAttributeFunc: func(input *ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error) {
				return nil, nil
			}}).
			ExpectInput("CreateSubnet", &ec2.CreateSubnetInput{
				AvailabilityZone: String("eu-west-1a"),
				CidrBlock:        String("10.10.10.0/24"),
				VpcId:            String("any-vpc-id"),
			}).
			ExpectInput("CreateTagsRequest", &ec2.CreateTagsInput{
				Resources: []*string{String("new-subnet-id")},
				Tags:      []*ec2.Tag{{Key: String("Name"), Value: String("my-subnet")}},
			}).
			ExpectInput("ModifySubnetAttribute", &ec2.ModifySubnetAttributeInput{
				MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{Value: Bool(true)},
				SubnetId:            String("new-subnet-id"),
			}).ExpectCommandResult("new-subnet-id").ExpectCalls("CreateSubnet", "CreateTagsRequest", "ModifySubnetAttribute").Run(t)
	})

	t.Run("update", func(t *testing.T) {
		Template("update subnet id=any-subnet-id public=true").Mock(&ec2Mock{
			ModifySubnetAttributeFunc: func(input *ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error) {
				return nil, nil
			}}).
			ExpectInput("ModifySubnetAttribute", &ec2.ModifySubnetAttributeInput{
				MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{Value: Bool(true)},
				SubnetId:            String("any-subnet-id"),
			}).ExpectCalls("ModifySubnetAttribute").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete subnet id=any-subnet-id").Mock(&ec2Mock{
			DeleteSubnetFunc: func(input *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteSubnet", &ec2.DeleteSubnetInput{
				SubnetId: String("any-subnet-id"),
			}).ExpectCalls("DeleteSubnet").Run(t)
	})
}
