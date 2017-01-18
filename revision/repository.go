package revision

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/badwolf/triple/node"
	"gopkg.in/src-d/go-git.v4"
)

// Repository represents the git repository containing RDF files (infra and access)
type Repository struct {
	gitRepository *git.Repository
	files         []string
	path          string
}

type fetchParameter int

const (
	GroupAll fetchParameter = iota
	GroupByDay
	GroupByWeek
	NoGroup
)

func OpenRepository(path string) (*Repository, error) {
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		if _, err := executeGitCommand(path, "init"); err != nil {
			return nil, err
		}
	}

	repo, err := git.NewFilesystemRepository(filepath.Join(path, ".git"))
	return &Repository{gitRepository: repo, path: path}, err
}

// LastDiffs list the last revisions for the files in parameters (if no file in parameter, for all repository files)
func (rr *Repository) LastDiffs(numberRevisions int, root *node.Node, param fetchParameter, files ...string) ([]*Diff, error) {
	var diffs []*Diff
	if len(files) == 0 {
		files = rr.files
	}

	revisions, err := rr.lastRevisions(numberRevisions)
	if err != nil {
		return diffs, err
	}
	return rr.generateDiffs(generateRevisionPairs(revisions, param), root, files)
}

func (rr *Repository) lastRevisions(n int) ([]*Revision, error) {
	var res []*Revision

	head, err := rr.gitRepository.Head()
	if err != nil {
		return res, nil //Empty repository
	}

	commit, err := rr.gitRepository.Commit(head.Hash())
	if err != nil {
		return res, err
	}

	commits, err := commit.History()
	for i, parent := range commits {
		if i >= n {
			break
		}
		res = append(res, NewRevision(parent))
	}

	if len(res) < n {
		res = append(res, initRevision)
	}

	return res, err
}

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
