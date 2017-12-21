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
	"fmt"
	"net/http"
	"os"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/logger"
)

func ResolveRegionFromEnv() (region string) {
	var sess *session.Session
	var err error

	if sess, err = newSessionResolver().resolve(); err == nil {
		region = awssdk.StringValue(sess.Config.Region)
	}

	if awsconfig.IsValidRegion(region) {
		fmt.Fprintf(os.Stderr, "Found existing AWS region '%s'. Setting it as your default region.\n", region)
	} else if sess != nil {
		if r, err := ec2metadata.New(sess).Region(); err == nil {
			fmt.Fprintf(os.Stderr, "Found AWS region '%s' from local EC2 instance metadata. Setting it as your default region.\n", r)
			region = r
		}
	}

	if !awsconfig.IsValidRegion(region) {
		region = awsconfig.StdinRegionSelector()
		fmt.Println()
	}

	return
}

type sessionResolver struct {
	region, profile                      string
	profileSetterCallback                func(val string) error
	httpClient                           *http.Client
	credentialHTTPClient                 *http.Client
	logger                               *logger.Logger
	enableRequestsFullLogging            bool
	enableNetworkMonitorRequestsHandlers bool
	enableCredentialResolvers            bool
}

func newSessionResolver() *sessionResolver {
	return &sessionResolver{
		credentialHTTPClient:  &http.Client{Timeout: 1 * time.Second},
		httpClient:            http.DefaultClient,
		profileSetterCallback: func(val string) error { return nil },
		logger:                logger.DiscardLogger,
	}
}

func (s *sessionResolver) withRegion(region string) *sessionResolver {
	s.region = region
	return s
}

func (s *sessionResolver) withProfile(profile string) *sessionResolver {
	s.profile = profile
	return s
}

func (s *sessionResolver) withCredentialResolvers() *sessionResolver {
	s.enableCredentialResolvers = true
	return s
}

func (s *sessionResolver) withProfileSetter(f func(val string) error) *sessionResolver {
	s.profileSetterCallback = f
	return s
}

func (s *sessionResolver) withLogger(l *logger.Logger) *sessionResolver {
	s.logger = l
	return s
}

func (s *sessionResolver) withNetworkMonitor(enableNetworkMonitor bool) *sessionResolver {
	s.enableNetworkMonitorRequestsHandlers = enableNetworkMonitor
	return s
}

func (s *sessionResolver) resolve() (*session.Session, error) {
	session, err := session.NewSessionWithOptions(session.Options{
		Config: awssdk.Config{
			Region:                        awssdk.String(s.region),
			HTTPClient:                    s.credentialHTTPClient,
			CredentialsChainVerboseErrors: awssdk.Bool(true),
		},
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		Profile:                 s.profile,
	})
	if err != nil {
		return nil, err
	}

	if s.enableRequestsFullLogging {
		session.Config = session.Config.WithLogLevel(awssdk.LogDebugWithHTTPBody)
	}

	session.Handlers.Retry.PushFront(func(req *request.Request) {
		if req.IsErrorThrottle() && s.logger != nil {
			s.logger.Verbosef("retrying %s: %s: %s", req.Operation.Name, req.Error.(awserr.Error).Code(), req.Error.(awserr.Error).Message())
		}
	})

	if s.enableNetworkMonitorRequestsHandlers {
		session.Handlers.Send.PushFront(func(r *request.Request) {
			DefaultNetworkMonitor.addRequest(r)
		})
		session.Handlers.Complete.PushBack(func(r *request.Request) {
			DefaultNetworkMonitor.setRequestEnd(r)
		})
	}

	if s.enableCredentialResolvers {
		session.Config.Credentials = credentials.NewCredentials(
			&credentials.ChainProvider{
				VerboseErrors: true,
				Providers: []credentials.Provider{
					&fileCacheProvider{
						creds:   session.Config.Credentials,
						profile: s.profile,
						log:     s.logger,
					},
					&credentialsPrompterProvider{
						profile: s.profile,
						out:     os.Stderr,
						profileSetterCallback: s.profileSetterCallback,
					},
				},
			})

		if _, err = session.Config.Credentials.Get(); err != nil {
			return session, err
		}
	}

	session.Config.HTTPClient = s.httpClient

	return session, nil
}
