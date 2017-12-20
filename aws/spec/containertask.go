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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/wallix/awless/logger"
)

type StartContainertask struct {
	_                         string `action:"start" entity:"containertask" awsAPI:"ecs"`
	logger                    *logger.Logger
	graph                     cloud.GraphAPI
	api                       ecsiface.ECSAPI
	Cluster                   *string `templateName:"cluster"`
	DesiredCount              *int64  `templateName:"desired-count"`
	Name                      *string `templateName:"name"`
	Type                      *string `templateName:"type"`
	Role                      *string `templateName:"role"`
	DeploymentName            *string `templateName:"deployment-name"`
	LoadBalancerContainerName *string `templateName:"loadbalancer.container-name"`
	LoadBalancerContainerPort *int64  `templateName:"loadbalancer.container-port"`
	LoadBalancerTargetgroup   *string `templateName:"loadbalancer.targetgroup"`
}

func (cmd *StartContainertask) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("cluster"), params.Key("desired-count"), params.Key("name"), params.Key("type"), params.Opt("deployment-name", "loadbalancer.container-name", "loadbalancer.container-port", "loadbalancer.targetgroup", "role")),
		params.Validators{
			"type": func(i interface{}, others map[string]interface{}) error {
				typ := fmt.Sprint(i)
				if typ != "task" && typ != "service" {
					return fmt.Errorf("expected any of [task service] but got %s", typ)
				}
				_, hasDepName := others["deployment-name"]
				if typ == "service" && !hasDepName {
					return errors.New("missing required param 'deployment-name' when type=service")
				}
				return nil
			},
		})
}

func (cmd *StartContainertask) ManualRun(renv env.Running) (interface{}, error) {
	switch StringValue(cmd.Type) {
	case "service":
		setters := []setter{
			{val: cmd.Cluster, fieldPath: "Cluster", fieldType: awsstr},
			{val: cmd.Name, fieldPath: "TaskDefinition", fieldType: awsstr},
			{val: cmd.DeploymentName, fieldPath: "ServiceName", fieldType: awsstr},
			{val: cmd.DesiredCount, fieldPath: "DesiredCount", fieldType: awsint64},
		}
		if cmd.Role != nil {
			setters = append(setters, setter{val: cmd.Role, fieldPath: "Role", fieldType: awsstr})
		}
		if cmd.LoadBalancerContainerName != nil {
			setters = append(setters, setter{val: cmd.LoadBalancerContainerName, fieldPath: "LoadBalancers[0]ContainerName", fieldType: awsslicestruct})
		}
		if cmd.LoadBalancerContainerPort != nil {
			setters = append(setters, setter{val: cmd.LoadBalancerContainerPort, fieldPath: "LoadBalancers[0]ContainerPort", fieldType: awsslicestructint64})
		}
		if cmd.LoadBalancerTargetgroup != nil {
			setters = append(setters, setter{val: cmd.LoadBalancerTargetgroup, fieldPath: "LoadBalancers[0]TargetGroupArn", fieldType: awsslicestruct})
		}

		call := &awsCall{
			fnName:  "ecs.CreateService",
			fn:      cmd.api.CreateService,
			logger:  cmd.logger,
			setters: setters,
		}

		return call.execute(&ecs.CreateServiceInput{})
	case "task":
		call := &awsCall{
			fnName: "ecs.RunTask",
			fn:     cmd.api.RunTask,
			logger: cmd.logger,
			setters: []setter{
				{val: cmd.Cluster, fieldPath: "Cluster", fieldType: awsstr},
				{val: cmd.Name, fieldPath: "TaskDefinition", fieldType: awsstr},
				{val: cmd.DesiredCount, fieldPath: "Count", fieldType: awsint64},
			},
		}

		output, err := call.execute(&ecs.RunTaskInput{})
		if err != nil {
			return nil, err
		}
		if len(output.(*ecs.RunTaskOutput).Failures) > 0 {
			return nil, fmt.Errorf("fail to run task: %s", aws.StringValue(output.(*ecs.RunTaskOutput).Failures[0].Reason))
		}
		if len(output.(*ecs.RunTaskOutput).Tasks) > 0 {
			return output, nil
		}
		return nil, fmt.Errorf("no task started successfully")
	}
	return nil, fmt.Errorf("invalid type '%s'", StringValue(cmd.Type))
}

