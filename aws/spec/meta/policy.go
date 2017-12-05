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
	"bytes"
	"fmt"
	"strings"

	"github.com/wallix/awless/template"
)

type AttachPolicyMeta struct {
	_        struct{} `action:"attach" entity:"policy"`
	Services *string  `templateName:"services" required:""`
}

func (m *AttachPolicyMeta) Match(action, entity string, paramKeys []string) bool {
	if action != "attach" && entity != "policy" {
		return false
	}
	return contains(paramKeys, "services")
}

func (m *AttachPolicyMeta) Resolve(params map[string]string) (*template.Template, error) {
	servicesStr := strings.TrimPrefix(params["services"], "[")
	servicesStr = strings.TrimSuffix(servicesStr, "]")
	servicesStr = strings.TrimSpace(servicesStr)
	if servicesStr == "" {
		return nil, nil
	}

	delete(params, "services")

	services := strings.Split(servicesStr, ",")
	var line bytes.Buffer
	line.WriteString("attach policy")
	for k, v := range params {
		line.WriteString(fmt.Sprintf(" %s=%s", k, v))
	}
	line.WriteString(" service=")
	var tplLines []string
	for _, service := range services {
		service = strings.TrimSpace(service)
		tplLines = append(tplLines, line.String()+service)
	}
	return template.Parse(strings.Join(tplLines, "\n"))
}
