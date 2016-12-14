package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/wallix/awless/cloud/aws"
)

func display(item interface{}, err error, format ...string) {
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
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 1, ' ', 0)
	aws.TabularDisplay(item, w)
	w.Flush()
}
