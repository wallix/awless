package config

import (
	"path/filepath"

	"github.com/wallix/awless/rdf"
)

func LoadInfraGraph() (*rdf.Graph, error) {
	var err error
	infra := rdf.NewGraph()
	if !AwlessFirstSync {
		infra, err = rdf.NewGraphFromFile(filepath.Join(RepoDir, InfraFilename))
		if err != nil {
			return infra, err
		}
	}
	return infra, nil
}

func LoadAccessGraph() (*rdf.Graph, error) {
	var err error
	access := rdf.NewGraph()
	if !AwlessFirstSync {
		access, err = rdf.NewGraphFromFile(filepath.Join(RepoDir, AccessFilename))
		if err != nil {
			return access, err
		}
	}
	return access, nil
}
