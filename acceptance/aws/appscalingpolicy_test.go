package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
)

func TestAppscalingPolicy(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create appscalingpolicy name=my-new-policy-name type=StepScaling resource=service/default/sample-webapp dimension=ecs:service:DesiredCount "+
			"service-namespace=ecs stepscaling-adjustment-type=ChangeInCapacity stepscaling-cooldown=12 stepscaling-aggregation-type=Average stepscaling-min-adjustment-magnitude=2 "+
			"stepscaling-adjustments=75::+1,0:25:-1").
			Mock(&applicationautoscalingMock{
				PutScalingPolicyFunc: func(param0 *applicationautoscaling.PutScalingPolicyInput) (*applicationautoscaling.PutScalingPolicyOutput, error) {
					return &applicationautoscaling.PutScalingPolicyOutput{PolicyARN: String("new-appscalingpolicy-arn")}, nil
				},
			}).ExpectInput("PutScalingPolicy", &applicationautoscaling.PutScalingPolicyInput{
			PolicyName:        String("my-new-policy-name"),
			PolicyType:        String("StepScaling"),
			ResourceId:        String("service/default/sample-webapp"),
			ScalableDimension: String("ecs:service:DesiredCount"),
			ServiceNamespace:  String("ecs"),
			StepScalingPolicyConfiguration: &applicationautoscaling.StepScalingPolicyConfiguration{
				AdjustmentType:         String("ChangeInCapacity"),
				Cooldown:               Int64(12),
				MetricAggregationType:  String("Average"),
				MinAdjustmentMagnitude: Int64(int64(2)),
				StepAdjustments: []*applicationautoscaling.StepAdjustment{
					{
						MetricIntervalLowerBound: Float64(75),
						MetricIntervalUpperBound: nil,
						ScalingAdjustment:        Int64(1),
					},
					{
						MetricIntervalLowerBound: Float64(0),
						MetricIntervalUpperBound: Float64(25),
						ScalingAdjustment:        Int64(-1),
					},
				},
			},
		}).
			ExpectCommandResult("new-appscalingpolicy-arn").ExpectCalls("PutScalingPolicy").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete appscalingpolicy name=my-policy-to-delete resource=service/default/sample-webapp dimension=ecs:service:DesiredCount service-namespace=ecs").
			Mock(&applicationautoscalingMock{
				DeleteScalingPolicyFunc: func(param0 *applicationautoscaling.DeleteScalingPolicyInput) (*applicationautoscaling.DeleteScalingPolicyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteScalingPolicy", &applicationautoscaling.DeleteScalingPolicyInput{
			PolicyName:        String("my-policy-to-delete"),
			ResourceId:        String("service/default/sample-webapp"),
			ScalableDimension: String("ecs:service:DesiredCount"),
			ServiceNamespace:  String("ecs"),
		}).
			ExpectCalls("DeleteScalingPolicy").Run(t)
	})
}
