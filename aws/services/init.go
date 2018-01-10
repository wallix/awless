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

	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync"
)

var (
	AccessService, InfraService, StorageService, MessagingService, DnsService, LambdaService, MonitoringService, CdnService, CloudformationService cloud.Service
)

func Init(profile, region string, extraConf map[string]interface{}, log *logger.Logger, profileSetterCallback func(val string) error, enableNetworkMonitor bool) error {
	if region == "" {
		return errors.New("empty AWS region. Set it with `awless config set aws.region`")
	}

	sb := newSessionResolver().withRegion(region).withProfile(profile).withNetworkMonitor(enableNetworkMonitor)
	sb = sb.withProfileSetter(profileSetterCallback).withLogger(log).withCredentialResolvers()

	sess, err := sb.resolve()
	if err != nil {
		return err
	}

	AccessService = NewAccess(sess, profile, extraConf, log)
	InfraService = NewInfra(sess, profile, extraConf, log)
	StorageService = NewStorage(sess, profile, extraConf, log)
	MessagingService = NewMessaging(sess, profile, extraConf, log)
	DnsService = NewDns(sess, profile, extraConf, log)
	LambdaService = NewLambda(sess, profile, extraConf, log)
	MonitoringService = NewMonitoring(sess, profile, extraConf, log)
	CdnService = NewCdn(sess, profile, extraConf, log)
	CloudformationService = NewCloudformation(sess, profile, extraConf, log)

	cloud.ServiceRegistry[InfraService.Name()] = InfraService
	cloud.ServiceRegistry[AccessService.Name()] = AccessService
	cloud.ServiceRegistry[StorageService.Name()] = StorageService
	cloud.ServiceRegistry[MessagingService.Name()] = MessagingService
	cloud.ServiceRegistry[DnsService.Name()] = DnsService
	cloud.ServiceRegistry[LambdaService.Name()] = LambdaService
	cloud.ServiceRegistry[MonitoringService.Name()] = MonitoringService
	cloud.ServiceRegistry[CdnService.Name()] = CdnService
	cloud.ServiceRegistry[CloudformationService.Name()] = CloudformationService

	awsspec.CommandFactory = &awsspec.AWSFactory{
		Log:  log,
		Sess: sess,
		Graph: &cloud.LazyGraph{LoadingFunc: func() cloud.GraphAPI {
			g, err := sync.LoadLocalGraphs(profile, region)
			if err != nil || g == nil {
				g = graph.NewGraph()
			}
			return g
		}},
	}

	return nil
}

func getBool(m map[string]interface{}, key string, def bool) bool {
	if b, ok := m[key].(bool); ok {
		return b
	}
	return def
}
