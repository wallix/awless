package inspectors

import (
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

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

func (p *BucketSizer) Services() []string {
	return []string{"storage"}
}

func (i *BucketSizer) Inspect(graphs ...*graph.Graph) error {
	if len(graphs) < 0 {
		return errors.New("no graph provided for")
	}

	g := graphs[0]

	i.buckets = make(map[string]*bucket)

	objects, err := g.GetAllResources(graph.Object)
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
		fmt.Fprintf(tabw, "%s\t%d\t%.5f Gb\t\n", name, bucket.objects, bucket.size/1e9)
	}

	fmt.Fprintf(tabw, "%s\t%s\t%.4f Gb\t\n", "", "", i.total/1e9)

	tabw.Flush()
}
