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

				var triplesExpectedExtra []*triple.Triple
				var triplesExpectedMissing []*triple.Triple

				if infra1 != infra2 {
					triplesExpectedExtra, err = loadTriplesFromFile(expectedExtraFilename)
					if err != nil {
						t.Fatalf("error '%s' while loading '%s'", err, expectedExtraFilename)
					}

					triplesExpectedMissing, err = loadTriplesFromFile(expectedMissingFilename)
					if err != nil {
						t.Fatalf("error '%s' while loading '%s'", err, expectedMissingFilename)
					}
				}

				extras, _ := extraGraph.allTriples()
				inter := IntersectTriples(triplesExpectedExtra, extras)
				if !(len(inter) == len(triplesExpectedExtra) && len(inter) == extraGraph.TriplesCount()) {
					t.Fatal("sets of triples are not equal: extra elements in one")
				}

				missings, _ := missingGraph.allTriples()
				inter = IntersectTriples(triplesExpectedMissing, missings)
				if !(len(inter) == len(triplesExpectedMissing) && len(inter) == missingGraph.TriplesCount()) {
					t.Fatal("sets of triples are not equal: extra elements in one")
				}
			}
		}
	}
}
