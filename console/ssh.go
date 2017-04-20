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
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh/terminal"
)

const awlessKnownHostsFile = "known_hosts"

type Credentials struct {
	IP      string
	User    string
	KeyPath string
}

func NewSSHClient(cred *Credentials, disableStrictHostKeyChecking bool) (*ssh.Client, error) {
	privateKey, err := ioutil.ReadFile(cred.KeyPath)
	if os.IsNotExist(err) {
		privateKey, err = ioutil.ReadFile(cred.KeyPath + ".pem")
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Cannot find SSH key at '%s'. You can add `-i ./path/to/key`", cred.KeyPath)
		}
		cred.KeyPath = cred.KeyPath + ".pem"
	}
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil && strings.Contains(err.Error(), "cannot decode encrypted private keys") {
		fmt.Fprintf(os.Stderr, "This SSH key is encrypted. Please enter passphrase for key '%s':", cred.KeyPath)
		var passphrase []byte
		passphrase, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(os.Stderr)
		signer, err = decryptSSHKey(privateKey, passphrase)
	}
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: cred.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout:         2 * time.Second,
		HostKeyCallback: checkHostKey,
	}
	if disableStrictHostKeyChecking {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	return ssh.Dial("tcp", cred.IP+":22", config)
}

func InteractiveTerminal(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, stderr)

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	width := GetTerminalWidth()
	if width == 0 {
		width = 100
	}
	height := GetTerminalHeight()
	if height == 0 {
		height = 100
	}
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		return err
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		return err
	}

	signalc := make(chan os.Signal)
	defer func() {
		signal.Reset()
		close(signalc)
	}()
	go propagateSignals(signalc, session, stdin)
	signal.Notify(signalc, os.Interrupt, os.Kill)
	return session.Wait()
}

var trustKeyFunc func(hostname string, remote net.Addr, key ssh.PublicKey, keyFileName string) bool = func(hostname string, remote net.Addr, key ssh.PublicKey, keyFileName string) bool {
	fmt.Printf("awless could not validate the authenticity of '%s' (unknown host)\n", hostname)
	fmt.Printf("%s public key fingerprint is %s.\n", key.Type(), ssh.FingerprintSHA256(key))
	fmt.Printf("Do you want to continue connecting and persist this key to '%s' (yes/no)? ", keyFileName)
	var yesorno string
	_, err := fmt.Scanln(&yesorno)
	if err != nil {
		return false
	}
	return strings.ToLower(yesorno) == "yes"
}

func checkHostKey(hostname string, remote net.Addr, key ssh.PublicKey) error {
	var knownHostsFiles []string
	var fileToAddKnownKey string

	opensshFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	if _, err := os.Stat(opensshFile); err == nil {
		knownHostsFiles = append(knownHostsFiles, opensshFile)
		fileToAddKnownKey = opensshFile
	}

	awlessFile := filepath.Join(os.Getenv("__AWLESS_HOME"), "known_hosts")
	if _, err := os.Stat(awlessFile); err == nil {
		knownHostsFiles = append(knownHostsFiles, awlessFile)
	}
	if fileToAddKnownKey == "" {
		fileToAddKnownKey = awlessFile
	}

	checkKnownHostFunc, err := knownhosts.New(knownHostsFiles...)
	if err != nil {
		return err
	}
	knownhostsErr := checkKnownHostFunc(hostname, remote, key)
	keyError, ok := knownhostsErr.(*knownhosts.KeyError)
	if !ok {
		return knownhostsErr
	}
	if len(keyError.Want) == 0 {
		if trustKeyFunc(hostname, remote, key, fileToAddKnownKey) {
			f, err := os.OpenFile(fileToAddKnownKey, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = f.WriteString(knownhosts.Line([]string{hostname}, key))
			return err
		} else {
			return errors.New("Host public key verification failed.")
		}
	}

	var knownKeyInfos string
	var knownKeyFiles []string
	for _, knownKey := range keyError.Want {
		knownKeyInfos += fmt.Sprintf("\n-> %s (%s key in %s:%d)", ssh.FingerprintSHA256(knownKey.Key), knownKey.Key.Type(), knownKey.Filename, knownKey.Line)
		knownKeyFiles = append(knownKeyFiles, fmt.Sprintf("'%s:%d'", knownKey.Filename, knownKey.Line))
	}

	return fmt.Errorf(`
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
AWLESS DETECTED THAT THE REMOTE HOST PUBLIC KEY HAS CHANGED
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

Host key for '%s' has changed and you did not disable strict host key checking.
Someone may be trying to intercept your connection (man-in-the-middle attack). Otherwise, the host key may have been changed.

The fingerprint for the %s key sent by the remote host is %s.
You persisted:%s

To get rid of this message, update %s`, hostname, key.Type(), ssh.FingerprintSHA256(key), knownKeyInfos, strings.Join(knownKeyFiles, ","))
}

func propagateSignals(signalc chan os.Signal, session *ssh.Session, stdin io.WriteCloser) {
	for s := range signalc {
		switch s {
		case os.Interrupt:
			fmt.Fprint(stdin, "\x03")
		}
	}
}

func decryptSSHKey(key []byte, password []byte) (ssh.Signer, error) {
	block, _ := pem.Decode(key)
	pem, err := x509.DecryptPEMBlock(block, password)
	if err != nil {
		return nil, err
	}
	sshkey, err := x509.ParsePKCS1PrivateKey(pem)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(sshkey)
}
