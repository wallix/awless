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
	"fmt"
	"time"

	"github.com/wallix/awless/cloud/graph"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/wallix/awless/logger"
)

type CreateScalinggroup struct {
	_                      string `action:"create" entity:"scalinggroup" awsAPI:"autoscaling" awsCall:"CreateAutoScalingGroup" awsInput:"autoscaling.CreateAutoScalingGroupInput" awsOutput:"autoscaling.CreateAutoScalingGroupOutput"`
	logger                 *logger.Logger
	graph                  cloudgraph.GraphAPI
	api                    autoscalingiface.AutoScalingAPI
	Name                   *string   `awsName:"AutoScalingGroupName" awsType:"awsstr" templateName:"name" required:""`
	Launchconfiguration    *string   `awsName:"LaunchConfigurationName" awsType:"awsstr" templateName:"launchconfiguration" required:""`
	MaxSize                *int64    `awsName:"MaxSize" awsType:"awsint64" templateName:"max-size" required:""`
	MinSize                *int64    `awsName:"MinSize" awsType:"awsint64" templateName:"min-size" required:""`
	Subnets                []*string `awsName:"VPCZoneIdentifier" awsType:"awscsvstr" templateName:"subnets" required:""`
	Cooldown               *int64    `awsName:"DefaultCooldown" awsType:"awsint64" templateName:"cooldown"`
	DesiredCapacity        *int64    `awsName:"DesiredCapacity" awsType:"awsint64" templateName:"desired-capacity"`
	HealthcheckGracePeriod *int64    `awsName:"HealthCheckGracePeriod" awsType:"awsint64" templateName:"healthcheck-grace-period"`
	HealthcheckType        *string   `awsName:"HealthCheckType" awsType:"awsstr" templateName:"healthcheck-type"`
	NewInstancesProtected  *bool     `awsName:"NewInstancesProtectedFromScaleIn" awsType:"awsbool" templateName:"new-instances-protected"`
	Targetgroups           []*string `awsName:"TargetGroupARNs" awsType:"awsstringslice" templateName:"targetgroups"`
}

func (cmd *CreateScalinggroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateScalinggroup) ExtractResult(i interface{}) string {
	return StringValue(cmd.Name)
}

type UpdateScalinggroup struct {
	_                      string `action:"update" entity:"scalinggroup" awsAPI:"autoscaling" awsCall:"UpdateAutoScalingGroup" awsInput:"autoscaling.UpdateAutoScalingGroupInput" awsOutput:"autoscaling.UpdateAutoScalingGroupOutput"`
	logger                 *logger.Logger
	graph                  cloudgraph.GraphAPI
	api                    autoscalingiface.AutoScalingAPI
	Name                   *string   `awsName:"AutoScalingGroupName" awsType:"awsstr" templateName:"name" required:""`
	Cooldown               *int64    `awsName:"DefaultCooldown" awsType:"awsint64" templateName:"cooldown"`
	DesiredCapacity        *int64    `awsName:"DesiredCapacity" awsType:"awsint64" templateName:"desired-capacity"`
	HealthcheckGracePeriod *int64    `awsName:"HealthCheckGracePeriod" awsType:"awsint64" templateName:"healthcheck-grace-period"`
	HealthcheckType        *string   `awsName:"HealthCheckType" awsType:"awsstr" templateName:"healthcheck-type"`
	Launchconfiguration    *string   `awsName:"LaunchConfigurationName" awsType:"awsstr" templateName:"launchconfiguration"`
	MaxSize                *int64    `awsName:"MaxSize" awsType:"awsint64" templateName:"max-size"`
	MinSize                *int64    `awsName:"MinSize" awsType:"awsint64" templateName:"min-size"`
	NewInstancesProtected  *bool     `awsName:"NewInstancesProtectedFromScaleIn" awsType:"awsbool" templateName:"new-instances-protected"`
	Subnets                []*string `awsName:"VPCZoneIdentifier" awsType:"awscsvstr" templateName:"subnets"`
}

func (cmd *UpdateScalinggroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type DeleteScalinggroup struct {
	_      string `action:"delete" entity:"scalinggroup" awsAPI:"autoscaling" awsCall:"DeleteAutoScalingGroup" awsInput:"autoscaling.DeleteAutoScalingGroupInput" awsOutput:"autoscaling.DeleteAutoScalingGroupOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    autoscalingiface.AutoScalingAPI
	Name   *string `awsName:"AutoScalingGroupName" awsType:"awsstr" templateName:"name" required:""`
	Force  *bool   `awsName:"ForceDelete" awsType:"awsbool" templateName:"force"`
}

func (cmd *DeleteScalinggroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type CheckScalinggroup struct {
	_       string `action:"check" entity:"scalinggroup" awsAPI:"autoscaling"`
	logger  *logger.Logger
	graph   cloudgraph.GraphAPI
	api     autoscalingiface.AutoScalingAPI
	Name    *string `templateName:"name" required:""`
	Count   *int64  `templateName:"count" required:""`
	Timeout *int64  `templateName:"timeout" required:""`
}

func (cmd *CheckScalinggroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (sg *CheckScalinggroup) ManualRun(map[string]interface{}) (interface{}, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{sg.Name},
	}

	sgName := StringValue(sg.Name)

	c := &checker{
		description: fmt.Sprintf("scalinggroup '%s'", sgName),
		timeout:     time.Duration(Int64AsIntValue(sg.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		checkName:   "count",
		fetchFunc: func() (string, error) {
			output, err := sg.api.DescribeAutoScalingGroups(input)
			if err != nil {
				return "", err
			}
			for _, group := range output.AutoScalingGroups {
				if StringValue(group.AutoScalingGroupName) == sgName {
					return fmt.Sprint(len(group.Instances)), nil
				}
			}
			return "", fmt.Errorf("scalinggroup %s not found", sgName)
		},
		expect: fmt.Sprint(Int64AsIntValue(sg.Count)),
		logger: sg.logger,
	}
	return nil, c.check()
}
