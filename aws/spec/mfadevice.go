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
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/logger"
)

type CreateMfadevice struct {
	_      string `action:"create" entity:"mfadevice" awsAPI:"iam"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Name   *string `templateName:"name"`
}

func (cmd *CreateMfadevice) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

func (cmd *CreateMfadevice) ManualRun(renv env.Running) (interface{}, error) {
	name := StringValue(cmd.Name)
	input := &iam.CreateVirtualMFADeviceInput{
		VirtualMFADeviceName: cmd.Name,
	}
	var err error

	start := time.Now()
	var output *iam.CreateVirtualMFADeviceOutput
	output, err = cmd.api.CreateVirtualMFADevice(input)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	cmd.logger.ExtraVerbosef("iam.CreateVirtualMFADevice call took %s", time.Since(start))

	cmd.logger.Infof("MFA virtual device created. Here is the secret to create virtual device: %s.", string(output.VirtualMFADevice.Base32StringSeed))
	cmd.logger.Infof("You can also use this QRCode:")

	qrCodeURI := fmt.Sprintf("otpauth://totp/AWS:%s?secret=%s", name, string(output.VirtualMFADevice.Base32StringSeed))
	qrcode, err := qr.Encode(qrCodeURI, qr.L, qr.Auto)
	if err != nil {
		return nil, fmt.Errorf("encode qrcode: %s", err)
	}
	qrCodeDisplaySize := 40
	qrcode, err = barcode.Scale(qrcode, qrCodeDisplaySize, qrCodeDisplaySize)
	if err != nil {
		return nil, fmt.Errorf("scale qrcode: %s", err)
	}
	displayQRCode(os.Stderr, qrcode)
	cmd.logger.Warning("This is your only opportunity to view the secret. You will not have access to the secret again after this step.\n")
	return output, nil
}

func (cmd *CreateMfadevice) ExtractResult(i interface{}) string {
	return StringValue(i.(*iam.CreateVirtualMFADeviceOutput).VirtualMFADevice.SerialNumber)
}

type DeleteMfadevice struct {
	_      string `action:"delete" entity:"mfadevice" awsAPI:"iam" awsCall:"DeleteVirtualMFADevice" awsInput:"iam.DeleteVirtualMFADeviceInput" awsOutput:"iam.DeleteVirtualMFADeviceOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Id     *string `awsName:"SerialNumber" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteMfadevice) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

var (
	awsConfigFilepath = filepath.Join(awsconfig.AWSHomeDir(), "config")
)

type AttachMfadevice struct {
	_        string `action:"attach" entity:"mfadevice" awsAPI:"iam" awsCall:"EnableMFADevice" awsInput:"iam.EnableMFADeviceInput" awsOutput:"iam.EnableMFADeviceOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      iamiface.IAMAPI
	Id       *string `awsName:"SerialNumber" awsType:"awsstr" templateName:"id"`
	User     *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
	MfaCode1 *string `awsName:"AuthenticationCode1" awsType:"aws6digitsstring" templateName:"mfa-code-1"`
	MfaCode2 *string `awsName:"AuthenticationCode2" awsType:"aws6digitsstring" templateName:"mfa-code-2"`
	NoPrompt *bool   `templateName:"no-prompt"`
}

func (cmd *AttachMfadevice) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("mfa-code-1"), params.Key("mfa-code-2"), params.Key("user"),
		params.Opt("no-prompt"),
	))
}

func (cmd *AttachMfadevice) AfterRun(renv env.Running, output interface{}) error {
	if !BoolValue(cmd.NoPrompt) {
		if promptConfirm("\nDo you want to create a profile for this MFA device in %s?", awsConfigFilepath) {
			roleArn, err := promptRole(cmd.api)
			for err != nil {
				if !promptConfirm("\nDo you want to create a profile for this MFA device in %s?", awsConfigFilepath) {
					return nil
				}
				roleArn, err = promptRole(cmd.api)
			}
			fmt.Fprintln(os.Stderr)
			srcProfile := promptStringWithDefault("Enter source profile used to assume role: (default) ", "default")

			mfaProfile := promptStringWithDefault("Enter new MFA profile name: (mfa) ", "mfa")

			config := fmt.Sprintf("\n[%s]\n"+
				"source_profile = %s\n"+
				"mfa_serial = %s\n"+
				"role_arn = %s", mfaProfile, srcProfile, StringValue(cmd.Id), roleArn)
			if promptConfirm("\n%s\n\nAppend this to '%s'?", config, awsConfigFilepath) {
				created, err := appendToAwsFile(config, awsConfigFilepath)
				if err != nil {
					cmd.logger.Error(err)
				} else {
					if created {
						fmt.Fprintf(os.Stderr, "\n\u2713 %s created", awsConfigFilepath)
					}
					fmt.Fprintf(os.Stderr, "\n\u2713 New profile '%s' for MFA device stored successfully in '%s'\n\n", mfaProfile, awsConfigFilepath)
					return nil
				}
			}
			fmt.Fprintf(os.Stderr, "Canceled modification of '%s'.\n\n", awsConfigFilepath)
			return nil
		}
		fmt.Fprintf(os.Stderr, "Canceled adding profile for MFA device to '%s'.\n\n", awsConfigFilepath)
		return nil
	}
	return nil
}

type DetachMfadevice struct {
	_      string `action:"detach" entity:"mfadevice" awsAPI:"iam" awsCall:"DeactivateMFADevice" awsInput:"iam.DeactivateMFADeviceInput" awsOutput:"iam.DeactivateMFADeviceOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    iamiface.IAMAPI
	Id     *string `awsName:"SerialNumber" awsType:"awsstr" templateName:"id"`
	User   *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
}

func (cmd *DetachMfadevice) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("user")))
}

