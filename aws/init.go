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
	"errors"
	"fmt"
	"net/http"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
)

var (
	AccessService, InfraService, StorageService, MessagingService, DnsService, LambdaService, MonitoringService, CdnService, CloudformationService cloud.Service
)

func InitServices(conf map[string]interface{}, log *logger.Logger) error {
	awsconf := config(conf)
	region := awsconf.region()
	if region == "" {
		return errors.New("empty AWS region. Set it with `awless config set aws.region`")
	}

	sess, err := initAWSSession(region, awsconf.profile())
	if err != nil {
		return err
	}
	//addDebugRequestsHandlers(sess)

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

	sess, err := initAWSSession(region, profile)
	if err != nil {
		return nil, err
	}

	drivLog := logger.DiscardLogger
	if len(log) > 0 {
		drivLog = log[0]
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

func initAWSSession(region, profile string) (*session.Session, error) {
	session, err := session.NewSessionWithOptions(session.Options{
		Config: awssdk.Config{
			Region:     awssdk.String(region),
			HTTPClient: &http.Client{Timeout: 2 * time.Second},
		},
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		Profile:                 profile,
	})

	//session.Config = session.Config.WithLogLevel(awssdk.LogDebugWithHTTPBody)
	if err != nil {
		return nil, err
	}

	if _, err = session.Config.Credentials.Get(); err != nil {
		fmt.Printf("Cannot resolve AWS credentials for profile '%s' (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY)\n", profile)
		creds := awsdriver.NewCredsPrompter(profile)
		if err := creds.Prompt(); err != nil {
			return nil, fmt.Errorf("prompting credentials: %s", err)
		}
		created, err := creds.Store()
		if err != nil {
			return nil, fmt.Errorf("storing credentials at '%s': %s", awsdriver.AWSCredFilepath, err)
		}
		if created {
			fmt.Printf("\n\u2713 %s created", awsdriver.AWSCredFilepath)
			fmt.Printf("\n\u2713 Credentials for profile '%s' stored successfully\n", creds.Profile)
		} else {
			fmt.Printf("\n\u2713 Credentials for profile '%s' stored successfully in %s\n", creds.Profile, awsdriver.AWSCredFilepath)
		}
		fmt.Printf("\u2713 Reloading session...\n\n")
		return initAWSSession(region, profile)
	}

	session.Config.HTTPClient = http.DefaultClient

	return session, nil
}

func addDebugRequestsHandlers(s *session.Session) {
	s.Handlers.Send.PushFront(func(r *request.Request) {
		DefaultNetworkMonitor.addRequest(r)
	})
	s.Handlers.Complete.PushBack(func(r *request.Request) {
		DefaultNetworkMonitor.setRequestEnd(r)
	})
}
