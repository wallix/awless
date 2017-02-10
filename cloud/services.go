package cloud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wallix/awless/graph"
)

var ErrFetchAccessDenied = errors.New("access denied to cloud resource")

type Service interface {
	Name() string
	Provider() string
	ProviderAPI() string
	ProviderRunnableAPI() interface{}
	ResourceTypes() []string
	FetchResources() (*graph.Graph, error)
	FetchByType(t string) (*graph.Graph, error)
}

var ServiceRegistry = make(map[string]Service)

func GetServiceForType(t string) (Service, error) {
	for _, srv := range ServiceRegistry {
		for _, typ := range srv.ResourceTypes() {
			if typ == t {
				return srv, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot find cloud service for resource type %s", t)
}

func PluralizeResource(singular string) string {
	if strings.HasSuffix(singular, "cy") {
		return strings.TrimSuffix(singular, "y") + "ies"
	}
	return singular + "s"
}
