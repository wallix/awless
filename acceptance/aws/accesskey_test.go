package awsat

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/aws/config"
	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestAccesskey(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		defer redirectStdErrToDevNull()()

		t.Run("no-prompt", func(t *testing.T) {

			Template("create accesskey user=jdoe no-prompt=true").
				Mock(&iamMock{
					CreateAccessKeyFunc: func(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
						return &iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{AccessKeyId: String("new-keypair-id")}}, nil
					},
				}).ExpectInput("CreateAccessKey", &iam.CreateAccessKeyInput{UserName: String("jdoe")}).
				ExpectCommandResult("new-keypair-id").ExpectCalls("CreateAccessKey").Run(t)
		})

		t.Run("save credentials", func(t *testing.T) {
			awsFolder, err := ioutil.TempDir("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(awsFolder)
			awsconfig.AWSHomeDir = func() string {
				return awsFolder
			}
			awsspec.AWSCredFilepath = filepath.Join(awsFolder, "credentials")
			Template("create accesskey user=jdoe save=true").
				Mock(&iamMock{
					CreateAccessKeyFunc: func(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
						return &iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{AccessKeyId: String("0123456EXAMPLE0123456EXAMPLE"), SecretAccessKey: String("MYSECRETKEYMYSECRETKEYMYSECRETKEYMYSECRETKEY")}}, nil
					},
				}).ExpectInput("CreateAccessKey", &iam.CreateAccessKeyInput{UserName: String("jdoe")}).
				ExpectCommandResult("0123456EXAMPLE0123456EXAMPLE").ExpectCalls("CreateAccessKey").Run(t)
			cred, err := ioutil.ReadFile(awsspec.AWSCredFilepath)
			if err != nil {
				t.Fatal(err)
			}
			expected := `
[jdoe]
aws_access_key_id = 0123456EXAMPLE0123456EXAMPLE
aws_secret_access_key = MYSECRETKEYMYSECRETKEYMYSECRETKEYMYSECRETKEY
`
			if got, want := string(cred), expected; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("with user", func(t *testing.T) {
			Template("delete accesskey id=ACCESSKEYID user=jdoe").
				Mock(&iamMock{
					DeleteAccessKeyFunc: func(param0 *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
						return nil, nil
					},
				}).ExpectInput("DeleteAccessKey", &iam.DeleteAccessKeyInput{
				UserName:    String("jdoe"),
				AccessKeyId: String("ACCESSKEYID"),
			}).ExpectCalls("DeleteAccessKey").Run(t)
		})
		t.Run("without user and id in local graph", func(t *testing.T) {
			g := graph.NewGraph()
			g.AddResource(resourcetest.AccessKey("ACCESSKEYID").Prop(properties.Username, "myusername").Build())
			g.AddResource(resourcetest.AccessKey("OTHERACCESSKEYID").Prop(properties.Username, "notthis").Build())
			Template("delete accesskey id=ACCESSKEYID").
				Mock(&iamMock{
					DeleteAccessKeyFunc: func(param0 *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
						return nil, nil
					},
				}).Graph(g).ExpectInput("DeleteAccessKey", &iam.DeleteAccessKeyInput{
				UserName:    String("myusername"),
				AccessKeyId: String("ACCESSKEYID"),
			}).ExpectCalls("DeleteAccessKey").Run(t)
		})
		t.Run("without user and id not in local graph", func(t *testing.T) {
			Template("delete accesskey id=ACCESSKEYID").
				Mock(&iamMock{
					DeleteAccessKeyFunc: func(param0 *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
						return nil, nil
					},
				}).ExpectInput("DeleteAccessKey", &iam.DeleteAccessKeyInput{
				AccessKeyId: String("ACCESSKEYID"),
			}).ExpectCalls("DeleteAccessKey").Run(t)
		})
		t.Run("with user without id", func(t *testing.T) {
			Template("delete accesskey user=myusername").
				Mock(&iamMock{
					DeleteAccessKeyFunc: func(param0 *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
						return nil, nil
					},
					ListAccessKeysFunc: func(param0 *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
						return &iam.ListAccessKeysOutput{AccessKeyMetadata: []*iam.AccessKeyMetadata{{AccessKeyId: String("ACCESSKEYID")}}}, nil
					},
				}).ExpectInput("DeleteAccessKey", &iam.DeleteAccessKeyInput{
				UserName:    String("myusername"),
				AccessKeyId: String("ACCESSKEYID"),
			}).ExpectInput("ListAccessKeys", &iam.ListAccessKeysInput{
				UserName: String("myusername"),
			}).ExpectCalls("ListAccessKeys", "DeleteAccessKey").Run(t)
		})
	})
}

func redirectStdErrToDevNull() func() {
	originalStdErr := os.Stderr
	toDefer := func() {
		os.Stderr = originalStdErr
	}
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	os.Stderr = devNull
	return toDefer
}
