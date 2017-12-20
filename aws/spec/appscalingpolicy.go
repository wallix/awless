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
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateAppscalingpolicy struct {
	_                                 string `action:"create" entity:"appscalingpolicy" awsAPI:"applicationautoscaling" awsCall:"PutScalingPolicy" awsInput:"applicationautoscaling.PutScalingPolicyInput" awsOutput:"applicationautoscaling.PutScalingPolicyOutput"`
	logger                            *logger.Logger
	graph                             cloud.GraphAPI
	api                               applicationautoscalingiface.ApplicationAutoScalingAPI
	Name                              *string   `awsName:"PolicyName" awsType:"awsstr" templateName:"name"`
	Type                              *string   `awsName:"PolicyType" awsType:"awsstr" templateName:"type"`
	Resource                          *string   `awsName:"ResourceId" awsType:"awsstr" templateName:"resource"`
	Dimension                         *string   `awsName:"ScalableDimension" awsType:"awsstr" templateName:"dimension"`
	ServiceNamespace                  *string   `awsName:"ServiceNamespace" awsType:"awsstr" templateName:"service-namespace"`
	StepscalingAdjustmentType         *string   `awsName:"StepScalingPolicyConfiguration.AdjustmentType" awsType:"awsstr" templateName:"stepscaling-adjustment-type"`
	StepscalingAdjustments            []*string `awsName:"StepScalingPolicyConfiguration.StepAdjustments" awsType:"awsstepadjustments" templateName:"stepscaling-adjustments"`
	StepscalingCooldown               *int64    `awsName:"StepScalingPolicyConfiguration.Cooldown" awsType:"awsint64" templateName:"stepscaling-cooldown"`
	StepscalingAggregationType        *string   `awsName:"StepScalingPolicyConfiguration.MetricAggregationType" awsType:"awsstr" templateName:"stepscaling-aggregation-type"`
	StepscalingMinAdjustmentMagnitude *int64    `awsName:"StepScalingPolicyConfiguration.MinAdjustmentMagnitude" awsType:"awsint64" templateName:"stepscaling-min-adjustment-magnitude"`
}

func (cmd *CreateAppscalingpolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("dimension"), params.Key("name"), params.Key("resource"), params.Key("service-namespace"), params.Key("stepscaling-adjustment-type"), params.Key("stepscaling-adjustments"), params.Key("type"),
		params.Opt("stepscaling-aggregation-type", "stepscaling-cooldown", "stepscaling-min-adjustment-magnitude"),
	))
}

func (cmd *CreateAppscalingpolicy) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*applicationautoscaling.PutScalingPolicyOutput).PolicyARN)
}

type DeleteAppscalingpolicy struct {
	_                string `action:"delete" entity:"appscalingpolicy" awsAPI:"applicationautoscaling" awsCall:"DeleteScalingPolicy" awsInput:"applicationautoscaling.DeleteScalingPolicyInput" awsOutput:"applicationautoscaling.DeleteScalingPolicyOutput"`
	logger           *logger.Logger
	graph            cloud.GraphAPI
	api              applicationautoscalingiface.ApplicationAutoScalingAPI
	Name             *string `awsName:"PolicyName" awsType:"awsstr" templateName:"name"`
	Resource         *string `awsName:"ResourceId" awsType:"awsstr" templateName:"resource"`
	Dimension        *string `awsName:"ScalableDimension" awsType:"awsstr" templateName:"dimension"`
	ServiceNamespace *string `awsName:"ServiceNamespace" awsType:"awsstr" templateName:"service-namespace"`
}

func (cmd *DeleteAppscalingpolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("dimension"), params.Key("name"), params.Key("resource"), params.Key("service-namespace")))
}
