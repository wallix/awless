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
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	stdsync "sync"
	"text/tabwriter"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/wallix/awless-scheduler/client"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/aws/doc"
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

var scheduleFlag bool
var scheduleRunInFlag string
var scheduleRevertInFlag string
var listRemoteTemplatesFlag bool

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVar(&listRemoteTemplatesFlag, "list", false, "List templates available at https://github.com/wallix/awless-templates")
	runCmd.Flags().BoolVar(&scheduleFlag, "schedule", false, "Schedule the execution of this template")
	runCmd.Flags().StringVar(&scheduleRunInFlag, "run-in", "", "Postpone the execution of this template")
	runCmd.Flags().StringVar(&scheduleRevertInFlag, "revert-in", "", "Schedule the revertion of this template")
	runCmd.Flags().MarkHidden("schedule")
	runCmd.Flags().MarkHidden("run-in")
	runCmd.Flags().MarkHidden("revert-in")
	for action, entities := range awsdriver.DriverSupportedActions() {
		cmd := createDriverCommands(action, entities)
		cmd.PersistentFlags().BoolVar(&scheduleFlag, "schedule", false, "Schedule the execution of this command")
		cmd.PersistentFlags().StringVar(&scheduleRunInFlag, "run-in", "", "Postpone the execution of this command")
		cmd.PersistentFlags().StringVar(&scheduleRevertInFlag, "revert-in", "", "Schedule the revertion of this command")
		cmd.PersistentFlags().MarkHidden("schedule")
		cmd.PersistentFlags().MarkHidden("run-in")
		cmd.PersistentFlags().MarkHidden("revert-in")
		RootCmd.AddCommand(cmd)
	}
}