func (cmd *StartContainertask) ExtractResult(i interface{}) string {
	switch ii := i.(type) {
	case *ecs.CreateServiceOutput:
		return StringValue(ii.Service.ServiceArn)
	case *ecs.RunTaskOutput:
		return StringValue(ii.Tasks[0].TaskArn)
	default:
		return ""
	}
}

type StopContainertask struct {
	_              string `action:"stop" entity:"containertask" awsAPI:"ecs"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            ecsiface.ECSAPI
	Cluster        *string `templateName:"cluster"`
	Type           *string `templateName:"type"`
	DeploymentName *string `templateName:"deployment-name"`
	RunArn         *string `templateName:"run-arn"`
}

func (cmd *StopContainertask) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("cluster"), params.Key("type"), params.Opt("deployment-name", "run-arn")),
		params.Validators{
			"type": func(i interface{}, others map[string]interface{}) error {
				typ := fmt.Sprint(i)
				if typ != "task" && typ != "service" {
					return fmt.Errorf("expected any of [task service] but got %s", typ)
				}
				_, hasDepName := others["deployment-name"]
				if typ == "service" && !hasDepName {
					return errors.New("missing required param 'deployment-name' when type=service")
				}
				_, hasRunARN := others["run-arn"]
				if typ == "task" && !hasRunARN {
					return errors.New("missing required param 'run-arn' when type=task")
				}
				return nil
			},
		})
}

func (cmd *StopContainertask) ManualRun(renv env.Running) (interface{}, error) {
	switch StringValue(cmd.Type) {
	case "service":
		call := &awsCall{
			fnName: "ecs.DeleteService",
			fn:     cmd.api.DeleteService,
			logger: cmd.logger,
			setters: []setter{
				{val: cmd.Cluster, fieldPath: "Cluster", fieldType: awsstr},
				{val: cmd.DeploymentName, fieldPath: "Service", fieldType: awsstr},
			},
		}
		return call.execute(&ecs.DeleteServiceInput{})
	case "task":
		call := &awsCall{
			fnName: "ecs.StopTask",
			fn:     cmd.api.StopTask,
			logger: cmd.logger,
			setters: []setter{
				{val: cmd.Cluster, fieldPath: "Cluster", fieldType: awsstr},
				{val: cmd.RunArn, fieldPath: "Task", fieldType: awsstr},
			},
		}

		return call.execute(&ecs.StopTaskInput{})
	}
	return nil, fmt.Errorf("invalid type '%s'", StringValue(cmd.Type))
}

type UpdateContainertask struct {
	_              string `action:"update" entity:"containertask" awsAPI:"ecs" awsCall:"UpdateService" awsInput:"ecs.UpdateServiceInput" awsOutput:"ecs.UpdateServiceOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            ecsiface.ECSAPI
	Cluster        *string `awsName:"Cluster" awsType:"awsstr" templateName:"cluster"`
	DeploymentName *string `awsName:"Service" awsType:"awsstr" templateName:"deployment-name"`
	DesiredCount   *int64  `awsName:"DesiredCount" awsType:"awsint64" templateName:"desired-count"`
	Name           *string `awsName:"TaskDefinition" awsType:"awsstr" templateName:"name"`
}

func (cmd *UpdateContainertask) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("cluster"), params.Key("deployment-name"),
		params.Opt("desired-count", "name"),
	))
}

