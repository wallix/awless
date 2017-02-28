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
	"net/http"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/cloud"
)

var (
	AccessService, InfraService, StorageService, NotificationService, QueueService cloud.Service

	SecuAPI Security
)

func InitSession(region, profile string) (*session.Session, error) {
	session, err := session.NewSessionWithOptions(session.Options{
		Config:                  awssdk.Config{Region: awssdk.String(region), HTTPClient: &http.Client{Timeout: 1 * time.Second}},
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		Profile:                 profile,
	})
	if err != nil {
		return nil, err
	}

	if _, err = session.Config.Credentials.Get(); err != nil {
		return nil, errors.New("Your AWS credentials seem undefined! AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY need to be exported in your CLI environment\nInstallation documentation is at https://github.com/wallix/awless/wiki/Installation")
	}
	session.Config.Credentials.Get()

	return session, nil
}

func InitServices(region, profile string) error {
	sess, err := InitSession(region, profile)
	if err != nil {
		return err
	}
	AccessService = NewAccess(sess)
	InfraService = NewInfra(sess)
	StorageService = NewStorage(sess)
	SecuAPI = NewSecu(sess)
	NotificationService = NewNotification(sess)
	QueueService = NewQueue(sess)

	cloud.ServiceRegistry[InfraService.Name()] = InfraService
	cloud.ServiceRegistry[AccessService.Name()] = AccessService
	cloud.ServiceRegistry[StorageService.Name()] = StorageService
	cloud.ServiceRegistry[NotificationService.Name()] = NotificationService
	cloud.ServiceRegistry[QueueService.Name()] = QueueService

	return nil
}
