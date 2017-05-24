/*
Copyright 2017 WALLIX

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

package awsdriver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mitchellh/ioprogress"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/logger"
)

const (
	notFoundState = "not-found"
)

func (d *Ec2Driver) Attach_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("attach securitygroup: missing required params 'id'")
	}

	_, hasInstance := params["instance"]

	if !hasInstance {
		return nil, errors.New("attach securitygroup: missing 'instance' param")
	}

	d.logger.Verbose("params dry run: attach securitygroup ok")
	return nil, nil
}

func (d *Ec2Driver) Attach_Securitygroup(params map[string]interface{}) (interface{}, error) {
	instance, hasInstance := params["instance"].(string)

	switch {
	case hasInstance:
		groups, err := d.fetchInstanceSecurityGroups(instance)
		if err != nil {
			return nil, fmt.Errorf("fetching securitygroups for instance %s: %s", instance, err)
		}

		groups = append(groups, fmt.Sprint(params["id"]))
		if len(groups) == 0 {
			d.logger.Errorf("AWS instances must have at least one securitygroup")
		}
		call := &driverCall{
			d:      d,
			fn:     d.ModifyInstanceAttribute,
			logger: d.logger,
			setters: []setter{
				{val: instance, fieldPath: "InstanceID", fieldType: awsstr},
				{val: groups, fieldPath: "Groups", fieldType: awsstringslice},
			},
			desc: "attach securitygroup",
		}
		return call.execute(&ec2.ModifyInstanceAttributeInput{})
	}

	return nil, errors.New("missing 'instance' param")
}

func (d *Ec2Driver) Detach_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("detach securitygroup: missing required params 'id'")
	}

	_, hasInstance := params["instance"]

	if !hasInstance {
		return nil, errors.New("detach securitygroup: missing 'instance' param")
	}

	d.logger.Verbose("params dry run: detach securitygroup ok")
	return nil, nil
}

func (d *Ec2Driver) Detach_Securitygroup(params map[string]interface{}) (interface{}, error) {
	instance, hasInstance := params["instance"].(string)

	switch {
	case hasInstance:
		groups, err := d.fetchInstanceSecurityGroups(instance)
		if err != nil {
			return nil, fmt.Errorf("fetching securitygroups for instance %s: %s", instance, err)
		}

		cleaned := removeString(groups, fmt.Sprint(params["id"]))

		if len(cleaned) == 0 {
			d.logger.Errorf("AWS instances must have at least one securitygroup")
		}
		call := &driverCall{
			d:      d,
			fn:     d.ModifyInstanceAttribute,
			logger: d.logger,
			setters: []setter{
				{val: instance, fieldPath: "InstanceID", fieldType: awsstr},
				{val: cleaned, fieldPath: "Groups", fieldType: awsstringslice},
			},
			desc: "detach securitygroup",
		}
		return call.execute(&ec2.ModifyInstanceAttributeInput{})
	}

	return nil, errors.New("missing 'instance' param")
}

func (d *Ec2Driver) fetchInstanceSecurityGroups(id string) ([]string, error) {
	params := &ec2.DescribeInstanceAttributeInput{
		Attribute:  aws.String("groupSet"),
		InstanceId: aws.String(id),
	}
	resp, err := d.DescribeInstanceAttribute(params)
	if err != nil {
		return nil, err
	}

	var groups []string
	for _, g := range resp.Groups {
		groups = append(groups, aws.StringValue(g.GroupId))
	}

	return groups, nil
}

func removeString(arr []string, s string) (out []string) {
	for _, e := range arr {
		if e != s {
			out = append(out, e)
		}
	}
	return
}

func (d *IamDriver) Create_Policy_DryRun(params map[string]interface{}) (interface{}, error) {
	_, effect := params["effect"]
	_, action := params["action"]
	_, resource := params["resource"]

	if !effect && !action && !resource {
		return nil, errors.New("create role: missing policy effect, action and resource values")
	}

	d.logger.Verbose("params dry run: create policy ok")
	return nil, nil
}

func (d *IamDriver) Create_Policy(params map[string]interface{}) (interface{}, error) {
	effect, _ := params["effect"].(string)
	resource, _ := params["resource"].(string)
	actions, multipleAction := params["action"].([]string)
	action, singleAction := params["action"].(string)

	if resource == "all" {
		resource = "*"
	}

	stat := policyStatement{Effect: strings.Title(effect), Resource: resource}

	if multipleAction {
		stat.Actions = actions
	}
	if singleAction {
		stat.Actions = []string{action}
	}

	policy := &policyBody{
		Version:   "2012-10-17",
		Statement: []policyStatement{stat},
	}

	b, err := json.MarshalIndent(policy, "", " ")
	if err != nil {
		return nil, errors.New("cannot marshal policy document")
	}

	d.logger.ExtraVerbosef("policy document json:\n%s\n", string(b))

	call := &driverCall{
		d:      d,
		desc:   "create policy",
		fn:     d.CreatePolicy,
		logger: d.logger,
		setters: []setter{
			{val: params["name"], fieldPath: "PolicyName", fieldType: awsstr},
			{val: string(b), fieldPath: "PolicyDocument", fieldType: awsstr},
		},
	}

	output, err := call.execute(&iam.CreatePolicyInput{})
	if err != nil {
		return nil, err
	}

	return aws.StringValue(output.(*iam.CreatePolicyOutput).Policy.Arn), nil
}

func (d *IamDriver) Delete_Role_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete role: missing required params 'name'")
	}

	d.logger.Verbose("params dry run: delete role ok")
	return nil, nil
}

func (d *IamDriver) Delete_Role(params map[string]interface{}) (interface{}, error) {
	d.Detach_Role(map[string]interface{}{
		"name":            params["name"],
		"instanceprofile": params["name"],
	})
	d.Delete_Instanceprofile(params)

	input := &iam.DeleteRoleInput{}

	if err := setFieldWithType(params["name"], input, "RoleName", awsstr); err != nil {
		return nil, err
	}

	start := time.Now()
	output, err := d.DeleteRole(input)
	if err != nil {
		return nil, fmt.Errorf("delete role: %s", err)
	}
	d.logger.ExtraVerbosef("iam.DeleteRole call took %s", time.Since(start))
	d.logger.Info("delete role done")

	return output, nil
}

func (d *IamDriver) Create_Role_DryRun(params map[string]interface{}) (interface{}, error) {
	_, pAccount := params["principal-account"]
	_, pService := params["principal-service"]
	_, pUser := params["principal-user"]

	if !pAccount && !pService && !pUser {
		return nil, errors.New("create role: missing principal (either a user, service or account)")
	}

	d.logger.Verbose("params dry run: create role ok")
	return nil, nil
}

type principal struct {
	AWS     interface{} `json:",omitempty"`
	Service interface{} `json:",omitempty"`
}

type policyStatement struct {
	Effect    string     `json:",omitempty"`
	Actions   []string   `json:"Action,omitempty"`
	Resource  string     `json:",omitempty"`
	Principal *principal `json:",omitempty"`
}

type policyBody struct {
	Version   string
	Statement []policyStatement
}

func (d *IamDriver) Create_Role(params map[string]interface{}) (interface{}, error) {
	pAccount, _ := params["principal-account"]
	pService, _ := params["principal-service"]
	pUser, _ := params["principal-user"]

	princ := new(principal)
	if pAccount != nil {
		princ.AWS = pAccount
	} else if pUser != nil {
		princ.AWS = pUser
	} else if pService != nil {
		princ.Service = pService
	}

	trust := &policyBody{
		Version:   "2012-10-17",
		Statement: []policyStatement{{Effect: "Allow", Actions: []string{"sts:AssumeRole"}, Principal: princ}},
	}

	b, err := json.MarshalIndent(trust, "", " ")
	if err != nil {
		return nil, errors.New("cannot marshal role trust policy document")
	}

	d.logger.ExtraVerbosef("role trust policy document json:\n%s\n", string(b))

	call := &driverCall{
		d:      d,
		desc:   "create role",
		fn:     d.CreateRole,
		logger: d.logger,
		setters: []setter{
			{val: params["name"], fieldPath: "RoleName", fieldType: awsstr},
			{val: string(b), fieldPath: "AssumeRolePolicyDocument", fieldType: awsstr},
		},
	}

	output, err := call.execute(&iam.CreateRoleInput{})
	if err != nil {
		return nil, err
	}
	role := output.(*iam.CreateRoleOutput).Role
	roleName := aws.StringValue(role.RoleName)

	d.Create_Instanceprofile(params)
	d.Attach_Role(map[string]interface{}{
		"name":            roleName,
		"instanceprofile": roleName,
	})
	if secs, ok := params["sleep-after"].(int); ok {
		d.logger.Infof("sleeping for %d seconds", secs)
		time.Sleep(time.Duration(secs) * time.Second)
	}

	return aws.StringValue(role.Arn), nil
}

func (d *IamDriver) Attach_Policy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("attach policy: missing required params 'arn'")
	}

	_, hasUser := params["user"]
	_, hasGroup := params["group"]
	_, hasRole := params["role"]

	if !hasUser && !hasGroup && !hasRole {
		return nil, errors.New("attach policy: missing one of 'user, group, role' param")
	}

	d.logger.Verbose("params dry run: attach policy ok")
	return nil, nil
}

func (d *IamDriver) Attach_Policy(params map[string]interface{}) (interface{}, error) {
	user, hasUser := params["user"]
	group, hasGroup := params["group"]
	role, hasRole := params["role"]

	call := &driverCall{
		d:      d,
		logger: d.logger,
		setters: []setter{
			{val: params["arn"], fieldPath: "PolicyArn", fieldType: awsstr},
		},
	}

	switch {
	case hasUser:
		call.desc = "attach policy to user"
		call.fn = d.AttachUserPolicy
		call.setters = append(call.setters, setter{val: user, fieldPath: "UserName", fieldType: awsstr})
		return call.execute(&iam.AttachUserPolicyInput{})
	case hasGroup:
		call.desc = "attach policy to group"
		call.fn = d.AttachGroupPolicy
		call.setters = append(call.setters, setter{val: group, fieldPath: "GroupName", fieldType: awsstr})
		return call.execute(&iam.AttachGroupPolicyInput{})
	case hasRole:
		call.desc = "attach policy to role"
		call.fn = d.AttachRolePolicy
		call.setters = append(call.setters, setter{val: role, fieldPath: "RoleName", fieldType: awsstr})
		return call.execute(&iam.AttachRolePolicyInput{})
	}

	return nil, errors.New("missing one of 'user, group, role' param")
}

func (d *IamDriver) Detach_Policy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("detach policy: missing required params 'arn'")
	}

	_, hasUser := params["user"]
	_, hasGroup := params["group"]
	_, hasRole := params["role"]

	if !hasUser && !hasGroup && !hasRole {
		return nil, errors.New("detach policy: missing one of 'user, group, role' param")
	}

	d.logger.Verbose("params dry run: detach policy ok")
	return nil, nil
}

func (d *IamDriver) Detach_Policy(params map[string]interface{}) (interface{}, error) {
	user, hasUser := params["user"]
	group, hasGroup := params["group"]
	role, hasRole := params["role"]

	call := &driverCall{
		d:      d,
		logger: d.logger,
		setters: []setter{
			{val: params["arn"], fieldPath: "PolicyArn", fieldType: awsstr},
		},
	}

	switch {
	case hasUser:
		call.desc = "detach policy from user"
		call.fn = d.DetachUserPolicy
		call.setters = append(call.setters, setter{val: user, fieldPath: "UserName", fieldType: awsstr})
		return call.execute(&iam.DetachUserPolicyInput{})
	case hasGroup:
		call.desc = "detach policy from group"
		call.fn = d.DetachGroupPolicy
		call.setters = append(call.setters, setter{val: group, fieldPath: "GroupName", fieldType: awsstr})
		return call.execute(&iam.DetachGroupPolicyInput{})
	case hasRole:
		call.desc = "detach policy from role"
		call.fn = d.DetachRolePolicy
		call.setters = append(call.setters, setter{val: role, fieldPath: "RoleName", fieldType: awsstr})
		return call.execute(&iam.DetachRolePolicyInput{})
	}

	return nil, errors.New("missing one of 'user, group, role' param")
}

func (d *IamDriver) Create_Accesskey_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["user"]; !ok {
		return nil, errors.New("create accesskey: missing required params 'user'")
	}

	d.logger.Verbose("params dry run: create accesskey ok")
	return nil, nil
}

func (d *IamDriver) Create_Accesskey(params map[string]interface{}) (interface{}, error) {
	input := &iam.CreateAccessKeyInput{}
	var err error

	// Required params
	err = setFieldWithType(params["user"], input, "UserName", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *iam.CreateAccessKeyOutput
	output, err = d.CreateAccessKey(input)

	if err != nil {
		return nil, fmt.Errorf("create accesskey: %s", err)
	}
	d.logger.ExtraVerbosef("iam.CreateAccessKey call took %s", time.Since(start))

	d.logger.Infof("Access key created. Here are the crendentials for user %s:", aws.StringValue(output.AccessKey.UserName))
	fmt.Println()
	fmt.Println(strings.Repeat("*", 64))
	fmt.Printf("aws_access_key_id = %s\n", aws.StringValue(output.AccessKey.AccessKeyId))
	fmt.Printf("aws_secret_access_key = %s\n", aws.StringValue(output.AccessKey.SecretAccessKey))
	fmt.Println(strings.Repeat("*", 64))
	fmt.Println()
	d.logger.Warning("This is your only opportunity to view the secret access keys.")
	d.logger.Warning("Save the user's new access key ID and secret access key in a safe and secure place.")
	d.logger.Warning("You will not have access to the secret keys again after this step.")

	id := aws.StringValue(output.AccessKey.AccessKeyId)

	d.logger.Infof("create accesskey '%s' done", id)
	return id, nil
}

func (d *Ec2Driver) Check_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("check instance: missing required params 'id'")
	}

	states := map[string]struct{}{
		"pending":       {},
		"running":       {},
		"shutting-down": {},
		"terminated":    {},
		"stopping":      {},
		"stopped":       {},
		notFoundState:   {},
	}

	if state, ok := params["state"].(string); !ok {
		return nil, errors.New("check instance: missing required params 'state'")
	} else {
		if _, stok := states[state]; !stok {
			return nil, fmt.Errorf("check instance: invalid state '%s'", state)
		}
	}

	if _, ok := params["timeout"]; !ok {
		return nil, errors.New("check instance: missing required params 'timeout'")
	}

	if _, ok := params["timeout"].(int); !ok {
		return nil, errors.New("check instance: timeout param is not int")
	}

	input := &ec2.DescribeInstancesInput{}
	input.DryRun = aws.Bool(true)

	// Required params
	err := setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	_, err = d.DescribeInstances(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("instance")
			d.logger.Verbose("dry run: check instance ok")
			return id, nil
		}
	}
	return nil, fmt.Errorf("dry run: check instance: %s", err)
}

func (d *Ec2Driver) Check_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DescribeInstancesInput{}

	// Required params
	err := setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}
	c := &checker{
		description: fmt.Sprintf("instance %s", params["id"]),
		timeout:     time.Duration(params["timeout"].(int)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := d.DescribeInstances(input)
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
							if aws.StringValue(inst.InstanceId) == params["id"] {
								return aws.StringValue(inst.State.Name), nil
							}
						}
					}
				}
			}
			return notFoundState, nil
		},
		expect: fmt.Sprint(params["state"]),
		logger: d.logger,
	}
	return nil, c.check()
}

func (d *Ec2Driver) Check_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("check securitygroup: missing required params 'id'")
	}

	states := map[string]struct{}{
		"unused": {},
	}

	if state, ok := params["state"].(string); !ok {
		return nil, errors.New("check securitygroup: missing required params 'state'")
	} else {
		if _, stok := states[state]; !stok {
			return nil, fmt.Errorf("check securitygroup: invalid state '%s'", state)
		}
	}

	if _, ok := params["timeout"]; !ok {
		return nil, errors.New("check securitygroup: missing required params 'timeout'")
	}

	if _, ok := params["timeout"].(int); !ok {
		return nil, errors.New("check securitygroup: timeout param is not int")
	}
	d.logger.Verbose("dry run: check instance ok")
	return nil, nil
}

func (d *Ec2Driver) Check_Securitygroup(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{Name: aws.String("group-id"), Values: []*string{aws.String(fmt.Sprint(params["id"]))}},
		},
	}

	c := &checker{
		description: fmt.Sprintf("securitygroup %s", params["id"]),
		timeout:     time.Duration(params["timeout"].(int)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := d.DescribeNetworkInterfaces(input)
			if err != nil {
				return "", err
			}
			if len(output.NetworkInterfaces) == 0 {
				return "unused", nil
			}
			var niIds []string
			for _, ni := range output.NetworkInterfaces {
				niIds = append(niIds, aws.StringValue(ni.NetworkInterfaceId))
			}
			return fmt.Sprintf("used by %s", strings.Join(niIds, ", ")), nil
		},
		expect: fmt.Sprint(params["state"]),
		logger: d.logger,
	}
	return nil, c.check()
}

func (d *Elbv2Driver) Check_Loadbalancer_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["id"]; !ok {
		return nil, errors.New("check loadbalancer: missing required params 'id'")
	}

	states := map[string]struct{}{
		"provisioning": {},
		"active":       {},
		"failed":       {},
		notFoundState:  {},
	}

	if state, ok := params["state"].(string); !ok {
		return nil, errors.New("check loadbalancer: missing required params 'state'")
	} else {
		if _, stok := states[state]; !stok {
			return nil, fmt.Errorf("check loadbalancer: invalid state '%s'", state)
		}
	}

	if _, ok := params["timeout"]; !ok {
		return nil, errors.New("check loadbalancer: missing required params 'timeout'")
	}

	d.logger.Verbose("params dry run: check loadbalancer ok")
	return nil, nil
}

func (d *Elbv2Driver) Check_Loadbalancer(params map[string]interface{}) (interface{}, error) {
	input := &elbv2.DescribeLoadBalancersInput{}

	// Required params
	err := setFieldWithType(params["id"], input, "LoadBalancerArns", awsstringslice)
	if err != nil {
		return nil, err
	}
	c := &checker{
		description: fmt.Sprintf("loadbalancer %s", params["id"]),
		timeout:     time.Duration(params["timeout"].(int)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := d.DescribeLoadBalancers(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "LoadBalancerNotFound" {
						return notFoundState, nil
					}
				} else {
					return "", err
				}
			} else {
				for _, lb := range output.LoadBalancers {
					if aws.StringValue(lb.LoadBalancerArn) == params["id"] {
						return aws.StringValue(lb.State.Code), nil
					}
				}
			}
			return notFoundState, nil
		},
		expect: fmt.Sprint(params["state"]),
		logger: d.logger,
	}
	return nil, c.check()
}

func (d *AutoscalingDriver) Check_Scalinggroup_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"].(string); !ok {
		return nil, errors.New("check scalinggroup: missing required params 'name'")
	}

	if _, ok := params["count"]; !ok {
		return nil, errors.New("check scalinggroup: missing required param 'count'")
	}

	if _, ok := params["timeout"].(int); !ok {
		return nil, errors.New("check scalinggroup: missing required int param 'timeout'")
	}

	d.logger.Verbose("params dry run: check scalinggroup ok")
	return nil, nil
}

func (d *AutoscalingDriver) Check_Scalinggroup(params map[string]interface{}) (interface{}, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{}

	// Required params
	err := setFieldWithType(params["name"], input, "AutoScalingGroupNames", awsstringslice)
	if err != nil {
		return nil, err
	}
	c := &checker{
		description: fmt.Sprintf("scalinggroup '%s'", params["name"]),
		timeout:     time.Duration(params["timeout"].(int)) * time.Second,
		frequency:   5 * time.Second,
		checkName:   "count",
		fetchFunc: func() (string, error) {
			output, err := d.DescribeAutoScalingGroups(input)
			if err != nil {
				return "", err
			}
			for _, group := range output.AutoScalingGroups {
				if aws.StringValue(group.AutoScalingGroupName) == params["name"] {
					count := len(group.Instances) + len(group.LoadBalancerNames) + len(group.TargetGroupARNs)
					return fmt.Sprint(count), nil
				}
			}
			return "", fmt.Errorf("scalinggroup %s not found", params["name"])
		},
		expect: fmt.Sprint(params["count"]),
		logger: d.logger,
	}
	return nil, c.check()
}

func (d *Ec2Driver) Create_Tag_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateTagsInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["resource"], input, "Resources", awsstringslice)
	if err != nil {
		return nil, err
	}
	input.Tags = []*ec2.Tag{{Key: aws.String(fmt.Sprint(params["key"])), Value: aws.String(fmt.Sprint(params["value"]))}}

	_, err = d.CreateTags(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("tag")
			d.logger.Verbosef("dry run: create tag '%s=%s' ok", params["key"], params["value"])
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: create tag: %s", err)
}

func (d *Ec2Driver) Create_Tag(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateTagsInput{}
	var err error

	// Required params
	err = setFieldWithType(params["resource"], input, "Resources", awsstringslice)
	if err != nil {
		return nil, err
	}
	input.Tags = []*ec2.Tag{{Key: aws.String(fmt.Sprint(params["key"])), Value: aws.String(fmt.Sprint(params["value"]))}}

	start := time.Now()
	var output *ec2.CreateTagsOutput
	output, err = d.CreateTags(input)
	if err != nil {
		return nil, fmt.Errorf("create tag: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.CreateTags call took %s", time.Since(start))
	d.logger.Infof("create tag '%s=%s' on '%s' done", params["key"], params["value"], params["resource"])
	return output, nil
}

func (d *Ec2Driver) Delete_Tag_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteTagsInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["resource"], input, "Resources", awsstringslice)
	if err != nil {
		return nil, err
	}
	input.Tags = []*ec2.Tag{{Key: aws.String(fmt.Sprint(params["key"])), Value: aws.String(fmt.Sprint(params["value"]))}}

	_, err = d.DeleteTags(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("tag")
			d.logger.Verbosef("dry run: delete tag '%s=%s' ok", params["key"], params["value"])
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete tag: %s", err)
}

func (d *Ec2Driver) Delete_Tag(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeleteTagsInput{}
	var err error

	// Required params
	err = setFieldWithType(params["resource"], input, "Resources", awsstringslice)
	if err != nil {
		return nil, err
	}
	input.Tags = []*ec2.Tag{{Key: aws.String(fmt.Sprint(params["key"])), Value: aws.String(fmt.Sprint(params["value"]))}}

	start := time.Now()
	var output *ec2.DeleteTagsOutput
	output, err = d.DeleteTags(input)
	if err != nil {
		return nil, fmt.Errorf("delete tag: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeleteTags call took %s", time.Since(start))
	d.logger.Infof("delete tag '%s=%s' on '%s' done", params["key"], params["value"], params["resource"])
	return output, nil
}

func (d *Ec2Driver) Create_Keypair_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ImportKeyPairInput{}

	input.DryRun = aws.Bool(true)
	err := setFieldWithType(params["name"], input, "KeyName", awsstr)
	if err != nil {
		return nil, err
	}

	if params["name"] == "" {
		return nil, fmt.Errorf("dry run: saving private key: empty 'name' parameter")
	}

	const keyDirEnv = "__AWLESS_KEYS_DIR"
	keyDir := os.Getenv(keyDirEnv)
	if keyDir == "" {
		return nil, fmt.Errorf("dry run: saving private key: empty env var '%s'", keyDirEnv)
	}

	privKeyPath := filepath.Join(keyDir, fmt.Sprint(params["name"])+".pem")
	_, err = os.Stat(privKeyPath)
	if err == nil {
		return nil, fmt.Errorf("dry run: saving private key: file already exists at path: %s", privKeyPath)
	}

	return nil, nil
}

func (d *Ec2Driver) Create_Keypair(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ImportKeyPairInput{}
	err := setFieldWithType(params["name"], input, "KeyName", awsstr)
	if err != nil {
		return nil, err
	}
	var encrypted bool
	var encryptedMsg string
	if enc, _ := params["encrypted"].(string); enc == "true" {
		encrypted = true
		encryptedMsg = " encrypted"
	}

	d.logger.Infof("Generating locally a%s RSA 4096 bits keypair...", encryptedMsg)
	pub, priv, err := console.GenerateSSHKeyPair(4096, encrypted)
	if err != nil {
		return nil, fmt.Errorf("generating key: %s", err)
	}
	privKeyPath := filepath.Join(os.Getenv("__AWLESS_KEYS_DIR"), fmt.Sprint(params["name"])+".pem")
	_, err = os.Stat(privKeyPath)
	if err == nil {
		return nil, fmt.Errorf("saving private key: file already exists at path: %s", privKeyPath)
	}
	err = ioutil.WriteFile(privKeyPath, priv, 0400)
	if err != nil {
		return nil, fmt.Errorf("saving private key: %s", err)
	}
	d.logger.Infof("4096 RSA keypair generated locally and stored%s in '%s'", encryptedMsg, privKeyPath)
	input.PublicKeyMaterial = pub

	output, err := d.ImportKeyPair(input)
	if err != nil {
		return nil, fmt.Errorf("create key: %s", err)
	}
	id := aws.StringValue(output.KeyName)
	d.logger.Infof("create keypair '%s' done", id)
	return aws.StringValue(output.KeyName), nil
}

func (d *Ec2Driver) Update_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	ipPerms, err := buildIpPermissionsFromParams(params)
	if err != nil {
		return nil, err
	}
	var input interface{}
	if action, ok := params["inbound"].(string); ok {
		switch action {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupIngressInput{DryRun: aws.Bool(true), IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupIngressInput{DryRun: aws.Bool(true), IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'inbound' parameter expect 'authorize' or 'revoke', got %s", action)
		}
	}
	if action, ok := params["outbound"].(string); ok {
		switch action {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupEgressInput{DryRun: aws.Bool(true), IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupEgressInput{DryRun: aws.Bool(true), IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'outbound' parameter expect 'authorize' or 'revoke', got %s", action)
		}
	}
	if input == nil {
		return nil, fmt.Errorf("expect either 'inbound' or 'outbound' parameter")
	}

	// Required params
	err = setFieldWithType(params["id"], input, "GroupId", awsstr)
	if err != nil {
		return nil, err
	}

	switch ii := input.(type) {
	case *ec2.AuthorizeSecurityGroupIngressInput:
		_, err = d.AuthorizeSecurityGroupIngress(ii)
	case *ec2.RevokeSecurityGroupIngressInput:
		_, err = d.RevokeSecurityGroupIngress(ii)
	case *ec2.AuthorizeSecurityGroupEgressInput:
		_, err = d.AuthorizeSecurityGroupEgress(ii)
	case *ec2.RevokeSecurityGroupEgressInput:
		_, err = d.RevokeSecurityGroupEgress(ii)
	}
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			d.logger.Verbose("dry run: update securitygroup ok")
			return nil, nil
		}
	}
	return nil, fmt.Errorf("dry run: update securitygroup: %s", err)
}

func (d *Ec2Driver) Update_Securitygroup(params map[string]interface{}) (interface{}, error) {
	ipPerms, err := buildIpPermissionsFromParams(params)
	if err != nil {
		return nil, err
	}
	var input interface{}
	if action, ok := params["inbound"].(string); ok {
		switch action {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupIngressInput{IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupIngressInput{IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'inbound' parameter expect 'authorize' or 'revoke', got %s", action)
		}
	}
	if action, ok := params["outbound"].(string); ok {
		switch action {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupEgressInput{IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupEgressInput{IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'outbound' parameter expect 'authorize' or 'revoke', got %s", action)
		}
	}
	if input == nil {
		return nil, fmt.Errorf("expect either 'inbound' or 'outbound' parameter")
	}

	// Required params
	err = setFieldWithType(params["id"], input, "GroupId", awsstr)
	if err != nil {
		return nil, err
	}

	var output interface{}
	switch ii := input.(type) {
	case *ec2.AuthorizeSecurityGroupIngressInput:
		output, err = d.AuthorizeSecurityGroupIngress(ii)
	case *ec2.RevokeSecurityGroupIngressInput:
		output, err = d.RevokeSecurityGroupIngress(ii)
	case *ec2.AuthorizeSecurityGroupEgressInput:
		output, err = d.AuthorizeSecurityGroupEgress(ii)
	case *ec2.RevokeSecurityGroupEgressInput:
		output, err = d.RevokeSecurityGroupEgress(ii)
	}
	if err != nil {
		return nil, fmt.Errorf("update securitygroup: %s", err)
	}

	d.logger.Info("update securitygroup done")
	return output, nil
}

func (d *S3Driver) Create_S3object_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["bucket"]; !ok {
		return nil, errors.New("create s3object: missing required params 'bucket'")
	}

	if _, ok := params["file"].(string); !ok {
		return nil, errors.New("create s3object: missing required string params 'file'")
	}

	stat, err := os.Stat(params["file"].(string))
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot find file '%s'", params["file"])
	}
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory", params["file"])
	}

	d.logger.Verbose("params dry run: create s3object ok")
	return nil, nil
}

type progressReadSeeker struct {
	file   *os.File
	reader *ioprogress.Reader
}

func newProgressReader(f *os.File) (*progressReadSeeker, error) {
	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	draw := func(progress, total int64) string {
		// &s3.PutObjectInput.Body will be read twice
		// once in memory and a second time for the HTTP upload
		// here we only display for the actual HTTP upload
		if progress > total {
			return ioprogress.DrawTextFormatBytes(progress/2, total)
		}
		return ""
	}

	reader := &ioprogress.Reader{
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, draw),
		Reader:   f,
		Size:     finfo.Size(),
	}

	return &progressReadSeeker{file: f, reader: reader}, nil
}

func (pr *progressReadSeeker) Read(p []byte) (int, error) {
	return pr.reader.Read(p)
}

func (pr *progressReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return pr.file.Seek(offset, whence)
}

func (d *S3Driver) Create_S3object(params map[string]interface{}) (interface{}, error) {
	input := &s3.PutObjectInput{}

	f, err := os.Open(params["file"].(string))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	progressR, err := newProgressReader(f)
	if err != nil {
		return nil, err
	}
	input.Body = progressR

	var fileName string
	if n, ok := params["name"].(string); ok && n != "" {
		fileName = n
	} else {
		_, fileName = filepath.Split(f.Name())
	}
	input.Key = aws.String(fileName)

	fileExt := filepath.Ext(f.Name())
	if mimeType := mime.TypeByExtension(fileExt); mimeType != "" {
		d.logger.ExtraVerbosef("setting object content-type to '%s'", mimeType)
		input.ContentType = aws.String(mimeType)
	}

	// Required params
	err = setFieldWithType(params["bucket"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}
	// Extra params
	if _, ok := params["acl"]; ok {
		err = setFieldWithType(params["acl"], input, "ACL", awsstr)
		if err != nil {
			return nil, err
		}
	}

	d.logger.Infof("uploading '%s'", fileName)

	_, err = d.PutObject(input)
	if err != nil {
		return nil, fmt.Errorf("create s3object: %s", err)
	}

	d.logger.Info("create s3object done")
	return fileName, nil
}

func (d *S3Driver) Update_Bucket_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"]; !ok {
		return nil, errors.New("update bucket: missing required param 'name'")
	}

	_, updatePublicWebsite := params["public-website"]
	if updatePublicWebsite {
		if _, err := strconv.ParseBool(fmt.Sprint(params["public-website"])); err != nil {
			return nil, fmt.Errorf("update bucket: 'public-website' is not a bool: %s", err)
		}
	}
	_, updateAcl := params["acl"]
	if !updatePublicWebsite && !updateAcl {
		return nil, fmt.Errorf("update bucket: must set either 'public-website' or 'acl'")
	}

	d.logger.Verbose("params dry run: update buclet ok")
	return nil, nil
}

func (d *S3Driver) Update_Bucket(params map[string]interface{}) (interface{}, error) {
	bucket := fmt.Sprint(params["name"])

	start := time.Now()

	if _, ok := params["acl"]; ok { // Update the canned ACL to apply to the bucket
		input := &s3.PutBucketAclInput{
			Bucket: aws.String(bucket),
		}
		err := setFieldWithType(params["acl"], input, "ACL", awsstr)
		if err != nil {
			return nil, err
		}
		_, err = d.PutBucketAcl(input)
		if err != nil {
			return nil, fmt.Errorf("update bucket: %s", err)
		}

		d.logger.ExtraVerbosef("s3.PutBucketAcl call took %s", time.Since(start))
		d.logger.Info("update bucket done")
		return nil, nil
	}

	if _, ok := params["public-website"]; ok { // Set/Unset this bucket as a public website
		publicWebsite, err := strconv.ParseBool(fmt.Sprint(params["public-website"]))
		if err != nil {
			return nil, fmt.Errorf("update bucket: 'public-website' is not a bool: %s", err)
		}
		if publicWebsite {
			input := &s3.PutBucketWebsiteInput{
				Bucket:               aws.String(bucket),
				WebsiteConfiguration: &s3.WebsiteConfiguration{},
			}
			if hostname, ok := params["redirect-hostname"].(string); ok {
				input.WebsiteConfiguration.RedirectAllRequestsTo = &s3.RedirectAllRequestsTo{HostName: aws.String(hostname)}
				if enforceHttps, found := params["enforce-https"]; found && strings.ToLower(fmt.Sprint(enforceHttps)) == "true" {
					input.WebsiteConfiguration.RedirectAllRequestsTo.Protocol = aws.String("https")
				}
			} else if index, ok := params["index-suffix"].(string); ok {
				input.WebsiteConfiguration.IndexDocument = &s3.IndexDocument{Suffix: aws.String(index)}
			} else {
				input.WebsiteConfiguration.IndexDocument = &s3.IndexDocument{Suffix: aws.String("index.html")}
			}

			_, err := d.PutBucketWebsite(input)

			if err != nil {
				return nil, fmt.Errorf("update bucket: %s", err)
			}
		} else {
			_, err := d.DeleteBucketWebsite(&s3.DeleteBucketWebsiteInput{
				Bucket: aws.String(bucket),
			})
			if err != nil {
				return nil, fmt.Errorf("update bucket: %s", err)
			}
		}

		d.logger.ExtraVerbosef("s3.PutBucketWebsite call took %s", time.Since(start))
		d.logger.Info("update bucket done")
		return nil, nil
	}
	return nil, nil
}

func (d *Route53Driver) Create_Record_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["zone"]; !ok {
		return nil, errors.New("create record: missing required params 'zone'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("create record: missing required params 'name'")
	}

	if _, ok := params["type"]; !ok {
		return nil, errors.New("create record: missing required params 'type'")
	}

	if _, ok := params["value"]; !ok {
		return nil, errors.New("create record: missing required params 'value'")
	}

	if _, ok := params["ttl"]; !ok {
		return nil, errors.New("create record: missing required params 'ttl'")
	}

	d.logger.Verbose("params dry run: create record ok")
	return nil, nil
}

func (d *Route53Driver) Create_Record(params map[string]interface{}) (interface{}, error) {
	input := &route53.ChangeResourceRecordSetsInput{}
	var err error
	// Required params
	err = setFieldWithType(params["zone"], input, "HostedZoneId", awsstr)
	if err != nil {
		return nil, err
	}
	resourceRecord := &route53.ResourceRecord{}
	change := &route53.Change{ResourceRecordSet: &route53.ResourceRecordSet{ResourceRecords: []*route53.ResourceRecord{resourceRecord}}}
	input.ChangeBatch = &route53.ChangeBatch{Changes: []*route53.Change{change}}
	err = setFieldWithType("CREATE", change, "Action", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], change, "ResourceRecordSet.Name", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["type"], change, "ResourceRecordSet.Type", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["ttl"], change, "ResourceRecordSet.TTL", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["value"], resourceRecord, "Value", awsstr)
	if err != nil {
		return nil, err
	}

	// Extra params
	if _, ok := params["comment"]; ok {
		err = setFieldWithType(params["comment"], input, "ChangeBatch.Comment", awsstr)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *route53.ChangeResourceRecordSetsOutput
	output, err = d.ChangeResourceRecordSets(input)

	if err != nil {
		return nil, fmt.Errorf("create record: %s", err)
	}
	d.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	d.logger.Info("create record done")
	return aws.StringValue(output.ChangeInfo.Id), nil
}

func (d *Route53Driver) Delete_Record_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["zone"]; !ok {
		return nil, errors.New("delete record: missing required params 'zone'")
	}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("delete record: missing required params 'name'")
	}

	if _, ok := params["type"]; !ok {
		return nil, errors.New("delete record: missing required params 'type'")
	}

	if _, ok := params["value"]; !ok {
		return nil, errors.New("delete record: missing required params 'value'")
	}

	if _, ok := params["ttl"]; !ok {
		return nil, errors.New("delete record: missing required params 'value'")
	}

	d.logger.Verbose("params dry run: delete record ok")
	return nil, nil
}

func (d *Route53Driver) Delete_Record(params map[string]interface{}) (interface{}, error) {
	input := &route53.ChangeResourceRecordSetsInput{}
	var err error
	// Required params
	err = setFieldWithType(params["zone"], input, "HostedZoneId", awsstr)
	if err != nil {
		return nil, err
	}
	resourceRecord := &route53.ResourceRecord{}
	change := &route53.Change{ResourceRecordSet: &route53.ResourceRecordSet{ResourceRecords: []*route53.ResourceRecord{resourceRecord}}}
	input.ChangeBatch = &route53.ChangeBatch{Changes: []*route53.Change{change}}
	err = setFieldWithType("DELETE", change, "Action", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["name"], change, "ResourceRecordSet.Name", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["type"], change, "ResourceRecordSet.Type", awsstr)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["ttl"], change, "ResourceRecordSet.TTL", awsint64)
	if err != nil {
		return nil, err
	}
	err = setFieldWithType(params["value"], resourceRecord, "Value", awsstr)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var output *route53.ChangeResourceRecordSetsOutput
	output, err = d.ChangeResourceRecordSets(input)

	if err != nil {
		return nil, fmt.Errorf("delete record: %s", err)
	}
	d.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	d.logger.Info("delete record done")
	return aws.StringValue(output.ChangeInfo.Id), nil
}

func buildIpPermissionsFromParams(params map[string]interface{}) ([]*ec2.IpPermission, error) {
	if _, ok := params["cidr"].(string); !ok {
		return nil, fmt.Errorf("invalid cidr '%v'", params["cidr"])
	}
	ipPerm := &ec2.IpPermission{
		IpRanges: []*ec2.IpRange{{CidrIp: aws.String(params["cidr"].(string))}},
	}
	if _, ok := params["protocol"].(string); !ok {
		return nil, fmt.Errorf("invalid protocol '%v'", params["protocol"])
	}
	p := params["protocol"].(string)
	if strings.Contains("any", p) {
		ipPerm.FromPort = aws.Int64(int64(-1))
		ipPerm.ToPort = aws.Int64(int64(-1))
		ipPerm.IpProtocol = aws.String("-1")
		return []*ec2.IpPermission{ipPerm}, nil
	}
	ipPerm.IpProtocol = aws.String(p)
	switch ports := params["portrange"].(type) {
	case int:
		ipPerm.FromPort = aws.Int64(int64(ports))
		ipPerm.ToPort = aws.Int64(int64(ports))
	case int64:
		ipPerm.FromPort = aws.Int64(ports)
		ipPerm.ToPort = aws.Int64(ports)
	case string:
		switch {
		case strings.Contains(ports, "any"):
			ipPerm.FromPort = aws.Int64(int64(-1))
			ipPerm.ToPort = aws.Int64(int64(-1))
		case strings.Contains(ports, "-"):
			from, err := strconv.ParseInt(strings.SplitN(ports, "-", 2)[0], 10, 64)
			if err != nil {
				return nil, err
			}
			to, err := strconv.ParseInt(strings.SplitN(ports, "-", 2)[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ipPerm.FromPort = aws.Int64(from)
			ipPerm.ToPort = aws.Int64(to)
		default:
			port, err := strconv.ParseInt(ports, 10, 64)
			if err != nil {
				return nil, err
			}
			ipPerm.FromPort = aws.Int64(port)
			ipPerm.ToPort = aws.Int64(port)
		}
	}

	return []*ec2.IpPermission{ipPerm}, nil
}

func (d *CloudwatchDriver) Attach_Alarm_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"].(string); !ok {
		return nil, errors.New("attach alarm: dry run: missing required params 'name'")
	}

	_, arn := params["action-arn"].(string)

	if !arn {
		return nil, errors.New("attach alarm: dry run: missing 'action-arn' param")
	}

	d.logger.Verbose("params dry run: attach alarm ok")
	return nil, nil
}

func (d *CloudwatchDriver) Attach_Alarm(params map[string]interface{}) (interface{}, error) {
	alarm, err := d.getAlarm(params)
	if err != nil {
		return nil, fmt.Errorf("attach alarm: %s", err)
	}
	alarm.AlarmActions = append(alarm.AlarmActions, aws.String(params["action-arn"].(string)))

	_, err = d.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		ActionsEnabled:                   alarm.ActionsEnabled,
		AlarmActions:                     alarm.AlarmActions,
		AlarmDescription:                 alarm.AlarmDescription,
		AlarmName:                        alarm.AlarmName,
		ComparisonOperator:               alarm.ComparisonOperator,
		Dimensions:                       alarm.Dimensions,
		EvaluateLowSampleCountPercentile: alarm.EvaluateLowSampleCountPercentile,
		EvaluationPeriods:                alarm.EvaluationPeriods,
		ExtendedStatistic:                alarm.ExtendedStatistic,
		InsufficientDataActions:          alarm.InsufficientDataActions,
		MetricName:                       alarm.MetricName,
		Namespace:                        alarm.Namespace,
		OKActions:                        alarm.OKActions,
		Period:                           alarm.Period,
		Statistic:                        alarm.Statistic,
		Threshold:                        alarm.Threshold,
		TreatMissingData:                 alarm.TreatMissingData,
		Unit:                             alarm.Unit,
	})
	if err != nil {
		return nil, fmt.Errorf("attach alarm: %s", err)
	}
	return nil, nil
}

func (d *CloudwatchDriver) Detach_Alarm_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["name"].(string); !ok {
		return nil, errors.New("detach alarm: dry run: missing required params 'name'")
	}

	_, arn := params["action-arn"].(string)

	if !arn {
		return nil, errors.New("detach alarm: dry run: missing 'action-arn' param")
	}

	d.logger.Verbose("params dry run: detach alarm ok")
	return nil, nil
}

func (d *CloudwatchDriver) Detach_Alarm(params map[string]interface{}) (interface{}, error) {
	alarm, err := d.getAlarm(params)
	if err != nil {
		return nil, fmt.Errorf("detach alarm: %s", err)
	}
	actionArn := params["action-arn"].(string)
	var found bool
	var updatedActions []*string
	for _, action := range alarm.AlarmActions {
		if aws.StringValue(action) == actionArn {
			found = true
		} else {
			updatedActions = append(updatedActions, action)
		}
	}
	if !found {
		return nil, fmt.Errorf("detach alarm: action '%s' is not attached to alarm actions of alarm %s", actionArn, aws.StringValue(alarm.AlarmName))
	}

	_, err = d.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		ActionsEnabled:                   alarm.ActionsEnabled,
		AlarmActions:                     updatedActions,
		AlarmDescription:                 alarm.AlarmDescription,
		AlarmName:                        alarm.AlarmName,
		ComparisonOperator:               alarm.ComparisonOperator,
		Dimensions:                       alarm.Dimensions,
		EvaluateLowSampleCountPercentile: alarm.EvaluateLowSampleCountPercentile,
		EvaluationPeriods:                alarm.EvaluationPeriods,
		ExtendedStatistic:                alarm.ExtendedStatistic,
		InsufficientDataActions:          alarm.InsufficientDataActions,
		MetricName:                       alarm.MetricName,
		Namespace:                        alarm.Namespace,
		OKActions:                        alarm.OKActions,
		Period:                           alarm.Period,
		Statistic:                        alarm.Statistic,
		Threshold:                        alarm.Threshold,
		TreatMissingData:                 alarm.TreatMissingData,
		Unit:                             alarm.Unit,
	})
	if err != nil {
		return nil, fmt.Errorf("detach alarm: %s", err)
	}
	return nil, nil
}

func (d *CloudwatchDriver) getAlarm(params map[string]interface{}) (*cloudwatch.MetricAlarm, error) {
	alarm, ok := params["name"].(string)
	if !ok {
		return nil, errors.New("missing required params 'name'")
	}
	out, err := d.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{AlarmNames: []*string{aws.String(alarm)}})
	if err != nil {
		return nil, err
	}
	if l := len(out.MetricAlarms); l == 0 {
		return nil, fmt.Errorf("alarm '%s' not found", alarm)
	} else if l > 1 {
		return nil, fmt.Errorf("%d alarms found with name '%s'", l, alarm)
	}
	return out.MetricAlarms[0], nil
}

func (d *Ec2Driver) Delete_Image_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeregisterImageInput{}
	input.DryRun = aws.Bool(true)
	var err error

	// Required params
	err = setFieldWithType(params["id"], input, "ImageId", awsstr)
	if err != nil {
		return nil, err
	}

	if del, ok := params["delete-snapshots"]; ok && fmt.Sprint(del) == "true" {
		var snaps []string
		if snaps, err = d.imageSnapshots(aws.StringValue(input.ImageId)); err != nil {
			return nil, err
		}
		if len(snaps) > 0 {
			d.logger.Infof("deleting image will also delete snapshot %s (prevent that by appending `delete-snapshots=false`)", strings.Join(snaps, ", "))
		}
	}

	_, err = d.DeregisterImage(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("image")
			d.logger.Verbose("dry run: delete image ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: delete image: %s", err)
}

func (d *Ec2Driver) Delete_Image(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DeregisterImageInput{}
	var err error

	err = setFieldWithType(params["id"], input, "ImageId", awsstr)
	if err != nil {
		return nil, err
	}

	var snaps []string
	if del, ok := params["delete-snapshots"]; ok && fmt.Sprint(del) == "true" {
		if snaps, err = d.imageSnapshots(aws.StringValue(input.ImageId)); err != nil {
			return nil, err
		}
	}

	start := time.Now()
	var output *ec2.DeregisterImageOutput
	output, err = d.DeregisterImage(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("delete image: deregister: %s", err)
	}
	d.logger.ExtraVerbosef("ec2.DeregisterImage call took %s", time.Since(start))
	d.logger.Info("delete image done")

	if del, ok := params["delete-snapshots"]; ok && fmt.Sprint(del) == "true" {
		for _, snap := range snaps {
			snapDelParams := map[string]interface{}{"id": snap}
			if _, err = d.Delete_Snapshot(snapDelParams); err != nil {
				return nil, fmt.Errorf("error while deleting snapshot %s: %s", snap, err)
			}
		}
	}
	return output, nil
}

func (d *Ec2Driver) imageSnapshots(id string) ([]string, error) {
	var snapshots []string
	imgs, err := d.DescribeImages(&ec2.DescribeImagesInput{ImageIds: []*string{aws.String(id)}})
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
		if snapshot := aws.StringValue(dev.Ebs.SnapshotId); snapshot != "" {
			snapshots = append(snapshots, snapshot)
		}
	}
	return snapshots, nil
}

func fakeDryRunId(entity string) string {
	suffix := rand.Intn(1e6)
	switch entity {
	case cloud.Instance:
		return fmt.Sprintf("i-%d", suffix)
	case cloud.Subnet:
		return fmt.Sprintf("subnet-%d", suffix)
	case cloud.Vpc:
		return fmt.Sprintf("vpc-%d", suffix)
	case cloud.Volume:
		return fmt.Sprintf("vol-%d", suffix)
	case cloud.SecurityGroup:
		return fmt.Sprintf("sg-%d", suffix)
	case cloud.InternetGateway:
		return fmt.Sprintf("igw-%d", suffix)
	default:
		return fmt.Sprintf("dryrunid-%d", suffix)
	}
}

type checker struct {
	description string
	timeout     time.Duration
	frequency   time.Duration
	fetchFunc   func() (string, error)
	expect      string
	logger      *logger.Logger
	checkName   string
}

func (c *checker) check() error {
	timer := time.NewTimer(c.timeout)
	if c.checkName == "" {
		c.checkName = "status"
	}
	defer timer.Stop()
	for {
		select {
		case <-time.After(c.frequency):
			got, err := c.fetchFunc()
			if err != nil {
				return fmt.Errorf("check %s: %s", c.description, err)
			}
			if got == c.expect {
				c.logger.Infof("check %s %s '%s' done", c.description, c.checkName, c.expect)
				return nil
			}
			c.logger.Infof("%s %s '%s', expect '%s', retry in %s (timeout %s).", c.description, c.checkName, got, c.expect, c.frequency, c.timeout)
		case <-timer.C:
			return fmt.Errorf("timeout of %s expired", c.timeout)
		}
	}
}
