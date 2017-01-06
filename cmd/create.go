package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	awscloud "github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/script"
	"github.com/wallix/awless/script/ast"
	"github.com/wallix/awless/script/driver/aws"
)

func init() {
	createCmd.AddCommand(createInstanceCmd)

	RootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create various type of resources by id: users, groups, instances, vpcs, ...",
}

var createInstanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"inst", "i"},
	Short:   "Create an instance",

	RunE: func(cmd *cobra.Command, args []string) error {
		temp := `create instance subnet={instance_subnet} count={instance_count} base={instance_image} type={instance_type}`

		scr, err := script.Parse(temp)
		if err != nil {
			return err
		}

		defaults := map[string]interface{}{
			"instance_type":  "t2.micro",
			"instance_image": "ami-9398d3e0",
			"instance_count": 1,
		}

		script.VisitExpressionNodes(scr, script.ResolveHolesWith(defaults))

		prompt := func(question string) interface{} {
			var resp string
			fmt.Printf("%s ? ", question)
			_, err := fmt.Scanln(&resp)
			if err != nil {
				return err
			}

			return resp
		}

		script.VisitExpressionNodes(scr, script.InteractiveResolveHoles(prompt))

		var yesorno string
		fmt.Print("\nDone. Params are:\n\n")
		script.VisitExpressionNodes(scr, func(expr *ast.ExpressionNode) {
			fmt.Print(expr.Params)
			fmt.Println()
		})
		fmt.Print("\n\nAbout to run? (y/n): ")
		_, err = fmt.Scanln(&yesorno)
		if err != nil {
			return err
		}

		if strings.TrimSpace(yesorno) == "y" {
			awsDriver := aws.NewDriver(awscloud.InfraService)
			awsDriver.SetLogger(log.New(os.Stdout, "[aws driver] ", log.Ltime))

			return script.Visit(scr, awsDriver)
		}

		return nil
	},
}
