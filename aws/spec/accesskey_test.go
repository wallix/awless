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
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wallix/awless/aws/config"
)

func TestCredentialsPrompter(t *testing.T) {
	tmpAWSDir, err := ioutil.TempDir("", "testAWS")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpAWSDir)
	awsconfig.AWSHomeDir = func() string { return tmpAWSDir }
	AWSCredFilepath = filepath.Join(awsconfig.AWSHomeDir(), "credentials")

	tcases := []struct {
		accessKey, secret string
		expectErr         string
	}{
		{accessKey: "AKIAIOSFODNN7EXAMPLE", secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectErr: ""},
		{accessKey: "", secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectErr: "empty access key"},
		{accessKey: "AKIAIOSFODNN7EXAMPLE", secret: "", expectErr: "empty secret access key"},
		{accessKey: "invalid", secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectErr: "invalid access key"},
		{accessKey: "AKIAIOSFODNN7EXAMPLE", secret: "invalid", expectErr: "invalid secret access key"},
		{accessKey: "$fordiden-chars!", secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", expectErr: "invalid access key"},
		{accessKey: "AKIAIOSFODNN7EXAMPLE", secret: "faaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaarTooooooooooooLooooooooooooonnnnnnnng", expectErr: "invalid secret access key"},
	}

	for i, tcase := range tcases {
		prompter := NewCredsPrompter("test-profile")
		prompter.Val.AccessKeyID = tcase.accessKey
		prompter.Val.SecretAccessKey = tcase.secret

		created, err := prompter.Store()
		if tcase.expectErr == "" {
			if err != nil {
				t.Fatalf("%d: expect no error, got %s", i+1, err)
			}
		} else {
			if err == nil {
				t.Fatalf("%d: expect error %s, got nil", i+1, tcase.expectErr)
			} else {
				if got, want := err.Error(), tcase.expectErr; !strings.Contains(got, want) {
					t.Fatalf("%d: expect error contains %s, got %s", i+1, want, got)
				}
			}
			continue
		}
		if got, want := created, false; got != want {
			t.Fatalf("%d: got %t, want %t", i+1, got, want)
		}
		credentials, err := ioutil.ReadFile(AWSCredFilepath)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}
		expectCredentials := fmt.Sprintf(`
[test-profile]
aws_access_key_id = %s
aws_secret_access_key = %s
`, tcase.accessKey, tcase.secret)
		if got, want := string(credentials), expectCredentials; got != want {
			t.Fatalf("%d: got\n%q\n, want\n%q\n", i+1, got, want)
		}

		os.Remove(path.Join(tmpAWSDir, "credentials"))
	}

}
