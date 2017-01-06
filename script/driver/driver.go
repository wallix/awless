package driver

import "log"

type Driver interface {
	Lookup(...string) DriverFn
	SetLogger(*log.Logger)
}

type DriverFn func(map[string]interface{}) (interface{}, error)
