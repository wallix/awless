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

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateNetworkinterface struct {
	_              string `action:"create" entity:"networkinterface" awsAPI:"ec2" awsCall:"CreateNetworkInterface" awsInput:"ec2.CreateNetworkInterfaceInput" awsOutput:"ec2.CreateNetworkInterfaceOutput" awsDryRun:""`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            ec2iface.EC2API
	Subnet         *string   `awsName:"SubnetId" awsType:"awsstr" templateName:"subnet"`
	Description    *string   `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Securitygroups []*string `awsName:"Groups" awsType:"awsstringslice" templateName:"securitygroups"`
	Privateip      *string   `awsName:"PrivateIpAddress" awsType:"awsstr" templateName:"privateip"`
}

func (cmd *CreateNetworkinterface) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("subnet"), params.Opt("description", "privateip", "securitygroups")),
		params.Validators{"privateip": params.IsIP},
	)
}

func (cmd *CreateNetworkinterface) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CreateNetworkInterfaceOutput).NetworkInterface.NetworkInterfaceId)
}

type DeleteNetworkinterface struct {
	_      string `action:"delete" entity:"networkinterface" awsAPI:"ec2" awsCall:"DeleteNetworkInterface" awsInput:"ec2.DeleteNetworkInterfaceInput" awsOutput:"ec2.DeleteNetworkInterfaceOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"NetworkInterfaceId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteNetworkinterface) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type AttachNetworkinterface struct {
	_           string `action:"attach" entity:"networkinterface" awsAPI:"ec2" awsCall:"AttachNetworkInterface" awsInput:"ec2.AttachNetworkInterfaceInput" awsOutput:"ec2.AttachNetworkInterfaceOutput" awsDryRun:""`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         ec2iface.EC2API
	Id          *string `awsName:"NetworkInterfaceId" awsType:"awsstr" templateName:"id"`
	Instance    *string `awsName:"InstanceId" awsType:"awsstr" templateName:"instance"`
	DeviceIndex *int64  `awsName:"DeviceIndex" awsType:"awsint64" templateName:"device-index"`
}

func (cmd *AttachNetworkinterface) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("device-index"), params.Key("id"), params.Key("instance")))
}

func (cmd *AttachNetworkinterface) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.AttachNetworkInterfaceOutput).AttachmentId)
}

type DetachNetworkinterface struct {
	_          string `action:"detach" entity:"networkinterface" awsAPI:"ec2" awsDryRun:"manual"`
	logger     *logger.Logger
	graph      cloud.GraphAPI
	api        ec2iface.EC2API
	Attachment *string `awsName:"AttachmentId" awsType:"awsstr" templateName:"attachment"`
	Instance   *string `awsName:"InstanceId" awsType:"awsstr" templateName:"instance"`
	Id         *string `awsName:"NetworkInterfaceId" awsType:"awsstr" templateName:"id"`
	Force      *bool   `awsName:"Force" awsType:"awsbool" templateName:"force"`
}

func (cmd *DetachNetworkinterface) ParamsSpec() params.Spec {
	return params.NewSpec(params.OnlyOneOf(
		params.AllOf(params.Key("instance"), params.Key("id")),
		params.Key("attachment"),
		params.Opt("force"),
	))
}

func (cmd *DetachNetworkinterface) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}

	input := &ec2.DetachNetworkInterfaceInput{}
	input.DryRun = Bool(true)

	if cmd.Attachment != nil {
		if err := setFieldWithType(cmd.Attachment, input, "AttachmentId", awsstr, renv.Context()); err != nil {
			return nil, err
		}
	} else if cmd.Instance != nil && cmd.Id != nil {
		attachId, err := cmd.findAttachmentBetweenInstanceAndNetworkInterface(cmd.Instance, cmd.Id)
		if err == nil && attachId != "" {
			input.SetAttachmentId(attachId)
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("either required 'attachment' or ('instance' and 'id')")
	}

	if cmd.Force != nil {
		if err := setFieldWithType(cmd.Force, input, "Force", awsbool, renv.Context()); err != nil {
			return nil, err
		}
	}

	_, err := cmd.api.DetachNetworkInterface(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("networkinterface")
			cmd.logger.Verbose("dry run: detach networkinterface ok")
			return id, nil
		}
	}

	return nil, err
}

func (cmd *DetachNetworkinterface) ManualRun(renv env.Running) (interface{}, error) {
	input := &ec2.DetachNetworkInterfaceInput{}

	if cmd.Attachment != nil {
		if err := setFieldWithType(cmd.Attachment, input, "AttachmentId", awsstr, renv.Context()); err != nil {
			return nil, err
		}
	} else if cmd.Instance != nil && cmd.Id != nil {
		attachId, err := cmd.findAttachmentBetweenInstanceAndNetworkInterface(cmd.Instance, cmd.Id)
		if err == nil && attachId != "" {
			input.SetAttachmentId(attachId)
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("detach networkinterface: either required 'attachment' or ('instance' and 'id')")
	}

	if cmd.Force != nil {
		if err := setFieldWithType(cmd.Force, input, "Force", awsbool, renv.Context()); err != nil {
			return nil, err
		}
	}

	start := time.Now()
	output, err := cmd.api.DetachNetworkInterface(input)
	cmd.logger.ExtraVerbosef("ec2.DetachNetworkInterface call took %s", time.Since(start))
	return output, err
}

type CheckNetworkinterface struct {
	_       string `action:"check" entity:"networkinterface" awsAPI:"ec2"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     ec2iface.EC2API
	Id      *string `templateName:"id"`
	State   *string `templateName:"state"`
	Timeout *int64  `templateName:"timeout"`
}

func (cmd *CheckNetworkinterface) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("id"), params.Key("state"), params.Key("timeout")),
		params.Validators{
			"state": params.IsInEnumIgnoreCase("available", "attaching", "detaching", "in-use", notFoundState),
		})
}

func (cmd *CheckNetworkinterface) ManualRun(renv env.Running) (interface{}, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{cmd.Id},
	}

	c := &checker{
		description: fmt.Sprintf("network interface %s", StringValue(cmd.Id)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.DescribeNetworkInterfaces(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "NetworkInterfaceNotFound" {
						return notFoundState, nil
					}
				} else {
					return "", err
				}
			} else {
				for _, neti := range output.NetworkInterfaces {
					if StringValue(neti.NetworkInterfaceId) == StringValue(cmd.Id) {
						return StringValue(neti.Status), nil
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

func (cmd *DetachNetworkinterface) findAttachmentBetweenInstanceAndNetworkInterface(instanceId, netInterfaceId *string) (string, error) {
	filters := &ec2.DescribeInstancesInput{}
	filters.SetFilters([]*ec2.Filter{
		{Name: String("network-interface.network-interface-id"), Values: []*string{netInterfaceId}},
		{Name: String("instance-id"), Values: []*string{instanceId}},
	})
	if out, err := cmd.api.DescribeInstances(filters); err != nil {
		return "", err
	} else if reserv := out.Reservations; len(reserv) == 1 && len(reserv[0].Instances) == 1 {
		for _, neti := range reserv[0].Instances[0].NetworkInterfaces {
			if StringValue(netInterfaceId) == StringValue(neti.NetworkInterfaceId) {
				return StringValue(neti.Attachment.AttachmentId), nil
			}
		}
	}
	return "", fmt.Errorf("not found: attachment between instance '%s' and network interface '%s'", StringValue(instanceId), StringValue(netInterfaceId))
}
