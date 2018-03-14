package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

func TestListener(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create listener actiontype=forward loadbalancer=arn:of:loadbalancer port=80 "+
			"protocol=HTTP targetgroup=arn:of:targetgroup certificate=arn:of:certificate sslpolicy=ELBSecurityPolicy-2016-08").Mock(&elbv2Mock{
			CreateListenerFunc: func(input *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
				return &elbv2.CreateListenerOutput{Listeners: []*elbv2.Listener{
					{ListenerArn: String("arn:of:new:listener")},
				}}, nil
			}}).
			ExpectInput("CreateListener", &elbv2.CreateListenerInput{
				DefaultActions: []*elbv2.Action{
					{Type: String("forward"), TargetGroupArn: String("arn:of:targetgroup")},
				},
				LoadBalancerArn: String("arn:of:loadbalancer"),
				Port:            Int64(80),
				Protocol:        String("HTTP"),
				Certificates:    []*elbv2.Certificate{{CertificateArn: String("arn:of:certificate")}},
				SslPolicy:       String("ELBSecurityPolicy-2016-08"),
			}).ExpectCalls("CreateListener").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach listener id=arn:listener certificate=arn:certificate:to:attach").
			Mock(&elbv2Mock{
				AddListenerCertificatesFunc: func(input *elbv2.AddListenerCertificatesInput) (*elbv2.AddListenerCertificatesOutput, error) {
					return nil, nil
				},
			}).ExpectInput("AddListenerCertificates", &elbv2.AddListenerCertificatesInput{
			Certificates: []*elbv2.Certificate{
				{CertificateArn: String("arn:certificate:to:attach")},
			},
			ListenerArn: String("arn:listener"),
		}).ExpectCalls("AddListenerCertificates").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete listener id=arn:of:listener:to:delete").Mock(&elbv2Mock{
			DeleteListenerFunc: func(input *elbv2.DeleteListenerInput) (*elbv2.DeleteListenerOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteListener", &elbv2.DeleteListenerInput{
				ListenerArn: String("arn:of:listener:to:delete"),
			}).ExpectCalls("DeleteListener").Run(t)
	})

}
