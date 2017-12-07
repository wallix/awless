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
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/wallix/awless/cloud/graph"
	"github.com/wallix/awless/logger"
)

type CreateLaunchconfiguration struct {
	_              string `action:"create" entity:"launchconfiguration" awsAPI:"autoscaling" awsCall:"CreateLaunchConfiguration" awsInput:"autoscaling.CreateLaunchConfigurationInput" awsOutput:"autoscaling.CreateLaunchConfigurationOutput"`
	logger         *logger.Logger
	graph          cloudgraph.GraphAPI
	api            autoscalingiface.AutoScalingAPI
	Image          *string   `awsName:"ImageId" awsType:"awsstr" templateName:"image" required:""`
	Type           *string   `awsName:"InstanceType" awsType:"awsstr" templateName:"type" required:""`
	Name           *string   `awsName:"LaunchConfigurationName" awsType:"awsstr" templateName:"name" required:""`
	Public         *bool     `awsName:"AssociatePublicIpAddress" awsType:"awsbool" templateName:"public"`
	Keypair        *string   `awsName:"KeyName" awsType:"awsstr" templateName:"keypair"`
	Userdata       *string   `awsName:"UserData" awsType:"awsfiletobase64" templateName:"userdata"`
	Securitygroups []*string `awsName:"SecurityGroups" awsType:"awsstringslice" templateName:"securitygroups"`
	Role           *string   `awsName:"IamInstanceProfile" awsType:"awsstr" templateName:"role"`
	Spotprice      *string   `awsName:"SpotPrice" awsType:"awsstr" templateName:"spotprice"`
	DistroQuery    *string   `awsType:"awsstr" templateName:"distro"`
}

func (cmd *CreateLaunchconfiguration) ValidateParams(params []string) ([]string, error) {
	return paramRule{
		tree:   allOf(oneOf(node("distro"), node("image")), node("type"), node("name")),
		extras: []string{"public", "keypair", "userdata", "securitygroups", "role", "spotprice"},
	}.verify(params)
}

func (cmd *CreateLaunchconfiguration) ConvertParams() ([]string, func(values map[string]interface{}) (map[string]interface{}, error)) {
	createInstance := CommandFactory.Build("createinstance")().(*CreateInstance)
	return []string{"distro"}, createInstance.convertDistroToAMI
}

func (cmd *CreateLaunchconfiguration) ExtractResult(i interface{}) string {
	return StringValue(cmd.Name)
}

type DeleteLaunchconfiguration struct {
	_      string `action:"delete" entity:"launchconfiguration" awsAPI:"autoscaling" awsCall:"DeleteLaunchConfiguration" awsInput:"autoscaling.DeleteLaunchConfigurationInput" awsOutput:"autoscaling.DeleteLaunchConfigurationOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    autoscalingiface.AutoScalingAPI
	Name   *string `awsName:"LaunchConfigurationName" awsType:"awsstr" templateName:"name" required:""`
}

func (cmd *DeleteLaunchconfiguration) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}
