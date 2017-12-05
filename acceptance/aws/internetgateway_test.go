package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestInternetGateway(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create internetgateway").
			Mock(&ec2Mock{
				CreateInternetGatewayFunc: func(param0 *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					return &ec2.CreateInternetGatewayOutput{InternetGateway: &ec2.InternetGateway{InternetGatewayId: String("new-internetgateway-id")}}, nil
				},
			}).ExpectInput("CreateInternetGateway", &ec2.CreateInternetGatewayInput{}).
			ExpectCommandResult("new-internetgateway-id").ExpectCalls("CreateInternetGateway").Run(t)
	})

	t.Run("create-attach meta", func(t *testing.T) {
		Template("create internetgateway vpc=my-vpc-id").
			Mock(&ec2Mock{
				CreateInternetGatewayFunc: func(param0 *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
					return &ec2.CreateInternetGatewayOutput{InternetGateway: &ec2.InternetGateway{InternetGatewayId: String("new-internetgateway-id")}}, nil
				},
				AttachInternetGatewayFunc: func(param0 *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					return nil, nil
				},
			}).ExpectInput("CreateInternetGateway", &ec2.CreateInternetGatewayInput{}).
			ExpectInput("AttachInternetGateway", &ec2.AttachInternetGatewayInput{
				InternetGatewayId: String("new-internetgateway-id"),
				VpcId:             String("my-vpc-id"),
			}).
			ExpectCommandResult("new-internetgateway-id").ExpectCalls("CreateInternetGateway", "AttachInternetGateway").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete internetgateway id=igw-1234").
			Mock(&ec2Mock{
				DeleteInternetGatewayFunc: func(param0 *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteInternetGateway", &ec2.DeleteInternetGatewayInput{InternetGatewayId: String("igw-1234")}).
			ExpectCalls("DeleteInternetGateway").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach internetgateway id=igw-1234 vpc=vpc-2345").
			Mock(&ec2Mock{
				AttachInternetGatewayFunc: func(param0 *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
					return nil, nil
				},
			}).ExpectInput("AttachInternetGateway", &ec2.AttachInternetGatewayInput{
			InternetGatewayId: String("igw-1234"),
			VpcId:             String("vpc-2345"),
		}).
			ExpectCalls("AttachInternetGateway").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach internetgateway id=igw-1234 vpc=vpc-2345").
			Mock(&ec2Mock{
				DetachInternetGatewayFunc: func(param0 *ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DetachInternetGateway", &ec2.DetachInternetGatewayInput{
			InternetGatewayId: String("igw-1234"),
			VpcId:             String("vpc-2345"),
		}).
			ExpectCalls("DetachInternetGateway").Run(t)
	})

}
