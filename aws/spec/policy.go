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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/wallix/awless/logger"
)

type CreatePolicy struct {
	_           string `action:"create" entity:"policy" awsAPI:"iam" awsCall:"CreatePolicy" awsInput:"iam.CreatePolicyInput" awsOutput:"iam.CreatePolicyOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         iamiface.IAMAPI
	Name        *string   `awsName:"PolicyName" awsType:"awsstr" templateName:"name"`
	Effect      *string   `templateName:"effect"`
	Action      []*string `templateName:"action"`
	Resource    []*string `templateName:"resource"`
	Description *string   `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Document    *string   `awsName:"PolicyDocument" awsType:"awsstr"`
	Conditions  []*string `templateName:"conditions"`
}

func (cmd *CreatePolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("action"), params.Key("effect"), params.Key("name"), params.Key("resource"),
		params.Opt("conditions", "description"),
	))
}

func (cmd *CreatePolicy) BeforeRun(renv env.Running) error {
	stat, err := buildStatementFromParams(cmd.Effect, cmd.Resource, cmd.Action, cmd.Conditions)
	if err != nil {
		return err
	}
	policy := &policyBody{
		Version:   "2012-10-17",
		Statement: []*policyStatement{stat},
	}

	b, err := json.MarshalIndent(policy, "", " ")
	if err != nil {
		return fmt.Errorf("cannot marshal policy document: %s", err)
	}
	cmd.Document = String(string(b))
	cmd.logger.ExtraVerbosef("policy document json:\n%s\n", string(b))
	return nil
}

func (cmd *CreatePolicy) ExtractResult(i interface{}) string {
	return StringValue(i.(*iam.CreatePolicyOutput).Policy.Arn)
}

type UpdatePolicy struct {
	_              string `action:"update" entity:"policy" awsAPI:"iam" awsCall:"CreatePolicyVersion" awsInput:"iam.CreatePolicyVersionInput" awsOutput:"iam.CreatePolicyVersionOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            iamiface.IAMAPI
	Arn            *string   `awsName:"PolicyArn" awsType:"awsstr" templateName:"arn"`
	Effect         *string   `templateName:"effect"`
	Action         []*string `templateName:"action"`
	Resource       []*string `templateName:"resource"`
	Conditions     []*string `templateName:"conditions"`
	Document       *string   `awsName:"PolicyDocument" awsType:"awsstr"`
	DefaultVersion *bool     `awsName:"SetAsDefault" awsType:"awsbool"`
}

func (cmd *UpdatePolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("action"), params.Key("arn"), params.Key("effect"), params.Key("resource"),
		params.Opt("conditions"),
	))
}

func (cmd *UpdatePolicy) BeforeRun(renv env.Running) error {
	document, err := cmd.getPolicyLastVersionDocument(cmd.Arn)
	if err != nil {
		return err
	}
	var defaultPolicyDocument *struct {
		Version    string             `json:",omitempty"`
		ID         string             `json:"Id,omitempty"`
		Statements []*json.RawMessage `json:"Statement,omitempty"`
	}

	if err = json.Unmarshal([]byte(document), &defaultPolicyDocument); err != nil {
		return err
	}
	stat, err := buildStatementFromParams(cmd.Effect, cmd.Resource, cmd.Action, cmd.Conditions)
	if err != nil {
		return err
	}

	var newStatement json.RawMessage
	if newStatement, err = json.Marshal(stat); err != nil {
		return err
	}
	defaultPolicyDocument.Statements = append(defaultPolicyDocument.Statements, &newStatement)

	b, err := json.MarshalIndent(defaultPolicyDocument, "", " ")
	if err != nil {
		return fmt.Errorf("cannot marshal policy document: %s", err)
	}
	cmd.Document = String(string(b))
	cmd.DefaultVersion = aws.Bool(true)
	cmd.logger.ExtraVerbosef("policy document json:\n%s\n", string(b))
	return nil
}

func (cmd *UpdatePolicy) getPolicyLastVersionDocument(arn *string) (string, error) {
	listVersionsInput := &iam.ListPolicyVersionsInput{
		PolicyArn: arn,
	}
	listVersionsOut, err := cmd.api.ListPolicyVersions(listVersionsInput)
	if err != nil {
		return "", err
	}
	var defaultVersion *iam.PolicyVersion
	for _, version := range listVersionsOut.Versions {
		if aws.BoolValue(version.IsDefaultVersion) {
			policyDetailInput := &iam.GetPolicyVersionInput{
				VersionId: version.VersionId,
				PolicyArn: arn,
			}
			var policyDetailOutput *iam.GetPolicyVersionOutput
			if policyDetailOutput, err = cmd.api.GetPolicyVersion(policyDetailInput); err != nil {
				return "", err
			}
			defaultVersion = policyDetailOutput.PolicyVersion
		}
	}
	if defaultVersion == nil {
		return "", fmt.Errorf("update policy: can not find default version for policy with arn '%s'", StringValue(arn))
	}
	document, err := url.QueryUnescape(aws.StringValue(defaultVersion.Document))
	if err != nil {
		return "", fmt.Errorf("decoding policy document: %s", err)
	}
	return document, nil
}

type DeletePolicy struct {
	_           string `action:"delete" entity:"policy" awsAPI:"iam"  awsCall:"DeletePolicy" awsInput:"iam.DeletePolicyInput" awsOutput:"iam.DeletePolicyOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         iamiface.IAMAPI
	Arn         *string `awsName:"PolicyArn" awsType:"awsstr" templateName:"arn"`
	AllVersions *bool   `templateName:"all-versions"`
}

func (cmd *DeletePolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("arn"),
		params.Opt("all-versions"),
	))
}

func (cmd *DeletePolicy) BeforeRun(renv env.Running) error {
	if BoolValue(cmd.AllVersions) {
		list, err := cmd.api.ListPolicyVersions(&iam.ListPolicyVersionsInput{PolicyArn: cmd.Arn})
		if err != nil {
			return fmt.Errorf("list all policy versions: %s", err)
		}
		for _, v := range list.Versions {
			if !aws.BoolValue(v.IsDefaultVersion) {
				cmd.logger.Verbosef("deleting version '%s' of policy '%s'", aws.StringValue(v.VersionId), StringValue(cmd.Arn))
				if _, err := cmd.api.DeletePolicyVersion(&iam.DeletePolicyVersionInput{PolicyArn: cmd.Arn, VersionId: v.VersionId}); err != nil {
					return fmt.Errorf("delete version %s: %s", aws.StringValue(v.VersionId), err)
				}
			}
		}
	}
	return nil
}

type AttachPolicy struct {
	_       string `action:"attach" entity:"policy" awsAPI:"iam"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     iamiface.IAMAPI
	Arn     *string `awsName:"PolicyArn" awsType:"awsstr" templateName:"arn"`
	User    *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
	Group   *string `awsName:"GroupName" awsType:"awsstr" templateName:"group"`
	Role    *string `awsName:"RoleName" awsType:"awsstr" templateName:"role"`
	Service *string `templateName:"service"`
	Access  *string `templateName:"access"`
}

func (cmd *AttachPolicy) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.AllOf(
		params.OnlyOneOf(params.Key("user"), params.Key("role"), params.Key("group")),
		params.OnlyOneOf(params.Key("arn"), params.AllOf(params.Key("access"), params.Key("service"))),
	))
	builder.AddReducer(transformAccessServiceToARN, "access", "service")
	return builder.Done()
}

func transformAccessServiceToARN(values map[string]interface{}) (map[string]interface{}, error) {
	service, hasService := values["service"].(string)
	access, hasAccess := values["access"].(string)

	if hasService && hasAccess {
		pol, err := lookupAWSPolicy(service, access)
		if err != nil {
			return values, err
		}
		return map[string]interface{}{"arn": pol.Arn}, nil
	} else {
		return nil, nil
	}
}

func (cmd *AttachPolicy) ManualRun(renv env.Running) (interface{}, error) {
	start := time.Now()
	switch {
	case cmd.User != nil:
		input := &iam.AttachUserPolicyInput{}
		input.PolicyArn = cmd.Arn
		input.UserName = cmd.User
		output, err := cmd.api.AttachUserPolicy(input)
		cmd.logger.ExtraVerbosef("ec2.AttachUserPolicy call took %s", time.Since(start))
		return output, err
	case cmd.Group != nil:
		input := &iam.AttachGroupPolicyInput{}
		input.PolicyArn = cmd.Arn
		input.GroupName = cmd.Group
		output, err := cmd.api.AttachGroupPolicy(input)
		cmd.logger.ExtraVerbosef("ec2.AttachGroupPolicy call took %s", time.Since(start))
		return output, err
	case cmd.Role != nil:
		input := &iam.AttachRolePolicyInput{}
		input.PolicyArn = cmd.Arn
		input.RoleName = cmd.Role
		output, err := cmd.api.AttachRolePolicy(input)
		cmd.logger.ExtraVerbosef("ec2.AttachRolePolicy call took %s", time.Since(start))
		return output, err
	default:
		return nil, errors.New("missing one of 'user, group, role' param")
	}
}

type DetachPolicy struct {
	_      string `action:"detach" entity:"policy" awsAPI:"iam"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Arn    *string `awsName:"PolicyArn" awsType:"awsstr" templateName:"arn"`
	User   *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
	Group  *string `awsName:"GroupName" awsType:"awsstr" templateName:"group"`
	Role   *string `awsName:"RoleName" awsType:"awsstr" templateName:"role"`
}

