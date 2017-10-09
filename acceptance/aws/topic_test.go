package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
)

func TestTopic(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create topic name=donald").Mock(&snsMock{
			CreateTopicFunc: func(input *sns.CreateTopicInput) (*sns.CreateTopicOutput, error) {
				return &sns.CreateTopicOutput{TopicArn: String("new-topic-arn")}, nil
			}}).ExpectInput("CreateTopic", &sns.CreateTopicInput{Name: String("donald")}).
			ExpectCommandResult("new-topic-arn").ExpectCalls("CreateTopic").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete topic id=any-topic-id").Mock(&snsMock{
			DeleteTopicFunc: func(input *sns.DeleteTopicInput) (*sns.DeleteTopicOutput, error) {
				return nil, nil
			}}).ExpectInput("DeleteTopic", &sns.DeleteTopicInput{
			TopicArn: String("any-topic-id"),
		}).ExpectCalls("DeleteTopic").Run(t)
	})
}
