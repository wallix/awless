package display

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wallix/awless/cloud/aws"
)

const truncateSize = 25

// PropertyDisplayer describe how to display a property in a table
type PropertyDisplayer struct {
	Property                string
	Label                   string
	ColoredValues           map[string]string
	DontTruncate            bool
	TruncateRight           bool
	CollapseIdenticalValues bool
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
