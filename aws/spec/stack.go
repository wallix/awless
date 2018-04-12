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
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/wallix/awless/logger"
	"gopkg.in/yaml.v2"
)

type CreateStack struct {
	_                     string `action:"create" entity:"stack" awsAPI:"cloudformation" awsCall:"CreateStack" awsInput:"cloudformation.CreateStackInput" awsOutput:"cloudformation.CreateStackOutput"`
	logger                *logger.Logger
	graph                 cloud.GraphAPI
	api                   cloudformationiface.CloudFormationAPI
	Name                  *string   `awsName:"StackName" awsType:"awsstr" templateName:"name"`
	TemplateFile          *string   `awsName:"TemplateBody" awsType:"awsfiletostring" templateName:"template-file"`
	Capabilities          []*string `awsName:"Capabilities" awsType:"awsstringslice" templateName:"capabilities"`
	DisableRollback       *bool     `awsName:"DisableRollback" awsType:"awsbool" templateName:"disable-rollback"`
	Notifications         []*string `awsName:"NotificationARNs" awsType:"awsstringslice" templateName:"notifications"`
	OnFailure             *string   `awsName:"OnFailure" awsType:"awsstr" templateName:"on-failure"`
	Parameters            []*string `awsName:"Parameters" awsType:"awsparameterslice" templateName:"parameters"`
	ResourceTypes         []*string `awsName:"ResourceTypes" awsType:"awsstringslice" templateName:"resource-types"`
	Role                  *string   `awsName:"RoleARN" awsType:"awsstr" templateName:"role"`
	PolicyFile            *string   `awsName:"StackPolicyBody" awsType:"awsfiletostring" templateName:"policy-file"`
	Timeout               *int64    `awsName:"TimeoutInMinutes" awsType:"awsint64" templateName:"timeout"`
	Tags                  []*string `awsName:"Tags" awsType:"awstagslice" templateName:"tags"`
	PolicyBody            *string   `awsName:"StackPolicyBody" awsType:"awsstr"`
	StackFile             *string   `templateName:"stack-file"`
	RollbackTriggers      []*string `awsName:"RollbackConfiguration.RollbackTriggers" awsType:"awsalarmrollbacktriggers" templateName:"rollback-triggers"`
	RollbackMonitoringMin *int64    `awsName:"RollbackConfiguration.MonitoringTimeInMinutes" awsType:"awsint64" templateName:"rollback-monitoring-min"`
}

func (cmd *CreateStack) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("name"), params.Key("template-file"), params.Opt("capabilities", "disable-rollback", "notifications", "on-failure", "parameters", "policy-file", "resource-types", "role", "stack-file", "tags", "timeout", "rollback-triggers", "rollback-monitoring-min")),
		params.Validators{"template-file": params.IsFilepath},
	)
}

func (cmd *CreateStack) ExtractResult(i interface{}) string {
	return StringValue(i.(*cloudformation.CreateStackOutput).StackId)
}

// Add StackFile support via BeforeRun hook
// https://github.com/wallix/awless/issues/145
// http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/continuous-delivery-codepipeline-cfn-artifacts.html
func (cmd *CreateStack) BeforeRun(renv env.Running) (err error) {
	cmd.Parameters, cmd.Tags, cmd.PolicyBody, err = processStackFile(cmd.StackFile, cmd.PolicyFile, cmd.Parameters, cmd.Tags)
	return
}

type UpdateStack struct {
	_                     string `action:"update" entity:"stack" awsAPI:"cloudformation" awsCall:"UpdateStack" awsInput:"cloudformation.UpdateStackInput" awsOutput:"cloudformation.UpdateStackOutput"`
	logger                *logger.Logger
	graph                 cloud.GraphAPI
	api                   cloudformationiface.CloudFormationAPI
	Name                  *string   `awsName:"StackName" awsType:"awsstr" templateName:"name"`
	Capabilities          []*string `awsName:"Capabilities" awsType:"awsstringslice" templateName:"capabilities"`
	Notifications         []*string `awsName:"NotificationARNs" awsType:"awsstringslice" templateName:"notifications"`
	Parameters            []*string `awsName:"Parameters" awsType:"awsparameterslice" templateName:"parameters"`
	ResourceTypes         []*string `awsName:"ResourceTypes" awsType:"awsstringslice" templateName:"resource-types"`
	Role                  *string   `awsName:"RoleARN" awsType:"awsstr" templateName:"role"`
	PolicyFile            *string   `awsName:"StackPolicyBody" awsType:"awsfiletostring" templateName:"policy-file"`
	PolicyUpdateFile      *string   `awsName:"StackPolicyDuringUpdateBody" awsType:"awsfiletostring" templateName:"policy-update-file"`
	TemplateFile          *string   `awsName:"TemplateBody" awsType:"awsfiletostring" templateName:"template-file"`
	UsePreviousTemplate   *bool     `awsName:"UsePreviousTemplate" awsType:"awsbool" templateName:"use-previous-template"`
	Tags                  []*string `awsName:"Tags" awsType:"awstagslice" templateName:"tags"`
	PolicyBody            *string   `awsName:"StackPolicyBody" awsType:"awsstr"`
	StackFile             *string   `templateName:"stack-file"`
	RollbackTriggers      []*string `awsName:"RollbackConfiguration.RollbackTriggers" awsType:"awsalarmrollbacktriggers" templateName:"rollback-triggers"`
	RollbackMonitoringMin *int64    `awsName:"RollbackConfiguration.MonitoringTimeInMinutes" awsType:"awsint64" templateName:"rollback-monitoring-min"`
}

