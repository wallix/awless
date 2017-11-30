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
	"net"
	"strconv"
	"strings"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateSecuritygroup struct {
	_           string `action:"create" entity:"securitygroup" awsAPI:"ec2" awsCall:"CreateSecurityGroup" awsInput:"ec2.CreateSecurityGroupInput" awsOutput:"ec2.CreateSecurityGroupOutput" awsDryRun:""`
	logger      *logger.Logger
	api         ec2iface.EC2API
	Name        *string `awsName:"GroupName" awsType:"awsstr" templateName:"name" required:""`
	Vpc         *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc" required:""`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description" required:""`
}

func (cmd *CreateSecuritygroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateSecuritygroup) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CreateSecurityGroupOutput).GroupId)
}

type UpdateSecuritygroup struct {
	_             string `action:"update" entity:"securitygroup" awsAPI:"ec2" awsDryRun:"manual"`
	logger        *logger.Logger
	api           ec2iface.EC2API
	Id            *string `templateName:"id" required:""`
	Protocol      *string `templateName:"protocol" required:""`
	CIDR          *string `templateName:"cidr"`
	Securitygroup *string `templateName:"securitygroup"`
	Inbound       *string `templateName:"inbound"`
	Outbound      *string `templateName:"outbound"`
	Portrange     *string `templateName:"portrange"`
}

func (cmd *UpdateSecuritygroup) ValidateParams(params []string) ([]string, error) {
	return paramRule{
		tree:   allOf(node("id"), node("protocol"), oneOfE(node("inbound"), node("outbound")), oneOf(node("cidr"), node("securitygroup"))),
		extras: []string{"portrange"},
	}.verify(params)
}

func (cmd *UpdateSecuritygroup) Validate_CIDR() error {
	_, _, err := net.ParseCIDR(StringValue(cmd.CIDR))
	return err
}

// Fail fast when protocol is TCP/UDP and port range is missing, instead of waiting
// for AWS server validation error:
//     InvalidParameterValue: Invalid value 'Must specify both from and to ports with TCP/UDP.' for portRange.
func (cmd *UpdateSecuritygroup) Validate_Protocol() error {
	if p := cmd.Protocol; p != nil {
		if isTCPorUDP(*p) && cmd.Portrange == nil {
			return errors.New("missing 'portrange' when protocol is TCP/UDP")
		}
	}
	return nil
}

func (cmd *UpdateSecuritygroup) Validate_Inbound() error {
	return NewEnumValidator("authorize", "revoke").Validate(cmd.Inbound)
}

func (cmd *UpdateSecuritygroup) Validate_Outbound() error {
	return NewEnumValidator("authorize", "revoke").Validate(cmd.Outbound)
}

func (cmd *UpdateSecuritygroup) DryRun(ctx, params map[string]interface{}) (interface{}, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %s", err)
	}
	ipPerms, err := cmd.buildIpPermissions()
	if err != nil {
		return nil, err
	}
	var input interface{}
	switch StringValue(cmd.Inbound) {
	case "authorize":
		input = &ec2.AuthorizeSecurityGroupIngressInput{DryRun: Bool(true), IpPermissions: ipPerms}
	case "revoke":
		input = &ec2.RevokeSecurityGroupIngressInput{DryRun: Bool(true), IpPermissions: ipPerms}
	}
	switch StringValue(cmd.Outbound) {
	case "authorize":
		input = &ec2.AuthorizeSecurityGroupEgressInput{DryRun: Bool(true), IpPermissions: ipPerms}
	case "revoke":
		input = &ec2.RevokeSecurityGroupEgressInput{DryRun: Bool(true), IpPermissions: ipPerms}
	}
	if input == nil {
		return nil, fmt.Errorf("expect either 'inbound' or 'outbound' parameter")
	}

	err = setFieldWithType(cmd.Id, input, "GroupId", awsstr)
	if err != nil {
		return nil, err
	}

	switch ii := input.(type) {
	case *ec2.AuthorizeSecurityGroupIngressInput:
		_, err = cmd.api.AuthorizeSecurityGroupIngress(ii)
	case *ec2.RevokeSecurityGroupIngressInput:
		_, err = cmd.api.RevokeSecurityGroupIngress(ii)
	case *ec2.AuthorizeSecurityGroupEgressInput:
		_, err = cmd.api.AuthorizeSecurityGroupEgress(ii)
	case *ec2.RevokeSecurityGroupEgressInput:
		_, err = cmd.api.RevokeSecurityGroupEgress(ii)
	}
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			cmd.logger.Verbose("dry run: update securitygroup ok")
			return nil, nil
		}
	}
	return nil, fmt.Errorf("dry run: update securitygroup: %s", err)
}

func (cmd *UpdateSecuritygroup) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	ipPerms, err := cmd.buildIpPermissions()
	if err != nil {
		return nil, err
	}
	var input interface{}
	if inbound := cmd.Inbound; inbound != nil {
		switch StringValue(inbound) {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupIngressInput{IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupIngressInput{IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'inbound' parameter expect 'authorize' or 'revoke', got %s", StringValue(inbound))
		}
	}
	if outbound := cmd.Outbound; outbound != nil {
		switch StringValue(outbound) {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupEgressInput{IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupEgressInput{IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'outbound' parameter expect 'authorize' or 'revoke', got %s", StringValue(outbound))
		}
	}
	if input == nil {
		return nil, fmt.Errorf("expect either 'inbound' or 'outbound' parameter")
	}

	// Required params
	err = setFieldWithType(cmd.Id, input, "GroupId", awsstr)
	if err != nil {
		return nil, err
	}

	var output interface{}
	start := time.Now()
	switch ii := input.(type) {
	case *ec2.AuthorizeSecurityGroupIngressInput:
		output, err = cmd.api.AuthorizeSecurityGroupIngress(ii)
		cmd.logger.ExtraVerbosef("ec2.AuthorizeSecurityGroupIngress call took %s", time.Since(start))
	case *ec2.RevokeSecurityGroupIngressInput:
		output, err = cmd.api.RevokeSecurityGroupIngress(ii)
		cmd.logger.ExtraVerbosef("ec2.RevokeSecurityGroupIngress call took %s", time.Since(start))
	case *ec2.AuthorizeSecurityGroupEgressInput:
		output, err = cmd.api.AuthorizeSecurityGroupEgress(ii)
		cmd.logger.ExtraVerbosef("ec2.AuthorizeSecurityGroupEgress call took %s", time.Since(start))
	case *ec2.RevokeSecurityGroupEgressInput:
		output, err = cmd.api.RevokeSecurityGroupEgress(ii)
		cmd.logger.ExtraVerbosef("ec2.RevokeSecurityGroupEgress call took %s", time.Since(start))
	}

	return output, err
}

type DeleteSecuritygroup struct {
	_      string `action:"delete" entity:"securitygroup" awsAPI:"ec2" awsCall:"DeleteSecurityGroup" awsInput:"ec2.DeleteSecurityGroupInput" awsOutput:"ec2.DeleteSecurityGroupOutput" awsDryRun:""`
	logger *logger.Logger
	api    ec2iface.EC2API
	Id     *string `awsName:"GroupId" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *DeleteSecuritygroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type CheckSecuritygroup struct {
	_       string `action:"check" entity:"securitygroup" awsAPI:"ec2"`
	logger  *logger.Logger
	api     ec2iface.EC2API
	Id      *string `templateName:"id" required:""`
	State   *string `templateName:"state" required:""`
	Timeout *int64  `templateName:"timeout" required:""`
}

func (cmd *CheckSecuritygroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CheckSecuritygroup) Validate_State() error {
	return NewEnumValidator("unused").Validate(cmd.State)
}

func (cmd *CheckSecuritygroup) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{Name: String("group-id"), Values: []*string{cmd.Id}},
		},
	}

	c := &checker{
		description: fmt.Sprintf("securitygroup %s", StringValue(cmd.Id)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.DescribeNetworkInterfaces(input)
			if err != nil {
				return "", err
			}
			if len(output.NetworkInterfaces) == 0 {
				return "unused", nil
			}
			var niIds []string
			for _, ni := range output.NetworkInterfaces {
				niIds = append(niIds, StringValue(ni.NetworkInterfaceId))
			}
			return fmt.Sprintf("used by %s", strings.Join(niIds, ", ")), nil
		},
		expect: StringValue(cmd.State),
		logger: cmd.logger,
	}
	return nil, c.check()
}

