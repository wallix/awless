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
uimitations under the License.
*/

package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	stdsync "sync"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	awscloud "github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver"
)

func init() {
	RootCmd.AddCommand(runCmd)
	for action, entities := range awscloud.DriverSupportedActions() {
		RootCmd.AddCommand(
			createDriverCommands(action, entities),
		)
	}
}

var runCmd = &cobra.Command{
	Use:               "run PATH",
	Short:             "Run a template given a filepath or a URL (prefixed with http)",
	Example:           "  awless run ~/templates/my-infra.txt\n  awless run https://raw.githubusercontent.com/wallix/awless-templates/master/create_vpc.awls\n  awless run repo:create_vpc",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
	PersistentPostRun: applyHooks(saveHistoryHook, verifyNewVersionHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing PATH arg (filepath or url)")
		}

		content, err := getTemplateText(args[0])
		exitOn(err)

		logger.Verbosef("Loaded template text:\n\n%s\n", content)

		templ, err := template.Parse(string(content))
		exitOn(err)

		extraParams, err := template.ParseParams(strings.Join(args[1:], " "))
		exitOn(err)

		exitOn(runTemplate(templ, config.Defaults, extraParams))

		return nil
	},
}

func missingHolesStdinFunc() func(string) interface{} {
	var count int
	return func(hole string) (response interface{}) {
		if count < 1 {
			fmt.Println("Please specify (Ctrl+C to quit, Tab for completion):")
		}

		var err error
		for response, err = askHole(hole); err != nil; response, err = askHole(hole) {
			logger.Errorf("invalid value: %s", err)
		}
		count++
		return
	}
}

func askHole(hole string) (interface{}, error) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("%s? ", hole),
		AutoComplete:    idAndNameCompleter(hole),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		exitOn(err)
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				os.Exit(0)
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case line == "":
			return nil, errors.New("empty")
		case !isQuoted(line) && !template.MatchStringParamValue(line):
			return nil, errors.New("string contains spaces or special characters: surround it with quotes")
		default:
			params, err := template.ParseParams(fmt.Sprintf("%s=%s", hole, line))
			if err != nil {
				return nil, err
			}
			return params[hole], nil
		}
	}
	return nil, nil
}

type onceLoader struct {
	g    *graph.Graph
	err  error
	once stdsync.Once
}

func (l *onceLoader) load() (*graph.Graph, error) {
	l.once.Do(func() {
		l.g, l.err = sync.LoadAllGraphs()
	})
	return l.g, l.err
}

var allGraphsOnce = &onceLoader{}

func idAndNameCompleter(hole string) readline.AutoCompleter {
	g, err := allGraphsOnce.load()
	if err != nil {
		exitOn(err)
	}

	types := strings.Split(hole, ".")
	resources, err := g.GetAllResources(types...)
	if err != nil {
		exitOn(err)
	}
	listAllResourcesIdAndName := func(s string) (suggest []string) {
		for _, res := range resources {
			id := res.Id()
			if !template.MatchStringParamValue(id) {
				id = "'" + id + "'"
			}
			if strings.Contains(id, s) {
				suggest = append(suggest, id)
			}
			if val, ok := res.Properties["Name"]; ok {
				switch val.(type) {
				case string:
					name := val.(string)
					if !template.MatchStringParamValue(name) {
						name = "'" + name + "'"
					}
					prefixed := fmt.Sprintf("@%s", name)
					if strings.Contains(prefixed, s) && name != "" {
						suggest = append(suggest, prefixed)
					}
				}
			}
		}

		sort.Strings(suggest)

		return
	}
	return readline.NewPrefixCompleter(readline.PcItemDynamic(listAllResourcesIdAndName))
}

