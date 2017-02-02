package cmd

import (
	"bytes"
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

var renderGreenFn = color.New(color.FgGreen).SprintFunc()
var renderRedFn = color.New(color.FgRed).SprintFunc()

func init() {
	RootCmd.AddCommand(runCmd)
	for action, entities := range aws.DriverSupportedActions() {
		RootCmd.AddCommand(
			createDriverCommands(action, entities),
		)
	}
}

var runCmd = &cobra.Command{
	Use:                "run",
	Short:              "Run an awless template file given as the only argument. Ex: awless run mycloud.awless",
	PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,

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

		exitOn(runTemplate(templ))

		return nil
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

	awsDriver := aws.NewDriver(
		awscloud.InfraService.ProviderRunnableAPI(),
		awscloud.AccessService.ProviderRunnableAPI(),
	)

	if verboseFlag {
		awsDriver.SetLogger(log.New(os.Stdout, "[aws driver] ", log.Ltime))
	}

	_, err = templ.Compile(awsDriver)
	exitOn(err)

	fmt.Println()
	fmt.Printf("%s\n", renderGreenFn(templ))
	fmt.Println()
	fmt.Print("Run verified operations above? (y/n): ")
	var yesorno string
	_, err = fmt.Scanln(&yesorno)

	if strings.TrimSpace(yesorno) == "y" {
		newTempl, _ := templ.Run(awsDriver)

		executed := template.NewTemplateExecution(newTempl)

		fmt.Println()
		printReport(executed)

		db, close := database.Current()
		defer close()

		db.AddTemplateExecution(executed)
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

				node, err := template.ParseStatement(text)
				exitOn(err)

				expr, ok := node.(*ast.ExpressionNode)
				if !ok {
					return errors.New("Expecting an template expression not a template declaration")
				}

				templName := fmt.Sprintf("%s%s", expr.Action, expr.Entity)
				templDef, ok := aws.AWSTemplatesDefinitions[templName]
				if !ok {
					exitOn(errors.New("command unsupported on inline mode"))
				}

				templ, err := template.Parse(templDef)
				if err != nil {
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
				Use:                entity,
				PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, checkStatsHook),
				PersistentPostRunE: saveHistoryHook,
				Short:              fmt.Sprintf("Use it to %s a %s", action, entity),
				RunE:               run(action, entity),
			},
		)
	}

	return actionCmd
}

func printReport(t *template.TemplateExecution) {
	for _, done := range t.Executed {
		var line bytes.Buffer
		if done.Err == "" {
			line.WriteString(renderGreenFn("[DONE] "))
		} else {
			line.WriteString(renderRedFn("[ERROR] "))
		}

		if done.Result != "" {
			line.WriteString(fmt.Sprintf("%s %s ", done.Result, renderGreenFn("<-")))
		}
		line.WriteString(fmt.Sprintf("%s", done.Line))

		if done.Err != "" {
			line.WriteString(fmt.Sprintf("\n\terror: %s", done.Err))
		}

		fmt.Println(line.String())
	}

	fmt.Println()
	fmt.Printf("Revert this template with `awless revert -i %s`\n", t.ID)
	fmt.Println()
	fmt.Println("(List executed template executions using `awless revert -l`)")
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
		} else {
			fmt.Printf("Alias '%s' not found in local snapshot. You might want to perform an `awless sync`\n", a)
		}
	}
	return nil
}
