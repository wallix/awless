package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/wallix/awless/cloud/aws"
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