type AttachContainertask struct {
	_               string `action:"attach" entity:"containertask" awsAPI:"ecs"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             ecsiface.ECSAPI
	Name            *string   `templateName:"name"`
	ContainerName   *string   `templateName:"container-name"`
	Image           *string   `templateName:"image"`
	MemoryHardLimit *int64    `templateName:"memory-hard-limit"`
	Commands        []*string `templateName:"command"`
	Env             []*string `templateName:"env"`
	Privileged      *bool     `templateName:"privileged"`
	Workdir         *string   `templateName:"workdir"`
	Ports           []*string `templateName:"ports"`
}

func (cmd *AttachContainertask) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("container-name"), params.Key("image"), params.Key("memory-hard-limit"), params.Key("name"),
		params.Opt("command", "env", "ports", "privileged", "workdir"),
	))
}

func (cmd *AttachContainertask) ManualRun(renv env.Running) (interface{}, error) {
	var taskDefinitionInput *ecs.RegisterTaskDefinitionInput
	taskDefinitionName := StringValue(cmd.Name)

	taskdefOutput, err := cmd.api.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: cmd.Name,
	})
	if awserr, ok := err.(awserr.Error); err != nil && ok {
		if awserr.Code() == "ClientException" && strings.Contains(strings.ToLower(awserr.Message()), "unable to describe task definition") {
			cmd.logger.Verbosef("service %s does not exist: creating service", taskDefinitionName)
			taskDefinitionInput = &ecs.RegisterTaskDefinitionInput{
				Family: aws.String(taskDefinitionName),
			}
		} else {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		taskDefinitionInput = &ecs.RegisterTaskDefinitionInput{
			ContainerDefinitions: taskdefOutput.TaskDefinition.ContainerDefinitions,
			Family:               taskdefOutput.TaskDefinition.Family,
			NetworkMode:          taskdefOutput.TaskDefinition.NetworkMode,
			PlacementConstraints: taskdefOutput.TaskDefinition.PlacementConstraints,
			TaskRoleArn:          taskdefOutput.TaskDefinition.TaskRoleArn,
			Volumes:              taskdefOutput.TaskDefinition.Volumes,
		}
	}

	container := &ecs.ContainerDefinition{}
	if err = setFieldWithType(cmd.ContainerName, container, "Name", awsstr); err != nil {
		return nil, err
	}
	if err = setFieldWithType(cmd.Image, container, "Image", awsstr); err != nil {
		return nil, err
	}
	if err = setFieldWithType(cmd.MemoryHardLimit, container, "Memory", awsint64); err != nil {
		return nil, err
	}
	if cmd.Commands != nil {
		switch len(cmd.Commands) {
		case 1:
			if err = setFieldWithType(strings.Split(StringValue(cmd.Commands[0]), " "), container, "Command", awsstringslice); err != nil {
				return nil, err
			}
		default:
			if err = setFieldWithType(cmd.Commands, container, "Command", awsstringslice); err != nil {
				return nil, err
			}
		}
	}
	if len(cmd.Env) > 0 {
		if err = setFieldWithType(cmd.Env, container, "Environment", awsecskeyvalue); err != nil {
			return nil, err
		}
	}
	if BoolValue(cmd.Privileged) {
		if err = setFieldWithType(true, container, "Privileged", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.Workdir != nil {
		if err = setFieldWithType(cmd.Workdir, container, "WorkingDirectory", awsstr); err != nil {
			return nil, err
		}
	}
	if len(cmd.Ports) > 0 {
		if err = setFieldWithType(cmd.Ports, container, "PortMappings", awsportmappings); err != nil {
			return nil, err
		}
	}

	taskDefinitionInput.ContainerDefinitions = append(taskDefinitionInput.ContainerDefinitions, container)

	start := time.Now()

	taskDefOutput, err := cmd.api.RegisterTaskDefinition(taskDefinitionInput)
	if err != nil {
		return nil, err
	}
	cmd.logger.ExtraVerbosef("ecs.RegisterTaskDefinitionOutput call took %s", time.Since(start))
	return taskDefOutput, nil
}

func (cmd *AttachContainertask) ExtractResult(i interface{}) string {
	return StringValue(i.(*ecs.RegisterTaskDefinitionOutput).TaskDefinition.TaskDefinitionArn)
}

type DetachContainertask struct {
	_             string `action:"detach" entity:"containertask" awsAPI:"ecs"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           ecsiface.ECSAPI
	Name          *string `templateName:"name"`
	ContainerName *string `templateName:"container-name"`
}

func (cmd *DetachContainertask) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("container-name"), params.Key("name")))
}

func (cmd *DetachContainertask) ManualRun(renv env.Running) (interface{}, error) {
	taskdefOutput, err := cmd.api.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: cmd.Name,
	})
	if err != nil {
		return nil, err
	}

	var containerDefinitions []*ecs.ContainerDefinition
	var found bool
	var containerNames []string
	for _, def := range taskdefOutput.TaskDefinition.ContainerDefinitions {
		name := aws.StringValue(def.Name)
		containerNames = append(containerNames, name)
		if name == StringValue(cmd.ContainerName) || aws.StringValue(def.Image) == StringValue(cmd.ContainerName) {
			found = true
		} else {
			containerDefinitions = append(containerDefinitions, def)
		}
	}
	if !found {
		return nil, fmt.Errorf("did not find any container called '%s': found: '%s'", StringValue(cmd.ContainerName), strings.Join(containerNames, "','"))
	}

	if len(containerDefinitions) > 0 { //At least one container remaining
		taskDefinitionInput := &ecs.RegisterTaskDefinitionInput{
			ContainerDefinitions: containerDefinitions,
			Family:               taskdefOutput.TaskDefinition.Family,
			NetworkMode:          taskdefOutput.TaskDefinition.NetworkMode,
			PlacementConstraints: taskdefOutput.TaskDefinition.PlacementConstraints,
			TaskRoleArn:          taskdefOutput.TaskDefinition.TaskRoleArn,
			Volumes:              taskdefOutput.TaskDefinition.Volumes,
		}
		start := time.Now()

		if _, err := cmd.api.RegisterTaskDefinition(taskDefinitionInput); err != nil {
			return nil, err
		}
		cmd.logger.ExtraVerbosef("ecs.RegisterTaskDefinition call took %s", time.Since(start))

	} else {
		cmd.logger.Verbosef("no container remaining in service %s: deleting service", StringValue(cmd.Name))
		taskDefinitionInput := &ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: taskdefOutput.TaskDefinition.TaskDefinitionArn,
		}
		start := time.Now()

		if _, err := cmd.api.DeregisterTaskDefinition(taskDefinitionInput); err != nil {
			return nil, err
		}
		cmd.logger.ExtraVerbosef("ecs.DeregisterTaskDefinition call took %s", time.Since(start))
	}

	return taskdefOutput, nil
}

