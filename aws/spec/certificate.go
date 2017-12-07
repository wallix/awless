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
	"fmt"
	"time"

	"github.com/wallix/awless/cloud/graph"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/wallix/awless/logger"
)

type CreateCertificate struct {
	_                 string `action:"create" entity:"certificate" awsAPI:"acm"`
	logger            *logger.Logger
	graph             cloudgraph.GraphAPI
	api               acmiface.ACMAPI
	Domains           []*string `templateName:"domains" required:""`
	ValidationDomains []*string `templateName:"validation-domains"`
}

func (cmd *CreateCertificate) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateCertificate) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	input := &acm.RequestCertificateInput{}
	domains := awssdk.StringValueSlice(cmd.Domains)
	if len(domains) == 0 {
		return nil, fmt.Errorf("'domains' must contain at least one element")
	}
	// Required params
	err := setFieldWithType(domains[0], input, "DomainName", awsstr, ctx)
	if err != nil {
		return nil, err
	}
	if len(domains) > 1 {
		if err = setFieldWithType(domains[1:], input, "SubjectAlternativeNames", awsstringslice, ctx); err != nil {
			return nil, err
		}
	}

	domainsToValidate := make(map[string]string)
	// Extra params
	if len(cmd.ValidationDomains) > 0 {
		var validationOptions []*acm.DomainValidationOption

		validation := awssdk.StringValueSlice(cmd.ValidationDomains)
		for i, validationDomain := range validation {
			if i >= len(domains) {
				return nil, fmt.Errorf("there is more validation-domains than certificate domains: %v", validation)
			}
			domainsToValidate[domains[i]] = validationDomain
			validationOptions = append(validationOptions, &acm.DomainValidationOption{DomainName: String(domains[i]), ValidationDomain: String(validationDomain)})
		}
		input.DomainValidationOptions = validationOptions
	}
	if len(domainsToValidate) < len(domains) {
		for i := len(domainsToValidate); i < len(domains); i++ {
			domainsToValidate[domains[i]] = domains[i]
		}
	}

	start := time.Now()
	var output *acm.RequestCertificateOutput
	output, err = cmd.api.RequestCertificate(input)
	if err != nil {
		return nil, err
	}
	cmd.logger.ExtraVerbosef("acm.RequestCertificate call took %s", time.Since(start))

	if len(domainsToValidate) > 0 {
		var helpMsg bytes.Buffer
		for domain, validationDomain := range domainsToValidate {
			helpMsg.WriteString(fmt.Sprintf("\n\t-> %s: {admin/administrator/hostmaster/postmaster/webmaster}@%s", domain, validationDomain))
		}
		cmd.logger.Warningf("validate your certificates by following the instructions sent by email to %s", helpMsg.String())
	}
	return output, nil
}

func (cmd *CreateCertificate) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*acm.RequestCertificateOutput).CertificateArn)
}

type DeleteCertificate struct {
	_      string `action:"delete" entity:"certificate" awsAPI:"acm" awsCall:"DeleteCertificate" awsInput:"acm.DeleteCertificateInput" awsOutput:"acm.DeleteCertificateOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    acmiface.ACMAPI
	Arn    *string `awsName:"CertificateArn" awsType:"awsstr" templateName:"arn" required:""`
}

func (cmd *DeleteCertificate) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type CheckCertificate struct {
	_       string `action:"check" entity:"certificate" awsAPI:"acm"`
	logger  *logger.Logger
	graph   cloudgraph.GraphAPI
	api     acmiface.ACMAPI
	Arn     *string `templateName:"arn" required:""`
	State   *string `templateName:"state" required:""`
	Timeout *int64  `templateName:"timeout" required:""`
}

func (cmd *CheckCertificate) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CheckCertificate) Validate_State() error {
	return NewEnumValidator("issued", "pending_validation", notFoundState).Validate(cmd.State)
}

func (cmd *CheckCertificate) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	input := &acm.DescribeCertificateInput{
		CertificateArn: cmd.Arn,
	}

	c := &checker{
		description: fmt.Sprintf("certificate %s", StringValue(cmd.Arn)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.DescribeCertificate(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "CertificateNotFound" {
						return notFoundState, nil
					}
				} else {
					return "", err
				}
			}
			if output.Certificate == nil {
				return notFoundState, nil
			}
			return StringValue(output.Certificate.Status), nil
		},
		expect: StringValue(cmd.State),
		logger: cmd.logger,
	}
	return nil, c.check()
}
