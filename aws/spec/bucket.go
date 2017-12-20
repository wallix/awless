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
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/wallix/awless/logger"
)

type CreateBucket struct {
	_      string `action:"create" entity:"bucket" awsAPI:"s3" awsCall:"CreateBucket" awsInput:"s3.CreateBucketInput" awsOutput:"s3.CreateBucketOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    s3iface.S3API
	Name   *string `awsName:"Bucket" awsType:"awsstr" templateName:"name"`
	Acl    *string `awsName:"ACL" awsType:"awsstr" templateName:"acl"`
}

func (cmd *CreateBucket) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"),
		params.Opt("acl"),
	))
}

func (cmd *CreateBucket) ExtractResult(i interface{}) string {
	return StringValue(cmd.Name)
}

type UpdateBucket struct {
	_                string `action:"update" entity:"bucket" awsAPI:"s3"`
	logger           *logger.Logger
	graph            cloud.GraphAPI
	api              s3iface.S3API
	Name             *string `templateName:"name"`
	Acl              *string `templateName:"acl"`
	PublicWebsite    *bool   `templateName:"public-website"`
	RedirectHostname *string `templateName:"redirect-hostname"`
	IndexSuffix      *string `templateName:"index-suffix"`
	EnforceHttps     *bool   `templateName:"enforce-https"`
}

func (cmd *UpdateBucket) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"),
		params.Opt("acl", "enforce-https", "index-suffix", "public-website", "redirect-hostname"),
	))
}

func (cmd *UpdateBucket) ManualRun(renv env.Running) (interface{}, error) {
	start := time.Now()

	if cmd.Acl != nil { // Update the canned ACL to apply to the bucket
		input := &s3.PutBucketAclInput{
			Bucket: cmd.Name,
		}

		if err := setFieldWithType(cmd.Acl, input, "ACL", awsstr); err != nil {
			return nil, err
		}

		if _, err := cmd.api.PutBucketAcl(input); err != nil {
			return nil, err
		}

		cmd.logger.ExtraVerbosef("s3.PutBucketAcl call took %s", time.Since(start))
		return nil, nil
	}

	if cmd.PublicWebsite != nil { // Set/Unset this bucket as a public website
		if BoolValue(cmd.PublicWebsite) {
			input := &s3.PutBucketWebsiteInput{
				Bucket:               cmd.Name,
				WebsiteConfiguration: &s3.WebsiteConfiguration{},
			}
			if cmd.RedirectHostname != nil {
				input.WebsiteConfiguration.RedirectAllRequestsTo = &s3.RedirectAllRequestsTo{HostName: cmd.RedirectHostname}
				if BoolValue(cmd.EnforceHttps) {
					input.WebsiteConfiguration.RedirectAllRequestsTo.Protocol = aws.String("https")
				}
			} else if cmd.IndexSuffix != nil {
				input.WebsiteConfiguration.IndexDocument = &s3.IndexDocument{Suffix: cmd.IndexSuffix}
			} else {
				input.WebsiteConfiguration.IndexDocument = &s3.IndexDocument{Suffix: aws.String("index.html")}
			}

			if _, err := cmd.api.PutBucketWebsite(input); err != nil {
				return nil, err
			}
		} else {
			if _, err := cmd.api.DeleteBucketWebsite(&s3.DeleteBucketWebsiteInput{Bucket: cmd.Name}); err != nil {
				return nil, err
			}
		}
		cmd.logger.ExtraVerbosef("s3.PutBucketWebsite call took %s", time.Since(start))
	}
	return nil, nil
}

type DeleteBucket struct {
	_      string `action:"delete" entity:"bucket" awsAPI:"s3" awsCall:"DeleteBucket" awsInput:"s3.DeleteBucketInput" awsOutput:"s3.DeleteBucketOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    s3iface.S3API
	Name   *string `awsName:"Bucket" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteBucket) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}
