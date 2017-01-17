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
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

type GraphDisplayer interface {
	sorter
	Print(io.Writer) error
	SetGraph(*rdf.Graph)
}

type sorter interface {
	sort(table)
	columns() []int
}

type Options struct {
	RdfType  rdf.ResourceType
	Format   string
	SortBy   []string
	MaxWidth int
}

func BuildGraphDisplayer(headers []ColumnDefinition, opts Options) GraphDisplayer {
	titlesIds := make(map[string]int)
	for i, h := range headers {
		titlesIds[strings.ToLower(h.title(false))] = i
	}
	sortBy := []string{"Id"}
	if len(opts.SortBy) > 0 {
		sortBy = opts.SortBy
	}

	sortIds, err := titlesToIDs(titlesIds, sortBy)
	if err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
	}

	switch opts.Format {
	case "csv":
		return &csvGraphDisplayer{sorter: &defaultSorter{sortBy: sortIds}, rdfType: opts.RdfType, headers: headers}
	case "table":
		return &tableGraphDisplayer{sorter: &defaultSorter{sortBy: sortIds}, rdfType: opts.RdfType, headers: headers, maxwidth: opts.MaxWidth}
	case "porcelain":
		return &porcelainGraphDisplayer{sorter: &defaultSorter{sortBy: sortIds}, rdfType: opts.RdfType, headers: headers}
	default:
		fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'table'\n", opts.Format)
		return &tableGraphDisplayer{sorter: &defaultSorter{sortBy: sortIds}, rdfType: opts.RdfType, headers: headers, maxwidth: opts.MaxWidth}
	}
}

type table [][]interface{}

type csvGraphDisplayer struct {
	sorter
	g       *rdf.Graph
	rdfType rdf.ResourceType
	headers []ColumnDefinition
}

func (d *csvGraphDisplayer) Print(w io.Writer) error {
	resources, err := aws.LoadResourcesFromGraph(d.g, d.rdfType)
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

func (d *csvGraphDisplayer) SetGraph(g *rdf.Graph) {
	d.g = g
}

type tableGraphDisplayer struct {
	sorter
	g        *rdf.Graph
	rdfType  rdf.ResourceType
	headers  []ColumnDefinition
	maxwidth int
}

func (d *tableGraphDisplayer) Print(w io.Writer) error {
	resources, err := aws.LoadResourcesFromGraph(d.g, d.rdfType)
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
			colW := columnWidth(j, values, h) + 2 // +2 (tables margin)
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

func (d *tableGraphDisplayer) SetGraph(g *rdf.Graph) {
	d.g = g
}

type porcelainGraphDisplayer struct {
	sorter
	g       *rdf.Graph
	rdfType rdf.ResourceType
	headers []ColumnDefinition
}

func (d *porcelainGraphDisplayer) Print(w io.Writer) error {
	resources, err := aws.LoadResourcesFromGraph(d.g, d.rdfType)
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

func (d *porcelainGraphDisplayer) SetGraph(g *rdf.Graph) {
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

func titlesToIDs(mapping map[string]int, titles []string) ([]int, error) {
	var ids []int
	for _, t := range titles {
		id, ok := mapping[strings.ToLower(t)]
		if !ok {
			return ids, fmt.Errorf("Invalid column name '%s'", t)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func columnWidth(j int, t table, h ColumnDefinition) int {
	w := 0
	for i := range t {
		c := utf8.RuneCountInString(h.format(t[i][j]))
		if c > w {
			w = c
		}
	}
	return w
}
