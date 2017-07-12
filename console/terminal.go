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
	"os"
	"os/signal"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func GetTerminalWidth() int {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0
	}
	return w
}

func GetTerminalHeight() int {
	_, h, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0
	}
	return h
}

func InteractiveTerminal(client *ssh.Client) error {
	defer client.Close()

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
