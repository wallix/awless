package display

import (
	"bytes"
	"fmt"
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
