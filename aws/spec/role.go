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
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/wallix/awless/logger"
)

type CreateRole struct {
	_                string `action:"create" entity:"role" awsAPI:"iam"`
	logger           *logger.Logger
	graph            cloud.GraphAPI
	api              iamiface.IAMAPI
	Name             *string   `awsName:"RoleName" awsType:"awsstr" templateName:"name" `
	PrincipalAccount *string   `templateName:"principal-account"`
	PrincipalUser    *string   `templateName:"principal-user"`
	PrincipalService *string   `templateName:"principal-service"`
	Conditions       []*string `templateName:"conditions"`
	SleepAfter       *int64    `templateName:"sleep-after"`
}

func (cmd *CreateRole) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"),
		params.Opt("conditions", "principal-account", "principal-service", "principal-user", "sleep-after"),
	))
}

func (cmd *CreateRole) ManualRun(renv env.Running) (interface{}, error) {
	princ := new(principal)
	if cmd.PrincipalAccount != nil {
		princ.AWS = StringValue(cmd.PrincipalAccount)
	} else if cmd.PrincipalUser != nil {
		princ.AWS = StringValue(cmd.PrincipalUser)
	} else if cmd.PrincipalService != nil {
		princ.Service = StringValue(cmd.PrincipalService)
	}

	stat, err := buildStatementFromParams(String("Allow"), nil, []*string{String("sts:AssumeRole")}, cmd.Conditions)
	if err != nil {
		return nil, err
	}
	stat.Principal = princ
	trust := &policyBody{
		Version:   "2012-10-17",
		Statement: []*policyStatement{stat},
	}

	b, err := json.MarshalIndent(trust, "", " ")
	if err != nil {
		return nil, errors.New("cannot marshal role trust policy document")
	}

	cmd.logger.ExtraVerbosef("role trust policy document json:\n%s\n", string(b))

	call := &awsCall{
		fnName: "iam.CreateRole",
		fn:     cmd.api.CreateRole,
		logger: cmd.logger,
		setters: []setter{
			{val: cmd.Name, fieldPath: "RoleName", fieldType: awsstr},
			{val: string(b), fieldPath: "AssumeRolePolicyDocument", fieldType: awsstr},
		},
	}

	output, err := call.execute(&iam.CreateRoleInput{})
	if err != nil {
		return nil, err
	}
	role := output.(*iam.CreateRoleOutput).Role

	createInstProfile := CommandFactory.Build("createinstanceprofile")().(*CreateInstanceprofile)
	createInstProfile.Name = cmd.Name
	createInstProfile.Run(renv, nil)

	attachRole := CommandFactory.Build("attachrole")().(*AttachRole)
	attachRole.Name = role.RoleName
	attachRole.Instanceprofile = role.RoleName
	attachRole.Run(renv, nil)

	if v := cmd.SleepAfter; v != nil {
		vv := Int64AsIntValue(v)
		cmd.logger.Infof("sleeping for %d seconds", vv)
		time.Sleep(time.Duration(vv) * time.Second)
	}

	return output, nil
}

func (cmd *CreateRole) ExtractResult(i interface{}) string {
	return StringValue(i.(*iam.CreateRoleOutput).Role.Arn)
}

type DeleteRole struct {
	_      string `action:"delete" entity:"role" awsAPI:"iam"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Name   *string `awsName:"RoleName" awsType:"awsstr" templateName:"name" `
}

func (cmd *DeleteRole) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

func (cmd *DeleteRole) ManualRun(renv env.Running) (interface{}, error) {
	detachRole := CommandFactory.Build("detachrole")().(*DetachRole)
	detachRole.Name = cmd.Name
	detachRole.Instanceprofile = cmd.Name
	detachRole.Run(renv, nil)

	deleteInstProfile := CommandFactory.Build("deleteinstanceprofile")().(*DeleteInstanceprofile)
	deleteInstProfile.Name = cmd.Name
	deleteInstProfile.Run(renv, nil)

	input := &iam.DeleteRoleInput{}
	if err := setFieldWithType(cmd.Name, input, "RoleName", awsstr); err != nil {
		return nil, err
	}

	start := time.Now()
	output, err := cmd.api.DeleteRole(input)
	cmd.logger.ExtraVerbosef("iam.DeleteRole call took %s", time.Since(start))
	return output, err
}

type AttachRole struct {
	_               string `action:"attach" entity:"role" awsAPI:"iam" awsCall:"AddRoleToInstanceProfile" awsInput:"iam.AddRoleToInstanceProfileInput" awsOutput:"iam.AddRoleToInstanceProfileOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             iamiface.IAMAPI
	Instanceprofile *string `awsName:"InstanceProfileName" awsType:"awsstr" templateName:"instanceprofile" `
	Name            *string `awsName:"RoleName" awsType:"awsstr" templateName:"name" `
}

func (cmd *AttachRole) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("instanceprofile"), params.Key("name")))
}

type DetachRole struct {
	_               string `action:"detach" entity:"role" awsAPI:"iam" awsCall:"RemoveRoleFromInstanceProfile" awsInput:"iam.RemoveRoleFromInstanceProfileInput" awsOutput:"iam.RemoveRoleFromInstanceProfileOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             iamiface.IAMAPI
	Instanceprofile *string `awsName:"InstanceProfileName" awsType:"awsstr" templateName:"instanceprofile" `
	Name            *string `awsName:"RoleName" awsType:"awsstr" templateName:"name" `
}

func (cmd *DetachRole) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("instanceprofile"), params.Key("name")))
}
