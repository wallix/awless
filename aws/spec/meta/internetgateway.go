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

package awsspecmeta

import (
	"fmt"

	"github.com/wallix/awless/template"
)

type CreateInternetgatewayMeta struct {
	_   struct{} `action:"create" entity:"internetgateway"`
	Vpc *string  `templateName:"vpc" required:""`
}

func (m *CreateInternetgatewayMeta) Match(action, entity string, paramKeys []string) bool {
	if action != "create" && entity != "internetgateway" {
		return false
	}
	return contains(paramKeys, "vpc")
}

func (m *CreateInternetgatewayMeta) Resolve(params map[string]string) (*template.Template, error) {
	return template.Parse(fmt.Sprintf("igw = create internetgateway\nattach internetgateway id=$igw vpc=%s", params["vpc"]))
}
