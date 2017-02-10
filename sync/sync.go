package sync

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	gosync "sync"

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
		name string
		gph  *graph.Graph
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
			g, err := srv.FetchResources()
			errorc <- &srvErr{name: srv.Name(), err: err}
			resultc <- &result{name: srv.Name(), gph: g}
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
					logger.Verbosef("sync: access denied to service %s", srvErr.name)
				} else if srvErr.err != nil {
					return graphs, srvErr.err
				}
			}
		case res, ok := <-resultc:
			if !ok {
				break Loop
			}
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
