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
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateAppscalingtarget struct {
	_                string `action:"create" entity:"appscalingtarget" awsAPI:"applicationautoscaling" awsCall:"RegisterScalableTarget" awsInput:"applicationautoscaling.RegisterScalableTargetInput" awsOutput:"applicationautoscaling.RegisterScalableTargetOutput"`
	logger           *logger.Logger
	graph            cloud.GraphAPI
	api              applicationautoscalingiface.ApplicationAutoScalingAPI
	MaxCapacity      *int64  `awsName:"MaxCapacity" awsType:"awsint64" templateName:"max-capacity"`
	MinCapacity      *int64  `awsName:"MinCapacity" awsType:"awsint64" templateName:"min-capacity"`
	Resource         *string `awsName:"ResourceId" awsType:"awsstr" templateName:"resource"`
	Role             *string `awsName:"RoleARN" awsType:"awsstr" templateName:"role"`
	Dimension        *string `awsName:"ScalableDimension" awsType:"awsstr" templateName:"dimension"`
	ServiceNamespace *string `awsName:"ServiceNamespace" awsType:"awsstr" templateName:"service-namespace"`
}

func (cmd *CreateAppscalingtarget) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("dimension"), params.Key("max-capacity"), params.Key("min-capacity"), params.Key("resource"), params.Key("role"), params.Key("service-namespace")))
}

type DeleteAppscalingtarget struct {
	_                string `action:"delete" entity:"appscalingtarget" awsAPI:"applicationautoscaling" awsCall:"DeregisterScalableTarget" awsInput:"applicationautoscaling.DeregisterScalableTargetInput" awsOutput:"applicationautoscaling.DeregisterScalableTargetOutput"`
	logger           *logger.Logger
	graph            cloud.GraphAPI
	api              applicationautoscalingiface.ApplicationAutoScalingAPI
	Resource         *string `awsName:"ResourceId" awsType:"awsstr" templateName:"resource"`
	Dimension        *string `awsName:"ScalableDimension" awsType:"awsstr" templateName:"dimension"`
	ServiceNamespace *string `awsName:"ServiceNamespace" awsType:"awsstr" templateName:"service-namespace"`
}

func (cmd *DeleteAppscalingtarget) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("dimension"), params.Key("resource"), params.Key("service-namespace")))
}
