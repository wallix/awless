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

package aws

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
)

type AwsDriver struct {
	iam iamiface.IAMAPI
	ec2 ec2iface.EC2API
	s3  s3iface.S3API

	dryRun bool
	logger *logger.Logger
}

func NewDriver(ec2, iam, s3 interface{}) *AwsDriver {
	return &AwsDriver{
		ec2:    ec2.(ec2iface.EC2API),
		iam:    iam.(iamiface.IAMAPI),
		s3:     s3.(s3iface.S3API),
		logger: logger.DiscardLogger,
	}
}

func (d *AwsDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *AwsDriver) SetLogger(l *logger.Logger) { d.logger = l }

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
