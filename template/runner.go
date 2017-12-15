package template

import (
	"errors"
	"fmt"
	"os"

	"github.com/wallix/awless/logger"
)

type Runner struct {
	Template                               *Template
	Locale, Profile, Message, TemplatePath string
	Log                                    *logger.Logger
	Fillers                                []map[string]interface{}
	AliasFunc                              func(entity, key, alias string) string
	MissingHolesFunc                       func(string, []string) interface{}
	CmdLookuper                            func(tokens ...string) interface{}
	Validators                             []Validator

	BeforeRun func(*TemplateExecution) (bool, error)
	AfterRun  func(*TemplateExecution) error
}

func (ru *Runner) Run() error {
	tplExec := &TemplateExecution{
		Template: ru.Template,
		Path:     ru.TemplatePath,
		Locale:   ru.Locale,
		Profile:  ru.Profile,
		Source:   ru.Template.String(),
	}
	tplExec.SetMessage(ru.Message)

	cenv := NewEnv().WithAliasFunc(ru.AliasFunc).WithMissingHolesFunc(ru.MissingHolesFunc).
		WithFillers(ru.Fillers...).WithLookupCommandFunc(ru.CmdLookuper).WithLog(ru.Log).Build()

	var err error
	tplExec.Template, cenv, err = Compile(tplExec.Template, cenv, NewRunnerCompileMode)
	if err != nil {
		return err
	}

	tplExec.Fillers = cenv.ProcessedFillers()

	errs := tplExec.Template.Validate(ru.Validators...)
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Warning(err)
		}
		fmt.Fprintln(os.Stderr)
	}

	if tplExec.IsOneLiner() {
		logger.Verbose("Dry running template ...")
	} else {
		logger.Info("Dry running template ...")
	}

	renv := NewRunEnv(cenv)
	if _, err = tplExec.Template.DryRun(renv); err != nil {
		switch t := err.(type) {
		case *Errors:
			errs, _ := t.Errors()
			for _, e := range errs {
				logger.Errorf(e.Error())
			}
		default:
			logger.Error(err)
		}
		return errors.New("Dry run failed")
	}

	ok, err := ru.BeforeRun(tplExec)
	if err != nil {
		return err
	}

	if ok {
		tplExec.Template, err = tplExec.Template.Run(renv)
		if err != nil {
			logger.Errorf("Running template error: %s", err)
		}
		if err := ru.AfterRun(tplExec); err != nil {
			return err
		}
	}

	if tplExec.Stats().KOCount > 0 {
		os.Exit(1)
	}

	return nil
}
