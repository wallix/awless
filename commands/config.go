package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
)

var keysOnly bool

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configListCmd.Flags().BoolVar(&keysOnly, "keys", false, "list only config keys")
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
}

var configCmd = &cobra.Command{
	Use:                "config",
	Short:              "get, set, unset or list configuration values",
	PersistentPreRunE:  initAwlessEnvHook,
	PersistentPostRunE: saveHistoryHook,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",

	Run: func(cmd *cobra.Command, args []string) {
		db, err, close := database.Current()
		exitOn(err)
		defer close()
		d, err := db.GetDefaults()
		exitOn(err)
		for k, v := range d {
			if keysOnly {
				fmt.Println(k)
			} else {
				fmt.Printf("%s: %v\t(%[2]T)\n", k, v)
			}
		}
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get {key}",
	Short: "Get a configuration value",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("not enough parameters")
		}
		db, err, close := database.Current()
		exitOn(err)
		defer close()
		d, ok := db.GetDefault(args[0])
		if !ok {
			fmt.Println("this parameter has not been set")
		} else {
			fmt.Printf("%v\n", d)
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set {key} {value}",
	Short: "Set or update a configuration value",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("not enough parameters")
		}
		key := strings.TrimSpace(args[0])
		var value string
		if len(args) == 1 {
			switch key {
			case "region":
				value = askRegion()
			default:
				fmt.Print("Value ? > ")
				fmt.Scan(&value)
			}
		} else {
			value = args[1]
			switch key {
			case "region":
				if !aws.IsValidRegion(value) {
					fmt.Println("Invalid region!")
					value = askRegion()
				}
			}
		}
		if value == "" {
			return fmt.Errorf("invalid empty value")
		}

		var i interface{}

		if num, nerr := strconv.Atoi(value); nerr == nil {
			i = num
		} else if b, berr := strconv.ParseBool(value); berr == nil {
			i = b
		} else {
			i = value
		}

		db, err, close := database.Current()
		exitOn(err)
		defer close()
		exitOn(db.SetDefault(key, i))

		return nil
	},
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset {key}",
	Short: "Unset a configuration value",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("not enough parameters")
		}
		db, err, close := database.Current()
		exitOn(err)
		defer close()
		_, ok := db.GetDefault(args[0])
		if !ok {
			fmt.Println("this parameter has not been set")
		} else {
			db.UnsetDefault(args[0])
		}
		return nil
	},
}

func askRegion() string {
	var region string
	fmt.Println("Please choose one region:")

	fmt.Println(strings.Join(aws.AllRegions(), ", "))
	fmt.Println()
	fmt.Print("Value ? > ")
	fmt.Scan(&region)
	for !aws.IsValidRegion(region) {
		fmt.Println("Invalid!")
		fmt.Print("Value ? > ")
		fmt.Scan(&region)
	}
	return region
}
