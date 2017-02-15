// DO NOT EDIT
// This file was automatically generated with go generate

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
	"fmt"
	"strings"
)

type TemplateDefinition struct {
	Action, Entity                           string
	requiredParams, extraParams, tagsMapping []string
}

func (def TemplateDefinition) String() string {
	var required []string
	for _, v := range def.Required() {
		required = append(required, fmt.Sprintf("%s = { %s.%s }", v, def.Entity, v))
	}
	var tags []string
	for _, v := range def.tagsMapping {
		tags = append(tags, fmt.Sprintf("%s = { %s.%s }", v, def.Entity, v))
	}
	return fmt.Sprintf("%s %s %s %s", def.Action, def.Entity, strings.Join(required, " "), strings.Join(tags, " "))
}

func (def TemplateDefinition) Required() []string {
	return def.requiredParams
}

func (def TemplateDefinition) Extra() []string {
	return def.extraParams
}

var AWSTemplatesDefinitions = map[string]TemplateDefinition{
	"createvpc": {
		Action:         "create",
		Entity:         "vpc",
		requiredParams: []string{"cidr"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletevpc": {
		Action:         "delete",
		Entity:         "vpc",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createsubnet": {
		Action:         "create",
		Entity:         "subnet",
		requiredParams: []string{"cidr", "vpc"},
		extraParams:    []string{"zone"},
		tagsMapping:    []string{},
	},
	"updatesubnet": {
		Action:         "update",
		Entity:         "subnet",
		requiredParams: []string{"id"},
		extraParams:    []string{"public-vms"},
		tagsMapping:    []string{},
	},
	"deletesubnet": {
		Action:         "delete",
		Entity:         "subnet",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createinstance": {
		Action:         "create",
		Entity:         "instance",
		requiredParams: []string{"image", "type", "count", "count", "subnet"},
		extraParams:    []string{"lock", "key", "ip", "group", "userdata"},
		tagsMapping:    []string{"name"},
	},
	"updateinstance": {
		Action:         "update",
		Entity:         "instance",
		requiredParams: []string{"id"},
		extraParams:    []string{"lock", "group", "type"},
		tagsMapping:    []string{},
	},
	"deleteinstance": {
		Action:         "delete",
		Entity:         "instance",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"startinstance": {
		Action:         "start",
		Entity:         "instance",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"stopinstance": {
		Action:         "stop",
		Entity:         "instance",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"checkinstance": {
		Action:         "check",
		Entity:         "instance",
		requiredParams: []string{"id", "state", "timeout"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createsecuritygroup": {
		Action:         "create",
		Entity:         "securitygroup",
		requiredParams: []string{"description", "name", "vpc"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"updatesecuritygroup": {
		Action:         "update",
		Entity:         "securitygroup",
		requiredParams: []string{"cidr", "id", "protocol"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletesecuritygroup": {
		Action:         "delete",
		Entity:         "securitygroup",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createvolume": {
		Action:         "create",
		Entity:         "volume",
		requiredParams: []string{"zone", "size"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletevolume": {
		Action:         "delete",
		Entity:         "volume",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"attachvolume": {
		Action:         "attach",
		Entity:         "volume",
		requiredParams: []string{"device", "instance", "id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createinternetgateway": {
		Action:         "create",
		Entity:         "internetgateway",
		requiredParams: []string{},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deleteinternetgateway": {
		Action:         "delete",
		Entity:         "internetgateway",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"attachinternetgateway": {
		Action:         "attach",
		Entity:         "internetgateway",
		requiredParams: []string{"id", "vpc"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"detachinternetgateway": {
		Action:         "detach",
		Entity:         "internetgateway",
		requiredParams: []string{"id", "vpc"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createroutetable": {
		Action:         "create",
		Entity:         "routetable",
		requiredParams: []string{"vpc"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deleteroutetable": {
		Action:         "delete",
		Entity:         "routetable",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"attachroutetable": {
		Action:         "attach",
		Entity:         "routetable",
		requiredParams: []string{"id", "subnet"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"detachroutetable": {
		Action:         "detach",
		Entity:         "routetable",
		requiredParams: []string{"association"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createroute": {
		Action:         "create",
		Entity:         "route",
		requiredParams: []string{"cidr", "gateway", "table"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deleteroute": {
		Action:         "delete",
		Entity:         "route",
		requiredParams: []string{"cidr", "table"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createtags": {
		Action:         "create",
		Entity:         "tags",
		requiredParams: []string{"resource"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createkeypair": {
		Action:         "create",
		Entity:         "keypair",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletekeypair": {
		Action:         "delete",
		Entity:         "keypair",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createuser": {
		Action:         "create",
		Entity:         "user",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deleteuser": {
		Action:         "delete",
		Entity:         "user",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"creategroup": {
		Action:         "create",
		Entity:         "group",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletegroup": {
		Action:         "delete",
		Entity:         "group",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"attachpolicy": {
		Action:         "attach",
		Entity:         "policy",
		requiredParams: []string{"arn", "user"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"detachpolicy": {
		Action:         "detach",
		Entity:         "policy",
		requiredParams: []string{"arn", "user"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createbucket": {
		Action:         "create",
		Entity:         "bucket",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletebucket": {
		Action:         "delete",
		Entity:         "bucket",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createstorageobject": {
		Action:         "create",
		Entity:         "storageobject",
		requiredParams: []string{"file", "bucket"},
		extraParams:    []string{"name"},
		tagsMapping:    []string{},
	},
	"deletestorageobject": {
		Action:         "delete",
		Entity:         "storageobject",
		requiredParams: []string{"bucket", "key"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
}
