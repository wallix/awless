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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/match"
	"github.com/wallix/awless/graph"
)

var (
	tableColWidth   = 30
	autowrapMaxSize = 35
)

type Displayer interface {
	Print(io.Writer) error
}

type sorter interface {
	sort(table)
	columns() []int
	symbol() string
}

type Builder struct {
	filters           []string
	tagFilters        []string
	tagKeyFilters     []string
	tagValueFilters   []string
	columnDefinitions []ColumnDefinition
	format            string
	rdfType           string
	sort              []int
	reverseSort       bool
	maxwidth          int
	dataSource        interface{}
	root              cloud.Resource
	noHeaders         bool
}

func (b *Builder) SetSource(i interface{}) *Builder {
	b.dataSource = i
	return b
}

func (b *Builder) buildQuery() (cloud.Query, error) {
	var matchers []cloud.Matcher
	for _, f := range b.filters {
		splits := strings.SplitN(f, "=", 2)
		if len(splits) == 2 {
			name, val := strings.TrimSpace(strings.Title(splits[0])), strings.TrimSpace(splits[1])
			key := ColumnDefinitions(b.columnDefinitions).resolveKey(name)

			if key != "" {
				matchers = append(matchers, match.Property(key, val).IgnoreCase().MatchString().Contains())
			} else {
				var allowed []string
				for _, h := range b.columnDefinitions {
					allowed = append(allowed, h.propKey())
				}
				return cloud.Query{}, fmt.Errorf("Invalid filter key '%s'. Expecting any of: %s. (Note: filter keys/values are case insensitive)", name, strings.Join(allowed, ", "))
			}
		}
	}

	for _, f := range b.tagFilters {
		splits := strings.SplitN(f, "=", 2)
		if len(splits) == 2 {
			key, val := strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
			matchers = append(matchers, match.Tag(key, val))
		}
	}

	for _, k := range b.tagKeyFilters {
		matchers = append(matchers, match.TagKey(k))
	}

	for _, v := range b.tagValueFilters {
		matchers = append(matchers, match.TagValue(v))
	}
	q := cloud.NewQuery(b.rdfType)
	if len(matchers) > 0 {
		q = cloud.NewQuery(b.rdfType).Match(match.And(matchers...))
	}

	return q, nil
}

func (b *Builder) Build() (Displayer, error) {
	base := fromGraphDisplayer{sorter: &defaultSorter{sortBy: b.sort, descending: b.reverseSort}, rdfType: b.rdfType, columnDefinitions: b.columnDefinitions, maxwidth: b.maxwidth, noHeaders: b.noHeaders}

	switch b.dataSource.(type) {
	case cloud.GraphAPI:
		if b.rdfType == "" {
			gph := b.dataSource.(cloud.GraphAPI)
			switch b.format {
			case "table":
				dis := &multiResourcesTableDisplayer{base}
				dis.setGraph(gph)
				return dis, nil
			case "json":
				dis := &multiResourcesJSONDisplayer{base}
				dis.setGraph(gph)
				return dis, nil
			case "porcelain":
				dis := &porcelainDisplayer{base}
				dis.setGraph(gph)
				return dis, nil
			default:
				fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'table'\n", b.format)
				dis := &multiResourcesTableDisplayer{base}
				dis.setGraph(gph)
				return dis, nil
			}
		}

		filteredGraph := b.dataSource.(cloud.GraphAPI)
		q, err := b.buildQuery()
		if err != nil {
			return nil, err
		}
		if filteredGraph, err = filteredGraph.FilterGraph(q); err != nil {
			return nil, err
		}

		switch b.format {
		case "csv":
			dis := &csvDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis, nil
		case "tsv":
			dis := &tsvDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis, nil
		case "json":
			dis := &jsonDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis, nil
		case "porcelain":
			dis := &porcelainDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis, nil
		case "table":
			dis := &tableDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis, nil
		default:
			fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'table'\n", b.format)
			dis := &tableDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis, nil
		}
	case cloud.Resource:
		dis := &tableResourceDisplayer{columnDefinitions: b.columnDefinitions, maxwidth: b.maxwidth}
		dis.SetResource(b.dataSource.(cloud.Resource))
		return dis, nil
	case *graph.Diff:
		base := fromDiffDisplayer{root: b.root}
		switch b.format {
		case "tree":
			dis := &diffTreeDisplayer{&base}
			dis.SetDiff(b.dataSource.(*graph.Diff))
			return dis, nil
		case "table":
			dis := &diffTableDisplayer{&base}
			dis.SetDiff(b.dataSource.(*graph.Diff))
			return dis, nil
		default:
			fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'tree'\n", b.format)
			dis := &diffTreeDisplayer{&base}
			dis.SetDiff(b.dataSource.(*graph.Diff))
			return dis, nil
		}
	}

	return nil, nil
}

