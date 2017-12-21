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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/logger"
)

const keyDirEnv = "__AWLESS_KEYS_DIR"

type CreateKeypair struct {
	_                 string `action:"create" entity:"keypair" awsAPI:"ec2" awsCall:"ImportKeyPair" awsInput:"ec2.ImportKeyPairInput" awsOutput:"ec2.ImportKeyPairOutput"`
	logger            *logger.Logger
	graph             cloud.GraphAPI
	api               ec2iface.EC2API
	Name              *string `awsName:"KeyName" awsType:"awsstr" templateName:"name"`
	Encrypted         *bool   `templateName:"encrypted"`
	PublicKeyMaterial []byte  `awsName:"PublicKeyMaterial" awsType:"awsbyteslice"`
}

func (cmd *CreateKeypair) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("name"), params.Opt("encrypted")),
		params.Validators{
			"name": func(i interface{}, others map[string]interface{}) error {
				keyDir := os.Getenv(keyDirEnv)
				if keyDir == "" {
					return fmt.Errorf("empty env var '%s'", keyDirEnv)
				}

				privKeyPath := filepath.Join(keyDir, fmt.Sprint(i)+".pem")
				if _, err := os.Stat(privKeyPath); err == nil {
					return fmt.Errorf("file already exists at path: %s", privKeyPath)
				}
				return nil
			},
		})
}

func (cmd *CreateKeypair) BeforeRun(renv env.Running) error {
	var encryptedMsg string
	var encrypted bool

	if BoolValue(cmd.Encrypted) {
		encrypted = true
		encryptedMsg = " encrypted"
	}

	privKeyPath := filepath.Join(os.Getenv(keyDirEnv), StringValue(cmd.Name)+".pem")
	if _, err := os.Stat(privKeyPath); err == nil {
		return fmt.Errorf("saving private key: file already exists at path: %s", privKeyPath)
	}

	cmd.logger.Infof("Generating locally%s 4096 RSA at %s", encryptedMsg, privKeyPath)
	start := time.Now()
	pub, priv, err := console.GenerateSSHKeyPair(4096, encrypted)
	cmd.logger.ExtraVerbosef("4096 bits key generation took %s", time.Since(start))
	if err != nil {
		return fmt.Errorf("generating key: %s", err)
	}
	if err = ioutil.WriteFile(privKeyPath, priv, 0400); err != nil {
		return fmt.Errorf("saving private key: %s", err)
	}
	cmd.PublicKeyMaterial = pub
	return nil
}

func (cmd *CreateKeypair) ExtractResult(i interface{}) string {
	return StringValue(i.(*ec2.ImportKeyPairOutput).KeyName)
}

type DeleteKeypair struct {
	_      string `action:"delete" entity:"keypair" awsAPI:"ec2" awsCall:"DeleteKeyPair" awsInput:"ec2.DeleteKeyPairInput" awsOutput:"ec2.DeleteKeyPairOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Name   *string `awsName:"KeyName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteKeypair) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}
