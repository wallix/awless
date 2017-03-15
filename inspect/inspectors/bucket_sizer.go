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
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

type BucketSizer struct {
	total   float64
	buckets map[string]*bucket
}

type bucket struct {
	objects int
	size    float64
}

func (*BucketSizer) Name() string {
	return "bucket_sizer"
}

func (i *BucketSizer) Inspect(g *graph.Graph) error {
	i.buckets = make(map[string]*bucket)

	objects, err := g.GetAllResources(cloud.Object)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		size := obj.Properties["Size"].(float64)
		i.total = i.total + size
		name := obj.Properties["BucketName"].(string)
		b := i.buckets[name]
		if b == nil {
			b = new(bucket)
			i.buckets[name] = b
		}
		b.size = b.size + size
		b.objects = b.objects + 1
	}

	return nil
}

func (i *BucketSizer) Print(w io.Writer) {
	tabw := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)

	fmt.Fprintln(tabw, "Bucket\tObject count\tS3 total storage\t")
	fmt.Fprintln(tabw, "--------\t----------\t-----------------\t")

	for name, bucket := range i.buckets {
		fmt.Fprintf(tabw, "%s\t%d\t%.4f Gb\t\n", name, bucket.objects, bucket.size/1e9)
	}

	fmt.Fprintf(tabw, "%s\t%s\t%.4f Gb\t\n", "", "", i.total/1e9)

	tabw.Flush()
}
