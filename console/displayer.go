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
}

type Builder struct {
	filters         []string
	tagFilters      []string
	tagKeyFilters   []string
	tagValueFilters []string
	headers         []ColumnDefinition
	format          string
	rdfType         string
	sort            []int
	maxwidth        int
	dataSource      interface{}
	root            *graph.Resource
}

func (b *Builder) SetSource(i interface{}) *Builder {
	b.dataSource = i
	return b
}

func (b *Builder) buildGraphFilters() (funcs []graph.FilterFn, err error) {
	for _, f := range b.filters {
		splits := strings.SplitN(f, "=", 2)
		if len(splits) == 2 {
			name, val := strings.TrimSpace(strings.Title(splits[0])), strings.TrimSpace(splits[1])
			key := ColumnDefinitions(b.headers).resolveKey(name)

			if key != "" {
				funcs = append(funcs, graph.BuildPropertyFilterFunc(key, val))
			} else {
				var allowed []string
				for _, h := range b.headers {
					allowed = append(allowed, h.propKey())
				}
				err = fmt.Errorf("Invalid filter key '%s'. Expecting any of: %s. (Note: filter keys/values are case insensitive)", name, strings.Join(allowed, ", "))
			}
		}

	}
	return
}

func (b *Builder) buildGraphTagFilters() (funcs []graph.FilterFn) {
	for _, f := range b.tagFilters {
		splits := strings.SplitN(f, "=", 2)
		if len(splits) == 2 {
			key, val := strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
			funcs = append(funcs, graph.BuildTagFilterFunc(key, val))
		}
	}
	return
}

func (b *Builder) buildGraphTagKeyFilters() (funcs []graph.FilterFn) {
	for _, k := range b.tagKeyFilters {
		funcs = append(funcs, graph.BuildTagKeyFilterFunc(k))
	}
	return
}

func (b *Builder) buildGraphTagValueFilters() (funcs []graph.FilterFn) {
	for _, v := range b.tagValueFilters {
		funcs = append(funcs, graph.BuildTagValueFilterFunc(v))
	}
	return
}

