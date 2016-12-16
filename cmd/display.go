package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

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
	w := tabwriter.NewWriter(os.Stdout, TABWRITERWIDTH, 1, 1, ' ', 0)
	aws.TabularDisplay(item, w)
	w.Flush()
}

func displayGraph(graph *rdf.Graph, resourceType string, properties []string, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w := tabwriter.NewWriter(os.Stdout, TABWRITERWIDTH, 1, 1, ' ', 0)
	var header bytes.Buffer
	for i, prop := range properties {
		header.WriteString(prop)
		if i < len(properties)-1 {
			header.WriteString("\t")
		}
	}
	fmt.Fprintln(w, header.String())
	nodes, err := graph.NodesForType(resourceType)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, node := range nodes {
		nodeProperties, err := cloud.LoadPropertiesTriples(graph, node)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var line bytes.Buffer
		for i, propName := range properties {
			line.WriteString(displayProperty(nodeProperties, propName))
			if i < len(properties)-1 {
				line.WriteString("\t")
			}
		}
		fmt.Fprintln(w, line.String())
	}

	w.Flush()
}

func displayProperty(properties cloud.Properties, name string) string {
	if s := strings.SplitN(name, "[].", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].([]interface{}); ok {
			return trucateToSize(displaySliceProperty(i, s[1]), TABWRITERWIDTH)
		}
	} else if s := strings.SplitN(name, "[]length", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].([]interface{}); ok {
			return trucateToSize(displaySliceLength(i), TABWRITERWIDTH)
		}
	} else if s := strings.SplitN(name, ".", 2); len(s) >= 2 {
		if i, ok := properties[s[0]].(map[string]interface{}); ok {
			return trucateToSize(displayAttributeProperty(i, s[1]), TABWRITERWIDTH)
		}
	} else {
		return trucateToSize(displayStringProperty(properties[name]), TABWRITERWIDTH)
	}
	return ""
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
