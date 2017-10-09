package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
	"github.com/wallix/awless/template"
)

func NewRunner(tpl *template.Template, msg, tplPath string, fillers ...map[string]interface{}) *template.Runner {
	runner := &template.Runner{}

	runner.Template = tpl
	runner.Locale = config.GetAWSRegion()
	runner.Profile = config.GetAWSProfile()
	runner.Log = logger.DefaultLogger
	runner.Message = msg
	runner.TemplatePath = tplPath
	runner.Fillers = fillers
	runner.AliasFunc = resolveAliasFunc
	runner.MissingHolesFunc = missingHolesStdinFunc()

	runner.Validators = []template.Validator{
		&template.UniqueNameValidator{LookupGraph: func(key string) (*graph.Graph, bool) {
			g := sync.LoadLocalGraphForService(awsservices.ServicePerResourceType[key], config.GetAWSRegion())
			return g, true
		}},
		&template.ParamIsSetValidator{Action: "create", Entity: "instance", Param: "keypair", WarningMessage: "This instance has no access keypair. You might not be able to connect to it. Use `awless create instance keypair=my-keypair ...`"},
	}

	runner.CmdLookuper = func(tokens ...string) interface{} {
		newCommandFunc := awsspec.CommandFactory.Build(strings.Join(tokens, ""))
		if newCommandFunc == nil {
			return nil
		}
		return newCommandFunc()
	}

	runner.BeforeRun = func(tplExec *template.TemplateExecution) (bool, error) {
		fmt.Printf("%s\n", renderGreenFn(tplExec.Template))

		var yesorno string
		if forceGlobalFlag {
			yesorno = "y"
		} else {
			fmt.Println()
			if isSchedulingMode() {
				fmt.Print("Confirm scheduling? (y/n): ")
			} else {
				fmt.Print("Confirm? (y/n): ")
			}
			if _, err := fmt.Scanln(&yesorno); err != nil {
				return false, err
			}
		}

		if strings.TrimSpace(yesorno) == "y" {
			me, err := awsservices.AccessService.(*awsservices.Access).GetIdentity()
			if err != nil {
				logger.Warningf("cannot resolve template author identity: %s", err)
			} else {
				tplExec.Author = me.ResourcePath
				logger.ExtraVerbosef("resolved template author: %s", tplExec.Author)
			}
			if isSchedulingMode() {
				return false, scheduleTemplate(tplExec.Template, scheduleRunInFlag, scheduleRevertInFlag)
			}
			return true, nil
		}

		return false, errors.New("Did not confirm template execution")
	}

	runner.AfterRun = func(tplExec *template.TemplateExecution) error {
		newDefaultTemplatePrinter(os.Stdout).print(tplExec)

		if tplExec.Message == "" {
			if tplExec.IsOneLiner() {
				tplExec.SetMessage(fmt.Sprintf("Run %s", tplExec.Template))
			} else if path := tplExec.Path; path != "" {
				stats := tplExec.Stats()
				if stats.KOCount > 0 {
					tplExec.SetMessage(fmt.Sprintf("Run %d/%d commands from %s", stats.OKCount, stats.CmdCount, path))
				} else {
					tplExec.SetMessage(fmt.Sprintf("Run %d commands from %s", stats.OKCount, path))
				}
			}
		}

		if err := database.Execute(func(db *database.DB) error {
			return db.AddTemplate(tplExec)
		}); err != nil {
			logger.Errorf("Cannot save executed template in awless logs: %s", err)
		}

		if template.IsRevertible(tplExec.Template) {
			fmt.Println()
			logger.Infof("Revert this template with `awless revert %s -r %s -p %s`", tplExec.Template.ID, config.GetAWSRegion(), config.GetAWSProfile())
		}

		runSyncFor(tplExec)

		return nil
	}

	return runner
}
