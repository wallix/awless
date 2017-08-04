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
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"text/tabwriter"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

var pricesURL = "http://ec2-price.com"

type Pricer struct {
	grandTotal float64
	grandCount int
	typeTotal  map[string]float64
	typePrice  map[string]float64
	typeCount  map[string]int
}

func (p *Pricer) Name() string {
	return "pricer"
}

func (p *Pricer) Inspect(g *graph.Graph) error {
	region, err := getRegion(g)
	if err != nil {
		return err
	}

	instances, err := g.GetAllResources(cloud.Instance)
	if err != nil {
		return err
	}

	p.typeCount = make(map[string]int)
	p.typePrice = make(map[string]float64)
	p.typeTotal = make(map[string]float64)
	p.grandCount = 0

	for _, inst := range instances {
		if inst.Properties["State"] == "running" && inst.Properties["Lifecycle"] != "spot" {
			typ := inst.Properties["Type"].(string)
			p.typePrice[typ] = 0.0
			p.typeCount[typ] = p.typeCount[typ] + 1
			p.grandCount += 1
		}
	}

	fmt.Printf("Fetching prices at %s for region %s\n\n", pricesURL, region)

	type result struct {
		typ   string
		price float64
	}

	var wg sync.WaitGroup
	resultC := make(chan result)

	for ty := range p.typePrice {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			price, err := fetchPrice(t, region)
			if err != nil {
				fmt.Printf("fetching price for '%s': %s\n", t, err)
				return
			}

			resultC <- result{typ: t, price: price}
		}(ty)
	}

	go func() {
		wg.Wait()
		close(resultC)
	}()

	for r := range resultC {
		p.typePrice[r.typ] = r.price
	}

	for typ, count := range p.typeCount {
		p.typeTotal[typ] = float64(count) * p.typePrice[typ]
		p.grandTotal = p.grandTotal + p.typeTotal[typ]
	}

	return nil
}

func (p *Pricer) Print(w io.Writer) {
	tabw := tabwriter.NewWriter(w, 5, 8, 2, '\t', 0)

	fmt.Fprintln(tabw, "Instance\tCount\tPrice ea.\tEstimated Total Cost (no EBS)")
	fmt.Fprintln(tabw, "Type\tRunning\tPer Hour\tPer Day\tPer Month\t")
	fmt.Fprintln(tabw, "--------\t--------\t---------\t---------\t---------\t")

	for instType, count := range p.typeCount {
		fmt.Fprintf(tabw, "%s\t%7d\t%8.5f\t%8.2f\t%8.2f\n", instType, count, p.typePrice[instType], p.typePrice[instType]*24, p.typeTotal[instType]*24*30)
	}

	fmt.Fprintln(tabw, "\t--------\t\t---------\t---------\t")
	fmt.Fprintf(tabw, "Grand Total\t%7d\t\t%8.2f\t%8.2f\n", p.grandCount, p.grandTotal*24,p.grandTotal*24*30)

	tabw.Flush()
}

func fetchPrice(instType, region string) (float64, error) {
	resp, err := http.PostForm(
		pricesURL,
		url.Values{"instance_type": {instType}, "location": {region}},
	)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	price, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return 0.0, err
	}

	return price, nil
}

func getRegion(g *graph.Graph) (string, error) {
	all, err := g.GetAllResources("region")
	if err != nil {
		return "", err
	}
	if len(all) < 1 {
		return "", errors.New("cannot resolve region from graph")
	}

	return all[0].Id(), nil
}
