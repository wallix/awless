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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/wallix/awless/cloud/graph"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateImage struct {
	_           string `action:"create" entity:"image" awsAPI:"ec2" awsCall:"CreateImage" awsInput:"ec2.CreateImageInput" awsOutput:"ec2.CreateImageOutput" awsDryRun:"true"`
	logger      *logger.Logger
	graph       cloudgraph.GraphAPI
	api         ec2iface.EC2API
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name" required:""`
	Instance    *string `awsName:"InstanceId" awsType:"awsstr" templateName:"instance" required:""`
	Reboot      *bool   `awsName:"NoReboot" awsType:"awsbool" templateName:"reboot"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateImage) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateImage) Validate_Name() (err error) {
	if name := cmd.Name; name != nil {
		if len(*name) < 3 {
			err = errors.New("should at least be 3 characters")
		}
	}
	return
}

func (cmd *CreateImage) BeforeRun(ctx map[string]interface{}) error {
	if reboot := cmd.Reboot; reboot != nil && *reboot {
		cmd.Reboot = nil
	} else {
		cmd.Reboot = Bool(true) // so that ec2.CreateImageInput.NoReboot = true and therefore by default no reboot from AWS
	}
	return nil
}

func (cmd *CreateImage) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CreateImageOutput).ImageId)
}

type UpdateImage struct {
	_            string `action:"update" entity:"image" awsAPI:"ec2" awsDryRun:"manual"`
	logger       *logger.Logger
	graph        cloudgraph.GraphAPI
	api          ec2iface.EC2API
	Id           *string   `awsName:"ImageId" awsType:"awsstr" templateName:"id" required:""`
	Groups       []*string `awsName:"UserGroups" awsType:"awsstringslice" templateName:"groups"`
	Accounts     []*string `awsName:"UserIds" awsType:"awsstringslice" templateName:"accounts"`
	Operation    *string   `awsName:"OperationType" awsType:"awsstr" templateName:"operation"`
	ProductCodes []*string `awsName:"ProductCodes" awsType:"awsstringslice" templateName:"product-codes"`
	Description  *string   `awsName:"Description" awsType:"awsstringattribute" templateName:"description"`
}

func (cmd *UpdateImage) ValidateParams(params []string) ([]string, error) {
	return paramRule{
		tree: allOf(node("id"), oneOf(allOf(node("operation"), oneOf(node("accounts"), node("groups"))), node("product-codes"), node("description"))),
	}.verify(params)
}

func (cmd *UpdateImage) prepareImageAttributeInput(ctx map[string]interface{}) (*ec2.ModifyImageAttributeInput, error) {
	input := &ec2.ModifyImageAttributeInput{}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifyImageAttributeInput: %s", err)
	}
	if cmd.Accounts != nil || cmd.Groups != nil {
		input.SetAttribute("launchPermission")
	}
	if cmd.ProductCodes != nil {
		input.SetAttribute("productCodes")
	}
	if cmd.Description != nil {
		input.SetAttribute("description")
	}
	return input, nil
}

func (cmd *UpdateImage) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	input, err := cmd.prepareImageAttributeInput(ctx)
	if err != nil {
		return nil, err
	}
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifyImageAttributeInput: %s", err)
	}
	start := time.Now()
	output, err := cmd.api.ModifyImageAttribute(input)
	cmd.logger.ExtraVerbosef("ec2.ModifyImageAttributeInput call took %s", time.Since(start))
	return output, err
}

func (cmd *UpdateImage) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}
	input, err := cmd.prepareImageAttributeInput(ctx)
	if err != nil {
		return nil, err
	}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, ctx); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.ModifyImageAttributeInput: %s", err)
	}

	start := time.Now()
	_, err = cmd.api.ModifyImageAttribute(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			cmd.logger.ExtraVerbosef("dry run: ec2.ec2.ModifyImageAttribute call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: update image ok")
			return fakeDryRunId("image"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

type CopyImage struct {
	_            string `action:"copy" entity:"image" awsAPI:"ec2" awsCall:"CopyImage" awsInput:"ec2.CopyImageInput" awsOutput:"ec2.CopyImageOutput" awsDryRun:""`
	logger       *logger.Logger
	graph        cloudgraph.GraphAPI
	api          ec2iface.EC2API
	Name         *string `awsName:"Name" awsType:"awsstr" templateName:"name" required:""`
	SourceId     *string `awsName:"SourceImageId" awsType:"awsstr" templateName:"source-id" required:""`
	SourceRegion *string `awsName:"SourceRegion" awsType:"awsstr" templateName:"source-region" required:""`
	Encrypted    *bool   `awsName:"Encrypted" awsType:"awsbool" templateName:"encrypted"`
	Description  *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CopyImage) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CopyImage) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CopyImageOutput).ImageId)
}

type ImportImage struct {
	_            string `action:"import" entity:"image" awsAPI:"ec2" awsCall:"ImportImage" awsInput:"ec2.ImportImageInput" awsOutput:"ec2.ImportImageOutput" awsDryRun:""`
	logger       *logger.Logger
	graph        cloudgraph.GraphAPI
	api          ec2iface.EC2API
	Architecture *string `awsName:"Architecture" awsType:"awsstr" templateName:"architecture"`
	Description  *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	License      *string `awsName:"LicenseType" awsType:"awsstr" templateName:"license"`
	Platform     *string `awsName:"Platform" awsType:"awsstr" templateName:"platform"`
	Role         *string `awsName:"RoleName" awsType:"awsstr" templateName:"role"`
	Snapshot     *string `awsName:"DiskContainers[0]SnapshotId" awsType:"awsslicestruct" templateName:"snapshot"`
	Url          *string `awsName:"DiskContainers[0]Url" awsType:"awsslicestruct" templateName:"url"`
	Bucket       *string `awsName:"DiskContainers[0]UserBucket.S3Bucket" awsType:"awsslicestruct" templateName:"bucket"`
	S3object     *string `awsName:"DiskContainers[0]UserBucket.S3Key" awsType:"awsslicestruct" templateName:"s3object"`
}

func (cmd *ImportImage) ValidateParams(params []string) ([]string, error) {
	return paramRule{
		tree:   oneOfE(node("snapshot"), node("url"), allOf(node("bucket"), node("s3object"))),
		extras: []string{"architecture", "description", "license", "platform", "role"},
	}.verify(params)
}

func (cmd *ImportImage) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.ImportImageOutput).ImportTaskId)
}

type DeleteImage struct {
	_               string `action:"delete" entity:"image" awsAPI:"ec2" awsDryRun:"manual"`
	logger          *logger.Logger
	graph           cloudgraph.GraphAPI
	api             ec2iface.EC2API
	Id              *string `templateName:"id" required:""`
	DeleteSnapshots *bool   `templateName:"delete-snapshots"`
}

func (cmd *DeleteImage) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *DeleteImage) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}
	input := &ec2.DeregisterImageInput{}
	input.DryRun = Bool(true)

	if err := setFieldWithType(cmd.Id, input, "ImageId", awsstr); err != nil {
		return nil, err
	}

	if BoolValue(cmd.DeleteSnapshots) {
		var snaps []string
		var err error
		if snaps, err = cmd.imageSnapshots(StringValue(input.ImageId)); err != nil {
			return nil, err
		}
		if len(snaps) > 0 {
			cmd.logger.Infof("deleting image will also delete snapshot %s (prevent that by appending `delete-snapshots=false`)", strings.Join(snaps, ", "))
		}
	}

	_, err := cmd.api.DeregisterImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("image")
			cmd.logger.Verbose("dry run: delete image ok")
			return id, nil
		}
	}

	return nil, err
}

func (cmd *DeleteImage) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	input := &ec2.DeregisterImageInput{}

	if err := setFieldWithType(cmd.Id, input, "ImageId", awsstr); err != nil {
		return nil, err
	}

	var snaps []string
	var err error
	if BoolValue(cmd.DeleteSnapshots) {
		if snaps, err = cmd.imageSnapshots(StringValue(input.ImageId)); err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.DeregisterImageOutput
	if output, err = cmd.api.DeregisterImage(input); err != nil {
		return nil, err
	}
	cmd.logger.ExtraVerbosef("ec2.DeregisterImage call took %s", time.Since(start))

	if BoolValue(cmd.DeleteSnapshots) {
		for _, snap := range snaps {
			deleteSnapshot := CommandFactory.Build("deletesnapshot")().(*DeleteSnapshot)
			deleteSnapshot.Id = String(snap)
			if errs := deleteSnapshot.ValidateCommand(nil, nil); len(errs) > 0 {
				return nil, fmt.Errorf("%v", errs)
			}
			if _, err := deleteSnapshot.Run(ctx, nil); err != nil {
				return nil, fmt.Errorf("delete snapshot %s: %s", snap, err)
			}
		}
	}
	return output, nil
}

func (cmd *DeleteImage) imageSnapshots(id string) ([]string, error) {
	var snapshots []string
	imgs, err := cmd.api.DescribeImages(&ec2.DescribeImagesInput{ImageIds: []*string{String(id)}})
	if err != nil {
		return snapshots, err
	}
	if len(imgs.Images) == 0 {
		return snapshots, fmt.Errorf("no image found with id '%s'", id)
	}
	if len(imgs.Images) > 1 {
		return snapshots, fmt.Errorf("multiple images found with id '%s'", id)
	}
	for _, dev := range imgs.Images[0].BlockDeviceMappings {
		if snapshot := StringValue(dev.Ebs.SnapshotId); snapshot != "" {
			snapshots = append(snapshots, snapshot)
		}
	}
	return snapshots, nil
}
