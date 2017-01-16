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
	// VPC
	{
		Action: "create", Entity: "vpc", Input: "CreateVpcInput", ApiMethod: "CreateVpc", OutputExtractor: "Vpc.VpcId",
		ParamsMapping: map[string]string{
			"CidrBlock": "cidr",
		},
	},
	{
		Action: "delete", Entity: "vpc", Input: "DeleteVpcInput", ApiMethod: "DeleteVpc",
		ParamsMapping: map[string]string{
			"VpcId": "id",
		},
	},

	// SUBNET
	{
		Action: "create", Entity: "subnet", Input: "CreateSubnetInput", ApiMethod: "CreateSubnet", OutputExtractor: "Subnet.SubnetId",
		ParamsMapping: map[string]string{
			"CidrBlock": "cidr",
			"VpcId":     "vpc",
		},
	},
	{
		Action: "delete", Entity: "subnet", Input: "DeleteSubnetInput", ApiMethod: "DeleteSubnet",
		ParamsMapping: map[string]string{
			"SubnetId": "id",
		},
	},

	// INSTANCES
	{
		Action: "create", Entity: "instance", Input: "RunInstancesInput", ApiMethod: "RunInstances", OutputExtractor: "Instances[0].InstanceId",
		ParamsMapping: map[string]string{
			"ImageId":      "base",
			"MaxCount":     "count",
			"MinCount":     "count",
			"InstanceType": "type",
			"SubnetId":     "subnet",
		},
	},
	{
		Action: "delete", Entity: "instance", Input: "TerminateInstancesInput", ApiMethod: "TerminateInstances",
		ParamsMapping: map[string]string{
			"InstanceIds": "id",
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

func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}_DryRun(params map[string]interface{}) (interface{}, error) {
  input := &ec2.{{ $def.Input }}{}

  input.DryRun = aws.Bool(true)
  {{ range $awsField, $field := $def.ParamsMapping }}
  setField(params["{{ $field }}"], input, "{{ $awsField }}")
  {{- end }}

  _, err := d.api.{{ $def.ApiMethod }}(input)
  if awsErr, ok := err.(awserr.Error); ok {
    if awsErr.Code() == "DryRunOperation" {
      d.logger.Println("dry run: {{ $def.Action }} {{ $def.Entity }} done")
      return nil, nil
    }
  }

  d.logger.Printf("dry run: {{ $def.Action }} {{ $def.Entity }} error: %s", err)
  return nil, err
}

func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}(params map[string]interface{}) (interface{}, error) {
  input := &ec2.{{ $def.Input }}{}
  {{ range $awsField, $field := $def.ParamsMapping }}
  setField(params["{{ $field }}"], input, "{{ $awsField }}")
  {{- end }}

  output, err := d.api.{{ $def.ApiMethod }}(input)
  if err != nil {
    d.logger.Printf("{{ $def.Action }} {{ $def.Entity }} error: %s", err)
    return nil, err
  }
  d.logger.Println("{{ $def.Action }} {{ $def.Entity }} done")
  {{ if eq $def.OutputExtractor "" }}
  return output, nil {{ else}}
  return aws.StringValue(output.{{ $def.OutputExtractor }}), nil {{ end }}
}
{{ end }}
`

func capitalize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	return strings.ToUpper(s)
}
