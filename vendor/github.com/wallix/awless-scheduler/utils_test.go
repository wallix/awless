package main

import (
	"errors"
	"io/ioutil"

	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver"
)

func newTemplateEnv(fn func(string) (template.Definition, bool)) *template.Env {
	env := template.NewEnv()
	env.DefLookupFunc = fn
	return env
}

func createTmpFSStore() store {
	dir, err := ioutil.TempDir("", "scheduler-")
	if err != nil {
		panic(err)
	}

	fs, err := NewFSStore(dir)
	if err != nil {
		panic(err)
	}

	return fs
}

type happyDriver struct {
}

func (*happyDriver) Lookup(...string) (driver.DriverFn, error) {
	return func(ctx driver.Context, params map[string]interface{}) (interface{}, error) {
		return params["name"], nil
	}, nil
}
func (*happyDriver) SetDryRun(bool)           {}
func (*happyDriver) SetLogger(*logger.Logger) {}

type failDriver struct {
}

func (*failDriver) Lookup(...string) (driver.DriverFn, error) {
	return func(ctx driver.Context, params map[string]interface{}) (interface{}, error) {
		return nil, errors.New("mock driver failure")
	}, nil
}
func (*failDriver) SetDryRun(bool)           {}
func (*failDriver) SetLogger(*logger.Logger) {}
