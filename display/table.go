package display

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const charsIdenticalValues = "//"
const tabWriterWidth = 30

// Table is used to represent an Asci art table
type Table struct {
	columnsHeaders    []string
	columnsDisplayers map[string]*PropertyDisplayer
	columns           map[string][]string
	nbRows            int
}

// NewTable creates a new Table from its header
func NewTable(headers []*PropertyDisplayer) *Table {
	var t Table
	t.columnsDisplayers = make(map[string]*PropertyDisplayer)
	for _, d := range headers {
		t.columnsHeaders = append(t.columnsHeaders, d.displayName())
		t.columnsDisplayers[d.displayName()] = d
	}
	t.columns = make(map[string][]string, len(headers))

	return &t
}

// AddRow adds a row in the table
func (t *Table) AddRow(row ...string) {
	for i, header := range t.columnsHeaders {
		if i < len(row) {
			t.columns[header] = append(t.columns[header], row[i])
		} else {
			t.columns[header] = append(t.columns[header], "")
		}
	}
	t.nbRows++
}

// AddValue adds a value in the header column
func (t *Table) AddValue(header, value string) {
	t.columns[header] = append(t.columns[header], value)
	if len(t.columns[header]) > t.nbRows {
		t.nbRows = len(t.columns[header])
	}
}

// Fprint displays the table in a writer
func (t *Table) Fprint(w io.Writer) {
	t.FprintColumns(w, t.columnsHeaders...)
}

// FprintColumns display some columns of the table in a writer
func (t *Table) FprintColumns(w io.Writer, headers ...string) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(headers)

	var previousFullrow []string
	for i := 0; i < t.nbRows; i++ {
		var row []string
		var fullrow []string
		for j, header := range headers {
			var v string
			var collapseIdenticalValues bool
			if i < len(t.columns[header]) {
				v = t.columnsDisplayers[header].display(t.columns[header][i])
				collapseIdenticalValues = t.columnsDisplayers[header].CollapseIdenticalValues
			}
			fullrow = append(fullrow, v)
			if collapseIdenticalValues && len(previousFullrow) > j && previousFullrow[j] != "" && previousFullrow[j] == v {
				row = append(row, charsIdenticalValues)
			} else {
				row = append(row, v)
			}
		}
		table.Append(row)
		previousFullrow = fullrow
	}

	table.Render()
}

// ColumnValues returns the values of a table column
func (t *Table) ColumnValues(header string) []string {
	return t.columns[header]
}

// FprintColumnValues displays the values of a table column in a writer
func (t *Table) FprintColumnValues(w io.Writer, header, sep string) {
	fmt.Fprintln(w, strings.Join(t.ColumnValues(header), sep))
}
