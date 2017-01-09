package display

import (
	"fmt"
	"strings"

	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

type Displayer interface {
	Print() string
	SetGraph(*rdf.Graph)
	SetHeaders([]Header)
}

type Options struct {
	RdfType, Format string
}

func BuildDisplayer(opts Options) Displayer {
	switch opts.Format {
	case "csv":
		return &csvDisplayer{rdfType: opts.RdfType}
	default:
		panic(fmt.Sprintf("unknown displayer for %s", opts.Format))
	}
}

type csvDisplayer struct {
	g       *rdf.Graph
	rdfType string
	headers []Header
}

func (d *csvDisplayer) Print() string {
	nodes, err := d.g.NodesForType(d.rdfType)
	if err != nil {
		panic(err)
	}

	var resources []*aws.Resource
	for _, node := range nodes {
		res := aws.NewAwsResource(node.ID().String(), d.rdfType)
		if err := res.UnmarshalFromGraph(d.g); err != nil {
			panic(err)
		}
		resources = append(resources, res)
	}

	var lines []string

	var head []string
	for _, h := range d.headers {
		head = append(head, h.title())
	}

	lines = append(lines, strings.Join(head, ", "))

	for _, res := range resources {
		var props []string
		for _, h := range d.headers {
			props = append(props, h.format(res.Properties()[h.propKey()]))
		}
		lines = append(lines, strings.Join(props, ", "))
	}

	return strings.Join(lines, "\n")
}

func (d *csvDisplayer) SetHeaders(headers []Header) {
	d.headers = headers
}

func (d *csvDisplayer) SetGraph(g *rdf.Graph) {
	d.g = g
}
