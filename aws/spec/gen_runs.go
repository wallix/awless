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

// DO NOT EDIT
// This file was automatically generated with go generate
package awsspec

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/env"
)

func NewAttachAlarm(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachAlarm {
	cmd := new(AttachAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *AttachAlarm) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachAlarm) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachAlarm) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *AttachAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachClassicLoadbalancer(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachClassicLoadbalancer {
	cmd := new(AttachClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elb.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachClassicLoadbalancer) SetApi(api elbiface.ELBAPI) {
	cmd.api = api
}

func (cmd *AttachClassicLoadbalancer) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachClassicLoadbalancer) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.RegisterInstancesWithLoadBalancerInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.RegisterInstancesWithLoadBalancerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RegisterInstancesWithLoadBalancer(input)
	renv.Log().ExtraVerbosef("elb.RegisterInstancesWithLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachClassicLoadbalancer) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("classicloadbalancer"), nil
}

func (cmd *AttachClassicLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachContainertask(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachContainertask {
	cmd := new(AttachContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *AttachContainertask) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachContainertask) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachContainertask) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *AttachContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachElasticip(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachElasticip {
	cmd := new(AttachElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachElasticip) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachElasticip) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachElasticip) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AssociateAddressInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateAddressInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AssociateAddress(input)
	renv.Log().ExtraVerbosef("ec2.AssociateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach elasticip '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachElasticip) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.AssociateAddressInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateAddressInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AssociateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AssociateAddress call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach elasticip ok")
			return fakeDryRunId("elasticip"), nil
		}
	}

	return nil, err
}

func (cmd *AttachElasticip) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachInstance {
	cmd := new(AttachInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachInstance) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *AttachInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.RegisterTargetsInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.RegisterTargetsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RegisterTargets(input)
	renv.Log().ExtraVerbosef("elbv2.RegisterTargets call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instance"), nil
}

func (cmd *AttachInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachInstanceprofile(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachInstanceprofile {
	cmd := new(AttachInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachInstanceprofile) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachInstanceprofile) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachInstanceprofile) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach instanceprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInstanceprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachInternetgateway(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachInternetgateway {
	cmd := new(AttachInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachInternetgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachInternetgateway) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachInternetgateway) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachInternetGatewayInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachInternetGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AttachInternetGateway(input)
	renv.Log().ExtraVerbosef("ec2.AttachInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach internetgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInternetgateway) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.AttachInternetGatewayInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachInternetGatewayInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AttachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AttachInternetGateway call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach internetgateway ok")
			return fakeDryRunId("internetgateway"), nil
		}
	}

	return nil, err
}

func (cmd *AttachInternetgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachListener(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachListener {
	cmd := new(AttachListener)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachListener) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *AttachListener) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachListener) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.AddListenerCertificatesInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.AddListenerCertificatesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AddListenerCertificates(input)
	renv.Log().ExtraVerbosef("elbv2.AddListenerCertificates call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach listener: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach listener '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach listener done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachListener) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("listener"), nil
}

func (cmd *AttachListener) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachMfadevice(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachMfadevice {
	cmd := new(AttachMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachMfadevice) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *AttachMfadevice) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachMfadevice) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.EnableMFADeviceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.EnableMFADeviceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.EnableMFADevice(input)
	renv.Log().ExtraVerbosef("iam.EnableMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach mfadevice '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachMfadevice) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("mfadevice"), nil
}

func (cmd *AttachMfadevice) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachNetworkinterface(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachNetworkinterface {
	cmd := new(AttachNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachNetworkinterface) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachNetworkinterface) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachNetworkInterfaceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachNetworkInterfaceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AttachNetworkInterface(input)
	renv.Log().ExtraVerbosef("ec2.AttachNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachNetworkinterface) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.AttachNetworkInterfaceInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachNetworkInterfaceInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AttachNetworkInterface(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AttachNetworkInterface call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach networkinterface ok")
			return fakeDryRunId("networkinterface"), nil
		}
	}

	return nil, err
}

func (cmd *AttachNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachPolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachPolicy {
	cmd := new(AttachPolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachPolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *AttachPolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachPolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachPolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *AttachPolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachRole(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachRole {
	cmd := new(AttachRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachRole) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *AttachRole) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachRole) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.AddRoleToInstanceProfileInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.AddRoleToInstanceProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AddRoleToInstanceProfile(input)
	renv.Log().ExtraVerbosef("iam.AddRoleToInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach role '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachRole) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("role"), nil
}

func (cmd *AttachRole) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachRoutetable(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachRoutetable {
	cmd := new(AttachRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachRoutetable) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachRoutetable) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachRoutetable) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AssociateRouteTableInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateRouteTableInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AssociateRouteTable(input)
	renv.Log().ExtraVerbosef("ec2.AssociateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach routetable '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachRoutetable) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.AssociateRouteTableInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateRouteTableInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AssociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AssociateRouteTable call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach routetable ok")
			return fakeDryRunId("routetable"), nil
		}
	}

	return nil, err
}

func (cmd *AttachRoutetable) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachSecuritygroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachSecuritygroup {
	cmd := new(AttachSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachSecuritygroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachSecuritygroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachSecuritygroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("securitygroup"), nil
}

func (cmd *AttachSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachUser(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachUser {
	cmd := new(AttachUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachUser) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *AttachUser) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachUser) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.AddUserToGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.AddUserToGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AddUserToGroup(input)
	renv.Log().ExtraVerbosef("iam.AddUserToGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach user '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachUser) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("user"), nil
}

func (cmd *AttachUser) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachVolume(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AttachVolume {
	cmd := new(AttachVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachVolume) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachVolume) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachVolumeInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachVolumeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AttachVolume(input)
	renv.Log().ExtraVerbosef("ec2.AttachVolume call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachVolume) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.AttachVolumeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachVolumeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AttachVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AttachVolume call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach volume ok")
			return fakeDryRunId("volume"), nil
		}
	}

	return nil, err
}

