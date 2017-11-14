package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

func TestTargetgroup(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create targetgroup name=new-tg port=80 protocol=HTTP vpc=any-vpc-id healthcheckinterval=2 healthcheckpath=/health healthcheckport=80 healthcheckprotocol=HTTP healthchecktimeout=180 healthythreshold=30 unhealthythreshold=10 matcher=OK").Mock(&elbv2Mock{
			CreateTargetGroupFunc: func(input *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
				return &elbv2.CreateTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{{TargetGroupArn: String("new-tg-arn")}},
				}, nil
			}}).ExpectInput("CreateTargetGroup", &elbv2.CreateTargetGroupInput{
			Name:     String("new-tg"),
			Port:     Int64(80),
			Protocol: String("HTTP"),
			VpcId:    String("any-vpc-id"),
			HealthCheckIntervalSeconds: Int64(2),
			HealthCheckPath:            String("/health"),
			HealthCheckPort:            String("80"),
			HealthCheckProtocol:        String("HTTP"),
			HealthCheckTimeoutSeconds:  Int64(180),
			HealthyThresholdCount:      Int64(30),
			UnhealthyThresholdCount:    Int64(10),
			Matcher: &elbv2.Matcher{
				HttpCode: String("OK"),
			},
		},
		).ExpectCommandResult("new-tg-arn").ExpectCalls("CreateTargetGroup").Run(t)
	})

	t.Run("update", func(t *testing.T) {
		Template("update targetgroup id=any-tg stickiness=ouech stickinessduration=ouechdur deregistrationdelay=yeap healthcheckinterval=2 healthcheckpath=/health healthcheckport=80 healthcheckprotocol=HTTP healthchecktimeout=180 healthythreshold=30 unhealthythreshold=10 matcher=OK").Mock(&elbv2Mock{
			ModifyTargetGroupAttributesFunc: func(input *elbv2.ModifyTargetGroupAttributesInput) (*elbv2.ModifyTargetGroupAttributesOutput, error) {
				return &elbv2.ModifyTargetGroupAttributesOutput{
					Attributes: []*elbv2.TargetGroupAttribute{},
				}, nil
			},
			ModifyTargetGroupFunc: func(input *elbv2.ModifyTargetGroupInput) (*elbv2.ModifyTargetGroupOutput, error) { return nil, nil }}).ExpectInput("ModifyTargetGroupAttributes", &elbv2.ModifyTargetGroupAttributesInput{
			TargetGroupArn: String("any-tg"),
			Attributes: []*elbv2.TargetGroupAttribute{
				{Key: String("stickiness.enabled"), Value: String("ouech")},
				{Key: String("stickiness.lb_cookie.duration_seconds"), Value: String("ouechdur")},
				{Key: String("deregistration_delay.timeout_seconds"), Value: String("yeap")},
			}}).ExpectInput("ModifyTargetGroup", &elbv2.ModifyTargetGroupInput{
			TargetGroupArn:             String("any-tg"),
			HealthCheckIntervalSeconds: Int64(2),
			HealthCheckPath:            String("/health"),
			HealthCheckPort:            String("80"),
			HealthCheckProtocol:        String("HTTP"),
			HealthCheckTimeoutSeconds:  Int64(180),
			HealthyThresholdCount:      Int64(30),
			UnhealthyThresholdCount:    Int64(10),
			Matcher: &elbv2.Matcher{
				HttpCode: String("OK"),
			},
		}).ExpectCalls("ModifyTargetGroupAttributes", "ModifyTargetGroup").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete targetgroup id=any-tg-arn").Mock(&elbv2Mock{
			DeleteTargetGroupFunc: func(input *elbv2.DeleteTargetGroupInput) (*elbv2.DeleteTargetGroupOutput, error) {
				return &elbv2.DeleteTargetGroupOutput{}, nil
			}}).ExpectInput("DeleteTargetGroup", &elbv2.DeleteTargetGroupInput{
			TargetGroupArn: String("any-tg-arn"),
		}).ExpectCalls("DeleteTargetGroup").Run(t)
	})
}
