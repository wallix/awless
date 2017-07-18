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

package awsservices

import (
	"errors"
	"fmt"

	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
)

var (
	AccessService, InfraService, StorageService, MessagingService, DnsService, LambdaService, MonitoringService, CdnService, CloudformationService cloud.Service
)

func Init(conf map[string]interface{}, log *logger.Logger, profileSetterCallback func(val string) error) error {
	awsconf := config(conf)
	region := awsconf.region()
	if region == "" {
		return errors.New("empty AWS region. Set it with `awless config set aws.region`")
	}

	sb := newSessionResolver().withRegion(region).withProfile(awsconf.profile())
	sb = sb.withProfileSetter(profileSetterCallback).withLogger(log).withCredentialResolvers()

	sess, err := sb.resolve()
	if err != nil {
		return err
	}

	AccessService = NewAccess(sess, awsconf, log)
	InfraService = NewInfra(sess, awsconf, log)
	StorageService = NewStorage(sess, awsconf, log)
	MessagingService = NewMessaging(sess, awsconf, log)
	DnsService = NewDns(sess, awsconf, log)
	LambdaService = NewLambda(sess, awsconf, log)
	MonitoringService = NewMonitoring(sess, awsconf, log)
	CdnService = NewCdn(sess, awsconf, log)
	CloudformationService = NewCloudformation(sess, awsconf, log)

	cloud.ServiceRegistry[InfraService.Name()] = InfraService
	cloud.ServiceRegistry[AccessService.Name()] = AccessService
	cloud.ServiceRegistry[StorageService.Name()] = StorageService
	cloud.ServiceRegistry[MessagingService.Name()] = MessagingService
	cloud.ServiceRegistry[DnsService.Name()] = DnsService
	cloud.ServiceRegistry[LambdaService.Name()] = LambdaService
	cloud.ServiceRegistry[MonitoringService.Name()] = MonitoringService
	cloud.ServiceRegistry[CdnService.Name()] = CdnService
	cloud.ServiceRegistry[CloudformationService.Name()] = CloudformationService

	return nil
}

func NewDriver(region, profile string, log ...*logger.Logger) (driver.Driver, error) {
	if !awsconfig.IsValidRegion(region) {
		return nil, fmt.Errorf("invalid region '%s' provided", region)
	}

	drivLog := logger.DiscardLogger
	if len(log) > 0 {
		drivLog = log[0]
	}

	sb := newSessionResolver().withRegion(region).withProfile(profile).withLogger(drivLog).withCredentialResolvers()

	sess, err := sb.resolve()
	if err != nil {
		return nil, err
	}

	awsconf := config(
		map[string]interface{}{"aws.region": region, "aws.profile": profile},
	)

	var drivers []driver.Driver
	drivers = append(drivers, NewAccess(sess, awsconf, drivLog).Drivers()...)
	drivers = append(drivers, NewInfra(sess, awsconf, drivLog).Drivers()...)
	drivers = append(drivers, NewStorage(sess, awsconf, drivLog).Drivers()...)
	drivers = append(drivers, NewMessaging(sess, awsconf, drivLog).Drivers()...)
	drivers = append(drivers, NewDns(sess, awsconf, drivLog).Drivers()...)
	drivers = append(drivers, NewLambda(sess, awsconf, drivLog).Drivers()...)
	drivers = append(drivers, NewMonitoring(sess, awsconf, drivLog).Drivers()...)
	drivers = append(drivers, NewCdn(sess, awsconf, drivLog).Drivers()...)
	drivers = append(drivers, NewCloudformation(sess, awsconf, drivLog).Drivers()...)

	return driver.NewMultiDriver(drivers...), nil
}
