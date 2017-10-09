package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestNetworkInterface(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create networkinterface subnet=sub-1234 description='my ni desc' securitygroups=sg-1234,sg-2345 privateip=127.0.0.1").
			Mock(&ec2Mock{
				CreateNetworkInterfaceFunc: func(param0 *ec2.CreateNetworkInterfaceInput) (*ec2.CreateNetworkInterfaceOutput, error) {
					return &ec2.CreateNetworkInterfaceOutput{NetworkInterface: &ec2.NetworkInterface{NetworkInterfaceId: String("new-networkinterface-id")}}, nil
				},
			}).ExpectInput("CreateNetworkInterface", &ec2.CreateNetworkInterfaceInput{
			SubnetId:         String("sub-1234"),
			Description:      String("my ni desc"),
			Groups:           []*string{String("sg-1234"), String("sg-2345")},
			PrivateIpAddress: String("127.0.0.1"),
		}).
			ExpectCommandResult("new-networkinterface-id").ExpectCalls("CreateNetworkInterface").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete networkinterface id=ni-1234").
			Mock(&ec2Mock{
				DeleteNetworkInterfaceFunc: func(param0 *ec2.DeleteNetworkInterfaceInput) (*ec2.DeleteNetworkInterfaceOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteNetworkInterface", &ec2.DeleteNetworkInterfaceInput{NetworkInterfaceId: String("ni-1234")}).
			ExpectCalls("DeleteNetworkInterface").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach networkinterface id=ni-1234 instance=i-2345 device-index=2").
			Mock(&ec2Mock{
				AttachNetworkInterfaceFunc: func(param0 *ec2.AttachNetworkInterfaceInput) (*ec2.AttachNetworkInterfaceOutput, error) {
					return &ec2.AttachNetworkInterfaceOutput{AttachmentId: String("attach-ni-id")}, nil
				},
			}).ExpectInput("AttachNetworkInterface", &ec2.AttachNetworkInterfaceInput{
			NetworkInterfaceId: String("ni-1234"),
			InstanceId:         String("i-2345"),
			DeviceIndex:        Int64(2),
		}).ExpectCommandResult("attach-ni-id").
			ExpectCalls("AttachNetworkInterface").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		t.Run("with attachment id", func(t *testing.T) {
			Template("detach networkinterface attachment=id-of-attachment").
				Mock(&ec2Mock{
					DetachNetworkInterfaceFunc: func(param0 *ec2.DetachNetworkInterfaceInput) (*ec2.DetachNetworkInterfaceOutput, error) {
						return nil, nil
					},
				}).ExpectInput("DetachNetworkInterface", &ec2.DetachNetworkInterfaceInput{
				AttachmentId: String("id-of-attachment"),
			}).
				ExpectCalls("DetachNetworkInterface").Run(t)
		})
		t.Run("with instance and network interface", func(t *testing.T) {
			Template("detach networkinterface id=ni-1234 instance=i-2345 force=true").
				Mock(&ec2Mock{
					DescribeInstancesFunc: func(param0 *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
						return &ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{{
							Instances: []*ec2.Instance{
								{NetworkInterfaces: []*ec2.InstanceNetworkInterface{
									{NetworkInterfaceId: String("ni-1234"), Attachment: &ec2.InstanceNetworkInterfaceAttachment{AttachmentId: String("my-attachment-id")}},
								}},
							},
						}}}, nil
					},
					DetachNetworkInterfaceFunc: func(param0 *ec2.DetachNetworkInterfaceInput) (*ec2.DetachNetworkInterfaceOutput, error) {
						return nil, nil
					},
				}).ExpectInput("DescribeInstances", &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					{Name: String("network-interface.network-interface-id"), Values: []*string{String("ni-1234")}},
					{Name: String("instance-id"), Values: []*string{String("i-2345")}},
				},
			}).ExpectInput("DetachNetworkInterface", &ec2.DetachNetworkInterfaceInput{
				AttachmentId: String("my-attachment-id"),
				Force:        Bool(true),
			}).
				ExpectCalls("DescribeInstances", "DetachNetworkInterface").Run(t)
		})
	})
	t.Run("check", func(t *testing.T) {
		Template("check networkinterface id=ni-1234 state=available timeout=1").
			Mock(&ec2Mock{
				DescribeNetworkInterfacesFunc: func(param0 *ec2.DescribeNetworkInterfacesInput) (*ec2.DescribeNetworkInterfacesOutput, error) {
					return &ec2.DescribeNetworkInterfacesOutput{
						NetworkInterfaces: []*ec2.NetworkInterface{
							{NetworkInterfaceId: String("ni-1234"), Status: String("available")},
						},
					}, nil
				},
			}).ExpectInput("DescribeNetworkInterfaces", &ec2.DescribeNetworkInterfacesInput{
			NetworkInterfaceIds: []*string{String("ni-1234")},
		}).ExpectCalls("DescribeNetworkInterfaces").Run(t)
	})
}
