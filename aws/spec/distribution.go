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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/wallix/awless/logger"
)

type CreateDistribution struct {
	_              string `action:"create" entity:"distribution" awsAPI:"cloudfront"`
	logger         *logger.Logger
	api            cloudfrontiface.CloudFrontAPI
	OriginDomain   *string   `templateName:"origin-domain" required:""`
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

func (cmd *CreateDistribution) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

var CallerReferenceFunc = func() string {
	return fmt.Sprint(time.Now().UTC().Unix())
}

func (cmd *CreateDistribution) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	originId := "orig_1"
	input := &cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			CallerReference: aws.String(CallerReferenceFunc()),
			Comment:         aws.String(" "),
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

	if err := setFieldWithType(cmd.OriginDomain, input, "DistributionConfig.Origins.Items[0].DomainName", awsstr); err != nil {
		return nil, err
	}
	if domain := aws.StringValue(input.DistributionConfig.Origins.Items[0].DomainName); strings.HasSuffix(domain, ".s3.amazonaws.com") || (strings.HasSuffix(domain, ".amazonaws.com") && strings.Contains(domain, ".s3-website-")) {
		input.DistributionConfig.Origins.Items[0].S3OriginConfig = &cloudfront.S3OriginConfig{OriginAccessIdentity: aws.String("")}
	}

	if cmd.Certificate != nil {
		if err := setFieldWithType(cmd.Certificate, input, "DistributionConfig.ViewerCertificate.ACMCertificateArn", awsstr); err != nil {
			return nil, err
		}
		if err := setFieldWithType("sni-only", input, "DistributionConfig.ViewerCertificate.SSLSupportMethod", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.Comment != nil {
		if err := setFieldWithType(cmd.Comment, input, "DistributionConfig.Comment", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.DefaultFile != nil {
		if err := setFieldWithType(cmd.DefaultFile, input, "DistributionConfig.DefaultRootObject", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.DomainAliases != nil {
		if err := setFieldWithType(cmd.DomainAliases, input, "DistributionConfig.Aliases.Items", awsstringslice); err != nil {
			return nil, err
		}
	}
	if cmd.Enable != nil {
		if err := setFieldWithType(cmd.Enable, input, "DistributionConfig.Enabled", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.ForwardCookies != nil {
		if err := setFieldWithType(cmd.ForwardCookies, input, "DistributionConfig.DefaultCacheBehavior.ForwardedValues.Cookies.Forward", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.ForwardQueries != nil {
		if err := setFieldWithType(cmd.ForwardQueries, input, "DistributionConfig.DefaultCacheBehavior.ForwardedValues.QueryString", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.HttpsBehaviour != nil {
		if err := setFieldWithType(cmd.HttpsBehaviour, input, "DistributionConfig.DefaultCacheBehavior.ViewerProtocolPolicy", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.MinTtl != nil {
		if err := setFieldWithType(cmd.MinTtl, input, "DistributionConfig.DefaultCacheBehavior.MinTTL", awsint64); err != nil {
			return nil, err
		}
	}
	if cmd.OriginPath != nil {
		if err := setFieldWithType(cmd.OriginPath, input, "DistributionConfig.Origins.Items[0].OriginPath", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.PriceClass != nil {
		if err := setFieldWithType(cmd.PriceClass, input, "DistributionConfig.PriceClass", awsstr); err != nil {
			return nil, err
		}
	}

	if aliases := input.DistributionConfig.Aliases; aliases != nil {
		aliases.Quantity = aws.Int64(int64(len(aliases.Items)))
	}

	start := time.Now()
	output, err := cmd.api.CreateDistribution(input)
	cmd.logger.ExtraVerbosef("cloudfront.CreateDistribution call took %s", time.Since(start))
	return output, err
}

func (cmd *CreateDistribution) ExtractResult(i interface{}) string {
	return StringValue(i.(*cloudfront.CreateDistributionOutput).Distribution.Id)
}

type CheckDistribution struct {
	_       string `action:"check" entity:"distribution" awsAPI:"cloudfront"`
	logger  *logger.Logger
	api     cloudfrontiface.CloudFrontAPI
	Id      *string `templateName:"id" required:""`
	State   *string `templateName:"state" required:""`
	Timeout *int64  `templateName:"timeout" required:""`
}

func (cmd *CheckDistribution) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CheckDistribution) Validate_State() error {
	return NewEnumValidator("deployed", "inprogress", notFoundState).Validate(cmd.State)
}

func (cmd *CheckDistribution) ManualRun(ctx map[string]interface{}) (interface{}, error) {
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
	_      string `action:"update" entity:"distribution" awsAPI:"cloudfront"`
	logger *logger.Logger
	api    cloudfrontiface.CloudFrontAPI
	Id     *string `awsName:"Id" awsType:"awsstr" templateName:"id" required:""`
	Enable *bool   `awsName:"DistributionConfig.Enabled" awsType:"awsbool" templateName:"enable" required:""`
}

func (cmd *UpdateDistribution) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *UpdateDistribution) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	distribOutput, err := cmd.api.GetDistribution(&cloudfront.GetDistributionInput{
		Id: cmd.Id,
	})
	if err != nil {
		return nil, err
	}
	distriToUpdate := distribOutput.Distribution
	etag := distribOutput.ETag
	if enabled := aws.BoolValue(distriToUpdate.DistributionConfig.Enabled); BoolValue(cmd.Enable) == enabled {
		cmd.logger.Infof("distribution '%s' is already enable=%t", StringValue(cmd.Id), enabled)
		return distribOutput, nil
	}

	input := &cloudfront.UpdateDistributionInput{IfMatch: etag, DistributionConfig: distriToUpdate.DistributionConfig}

	if err = setFieldWithType(cmd.Id, input, "Id", awsstr); err != nil {
		return nil, err
	}
	if err = setFieldWithType(cmd.Enable, input, "DistributionConfig.Enabled", awsbool); err != nil {
		return nil, err
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
	api    cloudfrontiface.CloudFrontAPI
	Id     *string `awsName:"Id" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *DeleteDistribution) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *DeleteDistribution) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	cmd.logger.Info("disabling distribution")
	updateDistribution := CommandFactory.Build("updatedistribution")().(*UpdateDistribution)
	updateDistribution.Id = cmd.Id
	updateDistribution.Enable = Bool(false)
	if errs := updateDistribution.ValidateCommand(nil, nil); len(errs) > 0 {
		return nil, fmt.Errorf("%v", errs)
	}

	var etag string
	if out, err := updateDistribution.Run(ctx, nil); err != nil {
		return nil, err
	} else if str, ok := out.(string); ok {
		etag = str
	}

	cmd.logger.Info("check distribution disabling has been propagated")
	checkDistribution := CommandFactory.Build("checkdistribution")().(*CheckDistribution)
	checkDistribution.Id = cmd.Id
	checkDistribution.State = String("Deployed")
	checkDistribution.Timeout = Int64(900)
	if errs := checkDistribution.ValidateCommand(nil, nil); len(errs) > 0 {
		return nil, fmt.Errorf("%v", errs)
	}

	if _, err := checkDistribution.Run(ctx, nil); err != nil {
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
