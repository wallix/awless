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

package repo

import (
	"bytes"
	"fmt"
	"os/exec"
)

type gitCmd struct {
	dir string
	env []string
}

func newGit(workdir string, envs ...string) *gitCmd {
	return &gitCmd{dir: workdir}
}

func (g *gitCmd) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.dir
	cmd.Env = g.env

	out, err := cmd.Output()
	if err != nil {
		var msg bytes.Buffer

		msg.WriteString(err.Error())
		msg.WriteByte('\n')
		if exitErr, ok := err.(*exec.ExitError); ok {
			msg.Write(exitErr.Stderr)
			msg.WriteByte('\n')
		}

		return "", fmt.Errorf("git error: %s", msg.String())
	}

	return string(out), nil
}
