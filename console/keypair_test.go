/*
Copyright 2017 WALLIX

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

package console

import (
	"strings"
	"testing"

	"github.com/wallix/awless/ssh"
)

func TestGenerateKeyPair(t *testing.T) {
	size := 1024
	pub, private, err := GenerateSSHKeyPair(size, false)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := strings.HasPrefix(string(pub), "ssh-rsa"), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := strings.HasPrefix(string(private), "-----BEGIN RSA PRIVATE KEY-----"), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := strings.HasSuffix(string(private), "-----END RSA PRIVATE KEY-----\n"), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	askPasswordFunc = func() ([]byte, error) {
		return []byte("my$rongP4sswrd"), nil
	}

	pub, private, err = GenerateSSHKeyPair(size, true)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := strings.HasPrefix(string(pub), "ssh-rsa"), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	_, err = ssh.DecryptSSHKey(private, []byte("my$rongP4sswrd"))
	if err != nil {
		t.Fatal(err)
	}
}
