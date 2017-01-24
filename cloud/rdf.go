package cloud

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/wallix/awless/graph"
)

func FetchRDFResources(service Service, resourceType graph.ResourceType) (*graph.Graph, error) {
	fnName := fmt.Sprintf("%sGraph", strings.Title(resourceType.PluralString()))
	method := reflect.ValueOf(service).MethodByName(fnName)
	if method.IsValid() && !method.IsNil() {
		methodI := method.Interface()
		if graphFn, ok := methodI.(func() (*graph.Graph, error)); ok {
			return graphFn()
		}
	}
	return nil, (fmt.Errorf("Unknown type of resource: %s", resourceType.String()))
}
