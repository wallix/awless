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

type NullRepo struct{}

func (NullRepo) Commit(files ...string) error         { return nil }
func (NullRepo) List() ([]*Rev, error)                { return nil, nil }
func (NullRepo) LoadRev(version string) (*Rev, error) { return nil, nil }
func (NullRepo) BaseDir() string                      { return "" }

type gitRepo struct {
	repo    *git.Repository
	basedir string
}

func BaseDir() string {
	return filepath.Join(os.Getenv("__AWLESS_HOME"), "aws", "rdf")
}

func New() (Repo, error) {
	dir := BaseDir()
	os.MkdirAll(dir, 0700)
	return newGitRepo(dir)
}

func newGitRepo(path string) (Repo, error) {
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		if _, err := git.PlainInit(path, false); err != nil {
			return nil, err
		}
	}

	repo, err := git.PlainOpen(path)
	return &gitRepo{repo: repo, basedir: path}, err
}

func (r *gitRepo) BaseDir() string {
	return r.basedir
}

func (r *gitRepo) List() ([]*Rev, error) {
	var all []*Rev

	iter, err := r.repo.CommitObjects()
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

	sort.Slice(all, func(i, j int) bool { return all[i].Date.Before(all[j].Date) })

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
		sort.Slice(v, func(i, j int) bool { return v[i].Date.After(v[j].Date) })
		reduce = append(reduce, v[0])
	}

	return reduce
}

func (r *gitRepo) LoadRev(version string) (*Rev, error) {
	rev := &Rev{Id: version}

	commit, err := r.repo.CommitObject(plumbing.NewHash(version))
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

func (r *gitRepo) Commit(relativePaths ...string) error {
	wt, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	for _, f := range relativePaths {
		if _, err := wt.Add(f); err != nil {
			return err
		}
	}

	msg := fmt.Sprintf("syncing %s", strings.Join(relativePaths, ", "))
	committer := &object.Signature{Name: "awlessCLI", When: time.Now(), Email: "git@awless.io"}

	_, err = wt.Commit(msg, &git.CommitOptions{Author: committer})
	return err
}
