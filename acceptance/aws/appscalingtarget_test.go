package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
)

func TestAppscalingtarget(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create appscalingtarget resource=service/default/sample-webapp service-namespace=ecs dimension=ecs:service:DesiredCount role=arn:of:my-scaling-role min-capacity=2 max-capacity=11").
			Mock(&applicationautoscalingMock{
				RegisterScalableTargetFunc: func(param0 *applicationautoscaling.RegisterScalableTargetInput) (*applicationautoscaling.RegisterScalableTargetOutput, error) {
					return nil, nil
				},
			}).ExpectInput("RegisterScalableTarget", &applicationautoscaling.RegisterScalableTargetInput{
			MaxCapacity:       Int64(11),
			MinCapacity:       Int64(2),
			ResourceId:        String("service/default/sample-webapp"),
			RoleARN:           String("arn:of:my-scaling-role"),
			ScalableDimension: String("ecs:service:DesiredCount"),
			ServiceNamespace:  String("ecs"),
		}).ExpectCalls("RegisterScalableTarget").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete appscalingtarget resource=service/default/sample-webapp service-namespace=ecs dimension=ecs:service:DesiredCount").
			Mock(&applicationautoscalingMock{
				DeregisterScalableTargetFunc: func(param0 *applicationautoscaling.DeregisterScalableTargetInput) (*applicationautoscaling.DeregisterScalableTargetOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeregisterScalableTarget", &applicationautoscaling.DeregisterScalableTargetInput{
			ResourceId:        String("service/default/sample-webapp"),
			ScalableDimension: String("ecs:service:DesiredCount"),
			ServiceNamespace:  String("ecs"),
		}).ExpectCalls("DeregisterScalableTarget").Run(t)
	})

}