type AttachSecuritygroup struct {
	_        string `action:"attach" entity:"securitygroup" awsAPI:"ec2"`
	logger   *logger.Logger
	api      ec2iface.EC2API
	Id       *string `templateName:"id" required:""`
	Instance *string `templateName:"instance" required:""`
}

func (cmd *AttachSecuritygroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *AttachSecuritygroup) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	groups, err := fetchInstanceSecurityGroups(cmd.api, StringValue(cmd.Instance))
	if err != nil {
		return nil, fmt.Errorf("fetching securitygroups for instance %s: %s", StringValue(cmd.Instance), err)
	}

	groups = append(groups, StringValue(cmd.Id))
	call := &awsCall{
		fnName: "ec2.ModifyInstanceAttribute",
		fn:     cmd.api.ModifyInstanceAttribute,
		logger: cmd.logger,
		setters: []setter{
			{val: cmd.Instance, fieldPath: "InstanceID", fieldType: awsstr},
			{val: groups, fieldPath: "Groups", fieldType: awsstringslice},
		},
	}
	return call.execute(&ec2.ModifyInstanceAttributeInput{})
}

type DetachSecuritygroup struct {
	_        string `action:"detach" entity:"securitygroup" awsAPI:"ec2"`
	logger   *logger.Logger
	api      ec2iface.EC2API
	Id       *string `templateName:"id" required:""`
	Instance *string `templateName:"instance" required:""`
}

