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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/logger"
)

type cachedCredential struct {
	credentials.Value
	Expiration time.Time
}

func (c *cachedCredential) isExpired() bool {
	return c.Expiration.Before(time.Now().UTC())
}

type fileCacheProvider struct {
	creds   *credentials.Credentials
	curr    *cachedCredential
	profile string
	log     *logger.Logger
}

func (f *fileCacheProvider) Retrieve() (credentials.Value, error) {
	awlessCache := os.Getenv("__AWLESS_CACHE")
	if awlessCache == "" {
		return f.creds.Get()
	}
	credFolder := filepath.Join(awlessCache, "credentials")
	fold := &folder{credFolder}
	credFile := fmt.Sprintf("aws-profile-%s.json", f.profile)

	if content, ok := fold.getFileContent(credFile); ok {
		var cached *cachedCredential
		if err := json.Unmarshal(content, &cached); err != nil {
			return credentials.Value{}, err
		}
		f.log.ExtraVerbosef("loading credentials from '%s'", filepath.Join(credFolder, credFile))
		if !cached.isExpired() {
			f.curr = cached
			return cached.Value, nil
		} else {
			f.creds.Expire()
		}
	}
	credValue, err := f.creds.Get()
	if err != nil {
		if batcherr, ok := err.(awserr.BatchedErrors); !ok || batcherr.Code() != "NoCredentialProviders" {
			if failure, ok := err.(awserr.RequestFailure); ok {
				f.log.Errorf("%s: %s\n", failure.Code(), failure.Message())
			} else {
				f.log.Errorf("%s\n", err)
			}
		}
		return credValue, err
	}

	switch credValue.ProviderName {
	case stscreds.ProviderName:
		cred := &cachedCredential{credValue, time.Now().UTC().Add(stscreds.DefaultDuration)}
		f.curr = cred
		content, err := json.Marshal(cred)
		if err != nil {
			return credValue, err
		}
		if err = fold.putFileContent(credFile, content); err != nil {
			return credValue, fmt.Errorf("error writing cache file: %s", err.Error())
		}
		f.log.ExtraVerbosef("credentials cached in '%s'", filepath.Join(credFolder, credFile))
		return credValue, nil
	}
	return credValue, nil
}

func (f *fileCacheProvider) IsExpired() bool {
	if f.curr != nil {
		return f.curr.isExpired()
	}
	return f.creds.IsExpired()
}

type folder struct {
	path string
}

func (f *folder) getFileContent(filename string) (content []byte, ok bool) {
	if _, err := os.Stat(f.path); err != nil {
		return
	}
	credPath := filepath.Join(f.path, filename)

	if _, readerr := os.Stat(credPath); readerr != nil {
		return
	}
	var err error
	if content, err = ioutil.ReadFile(credPath); err != nil {
		return
	}
	ok = true
	return
}

func (f *folder) putFileContent(filename string, content []byte) error {
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		os.MkdirAll(f.path, 0700)
	}

	return ioutil.WriteFile(filepath.Join(f.path, filename), content, 0600)
}

type credentialsPrompterProvider struct {
	profile               string
	out                   io.Writer
	profileSetterCallback func(val string) error
	retrieved             bool
}

func (c *credentialsPrompterProvider) Retrieve() (credentials.Value, error) {
	c.retrieved = false
	fmt.Fprintf(c.out, "Cannot resolve AWS credentials for profile '%s' (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY)", c.profile)
	creds := awsspec.NewCredsPrompter(c.profile)
	creds.ProfileSetterCallback = c.profileSetterCallback
	if err := creds.Prompt(); err != nil {
		return credentials.Value{}, fmt.Errorf("prompting credentials: %s", err)
	}
	created, err := creds.Store()
	if err != nil {
		return credentials.Value{}, fmt.Errorf("storing credentials at '%s': %s", awsspec.AWSCredFilepath, err)
	}
	if created {
		fmt.Fprintf(c.out, "\n\u2713 %s created", awsspec.AWSCredFilepath)
		fmt.Fprintf(c.out, "\n\u2713 Credentials for profile '%s' stored successfully\n", creds.Profile)
	} else {
		fmt.Fprintf(c.out, "\n\u2713 Credentials for profile '%s' stored successfully in %s\n", creds.Profile, awsspec.AWSCredFilepath)
	}
	c.retrieved = true
	return creds.Val, nil
}

func (c *credentialsPrompterProvider) IsExpired() bool {
	return !c.retrieved
}