var runCmd = &cobra.Command{
	Use:               "run PATH",
	Short:             "Run a template given a filepath or a URL (prefixed with http)",
	Example:           "  awless run ~/templates/my-infra.txt\n  awless run https://raw.githubusercontent.com/wallix/awless-templates/master/create_vpc.awls\n  awless run repo:create_vpc",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		if listRemoteTemplatesFlag {
			exitOn(listRemoteTemplates())
			return nil
		}
		if len(args) < 1 {
			return errors.New("missing PATH arg (filepath or url)")
		}

		content, err := getTemplateText(args[0])
		exitOn(err)

		logger.Verbosef("Loaded template text:\n\n%s\n", removeComments(content))

		templ, err := template.Parse(string(content))
		exitOn(err)

		extraParams, err := template.ParseParams(strings.Join(args[1:], " "))
		exitOn(err)

		tplExec := &template.TemplateExecution{
			Template: templ,
			Locale:   config.GetAWSRegion(),
			Source:   templ.String(),
		}

		exitOn(runTemplate(tplExec, config.Defaults, extraParams))

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

func runTemplate(tplExec *template.TemplateExecution, fillers ...map[string]interface{}) error {
	env := template.NewEnv()
	env.Log = logger.DefaultLogger
	env.AddFillers(fillers...)
	env.DefLookupFunc = awsdriver.AWSLookupDefinitions
	env.AliasFunc = resolveAliasFunc
	env.MissingHolesFunc = missingHolesStdinFunc()

	if len(env.Fillers) > 0 {
		logger.ExtraVerbosef("default/given holes fillers: %s", sprintProcessedParams(env.Fillers))
	}

	var err error
	tplExec.Template, env, err = template.Compile(tplExec.Template, env)
	exitOn(err)

	tplExec.Fillers = env.GetProcessedFillers()

	validateTemplate(tplExec.Template)

	var drivers []driver.Driver
	for _, s := range cloud.ServiceRegistry {
		drivers = append(drivers, s.Drivers()...)
	}
	awsDriver := driver.NewMultiDriver(drivers...)

	awsDriver.SetLogger(logger.DefaultLogger)

	if err = tplExec.Template.DryRun(awsDriver); err != nil {
		switch t := err.(type) {
		case *template.Errors:
			errs, _ := t.Errors()
			for _, e := range errs {
				logger.Errorf(e.Error())
			}
		}
		exitOn(errors.New("Dryrun failed"))
	}

	fmt.Printf("%s\n", renderGreenFn(tplExec.Template))

	var yesorno string
	if forceGlobalFlag {
		yesorno = "y"
	} else {
		fmt.Println()
		if scheduleFlag {
			fmt.Print("Confirm scheduling? (y/n): ")
		} else {
			fmt.Print("Confirm? (y/n): ")
		}
		_, err = fmt.Scanln(&yesorno)
		exitOn(err)
	}

	if strings.TrimSpace(yesorno) == "y" {
		me, err := aws.AccessService.(*aws.Access).GetIdentity()
		if err != nil {
			logger.Warningf("cannot resolve template author identity: %s", err)
		} else {
			tplExec.Author = me.ResourcePath
			logger.ExtraVerbosef("resolved template author: %s", tplExec.Author)
		}

		if scheduleFlag {
			exitOn(scheduleTemplate(tplExec.Template, scheduleRunInFlag, scheduleRevertInFlag))
			return nil
		}
		tplExec.Template, err = tplExec.Template.Run(awsDriver)
		if err != nil {
			logger.Errorf("Running template error: %s", err)
		}

		printer := template.NewDefaultPrinter(os.Stdout)
		printer.RenderKO = renderRedFn
		printer.RenderOK = renderGreenFn
		printer.Print(tplExec)

		if err = database.Execute(func(db *database.DB) error {
			return db.AddTemplate(tplExec)
		}); err != nil {
			logger.Errorf("Cannot save executed template in awless logs: %s", err)
		}

		if template.IsRevertible(tplExec.Template) {
			fmt.Println()
			logger.Infof("Revert this template with `awless revert %s`", tplExec.Template.ID)
		}

		runSyncFor(tplExec.Template)
	}

	return nil
}

func validateTemplate(tpl *template.Template) {
	unicityRule := &template.UniqueNameValidator{LookupGraph: func(key string) (*graph.Graph, bool) {
		g := sync.LoadCurrentLocalGraph(aws.ServicePerResourceType[key])
		return g, true
	}}

	errs := tpl.Validate(unicityRule, &template.ParamIsSetValidator{Action: "create", Entity: "instance", Param: "keypair", WarningMessage: "This instance has no access keypair. You might not be able to connect to it. Use `awless create instance keypair=my-keypair ...`"})

	if len(errs) > 0 {
		for _, err := range errs {
			logger.Warning(err)
		}
		fmt.Fprintln(os.Stderr)
	}
}

func createDriverCommands(action string, entities []string) *cobra.Command {
	actionCmd := &cobra.Command{
		Use:               fmt.Sprintf("%s ENTITY [param=value ...]", action),
		Short:             oneLinerShortDesc(action, entities),
		Long:              fmt.Sprintf("Allow to %s: %v", action, strings.Join(entities, ", ")),
		Annotations:       map[string]string{"one-liner": "true"},
		PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
		PersistentPostRun: applyHooks(verifyNewVersionHook),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing ENTITY")
			}

			invalidEntityErr := fmt.Errorf("invalid entity '%s'", args[0])

			resources := resolveResourceFromRef(args[0])
			if len(resources) != 1 {
				return invalidEntityErr
			}

			templDef, ok := awsdriver.AWSLookupDefinitions(fmt.Sprintf("%s%s", action, resources[0].Type()))
			if !ok {
				return invalidEntityErr
			}
			templ, err := suggestFixParsingError(templDef, args, invalidEntityErr)
			exitOn(err)

			tplExec := &template.TemplateExecution{
				Template: templ,
				Locale:   config.GetAWSRegion(),
				Source:   templ.String(),
			}

			exitOn(runTemplate(tplExec, config.Defaults))
			return nil
		},
	}

	for _, entity := range entities {
		templDef, ok := awsdriver.AWSLookupDefinitions(fmt.Sprintf("%s%s", action, entity))
		if !ok {
			exitOn(errors.New("command unsupported on inline mode"))
		}
		run := func(def template.Definition) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				text := fmt.Sprintf("%s %s %s", def.Action, def.Entity, strings.Join(args, " "))

				templ, err := template.Parse(text)
				if err != nil {
					templ, err = suggestFixParsingError(def, args, err)
					exitOn(err)
				}

				tplExec := &template.TemplateExecution{
					Template: templ,
					Locale:   config.GetAWSRegion(),
					Source:   templ.String(),
				}

				exitOn(runTemplate(tplExec, config.Defaults))
				return nil
			}
		}
		var apiStr string
		if api, ok := awsdriver.APIPerTemplateDefName[templDef.Name()]; ok {
			apiStr = fmt.Sprint(strings.ToUpper(api) + " ")
		}

		var requiredStr bytes.Buffer
		if len(templDef.Required()) > 0 {
			requiredStr.WriteString("\n\tRequired params:")
			for _, req := range templDef.Required() {
				requiredStr.WriteString(fmt.Sprintf("\n\t\t- %s", req))
				if d, ok := awsdoc.TemplateParamsDoc(templDef.Name(), req); ok {
					requiredStr.WriteString(fmt.Sprintf(": %s", d))
				}
			}
		}

		var extraStr bytes.Buffer
		if len(templDef.Extra()) > 0 {
			extraStr.WriteString("\n\tExtra params:")
			for _, ext := range templDef.Extra() {
				extraStr.WriteString(fmt.Sprintf("\n\t\t- %s", ext))
				if d, ok := awsdoc.TemplateParamsDoc(templDef.Name(), ext); ok {
					extraStr.WriteString(fmt.Sprintf(": %s", d))
				}
			}
		}

		var validArgs []string
		for _, param := range templDef.Required() {
			validArgs = append(validArgs, param+"=")
		}
		for _, param := range templDef.Extra() {
			validArgs = append(validArgs, param+"=")
		}
		actionCmd.AddCommand(
			&cobra.Command{
				Use:               templDef.Entity,
				PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook),
				PersistentPostRun: applyHooks(verifyNewVersionHook),
				Short:             fmt.Sprintf("%s a %s%s", strings.Title(action), apiStr, templDef.Entity),
				Long:              fmt.Sprintf("%s a %s%s%s%s", strings.Title(templDef.Action), apiStr, templDef.Entity, requiredStr.String(), extraStr.String()),
				RunE:              run(templDef),
				ValidArgs:         validArgs,
			},
		)
	}

	return actionCmd
}

