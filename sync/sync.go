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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	gosync "sync"
	"time"

	"runtime"

	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/sync/repo"
)

const fileExt = ".triples"

var DefaultSyncer Syncer

type Syncer interface {
	repo.Repo
	Sync(...cloud.Service) (map[string]*graph.Graph, error)
}

type syncer struct {
	repo.Repo
	logger *logger.Logger
}

func NewSyncer(l ...*logger.Logger) Syncer {
	repo, err := repo.New()
	if err != nil {
		panic(err)
	}

	s := &syncer{Repo: repo}

	if len(l) > 0 {
		s.logger = l[0]
	} else {
		s.logger = logger.DiscardLogger
	}

	return s
}

func (s *syncer) Sync(services ...cloud.Service) (map[string]*graph.Graph, error) {
	var workers gosync.WaitGroup

	type result struct {
		service cloud.Service
		gph     *graph.Graph
		start   time.Time
		err     error
	}

	resultc := make(chan *result, len(services))

	for _, service := range services {
		if service.IsSyncDisabled() {
			s.logger.Verbosef("sync: *disabled* for service %s", service.Name())
			continue
		}
		workers.Add(1)
		go func(srv cloud.Service) {
			defer workers.Done()
			start := time.Now()
			g, err := srv.FetchResources()
			resultc <- &result{service: srv, gph: g, start: start, err: err}
		}(service)
	}

	go func() {
		workers.Wait()
		close(resultc)
	}()

	var allErrors []error
	graphs := make(map[string]*graph.Graph)
	servicesByName := make(map[string]cloud.Service)
Loop:
	for {
		select {
		case res, ok := <-resultc:
			if !ok {
				break Loop
			}
			if res.err != nil {
				allErrors = append(allErrors, fmt.Errorf("syncing %s: %s", res.service.Name(), res.err))
			} else {
				s.logger.ExtraVerbosef("sync: fetched %s service took %s", res.service.Name(), time.Since(res.start))
			}
			if serv := res.service; serv != nil {
				servicesByName[serv.Name()] = serv
				if res.gph != nil {
					graphs[serv.Name()] = res.gph
				}
			}
		}
	}

	var filepaths []string

	for name, g := range graphs {
		serviceRegion := servicesByName[name].Region()
		serviceDir := filepath.Join(s.BaseDir(), serviceRegion)
		os.MkdirAll(serviceDir, 0700)

		filename := fmt.Sprintf("%s%s", name, fileExt)
		fullpath := filepath.Join(serviceDir, filename)
		f, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("opening %s: %s", fullpath, err))
			continue
		}
		if err := g.MarshalTo(f); err != nil {
			allErrors = append(allErrors, fmt.Errorf("marshal to %s: %s", fullpath, err))
		}

		filepaths = append(filepaths, filepath.Join(serviceRegion, filename))
		if err := f.Close(); err != nil {
			allErrors = append(allErrors, fmt.Errorf("closing file %s: %s", fullpath, err))
		}
	}

	if runtime.GOOS != "windows" { // https://github.com/wallix/awless/issues/119
		if err := s.Commit(filepaths...); err != nil {
			allErrors = append(allErrors, fmt.Errorf("committing %s: %s", strings.Join(filepaths, ", "), err))
		}
	}

	return graphs, concatErrors(allErrors)
}

func concatErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	lines := []string{"syncing errors:"}
	for _, err := range errs {
		lines = append(lines, fmt.Sprintf("\t\t%s", err))
	}

	return errors.New(strings.Join(lines, "\n"))
}

func LoadLocalGraphForService(serviceName, region string) *graph.Graph {
	regionDir := region
	if aws.IsGlobalService(serviceName) {
		regionDir = "global"
	}
	path := filepath.Join(repo.BaseDir(), regionDir, fmt.Sprintf("%s%s", serviceName, fileExt))
	g, err := graph.NewGraphFromFile(path)
	if err != nil {
		return graph.NewGraph()
	}
	return g
}

func LoadLocalGraphs(region string) (*graph.Graph, error) {
	var files []string
	globalFiles, _ := filepath.Glob(filepath.Join(repo.BaseDir(), "global", fmt.Sprintf("*%s", fileExt)))
	regionFiles, _ := filepath.Glob(filepath.Join(repo.BaseDir(), region, fmt.Sprintf("*%s", fileExt)))

	files = append(files, globalFiles...)
	files = append(files, regionFiles...)

	g := graph.NewGraph()

	var readers []io.Reader
	for _, f := range files {
		reader, err := os.Open(f)
		if err != nil {
			return g, fmt.Errorf("loading '%s': %s", f, err)
		}
		readers = append(readers, reader)
	}

	err := g.UnmarshalMultiple(readers...)
	return g, err
}

func LoadAllLocalGraphs() (*graph.Graph, error) {
	path := filepath.Join(repo.BaseDir(), "*", fmt.Sprintf("*%s", fileExt))
	files, _ := filepath.Glob(path)

	g := graph.NewGraph()

	var readers []io.Reader
	for _, f := range files {
		reader, err := os.Open(f)
		if err != nil {
			return g, fmt.Errorf("loading '%s': %s", f, err)
		}
		readers = append(readers, reader)
	}

	err := g.UnmarshalMultiple(readers...)
	return g, err
}
