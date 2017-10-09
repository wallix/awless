package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestTag(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create tag key=MyKey resource=any-resource-id value=Value").Mock(&ec2Mock{
			CreateTagsRequestFunc: func(input *ec2.CreateTagsInput) (req *request.Request, output *ec2.CreateTagsOutput) {
				output = &ec2.CreateTagsOutput{}
				req = request.New(aws.Config{}, metadata.ClientInfo{}, request.Handlers{}, nil, &request.Operation{}, input, output)
				return
			}}).
			ExpectInput("CreateTagsRequest", &ec2.CreateTagsInput{
				Resources: []*string{String("any-resource-id")},
				Tags:      []*ec2.Tag{{Key: String("MyKey"), Value: String("Value")}},
			}).ExpectCalls("CreateTagsRequest").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete tag key=MyKey resource=any-resource-id").Mock(&ec2Mock{
			DeleteTagsFunc: func(input *ec2.DeleteTagsInput) (*ec2.DeleteTagsOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteTags", &ec2.DeleteTagsInput{
				Resources: []*string{String("any-resource-id")},
				Tags:      []*ec2.Tag{{Key: String("MyKey")}},
			}).ExpectCalls("DeleteTags").Run(t)
	})
}
