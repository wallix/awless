package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestElasticip(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create elasticip domain=vpc").
			Mock(&ec2Mock{
				AllocateAddressFunc: func(param0 *ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
					return &ec2.AllocateAddressOutput{AllocationId: String("new-elasticip-allocation-id")}, nil
				},
			}).ExpectInput("AllocateAddress", &ec2.AllocateAddressInput{
			Domain: String("vpc"),
		}).
			ExpectCommandResult("new-elasticip-allocation-id").ExpectCalls("AllocateAddress").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by id", func(t *testing.T) {
			Template("delete elasticip id=eipalloc-0123456").
				Mock(&ec2Mock{
					ReleaseAddressFunc: func(param0 *ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
						return nil, nil
					},
				}).ExpectInput("ReleaseAddress", &ec2.ReleaseAddressInput{AllocationId: String("eipalloc-0123456")}).
				ExpectCalls("ReleaseAddress").Run(t)
		})
		t.Run("by ip", func(t *testing.T) {
			Template("delete elasticip ip=127.0.0.1").
				Mock(&ec2Mock{
					ReleaseAddressFunc: func(param0 *ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
						return nil, nil
					},
				}).ExpectInput("ReleaseAddress", &ec2.ReleaseAddressInput{PublicIp: String("127.0.0.1")}).
				ExpectCalls("ReleaseAddress").Run(t)
		})
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach elasticip id=eipalloc-0123456 instance=i-1234 networkinterface=eni-2345 privateip=10.0.0.42 allow-reassociation=true").
			Mock(&ec2Mock{
				AssociateAddressFunc: func(param0 *ec2.AssociateAddressInput) (*ec2.AssociateAddressOutput, error) {
					return &ec2.AssociateAddressOutput{AssociationId: String("ip-assoc-id")}, nil
				},
			}).ExpectInput("AssociateAddress", &ec2.AssociateAddressInput{
			AllocationId:       String("eipalloc-0123456"),
			InstanceId:         String("i-1234"),
			NetworkInterfaceId: String("eni-2345"),
			PrivateIpAddress:   String("10.0.0.42"),
			AllowReassociation: Bool(true),
		}).ExpectCommandResult("ip-assoc-id").ExpectCalls("AssociateAddress").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach elasticip association=ipassoc-12345").
			Mock(&ec2Mock{
				DisassociateAddressFunc: func(param0 *ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DisassociateAddress", &ec2.DisassociateAddressInput{
			AssociationId: String("ipassoc-12345"),
		}).
			ExpectCalls("DisassociateAddress").Run(t)
	})

}
