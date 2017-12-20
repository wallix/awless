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
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateScalingpolicy struct {
	_                   string `action:"create" entity:"scalingpolicy" awsAPI:"autoscaling" awsCall:"PutScalingPolicy" awsInput:"autoscaling.PutScalingPolicyInput" awsOutput:"autoscaling.PutScalingPolicyOutput"`
	logger              *logger.Logger
	graph               cloud.GraphAPI
	api                 autoscalingiface.AutoScalingAPI
	AdjustmentType      *string `awsName:"AdjustmentType" awsType:"awsstr" templateName:"adjustment-type"`
	Scalinggroup        *string `awsName:"AutoScalingGroupName" awsType:"awsstr" templateName:"scalinggroup"`
	Name                *string `awsName:"PolicyName" awsType:"awsstr" templateName:"name"`
	AdjustmentScaling   *int64  `awsName:"ScalingAdjustment" awsType:"awsint64" templateName:"adjustment-scaling"`
	Cooldown            *int64  `awsName:"Cooldown" awsType:"awsint64" templateName:"cooldown"`
	AdjustmentMagnitude *int64  `awsName:"MinAdjustmentMagnitude" awsType:"awsint64" templateName:"adjustment-magnitude"`
}

func (cmd *CreateScalingpolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("adjustment-scaling"), params.Key("adjustment-type"), params.Key("name"), params.Key("scalinggroup"),
		params.Opt("adjustment-magnitude", "cooldown"),
	))
}

func (cmd *CreateScalingpolicy) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*autoscaling.PutScalingPolicyOutput).PolicyARN)
}

type DeleteScalingpolicy struct {
	_      string `action:"delete" entity:"scalingpolicy" awsAPI:"autoscaling" awsCall:"DeletePolicy" awsInput:"autoscaling.DeletePolicyInput" awsOutput:"autoscaling.DeletePolicyOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    autoscalingiface.AutoScalingAPI
	Id     *string `awsName:"PolicyName" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteScalingpolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}
