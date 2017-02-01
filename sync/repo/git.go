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
