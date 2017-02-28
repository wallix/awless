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
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
)

func (d *AwsDriver) Check_Instance_DryRun(params map[string]interface{}) (interface{}, error) {
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

	_, err = d.ec2.DescribeInstances(input)
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

func (d *AwsDriver) Check_Instance(params map[string]interface{}) (interface{}, error) {
	input := &ec2.DescribeInstancesInput{}

	// Required params
	err := setFieldWithType(params["id"], input, "InstanceIds", awsstringslice)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(params["timeout"].(int)) * time.Second
	timer := time.NewTimer(timeout)

	for {
		select {
		case <-time.After(1 * time.Second):
			output, err := d.ec2.DescribeInstances(input)
			if err != nil {
				d.logger.Errorf("check instance error: %s", err)
				return nil, err
			}

			if res := output.Reservations; len(res) > 0 {
				if instances := output.Reservations[0].Instances; len(instances) > 0 {
					for _, inst := range instances {
						if aws.StringValue(inst.InstanceId) == params["id"] {
							if aws.StringValue(inst.State.Name) == params["state"] {
								d.logger.Verbosef("check instance status '%s' done", params["state"])
								timer.Stop()
								return nil, nil
							}
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

func (d *AwsDriver) Create_Tags_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateTagsInput{}

	input.DryRun = aws.Bool(true)
	input.Resources = append(input.Resources, aws.String(fmt.Sprint(params["resource"])))

	for k, v := range params {
		if k == "resource" {
			continue
		}
		input.Tags = append(input.Tags, &ec2.Tag{Key: aws.String(k), Value: aws.String(fmt.Sprint(v))})
	}
	_, err := d.ec2.CreateTags(input)

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

func (d *AwsDriver) Create_Tags(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateTagsInput{}

	input.Resources = append(input.Resources, aws.String(fmt.Sprint(params["resource"])))

	for k, v := range params {
		if k == "resource" {
			continue
		}
		input.Tags = append(input.Tags, &ec2.Tag{Key: aws.String(k), Value: aws.String(fmt.Sprint(v))})
	}
	_, err := d.ec2.CreateTags(input)

	if err != nil {
		d.logger.Errorf("create tags error: %s", err)
		return nil, err
	}
	d.logger.Verbose("create tags done")

	return nil, nil
}

func (d *AwsDriver) Create_Keypair_DryRun(params map[string]interface{}) (interface{}, error) {
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

	privKeyPath := filepath.Join(config.KeysDir, fmt.Sprint(params["name"])+".pem")
	_, err = os.Stat(privKeyPath)
	if err == nil {
		fileExist := fmt.Errorf("file already exists at path: %s", privKeyPath)
		d.logger.Errorf("dry run: saving private key error: %s", fileExist)
		return nil, fileExist
	}

	return nil, nil
}

func (d *AwsDriver) Create_Keypair(params map[string]interface{}) (interface{}, error) {
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
	privKeyPath := filepath.Join(config.KeysDir, fmt.Sprint(params["name"])+".pem")
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
	fmt.Printf("4096 RSA keypair generated locally and stored in '%s'\n", privKeyPath)
	input.PublicKeyMaterial = pub

	output, err := d.ec2.ImportKeyPair(input)
	if err != nil {
		d.logger.Errorf("create keypair error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.KeyName)
	d.logger.Infof("create keypair '%s' done", id)
	return aws.StringValue(output.KeyName), nil
}

func (d *AwsDriver) Update_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
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
		_, err = d.ec2.AuthorizeSecurityGroupIngress(ii)
	case *ec2.RevokeSecurityGroupIngressInput:
		_, err = d.ec2.RevokeSecurityGroupIngress(ii)
	case *ec2.AuthorizeSecurityGroupEgressInput:
		_, err = d.ec2.AuthorizeSecurityGroupEgress(ii)
	case *ec2.RevokeSecurityGroupEgressInput:
		_, err = d.ec2.RevokeSecurityGroupEgress(ii)
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

func (d *AwsDriver) Update_Securitygroup(params map[string]interface{}) (interface{}, error) {
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
		output, err = d.ec2.AuthorizeSecurityGroupIngress(ii)
	case *ec2.RevokeSecurityGroupIngressInput:
		output, err = d.ec2.RevokeSecurityGroupIngress(ii)
	case *ec2.AuthorizeSecurityGroupEgressInput:
		output, err = d.ec2.AuthorizeSecurityGroupEgress(ii)
	case *ec2.RevokeSecurityGroupEgressInput:
		output, err = d.ec2.RevokeSecurityGroupEgress(ii)
	}
	if err != nil {
		d.logger.Errorf("update securitygroup error: %s", err)
		return nil, err
	}

	d.logger.Verbose("update securitygroup done")
	return output, nil
}

// This function was auto generated
func (d *AwsDriver) Create_Storageobject_DryRun(params map[string]interface{}) (interface{}, error) {
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

// This function was auto generated
func (d *AwsDriver) Create_Storageobject(params map[string]interface{}) (interface{}, error) {
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

	output, err := d.s3.PutObject(input)
	if err != nil {
		d.logger.Errorf("create storageobject error: %s", err)
		return nil, err
	}

	d.logger.Verbose("create storageobject done")
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
	case "instance":
		return fmt.Sprintf("i-%d", suffix)
	case "volume":
		return fmt.Sprintf("vol-%d", suffix)
	case "securitygroup":
		return fmt.Sprintf("sg-%d", suffix)
	case graph.InternetGateway.String():
		return fmt.Sprintf("igw-%d", suffix)
	default:
		return fmt.Sprintf("dryrunid-%d", suffix)
	}
}
