package inspect

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/wallix/awless/graph"
)

var InspectorsRegister map[string]Inspector

func init() {
	all := []Inspector{
		&Pricer{},
	}

	InspectorsRegister = make(map[string]Inspector)

	for _, i := range all {
		InspectorsRegister[i.Name()] = i
	}
}

type Inspector interface {
	Name() string
	Inspect(*graph.Graph) error
	Print(io.Writer)
}

type Pricer struct {
	total float64
	count map[string]int
}

func (p *Pricer) Name() string {
	return "pricer"
}

func (p *Pricer) Inspect(g *graph.Graph) error {
	p.count = make(map[string]int)

	instances, err := g.GetAllResources(graph.Instance)
	if err != nil {
		return err
	}

	for _, inst := range instances {
		typ := inst.Properties()["Type"].(string)
		if price, ok := prices[typ]; ok {
			p.total = p.total + price
			p.count[typ] = p.count[typ] + 1
		} else {
			fmt.Printf("no price for instance of type %s", typ)
		}
	}

	return nil
}

func (p *Pricer) Print(w io.Writer) {
	tabw := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)

	fmt.Fprintln(tabw, "Instance\tCount\tEstimated total per day\t")
	fmt.Fprintln(tabw, "--------\t-----\t-----------------------\t")

	for instType, count := range p.count {
		fmt.Fprintf(tabw, "%s\t%d\t%s\t\n", instType, count, "")
	}

	fmt.Fprintf(tabw, "%s\t%s\t$%.2f\t\n", "", "", p.total*24)

	tabw.Flush()
}
