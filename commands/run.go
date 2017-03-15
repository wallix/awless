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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	awscloud "github.com/wallix/awless/aws"
	"github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver"
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
	Use:               "run FILEPATH",
	Short:             "Run a template given a filepath",
	Example:           "  awless run ~/templates/my-infra.txt",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
	PersistentPostRun: applyHooks(saveHistoryHook, verifyNewVersionHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing FILEPATH arg")
		}

		content, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		templ, err := template.Parse(string(content))
		exitOn(err)

		extraParams, err := template.ParseParams(strings.Join(args[1:], " "))
		exitOn(err)

		env := template.NewEnv()
		env.Log = logger.DefaultLogger
		env.AddFillers(config.Defaults, extraParams)
		env.MissingHolesFunc = missingHolesStdinFunc()
		env.DefLookupFunc = lookupTemplateDefinitionsFunc()

		exitOn(runTemplate(templ, env))

		return nil
	},
}

func missingHolesStdinFunc() func(string) interface{} {
	var count int
	return func(hole string) interface{} {
		if count < 1 {
			fmt.Println("Please specify (Ctrl+C to quit):")
		}
		var resp interface{}
		ask := func() error {
			fmt.Printf("%s ? ", hole)
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return err
			}
			line = strings.TrimSpace(line)
			if line == "" {
				return errors.New("empty")
			}
			params, err := template.ParseParams(fmt.Sprintf("%s=%s", hole, line))
			if err != nil {
				return err
			}
			resp = params[hole]
			return nil
		}
		for err := ask(); err != nil; err = ask() {
			logger.Errorf("invalid value: %s", err)
		}
		count++
		return resp
	}
}

func runTemplate(templ *template.Template, env *template.Env) error {
	if len(env.Fillers) > 0 {
		logger.Verbosef("default/given holes fillers: %s", sprintProcessedParams(env.Fillers))
	}

	var err error
	templ, env, err = template.Compile(templ, env)
	exitOn(err)

	validateTemplate(templ)

	var drivers []driver.Driver
	for _, s := range cloud.ServiceRegistry {
		drivers = append(drivers, s.Drivers()...)
	}
	awsDriver := driver.NewMultiDriver(drivers...)

	awsDriver.SetLogger(logger.DefaultLogger)

	_, err = templ.Compile(awsDriver)
	exitOn(err)

	fmt.Println()
	fmt.Printf("%s\n", renderGreenFn(templ))
	fmt.Println()
	fmt.Print("Confirm? (y/n): ")
	var yesorno string
	_, err = fmt.Scanln(&yesorno)

	if strings.TrimSpace(yesorno) == "y" {
		newTempl, err := templ.Run(awsDriver)

		fmt.Println()
		fmt.Println("Executed:")
		printer := template.NewPrinter(os.Stdout)
		printer.RenderKO = renderRedFn
		printer.RenderOK = renderGreenFn
		printer.Print(newTempl)

		db, err, close := database.Current()
		exitOn(err)
		defer close()

		db.AddTemplate(newTempl)
		if template.IsRevertible(newTempl) {
			fmt.Println()
			logger.Infof("Revert this template with `awless revert %s`", newTempl.ID)
		}

		if err == nil && !newTempl.HasErrors() {
			runSyncFor(newTempl)
		}
	}

	return nil
}

func validateTemplate(tpl *template.Template) {
	unicityRule := &template.UniqueNameValidator{LookupGraph: func(key string) (*graph.Graph, bool) {
		g := sync.LoadCurrentLocalGraph(awscloud.ServicePerResourceType[key])
		return g, true
	}}

	errs := tpl.Validate(unicityRule)

	if len(errs) > 0 {
		for _, err := range errs {
			logger.Warning(err)
		}
	}
}

