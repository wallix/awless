package display

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/wallix/awless/graph"
)

type Displayer interface {
	sorter
	Print(io.Writer) error
}

type GraphDisplayer interface {
	Displayer
	SetGraph(*graph.Graph)
}

type sorter interface {
	sort(table)
	columns() []int
}

type Builder struct {
	headers  []ColumnDefinition
	format   string
	rdfType  graph.ResourceType
	sort     []int
	maxwidth int
	source   interface{}
}

func (b *Builder) SetSource(i interface{}) *Builder {
	b.source = i
	return b
}

func (b *Builder) Build() Displayer {
	base := fromGraphDisplayer{sorter: &defaultSorter{sortBy: b.sort}, rdfType: b.rdfType, headers: b.headers, maxwidth: b.maxwidth}

	switch b.source.(type) {
	case *graph.Graph:
		switch b.format {
		case "csv":
			dis := &csvDisplayer{base}
			dis.SetGraph(b.source.(*graph.Graph))
			return dis
		case "porcelain":
			dis := &porcelainDisplayer{base}
			dis.SetGraph(b.source.(*graph.Graph))
			return dis
		case "table":
			dis := &tableDisplayer{base}
			dis.SetGraph(b.source.(*graph.Graph))
			return dis
		default:
			fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'table'\n", b.format)
			dis := &tableDisplayer{base}
			dis.SetGraph(b.source.(*graph.Graph))
			return dis
		}
	case *graph.Resource:
		dis := &tableResourceDisplayer{headers: b.headers}
		dis.SetResource(b.source.(*graph.Resource))
		return dis
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

func WithIDsOnly(only bool) optsFn {
	return func(b *Builder) *Builder {
		if only {
			b.headers = []ColumnDefinition{
				StringColumnDefinition{Prop: "Id"},
				StringColumnDefinition{Prop: "Name"},
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

func WithRdfType(rdfType graph.ResourceType) optsFn {
	return func(b *Builder) *Builder {
		b.rdfType = rdfType
		return b
	}
}

type table [][]interface{}

type fromGraphDisplayer struct {
	sorter
	g        *graph.Graph
	rdfType  graph.ResourceType
	headers  []ColumnDefinition
	maxwidth int
}

type csvDisplayer struct {
	fromGraphDisplayer
}

func (d *csvDisplayer) Print(w io.Writer) error {
	resources, err := graph.LoadResourcesFromGraph(d.g, d.rdfType)
	if err != nil {
		return err
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.headers))
		}
		for j, h := range d.headers {
			values[i][j] = res.Properties()[h.propKey()]
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

func (d *csvDisplayer) SetGraph(g *graph.Graph) {
	d.g = g
}

type tableDisplayer struct {
	fromGraphDisplayer
}

func (d *tableDisplayer) Print(w io.Writer) error {
	resources, err := graph.LoadResourcesFromGraph(d.g, d.rdfType)
	if err != nil {
		return err
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.headers))
		}
		for j, h := range d.headers {
			values[i][j] = res.Properties()[h.propKey()]
		}
	}

	d.sorter.sort(values)

	columnsToDisplay := d.headers
	if d.maxwidth != 0 {
		columnsToDisplay = []ColumnDefinition{}
		currentWidth := 0
		for j, h := range d.headers {
			colW := t(j, values, h) + 2 // +2 (tables margin)
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

func (d *tableDisplayer) SetGraph(g *graph.Graph) {
	d.g = g
}

type porcelainDisplayer struct {
	fromGraphDisplayer
}

func (d *porcelainDisplayer) Print(w io.Writer) error {
	resources, err := graph.LoadResourcesFromGraph(d.g, d.rdfType)
	if err != nil {
		return err
	}

	values := make(table, len(resources))
	for i, res := range resources {
		if v := values[i]; v == nil {
			values[i] = make([]interface{}, len(d.headers))
		}
		for j, h := range d.headers {
			values[i][j] = res.Properties()[h.propKey()]
		}
	}

	d.sorter.sort(values)

	var lines []string

	for i := range values {
		for j, h := range d.headers {
			val := h.format(values[i][j])
			if val != "" {
				lines = append(lines, val)
			}
		}
	}

	_, err = w.Write([]byte(strings.Join(lines, "\n")))
	return err
}

func (d *porcelainDisplayer) SetGraph(g *graph.Graph) {
	d.g = g
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
		if !ok {
			return ids, fmt.Errorf("Invalid column name '%s'", t)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func t(j int, t table, h ColumnDefinition) int {
	w := 0
	for i := range t {
		c := utf8.RuneCountInString(h.format(t[i][j]))
		if c > w {
			w = c
		}
	}
	return w
}
