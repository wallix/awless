package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
)

func init() {
	RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "Show, get or set configuration",

	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0: //List all
			d, err := database.Current.GetDefaults()
			exitOn(err)
			for k, v := range d {
				fmt.Printf("%s: %v\t(%[2]T)\n", k, v)
			}
			return nil
		case 1: //List one
			d, ok := database.Current.GetDefault(args[0])
			if !ok {
				fmt.Println("this parameter has not been set")
			} else {
				fmt.Printf("%v\n", d)
			}
			return nil
		case 2: //set
			var val interface{}
			val, err := strconv.Atoi(args[1])
			if err != nil {
				val = args[1]
			}

			err = database.Current.SetDefault(args[0], val)
			exitOn(err)
			return nil
		default:
			return fmt.Errorf("too many parameters")
		}
	},
}
