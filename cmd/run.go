package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	awscloud "github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/script"
	"github.com/wallix/awless/script/ast"
	"github.com/wallix/awless/script/driver/aws"
)

func init() {
	RootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run awless scripting. Either as one direct command or a given file",

	RunE: func(cmd *cobra.Command, args []string) error {
		var scrpt *script.Script
		var serr error

		if len(args) < 1 {
			return errors.New("missing awless script file path or awless script line")
		}

		if len(args) == 1 {
			content, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			if scrpt, serr = script.Parse(string(content)); serr != nil {
				return serr
			}
		} else {
			text := strings.Join(args, " ")

			if scrpt, serr = script.Parse(text); serr != nil {
				return serr
			}
			expr, ok := scrpt.Statements[0].(*ast.ExpressionNode)
			if !ok {
				return errors.New("Expecting an script expression not a script declaration")
			}
			templ := aws.AWSTemplates[expr.Action+expr.Entity]

			tmplScrpt, err := script.Parse(templ)
			if err != nil {
				return fmt.Errorf("internal error parsing known template\n`%s`\n%s", templ, err)
			}

			prompt := func(question string) interface{} {
				var resp string
				fmt.Printf("%s ? ", question)
				_, err := fmt.Scanln(&resp)
				if err != nil {
					return err
				}

				return resp
			}

			tmplScrpt.InteractiveResolveTemplate(prompt)
			scrpt = tmplScrpt
		}

		defaults, err := database.Current.GetDefaults()
		if err != nil {
			return err
		}
		scrpt.ResolveTemplate(defaults)

		awsDriver := aws.NewDriver(awscloud.InfraService)
		awsDriver.SetLogger(log.New(os.Stdout, "[aws driver] ", log.Ltime))

		if _, err := scrpt.Compile(awsDriver); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(scrpt)
		fmt.Println()
		fmt.Print("About to run compiled script above? (y/n): ")
		var yesorno string
		_, err = fmt.Scanln(&yesorno)

		if strings.TrimSpace(yesorno) == "y" {
			if executedScript, err := scrpt.Run(awsDriver); err != nil {
				return err
			} else {
				fmt.Println()
				fmt.Println(executedScript)
				fmt.Println()
				fmt.Println("Above script ran successfully")
			}
		}

		return nil
	},
}
