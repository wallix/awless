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

package aws

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/console"
)

func (d *IamDriver) Attach_Policy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("attach policy: missing required params 'arn'")
	}

	_, hasUser := params["user"]
	_, hasGroup := params["group"]

	if !hasUser && !hasGroup {
		return nil, errors.New("attach policy: missing one of 'user, group' param")
	}

	d.logger.Verbose("params dry run: attach policy ok")
	return nil, nil
}

func (d *IamDriver) Attach_Policy(params map[string]interface{}) (interface{}, error) {
	user, hasUser := params["user"]
	group, hasGroup := params["group"]

	switch {
	case hasUser:
		return performCall(d, "attach user", &iam.AttachUserPolicyInput{}, d.AttachUserPolicy, []setter{
			{val: params["arn"], fieldPath: "PolicyArn", fieldType: awsstr},
			{val: user, fieldPath: "UserName", fieldType: awsstr},
		}...)
	case hasGroup:
		return performCall(d, "attach user", &iam.AttachGroupPolicyInput{}, d.AttachGroupPolicy, []setter{
			{val: params["arn"], fieldPath: "PolicyArn", fieldType: awsstr},
			{val: group, fieldPath: "GroupName", fieldType: awsstr},
		}...)
	}

	return nil, errors.New("missing one of 'user, group' param")
}

func (d *IamDriver) Detach_Policy_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["arn"]; !ok {
		return nil, errors.New("detach policy: missing required params 'arn'")
	}

	_, hasUser := params["user"]
	_, hasGroup := params["group"]

	if !hasUser && !hasGroup {
		return nil, errors.New("detach policy: missing one of 'user, group' param")
	}

	d.logger.Verbose("params dry run: detach policy ok")
	return nil, nil
}

func (d *IamDriver) Detach_Policy(params map[string]interface{}) (interface{}, error) {
	user, hasUser := params["user"]
	group, hasGroup := params["group"]

	switch {
	case hasUser:
		return performCall(d, "detach user", &iam.DetachUserPolicyInput{}, d.DetachUserPolicy, []setter{
			{val: params["arn"], fieldPath: "PolicyArn", fieldType: awsstr},
			{val: user, fieldPath: "UserName", fieldType: awsstr},
		}...)
	case hasGroup:
		return performCall(d, "detach user", &iam.DetachGroupPolicyInput{}, d.DetachGroupPolicy, []setter{
			{val: params["arn"], fieldPath: "PolicyArn", fieldType: awsstr},
			{val: group, fieldPath: "GroupName", fieldType: awsstr},
		}...)
	}

	return nil, errors.New("missing one of 'user, group' param")
}

type setter struct {
	val       interface{}
	fieldPath string
	fieldType int
}

func performCall(d *IamDriver, desc string, input interface{}, fn interface{}, setters ...setter) (output interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			output = nil
			err = fmt.Errorf("%s", e)
		}
	}()

	for _, set := range setters {
		if err = setFieldWithType(set.val, input, set.fieldPath, set.fieldType); err != nil {
			return nil, err
		}
	}

	fnVal := reflect.ValueOf(fn)
	values := []reflect.Value{reflect.ValueOf(input)}

	start := time.Now()
	results := fnVal.Call(values)

	if err, ok := results[1].Interface().(error); ok && err != nil {
		d.logger.Errorf("%s error: %s", desc, err)
		return nil, err
	}

	d.logger.ExtraVerbosef("%s call took %s", desc, time.Since(start))
	d.logger.Verbosef("%s done", desc)

	output = results[0].Interface()

	return
}

func (d *Ec2Driver) Check_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DescribeInstancesInput{}
	input.DryRun = aws.Bool(true)

	for _, val := range []string{"state", "id", "timeout"} {
		if _, ok := params[val]; !ok {
			err := fmt.Errorf("check instance error: missing required param '%s'", val)
			d.logger.Errorf("%s", err)
			return nil, err
		}
	}

	if _, ok := params["timeout"].(int); !ok {
		err := errors.New("check instance error: timeout param is not int")
		d.logger.Errorf("%s", err)
		return nil, err
	}

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
			d.logger.Verbose("full dry run: check instance ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: check instance error: %s", err)
	return nil, err
}

func (d *Ec2Driver) Check_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DescribeInstancesInput{}

	// Required params
	err := setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(params["timeout"].(int)) * time.Second
	timer := time.NewTimer(timeout)
	retry := 5 * time.Second
	for {
		select {
		case <-time.After(retry):
			output, err := d.DescribeInstances(input)
			if err != nil {
				d.logger.Errorf("check instance error: %s", err)
				return nil, err
			}

			if res := output.Reservations; len(res) > 0 {
				if instances := output.Reservations[0].Instances; len(instances) > 0 {
					for _, inst := range instances {
						if aws.StringValue(inst.InstanceId) == params["id"] {
							currentStatus := aws.StringValue(inst.State.Name)
							if currentStatus == params["state"] {
								d.logger.Verbosef("check instance status '%s' done", params["state"])
								timer.Stop()
								return nil, nil
							}
							d.logger.Infof("instance status '%s', expect '%s', retry in %s (timeout %s).", currentStatus, params["state"], retry, timeout)
						}
					}
				}
			}

		case <-timer.C:
			err := fmt.Errorf("timeout of %s expired", timeout)
			d.logger.Errorf("%s", err)
			return nil, err
		}
	}
}

