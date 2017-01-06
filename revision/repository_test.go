package revision

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
)

func TestHasChanges(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitrepo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	rr, err := OpenRepository(dir)
	if err != nil {
		t.Fatal(err)
	}

	hasChanges, err := rr.hasChanges()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := hasChanges, false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	f, err := ioutil.TempFile(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString("test")

	_, fileName := filepath.Split(f.Name())

	if err = rr.addFile(fileName); err != nil {
		t.Fatal(err)
	}

	hasChanges, err = rr.hasChanges()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := hasChanges, true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	rr.commitIfChanges()

	hasChanges, err = rr.hasChanges()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := hasChanges, false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
}

func TestCommitsAndDiffs(t *testing.T) {
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

	lastsDiffs, err := rr.LastDiffs(10, root, NoGroup)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	f, err := ioutil.TempFile(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString("/a<1>  \"to\"@[] /b<1>\n")
	f.WriteString("/b<1>  \"to\"@[] /c<1>\n")

	_, fileName := filepath.Split(f.Name())
	if err = rr.addFile(fileName); err != nil {
		t.Fatal(err)
	}

	err = rr.commitIfChanges()
	if err != nil {
		t.Fatal(err)
	}

	lastsDiffs, err = rr.LastDiffs(10, root, NoGroup)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := len(lastsDiffs[0].GraphDiff.Deleted()), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	expect := []*triple.Triple{parseTriple("/a<1>  \"to\"@[] /b<1>"), parseTriple("/b<1>  \"to\"@[] /c<1>")}
	if got, want := lastsDiffs[0].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got, want := lastsDiffs[0].GraphDiff.FullGraph().MustMarshal(), "/a<1>	\"to\"@[]	/b<1>\n/b<1>	\"to\"@[]	/c<1>"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	f, err = os.OpenFile(f.Name(), os.O_RDWR+os.O_TRUNC, 0666) //empty test file
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("/c<1>  \"to\"@[] /d<1>\n")

	err = rr.commitIfChanges()
	if err != nil {
		t.Fatal(err)
	}

	lastsDiffs, err = rr.LastDiffs(10, root, NoGroup)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := lastsDiffs[0].GraphDiff.Deleted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	expect = []*triple.Triple{parseTriple("/c<1>  \"to\"@[] /d<1>")}
	if got, want := lastsDiffs[0].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	fullGraph := `/a<1>	"to"@[]	/b<1>
/b<1>	"to"@[]	/c<1>
/c<1>	"to"@[]	/d<1>`
	if got, want := lastsDiffs[0].GraphDiff.FullGraph().MustMarshal(), fullGraph; got != want {
		t.Fatalf("got --\n%s\n--, want --\n%s\n--", got, want)
	}

	err = rr.commitIfChanges()
	if err != nil {
		t.Fatal(err)
	}

	lastsDiffs, err = rr.LastDiffs(10, root, NoGroup)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	lastsDiffs, err = rr.LastDiffs(1, root, NoGroup)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	lastsDiffs, err = rr.LastDiffs(2, root, NoGroup)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	f2, err := ioutil.TempFile(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()
	f2.WriteString("/w<1>  \"to\"@[] /x<1>\n")
	f2.WriteString("/x<1>  \"to\"@[] /y<1>\n")
	_, fileName2 := filepath.Split(f2.Name())
	if err = rr.addFile(fileName2); err != nil {
		t.Fatal(err)
	}
	err = rr.commitIfChanges()
	if err != nil {
		t.Fatal(err)
	}

	//Test if fileName has an impact on lastdiff
	lastsDiffs, err = rr.LastDiffs(2, root, NoGroup, fileName)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs[0].GraphDiff.Deleted()), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := len(lastsDiffs[0].GraphDiff.Inserted()), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	fullGraph = `/c<1>	"to"@[]	/d<1>`
	if got, want := lastsDiffs[0].GraphDiff.FullGraph().MustMarshal(), fullGraph; got != want {
		t.Fatalf("got --\n%s\n--, want --\n%s\n--", got, want)
	}
	lastsDiffs, err = rr.LastDiffs(2, root, NoGroup, fileName2)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs[0].GraphDiff.Deleted()), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	fullGraph = `/w<1>	"to"@[]	/x<1>
/x<1>	"to"@[]	/y<1>`
	if got, want := lastsDiffs[0].GraphDiff.FullGraph().MustMarshal(), fullGraph; got != want {
		t.Fatalf("got --\n%s\n--, want --\n%s\n--", got, want)
	}
	expect = []*triple.Triple{parseTriple("/w<1>  \"to\"@[] /x<1>"), parseTriple("/x<1>  \"to\"@[] /y<1>")}
	if got, want := lastsDiffs[0].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	lastsDiffs, err = rr.LastDiffs(2, root, NoGroup)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(lastsDiffs[0].GraphDiff.Deleted()), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := lastsDiffs[0].GraphDiff.Inserted(), expect; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
