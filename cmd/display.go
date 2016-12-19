package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

type PropertyDisplayer struct {
	Property      string
	Label         string
	ColoredValues map[string]string
}

type StringColoredDisplayer struct {
}

func (p *PropertyDisplayer) DisplayName() string {
	if p.Label == "" {
		return p.Property
	}
	return p.Label
}

const TABWRITERWIDTH = 20

func display(item interface{}, err error, format ...string) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(format) > 0 {
		switch format[0] {
		case "raw":
			fmt.Println(item)
		default:
			lineDisplay(item)
		}
	} else {
		lineDisplay(item)
	}
}

func lineDisplay(item interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	aws.TabularDisplay(item, table)
	table.Render()
}

func displayGraph(graph *rdf.Graph, resourceType string, properties []PropertyDisplayer, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	var headers []string
	for _, prop := range properties {
		headers = append(headers, prop.DisplayName())
	}
	table.SetHeader(headers)
	nodes, err := graph.NodesForType(resourceType)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, node := range nodes {
		nodeProperties, err := aws.LoadPropertiesFromGraph(graph, node)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var line []string
		for _, prop := range properties {
			line = append(line, displayProperty(nodeProperties, prop))
		}
		table.Append(line)
	}

	table.Render()
}

func displayProperty(properties aws.Properties, propertyDisplayer PropertyDisplayer) string {
	var res string
	if s := strings.SplitN(propertyDisplayer.Property, "[].", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].([]interface{}); ok {
			res = displaySliceProperty(i, s[1])
		}
	} else if s := strings.SplitN(propertyDisplayer.Property, "[]length", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].([]interface{}); ok {
			res = displaySliceLength(i)
		}
	} else if s := strings.SplitN(propertyDisplayer.Property, ".", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].(map[string]interface{}); ok {
			res = displayAttributeProperty(i, s[1])
		}
	} else {
		res = displayStringProperty(properties[propertyDisplayer.Property])
	}
	if propertyDisplayer.ColoredValues != nil {
		return colorDisplay(trucateToSize(res, TABWRITERWIDTH), propertyDisplayer.ColoredValues)
	} else {
		return trucateToSize(res, TABWRITERWIDTH)
	}
}

func displayStringProperty(prop interface{}) string {
	switch pp := prop.(type) {
	case string:
		return pp
	case bool:
		return fmt.Sprint(pp)
	default:
		return ""
	}
}

func displaySliceProperty(prop []interface{}, key string) string {
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

func displaySliceLength(prop []interface{}) string {
	return strconv.Itoa(len(prop))
}

func displayAttributeProperty(attr map[string]interface{}, key string) string {
	return fmt.Sprint(attr[key])
}

func humanize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}
	return strings.ToUpper(s)
}

func trucateToSize(str string, maxSize int) string {
	if maxSize < 3 {
		return str
	}
	if len(str) > maxSize {
		len := len(str)
		return "..." + str[len-maxSize+3:len-1]
	}
	return str
}

func stringToColor(str string) color.Attribute {
	switch strings.ToLower(str) {
	case "red":
		return color.FgRed
	case "yellow":
		return color.FgYellow
	case "blue":
		return color.FgBlue
	case "green":
		return color.FgGreen
	default:
		return color.FgBlack
	}
}

func colorDisplay(str string, coloredValues map[string]string) string {
	col := coloredValues[str]
	if col != "" {
		return color.New(stringToColor(col)).SprintFunc()(str)
	}
	return str
}
