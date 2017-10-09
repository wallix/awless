package awsat

import (
	"fmt"
	"testing"

	"io/ioutil"

	"os"

	"github.com/aws/aws-sdk-go/service/lambda"
)

func TestFunction(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Run("from s3 file", func(t *testing.T) {
			Template("create function name=my-function-name handler=lambda_handler role=arn:of:function:role runtime=python3.6 "+
				"bucket=my-function-bucket object=my/function/object.zip objectversion=v3 description='this is my function' memory=128 publish=true timeout=60").
				Mock(&lambdaMock{
					CreateFunctionFunc: func(param0 *lambda.CreateFunctionInput) (*lambda.FunctionConfiguration, error) {
						return &lambda.FunctionConfiguration{FunctionArn: String("new-function-id")}, nil
					},
				}).ExpectInput("CreateFunction", &lambda.CreateFunctionInput{
				FunctionName: String("my-function-name"),
				Handler:      String("lambda_handler"),
				Role:         String("arn:of:function:role"),
				Runtime:      String("python3.6"),
				Code: &lambda.FunctionCode{
					S3Bucket:        String("my-function-bucket"),
					S3Key:           String("my/function/object.zip"),
					S3ObjectVersion: String("v3"),
				},
				Description: String("this is my function"),
				MemorySize:  Int64(128),
				Publish:     Bool(true),
				Timeout:     Int64(60),
			}).ExpectCommandResult("new-function-id").ExpectCalls("CreateFunction").Run(t)
		})
		t.Run("from zip file", func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				os.Remove(tmpFile.Name())
			}()
			ioutil.WriteFile(tmpFile.Name(), []byte("this is the content of my file"), 0777)
			Template(fmt.Sprintf("create function name=my-function-name handler=lambda_handler role=arn:of:function:role runtime=python3.6 "+
				"zipfile=%s", tmpFile.Name())).
				Mock(&lambdaMock{
					CreateFunctionFunc: func(param0 *lambda.CreateFunctionInput) (*lambda.FunctionConfiguration, error) {
						return &lambda.FunctionConfiguration{FunctionArn: String("new-function-id")}, nil
					},
				}).ExpectInput("CreateFunction", &lambda.CreateFunctionInput{
				FunctionName: String("my-function-name"),
				Handler:      String("lambda_handler"),
				Role:         String("arn:of:function:role"),
				Runtime:      String("python3.6"),
				Code: &lambda.FunctionCode{
					ZipFile: []byte("this is the content of my file"),
				},
			}).ExpectCommandResult("new-function-id").ExpectCalls("CreateFunction").Run(t)
		})
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete function id=function-to-delete version=v2").
			Mock(&lambdaMock{
				DeleteFunctionFunc: func(param0 *lambda.DeleteFunctionInput) (*lambda.DeleteFunctionOutput, error) { return nil, nil },
			}).ExpectInput("DeleteFunction", &lambda.DeleteFunctionInput{
			FunctionName: String("function-to-delete"),
			Qualifier:    String("v2"),
		}).
			ExpectCalls("DeleteFunction").Run(t)
	})
}
