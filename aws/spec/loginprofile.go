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

type CreateLoginprofile struct {
	_             string `action:"create" entity:"loginprofile" awsAPI:"iam" awsCall:"CreateLoginProfile" awsInput:"iam.CreateLoginProfileInput" awsOutput:"iam.CreateLoginProfileOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           iamiface.IAMAPI
	Username      *string `awsName:"UserName" awsType:"awsstr" templateName:"username"`
	Password      *string `awsName:"Password" awsType:"awsstr" templateName:"password"`
	PasswordReset *bool   `awsName:"PasswordResetRequired" awsType:"awsbool" templateName:"password-reset"`
}

func (cmd *CreateLoginprofile) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("password"), params.Key("username"),
		params.Opt("password-reset"),
	))
}

func (cmd *CreateLoginprofile) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*iam.CreateLoginProfileOutput).LoginProfile.UserName)
}

type UpdateLoginprofile struct {
	_             string `action:"update" entity:"loginprofile" awsAPI:"iam" awsCall:"UpdateLoginProfile" awsInput:"iam.UpdateLoginProfileInput" awsOutput:"iam.UpdateLoginProfileOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           iamiface.IAMAPI
	Username      *string `awsName:"UserName" awsType:"awsstr" templateName:"username"`
	Password      *string `awsName:"Password" awsType:"awsstr" templateName:"password"`
	PasswordReset *bool   `awsName:"PasswordResetRequired" awsType:"awsbool" templateName:"password-reset"`
}

func (cmd *UpdateLoginprofile) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("password"), params.Key("username"),
		params.Opt("password-reset"),
	))
}

type DeleteLoginprofile struct {
	_        string `action:"delete" entity:"loginprofile" awsAPI:"iam" awsCall:"DeleteLoginProfile" awsInput:"iam.DeleteLoginProfileInput" awsOutput:"iam.DeleteLoginProfileOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      iamiface.IAMAPI
	Username *string `awsName:"UserName" awsType:"awsstr" templateName:"username"`
}

func (cmd *DeleteLoginprofile) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("username")))
}
