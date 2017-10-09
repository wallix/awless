package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func TestQueue(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create queue name=my-new-queue-name delay=12 max-msg-size=42 retention-period=10 policy=my-policy "+
			"msg-wait=12 redrive-policy=my-redrive-policy visibility-timeout=180").
			Mock(&sqsMock{
				CreateQueueFunc: func(param0 *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
					return &sqs.CreateQueueOutput{QueueUrl: String("my-queue-url")}, nil
				},
			}).ExpectInput("CreateQueue", &sqs.CreateQueueInput{
			QueueName: String("my-new-queue-name"),
			Attributes: map[string]*string{
				"DelaySeconds":                  String("12"),
				"MaximumMessageSize":            String("42"),
				"MessageRetentionPeriod":        String("10"),
				"Policy":                        String("my-policy"),
				"ReceiveMessageWaitTimeSeconds": String("12"),
				"RedrivePolicy":                 String("my-redrive-policy"),
				"VisibilityTimeout":             String("180"),
			},
		}).
			ExpectCommandResult("my-queue-url").ExpectCalls("CreateQueue").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete queue url=queue-url-to-delete").
			Mock(&sqsMock{
				DeleteQueueFunc: func(param0 *sqs.DeleteQueueInput) (*sqs.DeleteQueueOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteQueue", &sqs.DeleteQueueInput{QueueUrl: String("queue-url-to-delete")}).
			ExpectCalls("DeleteQueue").Run(t)
	})
}
