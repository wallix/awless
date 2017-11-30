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
	"github.com/wallix/awless/logger"
)

func NewAttachAlarm(sess *session.Session, l ...*logger.Logger) *AttachAlarm {
	cmd := new(AttachAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	return cmd
}

func (cmd *AttachAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *AttachAlarm) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach alarm '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachAlarm) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachAlarm) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *AttachAlarm) ParamsHelp() string {
	return generateParamsHelp("attachalarm", structListParamsKeys(cmd))
}

func (cmd *AttachAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachContainertask(sess *session.Session, l ...*logger.Logger) *AttachContainertask {
	cmd := new(AttachContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	return cmd
}

func (cmd *AttachContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *AttachContainertask) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach containertask '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachContainertask) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachContainertask) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *AttachContainertask) ParamsHelp() string {
	return generateParamsHelp("attachcontainertask", structListParamsKeys(cmd))
}

func (cmd *AttachContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachElasticip(sess *session.Session, l ...*logger.Logger) *AttachElasticip {
	cmd := new(AttachElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *AttachElasticip) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachElasticip) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AssociateAddressInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateAddressInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AssociateAddress(input)
	cmd.logger.ExtraVerbosef("ec2.AssociateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach elasticip '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachElasticip) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachElasticip) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.AssociateAddressInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.AssociateAddressInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AssociateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.AssociateAddress call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: attach elasticip ok")
			return fakeDryRunId("elasticip"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *AttachElasticip) ParamsHelp() string {
	return generateParamsHelp("attachelasticip", structListParamsKeys(cmd))
}

func (cmd *AttachElasticip) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachInstance(sess *session.Session, l ...*logger.Logger) *AttachInstance {
	cmd := new(AttachInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *AttachInstance) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *AttachInstance) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.RegisterTargetsInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.RegisterTargetsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RegisterTargets(input)
	cmd.logger.ExtraVerbosef("elbv2.RegisterTargets call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach instance '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInstance) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachInstance) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instance"), nil
}

func (cmd *AttachInstance) ParamsHelp() string {
	return generateParamsHelp("attachinstance", structListParamsKeys(cmd))
}

func (cmd *AttachInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachInstanceprofile(sess *session.Session, l ...*logger.Logger) *AttachInstanceprofile {
	cmd := new(AttachInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *AttachInstanceprofile) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachInstanceprofile) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach instanceprofile '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInstanceprofile) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachInstanceprofile) ParamsHelp() string {
	return generateParamsHelp("attachinstanceprofile", structListParamsKeys(cmd))
}

func (cmd *AttachInstanceprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachInternetgateway(sess *session.Session, l ...*logger.Logger) *AttachInternetgateway {
	cmd := new(AttachInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *AttachInternetgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachInternetgateway) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachInternetGatewayInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachInternetGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AttachInternetGateway(input)
	cmd.logger.ExtraVerbosef("ec2.AttachInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach internetgateway '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInternetgateway) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachInternetgateway) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.AttachInternetGatewayInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.AttachInternetGatewayInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AttachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.AttachInternetGateway call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: attach internetgateway ok")
			return fakeDryRunId("internetgateway"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *AttachInternetgateway) ParamsHelp() string {
	return generateParamsHelp("attachinternetgateway", structListParamsKeys(cmd))
}

func (cmd *AttachInternetgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachMfadevice(sess *session.Session, l ...*logger.Logger) *AttachMfadevice {
	cmd := new(AttachMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *AttachMfadevice) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *AttachMfadevice) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.EnableMFADeviceInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.EnableMFADeviceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.EnableMFADevice(input)
	cmd.logger.ExtraVerbosef("iam.EnableMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach mfadevice '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachMfadevice) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachMfadevice) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("mfadevice"), nil
}

func (cmd *AttachMfadevice) ParamsHelp() string {
	return generateParamsHelp("attachmfadevice", structListParamsKeys(cmd))
}

func (cmd *AttachMfadevice) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachNetworkinterface(sess *session.Session, l ...*logger.Logger) *AttachNetworkinterface {
	cmd := new(AttachNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *AttachNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachNetworkinterface) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachNetworkInterfaceInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachNetworkInterfaceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AttachNetworkInterface(input)
	cmd.logger.ExtraVerbosef("ec2.AttachNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach networkinterface '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachNetworkinterface) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachNetworkinterface) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.AttachNetworkInterfaceInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.AttachNetworkInterfaceInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AttachNetworkInterface(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.AttachNetworkInterface call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: attach networkinterface ok")
			return fakeDryRunId("networkinterface"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *AttachNetworkinterface) ParamsHelp() string {
	return generateParamsHelp("attachnetworkinterface", structListParamsKeys(cmd))
}

func (cmd *AttachNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachPolicy(sess *session.Session, l ...*logger.Logger) *AttachPolicy {
	cmd := new(AttachPolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *AttachPolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *AttachPolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach policy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachPolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachPolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *AttachPolicy) ParamsHelp() string {
	return generateParamsHelp("attachpolicy", structListParamsKeys(cmd))
}

func (cmd *AttachPolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachRole(sess *session.Session, l ...*logger.Logger) *AttachRole {
	cmd := new(AttachRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *AttachRole) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *AttachRole) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.AddRoleToInstanceProfileInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.AddRoleToInstanceProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AddRoleToInstanceProfile(input)
	cmd.logger.ExtraVerbosef("iam.AddRoleToInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach role '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachRole) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachRole) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("role"), nil
}

func (cmd *AttachRole) ParamsHelp() string {
	return generateParamsHelp("attachrole", structListParamsKeys(cmd))
}

func (cmd *AttachRole) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachRoutetable(sess *session.Session, l ...*logger.Logger) *AttachRoutetable {
	cmd := new(AttachRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *AttachRoutetable) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachRoutetable) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AssociateRouteTableInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateRouteTableInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AssociateRouteTable(input)
	cmd.logger.ExtraVerbosef("ec2.AssociateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach routetable '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachRoutetable) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachRoutetable) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.AssociateRouteTableInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.AssociateRouteTableInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AssociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.AssociateRouteTable call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: attach routetable ok")
			return fakeDryRunId("routetable"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *AttachRoutetable) ParamsHelp() string {
	return generateParamsHelp("attachroutetable", structListParamsKeys(cmd))
}

func (cmd *AttachRoutetable) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachSecuritygroup(sess *session.Session, l ...*logger.Logger) *AttachSecuritygroup {
	cmd := new(AttachSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *AttachSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachSecuritygroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach securitygroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachSecuritygroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachSecuritygroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("securitygroup"), nil
}

func (cmd *AttachSecuritygroup) ParamsHelp() string {
	return generateParamsHelp("attachsecuritygroup", structListParamsKeys(cmd))
}

func (cmd *AttachSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachUser(sess *session.Session, l ...*logger.Logger) *AttachUser {
	cmd := new(AttachUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *AttachUser) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *AttachUser) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.AddUserToGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.AddUserToGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AddUserToGroup(input)
	cmd.logger.ExtraVerbosef("iam.AddUserToGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach user '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachUser) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachUser) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("user"), nil
}

func (cmd *AttachUser) ParamsHelp() string {
	return generateParamsHelp("attachuser", structListParamsKeys(cmd))
}

func (cmd *AttachUser) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAttachVolume(sess *session.Session, l ...*logger.Logger) *AttachVolume {
	cmd := new(AttachVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *AttachVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *AttachVolume) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachVolumeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachVolumeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AttachVolume(input)
	cmd.logger.ExtraVerbosef("ec2.AttachVolume call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("attach volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("attach volume '%s' done", extracted)
	} else {
		cmd.logger.Verbose("attach volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachVolume) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AttachVolume) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.AttachVolumeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.AttachVolumeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AttachVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.AttachVolume call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: attach volume ok")
			return fakeDryRunId("volume"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *AttachVolume) ParamsHelp() string {
	return generateParamsHelp("attachvolume", structListParamsKeys(cmd))
}

func (cmd *AttachVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewAuthenticateRegistry(sess *session.Session, l ...*logger.Logger) *AuthenticateRegistry {
	cmd := new(AuthenticateRegistry)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecr.New(sess)
	}
	return cmd
}

func (cmd *AuthenticateRegistry) SetApi(api ecriface.ECRAPI) {
	cmd.api = api
}

func (cmd *AuthenticateRegistry) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("authenticate registry: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("authenticate registry '%s' done", extracted)
	} else {
		cmd.logger.Verbose("authenticate registry done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AuthenticateRegistry) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *AuthenticateRegistry) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("registry"), nil
}

func (cmd *AuthenticateRegistry) ParamsHelp() string {
	return generateParamsHelp("authenticateregistry", structListParamsKeys(cmd))
}

func (cmd *AuthenticateRegistry) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckCertificate(sess *session.Session, l ...*logger.Logger) *CheckCertificate {
	cmd := new(CheckCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = acm.New(sess)
	}
	return cmd
}

func (cmd *CheckCertificate) SetApi(api acmiface.ACMAPI) {
	cmd.api = api
}

func (cmd *CheckCertificate) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check certificate '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckCertificate) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckCertificate) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("certificate"), nil
}

func (cmd *CheckCertificate) ParamsHelp() string {
	return generateParamsHelp("checkcertificate", structListParamsKeys(cmd))
}

func (cmd *CheckCertificate) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckDatabase(sess *session.Session, l ...*logger.Logger) *CheckDatabase {
	cmd := new(CheckDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	return cmd
}

func (cmd *CheckDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *CheckDatabase) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check database '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckDatabase) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckDatabase) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *CheckDatabase) ParamsHelp() string {
	return generateParamsHelp("checkdatabase", structListParamsKeys(cmd))
}

func (cmd *CheckDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckDistribution(sess *session.Session, l ...*logger.Logger) *CheckDistribution {
	cmd := new(CheckDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudfront.New(sess)
	}
	return cmd
}

func (cmd *CheckDistribution) SetApi(api cloudfrontiface.CloudFrontAPI) {
	cmd.api = api
}

func (cmd *CheckDistribution) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check distribution '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckDistribution) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckDistribution) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("distribution"), nil
}

func (cmd *CheckDistribution) ParamsHelp() string {
	return generateParamsHelp("checkdistribution", structListParamsKeys(cmd))
}

func (cmd *CheckDistribution) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckInstance(sess *session.Session, l ...*logger.Logger) *CheckInstance {
	cmd := new(CheckInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CheckInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckInstance) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check instance '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckInstance) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckInstance) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instance"), nil
}

func (cmd *CheckInstance) ParamsHelp() string {
	return generateParamsHelp("checkinstance", structListParamsKeys(cmd))
}

func (cmd *CheckInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckLoadbalancer(sess *session.Session, l ...*logger.Logger) *CheckLoadbalancer {
	cmd := new(CheckLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *CheckLoadbalancer) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *CheckLoadbalancer) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check loadbalancer '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckLoadbalancer) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckLoadbalancer) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loadbalancer"), nil
}

func (cmd *CheckLoadbalancer) ParamsHelp() string {
	return generateParamsHelp("checkloadbalancer", structListParamsKeys(cmd))
}

func (cmd *CheckLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckNatgateway(sess *session.Session, l ...*logger.Logger) *CheckNatgateway {
	cmd := new(CheckNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CheckNatgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckNatgateway) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check natgateway '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckNatgateway) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckNatgateway) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("natgateway"), nil
}

func (cmd *CheckNatgateway) ParamsHelp() string {
	return generateParamsHelp("checknatgateway", structListParamsKeys(cmd))
}

func (cmd *CheckNatgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckNetworkinterface(sess *session.Session, l ...*logger.Logger) *CheckNetworkinterface {
	cmd := new(CheckNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CheckNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckNetworkinterface) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check networkinterface '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckNetworkinterface) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckNetworkinterface) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("networkinterface"), nil
}

func (cmd *CheckNetworkinterface) ParamsHelp() string {
	return generateParamsHelp("checknetworkinterface", structListParamsKeys(cmd))
}

func (cmd *CheckNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckScalinggroup(sess *session.Session, l ...*logger.Logger) *CheckScalinggroup {
	cmd := new(CheckScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	return cmd
}

func (cmd *CheckScalinggroup) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *CheckScalinggroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check scalinggroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckScalinggroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckScalinggroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalinggroup"), nil
}

func (cmd *CheckScalinggroup) ParamsHelp() string {
	return generateParamsHelp("checkscalinggroup", structListParamsKeys(cmd))
}

func (cmd *CheckScalinggroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckSecuritygroup(sess *session.Session, l ...*logger.Logger) *CheckSecuritygroup {
	cmd := new(CheckSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CheckSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckSecuritygroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check securitygroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckSecuritygroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckSecuritygroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("securitygroup"), nil
}

func (cmd *CheckSecuritygroup) ParamsHelp() string {
	return generateParamsHelp("checksecuritygroup", structListParamsKeys(cmd))
}

func (cmd *CheckSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCheckVolume(sess *session.Session, l ...*logger.Logger) *CheckVolume {
	cmd := new(CheckVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CheckVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CheckVolume) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("check volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("check volume '%s' done", extracted)
	} else {
		cmd.logger.Verbose("check volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckVolume) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CheckVolume) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("volume"), nil
}

func (cmd *CheckVolume) ParamsHelp() string {
	return generateParamsHelp("checkvolume", structListParamsKeys(cmd))
}

func (cmd *CheckVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCopyImage(sess *session.Session, l ...*logger.Logger) *CopyImage {
	cmd := new(CopyImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CopyImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CopyImage) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CopyImageInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopyImageInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CopyImage(input)
	cmd.logger.ExtraVerbosef("ec2.CopyImage call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("copy image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("copy image '%s' done", extracted)
	} else {
		cmd.logger.Verbose("copy image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CopyImage) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CopyImage) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CopyImageInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CopyImageInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CopyImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CopyImage call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: copy image ok")
			return fakeDryRunId("image"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CopyImage) ParamsHelp() string {
	return generateParamsHelp("copyimage", structListParamsKeys(cmd))
}

func (cmd *CopyImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCopySnapshot(sess *session.Session, l ...*logger.Logger) *CopySnapshot {
	cmd := new(CopySnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CopySnapshot) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CopySnapshot) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CopySnapshotInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopySnapshotInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CopySnapshot(input)
	cmd.logger.ExtraVerbosef("ec2.CopySnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("copy snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("copy snapshot '%s' done", extracted)
	} else {
		cmd.logger.Verbose("copy snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CopySnapshot) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CopySnapshot) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CopySnapshotInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CopySnapshotInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CopySnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CopySnapshot call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: copy snapshot ok")
			return fakeDryRunId("snapshot"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CopySnapshot) ParamsHelp() string {
	return generateParamsHelp("copysnapshot", structListParamsKeys(cmd))
}

func (cmd *CopySnapshot) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateAccesskey(sess *session.Session, l ...*logger.Logger) *CreateAccesskey {
	cmd := new(CreateAccesskey)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *CreateAccesskey) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateAccesskey) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateAccessKeyInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateAccessKeyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateAccessKey(input)
	cmd.logger.ExtraVerbosef("iam.CreateAccessKey call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create accesskey: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create accesskey '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create accesskey done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAccesskey) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateAccesskey) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("accesskey"), nil
}

func (cmd *CreateAccesskey) ParamsHelp() string {
	return generateParamsHelp("createaccesskey", structListParamsKeys(cmd))
}

func (cmd *CreateAccesskey) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateAlarm(sess *session.Session, l ...*logger.Logger) *CreateAlarm {
	cmd := new(CreateAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	return cmd
}

func (cmd *CreateAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *CreateAlarm) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.PutMetricAlarmInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.PutMetricAlarmInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.PutMetricAlarm(input)
	cmd.logger.ExtraVerbosef("cloudwatch.PutMetricAlarm call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create alarm '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAlarm) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateAlarm) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *CreateAlarm) ParamsHelp() string {
	return generateParamsHelp("createalarm", structListParamsKeys(cmd))
}

func (cmd *CreateAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateAppscalingpolicy(sess *session.Session, l ...*logger.Logger) *CreateAppscalingpolicy {
	cmd := new(CreateAppscalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = applicationautoscaling.New(sess)
	}
	return cmd
}

func (cmd *CreateAppscalingpolicy) SetApi(api applicationautoscalingiface.ApplicationAutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateAppscalingpolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.PutScalingPolicyInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.PutScalingPolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.PutScalingPolicy(input)
	cmd.logger.ExtraVerbosef("applicationautoscaling.PutScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create appscalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create appscalingpolicy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create appscalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAppscalingpolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateAppscalingpolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("appscalingpolicy"), nil
}

func (cmd *CreateAppscalingpolicy) ParamsHelp() string {
	return generateParamsHelp("createappscalingpolicy", structListParamsKeys(cmd))
}

func (cmd *CreateAppscalingpolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateAppscalingtarget(sess *session.Session, l ...*logger.Logger) *CreateAppscalingtarget {
	cmd := new(CreateAppscalingtarget)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = applicationautoscaling.New(sess)
	}
	return cmd
}

func (cmd *CreateAppscalingtarget) SetApi(api applicationautoscalingiface.ApplicationAutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateAppscalingtarget) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.RegisterScalableTargetInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.RegisterScalableTargetInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RegisterScalableTarget(input)
	cmd.logger.ExtraVerbosef("applicationautoscaling.RegisterScalableTarget call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create appscalingtarget: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create appscalingtarget '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create appscalingtarget done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAppscalingtarget) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateAppscalingtarget) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("appscalingtarget"), nil
}

func (cmd *CreateAppscalingtarget) ParamsHelp() string {
	return generateParamsHelp("createappscalingtarget", structListParamsKeys(cmd))
}

func (cmd *CreateAppscalingtarget) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateBucket(sess *session.Session, l ...*logger.Logger) *CreateBucket {
	cmd := new(CreateBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	return cmd
}

func (cmd *CreateBucket) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *CreateBucket) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.CreateBucketInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.CreateBucketInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateBucket(input)
	cmd.logger.ExtraVerbosef("s3.CreateBucket call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create bucket '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateBucket) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateBucket) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("bucket"), nil
}

func (cmd *CreateBucket) ParamsHelp() string {
	return generateParamsHelp("createbucket", structListParamsKeys(cmd))
}

func (cmd *CreateBucket) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateCertificate(sess *session.Session, l ...*logger.Logger) *CreateCertificate {
	cmd := new(CreateCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = acm.New(sess)
	}
	return cmd
}

func (cmd *CreateCertificate) SetApi(api acmiface.ACMAPI) {
	cmd.api = api
}

func (cmd *CreateCertificate) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create certificate '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateCertificate) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateCertificate) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("certificate"), nil
}

func (cmd *CreateCertificate) ParamsHelp() string {
	return generateParamsHelp("createcertificate", structListParamsKeys(cmd))
}

func (cmd *CreateCertificate) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateContainercluster(sess *session.Session, l ...*logger.Logger) *CreateContainercluster {
	cmd := new(CreateContainercluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	return cmd
}

func (cmd *CreateContainercluster) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *CreateContainercluster) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.CreateClusterInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.CreateClusterInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateCluster(input)
	cmd.logger.ExtraVerbosef("ecs.CreateCluster call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create containercluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create containercluster '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create containercluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateContainercluster) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateContainercluster) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containercluster"), nil
}

func (cmd *CreateContainercluster) ParamsHelp() string {
	return generateParamsHelp("createcontainercluster", structListParamsKeys(cmd))
}

func (cmd *CreateContainercluster) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateDatabase(sess *session.Session, l ...*logger.Logger) *CreateDatabase {
	cmd := new(CreateDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	return cmd
}

func (cmd *CreateDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *CreateDatabase) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create database '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDatabase) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateDatabase) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *CreateDatabase) ParamsHelp() string {
	return generateParamsHelp("createdatabase", structListParamsKeys(cmd))
}

func (cmd *CreateDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateDbsubnetgroup(sess *session.Session, l ...*logger.Logger) *CreateDbsubnetgroup {
	cmd := new(CreateDbsubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	return cmd
}

func (cmd *CreateDbsubnetgroup) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *CreateDbsubnetgroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.CreateDBSubnetGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.CreateDBSubnetGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateDBSubnetGroup(input)
	cmd.logger.ExtraVerbosef("rds.CreateDBSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create dbsubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create dbsubnetgroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create dbsubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDbsubnetgroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateDbsubnetgroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("dbsubnetgroup"), nil
}

func (cmd *CreateDbsubnetgroup) ParamsHelp() string {
	return generateParamsHelp("createdbsubnetgroup", structListParamsKeys(cmd))
}

func (cmd *CreateDbsubnetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateDistribution(sess *session.Session, l ...*logger.Logger) *CreateDistribution {
	cmd := new(CreateDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudfront.New(sess)
	}
	return cmd
}

func (cmd *CreateDistribution) SetApi(api cloudfrontiface.CloudFrontAPI) {
	cmd.api = api
}

func (cmd *CreateDistribution) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create distribution '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDistribution) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateDistribution) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("distribution"), nil
}

func (cmd *CreateDistribution) ParamsHelp() string {
	return generateParamsHelp("createdistribution", structListParamsKeys(cmd))
}

func (cmd *CreateDistribution) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateElasticip(sess *session.Session, l ...*logger.Logger) *CreateElasticip {
	cmd := new(CreateElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateElasticip) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateElasticip) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AllocateAddressInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AllocateAddressInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.AllocateAddress(input)
	cmd.logger.ExtraVerbosef("ec2.AllocateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create elasticip '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateElasticip) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateElasticip) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.AllocateAddressInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.AllocateAddressInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.AllocateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.AllocateAddress call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create elasticip ok")
			return fakeDryRunId("elasticip"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateElasticip) ParamsHelp() string {
	return generateParamsHelp("createelasticip", structListParamsKeys(cmd))
}

func (cmd *CreateElasticip) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateFunction(sess *session.Session, l ...*logger.Logger) *CreateFunction {
	cmd := new(CreateFunction)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = lambda.New(sess)
	}
	return cmd
}

func (cmd *CreateFunction) SetApi(api lambdaiface.LambdaAPI) {
	cmd.api = api
}

func (cmd *CreateFunction) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &lambda.CreateFunctionInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in lambda.CreateFunctionInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateFunction(input)
	cmd.logger.ExtraVerbosef("lambda.CreateFunction call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create function: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create function '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create function done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateFunction) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateFunction) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("function"), nil
}

func (cmd *CreateFunction) ParamsHelp() string {
	return generateParamsHelp("createfunction", structListParamsKeys(cmd))
}

func (cmd *CreateFunction) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateGroup(sess *session.Session, l ...*logger.Logger) *CreateGroup {
	cmd := new(CreateGroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *CreateGroup) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateGroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateGroup(input)
	cmd.logger.ExtraVerbosef("iam.CreateGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create group: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create group '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create group done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateGroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateGroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("group"), nil
}

func (cmd *CreateGroup) ParamsHelp() string {
	return generateParamsHelp("creategroup", structListParamsKeys(cmd))
}

func (cmd *CreateGroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateImage(sess *session.Session, l ...*logger.Logger) *CreateImage {
	cmd := new(CreateImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateImage) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateImageInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateImageInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateImage(input)
	cmd.logger.ExtraVerbosef("ec2.CreateImage call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create image '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateImage) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateImage) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateImageInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateImageInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateImage call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create image ok")
			return fakeDryRunId("image"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateImage) ParamsHelp() string {
	return generateParamsHelp("createimage", structListParamsKeys(cmd))
}

func (cmd *CreateImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateInstance(sess *session.Session, l ...*logger.Logger) *CreateInstance {
	cmd := new(CreateInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateInstance) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.RunInstancesInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RunInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RunInstances(input)
	cmd.logger.ExtraVerbosef("ec2.RunInstances call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create instance '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInstance) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateInstance) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.RunInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.RunInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.RunInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.RunInstances call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateInstance) ParamsHelp() string {
	return generateParamsHelp("createinstance", structListParamsKeys(cmd))
}

func (cmd *CreateInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateInstanceprofile(sess *session.Session, l ...*logger.Logger) *CreateInstanceprofile {
	cmd := new(CreateInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *CreateInstanceprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateInstanceprofile) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateInstanceProfileInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateInstanceProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateInstanceProfile(input)
	cmd.logger.ExtraVerbosef("iam.CreateInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create instanceprofile '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInstanceprofile) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateInstanceprofile) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instanceprofile"), nil
}

func (cmd *CreateInstanceprofile) ParamsHelp() string {
	return generateParamsHelp("createinstanceprofile", structListParamsKeys(cmd))
}

func (cmd *CreateInstanceprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateInternetgateway(sess *session.Session, l ...*logger.Logger) *CreateInternetgateway {
	cmd := new(CreateInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateInternetgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateInternetgateway) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateInternetGatewayInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateInternetGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateInternetGateway(input)
	cmd.logger.ExtraVerbosef("ec2.CreateInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create internetgateway '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInternetgateway) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateInternetgateway) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateInternetGatewayInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateInternetGatewayInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateInternetGateway call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create internetgateway ok")
			return fakeDryRunId("internetgateway"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateInternetgateway) ParamsHelp() string {
	return generateParamsHelp("createinternetgateway", structListParamsKeys(cmd))
}

func (cmd *CreateInternetgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateKeypair(sess *session.Session, l ...*logger.Logger) *CreateKeypair {
	cmd := new(CreateKeypair)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateKeypair) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateKeypair) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ImportKeyPairInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ImportKeyPairInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ImportKeyPair(input)
	cmd.logger.ExtraVerbosef("ec2.ImportKeyPair call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create keypair: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create keypair '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create keypair done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateKeypair) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateKeypair) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("keypair"), nil
}

func (cmd *CreateKeypair) ParamsHelp() string {
	return generateParamsHelp("createkeypair", structListParamsKeys(cmd))
}

func (cmd *CreateKeypair) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateLaunchconfiguration(sess *session.Session, l ...*logger.Logger) *CreateLaunchconfiguration {
	cmd := new(CreateLaunchconfiguration)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	return cmd
}

func (cmd *CreateLaunchconfiguration) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateLaunchconfiguration) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.CreateLaunchConfigurationInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.CreateLaunchConfigurationInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateLaunchConfiguration(input)
	cmd.logger.ExtraVerbosef("autoscaling.CreateLaunchConfiguration call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create launchconfiguration: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create launchconfiguration '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create launchconfiguration done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLaunchconfiguration) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateLaunchconfiguration) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("launchconfiguration"), nil
}

func (cmd *CreateLaunchconfiguration) ParamsHelp() string {
	return generateParamsHelp("createlaunchconfiguration", structListParamsKeys(cmd))
}

func (cmd *CreateLaunchconfiguration) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateListener(sess *session.Session, l ...*logger.Logger) *CreateListener {
	cmd := new(CreateListener)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *CreateListener) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *CreateListener) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateListenerInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateListenerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateListener(input)
	cmd.logger.ExtraVerbosef("elbv2.CreateListener call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create listener: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create listener '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create listener done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateListener) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateListener) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("listener"), nil
}

func (cmd *CreateListener) ParamsHelp() string {
	return generateParamsHelp("createlistener", structListParamsKeys(cmd))
}

func (cmd *CreateListener) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateLoadbalancer(sess *session.Session, l ...*logger.Logger) *CreateLoadbalancer {
	cmd := new(CreateLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *CreateLoadbalancer) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *CreateLoadbalancer) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateLoadBalancerInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateLoadBalancerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateLoadBalancer(input)
	cmd.logger.ExtraVerbosef("elbv2.CreateLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create loadbalancer '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLoadbalancer) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateLoadbalancer) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loadbalancer"), nil
}

func (cmd *CreateLoadbalancer) ParamsHelp() string {
	return generateParamsHelp("createloadbalancer", structListParamsKeys(cmd))
}

func (cmd *CreateLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateLoginprofile(sess *session.Session, l ...*logger.Logger) *CreateLoginprofile {
	cmd := new(CreateLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *CreateLoginprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateLoginprofile) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateLoginProfileInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateLoginProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateLoginProfile(input)
	cmd.logger.ExtraVerbosef("iam.CreateLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create loginprofile '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLoginprofile) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateLoginprofile) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loginprofile"), nil
}

func (cmd *CreateLoginprofile) ParamsHelp() string {
	return generateParamsHelp("createloginprofile", structListParamsKeys(cmd))
}

func (cmd *CreateLoginprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateMfadevice(sess *session.Session, l ...*logger.Logger) *CreateMfadevice {
	cmd := new(CreateMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *CreateMfadevice) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateMfadevice) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create mfadevice '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateMfadevice) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateMfadevice) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("mfadevice"), nil
}

func (cmd *CreateMfadevice) ParamsHelp() string {
	return generateParamsHelp("createmfadevice", structListParamsKeys(cmd))
}

func (cmd *CreateMfadevice) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateNatgateway(sess *session.Session, l ...*logger.Logger) *CreateNatgateway {
	cmd := new(CreateNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateNatgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateNatgateway) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateNatGatewayInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateNatGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateNatGateway(input)
	cmd.logger.ExtraVerbosef("ec2.CreateNatGateway call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create natgateway '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateNatgateway) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateNatgateway) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("natgateway"), nil
}

func (cmd *CreateNatgateway) ParamsHelp() string {
	return generateParamsHelp("createnatgateway", structListParamsKeys(cmd))
}

func (cmd *CreateNatgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateNetworkinterface(sess *session.Session, l ...*logger.Logger) *CreateNetworkinterface {
	cmd := new(CreateNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateNetworkinterface) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateNetworkInterfaceInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateNetworkInterfaceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateNetworkInterface(input)
	cmd.logger.ExtraVerbosef("ec2.CreateNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create networkinterface '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateNetworkinterface) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateNetworkinterface) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateNetworkInterfaceInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateNetworkInterfaceInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateNetworkInterface(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateNetworkInterface call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create networkinterface ok")
			return fakeDryRunId("networkinterface"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateNetworkinterface) ParamsHelp() string {
	return generateParamsHelp("createnetworkinterface", structListParamsKeys(cmd))
}

func (cmd *CreateNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreatePolicy(sess *session.Session, l ...*logger.Logger) *CreatePolicy {
	cmd := new(CreatePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *CreatePolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreatePolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreatePolicyInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreatePolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreatePolicy(input)
	cmd.logger.ExtraVerbosef("iam.CreatePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create policy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreatePolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreatePolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *CreatePolicy) ParamsHelp() string {
	return generateParamsHelp("createpolicy", structListParamsKeys(cmd))
}

func (cmd *CreatePolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateQueue(sess *session.Session, l ...*logger.Logger) *CreateQueue {
	cmd := new(CreateQueue)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sqs.New(sess)
	}
	return cmd
}

func (cmd *CreateQueue) SetApi(api sqsiface.SQSAPI) {
	cmd.api = api
}

func (cmd *CreateQueue) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sqs.CreateQueueInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in sqs.CreateQueueInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateQueue(input)
	cmd.logger.ExtraVerbosef("sqs.CreateQueue call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create queue: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create queue '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create queue done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateQueue) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateQueue) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("queue"), nil
}

func (cmd *CreateQueue) ParamsHelp() string {
	return generateParamsHelp("createqueue", structListParamsKeys(cmd))
}

func (cmd *CreateQueue) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRecord(sess *session.Session, l ...*logger.Logger) *CreateRecord {
	cmd := new(CreateRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	return cmd
}

func (cmd *CreateRecord) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *CreateRecord) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create record '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRecord) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateRecord) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("record"), nil
}

func (cmd *CreateRecord) ParamsHelp() string {
	return generateParamsHelp("createrecord", structListParamsKeys(cmd))
}

func (cmd *CreateRecord) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRepository(sess *session.Session, l ...*logger.Logger) *CreateRepository {
	cmd := new(CreateRepository)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecr.New(sess)
	}
	return cmd
}

func (cmd *CreateRepository) SetApi(api ecriface.ECRAPI) {
	cmd.api = api
}

func (cmd *CreateRepository) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecr.CreateRepositoryInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ecr.CreateRepositoryInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateRepository(input)
	cmd.logger.ExtraVerbosef("ecr.CreateRepository call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create repository: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create repository '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create repository done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRepository) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateRepository) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("repository"), nil
}

func (cmd *CreateRepository) ParamsHelp() string {
	return generateParamsHelp("createrepository", structListParamsKeys(cmd))
}

func (cmd *CreateRepository) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRole(sess *session.Session, l ...*logger.Logger) *CreateRole {
	cmd := new(CreateRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *CreateRole) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateRole) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create role '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRole) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateRole) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("role"), nil
}

func (cmd *CreateRole) ParamsHelp() string {
	return generateParamsHelp("createrole", structListParamsKeys(cmd))
}

func (cmd *CreateRole) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRoute(sess *session.Session, l ...*logger.Logger) *CreateRoute {
	cmd := new(CreateRoute)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateRoute) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateRoute) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateRouteInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateRoute(input)
	cmd.logger.ExtraVerbosef("ec2.CreateRoute call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create route: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create route '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create route done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRoute) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateRoute) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateRouteInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateRouteInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateRoute call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create route ok")
			return fakeDryRunId("route"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateRoute) ParamsHelp() string {
	return generateParamsHelp("createroute", structListParamsKeys(cmd))
}

func (cmd *CreateRoute) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateRoutetable(sess *session.Session, l ...*logger.Logger) *CreateRoutetable {
	cmd := new(CreateRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateRoutetable) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateRoutetable) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateRouteTableInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteTableInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateRouteTable(input)
	cmd.logger.ExtraVerbosef("ec2.CreateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create routetable '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRoutetable) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateRoutetable) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateRouteTableInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateRouteTableInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateRouteTable call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create routetable ok")
			return fakeDryRunId("routetable"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateRoutetable) ParamsHelp() string {
	return generateParamsHelp("createroutetable", structListParamsKeys(cmd))
}

func (cmd *CreateRoutetable) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateS3object(sess *session.Session, l ...*logger.Logger) *CreateS3object {
	cmd := new(CreateS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	return cmd
}

func (cmd *CreateS3object) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *CreateS3object) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create s3object '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateS3object) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateS3object) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("s3object"), nil
}

func (cmd *CreateS3object) ParamsHelp() string {
	return generateParamsHelp("creates3object", structListParamsKeys(cmd))
}

func (cmd *CreateS3object) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateScalinggroup(sess *session.Session, l ...*logger.Logger) *CreateScalinggroup {
	cmd := new(CreateScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	return cmd
}

func (cmd *CreateScalinggroup) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateScalinggroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.CreateAutoScalingGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.CreateAutoScalingGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateAutoScalingGroup(input)
	cmd.logger.ExtraVerbosef("autoscaling.CreateAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create scalinggroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateScalinggroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateScalinggroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalinggroup"), nil
}

func (cmd *CreateScalinggroup) ParamsHelp() string {
	return generateParamsHelp("createscalinggroup", structListParamsKeys(cmd))
}

func (cmd *CreateScalinggroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateScalingpolicy(sess *session.Session, l ...*logger.Logger) *CreateScalingpolicy {
	cmd := new(CreateScalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	return cmd
}

func (cmd *CreateScalingpolicy) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *CreateScalingpolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.PutScalingPolicyInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.PutScalingPolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.PutScalingPolicy(input)
	cmd.logger.ExtraVerbosef("autoscaling.PutScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create scalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create scalingpolicy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create scalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateScalingpolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateScalingpolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalingpolicy"), nil
}

func (cmd *CreateScalingpolicy) ParamsHelp() string {
	return generateParamsHelp("createscalingpolicy", structListParamsKeys(cmd))
}

func (cmd *CreateScalingpolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateSecuritygroup(sess *session.Session, l ...*logger.Logger) *CreateSecuritygroup {
	cmd := new(CreateSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateSecuritygroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSecurityGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSecurityGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateSecurityGroup(input)
	cmd.logger.ExtraVerbosef("ec2.CreateSecurityGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create securitygroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSecuritygroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateSecuritygroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateSecurityGroupInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateSecurityGroupInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateSecurityGroup call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create securitygroup ok")
			return fakeDryRunId("securitygroup"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateSecuritygroup) ParamsHelp() string {
	return generateParamsHelp("createsecuritygroup", structListParamsKeys(cmd))
}

func (cmd *CreateSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateSnapshot(sess *session.Session, l ...*logger.Logger) *CreateSnapshot {
	cmd := new(CreateSnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateSnapshot) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateSnapshot) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSnapshotInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSnapshotInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateSnapshot(input)
	cmd.logger.ExtraVerbosef("ec2.CreateSnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create snapshot '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSnapshot) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateSnapshot) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateSnapshotInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateSnapshotInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateSnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateSnapshot call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create snapshot ok")
			return fakeDryRunId("snapshot"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateSnapshot) ParamsHelp() string {
	return generateParamsHelp("createsnapshot", structListParamsKeys(cmd))
}

func (cmd *CreateSnapshot) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateStack(sess *session.Session, l ...*logger.Logger) *CreateStack {
	cmd := new(CreateStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudformation.New(sess)
	}
	return cmd
}

func (cmd *CreateStack) SetApi(api cloudformationiface.CloudFormationAPI) {
	cmd.api = api
}

func (cmd *CreateStack) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.CreateStackInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.CreateStackInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateStack(input)
	cmd.logger.ExtraVerbosef("cloudformation.CreateStack call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create stack '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateStack) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateStack) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("stack"), nil
}

func (cmd *CreateStack) ParamsHelp() string {
	return generateParamsHelp("createstack", structListParamsKeys(cmd))
}

func (cmd *CreateStack) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateSubnet(sess *session.Session, l ...*logger.Logger) *CreateSubnet {
	cmd := new(CreateSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateSubnet) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateSubnet) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSubnetInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSubnetInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateSubnet(input)
	cmd.logger.ExtraVerbosef("ec2.CreateSubnet call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create subnet '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSubnet) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateSubnet) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateSubnetInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateSubnetInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateSubnet call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create subnet ok")
			return fakeDryRunId("subnet"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateSubnet) ParamsHelp() string {
	return generateParamsHelp("createsubnet", structListParamsKeys(cmd))
}

func (cmd *CreateSubnet) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateSubscription(sess *session.Session, l ...*logger.Logger) *CreateSubscription {
	cmd := new(CreateSubscription)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sns.New(sess)
	}
	return cmd
}

func (cmd *CreateSubscription) SetApi(api snsiface.SNSAPI) {
	cmd.api = api
}

func (cmd *CreateSubscription) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.SubscribeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.SubscribeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.Subscribe(input)
	cmd.logger.ExtraVerbosef("sns.Subscribe call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create subscription: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create subscription '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create subscription done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSubscription) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateSubscription) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("subscription"), nil
}

func (cmd *CreateSubscription) ParamsHelp() string {
	return generateParamsHelp("createsubscription", structListParamsKeys(cmd))
}

func (cmd *CreateSubscription) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateTag(sess *session.Session, l ...*logger.Logger) *CreateTag {
	cmd := new(CreateTag)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateTag) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateTag) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create tag: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create tag '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create tag done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTag) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateTag) ParamsHelp() string {
	return generateParamsHelp("createtag", structListParamsKeys(cmd))
}

func (cmd *CreateTag) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateTargetgroup(sess *session.Session, l ...*logger.Logger) *CreateTargetgroup {
	cmd := new(CreateTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *CreateTargetgroup) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *CreateTargetgroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateTargetGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateTargetGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateTargetGroup(input)
	cmd.logger.ExtraVerbosef("elbv2.CreateTargetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create targetgroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTargetgroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateTargetgroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("targetgroup"), nil
}

func (cmd *CreateTargetgroup) ParamsHelp() string {
	return generateParamsHelp("createtargetgroup", structListParamsKeys(cmd))
}

func (cmd *CreateTargetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateTopic(sess *session.Session, l ...*logger.Logger) *CreateTopic {
	cmd := new(CreateTopic)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sns.New(sess)
	}
	return cmd
}

func (cmd *CreateTopic) SetApi(api snsiface.SNSAPI) {
	cmd.api = api
}

func (cmd *CreateTopic) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.CreateTopicInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.CreateTopicInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateTopic(input)
	cmd.logger.ExtraVerbosef("sns.CreateTopic call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create topic: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create topic '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create topic done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTopic) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateTopic) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("topic"), nil
}

func (cmd *CreateTopic) ParamsHelp() string {
	return generateParamsHelp("createtopic", structListParamsKeys(cmd))
}

func (cmd *CreateTopic) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateUser(sess *session.Session, l ...*logger.Logger) *CreateUser {
	cmd := new(CreateUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *CreateUser) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *CreateUser) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateUserInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateUserInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateUser(input)
	cmd.logger.ExtraVerbosef("iam.CreateUser call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create user '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateUser) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateUser) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("user"), nil
}

func (cmd *CreateUser) ParamsHelp() string {
	return generateParamsHelp("createuser", structListParamsKeys(cmd))
}

func (cmd *CreateUser) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateVolume(sess *session.Session, l ...*logger.Logger) *CreateVolume {
	cmd := new(CreateVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateVolume) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateVolumeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVolumeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateVolume(input)
	cmd.logger.ExtraVerbosef("ec2.CreateVolume call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create volume '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateVolume) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateVolume) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateVolumeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateVolumeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateVolume call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create volume ok")
			return fakeDryRunId("volume"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateVolume) ParamsHelp() string {
	return generateParamsHelp("createvolume", structListParamsKeys(cmd))
}

func (cmd *CreateVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateVpc(sess *session.Session, l ...*logger.Logger) *CreateVpc {
	cmd := new(CreateVpc)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *CreateVpc) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *CreateVpc) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateVpcInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVpcInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateVpc(input)
	cmd.logger.ExtraVerbosef("ec2.CreateVpc call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create vpc: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create vpc '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create vpc done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateVpc) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateVpc) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateVpcInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateVpcInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.CreateVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateVpc call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create vpc ok")
			return fakeDryRunId("vpc"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateVpc) ParamsHelp() string {
	return generateParamsHelp("createvpc", structListParamsKeys(cmd))
}

func (cmd *CreateVpc) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewCreateZone(sess *session.Session, l ...*logger.Logger) *CreateZone {
	cmd := new(CreateZone)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	return cmd
}

func (cmd *CreateZone) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *CreateZone) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &route53.CreateHostedZoneInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in route53.CreateHostedZoneInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreateHostedZone(input)
	cmd.logger.ExtraVerbosef("route53.CreateHostedZone call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("create zone: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("create zone '%s' done", extracted)
	} else {
		cmd.logger.Verbose("create zone done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateZone) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *CreateZone) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("zone"), nil
}

func (cmd *CreateZone) ParamsHelp() string {
	return generateParamsHelp("createzone", structListParamsKeys(cmd))
}

func (cmd *CreateZone) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteAccesskey(sess *session.Session, l ...*logger.Logger) *DeleteAccesskey {
	cmd := new(DeleteAccesskey)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DeleteAccesskey) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteAccesskey) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteAccessKeyInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteAccessKeyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAccessKey(input)
	cmd.logger.ExtraVerbosef("iam.DeleteAccessKey call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete accesskey: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete accesskey '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete accesskey done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAccesskey) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteAccesskey) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("accesskey"), nil
}

func (cmd *DeleteAccesskey) ParamsHelp() string {
	return generateParamsHelp("deleteaccesskey", structListParamsKeys(cmd))
}

func (cmd *DeleteAccesskey) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteAlarm(sess *session.Session, l ...*logger.Logger) *DeleteAlarm {
	cmd := new(DeleteAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	return cmd
}

func (cmd *DeleteAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *DeleteAlarm) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.DeleteAlarmsInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.DeleteAlarmsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAlarms(input)
	cmd.logger.ExtraVerbosef("cloudwatch.DeleteAlarms call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete alarm '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAlarm) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteAlarm) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *DeleteAlarm) ParamsHelp() string {
	return generateParamsHelp("deletealarm", structListParamsKeys(cmd))
}

func (cmd *DeleteAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteAppscalingpolicy(sess *session.Session, l ...*logger.Logger) *DeleteAppscalingpolicy {
	cmd := new(DeleteAppscalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = applicationautoscaling.New(sess)
	}
	return cmd
}

func (cmd *DeleteAppscalingpolicy) SetApi(api applicationautoscalingiface.ApplicationAutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteAppscalingpolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.DeleteScalingPolicyInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.DeleteScalingPolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteScalingPolicy(input)
	cmd.logger.ExtraVerbosef("applicationautoscaling.DeleteScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete appscalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete appscalingpolicy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete appscalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAppscalingpolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteAppscalingpolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("appscalingpolicy"), nil
}

func (cmd *DeleteAppscalingpolicy) ParamsHelp() string {
	return generateParamsHelp("deleteappscalingpolicy", structListParamsKeys(cmd))
}

func (cmd *DeleteAppscalingpolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteAppscalingtarget(sess *session.Session, l ...*logger.Logger) *DeleteAppscalingtarget {
	cmd := new(DeleteAppscalingtarget)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = applicationautoscaling.New(sess)
	}
	return cmd
}

func (cmd *DeleteAppscalingtarget) SetApi(api applicationautoscalingiface.ApplicationAutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteAppscalingtarget) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.DeregisterScalableTargetInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.DeregisterScalableTargetInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeregisterScalableTarget(input)
	cmd.logger.ExtraVerbosef("applicationautoscaling.DeregisterScalableTarget call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete appscalingtarget: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete appscalingtarget '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete appscalingtarget done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAppscalingtarget) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteAppscalingtarget) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("appscalingtarget"), nil
}

func (cmd *DeleteAppscalingtarget) ParamsHelp() string {
	return generateParamsHelp("deleteappscalingtarget", structListParamsKeys(cmd))
}

func (cmd *DeleteAppscalingtarget) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteBucket(sess *session.Session, l ...*logger.Logger) *DeleteBucket {
	cmd := new(DeleteBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	return cmd
}

func (cmd *DeleteBucket) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *DeleteBucket) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.DeleteBucketInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.DeleteBucketInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteBucket(input)
	cmd.logger.ExtraVerbosef("s3.DeleteBucket call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete bucket '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteBucket) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteBucket) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("bucket"), nil
}

func (cmd *DeleteBucket) ParamsHelp() string {
	return generateParamsHelp("deletebucket", structListParamsKeys(cmd))
}

func (cmd *DeleteBucket) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteCertificate(sess *session.Session, l ...*logger.Logger) *DeleteCertificate {
	cmd := new(DeleteCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = acm.New(sess)
	}
	return cmd
}

func (cmd *DeleteCertificate) SetApi(api acmiface.ACMAPI) {
	cmd.api = api
}

func (cmd *DeleteCertificate) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &acm.DeleteCertificateInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in acm.DeleteCertificateInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCertificate(input)
	cmd.logger.ExtraVerbosef("acm.DeleteCertificate call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete certificate '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteCertificate) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteCertificate) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("certificate"), nil
}

func (cmd *DeleteCertificate) ParamsHelp() string {
	return generateParamsHelp("deletecertificate", structListParamsKeys(cmd))
}

func (cmd *DeleteCertificate) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteContainercluster(sess *session.Session, l ...*logger.Logger) *DeleteContainercluster {
	cmd := new(DeleteContainercluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	return cmd
}

func (cmd *DeleteContainercluster) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *DeleteContainercluster) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.DeleteClusterInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.DeleteClusterInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCluster(input)
	cmd.logger.ExtraVerbosef("ecs.DeleteCluster call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete containercluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete containercluster '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete containercluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteContainercluster) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteContainercluster) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containercluster"), nil
}

func (cmd *DeleteContainercluster) ParamsHelp() string {
	return generateParamsHelp("deletecontainercluster", structListParamsKeys(cmd))
}

func (cmd *DeleteContainercluster) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteContainertask(sess *session.Session, l ...*logger.Logger) *DeleteContainertask {
	cmd := new(DeleteContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	return cmd
}

func (cmd *DeleteContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *DeleteContainertask) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete containertask '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteContainertask) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteContainertask) ParamsHelp() string {
	return generateParamsHelp("deletecontainertask", structListParamsKeys(cmd))
}

func (cmd *DeleteContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteDatabase(sess *session.Session, l ...*logger.Logger) *DeleteDatabase {
	cmd := new(DeleteDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	return cmd
}

func (cmd *DeleteDatabase) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *DeleteDatabase) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.DeleteDBInstanceInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.DeleteDBInstanceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteDBInstance(input)
	cmd.logger.ExtraVerbosef("rds.DeleteDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete database '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDatabase) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteDatabase) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("database"), nil
}

func (cmd *DeleteDatabase) ParamsHelp() string {
	return generateParamsHelp("deletedatabase", structListParamsKeys(cmd))
}

func (cmd *DeleteDatabase) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteDbsubnetgroup(sess *session.Session, l ...*logger.Logger) *DeleteDbsubnetgroup {
	cmd := new(DeleteDbsubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = rds.New(sess)
	}
	return cmd
}

func (cmd *DeleteDbsubnetgroup) SetApi(api rdsiface.RDSAPI) {
	cmd.api = api
}

func (cmd *DeleteDbsubnetgroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.DeleteDBSubnetGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.DeleteDBSubnetGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteDBSubnetGroup(input)
	cmd.logger.ExtraVerbosef("rds.DeleteDBSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete dbsubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete dbsubnetgroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete dbsubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDbsubnetgroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteDbsubnetgroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("dbsubnetgroup"), nil
}

func (cmd *DeleteDbsubnetgroup) ParamsHelp() string {
	return generateParamsHelp("deletedbsubnetgroup", structListParamsKeys(cmd))
}

func (cmd *DeleteDbsubnetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteDistribution(sess *session.Session, l ...*logger.Logger) *DeleteDistribution {
	cmd := new(DeleteDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudfront.New(sess)
	}
	return cmd
}

func (cmd *DeleteDistribution) SetApi(api cloudfrontiface.CloudFrontAPI) {
	cmd.api = api
}

func (cmd *DeleteDistribution) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete distribution '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDistribution) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteDistribution) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("distribution"), nil
}

func (cmd *DeleteDistribution) ParamsHelp() string {
	return generateParamsHelp("deletedistribution", structListParamsKeys(cmd))
}

func (cmd *DeleteDistribution) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteElasticip(sess *session.Session, l ...*logger.Logger) *DeleteElasticip {
	cmd := new(DeleteElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteElasticip) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteElasticip) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ReleaseAddressInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ReleaseAddressInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ReleaseAddress(input)
	cmd.logger.ExtraVerbosef("ec2.ReleaseAddress call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete elasticip '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteElasticip) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteElasticip) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.ReleaseAddressInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.ReleaseAddressInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.ReleaseAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.ReleaseAddress call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete elasticip ok")
			return fakeDryRunId("elasticip"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteElasticip) ParamsHelp() string {
	return generateParamsHelp("deleteelasticip", structListParamsKeys(cmd))
}

func (cmd *DeleteElasticip) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteFunction(sess *session.Session, l ...*logger.Logger) *DeleteFunction {
	cmd := new(DeleteFunction)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = lambda.New(sess)
	}
	return cmd
}

func (cmd *DeleteFunction) SetApi(api lambdaiface.LambdaAPI) {
	cmd.api = api
}

func (cmd *DeleteFunction) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &lambda.DeleteFunctionInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in lambda.DeleteFunctionInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteFunction(input)
	cmd.logger.ExtraVerbosef("lambda.DeleteFunction call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete function: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete function '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete function done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteFunction) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteFunction) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("function"), nil
}

func (cmd *DeleteFunction) ParamsHelp() string {
	return generateParamsHelp("deletefunction", structListParamsKeys(cmd))
}

func (cmd *DeleteFunction) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteGroup(sess *session.Session, l ...*logger.Logger) *DeleteGroup {
	cmd := new(DeleteGroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DeleteGroup) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteGroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteGroup(input)
	cmd.logger.ExtraVerbosef("iam.DeleteGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete group: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete group '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete group done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteGroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteGroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("group"), nil
}

func (cmd *DeleteGroup) ParamsHelp() string {
	return generateParamsHelp("deletegroup", structListParamsKeys(cmd))
}

func (cmd *DeleteGroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteImage(sess *session.Session, l ...*logger.Logger) *DeleteImage {
	cmd := new(DeleteImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteImage) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete image '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteImage) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteImage) ParamsHelp() string {
	return generateParamsHelp("deleteimage", structListParamsKeys(cmd))
}

func (cmd *DeleteImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteInstance(sess *session.Session, l ...*logger.Logger) *DeleteInstance {
	cmd := new(DeleteInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteInstance) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.TerminateInstancesInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.TerminateInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.TerminateInstances(input)
	cmd.logger.ExtraVerbosef("ec2.TerminateInstances call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete instance '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInstance) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteInstance) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.TerminateInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.TerminateInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.TerminateInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.TerminateInstances call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteInstance) ParamsHelp() string {
	return generateParamsHelp("deleteinstance", structListParamsKeys(cmd))
}

func (cmd *DeleteInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteInstanceprofile(sess *session.Session, l ...*logger.Logger) *DeleteInstanceprofile {
	cmd := new(DeleteInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DeleteInstanceprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteInstanceprofile) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteInstanceProfileInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteInstanceProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteInstanceProfile(input)
	cmd.logger.ExtraVerbosef("iam.DeleteInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete instanceprofile '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInstanceprofile) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteInstanceprofile) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instanceprofile"), nil
}

func (cmd *DeleteInstanceprofile) ParamsHelp() string {
	return generateParamsHelp("deleteinstanceprofile", structListParamsKeys(cmd))
}

func (cmd *DeleteInstanceprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteInternetgateway(sess *session.Session, l ...*logger.Logger) *DeleteInternetgateway {
	cmd := new(DeleteInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteInternetgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteInternetgateway) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteInternetGatewayInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteInternetGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteInternetGateway(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete internetgateway '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInternetgateway) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteInternetgateway) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteInternetGatewayInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteInternetGatewayInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteInternetGateway call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete internetgateway ok")
			return fakeDryRunId("internetgateway"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteInternetgateway) ParamsHelp() string {
	return generateParamsHelp("deleteinternetgateway", structListParamsKeys(cmd))
}

func (cmd *DeleteInternetgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteKeypair(sess *session.Session, l ...*logger.Logger) *DeleteKeypair {
	cmd := new(DeleteKeypair)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteKeypair) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteKeypair) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteKeyPairInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteKeyPairInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteKeyPair(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteKeyPair call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete keypair: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete keypair '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete keypair done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteKeypair) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteKeypair) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteKeyPairInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteKeyPairInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteKeyPair(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteKeyPair call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete keypair ok")
			return fakeDryRunId("keypair"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteKeypair) ParamsHelp() string {
	return generateParamsHelp("deletekeypair", structListParamsKeys(cmd))
}

func (cmd *DeleteKeypair) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteLaunchconfiguration(sess *session.Session, l ...*logger.Logger) *DeleteLaunchconfiguration {
	cmd := new(DeleteLaunchconfiguration)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	return cmd
}

func (cmd *DeleteLaunchconfiguration) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteLaunchconfiguration) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeleteLaunchConfigurationInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeleteLaunchConfigurationInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLaunchConfiguration(input)
	cmd.logger.ExtraVerbosef("autoscaling.DeleteLaunchConfiguration call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete launchconfiguration: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete launchconfiguration '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete launchconfiguration done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLaunchconfiguration) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteLaunchconfiguration) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("launchconfiguration"), nil
}

func (cmd *DeleteLaunchconfiguration) ParamsHelp() string {
	return generateParamsHelp("deletelaunchconfiguration", structListParamsKeys(cmd))
}

func (cmd *DeleteLaunchconfiguration) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteListener(sess *session.Session, l ...*logger.Logger) *DeleteListener {
	cmd := new(DeleteListener)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *DeleteListener) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *DeleteListener) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteListenerInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteListenerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteListener(input)
	cmd.logger.ExtraVerbosef("elbv2.DeleteListener call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete listener: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete listener '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete listener done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteListener) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteListener) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("listener"), nil
}

func (cmd *DeleteListener) ParamsHelp() string {
	return generateParamsHelp("deletelistener", structListParamsKeys(cmd))
}

func (cmd *DeleteListener) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteLoadbalancer(sess *session.Session, l ...*logger.Logger) *DeleteLoadbalancer {
	cmd := new(DeleteLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *DeleteLoadbalancer) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *DeleteLoadbalancer) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteLoadBalancerInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteLoadBalancerInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLoadBalancer(input)
	cmd.logger.ExtraVerbosef("elbv2.DeleteLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete loadbalancer '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLoadbalancer) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteLoadbalancer) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loadbalancer"), nil
}

func (cmd *DeleteLoadbalancer) ParamsHelp() string {
	return generateParamsHelp("deleteloadbalancer", structListParamsKeys(cmd))
}

func (cmd *DeleteLoadbalancer) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteLoginprofile(sess *session.Session, l ...*logger.Logger) *DeleteLoginprofile {
	cmd := new(DeleteLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DeleteLoginprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteLoginprofile) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteLoginProfileInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteLoginProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLoginProfile(input)
	cmd.logger.ExtraVerbosef("iam.DeleteLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete loginprofile '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLoginprofile) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteLoginprofile) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loginprofile"), nil
}

func (cmd *DeleteLoginprofile) ParamsHelp() string {
	return generateParamsHelp("deleteloginprofile", structListParamsKeys(cmd))
}

func (cmd *DeleteLoginprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteMfadevice(sess *session.Session, l ...*logger.Logger) *DeleteMfadevice {
	cmd := new(DeleteMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DeleteMfadevice) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteMfadevice) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteVirtualMFADeviceInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteVirtualMFADeviceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVirtualMFADevice(input)
	cmd.logger.ExtraVerbosef("iam.DeleteVirtualMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete mfadevice '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteMfadevice) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteMfadevice) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("mfadevice"), nil
}

func (cmd *DeleteMfadevice) ParamsHelp() string {
	return generateParamsHelp("deletemfadevice", structListParamsKeys(cmd))
}

func (cmd *DeleteMfadevice) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteNatgateway(sess *session.Session, l ...*logger.Logger) *DeleteNatgateway {
	cmd := new(DeleteNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteNatgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteNatgateway) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteNatGatewayInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteNatGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteNatGateway(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteNatGateway call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete natgateway '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteNatgateway) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteNatgateway) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("natgateway"), nil
}

func (cmd *DeleteNatgateway) ParamsHelp() string {
	return generateParamsHelp("deletenatgateway", structListParamsKeys(cmd))
}

func (cmd *DeleteNatgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteNetworkinterface(sess *session.Session, l ...*logger.Logger) *DeleteNetworkinterface {
	cmd := new(DeleteNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteNetworkinterface) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteNetworkInterfaceInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteNetworkInterfaceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteNetworkInterface(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete networkinterface '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteNetworkinterface) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteNetworkinterface) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteNetworkInterfaceInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteNetworkInterfaceInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteNetworkInterface(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteNetworkInterface call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete networkinterface ok")
			return fakeDryRunId("networkinterface"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteNetworkinterface) ParamsHelp() string {
	return generateParamsHelp("deletenetworkinterface", structListParamsKeys(cmd))
}

func (cmd *DeleteNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeletePolicy(sess *session.Session, l ...*logger.Logger) *DeletePolicy {
	cmd := new(DeletePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DeletePolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeletePolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeletePolicyInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeletePolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeletePolicy(input)
	cmd.logger.ExtraVerbosef("iam.DeletePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete policy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeletePolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeletePolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *DeletePolicy) ParamsHelp() string {
	return generateParamsHelp("deletepolicy", structListParamsKeys(cmd))
}

func (cmd *DeletePolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteQueue(sess *session.Session, l ...*logger.Logger) *DeleteQueue {
	cmd := new(DeleteQueue)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sqs.New(sess)
	}
	return cmd
}

func (cmd *DeleteQueue) SetApi(api sqsiface.SQSAPI) {
	cmd.api = api
}

func (cmd *DeleteQueue) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sqs.DeleteQueueInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in sqs.DeleteQueueInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteQueue(input)
	cmd.logger.ExtraVerbosef("sqs.DeleteQueue call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete queue: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete queue '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete queue done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteQueue) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteQueue) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("queue"), nil
}

func (cmd *DeleteQueue) ParamsHelp() string {
	return generateParamsHelp("deletequeue", structListParamsKeys(cmd))
}

func (cmd *DeleteQueue) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRecord(sess *session.Session, l ...*logger.Logger) *DeleteRecord {
	cmd := new(DeleteRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	return cmd
}

func (cmd *DeleteRecord) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *DeleteRecord) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete record '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRecord) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteRecord) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("record"), nil
}

func (cmd *DeleteRecord) ParamsHelp() string {
	return generateParamsHelp("deleterecord", structListParamsKeys(cmd))
}

func (cmd *DeleteRecord) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRepository(sess *session.Session, l ...*logger.Logger) *DeleteRepository {
	cmd := new(DeleteRepository)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecr.New(sess)
	}
	return cmd
}

func (cmd *DeleteRepository) SetApi(api ecriface.ECRAPI) {
	cmd.api = api
}

func (cmd *DeleteRepository) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecr.DeleteRepositoryInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ecr.DeleteRepositoryInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRepository(input)
	cmd.logger.ExtraVerbosef("ecr.DeleteRepository call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete repository: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete repository '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete repository done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRepository) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteRepository) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("repository"), nil
}

func (cmd *DeleteRepository) ParamsHelp() string {
	return generateParamsHelp("deleterepository", structListParamsKeys(cmd))
}

func (cmd *DeleteRepository) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRole(sess *session.Session, l ...*logger.Logger) *DeleteRole {
	cmd := new(DeleteRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DeleteRole) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteRole) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete role '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRole) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteRole) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("role"), nil
}

func (cmd *DeleteRole) ParamsHelp() string {
	return generateParamsHelp("deleterole", structListParamsKeys(cmd))
}

func (cmd *DeleteRole) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRoute(sess *session.Session, l ...*logger.Logger) *DeleteRoute {
	cmd := new(DeleteRoute)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteRoute) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteRoute) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteRouteInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRoute(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteRoute call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete route: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete route '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete route done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRoute) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteRoute) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteRouteInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteRouteInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteRoute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteRoute call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete route ok")
			return fakeDryRunId("route"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteRoute) ParamsHelp() string {
	return generateParamsHelp("deleteroute", structListParamsKeys(cmd))
}

func (cmd *DeleteRoute) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteRoutetable(sess *session.Session, l ...*logger.Logger) *DeleteRoutetable {
	cmd := new(DeleteRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteRoutetable) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteRoutetable) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteRouteTableInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteTableInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRouteTable(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete routetable '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRoutetable) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteRoutetable) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteRouteTableInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteRouteTableInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteRouteTable call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete routetable ok")
			return fakeDryRunId("routetable"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteRoutetable) ParamsHelp() string {
	return generateParamsHelp("deleteroutetable", structListParamsKeys(cmd))
}

func (cmd *DeleteRoutetable) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteS3object(sess *session.Session, l ...*logger.Logger) *DeleteS3object {
	cmd := new(DeleteS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	return cmd
}

func (cmd *DeleteS3object) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *DeleteS3object) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.DeleteObjectInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.DeleteObjectInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteObject(input)
	cmd.logger.ExtraVerbosef("s3.DeleteObject call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete s3object '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteS3object) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteS3object) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("s3object"), nil
}

func (cmd *DeleteS3object) ParamsHelp() string {
	return generateParamsHelp("deletes3object", structListParamsKeys(cmd))
}

func (cmd *DeleteS3object) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteScalinggroup(sess *session.Session, l ...*logger.Logger) *DeleteScalinggroup {
	cmd := new(DeleteScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	return cmd
}

func (cmd *DeleteScalinggroup) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteScalinggroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeleteAutoScalingGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeleteAutoScalingGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAutoScalingGroup(input)
	cmd.logger.ExtraVerbosef("autoscaling.DeleteAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete scalinggroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteScalinggroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteScalinggroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalinggroup"), nil
}

func (cmd *DeleteScalinggroup) ParamsHelp() string {
	return generateParamsHelp("deletescalinggroup", structListParamsKeys(cmd))
}

func (cmd *DeleteScalinggroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteScalingpolicy(sess *session.Session, l ...*logger.Logger) *DeleteScalingpolicy {
	cmd := new(DeleteScalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	return cmd
}

func (cmd *DeleteScalingpolicy) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *DeleteScalingpolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeletePolicyInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeletePolicyInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeletePolicy(input)
	cmd.logger.ExtraVerbosef("autoscaling.DeletePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete scalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete scalingpolicy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete scalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteScalingpolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteScalingpolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalingpolicy"), nil
}

func (cmd *DeleteScalingpolicy) ParamsHelp() string {
	return generateParamsHelp("deletescalingpolicy", structListParamsKeys(cmd))
}

func (cmd *DeleteScalingpolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteSecuritygroup(sess *session.Session, l ...*logger.Logger) *DeleteSecuritygroup {
	cmd := new(DeleteSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteSecuritygroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSecurityGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSecurityGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSecurityGroup(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteSecurityGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete securitygroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSecuritygroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteSecuritygroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteSecurityGroupInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteSecurityGroupInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSecurityGroup(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteSecurityGroup call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete securitygroup ok")
			return fakeDryRunId("securitygroup"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteSecuritygroup) ParamsHelp() string {
	return generateParamsHelp("deletesecuritygroup", structListParamsKeys(cmd))
}

func (cmd *DeleteSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteSnapshot(sess *session.Session, l ...*logger.Logger) *DeleteSnapshot {
	cmd := new(DeleteSnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteSnapshot) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteSnapshot) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSnapshotInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSnapshotInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSnapshot(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteSnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete snapshot '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSnapshot) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteSnapshot) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteSnapshotInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteSnapshotInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSnapshot(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteSnapshot call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete snapshot ok")
			return fakeDryRunId("snapshot"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteSnapshot) ParamsHelp() string {
	return generateParamsHelp("deletesnapshot", structListParamsKeys(cmd))
}

func (cmd *DeleteSnapshot) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteStack(sess *session.Session, l ...*logger.Logger) *DeleteStack {
	cmd := new(DeleteStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudformation.New(sess)
	}
	return cmd
}

func (cmd *DeleteStack) SetApi(api cloudformationiface.CloudFormationAPI) {
	cmd.api = api
}

func (cmd *DeleteStack) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.DeleteStackInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.DeleteStackInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteStack(input)
	cmd.logger.ExtraVerbosef("cloudformation.DeleteStack call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete stack '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteStack) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteStack) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("stack"), nil
}

func (cmd *DeleteStack) ParamsHelp() string {
	return generateParamsHelp("deletestack", structListParamsKeys(cmd))
}

func (cmd *DeleteStack) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteSubnet(sess *session.Session, l ...*logger.Logger) *DeleteSubnet {
	cmd := new(DeleteSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteSubnet) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteSubnet) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSubnetInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSubnetInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSubnet(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteSubnet call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete subnet '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSubnet) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteSubnet) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteSubnetInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteSubnetInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSubnet(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteSubnet call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete subnet ok")
			return fakeDryRunId("subnet"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteSubnet) ParamsHelp() string {
	return generateParamsHelp("deletesubnet", structListParamsKeys(cmd))
}

func (cmd *DeleteSubnet) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteSubscription(sess *session.Session, l ...*logger.Logger) *DeleteSubscription {
	cmd := new(DeleteSubscription)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sns.New(sess)
	}
	return cmd
}

func (cmd *DeleteSubscription) SetApi(api snsiface.SNSAPI) {
	cmd.api = api
}

func (cmd *DeleteSubscription) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.UnsubscribeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.UnsubscribeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.Unsubscribe(input)
	cmd.logger.ExtraVerbosef("sns.Unsubscribe call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete subscription: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete subscription '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete subscription done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSubscription) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteSubscription) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("subscription"), nil
}

func (cmd *DeleteSubscription) ParamsHelp() string {
	return generateParamsHelp("deletesubscription", structListParamsKeys(cmd))
}

func (cmd *DeleteSubscription) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteTag(sess *session.Session, l ...*logger.Logger) *DeleteTag {
	cmd := new(DeleteTag)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteTag) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteTag) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete tag: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete tag '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete tag done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTag) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteTag) ParamsHelp() string {
	return generateParamsHelp("deletetag", structListParamsKeys(cmd))
}

func (cmd *DeleteTag) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteTargetgroup(sess *session.Session, l ...*logger.Logger) *DeleteTargetgroup {
	cmd := new(DeleteTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *DeleteTargetgroup) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *DeleteTargetgroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteTargetGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteTargetGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteTargetGroup(input)
	cmd.logger.ExtraVerbosef("elbv2.DeleteTargetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete targetgroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTargetgroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteTargetgroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("targetgroup"), nil
}

func (cmd *DeleteTargetgroup) ParamsHelp() string {
	return generateParamsHelp("deletetargetgroup", structListParamsKeys(cmd))
}

func (cmd *DeleteTargetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteTopic(sess *session.Session, l ...*logger.Logger) *DeleteTopic {
	cmd := new(DeleteTopic)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = sns.New(sess)
	}
	return cmd
}

func (cmd *DeleteTopic) SetApi(api snsiface.SNSAPI) {
	cmd.api = api
}

func (cmd *DeleteTopic) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.DeleteTopicInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.DeleteTopicInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteTopic(input)
	cmd.logger.ExtraVerbosef("sns.DeleteTopic call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete topic: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete topic '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete topic done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTopic) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteTopic) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("topic"), nil
}

func (cmd *DeleteTopic) ParamsHelp() string {
	return generateParamsHelp("deletetopic", structListParamsKeys(cmd))
}

func (cmd *DeleteTopic) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteUser(sess *session.Session, l ...*logger.Logger) *DeleteUser {
	cmd := new(DeleteUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DeleteUser) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DeleteUser) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteUserInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteUserInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteUser(input)
	cmd.logger.ExtraVerbosef("iam.DeleteUser call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete user '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteUser) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteUser) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("user"), nil
}

func (cmd *DeleteUser) ParamsHelp() string {
	return generateParamsHelp("deleteuser", structListParamsKeys(cmd))
}

func (cmd *DeleteUser) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteVolume(sess *session.Session, l ...*logger.Logger) *DeleteVolume {
	cmd := new(DeleteVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteVolume) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteVolumeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVolumeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVolume(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteVolume call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete volume '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteVolume) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteVolume) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteVolumeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteVolumeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteVolume call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete volume ok")
			return fakeDryRunId("volume"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteVolume) ParamsHelp() string {
	return generateParamsHelp("deletevolume", structListParamsKeys(cmd))
}

func (cmd *DeleteVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteVpc(sess *session.Session, l ...*logger.Logger) *DeleteVpc {
	cmd := new(DeleteVpc)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DeleteVpc) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DeleteVpc) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteVpcInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVpcInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVpc(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteVpc call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete vpc: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete vpc '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete vpc done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteVpc) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteVpc) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteVpcInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DeleteVpcInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DeleteVpc(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteVpc call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: delete vpc ok")
			return fakeDryRunId("vpc"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DeleteVpc) ParamsHelp() string {
	return generateParamsHelp("deletevpc", structListParamsKeys(cmd))
}

func (cmd *DeleteVpc) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDeleteZone(sess *session.Session, l ...*logger.Logger) *DeleteZone {
	cmd := new(DeleteZone)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	return cmd
}

func (cmd *DeleteZone) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *DeleteZone) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &route53.DeleteHostedZoneInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in route53.DeleteHostedZoneInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeleteHostedZone(input)
	cmd.logger.ExtraVerbosef("route53.DeleteHostedZone call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("delete zone: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("delete zone '%s' done", extracted)
	} else {
		cmd.logger.Verbose("delete zone done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteZone) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DeleteZone) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("zone"), nil
}

func (cmd *DeleteZone) ParamsHelp() string {
	return generateParamsHelp("deletezone", structListParamsKeys(cmd))
}

func (cmd *DeleteZone) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachAlarm(sess *session.Session, l ...*logger.Logger) *DetachAlarm {
	cmd := new(DetachAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	return cmd
}

func (cmd *DetachAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *DetachAlarm) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach alarm '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachAlarm) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachAlarm) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *DetachAlarm) ParamsHelp() string {
	return generateParamsHelp("detachalarm", structListParamsKeys(cmd))
}

func (cmd *DetachAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachContainertask(sess *session.Session, l ...*logger.Logger) *DetachContainertask {
	cmd := new(DetachContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	return cmd
}

func (cmd *DetachContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *DetachContainertask) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach containertask '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachContainertask) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachContainertask) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *DetachContainertask) ParamsHelp() string {
	return generateParamsHelp("detachcontainertask", structListParamsKeys(cmd))
}

func (cmd *DetachContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachElasticip(sess *session.Session, l ...*logger.Logger) *DetachElasticip {
	cmd := new(DetachElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DetachElasticip) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachElasticip) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DisassociateAddressInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateAddressInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DisassociateAddress(input)
	cmd.logger.ExtraVerbosef("ec2.DisassociateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach elasticip '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachElasticip) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachElasticip) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DisassociateAddressInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DisassociateAddressInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DisassociateAddress(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DisassociateAddress call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: detach elasticip ok")
			return fakeDryRunId("elasticip"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DetachElasticip) ParamsHelp() string {
	return generateParamsHelp("detachelasticip", structListParamsKeys(cmd))
}

func (cmd *DetachElasticip) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachInstance(sess *session.Session, l ...*logger.Logger) *DetachInstance {
	cmd := new(DetachInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *DetachInstance) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *DetachInstance) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeregisterTargetsInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeregisterTargetsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeregisterTargets(input)
	cmd.logger.ExtraVerbosef("elbv2.DeregisterTargets call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach instance '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInstance) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachInstance) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instance"), nil
}

func (cmd *DetachInstance) ParamsHelp() string {
	return generateParamsHelp("detachinstance", structListParamsKeys(cmd))
}

func (cmd *DetachInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachInstanceprofile(sess *session.Session, l ...*logger.Logger) *DetachInstanceprofile {
	cmd := new(DetachInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DetachInstanceprofile) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachInstanceprofile) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach instanceprofile '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInstanceprofile) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachInstanceprofile) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("instanceprofile"), nil
}

func (cmd *DetachInstanceprofile) ParamsHelp() string {
	return generateParamsHelp("detachinstanceprofile", structListParamsKeys(cmd))
}

func (cmd *DetachInstanceprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachInternetgateway(sess *session.Session, l ...*logger.Logger) *DetachInternetgateway {
	cmd := new(DetachInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DetachInternetgateway) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachInternetgateway) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DetachInternetGatewayInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachInternetGatewayInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DetachInternetGateway(input)
	cmd.logger.ExtraVerbosef("ec2.DetachInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach internetgateway '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInternetgateway) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachInternetgateway) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DetachInternetGatewayInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DetachInternetGatewayInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DetachInternetGateway(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DetachInternetGateway call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: detach internetgateway ok")
			return fakeDryRunId("internetgateway"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DetachInternetgateway) ParamsHelp() string {
	return generateParamsHelp("detachinternetgateway", structListParamsKeys(cmd))
}

func (cmd *DetachInternetgateway) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachMfadevice(sess *session.Session, l ...*logger.Logger) *DetachMfadevice {
	cmd := new(DetachMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DetachMfadevice) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DetachMfadevice) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeactivateMFADeviceInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeactivateMFADeviceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DeactivateMFADevice(input)
	cmd.logger.ExtraVerbosef("iam.DeactivateMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach mfadevice '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachMfadevice) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachMfadevice) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("mfadevice"), nil
}

func (cmd *DetachMfadevice) ParamsHelp() string {
	return generateParamsHelp("detachmfadevice", structListParamsKeys(cmd))
}

func (cmd *DetachMfadevice) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachNetworkinterface(sess *session.Session, l ...*logger.Logger) *DetachNetworkinterface {
	cmd := new(DetachNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DetachNetworkinterface) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachNetworkinterface) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach networkinterface '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachNetworkinterface) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachNetworkinterface) ParamsHelp() string {
	return generateParamsHelp("detachnetworkinterface", structListParamsKeys(cmd))
}

func (cmd *DetachNetworkinterface) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachPolicy(sess *session.Session, l ...*logger.Logger) *DetachPolicy {
	cmd := new(DetachPolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DetachPolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DetachPolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach policy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachPolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachPolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *DetachPolicy) ParamsHelp() string {
	return generateParamsHelp("detachpolicy", structListParamsKeys(cmd))
}

func (cmd *DetachPolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachRole(sess *session.Session, l ...*logger.Logger) *DetachRole {
	cmd := new(DetachRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DetachRole) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DetachRole) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.RemoveRoleFromInstanceProfileInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.RemoveRoleFromInstanceProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RemoveRoleFromInstanceProfile(input)
	cmd.logger.ExtraVerbosef("iam.RemoveRoleFromInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach role '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachRole) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachRole) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("role"), nil
}

func (cmd *DetachRole) ParamsHelp() string {
	return generateParamsHelp("detachrole", structListParamsKeys(cmd))
}

func (cmd *DetachRole) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachRoutetable(sess *session.Session, l ...*logger.Logger) *DetachRoutetable {
	cmd := new(DetachRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DetachRoutetable) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachRoutetable) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DisassociateRouteTableInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateRouteTableInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DisassociateRouteTable(input)
	cmd.logger.ExtraVerbosef("ec2.DisassociateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach routetable '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachRoutetable) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachRoutetable) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DisassociateRouteTableInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DisassociateRouteTableInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DisassociateRouteTable(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DisassociateRouteTable call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: detach routetable ok")
			return fakeDryRunId("routetable"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DetachRoutetable) ParamsHelp() string {
	return generateParamsHelp("detachroutetable", structListParamsKeys(cmd))
}

func (cmd *DetachRoutetable) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachSecuritygroup(sess *session.Session, l ...*logger.Logger) *DetachSecuritygroup {
	cmd := new(DetachSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DetachSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachSecuritygroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach securitygroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachSecuritygroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachSecuritygroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("securitygroup"), nil
}

func (cmd *DetachSecuritygroup) ParamsHelp() string {
	return generateParamsHelp("detachsecuritygroup", structListParamsKeys(cmd))
}

func (cmd *DetachSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachUser(sess *session.Session, l ...*logger.Logger) *DetachUser {
	cmd := new(DetachUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *DetachUser) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *DetachUser) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.RemoveUserFromGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.RemoveUserFromGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.RemoveUserFromGroup(input)
	cmd.logger.ExtraVerbosef("iam.RemoveUserFromGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach user '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachUser) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachUser) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("user"), nil
}

func (cmd *DetachUser) ParamsHelp() string {
	return generateParamsHelp("detachuser", structListParamsKeys(cmd))
}

func (cmd *DetachUser) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewDetachVolume(sess *session.Session, l ...*logger.Logger) *DetachVolume {
	cmd := new(DetachVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *DetachVolume) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *DetachVolume) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DetachVolumeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachVolumeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DetachVolume(input)
	cmd.logger.ExtraVerbosef("ec2.DetachVolume call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("detach volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("detach volume '%s' done", extracted)
	} else {
		cmd.logger.Verbose("detach volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachVolume) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *DetachVolume) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.DetachVolumeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.DetachVolumeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.DetachVolume(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.DetachVolume call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: detach volume ok")
			return fakeDryRunId("volume"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *DetachVolume) ParamsHelp() string {
	return generateParamsHelp("detachvolume", structListParamsKeys(cmd))
}

func (cmd *DetachVolume) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewImportImage(sess *session.Session, l ...*logger.Logger) *ImportImage {
	cmd := new(ImportImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *ImportImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *ImportImage) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ImportImageInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ImportImageInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ImportImage(input)
	cmd.logger.ExtraVerbosef("ec2.ImportImage call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("import image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("import image '%s' done", extracted)
	} else {
		cmd.logger.Verbose("import image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *ImportImage) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *ImportImage) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.ImportImageInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.ImportImageInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.ImportImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.ImportImage call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: import image ok")
			return fakeDryRunId("image"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *ImportImage) ParamsHelp() string {
	return generateParamsHelp("importimage", structListParamsKeys(cmd))
}

func (cmd *ImportImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStartAlarm(sess *session.Session, l ...*logger.Logger) *StartAlarm {
	cmd := new(StartAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	return cmd
}

func (cmd *StartAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *StartAlarm) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.EnableAlarmActionsInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.EnableAlarmActionsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.EnableAlarmActions(input)
	cmd.logger.ExtraVerbosef("cloudwatch.EnableAlarmActions call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("start alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("start alarm '%s' done", extracted)
	} else {
		cmd.logger.Verbose("start alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartAlarm) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *StartAlarm) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *StartAlarm) ParamsHelp() string {
	return generateParamsHelp("startalarm", structListParamsKeys(cmd))
}

func (cmd *StartAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStartContainertask(sess *session.Session, l ...*logger.Logger) *StartContainertask {
	cmd := new(StartContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	return cmd
}

func (cmd *StartContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *StartContainertask) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("start containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("start containertask '%s' done", extracted)
	} else {
		cmd.logger.Verbose("start containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartContainertask) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *StartContainertask) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *StartContainertask) ParamsHelp() string {
	return generateParamsHelp("startcontainertask", structListParamsKeys(cmd))
}

func (cmd *StartContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStartInstance(sess *session.Session, l ...*logger.Logger) *StartInstance {
	cmd := new(StartInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *StartInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *StartInstance) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.StartInstancesInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StartInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.StartInstances(input)
	cmd.logger.ExtraVerbosef("ec2.StartInstances call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("start instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("start instance '%s' done", extracted)
	} else {
		cmd.logger.Verbose("start instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartInstance) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *StartInstance) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.StartInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.StartInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.StartInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.StartInstances call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: start instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *StartInstance) ParamsHelp() string {
	return generateParamsHelp("startinstance", structListParamsKeys(cmd))
}

func (cmd *StartInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStopAlarm(sess *session.Session, l ...*logger.Logger) *StopAlarm {
	cmd := new(StopAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudwatch.New(sess)
	}
	return cmd
}

func (cmd *StopAlarm) SetApi(api cloudwatchiface.CloudWatchAPI) {
	cmd.api = api
}

func (cmd *StopAlarm) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.DisableAlarmActionsInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.DisableAlarmActionsInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.DisableAlarmActions(input)
	cmd.logger.ExtraVerbosef("cloudwatch.DisableAlarmActions call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("stop alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("stop alarm '%s' done", extracted)
	} else {
		cmd.logger.Verbose("stop alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopAlarm) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *StopAlarm) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("alarm"), nil
}

func (cmd *StopAlarm) ParamsHelp() string {
	return generateParamsHelp("stopalarm", structListParamsKeys(cmd))
}

func (cmd *StopAlarm) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStopContainertask(sess *session.Session, l ...*logger.Logger) *StopContainertask {
	cmd := new(StopContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	return cmd
}

func (cmd *StopContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *StopContainertask) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("stop containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("stop containertask '%s' done", extracted)
	} else {
		cmd.logger.Verbose("stop containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopContainertask) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *StopContainertask) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *StopContainertask) ParamsHelp() string {
	return generateParamsHelp("stopcontainertask", structListParamsKeys(cmd))
}

func (cmd *StopContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewStopInstance(sess *session.Session, l ...*logger.Logger) *StopInstance {
	cmd := new(StopInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *StopInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *StopInstance) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.StopInstancesInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StopInstancesInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.StopInstances(input)
	cmd.logger.ExtraVerbosef("ec2.StopInstances call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("stop instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("stop instance '%s' done", extracted)
	} else {
		cmd.logger.Verbose("stop instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopInstance) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *StopInstance) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.StopInstancesInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.StopInstancesInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.StopInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.StopInstances call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: stop instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *StopInstance) ParamsHelp() string {
	return generateParamsHelp("stopinstance", structListParamsKeys(cmd))
}

func (cmd *StopInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateBucket(sess *session.Session, l ...*logger.Logger) *UpdateBucket {
	cmd := new(UpdateBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	return cmd
}

func (cmd *UpdateBucket) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *UpdateBucket) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update bucket '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateBucket) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateBucket) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("bucket"), nil
}

func (cmd *UpdateBucket) ParamsHelp() string {
	return generateParamsHelp("updatebucket", structListParamsKeys(cmd))
}

func (cmd *UpdateBucket) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateContainertask(sess *session.Session, l ...*logger.Logger) *UpdateContainertask {
	cmd := new(UpdateContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ecs.New(sess)
	}
	return cmd
}

func (cmd *UpdateContainertask) SetApi(api ecsiface.ECSAPI) {
	cmd.api = api
}

func (cmd *UpdateContainertask) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.UpdateServiceInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.UpdateServiceInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.UpdateService(input)
	cmd.logger.ExtraVerbosef("ecs.UpdateService call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update containertask '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateContainertask) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateContainertask) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("containertask"), nil
}

func (cmd *UpdateContainertask) ParamsHelp() string {
	return generateParamsHelp("updatecontainertask", structListParamsKeys(cmd))
}

func (cmd *UpdateContainertask) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateDistribution(sess *session.Session, l ...*logger.Logger) *UpdateDistribution {
	cmd := new(UpdateDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudfront.New(sess)
	}
	return cmd
}

func (cmd *UpdateDistribution) SetApi(api cloudfrontiface.CloudFrontAPI) {
	cmd.api = api
}

func (cmd *UpdateDistribution) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update distribution '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateDistribution) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateDistribution) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("distribution"), nil
}

func (cmd *UpdateDistribution) ParamsHelp() string {
	return generateParamsHelp("updatedistribution", structListParamsKeys(cmd))
}

func (cmd *UpdateDistribution) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateImage(sess *session.Session, l ...*logger.Logger) *UpdateImage {
	cmd := new(UpdateImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *UpdateImage) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *UpdateImage) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update image '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateImage) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateImage) ParamsHelp() string {
	return generateParamsHelp("updateimage", structListParamsKeys(cmd))
}

func (cmd *UpdateImage) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateInstance(sess *session.Session, l ...*logger.Logger) *UpdateInstance {
	cmd := new(UpdateInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *UpdateInstance) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *UpdateInstance) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ModifyInstanceAttributeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifyInstanceAttributeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ModifyInstanceAttribute(input)
	cmd.logger.ExtraVerbosef("ec2.ModifyInstanceAttribute call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update instance '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateInstance) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateInstance) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.ModifyInstanceAttributeInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.ModifyInstanceAttributeInput: %s", err)
	}

	start := time.Now()
	_, err := cmd.api.ModifyInstanceAttribute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.ModifyInstanceAttribute call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: update instance ok")
			return fakeDryRunId("instance"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *UpdateInstance) ParamsHelp() string {
	return generateParamsHelp("updateinstance", structListParamsKeys(cmd))
}

func (cmd *UpdateInstance) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateLoginprofile(sess *session.Session, l ...*logger.Logger) *UpdateLoginprofile {
	cmd := new(UpdateLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *UpdateLoginprofile) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *UpdateLoginprofile) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.UpdateLoginProfileInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.UpdateLoginProfileInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.UpdateLoginProfile(input)
	cmd.logger.ExtraVerbosef("iam.UpdateLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update loginprofile '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateLoginprofile) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateLoginprofile) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("loginprofile"), nil
}

func (cmd *UpdateLoginprofile) ParamsHelp() string {
	return generateParamsHelp("updateloginprofile", structListParamsKeys(cmd))
}

func (cmd *UpdateLoginprofile) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdatePolicy(sess *session.Session, l ...*logger.Logger) *UpdatePolicy {
	cmd := new(UpdatePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = iam.New(sess)
	}
	return cmd
}

func (cmd *UpdatePolicy) SetApi(api iamiface.IAMAPI) {
	cmd.api = api
}

func (cmd *UpdatePolicy) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreatePolicyVersionInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreatePolicyVersionInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.CreatePolicyVersion(input)
	cmd.logger.ExtraVerbosef("iam.CreatePolicyVersion call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update policy '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdatePolicy) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdatePolicy) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("policy"), nil
}

func (cmd *UpdatePolicy) ParamsHelp() string {
	return generateParamsHelp("updatepolicy", structListParamsKeys(cmd))
}

func (cmd *UpdatePolicy) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateRecord(sess *session.Session, l ...*logger.Logger) *UpdateRecord {
	cmd := new(UpdateRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = route53.New(sess)
	}
	return cmd
}

func (cmd *UpdateRecord) SetApi(api route53iface.Route53API) {
	cmd.api = api
}

func (cmd *UpdateRecord) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update record '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateRecord) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateRecord) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("record"), nil
}

func (cmd *UpdateRecord) ParamsHelp() string {
	return generateParamsHelp("updaterecord", structListParamsKeys(cmd))
}

func (cmd *UpdateRecord) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateS3object(sess *session.Session, l ...*logger.Logger) *UpdateS3object {
	cmd := new(UpdateS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = s3.New(sess)
	}
	return cmd
}

func (cmd *UpdateS3object) SetApi(api s3iface.S3API) {
	cmd.api = api
}

func (cmd *UpdateS3object) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.PutObjectAclInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.PutObjectAclInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.PutObjectAcl(input)
	cmd.logger.ExtraVerbosef("s3.PutObjectAcl call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update s3object '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateS3object) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateS3object) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("s3object"), nil
}

func (cmd *UpdateS3object) ParamsHelp() string {
	return generateParamsHelp("updates3object", structListParamsKeys(cmd))
}

func (cmd *UpdateS3object) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateScalinggroup(sess *session.Session, l ...*logger.Logger) *UpdateScalinggroup {
	cmd := new(UpdateScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = autoscaling.New(sess)
	}
	return cmd
}

func (cmd *UpdateScalinggroup) SetApi(api autoscalingiface.AutoScalingAPI) {
	cmd.api = api
}

func (cmd *UpdateScalinggroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.UpdateAutoScalingGroupInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.UpdateAutoScalingGroup(input)
	cmd.logger.ExtraVerbosef("autoscaling.UpdateAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update scalinggroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateScalinggroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateScalinggroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("scalinggroup"), nil
}

func (cmd *UpdateScalinggroup) ParamsHelp() string {
	return generateParamsHelp("updatescalinggroup", structListParamsKeys(cmd))
}

func (cmd *UpdateScalinggroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateSecuritygroup(sess *session.Session, l ...*logger.Logger) *UpdateSecuritygroup {
	cmd := new(UpdateSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *UpdateSecuritygroup) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *UpdateSecuritygroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update securitygroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateSecuritygroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateSecuritygroup) ParamsHelp() string {
	return generateParamsHelp("updatesecuritygroup", structListParamsKeys(cmd))
}

func (cmd *UpdateSecuritygroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateStack(sess *session.Session, l ...*logger.Logger) *UpdateStack {
	cmd := new(UpdateStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = cloudformation.New(sess)
	}
	return cmd
}

func (cmd *UpdateStack) SetApi(api cloudformationiface.CloudFormationAPI) {
	cmd.api = api
}

func (cmd *UpdateStack) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.UpdateStackInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.UpdateStackInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.UpdateStack(input)
	cmd.logger.ExtraVerbosef("cloudformation.UpdateStack call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update stack '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateStack) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateStack) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("stack"), nil
}

func (cmd *UpdateStack) ParamsHelp() string {
	return generateParamsHelp("updatestack", structListParamsKeys(cmd))
}

func (cmd *UpdateStack) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateSubnet(sess *session.Session, l ...*logger.Logger) *UpdateSubnet {
	cmd := new(UpdateSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = ec2.New(sess)
	}
	return cmd
}

func (cmd *UpdateSubnet) SetApi(api ec2iface.EC2API) {
	cmd.api = api
}

func (cmd *UpdateSubnet) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ModifySubnetAttributeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifySubnetAttributeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ModifySubnetAttribute(input)
	cmd.logger.ExtraVerbosef("ec2.ModifySubnetAttribute call took %s", time.Since(start))
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update subnet '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateSubnet) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateSubnet) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("subnet"), nil
}

func (cmd *UpdateSubnet) ParamsHelp() string {
	return generateParamsHelp("updatesubnet", structListParamsKeys(cmd))
}

func (cmd *UpdateSubnet) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}

func NewUpdateTargetgroup(sess *session.Session, l ...*logger.Logger) *UpdateTargetgroup {
	cmd := new(UpdateTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if sess != nil {
		cmd.api = elbv2.New(sess)
	}
	return cmd
}

func (cmd *UpdateTargetgroup) SetApi(api elbv2iface.ELBV2API) {
	cmd.api = api
}

func (cmd *UpdateTargetgroup) Run(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(ctx); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(ctx)
	if err != nil {
		return nil, err
	}

	var extracted interface{}
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			cmd.logger.Warning("update targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		cmd.logger.Verbosef("update targetgroup '%s' done", extracted)
	} else {
		cmd.logger.Verbose("update targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(ctx, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateTargetgroup) ValidateCommand(params map[string]interface{}, refs []string) (errs []error) {
	if err := cmd.inject(params); err != nil {
		return []error{err}
	}
	if err := validateStruct(cmd, refs); err != nil {
		errs = append(errs, err)
	}

	return
}

func (cmd *UpdateTargetgroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	return fakeDryRunId("targetgroup"), nil
}

func (cmd *UpdateTargetgroup) ParamsHelp() string {
	return generateParamsHelp("updatetargetgroup", structListParamsKeys(cmd))
}

func (cmd *UpdateTargetgroup) inject(params map[string]interface{}) error {
	return structSetter(cmd, params)
}
