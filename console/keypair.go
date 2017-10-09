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
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var askPasswordFunc func() ([]byte, error) = func() ([]byte, error) {
	fmt.Fprint(os.Stderr, "This SSH key will be encrypted. Please enter new password:")
	for {
		pass, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return pass, err
		}
		fmt.Fprint(os.Stderr, "\nConfirm password:")
		pass2, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return pass, err
		}
		if !bytes.Equal(pass, pass2) {
			fmt.Fprint(os.Stderr, "\nPasswords are different. Please enter new password:")
			continue
		}
		fmt.Fprintln(os.Stderr)
		return pass, nil
	}
}

var GenerateSSHKeyPair = func(size int, encryptKey bool) ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, nil, err
	}
	var pemBlock *pem.Block
	var passwd []byte
	if encryptKey {
		passwd, err = askPasswordFunc()
		if err != nil {
			return nil, nil, err
		}
		if len(passwd) == 0 {
			fmt.Fprintln(os.Stderr, "Empty password given, the keypair will not be encrypted.")
		}
	}
	if len(passwd) != 0 {
		pemBlock, err = x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key), passwd, x509.PEMCipherAES256)
		if err != nil {
			return nil, nil, err
		}
	} else {
		pemBlock = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	}

	privPem := pem.EncodeToMemory(pemBlock)

	sshPub, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return ssh.MarshalAuthorizedKey(sshPub), privPem, nil
}
