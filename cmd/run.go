package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
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

			templName := fmt.Sprintf("%s%s", expr.Action, expr.Entity)
			templ, ok := aws.AWSTemplates[templName]
			if !ok {
				exitOn(errors.New("command unsupported on inline mode"))
			}

			scrpt, serr = script.Parse(templ)
			if serr != nil {
				exitOn(fmt.Errorf("internal error parsing known template\n`%s`\n%s", templ, serr))
			}

			scrpt.ResolveTemplate(expr.Params)
		}

		defaults, err := database.Current.GetDefaults()
		exitOn(err)

		scrpt.ResolveTemplate(defaults)

		prompt := func(question string) interface{} {
			var resp string
			fmt.Printf("%s ? ", question)
			_, err := fmt.Scanln(&resp)
			exitOn(err)

			return resp
		}
		scrpt.InteractiveResolveTemplate(prompt)

		awsDriver := aws.NewDriver(awscloud.InfraService)
		if verboseFlag {
			awsDriver.SetLogger(log.New(os.Stdout, "[aws driver] ", log.Ltime))
		}

		_, err = scrpt.Compile(awsDriver)
		exitOn(err)

		fmt.Println()
		fmt.Println(scrpt)
		fmt.Println()
		fmt.Print("Run compiled script above? (y/n): ")
		var yesorno string
		_, err = fmt.Scanln(&yesorno)

		if strings.TrimSpace(yesorno) == "y" {
			executedScript, err := scrpt.Run(awsDriver)
			exitOn(err)

			fmt.Println()
			green := color.New(color.FgGreen).SprintFunc()
			for _, stat := range executedScript.Statements {
				fmt.Printf("%s -> %s\n", stat, green("DONE"))
			}
		}

		return nil
	},
}
