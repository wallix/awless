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
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/wallix/awless/logger"
)

var CallerReferenceFunc = func() string {
	return fmt.Sprint(time.Now().UTC().Unix())
}

type CreateDistribution struct {
	_              string `action:"create" entity:"distribution" awsAPI:"cloudfront"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            cloudfrontiface.CloudFrontAPI
	OriginDomain   *string   `templateName:"origin-domain"`
	Certificate    *string   `templateName:"certificate"`
	Comment        *string   `templateName:"comment"`
	DefaultFile    *string   `templateName:"default-file"`
	DomainAliases  []*string `templateName:"domain-aliases"`
	Enable         *bool     `templateName:"enable"`
	ForwardCookies *string   `templateName:"forward-cookies"`
	ForwardQueries *bool     `templateName:"forward-queries"`
	HttpsBehaviour *string   `templateName:"https-behaviour"`
	OriginPath     *string   `templateName:"origin-path"`
	PriceClass     *string   `templateName:"price-class"`
	MinTtl         *int64    `templateName:"min-ttl"`
}

func (cmd *CreateDistribution) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("origin-domain"),
		params.Opt("certificate", "comment", "default-file", "domain-aliases", "enable", "forward-cookies", "forward-queries", "https-behaviour", "min-ttl", "origin-path", "price-class"),
	))
}

func (cmd *CreateDistribution) ManualRun(renv env.Running) (interface{}, error) {
	originId := "orig_1"
	input := &cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			CallerReference: aws.String(CallerReferenceFunc()),
			Comment:         cmd.OriginDomain,
			DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
				MinTTL: aws.Int64(0),
				ForwardedValues: &cloudfront.ForwardedValues{
					Cookies:     &cloudfront.CookiePreference{Forward: aws.String("all")},
					QueryString: aws.Bool(true),
				},
				TrustedSigners: &cloudfront.TrustedSigners{
					Enabled:  aws.Bool(false),
					Quantity: aws.Int64(0),
				},
				TargetOriginId:       aws.String(originId),
				ViewerProtocolPolicy: aws.String("allow-all"),
			},
			Enabled: aws.Bool(true),
			Origins: &cloudfront.Origins{
				Quantity: aws.Int64(1),
				Items: []*cloudfront.Origin{
					{Id: aws.String(originId)},
				},
			},
		},
	}

	if domain := StringValue(cmd.OriginDomain); strings.HasSuffix(domain, ".s3.amazonaws.com") || (strings.HasSuffix(domain, ".amazonaws.com") && strings.Contains(domain, ".s3-website-")) {
		input.DistributionConfig.Origins.Items[0].S3OriginConfig = &cloudfront.S3OriginConfig{OriginAccessIdentity: aws.String("")}
	}

	call := &awsCall{
		fnName: "cloudfront.CreateDistribution",
		fn:     cmd.api.CreateDistribution,
		logger: cmd.logger,
		setters: []setter{
			{val: cmd.OriginDomain, fieldPath: "DistributionConfig.Origins.Items[0].DomainName", fieldType: awsstr},
		},
	}

	if cmd.Certificate != nil {
		call.setters = append(call.setters, setter{val: cmd.Certificate, fieldPath: "DistributionConfig.ViewerCertificate.ACMCertificateArn", fieldType: awsstr})
		call.setters = append(call.setters, setter{val: "sni-only", fieldPath: "DistributionConfig.ViewerCertificate.SSLSupportMethod", fieldType: awsstr})
	}

	if cmd.Comment != nil {
		call.setters = append(call.setters, setter{val: cmd.Comment, fieldPath: "DistributionConfig.Comment", fieldType: awsstr})
	}
	if cmd.DefaultFile != nil {
		call.setters = append(call.setters, setter{val: cmd.DefaultFile, fieldPath: "DistributionConfig.DefaultRootObject", fieldType: awsstr})
	}
	if cmd.DomainAliases != nil {
		call.setters = append(call.setters, setter{val: cmd.DomainAliases, fieldPath: "DistributionConfig.Aliases.Items", fieldType: awsstringslice})
		call.setters = append(call.setters, setter{val: len(cmd.DomainAliases), fieldPath: "DistributionConfig.Aliases.Quantity", fieldType: awsint64})
	}
	if cmd.Enable != nil {
		call.setters = append(call.setters, setter{val: cmd.Enable, fieldPath: "DistributionConfig.Enabled", fieldType: awsbool})
	}
	if cmd.ForwardCookies != nil {
		call.setters = append(call.setters, setter{val: cmd.ForwardCookies, fieldPath: "DistributionConfig.DefaultCacheBehavior.ForwardedValues.Cookies.Forward", fieldType: awsstr})
	}
	if cmd.ForwardQueries != nil {
		call.setters = append(call.setters, setter{val: cmd.ForwardQueries, fieldPath: "DistributionConfig.DefaultCacheBehavior.ForwardedValues.QueryString", fieldType: awsbool})
	}
	if cmd.HttpsBehaviour != nil {
		call.setters = append(call.setters, setter{val: cmd.HttpsBehaviour, fieldPath: "DistributionConfig.DefaultCacheBehavior.ViewerProtocolPolicy", fieldType: awsstr})
	}
	if cmd.MinTtl != nil {
		call.setters = append(call.setters, setter{val: cmd.MinTtl, fieldPath: "DistributionConfig.DefaultCacheBehavior.MinTTL", fieldType: awsint64})
	}
	if cmd.OriginPath != nil {
		call.setters = append(call.setters, setter{val: cmd.OriginPath, fieldPath: "DistributionConfig.Origins.Items[0].OriginPath", fieldType: awsstr})
	}
	if cmd.PriceClass != nil {
		call.setters = append(call.setters, setter{val: cmd.PriceClass, fieldPath: "DistributionConfig.PriceClass", fieldType: awsstr})
	}

	return call.execute(input)
}

func (cmd *CreateDistribution) ExtractResult(i interface{}) string {
	return StringValue(i.(*cloudfront.CreateDistributionOutput).Distribution.Id)
}

type CheckDistribution struct {
	_       string `action:"check" entity:"distribution" awsAPI:"cloudfront"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     cloudfrontiface.CloudFrontAPI
	Id      *string `templateName:"id"`
	State   *string `templateName:"state"`
	Timeout *int64  `templateName:"timeout"`
}

