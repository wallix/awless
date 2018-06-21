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
	"regexp"
	"strings"

	"github.com/wallix/awless/cloud/match"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/logger"
)

type CreateAccesskey struct {
	_      string `action:"create" entity:"accesskey" awsAPI:"iam" awsCall:"CreateAccessKey" awsInput:"iam.CreateAccessKeyInput" awsOutput:"iam.CreateAccessKeyOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	User   *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
	Save   *bool   `templateName:"save"`
}

func (cmd *CreateAccesskey) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.AllOf(params.Key("user"),
		params.Opt("save", "no-prompt"),
	))
	builder.AddReducer(
		func(values map[string]interface{}) (map[string]interface{}, error) {
			if noPrompt, hasNoPrompt := values["no-prompt"]; hasNoPrompt {
				b, err := castBool(noPrompt)
				if err != nil {
					return nil, fmt.Errorf("no-prompt: %s", err)
				}
				return map[string]interface{}{"save": !b}, nil
			} else {
				return nil, nil
			}
		},
		"no-prompt",
	)
	return builder.Done()
}

func (cmd *CreateAccesskey) AfterRun(renv env.Running, output interface{}) error {
	accessKey := output.(*iam.CreateAccessKeyOutput).AccessKey
	if !BoolValue(cmd.Save) {
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
	}

	if cmd.Save != nil && !BoolValue(cmd.Save) {
		return nil
	}
	profile := StringValue(cmd.User)
	if !BoolValue(cmd.Save) {
		if !promptConfirm("Do you want to save these access keys in %s?", AWSCredFilepath) {
			return nil
		}
		profile = promptStringWithDefault("Entry profile name: ("+StringValue(cmd.User)+") ", profile)
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

	return nil
}

func (cmd *CreateAccesskey) ExtractResult(i interface{}) string {
	return StringValue(i.(*iam.CreateAccessKeyOutput).AccessKey.AccessKeyId)
}

type DeleteAccesskey struct {
	_      string `action:"delete" entity:"accesskey" awsAPI:"iam" awsCall:"DeleteAccessKey" awsInput:"iam.DeleteAccessKeyInput" awsOutput:"iam.DeleteAccessKeyOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Id     *string `awsName:"AccessKeyId" awsType:"awsstr" templateName:"id"`
	User   *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
}

func (cmd *DeleteAccesskey) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.AtLeastOneOf(params.Key("id"), params.Key("user")))
	builder.AddReducer(
		func(values map[string]interface{}) (map[string]interface{}, error) {
			user, hasUser := values["user"].(string)
			id, hasId := values["id"].(string)
			if !hasUser && hasId {
				r, err := cmd.graph.FindOne(cloud.NewQuery(cloud.AccessKey).Match(match.Property(properties.ID, id)))
				if err != nil || r == nil {
					return values, nil
				}
				if keyUser, ok := r.Property(properties.Username); ok {
					values["user"] = keyUser
				}
			} else if hasUser && !hasId {
				keys, err := cmd.api.ListAccessKeys(&iam.ListAccessKeysInput{
					UserName: String(user),
				})
				if err != nil {
					return values, fmt.Errorf("can not find access key for %s: %s", user, err)
				}
				switch len(keys.AccessKeyMetadata) {
				case 0:
					return values, fmt.Errorf("no access key found for %s:", user)
				case 1:
					values["id"] = StringValue(keys.AccessKeyMetadata[0].AccessKeyId)
				default:
					var keysStr []string
					for _, k := range keys.AccessKeyMetadata {
						keysStr = append(keysStr, fmt.Sprintf("%s (created on %s)", StringValue(k.AccessKeyId), aws.TimeValue(k.CreateDate).Format("2006/01/02 15:04:05")))
					}
					return values, fmt.Errorf("multiple access keys found for %s: %s", user, strings.Join(keysStr, ", "))
				}

			}
			return values, nil
		},
		"user", "id",
	)
	return builder.Done()
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

var accessKeysRegex = regexp.MustCompile("^[a-zA-Z0-9/+=]{20,60}$")

func (c *credentialsPrompter) Store() (bool, error) {
	var created bool

	if c.Val.SecretAccessKey == "" {
		return created, errors.New("given empty secret access key")
	}
	if !accessKeysRegex.MatchString(c.Val.SecretAccessKey) {
		return created, errors.New("given invalid secret access key")
	}
	if c.Val.AccessKeyID == "" {
		return created, errors.New("given empty access key")
	}
	if !accessKeysRegex.MatchString(c.Val.AccessKeyID) {
		return created, errors.New("given invalid access key")
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
