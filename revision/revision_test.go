package revision

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
)

func TestGenerateRevisionPairs(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitrepo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	root, err := node.NewNodeFromStrings("/region", "eu-west-1")
	if err != nil {
		t.Fatal(err)
	}

	rr, err := OpenRepository(dir)
	if err != nil {
		t.Fatal(err)
	}

	f, err := ioutil.TempFile(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString("/a<1>  \"to\"@[] /b<1>\n")

	_, fileName := filepath.Split(f.Name())
	if err = rr.addFile(fileName); err != nil {
		t.Fatal(err)
	}

	if err = rr.commitIfChanges(time.Now().Add(-7 * 24 * time.Hour)); err != nil {
		t.Fatal(err)
	}

	f.WriteString("/b<1>  \"to\"@[] /c<1>\n")
	if err = rr.commitIfChanges(time.Now().Add(-7 * 24 * time.Hour)); err != nil {
		t.Fatal(err)
	}

	f.WriteString("/c<1>  \"to\"@[] /d<1>\n")
	if err = rr.commitIfChanges(); err != nil {
		t.Fatal(err)
	}

	f.WriteString("/d<1>  \"to\"@[] /e<1>\n")
	if err = rr.commitIfChanges(); err != nil {
		t.Fatal(err)
	}
	f.WriteString("/e<1>  \"to\"@[] /f<1>\n")
	if err = rr.commitIfChanges(); err != nil {
		t.Fatal(err)
	}
	lastsDiffs, err := rr.LastDiffs(10, root, NoGroup)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 5; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	expect := []*triple.Triple{parseTriple("/e<1>  \"to\"@[] /f<1>")}
	if got, want := lastsDiffs[0].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	expect = []*triple.Triple{parseTriple("/a<1>  \"to\"@[] /b<1>")}
	if got, want := lastsDiffs[4].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	lastsDiffs, err = rr.LastDiffs(10, root, GroupAll)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	expect = []*triple.Triple{
		parseTriple("/a<1>  \"to\"@[] /b<1>"),
		parseTriple("/b<1>  \"to\"@[] /c<1>"),
		parseTriple("/c<1>  \"to\"@[] /d<1>"),
		parseTriple("/d<1>  \"to\"@[] /e<1>"),
		parseTriple("/e<1>  \"to\"@[] /f<1>"),
	}
	if got, want := lastsDiffs[0].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	lastsDiffs, err = rr.LastDiffs(10, root, GroupByDay)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	expect = []*triple.Triple{
		parseTriple("/d<1>  \"to\"@[] /e<1>"),
		parseTriple("/e<1>  \"to\"@[] /f<1>"),
	}
	if got, want := lastsDiffs[0].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", lastsDiffs[0], want)
	}
	expect = []*triple.Triple{
		parseTriple("/a<1>  \"to\"@[] /b<1>"),
		parseTriple("/b<1>  \"to\"@[] /c<1>"),
		parseTriple("/c<1>  \"to\"@[] /d<1>"),
	}
	if got, want := lastsDiffs[1].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	lastsDiffs, err = rr.LastDiffs(10, root, GroupByWeek)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	expect = []*triple.Triple{
		parseTriple("/d<1>  \"to\"@[] /e<1>"),
		parseTriple("/e<1>  \"to\"@[] /f<1>"),
	}
	if got, want := lastsDiffs[0].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	expect = []*triple.Triple{
		parseTriple("/a<1>  \"to\"@[] /b<1>"),
		parseTriple("/b<1>  \"to\"@[] /c<1>"),
		parseTriple("/c<1>  \"to\"@[] /d<1>"),
	}
	if got, want := lastsDiffs[1].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

}

func parseTriple(s string) *triple.Triple {
	t, err := triple.Parse(s, literal.DefaultBuilder())
	if err != nil {
		panic(err)
	}

	return t
}
