/* Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsspec

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/wallix/awless/cloud/graph"
	"github.com/wallix/awless/logger"
)

type CreateSubscription struct {
	_        string `action:"create" entity:"subscription" awsAPI:"sns" awsCall:"Subscribe" awsInput:"sns.SubscribeInput" awsOutput:"sns.SubscribeOutput"`
	logger   *logger.Logger
	graph    cloudgraph.GraphAPI
	api      snsiface.SNSAPI
	Topic    *string `awsName:"TopicArn" awsType:"awsstr" templateName:"topic" required:""`
	Endpoint *string `awsName:"Endpoint" awsType:"awsstr" templateName:"endpoint" required:""`
	Protocol *string `awsName:"Protocol" awsType:"awsstr" templateName:"protocol" required:""`
}

func (cmd *CreateSubscription) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateSubscription) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*sns.SubscribeOutput).SubscriptionArn)
}

type DeleteSubscription struct {
	_      string `action:"delete" entity:"subscription" awsAPI:"sns" awsCall:"Unsubscribe" awsInput:"sns.UnsubscribeInput" awsOutput:"sns.UnsubscribeOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    snsiface.SNSAPI
	Id     *string `awsName:"SubscriptionArn" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *DeleteSubscription) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}
