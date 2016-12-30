package display

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/wallix/awless/cloud/aws"
)

const truncateSize = 25

var (
	PropertiesDisplayer = AwlessResourcesDisplayer{
		Services: map[string]*ServiceDisplayer{
			aws.InfraServiceName: &ServiceDisplayer{
				Resources: map[string]*ResourceDisplayer{
					"instance": &ResourceDisplayer{
						Properties: map[string]*PropertyDisplayer{
							"Id":        {Property: "Id"},
							"Name":      {Property: "Tags[].Name", Label: "Name"},
							"State":     {Property: "State.Name", Label: "State", ColoredValues: map[string]string{"running": "green", "stopped": "red"}},
							"Type":      {Property: "Type"},
							"PublicIp":  {Property: "PublicIp"},
							"PrivateIp": {Property: "PrivateIp"},
						},
					},
					"vpc": &ResourceDisplayer{
						Properties: map[string]*PropertyDisplayer{
							"Id":        {Property: "Id"},
							"IsDefault": {Property: "IsDefault", Label: "Default", ColoredValues: map[string]string{"true": "green"}},
							"State":     {Property: "State"},
							"CidrBlock": {Property: "CidrBlock"},
						},
					},
					"subnet": &ResourceDisplayer{
						Properties: map[string]*PropertyDisplayer{
							"Id": {Property: "Id"},
							"MapPublicIpOnLaunch": {Property: "MapPublicIpOnLaunch", Label: "Public VMs", ColoredValues: map[string]string{"true": "red"}},
							"State":               {Property: "State", ColoredValues: map[string]string{"available": "green"}},
							"CidrBlock":           {Property: "CidrBlock"},
						},
					},
				},
			},
			aws.AccessServiceName: &ServiceDisplayer{
				Resources: map[string]*ResourceDisplayer{
					"user": &ResourceDisplayer{
						Properties: map[string]*PropertyDisplayer{
							"Id":               {Property: "Id"},
							"Name":             {Property: "Name"},
							"Arn":              {Property: "Arn"},
							"Path":             {Property: "Path"},
							"PasswordLastUsed": {Property: "PasswordLastUsed"},
						},
					},
					"role": &ResourceDisplayer{
						Properties: map[string]*PropertyDisplayer{
							"Id":         {Property: "Id"},
							"Name":       {Property: "Name"},
							"Arn":        {Property: "Arn"},
							"CreateDate": {Property: "CreateDate"},
							"Path":       {Property: "Path"},
						},
					},
					"policy": &ResourceDisplayer{
						Properties: map[string]*PropertyDisplayer{
							"Id":           {Property: "Id"},
							"Name":         {Property: "Name"},
							"Arn":          {Property: "Arn"},
							"Description":  {Property: "Description"},
							"isAttachable": {Property: "isAttachable"},
							"CreateDate":   {Property: "CreateDate"},
							"UpdateDate":   {Property: "UpdateDate"},
							"Path":         {Property: "Path"},
						},
					},
					"group": &ResourceDisplayer{
						Properties: map[string]*PropertyDisplayer{
							"Id":         {Property: "Id"},
							"Name":       {Property: "Name"},
							"Arn":        {Property: "Arn"},
							"CreateDate": {Property: "CreateDate"},
							"Path":       {Property: "Path"},
						},
					},
				},
			},
		},
	}
)

type AwlessResourcesDisplayer struct {
	Services map[string]*ServiceDisplayer
}

type ServiceDisplayer struct {
	Resources map[string]*ResourceDisplayer
}

type ResourceDisplayer struct {
	Properties map[string]*PropertyDisplayer
}

// PropertyDisplayer describe how to display a property in a table
type PropertyDisplayer struct {
	Property      string
	Label         string
	ColoredValues map[string]string
	DontTruncate  bool
	TruncateRight bool
}

func (p *PropertyDisplayer) displayName() string {
	if p.Label == "" {
		return p.Property
	}
	return p.Label
}

func (p *PropertyDisplayer) display(value string) string {
	if !p.DontTruncate {
		if p.TruncateRight {
			value = truncateRight(value, truncateSize)
		} else {
			value = truncateLeft(value, truncateSize)
		}
	}
	if p.ColoredValues != nil {
		return colorDisplay(value, p.ColoredValues)
	}
	return value
}

func (p *PropertyDisplayer) displayForceColor(value string, c color.Attribute) string {
	if !p.DontTruncate {
		if p.TruncateRight {
			value = truncateRight(value, truncateSize)
		} else {
			value = truncateLeft(value, truncateSize)
		}
	}
	return color.New(c).SprintFunc()(value)
}

func propertyValue(properties aws.Properties, displayName string) string {
	var res string
	if s := strings.SplitN(displayName, "[].", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].([]interface{}); ok {
			res = propertyValueSlice(i, s[1])
		}
	} else if s := strings.SplitN(displayName, "[]length", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].([]interface{}); ok {
			res = propertyValueSliceLength(i)
		}
	} else if s := strings.SplitN(displayName, ".", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].(map[string]interface{}); ok {
			res = propertyValueAttribute(i, s[1])
		}
	} else {
		res = propertyValueString(properties[displayName])
	}
	return res
}

func propertyValueString(prop interface{}) string {
	switch pp := prop.(type) {
	case string:
		return pp
	case bool:
		return fmt.Sprint(pp)
	default:
		return ""
	}
}

func propertyValueSlice(prop []interface{}, key string) string {
	for _, p := range prop {
		//map [key: result]
		if m, ok := p.(map[string]interface{}); ok && m[key] != nil {
			return fmt.Sprint(m[key])
		}

		//map["Key": key, "Value": result]
		if m, ok := p.(map[string]interface{}); ok && m["Key"] == key {
			return fmt.Sprint(m["Value"])
		}
	}
	return ""
}

func propertyValueSliceLength(prop []interface{}) string {
	return strconv.Itoa(len(prop))
}

func propertyValueAttribute(attr map[string]interface{}, key string) string {
	return fmt.Sprint(attr[key])
}
