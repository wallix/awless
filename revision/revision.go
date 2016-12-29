package revision

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
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
	Time      time.Time
	Commit    string
	GraphDiff *rdf.Diff
}

// CommitIfChanges creates a new git commit if there are changes in the infra and access RDF files
func CommitIfChanges(repositoryPath string, files ...string) error {
	rr, err := openRepository(repositoryPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		rr.addFile(file)
	}

	return rr.commitIfChanges()
}

// LastDiffs list the last numberCommits commits for the files in parmeters (if no file in parameter, for all repository files)
func LastDiffs(repositoryPath string, numberCommits int, files ...string) ([]*CommitDiff, error) {
	var diffs []*CommitDiff
	rr, err := openRepository(repositoryPath)
	if err != nil {
		return diffs, err
	}

	return rr.lastsDiffs(numberCommits, files...)
}

func openRepository(path string) (*Repository, error) {
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

func (rr *Repository) lastsDiffs(numberCommits int, files ...string) ([]*CommitDiff, error) {
	if len(files) == 0 {
		files = rr.files
	}
	var result []*CommitDiff
	head, err := rr.gitRepository.Head()
	if err != nil {
		return result, nil //Empty repository
	}

	headCommit, err := rr.gitRepository.LookupCommit(head.Target())
	if err != nil {
		return result, err
	}

	commit := headCommit
	for i := 0; i < numberCommits; i++ {
		numberParents := commit.ParentCount()
		var parent *git.Commit
		if numberParents > 1 {
			return result, fmt.Errorf("The %s commit has more than 1 parent (%d parents)", commit.Id().String(), numberParents)
		} else if numberParents == 1 {
			parent = commit.Parent(0)
		}
		diff, err := newCommitDiff(parent, commit, rr.gitRepository, files)
		if err != nil {
			return result, err
		}
		result = append(result, diff)
		if numberParents == 0 {
			break
		}
		commit = parent
	}

	return result, nil
}

func newCommitDiff(parent, commit *git.Commit, repo *git.Repository, forFiles []string) (*CommitDiff, error) {
	var parentTree *git.Tree
	var err error
	if parent != nil {
		parentTree, err = parent.Tree()
		if err != nil {
			return nil, err
		}
	}
	commitTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	gitDiff, err := repo.DiffTreeToTree(parentTree, commitTree, &git.DiffOptions{})
	if err != nil {
		return nil, err
	}
	parentGraph, err := gitTreeToGraph(parentTree, repo, forFiles)
	if err != nil {
		return nil, err
	}
	res := &CommitDiff{
		Time:      commit.Committer().When,
		Commit:    commit.Id().String(),
		GraphDiff: rdf.NewEmptyDiffFromGraph(parentGraph),
	}

	err = gitDiff.ForEach(res.appendDiffFunction(forFiles), git.DiffDetailLines)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (c *CommitDiff) appendDiffFunction(forFiles []string) func(delta git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
	return func(delta git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		return func(git.DiffHunk) (git.DiffForEachLineCallback, error) {
			return func(line git.DiffLine) error {
				if delta.Flags == git.DiffFlagNotBinary && (contains(forFiles, delta.NewFile.Path) || contains(forFiles, delta.OldFile.Path)) {
					if line.Origin == git.DiffLineAddition {
						str := strings.TrimSpace(line.Content)
						if str != "" {
							t, e := triple.Parse(str, literal.DefaultBuilder())
							if e != nil {
								return e
							}
							c.GraphDiff.AddInserted(t, rdf.ParentOfPredicate)
						}
					}
					if line.Origin == git.DiffLineDeletion {
						str := strings.TrimSpace(line.Content)
						if str != "" {
							t, e := triple.Parse(str, literal.DefaultBuilder())
							if e != nil {
								return e
							}
							c.GraphDiff.AddDeleted(t, rdf.ParentOfPredicate)
						}
					}
				}
				return nil
			}, nil
		}, nil
	}
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
