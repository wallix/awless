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
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	stdsync "sync"
	"text/tabwriter"

	"time"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/wallix/awless-scheduler/client"
	"github.com/wallix/awless/aws/doc"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/match"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/params"
)

var (
	scheduleRunInFlag       string
	scheduleRevertInFlag    string
	runLogMessage           string
	listRemoteTemplatesFlag bool
	noSuggestedParamsFlag   bool
	allSuggestedParamsFlag  bool
)

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVar(&listRemoteTemplatesFlag, "list", false, "List templates available at https://github.com/wallix/awless-templates")
	runCmd.Flags().StringVar(&scheduleRunInFlag, "run-in", "", "Postpone the execution of this template")
	runCmd.Flags().StringVar(&scheduleRevertInFlag, "revert-in", "", "Schedule the revertion of this template")
	runCmd.Flags().StringVarP(&runLogMessage, "message", "m", "", "Add a message for this template execution to be persisted in your logs")

	var actions []string
	for a := range awsspec.DriverSupportedActions {
		actions = append(actions, a)
	}
	sort.Strings(actions)

	for _, action := range actions {
		entities := awsspec.DriverSupportedActions[action]
		sort.Strings(entities)
		cmd := createDriverCommands(action, entities)
		cmd.PersistentFlags().StringVar(&scheduleRunInFlag, "run-in", "", "Postpone the execution of this command")
		cmd.PersistentFlags().StringVar(&scheduleRevertInFlag, "revert-in", "", "Schedule the revertion of this command")
		RootCmd.AddCommand(cmd)
	}
}

const maxMsgLen = 140

var runCmd = &cobra.Command{
	Use:               "run PATH",
	Short:             "Run a template given a filepath or URL",
	Example:           "  awless run ~/templates/my-infra.txt\n  awless run https://raw.githubusercontent.com/wallix/awless-templates/master/create_vpc.awls\n  awless run repo:create_vpc",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		if listRemoteTemplatesFlag {
			exitOn(listRemoteTemplates())
			return nil
		}
		if len(args) < 1 {
			return errors.New("missing PATH arg (filepath or url)")
		}

		if len(runLogMessage) > maxMsgLen {
			exitOn(fmt.Errorf("message to be persisted should not exceed %d characters", maxMsgLen))
		}

		content, fullPath, err := getTemplateText(args[0])
		exitOn(err)

		logger.Verbosef("Loaded template text:\n\n%s\n", removeComments(content))

		templ, err := template.Parse(string(content))
		exitOn(err)

		extraParams, err := template.ParseParams(strings.Join(args[1:], " "))
		exitOn(err)

		tplExec := &template.TemplateExecution{
			Template: templ,
			Path:     fullPath,
			Message:  strings.TrimSpace(runLogMessage),
			Locale:   config.GetAWSRegion(),
			Profile:  config.GetAWSProfile(),
			Source:   templ.String(),
		}

		exitOn(NewRunner(tplExec.Template, tplExec.Message, tplExec.Path, config.Defaults, extraParams).Run())

		return nil
	},
}

func missingHolesStdinFunc() func(string, []string, bool) string {
	var count int
	return func(hole string, paramPaths []string, optional bool) (response string) {
		if count < 1 {
			fmt.Println("Please specify (Ctrl+C to quit, Tab for completion, Enter to skip optionals):")
		}
		var docs, enums []string
		var typedParam *awsdoc.ParamType
		for _, param := range paramPaths {
			splits := strings.Split(param, ".")
			if len(splits) != 3 {
				continue
			}
			if doc, hasDoc := awsdoc.TemplateParamsDoc(splits[0], splits[1], splits[2]); hasDoc {
				docs = append(docs, doc)
			}
			if enum, hasEnum := awsdoc.EnumDoc[param]; hasEnum {
				enums = append(enums, enum...)
			}
			if tparam, has := awsdoc.ParamTypeDoc[param]; has {
				typedParam = tparam
			}
		}
		if len(docs) > 0 {
			fmt.Fprintln(os.Stderr, strings.Join(docs, "; ")+":")
		}

		autocomplete := holeAutoCompletion(allGraphsOnce.mustLoad(), paramPaths)
		if typedParam != nil {
			autocomplete = typedParamCompletionFunc(allGraphsOnce.mustLoad(), typedParam.ResourceType, typedParam.PropertyName)
		}

		if len(enums) > 0 {
			autocomplete = enumCompletionFunc(enums)
		}

		var promptSuffix string
		if optional {
			promptSuffix = " (optional)"
		}
		var err error
		for response, err = askHole(hole, promptSuffix, autocomplete); err != nil; response, err = askHole(hole, promptSuffix, autocomplete) {
			if optional {
				return ""
			}
			logger.Error(err)
		}
		count++
		return
	}
}

