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
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/fatih/color"
	"github.com/wallix/awless/logger"
)

type CreateMfadevice struct {
	_      string `action:"create" entity:"mfadevice" awsAPI:"iam"`
	logger *logger.Logger
	api    iamiface.IAMAPI
	Name   *string `templateName:"name" required:""`
}

func (cmd *CreateMfadevice) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateMfadevice) ManualRun(ctx map[string]interface{}) (interface{}, error) {
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
	api    iamiface.IAMAPI
	Id     *string `awsName:"SerialNumber" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *DeleteMfadevice) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type AttachMfadevice struct {
	_        string `action:"attach" entity:"mfadevice" awsAPI:"iam" awsCall:"EnableMFADevice" awsInput:"iam.EnableMFADeviceInput" awsOutput:"iam.EnableMFADeviceOutput"`
	logger   *logger.Logger
	api      iamiface.IAMAPI
	Id       *string `awsName:"SerialNumber" awsType:"awsstr" templateName:"id" required:""`
	User     *string `awsName:"UserName" awsType:"awsstr" templateName:"user" required:""`
	MfaCode1 *string `awsName:"AuthenticationCode1" awsType:"aws6digitsstring" templateName:"mfa-code-1" required:""`
	MfaCode2 *string `awsName:"AuthenticationCode2" awsType:"aws6digitsstring" templateName:"mfa-code-2" required:""`
}

func (cmd *AttachMfadevice) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type DetachMfadevice struct {
	_      string `action:"detach" entity:"mfadevice" awsAPI:"iam" awsCall:"DeactivateMFADevice" awsInput:"iam.DeactivateMFADeviceInput" awsOutput:"iam.DeactivateMFADeviceOutput"`
	logger *logger.Logger
	api    iamiface.IAMAPI
	Id     *string `awsName:"SerialNumber" awsType:"awsstr" templateName:"id" required:""`
	User   *string `awsName:"UserName" awsType:"awsstr" templateName:"user" required:""`
}

func (cmd *DetachMfadevice) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
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
