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

package driver

import (
	"errors"
	"fmt"

	"github.com/wallix/awless/logger"
)

var ErrDriverFnNotFound = errors.New("driver function not found")

type Driver interface {
	Lookup(...string) (DriverFn, error)
	SetDryRun(bool)
	SetLogger(*logger.Logger)
}

type DriverFn func(map[string]interface{}) (interface{}, error)

type MultiDriver struct {
	drivers []Driver
}

func NewMultiDriver(drivers ...Driver) Driver {
	return &MultiDriver{drivers: drivers}
}

func (d *MultiDriver) SetDryRun(dry bool) {
	for _, dr := range d.drivers {
		dr.SetDryRun(dry)
	}
}

func (d *MultiDriver) SetLogger(l *logger.Logger) {
	for _, dr := range d.drivers {
		dr.SetLogger(l)
	}
}

func (d *MultiDriver) Lookup(lookups ...string) (driverFn DriverFn, err error) {
	var funcs []DriverFn
	for _, dr := range d.drivers {
		fn, err := dr.Lookup(lookups...)
		if err == ErrDriverFnNotFound {
			continue
		}
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, fn)
	}
	switch len(funcs) {
	case 0:
		return nil, fmt.Errorf("function corresponding to '%v' not found in drivers", lookups)
	case 1:
		return funcs[0], nil
	default:
		return nil, fmt.Errorf("%d functions corresponding to '%v' found in drivers", len(funcs), lookups)
	}
}
