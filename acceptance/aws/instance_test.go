package awsat

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/wallix/awless/aws/spec"
)

func TestInstance(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Run("with distro", func(t *testing.T) {
			awsspec.DefaultImageResolverCache.Store("canonical:ubuntu:xenial:x86_64:hvm:ebs", []*awsspec.AwsImage{{Id: "ami-123456"}})

			Template("create instance distro=canonical name=myinstance subnet=sub_1 type=t2.nano count=1").
				Mock(&ec2Mock{
					RunInstancesFunc: func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
						return &ec2.Reservation{Instances: []*ec2.Instance{{InstanceId: String("new-instance-id")}}}, nil
					},
					CreateTagsRequestFunc: func(input *ec2.CreateTagsInput) (req *request.Request, output *ec2.CreateTagsOutput) {
						output = &ec2.CreateTagsOutput{}
						req = request.New(aws.Config{}, metadata.ClientInfo{}, request.Handlers{}, nil, &request.Operation{}, input, output)
						return
					},
				}).ExpectInput("RunInstances", &ec2.RunInstancesInput{
				SubnetId:     String("sub_1"),
				ImageId:      String("ami-123456"),
				InstanceType: String("t2.nano"),
				MinCount:     Int64(1),
				MaxCount:     Int64(1),
			}).ExpectInput("CreateTagsRequest", &ec2.CreateTagsInput{
				Resources: []*string{String("new-instance-id")},
				Tags: []*ec2.Tag{
					{Key: String("Name"), Value: String("myinstance")},
				},
			}).ExpectCommandResult("new-instance-id").ExpectCalls("RunInstances", "CreateTagsRequest").
				ExpectRevert("delete instance id=new-instance-id").Run(t)
		})

		t.Run("with user data", func(t *testing.T) {
			_, userdataFile, cleanup := generateTmpFile("this is my content with {{ .AWLESS.oneRef }} content")
			defer cleanup()

			Template("oneRef=awesome\n"+
				"create instance count=3 image=ami-1234 "+
				"name=myinstance subnet=sub_1 type=t2.nano keypair=mykp ip=10.2.3.4 "+
				"userdata="+userdataFile+" securitygroup=sg-1234 lock=true role=myrole").
				Mock(&ec2Mock{
					RunInstancesFunc: func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
						return &ec2.Reservation{Instances: []*ec2.Instance{{InstanceId: String("new-instance-id")}}}, nil
					},
					CreateTagsRequestFunc: func(input *ec2.CreateTagsInput) (req *request.Request, output *ec2.CreateTagsOutput) {
						output = &ec2.CreateTagsOutput{}
						req = request.New(aws.Config{}, metadata.ClientInfo{}, request.Handlers{}, nil, &request.Operation{}, input, output)
						return
					},
				}).ExpectInput("RunInstances", &ec2.RunInstancesInput{
				SubnetId:              String("sub_1"),
				ImageId:               String("ami-1234"),
				InstanceType:          String("t2.nano"),
				MinCount:              Int64(3),
				MaxCount:              Int64(3),
				KeyName:               String("mykp"),
				PrivateIpAddress:      String("10.2.3.4"),
				SecurityGroupIds:      []*string{String("sg-1234")},
				DisableApiTermination: Bool(true),
				IamInstanceProfile:    &ec2.IamInstanceProfileSpecification{Name: String("myrole")},
				UserData:              String(base64.StdEncoding.EncodeToString([]byte("this is my content with awesome content"))),
			}).ExpectInput("CreateTagsRequest", &ec2.CreateTagsInput{
				Resources: []*string{String("new-instance-id")},
				Tags: []*ec2.Tag{
					{Key: String("Name"), Value: String("myinstance")},
				},
			}).ExpectCommandResult("new-instance-id").ExpectCalls("RunInstances", "CreateTagsRequest").Run(t)
		})
	})

	t.Run("update", func(t *testing.T) {
		Template("update instance id=id-1234 type=t2.micro lock=true").Mock(&ec2Mock{
			ModifyInstanceAttributeFunc: func(param0 *ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error) {
				return nil, nil
			},
		}).ExpectInput("ModifyInstanceAttribute", &ec2.ModifyInstanceAttributeInput{
			InstanceId:            String("id-1234"),
			InstanceType:          &ec2.AttributeValue{Value: String("t2.micro")},
			DisableApiTermination: &ec2.AttributeBooleanValue{Value: Bool(true)},
		}).
			ExpectCalls("ModifyInstanceAttribute").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("one id", func(t *testing.T) {
			Template("delete instance id=id-1234").Mock(&ec2Mock{
				TerminateInstancesFunc: func(param0 *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) { return nil, nil },
			}).ExpectInput("TerminateInstances", &ec2.TerminateInstancesInput{InstanceIds: []*string{String("id-1234")}}).
				ExpectCalls("TerminateInstances").Run(t)
		})

		t.Run("multiple ids", func(t *testing.T) {
			Template("delete instance ids=id-1234,id-2345").Mock(&ec2Mock{
				TerminateInstancesFunc: func(param0 *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) { return nil, nil },
			}).ExpectInput("TerminateInstances", &ec2.TerminateInstancesInput{InstanceIds: []*string{String("id-1234"), String("id-2345")}}).
				ExpectCalls("TerminateInstances").Run(t)
		})
	})

	t.Run("start", func(t *testing.T) {
		t.Run("one id", func(t *testing.T) {
			Template("start instance id=id-1234").Mock(&ec2Mock{
				StartInstancesFunc: func(param0 *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
					return &ec2.StartInstancesOutput{
						StartingInstances: []*ec2.InstanceStateChange{{InstanceId: String("id-1234")}}}, nil
				},
			}).ExpectInput("StartInstances", &ec2.StartInstancesInput{InstanceIds: []*string{String("id-1234")}}).
				ExpectCalls("StartInstances").Run(t)
		})

		t.Run("multiple ids", func(t *testing.T) {
			Template("start instance ids=id-1234,id-2345").Mock(&ec2Mock{
				StartInstancesFunc: func(param0 *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
					return &ec2.StartInstancesOutput{
						StartingInstances: []*ec2.InstanceStateChange{{InstanceId: String("id-1234")}, {InstanceId: String("id-2345")}}}, nil
				},
			}).ExpectInput("StartInstances", &ec2.StartInstancesInput{InstanceIds: []*string{String("id-1234"), String("id-2345")}}).
				ExpectCalls("StartInstances").Run(t)
		})
	})

	t.Run("stop", func(t *testing.T) {
		t.Run("one id", func(t *testing.T) {
			Template("stop instance id=id-1234").Mock(&ec2Mock{
				StopInstancesFunc: func(param0 *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
					return &ec2.StopInstancesOutput{
						StoppingInstances: []*ec2.InstanceStateChange{{InstanceId: String("id-1234")}}}, nil
				},
			}).ExpectInput("StopInstances", &ec2.StopInstancesInput{InstanceIds: []*string{String("id-1234")}}).
				ExpectCalls("StopInstances").Run(t)
		})

		t.Run("multiple ids", func(t *testing.T) {
			Template("stop instance ids=id-1234,id-2345").Mock(&ec2Mock{
				StopInstancesFunc: func(param0 *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
					return &ec2.StopInstancesOutput{
						StoppingInstances: []*ec2.InstanceStateChange{{InstanceId: String("id-1234")}, {InstanceId: String("id-2345")}}}, nil
				},
			}).ExpectInput("StopInstances", &ec2.StopInstancesInput{InstanceIds: []*string{String("id-1234"), String("id-2345")}}).
				ExpectCalls("StopInstances").Run(t)
		})
	})

	t.Run("restart", func(t *testing.T) {
		t.Run("one id", func(t *testing.T) {
			Template("restart instance id=id-1234").Mock(&ec2Mock{
				RebootInstancesFunc: func(param0 *ec2.RebootInstancesInput) (*ec2.RebootInstancesOutput, error) {
					return &ec2.RebootInstancesOutput{}, nil
				},
			}).ExpectInput("RebootInstances", &ec2.RebootInstancesInput{InstanceIds: []*string{String("id-1234")}}).
				ExpectCalls("RebootInstances").Run(t)
		})

		t.Run("multiple ids", func(t *testing.T) {
			Template("restart instance ids=id-1234,id-2345").Mock(&ec2Mock{
				RebootInstancesFunc: func(param0 *ec2.RebootInstancesInput) (*ec2.RebootInstancesOutput, error) {
					return &ec2.RebootInstancesOutput{}, nil
				},
			}).ExpectInput("RebootInstances", &ec2.RebootInstancesInput{InstanceIds: []*string{String("id-1234"), String("id-2345")}}).
				ExpectCalls("RebootInstances").Run(t)
		})
	})

	t.Run("check", func(t *testing.T) {
		Template("check instance id=id-1234 state=running timeout=1").Mock(&ec2Mock{
			DescribeInstancesFunc: func(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
				return &ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{
					{Instances: []*ec2.Instance{{InstanceId: input.InstanceIds[0], State: &ec2.InstanceState{Name: String("running")}}}},
				}}, nil
			}}).ExpectInput("DescribeInstances", &ec2.DescribeInstancesInput{InstanceIds: []*string{String("id-1234")}}).
			ExpectCalls("DescribeInstances").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach instance id=id-1234 targetgroup=arn:of:target:group port=8080").Mock(&elbv2Mock{
			RegisterTargetsFunc: func(param0 *elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
				return nil, nil
			},
		}).ExpectInput("RegisterTargets", &elbv2.RegisterTargetsInput{
			TargetGroupArn: String("arn:of:target:group"),
			Targets: []*elbv2.TargetDescription{
				{Id: String("id-1234"), Port: Int64(8080)},
			},
		}).
			ExpectCalls("RegisterTargets").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach instance id=id-1234 targetgroup=arn:of:target:group").Mock(&elbv2Mock{
			DeregisterTargetsFunc: func(param0 *elbv2.DeregisterTargetsInput) (*elbv2.DeregisterTargetsOutput, error) {
				return nil, nil
			},
		}).ExpectInput("DeregisterTargets", &elbv2.DeregisterTargetsInput{
			TargetGroupArn: String("arn:of:target:group"),
			Targets: []*elbv2.TargetDescription{
				{Id: String("id-1234")},
			},
		}).
			ExpectCalls("DeregisterTargets").Run(t)
	})
}

func generateTmpFile(content string) (*os.File, string, func()) {
	file, err := ioutil.TempFile("", "awless-at-tmpfile")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(file.Name(), []byte(content), 0644); err != nil {
		panic(err)
	}

	cleanup := func() {
		file.Close()
		if err := os.Remove(file.Name()); err != nil {
			panic(err)
		}
	}
	return file, file.Name(), cleanup
}
