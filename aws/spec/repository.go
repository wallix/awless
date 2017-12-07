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
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/wallix/awless/cloud/graph"
	"github.com/wallix/awless/logger"
)

type CreateRepository struct {
	_      string `action:"create" entity:"repository" awsAPI:"ecr" awsCall:"CreateRepository" awsInput:"ecr.CreateRepositoryInput" awsOutput:"ecr.CreateRepositoryOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ecriface.ECRAPI
	Name   *string `awsName:"RepositoryName" awsType:"awsstr" templateName:"name" required:""`
}

func (cmd *CreateRepository) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateRepository) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ecr.CreateRepositoryOutput).Repository.RepositoryArn)
}

type DeleteRepository struct {
	_       string `action:"delete" entity:"repository" awsAPI:"ecr" awsCall:"DeleteRepository" awsInput:"ecr.DeleteRepositoryInput" awsOutput:"ecr.DeleteRepositoryOutput"`
	logger  *logger.Logger
	graph   cloudgraph.GraphAPI
	api     ecriface.ECRAPI
	Name    *string `awsName:"RepositoryName" awsType:"awsstr" templateName:"name" required:""`
	Force   *bool   `awsName:"Force" awsType:"awsbool" templateName:"force"`
	Account *string `awsName:"RegistryId" awsType:"awsstr" templateName:"account"`
}

func (cmd *DeleteRepository) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}