func (cmd *AttachVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAuthenticateRegistry(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *AuthenticateRegistry {
	cmd := new(AuthenticateRegistry)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecr.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AuthenticateRegistry) SetApi(api ecriface.ECRAPI) {
	cmd.api = api
}

func (cmd *AuthenticateRegistry) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AuthenticateRegistry) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("authenticate registry: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("authenticate registry '%s' done", extracted)
	} else {
		renv.Log().Verbose("authenticate registry done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AuthenticateRegistry) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("registry"), nil
}

func (cmd *AuthenticateRegistry) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckCertificate(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckCertificate {
	cmd := new(CheckCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = acm.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckCertificate) SetApi(api acmiface.ACMAPI) {
	cmd.api = api
}

func (cmd *CheckCertificate) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckCertificate) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check certificate '%s' done", extracted)
	} else {
		renv.Log().Verbose("check certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckCertificate) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("certificate"), nil
}

func (cmd *CheckCertificate) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckDatabase(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckDatabase {
	cmd := new(CheckDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *CheckDatabase) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckDatabase) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check database '%s' done", extracted)
	} else {
		renv.Log().Verbose("check database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckDatabase) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *CheckDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckDistribution(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckDistribution {
	cmd := new(CheckDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudfront.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckDistribution) SetApi(api cloudfrontiface.CloudFrontAPI) {
	cmd.api = api
}

func (cmd *CheckDistribution) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckDistribution) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check distribution '%s' done", extracted)
	} else {
		renv.Log().Verbose("check distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckDistribution) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("distribution"), nil
}

func (cmd *CheckDistribution) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckInstance {
	cmd := new(CheckInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("check instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instance"), nil
}

func (cmd *CheckInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckLoadbalancer(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckLoadbalancer {
	cmd := new(CheckLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckLoadbalancer) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *CheckLoadbalancer) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckLoadbalancer) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check loadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("check loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckLoadbalancer) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loadbalancer"), nil
}

func (cmd *CheckLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckNatgateway(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckNatgateway {
	cmd := new(CheckNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckNatgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckNatgateway) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckNatgateway) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check natgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("check natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckNatgateway) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("natgateway"), nil
}

func (cmd *CheckNatgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckNetworkinterface(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckNetworkinterface {
	cmd := new(CheckNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckNetworkinterface) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckNetworkinterface) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("check networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckNetworkinterface) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("networkinterface"), nil
}

func (cmd *CheckNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckScalinggroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckScalinggroup {
	cmd := new(CheckScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckScalinggroup) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *CheckScalinggroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckScalinggroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check scalinggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("check scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckScalinggroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalinggroup"), nil
}

func (cmd *CheckScalinggroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckSecuritygroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckSecuritygroup {
	cmd := new(CheckSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckSecuritygroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckSecuritygroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("check securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckSecuritygroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("securitygroup"), nil
}

func (cmd *CheckSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckVolume(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CheckVolume {
	cmd := new(CheckVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckVolume) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckVolume) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("check volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckVolume) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("volume"), nil
}

func (cmd *CheckVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCopyImage(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CopyImage {
	cmd := new(CopyImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CopyImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CopyImage) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CopyImage) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CopyImageInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopyImageInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CopyImage(input)
	renv.Log().ExtraVerbosef("ec2.CopyImage call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("copy image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("copy image '%s' done", extracted)
	} else {
		renv.Log().Verbose("copy image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CopyImage) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CopyImageInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopyImageInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CopyImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CopyImage call took %s", time.Since(start))
			renv.Log().Verbose("dry run: copy image ok")
			return fakeDryRunId("image"), nil
		}
	}

	return nil, err
}

func (cmd *CopyImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCopySnapshot(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CopySnapshot {
	cmd := new(CopySnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CopySnapshot) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CopySnapshot) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CopySnapshot) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CopySnapshotInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopySnapshotInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CopySnapshot(input)
	renv.Log().ExtraVerbosef("ec2.CopySnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("copy snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("copy snapshot '%s' done", extracted)
	} else {
		renv.Log().Verbose("copy snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CopySnapshot) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CopySnapshotInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopySnapshotInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CopySnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CopySnapshot call took %s", time.Since(start))
			renv.Log().Verbose("dry run: copy snapshot ok")
			return fakeDryRunId("snapshot"), nil
		}
	}

	return nil, err
}

func (cmd *CopySnapshot) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateAccesskey(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateAccesskey {
	cmd := new(CreateAccesskey)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateAccesskey) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateAccesskey) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateAccesskey) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateAccessKeyInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateAccessKeyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateAccessKey(input)
	renv.Log().ExtraVerbosef("iam.CreateAccessKey call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create accesskey: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create accesskey '%s' done", extracted)
	} else {
		renv.Log().Verbose("create accesskey done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAccesskey) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("accesskey"), nil
}

func (cmd *CreateAccesskey) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateAlarm(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateAlarm {
	cmd := new(CreateAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *CreateAlarm) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateAlarm) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.PutMetricAlarmInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.PutMetricAlarmInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.PutMetricAlarm(input)
	renv.Log().ExtraVerbosef("cloudwatch.PutMetricAlarm call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("create alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAlarm) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *CreateAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateAppscalingpolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateAppscalingpolicy {
	cmd := new(CreateAppscalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = applicationautoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateAppscalingpolicy) SetApi(api applicationautoscalingiface.ApplicationAutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateAppscalingpolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateAppscalingpolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.PutScalingPolicyInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.PutScalingPolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.PutScalingPolicy(input)
	renv.Log().ExtraVerbosef("applicationautoscaling.PutScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create appscalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create appscalingpolicy '%s' done", extracted)
	} else {
		renv.Log().Verbose("create appscalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAppscalingpolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("appscalingpolicy"), nil
}

func (cmd *CreateAppscalingpolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateAppscalingtarget(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateAppscalingtarget {
	cmd := new(CreateAppscalingtarget)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = applicationautoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateAppscalingtarget) SetApi(api applicationautoscalingiface.ApplicationAutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateAppscalingtarget) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateAppscalingtarget) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.RegisterScalableTargetInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.RegisterScalableTargetInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RegisterScalableTarget(input)
	renv.Log().ExtraVerbosef("applicationautoscaling.RegisterScalableTarget call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create appscalingtarget: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create appscalingtarget '%s' done", extracted)
	} else {
		renv.Log().Verbose("create appscalingtarget done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAppscalingtarget) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("appscalingtarget"), nil
}

func (cmd *CreateAppscalingtarget) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateBucket(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateBucket {
	cmd := new(CreateBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateBucket) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *CreateBucket) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateBucket) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.CreateBucketInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.CreateBucketInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateBucket(input)
	renv.Log().ExtraVerbosef("s3.CreateBucket call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create bucket '%s' done", extracted)
	} else {
		renv.Log().Verbose("create bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateBucket) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("bucket"), nil
}

func (cmd *CreateBucket) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateCertificate(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateCertificate {
	cmd := new(CreateCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = acm.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateCertificate) SetApi(api acmiface.ACMAPI) {
	cmd.api = api
}

func (cmd *CreateCertificate) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateCertificate) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create certificate '%s' done", extracted)
	} else {
		renv.Log().Verbose("create certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateCertificate) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("certificate"), nil
}

func (cmd *CreateCertificate) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateClassicLoadbalancer(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateClassicLoadbalancer {
	cmd := new(CreateClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elb.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateClassicLoadbalancer) SetApi(api elbiface.ELBAPI) {
	cmd.api = api
}

func (cmd *CreateClassicLoadbalancer) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateClassicLoadbalancer) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.CreateLoadBalancerInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.CreateLoadBalancerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateLoadBalancer(input)
	renv.Log().ExtraVerbosef("elb.CreateLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("create classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateClassicLoadbalancer) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("classicloadbalancer"), nil
}

func (cmd *CreateClassicLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateContainercluster(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateContainercluster {
	cmd := new(CreateContainercluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateContainercluster) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *CreateContainercluster) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateContainercluster) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.CreateClusterInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.CreateClusterInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateCluster(input)
	renv.Log().ExtraVerbosef("ecs.CreateCluster call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create containercluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create containercluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("create containercluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateContainercluster) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containercluster"), nil
}

func (cmd *CreateContainercluster) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateDatabase(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateDatabase {
	cmd := new(CreateDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *CreateDatabase) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateDatabase) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create database '%s' done", extracted)
	} else {
		renv.Log().Verbose("create database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDatabase) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *CreateDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateDbsubnetgroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateDbsubnetgroup {
	cmd := new(CreateDbsubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateDbsubnetgroup) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *CreateDbsubnetgroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateDbsubnetgroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.CreateDBSubnetGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.CreateDBSubnetGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateDBSubnetGroup(input)
	renv.Log().ExtraVerbosef("rds.CreateDBSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create dbsubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create dbsubnetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create dbsubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDbsubnetgroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("dbsubnetgroup"), nil
}

func (cmd *CreateDbsubnetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateDistribution(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateDistribution {
	cmd := new(CreateDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudfront.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateDistribution) SetApi(api cloudfrontiface.CloudFrontAPI) {
	cmd.api = api
}

func (cmd *CreateDistribution) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateDistribution) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create distribution '%s' done", extracted)
	} else {
		renv.Log().Verbose("create distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDistribution) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("distribution"), nil
}

func (cmd *CreateDistribution) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateElasticip(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateElasticip {
	cmd := new(CreateElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateElasticip) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateElasticip) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateElasticip) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AllocateAddressInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AllocateAddressInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AllocateAddress(input)
	renv.Log().ExtraVerbosef("ec2.AllocateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create elasticip '%s' done", extracted)
	} else {
		renv.Log().Verbose("create elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateElasticip) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.AllocateAddressInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AllocateAddressInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AllocateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AllocateAddress call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create elasticip ok")
			return fakeDryRunId("elasticip"), nil
		}
	}

	return nil, err
}

func (cmd *CreateElasticip) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateFunction(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateFunction {
	cmd := new(CreateFunction)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = lambda.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateFunction) SetApi(api lambdaiface.LambdaAPI) {
	cmd.api = api
}

func (cmd *CreateFunction) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateFunction) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &lambda.CreateFunctionInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in lambda.CreateFunctionInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateFunction(input)
	renv.Log().ExtraVerbosef("lambda.CreateFunction call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create function: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create function '%s' done", extracted)
	} else {
		renv.Log().Verbose("create function done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateFunction) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("function"), nil
}

func (cmd *CreateFunction) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateGroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateGroup {
	cmd := new(CreateGroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateGroup) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateGroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateGroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateGroup(input)
	renv.Log().ExtraVerbosef("iam.CreateGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create group: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create group '%s' done", extracted)
	} else {
		renv.Log().Verbose("create group done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateGroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("group"), nil
}

func (cmd *CreateGroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateImage(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateImage {
	cmd := new(CreateImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateImage) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateImage) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateImageInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateImageInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateImage(input)
	renv.Log().ExtraVerbosef("ec2.CreateImage call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create image '%s' done", extracted)
	} else {
		renv.Log().Verbose("create image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateImage) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateImageInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateImageInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateImage call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create image ok")
			return fakeDryRunId("image"), nil
		}
	}

	return nil, err
}

func (cmd *CreateImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateInstance {
	cmd := new(CreateInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.RunInstancesInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RunInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RunInstances(input)
	renv.Log().ExtraVerbosef("ec2.RunInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("create instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.RunInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RunInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.RunInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.RunInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, err
}

func (cmd *CreateInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateInstanceprofile(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateInstanceprofile {
	cmd := new(CreateInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateInstanceprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateInstanceprofile) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateInstanceprofile) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateInstanceProfileInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateInstanceProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateInstanceProfile(input)
	renv.Log().ExtraVerbosef("iam.CreateInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create instanceprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("create instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInstanceprofile) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instanceprofile"), nil
}

func (cmd *CreateInstanceprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateInternetgateway(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateInternetgateway {
	cmd := new(CreateInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateInternetgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateInternetgateway) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateInternetgateway) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateInternetGatewayInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateInternetGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateInternetGateway(input)
	renv.Log().ExtraVerbosef("ec2.CreateInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create internetgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("create internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInternetgateway) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateInternetGatewayInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateInternetGatewayInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateInternetGateway call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create internetgateway ok")
			return fakeDryRunId("internetgateway"), nil
		}
	}

	return nil, err
}

func (cmd *CreateInternetgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateKeypair(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateKeypair {
	cmd := new(CreateKeypair)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateKeypair) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateKeypair) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateKeypair) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ImportKeyPairInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ImportKeyPairInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ImportKeyPair(input)
	renv.Log().ExtraVerbosef("ec2.ImportKeyPair call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create keypair: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create keypair '%s' done", extracted)
	} else {
		renv.Log().Verbose("create keypair done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateKeypair) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("keypair"), nil
}

func (cmd *CreateKeypair) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateLaunchconfiguration(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateLaunchconfiguration {
	cmd := new(CreateLaunchconfiguration)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateLaunchconfiguration) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateLaunchconfiguration) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateLaunchconfiguration) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.CreateLaunchConfigurationInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.CreateLaunchConfigurationInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateLaunchConfiguration(input)
	renv.Log().ExtraVerbosef("autoscaling.CreateLaunchConfiguration call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create launchconfiguration: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create launchconfiguration '%s' done", extracted)
	} else {
		renv.Log().Verbose("create launchconfiguration done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLaunchconfiguration) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("launchconfiguration"), nil
}

func (cmd *CreateLaunchconfiguration) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateListener(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateListener {
	cmd := new(CreateListener)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateListener) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *CreateListener) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateListener) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateListenerInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateListenerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateListener(input)
	renv.Log().ExtraVerbosef("elbv2.CreateListener call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create listener: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create listener '%s' done", extracted)
	} else {
		renv.Log().Verbose("create listener done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateListener) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("listener"), nil
}

func (cmd *CreateListener) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateLoadbalancer(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateLoadbalancer {
	cmd := new(CreateLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateLoadbalancer) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *CreateLoadbalancer) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateLoadbalancer) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateLoadBalancerInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateLoadBalancerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateLoadBalancer(input)
	renv.Log().ExtraVerbosef("elbv2.CreateLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create loadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("create loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLoadbalancer) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loadbalancer"), nil
}

func (cmd *CreateLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateLoginprofile(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateLoginprofile {
	cmd := new(CreateLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateLoginprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateLoginprofile) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateLoginprofile) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateLoginProfileInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateLoginProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateLoginProfile(input)
	renv.Log().ExtraVerbosef("iam.CreateLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create loginprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("create loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLoginprofile) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loginprofile"), nil
}

func (cmd *CreateLoginprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateMfadevice(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateMfadevice {
	cmd := new(CreateMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateMfadevice) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateMfadevice) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateMfadevice) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create mfadevice '%s' done", extracted)
	} else {
		renv.Log().Verbose("create mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateMfadevice) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("mfadevice"), nil
}

func (cmd *CreateMfadevice) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateNatgateway(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateNatgateway {
	cmd := new(CreateNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateNatgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateNatgateway) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateNatgateway) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateNatGatewayInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateNatGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateNatGateway(input)
	renv.Log().ExtraVerbosef("ec2.CreateNatGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create natgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("create natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateNatgateway) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("natgateway"), nil
}

func (cmd *CreateNatgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateNetworkinterface(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateNetworkinterface {
	cmd := new(CreateNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateNetworkinterface) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateNetworkinterface) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateNetworkInterfaceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateNetworkInterfaceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateNetworkInterface(input)
	renv.Log().ExtraVerbosef("ec2.CreateNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("create networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateNetworkinterface) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateNetworkInterfaceInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateNetworkInterfaceInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateNetworkInterface(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateNetworkInterface call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create networkinterface ok")
			return fakeDryRunId("networkinterface"), nil
		}
	}

	return nil, err
}

func (cmd *CreateNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreatePolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreatePolicy {
	cmd := new(CreatePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreatePolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreatePolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreatePolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreatePolicyInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreatePolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreatePolicy(input)
	renv.Log().ExtraVerbosef("iam.CreatePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("create policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreatePolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *CreatePolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateQueue(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateQueue {
	cmd := new(CreateQueue)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sqs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateQueue) SetApi(api sqsiface.SQSAPI) {
	cmd.api = api
}

func (cmd *CreateQueue) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateQueue) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sqs.CreateQueueInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in sqs.CreateQueueInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateQueue(input)
	renv.Log().ExtraVerbosef("sqs.CreateQueue call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create queue: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create queue '%s' done", extracted)
	} else {
		renv.Log().Verbose("create queue done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateQueue) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("queue"), nil
}

func (cmd *CreateQueue) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRecord(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateRecord {
	cmd := new(CreateRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRecord) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *CreateRecord) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRecord) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create record '%s' done", extracted)
	} else {
		renv.Log().Verbose("create record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRecord) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("record"), nil
}

func (cmd *CreateRecord) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRepository(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateRepository {
	cmd := new(CreateRepository)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecr.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRepository) SetApi(api ecriface.ECRAPI) {
	cmd.api = api
}

func (cmd *CreateRepository) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRepository) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecr.CreateRepositoryInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ecr.CreateRepositoryInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateRepository(input)
	renv.Log().ExtraVerbosef("ecr.CreateRepository call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create repository: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create repository '%s' done", extracted)
	} else {
		renv.Log().Verbose("create repository done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRepository) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("repository"), nil
}

func (cmd *CreateRepository) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRole(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateRole {
	cmd := new(CreateRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRole) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateRole) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRole) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create role '%s' done", extracted)
	} else {
		renv.Log().Verbose("create role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRole) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("role"), nil
}

func (cmd *CreateRole) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRoute(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateRoute {
	cmd := new(CreateRoute)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRoute) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateRoute) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRoute) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateRouteInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateRoute(input)
	renv.Log().ExtraVerbosef("ec2.CreateRoute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create route: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create route '%s' done", extracted)
	} else {
		renv.Log().Verbose("create route done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRoute) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateRouteInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateRoute call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create route ok")
			return fakeDryRunId("route"), nil
		}
	}

	return nil, err
}

func (cmd *CreateRoute) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRoutetable(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateRoutetable {
	cmd := new(CreateRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRoutetable) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateRoutetable) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRoutetable) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateRouteTableInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteTableInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateRouteTable(input)
	renv.Log().ExtraVerbosef("ec2.CreateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create routetable '%s' done", extracted)
	} else {
		renv.Log().Verbose("create routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRoutetable) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateRouteTableInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteTableInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateRouteTable call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create routetable ok")
			return fakeDryRunId("routetable"), nil
		}
	}

	return nil, err
}

func (cmd *CreateRoutetable) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateS3object(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateS3object {
	cmd := new(CreateS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateS3object) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *CreateS3object) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateS3object) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create s3object '%s' done", extracted)
	} else {
		renv.Log().Verbose("create s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateS3object) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("s3object"), nil
}

func (cmd *CreateS3object) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateScalinggroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateScalinggroup {
	cmd := new(CreateScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateScalinggroup) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateScalinggroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateScalinggroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.CreateAutoScalingGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.CreateAutoScalingGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateAutoScalingGroup(input)
	renv.Log().ExtraVerbosef("autoscaling.CreateAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create scalinggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateScalinggroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalinggroup"), nil
}

func (cmd *CreateScalinggroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateScalingpolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateScalingpolicy {
	cmd := new(CreateScalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateScalingpolicy) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateScalingpolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateScalingpolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.PutScalingPolicyInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.PutScalingPolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.PutScalingPolicy(input)
	renv.Log().ExtraVerbosef("autoscaling.PutScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create scalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create scalingpolicy '%s' done", extracted)
	} else {
		renv.Log().Verbose("create scalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateScalingpolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalingpolicy"), nil
}

func (cmd *CreateScalingpolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateSecuritygroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateSecuritygroup {
	cmd := new(CreateSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateSecuritygroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSecuritygroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSecurityGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSecurityGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateSecurityGroup(input)
	renv.Log().ExtraVerbosef("ec2.CreateSecurityGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSecuritygroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateSecurityGroupInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSecurityGroupInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateSecurityGroup call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create securitygroup ok")
			return fakeDryRunId("securitygroup"), nil
		}
	}

	return nil, err
}

func (cmd *CreateSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateSnapshot(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateSnapshot {
	cmd := new(CreateSnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSnapshot) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateSnapshot) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSnapshot) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSnapshotInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSnapshotInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateSnapshot(input)
	renv.Log().ExtraVerbosef("ec2.CreateSnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create snapshot '%s' done", extracted)
	} else {
		renv.Log().Verbose("create snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSnapshot) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateSnapshotInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSnapshotInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateSnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateSnapshot call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create snapshot ok")
			return fakeDryRunId("snapshot"), nil
		}
	}

	return nil, err
}

func (cmd *CreateSnapshot) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateStack(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateStack {
	cmd := new(CreateStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudformation.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateStack) SetApi(api cloudformationiface.CloudFormationAPI) {
	cmd.api = api
}

func (cmd *CreateStack) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateStack) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.CreateStackInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.CreateStackInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateStack(input)
	renv.Log().ExtraVerbosef("cloudformation.CreateStack call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create stack '%s' done", extracted)
	} else {
		renv.Log().Verbose("create stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateStack) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("stack"), nil
}

func (cmd *CreateStack) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateSubnet(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateSubnet {
	cmd := new(CreateSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSubnet) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateSubnet) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSubnet) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSubnetInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSubnetInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateSubnet(input)
	renv.Log().ExtraVerbosef("ec2.CreateSubnet call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create subnet '%s' done", extracted)
	} else {
		renv.Log().Verbose("create subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSubnet) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateSubnetInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSubnetInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateSubnet call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create subnet ok")
			return fakeDryRunId("subnet"), nil
		}
	}

	return nil, err
}

func (cmd *CreateSubnet) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateSubscription(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateSubscription {
	cmd := new(CreateSubscription)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sns.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSubscription) SetApi(api snsiface.SNSAPI) {
	cmd.api = api
}

func (cmd *CreateSubscription) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSubscription) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.SubscribeInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.SubscribeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.Subscribe(input)
	renv.Log().ExtraVerbosef("sns.Subscribe call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create subscription: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create subscription '%s' done", extracted)
	} else {
		renv.Log().Verbose("create subscription done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSubscription) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("subscription"), nil
}

func (cmd *CreateSubscription) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateTag(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateTag {
	cmd := new(CreateTag)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateTag) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateTag) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateTag) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create tag: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create tag '%s' done", extracted)
	} else {
		renv.Log().Verbose("create tag done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTag) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateTargetgroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateTargetgroup {
	cmd := new(CreateTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateTargetgroup) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *CreateTargetgroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateTargetgroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateTargetGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateTargetGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateTargetGroup(input)
	renv.Log().ExtraVerbosef("elbv2.CreateTargetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create targetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTargetgroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("targetgroup"), nil
}

func (cmd *CreateTargetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateTopic(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateTopic {
	cmd := new(CreateTopic)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sns.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateTopic) SetApi(api snsiface.SNSAPI) {
	cmd.api = api
}

func (cmd *CreateTopic) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateTopic) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.CreateTopicInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.CreateTopicInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateTopic(input)
	renv.Log().ExtraVerbosef("sns.CreateTopic call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create topic: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create topic '%s' done", extracted)
	} else {
		renv.Log().Verbose("create topic done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTopic) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("topic"), nil
}

func (cmd *CreateTopic) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateUser(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateUser {
	cmd := new(CreateUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateUser) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateUser) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateUser) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateUserInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateUserInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateUser(input)
	renv.Log().ExtraVerbosef("iam.CreateUser call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create user '%s' done", extracted)
	} else {
		renv.Log().Verbose("create user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateUser) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("user"), nil
}

func (cmd *CreateUser) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateVolume(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateVolume {
	cmd := new(CreateVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateVolume) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateVolume) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateVolumeInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVolumeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateVolume(input)
	renv.Log().ExtraVerbosef("ec2.CreateVolume call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("create volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateVolume) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateVolumeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVolumeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateVolume call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create volume ok")
			return fakeDryRunId("volume"), nil
		}
	}

	return nil, err
}

func (cmd *CreateVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateVpc(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateVpc {
	cmd := new(CreateVpc)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateVpc) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateVpc) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateVpc) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateVpcInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVpcInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateVpc(input)
	renv.Log().ExtraVerbosef("ec2.CreateVpc call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create vpc: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create vpc '%s' done", extracted)
	} else {
		renv.Log().Verbose("create vpc done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateVpc) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateVpcInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVpcInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateVpc call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create vpc ok")
			return fakeDryRunId("vpc"), nil
		}
	}

	return nil, err
}

func (cmd *CreateVpc) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateZone(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *CreateZone {
	cmd := new(CreateZone)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateZone) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *CreateZone) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateZone) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &route53.CreateHostedZoneInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in route53.CreateHostedZoneInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateHostedZone(input)
	renv.Log().ExtraVerbosef("route53.CreateHostedZone call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create zone: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create zone '%s' done", extracted)
	} else {
		renv.Log().Verbose("create zone done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateZone) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("zone"), nil
}

func (cmd *CreateZone) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteAccesskey(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteAccesskey {
	cmd := new(DeleteAccesskey)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteAccesskey) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteAccesskey) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteAccesskey) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteAccessKeyInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteAccessKeyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAccessKey(input)
	renv.Log().ExtraVerbosef("iam.DeleteAccessKey call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete accesskey: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete accesskey '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete accesskey done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAccesskey) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("accesskey"), nil
}

func (cmd *DeleteAccesskey) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteAlarm(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteAlarm {
	cmd := new(DeleteAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *DeleteAlarm) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteAlarm) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.DeleteAlarmsInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.DeleteAlarmsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAlarms(input)
	renv.Log().ExtraVerbosef("cloudwatch.DeleteAlarms call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAlarm) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *DeleteAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteAppscalingpolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteAppscalingpolicy {
	cmd := new(DeleteAppscalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = applicationautoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteAppscalingpolicy) SetApi(api applicationautoscalingiface.ApplicationAutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteAppscalingpolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteAppscalingpolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.DeleteScalingPolicyInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.DeleteScalingPolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteScalingPolicy(input)
	renv.Log().ExtraVerbosef("applicationautoscaling.DeleteScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete appscalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete appscalingpolicy '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete appscalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAppscalingpolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("appscalingpolicy"), nil
}

func (cmd *DeleteAppscalingpolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteAppscalingtarget(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteAppscalingtarget {
	cmd := new(DeleteAppscalingtarget)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = applicationautoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteAppscalingtarget) SetApi(api applicationautoscalingiface.ApplicationAutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteAppscalingtarget) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteAppscalingtarget) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.DeregisterScalableTargetInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.DeregisterScalableTargetInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeregisterScalableTarget(input)
	renv.Log().ExtraVerbosef("applicationautoscaling.DeregisterScalableTarget call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete appscalingtarget: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete appscalingtarget '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete appscalingtarget done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAppscalingtarget) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("appscalingtarget"), nil
}

func (cmd *DeleteAppscalingtarget) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteBucket(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteBucket {
	cmd := new(DeleteBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteBucket) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *DeleteBucket) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteBucket) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.DeleteBucketInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.DeleteBucketInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteBucket(input)
	renv.Log().ExtraVerbosef("s3.DeleteBucket call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete bucket '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteBucket) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("bucket"), nil
}

func (cmd *DeleteBucket) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteCertificate(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteCertificate {
	cmd := new(DeleteCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = acm.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteCertificate) SetApi(api acmiface.ACMAPI) {
	cmd.api = api
}

func (cmd *DeleteCertificate) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteCertificate) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &acm.DeleteCertificateInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in acm.DeleteCertificateInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCertificate(input)
	renv.Log().ExtraVerbosef("acm.DeleteCertificate call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete certificate '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteCertificate) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("certificate"), nil
}

func (cmd *DeleteCertificate) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteClassicLoadbalancer(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteClassicLoadbalancer {
	cmd := new(DeleteClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elb.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteClassicLoadbalancer) SetApi(api elbiface.ELBAPI) {
	cmd.api = api
}

func (cmd *DeleteClassicLoadbalancer) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteClassicLoadbalancer) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.DeleteLoadBalancerInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.DeleteLoadBalancerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLoadBalancer(input)
	renv.Log().ExtraVerbosef("elb.DeleteLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteClassicLoadbalancer) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("classicloadbalancer"), nil
}

func (cmd *DeleteClassicLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteContainercluster(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteContainercluster {
	cmd := new(DeleteContainercluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteContainercluster) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *DeleteContainercluster) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteContainercluster) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.DeleteClusterInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.DeleteClusterInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCluster(input)
	renv.Log().ExtraVerbosef("ecs.DeleteCluster call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete containercluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete containercluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete containercluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteContainercluster) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containercluster"), nil
}

func (cmd *DeleteContainercluster) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteContainertask(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteContainertask {
	cmd := new(DeleteContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *DeleteContainertask) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteContainertask) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteDatabase(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteDatabase {
	cmd := new(DeleteDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *DeleteDatabase) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteDatabase) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.DeleteDBInstanceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.DeleteDBInstanceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteDBInstance(input)
	renv.Log().ExtraVerbosef("rds.DeleteDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete database '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDatabase) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *DeleteDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteDbsubnetgroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteDbsubnetgroup {
	cmd := new(DeleteDbsubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteDbsubnetgroup) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *DeleteDbsubnetgroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteDbsubnetgroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.DeleteDBSubnetGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.DeleteDBSubnetGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteDBSubnetGroup(input)
	renv.Log().ExtraVerbosef("rds.DeleteDBSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete dbsubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete dbsubnetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete dbsubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDbsubnetgroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("dbsubnetgroup"), nil
}

func (cmd *DeleteDbsubnetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteDistribution(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteDistribution {
	cmd := new(DeleteDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudfront.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteDistribution) SetApi(api cloudfrontiface.CloudFrontAPI) {
	cmd.api = api
}

func (cmd *DeleteDistribution) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteDistribution) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete distribution '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDistribution) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("distribution"), nil
}

func (cmd *DeleteDistribution) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteElasticip(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteElasticip {
	cmd := new(DeleteElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteElasticip) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteElasticip) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteElasticip) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ReleaseAddressInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ReleaseAddressInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ReleaseAddress(input)
	renv.Log().ExtraVerbosef("ec2.ReleaseAddress call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete elasticip '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteElasticip) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.ReleaseAddressInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ReleaseAddressInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.ReleaseAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.ReleaseAddress call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete elasticip ok")
			return fakeDryRunId("elasticip"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteElasticip) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteFunction(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteFunction {
	cmd := new(DeleteFunction)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = lambda.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteFunction) SetApi(api lambdaiface.LambdaAPI) {
	cmd.api = api
}

func (cmd *DeleteFunction) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteFunction) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &lambda.DeleteFunctionInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in lambda.DeleteFunctionInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteFunction(input)
	renv.Log().ExtraVerbosef("lambda.DeleteFunction call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete function: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete function '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete function done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteFunction) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("function"), nil
}

func (cmd *DeleteFunction) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteGroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteGroup {
	cmd := new(DeleteGroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteGroup) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteGroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteGroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteGroup(input)
	renv.Log().ExtraVerbosef("iam.DeleteGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete group: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete group '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete group done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteGroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("group"), nil
}

func (cmd *DeleteGroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteImage(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteImage {
	cmd := new(DeleteImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteImage) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteImage) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete image '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteInstance {
	cmd := new(DeleteInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.TerminateInstancesInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.TerminateInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.TerminateInstances(input)
	renv.Log().ExtraVerbosef("ec2.TerminateInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.TerminateInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.TerminateInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.TerminateInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.TerminateInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteInstanceprofile(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteInstanceprofile {
	cmd := new(DeleteInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteInstanceprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteInstanceprofile) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteInstanceprofile) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteInstanceProfileInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteInstanceProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteInstanceProfile(input)
	renv.Log().ExtraVerbosef("iam.DeleteInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete instanceprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInstanceprofile) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instanceprofile"), nil
}

func (cmd *DeleteInstanceprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteInternetgateway(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteInternetgateway {
	cmd := new(DeleteInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteInternetgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteInternetgateway) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteInternetgateway) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteInternetGatewayInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteInternetGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteInternetGateway(input)
	renv.Log().ExtraVerbosef("ec2.DeleteInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete internetgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInternetgateway) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteInternetGatewayInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteInternetGatewayInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteInternetGateway call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete internetgateway ok")
			return fakeDryRunId("internetgateway"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteInternetgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteKeypair(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteKeypair {
	cmd := new(DeleteKeypair)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteKeypair) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteKeypair) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteKeypair) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteKeyPairInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteKeyPairInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteKeyPair(input)
	renv.Log().ExtraVerbosef("ec2.DeleteKeyPair call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete keypair: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete keypair '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete keypair done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteKeypair) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteKeyPairInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteKeyPairInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteKeyPair(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteKeyPair call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete keypair ok")
			return fakeDryRunId("keypair"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteKeypair) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteLaunchconfiguration(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteLaunchconfiguration {
	cmd := new(DeleteLaunchconfiguration)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteLaunchconfiguration) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteLaunchconfiguration) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteLaunchconfiguration) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeleteLaunchConfigurationInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeleteLaunchConfigurationInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLaunchConfiguration(input)
	renv.Log().ExtraVerbosef("autoscaling.DeleteLaunchConfiguration call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete launchconfiguration: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete launchconfiguration '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete launchconfiguration done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLaunchconfiguration) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("launchconfiguration"), nil
}

func (cmd *DeleteLaunchconfiguration) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteListener(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteListener {
	cmd := new(DeleteListener)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteListener) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *DeleteListener) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteListener) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteListenerInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteListenerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteListener(input)
	renv.Log().ExtraVerbosef("elbv2.DeleteListener call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete listener: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete listener '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete listener done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteListener) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("listener"), nil
}

func (cmd *DeleteListener) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteLoadbalancer(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteLoadbalancer {
	cmd := new(DeleteLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteLoadbalancer) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *DeleteLoadbalancer) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteLoadbalancer) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteLoadBalancerInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteLoadBalancerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLoadBalancer(input)
	renv.Log().ExtraVerbosef("elbv2.DeleteLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete loadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLoadbalancer) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loadbalancer"), nil
}

func (cmd *DeleteLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteLoginprofile(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteLoginprofile {
	cmd := new(DeleteLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteLoginprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteLoginprofile) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteLoginprofile) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteLoginProfileInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteLoginProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLoginProfile(input)
	renv.Log().ExtraVerbosef("iam.DeleteLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete loginprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLoginprofile) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loginprofile"), nil
}

func (cmd *DeleteLoginprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteMfadevice(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteMfadevice {
	cmd := new(DeleteMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteMfadevice) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteMfadevice) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteMfadevice) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteVirtualMFADeviceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteVirtualMFADeviceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVirtualMFADevice(input)
	renv.Log().ExtraVerbosef("iam.DeleteVirtualMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete mfadevice '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteMfadevice) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("mfadevice"), nil
}

func (cmd *DeleteMfadevice) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteNatgateway(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteNatgateway {
	cmd := new(DeleteNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteNatgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteNatgateway) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteNatgateway) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteNatGatewayInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteNatGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteNatGateway(input)
	renv.Log().ExtraVerbosef("ec2.DeleteNatGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete natgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteNatgateway) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("natgateway"), nil
}

func (cmd *DeleteNatgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteNetworkinterface(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteNetworkinterface {
	cmd := new(DeleteNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteNetworkinterface) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteNetworkinterface) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteNetworkInterfaceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteNetworkInterfaceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteNetworkInterface(input)
	renv.Log().ExtraVerbosef("ec2.DeleteNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteNetworkinterface) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteNetworkInterfaceInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteNetworkInterfaceInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteNetworkInterface(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteNetworkInterface call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete networkinterface ok")
			return fakeDryRunId("networkinterface"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeletePolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeletePolicy {
	cmd := new(DeletePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeletePolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeletePolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeletePolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeletePolicyInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeletePolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeletePolicy(input)
	renv.Log().ExtraVerbosef("iam.DeletePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeletePolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *DeletePolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteQueue(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteQueue {
	cmd := new(DeleteQueue)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sqs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteQueue) SetApi(api sqsiface.SQSAPI) {
	cmd.api = api
}

func (cmd *DeleteQueue) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteQueue) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sqs.DeleteQueueInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in sqs.DeleteQueueInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteQueue(input)
	renv.Log().ExtraVerbosef("sqs.DeleteQueue call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete queue: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete queue '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete queue done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteQueue) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("queue"), nil
}

func (cmd *DeleteQueue) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRecord(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRecord {
	cmd := new(DeleteRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRecord) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *DeleteRecord) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRecord) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete record '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRecord) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("record"), nil
}

func (cmd *DeleteRecord) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRepository(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRepository {
	cmd := new(DeleteRepository)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecr.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRepository) SetApi(api ecriface.ECRAPI) {
	cmd.api = api
}

func (cmd *DeleteRepository) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRepository) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecr.DeleteRepositoryInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ecr.DeleteRepositoryInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRepository(input)
	renv.Log().ExtraVerbosef("ecr.DeleteRepository call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete repository: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete repository '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete repository done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRepository) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("repository"), nil
}

func (cmd *DeleteRepository) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRole(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRole {
	cmd := new(DeleteRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRole) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteRole) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRole) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete role '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRole) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("role"), nil
}

func (cmd *DeleteRole) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRoute(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRoute {
	cmd := new(DeleteRoute)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRoute) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteRoute) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRoute) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteRouteInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRoute(input)
	renv.Log().ExtraVerbosef("ec2.DeleteRoute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete route: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete route '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete route done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRoute) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteRouteInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteRoute call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete route ok")
			return fakeDryRunId("route"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteRoute) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRoutetable(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRoutetable {
	cmd := new(DeleteRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRoutetable) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteRoutetable) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRoutetable) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteRouteTableInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteTableInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRouteTable(input)
	renv.Log().ExtraVerbosef("ec2.DeleteRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete routetable '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRoutetable) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteRouteTableInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteTableInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteRouteTable call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete routetable ok")
			return fakeDryRunId("routetable"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteRoutetable) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteS3object(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteS3object {
	cmd := new(DeleteS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteS3object) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *DeleteS3object) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteS3object) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.DeleteObjectInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.DeleteObjectInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteObject(input)
	renv.Log().ExtraVerbosef("s3.DeleteObject call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete s3object '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteS3object) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("s3object"), nil
}

func (cmd *DeleteS3object) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteScalinggroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteScalinggroup {
	cmd := new(DeleteScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteScalinggroup) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteScalinggroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteScalinggroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeleteAutoScalingGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeleteAutoScalingGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAutoScalingGroup(input)
	renv.Log().ExtraVerbosef("autoscaling.DeleteAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete scalinggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteScalinggroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalinggroup"), nil
}

func (cmd *DeleteScalinggroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteScalingpolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteScalingpolicy {
	cmd := new(DeleteScalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteScalingpolicy) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteScalingpolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteScalingpolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeletePolicyInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeletePolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeletePolicy(input)
	renv.Log().ExtraVerbosef("autoscaling.DeletePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete scalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete scalingpolicy '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete scalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteScalingpolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalingpolicy"), nil
}

func (cmd *DeleteScalingpolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteSecuritygroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSecuritygroup {
	cmd := new(DeleteSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteSecuritygroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSecuritygroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSecurityGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSecurityGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSecurityGroup(input)
	renv.Log().ExtraVerbosef("ec2.DeleteSecurityGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSecuritygroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteSecurityGroupInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSecurityGroupInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteSecurityGroup call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete securitygroup ok")
			return fakeDryRunId("securitygroup"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteSnapshot(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSnapshot {
	cmd := new(DeleteSnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSnapshot) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteSnapshot) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSnapshot) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSnapshotInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSnapshotInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSnapshot(input)
	renv.Log().ExtraVerbosef("ec2.DeleteSnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete snapshot '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSnapshot) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteSnapshotInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSnapshotInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteSnapshot call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete snapshot ok")
			return fakeDryRunId("snapshot"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteSnapshot) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteStack(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteStack {
	cmd := new(DeleteStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudformation.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteStack) SetApi(api cloudformationiface.CloudFormationAPI) {
	cmd.api = api
}

func (cmd *DeleteStack) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteStack) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.DeleteStackInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.DeleteStackInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteStack(input)
	renv.Log().ExtraVerbosef("cloudformation.DeleteStack call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete stack '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteStack) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("stack"), nil
}

func (cmd *DeleteStack) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteSubnet(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSubnet {
	cmd := new(DeleteSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSubnet) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteSubnet) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSubnet) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSubnetInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSubnetInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSubnet(input)
	renv.Log().ExtraVerbosef("ec2.DeleteSubnet call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete subnet '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSubnet) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteSubnetInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSubnetInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteSubnet call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete subnet ok")
			return fakeDryRunId("subnet"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteSubnet) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteSubscription(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSubscription {
	cmd := new(DeleteSubscription)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sns.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSubscription) SetApi(api snsiface.SNSAPI) {
	cmd.api = api
}

func (cmd *DeleteSubscription) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSubscription) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.UnsubscribeInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.UnsubscribeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.Unsubscribe(input)
	renv.Log().ExtraVerbosef("sns.Unsubscribe call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete subscription: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete subscription '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete subscription done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSubscription) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("subscription"), nil
}

func (cmd *DeleteSubscription) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteTag(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteTag {
	cmd := new(DeleteTag)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteTag) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteTag) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteTag) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete tag: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete tag '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete tag done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTag) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteTargetgroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteTargetgroup {
	cmd := new(DeleteTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteTargetgroup) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *DeleteTargetgroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteTargetgroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteTargetGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteTargetGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteTargetGroup(input)
	renv.Log().ExtraVerbosef("elbv2.DeleteTargetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete targetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTargetgroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("targetgroup"), nil
}

func (cmd *DeleteTargetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteTopic(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteTopic {
	cmd := new(DeleteTopic)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sns.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteTopic) SetApi(api snsiface.SNSAPI) {
	cmd.api = api
}

func (cmd *DeleteTopic) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteTopic) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.DeleteTopicInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.DeleteTopicInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteTopic(input)
	renv.Log().ExtraVerbosef("sns.DeleteTopic call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete topic: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete topic '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete topic done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTopic) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("topic"), nil
}

func (cmd *DeleteTopic) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteUser(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteUser {
	cmd := new(DeleteUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteUser) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteUser) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteUser) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteUserInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteUserInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteUser(input)
	renv.Log().ExtraVerbosef("iam.DeleteUser call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete user '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteUser) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("user"), nil
}

func (cmd *DeleteUser) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteVolume(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteVolume {
	cmd := new(DeleteVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteVolume) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteVolume) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteVolumeInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVolumeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVolume(input)
	renv.Log().ExtraVerbosef("ec2.DeleteVolume call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteVolume) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteVolumeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVolumeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteVolume call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete volume ok")
			return fakeDryRunId("volume"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteVpc(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteVpc {
	cmd := new(DeleteVpc)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteVpc) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteVpc) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteVpc) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteVpcInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVpcInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVpc(input)
	renv.Log().ExtraVerbosef("ec2.DeleteVpc call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete vpc: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete vpc '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete vpc done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteVpc) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteVpcInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVpcInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteVpc call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete vpc ok")
			return fakeDryRunId("vpc"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteVpc) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteZone(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DeleteZone {
	cmd := new(DeleteZone)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteZone) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *DeleteZone) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteZone) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &route53.DeleteHostedZoneInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in route53.DeleteHostedZoneInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteHostedZone(input)
	renv.Log().ExtraVerbosef("route53.DeleteHostedZone call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete zone: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete zone '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete zone done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteZone) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("zone"), nil
}

func (cmd *DeleteZone) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachAlarm(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachAlarm {
	cmd := new(DetachAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *DetachAlarm) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachAlarm) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachAlarm) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *DetachAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachClassicLoadbalancer(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachClassicLoadbalancer {
	cmd := new(DetachClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elb.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachClassicLoadbalancer) SetApi(api elbiface.ELBAPI) {
	cmd.api = api
}

func (cmd *DetachClassicLoadbalancer) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachClassicLoadbalancer) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.DeregisterInstancesFromLoadBalancerInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.DeregisterInstancesFromLoadBalancerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeregisterInstancesFromLoadBalancer(input)
	renv.Log().ExtraVerbosef("elb.DeregisterInstancesFromLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachClassicLoadbalancer) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("classicloadbalancer"), nil
}

func (cmd *DetachClassicLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachContainertask(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachContainertask {
	cmd := new(DetachContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *DetachContainertask) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachContainertask) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachContainertask) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *DetachContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachElasticip(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachElasticip {
	cmd := new(DetachElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachElasticip) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachElasticip) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachElasticip) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DisassociateAddressInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateAddressInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DisassociateAddress(input)
	renv.Log().ExtraVerbosef("ec2.DisassociateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach elasticip '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachElasticip) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DisassociateAddressInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateAddressInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DisassociateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DisassociateAddress call took %s", time.Since(start))
			renv.Log().Verbose("dry run: detach elasticip ok")
			return fakeDryRunId("elasticip"), nil
		}
	}

	return nil, err
}

func (cmd *DetachElasticip) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachInstance {
	cmd := new(DetachInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachInstance) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *DetachInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeregisterTargetsInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeregisterTargetsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeregisterTargets(input)
	renv.Log().ExtraVerbosef("elbv2.DeregisterTargets call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instance"), nil
}

func (cmd *DetachInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachInstanceprofile(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachInstanceprofile {
	cmd := new(DetachInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachInstanceprofile) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachInstanceprofile) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachInstanceprofile) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach instanceprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInstanceprofile) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instanceprofile"), nil
}

func (cmd *DetachInstanceprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachInternetgateway(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachInternetgateway {
	cmd := new(DetachInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachInternetgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachInternetgateway) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachInternetgateway) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DetachInternetGatewayInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachInternetGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DetachInternetGateway(input)
	renv.Log().ExtraVerbosef("ec2.DetachInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach internetgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInternetgateway) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DetachInternetGatewayInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachInternetGatewayInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DetachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DetachInternetGateway call took %s", time.Since(start))
			renv.Log().Verbose("dry run: detach internetgateway ok")
			return fakeDryRunId("internetgateway"), nil
		}
	}

	return nil, err
}

func (cmd *DetachInternetgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachMfadevice(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachMfadevice {
	cmd := new(DetachMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachMfadevice) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DetachMfadevice) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachMfadevice) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeactivateMFADeviceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeactivateMFADeviceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeactivateMFADevice(input)
	renv.Log().ExtraVerbosef("iam.DeactivateMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach mfadevice '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachMfadevice) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("mfadevice"), nil
}

func (cmd *DetachMfadevice) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachNetworkinterface(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachNetworkinterface {
	cmd := new(DetachNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachNetworkinterface) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachNetworkinterface) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachPolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachPolicy {
	cmd := new(DetachPolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachPolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DetachPolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachPolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachPolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *DetachPolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachRole(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachRole {
	cmd := new(DetachRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachRole) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DetachRole) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachRole) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.RemoveRoleFromInstanceProfileInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.RemoveRoleFromInstanceProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RemoveRoleFromInstanceProfile(input)
	renv.Log().ExtraVerbosef("iam.RemoveRoleFromInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach role '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachRole) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("role"), nil
}

func (cmd *DetachRole) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachRoutetable(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachRoutetable {
	cmd := new(DetachRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachRoutetable) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachRoutetable) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachRoutetable) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DisassociateRouteTableInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateRouteTableInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DisassociateRouteTable(input)
	renv.Log().ExtraVerbosef("ec2.DisassociateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach routetable '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachRoutetable) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DisassociateRouteTableInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateRouteTableInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DisassociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DisassociateRouteTable call took %s", time.Since(start))
			renv.Log().Verbose("dry run: detach routetable ok")
			return fakeDryRunId("routetable"), nil
		}
	}

	return nil, err
}

func (cmd *DetachRoutetable) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachSecuritygroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachSecuritygroup {
	cmd := new(DetachSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachSecuritygroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachSecuritygroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachSecuritygroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("securitygroup"), nil
}

func (cmd *DetachSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachUser(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachUser {
	cmd := new(DetachUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachUser) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DetachUser) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachUser) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.RemoveUserFromGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.RemoveUserFromGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RemoveUserFromGroup(input)
	renv.Log().ExtraVerbosef("iam.RemoveUserFromGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach user '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachUser) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("user"), nil
}

func (cmd *DetachUser) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachVolume(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *DetachVolume {
	cmd := new(DetachVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachVolume) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachVolume) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DetachVolumeInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachVolumeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DetachVolume(input)
	renv.Log().ExtraVerbosef("ec2.DetachVolume call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachVolume) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DetachVolumeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachVolumeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DetachVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DetachVolume call took %s", time.Since(start))
			renv.Log().Verbose("dry run: detach volume ok")
			return fakeDryRunId("volume"), nil
		}
	}

	return nil, err
}

func (cmd *DetachVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewImportImage(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *ImportImage {
	cmd := new(ImportImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *ImportImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *ImportImage) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *ImportImage) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ImportImageInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ImportImageInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ImportImage(input)
	renv.Log().ExtraVerbosef("ec2.ImportImage call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("import image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("import image '%s' done", extracted)
	} else {
		renv.Log().Verbose("import image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *ImportImage) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.ImportImageInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ImportImageInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.ImportImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.ImportImage call took %s", time.Since(start))
			renv.Log().Verbose("dry run: import image ok")
			return fakeDryRunId("image"), nil
		}
	}

	return nil, err
}

func (cmd *ImportImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewRestartDatabase(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *RestartDatabase {
	cmd := new(RestartDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *RestartDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *RestartDatabase) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *RestartDatabase) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.RebootDBInstanceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.RebootDBInstanceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RebootDBInstance(input)
	renv.Log().ExtraVerbosef("rds.RebootDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("restart database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("restart database '%s' done", extracted)
	} else {
		renv.Log().Verbose("restart database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *RestartDatabase) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *RestartDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewRestartInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *RestartInstance {
	cmd := new(RestartInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *RestartInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *RestartInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *RestartInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.RebootInstancesInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RebootInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RebootInstances(input)
	renv.Log().ExtraVerbosef("ec2.RebootInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("restart instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("restart instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("restart instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *RestartInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.RebootInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RebootInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.RebootInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.RebootInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: restart instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, err
}

func (cmd *RestartInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStartAlarm(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *StartAlarm {
	cmd := new(StartAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *StartAlarm) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartAlarm) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.EnableAlarmActionsInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.EnableAlarmActionsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.EnableAlarmActions(input)
	renv.Log().ExtraVerbosef("cloudwatch.EnableAlarmActions call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("start alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartAlarm) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *StartAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStartContainertask(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *StartContainertask {
	cmd := new(StartContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *StartContainertask) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartContainertask) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("start containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartContainertask) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *StartContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStartDatabase(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *StartDatabase {
	cmd := new(StartDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *StartDatabase) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartDatabase) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.StartDBInstanceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.StartDBInstanceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.StartDBInstance(input)
	renv.Log().ExtraVerbosef("rds.StartDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start database '%s' done", extracted)
	} else {
		renv.Log().Verbose("start database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartDatabase) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *StartDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStartInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *StartInstance {
	cmd := new(StartInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *StartInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.StartInstancesInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StartInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.StartInstances(input)
	renv.Log().ExtraVerbosef("ec2.StartInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("start instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.StartInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StartInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.StartInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.StartInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: start instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, err
}

func (cmd *StartInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStopAlarm(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *StopAlarm {
	cmd := new(StopAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *StopAlarm) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopAlarm) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.DisableAlarmActionsInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.DisableAlarmActionsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DisableAlarmActions(input)
	renv.Log().ExtraVerbosef("cloudwatch.DisableAlarmActions call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopAlarm) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *StopAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStopContainertask(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *StopContainertask {
	cmd := new(StopContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *StopContainertask) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopContainertask) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopContainertask) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *StopContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStopDatabase(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *StopDatabase {
	cmd := new(StopDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *StopDatabase) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopDatabase) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.StopDBInstanceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.StopDBInstanceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.StopDBInstance(input)
	renv.Log().ExtraVerbosef("rds.StopDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop database '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopDatabase) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *StopDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStopInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *StopInstance {
	cmd := new(StopInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *StopInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.StopInstancesInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StopInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.StopInstances(input)
	renv.Log().ExtraVerbosef("ec2.StopInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.StopInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StopInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.StopInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.StopInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: stop instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, err
}

func (cmd *StopInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateBucket(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateBucket {
	cmd := new(UpdateBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateBucket) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *UpdateBucket) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateBucket) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update bucket '%s' done", extracted)
	} else {
		renv.Log().Verbose("update bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateBucket) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("bucket"), nil
}

func (cmd *UpdateBucket) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateClassicLoadbalancer(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateClassicLoadbalancer {
	cmd := new(UpdateClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elb.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateClassicLoadbalancer) SetApi(api elbiface.ELBAPI) {
	cmd.api = api
}

func (cmd *UpdateClassicLoadbalancer) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateClassicLoadbalancer) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.ConfigureHealthCheckInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.ConfigureHealthCheckInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ConfigureHealthCheck(input)
	renv.Log().ExtraVerbosef("elb.ConfigureHealthCheck call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("update classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateClassicLoadbalancer) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("classicloadbalancer"), nil
}

func (cmd *UpdateClassicLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateContainertask(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateContainertask {
	cmd := new(UpdateContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *UpdateContainertask) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateContainertask) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.UpdateServiceInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.UpdateServiceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.UpdateService(input)
	renv.Log().ExtraVerbosef("ecs.UpdateService call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("update containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateContainertask) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *UpdateContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateDistribution(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateDistribution {
	cmd := new(UpdateDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudfront.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateDistribution) SetApi(api cloudfrontiface.CloudFrontAPI) {
	cmd.api = api
}

func (cmd *UpdateDistribution) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateDistribution) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update distribution '%s' done", extracted)
	} else {
		renv.Log().Verbose("update distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateDistribution) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("distribution"), nil
}

func (cmd *UpdateDistribution) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateImage(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateImage {
	cmd := new(UpdateImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *UpdateImage) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateImage) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update image '%s' done", extracted)
	} else {
		renv.Log().Verbose("update image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateInstance(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateInstance {
	cmd := new(UpdateInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *UpdateInstance) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateInstance) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ModifyInstanceAttributeInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifyInstanceAttributeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ModifyInstanceAttribute(input)
	renv.Log().ExtraVerbosef("ec2.ModifyInstanceAttribute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("update instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateInstance) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.ModifyInstanceAttributeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifyInstanceAttributeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.ModifyInstanceAttribute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.ModifyInstanceAttribute call took %s", time.Since(start))
			renv.Log().Verbose("dry run: update instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, err
}

func (cmd *UpdateInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateLoginprofile(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateLoginprofile {
	cmd := new(UpdateLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateLoginprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *UpdateLoginprofile) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateLoginprofile) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.UpdateLoginProfileInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.UpdateLoginProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.UpdateLoginProfile(input)
	renv.Log().ExtraVerbosef("iam.UpdateLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update loginprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("update loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateLoginprofile) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loginprofile"), nil
}

func (cmd *UpdateLoginprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdatePolicy(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdatePolicy {
	cmd := new(UpdatePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdatePolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *UpdatePolicy) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdatePolicy) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreatePolicyVersionInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreatePolicyVersionInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreatePolicyVersion(input)
	renv.Log().ExtraVerbosef("iam.CreatePolicyVersion call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("update policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdatePolicy) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *UpdatePolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateRecord(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateRecord {
	cmd := new(UpdateRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateRecord) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *UpdateRecord) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateRecord) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update record '%s' done", extracted)
	} else {
		renv.Log().Verbose("update record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateRecord) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("record"), nil
}

func (cmd *UpdateRecord) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateS3object(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateS3object {
	cmd := new(UpdateS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateS3object) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *UpdateS3object) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateS3object) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.PutObjectAclInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.PutObjectAclInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.PutObjectAcl(input)
	renv.Log().ExtraVerbosef("s3.PutObjectAcl call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update s3object '%s' done", extracted)
	} else {
		renv.Log().Verbose("update s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateS3object) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("s3object"), nil
}

func (cmd *UpdateS3object) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateScalinggroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateScalinggroup {
	cmd := new(UpdateScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateScalinggroup) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *UpdateScalinggroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateScalinggroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.UpdateAutoScalingGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.UpdateAutoScalingGroup(input)
	renv.Log().ExtraVerbosef("autoscaling.UpdateAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update scalinggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("update scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateScalinggroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalinggroup"), nil
}

func (cmd *UpdateScalinggroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateSecuritygroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateSecuritygroup {
	cmd := new(UpdateSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *UpdateSecuritygroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateSecuritygroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("update securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateStack(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateStack {
	cmd := new(UpdateStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudformation.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateStack) SetApi(api cloudformationiface.CloudFormationAPI) {
	cmd.api = api
}

func (cmd *UpdateStack) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateStack) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.UpdateStackInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.UpdateStackInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.UpdateStack(input)
	renv.Log().ExtraVerbosef("cloudformation.UpdateStack call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update stack '%s' done", extracted)
	} else {
		renv.Log().Verbose("update stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateStack) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("stack"), nil
}

func (cmd *UpdateStack) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateSubnet(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateSubnet {
	cmd := new(UpdateSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateSubnet) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *UpdateSubnet) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateSubnet) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ModifySubnetAttributeInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifySubnetAttributeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ModifySubnetAttribute(input)
	renv.Log().ExtraVerbosef("ec2.ModifySubnetAttribute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update subnet '%s' done", extracted)
	} else {
		renv.Log().Verbose("update subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateSubnet) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("subnet"), nil
}

func (cmd *UpdateSubnet) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateTargetgroup(sess *session.Session, g cloud.GraphAPI, l ...*logger.Logger) *UpdateTargetgroup {
	cmd := new(UpdateTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateTargetgroup) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *UpdateTargetgroup) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateTargetgroup) run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update targetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("update targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateTargetgroup) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("targetgroup"), nil
}

func (cmd *UpdateTargetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}
