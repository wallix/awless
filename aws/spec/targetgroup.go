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
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/wallix/awless/logger"
)

type CreateTargetgroup struct {
	_                   string `action:"create" entity:"targetgroup" awsAPI:"elbv2" awsCall:"CreateTargetGroup" awsInput:"elbv2.CreateTargetGroupInput" awsOutput:"elbv2.CreateTargetGroupOutput"`
	logger              *logger.Logger
	graph               cloud.GraphAPI
	api                 elbv2iface.ELBV2API
	Name                *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Port                *int64  `awsName:"Port" awsType:"awsint64" templateName:"port"`
	Protocol            *string `awsName:"Protocol" awsType:"awsstr" templateName:"protocol"`
	Vpc                 *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
	Healthcheckinterval *int64  `awsName:"HealthCheckIntervalSeconds" awsType:"awsint64" templateName:"healthcheckinterval"`
	Healthcheckpath     *string `awsName:"HealthCheckPath" awsType:"awsstr" templateName:"healthcheckpath"`
	Healthcheckport     *string `awsName:"HealthCheckPort" awsType:"awsstr" templateName:"healthcheckport"`
	Healthcheckprotocol *string `awsName:"HealthCheckProtocol" awsType:"awsstr" templateName:"healthcheckprotocol"`
	Healthchecktimeout  *int64  `awsName:"HealthCheckTimeoutSeconds" awsType:"awsint64" templateName:"healthchecktimeout"`
	Healthythreshold    *int64  `awsName:"HealthyThresholdCount" awsType:"awsint64" templateName:"healthythreshold"`
	Unhealthythreshold  *int64  `awsName:"UnhealthyThresholdCount" awsType:"awsint64" templateName:"unhealthythreshold"`
	Matcher             *string `awsName:"Matcher.HttpCode" awsType:"awsstr" templateName:"matcher"`
}

func (cmd *CreateTargetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Key("port"), params.Key("protocol"), params.Key("vpc"),
		params.Opt("healthcheckinterval", "healthcheckpath", "healthcheckport", "healthcheckprotocol", "healthchecktimeout", "healthythreshold", "matcher", "unhealthythreshold"),
	))
}

func (cmd *CreateTargetgroup) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*elbv2.CreateTargetGroupOutput).TargetGroups[0].TargetGroupArn)
}

type UpdateTargetgroup struct {
	_                   string `action:"update" entity:"targetgroup" awsAPI:"elbv2"`
	logger              *logger.Logger
	graph               cloud.GraphAPI
	api                 elbv2iface.ELBV2API
	Id                  *string `awsName:"TargetGroupArn" awsType:"awsstr" templateName:"id"`
	Deregistrationdelay *string `awsType:"awsstr" templateName:"deregistrationdelay"`
	Stickiness          *string `awsType:"awsstr" templateName:"stickiness"`
	Stickinessduration  *string `awsType:"awsstr" templateName:"stickinessduration"`
	Healthcheckinterval *int64  `awsName:"HealthCheckIntervalSeconds" awsType:"awsint64" templateName:"healthcheckinterval"`
	Healthcheckpath     *string `awsName:"HealthCheckPath" awsType:"awsstr" templateName:"healthcheckpath"`
	Healthcheckport     *string `awsName:"HealthCheckPort" awsType:"awsstr" templateName:"healthcheckport"`
	Healthcheckprotocol *string `awsName:"HealthCheckProtocol" awsType:"awsstr" templateName:"healthcheckprotocol"`
	Healthchecktimeout  *int64  `awsName:"HealthCheckTimeoutSeconds" awsType:"awsint64" templateName:"healthchecktimeout"`
	Healthythreshold    *int64  `awsName:"HealthyThresholdCount" awsType:"awsint64" templateName:"healthythreshold"`
	Unhealthythreshold  *int64  `awsName:"UnhealthyThresholdCount" awsType:"awsint64" templateName:"unhealthythreshold"`
	Matcher             *string `awsName:"Matcher.HttpCode" awsType:"awsstr" templateName:"matcher"`
}