func (cmd *DetachPolicy) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.AllOf(
		params.OnlyOneOf(params.Key("user"), params.Key("role"), params.Key("group")),
		params.OnlyOneOf(params.Key("arn"), params.AllOf(params.Key("access"), params.Key("service"))),
	))
	builder.AddReducer(transformAccessServiceToARN, "access", "service")
	return builder.Done()
}

func (cmd *DetachPolicy) ManualRun(renv env.Running) (interface{}, error) {
	start := time.Now()
	switch {
	case cmd.User != nil:
		input := &iam.DetachUserPolicyInput{}
		input.PolicyArn = cmd.Arn
		input.UserName = cmd.User
		output, err := cmd.api.DetachUserPolicy(input)
		cmd.logger.ExtraVerbosef("ec2.DetachUserPolicy call took %s", time.Since(start))
		return output, err
	case cmd.Group != nil:
		input := &iam.DetachGroupPolicyInput{}
		input.PolicyArn = cmd.Arn
		input.GroupName = cmd.Group
		output, err := cmd.api.DetachGroupPolicy(input)
		cmd.logger.ExtraVerbosef("ec2.DetachGroupPolicy call took %s", time.Since(start))
		return output, err
	case cmd.Role != nil:
		input := &iam.DetachRolePolicyInput{}
		input.PolicyArn = cmd.Arn
		input.RoleName = cmd.Role
		output, err := cmd.api.DetachRolePolicy(input)
		cmd.logger.ExtraVerbosef("ec2.DetachRolePolicy call took %s", time.Since(start))
		return output, err
	default:
		return nil, errors.New("missing one of 'user, group, role' param")
	}
}

type policyBody struct {
	Version   string
	Statement []*policyStatement
}

type policyStatement struct {
	Effect     string           `json:",omitempty"`
	Actions    []string         `json:"Action,omitempty"`
	Resources  []string         `json:"Resource,omitempty"`
	Principal  *principal       `json:",omitempty"`
	Conditions policyConditions `json:"Condition,omitempty"`
}

type principal struct {
	AWS     interface{} `json:",omitempty"`
	Service interface{} `json:",omitempty"`
}

type policyCondition struct {
	Type  string
	Key   string
	Value string
}

func buildStatementFromParams(effect *string, resource, action, condition []*string) (*policyStatement, error) {
	stat := &policyStatement{Effect: strings.Title(StringValue(effect))}
	if resource != nil {
		res := castStringSlice(resource)
		if len(res) == 1 && res[0] == "all" {
			res[0] = "*"
		}
		stat.Resources = res
	}

	if action != nil {
		stat.Actions = castStringSlice(action)
	}
	if condition != nil {
		condStr := castStringSlice(condition)
		for _, str := range condStr {
			cond, err := parseCondition(str)
			if err != nil {
				return stat, err
			}
			stat.Conditions = append(stat.Conditions, cond)
		}
	}
	return stat, nil
}

type policyConditions []*policyCondition

func (c *policyConditions) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("\"\""), nil
	}
	var buff bytes.Buffer
	buff.WriteRune('{')
	for i, cond := range *c {
		buff.WriteString(fmt.Sprintf("\"%s\":{\"%s\":\"%s\"}", cond.Type, cond.Key, cond.Value))
		if i < len(*c)-1 {
			buff.WriteRune(',')
		}
	}
	buff.WriteRune('}')
	return buff.Bytes(), nil
}

var conditionRegex = regexp.MustCompile("^([a-zA-Z0-9:_\\-\\[\\]\\*]+)(==|!=|=~|!~|<=|>=|<|>)(.*)$")

func parseCondition(condition string) (*policyCondition, error) {
	matches := conditionRegex.FindStringSubmatch(condition)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid condition '%s'", condition)
	}
	key, operator, value := matches[1], matches[2], matches[3]
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") && len(value) >= 2 {
		value = value[1 : len(value)-1]
	}
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) >= 2 {
		value = value[1 : len(value)-1]
	}

	if strings.ToLower(value) == "null" {
		switch operator {
		case "==":
			return &policyCondition{Type: "Null", Key: key, Value: "true"}, nil
		case "!=":
			return &policyCondition{Type: "Null", Key: key, Value: "false"}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for null value '%s', expected either '==' or '!='", operator, value)
		}
	} else if strings.HasPrefix(value, "arn:") {
		switch operator {
		case "==":
			return &policyCondition{Type: "ArnEquals", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "ArnNotEquals", Key: key, Value: value}, nil
		case "=~":
			return &policyCondition{Type: "ArnLike", Key: key, Value: value}, nil
		case "!~":
			return &policyCondition{Type: "ArnNotLike", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for arn value '%s', expected either '==', '!=', '=~' or '!~'", operator, value)
		}
	} else if _, _, cidrErr := net.ParseCIDR(value); cidrErr == nil || net.ParseIP(value) != nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "IpAddress", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "NotIpAddress", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for IP value '%s', expected either '==' or '!='", operator, value)
		}
	} else if _, err := time.Parse("2006-01-02T15:04:05Z", value); err == nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "DateEquals", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "DateNotEquals", Key: key, Value: value}, nil
		case "<":
			return &policyCondition{Type: "DateLessThan", Key: key, Value: value}, nil
		case "<=":
			return &policyCondition{Type: "DateLessThanEquals", Key: key, Value: value}, nil
		case ">":
			return &policyCondition{Type: "DateGreaterThan", Key: key, Value: value}, nil
		case ">=":
			return &policyCondition{Type: "DateGreaterThanEquals", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for date value '%s', expected either '==', '!=', '>', '>=', '<' or '<='", operator, value)
		}
	} else if _, err := strconv.Atoi(value); err == nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "NumericEquals", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "NumericNotEquals", Key: key, Value: value}, nil
		case "<":
			return &policyCondition{Type: "NumericLessThan", Key: key, Value: value}, nil
		case "<=":
			return &policyCondition{Type: "NumericLessThanEquals", Key: key, Value: value}, nil
		case ">":
			return &policyCondition{Type: "NumericGreaterThan", Key: key, Value: value}, nil
		case ">=":
			return &policyCondition{Type: "NumericGreaterThanEquals", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for int value '%s', expected either '==', '!=', '>', '>=', '<' or '<='", operator, value)
		}
	} else if b, err := strconv.ParseBool(value); err == nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "Bool", Key: key, Value: fmt.Sprint(b)}, nil
		case "!=":
			return &policyCondition{Type: "Bool", Key: key, Value: fmt.Sprint(!b)}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for bool value '%s', expected either '==' or '!='", operator, value)
		}
	} else if _, err := base64.StdEncoding.DecodeString(value); value != "" && err == nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "BinaryEquals", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for binary value '%s', expected '=='", operator, value)
		}
	} else {
		switch operator {
		case "==":
			return &policyCondition{Type: "StringEquals", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "StringNotEquals", Key: key, Value: value}, nil
		case "=~":
			return &policyCondition{Type: "StringLike", Key: key, Value: value}, nil
		case "!~":
			return &policyCondition{Type: "StringNotLike", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for string value '%s', expected either '==', '!=', '=~' or '!~'", operator, value)
		}
	}
}

func lookupAWSPolicy(service, access string) (*policy, error) {
	if access != "readonly" && access != "full" {
		return nil, errors.New("looking up AWS policies: access value can only be 'readonly' or 'full'")
	}

	var suggestions []string
	for _, p := range awsPolicies {
		name := strings.ToLower(p.Name)
		match := fmt.Sprintf("%s%s", strings.ToLower(service), strings.ToLower(access))
		if strings.Contains(name, match) {
			return p, nil
		}
		if strings.Contains(name, strings.ToLower(service)) {
			suggestions = append(suggestions, fmt.Sprintf("\t\tarn=%s", p.Arn))
		}
	}

	errBuff := bytes.NewBufferString(fmt.Sprintf("No AWS policy matching service '%s' and access '%s'", service, access))
	if len(suggestions) > 0 {
		errBuff.WriteString(". Try using the full ARN of those potential matches:\n")
		errBuff.WriteString(strings.Join(suggestions, "\n"))
	}

	return nil, errors.New(errBuff.String())
}

type policy struct {
	Name string `json:"PolicyName"`
	Id   string `json:"PolicyId"`
	Arn  string `json:"Arn"`
}