func (cmd *UpdateStack) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"),
		params.Opt("capabilities", "notifications", "parameters", "policy-file", "policy-update-file", "resource-types", "role", "stack-file", "tags", "template-file", "use-previous-template", "rollback-triggers", "rollback-monitoring-min"),
	))
}

func (cmd *UpdateStack) ExtractResult(i interface{}) string {
	return StringValue(i.(*cloudformation.UpdateStackOutput).StackId)
}

// Add StackFile support via BeforeRun hook
// https://github.com/wallix/awless/issues/145
// http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/continuous-delivery-codepipeline-cfn-artifacts.html
func (cmd *UpdateStack) BeforeRun(renv env.Running) (err error) {
	cmd.Parameters, cmd.Tags, cmd.PolicyBody, err = processStackFile(cmd.StackFile, cmd.PolicyFile, cmd.Parameters, cmd.Tags)
	return
}

type stackFile struct {
	Parameters  map[string]string      `yaml:"Parameters"`
	Tags        map[string]string      `yaml:"Tags"`
	StackPolicy map[string]interface{} `yaml:"StackPolicy"`
}

func processStackFile(stackFilePath, policyFile *string, parameters, tags []*string) (newParams, newTags []*string, policyData *string, err error) {
	if stackFilePath == nil {
		return parameters, tags, nil, nil
	}

	data, err := readStackFile(*stackFilePath)
	if err != nil {
		return nil, nil, nil, err
	}

	if data == nil {
		return parameters, tags, nil, nil
	}

	newParams = mergeCliAndFileValues(data.Parameters, parameters)
	newTags = mergeCliAndFileValues(data.Tags, tags)

	if policyFile == nil && data.StackPolicy != nil {
		policyBytes, err := json.Marshal(data.StackPolicy)
		if err != nil {
			return nil, nil, nil, err
		}

		policyData = String(string(policyBytes))
	}

	return newParams, newTags, policyData, nil
}

func readStackFile(p string) (sf *stackFile, err error) {
	var file []byte
	file, err = ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	switch path.Ext(p) {
	case ".json":
		err = json.Unmarshal(file, &sf)
		if err != nil {
			// Result error message:
			// [info]    KO update stack
			// before run:
			// json: unmarshal errors:
			//   invalid character '}' looking for beginning of object key string
			return nil, fmt.Errorf("\njson: unmarshal errors:\n  %s", err)
		}
	case ".yml", ".yaml":
		err = yaml.Unmarshal(file, &sf)
		if err != nil {
			// Result error message:
			// [info]    KO update stack
			// before run:
			// yaml: unmarshal errors:
			//   line 1: cannot unmarshal !!str `lalla` into awsspec.stackFile
			return nil, fmt.Errorf("\n%s", err)
		}
	default:
		return nil, fmt.Errorf("Unknown StackFile format %q. Should be \".json\", \".yml\" or \".yaml\"", path.Ext(p))
	}
	return sf, err
}

// mergeCliAndFileValues is the helper func used to merge tags or parameters
// supplied with CLI and StackFile with higher priority for values passed via CLI
// example:
// via cli passed next parameters:
//   Test1=a
//   Test2=b
// via StackFile passed next parameters:
//   Test2=x
//   Test3=y
// after merge result will be:
//   Test1=a
//   Test2=b
//   Test3=y
func mergeCliAndFileValues(valMap map[string]string, valSlice []*string) (resSlice []*string) {
	// if values map are absent in StackFile
	// just return slice of CLI values
	if valMap == nil {
		return valSlice
	}

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

	mapKeys := make([]string, 0, len(valMap))
	for k := range valMap {
		mapKeys = append(mapKeys, k)
	}

	// soring map keys, so we have predictable values order for tests
	sort.Strings(mapKeys)

	// building final parameters list in the expected
	// "awsparameterslice" format
	for _, k := range mapKeys {
		p := strings.Join([]string{k, valMap[k]}, ":")
		resSlice = append(resSlice, &p)
	}

	return resSlice
}

type DeleteStack struct {
	_               string `action:"delete" entity:"stack" awsAPI:"cloudformation" awsCall:"DeleteStack" awsInput:"cloudformation.DeleteStackInput" awsOutput:"cloudformation.DeleteStackOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             cloudformationiface.CloudFormationAPI
	Name            *string   `awsName:"StackName" awsType:"awsstr" templateName:"name"`
	RetainResources []*string `awsName:"RetainResources" awsType:"awsstringslice" templateName:"retain-resources"`
}

func (cmd *DeleteStack) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"),
		params.Opt("retain-resources"),
	))
}
