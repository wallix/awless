package revision

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

var ErrGitNotFound = errors.New("git: executable has not been found")

func executeGitCommand(dir string, command ...string) (string, error) {
	return executeGitCommandWithEnv(dir, []string{}, command...)
}

func executeGitCommandWithEnv(dir string, env []string, command ...string) (string, error) {
	git, err := exec.LookPath("git")
	if err != nil {
		return "", ErrGitNotFound
	}
	cmd := exec.Command(git, command...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = env
	err = cmd.Run()
	if err != nil || stderr.String() != "" {
		return "", fmt.Errorf("git error: %s: %s", err.Error(), stderr.String())
	}
	return stdout.String(), nil
}
