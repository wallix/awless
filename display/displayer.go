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
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/olekukonko/tablewriter"
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
	headers  []ColumnDefinition
	format   string
	rdfType  graph.ResourceType
	sort     []int
	maxwidth int
	dataSource   interface{}
	root     *node.Node
}

func (b *Builder) SetSource(i interface{}) *Builder {
	b.dataSource = i
	return b
}

func (b *Builder) Build() Displayer {
	base := fromGraphDisplayer{sorter: &defaultSorter{sortBy: b.sort}, rdfType: b.rdfType, headers: b.headers, maxwidth: b.maxwidth}

	switch b.dataSource.(type) {
	case *graph.Graph:
		if b.rdfType == "" {
			switch b.format {
			case "table":
				dis := &multiResourcesTableDisplayer{base}
				dis.setGraph(b.dataSource.(*graph.Graph))
				return dis
			case "porcelain":
				dis := &porcelainDisplayer{base}
				dis.setGraph(b.dataSource.(*graph.Graph))
				return dis
			default:
				fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'table'\n", b.format)
				dis := &multiResourcesTableDisplayer{base}
				dis.setGraph(b.dataSource.(*graph.Graph))
				return dis
			}
		}
		switch b.format {
		case "csv":
			dis := &csvDisplayer{base}
			dis.setGraph(b.dataSource.(*graph.Graph))
			return dis
		case "porcelain":
			dis := &porcelainDisplayer{base}
			dis.setGraph(b.dataSource.(*graph.Graph))
			return dis
		case "table":
			dis := &tableDisplayer{base}
			dis.setGraph(b.dataSource.(*graph.Graph))
			return dis
		default:
			fmt.Fprintf(os.Stderr, "unknown format '%s', display as 'table'\n", b.format)
			dis := &tableDisplayer{base}
			dis.setGraph(b.dataSource.(*graph.Graph))
			return dis
		}
	case *graph.Resource:
		dis := &tableResourceDisplayer{headers: b.headers}
		dis.SetResource(b.dataSource.(*graph.Resource))
		return dis
	case *graph.Diff:
		dis := &diffTableDisplayer{root: b.root}
		dis.SetDiff(b.dataSource.(*graph.Diff))
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

func WithRdfType(rdfType graph.ResourceType) optsFn {
	return func(b *Builder) *Builder {
		b.rdfType = rdfType
		return b
	}
}

func WithRootNode(root *node.Node) optsFn {
	return func(b *Builder) *Builder {
		b.root = root
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

func (d *fromGraphDisplayer) setGraph(g *graph.Graph) {
	d.g = g
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

type porcelainDisplayer struct {
	fromGraphDisplayer
}

func (d *porcelainDisplayer) Print(w io.Writer) error {
	var types []graph.ResourceType
	if d.rdfType == "" {
		for t := range DefaultsColumnDefinitions {
			types = append(types, t)
		}
	} else {
		types = append(types, d.rdfType)
	}

	var values table
	for _, t := range types {
		resources, err := graph.LoadResourcesFromGraph(d.g, t)
		if err != nil {
			return err
		}

		for _, res := range resources {
			var row = make([]interface{}, len(d.headers))
			for j, h := range d.headers {
				row[j] = res.Properties()[h.propKey()]
			}
			values = append(values, row)
		}
	}

	d.sorter.sort(values)

	var lines []string

	for i := range values {
		for j, _ := range d.headers {
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
		resources, err := graph.LoadResourcesFromGraph(d.g, t)
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
				row[0] = t.String()
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

type diffTableDisplayer struct {
	root *node.Node
	diff *graph.Diff
}

func (d *diffTableDisplayer) Print(w io.Writer) error {
	var values table
	err := d.diff.FullGraph().VisitUnique(d.root, func(g *graph.Graph, n *node.Node, distance int) error {
		var lit *literal.Literal
		diffTriples, err := g.TriplesInDiff(n)
		if len(diffTriples) > 0 && err == nil {
			lit, _ = diffTriples[0].Object().Literal()
		}
		nCommon, nInserted, nDeleted := graph.InitFromRdfNode(n), graph.InitFromRdfNode(n), graph.InitFromRdfNode(n)

		err = nCommon.UnmarshalFromGraph(&graph.Graph{d.diff.CommonGraph()})
		if err != nil {
			return err
		}

		err = nInserted.UnmarshalFromGraph(&graph.Graph{d.diff.InsertedGraph()})
		if err != nil {
			return err
		}

		err = nDeleted.UnmarshalFromGraph(&graph.Graph{d.diff.DeletedGraph()})
		if err != nil {
			return err
		}

		var displayProperties, propsChanges, rNew bool
		var rName string

		var litString string
		if lit != nil {
			litString, _ = lit.Text()
		}

		switch litString {
		case "extra":
			rNew = true
			rName = nameOrID(nInserted)
		case "missing":
			rName = nameOrID(nDeleted)
			values = append(values, []interface{}{
				graph.NewResourceType(n.Type()).String(),
				color.New(color.FgRed).SprintFunc()("- " + rName),
				"",
				"",
			})
		default:
			rName = nameOrID(nCommon)
			displayProperties = true
		}
		if displayProperties {
			propsChanges, err = addProperties(&values, nCommon.Type(), rName, rNew, nInserted.Properties(), nDeleted.Properties())
			if err != nil {
				return err
			}
		}
		if !propsChanges && rNew {
			values = append(values, []interface{}{
				graph.NewResourceType(n.Type()).String(),
				color.New(color.FgGreen).SprintFunc()("+ " + n.ID().String()),
				"",
				"",
			})
		}
		return nil
	})
	if err != nil {
		return err
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

func addProperties(values *table, rType graph.ResourceType, rName string, rNew bool, insertedProps, deletedProps graph.Properties) (bool, error) {
	changes := false

	for prop, val := range insertedProps {
		var header ColumnDefinition
		for _, h := range DefaultsColumnDefinitions[rType] {
			if h.propKey() == prop {
				header = h
			}
		}
		if header == nil {
			header = &StringColumnDefinition{Prop: prop}
		}
		resourceDisplayF := fmt.Sprint
		if rNew {
			resourceDisplayF = func(i ...interface{}) string { return color.New(color.FgGreen).SprintFunc()("+ " + fmt.Sprint(i...)) }
		}
		(*values) = append((*values), []interface{}{rType.String(),
			resourceDisplayF(rName),
			prop,
			color.New(color.FgGreen).SprintFunc()("+ " + fmt.Sprint(val)),
		})
		changes = true
	}

	for prop, val := range deletedProps {
		var header ColumnDefinition
		for _, h := range DefaultsColumnDefinitions[rType] {
			if h.propKey() == prop {
				header = h
			}
		}
		if header == nil {
			header = &StringColumnDefinition{Prop: prop}
		}
		resourceDisplayF := fmt.Sprint
		if rNew {
			resourceDisplayF = func(i ...interface{}) string { return color.New(color.FgRed).SprintFunc()("- " + fmt.Sprint(i...)) }
		}
		(*values) = append((*values), []interface{}{rType.String(),
			resourceDisplayF(rName),
			prop,
			color.New(color.FgRed).SprintFunc()("- " + fmt.Sprint(val)),
		})
		changes = true
	}

	return changes, nil
}

func (d *diffTableDisplayer) SetDiff(diff *graph.Diff) {
	d.diff = diff
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

func nameOrID(res *graph.Resource) string {
	if name, ok := res.Properties()["Name"]; ok && name != "" {
		return fmt.Sprint(name)
	}
	if id, ok := res.Properties()["Id"]; ok && id != "" {
		return fmt.Sprint(id)
	}
	return res.Id()
}
