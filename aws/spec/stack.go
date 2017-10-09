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
	"errors"
	"os"

	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/wallix/awless/logger"
)

type CreateStack struct {
	_               string `action:"create" entity:"stack" awsAPI:"cloudformation" awsCall:"CreateStack" awsInput:"cloudformation.CreateStackInput" awsOutput:"cloudformation.CreateStackOutput"`
	logger          *logger.Logger
	api             cloudformationiface.CloudFormationAPI
	Name            *string   `awsName:"StackName" awsType:"awsstr" templateName:"name" required:""`
	TemplateFile    *string   `awsName:"TemplateBody" awsType:"awsfiletostring" templateName:"template-file" required:""`
	Capabilities    []*string `awsName:"Capabilities" awsType:"awsstringslice" templateName:"capabilities"`
	DisableRollback *bool     `awsName:"DisableRollback" awsType:"awsbool" templateName:"disable-rollback"`
	Notifications   []*string `awsName:"NotificationARNs" awsType:"awsstringslice" templateName:"notifications"`
	OnFailure       *string   `awsName:"OnFailure" awsType:"awsstr" templateName:"on-failure"`
	Parameters      []*string `awsName:"Parameters" awsType:"awsparameterslice" templateName:"parameters"`
	ResourceTypes   []*string `awsName:"ResourceTypes" awsType:"awsstringslice" templateName:"resource-types"`
	Role            *string   `awsName:"RoleARN" awsType:"awsstr" templateName:"role"`
	PolicyFile      *string   `awsName:"StackPolicyBody" awsType:"awsfiletostring" templateName:"policy-file"`
	Timeout         *int64    `awsName:"TimeoutInMinutes" awsType:"awsint64" templateName:"timeout"`
}

func (cmd *CreateStack) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateStack) Validate_TemplateFile() error {
	if _, err := os.Stat(StringValue(cmd.TemplateFile)); err != nil {
		return errors.New(strings.TrimLeft(err.Error(), "stat "))
	}
	return nil
}

func (cmd *CreateStack) ExtractResult(i interface{}) string {
	return StringValue(i.(*cloudformation.CreateStackOutput).StackId)
}

type UpdateStack struct {
	_                   string `action:"update" entity:"stack" awsAPI:"cloudformation" awsCall:"UpdateStack" awsInput:"cloudformation.UpdateStackInput" awsOutput:"cloudformation.UpdateStackOutput"`
	logger              *logger.Logger
	api                 cloudformationiface.CloudFormationAPI
	Name                *string   `awsName:"StackName" awsType:"awsstr" templateName:"name" required:""`
	Capabilities        []*string `awsName:"Capabilities" awsType:"awsstringslice" templateName:"capabilities"`
	Notifications       []*string `awsName:"NotificationARNs" awsType:"awsstringslice" templateName:"notifications"`
	Parameters          []*string `awsName:"Parameters" awsType:"awsparameterslice" templateName:"parameters"`
	ResourceTypes       []*string `awsName:"ResourceTypes" awsType:"awsstringslice" templateName:"resource-types"`
	Role                *string   `awsName:"RoleARN" awsType:"awsstr" templateName:"role"`
	PolicyFile          *string   `awsName:"StackPolicyBody" awsType:"awsfiletostring" templateName:"policy-file"`
	PolicyUpdateFile    *string   `awsName:"StackPolicyDuringUpdateBody" awsType:"awsfiletostring" templateName:"policy-update-file"`
	TemplateFile        *string   `awsName:"TemplateBody" awsType:"awsfiletostring" templateName:"template-file"`
	UsePreviousTemplate *bool     `awsName:"UsePreviousTemplate" awsType:"awsbool" templateName:"use-previous-template"`
}

func (cmd *UpdateStack) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *UpdateStack) ExtractResult(i interface{}) string {
	return StringValue(i.(*cloudformation.UpdateStackOutput).StackId)
}

type DeleteStack struct {
	_               string `action:"delete" entity:"stack" awsAPI:"cloudformation" awsCall:"DeleteStack" awsInput:"cloudformation.DeleteStackInput" awsOutput:"cloudformation.DeleteStackOutput"`
	logger          *logger.Logger
	api             cloudformationiface.CloudFormationAPI
	Name            *string   `awsName:"StackName" awsType:"awsstr" templateName:"name" required:""`
	RetainResources []*string `awsName:"RetainResources" awsType:"awsstringslice" templateName:"retain-resources"`
}

func (cmd *DeleteStack) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}
