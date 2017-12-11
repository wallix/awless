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
)

var pricesURL = "http://ec2-price.com"

type Pricer struct {
	total float64
	count map[string]int
}

func (p *Pricer) Name() string {
	return "pricer"
}

func (p *Pricer) Inspect(g cloud.GraphAPI) error {
	region, err := getRegion(g)
	if err != nil {
		return err
	}

	instances, err := g.Find(cloud.NewQuery(cloud.Instance))
	if err != nil {
		return err
	}

	p.count = make(map[string]int)
	pricePerType := make(map[string]float64)

	for _, inst := range instances {
		typ := inst.Properties()["Type"].(string)
		pricePerType[typ] = 0.0
		p.count[typ] = p.count[typ] + 1
	}

	fmt.Printf("Fetching prices at %s for region %s\n\n", pricesURL, region)

	type result struct {
		typ   string
		price float64
	}

	var wg sync.WaitGroup
	resultC := make(chan result)

	for ty := range pricePerType {
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
		pricePerType[r.typ] = r.price
	}

	for typ, count := range p.count {
		p.total = p.total + (float64(count) * pricePerType[typ])
	}

	return nil
}

func (p *Pricer) Print(w io.Writer) {
	tabw := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)

	fmt.Fprintln(tabw, "Instance\tCount\tEstimated total/day (no EBS)\t")
	fmt.Fprintln(tabw, "--------\t-----\t----------------------------\t")

	for instType, count := range p.count {
		fmt.Fprintf(tabw, "%s\t%d\t%s\t\n", instType, count, "")
	}

	fmt.Fprintf(tabw, "%s\t%s\t$%.2f\t\n", "", "", p.total*24)

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
	if err != nil {
		return 0.0, err
	}
	price, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return 0.0, err
	}

	return price, nil
}

func getRegion(g cloud.GraphAPI) (string, error) {
	all, err := g.Find(cloud.NewQuery("region"))
	if err != nil {
		return "", err
	}
	if len(all) < 1 {
		return "", errors.New("cannot resolve region from graph")
	}

	return all[0].Id(), nil
}