type DeleteContainertask struct {
	_           string `action:"delete" entity:"containertask" awsAPI:"ecs" awsDryRun:"manual"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         ecsiface.ECSAPI
	Name        *string `templateName:"name"`
	AllVersions *bool   `templateName:"all-versions"`
}

func (cmd *DeleteContainertask) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"),
		params.Opt("all-versions"),
	))
}

func (cmd *DeleteContainertask) dryRun(renv env.Running, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}
	taskDefinitionName := StringValue(cmd.Name)

	taskDefOutput, err := cmd.api.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
		FamilyPrefix: cmd.Name,
	})
	if err != nil {
		return nil, err
	}
	switch len(taskDefOutput.TaskDefinitionArns) {
	case 0:
		return nil, fmt.Errorf("no containertask found with name '%s'", taskDefinitionName)
	case 1:
		cmd.logger.Verbosef("only one version found for containertask '%s', will delete '%s'.", taskDefinitionName, aws.StringValue(taskDefOutput.TaskDefinitionArns[0]))
	default:
		if BoolValue(cmd.AllVersions) {
			cmd.logger.Warningf("multiple versions found for containertask '%s'. Will delete '%s'", taskDefinitionName, strings.Join(aws.StringValueSlice(taskDefOutput.TaskDefinitionArns), "','"))
		} else {
			cmd.logger.Infof("multiple versions found for containertask '%s'", taskDefinitionName)
			cmd.logger.Warningf("will delete only latest version: '%s'. Add param `all-versions=true` to delete all versions", aws.StringValue(taskDefOutput.TaskDefinitionArns[len(taskDefOutput.TaskDefinitionArns)-1]))
		}
	}
	return nil, nil
}

func (cmd *DeleteContainertask) ManualRun(renv env.Running) (interface{}, error) {
	taskDefinitionName := StringValue(cmd.Name)

	if BoolValue(cmd.AllVersions) {
		taskDefOutput, err := cmd.api.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
			FamilyPrefix: aws.String(taskDefinitionName),
		})
		if err != nil {
			return nil, err
		}
		for _, task := range taskDefOutput.TaskDefinitionArns {
			cmd.logger.ExtraVerbosef("deleting '%s'", aws.StringValue(task))
			start := time.Now()
			if _, err := cmd.api.DeregisterTaskDefinition(&ecs.DeregisterTaskDefinitionInput{TaskDefinition: task}); err != nil {
				return nil, fmt.Errorf("deregister task definition: %s", err)
			}
			cmd.logger.ExtraVerbosef("ecs.DeregisterTaskDefinition call took %s", time.Since(start))
		}
	} else {
		taskDefOutput, err := cmd.api.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String(taskDefinitionName),
		})
		if err != nil {
			return nil, err
		}
		cmd.logger.ExtraVerbosef("deleting '%s'", aws.StringValue(taskDefOutput.TaskDefinition.TaskDefinitionArn))
		start := time.Now()
		if _, err := cmd.api.DeregisterTaskDefinition(&ecs.DeregisterTaskDefinitionInput{TaskDefinition: taskDefOutput.TaskDefinition.TaskDefinitionArn}); err != nil {
			return nil, fmt.Errorf("deregister task definition: %s", err)
		}
		cmd.logger.ExtraVerbosef("ecs.DeregisterTaskDefinition call took %s", time.Since(start))
	}
	return nil, nil
}