func displayQRCode(w io.Writer, qrCode barcode.Barcode) {
	white := color.New(color.BgWhite)
	black := color.New(color.BgBlack)
	for x := 0; x < qrCode.Bounds().Dx(); x++ {
		for y := 0; y < qrCode.Bounds().Dy(); y++ {
			r32, g32, b32, _ := qrCode.At(x, y).RGBA()
			r, g, b := int(r32>>8), int(g32>>8), int(b32>>8)
			if (r+g+b)/3 > 180 {
				white.Fprint(w, "  ")
			} else {
				black.Fprint(w, "  ")
			}
		}
		fmt.Fprintln(w)
	}
}

func promptStringWithDefault(msg, def string) (res string) {
	fmt.Fprintf(os.Stderr, msg)
	fmt.Scanln(&res)
	res = strings.TrimSpace(res)
	if res == "" {
		res = def
	}
	return
}

func promptRole(api iamiface.IAMAPI) (string, error) {
	rolesNameToArn := make(map[string]string)

	err := api.ListRolesPages(&iam.ListRolesInput{}, func(out *iam.ListRolesOutput, lastPage bool) bool {
		for _, role := range out.Roles {
			rolesNameToArn[StringValue(role.RoleName)] = StringValue(role.Arn)
		}
		return out.Marker != nil
	})

	if err == nil && len(rolesNameToArn) > 0 {
		var roles []readline.PrefixCompleterInterface
		for name, arn := range rolesNameToArn {
			roles = append(roles, readline.PcItem(name))
			roles = append(roles, readline.PcItem(arn))
		}
		var roleCompleter = readline.NewPrefixCompleter(roles...)
		fmt.Fprint(os.Stderr, "Please specify the role (name or ARN) to assume with this MFA device: (Tab for completion) \n")
		rl, err := readline.NewEx(&readline.Config{
			Prompt:       "> ",
			AutoComplete: roleCompleter,
		})
		if err != nil {
			return "", fmt.Errorf("error while selecting role: %s", err)
		}
		defer rl.Close()
		role, err := rl.Readline()
		if err != nil {
			return "", fmt.Errorf("error while selecting role: %s", err)
		}
		if arn, isName := rolesNameToArn[strings.TrimSpace(role)]; isName {
			return arn, nil
		}
		return role, nil
	}
	//No permission to list roles:
	var roleArn string
	fmt.Fprint(os.Stderr, "Please specify the role ARN to assume with this MFA device:")
	fmt.Scanln(&roleArn)
	roleArn = strings.TrimSpace(roleArn)
	if roleArn == "" {
		return roleArn, errors.New("Role cannot be empty")
	}
	return roleArn, nil
}
