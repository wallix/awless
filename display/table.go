package display

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const charsIdenticalValues = "//"
const ascSymbol = " ▲"
const tabWriterWidth = 30

// Table is used to represent an Asci art table
type Table struct {
	columnsHeaders      []string
	columnsDisplayers   map[string]*PropertyDisplayer
	columns             map[string][]string
	nbRows              int
	sortByColumns       []string
	MergeIdenticalCells bool
}

// NewTable creates a new Table from its header
func NewTable(headers []*PropertyDisplayer) *Table {
	var t Table
	t.columnsDisplayers = make(map[string]*PropertyDisplayer)
	for _, d := range headers {
		t.columnsHeaders = append(t.columnsHeaders, cleanColumnName(d.displayName()))
		t.columnsDisplayers[cleanColumnName(d.displayName())] = d
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
	t.columns[cleanColumnName(header)] = append(t.columns[cleanColumnName(header)], value)
	if len(t.columns[cleanColumnName(header)]) > t.nbRows {
		t.nbRows = len(t.columns[cleanColumnName(header)])
	}
}

//SetSortBy sets the columns that will be used to sort the rows of the table
func (t *Table) SetSortBy(columns ...string) {
	t.sortByColumns = []string{}
	for _, c := range columns {
		c = cleanColumnName(c)
		if _, ok := t.columnsDisplayers[c]; ok {
			t.sortByColumns = append(t.sortByColumns, c)
		}
	}
}

// Fprint displays the table in a writer
func (t *Table) Fprint(w io.Writer) {
	sort.Sort(byColumns(*t))
	t.FprintColumns(w, t.columnsHeaders...)
}

// FprintColumns display some columns of the table in a writer
func (t *Table) FprintColumns(w io.Writer, headers ...string) {
	table := tablewriter.NewWriter(w)
	if t.MergeIdenticalCells {
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
	}
	var displayHeaders []string
	for _, h := range headers {
		if len(t.sortByColumns) >= 1 && cleanColumnName(t.sortByColumns[0]) == h {
			displayHeaders = append(displayHeaders, h+ascSymbol)
		} else {
			displayHeaders = append(displayHeaders, h)
		}
	}
	table.SetHeader(displayHeaders)

	for i := 0; i < t.nbRows; i++ {
		var row []string
		for _, header := range headers {
			header = cleanColumnName(header)
			var v string
			if i < len(t.columns[header]) {
				v = t.columnsDisplayers[header].display(t.columns[header][i])
			}
			row = append(row, v)
		}
		table.Append(row)
	}

	table.Render()
}

// ColumnValues returns the values of a table column
func (t *Table) ColumnValues(header string) []string {
	sort.Sort(byColumns(*t))
	res := t.columns[cleanColumnName(header)]
	removeDuplicatesOfSlice(&res)
	return res
}

// FprintColumnValues displays the values of a table column in a writer
func (t *Table) FprintColumnValues(w io.Writer, header, sep string) {
	sort.Sort(byColumns(*t))
	fmt.Fprintln(w, strings.Join(t.ColumnValues(header), sep))
}

func removeDuplicatesOfSlice(data *[]string) {
	length := len(*data) - 1
	for i := 0; i < length; i++ {
		for j := i + 1; j <= length; j++ {
			if (*data)[i] == (*data)[j] {
				(*data)[j] = (*data)[length]
				(*data) = (*data)[0:length]
				length--
				j--
			}
		}
	}
}

func cleanColumnName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

type byColumns Table

func (c byColumns) Len() int { return c.nbRows }
func (c byColumns) Swap(i, j int) {
	for _, col := range c.columns {
		col[i], col[j] = col[j], col[i]
	}
}
func (c byColumns) Less(i, j int) bool {
	for _, col := range c.sortByColumns {
		if c.columns[col][i] == c.columns[col][j] {
			continue
		}
		return c.columns[col][i] < c.columns[col][j]
	}
	return false
}
