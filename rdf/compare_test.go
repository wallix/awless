package rdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/badwolf/triple"
)

type infraFile struct {
	path, name string
}

func (i *infraFile) String() string {
	return fmt.Sprintf("[infra testfile at %s]", i.path)
}

func TestCompareInfras(t *testing.T) {
	infras := []*infraFile{}

	filepath.Walk("testdata/infra/", rdfFilesOnly(&infras))

	for _, infra := range infras {
		infraGraph, err := NewGraphFromFile(infra.path)
		if err != nil {
			t.Fatal(err)
		}
		for _, otherInfra := range infras {
			otherGraph, err := NewGraphFromFile(otherInfra.path)
			if err != nil {
				t.Fatal(err)
			}

			extraGraph, missingGraph, err := Compare("eu-west-1", infraGraph, otherGraph)
			if err != nil {
				t.Fatal(err)
			}

			expectedExtras := NewGraph()
			expectedMissings := NewGraph()

			if infra.name != otherInfra.name {
				extrasFilename, missingsFilename := expectedFilepathForComparison(infra.name, otherInfra.name)

				if expectedExtras, err = NewGraphFromFile(extrasFilename); err != nil {
					if _, ok := err.(*os.PathError); !ok {
						t.Fatal(err)
					} else {
						continue
					}
				}

				if expectedMissings, err = NewGraphFromFile(missingsFilename); err != nil {
					if _, ok := err.(*os.PathError); !ok {
						t.Fatal(err)
					} else {
						continue
					}
				}
			}

			if got, want := extraGraph.MustMarshal(), expectedExtras.MustMarshal(); got != want {
				t.Fatalf("\n%s: got\n%s\n\n%s: want\n%s\n\n", infra, got, otherInfra, want)
			}
			if got, want := missingGraph.MustMarshal(), expectedMissings.MustMarshal(); got != want {
				t.Fatalf("\n%s: got\n%s\n\n%s: want\n%s\n\n", infra, got, otherInfra, want)
			}

		}
	}
}

func TestIntersectTriples(t *testing.T) {
	var a, b, expect []*triple.Triple

	a = append(a, parseTriple("/a<1>  \"to\"@[] /b<1>"))
	a = append(a, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	a = append(a, parseTriple("/a<3>  \"to\"@[] /b<3>"))
	a = append(a, parseTriple("/a<4>  \"to\"@[] /b<4>"))

	b = append(b, parseTriple("/a<0>  \"to\"@[] /b<0>"))
	b = append(b, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	b = append(b, parseTriple("/a<3>  \"to\"@[] /b<3>"))
	b = append(b, parseTriple("/a<5>  \"to\"@[] /b<5>"))
	b = append(b, parseTriple("/a<6>  \"to\"@[] /b<6>"))

	result := intersectTriples(a, b)
	expect = append(expect, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	expect = append(expect, parseTriple("/a<3>  \"to\"@[] /b<3>"))

	if got, want := marshalTriples(result), marshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestSubstractTriples(t *testing.T) {
	var a, b, expect []*triple.Triple

	a = append(a, parseTriple("/a<1>  \"to\"@[] /b<1>"))
	a = append(a, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	a = append(a, parseTriple("/a<3>  \"to\"@[] /b<3>"))
	a = append(a, parseTriple("/a<4>  \"to\"@[] /b<4>"))

	b = append(b, parseTriple("/a<0>  \"to\"@[] /b<0>"))
	b = append(b, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	b = append(b, parseTriple("/a<3>  \"to\"@[] /b<3>"))
	b = append(b, parseTriple("/a<5>  \"to\"@[] /b<5>"))
	b = append(b, parseTriple("/a<6>  \"to\"@[] /b<6>"))

	result := substractTriples(a, b)
	expect = append(expect, parseTriple("/a<1>  \"to\"@[] /b<1>"))
	expect = append(expect, parseTriple("/a<4>  \"to\"@[] /b<4>"))

	if got, want := marshalTriples(result), marshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	result = substractTriples(b, a)
	expect = []*triple.Triple{}
	expect = append(expect, parseTriple("/a<0>  \"to\"@[] /b<0>"))
	expect = append(expect, parseTriple("/a<5>  \"to\"@[] /b<5>"))
	expect = append(expect, parseTriple("/a<6>  \"to\"@[] /b<6>"))

	if got, want := marshalTriples(result), marshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func rdfFilesOnly(collect *[]*infraFile) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if info.Mode().IsRegular() && ext == ".rdf" {
			name := strings.TrimSuffix(filepath.Base(path), ext)
			*collect = append(*collect, &infraFile{path: path, name: name})
		}
		return nil
	}
}

func expectedFilepathForComparison(first, second string) (extras string, missings string) {
	extras = fmt.Sprintf("testdata/diff_extra/expected_%s_%s.rdf", first, second)
	missings = fmt.Sprintf("testdata/diff_missing/expected_%s_%s.rdf", first, second)
	return
}
