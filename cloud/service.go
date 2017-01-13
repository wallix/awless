package cloud

import "github.com/wallix/awless/rdf"

type Service interface {
	FetchRDFResources(rdf.ResourceType) (*rdf.Graph, error)
}
