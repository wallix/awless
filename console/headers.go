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

package console

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/wallix/awless/graph"
)

type TimeFormat int

const (
	Humanize TimeFormat = iota
	Basic
	Short
)

type ColumnDefinition interface {
	propKey() string
	title(...string) string
	format(i interface{}) string
}

type ColumnDefinitions []ColumnDefinition

func (d ColumnDefinitions) resolveKey(name string) string {
	low := strings.ToLower(name)
	for _, def := range d {
		switch low {
		case strings.ToLower(def.propKey()), strings.ToLower(def.title()):
			return def.propKey()
		}
	}
	return ""
}

type StringColumnDefinition struct {
	Prop, Friendly string
}

func (h StringColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	return fmt.Sprint(i)
}
func (h StringColumnDefinition) propKey() string { return h.Prop }
func (h StringColumnDefinition) title(suffix ...string) string {
	t := h.Friendly
	if t == "" {
		t = h.Prop
	}
	if len(suffix) > 0 {
		t += suffix[0]
	}
	return t
}

type ColoredValueColumnDefinition struct {
	StringColumnDefinition
	ColoredValues map[string]color.Attribute
}

func (h ColoredValueColumnDefinition) format(i interface{}) string {
	str := h.StringColumnDefinition.format(i)
	col, ok := h.ColoredValues[str]
	if ok {
		return color.New(col).SprintFunc()(str)
	}
	return str
}

type ARNLastValueColumnDefinition struct {
	StringColumnDefinition
	Separator string
}

func (h ARNLastValueColumnDefinition) format(i interface{}) string {
	str := h.StringColumnDefinition.format(i)
	splits := strings.Split(str, h.Separator)
	if len(splits) > 1 {
		return splits[len(splits)-1]
	}
	return str
}

func ToShortArn(s string) string {
	index := strings.LastIndex(s, ":")
	if index > 0 {
		return s[index+1:]
	}
	return s
}

type SliceColumnDefinition struct {
	ForEach func(string) string
	StringColumnDefinition
}

func (h SliceColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Slice {
		return fmt.Sprintf("invalid slice: %T", i)
	}
	var buf bytes.Buffer
	for i := 0; i < value.Len(); i++ {
		s := fmt.Sprint(value.Index(i).Interface())
		if h.ForEach != nil {
			s = h.ForEach(s)
		}
		buf.WriteString(s)
		if i < value.Len()-1 {
			buf.WriteRune(' ')
		}
	}
	return buf.String()
}

type KeyValuesColumnDefinition struct {
	StringColumnDefinition
}

func (h KeyValuesColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	ii, ok := i.([]*graph.KeyValue)
	if !ok {
		return fmt.Sprintf("invalid keyvalue, got %T", i)
	}
	var b bytes.Buffer
	for i, kv := range ii {
		b.WriteString(fmt.Sprintf("%s:%s", color.CyanString(kv.KeyName), kv.Value))
		if i < len(ii)-1 {
			b.WriteString(" ")
		}
	}
	return b.String()
}

type TimeColumnDefinition struct {
	StringColumnDefinition
	Format TimeFormat
}

func (h TimeColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	ii, ok := i.(time.Time)
	if !ok {
		return "invalid time"
	}
	switch h.Format {
	case Humanize:
		return HumanizeTime(ii)
	case Short:
		return ii.Format("1/2/06 15:04")
	default:
		return ii.Format("Mon, Jan 2, 2006 15:04")
	}
}

type StorageColumnDefinition struct {
	StringColumnDefinition
	Unit storageUnit
}

func (h StorageColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Uint || val.Kind() == reflect.Uint64 {
		return HumanizeStorage(val.Uint(), h.Unit)
	}
	if val.Kind() == reflect.Int || val.Kind() == reflect.Int64 {
		return HumanizeStorage(uint64(val.Int()), h.Unit)
	}
	return "invalid size"
}

type FirewallRulesColumnDefinition struct {
	StringColumnDefinition
}

func (h FirewallRulesColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	ii, ok := i.([]*graph.FirewallRule)
	if !ok {
		return "invalid rules"
	}
	var w bytes.Buffer

	for _, r := range ii {
		w.WriteString("[")
		var netStrings []string
		for _, net := range r.IPRanges {
			netStrings = append(netStrings, net.String())
		}
		for _, src := range r.Sources {
			netStrings = append(netStrings, src)
		}
		w.WriteString(strings.Join(netStrings, ";"))

		w.WriteString("](")

		switch {
		case r.Protocol == "any":
			w.WriteString(r.Protocol)
		case r.PortRange.Any:
			w.WriteString(fmt.Sprintf("%s:any", r.Protocol))
		case r.PortRange.FromPort == r.PortRange.ToPort:
			w.WriteString(fmt.Sprintf("%s:%d", r.Protocol, r.PortRange.FromPort))
		default:
			w.WriteString(fmt.Sprintf("%s:%d-%d", r.Protocol, r.PortRange.FromPort, r.PortRange.ToPort))
		}

		w.WriteString(") ")
	}
	return w.String()
}

type RoutesColumnDefinition struct {
	StringColumnDefinition
}

func (h RoutesColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	ii, ok := i.([]*graph.Route)
	if !ok {
		return "invalid routes"
	}
	var w bytes.Buffer

	for _, r := range ii {
		if r.Destination != nil {
			w.WriteString(r.Destination.String())
		}
		if r.DestinationIPv6 != nil && r.Destination != nil {
			w.WriteString("+")
		}
		if r.DestinationIPv6 != nil {
			w.WriteString(r.DestinationIPv6.String())
		}
		w.WriteString("->")
		if len(r.Targets) > 1 {
			w.WriteString("[")
		}
		for _, t := range r.Targets {
			switch t.Type {
			case graph.EgressOnlyInternetGatewayTarget:
				w.WriteString("inbound-internget-gw")
			case graph.GatewayTarget:
				w.WriteString("gw")
			case graph.InstanceTarget:
				w.WriteString("inst")
			case graph.NatTarget:
				w.WriteString("nat")
			case graph.NetworkInterfaceTarget:
				w.WriteString("ni")
			case graph.VpcPeeringConnectionTarget:
				w.WriteString("vpc")
			default:
				w.WriteString("unknown")
			}
			w.WriteString(":")
			w.WriteString(t.Ref)
			w.WriteString(" ")
		}
		if len(r.Targets) > 1 {
			w.WriteString("] ")
		}
	}
	return w.String()
}

type GrantsColumnDefinition struct {
	StringColumnDefinition
}

func (h GrantsColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	ii, ok := i.([]*graph.Grant)
	if !ok {
		return "invalid grants"
	}
	var w bytes.Buffer

	for _, g := range ii {
		w.WriteString(g.Permission)
		w.WriteString("[")
		switch g.Grantee.GranteeType {
		case "CanonicalUser":
			w.WriteString("user:")
			if g.Grantee.GranteeDisplayName != "" {
				w.WriteString(g.Grantee.GranteeDisplayName)
			} else {
				w.WriteString(g.Grantee.GranteeID)
			}
		case "Group":
			w.WriteString("group:")
			w.WriteString(g.Grantee.GranteeID)

		default:
			w.WriteString(g.Grantee.GranteeType)
			w.WriteString(":")
			w.WriteString(g.Grantee.GranteeID)

		}
		w.WriteString("]")
		w.WriteString(" ")
	}
	return w.String()
}
