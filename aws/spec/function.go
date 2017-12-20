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
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateFunction struct {
	_             string `action:"create" entity:"function" awsAPI:"lambda" awsCall:"CreateFunction" awsInput:"lambda.CreateFunctionInput" awsOutput:"lambda.FunctionConfiguration"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           lambdaiface.LambdaAPI
	Name          *string `awsName:"FunctionName" awsType:"awsstr" templateName:"name"`
	Handler       *string `awsName:"Handler" awsType:"awsstr" templateName:"handler"`
	Role          *string `awsName:"Role" awsType:"awsstr" templateName:"role"`
	Runtime       *string `awsName:"Runtime" awsType:"awsstr" templateName:"runtime"`
	Bucket        *string `awsName:"Code.S3Bucket" awsType:"awsstr" templateName:"bucket"`
	Object        *string `awsName:"Code.S3Key" awsType:"awsstr" templateName:"object"`
	Objectversion *string `awsName:"Code.S3ObjectVersion" awsType:"awsstr" templateName:"objectversion"`
	Zipfile       *string `awsName:"Code.ZipFile" awsType:"awsfiletobyteslice" templateName:"zipfile"`
	Description   *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Memory        *int64  `awsName:"MemorySize" awsType:"awsint64" templateName:"memory"`
	Publish       *bool   `awsName:"Publish" awsType:"awsbool" templateName:"publish"`
	Timeout       *int64  `awsName:"Timeout" awsType:"awsint64" templateName:"timeout"`
}

func (cmd *CreateFunction) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("handler"), params.Key("name"), params.Key("role"), params.Key("runtime"),
		params.Opt("bucket", "description", "memory", "object", "objectversion", "publish", "timeout", "zipfile"),
	))
}

func (cmd *CreateFunction) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*lambda.FunctionConfiguration).FunctionArn)
}

type DeleteFunction struct {
	_       string `action:"delete" entity:"function" awsAPI:"lambda" awsCall:"DeleteFunction" awsInput:"lambda.DeleteFunctionInput" awsOutput:"lambda.DeleteFunctionOutput"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     lambdaiface.LambdaAPI
	Id      *string `awsName:"FunctionName" awsType:"awsstr" templateName:"id"`
	Version *string `awsName:"Qualifier" awsType:"awsstr" templateName:"version"`
}

func (cmd *DeleteFunction) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"),
		params.Opt("version"),
	))
}
