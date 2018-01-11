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

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"
)

type CreateInstance struct {
	_              string `action:"create" entity:"instance" awsAPI:"ec2" awsCall:"RunInstances" awsInput:"ec2.RunInstancesInput" awsOutput:"ec2.Reservation" awsDryRun:""`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            ec2iface.EC2API
	Image          *string   `awsName:"ImageId" awsType:"awsstr" templateName:"image"`
	Count          *int64    `awsName:"MaxCount,MinCount" awsType:"awsin64" templateName:"count"`
	Type           *string   `awsName:"InstanceType" awsType:"awsstr" templateName:"type"`
	Name           *string   `templateName:"name"`
	Subnet         *string   `awsName:"SubnetId" awsType:"awsstr" templateName:"subnet"`
	Keypair        *string   `awsName:"KeyName" awsType:"awsstr" templateName:"keypair"`
	PrivateIP      *string   `awsName:"PrivateIpAddress" awsType:"awsstr" templateName:"ip"`
	UserData       *string   `awsName:"UserData" awsType:"awsuserdatatobase64" templateName:"userdata"`
	SecurityGroups []*string `awsName:"SecurityGroupIds" awsType:"awsstringslice" templateName:"securitygroup"`
	Lock           *bool     `awsName:"DisableApiTermination" awsType:"awsbool" templateName:"lock"`
	Role           *string   `awsName:"IamInstanceProfile.Name" awsType:"awsstr" templateName:"role"`
	DistroQuery    *string   `awsType:"awsstr" templateName:"distro"`
}

func (cmd *CreateInstance) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(
		params.AllOf(params.OnlyOneOf(params.Key("distro"), params.Key("image")),
			params.Key("count"), params.Key("type"), params.Key("name"), params.Key("subnet"),
			params.Opt(params.Suggested("keypair", "securitygroup"), "ip", "userdata", "lock", "role"),
		),
		params.Validators{"ip": params.IsIP},
	)
	builder.AddReducer(cmd.convertDistroToAMI, "distro")
	return builder.Done()
}

func (cmd *CreateInstance) convertDistroToAMI(values map[string]interface{}) (map[string]interface{}, error) {
	if distro, ok := values["distro"].(string); ok {
		query, err := ParseImageQuery(distro)
		if err != nil {
			return nil, fmt.Errorf("distro: %s", err)
		}
		resolver := ImageResolver(cmd.api.DescribeImages)
		cmd.logger.Verbosef("Searching for bare community distro: '%s' expanded to '%s'", distro, query)
		images, fromCache, err := resolver.Resolve(query)
		if err != nil {
			return nil, fmt.Errorf("distro: %s", err)
		}
		if len(images) > 0 {
			var caching string
			if fromCache {
				caching = " from cache"
			}
			cmd.logger.Infof("Image %s resolved%s for distro '%s' (expanded to '%s')", images[0].Id, caching, distro, query)
			return map[string]interface{}{"image": images[0].Id}, nil
		} else {
			return nil, fmt.Errorf("distro: no image id found for query '%s'", query)
		}
	}
	return nil, nil
}

func (cmd *CreateInstance) ExtractResult(i interface{}) string {
	return StringValue(i.(*ec2.Reservation).Instances[0].InstanceId)
}

func (cmd *CreateInstance) AfterRun(renv env.Running, output interface{}) error {
	return createNameTag(String(cmd.ExtractResult(output)), cmd.Name, renv)
}

type UpdateInstance struct {
	_      string `action:"update" entity:"instance" awsAPI:"ec2" awsCall:"ModifyInstanceAttribute" awsInput:"ec2.ModifyInstanceAttributeInput" awsOutput:"ec2.ModifyInstanceAttributeOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"InstanceId" awsType:"awsstr" templateName:"id"`
	Type   *string `awsName:"InstanceType.Value" awsType:"awsstr" templateName:"type"`
	Lock   *bool   `awsName:"DisableApiTermination" awsType:"awsboolattribute" templateName:"lock"`
}

func (cmd *UpdateInstance) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Opt("lock", "type")))
}

type DeleteInstance struct {
	_      string `action:"delete" entity:"instance" awsAPI:"ec2" awsCall:"TerminateInstances" awsInput:"ec2.TerminateInstancesInput" awsOutput:"ec2.TerminateInstancesOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	IDs    []*string `awsName:"InstanceIds" awsType:"awsstringslice" templateName:"ids"`
}

func (cmd *DeleteInstance) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.OnlyOneOf(params.Key("ids"), params.Key("id")))
	builder.AddReducer(idToIds, "id")
	return builder.Done()
}

type StartInstance struct {
	_      string `action:"start" entity:"instance" awsAPI:"ec2" awsCall:"StartInstances" awsInput:"ec2.StartInstancesInput" awsOutput:"ec2.StartInstancesOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     []*string `awsName:"InstanceIds" awsType:"awsstringslice" templateName:"ids"`
}

