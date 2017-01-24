package cloud

import "github.com/wallix/awless/graph"

type Service interface {
	FetchRDFResources(graph.ResourceType) (*graph.Graph, error)
}