func (cmd *DetachSecuritygroup) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *DetachSecuritygroup) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	groups, err := fetchInstanceSecurityGroups(cmd.api, StringValue(cmd.Instance))
	if err != nil {
		return nil, fmt.Errorf("fetching securitygroups for instance %s: %s", StringValue(cmd.Instance), err)
	}

	cleaned := removeString(groups, StringValue(cmd.Id))
	if len(cleaned) == 0 {
		cmd.logger.Errorf("AWS instances must have at least one securitygroup")
	}
	call := &awsCall{
		fnName: "ec2.ModifyInstanceAttribute",
		fn:     cmd.api.ModifyInstanceAttribute,
		logger: cmd.logger,
		setters: []setter{
			{val: cmd.Instance, fieldPath: "InstanceID", fieldType: awsstr},
			{val: cleaned, fieldPath: "Groups", fieldType: awsstringslice},
		},
	}
	return call.execute(&ec2.ModifyInstanceAttributeInput{})
}

func (cmd *UpdateSecuritygroup) buildIpPermissions() ([]*ec2.IpPermission, error) {
	ipPerm := &ec2.IpPermission{}
	if cidr := cmd.CIDR; cidr != nil {
		ipPerm.IpRanges = []*ec2.IpRange{{CidrIp: cidr}}
	} else if secgroup := cmd.Securitygroup; secgroup != nil {
		ipPerm.UserIdGroupPairs = []*ec2.UserIdGroupPair{{GroupId: secgroup}}
	} else {
		return nil, errors.New("missing either 'cidr' or 'securitygroup' parameter")
	}

	p := StringValue(cmd.Protocol)
	if strings.Contains("any", p) {
		ipPerm.FromPort = Int64(-1)
		ipPerm.ToPort = Int64(-1)
		ipPerm.IpProtocol = String("-1")
		return []*ec2.IpPermission{ipPerm}, nil
	}
	ipPerm.IpProtocol = String(p)

	if pRange := cmd.Portrange; pRange != nil {
		ports := *pRange
		switch {
		case strings.Contains(ports, "any"):
			if isTCPorUDP(p) {
				ipPerm.FromPort = Int64(int64(0))
				ipPerm.ToPort = Int64(int64(65535))
			} else {
				ipPerm.FromPort = Int64(int64(-1))
				ipPerm.ToPort = Int64(int64(-1))
			}
		case strings.Contains(ports, "-"):
			from, err := strconv.ParseInt(strings.SplitN(ports, "-", 2)[0], 10, 64)
			if err != nil {
				return nil, err
			}
			to, err := strconv.ParseInt(strings.SplitN(ports, "-", 2)[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ipPerm.FromPort = Int64(from)
			ipPerm.ToPort = Int64(to)
		default:
			port, err := strconv.ParseInt(ports, 10, 64)
			if err != nil {
				return nil, err
			}
			ipPerm.FromPort = Int64(port)
			ipPerm.ToPort = Int64(port)
		}
	}

	return []*ec2.IpPermission{ipPerm}, nil
}

func isTCPorUDP(p string) bool {
	return strings.ToLower(p) == "tcp" || strings.ToLower(p) == "udp"
}

func fetchInstanceSecurityGroups(api ec2iface.EC2API, id string) ([]string, error) {
	params := &ec2.DescribeInstanceAttributeInput{
		Attribute:  String("groupSet"),
		InstanceId: String(id),
	}
	resp, err := api.DescribeInstanceAttribute(params)
	if err != nil {
		return nil, err
	}

	var groups []string
	for _, g := range resp.Groups {
		groups = append(groups, StringValue(g.GroupId))
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