type optsFn func(b *Builder) *Builder

func BuildOptions(opts ...optsFn) *Builder {
	b := &Builder{}

	b.sort = []int{0}
	b.format = "table"

	for _, fn := range opts {
		fn(b)
	}

	if len(b.columnDefinitions) == 0 {
		b.columnDefinitions = DefaultsColumnDefinitions[b.rdfType]
	}

	return b
}

func WithFormat(format string) optsFn {
	return func(b *Builder) *Builder {
		b.format = format
		return b
	}
}

func WithColumns(properties []string) optsFn {
	return func(b *Builder) *Builder {
		if len(properties) == 0 {
			properties = ColumnsInListing[b.rdfType]
		}
		var columns []ColumnDefinition
		for _, p := range properties {
			var found bool
			for _, definition := range DefaultsColumnDefinitions[b.rdfType] {
				if strings.ToLower(p) == strings.ToLower(definition.propKey()) || strings.ToLower(p) == strings.ToLower(definition.title()) {
					found = true
					columns = append(columns, definition)
					continue
				}
			}
			if !found {
				columns = append(columns, StringColumnDefinition{Prop: strings.Title(p)})
			}
		}
		b.columnDefinitions = columns
		return b
	}
}

func WithColumnDefinitions(definitions []ColumnDefinition) optsFn {
	return func(b *Builder) *Builder {
		b.columnDefinitions = definitions
		return b
	}
}

func WithFilters(fs []string) optsFn {
	return func(b *Builder) *Builder {
		b.filters = fs
		return b
	}
}

func WithTagFilters(fs []string) optsFn {
	return func(b *Builder) *Builder {
		b.tagFilters = fs
		return b
	}
}

func WithTagKeyFilters(fs []string) optsFn {
	return func(b *Builder) *Builder {
		b.tagKeyFilters = fs
		return b
	}
}

func WithTagValueFilters(fs []string) optsFn {
	return func(b *Builder) *Builder {
		b.tagValueFilters = fs
		return b
	}
}

func WithIDsOnly(only bool) optsFn {
	return func(b *Builder) *Builder {
		if only {
			b.columnDefinitions = []ColumnDefinition{
				&StringColumnDefinition{Prop: "ID"},
				&StringColumnDefinition{Prop: "Name"},
			}
			b.format = "porcelain"
		}

		return b
	}
}

func WithSortBy(sortingBy ...string) optsFn {
	return func(b *Builder) *Builder {
		indexes, err := resolveSortIndexes(b.columnDefinitions, sortingBy...)
		if err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		}

		b.sort = indexes

		return b
	}
}

func WithReverseSort(r bool) optsFn {
	return func(b *Builder) *Builder {
		b.reverseSort = r
		return b
	}
}

func WithMaxWidth(maxwidth int) optsFn {
	return func(b *Builder) *Builder {
		b.maxwidth = maxwidth
		return b
	}
}

func WithRdfType(rdfType string) optsFn {
	return func(b *Builder) *Builder {
		b.rdfType = rdfType
		return b
	}
}

func WithRootNode(root cloud.Resource) optsFn {
	return func(b *Builder) *Builder {
		b.root = root
		return b
	}
}

func WithNoHeaders(nh bool) optsFn {
	return func(b *Builder) *Builder {
		b.noHeaders = nh
		return b
	}
}

type table [][]interface{}

type fromGraphDisplayer struct {
	sorter
	g                 cloud.GraphAPI
	rdfType           string
	columnDefinitions []ColumnDefinition
	maxwidth          int
	noHeaders         bool
}

func (d *fromGraphDisplayer) setGraph(g cloud.GraphAPI) {
	d.g = g
}

type csvDisplayer struct {
	fromGraphDisplayer
}

