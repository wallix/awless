package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// OpenRepository opens a new or existing git repository
func OpenRepository(path string) (*Repository, error) {
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		if _, err := executeGitCommand(path, "init"); err != nil {
			return nil, err
		}
	}

	repo, err := git.NewFilesystemRepository(filepath.Join(path, ".git"))
	return &Repository{gitRepository: repo, path: path}, err
}

// CommitIfChanges creates a new git commit if there are changes in the infra and access RDF files
func (rr *Repository) CommitIfChanges(files ...string) error {
	for _, file := range files {
		rr.addFile(file)
	}

	return rr.commitIfChanges()
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

func (rr *Repository) hasChanges() (bool, error) {
	stdout, err := executeGitCommand(rr.path, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return !(strings.TrimSpace(stdout) == ""), nil
}

func (rr *Repository) addFile(path string) error {
	rr.files = append(rr.files, path)
	return nil
}

func (rr *Repository) commitIfChanges(overwriteTime ...time.Time) error {
	for _, filePath := range rr.files {
		if _, err := executeGitCommand(rr.path, "add", filePath); err != nil {
			return err
		}
	}

	if hasChanges, e := rr.hasChanges(); e != nil {
		return e
	} else if !hasChanges {
		return nil
	}
	var env []string
	if len(overwriteTime) != 0 {
		env = []string{fmt.Sprintf("GIT_AUTHOR_DATE=%s", overwriteTime[0]), fmt.Sprintf("GIT_COMMITTER_DATE=%s", overwriteTime[0])}
	}

	if _, err := executeGitCommandWithEnv(rr.path, env, "-c", "user.name='awless'", "-c", "user.email='git@awless.io'", "commit", "-m", "new sync"); err != nil {
		return err
	}

	return nil
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
