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
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateDbsubnetgroup struct {
	_           string `action:"create" entity:"dbsubnetgroup" awsAPI:"rds" awsCall:"CreateDBSubnetGroup" awsInput:"rds.CreateDBSubnetGroupInput" awsOutput:"rds.CreateDBSubnetGroupOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         rdsiface.RDSAPI
	Name        *string   `awsName:"DBSubnetGroupName" awsType:"awsstr" templateName:"name"`
	Description *string   `awsName:"DBSubnetGroupDescription" awsType:"awsstr" templateName:"description"`
	Subnets     []*string `awsName:"SubnetIds" awsType:"awsstringslice" templateName:"subnets"`
}

func (cmd *CreateDbsubnetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("description"), params.Key("name"), params.Key("subnets")))
}

func (cmd *CreateDbsubnetgroup) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*rds.CreateDBSubnetGroupOutput).DBSubnetGroup.DBSubnetGroupName)
}

type DeleteDbsubnetgroup struct {
	_      string `action:"delete" entity:"dbsubnetgroup" awsAPI:"rds" awsCall:"DeleteDBSubnetGroup" awsInput:"rds.DeleteDBSubnetGroupInput" awsOutput:"rds.DeleteDBSubnetGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    rdsiface.RDSAPI
	Name   *string `awsName:"DBSubnetGroupName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteDbsubnetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}
