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

// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"github.com/wallix/awless/template"
)

var AWSTemplatesDefinitions = map[string]template.TemplateDefinition{
	"createvpc": {
		Action:         "create",
		Entity:         "vpc",
		Api:            "ec2",
		RequiredParams: []string{"cidr"},
		ExtraParams:    []string{"name"},
	},
	"deletevpc": {
		Action:         "delete",
		Entity:         "vpc",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"createsubnet": {
		Action:         "create",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []string{"cidr", "vpc"},
		ExtraParams:    []string{"zone", "name"},
	},
	"updatesubnet": {
		Action:         "update",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{"public"},
	},
	"deletesubnet": {
		Action:         "delete",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"createinstance": {
		Action:         "create",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []string{"image", "count", "count", "type", "subnet", "name"},
		ExtraParams:    []string{"key", "ip", "userdata", "group", "lock"},
	},
	"updateinstance": {
		Action:         "update",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{"type", "group", "lock"},
	},
	"deleteinstance": {
		Action:         "delete",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"startinstance": {
		Action:         "start",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"stopinstance": {
		Action:         "stop",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"checkinstance": {
		Action:         "check",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []string{"id", "state", "timeout"},
		ExtraParams:    []string{},
	},
	"createsecuritygroup": {
		Action:         "create",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []string{"name", "vpc", "description"},
		ExtraParams:    []string{},
	},
	"updatesecuritygroup": {
		Action:         "update",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []string{"id", "cidr", "protocol"},
		ExtraParams:    []string{"inbound", "outbound", "portrange"},
	},
	"deletesecuritygroup": {
		Action:         "delete",
		Entity:         "securitygroup",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"createvolume": {
		Action:         "create",
		Entity:         "volume",
		Api:            "ec2",
		RequiredParams: []string{"zone", "size"},
		ExtraParams:    []string{},
	},
	"deletevolume": {
		Action:         "delete",
		Entity:         "volume",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"attachvolume": {
		Action:         "attach",
		Entity:         "volume",
		Api:            "ec2",
		RequiredParams: []string{"device", "id", "instance"},
		ExtraParams:    []string{},
	},
	"detachvolume": {
		Action:         "detach",
		Entity:         "volume",
		Api:            "ec2",
		RequiredParams: []string{"device", "id", "instance"},
		ExtraParams:    []string{"force"},
	},
	"createinternetgateway": {
		Action:         "create",
		Entity:         "internetgateway",
		Api:            "ec2",
		RequiredParams: []string{},
		ExtraParams:    []string{},
	},
	"deleteinternetgateway": {
		Action:         "delete",
		Entity:         "internetgateway",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"attachinternetgateway": {
		Action:         "attach",
		Entity:         "internetgateway",
		Api:            "ec2",
		RequiredParams: []string{"id", "vpc"},
		ExtraParams:    []string{},
	},
	"detachinternetgateway": {
		Action:         "detach",
		Entity:         "internetgateway",
		Api:            "ec2",
		RequiredParams: []string{"id", "vpc"},
		ExtraParams:    []string{},
	},
	"createroutetable": {
		Action:         "create",
		Entity:         "routetable",
		Api:            "ec2",
		RequiredParams: []string{"vpc"},
		ExtraParams:    []string{},
	},
	"deleteroutetable": {
		Action:         "delete",
		Entity:         "routetable",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"attachroutetable": {
		Action:         "attach",
		Entity:         "routetable",
		Api:            "ec2",
		RequiredParams: []string{"id", "subnet"},
		ExtraParams:    []string{},
	},
	"detachroutetable": {
		Action:         "detach",
		Entity:         "routetable",
		Api:            "ec2",
		RequiredParams: []string{"association"},
		ExtraParams:    []string{},
	},
	"createroute": {
		Action:         "create",
		Entity:         "route",
		Api:            "ec2",
		RequiredParams: []string{"table", "cidr", "gateway"},
		ExtraParams:    []string{},
	},
	"deleteroute": {
		Action:         "delete",
		Entity:         "route",
		Api:            "ec2",
		RequiredParams: []string{"table", "cidr"},
		ExtraParams:    []string{},
	},
	"createtag": {
		Action:         "create",
		Entity:         "tag",
		Api:            "ec2",
		RequiredParams: []string{"resource", "key", "value"},
		ExtraParams:    []string{},
	},
	"createkeypair": {
		Action:         "create",
		Entity:         "keypair",
		Api:            "ec2",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"deletekeypair": {
		Action:         "delete",
		Entity:         "keypair",
		Api:            "ec2",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"createloadbalancer": {
		Action:         "create",
		Entity:         "loadbalancer",
		Api:            "elbv2",
		RequiredParams: []string{"name", "subnets"},
		ExtraParams:    []string{"iptype", "scheme", "groups"},
	},
	"deleteloadbalancer": {
		Action:         "delete",
		Entity:         "loadbalancer",
		Api:            "elbv2",
		RequiredParams: []string{"arn"},
		ExtraParams:    []string{},
	},
	"createlistener": {
		Action:         "create",
		Entity:         "listener",
		Api:            "elbv2",
		RequiredParams: []string{"actiontype", "target", "loadbalancer", "port", "protocol"},
		ExtraParams:    []string{"certificate", "sslpolicy"},
	},
	"deletelistener": {
		Action:         "delete",
		Entity:         "listener",
		Api:            "elbv2",
		RequiredParams: []string{"arn"},
		ExtraParams:    []string{},
	},
	"createtargetgroup": {
		Action:         "create",
		Entity:         "targetgroup",
		Api:            "elbv2",
		RequiredParams: []string{"name", "port", "protocol", "vpc"},
		ExtraParams:    []string{"healthcheckinterval", "healthcheckpath", "healthcheckport", "healthcheckprotocol", "healthchecktimeout", "healthythreshold", "unhealthythreshold", "matcher"},
	},
	"deletetargetgroup": {
		Action:         "delete",
		Entity:         "targetgroup",
		Api:            "elbv2",
		RequiredParams: []string{"arn"},
		ExtraParams:    []string{},
	},
	"createuser": {
		Action:         "create",
		Entity:         "user",
		Api:            "iam",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"deleteuser": {
		Action:         "delete",
		Entity:         "user",
		Api:            "iam",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"attachuser": {
		Action:         "attach",
		Entity:         "user",
		Api:            "iam",
		RequiredParams: []string{"group", "name"},
		ExtraParams:    []string{},
	},
	"detachuser": {
		Action:         "detach",
		Entity:         "user",
		Api:            "iam",
		RequiredParams: []string{"group", "name"},
		ExtraParams:    []string{},
	},
	"creategroup": {
		Action:         "create",
		Entity:         "group",
		Api:            "iam",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"deletegroup": {
		Action:         "delete",
		Entity:         "group",
		Api:            "iam",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"attachpolicy": {
		Action:         "attach",
		Entity:         "policy",
		Api:            "iam",
		RequiredParams: []string{"arn"},
		ExtraParams:    []string{"user", "group"},
	},
	"detachpolicy": {
		Action:         "detach",
		Entity:         "policy",
		Api:            "iam",
		RequiredParams: []string{"arn"},
		ExtraParams:    []string{"user", "group"},
	},
	"createbucket": {
		Action:         "create",
		Entity:         "bucket",
		Api:            "s3",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"deletebucket": {
		Action:         "delete",
		Entity:         "bucket",
		Api:            "s3",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"createstorageobject": {
		Action:         "create",
		Entity:         "storageobject",
		Api:            "s3",
		RequiredParams: []string{"bucket", "file"},
		ExtraParams:    []string{"name"},
	},
	"deletestorageobject": {
		Action:         "delete",
		Entity:         "storageobject",
		Api:            "s3",
		RequiredParams: []string{"bucket", "key"},
		ExtraParams:    []string{},
	},
	"createtopic": {
		Action:         "create",
		Entity:         "topic",
		Api:            "sns",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{},
	},
	"deletetopic": {
		Action:         "delete",
		Entity:         "topic",
		Api:            "sns",
		RequiredParams: []string{"arn"},
		ExtraParams:    []string{},
	},
	"createsubscription": {
		Action:         "create",
		Entity:         "subscription",
		Api:            "sns",
		RequiredParams: []string{"topic", "endpoint", "protocol"},
		ExtraParams:    []string{},
	},
	"deletesubscription": {
		Action:         "delete",
		Entity:         "subscription",
		Api:            "sns",
		RequiredParams: []string{"arn"},
		ExtraParams:    []string{},
	},
	"createqueue": {
		Action:         "create",
		Entity:         "queue",
		Api:            "sqs",
		RequiredParams: []string{"name"},
		ExtraParams:    []string{"delay", "maxMsgSize", "retentionPeriod", "policy", "msgWait", "redrivePolicy", "visibilityTimeout"},
	},
	"deletequeue": {
		Action:         "delete",
		Entity:         "queue",
		Api:            "sqs",
		RequiredParams: []string{"url"},
		ExtraParams:    []string{},
	},
	"createzone": {
		Action:         "create",
		Entity:         "zone",
		Api:            "route53",
		RequiredParams: []string{"callerreference", "name"},
		ExtraParams:    []string{"delegationsetid", "comment", "isprivate", "vpcid", "vpcregion"},
	},
	"deletezone": {
		Action:         "delete",
		Entity:         "zone",
		Api:            "route53",
		RequiredParams: []string{"id"},
		ExtraParams:    []string{},
	},
	"createrecord": {
		Action:         "create",
		Entity:         "record",
		Api:            "route53",
		RequiredParams: []string{"zone", "name", "type", "value", "ttl"},
		ExtraParams:    []string{"comment"},
	},
	"deleterecord": {
		Action:         "delete",
		Entity:         "record",
		Api:            "route53",
		RequiredParams: []string{"zone", "name", "type", "value", "ttl"},
		ExtraParams:    []string{},
	},
}

func DriverSupportedActions() map[string][]string {
	supported := make(map[string][]string)
	supported["create"] = append(supported["create"], "vpc")
	supported["delete"] = append(supported["delete"], "vpc")
	supported["create"] = append(supported["create"], "subnet")
	supported["update"] = append(supported["update"], "subnet")
	supported["delete"] = append(supported["delete"], "subnet")
	supported["create"] = append(supported["create"], "instance")
	supported["update"] = append(supported["update"], "instance")
	supported["delete"] = append(supported["delete"], "instance")
	supported["start"] = append(supported["start"], "instance")
	supported["stop"] = append(supported["stop"], "instance")
	supported["check"] = append(supported["check"], "instance")
	supported["create"] = append(supported["create"], "securitygroup")
	supported["update"] = append(supported["update"], "securitygroup")
	supported["delete"] = append(supported["delete"], "securitygroup")
	supported["create"] = append(supported["create"], "volume")
	supported["delete"] = append(supported["delete"], "volume")
	supported["attach"] = append(supported["attach"], "volume")
	supported["detach"] = append(supported["detach"], "volume")
	supported["create"] = append(supported["create"], "internetgateway")
	supported["delete"] = append(supported["delete"], "internetgateway")
	supported["attach"] = append(supported["attach"], "internetgateway")
	supported["detach"] = append(supported["detach"], "internetgateway")
	supported["create"] = append(supported["create"], "routetable")
	supported["delete"] = append(supported["delete"], "routetable")
	supported["attach"] = append(supported["attach"], "routetable")
	supported["detach"] = append(supported["detach"], "routetable")
	supported["create"] = append(supported["create"], "route")
	supported["delete"] = append(supported["delete"], "route")
	supported["create"] = append(supported["create"], "tag")
	supported["create"] = append(supported["create"], "keypair")
	supported["delete"] = append(supported["delete"], "keypair")
	supported["create"] = append(supported["create"], "loadbalancer")
	supported["delete"] = append(supported["delete"], "loadbalancer")
	supported["create"] = append(supported["create"], "listener")
	supported["delete"] = append(supported["delete"], "listener")
	supported["create"] = append(supported["create"], "targetgroup")
	supported["delete"] = append(supported["delete"], "targetgroup")
	supported["create"] = append(supported["create"], "user")
	supported["delete"] = append(supported["delete"], "user")
	supported["attach"] = append(supported["attach"], "user")
	supported["detach"] = append(supported["detach"], "user")
	supported["create"] = append(supported["create"], "group")
	supported["delete"] = append(supported["delete"], "group")
	supported["attach"] = append(supported["attach"], "policy")
	supported["detach"] = append(supported["detach"], "policy")
	supported["create"] = append(supported["create"], "bucket")
	supported["delete"] = append(supported["delete"], "bucket")
	supported["create"] = append(supported["create"], "storageobject")
	supported["delete"] = append(supported["delete"], "storageobject")
	supported["create"] = append(supported["create"], "topic")
	supported["delete"] = append(supported["delete"], "topic")
	supported["create"] = append(supported["create"], "subscription")
	supported["delete"] = append(supported["delete"], "subscription")
	supported["create"] = append(supported["create"], "queue")
	supported["delete"] = append(supported["delete"], "queue")
	supported["create"] = append(supported["create"], "zone")
	supported["delete"] = append(supported["delete"], "zone")
	supported["create"] = append(supported["create"], "record")
	supported["delete"] = append(supported["delete"], "record")
	return supported
}