func askHole(hole, promptSuffix string, autocomplete readline.AutoCompleter) (string, error) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          renderCyanBoldFn(hole+"?") + renderYellowFn(promptSuffix) + " ",
		AutoComplete:    autocomplete,
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
		if line != "" {
			return line, nil
		}
		return "", errors.New("required value")
	}
	return "", errors.New("required value")
}

type onceLoader struct {
	g    cloud.GraphAPI
	err  error
	once stdsync.Once
}

func (l *onceLoader) mustLoad() cloud.GraphAPI {
	l.once.Do(func() {
		l.g, l.err = sync.LoadLocalGraphs(config.GetAWSProfile(), config.GetAWSRegion())
	})
	exitOn(l.err)
	return l.g
}

var allGraphsOnce = &onceLoader{}

func createDriverCommands(action string, entities []string) *cobra.Command {
	actionCmd := &cobra.Command{
		Use:               fmt.Sprintf("%s ENTITY [param=value ...]", action),
		Short:             oneLinerShortDesc(action, entities),
		Long:              fmt.Sprintf("Allow to %s: %v", action, strings.Join(entities, ", ")),
		Annotations:       map[string]string{"one-liner": "true"},
		PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
		PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing ENTITY")
			}

			invalidEntityErr := fmt.Errorf("invalid entity '%s'", args[0])

			_, resources, matchingProperty := resolveResourceFromRefInCurrentRegion(args[0])
			if len(resources) != 1 {
				return invalidEntityErr
			}

			templDef, ok := awsspec.AWSLookupDefinitions(fmt.Sprintf("%s%s", action, resources[0].Type()))
			if !ok {
				return invalidEntityErr
			}
			templ, err := suggestFixParsingError(templDef, args, matchingProperty, invalidEntityErr)
			exitOn(err)

			tplExec := &template.TemplateExecution{
				Template: templ,
				Locale:   config.GetAWSRegion(),
				Profile:  config.GetAWSProfile(),
				Source:   templ.String(),
			}

			exitOn(NewRunner(tplExec.Template, tplExec.Message, tplExec.Path).Run())
			return nil
		},
	}

	for _, entity := range entities {
		templDef, ok := awsspec.AWSLookupDefinitions(fmt.Sprintf("%s%s", action, entity))
		if !ok {
			exitOn(errors.New("command unsupported on inline mode"))
		}
		run := func(def awsspec.Definition) func(cmd *cobra.Command, args []string) error {
			return func(cmd *cobra.Command, args []string) error {
				text := fmt.Sprintf("%s %s %s", def.Action, def.Entity, strings.Join(args, " "))

				templ, err := template.Parse(text)
				if err != nil {
					_, resources, matchingProperty := resolveResourceFromRefInCurrentRegion(args[0])
					if len(resources) != 1 {
						exitOn(err)
					}
					templ, err = suggestFixParsingError(def, args, matchingProperty, err)
					exitOn(err)
				}

				tplExec := &template.TemplateExecution{
					Template: templ,
					Locale:   config.GetAWSRegion(),
					Profile:  config.GetAWSProfile(),
					Source:   templ.String(),
				}

				exitOn(NewRunner(tplExec.Template, tplExec.Message, tplExec.Path, config.Defaults).Run())
				return nil
			}
		}
		var apiStr string
		if api, ok := awsspec.APIPerTemplateDefName[templDef.Action+templDef.Entity]; ok {
			apiStr = fmt.Sprint(strings.ToUpper(api) + " ")
		}

		var paramsStr bytes.Buffer
		allParams, optParams, _ := params.List(templDef.Params)
		tab := tabwriter.NewWriter(&paramsStr, 0, 0, 3, '.', 0)
		for _, p := range allParams {
			fmt.Fprintf(tab, "  %s\t", p)
			if d, ok := awsdoc.TemplateParamsDocWithEnums(templDef.Action, templDef.Entity, p); ok {
				fmt.Fprintf(tab, " %s", d)
			}
			fmt.Fprintln(tab)
		}
		for _, p := range optParams {
			fmt.Fprintf(tab, "  [%s]\t", p)
			if d, ok := awsdoc.TemplateParamsDocWithEnums(templDef.Action, templDef.Entity, p); ok {
				fmt.Fprintf(tab, " %s", d)
			}
			fmt.Fprintln(tab)
		}
		tab.Flush()

		var validArgs []string
		for _, param := range append(allParams, optParams...) {
			validArgs = append(validArgs, param+"=")
		}
		currentCmd := &cobra.Command{
			Use:               fmt.Sprintf("%s [param=value ...]", templDef.Entity),
			PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
			PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),
			Short:             awsdoc.AwlessCommandDefinitionsDoc(action, templDef.Entity, fmt.Sprintf("%s a %s%s", strings.Title(action), apiStr, templDef.Entity)),
			Long:              fmt.Sprintf("Params: \n%s\nParams patterns:\n  %s", paramsStr.String(), templDef.Params),
			Example:           awsdoc.AwlessExamplesDoc(action, templDef.Entity),
			RunE:              run(templDef),
			ValidArgs:         validArgs,
		}
		currentCmd.SetUsageTemplate(customCommandUsageTemplate)
		currentCmd.SetHelpTemplate(`{{with .Short}}{{. | trimTrailingWhitespaces}}
{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)
		currentCmd.Flags().BoolVar(&noSuggestedParamsFlag, "prompt-only-required", false, "Prompt only required parameters")
		currentCmd.Flags().BoolVarP(&allSuggestedParamsFlag, "prompt-all", "a", false, "Prompt all non-provided parameters")

		actionCmd.AddCommand(currentCmd)
	}

	return actionCmd
}

func runSyncFor(tplExec *template.TemplateExecution) {
	if !config.GetAutosync() {
		return
	}

	if tplExec.Stats().AllKO() {
		return
	}

	apis := tplExec.Template.UniqueDefinitions(awsspec.APIPerTemplateDefName)

	services := awsservices.GetCloudServicesForAPIs(apis...)

	if !noSyncGlobalFlag {
		go func() { // allow to only display this verbose line only if taking more than 1 second before exiting CLI
			time.Sleep(2 * time.Second)
			logger.Infof("Resyncing %s ... (disable with --no-sync global flag)", joinSentence(cloud.Services(services).Names()))
		}()
	}
	if _, err := sync.DefaultSyncer.Sync(services...); err != nil {
		logger.ExtraVerbosef(err.Error())
	}
}

func resolveAliasFunc(paramPath, alias string) string {
	splits := strings.Split(paramPath, ".")
	if len(splits) != 3 {
		logger.Errorf("resolve alias: invalid param path: %s", paramPath)
		return ""
	}
	entity, key := splits[1], splits[2]
	var typedParam *awsdoc.ParamType
	if tparam, has := awsdoc.ParamTypeDoc[paramPath]; has {
		typedParam = tparam
	}

	gph, err := sync.LoadLocalGraphs(config.GetAWSProfile(), config.GetAWSRegion())
	if err != nil {
		fmt.Printf("resolve alias '%s': cannot load local graphs for region %s: %s\n", alias, config.GetAWSRegion(), err)
		return ""
	}
	resType := key
	if typedParam != nil {
		resType = typedParam.ResourceType
	} else {
		if strings.Contains(key, "id") {
			resType = entity
		}
	}

	resources, err := gph.Find(cloud.NewQuery(resType).Match(match.And(match.Property("Name", alias))))
	if err != nil {
		return ""
	}
	var matchingResource cloud.Resource
	switch len(resources) {
	case 1:
		matchingResource = resources[0]
	default:
		resources, err := gph.FindWithProperties(map[string]interface{}{"Name": alias})
		if err != nil {
			return ""
		}
		if len(resources) > 0 {
			matchingResource = resources[0]
		}
	}
	if matchingResource == nil {
		return ""
	}
	if typedParam != nil {
		if prop, ok := matchingResource.Properties()[typedParam.PropertyName].(string); ok {
			return prop
		}
	}

	return matchingResource.Id()
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
	Title, Name, MinimalVersion string
	Tags                        []string
}

func getTemplateText(path string) (content []byte, expanded string, err error) {
	if strings.HasPrefix(path, "repo:") {
		path = fmt.Sprintf("%s/%s", DEFAULT_REPO_PREFIX, strings.TrimPrefix(path[5:], "/"))
		path = fmt.Sprintf("%s%s", strings.TrimSuffix(path, FILE_EXT), FILE_EXT)
	}

	expanded = path

	if strings.HasPrefix(path, "http") {
		logger.ExtraVerbosef("fetching remote template at '%s'", path)
		content, err = readHttpContent(path)
	} else {
		f, ferr := os.Open(path)
		if ferr != nil {
			return nil, "", ferr
		}
		defer f.Close()

		var perr error
		expanded, perr = filepath.Abs(f.Name())
		if perr != nil {
			expanded = path
		}
		content, err = ioutil.ReadAll(f)
	}

	if err != nil {
		return content, expanded, err
	}

	requiredVersion, ok := detectMinimalVersionInTemplate(content)
	if ok {
		comp, _ := config.CompareSemver(requiredVersion, config.Version)
		if comp > 0 {
			return content, expanded, fmt.Errorf("This template has metadata indicating to be parsed with at least awless version %s. Your current version is %s", requiredVersion, config.Version)
		}
	}

	return content, expanded, nil
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

var (
	minimalVersionRegex = regexp.MustCompile(`^# *MinimalVersion: *(v?\d{1,3}\.\d{1,3}\.\d{1,3}).*$`)
)

