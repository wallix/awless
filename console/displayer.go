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
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

type Displayer interface {
	Print(io.Writer) error
}

type sorter interface {
	sort(table)
	columns() []int
}

type Builder struct {
	filters    []string
	headers    []ColumnDefinition
	format     string
	rdfType    string
	sort       []int
	maxwidth   int
	dataSource interface{}
	root       *graph.Resource
}

func (b *Builder) SetSource(i interface{}) *Builder {
	b.dataSource = i
	return b
}

func (b *Builder) buildGraphFilters() (funcs []graph.FilterFn) {
	for _, f := range b.filters {
		splits := strings.SplitN(f, "=", 2)
		if len(splits) == 2 {
			name, val := strings.TrimSpace(strings.Title(splits[0])), strings.TrimSpace(splits[1])
			key := ColumnDefinitions(b.headers).resolveKey(name)

			if key != "" {
				funcs = append(funcs, graph.BuildPropertyFilterFunc(key, val))
			}
		}

	}
	return
}

func (b *Builder) Build() Displayer {
	base := fromGraphDisplayer{sorter: &defaultSorter{sortBy: b.sort}, rdfType: b.rdfType, headers: b.headers, maxwidth: b.maxwidth}

	switch b.dataSource.(type) {
	case *graph.Graph:
		gph := b.dataSource.(*graph.Graph)
		filteredGraph, _ := gph.Filter(b.rdfType, b.buildGraphFilters()...)

		if b.rdfType == "" {
			switch b.format {
			case "table":
				dis := &multiResourcesTableDisplayer{base}
				dis.setGraph(gph)
				return dis
			case "json":
				dis := &multiResourcesJSONDisplayer{base}
				dis.setGraph(gph)
				return dis
			case "porcelain":
				dis := &porcelainDisplayer{base}
				dis.setGraph(gph)
				return dis
			default:
				fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'table'\n", b.format)
				dis := &multiResourcesTableDisplayer{base}
				dis.setGraph(gph)
				return dis
			}
		}
		switch b.format {
		case "csv":
			dis := &csvDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis
		case "tsv":
			dis := &tsvDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis
		case "json":
			dis := &jsonDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis
		case "porcelain":
			dis := &porcelainDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis
		case "table":
			dis := &tableDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis
		default:
			fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'table'\n", b.format)
			dis := &tableDisplayer{base}
			dis.setGraph(filteredGraph)
			return dis
		}
	case *graph.Resource:
		dis := &tableResourceDisplayer{headers: b.headers}
		dis.SetResource(b.dataSource.(*graph.Resource))
		return dis
	case *graph.Diff:
		base := fromDiffDisplayer{root: b.root}
		switch b.format {
		case "tree":
			dis := &diffTreeDisplayer{&base}
			dis.SetDiff(b.dataSource.(*graph.Diff))
			return dis
		case "table":
			dis := &diffTableDisplayer{&base}
			dis.SetDiff(b.dataSource.(*graph.Diff))
			return dis
		default:
			fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'tree'\n", b.format)
			dis := &diffTreeDisplayer{&base}
			dis.SetDiff(b.dataSource.(*graph.Diff))
			return dis
		}
	}

	return nil
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

func WithIDsOnly(only bool) optsFn {
	return func(b *Builder) *Builder {
		if only {
			b.headers = []ColumnDefinition{
				&StringColumnDefinition{Prop: "Id"},
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

	lines = append(lines, strings.Join(head, ", "))

	for i := range values {
		var props []string
		for j, h := range d.headers {
			props = append(props, h.format(values[i][j]))
		}
		lines = append(lines, strings.Join(props, ", "))
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

	tabw := tabwriter.NewWriter(w, 8, 8, 0, '\t', 0)

	fmt.Fprintln(tabw, strings.Join(head, "\t"))

	for i := range values {
		var props []string
		for j, h := range d.headers {
			props = append(props, h.format(values[i][j]))
		}
		fmt.Fprintln(tabw, strings.Join(props, "\t"))
	}

	return tabw.Flush()
}

type jsonDisplayer struct {
	fromGraphDisplayer
}

func (d *jsonDisplayer) Print(w io.Writer) error {
	resources, err := d.g.GetAllResources(d.rdfType)
	if err != nil {
		return err
	}

	sort.Sort(graph.ResourceById(resources))

	var props []graph.Properties
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
	columnsToDisplay := d.headers
	if d.maxwidth != 0 {
		columnsToDisplay = []ColumnDefinition{}
		currentWidth := 0
		for j, h := range d.headers {
			colW := t(j, values, h) + 3 // +3 (tables margin)
			if currentWidth+colW > d.maxwidth {
				break
			}
			currentWidth += colW
			columnsToDisplay = append(columnsToDisplay, h)
		}
	}

	markColumnAsc := -1
	if len(d.sorter.columns()) > 0 {
		markColumnAsc = d.sorter.columns()[0]
	}

	table := tablewriter.NewWriter(w)
	var displayHeaders []string
	for i, h := range columnsToDisplay {
		displayHeaders = append(displayHeaders, h.title(i == markColumnAsc))
	}
	table.SetHeader(displayHeaders)

	for i := range values {
		var props []string
		for j, h := range columnsToDisplay {
			props = append(props, h.format(values[i][j]))
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
		var props []graph.Properties
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
			added := rem.Properties.Subtract(common.Properties)
			for k, v := range added {
				values = append(values, []interface{}{
					resType, naming, k, color.New(color.FgGreen).SprintFunc()("+ " + fmt.Sprint(v)),
				})
			}

			deleted := common.Properties.Subtract(rem.Properties)
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
			g.AddResource(res)
			previous := res
			for _, parent := range parents {
				g.AddResource(parent)
				g.AddParentRelation(parent, previous)
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
		var tabs bytes.Buffer
		for i := 0; i < distance; i++ {
			tabs.WriteByte('\t')
		}

		switch res.Meta["diff"] {
		case "extra":
			color.Set(color.FgGreen)
			fmt.Fprintf(w, "+%s%s, %s\n", tabs.String(), res.Type(), res.Id())
			color.Unset()
		case "missing":
			color.Set(color.FgRed)
			fmt.Fprintf(w, "-%s%s, %s\n", tabs.String(), res.Type(), res.Id())
			color.Unset()
		default:
			fmt.Fprintf(w, "%s%s, %s\n", tabs.String(), res.Type(), res.Id())
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
		return aa.Before(bb)
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

func t(j int, t table, h ColumnDefinition) int {
	w := 0
	for i := range t {
		words := get_words_from(h.format(t[i][j]))
		for _, word := range words {
			c := utf8.RuneCountInString(word)
			if c > w {
				w = c
			}
		}
	}
	return w
}

func get_words_from(text string) []string {
	return regexp.MustCompile("(\\b[^\\s]+\\b)").FindAllString(text, -1)
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
