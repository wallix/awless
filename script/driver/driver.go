package driver

import "log"

type Driver interface {
	Lookup(...string) DriverFn
	SetLogger(*log.Logger)
}

type DriverFn func(map[string]interface{}) error

var (
	CREATE = "CREATE"
	DELETE = "DELETE"
	VERIFY = "VERIFY"

	INSTANCE = "INSTANCE"
	SUBNET   = "SUBNET"
	VPC      = "VPC"

	COUNT      = "COUNT"
	BASE       = "BASE"
	TYPE       = "TYPE"
	CIDR       = "CIDR"
	REFERENCES = "REFERENCES"
	REF        = "REF"
)
