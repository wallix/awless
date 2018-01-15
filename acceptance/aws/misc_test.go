package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestReducingReferences(t *testing.T) {
	t.Run("reference created from command run", func(t *testing.T) {
		tcases := []struct {
			template      string
			expStoppedIds []*string
		}{
			{template: `inst = create instance name=myinstance count=1 image=ami-12345 subnet=sub-1234 type=t2.nano
	stop instance id=$inst`, expStoppedIds: []*string{String("new-instance-id")}},
			{template: `inst = create instance name=myinstance count=1 image=ami-12345 subnet=sub-1234 type=t2.nano
	stop instance id=[id-1234,$inst,id-2345]`, expStoppedIds: []*string{String("id-1234"), String("new-instance-id"), String("id-2345")}},
		}
		for _, tcase := range tcases {
			Template(tcase.template).Mock(&ec2Mock{
				RunInstancesFunc: func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
					return &ec2.Reservation{Instances: []*ec2.Instance{{InstanceId: String("new-instance-id")}}}, nil
				},
				CreateTagsRequestFunc: func(input *ec2.CreateTagsInput) (req *request.Request, output *ec2.CreateTagsOutput) {
					output = &ec2.CreateTagsOutput{}
					req = request.New(aws.Config{}, metadata.ClientInfo{}, request.Handlers{}, nil, &request.Operation{}, input, output)
					return
				},
				StopInstancesFunc: func(input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
					return &ec2.StopInstancesOutput{StoppingInstances: []*ec2.InstanceStateChange{{InstanceId: String("new-instance-id")}}}, nil
				},
			}).IgnoreInput("RunInstances", "CreateTagsRequest").ExpectInput("StopInstances", &ec2.StopInstancesInput{
				InstanceIds: tcase.expStoppedIds,
			}).ExpectCalls("RunInstances", "CreateTagsRequest", "StopInstances").Run(t)
		}
	})
	t.Run("reference from inlined variable", func(t *testing.T) {
		Template(`inst = i-1234
		stop instance id=$inst`).Mock(&ec2Mock{StopInstancesFunc: func(input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
			return &ec2.StopInstancesOutput{StoppingInstances: []*ec2.InstanceStateChange{{InstanceId: String("i-1234")}}}, nil
		},
		}).ExpectInput("StopInstances", &ec2.StopInstancesInput{
			InstanceIds: []*string{String("i-1234")},
		}).ExpectCalls("StopInstances").Run(t)
	})

}
