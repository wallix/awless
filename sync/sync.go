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
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	SetLogger(*logger.Logger)
}

type syncer struct {
	repo.Repo
	dryrun bool
	logger *logger.Logger
}

func NewSyncer(dryrun bool) Syncer {
	repo, err := repo.New()
	if err != nil {
		panic(err)
	}

	return &syncer{Repo: repo, dryrun: dryrun, logger: logger.DiscardLogger}
}

func (s *syncer) SetLogger(l *logger.Logger) { s.logger = l }

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

Loop:
	for {
		select {
		case srvErr, ok := <-errorc:
			if ok {
				if srvErr.err == cloud.ErrFetchAccessDenied {
					logger.Errorf("sync: access denied to service %s", srvErr.name)
				} else if srvErr.err != nil {
					return graphs, srvErr.err
				}
			}
		case res, ok := <-resultc:
			if !ok {
				break Loop
			}
			logger.ExtraVerbosef("sync: fetched %s service took %s", res.name, time.Since(res.start))
			graphs[res.name] = res.gph
		}
	}

	if !s.dryrun {
		var filenames []string

		for name, g := range graphs {
			filename := fmt.Sprintf("%s.rdf", name)
			tofile, err := g.Marshal()
			if err != nil {
				return graphs, err
			}
			if err = ioutil.WriteFile(filepath.Join(config.RepoDir, filename), tofile, 0600); err != nil {
				return graphs, err
			}
			filenames = append(filenames, filename)
		}

		if err := s.Commit(filenames...); err != nil {
			return graphs, err
		}
	}

	return graphs, nil
}

func LoadCurrentLocalGraph(serviceName string) *graph.Graph {
	path := filepath.Join(config.RepoDir, fmt.Sprintf("%s.rdf", serviceName))
	g, err := graph.NewGraphFromFile(path)
	if err != nil {
		return graph.NewGraph()
	}
	return g
}
