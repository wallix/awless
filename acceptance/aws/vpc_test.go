package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestVPC(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create vpc name=myvpc cidr=10.0.0.0/16").Mock(&ec2Mock{
			CreateVpcFunc: func(input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
				return &ec2.CreateVpcOutput{Vpc: &ec2.Vpc{VpcId: String("new-vpc-id")}}, nil
			},
			CreateTagsRequestFunc: func(input *ec2.CreateTagsInput) (req *request.Request, output *ec2.CreateTagsOutput) {
				output = &ec2.CreateTagsOutput{}
				req = request.New(aws.Config{}, metadata.ClientInfo{}, request.Handlers{}, nil, &request.Operation{}, input, output)
				return
			}}).
			ExpectInput("CreateVpc", &ec2.CreateVpcInput{CidrBlock: String("10.0.0.0/16")}).
			ExpectInput("CreateTagsRequest", &ec2.CreateTagsInput{
				Resources: []*string{String("new-vpc-id")},
				Tags: []*ec2.Tag{
					{Key: String("Name"), Value: String("myvpc")},
				},
			}).ExpectCommandResult("new-vpc-id").ExpectCalls("CreateVpc", "CreateTagsRequest").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete vpc id=any-vpc-id").Mock(&ec2Mock{
			DeleteVpcFunc: func(input *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
				return &ec2.DeleteVpcOutput{}, nil
			}},
		).ExpectInput("DeleteVpc", &ec2.DeleteVpcInput{VpcId: String("any-vpc-id")}).
			ExpectCalls("DeleteVpc").Run(t)
	})
}
