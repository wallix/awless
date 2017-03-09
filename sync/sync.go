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

package sync

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	gosync "sync"
	"time"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync/repo"
)

var DefaultSyncer Syncer

type Syncer interface {
	repo.Repo
	Sync(...cloud.Service) (map[string]*graph.Graph, error)
}

type syncer struct {
	repo.Repo
	logger *logger.Logger
}

func NewSyncer(l *logger.Logger) Syncer {
	repo, err := repo.New()
	if err != nil {
		panic(err)
	}

	if l == nil {
		l = logger.DiscardLogger
	}

	return &syncer{Repo: repo, logger: l}
}

func (s *syncer) Sync(services ...cloud.Service) (map[string]*graph.Graph, error) {
	graphs := make(map[string]*graph.Graph)
	var workers gosync.WaitGroup

	type result struct {
		name  string
		gph   *graph.Graph
		start time.Time
	}

	type srvErr struct {
		name string
		err  error
	}

	resultc := make(chan *result, len(services))
	errorc := make(chan *srvErr, len(services))

	for _, service := range services {
		workers.Add(1)
		go func(srv cloud.Service) {
			defer workers.Done()
			start := time.Now()
			g, err := srv.FetchResources()
			errorc <- &srvErr{name: srv.Name(), err: err}
			resultc <- &result{name: srv.Name(), gph: g, start: start}
		}(service)
	}

	go func() {
		workers.Wait()
		close(errorc)
		close(resultc)
	}()

	var allErrors []error

Loop:
	for {
		select {
		case srvErr, ok := <-errorc:
			if ok && srvErr.err != nil {
				allErrors = append(allErrors, fmt.Errorf("syncing %s: %s", srvErr.name, srvErr.err))
			}
		case res, ok := <-resultc:
			if !ok {
				break Loop
			}
			logger.ExtraVerbosef("sync: fetched %s service took %s", res.name, time.Since(res.start))
			graphs[res.name] = res.gph
		}
	}

	var filenames []string

	for name, g := range graphs {
		filename := fmt.Sprintf("%s.rdf", name)
		tofile, err := g.Marshal()
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("marshal %s: %s", filename, err))
		}
		filepath := filepath.Join(config.RepoDir, filename)
		if err = ioutil.WriteFile(filepath, tofile, 0600); err != nil {
			allErrors = append(allErrors, fmt.Errorf("writing %s: %s", filepath, err))
		}
		filenames = append(filenames, filename)
	}

	if err := s.Commit(filenames...); err != nil {
		allErrors = append(allErrors, fmt.Errorf("commit %s: %s", strings.Join(filenames, ", "), err))
	}

	return graphs, concatErrors(allErrors)
}

func concatErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var lines []string
	for _, err := range errs {
		lines = append(lines, err.Error())
	}

	return errors.New(strings.Join(lines, "\n"))
}

func LoadCurrentLocalGraph(serviceName string) *graph.Graph {
	path := filepath.Join(config.RepoDir, fmt.Sprintf("%s.rdf", serviceName))
	g, err := graph.NewGraphFromFile(path)
	if err != nil {
		return graph.NewGraph()
	}
	return g
}

func LoadAllGraphs() (*graph.Graph, error) {
	path := filepath.Join(config.RepoDir, "*.rdf")
	files, _ := filepath.Glob(path)

	g := graph.NewGraph()

	var buff bytes.Buffer
	for _, f := range files {
		reader, err := os.Open(f)
		if err != nil {
			return g, fmt.Errorf("loading '%s': %s", f, err)
		}
		io.Copy(&buff, reader)
		buff.WriteByte('\n')
	}

	g.Unmarshal(buff.Bytes())
	return g, nil
}
