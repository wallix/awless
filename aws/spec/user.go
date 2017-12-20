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
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateUser struct {
	_      string `action:"create" entity:"user" awsAPI:"iam" awsCall:"CreateUser" awsInput:"iam.CreateUserInput" awsOutput:"iam.CreateUserOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Name   *string `awsName:"UserName" awsType:"awsstr" templateName:"name"`
}

func (cmd *CreateUser) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

func (cmd *CreateUser) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*iam.CreateUserOutput).User.UserId)
}

type DeleteUser struct {
	_      string `action:"delete" entity:"user" awsAPI:"iam" awsCall:"DeleteUser" awsInput:"iam.DeleteUserInput" awsOutput:"iam.DeleteUserOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Name   *string `awsName:"UserName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteUser) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type AttachUser struct {
	_      string `action:"attach" entity:"user" awsAPI:"iam" awsCall:"AddUserToGroup" awsInput:"iam.AddUserToGroupInput" awsOutput:"iam.AddUserToGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Group  *string `awsName:"GroupName" awsType:"awsstr" templateName:"group"`
	Name   *string `awsName:"UserName" awsType:"awsstr" templateName:"name"`
}

func (cmd *AttachUser) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("group"), params.Key("name")))
}

type DetachUser struct {
	_      string `action:"detach" entity:"user" awsAPI:"iam" awsCall:"RemoveUserFromGroup" awsInput:"iam.RemoveUserFromGroupInput" awsOutput:"iam.RemoveUserFromGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Group  *string `awsName:"GroupName" awsType:"awsstr" templateName:"group"`
	Name   *string `awsName:"UserName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DetachUser) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("group"), params.Key("name")))
}