func (b *Builder) Build() (Displayer, error) {
	base := fromGraphDisplayer{sorter: &defaultSorter{sortBy: b.sort}, rdfType: b.rdfType, headers: b.headers, maxwidth: b.maxwidth}

	switch b.dataSource.(type) {
	case *graph.Graph:
		if b.rdfType == "" {
			gph := b.dataSource.(*graph.Graph)
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

		filteredGraph := b.dataSource.(*graph.Graph)

		if filters, err := b.buildGraphFilters(); len(filters) > 0 && err == nil {
			filteredGraph, err = filteredGraph.Filter(b.rdfType, filters...)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		var ferr error
		if filters := b.buildGraphTagFilters(); len(filters) > 0 {
			filteredGraph, ferr = filteredGraph.Filter(b.rdfType, filters...)
			if ferr != nil {
				return nil, ferr
			}
		}
		if filters := b.buildGraphTagKeyFilters(); len(filters) > 0 {
			filteredGraph, ferr = filteredGraph.OrFilter(b.rdfType, filters...)
			if ferr != nil {
				return nil, ferr
			}
		}
		if filters := b.buildGraphTagValueFilters(); len(filters) > 0 {
			filteredGraph, ferr = filteredGraph.OrFilter(b.rdfType, filters...)
			if ferr != nil {
				return nil, ferr
			}
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
	case *graph.Resource:
		dis := &tableResourceDisplayer{headers: b.headers, maxwidth: b.maxwidth}
		dis.SetResource(b.dataSource.(*graph.Resource))
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

	if len(b.headers) == 0 {
		b.headers = DefaultsColumnDefinitions[b.rdfType]
	}

	return b
}

func WithFormat(format string) optsFn {
	return func(b *Builder) *Builder {
		b.format = format
		return b
	}
}

func WithHeaders(h []ColumnDefinition) optsFn {
	return func(b *Builder) *Builder {
		b.headers = h
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
			b.headers = []ColumnDefinition{
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
		indexes, err := resolveSortIndexes(b.headers, sortingBy...)
		if err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		}

		b.sort = indexes

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

func WithRootNode(root *graph.Resource) optsFn {
	return func(b *Builder) *Builder {
		b.root = root
		return b
	}
}

type table [][]interface{}

type fromGraphDisplayer struct {
	sorter
	g        *graph.Graph
	rdfType  string
	headers  []ColumnDefinition
	maxwidth int
}

func (d *fromGraphDisplayer) setGraph(g *graph.Graph) {
	d.g = g
}

type csvDisplayer struct {
	fromGraphDisplayer
}

func (d *csvDisplayer) Print(w io.Writer) error {
	resources, err := d.g.GetAllResources(d.rdfType)
	if err != nil {
		return err
	}

	if len(d.headers) == 0 {
		return nil
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.headers))
		}
		for j, h := range d.headers {
			values[i][j] = res.Properties[h.propKey()]
		}
	}

	d.sorter.sort(values)

	var lines []string

	var head []string
	for _, h := range d.headers {
		head = append(head, h.title(false))
	}

	lines = append(lines, strings.Join(head, ","))

	for i := range values {
		var props []string
		for j, h := range d.headers {
			props = append(props, h.format(values[i][j]))
		}
		lines = append(lines, strings.Join(props, ","))
	}

	_, err = w.Write([]byte(strings.Join(lines, "\n")))
	return err
}

type tsvDisplayer struct {
	fromGraphDisplayer
}

func (d *tsvDisplayer) Print(w io.Writer) error {
	color.NoColor = true // as default tabwriter does not play nice with the color library

	resources, err := d.g.GetAllResources(d.rdfType)
	if err != nil {
		return err
	}

	if len(d.headers) == 0 {
		return nil
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.headers))
		}
		for j, h := range d.headers {
			values[i][j] = res.Properties[h.propKey()]
		}
	}

	d.sorter.sort(values)

	var head []string
	for _, h := range d.headers {
		head = append(head, h.title(false))
	}

	fmt.Fprintln(w, strings.Join(head, "\t"))

	for i := range values {
		var props []string
		for j, h := range d.headers {
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
	resources, err := d.g.GetAllResources(d.rdfType)
	if err != nil {
		return err
	}

	sort.Slice(resources, func(i, j int) bool { return resources[i].Id() < resources[j].Id() })

	var props []map[string]interface{}
	for _, res := range resources {
		props = append(props, res.Properties)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")

	return enc.Encode(props)
}

type tableDisplayer struct {
	fromGraphDisplayer
}

func (d *tableDisplayer) Print(w io.Writer) error {
	resources, err := d.g.GetAllResources(d.rdfType)
	if err != nil {
		return err
	}
	if len(resources) == 0 {
		w.Write([]byte("No results found.\n"))
		return nil
	}
	if len(d.headers) == 0 {
		w.Write([]byte("No columns to display.\n"))
		return nil
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.headers))
		}
		for j, h := range d.headers {
			values[i][j] = res.Properties[h.propKey()]
		}
	}

	d.sorter.sort(values)
	markColumnAsc := -1
	if len(d.sorter.columns()) > 0 {
		markColumnAsc = d.sorter.columns()[0]
	}

	columnsToDisplay := d.headers
	if d.maxwidth != 0 {
		columnsToDisplay = []ColumnDefinition{}
		currentWidth := 1 // first border
		for j, h := range d.headers {
			colW := colWidth(j, values, h, j == markColumnAsc) + 3 // +3 (tables margin + border)
			if currentWidth+colW > d.maxwidth {
				break
			}
			currentWidth += colW
			columnsToDisplay = append(columnsToDisplay, h)
		}
	}

	table := tablewriter.NewWriter(w)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(tableColWidth)
	var displayHeaders []string
	for i, h := range columnsToDisplay {
		displayHeaders = append(displayHeaders, h.title(i == markColumnAsc))
	}
	table.SetHeader(displayHeaders)

	wraper := autoWraper{maxWidth: autowrapMaxSize, wrappingChar: " "}
	for i := range values {
		var props []string
		for j, h := range columnsToDisplay {
			val := h.format(values[i][j])
			props = append(props, wraper.Wrap(val))
		}
		table.Append(props)
	}

	table.Render()
	if len(columnsToDisplay) < len(d.headers) {
		var hiddenColumns []string
		for i := len(columnsToDisplay); i < len(d.headers); i++ {
			hiddenColumns = append(hiddenColumns, "'"+d.headers[i].title(false)+"'")
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
		resources, err := d.g.GetAllResources(t)
		if err != nil {
			return err
		}

		for _, res := range resources {
			var row = make([]interface{}, len(d.headers))
			for j, h := range d.headers {
				row[j] = res.Properties[h.propKey()]
			}
			values = append(values, row)
		}
	}

	d.sorter.sort(values)

	var lines []string

	for i := range values {
		for j := range d.headers {
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
		resources, err := d.g.GetAllResources(t)
		if err != nil {
			return err
		}
		for _, res := range resources {
			for prop, val := range res.Properties {
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
				row[2] = header.title(false)
				row[3] = header.format(val)
				values = append(values, row[:])
			}
		}
	}
	sort.Sort(byCols{table: values, sortBy: []int{0, 1, 2, 3}})

	table := tablewriter.NewWriter(w)
	table.SetAutoMergeCells(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(tableColWidth)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetHeader([]string{"Type" + ascSymbol, "Name/Id", "Property", "Value"})

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
	var resources []*graph.Resource
	var err error

	all := make(map[string]interface{})
	for t := range DefaultsColumnDefinitions {
		resources, err = d.g.GetAllResources(t)
		if err != nil {
			return err
		}
		var props []map[string]interface{}
		for _, res := range resources {
			props = append(props, res.Properties)
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
	root *graph.Resource
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

	fromCommons := make(map[string]*graph.Resource)
	toCommons := make(map[string]*graph.Resource)
	each := func(res *graph.Resource, distance int) error {
		switch res.Meta["diff"] {
		case "extra":
			values = append(values, []interface{}{
				res.Type(), color.New(color.FgRed).SprintFunc()("- " + nameOrID(res)), "", "",
			})
		default:
			fromCommons[res.Id()] = res
		}
		return nil
	}
	err := d.diff.FromGraph().Accept(&graph.ChildrenVisitor{From: d.root, Each: each})
	if err != nil {
		return err
	}

	each = func(res *graph.Resource, distance int) error {
		switch res.Meta["diff"] {
		case "extra":
			values = append(values, []interface{}{
				res.Type(), color.New(color.FgGreen).SprintFunc()("+ " + nameOrID(res)), "", "",
			})
		default:
			toCommons[res.Id()] = res
		}
		return nil
	}
	err = d.diff.ToGraph().Accept(&graph.ChildrenVisitor{From: d.root, Each: each})
	if err != nil {
		return err
	}

	for _, common := range fromCommons {
		resType := common.Type()
		naming := nameOrID(common)

		if rem, ok := toCommons[common.Id()]; ok {
			added := graph.Subtract(rem.Properties, common.Properties)
			for k, v := range added {
				values = append(values, []interface{}{
					resType, naming, k, color.New(color.FgGreen).SprintFunc()("+ " + fmt.Sprint(v)),
				})
			}

			deleted := graph.Subtract(common.Properties, rem.Properties)
			for k, v := range deleted {
				values = append(values, []interface{}{
					resType, naming, k, color.New(color.FgRed).SprintFunc()("- " + fmt.Sprint(v)),
				})
			}
		}
	}

	sort.Sort(byCols{table: values, sortBy: []int{0, 1, 2, 3}})
	table := tablewriter.NewWriter(w)
	table.SetAutoMergeCells(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetHeader([]string{"Type" + ascSymbol, "Name/Id", "Property", "Value"})

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
		switch res.Meta["diff"] {
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

	err := d.diff.MergedGraph().Accept(&graph.ChildrenVisitor{From: d.root, Each: each, IncludeFrom: true})
	if err != nil {
		return err
	}

	each = func(res *graph.Resource, distance int) error {
		tabs := strings.Repeat("\t", distance)

		switch res.Meta["diff"] {
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

	err = g.Accept(&graph.ChildrenVisitor{From: d.root, Each: each, IncludeFrom: true})
	if err != nil {
		return err
	}

	return nil
}

type defaultSorter struct {
	sortBy []int
}

func (d *defaultSorter) sort(lines table) {
	sort.Sort(byCols{table: lines, sortBy: d.sortBy})
}

func (d *defaultSorter) columns() []int {
	return d.sortBy
}

type byCols struct {
	table  table
	sortBy []int
}

func (b byCols) Len() int { return len(b.table) }
func (b byCols) Swap(i, j int) {
	b.table[i], b.table[j] = b.table[j], b.table[i]
}
func (b byCols) Less(i, j int) bool {
	for _, col := range b.sortBy {
		if reflect.DeepEqual(b.table[i][col], b.table[j][col]) {
			continue
		}
		return valueLowerOrEqual(b.table[i][col], b.table[j][col])
	}
	return false
}

func valueLowerOrEqual(a, b interface{}) bool {
	if a == b {
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
		normalized[strings.ToLower(h.title(false))] = i
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

func colWidth(j int, t table, h ColumnDefinition, hasSortSign bool) int {
	max := tablewriter.DisplayWidth(h.title(hasSortSign))
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

func nameOrID(res *graph.Resource) string {
	if name, ok := res.Properties["Name"]; ok && name != "" {
		return fmt.Sprint(name)
	}
	if id, ok := res.Properties["Id"]; ok && id != "" {
		return fmt.Sprint(id)
	}
	return res.Id()
}
