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

	"github.com/wallix/awless/cloud/graph"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/wallix/awless/logger"
)

type CreateRecord struct {
	_       string `action:"create" entity:"record" awsAPI:"route53"`
	logger  *logger.Logger
	graph   cloudgraph.GraphAPI
	api     route53iface.Route53API
	Zone    *string `templateName:"zone" required:""`
	Name    *string `templateName:"name" required:""`
	Type    *string `templateName:"type" required:""`
	Value   *string `templateName:"value" required:""`
	Ttl     *int64  `templateName:"ttl" required:""`
	Comment *string `templateName:"comment"`
}

func (cmd *CreateRecord) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateRecord) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	start := time.Now()
	output, err := changeResourceRecordSets(cmd.api, String("CREATE"), cmd.Zone, cmd.Name, cmd.Type, cmd.Value, cmd.Comment, cmd.Ttl)
	cmd.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	return output, err
}

func (cmd *CreateRecord) ExtractResult(i interface{}) string {
	return StringValue(i.(*route53.ChangeResourceRecordSetsOutput).ChangeInfo.Id)
}

type UpdateRecord struct {
	_      string `action:"update" entity:"record" awsAPI:"route53"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    route53iface.Route53API
	Zone   *string `templateName:"zone" required:""`
	Name   *string `templateName:"name" required:""`
	Type   *string `templateName:"type" required:""`
	Value  *string `templateName:"value" required:""`
	Ttl    *int64  `templateName:"ttl" required:""`
}

func (cmd *UpdateRecord) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *UpdateRecord) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	start := time.Now()
	output, err := changeResourceRecordSets(cmd.api, String("UPSERT"), cmd.Zone, cmd.Name, cmd.Type, cmd.Value, nil, cmd.Ttl)
	cmd.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	return output, err
}

func (cmd *UpdateRecord) ExtractResult(i interface{}) string {
	return StringValue(i.(*route53.ChangeResourceRecordSetsOutput).ChangeInfo.Id)
}

type DeleteRecord struct {
	_      string `action:"delete" entity:"record" awsAPI:"route53"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    route53iface.Route53API
	Zone   *string `templateName:"zone" required:""`
	Name   *string `templateName:"name" required:""`
	Type   *string `templateName:"type" required:""`
	Value  *string `templateName:"value" required:""`
	Ttl    *int64  `templateName:"ttl" required:""`
}

func (cmd *DeleteRecord) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *DeleteRecord) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	start := time.Now()
	output, err := changeResourceRecordSets(cmd.api, String("DELETE"), cmd.Zone, cmd.Name, cmd.Type, cmd.Value, nil, cmd.Ttl)
	cmd.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	return output, err
}

func (cmd *DeleteRecord) ExtractResult(i interface{}) string {
	return StringValue(i.(*route53.ChangeResourceRecordSetsOutput).ChangeInfo.Id)
}

func changeResourceRecordSets(api route53iface.Route53API, action, zone, name, recordType, value, comment *string, ttl *int64) (*route53.ChangeResourceRecordSetsOutput, error) {
	input := &route53.ChangeResourceRecordSetsInput{}
	var err error
	// Required params
	err = setFieldWithType(zone, input, "HostedZoneId", awsstr)
	if err != nil {
		return nil, err
	}
	resourceRecord := &route53.ResourceRecord{}
	change := &route53.Change{ResourceRecordSet: &route53.ResourceRecordSet{ResourceRecords: []*route53.ResourceRecord{resourceRecord}}}
	input.ChangeBatch = &route53.ChangeBatch{Changes: []*route53.Change{change}}
	if err = setFieldWithType(action, change, "Action", awsstr); err != nil {
		return nil, err
	}
	if err = setFieldWithType(name, change, "ResourceRecordSet.Name", awsstr); err != nil {
		return nil, err
	}
	if err = setFieldWithType(recordType, change, "ResourceRecordSet.Type", awsstr); err != nil {
		return nil, err
	}
	if err = setFieldWithType(ttl, change, "ResourceRecordSet.TTL", awsint64); err != nil {
		return nil, err
	}
	if err = setFieldWithType(value, resourceRecord, "Value", awsstr); err != nil {
		return nil, err
	}

	// Extra params
	if comment != nil {
		if err = setFieldWithType(comment, input, "ChangeBatch.Comment", awsstr); err != nil {
			return nil, err
		}
	}

	return api.ChangeResourceRecordSets(input)
}