func (cmd *CheckDistribution) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("id"), params.Key("state"), params.Key("timeout")),
		params.Validators{
			"state": params.IsInEnumIgnoreCase("deployed", "inprogress", notFoundState),
		})
}

func (cmd *CheckDistribution) ManualRun(renv env.Running) (interface{}, error) {
	input := &cloudfront.GetDistributionInput{
		Id: cmd.Id,
	}

	c := &checker{
		description: fmt.Sprintf("distribution %s", StringValue(cmd.Id)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.GetDistribution(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "NoSuchDistribution" {
						return notFoundState, nil
					}
					return "", awserr
				} else {
					return "", err
				}
			} else {
				return aws.StringValue(output.Distribution.Status), nil
			}
		},
		expect: StringValue(cmd.State),
		logger: cmd.logger,
	}
	return nil, c.check()
}

type UpdateDistribution struct {
	_              string `action:"update" entity:"distribution" awsAPI:"cloudfront"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            cloudfrontiface.CloudFrontAPI
	Id             *string   `awsName:"Id" awsType:"awsstr" templateName:"id"`
	OriginDomain   *string   `templateName:"origin-domain"`
	Certificate    *string   `templateName:"certificate"`
	Comment        *string   `templateName:"comment"`
	DefaultFile    *string   `templateName:"default-file"`
	DomainAliases  []*string `templateName:"domain-aliases"`
	Enable         *bool     `templateName:"enable"`
	ForwardCookies *string   `templateName:"forward-cookies"`
	ForwardQueries *bool     `templateName:"forward-queries"`
	HttpsBehaviour *string   `templateName:"https-behaviour"`
	OriginPath     *string   `templateName:"origin-path"`
	PriceClass     *string   `templateName:"price-class"`
	MinTtl         *int64    `templateName:"min-ttl"`
}

func (cmd *UpdateDistribution) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"),
		params.Opt("certificate", "comment", "default-file", "domain-aliases", "enable", "forward-cookies", "forward-queries", "https-behaviour", "min-ttl", "origin-domain", "origin-path", "price-class"),
	))
}

func (cmd *UpdateDistribution) ManualRun(renv env.Running) (interface{}, error) {
	distribOutput, err := cmd.api.GetDistribution(&cloudfront.GetDistributionInput{
		Id: cmd.Id,
	})
	if err != nil {
		return nil, err
	}
	distriToUpdate := distribOutput.Distribution
	configToUpdate := distriToUpdate.DistributionConfig
	etag := distribOutput.ETag
	beforeUpdate := distribOutput.Distribution.DistributionConfig.String()

	input := &cloudfront.UpdateDistributionInput{
		IfMatch:            etag,
		DistributionConfig: distriToUpdate.DistributionConfig,
	}

	if err = setFieldWithType(cmd.Id, input, "Id", awsstr); err != nil {
		return nil, err
	}
	if cmd.Enable != nil && BoolValue(cmd.Enable) != BoolValue(distriToUpdate.DistributionConfig.Enabled) {
		if err = setFieldWithType(cmd.Enable, input, "DistributionConfig.Enabled", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.OriginDomain != nil || cmd.OriginPath != nil {
		if configToUpdate.Origins == nil || len(configToUpdate.Origins.Items) == 0 {
			configToUpdate.Origins = &cloudfront.Origins{
				Quantity: aws.Int64(1),
				Items: []*cloudfront.Origin{
					{Id: aws.String("orig_1")},
				},
			}
		}
		if cmd.OriginDomain != nil {
			if err = setFieldWithType(cmd.OriginDomain, input, "DistributionConfig.Origins.Items[0].DomainName", awsstr); err != nil {
				return nil, err
			}
			if domain := aws.StringValue(input.DistributionConfig.Origins.Items[0].DomainName); strings.HasSuffix(domain, ".s3.amazonaws.com") || (strings.HasSuffix(domain, ".amazonaws.com") && strings.Contains(domain, ".s3-website-")) {
				input.DistributionConfig.Origins.Items[0].S3OriginConfig = &cloudfront.S3OriginConfig{OriginAccessIdentity: aws.String("")}
			}
		}

		if cmd.OriginPath != nil {
			if err = setFieldWithType(cmd.OriginPath, input, "DistributionConfig.Origins.Items[0].OriginPath", awsstr); err != nil {
				return nil, err
			}
		}
	}

	if cmd.Certificate != nil {
		if err = setFieldWithType(cmd.Certificate, input, "DistributionConfig.ViewerCertificate.ACMCertificateArn", awsstr); err != nil {
			return nil, err
		}
		if err = setFieldWithType("sni-only", input, "DistributionConfig.ViewerCertificate.SSLSupportMethod", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.Comment != nil {
		if err = setFieldWithType(cmd.Comment, input, "DistributionConfig.Comment", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.DefaultFile != nil {
		if err = setFieldWithType(cmd.DefaultFile, input, "DistributionConfig.DefaultRootObject", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.DomainAliases != nil {
		if err = setFieldWithType(cmd.DomainAliases, input, "DistributionConfig.Aliases.Items", awsstringslice); err != nil {
			return nil, err
		}
	}
	if cmd.Enable != nil {
		if err = setFieldWithType(cmd.Enable, input, "DistributionConfig.Enabled", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.ForwardCookies != nil {
		if err = setFieldWithType(cmd.ForwardCookies, input, "DistributionConfig.DefaultCacheBehavior.ForwardedValues.Cookies.Forward", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.ForwardQueries != nil {
		if err = setFieldWithType(cmd.ForwardQueries, input, "DistributionConfig.DefaultCacheBehavior.ForwardedValues.QueryString", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.HttpsBehaviour != nil {
		if err = setFieldWithType(cmd.HttpsBehaviour, input, "DistributionConfig.DefaultCacheBehavior.ViewerProtocolPolicy", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.MinTtl != nil {
		if err = setFieldWithType(cmd.MinTtl, input, "DistributionConfig.DefaultCacheBehavior.MinTTL", awsint64); err != nil {
			return nil, err
		}
	}
	if cmd.PriceClass != nil {
		if err = setFieldWithType(cmd.PriceClass, input, "DistributionConfig.PriceClass", awsstr); err != nil {
			return nil, err
		}
	}

	if aliases := input.DistributionConfig.Aliases; aliases != nil {
		aliases.Quantity = aws.Int64(int64(len(aliases.Items)))
	}

	if beforeUpdate == input.DistributionConfig.String() {
		cmd.logger.Infof("no property has been changed to distribution '%s'", StringValue(cmd.Id))
		return distribOutput, nil
	}

	start := time.Now()
	var output *cloudfront.UpdateDistributionOutput
	output, err = cmd.api.UpdateDistribution(input)
	cmd.logger.ExtraVerbosef("cloudfront.UpdateDistribution call took %s", time.Since(start))
	return output, err
}

func (cmd *UpdateDistribution) ExtractResult(i interface{}) string {
	switch ii := i.(type) {
	case *cloudfront.GetDistributionOutput:
		return StringValue(ii.ETag)
	case *cloudfront.UpdateDistributionOutput:
		return StringValue(ii.ETag)
	default:
		return ""
	}
}

type DeleteDistribution struct {
	_      string `action:"delete" entity:"distribution" awsAPI:"cloudfront"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    cloudfrontiface.CloudFrontAPI
	Id     *string `awsName:"Id" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteDistribution) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

func (cmd *DeleteDistribution) ManualRun(renv env.Running) (interface{}, error) {
	cmd.logger.Info("disabling distribution")
	updateDistribution := CommandFactory.Build("updatedistribution")().(*UpdateDistribution)
	entries := map[string]interface{}{
		"id":     cmd.Id,
		"enable": false,
	}
	if err := params.Validate(updateDistribution.ParamsSpec().Validators(), entries); err != nil {
		return nil, err
	}

	var etag string
	if out, err := updateDistribution.Run(renv, entries); err != nil {
		return nil, err
	} else if str, ok := out.(string); ok {
		etag = str
	}

	cmd.logger.Info("check distribution disabling has been propagated")
	checkDistribution := CommandFactory.Build("checkdistribution")().(*CheckDistribution)
	entries = map[string]interface{}{
		"id":      cmd.Id,
		"state":   "Deployed",
		"timeout": 1800,
	}
	if err := params.Validate(checkDistribution.ParamsSpec().Validators(), entries); err != nil {
		return nil, err
	}

	if _, err := checkDistribution.Run(renv, entries); err != nil {
		return nil, err
	}

	input := &cloudfront.DeleteDistributionInput{IfMatch: aws.String(fmt.Sprint(etag))}

	if err := setFieldWithType(cmd.Id, input, "Id", awsstr); err != nil {
		return nil, err
	}

	start := time.Now()
	output, err := cmd.api.DeleteDistribution(input)
	cmd.logger.ExtraVerbosef("cloudfront.DeleteDistribution call took %s", time.Since(start))
	return output, err
}
