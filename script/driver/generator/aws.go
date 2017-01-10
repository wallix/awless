//go:generate go run $GOFILE
package main

import (
	"html/template"
	"os"
	"strings"
)

func main() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"capitalize": capitalize,
	}).Parse(funcTempl)

	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("../aws/driver_funcs.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	err = templ.Execute(f, definitions)
	if err != nil {
		panic(err)
	}
}

var definitions = []struct {
	ParamsMapping                     map[string]string
	Action, Entity                    string
	Input, ApiMethod, OutputExtractor string
}{
	{
		Action: "create", Entity: "vpc", Input: "ec2.CreateVpcInput", ApiMethod: "CreateVpc", OutputExtractor: "Vpc.VpcId",
		ParamsMapping: map[string]string{
			"CidrBlock": "cidr",
		},
	},
	{
		Action: "create", Entity: "subnet", Input: "ec2.CreateSubnetInput", ApiMethod: "CreateSubnet", OutputExtractor: "Subnet.SubnetId",
		ParamsMapping: map[string]string{
			"CidrBlock": "cidr",
			"VpcId":     "vpc",
		},
	},
	{
		Action: "create", Entity: "instance", Input: "ec2.RunInstancesInput", ApiMethod: "RunInstances", OutputExtractor: "Instances[0].InstanceId",
		ParamsMapping: map[string]string{
			"ImageId":      "base",
			"MaxCount":     "count",
			"MinCount":     "count",
			"InstanceType": "type",
			"SubnetId":     "subnet",
		},
	},
}

const funcTempl = `// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
  "fmt"
  "io/ioutil"
  "log"
  "reflect"
  "strings"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/service/ec2/ec2iface"
  "github.com/wallix/awless/script/driver"
)
{{ range $index, $def := . }}
func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}(params map[string]interface{}) (interface{}, error) {
  input := &{{ $def.Input }}{}

  {{ range $awsField, $field := $def.ParamsMapping }}
  setField(params["{{ $field }}"], input, "{{ $awsField }}")
  {{- end }}

  output, err := d.api.{{ $def.ApiMethod }}(input)
  if err != nil {
    d.logger.Printf("{{ $def.Action }} {{ $def.Entity }} error: %s", err)
    return nil, err
  }
  d.logger.Println("{{ $def.Action }} {{ $def.Entity }} done")

  return aws.StringValue(output.{{ $def.OutputExtractor }}), nil
}
{{ end }}
`

func capitalize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	return strings.ToUpper(s)
}
