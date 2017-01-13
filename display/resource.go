package display

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/wallix/awless/cloud/aws"
)

type ResourceDisplayer interface {
	Print() string
	SetResource(*aws.Resource)
}

func BuildResourceDisplayer(headers []ColumnDefinition, opts Options) ResourceDisplayer {
	switch opts.Format {
	case "table":
		return &tableResourceDisplayer{headers: headers}
	default:
		panic(fmt.Sprintf("unknown displayer for %s", opts.Format))
	}
}

type tableResourceDisplayer struct {
	r       *aws.Resource
	headers []ColumnDefinition
}

func (d *tableResourceDisplayer) Print() string {
	var w bytes.Buffer

	values := make(table, len(d.r.Properties()))

	i := 0
	for prop, val := range d.r.Properties() {
		var header ColumnDefinition
		for _, h := range d.headers {
			if h.propKey() == prop {
				header = h
			}
		}
		if header == nil {
			header = StringColumnDefinition{Prop: prop}
		}

		if v := values[i]; v == nil {
			values[i] = make([]interface{}, 2)
		}
		values[i][0] = header.title(false)
		values[i][1] = header.format(val)
		i++
	}

	sort.Sort(byCols{table: values, sortBy: []int{0}})

	table := tablewriter.NewWriter(&w)
	table.SetHeader([]string{"Property" + ascSymbol, "Value"})

	for i := range values {
		table.Append([]string{fmt.Sprint(values[i][0]), fmt.Sprint(values[i][1])})
	}

	table.Render()

	return w.String()
}

func (d *tableResourceDisplayer) SetResource(r *aws.Resource) {
	d.r = r
}
