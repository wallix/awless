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

	env := NewEnv()
	env.Log = ru.Log
	env.AddFillers(ru.Fillers...)
	env.AliasFunc = ru.AliasFunc
	env.MissingHolesFunc = ru.MissingHolesFunc
	env.Lookuper = ru.CmdLookuper

	var err error
	tplExec.Template, env, err = Compile(tplExec.Template, env, NewRunnerCompileMode)
	if err != nil {
		return err
	}

	tplExec.Fillers = env.GetProcessedFillers()

	errs := tplExec.Template.Validate(ru.Validators...)
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Warning(err)
		}
		fmt.Fprintln(os.Stderr)
	}

	logger.Info("Dry running template ...")
	env.IsDryRun = true
	if _, err = tplExec.Template.Run(env); err != nil {
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
	env.IsDryRun = false

	ok, err := ru.BeforeRun(tplExec)
	if err != nil {
		return err
	}

	if ok {
		tplExec.Template, err = tplExec.Template.Run(env)
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
