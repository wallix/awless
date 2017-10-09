package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
)

func TestSubscription(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create subscription topic=any-topic endpoint=any-endpoint protocol=HTTP").Mock(&snsMock{
			SubscribeFunc: func(input *sns.SubscribeInput) (*sns.SubscribeOutput, error) {
				return &sns.SubscribeOutput{SubscriptionArn: String("subscription-arn")}, nil
			}}).
			ExpectInput("Subscribe", &sns.SubscribeInput{
				Endpoint: String("any-endpoint"),
				Protocol: String("HTTP"),
				TopicArn: String("any-topic"),
			}).ExpectCommandResult("subscription-arn").ExpectCalls("Subscribe").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete subscription id=any-subscription-arn").Mock(&snsMock{
			UnsubscribeFunc: func(input *sns.UnsubscribeInput) (*sns.UnsubscribeOutput, error) {
				return nil, nil
			}}).
			ExpectInput("Unsubscribe", &sns.UnsubscribeInput{
				SubscriptionArn: String("any-subscription-arn"),
			}).ExpectCalls("Unsubscribe").Run(t)
	})
}