func (d *Ec2Driver) Create_Tags_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateTagsInput{}

	input.DryRun = aws.Bool(true)
	input.Resources = append(input.Resources, aws.String(fmt.Sprint(params["resource"])))

	for k, v := range params {
		if k == "resource" {
			continue
		}
		input.Tags = append(input.Tags, &ec2.Tag{Key: aws.String(k), Value: aws.String(fmt.Sprint(v))})
	}
	_, err := d.CreateTags(input)

	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			d.logger.Verbose("full dry run: create tags ok")
			return nil, nil
		}
	}

	d.logger.Errorf("dry run: create tags error: %s", err)
	return nil, err
}

func (d *Ec2Driver) Create_Tags(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateTagsInput{}

	input.Resources = append(input.Resources, aws.String(fmt.Sprint(params["resource"])))

	for k, v := range params {
		if k == "resource" {
			continue
		}
		input.Tags = append(input.Tags, &ec2.Tag{Key: aws.String(k), Value: aws.String(fmt.Sprint(v))})
	}
	_, err := d.CreateTags(input)

	if err != nil {
		d.logger.Errorf("create tags error: %s", err)
		return nil, err
	}
	d.logger.Verbose("create tags done")

	return nil, nil
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
			d.logger.Verbose("full dry run: create tag ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: create tag error: %s", err)
	return nil, err
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
		d.logger.Errorf("create tag error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("ec2.CreateTags call took %s", time.Since(start))
	d.logger.Verbose("create tag done")
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
		err = fmt.Errorf("empty 'name' parameter")
		d.logger.Errorf("dry run: saving private key error: %s", err)
		return nil, err
	}

	const keyDirEnv = "__AWLESS_KEYS_DIR"
	keyDir := os.Getenv(keyDirEnv)
	if keyDir == "" {
		err = fmt.Errorf("empty env var '%s'", keyDirEnv)
		d.logger.Errorf("dry run: saving private key error: %s", err)
		return nil, err
	}

	privKeyPath := filepath.Join(keyDir, fmt.Sprint(params["name"])+".pem")
	_, err = os.Stat(privKeyPath)
	if err == nil {
		fileExist := fmt.Errorf("file already exists at path: %s", privKeyPath)
		d.logger.Errorf("dry run: saving private key error: %s", fileExist)
		return nil, fileExist
	}

	return nil, nil
}

func (d *Ec2Driver) Create_Keypair(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ImportKeyPairInput{}
	err := setFieldWithType(params["name"], input, "KeyName", awsstr)
	if err != nil {
		return nil, err
	}

	d.logger.Info("Generating locally a RSA 4096 bits keypair...")
	pub, priv, err := console.GenerateSSHKeyPair(4096)
	if err != nil {
		d.logger.Errorf("generating keypair error: %s", err)
		return nil, err
	}
	privKeyPath := filepath.Join(os.Getenv("__AWLESS_KEYS_DIR"), fmt.Sprint(params["name"])+".pem")
	_, err = os.Stat(privKeyPath)
	if err == nil {
		fileExist := fmt.Errorf("file already exists at path: %s", privKeyPath)
		d.logger.Errorf("saving private key error: %s", fileExist)
		return nil, fileExist
	}
	err = ioutil.WriteFile(privKeyPath, priv, 0400)
	if err != nil {
		d.logger.Errorf("saving private key error: %s", err)
		return nil, err
	}
	d.logger.Infof("4096 RSA keypair generated locally and stored in '%s'", privKeyPath)
	input.PublicKeyMaterial = pub

	output, err := d.ImportKeyPair(input)
	if err != nil {
		d.logger.Errorf("create keypair error: %s", err)
		return nil, err
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
			d.logger.Verbose("full dry run: update securitygroup ok")
			return nil, nil
		}
	}

	d.logger.Errorf("dry run: update securitygroup error: %s", err)
	return nil, err
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
		d.logger.Errorf("update securitygroup error: %s", err)
		return nil, err
	}

	d.logger.Verbose("update securitygroup done")
	return output, nil
}

func (d *S3Driver) Create_Storageobject_DryRun(params map[string]interface{}) (interface{}, error) {
	if _, ok := params["bucket"]; !ok {
		return nil, errors.New("create storageobject: missing required params 'bucket'")
	}

	if _, ok := params["file"].(string); !ok {
		return nil, errors.New("create storageobject: missing required string params 'file'")
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

	d.logger.Verbose("params dry run: create storageobject ok")
	return nil, nil
}

func (d *S3Driver) Create_Storageobject(params map[string]interface{}) (interface{}, error) {
	input := &s3.PutObjectInput{}

	f, err := os.Open(params["file"].(string))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	input.Body = f

	var fileName string
	if n, ok := params["name"].(string); ok && n != "" {
		fileName = n
	} else {
		fileName = f.Name()
	}
	input.Key = aws.String(fileName)

	// Required params
	err = setFieldWithType(params["bucket"], input, "Bucket", awsstr)
	if err != nil {
		return nil, err
	}

	output, err := d.PutObject(input)
	if err != nil {
		d.logger.Errorf("create storageobject error: %s", err)
		return nil, err
	}

	d.logger.Verbose("create storageobject done")
	return output, nil
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
	output = output
	if err != nil {
		d.logger.Errorf("create record error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	d.logger.Verbose("create record done")
	return output, nil
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
	output = output
	if err != nil {
		d.logger.Errorf("delete record error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	d.logger.Verbose("delete record done")
	return output, nil
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

func fakeDryRunId(entity string) string {
	suffix := rand.Intn(1e6)
	switch entity {
	case cloud.Instance:
		return fmt.Sprintf("i-%d", suffix)
	case cloud.Subnet:
		return fmt.Sprintf("subnet-%d", suffix)
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
