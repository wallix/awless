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
	"createvpc": TemplateDefinition{
		Action:         "create",
		Entity:         "vpc",
		requiredParams: []string{"cidr"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletevpc": TemplateDefinition{
		Action:         "delete",
		Entity:         "vpc",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createsubnet": TemplateDefinition{
		Action:         "create",
		Entity:         "subnet",
		requiredParams: []string{"cidr", "vpc"},
		extraParams:    []string{"zone"},
		tagsMapping:    []string{},
	},
	"updatesubnet": TemplateDefinition{
		Action:         "update",
		Entity:         "subnet",
		requiredParams: []string{"id"},
		extraParams:    []string{"public-vms"},
		tagsMapping:    []string{},
	},
	"deletesubnet": TemplateDefinition{
		Action:         "delete",
		Entity:         "subnet",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createinstance": TemplateDefinition{
		Action:         "create",
		Entity:         "instance",
		requiredParams: []string{"image", "type", "count", "count", "subnet"},
		extraParams:    []string{"lock", "key", "ip", "group", "userdata"},
		tagsMapping:    []string{"name"},
	},
	"updateinstance": TemplateDefinition{
		Action:         "update",
		Entity:         "instance",
		requiredParams: []string{"id"},
		extraParams:    []string{"lock", "group", "type"},
		tagsMapping:    []string{},
	},
	"deleteinstance": TemplateDefinition{
		Action:         "delete",
		Entity:         "instance",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"startinstance": TemplateDefinition{
		Action:         "start",
		Entity:         "instance",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"stopinstance": TemplateDefinition{
		Action:         "stop",
		Entity:         "instance",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"checkinstance": TemplateDefinition{
		Action:         "check",
		Entity:         "instance",
		requiredParams: []string{"id", "state", "timeout"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createsecuritygroup": TemplateDefinition{
		Action:         "create",
		Entity:         "securitygroup",
		requiredParams: []string{"description", "name", "vpc"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"updatesecuritygroup": TemplateDefinition{
		Action:         "update",
		Entity:         "securitygroup",
		requiredParams: []string{"cidr", "id", "protocol"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletesecuritygroup": TemplateDefinition{
		Action:         "delete",
		Entity:         "securitygroup",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createvolume": TemplateDefinition{
		Action:         "create",
		Entity:         "volume",
		requiredParams: []string{"zone", "size"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletevolume": TemplateDefinition{
		Action:         "delete",
		Entity:         "volume",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"attachvolume": TemplateDefinition{
		Action:         "attach",
		Entity:         "volume",
		requiredParams: []string{"device", "instance", "id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createinternetgateway": TemplateDefinition{
		Action:         "create",
		Entity:         "internetgateway",
		requiredParams: []string{},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deleteinternetgateway": TemplateDefinition{
		Action:         "delete",
		Entity:         "internetgateway",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"attachinternetgateway": TemplateDefinition{
		Action:         "attach",
		Entity:         "internetgateway",
		requiredParams: []string{"id", "vpc"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"detachinternetgateway": TemplateDefinition{
		Action:         "detach",
		Entity:         "internetgateway",
		requiredParams: []string{"id", "vpc"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createroutetable": TemplateDefinition{
		Action:         "create",
		Entity:         "routetable",
		requiredParams: []string{"vpc"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deleteroutetable": TemplateDefinition{
		Action:         "delete",
		Entity:         "routetable",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"attachroutetable": TemplateDefinition{
		Action:         "attach",
		Entity:         "routetable",
		requiredParams: []string{"id", "subnet"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"detachroutetable": TemplateDefinition{
		Action:         "detach",
		Entity:         "routetable",
		requiredParams: []string{"association"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createroute": TemplateDefinition{
		Action:         "create",
		Entity:         "route",
		requiredParams: []string{"cidr", "gateway", "table"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deleteroute": TemplateDefinition{
		Action:         "delete",
		Entity:         "route",
		requiredParams: []string{"cidr", "table"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createtags": TemplateDefinition{
		Action:         "create",
		Entity:         "tags",
		requiredParams: []string{"resource"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createkeypair": TemplateDefinition{
		Action:         "create",
		Entity:         "keypair",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletekeypair": TemplateDefinition{
		Action:         "delete",
		Entity:         "keypair",
		requiredParams: []string{"id"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createuser": TemplateDefinition{
		Action:         "create",
		Entity:         "user",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deleteuser": TemplateDefinition{
		Action:         "delete",
		Entity:         "user",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"creategroup": TemplateDefinition{
		Action:         "create",
		Entity:         "group",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletegroup": TemplateDefinition{
		Action:         "delete",
		Entity:         "group",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"attachpolicy": TemplateDefinition{
		Action:         "attach",
		Entity:         "policy",
		requiredParams: []string{"arn", "user"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"detachpolicy": TemplateDefinition{
		Action:         "detach",
		Entity:         "policy",
		requiredParams: []string{"arn", "user"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createbucket": TemplateDefinition{
		Action:         "create",
		Entity:         "bucket",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"deletebucket": TemplateDefinition{
		Action:         "delete",
		Entity:         "bucket",
		requiredParams: []string{"name"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
	"createstorageobject": TemplateDefinition{
		Action:         "create",
		Entity:         "storageobject",
		requiredParams: []string{"file", "bucket"},
		extraParams:    []string{"name"},
		tagsMapping:    []string{},
	},
	"deletestorageobject": TemplateDefinition{
		Action:         "delete",
		Entity:         "storageobject",
		requiredParams: []string{"bucket", "key"},
		extraParams:    []string{},
		tagsMapping:    []string{},
	},
}
