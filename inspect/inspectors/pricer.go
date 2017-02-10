/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package inspectors

import (
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/wallix/awless/graph"
)

type Pricer struct {
	total float64
	count map[string]int
}

func (p *Pricer) Name() string {
	return "pricer"
}

func (p *Pricer) Services() []string {
	return []string{"infra"}
}

func (p *Pricer) Inspect(graphs ...*graph.Graph) error {
	if len(graphs) < 0 {
		return errors.New("no graph provided for")
	}

	g := graphs[0]

	p.count = make(map[string]int)

	instances, err := g.GetAllResources(graph.Instance)
	if err != nil {
		return err
	}

	for _, inst := range instances {
		typ := inst.Properties["Type"].(string)
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

// This prices serve as an example as they are only valid for eu-west-1 region
var prices = map[string]float64{
	"t2.nano":     0.0063,
	"t2.micro":    0.013,
	"t2.small":    0.025,
	"t2.medium":   0.05,
	"t2.large":    0.101,
	"t2.xlarge":   0.202,
	"t2.2xlarge":  0.404,
	"m4.large":    0.119,
	"m4.xlarge":   0.238,
	"m4.2xlarge":  0.475,
	"m4.4xlarge":  0.95,
	"m4.10xlarge": 2.377,
	"m4.16xlarge": 3.803,
	"m3.medium":   0.073,
	"m3.large":    0.146,
	"m3.xlarge":   0.293,
	"m3.2xlarge":  0.585,
	"c4.large":    0.113,
	"c4.xlarge":   0.226,
	"c4.2xlarge":  0.453,
	"c4.4xlarge":  0.905,
	"c4.8xlarge":  1.811,
	"c3.large":    0.12,
	"c3.xlarge":   0.239,
	"c3.2xlarge":  0.478,
	"c3.4xlarge":  0.956,
	"c3.8xlarge":  1.912,
	"p2.xlarge":   0.972,
	"p2.8xlarge":  7.776,
	"p2.16xlarge": 15.552,
	"g2.2xlarge":  0.702,
	"g2.8xlarge":  2.808,
	"x1.16xlarge": 8.003,
	"x1.32xlarge": 16.006,
	"r3.large":    0.185,
	"r3.xlarge":   0.371,
	"r3.2xlarge":  0.741,
	"r3.4xlarge":  1.482,
	"r3.8xlarge":  2.964,
	"r4.large":    0.148,
	"r4.xlarge":   0.296,
	"r4.2xlarge":  0.593,
	"r4.4xlarge":  1.186,
	"r4.8xlarge":  2.371,
	"r4.16xlarge": 4.742,
	"i2.xlarge":   0.938,
	"i2.2xlarge":  1.876,
	"i2.4xlarge":  3.751,
	"i2.8xlarge":  7.502,
	"d2.xlarge":   0.735,
	"d2.2xlarge":  1.47,
	"d2.4xlarge":  2.94,
	"d2.8xlarge":  5.88,
}
