package sync

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	gosync "sync"

	"github.com/wallix/awless/cloud/aws"
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
	region        string
	infraService  *aws.Infra
	accessService *aws.Access
}

func NewSyncer(region string, inf *aws.Infra, access *aws.Access) Syncer {
	repo, err := repo.New()
	if err != nil {
		panic(err)
	}

	return &syncer{
		Repo:          repo,
		region:        region,
		infraService:  inf,
		accessService: access,
	}
}

func (s *syncer) Sync() (*graph.Graph, *graph.Graph, error) {
	var wg gosync.WaitGroup

	type results struct {
		awsData interface{}
		err     error
	}

	resultc := make(chan results)

	wg.Add(1)
	go func() {
		defer wg.Done()
		res := results{}
		res.awsData, res.err = s.infraService.FetchAwsInfra()
		resultc <- res
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res := results{}
		res.awsData, res.err = s.accessService.FetchAwsAccess()
		resultc <- res
	}()

	go func() {
		wg.Wait()
		close(resultc)
	}()

	var awsInfra *aws.AwsInfra
	var awsAccess *aws.AwsAccess

	for res := range resultc {
		if res.err != nil {
			return nil, nil, res.err
		}

		switch res.awsData.(type) {
		case *aws.AwsInfra:
			awsInfra = res.awsData.(*aws.AwsInfra)
		case *aws.AwsAccess:
			awsAccess = res.awsData.(*aws.AwsAccess)
		default:
			return nil, nil, fmt.Errorf("unexpected returned type %T", res.awsData)
		}
	}

	infrag, err := aws.BuildAwsInfraGraph(s.region, awsInfra)

	tofile, err := infrag.Marshal()
	if err != nil {
		return nil, nil, err
	}
	if err = ioutil.WriteFile(filepath.Join(config.RepoDir, config.InfraFilename), tofile, 0600); err != nil {
		return nil, nil, err
	}

	accessg, err := aws.BuildAwsAccessGraph(s.region, awsAccess)

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
