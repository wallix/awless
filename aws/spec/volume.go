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
	"fmt"
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateVolume struct {
	_                string `action:"create" entity:"volume" awsAPI:"ec2" awsCall:"CreateVolume" awsInput:"ec2.CreateVolumeInput" awsOutput:"ec2.Volume" awsDryRun:""`
	logger           *logger.Logger
	graph            cloud.GraphAPI
	api              ec2iface.EC2API
	Availabilityzone *string `awsName:"AvailabilityZone" awsType:"awsstr" templateName:"availabilityzone"`
	Size             *int64  `awsName:"Size" awsType:"awsint64" templateName:"size"`
}

func (cmd *CreateVolume) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("availabilityzone"), params.Key("size")))
}

func (cmd *CreateVolume) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.Volume).VolumeId)
}

type CheckVolume struct {
	_       string `action:"check" entity:"volume" awsAPI:"ec2"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     ec2iface.EC2API
	Id      *string `templateName:"id"`
	State   *string `templateName:"state"`
	Timeout *int64  `templateName:"timeout"`
}

func (cmd *CheckVolume) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("id"), params.Key("state"), params.Key("timeout")),
		params.Validators{
			"state": params.IsInEnumIgnoreCase("available", "in-use", notFoundState),
		},
	)
}

func (cmd *CheckVolume) ManualRun(renv env.Running) (interface{}, error) {
	input := &ec2.DescribeVolumesInput{VolumeIds: []*string{cmd.Id}}

	c := &checker{
		description: fmt.Sprintf("volume %s", StringValue(cmd.Id)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.DescribeVolumes(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "VolumeNotFound" {
						return notFoundState, nil
					}
				} else {
					return "", err
				}
			} else {
				for _, vol := range output.Volumes {
					if StringValue(vol.VolumeId) == StringValue(cmd.Id) {
						return StringValue(vol.State), nil
					}
				}
			}
			return notFoundState, nil
		},
		expect: StringValue(cmd.State),
		logger: cmd.logger,
	}
	return nil, c.check()
}

type DeleteVolume struct {
	_      string `action:"delete" entity:"volume" awsAPI:"ec2" awsCall:"DeleteVolume" awsInput:"ec2.DeleteVolumeInput" awsOutput:"ec2.DeleteVolumeOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"VolumeId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteVolume) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type AttachVolume struct {
	_        string `action:"attach" entity:"volume" awsAPI:"ec2" awsCall:"AttachVolume" awsInput:"ec2.AttachVolumeInput" awsOutput:"ec2.VolumeAttachment" awsDryRun:""`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      ec2iface.EC2API
	Device   *string `awsName:"Device" awsType:"awsstr" templateName:"device"`
	Id       *string `awsName:"VolumeId" awsType:"awsstr" templateName:"id"`
	Instance *string `awsName:"InstanceId" awsType:"awsstr" templateName:"instance"`
}

func (cmd *AttachVolume) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("device"), params.Key("id"), params.Key("instance")))
}
func (cmd *AttachVolume) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.VolumeAttachment).VolumeId)
}

type DetachVolume struct {
	_        string `action:"detach" entity:"volume" awsAPI:"ec2" awsCall:"DetachVolume" awsInput:"ec2.DetachVolumeInput" awsOutput:"ec2.VolumeAttachment" awsDryRun:""`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      ec2iface.EC2API
	Device   *string `awsName:"Device" awsType:"awsstr" templateName:"device"`
	Id       *string `awsName:"VolumeId" awsType:"awsstr" templateName:"id"`
	Instance *string `awsName:"InstanceId" awsType:"awsstr" templateName:"instance"`
	Force    *bool   `awsName:"Force" awsType:"awsbool" templateName:"force"`
}

func (cmd *DetachVolume) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("device"), params.Key("id"), params.Key("instance"),
		params.Opt("force"),
	))
}

func (cmd *DetachVolume) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.VolumeAttachment).VolumeId)
}
