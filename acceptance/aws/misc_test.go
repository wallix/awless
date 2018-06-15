package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

func TestReducingReferences(t *testing.T) {
	t.Run("reference created from command run", func(t *testing.T) {
		tcases := []struct {
			template      string
			expRevert     string
			expStoppedIds []*string
		}{
			{
				template: `inst = create instance name=myinstance count=1 image=ami-12345 subnet=sub-1234 type=t2.nano
	stop instance id=$inst`,
				expStoppedIds: []*string{String("new-instance-id")},
				expRevert: `check instance id=new-instance-id state=stopped timeout=180
start instance ids=new-instance-id
delete instance id=new-instance-id`,
			},
			{
				template: `inst = create instance name=myinstance count=1 image=ami-12345 subnet=sub-1234 type=t2.nano
	stop instance id=[id-1234,$inst,id-2345]`,
				expStoppedIds: []*string{String("id-1234"), String("new-instance-id"), String("id-2345")},
				expRevert: `check instance id=id-1234 state=stopped timeout=180
check instance id=new-instance-id state=stopped timeout=180
check instance id=id-2345 state=stopped timeout=180
start instance ids=[id-1234,new-instance-id,id-2345]
delete instance id=new-instance-id`,
			},
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
			}).ExpectCalls("RunInstances", "CreateTagsRequest", "StopInstances").ExpectRevert(tcase.expRevert).Run(t)
		}
	})
	t.Run("reference from inlined variable", func(t *testing.T) {
		Template(`inst = i-1234
		stop instance id=$inst`).Mock(&ec2Mock{StopInstancesFunc: func(input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
			return &ec2.StopInstancesOutput{StoppingInstances: []*ec2.InstanceStateChange{{InstanceId: String("i-1234")}}}, nil
		},
		}).ExpectInput("StopInstances", &ec2.StopInstancesInput{
			InstanceIds: []*string{String("i-1234")},
		}).ExpectCalls("StopInstances").ExpectRevert("check instance id=i-1234 state=stopped timeout=180\nstart instance ids=i-1234").Run(t)
	})

	t.Run("reverting with multiple ids", func(t *testing.T) {
		Template(`inst = create instance name=myinstance count=1 image=ami-12345 subnet=sub-1234 type=t2.nano
start instance id=[id-1234,$inst,id-2345]`).Mock(&ec2Mock{
			RunInstancesFunc: func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
				return &ec2.Reservation{Instances: []*ec2.Instance{{InstanceId: String("new-instance-id")}}}, nil
			},
			CreateTagsRequestFunc: func(input *ec2.CreateTagsInput) (req *request.Request, output *ec2.CreateTagsOutput) {
				output = &ec2.CreateTagsOutput{}
				req = request.New(aws.Config{}, metadata.ClientInfo{}, request.Handlers{}, nil, &request.Operation{}, input, output)
				return
			},
			StartInstancesFunc: func(input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
				return &ec2.StartInstancesOutput{StartingInstances: []*ec2.InstanceStateChange{{InstanceId: String("new-instance-id")}}}, nil
			},
		}).IgnoreInput("RunInstances", "CreateTagsRequest").ExpectInput("StartInstances", &ec2.StartInstancesInput{
			InstanceIds: []*string{String("id-1234"), String("new-instance-id"), String("id-2345")},
		}).ExpectCalls("RunInstances", "CreateTagsRequest", "StartInstances").
			ExpectRevert(`check instance id=id-1234 state=running timeout=180
check instance id=new-instance-id state=running timeout=180
check instance id=id-2345 state=running timeout=180
stop instance ids=[id-1234,new-instance-id,id-2345]
delete instance id=new-instance-id`).Run(t)
	})
}

func TestWithMissingRequiredParams(t *testing.T) {
	Template("create user").Mock(&iamMock{
		CreateUserFunc: func(input *iam.CreateUserInput) (*iam.CreateUserOutput, error) {
			return &iam.CreateUserOutput{User: &iam.User{UserId: String("new-user-id")}}, nil
		}}).Fillers(map[string]string{"user.name": "donald"}).ExpectInput("CreateUser", &iam.CreateUserInput{
		UserName: String("donald"),
	}).ExpectCommandResult("new-user-id").ExpectCalls("CreateUser").Run(t)
}

func TestConcatenationWithReferences(t *testing.T) {
	Template(`mfaResource = "arn:aws:iam::" + {account.id} + ":mfa/${aws:username}"
userResource = "arn:aws:iam::" + {account.id} + ":user/${aws:username}"

policy = create policy name=ManageOwnMFADevice effect=allow action=[iam:CreateVirtualMFADevice,iam:EnableMFADevice,iam:ResyncMFADevice,iam:DeleteVirtualMFADevice] resource=$mfaResource,$userResource`).
		Fillers(map[string]string{"account.id": "0123456789"}).
		Mock(&iamMock{CreatePolicyFunc: func(input *iam.CreatePolicyInput) (*iam.CreatePolicyOutput, error) {
			return &iam.CreatePolicyOutput{Policy: &iam.Policy{Arn: String("arn:new-policy-arn")}}, nil
		},
		}).ExpectInput("CreatePolicy", &iam.CreatePolicyInput{
		PolicyName: String("ManageOwnMFADevice"),
		PolicyDocument: String(`{
 "Version": "2012-10-17",
 "Statement": [
  {
   "Effect": "Allow",
   "Action": [
    "iam:CreateVirtualMFADevice",
    "iam:EnableMFADevice",
    "iam:ResyncMFADevice",
    "iam:DeleteVirtualMFADevice"
   ],
   "Resource": [
    "arn:aws:iam::0123456789:mfa/${aws:username}",
    "arn:aws:iam::0123456789:user/${aws:username}"
   ]
  }
 ]
}`)}).ExpectCalls("CreatePolicy").ExpectRevert("delete policy all-versions=true arn=arn:new-policy-arn").Run(t)
}
