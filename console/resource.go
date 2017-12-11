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

	"github.com/olekukonko/tablewriter"
	"github.com/wallix/awless/cloud"
)

type tableResourceDisplayer struct {
	maxwidth          int
	r                 cloud.Resource
	columnDefinitions []ColumnDefinition
}

func (d *tableResourceDisplayer) Print(w io.Writer) error {
	values := make(table, len(d.r.Properties()))

	i := 0
	propertyNameMaxWith := 13
	for prop, val := range d.r.Properties() {
		var header ColumnDefinition
		for _, h := range d.columnDefinitions {
			if h.propKey() == prop {
				header = h
			}
		}
		if header == nil {
			header = &StringColumnDefinition{Prop: prop}
		}

		if v := values[i]; v == nil {
			values[i] = make([]interface{}, 2)
		}
		values[i][0] = header.title()
		if l := len(header.title()); l > propertyNameMaxWith {
			propertyNameMaxWith = l
		}
		values[i][1] = header.format(val)
		i++
	}

	ds := defaultSorter{sortBy: []int{0}}
	ds.sort(values)

	valueColumnMaxwidth := d.maxwidth - (propertyNameMaxWith + 7) // ( = border + 2 * margin + border + 2 * margin + border)
	if valueColumnMaxwidth <= 0 {
		valueColumnMaxwidth = 50
	}

	table := tablewriter.NewWriter(w)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetColWidth(valueColumnMaxwidth)
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Property" + ds.symbol(), "Value"})

	wraper := autoWraper{maxWidth: valueColumnMaxwidth, wrappingChar: " "}

	for i := range values {
		if val := fmt.Sprint(values[i][1]); val != "" {
			table.Append([]string{fmt.Sprint(values[i][0]), wraper.Wrap(val)})
		}
	}

	table.Render()

	return nil
}

func (d *tableResourceDisplayer) SetResource(r cloud.Resource) {
	d.r = r
}
