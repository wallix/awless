/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud"
	awscloud "github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/cloud/aws/validation"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
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
	PersistentPreRun:   applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
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

		exitOn(runTemplate(templ, getCurrentDefaults()))

		return nil
	},
}

func runTemplate(templ *template.Template, defaults map[string]interface{}) error {
	resolved, err := templ.ResolveTemplate(defaults)
	exitOn(err)
	logger.Infof("used default params: %s (list and set defaults with `awless config`)", sprintProcessedParams(resolved))

	if len(templ.GetHoles()) > 0 {
		fmt.Println("\nMissing required params (Ctrl+C to quit):")
		prompt := func(question string) interface{} {
			var resp string
			for {
				fmt.Printf("%s ? ", question)
				_, err := fmt.Scanln(&resp)
				if err == nil {
					break
				}
				logger.Error("invalid value:", err)
			}
			return resp
		}
		templ.InteractiveResolveTemplate(prompt)
	}

	awsDriver := aws.NewDriver(
		awscloud.InfraService.ProviderRunnableAPI(),
		awscloud.AccessService.ProviderRunnableAPI(),
		awscloud.StorageService.ProviderRunnableAPI(),
	)
	awsDriver.SetLogger(logger.DefaultLogger)

	for _, st := range templ.Statements {
		if validators, ok := validation.ValidatorsPerActions[st.Action()]; ok {
			graph := sync.LoadCurrentLocalGraph(awscloud.ServicePerResourceType[st.Entity()])
			for _, v := range validators {
				if err := v.Validate(graph, st.Params()); err != nil {
					return err
				}
			}
		}
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
		newTempl, err := templ.Run(awsDriver)

		executed := template.NewTemplateExecution(newTempl)

		fmt.Println()
		printReport(executed)

		db, err, close := database.Current()
		exitOn(err)
		defer close()

		db.AddTemplateExecution(executed)

		if err == nil && !executed.HasErrors() {
			if autoSync, ok := defaults[database.SyncAuto]; ok && autoSync.(bool) {
				runSync(newTempl.GetEntitiesSet())
			}
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
		templDef, ok := aws.AWSTemplatesDefinitions[fmt.Sprintf("%s%s", action, entity)]
		if !ok {
			exitOn(errors.New("command unsupported on inline mode"))
		}

		run := func(def aws.TemplateDefinition) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				text := fmt.Sprintf("%s %s %s", def.Action, def.Entity, strings.Join(args, " "))

				node, err := template.ParseStatement(text)
				exitOn(err)

				expr, ok := node.(*ast.ExpressionNode)
				if !ok {
					return errors.New("Expecting a template expression not a template declaration")
				}

				templ, err := template.Parse(templDef.String())
				if err != nil {
					exitOn(fmt.Errorf("internal error parsing template definition\n`%s`\n%s", templDef, err))
				}
				logger.Verbosef("template definition: %s", templDef)

				for k, v := range expr.Params {
					if !strings.Contains(k, ".") {
						expr.Params[fmt.Sprintf("%s.%s", expr.Entity, k)] = v
						delete(expr.Params, k)
					}
				}

				addAliasesToParams(expr)
				resolved, err := templ.ResolveTemplate(expr.Params)
				exitOn(err)
				logger.Infof("used provided params: %s.", sprintProcessedParams(resolved))
				templ.MergeParams(expr.Params)

				exitOn(runTemplate(templ, getCurrentDefaults()))
				return nil
			}
		}

		actionCmd.AddCommand(
			&cobra.Command{
				Use:                templDef.Entity,
				PersistentPreRun:   applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
				PersistentPostRunE: saveHistoryHook,
				Short:              fmt.Sprintf("%s a %s", strings.Title(action), templDef.Entity),
				Long:               fmt.Sprintf("%s a %s\n\tRequired params: %s\n\tExtra params: %s", strings.Title(templDef.Action), templDef.Entity, strings.Join(templDef.Required(), ", "), strings.Join(templDef.Extra(), ", ")),
				RunE:               run(templDef),
			},
		)
	}

	return actionCmd
}

func runSync(entities []string) {
	var services []cloud.Service

	for _, entity := range entities {
		srv, err := cloud.GetServiceForType(entity)
		exitOn(err)
		services = append(services, srv)
	}

	if _, err := sync.DefaultSyncer.Sync(services...); err != nil {
		logger.Errorf("error while synching for %s\n", strings.Join(entities, ", "))
	} else if verboseFlag {
		logger.Infof("performed sync for %s", strings.Join(entities, ", "))
	}
}

func getCurrentDefaults() map[string]interface{} {
	db, err, dbclose := database.Current()
	exitOn(err)
	defaults, err := db.GetDefaults()
	exitOn(err)
	dbclose()
	return defaults
}

func printReport(t *template.TemplateExecution) {
	for _, done := range t.Executed {
		var line bytes.Buffer
		if done.Result != "" {
			line.WriteString(fmt.Sprintf("%s %s ", done.Result, renderGreenFn("<-")))
		}
		line.WriteString(fmt.Sprintf("%s", done.Line))

		if done.Err != "" {
			line.WriteString(fmt.Sprintf("\n\terror: %s", done.Err))
		}

		if done.Err == "" {
			logger.Info(line.String())
		} else {
			logger.Error(line.String())
		}
	}

	if t.IsRevertible() {
		logger.Infof("revert this template with `awless revert %s`", t.ID)
	}
}

func addAliasesToParams(expr *ast.ExpressionNode) error {
	for k, v := range expr.Aliases {
		if !strings.Contains(k, ".") {
			expr.Aliases[fmt.Sprintf("%s.%s", expr.Entity, k)] = v
			delete(expr.Aliases, k)
		}
	}

	graphForResource := sync.LoadCurrentLocalGraph(awscloud.ServicePerResourceType[expr.Entity])

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
		if id, ok := a.ResolveToId(graphForResource, rT); ok {
			expr.Params[k] = id
		} else {
			logger.Infof("Alias '%s' not found in local snapshot. You might want to perform an `awless sync`\n", a)
		}
	}
	return nil
}

func sprintProcessedParams(processed map[string]interface{}) string {
	if len(processed) == 0 {
		return "<none>"
	}
	var str []string
	for k, v := range processed {
		str = append(str, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(str, ", ")
}
