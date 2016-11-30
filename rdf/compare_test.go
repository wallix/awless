package rdf

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/badwolf/triple"
)

var regionID string

func init() {
	regionID = "eu-west-1"
}

func TestCompare(t *testing.T) {
	var infras []string
	files, _ := ioutil.ReadDir("./testdata/infra/")
	for _, f := range files {
		if extIndex := strings.LastIndex(f.Name(), "."); extIndex >= 0 && f.Name()[extIndex+1:] == "rdf" {
			infras = append(infras, f.Name()[:extIndex])
		}
	}

	for _, infra1 := range infras {
		filename1 := "testdata/infra/" + infra1 + ".rdf"
		graphInfra1, err1 := NewGraphFromFile(filename1)
		if err1 != nil {
			t.Fatalf("error '%s' while loading '%s'", err1, filename1)
		}
		for _, infra2 := range infras {
			expectedExtraFilename := "testdata/diff_extra/expected_" + infra1 + "_" + infra2 + ".rdf"
			expectedMissingFilename := "testdata/diff_missing/expected_" + infra1 + "_" + infra2 + ".rdf"
			if _, err := os.Stat(expectedExtraFilename); infra1 != infra2 && os.IsNotExist(err) {
				t.Logf("There is no test data for comparison from %s to %s", infra1, infra2)
			} else {

				filename2 := "testdata/infra/" + infra2 + ".rdf"
				graphInfra2, err2 := NewGraphFromFile(filename2)
				if err2 != nil {
					t.Fatalf("error '%s' while loading '%s'", err2, filename2)
				}

				extraGraph, missingGraph, err := Compare(regionID, graphInfra1, graphInfra2)
				if err != nil {
					t.Fatalf("error while comparing triples : %s", err)
				}

				expectedExtras := NewGraph()
				expectedMissings := NewGraph()

				if infra1 != infra2 {
					expectedExtras, err = NewGraphFromFile(expectedExtraFilename)
					if err != nil {
						t.Fatalf("error '%s' while loading '%s'", err, expectedExtraFilename)
					}

					expectedMissings, err = NewGraphFromFile(expectedMissingFilename)
					if err != nil {
						t.Fatalf("error '%s' while loading '%s'", err, expectedMissingFilename)
					}
				}

				if got, want := extraGraph.FlushString(), expectedExtras.FlushString(); got != want {
					t.Errorf("\nfor %s,%s: got\n%s\n\nwant\n%s\n\n", infra1, infra2, got, want)
				}
				if got, want := missingGraph.FlushString(), expectedMissings.FlushString(); got != want {
					t.Errorf("\nfor %s,%s: got\n%s\n\nwant\n%s\n\n", infra1, infra2, got, want)
				}
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
