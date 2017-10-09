package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
)

/*
	Callerreference *string `awsName:"CallerReference" awsType:"awsstr" templateName:"callerreference" required:""`
	Name            *string `awsName:"Name" awsType:"awsstr" templateName:"name" required:""`
	Delegationsetid *string `awsName:"DelegationSetId" awsType:"awsstr" templateName:"delegationsetid"`
	Comment         *string `awsName:"HostedZoneConfig.Comment" awsType:"awsstr" templateName:"comment"`
	Isprivate       *bool   `awsName:"HostedZoneConfig.PrivateZone" awsType:"awsbool" templateName:"isprivate"`
	Vpcid           *string `awsName:"VPC.VPCId" awsType:"awsstr" templateName:"vpcid"`
	Vpcregion       *string `awsName:"VPC.VPCRegion" awsType:"awsstr" templateName:"vpcregion"`

*/

func TestZone(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template(`create zone callerreference=caller name=new-zone delegationsetid=1234 comment="new zone" isprivate=true vpcid=any-vpc vpcregion=us-west-2`).Mock(&route53Mock{
			CreateHostedZoneFunc: func(input *route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
				return &route53.CreateHostedZoneOutput{
					HostedZone: &route53.HostedZone{Id: String("new-zone-id")},
				}, nil
			},
		}).ExpectInput("CreateHostedZone", &route53.CreateHostedZoneInput{
			CallerReference: String("caller"),
			DelegationSetId: String("1234"),
			HostedZoneConfig: &route53.HostedZoneConfig{
				Comment:     String("new zone"),
				PrivateZone: Bool(true),
			},
			Name: String("new-zone"),
			VPC:  &route53.VPC{VPCId: String("any-vpc"), VPCRegion: String("us-west-2")},
		}).ExpectCalls("CreateHostedZone").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete zone id=any-zone-id").Mock(&route53Mock{
			DeleteHostedZoneFunc: func(input *route53.DeleteHostedZoneInput) (*route53.DeleteHostedZoneOutput, error) {
				return nil, nil
			},
		}).ExpectInput("DeleteHostedZone", &route53.DeleteHostedZoneInput{
			Id: String("any-zone-id"),
		}).ExpectCalls("DeleteHostedZone").Run(t)
	})
}
