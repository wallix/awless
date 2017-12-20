/* Copyright 2017 WALLIX

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

package awsspec

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/wallix/awless/logger"
)

type AuthenticateRegistry struct {
	_                string `action:"authenticate" entity:"registry" awsAPI:"ecr"`
	logger           *logger.Logger
	graph            cloud.GraphAPI
	api              ecriface.ECRAPI
	Accounts         []*string `templateName:"accounts"`
	NoConfirm        *bool     `templateName:"no-confirm"`
	DisableDockerCmd *bool     `templateName:"no-docker-login"`
}

func (cmd *AuthenticateRegistry) ParamsSpec() params.Spec {
	return params.NewSpec(params.AtLeastOneOf(params.Key("accounts"), params.Key("no-confirm"), params.Key("no-docker-login")))
}

func (cmd *AuthenticateRegistry) ManualRun(renv env.Running) (interface{}, error) {
	input := &ecr.GetAuthorizationTokenInput{}
	var err error

	// Extra params
	if len(cmd.Accounts) > 0 {
		err = setFieldWithType(cmd.Accounts, input, "RegistryIds", awsstringslice)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	output, err := cmd.api.GetAuthorizationToken(input)
	if err != nil {
		return nil, err
	}
	cmd.logger.ExtraVerbosef("ecr.GetAuthorizationToken call took %s", time.Since(start))
	for _, auth := range output.AuthorizationData {
		token := aws.StringValue(auth.AuthorizationToken)
		decoded, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, err
		}
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			return nil, fmt.Errorf("invalid authorization token: expect user:password, got %s", decoded)
		}
		torun := []string{"docker", "login", "--username", credentials[0], "--password", credentials[1], StringValue(auth.ProxyEndpoint)}

		if BoolValue(cmd.DisableDockerCmd) {
			cmd.logger.Infof("Docker authentication command:\n%s", strings.Join(torun, " "))
		} else {
			confirm := !(BoolValue(cmd.NoConfirm))
			if confirm {
				fmt.Fprintf(os.Stderr, "\nDocker authentication command:\n\n%s\n\nDo you want to run this command:(y/n)? ", strings.Join(torun, " "))
				var yesorno string
				_, err := fmt.Scanln(&yesorno)
				if err != nil {
					return nil, err
				}
				if strings.ToLower(yesorno) != "y" {
					return nil, nil
				}
			}
			dockerCmd := exec.Command("docker", torun[1:]...)
			out, err := dockerCmd.Output()
			if err != nil {
				if e, ok := err.(*exec.ExitError); ok {
					return nil, fmt.Errorf("error running docker command: %s", e.Stderr)
				}
				return nil, fmt.Errorf("error running docker command: %s", err)
			}
			if len(out) > 0 {
				cmd.logger.Info(string(out))
			}
		}
	}

	return nil, nil
}