func (d *csvDisplayer) Print(w io.Writer) error {
	resources, err := d.g.Find(cloud.NewQuery(d.rdfType))
	if err != nil {
		return err
	}

	if len(d.columnDefinitions) == 0 {
		return nil
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.columnDefinitions))
		}
		for j, h := range d.columnDefinitions {
			values[i][j] = res.Properties()[h.propKey()]
		}
	}

	d.sorter.sort(values)

	var buff bytes.Buffer

	var head []string
	for _, h := range d.columnDefinitions {
		head = append(head, h.title())
	}

	if !d.noHeaders {
		buff.WriteString(strings.Join(head, ",") + "\n")
	}

	for i := range values {
		var props []string
		for j, h := range d.columnDefinitions {
			val := h.format(values[i][j])
			if strings.ContainsAny(val, ",\n\"") {
				val = strings.Replace(val, "\"", "\"\"", -1) // Replace " in val by "" (cf https://tools.ietf.org/html/rfc4180)
				val = "\"" + val + "\""
			}
			props = append(props, val)
		}
		buff.WriteString(strings.Join(props, ",") + "\n")
	}

	_, err = w.Write(buff.Bytes())
	return err
}

type tsvDisplayer struct {
	fromGraphDisplayer
}

func (d *tsvDisplayer) Print(w io.Writer) error {
	color.NoColor = true // as default tabwriter does not play nice with the color library

	resources, err := d.g.Find(cloud.NewQuery(d.rdfType))
	if err != nil {
		return err
	}

	if len(d.columnDefinitions) == 0 {
		return nil
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.columnDefinitions))
		}
		for j, h := range d.columnDefinitions {
			values[i][j] = res.Properties()[h.propKey()]
		}
	}

	d.sorter.sort(values)

	var head []string
	for _, h := range d.columnDefinitions {
		head = append(head, h.title())
	}

	if !d.noHeaders {
		fmt.Fprintln(w, strings.Join(head, "\t"))
	}

	for i := range values {
		var props []string
		for j, h := range d.columnDefinitions {
			props = append(props, h.format(values[i][j]))
		}
		fmt.Fprintln(w, strings.Join(props, "\t"))
	}

	return nil
}

type jsonDisplayer struct {
	fromGraphDisplayer
}

func (d *jsonDisplayer) Print(w io.Writer) error {
	resources, err := d.g.Find(cloud.NewQuery(d.rdfType))
	if err != nil {
		return err
	}

	sort.Slice(resources, func(i, j int) bool { return resources[i].Id() < resources[j].Id() })

	var props []map[string]interface{}
	for _, res := range resources {
		props = append(props, res.Properties())
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")

	return enc.Encode(props)
}

type tableDisplayer struct {
	fromGraphDisplayer
}

func (d *tableDisplayer) Print(w io.Writer) error {
	resources, err := d.g.Find(cloud.NewQuery(d.rdfType))
	if err != nil {
		return err
	}
	if len(resources) == 0 {
		w.Write([]byte("No results found.\n"))
		return nil
	}
	if len(d.columnDefinitions) == 0 {
		w.Write([]byte("No columns to display.\n"))
		return nil
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.columnDefinitions))
		}
		for j, h := range d.columnDefinitions {
			values[i][j] = res.Properties()[h.propKey()]
		}
	}

	d.sorter.sort(values)
	markColumnAsc := -1
	if len(d.sorter.columns()) > 0 {
		markColumnAsc = d.sorter.columns()[0]
	}

	columnsToDisplay := d.columnDefinitions
	maxWidthNoWraping := 1
	if d.maxwidth != 0 {
		columnsToDisplay = []ColumnDefinition{}
		currentWidth := 1 // first border
		for j, h := range d.columnDefinitions {
			var symbol string
			if markColumnAsc == j {
				symbol = d.sorter.symbol()
			}
			colW := colWidth(j, values, h, symbol) + 3 // +3 (tables margin + border)
			if currentWidth+colW > d.maxwidth {
				break
			}
			currentWidth += colW
			maxWidthNoWraping += colWidthNoWraping(j, values, h, symbol) + 3
			columnsToDisplay = append(columnsToDisplay, h)
		}
	}

	table := tablewriter.NewWriter(w)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(tableColWidth)
	if !d.noHeaders {
		var displayHeaders []string
		for i, h := range columnsToDisplay {
			var symbol string
			if markColumnAsc == i {
				symbol = d.sorter.symbol()
			}
			displayHeaders = append(displayHeaders, h.title(symbol))
		}
		table.SetHeader(displayHeaders)
	}

	var enableWraping bool
	if d.maxwidth <= maxWidthNoWraping {
		enableWraping = true
	}

	wraper := autoWraper{maxWidth: autowrapMaxSize, wrappingChar: " "}
	for i := range values {
		var props []string
		for j, h := range columnsToDisplay {
			val := h.format(values[i][j])
			if enableWraping {
				props = append(props, wraper.Wrap(val))
			} else {
				props = append(props, val)
			}

		}
		table.Append(props)
	}

	table.Render()
	if len(columnsToDisplay) < len(d.columnDefinitions) {
		var hiddenColumns []string
		for i := len(columnsToDisplay); i < len(d.columnDefinitions); i++ {
			hiddenColumns = append(hiddenColumns, "'"+d.columnDefinitions[i].title()+"'")
		}
		if len(hiddenColumns) == 1 {
			fmt.Fprint(w, color.New(color.FgRed).SprintfFunc()("Column truncated to fit terminal: %s\n", hiddenColumns[0]))
		} else {
			fmt.Fprint(w, color.New(color.FgRed).SprintfFunc()("Columns truncated to fit terminal: %s\n", strings.Join(hiddenColumns, ", ")))
		}
	}
	return nil
}

