package store

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
		triplesInfra1, err1 := loadTriplesFromFile(filename1)
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
				triplesInfra2, err2 := loadTriplesFromFile(filename2)
				if err2 != nil {
					t.Fatalf("error '%s' while loading '%s'", err2, filename2)
				}

				extraTriples, missingTriples, err := Compare(regionID, triplesInfra1, triplesInfra2)
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

				if got, want := SortLines(MarshalTriples(extraTriples)), SortLines(MarshalTriples(triplesExpectedExtra)); got != want {
					t.Errorf("error on extras from infra '%s' to infra '%s'. \ngot \n%s\nwant \n%s\n", infra1, infra2, got, want)
				}
				if got, want := SortLines(MarshalTriples(missingTriples)), SortLines(MarshalTriples(triplesExpectedMissing)); got != want {
					t.Errorf("error on missings from infra '%s' to infra '%s'.\ngot \n%s\nwant \n%s\n", infra1, infra2, got, want)
				}
			}
		}
	}
}

func loadTriplesFromFile(filepath string) ([]*triple.Triple, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	str := string(content)
	if str == "" || str == "\n" {
		return []*triple.Triple{}, nil
	} else {
		triples, err := UnmarshalTriples(str)

		if err != nil {
			return nil, err
		}
		return triples, nil
	}
}
