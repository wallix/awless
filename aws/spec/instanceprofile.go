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

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/wallix/awless/logger"
)

type CreateInstanceprofile struct {
	_      string `action:"create" entity:"instanceprofile" awsAPI:"iam" awsCall:"CreateInstanceProfile" awsInput:"iam.CreateInstanceProfileInput" awsOutput:"iam.CreateInstanceProfileOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Name   *string `awsName:"InstanceProfileName" awsType:"awsstr" templateName:"name"`
}

func (cmd *CreateInstanceprofile) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type DeleteInstanceprofile struct {
	_      string `action:"delete" entity:"instanceprofile" awsAPI:"iam" awsCall:"DeleteInstanceProfile" awsInput:"iam.DeleteInstanceProfileInput" awsOutput:"iam.DeleteInstanceProfileOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Name   *string `awsName:"InstanceProfileName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteInstanceprofile) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type AttachInstanceprofile struct {
	_        string `action:"attach" entity:"instanceprofile" awsAPI:"ec2" awsDryRun:"manual"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      ec2iface.EC2API
	Instance *string `awsName:"InstanceId" awsType:"awsstr" templateName:"instance"`
	Name     *string `awsName:"IamInstanceProfile.Name" awsType:"awsstr" templateName:"name"`
	Replace  *bool   `templateName:"replace"`
}

func (cmd *AttachInstanceprofile) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("instance"), params.Key("name"),
		params.Opt("replace"),
	))
}

func (cmd *AttachInstanceprofile) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}
	if BoolValue(cmd.Replace) {
		in := &ec2.DescribeIamInstanceProfileAssociationsInput{
			Filters: []*ec2.Filter{
				{Name: String("instance-id"), Values: []*string{cmd.Instance}},
			},
		}
		out, err := cmd.api.DescribeIamInstanceProfileAssociations(in)
		if err != nil {
			return nil, fmt.Errorf("replace mode on: cannot get: %s", err)
		}
		if assocs := out.IamInstanceProfileAssociations; len(assocs) > 0 {
			for _, ass := range assocs {
				cmd.logger.ExtraVerbosef("dry run: attach instanceprofile: existing instanceprofile %s (state: %s) on instance %s", StringValue(ass.IamInstanceProfile.Id), StringValue(ass.State), StringValue(cmd.Instance))
			}
		}
	}
	cmd.logger.Verbose("params dry run: attach instanceprofile ok")
	return fakeDryRunId("instanceprofile"), nil
}

func (cmd *AttachInstanceprofile) ManualRun(renv env.Running) (interface{}, error) {
	instanceId := StringValue(cmd.Instance)
	profileName := StringValue(cmd.Name)

	if BoolValue(cmd.Replace) {
		out, err := cmd.api.DescribeIamInstanceProfileAssociations(
			&ec2.DescribeIamInstanceProfileAssociationsInput{
				Filters: []*ec2.Filter{
					{Name: String("instance-id"), Values: []*string{String(instanceId)}},
					{Name: String("state"), Values: []*string{String("associated")}},
				},
			})
		if err != nil {
			logger.Warningf("attach instanceprofile: replace mode on: dry run was ok but now cannot get instance profile association")
		}
		assoc := out.IamInstanceProfileAssociations
		if len(assoc) > 0 {
			start := time.Now()
			assocId := StringValue(assoc[0].AssociationId)
			assocInstId := StringValue(assoc[0].InstanceId)
			oldProfileArn := StringValue(assoc[0].IamInstanceProfile.Arn)
			cmd.logger.ExtraVerbosef("attach profile: found existing profile to replace with %s", profileName)
			if assocInstId == instanceId {
				out, err := cmd.api.ReplaceIamInstanceProfileAssociation(
					&ec2.ReplaceIamInstanceProfileAssociationInput{
						AssociationId: String(assocId),
						IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
							Name: String(profileName),
						},
					})
				if err != nil {
					return nil, fmt.Errorf("attach instanceprofile: replace mode on: cannot replace with new instance profile: %s", err)
				}

				cmd.logger.Verbosef("attach profile: replaced profile '%s' with '%s' on instance %s", oldProfileArn, profileName, instanceId)
				cmd.logger.ExtraVerbosef("ec2.ReplaceIamInstanceProfileAssociation call took %s", time.Since(start))

				return out, nil
			}
		}
	}

	input := &ec2.AssociateIamInstanceProfileInput{}
	if err := setFieldWithType(instanceId, input, "InstanceId", awsstr, renv.Context()); err != nil {
		return nil, err
	}
	if err := setFieldWithType(profileName, input, "IamInstanceProfile.Name", awsstr, renv.Context()); err != nil {
		return nil, err
	}

	start := time.Now()
	output, err := cmd.api.AssociateIamInstanceProfile(input)
	cmd.logger.ExtraVerbosef("ec2.AssociateIamInstanceProfile call took %s", time.Since(start))
	return output, err
}

type DetachInstanceprofile struct {
	_        string `action:"detach" entity:"instanceprofile" awsAPI:"ec2"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      ec2iface.EC2API
	Instance *string `templateName:"instance"`
	Name     *string `templateName:"name"`
}

func (cmd *DetachInstanceprofile) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("instance"), params.Key("name")))
}

func (cmd *DetachInstanceprofile) ManualRun(renv env.Running) (interface{}, error) {
	instanceId := StringValue(cmd.Instance)
	profileName := StringValue(cmd.Name)

	out, err := cmd.api.DescribeIamInstanceProfileAssociations(
		&ec2.DescribeIamInstanceProfileAssociationsInput{
			Filters: []*ec2.Filter{
				{Name: String("instance-id"), Values: []*string{String(instanceId)}},
			},
		})
	if err != nil {
		return nil, fmt.Errorf("cannot list profile on instance %s: %s", instanceId, err)
	}

	assocs := out.IamInstanceProfileAssociations
	if len(assocs) < 1 {
		cmd.logger.Infof("detach instanceprofile: nothing to be detached on instance %s", instanceId)
		return nil, nil
	}

	var lastId string
	for _, ass := range assocs {
		if strings.Contains(StringValue(ass.IamInstanceProfile.Arn), profileName) {
			input := &ec2.DisassociateIamInstanceProfileInput{
				AssociationId: ass.AssociationId,
			}

			start := time.Now()
			output, err := cmd.api.DisassociateIamInstanceProfile(input)
			if err != nil {
				return nil, err
			}
			cmd.logger.ExtraVerbosef("ec2.DisassociateIamInstanceProfile call took %s", time.Since(start))
			id := StringValue(output.IamInstanceProfileAssociation.IamInstanceProfile.Id)
			lastId = id
		}
	}

	return lastId, nil
}

func (cmd *DetachInstanceprofile) ExtractResult(i interface{}) string {
	if i != nil {
		return fmt.Sprint(i)
	}
	return ""
}
