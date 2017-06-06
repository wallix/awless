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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wallix/awless/graph"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Rev struct {
	Id   string
	Date time.Time

	Infra  *graph.Graph
	Access *graph.Graph
}

func (r *Rev) DateString() string {
	return r.Date.Format("Mon Jan 2 15:04:05")
}

type Repo interface {
	Commit(files ...string) error
	List() ([]*Rev, error)
	LoadRev(version string) (*Rev, error)
	BaseDir() string
}

type noRevisionRepo struct {
	basedir string
}

func (*noRevisionRepo) Commit(files ...string) error         { return nil }
func (*noRevisionRepo) LoadRev(version string) (*Rev, error) { return &Rev{}, nil }
func (*noRevisionRepo) List() ([]*Rev, error)                { return nil, nil }
func (r *noRevisionRepo) BaseDir() string                    { return r.basedir }

type gitRepo struct {
	repo    *git.Repository
	files   []string
	basedir string
}

func BaseDir() string {
	return filepath.Join(os.Getenv("__AWLESS_HOME"), "aws", "rdf")
}

func New() (Repo, error) {
	dir := BaseDir()
	os.MkdirAll(dir, 0700)

	if IsGitInstalled() {
		return newGitRepo(dir)
	}

	return &noRevisionRepo{dir}, nil
}

func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func newGitRepo(path string) (Repo, error) {
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		if _, err := newGit(path).run("init"); err != nil {
			return nil, err
		}
	}

	repo, err := git.NewFilesystemRepository(filepath.Join(path, ".git"))
	return &gitRepo{repo: repo, basedir: path}, err
}

func (r *gitRepo) BaseDir() string {
	return r.basedir
}

func (r *gitRepo) List() ([]*Rev, error) {
	var all []*Rev

	iter, err := r.repo.Commits()
	if err != nil {
		return all, err
	}
	defer iter.Close()

	for {
		commit, err := iter.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Sprintf("error listing repo revisions: %s", err))
		}

		all = append(all, &Rev{Id: commit.Hash.String(), Date: commit.Committer.When})
	}

	sort.Sort(revsByDate(all))

	return all, nil
}

func reduceToLastRevOfEachDay(revs []*Rev) []*Rev {
	perDay := make(map[string][]*Rev)

	for _, rev := range revs {
		day := rev.Date.Format("2006-01-02")
		perDay[day] = append(perDay[day], rev)
	}

	reduce := []*Rev{}
	for _, v := range perDay {
		sort.Sort(sort.Reverse(revsByDate(v)))
		reduce = append(reduce, v[0])
	}

	return reduce
}

func (r *gitRepo) LoadRev(version string) (*Rev, error) {
	rev := &Rev{Id: version}

	commit, err := r.repo.Commit(plumbing.NewHash(version))
	if err != nil {
		return nil, err
	}

	rev.Date = commit.Committer.When

	rev.Infra = graph.NewGraph()
	rev.Access = graph.NewGraph()

	if err := unmarshalIntoGraph(rev.Infra, commit, "infra.triples"); err != nil {
		return rev, err
	}
	if err := unmarshalIntoGraph(rev.Access, commit, "access.triples"); err != nil {
		return rev, err
	}

	return rev, nil
}

func unmarshalIntoGraph(g *graph.Graph, commit *object.Commit, filename string) error {
	f, err := commit.File(filename)
	if err != nil && err != object.ErrFileNotFound {
		return err
	} else if err == nil {
		contents, err := f.Contents()
		if err != nil {
			return err
		}
		g.Unmarshal([]byte(contents))
	}
	return nil
}

func (r *gitRepo) Commit(files ...string) error {
	for _, path := range files {
		r.files = append(r.files, path)
	}

	for _, path := range r.files {
		if _, err := newGit(r.BaseDir()).run("add", path); err != nil {
			return err
		}
	}

	if hasChanges, err := r.hasChanges(); err != nil {
		return err
	} else if !hasChanges {
		return nil
	}

	_, err := newGit(r.BaseDir()).run(
		append(awlessCommitter, "commit", "-m", fmt.Sprintf("syncing %s", strings.Join(files, ", ")))...,
	)

	return err
}

var awlessCommitter = []string{"-c", "user.name='awless'", "-c", "user.email='git@awless.io'"}

func (r *gitRepo) hasChanges() (bool, error) {
	stdout, err := newGit(r.BaseDir()).run("status", "--porcelain")
	if err != nil {
		return false, err
	}

	return !(strings.TrimSpace(stdout) == ""), nil
}