func (cmd *UpdateTargetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"),
		params.Opt("deregistrationdelay", "healthcheckinterval", "healthcheckpath", "healthcheckport", "healthcheckprotocol", "healthchecktimeout", "healthythreshold", "matcher", "stickiness", "stickinessduration", "unhealthythreshold"),
	))
}

func (tg *UpdateTargetgroup) ManualRun(renv env.Running) (interface{}, error) {
	tgArn := StringValue(tg.Id)

	attrsInput := &elbv2.ModifyTargetGroupAttributesInput{}
	var areTargetAttrsModified bool

	if v := tg.Stickiness; v != nil {
		attrsInput.Attributes = append(attrsInput.Attributes, &elbv2.TargetGroupAttribute{
			Key:   String("stickiness.enabled"),
			Value: v,
		})
		areTargetAttrsModified = true
	}
	if v := tg.Stickinessduration; v != nil {
		attrsInput.Attributes = append(attrsInput.Attributes, &elbv2.TargetGroupAttribute{
			Key:   String("stickiness.lb_cookie.duration_seconds"),
			Value: v,
		})
		areTargetAttrsModified = true
	}
	if v := tg.Deregistrationdelay; v != nil {
		attrsInput.Attributes = append(attrsInput.Attributes, &elbv2.TargetGroupAttribute{
			Key:   String("deregistration_delay.timeout_seconds"),
			Value: v,
		})
		areTargetAttrsModified = true
	}

	var err error

	if areTargetAttrsModified {
		if err = setFieldWithType(tgArn, attrsInput, "TargetGroupArn", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		start := time.Now()
		if _, err = tg.api.ModifyTargetGroupAttributes(attrsInput); err != nil {
			return nil, err
		}
		tg.logger.ExtraVerbosef("elbv2.ModifyTargetGroupAttributes call took %s", time.Since(start))
	}

	input := &elbv2.ModifyTargetGroupInput{}
	var isTargetGroupModified bool

	if v := tg.Healthcheckinterval; v != nil {
		if err = setFieldWithType(v, input, "HealthCheckIntervalSeconds", awsint64, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := tg.Healthcheckpath; v != nil {
		if err = setFieldWithType(v, input, "HealthCheckPath", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := tg.Healthcheckport; v != nil {
		if err = setFieldWithType(v, input, "HealthCheckPort", awsstr, renv.Context()); err != nil {
			return nil, err
		}
	}
	if v := tg.Healthcheckprotocol; v != nil {
		if err = setFieldWithType(v, input, "HealthCheckProtocol", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := tg.Healthchecktimeout; v != nil {
		if err = setFieldWithType(v, input, "HealthCheckTimeoutSeconds", awsint64, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := tg.Healthythreshold; v != nil {
		if err = setFieldWithType(v, input, "HealthyThresholdCount", awsint64, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := tg.Unhealthythreshold; v != nil {
		if err = setFieldWithType(v, input, "UnhealthyThresholdCount", awsint64, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := tg.Matcher; v != nil {
		if err = setFieldWithType(v, input, "Matcher.HttpCode", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}

	if isTargetGroupModified {
		if err = setFieldWithType(tgArn, input, "TargetGroupArn", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		start := time.Now()
		output, err := tg.api.ModifyTargetGroup(input)
		tg.logger.ExtraVerbosef("elbv2.ModifyTargetGroup call took %s", time.Since(start))
		return output, err
	}
	return nil, nil
}

type DeleteTargetgroup struct {
	_      string `action:"delete" entity:"targetgroup" awsAPI:"elbv2" awsCall:"DeleteTargetGroup" awsInput:"elbv2.DeleteTargetGroupInput" awsOutput:"elbv2.DeleteTargetGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    elbv2iface.ELBV2API
	Id     *string `awsName:"TargetGroupArn" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteTargetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}
