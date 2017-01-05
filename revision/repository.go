package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/badwolf/triple/node"
	git "github.com/libgit2/git2go"
)

// Repository represents the git repository containing RDF files (infra and access)
type Repository struct {
	gitRepository *git.Repository
	index         *git.Index
	files         []string
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
		if _, err := git.InitRepository(path, false); err != nil {
			return nil, err
		}
	}

	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}

	idx, err := repo.Index()
	if err != nil {
		return nil, err
	}
	return &Repository{gitRepository: repo, index: idx}, nil
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
	status, err := rr.gitRepository.StatusList(&git.StatusOptions{})
	if err != nil {
		return false, err
	}
	changes, err := status.EntryCount()
	if err != nil {
		return false, err
	}
	return (changes != 0), nil
}

func (rr *Repository) addFile(path string) error {
	rr.files = append(rr.files, path)
	return rr.index.AddByPath(path)
}

func (rr *Repository) commitIfChanges(overwriteTime ...time.Time) error {
	for _, filePath := range rr.files {
		if err := rr.index.AddByPath(filePath); err != nil {
			return err
		}
	}

	treeID, err := rr.index.WriteTree()
	if err != nil {
		return err
	}

	if err = rr.index.Write(); err != nil {
		return err
	}

	if hasChanges, e := rr.hasChanges(); e != nil {
		return e
	} else if !hasChanges {
		return nil
	}

	tree, err := rr.gitRepository.LookupTree(treeID)
	if err != nil {
		return err
	}

	var parents []*git.Commit

	head, err := rr.gitRepository.Head()
	if err == nil {
		headCommit, e := rr.gitRepository.LookupCommit(head.Target())
		if e != nil {
			return e
		}
		parents = append(parents, headCommit)
	}
	time := time.Now()
	if len(overwriteTime) > 0 {
		time = overwriteTime[0]
	}
	sig := &git.Signature{Name: "awless", Email: "git@awless.io", When: time}
	if _, err = rr.gitRepository.CreateCommit("HEAD", sig, sig, "new sync", tree, parents...); err != nil {
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

	commit, err := rr.gitRepository.LookupCommit(head.Target())
	if err != nil {
		return res, err
	}

	res = append(res, NewRevision(commit))

	for i := 0; i < n; i++ {
		numberParents := commit.ParentCount()
		var parent *git.Commit

		if numberParents > 1 {
			return res, fmt.Errorf("The %s commit has more than 1 parent (%d parents)", commit.Id().String(), numberParents)
		}
		if numberParents == 1 {
			parent = commit.Parent(0)
			commit = parent
		}

		if parent == nil {
			res = append(res, initRevision)
			break
		}
		res = append(res, NewRevision(parent))

		if numberParents == 0 {
			break
		}
	}

	return res, err
}
