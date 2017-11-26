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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

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
	Tags                []*string `awsName:"Tags" awsType:"awstagslice" templateName:"tags"`
	PolicyBody          *string   `awsName:"StackPolicyBody" awsType:"awsstr"`
	StackFile           *string   `templateName:"stack-file"`
}

func (cmd *UpdateStack) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *UpdateStack) ExtractResult(i interface{}) string {
	return StringValue(i.(*cloudformation.UpdateStackOutput).StackId)
}

func (cmd *UpdateStack) BeforeRun(ctx map[string]interface{}) error {
	if cmd.StackFile == nil {
		return nil
	}

	type stackFile struct {
		Parameters  map[string]string
		Tags        map[string]string
		StackPolicy map[string]interface{}
	}

	data := &stackFile{}

	file, err := ioutil.ReadFile(*cmd.StackFile)
	if err != nil {
		return err
	}

	switch path.Ext(*cmd.StackFile) {
	case ".json":
		err = json.Unmarshal(file, &data)
		cmd.logger.Infof("Reading JSON file %s", *cmd.StackFile)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(file, &data)
		cmd.logger.Infof("Reading YML file %s", *cmd.StackFile)
	default:
		return fmt.Errorf("Unknown format %s", path.Ext(*cmd.StackFile))
	}

	if err != nil {
		return err
	}

	cmd.Parameters = mergeCliAndFileValues(data.Parameters, cmd.Parameters)
	cmd.Tags = mergeCliAndFileValues(data.Tags, cmd.Tags)

	// use PolicyBody only when PolicyFile isn't specified
	if cmd.PolicyFile == nil {
		policyBytes, err := json.Marshal(data.StackPolicy)
		if err != nil {
			return err
		}

		policyStr := string(policyBytes)
		cmd.PolicyBody = &policyStr
	}

	return nil
}

// mergeCliAndFileValues is the helper func used to merge tags or parameters
// supplied with CLI and StackFile with higher priority for values passed via CLI
func mergeCliAndFileValues(valMap map[string]string, valSlice []*string) (resSlice []*string) {
	val := make(map[string]string)

	// building map of parameters passed from cli
	for _, v := range valSlice {
		splits := strings.SplitN(*v, ":", 2)
		if len(splits) == 2 {
			val[splits[0]] = splits[1]
		}
	}

	// adding/overwritting values from cli
	// to the files values map
	for k, v := range val {
		valMap[k] = v
	}

	// building final parameters list in the expected
	// "awsparameterslice" format
	for k, v := range valMap {
		p := strings.Join([]string{k, v}, ":")
		resSlice = append(resSlice, &p)
	}

	return resSlice
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
