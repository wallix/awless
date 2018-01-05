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
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/match"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/wallix/awless/logger"
)

type CreateRecord struct {
	_       string `action:"create" entity:"record" awsAPI:"route53"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     route53iface.Route53API
	Zone    *string   `templateName:"zone"`
	Name    *string   `templateName:"name"`
	Type    *string   `templateName:"type"`
	Values  []*string `templateName:"values"`
	Ttl     *int64    `templateName:"ttl"`
	Comment *string   `templateName:"comment"`
}

func (cmd *CreateRecord) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.AllOf(params.Key("name"), params.Key("ttl"), params.Key("type"), params.OnlyOneOf(params.Key("values"), params.Key("value")), params.Key("zone"),
		params.Opt("comment"),
	))
	builder.AddReducer(valueToValues, "value")
	return builder.Done()
}

func (cmd *CreateRecord) ManualRun(renv env.Running) (interface{}, error) {
	start := time.Now()
	output, err := changeResourceRecordSets(cmd.api, String("CREATE"), cmd.Zone, cmd.Name, cmd.Type, cmd.Values, cmd.Comment, cmd.Ttl)
	cmd.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	return output, err
}

func (cmd *CreateRecord) ExtractResult(i interface{}) string {
	return StringValue(i.(*route53.ChangeResourceRecordSetsOutput).ChangeInfo.Id)
}

type UpdateRecord struct {
	_      string `action:"update" entity:"record" awsAPI:"route53"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    route53iface.Route53API
	Zone   *string   `templateName:"zone"`
	Name   *string   `templateName:"name"`
	Type   *string   `templateName:"type"`
	Values []*string `templateName:"values"`
	Ttl    *int64    `templateName:"ttl"`
}

func (cmd *UpdateRecord) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.AllOf(params.Key("name"), params.Key("ttl"), params.Key("type"), params.OnlyOneOf(params.Key("values"), params.Key("value")), params.Key("zone")))
	builder.AddReducer(valueToValues, "value")
	return builder.Done()
}

func (cmd *UpdateRecord) ManualRun(renv env.Running) (interface{}, error) {
	start := time.Now()
	output, err := changeResourceRecordSets(cmd.api, String("UPSERT"), cmd.Zone, cmd.Name, cmd.Type, cmd.Values, nil, cmd.Ttl)
	cmd.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	return output, err
}

func (cmd *UpdateRecord) ExtractResult(i interface{}) string {
	return StringValue(i.(*route53.ChangeResourceRecordSetsOutput).ChangeInfo.Id)
}

type DeleteRecord struct {
	_      string `action:"delete" entity:"record" awsAPI:"route53"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    route53iface.Route53API
	Zone   *string   `templateName:"zone"`
	Name   *string   `templateName:"name"`
	Type   *string   `templateName:"type"`
	Values []*string `templateName:"values"`
	Ttl    *int64    `templateName:"ttl"`
}

func (cmd *DeleteRecord) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(
		params.OnlyOneOf(
			params.AllOf(params.Key("name"), params.Key("ttl"), params.Key("type"), params.OnlyOneOf(params.Key("values"), params.Key("value")), params.Key("zone")),
			params.AllOf(params.Key("id")),
		),
	)
	builder.AddReducer(valueToValues, "value")
	builder.AddReducer(
		func(values map[string]interface{}) (map[string]interface{}, error) {
			id, hasId := values["id"].(string)
			if hasId {
				delete(values, "id")
				r, err := cmd.graph.FindOne(cloud.NewQuery(cloud.Record).Match(match.Property(properties.ID, id)))
				if err != nil {
					return values, fmt.Errorf("can not find record for %s: %s", id, err)
				}
				if r == nil {
					return values, fmt.Errorf("record not found with id '%s' in local model ", id)
				}
				if name, ok := r.Property(properties.Name); ok {
					values["name"] = name
				}
				if ttl, ok := r.Property(properties.TTL); ok {
					values["ttl"] = ttl
				}
				if t, ok := r.Property(properties.Type); ok {
					values["type"] = t
				}
				if rec, ok := r.Property(properties.Records); ok {
					values["values"] = rec
				}
				parents, err := cmd.graph.ResourceRelations(r, rdf.ParentOf, false)
				if err != nil {
					return values, fmt.Errorf("cannot get record's zone: %s", err)
				}
				if len(parents) != 1 || parents[0].Type() != cloud.Zone {
					return values, fmt.Errorf("record is not in a zone, got %v ", parents)
				}
				values["zone"] = parents[0].Id()
			}
			return values, nil
		},
		"id",
	)
	return builder.Done()
}

func (cmd *DeleteRecord) ManualRun(renv env.Running) (interface{}, error) {
	start := time.Now()
	output, err := changeResourceRecordSets(cmd.api, String("DELETE"), cmd.Zone, cmd.Name, cmd.Type, cmd.Values, nil, cmd.Ttl)
	cmd.logger.ExtraVerbosef("route53.ChangeResourceRecordSets call took %s", time.Since(start))
	return output, err
}

func (cmd *DeleteRecord) ExtractResult(i interface{}) string {
	return StringValue(i.(*route53.ChangeResourceRecordSetsOutput).ChangeInfo.Id)
}

func changeResourceRecordSets(api route53iface.Route53API, action, zone, name, recordType *string, values []*string, comment *string, ttl *int64) (*route53.ChangeResourceRecordSetsOutput, error) {
	input := &route53.ChangeResourceRecordSetsInput{}
	var err error
	// Required params
	err = setFieldWithType(zone, input, "HostedZoneId", awsstr)
	if err != nil {
		return nil, err
	}
	change := &route53.Change{ResourceRecordSet: &route53.ResourceRecordSet{}}
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
	for _, value := range values {
		resourceRecord := &route53.ResourceRecord{}
		if err = setFieldWithType(value, resourceRecord, "Value", awsstr); err != nil {
			return nil, err
		}
		change.ResourceRecordSet.ResourceRecords = append(change.ResourceRecordSet.ResourceRecords, resourceRecord)
	}

	// Extra params
	if comment != nil {
		if err = setFieldWithType(comment, input, "ChangeBatch.Comment", awsstr); err != nil {
			return nil, err
		}
	}

	return api.ChangeResourceRecordSets(input)
}

func valueToValues(values map[string]interface{}) (map[string]interface{}, error) {
	if value, hasValue := values["value"]; hasValue {
		return map[string]interface{}{"values": value}, nil
	} else {
		return nil, nil
	}
}
