package sync

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	gosync "sync"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync/repo"
)

var DefaultSyncer Syncer

type Syncer interface {
	repo.Repo
	Sync() (*graph.Graph, *graph.Graph, error)
}

type syncer struct {
	repo.Repo
	region                      string
	infraService, accessService cloud.Service
}

func NewSyncer(region string, services ...cloud.Service) Syncer {
	repo, err := repo.New()
	if err != nil {
		panic(err)
	}

	syncer := &syncer{
		Repo:   repo,
		region: region,
	}

	for _, service := range services {
		switch service.Name() {
		case "infra":
			syncer.infraService = service
		case "access":
			syncer.accessService = service
		default:
			panic(fmt.Sprintf("syncer: cannot init: unexpected service name %s", service.Name()))
		}
	}

	return syncer
}

func (s *syncer) Sync() (*graph.Graph, *graph.Graph, error) {
	var wg gosync.WaitGroup
	var infrag, accessg *graph.Graph

	errorc := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		infrag, err = s.infraService.FetchResources()
		errorc <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		accessg, err = s.accessService.FetchResources()
		errorc <- err
	}()

	go func() {
		wg.Wait()
		close(errorc)
	}()

	for err := range errorc {
		if err != nil {
			return nil, nil, err
		}
	}

	tofile, err := infrag.Marshal()
	if err != nil {
		return nil, nil, err
	}
	if err = ioutil.WriteFile(filepath.Join(config.RepoDir, config.InfraFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	tofile, err = accessg.Marshal()
	if err != nil {
		return nil, nil, err
	}
	if err := ioutil.WriteFile(filepath.Join(config.RepoDir, config.AccessFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	if err := s.Commit(config.InfraFilename, config.AccessFilename); err != nil {
		return nil, nil, err
	}

	return infrag, accessg, nil
}