var awsPolicies = []*policy{
	{
		Name: "AWSDirectConnectReadOnlyAccess",
		Id:   "ANPAI23HZ27SI6FQMGNQ2",
		Arn:  "arn:aws:iam::aws:policy/AWSDirectConnectReadOnlyAccess",
	},
	{
		Name: "AmazonGlacierReadOnlyAccess",
		Id:   "ANPAI2D5NJKMU274MET4E",
		Arn:  "arn:aws:iam::aws:policy/AmazonGlacierReadOnlyAccess",
	},
	{
		Name: "AWSMarketplaceFullAccess",
		Id:   "ANPAI2DV5ULJSO2FYVPYG",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceFullAccess",
	},
	{
		Name: "AutoScalingConsoleReadOnlyAccess",
		Id:   "ANPAI3A7GDXOYQV3VUQMK",
		Arn:  "arn:aws:iam::aws:policy/AutoScalingConsoleReadOnlyAccess",
	},
	{
		Name: "AmazonDMSRedshiftS3Role",
		Id:   "ANPAI3CCUQ4U5WNC5F6B6",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonDMSRedshiftS3Role",
	},
	{
		Name: "AWSQuickSightListIAM",
		Id:   "ANPAI3CH5UUWZN4EKGILO",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSQuickSightListIAM",
	},
	{
		Name: "AWSHealthFullAccess",
		Id:   "ANPAI3CUMPCPEUPCSXC4Y",
		Arn:  "arn:aws:iam::aws:policy/AWSHealthFullAccess",
	},
	{
		Name: "AmazonRDSFullAccess",
		Id:   "ANPAI3R4QMOG6Q5A4VWVG",
		Arn:  "arn:aws:iam::aws:policy/AmazonRDSFullAccess",
	},
	{
		Name: "SupportUser",
		Id:   "ANPAI3V4GSSN5SJY3P2RO",
		Arn:  "arn:aws:iam::aws:policy/job-function/SupportUser",
	},
	{
		Name: "AmazonEC2FullAccess",
		Id:   "ANPAI3VAJF5ZCRZ7MCQE6",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2FullAccess",
	},
	{
		Name: "AWSElasticBeanstalkReadOnlyAccess",
		Id:   "ANPAI47KNGXDAXFD4SDHG",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkReadOnlyAccess",
	},
	{
		Name: "AWSCertificateManagerReadOnly",
		Id:   "ANPAI4GSWX6S4MESJ3EWC",
		Arn:  "arn:aws:iam::aws:policy/AWSCertificateManagerReadOnly",
	},
	{
		Name: "AWSQuicksightAthenaAccess",
		Id:   "ANPAI4JB77JXFQXDWNRPM",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSQuicksightAthenaAccess",
	},
	{
		Name: "AWSCodeCommitPowerUser",
		Id:   "ANPAI4UIINUVGB5SEC57G",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeCommitPowerUser",
	},
	{
		Name: "AWSCodeCommitFullAccess",
		Id:   "ANPAI4VCZ3XPIZLQ5NZV2",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeCommitFullAccess",
	},
	{
		Name: "IAMSelfManageServiceSpecificCredentials",
		Id:   "ANPAI4VT74EMXK2PMQJM2",
		Arn:  "arn:aws:iam::aws:policy/IAMSelfManageServiceSpecificCredentials",
	},
	{
		Name: "AmazonSQSFullAccess",
		Id:   "ANPAI65L554VRJ33ECQS6",
		Arn:  "arn:aws:iam::aws:policy/AmazonSQSFullAccess",
	},
	{
		Name: "AWSLambdaFullAccess",
		Id:   "ANPAI6E2CYYMI4XI7AA5K",
		Arn:  "arn:aws:iam::aws:policy/AWSLambdaFullAccess",
	},
	{
		Name: "AWSIoTLogging",
		Id:   "ANPAI6R6Z2FHHGS454W7W",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSIoTLogging",
	},
	{
		Name: "AmazonEC2RoleforSSM",
		Id:   "ANPAI6TL3SMY22S4KMMX6",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM",
	},
	{
		Name: "AWSCloudHSMRole",
		Id:   "ANPAI7QIUU4GC66SF26WE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSCloudHSMRole",
	},
	{
		Name: "IAMFullAccess",
		Id:   "ANPAI7XKCFMBPM3QQRRVQ",
		Arn:  "arn:aws:iam::aws:policy/IAMFullAccess",
	},
	{
		Name: "AmazonInspectorFullAccess",
		Id:   "ANPAI7Y6NTA27NWNA5U5E",
		Arn:  "arn:aws:iam::aws:policy/AmazonInspectorFullAccess",
	},
	{
		Name: "AmazonElastiCacheFullAccess",
		Id:   "ANPAIA2V44CPHAUAAECKG",
		Arn:  "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess",
	},
	{
		Name: "AWSAgentlessDiscoveryService",
		Id:   "ANPAIA3DIL7BYQ35ISM4K",
		Arn:  "arn:aws:iam::aws:policy/AWSAgentlessDiscoveryService",
	},
	{
		Name: "AWSXrayWriteOnlyAccess",
		Id:   "ANPAIAACM4LMYSRGBCTM6",
		Arn:  "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess",
	},
	{
		Name: "AutoScalingReadOnlyAccess",
		Id:   "ANPAIAFWUVLC2LPLSFTFG",
		Arn:  "arn:aws:iam::aws:policy/AutoScalingReadOnlyAccess",
	},
	{
		Name: "AutoScalingFullAccess",
		Id:   "ANPAIAWRCSJDDXDXGPCFU",
		Arn:  "arn:aws:iam::aws:policy/AutoScalingFullAccess",
	},
	{
		Name: "AmazonEC2RoleforAWSCodeDeploy",
		Id:   "ANPAIAZKXZ27TAJ4PVWGK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforAWSCodeDeploy",
	},
	{
		Name: "AWSMobileHub_ReadOnly",
		Id:   "ANPAIBXVYVL3PWQFBZFGW",
		Arn:  "arn:aws:iam::aws:policy/AWSMobileHub_ReadOnly",
	},
	{
		Name: "CloudWatchEventsBuiltInTargetExecutionAccess",
		Id:   "ANPAIC5AQ5DATYSNF4AUM",
		Arn:  "arn:aws:iam::aws:policy/service-role/CloudWatchEventsBuiltInTargetExecutionAccess",
	},
	{
		Name: "AmazonCloudDirectoryReadOnlyAccess",
		Id:   "ANPAICMSZQGR3O62KMD6M",
		Arn:  "arn:aws:iam::aws:policy/AmazonCloudDirectoryReadOnlyAccess",
	},
	{
		Name: "AWSOpsWorksFullAccess",
		Id:   "ANPAICN26VXMXASXKOQCG",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksFullAccess",
	},
	{
		Name: "AWSOpsWorksCMInstanceProfileRole",
		Id:   "ANPAICSU3OSHCURP2WIZW",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksCMInstanceProfileRole",
	},
	{
		Name: "AWSCodePipelineApproverAccess",
		Id:   "ANPAICXNWK42SQ6LMDXM2",
		Arn:  "arn:aws:iam::aws:policy/AWSCodePipelineApproverAccess",
	},
	{
		Name: "AWSApplicationDiscoveryAgentAccess",
		Id:   "ANPAICZIOVAGC6JPF3WHC",
		Arn:  "arn:aws:iam::aws:policy/AWSApplicationDiscoveryAgentAccess",
	},
	{
		Name: "ViewOnlyAccess",
		Id:   "ANPAID22R6XPJATWOFDK6",
		Arn:  "arn:aws:iam::aws:policy/job-function/ViewOnlyAccess",
	},
	{
		Name: "AmazonElasticMapReduceRole",
		Id:   "ANPAIDI2BQT2LKXZG36TW",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceRole",
	},
	{
		Name: "AmazonRoute53DomainsReadOnlyAccess",
		Id:   "ANPAIDRINP6PPTRXYVQCI",
		Arn:  "arn:aws:iam::aws:policy/AmazonRoute53DomainsReadOnlyAccess",
	},
	{
		Name: "AWSOpsWorksRole",
		Id:   "ANPAIDUTMOKHJFAPJV45W",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSOpsWorksRole",
	},
	{
		Name: "ApplicationAutoScalingForAmazonAppStreamAccess",
		Id:   "ANPAIEL3HJCCWFVHA6KPG",
		Arn:  "arn:aws:iam::aws:policy/service-role/ApplicationAutoScalingForAmazonAppStreamAccess",
	},
	{
		Name: "AmazonEC2ContainerRegistryFullAccess",
		Id:   "ANPAIESRL7KD7IIVF6V4W",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess",
	},
	{
		Name: "SimpleWorkflowFullAccess",
		Id:   "ANPAIFE3AV6VE7EANYBVM",
		Arn:  "arn:aws:iam::aws:policy/SimpleWorkflowFullAccess",
	},
	{
		Name: "AmazonS3FullAccess",
		Id:   "ANPAIFIR6V6BVTRAHWINE",
		Arn:  "arn:aws:iam::aws:policy/AmazonS3FullAccess",
	},
	{
		Name: "AWSStorageGatewayReadOnlyAccess",
		Id:   "ANPAIFKCTUVOPD5NICXJK",
		Arn:  "arn:aws:iam::aws:policy/AWSStorageGatewayReadOnlyAccess",
	},
	{
		Name: "Billing",
		Id:   "ANPAIFTHXT6FFMIRT7ZEA",
		Arn:  "arn:aws:iam::aws:policy/job-function/Billing",
	},
	{
		Name: "QuickSightAccessForS3StorageManagementAnalyticsReadOnly",
		Id:   "ANPAIFWG3L3WDMR4I7ZJW",
		Arn:  "arn:aws:iam::aws:policy/service-role/QuickSightAccessForS3StorageManagementAnalyticsReadOnly",
	},
	{
		Name: "AmazonEC2ContainerRegistryReadOnly",
		Id:   "ANPAIFYZPA37OOHVIH7KQ",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
	},
	{
		Name: "AmazonElasticMapReduceforEC2Role",
		Id:   "ANPAIGALS5RCDLZLB3PGS",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceforEC2Role",
	},
	{
		Name: "DatabaseAdministrator",
		Id:   "ANPAIGBMAW4VUQKOQNVT6",
		Arn:  "arn:aws:iam::aws:policy/job-function/DatabaseAdministrator",
	},
	{
		Name: "AmazonRedshiftReadOnlyAccess",
		Id:   "ANPAIGD46KSON64QBSEZM",
		Arn:  "arn:aws:iam::aws:policy/AmazonRedshiftReadOnlyAccess",
	},
	{
		Name: "AmazonEC2ReadOnlyAccess",
		Id:   "ANPAIGDT4SV4GSETWTBZK",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess",
	},
	{
		Name: "AWSXrayReadOnlyAccess",
		Id:   "ANPAIH4OFXWPS6ZX6OPGQ",
		Arn:  "arn:aws:iam::aws:policy/AWSXrayReadOnlyAccess",
	},
	{
		Name: "AWSElasticBeanstalkEnhancedHealth",
		Id:   "ANPAIH5EFJNMOGUUTKLFE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSElasticBeanstalkEnhancedHealth",
	},
	{
		Name: "AmazonElasticMapReduceReadOnlyAccess",
		Id:   "ANPAIHP6NH2S6GYFCOINC",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticMapReduceReadOnlyAccess",
	},
	{
		Name: "AWSDirectoryServiceReadOnlyAccess",
		Id:   "ANPAIHWYO6WSDNCG64M2W",
		Arn:  "arn:aws:iam::aws:policy/AWSDirectoryServiceReadOnlyAccess",
	},
	{
		Name: "AmazonVPCReadOnlyAccess",
		Id:   "ANPAIICZJNOJN36GTG6CM",
		Arn:  "arn:aws:iam::aws:policy/AmazonVPCReadOnlyAccess",
	},
	{
		Name: "CloudWatchEventsReadOnlyAccess",
		Id:   "ANPAIILJPXXA6F7GYLYBS",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchEventsReadOnlyAccess",
	},
	{
		Name: "AmazonAPIGatewayInvokeFullAccess",
		Id:   "ANPAIIWAX2NOOQJ4AIEQ6",
		Arn:  "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess",
	},
	{
		Name: "AmazonKinesisAnalyticsReadOnly",
		Id:   "ANPAIJIEXZAFUK43U7ARK",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisAnalyticsReadOnly",
	},
	{
		Name: "AmazonMobileAnalyticsFullAccess",
		Id:   "ANPAIJIKLU2IJ7WJ6DZFG",
		Arn:  "arn:aws:iam::aws:policy/AmazonMobileAnalyticsFullAccess",
	},
	{
		Name: "AWSMobileHub_FullAccess",
		Id:   "ANPAIJLU43R6AGRBK76DM",
		Arn:  "arn:aws:iam::aws:policy/AWSMobileHub_FullAccess",
	},
	{
		Name: "AmazonAPIGatewayPushToCloudWatchLogs",
		Id:   "ANPAIK4GFO7HLKYN64ASK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs",
	},
	{
		Name: "AWSDataPipelineRole",
		Id:   "ANPAIKCP6XS3ESGF4GLO2",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSDataPipelineRole",
	},
	{
		Name: "CloudWatchFullAccess",
		Id:   "ANPAIKEABORKUXN6DEAZU",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchFullAccess",
	},
	{
		Name: "ServiceCatalogAdminFullAccess",
		Id:   "ANPAIKTX42IAS75B7B7BY",
		Arn:  "arn:aws:iam::aws:policy/ServiceCatalogAdminFullAccess",
	},
	{
		Name: "AmazonRDSDirectoryServiceAccess",
		Id:   "ANPAIL4KBY57XWMYUHKUU",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonRDSDirectoryServiceAccess",
	},
	{
		Name: "AWSCodePipelineReadOnlyAccess",
		Id:   "ANPAILFKZXIBOTNC5TO2Q",
		Arn:  "arn:aws:iam::aws:policy/AWSCodePipelineReadOnlyAccess",
	},
	{
		Name: "ReadOnlyAccess",
		Id:   "ANPAILL3HVNFSB6DCOWYQ",
		Arn:  "arn:aws:iam::aws:policy/ReadOnlyAccess",
	},
	{
		Name: "AmazonMachineLearningBatchPredictionsAccess",
		Id:   "ANPAILOI4HTQSFTF3GQSC",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningBatchPredictionsAccess",
	},
	{
		Name: "AmazonRekognitionReadOnlyAccess",
		Id:   "ANPAILWSUHXUY4ES43SA4",
		Arn:  "arn:aws:iam::aws:policy/AmazonRekognitionReadOnlyAccess",
	},
	{
		Name: "AWSCodeDeployReadOnlyAccess",
		Id:   "ANPAILZHHKCKB4NE7XOIQ",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeDeployReadOnlyAccess",
	},
	{
		Name: "CloudSearchFullAccess",
		Id:   "ANPAIM6OOWKQ7L7VBOZOC",
		Arn:  "arn:aws:iam::aws:policy/CloudSearchFullAccess",
	},
	{
		Name: "AWSCloudHSMFullAccess",
		Id:   "ANPAIMBQYQZM7F63DA2UU",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudHSMFullAccess",
	},
	{
		Name: "AmazonEC2SpotFleetAutoscaleRole",
		Id:   "ANPAIMFFRMIOBGDP2TAVE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetAutoscaleRole",
	},
	{
		Name: "AWSCodeBuildDeveloperAccess",
		Id:   "ANPAIMKTMR34XSBQW45HS",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeBuildDeveloperAccess",
	},
	{
		Name: "AmazonEC2SpotFleetRole",
		Id:   "ANPAIMRTKHWK7ESSNETSW",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetRole",
	},
	{
		Name: "AWSDataPipeline_PowerUser",
		Id:   "ANPAIMXGLVY6DVR24VTYS",
		Arn:  "arn:aws:iam::aws:policy/AWSDataPipeline_PowerUser",
	},
	{
		Name: "AmazonElasticTranscoderJobsSubmitter",
		Id:   "ANPAIN5WGARIKZ3E2UQOU",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticTranscoderJobsSubmitter",
	},
	{
		Name: "AWSCodeStarServiceRole",
		Id:   "ANPAIN6D4M2KD3NBOC4M4",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSCodeStarServiceRole",
	},
	{
		Name: "AWSDirectoryServiceFullAccess",
		Id:   "ANPAINAW5ANUWTH3R4ANI",
		Arn:  "arn:aws:iam::aws:policy/AWSDirectoryServiceFullAccess",
	},
	{
		Name: "AmazonDynamoDBFullAccess",
		Id:   "ANPAINUGF2JSOSUY76KYA",
		Arn:  "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
	},
	{
		Name: "AmazonSESReadOnlyAccess",
		Id:   "ANPAINV2XPFRMWJJNSCGI",
		Arn:  "arn:aws:iam::aws:policy/AmazonSESReadOnlyAccess",
	},
	{
		Name: "AWSWAFReadOnlyAccess",
		Id:   "ANPAINZVDMX2SBF7EU2OC",
		Arn:  "arn:aws:iam::aws:policy/AWSWAFReadOnlyAccess",
	},
	{
		Name: "AutoScalingNotificationAccessRole",
		Id:   "ANPAIO2VMUPGDC5PZVXVA",
		Arn:  "arn:aws:iam::aws:policy/service-role/AutoScalingNotificationAccessRole",
	},
	{
		Name: "AmazonMechanicalTurkReadOnly",
		Id:   "ANPAIO5IY3G3WXSX5PPRM",
		Arn:  "arn:aws:iam::aws:policy/AmazonMechanicalTurkReadOnly",
	},
	{
		Name: "AmazonKinesisReadOnlyAccess",
		Id:   "ANPAIOCMTDT5RLKZ2CAJO",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisReadOnlyAccess",
	},
	{
		Name: "AWSCodeDeployFullAccess",
		Id:   "ANPAIONKN3TJZUKXCHXWC",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeDeployFullAccess",
	},
	{
		Name: "CloudWatchActionsEC2Access",
		Id:   "ANPAIOWD4E3FVSORSZTGU",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchActionsEC2Access",
	},
	{
		Name: "AWSLambdaDynamoDBExecutionRole",
		Id:   "ANPAIP7WNAGMIPYNW4WQG",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaDynamoDBExecutionRole",
	},
	{
		Name: "AmazonRoute53DomainsFullAccess",
		Id:   "ANPAIPAFBMIYUILMOKL6G",
		Arn:  "arn:aws:iam::aws:policy/AmazonRoute53DomainsFullAccess",
	},
	{
		Name: "AmazonElastiCacheReadOnlyAccess",
		Id:   "ANPAIPDACSNQHSENWAKM2",
		Arn:  "arn:aws:iam::aws:policy/AmazonElastiCacheReadOnlyAccess",
	},
	{
		Name: "AmazonAthenaFullAccess",
		Id:   "ANPAIPJMLMD4C7RYZ6XCK",
		Arn:  "arn:aws:iam::aws:policy/AmazonAthenaFullAccess",
	},
	{
		Name: "AmazonElasticFileSystemReadOnlyAccess",
		Id:   "ANPAIPN5S4NE5JJOKVC4Y",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticFileSystemReadOnlyAccess",
	},
	{
		Name: "CloudFrontFullAccess",
		Id:   "ANPAIPRV52SH6HDCCFY6U",
		Arn:  "arn:aws:iam::aws:policy/CloudFrontFullAccess",
	},
	{
		Name: "AmazonMachineLearningRoleforRedshiftDataSource",
		Id:   "ANPAIQ5UDYYMNN42BM4AK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonMachineLearningRoleforRedshiftDataSource",
	},
	{
		Name: "AmazonMobileAnalyticsNon-financialReportAccess",
		Id:   "ANPAIQLKQ4RXPUBBVVRDE",
		Arn:  "arn:aws:iam::aws:policy/AmazonMobileAnalyticsNon-financialReportAccess",
	},
	{
		Name: "AWSCloudTrailFullAccess",
		Id:   "ANPAIQNUJTQYDRJPC3BNK",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudTrailFullAccess",
	},
	{
		Name: "AmazonCognitoDeveloperAuthenticatedIdentities",
		Id:   "ANPAIQOKZ5BGKLCMTXH4W",
		Arn:  "arn:aws:iam::aws:policy/AmazonCognitoDeveloperAuthenticatedIdentities",
	},
	{
		Name: "AWSConfigRole",
		Id:   "ANPAIQRXRDRGJUA33ELIO",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSConfigRole",
	},
	{
		Name: "AmazonAppStreamServiceAccess",
		Id:   "ANPAISBRZ7LMMCBYEF3SE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonAppStreamServiceAccess",
	},
	{
		Name: "AmazonRedshiftFullAccess",
		Id:   "ANPAISEKCHH4YDB46B5ZO",
		Arn:  "arn:aws:iam::aws:policy/AmazonRedshiftFullAccess",
	},
	{
		Name: "AmazonZocaloReadOnlyAccess",
		Id:   "ANPAISRCSSJNS3QPKZJPM",
		Arn:  "arn:aws:iam::aws:policy/AmazonZocaloReadOnlyAccess",
	},
	{
		Name: "AWSCloudHSMReadOnlyAccess",
		Id:   "ANPAISVCBSY7YDBOT67KE",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudHSMReadOnlyAccess",
	},
	{
		Name: "SystemAdministrator",
		Id:   "ANPAITJPEZXCYCBXANDSW",
		Arn:  "arn:aws:iam::aws:policy/job-function/SystemAdministrator",
	},
	{
		Name: "AmazonEC2ContainerServiceEventsRole",
		Id:   "ANPAITKFNIUAG27VSYNZ4",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceEventsRole",
	},
	{
		Name: "AmazonRoute53ReadOnlyAccess",
		Id:   "ANPAITOYK2ZAOQFXV2JNC",
		Arn:  "arn:aws:iam::aws:policy/AmazonRoute53ReadOnlyAccess",
	},
	{
		Name: "AmazonEC2ReportsAccess",
		Id:   "ANPAIU6NBZVF2PCRW36ZW",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ReportsAccess",
	},
	{
		Name: "AmazonEC2ContainerServiceAutoscaleRole",
		Id:   "ANPAIUAP3EGGGXXCPDQKK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceAutoscaleRole",
	},
	{
		Name: "AWSBatchServiceRole",
		Id:   "ANPAIUETIXPCKASQJURFE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSBatchServiceRole",
	},
	{
		Name: "AWSElasticBeanstalkWebTier",
		Id:   "ANPAIUF4325SJYOREKW3A",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkWebTier",
	},
	{
		Name: "AmazonSQSReadOnlyAccess",
		Id:   "ANPAIUGSSQY362XGCM6KW",
		Arn:  "arn:aws:iam::aws:policy/AmazonSQSReadOnlyAccess",
	},
	{
		Name: "AWSMobileHub_ServiceUseOnly",
		Id:   "ANPAIUHPQXBDZUWOP3PSK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSMobileHub_ServiceUseOnly",
	},
	{
		Name: "AmazonKinesisFullAccess",
		Id:   "ANPAIVF32HAMOXCUYRAYE",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisFullAccess",
	},
	{
		Name: "AmazonMachineLearningReadOnlyAccess",
		Id:   "ANPAIW5VYBCGEX56JCINC",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningReadOnlyAccess",
	},
	{
		Name: "AmazonRekognitionFullAccess",
		Id:   "ANPAIWDAOK6AIFDVX6TT6",
		Arn:  "arn:aws:iam::aws:policy/AmazonRekognitionFullAccess",
	},
	{
		Name: "RDSCloudHsmAuthorizationRole",
		Id:   "ANPAIWKFXRLQG2ROKKXLE",
		Arn:  "arn:aws:iam::aws:policy/service-role/RDSCloudHsmAuthorizationRole",
	},
	{
		Name: "AmazonMachineLearningFullAccess",
		Id:   "ANPAIWKW6AGSGYOQ5ERHC",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningFullAccess",
	},
	{
		Name: "AdministratorAccess",
		Id:   "ANPAIWMBCKSKIEE64ZLYK",
		Arn:  "arn:aws:iam::aws:policy/AdministratorAccess",
	},
	{
		Name: "AmazonMachineLearningRealTimePredictionOnlyAccess",
		Id:   "ANPAIWMCNQPRWMWT36GVQ",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningRealTimePredictionOnlyAccess",
	},
	{
		Name: "AWSConfigUserAccess",
		Id:   "ANPAIWTTSFJ7KKJE3MWGA",
		Arn:  "arn:aws:iam::aws:policy/AWSConfigUserAccess",
	},
	{
		Name: "AWSIoTConfigAccess",
		Id:   "ANPAIWWGD4LM4EMXNRL7I",
		Arn:  "arn:aws:iam::aws:policy/AWSIoTConfigAccess",
	},
	{
		Name: "SecurityAudit",
		Id:   "ANPAIX2T3QCXHR2OGGCTO",
		Arn:  "arn:aws:iam::aws:policy/SecurityAudit",
	},
	{
		Name: "AWSCodeStarFullAccess",
		Id:   "ANPAIXI233TFUGLZOJBEC",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeStarFullAccess",
	},
	{
		Name: "AWSDataPipeline_FullAccess",
		Id:   "ANPAIXOFIG7RSBMRPHXJ4",
		Arn:  "arn:aws:iam::aws:policy/AWSDataPipeline_FullAccess",
	},
	{
		Name: "AmazonDynamoDBReadOnlyAccess",
		Id:   "ANPAIY2XFNA232XJ6J7X2",
		Arn:  "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
	},
	{
		Name: "AutoScalingConsoleFullAccess",
		Id:   "ANPAIYEN6FJGYYWJFFCZW",
		Arn:  "arn:aws:iam::aws:policy/AutoScalingConsoleFullAccess",
	},
	{
		Name: "AmazonSNSReadOnlyAccess",
		Id:   "ANPAIZGQCQTFOFPMHSB6W",
		Arn:  "arn:aws:iam::aws:policy/AmazonSNSReadOnlyAccess",
	},
	{
		Name: "AmazonElasticMapReduceFullAccess",
		Id:   "ANPAIZP5JFP3AMSGINBB2",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticMapReduceFullAccess",
	},
	{
		Name: "AmazonS3ReadOnlyAccess",
		Id:   "ANPAIZTJ4DXE7G6AGAE6M",
		Arn:  "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
	},
	{
		Name: "AWSElasticBeanstalkFullAccess",
		Id:   "ANPAIZYX2YLLBW2LJVUFW",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkFullAccess",
	},
	{
		Name: "AmazonWorkSpacesAdmin",
		Id:   "ANPAJ26AU6ATUQCT5KVJU",
		Arn:  "arn:aws:iam::aws:policy/AmazonWorkSpacesAdmin",
	},
	{
		Name: "AWSCodeDeployRole",
		Id:   "ANPAJ2NKMKD73QS5NBFLA",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole",
	},
	{
		Name: "AmazonSESFullAccess",
		Id:   "ANPAJ2P4NXCHAT7NDPNR4",
		Arn:  "arn:aws:iam::aws:policy/AmazonSESFullAccess",
	},
	{
		Name: "CloudWatchLogsReadOnlyAccess",
		Id:   "ANPAJ2YIYDYSNNEHK3VKW",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchLogsReadOnlyAccess",
	},
	{
		Name: "AmazonKinesisFirehoseReadOnlyAccess",
		Id:   "ANPAJ36NT645INW4K24W6",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisFirehoseReadOnlyAccess",
	},
	{
		Name: "AWSOpsWorksRegisterCLI",
		Id:   "ANPAJ3AB5ZBFPCQGTVDU4",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksRegisterCLI",
	},
	{
		Name: "AmazonDynamoDBFullAccesswithDataPipeline",
		Id:   "ANPAJ3ORT7KDISSXGHJXA",
		Arn:  "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccesswithDataPipeline",
	},
	{
		Name: "AmazonEC2RoleforDataPipelineRole",
		Id:   "ANPAJ3Z5I2WAJE5DN2J36",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforDataPipelineRole",
	},
	{
		Name: "CloudWatchLogsFullAccess",
		Id:   "ANPAJ3ZGNWK2R5HW5BQFO",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess",
	},
	{
		Name: "AWSElasticBeanstalkMulticontainerDocker",
		Id:   "ANPAJ45SBYG72SD6SHJEY",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkMulticontainerDocker",
	},
	{
		Name: "AmazonElasticTranscoderFullAccess",
		Id:   "ANPAJ4D5OJU75P5ZJZVNY",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticTranscoderFullAccess",
	},
	{
		Name: "IAMUserChangePassword",
		Id:   "ANPAJ4L4MM2A7QIEB56MS",
		Arn:  "arn:aws:iam::aws:policy/IAMUserChangePassword",
	},
	{
		Name: "AmazonAPIGatewayAdministrator",
		Id:   "ANPAJ4PT6VY5NLKTNUYSI",
		Arn:  "arn:aws:iam::aws:policy/AmazonAPIGatewayAdministrator",
	},
	{
		Name: "ServiceCatalogEndUserAccess",
		Id:   "ANPAJ56OMCO72RI4J5FSA",
		Arn:  "arn:aws:iam::aws:policy/ServiceCatalogEndUserAccess",
	},
	{
		Name: "AmazonPollyReadOnlyAccess",
		Id:   "ANPAJ5FENL3CVPL2FPDLA",
		Arn:  "arn:aws:iam::aws:policy/AmazonPollyReadOnlyAccess",
	},
	{
		Name: "AmazonMobileAnalyticsWriteOnlyAccess",
		Id:   "ANPAJ5TAWBBQC2FAL3G6G",
		Arn:  "arn:aws:iam::aws:policy/AmazonMobileAnalyticsWriteOnlyAccess",
	},
	{
		Name: "AmazonEC2SpotFleetTaggingRole",
		Id:   "ANPAJ5U6UMLCEYLX5OLC4",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetTaggingRole",
	},
	{
		Name: "DataScientist",
		Id:   "ANPAJ5YHI2BQW7EQFYDXS",
		Arn:  "arn:aws:iam::aws:policy/job-function/DataScientist",
	},
	{
		Name: "AWSMarketplaceMeteringFullAccess",
		Id:   "ANPAJ65YJPG7CC7LDXNA6",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceMeteringFullAccess",
	},
	{
		Name: "AWSOpsWorksCMServiceRole",
		Id:   "ANPAJ6I6MPGJE62URSHCO",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSOpsWorksCMServiceRole",
	},
	{
		Name: "AWSConnector",
		Id:   "ANPAJ6YATONJHICG3DJ3U",
		Arn:  "arn:aws:iam::aws:policy/AWSConnector",
	},
	{
		Name: "AWSBatchFullAccess",
		Id:   "ANPAJ7K2KIWB3HZVK3CUO",
		Arn:  "arn:aws:iam::aws:policy/AWSBatchFullAccess",
	},
	{
		Name: "ServiceCatalogAdminReadOnlyAccess",
		Id:   "ANPAJ7XOUSS75M4LIPKO4",
		Arn:  "arn:aws:iam::aws:policy/ServiceCatalogAdminReadOnlyAccess",
	},
	{
		Name: "AmazonSSMFullAccess",
		Id:   "ANPAJA7V6HI4ISQFMDYAG",
		Arn:  "arn:aws:iam::aws:policy/AmazonSSMFullAccess",
	},
	{
		Name: "AWSCodeCommitReadOnly",
		Id:   "ANPAJACNSXR7Z2VLJW3D6",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeCommitReadOnly",
	},
	{
		Name: "AmazonEC2ContainerServiceFullAccess",
		Id:   "ANPAJALOYVTPDZEMIACSM",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ContainerServiceFullAccess",
	},
	{
		Name: "AmazonCognitoReadOnly",
		Id:   "ANPAJBFTRZD2GQGJHSVQK",
		Arn:  "arn:aws:iam::aws:policy/AmazonCognitoReadOnly",
	},
	{
		Name: "AmazonDMSCloudWatchLogsRole",
		Id:   "ANPAJBG7UXZZXUJD3TDJE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonDMSCloudWatchLogsRole",
	},
	{
		Name: "AWSApplicationDiscoveryServiceFullAccess",
		Id:   "ANPAJBNJEA6ZXM2SBOPDU",
		Arn:  "arn:aws:iam::aws:policy/AWSApplicationDiscoveryServiceFullAccess",
	},
	{
		Name: "AmazonVPCFullAccess",
		Id:   "ANPAJBWPGNOVKZD3JI2P2",
		Arn:  "arn:aws:iam::aws:policy/AmazonVPCFullAccess",
	},
	{
		Name: "AWSImportExportFullAccess",
		Id:   "ANPAJCQCT4JGTLC6722MQ",
		Arn:  "arn:aws:iam::aws:policy/AWSImportExportFullAccess",
	},
	{
		Name: "AmazonMechanicalTurkFullAccess",
		Id:   "ANPAJDGCL5BET73H5QIQC",
		Arn:  "arn:aws:iam::aws:policy/AmazonMechanicalTurkFullAccess",
	},
	{
		Name: "AmazonEC2ContainerRegistryPowerUser",
		Id:   "ANPAJDNE5PIHROIBGGDDW",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser",
	},
	{
		Name: "AmazonMachineLearningCreateOnlyAccess",
		Id:   "ANPAJDRUNIC2RYAMAT3CK",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningCreateOnlyAccess",
	},
	{
		Name: "AWSCloudTrailReadOnlyAccess",
		Id:   "ANPAJDU7KJADWBSEQ3E7S",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudTrailReadOnlyAccess",
	},
	{
		Name: "AWSLambdaExecute",
		Id:   "ANPAJE5FX7FQZSU5XAKGO",
		Arn:  "arn:aws:iam::aws:policy/AWSLambdaExecute",
	},
	{
		Name: "AWSIoTRuleActions",
		Id:   "ANPAJEZ6FS7BUZVUHMOKY",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSIoTRuleActions",
	},
	{
		Name: "AWSQuickSightDescribeRedshift",
		Id:   "ANPAJFEM6MLSLTW4ZNBW2",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSQuickSightDescribeRedshift",
	},
	{
		Name: "VMImportExportRoleForAWSConnector",
		Id:   "ANPAJFLQOOJ6F5XNX4LAW",
		Arn:  "arn:aws:iam::aws:policy/service-role/VMImportExportRoleForAWSConnector",
	},
	{
		Name: "AWSCodePipelineCustomActionAccess",
		Id:   "ANPAJFW5Z32BTVF76VCYC",
		Arn:  "arn:aws:iam::aws:policy/AWSCodePipelineCustomActionAccess",
	},
	{
		Name: "AWSOpsWorksInstanceRegistration",
		Id:   "ANPAJG3LCPVNI4WDZCIMU",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksInstanceRegistration",
	},
	{
		Name: "AmazonCloudDirectoryFullAccess",
		Id:   "ANPAJG3XQK77ATFLCF2CK",
		Arn:  "arn:aws:iam::aws:policy/AmazonCloudDirectoryFullAccess",
	},
	{
		Name: "AWSStorageGatewayFullAccess",
		Id:   "ANPAJG5SSPAVOGK3SIDGU",
		Arn:  "arn:aws:iam::aws:policy/AWSStorageGatewayFullAccess",
	},
	{
		Name: "AmazonLexReadOnly",
		Id:   "ANPAJGBI5LSMAJNDGBNAM",
		Arn:  "arn:aws:iam::aws:policy/AmazonLexReadOnly",
	},
	{
		Name: "AmazonElasticTranscoderReadOnlyAccess",
		Id:   "ANPAJGPP7GPMJRRJMEP3Q",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticTranscoderReadOnlyAccess",
	},
	{
		Name: "AWSIoTConfigReadOnlyAccess",
		Id:   "ANPAJHENEMXGX4XMFOIOI",
		Arn:  "arn:aws:iam::aws:policy/AWSIoTConfigReadOnlyAccess",
	},
	{
		Name: "AmazonWorkMailReadOnlyAccess",
		Id:   "ANPAJHF7J65E2QFKCWAJM",
		Arn:  "arn:aws:iam::aws:policy/AmazonWorkMailReadOnlyAccess",
	},
	{
		Name: "AmazonDMSVPCManagementRole",
		Id:   "ANPAJHKIGMBQI4AEFFSYO",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonDMSVPCManagementRole",
	},
	{
		Name: "AWSLambdaKinesisExecutionRole",
		Id:   "ANPAJHOLKJPXV4GBRMJUQ",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaKinesisExecutionRole",
	},
	{
		Name: "ResourceGroupsandTagEditorReadOnlyAccess",
		Id:   "ANPAJHXQTPI5I5JKAIU74",
		Arn:  "arn:aws:iam::aws:policy/ResourceGroupsandTagEditorReadOnlyAccess",
	},
	{
		Name: "AmazonSSMAutomationRole",
		Id:   "ANPAJIBQCTBCXD2XRNB6W",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonSSMAutomationRole",
	},
	{
		Name: "ServiceCatalogEndUserFullAccess",
		Id:   "ANPAJIW7AFFOONVKW75KU",
		Arn:  "arn:aws:iam::aws:policy/ServiceCatalogEndUserFullAccess",
	},
	{
		Name: "AWSStepFunctionsConsoleFullAccess",
		Id:   "ANPAJIYC52YWRX6OSMJWK",
		Arn:  "arn:aws:iam::aws:policy/AWSStepFunctionsConsoleFullAccess",
	},
	{
		Name: "AWSCodeBuildReadOnlyAccess",
		Id:   "ANPAJIZZWN6557F5HVP2K",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeBuildReadOnlyAccess",
	},
	{
		Name: "AmazonMachineLearningManageRealTimeEndpointOnlyAccess",
		Id:   "ANPAJJL3PC3VCSVZP6OCI",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningManageRealTimeEndpointOnlyAccess",
	},
	{
		Name: "CloudWatchEventsInvocationAccess",
		Id:   "ANPAJJXD6JKJLK2WDLZNO",
		Arn:  "arn:aws:iam::aws:policy/service-role/CloudWatchEventsInvocationAccess",
	},
	{
		Name: "CloudFrontReadOnlyAccess",
		Id:   "ANPAJJZMNYOTZCNQP36LG",
		Arn:  "arn:aws:iam::aws:policy/CloudFrontReadOnlyAccess",
	},
	{
		Name: "AmazonSNSRole",
		Id:   "ANPAJK5GQB7CIK7KHY2GA",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonSNSRole",
	},
	{
		Name: "AmazonMobileAnalyticsFinancialReportAccess",
		Id:   "ANPAJKJHO2R27TXKCWBU4",
		Arn:  "arn:aws:iam::aws:policy/AmazonMobileAnalyticsFinancialReportAccess",
	},
	{
		Name: "AWSElasticBeanstalkService",
		Id:   "ANPAJKQ5SN74ZQ4WASXBM",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSElasticBeanstalkService",
	},
	{
		Name: "IAMReadOnlyAccess",
		Id:   "ANPAJKSO7NDY4T57MWDSQ",
		Arn:  "arn:aws:iam::aws:policy/IAMReadOnlyAccess",
	},
	{
		Name: "AmazonRDSReadOnlyAccess",
		Id:   "ANPAJKTTTYV2IIHKLZ346",
		Arn:  "arn:aws:iam::aws:policy/AmazonRDSReadOnlyAccess",
	},
	{
		Name: "AmazonCognitoPowerUser",
		Id:   "ANPAJKW5H2HNCPGCYGR6Y",
		Arn:  "arn:aws:iam::aws:policy/AmazonCognitoPowerUser",
	},
	{
		Name: "AmazonElasticFileSystemFullAccess",
		Id:   "ANPAJKXTMNVQGIDNCKPBC",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticFileSystemFullAccess",
	},
	{
		Name: "ServerMigrationConnector",
		Id:   "ANPAJKZRWXIPK5HSG3QDQ",
		Arn:  "arn:aws:iam::aws:policy/ServerMigrationConnector",
	},
	{
		Name: "AmazonZocaloFullAccess",
		Id:   "ANPAJLCDXYRINDMUXEVL6",
		Arn:  "arn:aws:iam::aws:policy/AmazonZocaloFullAccess",
	},
	{
		Name: "AWSLambdaReadOnlyAccess",
		Id:   "ANPAJLDG7J3CGUHFN4YN6",
		Arn:  "arn:aws:iam::aws:policy/AWSLambdaReadOnlyAccess",
	},
	{
		Name: "AWSAccountUsageReportAccess",
		Id:   "ANPAJLIB4VSBVO47ZSBB6",
		Arn:  "arn:aws:iam::aws:policy/AWSAccountUsageReportAccess",
	},
	{
		Name: "AWSMarketplaceGetEntitlements",
		Id:   "ANPAJLPIMQE4WMHDC2K7C",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceGetEntitlements",
	},
	{
		Name: "AmazonEC2ContainerServiceforEC2Role",
		Id:   "ANPAJLYJCVHC7TQHCSQDS",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role",
	},
	{
		Name: "AmazonAppStreamFullAccess",
		Id:   "ANPAJLZZXU2YQVGL4QDNC",
		Arn:  "arn:aws:iam::aws:policy/AmazonAppStreamFullAccess",
	},
	{
		Name: "AWSIoTDataAccess",
		Id:   "ANPAJM2KI2UJDR24XPS2K",
		Arn:  "arn:aws:iam::aws:policy/AWSIoTDataAccess",
	},
	{
		Name: "AmazonESFullAccess",
		Id:   "ANPAJM6ZTCU24QL5PZCGC",
		Arn:  "arn:aws:iam::aws:policy/AmazonESFullAccess",
	},
	{
		Name: "ServerMigrationServiceRole",
		Id:   "ANPAJMBH3M6BO63XFW2D4",
		Arn:  "arn:aws:iam::aws:policy/service-role/ServerMigrationServiceRole",
	},
	{
		Name: "AWSWAFFullAccess",
		Id:   "ANPAJMIKIAFXZEGOLRH7C",
		Arn:  "arn:aws:iam::aws:policy/AWSWAFFullAccess",
	},
	{
		Name: "AmazonKinesisFirehoseFullAccess",
		Id:   "ANPAJMZQMTZ7FRBFHHAHI",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisFirehoseFullAccess",
	},
	{
		Name: "CloudWatchReadOnlyAccess",
		Id:   "ANPAJN23PDQP7SZQAE3QE",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchReadOnlyAccess",
	},
	{
		Name: "AWSLambdaBasicExecutionRole",
		Id:   "ANPAJNCQGXC42545SKXIK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
	},
	{
		Name: "ResourceGroupsandTagEditorFullAccess",
		Id:   "ANPAJNOS54ZFXN4T2Y34A",
		Arn:  "arn:aws:iam::aws:policy/ResourceGroupsandTagEditorFullAccess",
	},
	{
		Name: "AWSKeyManagementServicePowerUser",
		Id:   "ANPAJNPP7PPPPMJRV2SA4",
		Arn:  "arn:aws:iam::aws:policy/AWSKeyManagementServicePowerUser",
	},
	{
		Name: "AWSImportExportReadOnlyAccess",
		Id:   "ANPAJNTV4OG52ESYZHCNK",
		Arn:  "arn:aws:iam::aws:policy/AWSImportExportReadOnlyAccess",
	},
	{
		Name: "AmazonElasticTranscoderRole",
		Id:   "ANPAJNW3WMKVXFJ2KPIQ2",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonElasticTranscoderRole",
	},
	{
		Name: "AmazonEC2ContainerServiceRole",
		Id:   "ANPAJO53W2XHNACG7V77Q",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceRole",
	},
	{
		Name: "AWSDeviceFarmFullAccess",
		Id:   "ANPAJO7KEDP4VYJPNT5UW",
		Arn:  "arn:aws:iam::aws:policy/AWSDeviceFarmFullAccess",
	},
	{
		Name: "AmazonSSMReadOnlyAccess",
		Id:   "ANPAJODSKQGGJTHRYZ5FC",
		Arn:  "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess",
	},
	{
		Name: "AWSStepFunctionsReadOnlyAccess",
		Id:   "ANPAJONHB2TJQDJPFW5TM",
		Arn:  "arn:aws:iam::aws:policy/AWSStepFunctionsReadOnlyAccess",
	},
	{
		Name: "AWSMarketplaceRead-only",
		Id:   "ANPAJOOM6LETKURTJ3XZ2",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceRead-only",
	},
	{
		Name: "AWSCodePipelineFullAccess",
		Id:   "ANPAJP5LH77KSAT2KHQGG",
		Arn:  "arn:aws:iam::aws:policy/AWSCodePipelineFullAccess",
	},
	{
		Name: "AWSGreengrassResourceAccessRolePolicy",
		Id:   "ANPAJPKEIMB6YMXDEVRTM",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSGreengrassResourceAccessRolePolicy",
	},
	{
		Name: "NetworkAdministrator",
		Id:   "ANPAJPNMADZFJCVPJVZA2",
		Arn:  "arn:aws:iam::aws:policy/job-function/NetworkAdministrator",
	},
	{
		Name: "AmazonWorkSpacesApplicationManagerAdminAccess",
		Id:   "ANPAJPRL4KYETIH7XGTSS",
		Arn:  "arn:aws:iam::aws:policy/AmazonWorkSpacesApplicationManagerAdminAccess",
	},
	{
		Name: "AmazonDRSVPCManagement",
		Id:   "ANPAJPXIBTTZMBEFEX6UA",
		Arn:  "arn:aws:iam::aws:policy/AmazonDRSVPCManagement",
	},
	{
		Name: "AWSXrayFullAccess",
		Id:   "ANPAJQBYG45NSJMVQDB2K",
		Arn:  "arn:aws:iam::aws:policy/AWSXrayFullAccess",
	},
	{
		Name: "AWSElasticBeanstalkWorkerTier",
		Id:   "ANPAJQDLBRSJVKVF4JMSK",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkWorkerTier",
	},
	{
		Name: "AWSDirectConnectFullAccess",
		Id:   "ANPAJQF2QKZSK74KTIHOW",
		Arn:  "arn:aws:iam::aws:policy/AWSDirectConnectFullAccess",
	},
	{
		Name: "AWSCodeBuildAdminAccess",
		Id:   "ANPAJQJGIOIE3CD2TQXDS",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeBuildAdminAccess",
	},
	{
		Name: "AmazonKinesisAnalyticsFullAccess",
		Id:   "ANPAJQOSKHTXP43R7P5AC",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisAnalyticsFullAccess",
	},
	{
		Name: "AWSAccountActivityAccess",
		Id:   "ANPAJQRYCWMFX5J3E333K",
		Arn:  "arn:aws:iam::aws:policy/AWSAccountActivityAccess",
	},
	{
		Name: "AmazonGlacierFullAccess",
		Id:   "ANPAJQSTZJWB2AXXAKHVQ",
		Arn:  "arn:aws:iam::aws:policy/AmazonGlacierFullAccess",
	},
	{
		Name: "AmazonWorkMailFullAccess",
		Id:   "ANPAJQVKNMT7SVATQ4AUY",
		Arn:  "arn:aws:iam::aws:policy/AmazonWorkMailFullAccess",
	},
	{
		Name: "AWSMarketplaceManageSubscriptions",
		Id:   "ANPAJRDW2WIFN7QLUAKBQ",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceManageSubscriptions",
	},
	{
		Name: "AWSElasticBeanstalkCustomPlatformforEC2Role",
		Id:   "ANPAJRVFXSS6LEIQGBKDY",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkCustomPlatformforEC2Role",
	},
	{
		Name: "AWSSupportAccess",
		Id:   "ANPAJSNKQX2OW67GF4S7E",
		Arn:  "arn:aws:iam::aws:policy/AWSSupportAccess",
	},
	{
		Name: "AmazonElasticMapReduceforAutoScalingRole",
		Id:   "ANPAJSVXG6QHPE6VHDZ4Q",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceforAutoScalingRole",
	},
	{
		Name: "AWSLambdaInvocation-DynamoDB",
		Id:   "ANPAJTHQ3EKCQALQDYG5G",
		Arn:  "arn:aws:iam::aws:policy/AWSLambdaInvocation-DynamoDB",
	},
	{
		Name: "IAMUserSSHKeys",
		Id:   "ANPAJTSHUA4UXGXU7ANUA",
		Arn:  "arn:aws:iam::aws:policy/IAMUserSSHKeys",
	},
	{
		Name: "AWSIoTFullAccess",
		Id:   "ANPAJU2FPGG6PQWN72V2G",
		Arn:  "arn:aws:iam::aws:policy/AWSIoTFullAccess",
	},
	{
		Name: "AWSQuickSightDescribeRDS",
		Id:   "ANPAJU5J6OAMCJD3OO76O",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSQuickSightDescribeRDS",
	},
	{
		Name: "AWSConfigRulesExecutionRole",
		Id:   "ANPAJUB3KIKTA4PU4OYAA",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSConfigRulesExecutionRole",
	},
	{
		Name: "AmazonESReadOnlyAccess",
		Id:   "ANPAJUDMRLOQ7FPAR46FQ",
		Arn:  "arn:aws:iam::aws:policy/AmazonESReadOnlyAccess",
	},
	{
		Name: "AWSCodeDeployDeployerAccess",
		Id:   "ANPAJUWEPOMGLMVXJAPUI",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeDeployDeployerAccess",
	},
	{
		Name: "AmazonPollyFullAccess",
		Id:   "ANPAJUZOYQU6XQYPR7EWS",
		Arn:  "arn:aws:iam::aws:policy/AmazonPollyFullAccess",
	},
	{
		Name: "AmazonSSMMaintenanceWindowRole",
		Id:   "ANPAJV3JNYSTZ47VOXYME",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonSSMMaintenanceWindowRole",
	},
	{
		Name: "AmazonRDSEnhancedMonitoringRole",
		Id:   "ANPAJV7BS425S4PTSSVGK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole",
	},
	{
		Name: "AmazonLexFullAccess",
		Id:   "ANPAJVLXDHKVC23HRTKSI",
		Arn:  "arn:aws:iam::aws:policy/AmazonLexFullAccess",
	},
	{
		Name: "AWSLambdaVPCAccessExecutionRole",
		Id:   "ANPAJVTME3YLVNL72YR2K",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
	},
	{
		Name: "AmazonLexRunBotsOnly",
		Id:   "ANPAJVZGB5CM3N6YWJHBE",
		Arn:  "arn:aws:iam::aws:policy/AmazonLexRunBotsOnly",
	},
	{
		Name: "AmazonSNSFullAccess",
		Id:   "ANPAJWEKLCXXUNT2SOLSG",
		Arn:  "arn:aws:iam::aws:policy/AmazonSNSFullAccess",
	},
	{
		Name: "CloudSearchReadOnlyAccess",
		Id:   "ANPAJWPLX7N7BCC3RZLHW",
		Arn:  "arn:aws:iam::aws:policy/CloudSearchReadOnlyAccess",
	},
	{
		Name: "AWSGreengrassFullAccess",
		Id:   "ANPAJWPV6OBK4QONH4J3O",
		Arn:  "arn:aws:iam::aws:policy/AWSGreengrassFullAccess",
	},
	{
		Name: "AWSCloudFormationReadOnlyAccess",
		Id:   "ANPAJWVBEE4I2POWLODLW",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudFormationReadOnlyAccess",
	},
	{
		Name: "AmazonRoute53FullAccess",
		Id:   "ANPAJWVDLG5RPST6PHQ3A",
		Arn:  "arn:aws:iam::aws:policy/AmazonRoute53FullAccess",
	},
	{
		Name: "AWSLambdaRole",
		Id:   "ANPAJX4DPCRGTC4NFDUXI",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaRole",
	},
	{
		Name: "AWSLambdaENIManagementAccess",
		Id:   "ANPAJXAW2Q3KPTURUT2QC",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaENIManagementAccess",
	},
	{
		Name: "AWSOpsWorksCloudWatchLogs",
		Id:   "ANPAJXFIK7WABAY5CPXM4",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksCloudWatchLogs",
	},
	{
		Name: "AmazonAppStreamReadOnlyAccess",
		Id:   "ANPAJXIFDGB4VBX23DX7K",
		Arn:  "arn:aws:iam::aws:policy/AmazonAppStreamReadOnlyAccess",
	},
	{
		Name: "AWSStepFunctionsFullAccess",
		Id:   "ANPAJXKA6VP3UFBVHDPPA",
		Arn:  "arn:aws:iam::aws:policy/AWSStepFunctionsFullAccess",
	},
	{
		Name: "AmazonInspectorReadOnlyAccess",
		Id:   "ANPAJXQNTHTEJ2JFRN2SE",
		Arn:  "arn:aws:iam::aws:policy/AmazonInspectorReadOnlyAccess",
	},
	{
		Name: "AWSCertificateManagerFullAccess",
		Id:   "ANPAJYCHABBP6VQIVBCBQ",
		Arn:  "arn:aws:iam::aws:policy/AWSCertificateManagerFullAccess",
	},
	{
		Name: "PowerUserAccess",
		Id:   "ANPAJYRXTHIB4FOVS3ZXS",
		Arn:  "arn:aws:iam::aws:policy/PowerUserAccess",
	},
	{
		Name: "CloudWatchEventsFullAccess",
		Id:   "ANPAJZLOYLNHESMYOJAFU",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchEventsFullAccess",
	},
}
