package aws

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/script/driver"
)

type AwsDriver struct {
	api    ec2iface.EC2API
	dryRun bool
	logger *log.Logger
}

func NewDriver(api ec2iface.EC2API) *AwsDriver {
	return &AwsDriver{
		api:    api,
		logger: log.New(ioutil.Discard, "", 0),
	}
}

func (d *AwsDriver) SetDryRun(dry bool)      { d.dryRun = dry }
func (d *AwsDriver) SetLogger(l *log.Logger) { d.logger = l }

func (d *AwsDriver) Lookup(lookups ...string) driver.DriverFn {
	if len(lookups) < 2 {
		panic("need at least 2 string to lookup driver method")
	}

	var format string
	if d.dryRun {
		format = "%s_%s_DryRun"
	} else {
		format = "%s_%s"
	}

	fnName := fmt.Sprintf(format, humanize(lookups[0]), humanize(lookups[1]))
	method := reflect.ValueOf(d).MethodByName(fnName).Interface()

	driverFn, converted := method.(func(map[string]interface{}) (interface{}, error))
	if !converted {
		panic(fmt.Sprintf("method '%s' found on '%T' is not a driver function", fnName, d))
	}

	return driverFn
}

func humanize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	return strings.ToUpper(s)
}
