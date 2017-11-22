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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/logger"
)

type CreateAccesskey struct {
	_        string `action:"create" entity:"accesskey" awsAPI:"iam" awsCall:"CreateAccessKey" awsInput:"iam.CreateAccessKeyInput" awsOutput:"iam.CreateAccessKeyOutput"`
	logger   *logger.Logger
	api      iamiface.IAMAPI
	User     *string `awsName:"UserName" awsType:"awsstr" templateName:"user" required:""`
	NoPrompt *bool   `templateName:"no-prompt"`
}

func (cmd *CreateAccesskey) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateAccesskey) AfterRun(ctx map[string]interface{}, output interface{}) error {
	accessKey := output.(*iam.CreateAccessKeyOutput).AccessKey
	cmd.logger.Infof("Access key created. Here are the crendentials for user %s:", aws.StringValue(accessKey.UserName))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, strings.Repeat("*", 64))
	fmt.Fprintf(os.Stderr, "aws_access_key_id = %s\n", aws.StringValue(accessKey.AccessKeyId))
	fmt.Fprintf(os.Stderr, "aws_secret_access_key = %s\n", aws.StringValue(accessKey.SecretAccessKey))
	fmt.Fprintln(os.Stderr, strings.Repeat("*", 64))
	fmt.Fprintln(os.Stderr)
	cmd.logger.Warning("This is your only opportunity to view the secret access keys.")
	cmd.logger.Warning("Save the user's new access key ID and secret access key in a safe and secure place.")
	cmd.logger.Warning("You will not have access to the secret keys again after this step.\n")

	if !BoolValue(cmd.NoPrompt) {
		if promptConfirm("Do you want to save these access keys in %s?", AWSCredFilepath) {
			var profile string
			fmt.Print("Entry profile name (will default to AWS 'default') ? ")
			fmt.Scanln(&profile)
			profile = strings.TrimSpace(profile)
			if profile == "" {
				profile = "default"
			}
			creds := NewCredsPrompter(profile)
			creds.Val.AccessKeyID = aws.StringValue(accessKey.AccessKeyId)
			creds.Val.SecretAccessKey = aws.StringValue(accessKey.SecretAccessKey)
			created, err := creds.Store()
			if err != nil {
				logger.Errorf("cannot store access keys: %s", err)
			} else {
				if created {
					fmt.Fprintf(os.Stderr, "\n\u2713 %s created", AWSCredFilepath)
				}
				fmt.Fprintf(os.Stderr, "\n\u2713 Credentials for profile '%s' stored successfully in %s\n\n", creds.Profile, AWSCredFilepath)
			}
		}
	}

	return nil
}

func (cmd *CreateAccesskey) ExtractResult(i interface{}) string {
	return StringValue(i.(*iam.CreateAccessKeyOutput).AccessKey.AccessKeyId)
}

type DeleteAccesskey struct {
	_      string `action:"delete" entity:"accesskey" awsAPI:"iam" awsCall:"DeleteAccessKey" awsInput:"iam.DeleteAccessKeyInput" awsOutput:"iam.DeleteAccessKeyOutput"`
	logger *logger.Logger
	api    iamiface.IAMAPI
	Id     *string `awsName:"AccessKeyId" awsType:"awsstr" templateName:"id" required:""`
	User   *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
}

func (cmd *DeleteAccesskey) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

var (
	AWSCredFilepath = filepath.Join(awsconfig.AWSHomeDir(), "credentials")
)

type credentialsPrompter struct {
	Profile               string
	Val                   credentials.Value
	ProfileSetterCallback func(string) error
}

func NewCredsPrompter(profile string) *credentialsPrompter {
	return &credentialsPrompter{Profile: profile, ProfileSetterCallback: func(string) error { return nil }}
}

func (c *credentialsPrompter) Prompt() error {
	token := "and choose a profile name"
	if c.HasProfile() {
		token = fmt.Sprintf("for profile '%s'", c.Profile)
	}
	fmt.Printf("\nPlease enter access keys %s (stored at %s):\n", token, AWSCredFilepath)

	promptUntilNonEmpty("AWS Access Key ID? ", &c.Val.AccessKeyID)
	promptUntilNonEmpty("AWS Secret Access Key? ", &c.Val.SecretAccessKey)
	if c.HasProfile() {
		promptToOverride(fmt.Sprintf("Change your profile name (or just press Enter to keep '%s')? ", c.Profile), &c.Profile)
	} else {
		c.Profile = "default"
		promptToOverride("Choose a profile name (or just press Enter to have AWS 'default')? ", &c.Profile)
	}

	if c.ProfileSetterCallback != nil {
		c.ProfileSetterCallback(c.Profile)
	}

	return nil
}

func (c *credentialsPrompter) Store() (bool, error) {
	var created bool

	if c.Val.SecretAccessKey == "" {
		return created, errors.New("given empty secret access key")
	}
	if c.Val.AccessKeyID == "" {
		return created, errors.New("given empty access key")
	}
	return appendToAwsFile(
		fmt.Sprintf("\n[%s]\naws_access_key_id = %s\naws_secret_access_key = %s\n", c.Profile, c.Val.AccessKeyID, c.Val.SecretAccessKey),
		AWSCredFilepath,
	)
}

func appendToAwsFile(content string, awsFilePath string) (bool, error) {
	var created bool
	if awsHomeDirMissing() {
		if err := os.MkdirAll(awsconfig.AWSHomeDir(), 0700); err != nil {
			return created, fmt.Errorf("creating '%s' : %s", awsconfig.AWSHomeDir(), err)
		}
		created = true
	}

	f, err := os.OpenFile(awsFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return created, fmt.Errorf("appending to '%s': %s", awsFilePath, err)
	}

	if _, err := fmt.Fprintf(f, content); err != nil {
		return created, err
	}

	return created, nil
}

func promptConfirm(msg string, a ...interface{}) bool {
	var yesorno string
	fmt.Fprintf(os.Stderr, "%s [y/N] ", fmt.Sprintf(msg, a...))
	fmt.Scanln(&yesorno)
	if y := strings.TrimSpace(strings.ToLower(yesorno)); y == "y" || y == "yes" {
		return true
	}
	return false
}

func (c *credentialsPrompter) HasProfile() bool {
	return strings.TrimSpace(c.Profile) != ""
}

func promptToOverride(question string, v *string) {
	fmt.Print(question)
	var override string
	fmt.Scanln(&override)
	if strings.TrimSpace(override) != "" {
		*v = override
		return
	}
}

func promptUntilNonEmpty(question string, v *string) {
	ask := func(v *string) bool {
		fmt.Print(question)
		_, err := fmt.Scanln(v)
		if err == nil && strings.TrimSpace(*v) != "" {
			return false
		}
		if err != nil {
			fmt.Printf("Error: %s. Retry please...\n", err)
		}
		return true
	}
	for ask(v) {
	}
	return
}

func awsHomeDirMissing() bool {
	_, err := os.Stat(awsconfig.AWSHomeDir())
	return os.IsNotExist(err)
}
