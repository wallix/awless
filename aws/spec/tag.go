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
	"strings"
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateTag struct {
	_        string `action:"create" entity:"tag" awsAPI:"ec2" awsDryRun:"manual"` //  awsCall:"CreateTags" awsInput:"ec2.CreateTagsInput" awsOutput:"ec2.CreateTagsOutput"
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      ec2iface.EC2API
	Resource *string `awsName:"Resources" awsType:"awsstringslice" templateName:"resource"`
	Key      *string `templateName:"key"`
	Value    *string `templateName:"value"`
}

func (cmd *CreateTag) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("key"), params.Key("resource"), params.Key("value")))
}

func (cmd *CreateTag) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %s", err)
	}

	input := &ec2.CreateTagsInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateTagsInput: %s", err)
	}
	input.Tags = []*ec2.Tag{{Key: cmd.Key, Value: cmd.Value}}

	start := time.Now()
	_, err := cmd.api.CreateTags(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateTags call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create tag ok")
			return fakeDryRunId("tag"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %s", err)
}

func (cmd *CreateTag) ManualRun(renv env.Running) (interface{}, error) {
	input := &ec2.CreateTagsInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateTagsInput: %s", err)
	}
	input.Tags = []*ec2.Tag{{Key: cmd.Key, Value: cmd.Value}}

	start := time.Now()
	req, _ := cmd.api.CreateTagsRequest(input)
	req.Retryer = createTagRetryer{}
	if err := req.Send(); err != nil {
		return nil, err
	}
	cmd.logger.ExtraVerbosef("ec2.CreateTags call took %s", time.Since(start))
	return nil, nil
}

type DeleteTag struct {
	_        string `action:"delete" entity:"tag" awsAPI:"ec2" awsDryRun:"manual"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      ec2iface.EC2API
	Resource *string `awsName:"Resources" awsType:"awsstringslice" templateName:"resource"`
	Key      *string `templateName:"key"`
	Value    *string `templateName:"value"`
}

func (cmd *DeleteTag) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("key"), params.Key("resource"),
		params.Opt("value"),
	))
}

func (cmd *DeleteTag) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DeleteTagsInput{}
	input.SetDryRun(true)
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteTagsInput: %s", err)
	}
	input.Tags = []*ec2.Tag{{Key: cmd.Key, Value: cmd.Value}}

	start := time.Now()
	_, err := cmd.api.DeleteTags(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteTags call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create tag ok")
			return fakeDryRunId("tag"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteTag) ManualRun(renv env.Running) (interface{}, error) {
	input := &ec2.DeleteTagsInput{}
	if err := structInjector(cmd, input, renv.Context()); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteTagsInput: %s", err)
	}
	input.Tags = []*ec2.Tag{{Key: cmd.Key, Value: cmd.Value}}

	start := time.Now()
	_, err := cmd.api.DeleteTags(input)
	cmd.logger.ExtraVerbosef("ec2.DeleteTags call took %s", time.Since(start))
	return nil, err
}

func createNameTag(resource, name *string, renv env.Running) error {
	createTag := CommandFactory.Build("createtag")().(*CreateTag)
	entries := map[string]interface{}{
		"key":      "Name",
		"value":    name,
		"resource": resource,
	}
	if err := params.Validate(createTag.ParamsSpec().Validators(), entries); err != nil {
		return err
	}
	_, err := createTag.Run(renv, entries)
	return err
}

type createTagRetryer struct {
	client.DefaultRetryer
}

func (d createTagRetryer) MaxRetries() int { return 5 }
func (d createTagRetryer) ShouldRetry(r *request.Request) bool {
	if d.DefaultRetryer.ShouldRetry(r) || !(r.HTTPResponse.StatusCode < 300 && r.HTTPResponse.StatusCode >= 200) {
		return true
	}

	return false
}