func runTemplate(templ *template.Template, fillers ...map[string]interface{}) error {
	env := template.NewEnv()
	env.Log = logger.DefaultLogger
	env.AddFillers(fillers...)
	env.DefLookupFunc = lookupDefinitionsFunc
	env.AliasFunc = resolveAliasFunc
	env.MissingHolesFunc = missingHolesStdinFunc()

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

	errs := templ.DryRun(awsDriver)
	if len(errs) > 0 {
		for _, dryRunErr := range errs {
			logger.Errorf(dryRunErr.Error())
		}
		exitOn(errors.New("Dryrun failed"))
	}

	fmt.Printf("%s\n", renderGreenFn(templ))

	var yesorno string
	if forceGlobalFlag {
		yesorno = "y"
	} else {
		fmt.Println()
		fmt.Print("Confirm? (y/n): ")
		_, err = fmt.Scanln(&yesorno)
		exitOn(err)
	}

	if strings.TrimSpace(yesorno) == "y" {
		newTempl, err := templ.Run(awsDriver)
		if err != nil {
			logger.Errorf("Running template error: %s", err)
		}

		printer := template.NewDefaultPrinter(os.Stdout)
		printer.RenderKO = renderRedFn
		printer.RenderOK = renderGreenFn
		printer.Print(newTempl)

		db, err, close := database.Current()
		exitOn(err)
		defer close()

		err = db.AddTemplate(newTempl)
		if err != nil {
			logger.Errorf("Cannot save executed template in awless logs: %s", err)
		}
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
		g := sync.LoadCurrentLocalGraph(aws.ServicePerResourceType[key])
		return g, true
	}}

	errs := tpl.Validate(unicityRule, &template.ParamIsSetValidator{Action: "create", Entity: "instance", Param: "key", WarningMessage: "This instance has no access key. You might not be able to connect to it. Use `awless create instance key=my-key ...`"})

	if len(errs) > 0 {
		for _, err := range errs {
			logger.Warning(err)
		}
		fmt.Fprintln(os.Stderr)
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
		templDef, ok := lookupDefinitionsFunc(fmt.Sprintf("%s%s", action, entity))
		if !ok {
			exitOn(errors.New("command unsupported on inline mode"))
		}
		run := func(def template.Definition) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				text := fmt.Sprintf("%s %s %s", def.Action, def.Entity, strings.Join(args, " "))

				templ, err := template.Parse(text)
				exitOn(err)

				exitOn(runTemplate(templ, config.Defaults))
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

func lookupDefinitionsFunc(key string) (t template.Definition, ok bool) {
	t, ok = awscloud.AWSTemplatesDefinitions[key]
	return
}

func runSyncFor(tpl *template.Template) {
	if !config.GetAutosync() {
		return
	}

	defs := tpl.UniqueDefinitions(lookupDefinitionsFunc)

	services := aws.GetCloudServicesForAPIs(defs.Map(
		func(d template.Definition) string { return d.Api },
	)...)

	if _, err := sync.DefaultSyncer.Sync(services...); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Verbosef("performed sync for %s", strings.Join(cloud.Services(services).Names(), ", "))
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

func resolveAliasFunc(entity, key, alias string) string {
	gph := sync.LoadCurrentLocalGraph(aws.ServicePerResourceType[entity])
	resType := key
	if strings.Contains(key, "id") {
		resType = entity
	}

	resources, err := gph.ResolveResources(&graph.And{Resolvers: []graph.Resolver{&graph.ByProperty{Key: "Name", Value: alias}, &graph.ByType{Typ: resType}}})
	if err != nil {
		return ""
	}
	switch len(resources) {
	case 1:
		return resources[0].Id()
	default:
		resources, err := gph.ResolveResources(&graph.And{Resolvers: []graph.Resolver{&graph.ByProperty{Key: "Name", Value: alias}}})
		if err != nil {
			return ""
		}
		if len(resources) > 0 {
			return resources[0].Id()
		}
	}

	return ""
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

const (
	DEFAULT_REPO_PREFIX = "https://raw.githubusercontent.com/wallix/awless-templates/master"
	FILE_EXT            = ".awls"
)

func getTemplateText(path string) ([]byte, error) {
	if strings.HasPrefix(path, "repo:") {
		path = fmt.Sprintf("%s/%s", DEFAULT_REPO_PREFIX, strings.TrimPrefix(path[5:], "/"))
		path = fmt.Sprintf("%s%s", strings.TrimSuffix(path, FILE_EXT), FILE_EXT)
	}

	if strings.HasPrefix(path, "http") {
		logger.ExtraVerbosef("fetching remote template at '%s'", path)
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("'%s' when fetching template at '%s'", resp.Status, path)
		}

		return ioutil.ReadAll(resp.Body)
	}

	return ioutil.ReadFile(path)
}

func isQuoted(s string) bool {
	if strings.HasPrefix(s, "@") {
		return isQuoted(s[1:])
	}
	return (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) || strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")
}
