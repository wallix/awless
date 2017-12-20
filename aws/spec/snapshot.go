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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateSnapshot struct {
	_           string `action:"create" entity:"snapshot" awsAPI:"ec2" awsCall:"CreateSnapshot" awsInput:"ec2.CreateSnapshotInput" awsOutput:"ec2.Snapshot" awsDryRun:""`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         ec2iface.EC2API
	Volume      *string `awsName:"VolumeId" awsType:"awsstr" templateName:"volume"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateSnapshot) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("volume"),
		params.Opt("description"),
	))
}

func (cmd *CreateSnapshot) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.Snapshot).SnapshotId)
}

type DeleteSnapshot struct {
	_      string `action:"delete" entity:"snapshot" awsAPI:"ec2" awsCall:"DeleteSnapshot" awsInput:"ec2.DeleteSnapshotInput" awsOutput:"ec2.DeleteSnapshotOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"SnapshotId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteSnapshot) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type CopySnapshot struct {
	_            string `action:"copy" entity:"snapshot" awsAPI:"ec2" awsCall:"CopySnapshot" awsInput:"ec2.CopySnapshotInput" awsOutput:"ec2.CopySnapshotOutput" awsDryRun:""`
	logger       *logger.Logger
	graph        cloud.GraphAPI
	api          ec2iface.EC2API
	SourceId     *string `awsName:"SourceSnapshotId" awsType:"awsstr" templateName:"source-id"`
	SourceRegion *string `awsName:"SourceRegion" awsType:"awsstr" templateName:"source-region"`
	Encrypted    *bool   `awsName:"Encrypted" awsType:"awsbool" templateName:"encrypted"`
	Description  *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CopySnapshot) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("source-id"), params.Key("source-region"),
		params.Opt("description", "encrypted"),
	))
}

func (cmd *CopySnapshot) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CopySnapshotOutput).SnapshotId)
}
