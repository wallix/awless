package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestRoute(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create route table=table-id cidr=10.0.0.0/16 gateway=igw-id").
			Mock(&ec2Mock{
				CreateRouteFunc: func(param0 *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
					return nil, nil
				},
			}).ExpectInput("CreateRoute", &ec2.CreateRouteInput{
			RouteTableId:         String("table-id"),
			DestinationCidrBlock: String("10.0.0.0/16"),
			GatewayId:            String("igw-id"),
		}).ExpectCalls("CreateRoute").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete route table=table-id cidr=10.0.0.0/16").
			Mock(&ec2Mock{
				DeleteRouteFunc: func(param0 *ec2.DeleteRouteInput) (*ec2.DeleteRouteOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteRoute", &ec2.DeleteRouteInput{
			RouteTableId:         String("table-id"),
			DestinationCidrBlock: String("10.0.0.0/16"),
		}).ExpectCalls("DeleteRoute").Run(t)
	})

}