type porcelainDisplayer struct {
	fromGraphDisplayer
}

func (d *porcelainDisplayer) Print(w io.Writer) error {
	var types []string

	if d.rdfType == "" {
		for t := range DefaultsColumnDefinitions {
			types = append(types, t)
		}
	} else {
		types = append(types, d.rdfType)
	}

	var values table
	for _, t := range types {
		resources, err := d.g.Find(cloud.NewQuery(t))
		if err != nil {
			return err
		}

		for _, res := range resources {
			var row = make([]interface{}, len(d.columnDefinitions))
			for j, h := range d.columnDefinitions {
				row[j] = res.Properties()[h.propKey()]
			}
			values = append(values, row)
		}
	}

	d.sorter.sort(values)

	var lines []string

	for i := range values {
		for j := range d.columnDefinitions {
			v := values[i][j]
			if v != nil {
				val := fmt.Sprint(v)
				if val != "" {
					lines = append(lines, val)
				}
			}

		}
	}

	_, err := w.Write([]byte(strings.Join(lines, "\n")))
	return err
}

type multiResourcesTableDisplayer struct {
	fromGraphDisplayer
}

func (d *multiResourcesTableDisplayer) Print(w io.Writer) error {
	var values table

	for t, propDefs := range DefaultsColumnDefinitions {
		resources, err := d.g.Find(cloud.NewQuery(t))
		if err != nil {
			return err
		}
		for _, res := range resources {
			for prop, val := range res.Properties() {
				var header ColumnDefinition
				for _, h := range propDefs {
					if h.propKey() == prop {
						header = h
					}
				}
				if header == nil {
					header = &StringColumnDefinition{Prop: prop}
				}
				var row [4]interface{}
				row[0] = t
				row[1] = nameOrID(res)
				row[2] = header.title()
				row[3] = header.format(val)
				values = append(values, row[:])
			}
		}
	}

	ds := defaultSorter{sortBy: []int{0, 1, 2, 3}}
	ds.sort(values)

	table := tablewriter.NewWriter(w)
	table.SetAutoMergeCells(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(tableColWidth)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetHeader([]string{"Type" + ds.symbol(), "Name/Id", "Property", "Value"})

	wraper := autoWraper{maxWidth: autowrapMaxSize, wrappingChar: " "}

	for i := range values {
		row := make([]string, len(values[i]))
		for j := range values[i] {
			row[j] = wraper.Wrap(fmt.Sprint(values[i][j]))
		}
		table.Append(row)
	}

	table.Render()

	return nil
}

type multiResourcesJSONDisplayer struct {
	fromGraphDisplayer
}

func (d *multiResourcesJSONDisplayer) Print(w io.Writer) error {
	var resources []cloud.Resource
	var err error

	all := make(map[string]interface{})
	for t := range DefaultsColumnDefinitions {
		resources, err = d.g.Find(cloud.NewQuery(t))
		if err != nil {
			return err
		}
		var props []map[string]interface{}
		for _, res := range resources {
			props = append(props, res.Properties())
		}
		if len(resources) > 0 {
			all[cloud.PluralizeResource(t)] = props
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")

	return enc.Encode(all)
}

type fromDiffDisplayer struct {
	root cloud.Resource
	diff *graph.Diff
}

func (d *fromDiffDisplayer) SetDiff(diff *graph.Diff) {
	d.diff = diff
}

type diffTableDisplayer struct {
	*fromDiffDisplayer
}

func (d *diffTableDisplayer) Print(w io.Writer) error {
	var values table

	fromCommons := make(map[string]cloud.Resource)
	toCommons := make(map[string]cloud.Resource)
	each := func(res *graph.Resource, distance int) error {
		diff, _ := res.Meta("diff")
		switch diff {
		case "extra":
			values = append(values, []interface{}{
				res.Type(), color.New(color.FgRed).SprintFunc()("- " + nameOrID(res)), "", "",
			})
		default:
			fromCommons[res.Id()] = res
		}
		return nil
	}
	err := d.diff.FromGraph().Accept(&graph.ChildrenVisitor{From: d.root.(*graph.Resource), Each: each})
	if err != nil {
		return err
	}

	each = func(res *graph.Resource, distance int) error {
		meta, _ := res.Meta("diff")
		switch meta {
		case "extra":
			values = append(values, []interface{}{
				res.Type(), color.New(color.FgGreen).SprintFunc()("+ " + nameOrID(res)), "", "",
			})
		default:
			toCommons[res.Id()] = res
		}
		return nil
	}
	err = d.diff.ToGraph().Accept(&graph.ChildrenVisitor{From: d.root.(*graph.Resource), Each: each})
	if err != nil {
		return err
	}

	for _, common := range fromCommons {
		resType := common.Type()
		naming := nameOrID(common)

		if rem, ok := toCommons[common.Id()]; ok {
			added := graph.Subtract(rem.Properties(), common.Properties())
			for k, v := range added {
				values = append(values, []interface{}{
					resType, naming, k, color.New(color.FgGreen).SprintFunc()("+ " + fmt.Sprint(v)),
				})
			}

			deleted := graph.Subtract(common.Properties(), rem.Properties())
			for k, v := range deleted {
				values = append(values, []interface{}{
					resType, naming, k, color.New(color.FgRed).SprintFunc()("- " + fmt.Sprint(v)),
				})
			}
		}
	}

	ds := defaultSorter{sortBy: []int{0, 1, 2, 3}}
	ds.sort(values)

	table := tablewriter.NewWriter(w)
	table.SetAutoMergeCells(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetHeader([]string{"Type" + ds.symbol(), "Name/Id", "Property", "Value"})

	for i := range values {
		row := make([]string, len(values[i]))
		for j := range values[i] {
			row[j] = fmt.Sprint(values[i][j])
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

type diffTreeDisplayer struct {
	*fromDiffDisplayer
}

func (d *diffTreeDisplayer) Print(w io.Writer) error {
	g := graph.NewGraph()

	each := func(res *graph.Resource, distance int) error {
		meta, _ := res.Meta("diff")
		switch meta {
		case "extra", "missing":
			var parents []*graph.Resource
			err := d.diff.MergedGraph().Accept(&graph.ParentsVisitor{From: res, Each: graph.VisitorCollectFunc(&parents)})
			if err != nil {
				return err
			}
			if err := g.AddResource(res); err != nil {
				return err
			}
			previous := res
			for _, parent := range parents {
				if err := g.AddResource(parent); err != nil {
					return err
				}
				if err := g.AddParentRelation(parent, previous); err != nil {
					return err
				}
				previous = parent
			}
		}
		return nil
	}

	err := d.diff.MergedGraph().Accept(&graph.ChildrenVisitor{From: d.root.(*graph.Resource), Each: each, IncludeFrom: true})
	if err != nil {
		return err
	}

	each = func(res *graph.Resource, distance int) error {
		tabs := strings.Repeat("\t", distance)
		diffMeta, _ := res.Meta("diff")
		switch diffMeta {
		case "extra":
			color.Set(color.FgGreen)
			fmt.Fprintf(w, "+%s%s, %s\n", tabs, res.Type(), res.Id())
			color.Unset()
		case "missing":
			color.Set(color.FgRed)
			fmt.Fprintf(w, "-%s%s, %s\n", tabs, res.Type(), res.Id())
			color.Unset()
		default:
			fmt.Fprintf(w, "%s%s, %s\n", tabs, res.Type(), res.Id())
		}
		return nil
	}

	err = g.Accept(&graph.ChildrenVisitor{From: d.root.(*graph.Resource), Each: each, IncludeFrom: true})
	if err != nil {
		return err
	}

	return nil
}

type defaultSorter struct {
	sortBy     []int
	descending bool
}

func (d *defaultSorter) sort(lines table) {
	var compare func(i, j int) bool
	if d.descending {
		compare = func(j, i int) bool {
			for _, col := range d.sortBy {
				if reflect.DeepEqual(lines[i][col], lines[j][col]) {
					continue
				}
				return valueLowerOrEqual(lines[i][col], lines[j][col])
			}
			return false
		}
	} else {
		compare = func(i, j int) bool {
			for _, col := range d.sortBy {
				if reflect.DeepEqual(lines[i][col], lines[j][col]) {
					continue
				}
				return valueLowerOrEqual(lines[i][col], lines[j][col])
			}
			return false
		}
	}
	sort.Slice(lines, compare)
}

func (d *defaultSorter) columns() []int {
	return d.sortBy
}

func (d *defaultSorter) symbol() string {
	if d.descending {
		return " ▼"
	}
	return " ▲"
}

func valueLowerOrEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if aTyp, bTyp := reflect.TypeOf(a), reflect.TypeOf(b); aTyp != nil &&
		aTyp.Comparable() && bTyp != nil && bTyp.Comparable() && a == b {
		return true
	}
	if a == nil {
		return true
	}
	if b == nil {
		return false
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		panic(fmt.Sprintf("can not compare values of type %T and %T", a, b))
	}
	switch a.(type) {
	case int:
		aa := a.(int)
		bb := b.(int)
		return aa <= bb
	case float64:
		aa := a.(float64)
		bb := b.(float64)
		return aa <= bb
	case string:
		aa := a.(string)
		bb := b.(string)
		return aa <= bb
	case time.Time:
		aa := a.(time.Time)
		bb := b.(time.Time)
		return aa.After(bb)
	case []string, []int:
		return fmt.Sprint(a) <= fmt.Sprint(b)
	default:
		panic(fmt.Sprintf("can not compare values of type %T", a))
	}
}

func resolveSortIndexes(headers []ColumnDefinition, sortingBy ...string) ([]int, error) {
	sortBy := []string{"id"}
	if len(sortingBy) > 0 {
		sortBy = sortingBy
	}

	normalized := make(map[string]int)
	for i, h := range headers {
		normalized[strings.ToLower(h.title())] = i
	}

	var ids []int
	for _, t := range sortBy {
		id, ok := normalized[strings.ToLower(t)]
		if !ok && strings.ToLower(t) != "id" {
			return ids, fmt.Errorf("Invalid column name '%s'", t)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func colWidth(j int, t table, h ColumnDefinition, sortSymbol string) int {
	max := tablewriter.DisplayWidth(h.title(sortSymbol))
	wraper := autoWraper{maxWidth: autowrapMaxSize, wrappingChar: " "}
	for i := range t {
		val := wraper.Wrap(h.format(t[i][j]))
		valLen := tablewriter.DisplayWidth(val)
		if valLen > tableColWidth {
			if tableColWidth > max {
				max = tableColWidth
			}
		}
		lines, _ := tablewriter.WrapString(val, tableColWidth)
		for _, line := range lines {
			width := tablewriter.DisplayWidth(line)
			if width > max {
				max = width
			}
		}
	}
	return max
}

func colWidthNoWraping(j int, t table, h ColumnDefinition, sortSymbol string) int {
	max := tablewriter.DisplayWidth(h.title(sortSymbol))
	for i := range t {
		val := h.format(t[i][j])
		valLen := tablewriter.DisplayWidth(val)
		if valLen > tableColWidth {
			if tableColWidth > max {
				max = tableColWidth
			}
		}
		lines, _ := tablewriter.WrapString(val, tableColWidth)
		for _, line := range lines {
			width := tablewriter.DisplayWidth(line)
			if width > max {
				max = width
			}
		}
	}
	return max
}

func nameOrID(res cloud.Resource) string {
	if name, ok := res.Property("Name"); ok && name != "" {
		return fmt.Sprint(name)
	}
	if id, ok := res.Property("Id"); ok && id != "" {
		return fmt.Sprint(id)
	}
	return res.Id()
}
