package config

import (
	"path/filepath"

	"github.com/wallix/awless/graph"
)

func LoadInfraGraph() (*graph.Graph, error) {
	var err error
	infra := graph.NewGraph()
	if !AwlessFirstSync {
		infra, err = graph.NewGraphFromFile(filepath.Join(RepoDir, InfraFilename))
		if err != nil {
			return infra, err
		}
	}
	return infra, nil
}

func LoadAccessGraph() (*graph.Graph, error) {
	var err error
	access := graph.NewGraph()
	if !AwlessFirstSync {
		access, err = graph.NewGraphFromFile(filepath.Join(RepoDir, AccessFilename))
		if err != nil {
			return access, err
		}
	}
	return access, nil
}
