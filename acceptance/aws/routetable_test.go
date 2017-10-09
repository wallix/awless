package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestRouteTable(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create routetable vpc=vpc-1234").
			Mock(&ec2Mock{
				CreateRouteTableFunc: func(param0 *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
					return &ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: String("new-routetable-id")}}, nil
				},
			}).ExpectInput("CreateRouteTable", &ec2.CreateRouteTableInput{VpcId: String("vpc-1234")}).
			ExpectCommandResult("new-routetable-id").ExpectCalls("CreateRouteTable").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete routetable id=rt-1234").
			Mock(&ec2Mock{
				DeleteRouteTableFunc: func(param0 *ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteRouteTable", &ec2.DeleteRouteTableInput{RouteTableId: String("rt-1234")}).
			ExpectCalls("DeleteRouteTable").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach routetable id=my-rt-id subnet=my-subnet-id").
			Mock(&ec2Mock{
				AssociateRouteTableFunc: func(param0 *ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
					return &ec2.AssociateRouteTableOutput{AssociationId: String("new-assoc-id")}, nil
				},
			}).ExpectInput("AssociateRouteTable", &ec2.AssociateRouteTableInput{
			RouteTableId: String("my-rt-id"),
			SubnetId:     String("my-subnet-id"),
		}).ExpectCommandResult("new-assoc-id").ExpectCalls("AssociateRouteTable").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach routetable association=assoc-2345").
			Mock(&ec2Mock{
				DisassociateRouteTableFunc: func(param0 *ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DisassociateRouteTable", &ec2.DisassociateRouteTableInput{AssociationId: String("assoc-2345")}).
			ExpectCalls("DisassociateRouteTable").Run(t)
	})
}