func createDriverCommands(action string, entities []string) *cobra.Command {
	actionCmd := &cobra.Command{
		Use:         action,
		Short:       oneLinerShortDesc(action, entities),
		Long:        fmt.Sprintf("Allow to %s: %v", action, strings.Join(entities, ", ")),
		Annotations: map[string]string{"one-liner": "true"},
	}

	for _, entity := range entities {
		templDef, ok := lookupTemplateDefinitionsFunc()(fmt.Sprintf("%s%s", action, entity))
		if !ok {
			exitOn(errors.New("command unsupported on inline mode"))
		}
		run := func(def template.TemplateDefinition) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				text := fmt.Sprintf("%s %s %s", def.Action, def.Entity, strings.Join(args, " "))

				templ, err := template.Parse(text)
				exitOn(err)

				env := template.NewEnv()
				env.Log = logger.DefaultLogger
				env.AddFillers(config.Defaults)
				env.DefLookupFunc = lookupTemplateDefinitionsFunc()
				env.AliasFunc = resolveAliasFunc(def.Entity)
				env.MissingHolesFunc = missingHolesStdinFunc()

				exitOn(runTemplate(templ, env))
				return nil
			}
		}

		actionCmd.AddCommand(
			&cobra.Command{
				Use:               templDef.Entity,
				PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
				PersistentPostRun: applyHooks(saveHistoryHook, verifyNewVersionHook),
				Short:             fmt.Sprintf("%s a %s", strings.Title(action), templDef.Entity),
				Long:              fmt.Sprintf("%s a %s\n\tRequired params: %s\n\tExtra params: %s", strings.Title(templDef.Action), templDef.Entity, strings.Join(templDef.Required(), ", "), strings.Join(templDef.Extra(), ", ")),
				RunE:              run(templDef),
			},
		)
	}

	return actionCmd
}

func lookupTemplateDefinitionsFunc() template.LookupTemplateDefFunc {
	return func(key string) (t template.TemplateDefinition, ok bool) {
		t, ok = aws.AWSTemplatesDefinitions[key]
		return
	}
}

func runSyncFor(tpl *template.Template) {
	if !config.GetAutosync() {
		return
	}

	collector := &template.CollectDefinitions{L: lookupTemplateDefinitionsFunc()}
	tpl.Visit(collector)

	uniqueNames := make(map[string]bool)
	for _, def := range collector.C {
		name, ok := awscloud.ServicePerAPI[def.Api]
		if ok {
			uniqueNames[name] = true
		}
	}

	var srvNames []string
	for name := range uniqueNames {
		srvNames = append(srvNames, name)
	}

	var services []cloud.Service
	for _, name := range srvNames {
		srv, ok := cloud.ServiceRegistry[name]
		if !ok {
			logger.Errorf("internal: cannot resolve service name '%s'", name)
		} else {
			services = append(services, srv)
		}
	}

	if _, err := sync.DefaultSyncer.Sync(services...); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Verbosef("performed sync for %s", strings.Join(srvNames, ", "))
	}
}

func printReport(t *template.Template) {
	for _, done := range t.CommandNodesIterator() {
		var line bytes.Buffer
		if done.CmdResult != "" {
			line.WriteString(fmt.Sprintf("%s %s ", done.CmdResult, renderGreenFn("<-")))
		}
		line.WriteString(fmt.Sprintf("%s", done.String()))

		if done.CmdErr != nil {
			line.WriteString(fmt.Sprintf("\n\terror: %s", done.CmdErr))
		}

		if done.CmdErr == nil {
			logger.Info(line.String())
		} else {
			logger.Error(line.String())
		}
	}

	if template.IsRevertible(t) {
		logger.Infof("revert this template with `awless revert %s`", t.ID)
	}
}

func resolveAliasFunc(entity string) func(k, v string) string {
	gph := sync.LoadCurrentLocalGraph(awscloud.ServicePerResourceType[entity])

	return func(key, alias string) string {
		resType := key
		if strings.Contains(key, "id") {
			resType = entity
		}
		a := graph.Alias(alias)
		if id, ok := a.ResolveToId(gph, resType); ok {
			return id
		}
		return ""
	}

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

func oneLinerShortDesc(action string, entities []string) string {
	if len(entities) > 5 {
		return fmt.Sprintf("%s, \u2026 (see `awless %s -h` for more)", strings.Join(entities[0:5], ", "), action)
	} else {
		return strings.Join(entities, ", ")
	}

}