func runSyncFor(tpl *template.Template) {
	if !config.GetAutosync() {
		return
	}

	defs := tpl.UniqueDefinitions(awsdriver.AWSLookupDefinitions)

	services := aws.GetCloudServicesForAPIs(defs.Map(
		func(d template.Definition) string { return d.Api },
	)...)

	if _, err := sync.DefaultSyncer.Sync(services...); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Verbosef("performed sync for %s", strings.Join(cloud.Services(services).Names(), ", "))
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
	FILE_EXT            = ".aws"
)

type templateMetadata struct {
	Title, Name string
	Tags        []string
}

func getTemplateText(path string) ([]byte, error) {
	if strings.HasPrefix(path, "repo:") {
		path = fmt.Sprintf("%s/%s", DEFAULT_REPO_PREFIX, strings.TrimPrefix(path[5:], "/"))
		path = fmt.Sprintf("%s%s", strings.TrimSuffix(path, FILE_EXT), FILE_EXT)
	}

	if strings.HasPrefix(path, "http") {
		logger.ExtraVerbosef("fetching remote template at '%s'", path)
		return readHttpContent(path)
	}
	return ioutil.ReadFile(path)
}

func removeComments(b []byte) []byte {
	scn := bufio.NewScanner(bytes.NewReader(b))
	var cleaned bytes.Buffer
	for scn.Scan() {
		line := scn.Text()
		if comment, _ := regexp.MatchString(`^\s*#`, line); comment {
			continue
		}
		cleaned.WriteString(line)
		cleaned.WriteByte('\n')
	}

	return cleaned.Bytes()
}

func listRemoteTemplates() error {
	manifestFile, err := readHttpContent(DEFAULT_REPO_PREFIX + "/manifest.json")
	if err != nil {
		return err
	}
	var remoteTemplates []*templateMetadata
	if err = json.Unmarshal(manifestFile, &remoteTemplates); err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Title\tTags\tRun it with")
	fmt.Fprintln(w, "-----\t----\t-----------")
	for _, tpl := range remoteTemplates {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\tawless run repo:%s -v", tpl.Title, strings.Join(tpl.Tags, ","), tpl.Name))
	}
	w.Flush()
	return nil
}

func readHttpContent(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("'%s' when fetching '%s'", resp.Status, path)
	}

	return ioutil.ReadAll(resp.Body)
}

func isQuoted(s string) bool {
	if strings.HasPrefix(s, "@") {
		return isQuoted(s[1:])
	}
	return (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) || strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")
}

func scheduleTemplate(t *template.Template, runIn, revertIn string) error {
	schedClient, err := client.New(config.GetSchedulerURL())
	if err != nil {
		return err
	}
	logger.Verbosef("sending template to scheduler %s", schedClient.ServiceURL)

	if err := schedClient.Post(client.Form{
		Region:   config.GetAWSRegion(),
		RunIn:    runIn,
		RevertIn: revertIn,
		Template: t.String(),
	}); err != nil {
		return fmt.Errorf("Cannot schedule template: %s", err)
	}

	logger.Info("template scheduled successfully")

	return nil
}

func suggestFixParsingError(def template.Definition, args []string, defaultErr error) (*template.Template, error) {
	if len(def.Required()) != 1 || len(args) != 1 {
		return nil, defaultErr
	}

	suggestText := fmt.Sprintf("%s %s %s=%s", def.Action, def.Entity, def.Required()[0], args[0])

	fmt.Printf("Did you mean `awless %s` (y/n)? ", suggestText)
	var yesorno string
	_, err := fmt.Scanln(&yesorno)
	if err != nil {
		return nil, defaultErr
	}

	if yesorno = strings.ToLower(strings.TrimSpace(yesorno)); !(yesorno == "y" || yesorno == "yes") {
		return nil, defaultErr
	}

	templ, err := template.Parse(suggestText)
	if err != nil {
		return templ, err
	}

	return templ, nil
}