func detectMinimalVersionInTemplate(content []byte) (string, bool) {
	scn := bufio.NewScanner(bytes.NewReader(content))
	for scn.Scan() {
		matches := minimalVersionRegex.FindStringSubmatch(scn.Text())
		if len(matches) > 1 {
			logger.ExtraVerbosef("detected minimal version %s in templates", matches[1])
			return matches[1], true
		}
	}
	return "", false
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
		if tpl.MinimalVersion == "" {
			tpl.MinimalVersion = config.Version
		}
		if comp, err := config.CompareSemver(tpl.MinimalVersion, config.Version); comp < 1 && err == nil {
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\tawless run repo:%s -v", tpl.Title, strings.Join(tpl.Tags, ","), tpl.Name))
		}
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

func isCSV(s string) bool {
	if strings.HasPrefix(s, "[") {
		if !strings.HasSuffix(s, "]") {
			return false
		}
		s = s[1 : len(s)-1]
	}
	for _, split := range strings.Split(s, ",") {
		if !template.MatchStringParamValue(split) {
			return false
		}
	}
	return true
}

func scheduleTemplate(t *template.Template, runIn, revertIn string) error {
	schedClient, err := client.New(config.GetSchedulerURL())
	if err != nil {
		return fmt.Errorf("cannot connect to scheduler: %s", err)
	}
	logger.Verbosef("sending template to scheduler %s", schedClient.ServiceURL)

	if err := schedClient.Post(client.Form{
		Region:   config.GetAWSRegion(),
		RunIn:    runIn,
		RevertIn: revertIn,
		Template: t.String(),
	}); err != nil {
		return fmt.Errorf("cannot schedule template: %s", err)
	}

	logger.Info("template scheduled successfully")

	return nil
}

func suggestFixParsingError(def awsspec.Definition, args []string, matchingProperty string, defaultErr error) (*template.Template, error) {
	if len(def.Params.Required()) != 1 || len(args) != 1 {
		return nil, defaultErr
	}
	propKey := def.Params.Required()[0]
	propValue := args[0]
	if matchingProperty == properties.Name && !strings.HasPrefix(propValue, "@") && !strings.HasSuffix(propKey, "name") {
		propValue = "@" + propValue
	}

	suggestText := fmt.Sprintf("%s %s %s=%s", def.Action, def.Entity, propKey, propValue)

	if !promptConfirmDefaultYes("Did you mean `awless %s` ? ", suggestText) {
		return nil, defaultErr
	}

	templ, err := template.Parse(suggestText)
	if err != nil {
		return templ, err
	}

	return templ, nil
}

func isSchedulingMode() bool {
	runin := strings.TrimSpace(scheduleRunInFlag)
	revertin := strings.TrimSpace(scheduleRevertInFlag)

	if runin != "" || revertin != "" {
		return true
	}
	return false
}

func joinSentence(arr []string) string {
	sep := ", "
	if ln := len(arr); ln > 1 {
		return fmt.Sprintf("%s and %s", strings.Join(arr[:ln-1], sep), arr[ln-1])
	}
	return strings.Join(arr, sep)
}

const customCommandUsageTemplate = `
Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}

{{.Long}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func promptConfirmDefaultYes(msg string, a ...interface{}) bool {
	var yesorno string
	fmt.Fprintf(os.Stderr, "%s [Y/n] ", fmt.Sprintf(msg, a...))
	fmt.Scanln(&yesorno)
	if y := strings.TrimSpace(strings.ToLower(yesorno)); y == "y" || y == "yes" || y == "" {
		return true
	}
	return false
}
