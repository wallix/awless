package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/badwolf/triple/node"
	git "github.com/libgit2/git2go"
	"github.com/wallix/awless/rdf"
)

// Repository represents the git repository containing RDF files (infra and access)
type Repository struct {
	gitRepository *git.Repository
	index         *git.Index
	files         []string
}

// CommitDiff represents the deleted/inserted RDF triples of a git commit
type CommitDiff struct {
	ParentTime   time.Time
	ChildTime    time.Time
	ParentCommit string
	ChildCommit  string
	GraphDiff    *rdf.Diff
}

type fetchParameter int

const (
	GroupAll fetchParameter = iota
	GroupByDay
	GroupByWeek
	NoGroup
)

type commitPair struct {
	parent *git.Commit
	child  *git.Commit
}

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

// LastDiffs list the last numberCommits commits for the files in parmeters (if no file in parameter, for all repository files)
func (rr *Repository) LastDiffs(numberCommits int, root *node.Node, param fetchParameter, files ...string) ([]*CommitDiff, error) {
	var diffs []*CommitDiff
	if len(files) == 0 {
		files = rr.files
	}

	commits, err := rr.lastCommits(numberCommits)
	if err != nil {
		return diffs, err
	}
	return generateCommitDiffs(generateCommitPairs(commits, param), rr.gitRepository, root, files)
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

func (rr *Repository) commitIfChanges() error {
	for _, filePath := range rr.files {
		if err := rr.index.AddByPath(filePath); err != nil {
			return err
		}
	}

	treeID, err := rr.index.WriteTree()
	if err != nil {
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

	sig := &git.Signature{Name: "awless", Email: "git@awless.io", When: time.Now()}
	if _, err = rr.gitRepository.CreateCommit("HEAD", sig, sig, "new sync", tree, parents...); err != nil {
		return err
	}

	return nil
}

func (rr *Repository) lastCommits(n int) ([]*git.Commit, error) {
	var res []*git.Commit
	head, err := rr.gitRepository.Head()
	if err != nil {
		return res, nil //Empty repository
	}

	commit, err := rr.gitRepository.LookupCommit(head.Target())
	if err != nil {
		return res, err
	}
	res = append(res, commit)

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
		res = append(res, parent)

		if numberParents == 0 {
			break
		}
	}

	return res, err
}

func generateCommitPairs(commits []*git.Commit, param fetchParameter) []*commitPair {
	var res []*commitPair

	switch param {
	case GroupAll:
		return []*commitPair{{parent: commits[len(commits)-1], child: commits[0]}}
	case GroupByDay:
		if len(commits) == 0 {
			return res
		}
		commit := commits[0]
		previousAddedCommit := commit
		time := commits[0].Committer().When
		for i := 1; i < len(commits); i++ {
			newCommit := commits[i]
			if newCommit != nil && time.Sub(newCommit.Committer().When).Hours() > 24. {
				res = append(res, &commitPair{parent: newCommit, child: commit})
				time = newCommit.Committer().When
				previousAddedCommit = newCommit
			}
			commit = newCommit
		}
		res = append(res, &commitPair{parent: commit, child: previousAddedCommit})
	case GroupByWeek:
		if len(commits) == 0 {
			return res
		}
		commit := commits[0]
		previousAddedCommit := commit
		time := commits[0].Committer().When
		for i := 1; i < len(commits); i++ {
			newCommit := commits[i]
			if newCommit != nil && time.Sub(newCommit.Committer().When).Hours() > 7*24. {
				res = append(res, &commitPair{parent: newCommit, child: commit})
				time = newCommit.Committer().When
				previousAddedCommit = newCommit
			}
			commit = newCommit
		}
		res = append(res, &commitPair{parent: commit, child: previousAddedCommit})
	default:
		for i := 0; i < len(commits)-1; i++ {
			res = append(res, &commitPair{parent: commits[i+1], child: commits[i]})
		}
	}
	return res
}

func generateCommitDiffs(pairs []*commitPair, repo *git.Repository, root *node.Node, forFiles []string) ([]*CommitDiff, error) {
	var res []*CommitDiff
	for _, commitPair := range pairs {
		diff, err := newCommitDiff(commitPair.parent, commitPair.child, repo, root, forFiles)
		if err != nil {
			return res, err
		}
		res = append(res, diff)
	}
	return res, nil
}

func newCommitDiff(parent, commit *git.Commit, repo *git.Repository, root *node.Node, forFiles []string) (*CommitDiff, error) {
	var parentTree *git.Tree
	var err error
	var parentTime time.Time
	var parentCommit string
	if parent != nil {
		parentTime = parent.Committer().When
		parentCommit = parent.Id().String()
		parentTree, err = parent.Tree()
		if err != nil {
			return nil, err
		}
	}
	commitTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	parentGraph, err := gitTreeToGraph(parentTree, repo, forFiles)
	if err != nil {
		return nil, err
	}
	commitGraph, err := gitTreeToGraph(commitTree, repo, forFiles)
	if err != nil {
		return nil, err
	}
	diff, err := rdf.DefaultDiffer.Run(root, commitGraph, parentGraph)
	if err != nil {
		return nil, err
	}

	res := &CommitDiff{
		ChildTime:    commit.Committer().When,
		ChildCommit:  commit.Id().String(),
		ParentTime:   parentTime,
		ParentCommit: parentCommit,
		GraphDiff:    diff,
	}

	return res, nil
}

func gitTreeToGraph(tree *git.Tree, repo *git.Repository, files []string) (*rdf.Graph, error) {
	g := rdf.NewGraph()
	if tree == nil {
		return g, nil
	}
	nbEntries := tree.EntryCount()
	for i := uint64(0); i < nbEntries; i++ {
		entry := tree.EntryByIndex(i)
		if entry.Type == git.ObjectBlob && contains(files, entry.Name) {
			blob, err := repo.LookupBlob(entry.Id)
			if err != nil {
				return g, err
			}
			err = g.Unmarshal(blob.Contents())
			if err != nil {
				return g, err
			}
		}
	}
	return g, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
