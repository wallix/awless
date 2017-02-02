package cloud

import "github.com/wallix/awless/graph"

type Service interface {
	Name() string
	Provider() string
	ProviderRunnableAPI() interface{}
	ResourceTypes() []string
	FetchResources() (*graph.Graph, error)
	FetchByType(t string) (*graph.Graph, error)
}
