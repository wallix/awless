package display

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

const ascSymbol = " ▲"

//const truncateSize = 25

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
	Prop, Friendly string
	DontTruncate   bool
	TruncateRight  bool
	TruncateSize   int
}

func (h StringColumnDefinition) format(i interface{}) string {
	if i == nil {
		return ""
	}
	if !h.DontTruncate {
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
