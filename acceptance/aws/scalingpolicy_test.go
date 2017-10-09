package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

/*
	AdjustmentType      *string `awsName:"AdjustmentType" awsType:"awsstr" templateName:"adjustment-type" required:""`
	Scalinggroup        *string `awsName:"AutoScalingGroupName" awsType:"awsstr" templateName:"scalinggroup" required:""`
	Name                *string `awsName:"PolicyName" awsType:"awsstr" templateName:"name" required:""`
	AdjustmentScaling   *int64  `awsName:"ScalingAdjustment" awsType:"awsint64" templateName:"adjustment-scaling" required:""`
	Cooldown            *int64  `awsName:"Cooldown" awsType:"awsint64" templateName:"cooldown"`
	AdjustmentMagnitude *int64  `awsName:"MinAdjustmentMagnitude" awsType:"awsint64" templateName:"adjustment-magnitude"`

*/

func TestScalingPolicy(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create scalingpolicy adjustment-type=adj scalinggroup=scaler name=any-policy adjustment-scaling=3 cooldown=5 adjustment-magnitude=12").Mock(&autoscalingMock{
			PutScalingPolicyFunc: func(input *autoscaling.PutScalingPolicyInput) (*autoscaling.PutScalingPolicyOutput, error) {
				return &autoscaling.PutScalingPolicyOutput{PolicyARN: String("new-policy-arn")}, nil
			}}).
			ExpectInput("PutScalingPolicy", &autoscaling.PutScalingPolicyInput{
				AdjustmentType:         String("adj"),
				AutoScalingGroupName:   String("scaler"),
				PolicyName:             String("any-policy"),
				Cooldown:               Int64(5),
				ScalingAdjustment:      Int64(3),
				MinAdjustmentMagnitude: Int64(12),
			}).ExpectCommandResult("new-policy-arn").ExpectCalls("PutScalingPolicy").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete scalingpolicy id=any-scalingpolicy").Mock(&autoscalingMock{
			DeletePolicyFunc: func(input *autoscaling.DeletePolicyInput) (*autoscaling.DeletePolicyOutput, error) {
				return nil, nil
			}}).ExpectInput("DeletePolicy", &autoscaling.DeletePolicyInput{
			PolicyName: String("any-scalingpolicy"),
		}).ExpectCalls("DeletePolicy").Run(t)
	})
}
