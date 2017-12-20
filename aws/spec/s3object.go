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
	"mime"
	"os"
	"path/filepath"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/mitchellh/ioprogress"
	"github.com/wallix/awless/logger"
)

type CreateS3object struct {
	_      string `action:"create" entity:"s3object" awsAPI:"s3"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    s3iface.S3API
	Bucket *string `awsName:"Bucket" awsType:"awsstr" templateName:"bucket"`
	File   *string `awsName:"Body" awsType:"awsstr" templateName:"file"`
	Name   *string `awsName:"Key" awsType:"awsstr" templateName:"name"`
	Acl    *string `awsName:"ACL" awsType:"awsstr" templateName:"acl"`
}

func (cmd *CreateS3object) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("bucket"), params.Key("file"), params.Opt("acl", "name")),
		params.Validators{"file": params.IsFilepath},
	)
}

func (cmd *CreateS3object) ManualRun(env.Running) (interface{}, error) {
	input := &s3.PutObjectInput{}

	f, err := os.Open(StringValue(cmd.File))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	progressR, err := ProgressBarFactory(f)
	if err != nil {
		return nil, err
	}
	input.Body = progressR

	var fileName string
	if n := StringValue(cmd.Name); n != "" {
		fileName = n
	} else {
		_, fileName = filepath.Split(f.Name())
	}
	input.Key = aws.String(fileName)

	fileExt := filepath.Ext(f.Name())
	if mimeType := mime.TypeByExtension(fileExt); mimeType != "" {
		cmd.logger.ExtraVerbosef("setting object content-type to '%s'", mimeType)
		input.ContentType = aws.String(mimeType)
	}

	if err = setFieldWithType(cmd.Bucket, input, "Bucket", awsstr); err != nil {
		return nil, err
	}

	if v := cmd.Acl; v != nil {
		if err = setFieldWithType(v, input, "ACL", awsstr); err != nil {
			return nil, err
		}
	}

	cmd.logger.Infof("uploading '%s'", fileName)

	if _, err = cmd.api.PutObject(input); err != nil {
		return nil, err
	}

	return fileName, nil
}

func (cmd *CreateS3object) ExtractResult(i interface{}) string {
	return i.(string)
}

type UpdateS3object struct {
	_       string `action:"update" entity:"s3object" awsAPI:"s3" awsCall:"PutObjectAcl" awsInput:"s3.PutObjectAclInput" awsOutput:"s3.PutObjectAclOutput"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     s3iface.S3API
	Bucket  *string `awsName:"Bucket" awsType:"awsstr" templateName:"bucket"`
	Name    *string `awsName:"Key" awsType:"awsstr" templateName:"name"`
	Acl     *string `awsName:"ACL" awsType:"awsstr" templateName:"acl"`
	Version *string `awsName:"VersionId" awsType:"awsstr" templateName:"version"`
}

func (cmd *UpdateS3object) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("acl"), params.Key("bucket"), params.Key("name"),
		params.Opt("version"),
	))
}

type DeleteS3object struct {
	_      string `action:"delete" entity:"s3object" awsAPI:"s3" awsCall:"DeleteObject" awsInput:"s3.DeleteObjectInput" awsOutput:"s3.DeleteObjectOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    s3iface.S3API
	Bucket *string `awsName:"Bucket" awsType:"awsstr" templateName:"bucket"`
	Name   *string `awsName:"Key" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteS3object) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("bucket"), params.Key("name")))
}

type ProgressReadSeeker struct {
	file   *os.File
	reader *ioprogress.Reader
}

func NewProgressReader(f *os.File) (*ProgressReadSeeker, error) {
	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	draw := func(progress, total int64) string {
		// &s3.PutObjectInput.Body will be read twice
		// once in memory and a second time for the HTTP upload
		// here we only display for the actual HTTP upload
		if progress > total {
			return ioprogress.DrawTextFormatBytes(progress/2, total)
		}
		return ""
	}

	reader := &ioprogress.Reader{
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, draw),
		Reader:   f,
		Size:     finfo.Size(),
	}

	return &ProgressReadSeeker{file: f, reader: reader}, nil
}

func (pr *ProgressReadSeeker) Read(p []byte) (int, error) {
	return pr.reader.Read(p)
}

func (pr *ProgressReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return pr.file.Seek(offset, whence)
}

// Allow to control for testing
var ProgressBarFactory func(*os.File) (*ProgressReadSeeker, error)

func init() {
	ProgressBarFactory = NewProgressReader
}
