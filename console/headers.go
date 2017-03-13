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

const ascSymbol = " ▲"

const truncateSize = 25

type TimeFormat int

const (
	Humanize TimeFormat = iota
	Basic
	Short
)

type ColumnDefinition interface {
	propKey() string
	title(bool) string
	format(i interface{}) string
}

type ColumnDefinitions []ColumnDefinition

func (d ColumnDefinitions) resolveKey(name string) string {
	low := strings.ToLower(name)
	for _, def := range d {
		if low == strings.ToLower(def.propKey()) {
			return def.propKey()
		}
		switch def.(type) {
		case StringColumnDefinition:
			sdef := def.(StringColumnDefinition)
			if low == strings.ToLower(sdef.Friendly) {
				return def.propKey()
			}
		}
	}

	return ""
}

type StringColumnDefinition struct {
	Prop, Friendly  string
	DisableTruncate bool
	TruncateRight   bool
	TruncateSize    int
}

func (h StringColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	if !h.DisableTruncate {
		size := h.TruncateSize
		if size == 0 {
			size = truncateSize
		}
		if h.TruncateRight {
			return truncateRight(fmt.Sprint(i), size)
		} else {
			return truncateLeft(fmt.Sprint(i), size)
		}
	}
	return fmt.Sprint(i)
}
func (h StringColumnDefinition) propKey() string { return h.Prop }
func (h StringColumnDefinition) title(displayAscSymbol bool) string {
	t := h.Friendly
	if t == "" {
		t = h.Prop
	}
	if displayAscSymbol {
		t += ascSymbol
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

type SliceColumnDefinition struct {
	StringColumnDefinition
}

func (h SliceColumnDefinition) format(i interface{}) string {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Slice {
		return "invalid slice"
	}
	var buf bytes.Buffer
	for i := 0; i < value.Len(); i++ {
		buf.WriteString(fmt.Sprint(value.Index(i).Interface()))
		buf.WriteRune('\n')
	}
	return buf.String()
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
		return humanizeTime(ii)
	case Short:
		return ii.Format("1/2/06 15:04")
	default:
		return ii.Format("Mon, Jan 2, 2006 15:04")
	}
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
		var netStrings []string
		for _, net := range r.IPRanges {
			ones, _ := net.Mask.Size()
			if ones == 0 {
				netStrings = append(netStrings, "any")
			} else {
				netStrings = append(netStrings, net.String())
			}
		}
		w.WriteString(strings.Join(netStrings, ","))

		w.WriteString("(")

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
		switch g.GranteeType {
		case "CanonicalUser":
			w.WriteString("user:")
			if g.GranteeDisplayName != "" {
				w.WriteString(g.GranteeDisplayName)
			} else {
				w.WriteString(g.GranteeID)
			}
		case "Group":
			w.WriteString("group:")
			w.WriteString(g.GranteeID)

		default:
			w.WriteString(g.GranteeType)
			w.WriteString(":")
			w.WriteString(g.GranteeID)

		}
		w.WriteString("]")
		w.WriteString(" ")
	}
	return w.String()
}
