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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"golang.org/x/crypto/ssh"
)

type Credentials struct {
	IP      string
	User    string
	KeyPath string
}

func NewSSHClient(cred *Credentials) (*ssh.Client, error) {
	privateKey, err := ioutil.ReadFile(cred.KeyPath)
	if os.IsNotExist(err) {
		privateKey, err = ioutil.ReadFile(cred.KeyPath + ".pem")
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Could not find a SSH key at path '%s'.", cred.KeyPath)
		}
		cred.KeyPath = cred.KeyPath + ".pem"
	}
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: cred.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout: 2 * time.Second,
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

func propagateSignals(signalc chan os.Signal, session *ssh.Session, stdin io.WriteCloser) {
	for s := range signalc {
		switch s {
		case os.Interrupt:
			fmt.Fprint(stdin, "\x03")
		}
	}
}