func (cmd *StartInstance) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.OnlyOneOf(params.Key("ids"), params.Key("id")))
	builder.AddReducer(idToIds, "id")
	return builder.Done()
}

func (cmd *StartInstance) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.StartInstancesOutput).StartingInstances[0].InstanceId)
}

type StopInstance struct {
	_      string `action:"stop" entity:"instance" awsAPI:"ec2" awsCall:"StopInstances" awsInput:"ec2.StopInstancesInput" awsOutput:"ec2.StopInstancesOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     []*string `awsName:"InstanceIds" awsType:"awsstringslice" templateName:"ids"`
}

func (cmd *StopInstance) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.OnlyOneOf(params.Key("ids"), params.Key("id")))
	builder.AddReducer(idToIds, "id")
	return builder.Done()
}

func (cmd *StopInstance) ExtractResult(i interface{}) string {
	return StringValue(i.(*ec2.StopInstancesOutput).StoppingInstances[0].InstanceId)
}

type RestartInstance struct {
	_      string `action:"restart" entity:"instance" awsAPI:"ec2" awsCall:"RebootInstances" awsInput:"ec2.RebootInstancesInput" awsOutput:"ec2.RebootInstancesOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     []*string `awsName:"InstanceIds" awsType:"awsstringslice" templateName:"ids"`
}

func (cmd *RestartInstance) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.OnlyOneOf(params.Key("ids"), params.Key("id")))
	builder.AddReducer(idToIds, "id")
	return builder.Done()
}

const (
	notFoundState = "not-found"
)

type CheckInstance struct {
	_       string `action:"check" entity:"instance" awsAPI:"ec2"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     ec2iface.EC2API
	Id      *string `templateName:"id"`
	State   *string `templateName:"state"`
	Timeout *int64  `templateName:"timeout"`
}

func (cmd *CheckInstance) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("id"), params.Key("state"), params.Key("timeout")),
		params.Validators{
			"state": params.IsInEnumIgnoreCase("pending", "running", "shutting-down", "terminated", "stopping", "stopped", notFoundState),
		},
	)
}

func (cmd *CheckInstance) ManualRun(renv env.Running) (interface{}, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{cmd.Id},
	}

	c := &checker{
		description: fmt.Sprintf("instance %s", StringValue(cmd.Id)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.DescribeInstances(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "InstanceNotFound" {
						return notFoundState, nil
					}
				} else {
					return "", err
				}
			} else {
				if res := output.Reservations; len(res) > 0 {
					if instances := output.Reservations[0].Instances; len(instances) > 0 {
						for _, inst := range instances {
							if StringValue(inst.InstanceId) == StringValue(cmd.Id) {
								return StringValue(inst.State.Name), nil
							}
						}
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

type AttachInstance struct {
	_           string `action:"attach" entity:"instance" awsAPI:"elbv2" awsCall:"RegisterTargets" awsInput:"elbv2.RegisterTargetsInput" awsOutput:"elbv2.RegisterTargetsOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         elbv2iface.ELBV2API
	Targetgroup *string `awsName:"TargetGroupArn" awsType:"awsstr" templateName:"targetgroup"`
	Id          *string `awsName:"Targets[0]Id" awsType:"awsslicestruct" templateName:"id"`
	Port        *int64  `awsName:"Targets[0]Port" awsType:"awsslicestructint64" templateName:"port"`
}

func (cmd *AttachInstance) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("targetgroup"),
		params.Opt("port"),
	))
}

type DetachInstance struct {
	_           string `action:"detach" entity:"instance" awsAPI:"elbv2" awsCall:"DeregisterTargets" awsInput:"elbv2.DeregisterTargetsInput" awsOutput:"elbv2.DeregisterTargetsOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         elbv2iface.ELBV2API
	Targetgroup *string `awsName:"TargetGroupArn" awsType:"awsstr" templateName:"targetgroup"`
	Id          *string `awsName:"Targets[0]Id" awsType:"awsslicestruct" templateName:"id"`
}

func (cmd *DetachInstance) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("targetgroup")))
}

func idToIds(values map[string]interface{}) (map[string]interface{}, error) {
	if id, hasID := values["id"]; hasID {
		return map[string]interface{}{"ids": id}, nil
	} else {
		return nil, nil
	}
}
