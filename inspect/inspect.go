package inspect

import (
	"io"

	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/inspect/inspectors"
)

var InspectorsRegister map[string]Inspector

func init() {
	all := []Inspector{
		&inspectors.Pricer{}, &inspectors.BucketSizer{},
	}

	InspectorsRegister = make(map[string]Inspector)

	for _, i := range all {
		InspectorsRegister[i.Name()] = i
	}
}

type Inspector interface {
	Inspect(...*graph.Graph) error
	Print(io.Writer)
	Name() string
	Services() []string
}
