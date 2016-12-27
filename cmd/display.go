package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/stats"
)

func displayItem(item interface{}, err error, format ...string) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(format) > 0 {
		switch format[0] {
		case "raw":
			fmt.Println(item)
		default:
			lineDisplay(item)
		}
	} else {
		lineDisplay(item)
	}
}

func lineDisplay(item interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	aws.TabularDisplay(item, table)
	table.Render()
}

func displayAliases(aliases stats.Aliases, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if listOnlyIDs {
		names := make([]string, 0, len(aliases))
		for name := range aliases {
			names = append(names, name)
		}
		fmt.Fprintln(os.Stdout, strings.Join(names, " "))
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Alias of"})
		for name, target := range aliases {
			table.Append([]string{name, target})
		}
		table.Render()
	}
}

func generateString(ch rune, nb int) string {
	var res bytes.Buffer
	for i := 0; i < nb; i++ {
		res.WriteRune(ch)
	}
	return res.String()
}
