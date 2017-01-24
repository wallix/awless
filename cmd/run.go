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
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/ast"
	"github.com/wallix/awless/template/driver/aws"
)

func init() {
	RootCmd.AddCommand(runCmd)
	for action, entities := range aws.DriverSupportedActions() {
		RootCmd.AddCommand(
			createDriverCommands(action, entities),
		)
	}
}

var runCmd = &cobra.Command{
	Use:               "run",
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,
	Short:             "Run an awless template file given as the only argument. Ex: awless run mycloud.awless",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing awless template file path")
		}

		content, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		templ, err := template.Parse(string(content))
		exitOn(err)

		return runTemplate(templ)
	},
}

func runTemplate(templ *template.Template) error {
	db, dbclose := database.Current()
	defaults, err := db.GetDefaults()
	exitOn(err)
	dbclose()

	templ.ResolveTemplate(defaults)

	prompt := func(question string) interface{} {
		var resp string
		fmt.Printf("%s ? ", question)
		_, err = fmt.Scanln(&resp)
		exitOn(err)

		return resp
	}
	templ.InteractiveResolveTemplate(prompt)

	awsDriver := aws.NewDriver(awscloud.InfraService)
	if verboseFlag {
		awsDriver.SetLogger(log.New(os.Stdout, "[aws driver] ", log.Ltime))
	}

	_, err = templ.Compile(awsDriver)
	exitOn(err)

	green := color.New(color.FgGreen).SprintFunc()

	fmt.Println()
	fmt.Printf("%s\n", green(templ))
	fmt.Println()
	fmt.Print("Run compiled template above? (y/n): ")
	var yesorno string
	_, err = fmt.Scanln(&yesorno)

	if strings.TrimSpace(yesorno) == "y" {
		executedTemplate, err := templ.Run(awsDriver)
		exitOn(err)

		fmt.Println()
		for _, stat := range executedTemplate.Statements {
			fmt.Printf("%s -> %s\n", stat, green("DONE"))
		}
	}

	return nil
}

func createDriverCommands(action string, entities []string) *cobra.Command {
	actionCmd := &cobra.Command{
		Use:   action,
		Short: fmt.Sprintf("Allow to %s: %v", action, strings.Join(entities, ", ")),
	}

	for _, entity := range entities {
		run := func(act, ent string) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				text := fmt.Sprintf("%s %s %s", act, ent, strings.Join(args, " "))

				templ, err := template.Parse(text)
				exitOn(err)

				expr, ok := templ.Statements[0].(*ast.ExpressionNode)
				if !ok {
					return errors.New("Expecting an template expression not a template declaration")
				}

				templName := fmt.Sprintf("%s%s", expr.Action, expr.Entity)
				templDef, ok := aws.AWSDriverTemplates[templName]
				if !ok {
					exitOn(errors.New("command unsupported on inline mode"))
				}

				if templ, err = template.Parse(templDef); err != nil {
					exitOn(fmt.Errorf("internal error parsing template definition\n`%s`\n%s", templDef, err))
				}

				for k, v := range expr.Params {
					if !strings.Contains(k, ".") {
						expr.Params[fmt.Sprintf("%s.%s", expr.Entity, k)] = v
						delete(expr.Params, k)
					}
				}

				addAliasesToParams(expr)

				templ.ResolveTemplate(expr.Params)

				templ.MergeParams(expr.Params)

				return runTemplate(templ)
			}
		}

		actionCmd.AddCommand(
			&cobra.Command{
				Use:               entity,
				PersistentPreRun:  initCloudServicesFn,
				PersistentPostRun: saveHistoryFn,
				Short:             fmt.Sprintf("Use it to %s a %s", action, entity),
				RunE:              run(action, entity),
			},
		)
	}

	return actionCmd
}

func addAliasesToParams(expr *ast.ExpressionNode) error {
	for k, v := range expr.Aliases {
		if !strings.Contains(k, ".") {
			expr.Aliases[fmt.Sprintf("%s.%s", expr.Entity, k)] = v
			delete(expr.Aliases, k)
		}
	}

	infra, err := config.LoadInfraGraph()
	exitOn(err)
	access, err := config.LoadAccessGraph()
	exitOn(err)

	for k, v := range expr.Aliases {
		if !strings.Contains(k, ".") {
			return fmt.Errorf("invalid alias key (no '.') %s", k)
		}
		var t string
		if strings.Split(k, ".")[1] == "id" {
			t = strings.Split(k, ".")[0]
		} else {
			t = strings.Split(k, ".")[1]
		}
		rT := graph.ResourceType(t)
		a := graph.Alias(v)
		if id, ok := a.ResolveToId(infra, rT); ok {
			expr.Params[k] = id
		} else if id, ok := a.ResolveToId(access, rT); ok {
			expr.Params[k] = id
		}
	}
	return nil
}
