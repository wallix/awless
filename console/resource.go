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
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/wallix/awless/graph"
)

type tableResourceDisplayer struct {
	r       *graph.Resource
	headers []ColumnDefinition
}

func (d *tableResourceDisplayer) Print(w io.Writer) error {
	values := make(table, len(d.r.Properties))

	i := 0
	for prop, val := range d.r.Properties {
		var header ColumnDefinition
		for _, h := range d.headers {
			if h.propKey() == prop {
				header = h
			}
		}
		if header == nil {
			header = &StringColumnDefinition{Prop: prop, DisableTruncate: true}
		} else if strheader, ok := header.(StringColumnDefinition); ok {
			header = &StringColumnDefinition{Prop: strheader.Prop, Friendly: strheader.Friendly, DisableTruncate: true}
		}

		if v := values[i]; v == nil {
			values[i] = make([]interface{}, 2)
		}
		values[i][0] = header.title(false)
		values[i][1] = header.format(val)
		i++
	}

	sort.Sort(byCols{table: values, sortBy: []int{0}})

	table := tablewriter.NewWriter(w)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Property" + ascSymbol, "Value"})

	for i := range values {
		if val := fmt.Sprint(values[i][1]); val != "" {
			table.Append([]string{fmt.Sprint(values[i][0]), val})
		}
	}

	table.Render()

	return nil
}

func (d *tableResourceDisplayer) SetResource(r *graph.Resource) {
	d.r = r
}
