package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestNATGateway(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create natgateway elasticip-id=eip-12345 subnet=sub-23456").
			Mock(&ec2Mock{
				CreateNatGatewayFunc: func(param0 *ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
					return &ec2.CreateNatGatewayOutput{NatGateway: &ec2.NatGateway{NatGatewayId: String("new-natgateway-id")}}, nil
				},
			}).ExpectInput("CreateNatGateway", &ec2.CreateNatGatewayInput{
			AllocationId: String("eip-12345"),
			SubnetId:     String("sub-23456"),
		}).
			ExpectCommandResult("new-natgateway-id").ExpectCalls("CreateNatGateway").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete natgateway id=ngw-1234").
			Mock(&ec2Mock{
				DeleteNatGatewayFunc: func(param0 *ec2.DeleteNatGatewayInput) (*ec2.DeleteNatGatewayOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteNatGateway", &ec2.DeleteNatGatewayInput{NatGatewayId: String("ngw-1234")}).
			ExpectCalls("DeleteNatGateway").Run(t)
	})

	t.Run("check", func(t *testing.T) {
		Template("check natgateway id=ngw-1234 state=available timeout=1").
			Mock(&ec2Mock{
				DescribeNatGatewaysFunc: func(param0 *ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
					return &ec2.DescribeNatGatewaysOutput{NatGateways: []*ec2.NatGateway{
						{NatGatewayId: String("ngw-1234"), State: String("available")},
					}}, nil
				},
			}).ExpectInput("DescribeNatGateways", &ec2.DescribeNatGatewaysInput{NatGatewayIds: []*string{String("ngw-1234")}}).
			ExpectCalls("DescribeNatGateways").Run(t)
	})

}
