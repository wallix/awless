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

package main

import (
	"text/template"

	"github.com/wallix/awless/gen/aws"
)

func generateProperties() {
	templ, err := template.New("properties").Parse(propertiesTempl)
	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, aws.PropertiesDefinitions, CLOUD_PROPERTIES_DIR, "gen_properties.go")
}

func generateRDFProperties() {
	templ, err := template.New("rdfProperties").Parse(rdfPropertiesTempl)
	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, aws.PropertiesDefinitions, CLOUD_RDF_DIR, "gen_rdf.go")
}

const propertiesTempl = `/* Copyright 2017 WALLIX

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
package properties

const (
  {{- range $, $prop := . }}
  {{ $prop.AwlessLabel }} = "{{ $prop.AwlessLabel }}"
  {{- end }}
)

`

const rdfPropertiesTempl = `/* Copyright 2017 WALLIX

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
package rdf

import "github.com/wallix/awless/cloud/properties"

const (
  {{- range $, $prop := . }}
  {{ $prop.AwlessLabel }} = "{{ $prop.RDFLabel }}"
  {{- end }}
)

func init() {
	Labels = map[string]string{
	{{- range $, $prop := . }}
	properties.{{ $prop.AwlessLabel }}: {{ $prop.AwlessLabel }},
	{{- end }}
	}
}

var Properties = RDFProperties{
  {{- range $, $prop := . }}
  {{ $prop.AwlessLabel }}: {ID: {{ $prop.AwlessLabel }}, RdfType: "{{ $prop.RDFType }}", RdfsLabel: "{{ $prop.AwlessLabel }}", RdfsDefinedBy: "{{ $prop.RdfsDefinedBy }}", RdfsDataType: "{{ $prop.RdfsDataType }}"},
  {{- end }}
}

`
