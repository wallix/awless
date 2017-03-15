package template

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/oklog/ulid"
)

type renderFunc func(...interface{}) string

func renderNoop(s ...interface{}) string { return fmt.Sprint(s) }

type Printer struct {
	w io.Writer

	IncludeMeta bool
	RenderOK    renderFunc
	RenderKO    renderFunc
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{
		w:        w,
		RenderOK: renderNoop,
		RenderKO: renderNoop,
	}
}

func (p *Printer) Print(t *Template) {
	buff := bufio.NewWriter(p.w)

	if p.IncludeMeta {
		buff.WriteString(fmt.Sprintf("Date: %s", parseULIDDate(t.ID)))

		if IsRevertible(t) {
			buff.WriteString(fmt.Sprintf(", RevertID: %s", t.ID))
		} else {
			buff.WriteString(", RevertID: <not revertible>")
		}
		buff.WriteString("\n")
	}

	tabw := tabwriter.NewWriter(buff, 0, 8, 0, '\t', 0)
	for _, cmd := range t.CommandNodesIterator() {
		var result, status string

		exec := fmt.Sprintf("%s", cmd.String())

		if cmd.CmdErr != nil {
			status = "KO"
		} else {
			status = "OK"
		}

		if v, ok := cmd.CmdResult.(string); ok && v != "" {
			result = fmt.Sprintf("[%s]", v)
		}

		var line string
		if cmd.CmdErr != nil {
			line = fmt.Sprintf("    %s\t%s\t%s\t", p.RenderKO(status), exec, result)
		} else {
			line = fmt.Sprintf("    %s\t%s\t%s\t", p.RenderOK(status), exec, result)
		}

		fmt.Fprintln(tabw, line)
		if cmd.CmdErr != nil {
			for _, err := range formatMultiLineErrMsg(cmd.CmdErr.Error()) {
				fmt.Fprintf(tabw, "%s\t%s\n", "", err)
			}
		}

	}

	tabw.Flush()
	buff.Flush()
}

func formatMultiLineErrMsg(msg string) []string {
	notabs := strings.Replace(msg, "\t", "", -1)
	var indented []string
	for _, line := range strings.Split(notabs, "\n") {
		indented = append(indented, fmt.Sprintf("    %s", line))
	}
	return indented
}

func parseULIDDate(uid string) string {
	parsed, err := ulid.Parse(uid)
	if err != nil {
		panic(err)
	}

	date := time.Unix(int64(parsed.Time())/int64(1000), time.Nanosecond.Nanoseconds())

	return date.Format(time.Stamp)
}
